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

	"github.com/JasonKhew96/twiscraper/entity"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/pkg/errors"
)

var allMdV2 = []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!", "\\"}
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

	if !strings.HasSuffix(parsedUrl.Host, "twitter.com") && !strings.HasSuffix(parsedUrl.Host, "x.com") {
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

func tweet2Caption(tweet *entity.ParsedTweet) string {
	caption := fmt.Sprintf("%s\n\n%s", EscapeMarkdownV2(strings.ReplaceAll(tweet.FullText, "＃", "#")), EscapeMarkdownV2(tweet.Url))
	for _, mention := range tweet.Entities.UserMentions {
		caption = strings.Replace(caption, "@"+EscapeMarkdownV2(mention.ScreenName), fmt.Sprintf(`[@%s](https://x\.com/%s)`, EscapeMarkdownV2(mention.ScreenName), EscapeMarkdownV2(mention.ScreenName)), 1)
	}
	for _, hashtag := range tweet.Entities.Hashtags {
		caption = strings.Replace(caption, "\\#"+EscapeMarkdownV2(hashtag), fmt.Sprintf(`[\#%s](https://x\.com/hashtag/%s)`, EscapeMarkdownV2(hashtag), EscapeMarkdownV2(url.QueryEscape(hashtag))), 1)
	}
	return caption
}

func tweet2InputMedias(tweet *entity.ParsedTweet, caption string) []gotgbot.InputMedia {
	inputMedia := []gotgbot.InputMedia{}
	if len(tweet.Entities.Media) > 0 {
		for index, media := range tweet.Entities.Media {
			switch v := media.(type) {
			case entity.ParsedMediaPhoto:
				if v.AltText != "" {
					caption += fmt.Sprintf("\n\n\\[%d\\] %s", index+1, EscapeMarkdownV2(v.AltText))
				}
			case entity.ParsedMediaVideo:
				if v.AltText != "" {
					caption += fmt.Sprintf("\n\n\\[%d\\] %s", index+1, EscapeMarkdownV2(v.AltText))
				}
			}
		}
		for i, media := range tweet.Entities.Media {
			c := ""
			if len(inputMedia) == 0 {
				c = caption
			}
			var fn string
			switch v := media.(type) {
			case entity.ParsedMediaPhoto:
				newUrl := clearUrlQueries(v.Url)
				splits := strings.Split(newUrl, ".")
				ext := splits[len(splits)-1]
				fn = fmt.Sprintf("%s_%02d.%s", tweet.TweetId, i+1, ext)
				if ext == "jpg" || ext == "jpeg" || ext == "png" {
					newUrl = strings.TrimSuffix(newUrl, "."+ext) + "?format=" + ext + "&name=large"
				}
				var media gotgbot.InputFileOrString
				buf, err := downloadToBuffer(newUrl, fn)
				if err != nil {
					log.Println(err)
					media = gotgbot.InputFileByURL(newUrl)
				} else {
					media = buf
				}
				inputMedia = append(inputMedia, gotgbot.InputMediaPhoto{
					Media:     media,
					Caption:   c,
					ParseMode: "MarkdownV2",
				})
			case entity.ParsedMediaVideo:
				newUrl := clearUrlQueries(v.Url)
				splits := strings.Split(newUrl, ".")
				ext := splits[len(splits)-1]
				fn = fmt.Sprintf("%s_%02d.%s", tweet.TweetId, i+1, ext)
				var media gotgbot.InputFileOrString
				buf, err := downloadToBuffer(newUrl, fn)
				if err != nil {
					log.Println(err)
					media = gotgbot.InputFileByURL(newUrl)
				} else {
					media = buf
				}
				width := int64(v.Width)
				height := int64(v.Height)
				duration := int64(v.DurationMs / 1000)
				if len(tweet.Entities.Media) == 1 && v.IsAnimatedGif {
					inputMedia = append(inputMedia, gotgbot.InputMediaAnimation{
						Media:     media,
						Caption:   c,
						Width:     width,
						Height:    height,
						Duration:  duration,
						ParseMode: "MarkdownV2",
					})
				} else {
					inputMedia = append(inputMedia, gotgbot.InputMediaVideo{
						Media:     media,
						Cover:     v.ThumbUrl,
						Caption:   c,
						Width:     width,
						Height:    height,
						Duration:  duration,
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

type FileTooLargeError struct{}

func (e *FileTooLargeError) Error() string {
	return "file too large"
}

func downloadToBuffer(url, fn string) (gotgbot.InputFile, error) {
	defaultClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := defaultClient.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download file")
	}

	if resp.ContentLength > 50*1024*1024 {
		return nil, &FileTooLargeError{}
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return gotgbot.InputFileByReader(fn, buf), nil
}
