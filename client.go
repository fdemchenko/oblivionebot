package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const urlScheduler = "http://r.polissiauniver.edu.ua/cgi-bin/timetable.cgi"

type ScheduleClient struct {
	Client *http.Client
}

func (sc *ScheduleClient) GetScheduleHTML(request ScheduleRequest) ([]byte, error) {
	encodedGroup, _, err := transform.Bytes(charmap.Windows1251.NewEncoder(), []byte(request.Group))
	if err != nil {
		return nil, err
	}

	formValues := url.Values{}
	formValues.Set("sdate", request.Start.Format("02.01.2006"))
	formValues.Set("edate", request.End.Format("02.01.2006"))

	requestData := strings.NewReader(fmt.Sprintf("group=%s&%s", url.QueryEscape(string(encodedGroup)), formValues.Encode()))

	req, err := http.NewRequest("POST", urlScheduler, requestData)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
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
