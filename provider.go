package main

import (
	"bytes"
	"time"

	"golang.org/x/net/html"
)

var parsingPath = []string{"body", "#wrap", "div", "div", "div", ".container", ".row", "div"}

type ScheduleProvider struct {
	client   ScheduleClient
	cache    *Cache[string, []WorkingDay]
	cacheTTL time.Duration
}

func NewScheduleProvider(client ScheduleClient, cacheTTL time.Duration) *ScheduleProvider {
	return &ScheduleProvider{
		client:   client,
		cache:    NewCache[string, []WorkingDay](),
		cacheTTL: cacheTTL,
	}
}

func (sp *ScheduleProvider) GetSchedule(start, end time.Time) ([]WorkingDay, error) {
	cacheKey := start.Format("02.01") + "-" + end.Format("02.01")
	if workingDays, exists := sp.cache.Get(cacheKey); exists {
		return workingDays, nil
	}

	workingDays, err := sp.fetchSchedule(start, end)
	if err != nil {
		return nil, err
	}

	sp.cache.Set(cacheKey, workingDays, sp.cacheTTL)
	return workingDays, nil
}

func (sp *ScheduleProvider) fetchSchedule(start, end time.Time) ([]WorkingDay, error) {
	htmlBytes, err := sp.client.GetScheduleHTML(start, end)
	if err != nil {
		return nil, err
	}

	document, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, err
	}

	workingDayNodes := getNodesByTreePath(document.FirstChild.NextSibling, parsingPath)
	var workingDays []WorkingDay
	for _, worworkingDayNode := range workingDayNodes {
		workingDay, err := processWorkingDayNode(worworkingDayNode)
		if err != nil {
			return nil, err
		}
		workingDays = append(workingDays, *workingDay)
	}

	return workingDays, nil
}
