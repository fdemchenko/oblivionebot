package main

import (
	"fmt"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Provider interface {
	GetSchedule(ScheduleRequest) ([]WorkingDay, error)
}

type Handlers struct {
	scheduleProvider Provider
}

func NewHandlers(provider Provider) *Handlers {
	return &Handlers{scheduleProvider: provider}
}

func (h *Handlers) weekScheduleHandler(ctx tele.Context) error {
	startDate := time.Now().Add(-time.Duration(time.Now().Weekday()) * time.Hour * 24)
	endDate := startDate.Add(time.Hour * 24 * 5)

	workingDays, err := h.scheduleProvider.GetSchedule(ScheduleRequest{Start: startDate, End: endDate, Group: GROUP})
	if err != nil {
		return ctx.Send("На жаль відбулася помилка, неможливо отримати розклад.")
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

	return ctx.Send(message.String())
}
