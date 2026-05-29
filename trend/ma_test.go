// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"context"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

type maWithoutContext[T helper.Number] struct{}

func (m *maWithoutContext[T]) Compute(c <-chan T) <-chan T {
	return c
}

func (m *maWithoutContext[T]) IdlePeriod() int {
	return 0
}

func (m *maWithoutContext[T]) String() string {
	return "maWithoutContext"
}

func TestComputeMaWithContextFallback(t *testing.T) {
	input := helper.SliceToChan([]int{1, 2, 3})
	ma := &maWithoutContext[int]{}

	result := trend.ComputeMaWithContext(context.Background(), ma, input)

	actual := helper.ChanToSlice(result)
	expected := []int{1, 2, 3}

	if len(actual) != len(expected) {
		t.Fatalf("got %v, expected %v", actual, expected)
	}
	for i, v := range actual {
		if v != expected[i] {
			t.Fatalf("got %v, expected %v", actual, expected)
		}
	}
}
