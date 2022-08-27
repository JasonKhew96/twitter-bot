package main

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("start")
	bot, err := New()
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	go bot.worker()

	go bot.loop()

	bot.initBot()
}

func (bot *bot) loop() {
	if err := bot.cleanup(); err != nil {
		log.Fatal(err)
	}
	if err := bot.newLoop(); err != nil {
		log.Fatal(err)
	}
	time.AfterFunc(5*time.Minute, bot.loop)
}
