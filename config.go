package main

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type Config struct {
	DatabaseUrl          string
	TwitterCookie        string
	XCsrfToken           string
	TelegramBotToken     string
	ChannelChatID        int64
	GroupChatID          int64
	OwnerID              int64
	PopularTweetFactor   int
	PopularRetweetFactor int
	BotApiUrl            string
}

func loadConfig() (*Config, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	twitterCookie := os.Getenv("TWITTER_COOKIE")
	if twitterCookie == "" {
		return nil, errors.New("TWITTER_COOKIE is not set")
	}

	xcsrfToken := os.Getenv("XCSRF_TOKEN")
	if xcsrfToken == "" {
		return nil, errors.New("XCSRF_TOKEN is not set")
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}

	channelChatIDStr := os.Getenv("CHANNEL_CHAT_ID")
	if channelChatIDStr == "" {
		return nil, errors.New("CHANNEL_CHAT_ID is not set")
	}
	channelChatID, err := strconv.ParseInt(channelChatIDStr, 10, 64)
	if channelChatID == 0 || err != nil {
		return nil, errors.Wrap(err, "CHANNEL_CHAT_ID is not a number")
	}

	groupChatIDStr := os.Getenv("GROUP_CHAT_ID")
	if groupChatIDStr == "" {
		return nil, errors.New("GROUP_CHAT_ID is not set")
	}
	groupChatID, err := strconv.ParseInt(groupChatIDStr, 10, 64)
	if groupChatID == 0 || err != nil {
		return nil, errors.Wrap(err, "GROUP_CHAT_ID is not a number")
	}

	ownerIDStr := os.Getenv("OWNER_ID")
	if ownerIDStr == "" {
		return nil, errors.New("OWNER_ID is not set")
	}
	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 64)
	if ownerID == 0 || err != nil {
		return nil, errors.Wrap(err, "OWNER_ID is not a number")
	}

	popularTweetFactorStr := os.Getenv("POPULAR_TWEET_FACTOR")
	if popularTweetFactorStr == "" {
		return nil, errors.New("POPULAR_TWEET_FACTOR is not set")
	}
	popularTweetFactor, err := strconv.Atoi(popularTweetFactorStr)
	if popularTweetFactor == 0 || err != nil {
		return nil, errors.Wrap(err, "POPULAR_TWEET_FACTOR is not a number")
	}

	popularRetweetFactorStr := os.Getenv("POPULAR_RETWEET_FACTOR")
	if popularRetweetFactorStr == "" {
		return nil, errors.New("POPULAR_RETWEET_FACTOR is not set")
	}
	popularRetweetFactor, err := strconv.Atoi(popularRetweetFactorStr)
	if popularRetweetFactor == 0 || err != nil {
		return nil, errors.Wrap(err, "POPULAR_RETWEET_FACTOR is not a number")
	}
	botApiUrl := os.Getenv("BOT_API_URL")

	return &Config{
		DatabaseUrl:          databaseUrl,
		TwitterCookie:        twitterCookie,
		XCsrfToken:           xcsrfToken,
		TelegramBotToken:     telegramBotToken,
		ChannelChatID:        channelChatID,
		GroupChatID:          groupChatID,
		OwnerID:              ownerID,
		PopularTweetFactor:   popularTweetFactor,
		PopularRetweetFactor: popularRetweetFactor,
		BotApiUrl:            botApiUrl,
	}, nil
}
