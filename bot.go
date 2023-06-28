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

	"github.com/JasonKhew96/twiscraper"
	"github.com/JasonKhew96/twiscraper/entity"
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
	medias  []entity.ParsedMedia
}

type Job struct {
	inputMedias []gotgbot.InputMedia
	cache       *twiCache
}

type bot struct {
	db     *sql.DB
	twit   *twiscraper.Scraper
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
			"uid"			BIGINT NOT NULL UNIQUE PRIMARY KEY
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

	twit, err := twiscraper.New(&twiscraper.ScraperOptions{
		Delay:      3 * time.Second,
		Cookie:     config.TwitterCookie,
		XCsrfToken: config.XCsrfToken,
		Timeout:    time.Minute,
	})
	if err != nil {
		return nil, err
	}

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
	updater := ext.NewUpdater(nil)
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
		var msgs []gotgbot.Message
		var err error
		if len(job.inputMedias) > 1 {
			msgs, err = bot.tg.SendMediaGroup(bot.channelChatID, job.inputMedias, nil)
		} else {
			var msg *gotgbot.Message
			switch job.inputMedias[0].(type) {
			case gotgbot.InputMediaPhoto:
				caption := job.inputMedias[0].(gotgbot.InputMediaPhoto).Caption
				msg, err = bot.tg.SendPhoto(bot.channelChatID, job.inputMedias[0].GetMedia(), &gotgbot.SendPhotoOpts{
					Caption:   caption,
					ParseMode: "MarkdownV2",
				})
			case gotgbot.InputMediaVideo:
				caption := job.inputMedias[0].(gotgbot.InputMediaVideo).Caption
				msg, err = bot.tg.SendVideo(bot.channelChatID, job.inputMedias[0].GetMedia(), &gotgbot.SendVideoOpts{
					Caption:   caption,
					ParseMode: "MarkdownV2",
				})
			case gotgbot.InputMediaAnimation:
				caption := job.inputMedias[0].(gotgbot.InputMediaAnimation).Caption
				msg, err = bot.tg.SendAnimation(bot.channelChatID, job.inputMedias[0].GetMedia(), &gotgbot.SendAnimationOpts{
					Caption:   caption,
					ParseMode: "MarkdownV2",
				})
			default:
				log.Println("unknown media type ", job.inputMedias[0])
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("unknown media type in worker\n\n%+v", job.inputMedias), nil)
			}
			msgs = append(msgs, *msg)
		}
		if err != nil {
			log.Println(err)
			bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("%+v\n\n%+v", err.Error(), job.inputMedias), nil)
		} else if len(job.cache.medias) > 0 {
			bot.caches[msgs[0].MessageId] = job.cache
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
		profile, err := bot.twit.GetUserByScreenName(username)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error GetProfile %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		uid, err := strconv.ParseInt(profile.UserId, 10, 64)
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
		if err := bot.twit.Follow(profile.ScreenName); err != nil {
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
		profile, err := bot.twit.GetUserByScreenName(username)
		if err != nil {
			_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text:      fmt.Sprintf("Error GetProfile %s", err.Error()),
				ShowAlert: true,
				CacheTime: 60,
			})
			return err
		}
		uid, err := strconv.ParseInt(profile.UserId, 10, 64)
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
		if err := bot.twit.UnFollow(profile.ScreenName); err != nil {
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
	if ctx.EffectiveUser.Id != bot.ownerID {
		return nil
	}
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
	tweet, err := bot.twit.GetTweetDetail(tweetID)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return err
	}
	caption := tweet2Caption(tweet)
	inputMedias := tweet2InputMedias(tweet, caption)

	if len(inputMedias) > 1 {
		_, err = b.SendMediaGroup(ctx.EffectiveChat.Id, inputMedias, nil)
	} else if len(inputMedias) == 1 {
		inputMedia := inputMedias[0]
		switch media := inputMedia.(type) {
		case gotgbot.InputMediaPhoto:
			_, err = b.SendPhoto(ctx.EffectiveChat.Id, inputMedia.GetMedia(), &gotgbot.SendPhotoOpts{
				Caption:          media.Caption,
				ParseMode:        "MarkdownV2",
				ReplyToMessageId: ctx.EffectiveMessage.MessageId,
			})
		case gotgbot.InputMediaVideo:
			_, err = b.SendVideo(ctx.EffectiveChat.Id, inputMedia.GetMedia(), &gotgbot.SendVideoOpts{
				Caption:          media.Caption,
				ParseMode:        "MarkdownV2",
				ReplyToMessageId: ctx.EffectiveMessage.MessageId,
			})
		case gotgbot.InputMediaAnimation:
			_, err = b.SendAnimation(ctx.EffectiveChat.Id, inputMedia.GetMedia(), &gotgbot.SendAnimationOpts{
				Caption:          media.Caption,
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
	for _, media := range tweet.Entities.Media {
		switch v := media.(type) {
		case entity.ParsedMediaPhoto:
			urlList = append(urlList, clearUrlQueries(v.Url))
		case entity.ParsedMediaVideo:
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
				case entity.ParsedMediaPhoto:
					newUrl = clearUrlQueries(v.Url)
					splits := strings.Split(newUrl, ".")
					ext := splits[len(splits)-1]
					fn = fmt.Sprintf("%s_%02d.%s", c.tweetId, i+1, ext)
					if ext == "jpg" || ext == "jpeg" || ext == "png" {
						newUrl = strings.TrimSuffix(newUrl, "."+ext) + "?format=" + ext + "&name=orig"
					}
				case entity.ParsedMediaVideo:
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
				log.Println(err)
				_, err = bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("%+v\n\n%+v\n\n%s", err.Error(), inputMedia, ctx.Message.Entities[len(ctx.Message.Entities)-1].Url), nil)
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

	profile, err := bot.twit.GetUserByScreenName(twitterUrl.Username)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error GetProfile %s", err.Error()), nil)
		return err
	}

	uid, err := strconv.ParseInt(profile.UserId, 10, 64)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error ParseInt %s", err.Error()), nil)
		return err
	}

	if _, err := models.Unfolloweds(models.UnfollowedWhere.UID.EQ(uid)).DeleteAll(context.Background(), bot.db); err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error DeleteAll %s", err.Error()), nil)
		return err
	}

	if err := bot.twit.Follow(profile.ScreenName); err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error %s", err.Error()), nil)
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Following https://twitter.com/%s", twitterUrl.Username), &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Follow",
						CallbackData: "follow." + profile.ScreenName,
					},
					{
						Text:         "Unfollow",
						CallbackData: "unfollow." + profile.ScreenName,
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

	profile, err := bot.twit.GetUserByScreenName(twitterUrl.Username)
	if err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error GetProfile %s", err.Error()), nil)
		return err
	}

	uid, err := strconv.ParseInt(profile.UserId, 10, 64)
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

	if err := bot.twit.UnFollow(profile.ScreenName); err != nil {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Error %s", err.Error()), nil)
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Unfollowed https://twitter.com/%s", twitterUrl.Username), &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Follow",
						CallbackData: "follow." + profile.ScreenName,
					},
					{
						Text:         "Unfollow",
						CallbackData: "unfollow." + profile.ScreenName,
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

func (bot *bot) insertTweet(tweet *entity.ParsedTweet) error {
	id, err := strconv.ParseInt(tweet.TweetId, 10, 64)
	if err != nil {
		return err
	}

	uid, err := strconv.ParseInt(tweet.ParsedUser.UserId, 10, 64)
	if err != nil {
		return err
	}

	medias := ""
	if len(tweet.Entities.Media) > 0 {
		var urlList []string
		for _, media := range tweet.Entities.Media {
			switch v := media.(type) {
			case entity.ParsedMediaPhoto:
				urlList = append(urlList, clearUrlQueries(v.Url))
			case entity.ParsedMediaVideo:
				urlList = append(urlList, clearUrlQueries(v.Url))
			}
		}
		medias = strings.Join(urlList, "|")
	}

	t := models.Tweet{
		ID:       id,
		Likes:    int64(tweet.FavouriteCount),
		Retweets: int64(tweet.RetweetedCount),
		Replies:  int64(tweet.ReplyCount),
		Medias:   medias,
		Text:     null.StringFrom(tweet.FullText),
		// HTML:      null.StringFrom(tweet.HTML),
		Timestamp: tweet.CreatedAt,
		URL:       tweet.Url,
		UID:       uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return t.Insert(context.Background(), bot.db, boil.Infer())
}

func isMedia(tweet entity.ParsedTweet) bool {
	return len(tweet.Entities.Media) > 0
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

func isRepost(tweet *entity.ParsedTweet) bool {
	forbiddenHashTags := []string{
		"フォロー",
		"フォロワー",
		"連休",
		"見た人",
		"自分が",
		"晒そう",
		"晒す",
		"貼る",
	}
	forbiddenRegexHashTags := []string{
		`^いい\W+の日$`,
		`を(見|み)せてください$`,
		`見てみましょう$`,
		`^自分の`,
		`^今までで`,
		`^太ももは`,
		`^見た`,
		`^今(年|月)`,
		`^あなたの`,
		`^(春|夏|秋|冬)が終わり`,
		`^独学でここまで`,
		`一本勝負$`,
		`^みんなさん`,
		`(?i)^aiart(work|community)?$`,
		`(?i)^midjourney$`,
		`(?i)^(stable|waifu)diffusion(art)?$`,
		`(?i)^dreambooth$`,
		`(?i)^novelai$`,
		`(?i)^AIイラスト$`,
	}
	forbiddenTexts := []string{"再掲", "過去絵", "去年", "あなたのサークル"}
	forbiddenRegexTexts := []string{`(?i)\bwip\b`}

	for _, hashTag := range tweet.Entities.Hashtags {
		for _, f := range forbiddenHashTags {
			if strings.Contains(hashTag, f) {
				return true
			}
		}
		for _, f := range forbiddenRegexHashTags {
			if regexp.MustCompile(f).MatchString(hashTag) {
				return true
			}
		}
	}

	for _, forbiddenText := range forbiddenTexts {
		if strings.Contains(tweet.FullText, forbiddenText) {
			return true
		}
	}
	for _, forbiddenRegexText := range forbiddenRegexTexts {
		if regexp.MustCompile(forbiddenRegexText).MatchString(tweet.FullText) {
			return true
		}
	}
	return false
}

// guess
func isIllustrator(text string) bool {
	keyword := []string{"pixiv", "skeb", "potofu", "fanbox", "patreon", "rkgk"}
	textLower := strings.ToLower(text)
	for _, k := range keyword {
		if strings.Contains(textLower, k) {
			return true
		}
	}
	return false
}

func (bot *bot) processRetweet(tweet *entity.ParsedTweet, retweetUserId string) error {
	id, err := strconv.ParseInt(tweet.TweetId, 10, 64)
	if err != nil {
		return err
	}
	if d, err := bot.getTweetById(id); err == nil && d != nil {
		return nil
	}

	isMentioned := false
	for _, mention := range tweet.Entities.UserMentions {
		if mention.UserId == retweetUserId {
			isMentioned = true
			break
		}
	}

	if isMentioned {
		if !bot.isPopularTweet(tweet.CreatedAt, tweet.FavouriteCount) {
			return nil
		}
	} else {
		if !bot.isPopularRetweet(tweet.CreatedAt, tweet.FavouriteCount) {
			return nil
		}
	}

	uid, err := strconv.ParseInt(tweet.ParsedUser.UserId, 10, 64)
	if err != nil {
		return err
	}

	if !isMentioned {
		count, err := models.Unfolloweds(models.UnfollowedWhere.UID.EQ(uid)).Count(context.Background(), bot.db)
		if err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		if !isIllustrator(tweet.ParsedUser.Description) && !isIllustrator(tweet.ParsedUser.Url) {
			return nil
		}

		if !tweet.ParsedUser.IsFollowing {
			log.Println("Suggest", tweet.FavouriteCount, tweet.Views, tweet.Url)
			if _, err := bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("Followed https://twitter.com/%s", tweet.ParsedUser.ScreenName), &gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{
							{
								Text:         "Follow",
								CallbackData: "follow." + tweet.ParsedUser.ScreenName,
							},
							{
								Text:         "Unfollow",
								CallbackData: "unfollow." + tweet.ParsedUser.ScreenName,
							},
						},
					},
				},
			}); err != nil {
				log.Println(err)
			}
			if err := bot.twit.Follow(tweet.ParsedUser.ScreenName); err != nil {
				return err
			}
		}
	}

	if err := bot.insertTweet(tweet); err != nil {
		return err
	}

	if isRepost(tweet) {
		return nil
	}

	log.Println("retweet", tweet.FavouriteCount, tweet.Views, tweet.Url)

	caption := tweet2Caption(tweet)
	inputMedias := tweet2InputMedias(tweet, caption)

	bot.jobs <- Job{
		inputMedias: inputMedias,
		cache: &twiCache{
			tweetId: tweet.TweetId,
			medias:  tweet.Entities.Media,
		},
	}

	return nil
}

