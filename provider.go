package main

import (
	"bytes"
	"time"

	"golang.org/x/net/html"
)

var parsingPath = []string{"body", "#wrap", "div", "div", "div", ".container", ".row", "div"}

type ScheduleProvider struct {
	client ScheduleClient
}

func NewScheduleProvider(client ScheduleClient) *ScheduleProvider {
	return &ScheduleProvider{client: client}
}

func (sp *ScheduleProvider) GetSchedule(start, end time.Time) ([]WorkingDay, error) {
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
