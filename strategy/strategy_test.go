// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"context"
	"runtime"
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

func TestActionSourcesWithContext(t *testing.T) {
	date1, _ := time.Parse("2006-01-02", "2021-01-01")
	date2, _ := time.Parse("2006-01-02", "2021-01-02")

	snapshots := helper.SliceToChan([]*asset.Snapshot{
		{Date: date1, Close: 100},
		{Date: date2, Close: 101},
	})

	strategies := []strategy.Strategy{
		strategy.NewBuyAndHoldStrategy(),
		strategy.NewBuyAndHoldStrategy(),
	}

	sources := strategy.ActionSourcesWithContext(context.Background(), strategies, snapshots)

	if len(sources) != len(strategies) {
		t.Fatalf("expected %d sources, got %d", len(strategies), len(sources))
	}

	for _, source := range sources {
		count := 0
		for range source {
			count++
		}

		if count != 2 {
			t.Fatalf("expected 2 actions, got %d", count)
		}
	}
}

func TestActionSourcesWithContextCancellation(t *testing.T) {
	runtime.GC()
	baseline := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	snapshots := make(chan *asset.Snapshot)

	strategies := []strategy.Strategy{
		strategy.NewBuyAndHoldStrategy(),
		strategy.NewBuyAndHoldStrategy(),
	}

	sources := strategy.ActionSourcesWithContext(ctx, strategies, snapshots)

	cancel()

	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	current := runtime.NumGoroutine()
	if current > baseline+2 {
		t.Fatalf("Goroutine leak detected. Baseline: %d, Current: %d", baseline, current)
	}

	for _, source := range sources {
		if _, ok := <-source; ok {
			t.Fatal("Source channel should be closed after cancellation")
		}
	}
}
