// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventDate(t *testing.T) {
	tests := []struct {
		event    calendarEvent
		expected string
	}{
		{
			event: calendarEvent{
				AnnouncementDatetimeUTC: "2026-07-14T12:30:00Z",
			},
			expected: "2026-07-14",
		},
		{
			event: calendarEvent{
				Date: "2026-07-15",
			},
			expected: "2026-07-15",
		},
		{
			event: calendarEvent{
				Date: "2026",
			},
			expected: "2026",
		},
	}

	for _, test := range tests {
		actual := eventDate(test.event)
		if actual != test.expected {
			t.Errorf("actual %s expected %s", actual, test.expected)
		}
	}
}

func TestTopTierEvents(t *testing.T) {
	events := []calendarEvent{
		{
			Name:       "Low Tier Event",
			Date:       "2026-07-15",
			MarketTier: 3,
		},
		{
			Name:       "High Tier Event 1",
			Date:       "2026-07-16",
			MarketTier: 1,
		},
		{
			Name:               "Top Tier Currency Event",
			Date:               "2026-07-14",
			MarketTier:         2,
			TopTierForCurrency: true,
		},
	}

	filtered := topTierEvents(events)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 events, got %d", len(filtered))
	}

	// Verify sorting and filtering
	if filtered[0].Name != "Top Tier Currency Event" {
		t.Errorf("expected first event to be Top Tier Currency Event, got %s", filtered[0].Name)
	}
	if filtered[1].Name != "High Tier Event 1" {
		t.Errorf("expected second event to be High Tier Event 1, got %s", filtered[1].Name)
	}
}

func TestFetchCalendar(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/calendar/USD" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("start_date") != "2026-07-01" {
			t.Errorf("unexpected start_date: %s", r.URL.Query().Get("start_date"))
		}
		if r.URL.Query().Get("end_date") != "2026-07-20" {
			t.Errorf("unexpected end_date: %s", r.URL.Query().Get("end_date"))
		}

		response := calendarResponse{
			Data: []calendarEvent{
				{
					Name:       "Core Inflation",
					Date:       "2026-07-14",
					MarketTier: 1,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	oldHost := apiHost
	apiHost = server.URL
	defer func() { apiHost = oldHost }()

	events, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
	if err != nil {
		t.Fatalf("failed to fetch calendar: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].Name != "Core Inflation" {
		t.Errorf("expected Core Inflation, got %s", events[0].Name)
	}
}

func TestFetchCalendarError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	oldHost := apiHost
	apiHost = server.URL
	defer func() { apiHost = oldHost }()

	_, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMainFunction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := calendarResponse{
			Data: []calendarEvent{
				{
					Name:       "Core Inflation",
					Date:       "2026-07-14",
					MarketTier: 1,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	oldHost := apiHost
	apiHost = server.URL
	defer func() { apiHost = oldHost }()

	// Verify main runs without panic when API returns successful result
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main panicked: %v", r)
		}
	}()
	main()
}
