package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const urlScheduler = "http://r.polissiauniver.edu.ua/cgi-bin/timetable.cgi"
const group = "КН-24"

type ScheduleClient struct {
	Client *http.Client
}

func (sc *ScheduleClient) GetScheduleHTML(start, end time.Time) ([]byte, error) {
	encodedGroup, _, err := transform.Bytes(charmap.Windows1251.NewEncoder(), []byte(group))
	if err != nil {
		log.Fatal(err)
	}

	formValues := url.Values{}
	formValues.Set("sdate", start.Format("02.01.2006"))
	formValues.Set("edate", end.Format("02.01.2006"))

	requestData := strings.NewReader(fmt.Sprintf("group=%s&%s", url.QueryEscape(string(encodedGroup)), formValues.Encode()))

	request, err := http.NewRequest("POST", urlScheduler, requestData)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("http error: status code - %d", res.StatusCode)
	}

	reader := transform.NewReader(res.Body, charmap.Windows1251.NewDecoder())

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
