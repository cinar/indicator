// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"context"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

type strategyWithoutContext struct{}

func (s *strategyWithoutContext) Name() string {
	return "strategyWithoutContext"
}

func (s *strategyWithoutContext) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	actions := make(chan strategy.Action)
	go func() {
		defer close(actions)
		for range snapshots {
			actions <- strategy.Buy
		}
	}()
	return actions
}

func (s *strategyWithoutContext) Report(snapshots <-chan *asset.Snapshot) *helper.Report {
	return &helper.Report{}
}

func TestComputeStrategyWithContextFallback(t *testing.T) {
	date1, _ := time.Parse("2006-01-02", "2021-01-01")
	date2, _ := time.Parse("2006-01-02", "2021-01-02")

	snapshots := helper.SliceToChan([]*asset.Snapshot{
		{Date: date1, Close: 100},
		{Date: date2, Close: 101},
	})

	s := &strategyWithoutContext{}
	actions := strategy.ComputeStrategyWithContext(context.Background(), s, snapshots)

	count := 0
	for range actions {
		count++
	}

	if count != 2 {
		t.Fatalf("expected 2 actions, got %d", count)
	}
}
