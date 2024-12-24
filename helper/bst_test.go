// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

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
	minimums := []int{2, 1, 1, 1, 0, 0, 0, 0, -1, -1}
	maximums := []int{2, 2, 3, 4, 4, 6, 6, 10, 10, 10}

	bst := helper.NewBst[int]()

	for i, n := range input {
		bst.Insert(n)

		minimum := bst.Min()
		if minimum != minimums[i] {
			t.Fatalf("actual minimum %v expeceted minimum %v", minimum, minimums)
		}

		maximum := bst.Max()
		if maximum != maximums[i] {
			t.Fatalf("actual maximum %v expeceted maximum %v", maximum, maximums)
		}
	}
}

func TestEmptyMin(t *testing.T) {
	bst := helper.NewBst[int]()

	minValue := bst.Min()
	if minValue != 0 {
		t.Fatalf("actual min %v expected min 0", minValue)
	}
}

func TestEmptyMax(t *testing.T) {
	bst := helper.NewBst[int]()

	maxValue := bst.Max()
	if maxValue != 0 {
		t.Fatalf("actual max %v expected max 0", maxValue)
	}
}
