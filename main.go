package main

import (
	"log/slog"
	"net/http"
	"os"
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
	provider := NewScheduleProvider(client, time.Hour, logger)

	handlers := NewHandlers(provider)

	bot.Handle("/week", handlers.weekScheduleHandler, NewLogginsMiddleware(logger))
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
