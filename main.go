package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

var UkraineLocation *time.Location

const GROUP = "КН-24"

func init() {
	loc, err := time.LoadLocation("Europe/Kyiv")
	if err != nil {
		panic(err)
	}
	UkraineLocation = loc
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	botSettings := tele.Settings{
		Token:  os.Getenv("OBLIVIONE_TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(botSettings)
	if err != nil {
		logger.Error("failed to init telegram bot", slog.String("err", err.Error()))
		os.Exit(1)
	}

	commands := []tele.Command{
		{Text: "/week", Description: "Get schedule for current week"},
	}

	err = bot.SetCommands(commands)
	if err != nil {
		logger.Error("failed to set bot commands", slog.String("err", err.Error()))
	}

	client := ScheduleClient{Client: http.DefaultClient}
	provider := NewScheduleProvider(client, time.Hour)

	bot.Handle("/week", func(c tele.Context) error {
		startDate := time.Now().Add(-time.Duration(time.Now().Weekday()) * time.Hour * 24)
		endDate := startDate.Add(time.Hour * 24 * 5)

		workingDays, err := provider.GetSchedule(ScheduleRequest{Start: startDate, End: endDate, Group: GROUP})
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
	}, NewLogginsMiddleware(logger))

	bot.Start()
}

func NewLogginsMiddleware(logger *slog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			var userID int64
			var userName string = "<no name>"
			if ctx.Sender() != nil {
				userID = ctx.Sender().ID
				userName = ctx.Sender().FirstName
			}

			logger.Info("request processed",
				slog.Int64("userID", userID),
				slog.String("userName", userName),
				slog.String("command", ctx.Text()),
			)

			return next(ctx)
		}
	}
}
