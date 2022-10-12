package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"twitter-bot/models"

	twitterscraper "github.com/JasonKhew96/twitter-scraper"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type twiCache struct {
	tweetId string
	medias  []twitterscraper.Media
}

type Job struct {
	inputMedias []gotgbot.InputMedia
	cache       *twiCache
}

type bot struct {
	db     *sql.DB
	twit   *twitterscraper.Scraper
	tg     *gotgbot.Bot
	caches map[int64]*twiCache
	jobs   chan Job

	errCount int

	channelChatID int64
	groupChatID   int64
	ownerID       int64

	popularTweetFactor   int
	popularRetweetFactor int
}

func New() (*bot, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	/*
		CREATE TABLE "unfollowed" (
			"uid"			BIGINT NOT NULL UNIQUE PRIMARY KEY,
		);
		CREATE TABLE "tweets" (
			"id"			BIGINT NOT NULL UNIQUE PRIMARY KEY,
			"likes"			BIGINT NOT NULL,
			"retweets"		BIGINT NOT NULL,
			"replies"		BIGINT NOT NULL,
			"medias"		TEXT NOT NULL,
			"text"			TEXT,
			"html"			TEXT,
			"timestamp"		TIMESTAMP NOT NULL,
			"url"			TEXT NOT NULL,
			"uid"			BIGINT NOT NULL,
			"created_at"	TIMESTAMP NOT NULL,
			"updated_at"	TIMESTAMP NOT NULL
		);
	*/

	twit := twitterscraper.New().WithReplies(false).WithDelay(1).WithClientTimeout(time.Minute)
	twit.WithCookie(config.TwitterCookie)
	twit.WithXCsrfToken(config.XCsrfToken)

	b, err := gotgbot.NewBot(config.TelegramBotToken, &gotgbot.BotOpts{
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: time.Minute,
		},
	})
	if err != nil {
		return nil, err
	}

	return &bot{
		db:                   db,
		twit:                 twit,
		tg:                   b,
		caches:               make(map[int64]*twiCache),
		jobs:                 make(chan Job),
		errCount:             0,
		channelChatID:        config.ChannelChatID,
		groupChatID:          config.GroupChatID,
		ownerID:              config.OwnerID,
		popularTweetFactor:   config.PopularTweetFactor,
		popularRetweetFactor: config.PopularRetweetFactor,
	}, nil
}

func (bot *bot) Close() {
	bot.db.Close()
	bot.tg.Close(nil)
}

func (bot *bot) initBot() error {
	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		DispatcherOpts: ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update: ", err.Error())
				return ext.DispatcherActionNoop
			},
		},
	})
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		return msg.Chat.Id == bot.groupChatID
	}, bot.handleChatMessages))
	dispatcher.AddHandler(handlers.NewCommand("follow", bot.commandFollow))
	dispatcher.AddHandler(handlers.NewCommand("unfollow", bot.commandUnfollow))
	dispatcher.AddHandler(handlers.NewMessage(message.Private, bot.handlePrivateMessages))
	dispatcher.AddHandler(handlers.NewCallback(func(cq *gotgbot.CallbackQuery) bool {
		return cq.From.Id == bot.ownerID
	}, bot.handleCallbackData))

	// Start receiving updates.
	err := updater.StartPolling(bot.tg, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 60,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to start polling")
	}
	log.Printf("%s has been started...\n", bot.tg.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()

	return nil
}

func (bot *bot) worker() {
	for job := range bot.jobs {
		if msg, err := bot.tg.SendMediaGroup(bot.channelChatID, job.inputMedias, nil); err != nil {
			log.Println(err)
			bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("%+v\n\n%+v", err.Error(), job.inputMedias), nil)
		} else if len(job.cache.medias) > 0 {
			bot.caches[msg[0].MessageId] = job.cache
		}
		time.Sleep(10 * time.Second)
	}
}

