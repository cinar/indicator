package helper

import (
	"testing"
)

func TestMaxSince(t *testing.T) {
	input := SliceToChan([]int{48, 49, 47, 52, 52, 52, 53, 50, 55})
	expected := SliceToChan([]int{0, 0, 1, 0, 1, 2, 0, 1, 0})
	actual := MaxSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMaxSinceAtEnd(t *testing.T) {
	input := SliceToChan([]int{48, 49, 47, 52, 52, 52, 52})
	expected := SliceToChan([]int{0, 0, 1, 0, 1, 2, 2})
	actual := MaxSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
