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

func TestCancellationWithApplyWindow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int)

	applied := helper.ApplyWithContext(ctx, input, func(v int) int { return v })
	windowed := helper.WindowWithContext(ctx, input, func(v []int, i int) int { return v[i] }, 3)

	cancel()

	time.Sleep(50 * time.Millisecond)

	_, ok := <-applied
	if ok {
		t.Fatal("Applied channel should be closed after cancellation")
	}
	_, ok = <-windowed
	if ok {
		t.Fatal("Windowed channel should be closed after cancellation")
	}
}

func TestCancellationWithEchoSkip(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int)

	echoed := helper.EchoWithContext(ctx, input, 3, 1)
	skipped := helper.SkipWithContext(ctx, input, 2)
	head := helper.HeadWithContext(ctx, input, 2)
	first := helper.FirstWithContext(ctx, input, 2)
	last := helper.LastWithContext(ctx, input, 2)

	cancel()

	time.Sleep(50 * time.Millisecond)

	if _, ok := <-echoed; ok {
		t.Fatal("Echo channel should be closed after cancellation")
	}
	if _, ok := <-skipped; ok {
		t.Fatal("Skipped channel should be closed after cancellation")
	}
	if _, ok := <-head; ok {
		t.Fatal("Head channel should be closed after cancellation")
	}
	if _, ok := <-first; ok {
		t.Fatal("First channel should be closed after cancellation")
	}
	if _, ok := <-last; ok {
		t.Fatal("Last channel should be closed after cancellation")
	}
}