func (bot *bot) handleCallbackData(b *gotgbot.Bot, ctx *ext.Context) error {
	if !strings.Contains(ctx.CallbackQuery.Data, "follow.") && !strings.Contains(ctx.CallbackQuery.Data, "unfollow.") {
		_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "Wrong data",
			ShowAlert: true,
			CacheTime: 60,
		})
		return err
	}
	username := strings.Split(ctx.CallbackQuery.Data, ".")[1]
	switch {
	case strings.HasPrefix(ctx.CallbackQuery.Data, "follow."):
		profile, err := bot.twit.GetProfile(username)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error GetProfile %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		uid, err := strconv.ParseInt(profile.UserID, 10, 64)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error ParseInt %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		if _, err := models.Unfolloweds(models.UnfollowedWhere.UID.EQ(uid)).DeleteAll(context.Background(), bot.db); err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error DeleteAll %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		_, err = bot.twit.Follow(profile.Username)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error Follow %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		_, err = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      fmt.Sprintf("Followed https://twitter.com/%s", username),
			ShowAlert: true,
			CacheTime: 60,
		})
		return err
	case strings.HasPrefix(ctx.CallbackQuery.Data, "unfollow."):
		profile, err := bot.twit.GetProfile(username)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error GetProfile %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		uid, err := strconv.ParseInt(profile.UserID, 10, 64)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error ParseInt %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		t := models.Unfollowed{
			UID: uid,
		}
		if err := t.Insert(context.Background(), bot.db, boil.Infer()); err != nil {
			bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("Error Insert %s", err.Error()), nil)
		}
		_, err = bot.twit.Unfollow(profile.Username)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error Unfollow %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}

		_, err = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      fmt.Sprintf("Unfollowed https://twitter.com/%s", username),
			ShowAlert: true,
			CacheTime: 60,
		})
		return err
	default:
		break
	}
	return nil
}

func (bot *bot) handlePrivateMessages(b *gotgbot.Bot, ctx *ext.Context) error {
	if !(ctx.EffectiveMessage.Text != "" && len(ctx.EffectiveMessage.Entities) > 0 && ctx.EffectiveMessage.Entities[0].Type == "url") {
		return nil
	}
	str := ctx.EffectiveMessage.Text
	url, err := parseTwitterUrl(str)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return err
	}
	if url.TweetID == "" {
		return nil
	}
	tweetID := url.TweetID
	tweet, err := bot.twit.GetTweet(tweetID)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return err
	}
	caption := tweet2Caption(tweet)
	inputMedias := tweet2InputMedias(tweet, caption)

	if len(inputMedias) > 1 {
		_, err = b.SendMediaGroup(ctx.EffectiveChat.Id, inputMedias, nil)
	} else if len(inputMedias) == 1 {
		switch inputMedias[0].(type) {
		case gotgbot.InputMediaPhoto:
			_, err = b.SendPhoto(ctx.EffectiveChat.Id, inputMedias[0].GetMedia(), &gotgbot.SendPhotoOpts{
				Caption:          caption,
				ParseMode:        "MarkdownV2",
				ReplyToMessageId: ctx.EffectiveMessage.MessageId,
			})
		case gotgbot.InputMediaVideo:
			_, err = b.SendVideo(ctx.EffectiveChat.Id, inputMedias[0].GetMedia(), &gotgbot.SendVideoOpts{
				Caption:          caption,
				ParseMode:        "MarkdownV2",
				ReplyToMessageId: ctx.EffectiveMessage.MessageId,
			})
		}
	} else {
		_, err = b.SendMessage(ctx.EffectiveChat.Id, caption, &gotgbot.SendMessageOpts{
			DisableWebPagePreview: true,
			ParseMode:             "MarkdownV2",
			ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
		})
	}
	if err != nil {
		log.Println(err)
		_, err = ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return err
	}

	urlList := []string{}
	// urlList = append(urlList, tweet.Photos...)
	for _, media := range tweet.Medias {
		switch v := media.(type) {
		case twitterscraper.MediaPhoto:
			urlList = append(urlList, clearUrlQueries(v.Url))
		case twitterscraper.MediaVideo:
			urlList = append(urlList, clearUrlQueries(v.Url))
		}
	}

	_, err = ctx.EffectiveMessage.Reply(b, strings.Join(urlList, "\n"), &gotgbot.SendMessageOpts{
		DisableWebPagePreview: true,
	})
	return err
}

