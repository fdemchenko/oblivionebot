package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	startDate := time.Now().Add(-time.Duration(time.Now().Weekday()) * time.Hour * 24)
	endDate := startDate.Add(time.Hour * 24 * 5)

	client := ScheduleClient{Client: http.DefaultClient}
	provider := NewScheduleProvider(client)

	bot.Handle("/week", func(c tele.Context) error {
		workingDays, err := provider.GetSchedule(startDate, endDate)
		if err != nil {
			return c.Send("На жаль відбулася помилка, неможливо отримати розклад.")
		}

		var message strings.Builder
		for _, day := range workingDays {
			message.WriteString(fmt.Sprintf("%s\n", day.DayOfWeekName))
			for _, lesson := range day.Classes {
				startTime := lesson.StartTime.Format("15:04")
				endTime := lesson.EndTime.Format("15:04")
				message.WriteString(fmt.Sprintf("%s-%s: %s\n", startTime, endTime, lesson.Title))
				message.WriteString(fmt.Sprintf("%s, %s\n", lesson.Lecturer, lesson.Room))
			}
			message.WriteRune('\n')
		}

		return c.Send(message.String())
	})

	bot.Start()
}
