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
	"github.com/cinar/indicator/v2/examples/trend"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestOrStrategy(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	results, err := helper.ReadFromCsvFile[strategy.Result]("testdata/or.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.Map(results, func(r *strategy.Result) strategy.Action { return r.Action })

	or := strategy.NewOrStrategy("Or Strategy")
	or.Strategies = append(or.Strategies, strategy.NewBuyAndHoldStrategy(), trend.NewMacdStrategy())

	actual := or.Compute(snapshots)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrStrategyNoStrategies(t *testing.T) {
	snapshots := helper.SliceToChan([]*asset.Snapshot{})

	or := strategy.NewOrStrategy("Or Strategy")
	actual := or.Compute(snapshots)

	select {
	case _, ok := <-actual:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout - result channel never closed")
	}
}

func TestOrStrategyCancellation(t *testing.T) {
	runtime.GC()
	baseline := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	snapshots := make(chan *asset.Snapshot)

	or := strategy.NewOrStrategy("Or Strategy")
	or.Strategies = append(or.Strategies, strategy.NewBuyAndHoldStrategy(), strategy.NewBuyAndHoldStrategy())

	actual := or.ComputeWithContext(ctx, snapshots)

	cancel()

	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	current := runtime.NumGoroutine()
	if current > baseline+2 {
		t.Fatalf("Goroutine leak detected. Baseline: %d, Current: %d", baseline, current)
	}

	if _, ok := <-actual; ok {
		t.Fatal("Or strategy channel should be closed after cancellation")
	}
}

func TestOrStrategyReport(t *testing.T) {
	snapshots, err := helper.ReadFromCsvFile[asset.Snapshot]("testdata/repository/brk-b.csv")
	if err != nil {
		t.Fatal(err)
	}

	or := strategy.NewOrStrategy("Or Strategy")
	or.Strategies = append(or.Strategies, strategy.NewBuyAndHoldStrategy(), trend.NewMacdStrategy())

	report := or.Report(snapshots)

	fileName := "or.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}
