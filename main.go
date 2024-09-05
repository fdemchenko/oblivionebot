package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

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

	bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send(c.Message().Text)
	})

	bot.Start()
}
