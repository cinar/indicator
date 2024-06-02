// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
)

func TestTiingoRepositoryAssets(t *testing.T) {
	repository := asset.NewTiingoRepository("1234")

	_, err := repository.Assets()
	if err != errors.ErrUnsupported {
		t.Fatal(err)
	}
}

func TestTiingoRepositoryGet(t *testing.T) {
	data := []asset.TiingoEndOfDay{
		{
			Date:      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			AdjOpen:   10,
			AdjHigh:   30,
			AdjLow:    5,
			AdjClose:  20,
			AdjVolume: 100,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		_, err = w.Write(body)
		if err != nil {
			t.Fatal(err)
		}
	}))

	repository := asset.NewTiingoRepository("1234")
	repository.BaseURL = server.URL

	snapshots, err := repository.Get("A")
	if err != nil {
		t.Fatal(err)
	}

	expected := data[0].ToSnapshot()
	actual := <-snapshots

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestTiingoRepositoryGetNotReachable(t *testing.T) {
	repository := asset.NewTiingoRepository("1234")
	repository.BaseURL = "abcd://a.b.c.d"

	_, err := repository.Get("A")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTiingoRepositoryGetInvalid(t *testing.T) {
	response := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, response)
	}))

	repository := asset.NewTiingoRepository("1234")
	repository.BaseURL = server.URL

	_, err := repository.Get("A")
	if err != nil {
		t.Fatal(err)
	}

	response = "["

	_, err = repository.Get("A")
	if err != nil {
		t.Fatal(err)
	}

	response = "[{}"

	_, err = repository.Get("A")
	if err != nil {
		t.Fatal(err)
	}

	response = "[0]"

	_, err = repository.Get("A")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTiingoRepositoryLastDate(t *testing.T) {
	meta := asset.TiingoMeta{
		Ticker:       "A",
		Name:         "N",
		ExchangeCode: "E",
		StartDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		Description:  "D",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(meta)
		if err != nil {
			t.Fatal(err)
		}

		_, err = w.Write(body)
		if err != nil {
			t.Fatal(err)
		}
	}))

	repository := asset.NewTiingoRepository("1234")
	repository.BaseURL = server.URL

	lastDate, err := repository.LastDate("A")
	if err != nil {
		t.Fatal(err)
	}

	if lastDate != meta.EndDate {
		t.Fatalf("actual %v expected %v", lastDate, meta.EndDate)
	}
}

func TestTiingoRepositoryLastDateInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	repository := asset.NewTiingoRepository("1234")
	repository.BaseURL = server.URL

	_, err := repository.LastDate("A")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTiingoRepositoryLastDateNotReachable(t *testing.T) {
	repository := asset.NewTiingoRepository("1234")
	repository.BaseURL = "abcd://a.b.c.d"

	_, err := repository.LastDate("A")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTiingoRepositoryAppend(t *testing.T) {
	repository := asset.NewTiingoRepository("1234")

	err := repository.Append("A", nil)
	if err != errors.ErrUnsupported {
		t.Fatal(err)
	}
}
