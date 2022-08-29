package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	twitterscraper "github.com/JasonKhew96/twitter-scraper"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/pkg/errors"
)

type TwitterUrl struct {
	Username string
	UserID   string
	TweetID  string
}

func parseTwitterUrl(rawText string) (*TwitterUrl, error) {
	usernameOnlyRegex, err := regexp.Compile(`^([a-zA-Z0-9_]{1,15})$`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile username regex")
	}
	if usernameOnlyRegex.MatchString(rawText) {
		return &TwitterUrl{
			Username: rawText,
		}, nil
	}

	parsedUrl, err := url.Parse(rawText)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	if !strings.HasSuffix(parsedUrl.Host, "twitter.com") {
		return nil, errors.New("url is not a twitter url")
	}

	usernameRegex, err := regexp.Compile(`^\/([a-zA-Z0-9_]{1,15})(\/)?$`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile regex")
	}
	if usernameRegex.MatchString(parsedUrl.Path) {
		return &TwitterUrl{
			Username: usernameRegex.FindStringSubmatch(parsedUrl.Path)[1],
		}, nil
	}

	statusRegex, err := regexp.Compile(`^\/([a-zA-Z0-9_]{1,15})\/status\/(\d+)(\/)?$`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile regex")
	}

	if statusRegex.MatchString(parsedUrl.Path) {
		matches := statusRegex.FindStringSubmatch(parsedUrl.Path)
		return &TwitterUrl{
			Username: matches[1],
			TweetID:  matches[2],
		}, nil
	}

	statusWebRegex, err := regexp.Compile(`^\/i\/web\/status\/(\d+)(\/)?$`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile regex")
	}

	if statusWebRegex.MatchString(parsedUrl.Path) {
		matches := statusWebRegex.FindStringSubmatch(parsedUrl.Path)
		return &TwitterUrl{
			TweetID: matches[1],
		}, nil
	}

	if parsedUrl.Path == "/intent/user" {
		return &TwitterUrl{
			UserID: parsedUrl.Query().Get("user_id"),
		}, nil
	}

	return nil, errors.New("url is not a twitter url")
}

func tweet2Caption(tweet *twitterscraper.Tweet) string {
	caption := fmt.Sprintf("%s\n\n%s", EscapeMarkdownV2(strings.ReplaceAll(tweet.Text, "＃", "#")), EscapeMarkdownV2(tweet.PermanentURL))
	for _, mention := range tweet.Mentions {
		caption = strings.Replace(caption, "@"+EscapeMarkdownV2(mention), fmt.Sprintf(`[@%s](https://twitter\.com/%s)`, EscapeMarkdownV2(mention), EscapeMarkdownV2(mention)), 1)
	}
	for _, hashtag := range tweet.Hashtags {
		caption = strings.Replace(caption, "\\#"+EscapeMarkdownV2(hashtag), fmt.Sprintf(`[\#%s](https://twitter\.com/hashtag/%s)`, EscapeMarkdownV2(hashtag), EscapeMarkdownV2(hashtag)), 1)
	}
	return caption
}

func tweet2InputMedia(tweet *twitterscraper.Tweet, caption string) []gotgbot.InputMedia {
	inputMedia := []gotgbot.InputMedia{}
	if len(tweet.Videos) > 0 {
		for _, v := range tweet.Videos {
			c := ""
			if len(inputMedia) == 0 {
				c = caption
			}
			inputMedia = append(inputMedia, gotgbot.InputMediaVideo{
				Media:     v.URL,
				Caption:   c,
				ParseMode: "MarkdownV2",
			})
		}
	} else {
		for _, p := range tweet.Photos {
			c := ""
			if len(inputMedia) == 0 {
				c = caption
			}
			urlSplit := strings.Split(p, ".")
			newUrl := fmt.Sprintf("%s?format=%s&name=medium", p, urlSplit[len(urlSplit)-1])
			inputMedia = append(inputMedia, gotgbot.InputMediaPhoto{
				Media:     newUrl,
				Caption:   c,
				ParseMode: "MarkdownV2",
			})
		}
	}
	return inputMedia
}
