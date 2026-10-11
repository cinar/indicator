// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"context"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

func TestSince(t *testing.T) {
	input := helper.SliceToChan([]int{1, 1, 2, 2, 2, 1, 2, 3, 3, 4})
	expected := helper.SliceToChan([]int{0, 1, 0, 1, 2, 0, 0, 0, 1, 0})

	actual := helper.Since[int, int](input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSinceWithContext(t *testing.T) {
	input := helper.SliceToChan([]int{1, 1, 2, 2, 2, 1, 2, 3, 3, 4})
	expected := helper.SliceToChan([]int{0, 1, 0, 1, 2, 0, 0, 0, 1, 0})

	actual := helper.SinceWithContext[int, int](context.Background(), input)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSinceWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int)

	since := helper.SinceWithContext[int, int](ctx, input)

	cancel()

	time.Sleep(50 * time.Millisecond)

	if _, ok := <-since; ok {
		t.Fatal("Since channel should be closed after cancellation")
	}
}