func (bot *bot) handleChatMessages(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.Message.SenderChat != nil && ctx.Message.SenderChat.Id != bot.channelChatID {
		return nil
	}
	if c, ok := bot.caches[ctx.Message.ForwardFromMessageId]; ok {
		defer delete(bot.caches, ctx.Message.ForwardFromMessageId)
		if len(c.medias) > 0 {
			var inputMedia []gotgbot.InputMedia
			for i, media := range c.medias {
				var newUrl string
				var fn string
				switch v := media.(type) {
				case twitterscraper.MediaPhoto:
					newUrl = clearUrlQueries(v.Url)
					splits := strings.Split(newUrl, ".")
					ext := splits[len(splits)-1]
					fn = fmt.Sprintf("%s_%02d.%s", c.tweetId, i+1, ext)
					if ext == "jpg" || ext == "jpeg" || ext == "png" {
						newUrl = strings.TrimRight(newUrl, "."+ext) + "?format=" + ext + "&name=orig"
					}
				case twitterscraper.MediaVideo:
					newUrl = clearUrlQueries(v.Url)
					splits := strings.Split(newUrl, ".")
					ext := splits[len(splits)-1]
					fn = fmt.Sprintf("%s_%02d.%s", c.tweetId, i+1, ext)
				}

				var media gotgbot.InputFile
				buf, err := downloadToBuffer(newUrl, fn)
				if err != nil {
					log.Println(err)
					media = newUrl
				} else {
					media = buf
				}

				inputMedia = append(inputMedia, gotgbot.InputMediaDocument{
					Caption: newUrl,
					Media:   media,
				})
			}
			if _, err := b.SendMediaGroup(ctx.Message.Chat.Id, inputMedia, &gotgbot.SendMediaGroupOpts{
				ReplyToMessageId: ctx.Message.MessageId,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bot *bot) commandFollow(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveUser.Id != bot.ownerID {
		return nil
	}
	if ctx.EffectiveMessage.Text == "" {
		return nil
	}
	s := strings.Split(ctx.EffectiveMessage.Text, " ")
	if len(s) != 2 {
		_, err := ctx.EffectiveMessage.Reply(b, "Invalid command format.\nUsage:\n/follow <username>\n/follow <url>", nil)
		return err
	}
	log.Println(ctx.EffectiveMessage.Text)

	twitterUrl, err := parseTwitterUrl(s[1])
	if err != nil {
		_, err := ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return err
	}

	if twitterUrl.Username == "" {
		_, err = ctx.EffectiveMessage.Reply(b, "Invalid twitter username", nil)
		return err
	}

	profile, err := bot.twit.GetProfile(twitterUrl.Username)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error GetProfile %s", err.Error()), nil)
		return err
	}

	uid, err := strconv.ParseInt(profile.UserID, 10, 64)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error ParseInt %s", err.Error()), nil)
		return err
	}

	if _, err := models.Unfolloweds(models.UnfollowedWhere.UID.EQ(uid)).DeleteAll(context.Background(), bot.db); err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error DeleteAll %s", err.Error()), nil)
		return err
	}

	_, err = bot.twit.Follow(profile.Username)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error %s", err.Error()), nil)
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Following https://twitter.com/%s", twitterUrl.Username), &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Follow",
						CallbackData: "follow." + profile.Username,
					},
					{
						Text:         "Unfollow",
						CallbackData: "unfollow." + profile.Username,
					},
				},
			},
		},
	})
	return err
}

func (bot *bot) commandUnfollow(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveUser.Id != bot.ownerID {
		return nil
	}
	if ctx.EffectiveMessage.Text == "" {
		return nil
	}
	s := strings.Split(ctx.EffectiveMessage.Text, " ")
	if len(s) != 2 {
		_, err := ctx.EffectiveMessage.Reply(b, "Invalid command format.\nUsage:\n/unfollow <username>\n/unfollow <url>", nil)
		return err
	}

	log.Println(ctx.EffectiveMessage.Text)

	twitterUrl, err := parseTwitterUrl(s[1])
	if err != nil {
		_, err := ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return err
	}

	if twitterUrl.Username == "" {
		_, err = ctx.EffectiveMessage.Reply(b, "Invalid twitter username", nil)
		return err
	}

	profile, err := bot.twit.GetProfile(twitterUrl.Username)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error GetProfile %s", err.Error()), nil)
		return err
	}

	uid, err := strconv.ParseInt(profile.UserID, 10, 64)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error ParseInt %s", err.Error()), nil)
		return err
	}

	t := models.Unfollowed{
		UID: uid,
	}
	if err := t.Insert(context.Background(), bot.db, boil.Infer()); err != nil {
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error %s", err.Error()), nil)
	}

	_, err = bot.twit.Unfollow(profile.Username)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error %s", err.Error()), nil)
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Unfollowed https://twitter.com/%s", twitterUrl.Username), &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Follow",
						CallbackData: "follow." + profile.Username,
					},
					{
						Text:         "Unfollow",
						CallbackData: "unfollow." + profile.Username,
					},
				},
			},
		},
	})
	return err
}

func (bot *bot) getTweetById(id int64) (*models.Tweet, error) {
	return models.Tweets(qm.Where("id=?", id)).One(context.Background(), bot.db)
}

