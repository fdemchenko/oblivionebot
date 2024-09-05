package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

var UkraineLocation *time.Location

func init() {
	loc, err := time.LoadLocation("Europe/Kyiv")
	if err != nil {
		panic(err)
	}
	UkraineLocation = loc
}

func main() {
	botSettings := tele.Settings{
		Token:  os.Getenv("OBLIVIONE_TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(botSettings)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/week", func(c tele.Context) error {
		return c.Send(c.Message().Text)
	})

	bot.Start()
}
