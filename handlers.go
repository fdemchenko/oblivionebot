package main

import (
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

const ERROR_MESSAGE = "На жаль відбулася помилка, неможливо отримати розклад😔"

type Provider interface {
	GetWeekSchedule(group string) ([]WorkingDay, error)
	GetDaySchedule(day time.Time, group string) (*WorkingDay, error)
}

type Handlers struct {
	scheduleProvider Provider
}

func NewHandlers(provider Provider) *Handlers {
	return &Handlers{scheduleProvider: provider}
}

func (h *Handlers) weekScheduleHandler(ctx tele.Context) error {
	workingDays, err := h.scheduleProvider.GetWeekSchedule(GROUP)
	if err != nil {
		return ctx.Send(ERROR_MESSAGE)
	}

	if len(workingDays) == 0 {
		return ctx.Send("Пар немає🥳")
	}

	var message strings.Builder
	for _, day := range workingDays {
		message.WriteString(day.String())
		message.WriteRune('\n')
	}

	return ctx.Send(message.String())
}

func (h *Handlers) todaySchedulehandler(ctx tele.Context) error {
	workingDay, err := h.scheduleProvider.GetDaySchedule(time.Now(), GROUP)
	if err != nil {
		return ctx.Send(ERROR_MESSAGE)
	}

	if workingDay == nil {
		return ctx.Send("Сьогодні пар немає🥳")
	}

	return ctx.Send(workingDay.String())
}

func (h *Handlers) tomorrowSchedulehandler(ctx tele.Context) error {
	workingDay, err := h.scheduleProvider.GetDaySchedule(time.Now().Add(time.Hour*24).In(UkraineLocation), GROUP)
	if err != nil {
		return ctx.Send(ERROR_MESSAGE)
	}

	if workingDay == nil {
		return ctx.Send("Завтра пар немає🥳")
	}

	return ctx.Send(workingDay.String())
}
