// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// TiingoMeta is the response from the meta endpoint.
// https://www.tiingo.com/documentation/end-of-day
type TiingoMeta struct {
	// Ticker related to the asset.
	Ticker string `json:"ticker"`

	// Name is the full name of the asset.
	Name string `json:"name"`

	// ExchangeCode is the exchange where the asset is listed on.
	ExchangeCode string `json:"exchangeCode"`

	// Description is the description of the asset.
	Description string `json:"description"`

	// StartDate is the earliest date for the asset data.
	StartDate time.Time `json:"startDate"`

	// EndDate is the latest date for the asset data.
	EndDate time.Time `json:"endDate"`
}

// TiingoEndOfDay is the repose from the end-of-day endpoint.
// https://www.tiingo.com/documentation/end-of-day
type TiingoEndOfDay struct {
	// Date is the date this data pertains to.
	Date time.Time `json:"date"`

	// Open is the opening price.
	Open float64 `json:"open"`

	// High is the highest price.
	High float64 `json:"high"`

	// Low is the lowest price.
	Low float64 `json:"low"`

	// Close is the closing price.
	Close float64 `json:"close"`

	// Volume is the total volume.
	Volume int64 `json:"volume"`

	// AdjOpen is the adjusted opening price.
	AdjOpen float64 `json:"adjOpen"`

	// AdjHigh is the adjusted highest price.
	AdjHigh float64 `json:"adjHigh"`

	// AdjLow is the adjusted lowest price.
	AdjLow float64 `json:"adjLow"`

	// AdjClose is the adjusted closing price.
	AdjClose float64 `json:"adjClose"`

	// AdjVolume is the adjusted total volume.
	AdjVolume int64 `json:"adjVolume"`

	// Dividend is the dividend paid out.
	Dividend float64 `json:"divCash"`

	// Split to adjust values after a split.
	Split float64 `json:"splitFactor"`
}

// ToSnapshot converts the Tiingo end-of-day to a snapshot.
func (e *TiingoEndOfDay) ToSnapshot() *Snapshot {
	return &Snapshot{
		Date:   e.Date,
		Open:   e.AdjOpen,
		High:   e.AdjHigh,
		Low:    e.AdjLow,
		Close:  e.AdjClose,
		Volume: float64(e.AdjVolume),
	}
}

// TiingoRepository provides access to financial market data, retrieving
// asset snapshots, by interacting with the Tiingo Stock & Financial
// Markets API. To use this repository, you'll need a valid API key
// from https://www.tiingo.com.
type TiingoRepository struct {
	Repository

	// apiKey is the Tiingo API key.
	apiKey string

	// Client is the HTTP client.
	client *http.Client

	// BaseURL is the Tiingo API URL.
	BaseURL string

	// Logger is the slog logger instance.
	Logger *slog.Logger
}

// NewTiingoRepository initializes a file system repository with
// the given API key.
func NewTiingoRepository(apiKey string) *TiingoRepository {
	return &TiingoRepository{
		apiKey:  apiKey,
		client:  &http.Client{},
		BaseURL: "https://api.tiingo.com",
		Logger:  slog.Default(),
	}
}

// Assets returns the names of all assets in the repository.
func (*TiingoRepository) Assets() ([]string, error) {
	return nil, errors.ErrUnsupported
}

// Get attempts to return a channel of snapshots for the asset with the given name.
func (r *TiingoRepository) Get(name string) (<-chan *Snapshot, error) {
	return r.GetSince(name, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
}

// GetSince attempts to return a channel of snapshots for the asset with the given name since the given date.
func (r *TiingoRepository) GetSince(name string, date time.Time) (<-chan *Snapshot, error) {
	url := fmt.Sprintf("%s/tiingo/daily/%s/prices?startDate=%s&token=%s",
		r.BaseURL,
		name,
		date.Format("2006-01-02"),
		r.apiKey)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with %s", res.Status)
	}

	snapshots := make(chan *Snapshot)

	go func() {
		defer close(snapshots)

		decoder := json.NewDecoder(res.Body)

		_, err = decoder.Token()
		if err != nil {
			r.Logger.Error("Unable to read token.", "error", err)
			return
		}

		for decoder.More() {
			var data TiingoEndOfDay

			err = decoder.Decode(&data)
			if err != nil {
				r.Logger.Error("Unable to decode data.", "error", err)
				break
			}

			snapshots <- data.ToSnapshot()
		}

		_, err = decoder.Token()
		if err != nil {
			r.Logger.Error("GetSince failed.", "error", err)
			return
		}

		err = res.Body.Close()
		if err != nil {
			r.Logger.Error("Unable to close respose.", "error", err)
		}
	}()

	return snapshots, nil
}

// LastDate returns the date of the last snapshot for the asset with the given name.
func (r *TiingoRepository) LastDate(name string) (time.Time, error) {
	var lastDate time.Time

	url := fmt.Sprintf("%s/tiingo/daily/%s?token=%s", r.BaseURL, name, r.apiKey)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return lastDate, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return lastDate, err
	}

	if res.StatusCode != 200 {
		return lastDate, fmt.Errorf("request failed with %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return lastDate, err
	}

	err = res.Body.Close()
	if err != nil {
		return lastDate, err
	}

	var meta TiingoMeta

	err = json.Unmarshal(body, &meta)
	if err != nil {
		return lastDate, err
	}

	return meta.EndDate, nil
}

// Append adds the given snapshows to the asset with the given name.
func (*TiingoRepository) Append(_ string, _ <-chan *Snapshot) error {
	return errors.ErrUnsupported
}
