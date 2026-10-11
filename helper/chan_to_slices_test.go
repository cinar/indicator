// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

func TestChanToSlicesTwoChannels(t *testing.T) {
	input := []int{1, 2, 3, 4}
	c := helper.SliceToChan(input)

	outputs := helper.DuplicateWithContext(context.Background(), c, 2)

	actual := helper.ChanToSlices(outputs...)

	expected := [][]int{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestChanToSlicesMultipleChannels(t *testing.T) {
	input := []int{5, 10, 15}
	c := helper.SliceToChan(input)

	const count = 4

	outputs := helper.DuplicateWithContext(context.Background(), c, count)

	actual := helper.ChanToSlices(outputs...)

	if len(actual) != count {
		t.Fatalf("actual length %d expected %d", len(actual), count)
	}

	for i, slice := range actual {
		if !reflect.DeepEqual(slice, input) {
			t.Fatalf("index %d actual %v expected %v", i, slice, input)
		}
	}
}

func TestChanToSlicesDoesNotDeadlockOnUnbufferedChannels(t *testing.T) {
	runtime.GC()
	baseline := runtime.NumGoroutine()

	first := make(chan int)
	second := make(chan int)
	third := make(chan int)

	go func() {
		defer close(first)
		defer close(second)
		defer close(third)

		for i := 0; i < 5; i++ {
			first <- i
			second <- i
			third <- i
		}
	}()

	done := make(chan [][]int)

	go func() {
		done <- helper.ChanToSlices(first, second, third)
	}()

	select {
	case actual := <-done:
		expected := []int{0, 1, 2, 3, 4}

		for i, slice := range actual {
			if !reflect.DeepEqual(slice, expected) {
				t.Fatalf("index %d actual %v expected %v", i, slice, expected)
			}
		}

	case <-time.After(time.Second):
		t.Fatal("ChanToSlices deadlocked draining unbuffered channels")
	}

	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	current := runtime.NumGoroutine()
	if current > baseline+2 {
		t.Fatalf("Goroutine leak detected. Baseline: %d, Current: %d", baseline, current)
	}
}
