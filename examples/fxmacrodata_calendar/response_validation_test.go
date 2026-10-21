// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type deadlineTransport func(*http.Request) (*http.Response, error)

func (transport deadlineTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func TestFetchCalendarHasBoundedRequestDeadline(t *testing.T) {
	original := http.DefaultTransport
	http.DefaultTransport = deadlineTransport(func(request *http.Request) (*http.Response, error) {
		deadline, ok := request.Context().Deadline()
		if !ok || time.Until(deadline) <= 0 || time.Until(deadline) > time.Minute {
			t.Error("expected a positive request deadline of at most one minute")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"data":[]}`)), Header: make(http.Header)}, nil
	})
	t.Cleanup(func() { http.DefaultTransport = original })
	if _, err := fetchCalendar("USD", "2026-07-01", "2026-07-20"); err != nil {
		t.Fatal(err)
	}
}

func calendarBody(t *testing.T, body string) ([]calendarEvent, error) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(writer, body)
	}))
	defer server.Close()
	original := apiHost
	apiHost = server.URL
	defer func() { apiHost = original }()
	return fetchCalendar("USD", "2026-07-01", "2026-07-20")
}

func TestFetchCalendarRejectsErrorWithData(t *testing.T) {
	_, err := calendarBody(t, `{"error":"private-provider-marker","data":[]}`)
	if err == nil || strings.Contains(err.Error(), "private-provider-marker") {
		t.Fatal("expected a sanitized error for an error envelope")
	}
}

func TestFetchCalendarRejectsDetailWithData(t *testing.T) {
	_, err := calendarBody(t, `{"detail":{"message":"private-provider-marker"},"data":[]}`)
	if err == nil || strings.Contains(err.Error(), "private-provider-marker") {
		t.Fatal("expected a sanitized error for a detail envelope")
	}
}

func TestFetchCalendarRejectsMissingData(t *testing.T) {
	if _, err := calendarBody(t, `{}`); err == nil {
		t.Fatal("expected an error for missing data")
	}
}

func TestFetchCalendarRejectsNullData(t *testing.T) {
	if _, err := calendarBody(t, `{"data":null}`); err == nil {
		t.Fatal("expected an error for null data")
	}
}

func TestFetchCalendarRejectsNullDocument(t *testing.T) {
	if _, err := calendarBody(t, `null`); err == nil {
		t.Fatal("expected an error for a null document")
	}
}

func TestFetchCalendarRejectsTrailingDocument(t *testing.T) {
	if _, err := calendarBody(t, `{"data":[]} {"data":[]}`); err == nil {
		t.Fatal("expected an error for a second JSON document")
	}
}

func TestFetchCalendarRejectsTrailingText(t *testing.T) {
	if _, err := calendarBody(t, `{"data":[]} invalid`); err == nil {
		t.Fatal("expected an error for trailing invalid text")
	}
}

func TestFetchCalendarPreservesGenuineEmptyData(t *testing.T) {
	rows, err := calendarBody(t, `{"data":[]}`)
	if err != nil || rows == nil || len(rows) != 0 {
		t.Fatal("expected a genuine empty data list to succeed")
	}
}
