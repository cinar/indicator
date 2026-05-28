// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

func TestCancellationLeaks(t *testing.T) {
	runtime.GC()
	baseline := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int)

	dup := helper.DuplicateWithContext(ctx, input, 2)
	mapped0 := helper.MapWithContext(ctx, dup[0], func(x int) int { return x * 2 })
	mapped1 := helper.MapWithContext(ctx, dup[1], func(x int) int { return x + 1 })
	filtered := helper.FilterWithContext(ctx, mapped0, func(x int) bool { return x > 2 })
	count := helper.CountWithContext(ctx, 0, filtered)

	cancel()

	// Wait for goroutines to tear down
	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	current := runtime.NumGoroutine()
	if current > baseline+2 {
		t.Fatalf("Goroutine leak detected. Baseline: %d, Current: %d", baseline, current)
	}

	_, ok := <-count
	if ok {
		t.Fatal("Count channel should be closed after cancellation")
	}
	_, ok = <-mapped1
	if ok {
		t.Fatal("Mapped1 channel should be closed after cancellation")
	}
}
