// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

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

		if ring.At(j) != n {
			t.Fatalf("actual %v expected %v", ring.At(j), n)
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
