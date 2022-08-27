package main

import (
	"net/url"
	"regexp"
	"strings"

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
