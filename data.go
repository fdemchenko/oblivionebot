package main

import (
	"errors"
	"time"
)

var ErrNoLesson = errors.New("there is no lesson")

type Class struct {
	StartTime, EndTime time.Time
	Number             int
	Lecturer           string
	Title              string
	Room               string
	Groups             string
}

type WorkingDay struct {
	DayOfWeek     time.Time
	DayOfWeekName string
	Classes       []Class
}