func (bot *bot) insertTweet(tweet *twitterscraper.Tweet) error {
	id, err := strconv.ParseInt(tweet.ID, 10, 64)
	if err != nil {
		return err
	}

	uid, err := strconv.ParseInt(tweet.UserID, 10, 64)
	if err != nil {
		return err
	}

	medias := ""
	if len(tweet.Medias) > 0 {
		var urlList []string
		for _, media := range tweet.Medias {
			switch v := media.(type) {
			case twitterscraper.MediaPhoto:
				urlList = append(urlList, clearUrlQueries(v.Url))
			case twitterscraper.MediaVideo:
				urlList = append(urlList, clearUrlQueries(v.Url))
			}
		}
		medias = strings.Join(urlList, "|")
	}

	t := models.Tweet{
		ID:        id,
		Likes:     int64(tweet.Likes),
		Retweets:  int64(tweet.Retweets),
		Replies:   int64(tweet.Replies),
		Medias:    medias,
		Text:      null.StringFrom(tweet.Text),
		HTML:      null.StringFrom(tweet.HTML),
		Timestamp: time.Unix(tweet.Timestamp, 0),
		URL:       tweet.PermanentURL,
		UID:       uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return t.Insert(context.Background(), bot.db, boil.Infer())
}

func isMedia(tweet *twitterscraper.Tweet) bool {
	return len(tweet.Medias) > 0
}

func (bot *bot) isPopularRetweet(t time.Time, likes int) bool {
	sinceHours := int(math.Floor(time.Since(t).Hours()))

	for h := 1; h <= 24*3; h++ {
		if sinceHours <= h && likes >= h*bot.popularRetweetFactor {
			return true
		}
	}

	return false
}

func (bot *bot) isPopularTweet(t time.Time, likes int) bool {
	sinceHours := int(math.Floor(time.Since(t).Hours()))

	for h := 1; h <= 24*3; h++ {
		if sinceHours <= h && likes >= h*bot.popularTweetFactor {
			return true
		}
	}

	return false
}

func isRepost(tweet *twitterscraper.Tweet) bool {
	forbiddenHashTags := []string{"フォロー", "フォロワー", "連休", "見た人", "自分が", "晒そう", "晒す", "貼る"}
	// forbiddenRegexHashTags := []string{`\d{4}年自分が選ぶ今年[上下]半期の\d枚`, `今[年月]描いた絵を晒そう`, "^自分が"}
	forbiddenTexts := []string{"再掲", "過去絵", "去年"}
	forbiddenRegexTexts := []string{`(?i)\bwip\b`}

	for _, hashTag := range tweet.Hashtags {
		for _, f := range forbiddenHashTags {
			if strings.Contains(hashTag, f) {
				return true
			}
		}
		// for _, f := range forbiddenRegexHashTags {
		// 	if regexp.MustCompile(f).MatchString(hashTag) {
		// 		return true
		// 	}
		// }
	}

	for _, forbiddenText := range forbiddenTexts {
		if strings.Contains(tweet.Text, forbiddenText) {
			return true
		}
	}
	for _, forbiddenRegexText := range forbiddenRegexTexts {
		if regexp.MustCompile(forbiddenRegexText).MatchString(tweet.Text) {
			return true
		}
	}
	return false
}

// guess
func isIllustrator(text string) bool {
	keyword := []string{"pixiv", "skeb", "potofu"}
	textLower := strings.ToLower(text)
	for _, k := range keyword {
		if strings.Contains(textLower, k) {
			return true
		}
	}
	return false
}

func (bot *bot) processRetweet(tweet *twitterscraper.Tweet) error {
	if !bot.isPopularRetweet(tweet.TimeParsed, tweet.Likes) {
		return nil
	}

	id, err := strconv.ParseInt(tweet.ID, 10, 64)
	if err != nil {
		return err
	}
	if d, err := bot.getTweetById(id); err == nil && d != nil {
		return nil
	}

	user, err := bot.twit.GetProfile(tweet.Username)
	if err != nil {
		return err
	}

	uid, err := strconv.ParseInt(user.UserID, 10, 64)
	if err != nil {
		return err
	}

	count, err := models.Unfolloweds(models.UnfollowedWhere.UID.EQ(uid)).Count(context.Background(), bot.db)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	if !isIllustrator(user.Biography) && !isIllustrator(user.Website) {
		return nil
	}
	if !user.IsFollowing {
		log.Println("Suggest", tweet.PermanentURL)
		if _, err := bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("Followed https://twitter.com/%s", tweet.Username), &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Follow",
							CallbackData: "follow." + tweet.Username,
						},
						{
							Text:         "Unfollow",
							CallbackData: "unfollow." + tweet.Username,
						},
					},
				},
			},
		}); err != nil {
			log.Println(err)
		}
		if _, err := bot.twit.Follow(tweet.Username); err != nil {
			return err
		}
	}

	if err := bot.insertTweet(tweet); err != nil {
		return err
	}

	if isRepost(tweet) {
		return nil
	}

	log.Println(tweet.Likes, tweet.SensitiveContent, tweet.PermanentURL)

	caption := tweet2Caption(tweet)
	inputMedia := tweet2InputMedias(tweet, caption)

	bot.jobs <- Job{
		inputMedias: inputMedia,
		cache: &twiCache{
			tweetId: tweet.ID,
			medias:  tweet.Medias,
		},
	}

	return nil
}

