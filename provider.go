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

func (sp *ScheduleProvider) GetSchedule(request ScheduleRequest) ([]WorkingDay, error) {
	cacheKey := request.Start.Format("02.01") + "-" + request.End.Format("02.01")
	if workingDays, exists := sp.cache.Get(cacheKey); exists {
		return workingDays, nil
	}

	workingDays, err := sp.fetchSchedule(request)
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
