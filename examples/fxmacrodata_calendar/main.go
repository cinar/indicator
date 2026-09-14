// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// Package main provides a self-contained example of fetching and filtering
// macroeconomic announcement events from the FXMacroData public release-calendar endpoint.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

type calendarResponse struct {
	Data []calendarEvent `json:"data"`
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
	endpoint := fmt.Sprintf("%s/api/v1/calendar/%s", apiHost, currency)
	params := url.Values{}
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)

	res, err := http.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with %s", res.Status)
	}

	var payload calendarResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
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
