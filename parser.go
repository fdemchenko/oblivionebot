package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Selectors struct {
	ID      string
	Classes []string
}

func getNodesByTreePath(root *html.Node, path []string) []*html.Node {
	var nodes []*html.Node

	var f func(*html.Node, []string)
	f = func(n *html.Node, targets []string) {
		if len(targets) == 0 {
			nodes = append(nodes, n)
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if checkNodeForTarget(c, targets[0]) {
				f(c, targets[1:])
			}
		}
	}
	f(root, path)
	return nodes
}

func checkNodeForTarget(node *html.Node, target string) bool {
	if node.Type != html.ElementNode {
		return false
	}

	selectors := Selectors{}

	for _, attr := range node.Attr {
		if attr.Key == "class" {
			selectors.Classes = strings.Split(attr.Val, " ")
		} else if attr.Key == "id" {
			selectors.ID = attr.Val
		}
	}

	if strings.HasPrefix(target, ".") {
		return slices.Contains(selectors.Classes, strings.TrimPrefix(target, "."))
	}
	if strings.HasPrefix(target, "#") {
		return strings.TrimPrefix(target, "#") == selectors.ID
	}
	return node.Data == target
}

func processWorkingDayNode(node *html.Node) (*WorkingDay, error) {
	dateString := node.FirstChild.FirstChild.Data
	dayString := node.FirstChild.LastChild.FirstChild.Data

	date, err := time.ParseInLocation("02.01.2006", strings.TrimSpace(dateString), UkraineLocation)
	if err != nil {
		return nil, err
	}

	var table *html.Node
	for table = node.FirstChild; table.Data != "table"; table = table.NextSibling {
	}

	var classes []Class
	for tableRow := table.FirstChild.FirstChild; tableRow != nil; tableRow = tableRow.NextSibling {
		lesson, err := processLessonNode(tableRow, dateString)
		if err != nil {
			if errors.Is(err, ErrNoLesson) {
				continue
			} else {
				return nil, err
			}
		}
		classes = append(classes, *lesson)
	}

	return &WorkingDay{DayOfWeek: date, Classes: classes, DayOfWeekName: dayString}, nil
}

func processLessonNode(node *html.Node, day string) (*Class, error) {
	classNoString := node.FirstChild.FirstChild.Data
	startTimeString := node.FirstChild.NextSibling.FirstChild.Data
	endTimeString := node.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.Data

	startTime, err := time.ParseInLocation("02.01.2006 15:04", fmt.Sprintf("%s %s", day, startTimeString), UkraineLocation)
	if err != nil {
		return nil, err
	}
	endTime, err := time.ParseInLocation("02.01.2006 15:04", fmt.Sprintf("%s %s", day, endTimeString), UkraineLocation)
	if err != nil {
		return nil, err
	}

	lessonNo, err := strconv.Atoi(classNoString)
	if err != nil {
		return nil, err
	}

	lessonInfoNode := node.FirstChild.NextSibling.NextSibling
	if strings.TrimSpace(lessonInfoNode.FirstChild.Data) == "" {
		return nil, ErrNoLesson
	}

	var infos []string

	for infoNode := lessonInfoNode.FirstChild; infoNode != nil; infoNode = infoNode.NextSibling {
		if infoNode.Data != "br" {
			infos = append(infos, strings.TrimSpace(infoNode.Data))
		}
	}
	var groups string
	if len(infos) > 3 {
		groups = infos[2]
	}

	return &Class{
		Number:    lessonNo,
		StartTime: startTime,
		EndTime:   endTime,
		Lecturer:  infos[len(infos)-1],
		Title:     infos[0],
		Room:      infos[1],
		Groups:    groups,
	}, nil
}
