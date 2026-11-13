// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func useCalendarHost(t *testing.T, host string) {
	t.Helper()
	previous := apiHost
	apiHost = host
	t.Cleanup(func() { apiHost = previous })
}

func TestCalendarCurrencyPathEscaping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/api/v1/calendar/US%2FD%3Fextra=1%23part" ||
			r.URL.Query().Get("extra") != "" || r.URL.Query().Get("start_date") != "2026-07-01" ||
			r.URL.Query().Get("end_date") != "2026-07-20" {
			t.Error("currency must remain one encoded path segment with the original date query")
		}
		_, _ = io.WriteString(w, `{"data":[]}`)
	}))
	defer server.Close()
	useCalendarHost(t, server.URL)
	if _, err := fetchCalendar("US/D?extra=1#part", "2026-07-01", "2026-07-20"); err != nil {
		t.Fatal(err)
	}
}

func TestCalendarRejectsUnsafeAPIHost(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		_, _ = io.WriteString(w, `{"data":[]}`)
	}))
	defer server.Close()
	previous := apiHost
	t.Cleanup(func() { apiHost = previous })
	const marker = "private-provider-marker"
	for name, host := range map[string]string{
		"empty":          "",
		"non-HTTP":       "ftp://127.0.0.1/v1",
		"missing host":   "https:///v1",
		"credentials":    strings.Replace(server.URL, "http://", "http://user:"+marker+"@", 1),
		"query":          server.URL + "?token=" + marker,
		"fragment":       server.URL + "#" + marker,
		"control":        server.URL + "\n",
		"invalid escape": server.URL + "/%xx",
	} {
		t.Run(name, func(t *testing.T) {
			apiHost = host
			before := calls.Load()
			_, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
			if err == nil || err.Error() != "calendar API host must be HTTP(S), with a host and no credentials, query or fragment" || errors.Unwrap(err) != nil {
				t.Error("expected a sanitized invalid API-host error")
			}
			if calls.Load() != before {
				t.Error("invalid host reached the transport")
			}
		})
	}
}

func TestCalendarRejectsRedirect(t *testing.T) {
	var followed atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		followed.Add(1)
		_, _ = io.WriteString(w, `{"data":[]}`)
	}))
	defer target.Close()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer origin.Close()
	useCalendarHost(t, origin.URL)
	_, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
	if err == nil || err.Error() != "calendar request failed with HTTP 302" || followed.Load() != 0 {
		t.Fatal("expected a numeric redirect error without following the redirect")
	}
}

func TestCalendarNetworkDiagnosticIsSanitized(t *testing.T) {
	previous := http.DefaultTransport
	http.DefaultTransport = deadlineTransport(func(_ *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("private-provider-marker")
	})
	t.Cleanup(func() { http.DefaultTransport = previous })
	_, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
	if err == nil || err.Error() != "calendar request failed" || errors.Unwrap(err) != nil {
		t.Fatal("expected a sanitized transport error without a provider cause")
	}
}

func TestCalendarStatusReasonIsSanitized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		connection, buffer, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		defer connection.Close()
		_, _ = buffer.WriteString("HTTP/1.1 500 private-provider-marker\r\nContent-Length: 0\r\n\r\n")
		_ = buffer.Flush()
	}))
	defer server.Close()
	useCalendarHost(t, server.URL)
	_, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
	if err == nil || err.Error() != "calendar request failed with HTTP 500" {
		t.Fatal("expected numeric HTTP status without the provider's reason phrase")
	}
}

func TestCalendarDecodeDiagnosticIsSanitized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"data":[{"market_tier":18446744073709551616}]}`)
	}))
	defer server.Close()
	useCalendarHost(t, server.URL)
	_, err := fetchCalendar("USD", "2026-07-01", "2026-07-20")
	if err == nil || err.Error() != "calendar response could not be decoded" || errors.Unwrap(err) != nil {
		t.Fatal("expected a sanitized JSON decoding error without the raw value or cause")
	}
}

func TestCalendarTrailingSlashHost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/calendar/USD" {
			t.Error("expected one slash between the configured base and endpoint")
		}
		_, _ = io.WriteString(w, `{"data":[]}`)
	}))
	defer server.Close()
	useCalendarHost(t, server.URL+"/")
	if _, err := fetchCalendar("USD", "2026-07-01", "2026-07-20"); err != nil {
		t.Fatal(err)
	}
}
