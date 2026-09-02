// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestRing(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{0, 0, 0, 0, 1, 2, 3, 4, 5, 6}

	ring := helper.NewRing[int](4)

	for i, n := range input {
		actual := ring.Put(n)
		if actual != expected[i] {
			t.Fatalf("actual %v expected %v", actual, expected[i])
		}
	}
}

func TestRingEmpty(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	size := 4

	ring := helper.NewRing[int](size)

	if !ring.IsEmpty() {
		t.Fatal("not empty")
	}

	for i, n := range input {
		ring.Put(n)

		if ring.IsEmpty() {
			t.Fatal("is empty")
		}

		j := i
		if j >= size {
			j = size - 1
		}

		actual, ok := ring.At(j)
		if !ok {
			t.Fatalf("At(%d) not ok", j)
		}
		if actual != n {
			t.Fatalf("actual %v expected %v", actual, n)
		}
	}

	if !ring.IsFull() {
		t.Fatal("not full")
	}

	for i := 0; i < size; i++ {
		_, ok := ring.Get()
		if !ok {
			t.Fatal("is empty")
		}
	}

	if !ring.IsEmpty() {
		t.Fatal("not empty")
	}

	_, ok := ring.Get()
	if ok {
		t.Fatal("not empty")
	}
}

func TestRingAtPartiallyFilled(t *testing.T) {
	ring := helper.NewRing[int](5)

	ring.Put(1)
	ring.Put(2)

	for i, expected := range []int{1, 2} {
		actual, ok := ring.At(i)
		if !ok {
			t.Fatalf("At(%d) not ok", i)
		}
		if actual != expected {
			t.Fatalf("actual %v expected %v", actual, expected)
		}
	}

	for i := 2; i < 5; i++ {
		if _, ok := ring.At(i); ok {
			t.Fatalf("At(%d) should not be ok", i)
		}
	}
}

func TestRingZeroSizeClampedToOne(t *testing.T) {
	ring := helper.NewRing[int](0)

	if ring.Put(1) != 0 {
		t.Fatal("expected zero value evicted on first put")
	}

	actual, ok := ring.At(0)
	if !ok {
		t.Fatal("At(0) not ok")
	}
	if actual != 1 {
		t.Fatalf("actual %v expected %v", actual, 1)
	}

	if !ring.IsFull() {
		t.Fatal("expected ring of clamped size 1 to be full after one put")
	}
}

func TestRingNegativeSizeClampedToOne(t *testing.T) {
	ring := helper.NewRing[int](-1)

	ring.Put(1)
	actual := ring.Put(2)

	if actual != 1 {
		t.Fatalf("actual %v expected %v", actual, 1)
	}
}

func TestRingAtFull(t *testing.T) {
	ring := helper.NewRing[int](5)

	for i := 1; i <= 5; i++ {
		ring.Put(i)
	}

	for i, expected := range []int{1, 2, 3, 4, 5} {
		actual, ok := ring.At(i)
		if !ok {
			t.Fatalf("At(%d) not ok", i)
		}
		if actual != expected {
			t.Fatalf("actual %v expected %v", actual, expected)
		}
	}
}
