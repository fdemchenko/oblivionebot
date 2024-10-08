package main

import (
	"bytes"
	"log/slog"
	"time"

	"golang.org/x/net/html"
)

var parsingPath = []string{"body", "#wrap", "div", "div", "div", ".container", ".row", "div"}

type ScheduleProvider struct {
	client   ScheduleClient
	cache    *Cache[string, []WorkingDay]
	cacheTTL time.Duration
	logger   *slog.Logger
}

func NewScheduleProvider(client ScheduleClient, cacheTTL time.Duration, logger *slog.Logger) *ScheduleProvider {
	return &ScheduleProvider{
		client:   client,
		logger:   logger,
		cache:    NewCache[string, []WorkingDay](),
		cacheTTL: cacheTTL,
	}
}

func (sp *ScheduleProvider) GetDaySchedule(day time.Time, group string) (*WorkingDay, error) {
	weekSchedule, err := sp.GetWeekSchedule(group)
	if err != nil {
		return nil, err
	}

	for _, workingDay := range weekSchedule {
		if truncateToDays(truncateToDays(day.In(UkraineLocation))).Equal(truncateToDays(workingDay.DayOfWeek)) {
			return &workingDay, nil
		}
	}
	return nil, nil
}

func (sp *ScheduleProvider) GetWeekSchedule(group string) ([]WorkingDay, error) {
	startDate := time.Now().Add(-time.Duration((time.Now().In(UkraineLocation).Weekday()+1)%7) * time.Hour * 24)
	endDate := startDate.Add(time.Hour * 24 * 5)

	cacheKey := startDate.Format("02.01") + "-" + endDate.Format("02.01")
	if workingDays, exists := sp.cache.Get(cacheKey); exists {
		return workingDays, nil
	}

	workingDays, err := sp.fetchSchedule(ScheduleRequest{Start: startDate, End: endDate, Group: group})
	if err != nil {
		sp.logger.Error("failed to fetch schedule", slog.String("err", err.Error()))
		return nil, err
	}

	sp.cache.Set(cacheKey, workingDays, sp.cacheTTL)
	return workingDays, nil
}

func (sp *ScheduleProvider) fetchSchedule(request ScheduleRequest) ([]WorkingDay, error) {
	htmlBytes, err := sp.client.GetScheduleHTML(request)
	if err != nil {
		sp.logger.Error("failed to fetch schedule page markup", slog.String("err", err.Error()))
		return nil, err
	}

	document, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		sp.logger.Error("failed to parse html markup", slog.String("err", err.Error()))
		return nil, err
	}

	workingDayNodes := getNodesByTreePath(document.FirstChild.NextSibling, parsingPath)
	var workingDays []WorkingDay
	for _, worworkingDayNode := range workingDayNodes {
		workingDay, err := processWorkingDayNode(worworkingDayNode)
		if err != nil {
			sp.logger.Error("failed to process working day html node", slog.String("err", err.Error()))
			return nil, err
		}
		workingDays = append(workingDays, *workingDay)
	}

	return workingDays, nil
}
