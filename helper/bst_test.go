// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestInsertAndContains(t *testing.T) {
	input := []int{2, 1, 3, 4, 0, 6, 6, 10, -1, 9}

	bst := helper.NewBst[int]()

	for _, n := range input {
		bst.Insert(n)
	}

	for _, n := range input {
		if !bst.Contains(n) {
			t.Fatalf("value %v not found", n)
		}
	}
}

func TestInsertAndRemove(t *testing.T) {
	input := []int{2, 1, 3, 4, 0, 6, 6, 10, -1, 9}

	bst := helper.NewBst[int]()

	for _, n := range input {
		bst.Insert(n)
	}

	for _, n := range input {
		if !bst.Remove(n) {
			t.Fatalf("value %v not found", n)
		}
	}
}

func TestRemoveNonExistentValue(t *testing.T) {
	input := []int{2, 1, 3, 4, 0, 6, 6, 10, -1, 9}

	bst := helper.NewBst[int]()

	for _, n := range input {
		bst.Insert(n)
	}

	if bst.Remove(8) {
		t.Fatal("non existent value removed")
	}
}

func TestMinAndMax(t *testing.T) {
	input := []int{2, 1, 3, 4, 0, 6, 6, 10, -1, 9}
	mins := []int{2, 1, 1, 1, 0, 0, 0, 0, -1, -1}
	maxs := []int{2, 2, 3, 4, 4, 6, 6, 10, 10, 10}

	bst := helper.NewBst[int]()

	for i, n := range input {
		bst.Insert(n)

		min := bst.Min()
		if min != mins[i] {
			t.Fatalf("actual min %v expeceted min %v", min, mins)
		}

		max := bst.Max()
		if max != maxs[i] {
			t.Fatalf("actual min %v expeceted min %v", max, maxs)
		}
	}
}

func TestEmptyMin(t *testing.T) {
	bst := helper.NewBst[int]()

	min := bst.Min()
	if min != 0 {
		t.Fatalf("actual min %v expected min 0", min)
	}
}

func TestEmptyMax(t *testing.T) {
	bst := helper.NewBst[int]()

	max := bst.Max()
	if max != 0 {
		t.Fatalf("actual max %v expected max 0", max)
	}
}
