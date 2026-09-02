// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// closeTrackingBody wraps an io.ReadCloser and records whether Close was called.
type closeTrackingBody struct {
	io.ReadCloser
	closed bool
}

func (b *closeTrackingBody) Close() error {
	b.closed = true
	return b.ReadCloser.Close()
}

// bodyTrackingTransport wraps a RoundTripper and replaces the last response's
// body with a closeTrackingBody, so tests can assert it was closed.
type bodyTrackingTransport struct {
	base http.RoundTripper
	body *closeTrackingBody
}

func (t *bodyTrackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	t.body = &closeTrackingBody{ReadCloser: res.Body}
	res.Body = t.body

	return res, nil
}

// TestTiingoRepositoryLastDateClosesBodyOnStatusError verifies that
// LastDate closes the response body even when it returns early because
// of a non-200 status code, i.e. before the body is ever read.
func TestTiingoRepositoryLastDateClosesBodyOnStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	transport := &bodyTrackingTransport{base: http.DefaultTransport}

	repository := NewTiingoRepository("1234")
	repository.BaseURL = server.URL
	repository.client = &http.Client{Transport: transport}

	_, err := repository.LastDate("A")
	if err == nil {
		t.Fatal("expected error")
	}

	if transport.body == nil {
		t.Fatal("expected a tracked response body")
	}

	if !transport.body.closed {
		t.Fatal("expected response body to be closed")
	}
}

// TestTiingoRepositoryGetSinceClosesBodyOnDecodeError verifies that
// GetSince's background goroutine closes the response body even when
// decoding fails on the very first token, one of the early-return paths
// that previously left the body open.
func TestTiingoRepositoryGetSinceClosesBodyOnDecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	transport := &bodyTrackingTransport{base: http.DefaultTransport}

	repository := NewTiingoRepository("1234")
	repository.BaseURL = server.URL
	repository.client = &http.Client{Transport: transport}

	snapshots, err := repository.GetSince("A", time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	// Drain the channel to completion. Because the body is closed via a
	// defer declared after (and therefore run before) close(snapshots),
	// observing the channel close here guarantees the body is already
	// closed.
	for range snapshots {
	}

	if transport.body == nil {
		t.Fatal("expected a tracked response body")
	}

	if !transport.body.closed {
		t.Fatal("expected response body to be closed")
	}
}