func (bot *bot) processTweet(tweet *entity.ParsedTweet) error {
	if !bot.isPopularTweet(tweet.CreatedAt, tweet.FavouriteCount) {
		return nil
	}

	id, err := strconv.ParseInt(tweet.TweetId, 10, 64)
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

	log.Println("tweet", tweet.FavouriteCount, tweet.Views, tweet.Url)

	caption := tweet2Caption(tweet)
	inputMedias := tweet2InputMedias(tweet, caption)

	bot.jobs <- Job{
		inputMedias: inputMedias,
		cache: &twiCache{
			tweetId: tweet.TweetId,
			medias:  tweet.Entities.Media,
		},
	}

	return nil
}

func (bot *bot) newLoop() error {
	for tweet := range bot.twit.GetHomeTimeline(context.Background(), 20*100) {
		if tweet.Error != nil {
			bot.errCount++
			log.Println("GetHomeTimeline error", tweet.Error)
			bot.tg.SendMessage(bot.ownerID, "GetHomeTimeline error", nil)
			time.Sleep(time.Minute)
			break
		}
		bot.errCount = 0

		if !isMedia(tweet.ParsedTweet) {
			continue
		}

		if tweet.ParsedTweet.IsRetweet && tweet.ParsedTweet.RetweetedTweet.ParsedUser.UserId == tweet.ParsedTweet.ParsedUser.UserId {
			continue
		}

		if tweet.ParsedTweet.IsRetweet {
			err := bot.processRetweet(tweet.ParsedTweet.RetweetedTweet, tweet.ParsedTweet.ParsedUser.UserId)
			if err != nil {
				bot.errCount++
				log.Println("processRetweet error", tweet.ParsedTweet.Url, err)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processRetweet error %s", tweet.ParsedTweet.Url), nil)
				time.Sleep(time.Minute)
				continue
			}
		} else if tweet.ParsedTweet.IsRecommended || !tweet.ParsedTweet.ParsedUser.IsFollowing {
			err := bot.processRetweet(&tweet.ParsedTweet, "")
			if err != nil {
				bot.errCount++
				log.Println("processRetweet recommended error", tweet.ParsedTweet.Url, err)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processRetweet recommended error %s", tweet.ParsedTweet.Url), nil)
				time.Sleep(time.Minute)
				continue
			}
		} else {
			err := bot.processTweet(&tweet.ParsedTweet)
			if err != nil {
				bot.errCount++
				log.Println("processTweet error", tweet.ParsedTweet.Url, err)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processTweet error %s", tweet.ParsedTweet.Url), nil)
				time.Sleep(time.Minute)
				continue
			}
		}
	}
	if bot.errCount > 5 {
		return fmt.Errorf("TOO MUCH ERROR")
	}
	time.Sleep(5 * time.Second)
	for tweet := range bot.twit.GetHomeLatestTimeline(context.Background(), 20*100) {
		if tweet.Error != nil {
			bot.errCount++
			log.Println("GetHomeLatestTimeline error", tweet.Error)
			bot.tg.SendMessage(bot.ownerID, "GetHomeLatestTimeline error", nil)
			time.Sleep(time.Minute)
			break
		}
		bot.errCount = 0

		if !isMedia(tweet.ParsedTweet) {
			continue
		}

		if tweet.ParsedTweet.IsRetweet && tweet.ParsedTweet.RetweetedTweet.ParsedUser.UserId == tweet.ParsedTweet.ParsedUser.UserId {
			continue
		}

		if tweet.ParsedTweet.IsRetweet {
			err := bot.processRetweet(tweet.ParsedTweet.RetweetedTweet, tweet.ParsedTweet.ParsedUser.UserId)
			if err != nil {
				bot.errCount++
				log.Println("processRetweet error", tweet.ParsedTweet.Url, err)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processRetweet error %s", tweet.ParsedTweet.Url), nil)
				time.Sleep(time.Minute)
				continue
			}
		} else if tweet.ParsedTweet.IsRecommended || !tweet.ParsedTweet.ParsedUser.IsFollowing {
			err := bot.processRetweet(&tweet.ParsedTweet, "")
			if err != nil {
				bot.errCount++
				log.Println("processRetweet recommended error", tweet.ParsedTweet.Url, err)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processRetweet recommended error %s", tweet.ParsedTweet.Url), nil)
				time.Sleep(time.Minute)
				continue
			}
		} else {
			err := bot.processTweet(&tweet.ParsedTweet)
			if err != nil {
				bot.errCount++
				log.Println("processTweet error", tweet.ParsedTweet.Url, err)
				bot.tg.SendMessage(bot.ownerID, fmt.Sprintf("processTweet error %s", tweet.ParsedTweet.Url), nil)
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
	count, err := models.Tweets(models.TweetWhere.CreatedAt.LT(time.Now().Add(-90*24*time.Hour))).DeleteAll(context.Background(), bot.db)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Deleted %d old tweet(s)", count)
	}
	return nil
}
