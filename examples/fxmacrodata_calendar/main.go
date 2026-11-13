// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// Package main provides a self-contained example of fetching and filtering
// macroeconomic announcement events from the FXMacroData public release-calendar endpoint.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type calendarResponse struct {
	Data   []calendarEvent `json:"data"`
	Error  json.RawMessage `json:"error,omitempty"`
	Detail json.RawMessage `json:"detail,omitempty"`
}

type calendarEvent struct {
	Name                    string `json:"name"`
	Date                    string `json:"date"`
	AnnouncementDatetimeUTC string `json:"announcement_datetime_utc"`
	MarketTier              int    `json:"market_tier"`
	TopTierForCurrency      bool   `json:"top_tier_for_currency"`
}

func main() {
	events, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
	if err != nil {
		panic(err)
	}

	fmt.Println("Top-tier USD macro blackout dates:")
	for _, event := range topTierEvents(events) {
		fmt.Printf("  %s: %s\n", eventDate(event), event.Name)
	}
}

var apiHost = "https://fxmacrodata.com"

func fetchCalendar(currency, startDate, endDate string) ([]calendarEvent, error) {
	base, err := url.Parse(apiHost)
	if err != nil || base.Hostname() == "" || base.User != nil ||
		(base.Scheme != "https" && base.Scheme != "http") ||
		strings.ContainsAny(apiHost, "?#") {
		return nil, fmt.Errorf("calendar API host must be HTTP(S), with a host and no credentials, query or fragment")
	}
	endpoint := strings.TrimRight(base.String(), "/") + "/api/v1/calendar/" + url.PathEscape(currency)
	params := url.Values{}
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)

	client := http.Client{
		Timeout: 20 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("calendar request failed")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("calendar request failed with HTTP %d", res.StatusCode)
	}

	var payload calendarResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("calendar response could not be decoded")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("calendar response must contain one JSON value")
	}
	for _, field := range []json.RawMessage{payload.Error, payload.Detail} {
		if len(field) > 0 && !bytes.Equal(bytes.TrimSpace(field), []byte("null")) {
			return nil, fmt.Errorf("calendar response contains an error")
		}
	}
	if payload.Data == nil {
		return nil, fmt.Errorf("calendar response must contain a data list")
	}
	return payload.Data, nil
}

func topTierEvents(events []calendarEvent) []calendarEvent {
	filtered := make([]calendarEvent, 0)
	for _, event := range events {
		if event.TopTierForCurrency || event.MarketTier == 1 {
			filtered = append(filtered, event)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return eventDate(filtered[i]) < eventDate(filtered[j])
	})
	return filtered
}

func eventDate(event calendarEvent) string {
	value := event.AnnouncementDatetimeUTC
	if value == "" {
		value = event.Date
	}
	if len(value) > 10 {
		return value[:10]
	}
	return value
}
