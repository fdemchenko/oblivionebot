package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var ErrNoLesson = errors.New("there is no lesson")

type Class struct {
	StartTime, EndTime time.Time
	Number             int
	Lecturer           []string
	Title              []string
	Room               []string
	Groups             []string
}

type WorkingDay struct {
	DayOfWeek     time.Time
	DayOfWeekName string
	Classes       []Class
}

func (wd *WorkingDay) String() string {
	var workingDayString strings.Builder

	today := time.Now().In(UkraineLocation)

	workingDayString.WriteString(wd.DayOfWeekName)
	if truncateToDays(today).Equal(truncateToDays(wd.DayOfWeek)) {
		workingDayString.WriteString("    <- сьогодні")
	}
	workingDayString.WriteRune('\n')

	for i, lesson := range wd.Classes {
		var classEmoji string
		if rand.Intn(2) == 0 {
			classEmoji = "💻"
		} else {
			classEmoji = "📚"
		}
		startTime := lesson.StartTime.Format("15:04")
		endTime := lesson.EndTime.Format("15:04")
		if len(lesson.Lecturer) == 2 {
			workingDayString.WriteString(fmt.Sprintf("%s%s-%s:\n", classEmoji, startTime, endTime))
			workingDayString.WriteString(fmt.Sprintf("    %s: %s, %s, %s\n", lesson.Groups[0], lesson.Title[0], lesson.Lecturer[0], lesson.Room[0]))
			workingDayString.WriteString(fmt.Sprintf("    %s: %s, %s, %s", lesson.Groups[1], lesson.Title[1], lesson.Lecturer[1], lesson.Room[1]))
		} else {
			workingDayString.WriteString(fmt.Sprintf("%s%s-%s: %s\n", classEmoji, startTime, endTime, lesson.Title[0]))
			if lesson.Groups[0] != "" {
				workingDayString.WriteString(fmt.Sprintf("%s, %s (%s)", lesson.Lecturer[0], lesson.Room[0], lesson.Groups[0]))
			} else {
				workingDayString.WriteString(fmt.Sprintf("%s, %s", lesson.Lecturer[0], lesson.Room[0]))
			}
		}
		workingDayString.WriteRune('\n')
		if i < len(wd.Classes)-1 {
			workingDayString.WriteString("⎯⎯⎯⎯⎯⎯⎯⎯⎯")
		}
		workingDayString.WriteRune('\n')
	}
	return workingDayString.String()
}

type ScheduleRequest struct {
	Start time.Time
	End   time.Time
	Group string
}

func truncateToDays(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