func (bot *bot) processTweet(tweet *twitterscraper.Tweet) error {
	if !bot.isPopularTweet(tweet.TimeParsed, tweet.Likes) {
		return nil
	}

	id, err := strconv.ParseInt(tweet.ID, 10, 64)
	if err != nil {
		return err
	}
	if d, err := bot.getTweetById(id); err == nil && d != nil {
		return nil
	}
	if err := bot.insertTweet(tweet); err != nil {
		return err
	}

	if isRepost(tweet) {
		return nil
	}

	log.Println(tweet.Likes, tweet.SensitiveContent, tweet.PermanentURL)

	caption := tweet2Caption(tweet)
	inputMedia := tweet2InputMedias(tweet, caption)

	bot.jobs <- Job{
		inputMedias: inputMedia,
		cache: &twiCache{
			tweetId: tweet.ID,
			medias:  tweet.Medias,
		},
	}

	return nil
}

func (bot *bot) newLoop() error {
	for tweet := range bot.twit.GetHomeTimeline(context.Background(), 20*50) {
		if tweet.Error != nil {
			bot.errCount++
			log.Println("GetHomeTimeline error")
			bot.tg.SendMessage(bot.ownerID, "GetHomeTimeline error", nil)
			time.Sleep(time.Minute)
			break
		}
		bot.errCount = 0

		if tweet.IsRetweet && tweet.RetweetedStatus.UserID != tweet.UserID && (isMedia(tweet.RetweetedStatus)) {
			err := bot.processRetweet(tweet.RetweetedStatus)
			if err != nil {
				bot.errCount++
				log.Println("processRetweet error", tweet.PermanentURL)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processRetweet error %s", tweet.PermanentURL), nil)
				time.Sleep(time.Minute)
				continue
			}
		} else if isMedia(&tweet.Tweet) {
			err := bot.processTweet(&tweet.Tweet)
			if err != nil {
				bot.errCount++
				log.Println("processTweet error", tweet.PermanentURL)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processTweet error %s", tweet.PermanentURL), nil)
				time.Sleep(time.Minute)
				continue
			}
		}
	}
	if bot.errCount > 5 {
		return fmt.Errorf("TOO MUCH ERROR")
	}
	time.Sleep(5 * time.Second)
	for tweet := range bot.twit.GetHomeLatestTimeline(context.Background(), 20*50) {
		if tweet.Error != nil {
			bot.errCount++
			log.Println("GetHomeLatestTimeline error")
			bot.tg.SendMessage(bot.ownerID, "GetHomeLatestTimeline error", nil)
			time.Sleep(time.Minute)
			break
		}
		bot.errCount = 0

		if tweet.IsRetweet && tweet.RetweetedStatus.UserID != tweet.UserID && (isMedia(tweet.RetweetedStatus)) {
			err := bot.processRetweet(tweet.RetweetedStatus)
			if err != nil {
				bot.errCount++
				log.Println("processRetweet error", tweet.PermanentURL)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processRetweet error %s", tweet.PermanentURL), nil)
				time.Sleep(time.Minute)
				continue
			}
		} else if isMedia(&tweet.Tweet) {
			err := bot.processTweet(&tweet.Tweet)
			if err != nil {
				bot.errCount++
				log.Println("processTweet error", tweet.PermanentURL)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processTweet error %s", tweet.PermanentURL), nil)
				time.Sleep(time.Minute)
				continue
			}
		}
	}
	if bot.errCount > 5 {
		return fmt.Errorf("TOO MUCH ERROR")
	}
	return nil
}

func (bot *bot) cleanup() error {
	count, err := models.Tweets(models.TweetWhere.CreatedAt.LT(time.Now().Add(-30*24*time.Hour))).DeleteAll(context.Background(), bot.db)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Deleted %d old tweet(s)", count)
	}
	return nil
}
