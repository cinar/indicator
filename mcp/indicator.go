package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/registry"
)

// IndicatorRequest defines the structure for an indicator request: which
// indicator to compute, its optional parameters, and the OHLCV data to
// compute it over.
type IndicatorRequest struct {
	Indicator string `json:"indicator"`

	// Params are optional "Name=Value" overrides for the indicator's
	// parameters, e.g. ["Period=14"]. Unset parameters keep their
	// default value.
	Params []string `json:"params,omitempty"`

	Data OhlcvData `json:"data"`
}

// IndicatorResponse defines the structure of the JSON response for an
// indicator computation: the dates the values line up with, and each of
// the indicator's named output series.
type IndicatorResponse struct {
	Dates  []string             `json:"dates"`
	Series map[string][]float64 `json:"series"`
}

// parseIndicatorParams converts "Name=Value" strings into the map that
// registry.Run expects.
func parseIndicatorParams(params []string) (map[string]string, error) {
	parsed := make(map[string]string, len(params))

	for _, param := range params {
		name, value, ok := strings.Cut(param, "=")
		if !ok {
			return nil, fmt.Errorf("invalid param %q, expected Name=Value", param)
		}

		parsed[name] = value
	}

	return parsed, nil
}

// runIndicator computes the requested indicator over the given OHLCV data
// and collects it into a JSON-friendly response.
//
// It shares the indicator lookup, parameter binding, and computation logic
// with the indicator-cli command through the registry package, so both
// front ends stay in sync as indicators are added.
func runIndicator(ctx context.Context, req IndicatorRequest) (*IndicatorResponse, error) {
	dateArray := toTimeArray(req.Data.Date)

	if len(req.Data.Opening) != len(dateArray) || len(req.Data.Closing) != len(dateArray) ||
		len(req.Data.High) != len(dateArray) || len(req.Data.Low) != len(dateArray) ||
		len(req.Data.Volume) != len(dateArray) {
		return nil, fmt.Errorf("opening, closing, high, low, and volume must all have the same length as date, got %d, %d, %d, %d, %d, %d",
			len(dateArray), len(req.Data.Opening), len(req.Data.Closing), len(req.Data.High), len(req.Data.Low), len(req.Data.Volume))
	}

	params, err := parseIndicatorParams(req.Params)
	if err != nil {
		return nil, err
	}

	snapshots := make(chan *asset.Snapshot)

	go func() {
		defer close(snapshots)

		for i := range dateArray {
			snapshots <- &asset.Snapshot{
				Date:   dateArray[i],
				Open:   req.Data.Opening[i],
				Close:  req.Data.Closing[i],
				High:   req.Data.High[i],
				Low:    req.Data.Low[i],
				Volume: req.Data.Volume[i],
			}
		}
	}()

	result, err := registry.Run(ctx, req.Indicator, params, snapshots)
	if err != nil {
		return nil, err
	}

	response := &IndicatorResponse{
		Series: make(map[string][]float64, len(result.Series)),
	}

	for date := range result.Dates {
		response.Dates = append(response.Dates, date.Format("2006-01-02"))

		for _, series := range result.Series {
			value, ok := <-series.Values
			if !ok {
				return nil, fmt.Errorf("indicator %q produced fewer values than dates", req.Indicator)
			}

			response.Series[series.Name] = append(response.Series[series.Name], value)
		}
	}

	return response, nil
}
