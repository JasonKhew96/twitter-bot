package main

import (
	"log"
	"time"

	_ "modernc.org/sqlite"
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
	count, err := bot.newLoop()
	if err != nil {
		log.Fatal(err)
	}
	wait := 5 * time.Minute
	if count == 0 {
		wait = 15 * time.Minute
	}

	time.AfterFunc(wait, bot.loop)
}
