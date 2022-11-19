package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	twitterscraper "github.com/JasonKhew96/twitter-scraper"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/pkg/errors"
)

var allMdV2 = []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
var mdV2Repl = strings.NewReplacer(func() (out []string) {
	for _, x := range allMdV2 {
		out = append(out, x, "\\"+x)
	}
	return out
}()...)

func EscapeMarkdownV2(s string) string {
	return mdV2Repl.Replace(s)
}

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

func tweet2InputMedias(tweet *twitterscraper.Tweet, caption string) []gotgbot.InputMedia {
	inputMedia := []gotgbot.InputMedia{}
	if len(tweet.Medias) > 0 {
		for _, media := range tweet.Medias {
			c := ""
			if len(inputMedia) == 0 {
				c = caption
			}
			switch v := media.(type) {
			case twitterscraper.MediaPhoto:
				if v.Alt != "" {
					c += fmt.Sprintf("\n\n%s", EscapeMarkdownV2(v.Alt))
				}
				newUrl := clearUrlQueries(v.Url)
				var media gotgbot.InputFile
				buf, err := downloadToBuffer(newUrl, "")
				if err != nil {
					log.Println(err)
					media = newUrl
				} else {
					media = buf
				}
				inputMedia = append(inputMedia, gotgbot.InputMediaPhoto{
					Media:     media,
					Caption:   c,
					ParseMode: "MarkdownV2",
				})
			case twitterscraper.MediaVideo:
				if v.Alt != "" {
					c += fmt.Sprintf("\n\n%s", EscapeMarkdownV2(v.Alt))
				}
				newUrl := clearUrlQueries(v.Url)
				var media gotgbot.InputFile
				buf, err := downloadToBuffer(newUrl, "")
				if err != nil {
					log.Println(err)
					media = newUrl
				} else {
					media = buf
				}
				if len(tweet.Medias) == 1 && v.IsAnimatedGif {
					inputMedia = append(inputMedia, gotgbot.InputMediaAnimation{
						Media:     media,
						Caption:   c,
						ParseMode: "MarkdownV2",
					})
				} else {
					inputMedia = append(inputMedia, gotgbot.InputMediaVideo{
						Media:     media,
						Caption:   c,
						ParseMode: "MarkdownV2",
					})
				}
			}
		}
	}
	return inputMedia
}

func clearUrlQueries(link string) string {
	newUrl := link
	if tmp, err := url.Parse(newUrl); err == nil {
		tmp.RawQuery = ""
		newUrl = tmp.String()
	}
	return newUrl
}

func downloadToBuffer(url, fn string) (*gotgbot.NamedFile, error) {
	defaultClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := defaultClient.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download file")
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return &gotgbot.NamedFile{
		File:     buf,
		FileName: fn,
	}, nil
}
