package helper

import (
	"slices"
	"testing"
)

func TestSlicesReverseSimple(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{4, 3, 2, 1}
	actual := make([]int, 0, len(input))
	SlicesReverse(input, 0, func(i int) bool {
		actual = append(actual, i)
		return true
	})
	if !slices.Equal(actual, expected) {
		t.Fatal("not equal")
	}
}

func TestSlicesReverseMiddle(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 1, 4, 3}
	actual := make([]int, 0, len(input))
	SlicesReverse(input, 2, func(i int) bool {
		actual = append(actual, i)
		return true
	})
	if !slices.Equal(actual, expected) {
		t.Fatal("not equal")
	}
}
