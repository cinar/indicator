package helper

import (
	"math"
	"testing"
)

func TestMinSince(t *testing.T) {
	input := SliceToChan([]int{48, 50, 50, 50, 49, 49, 51})
	expected := SliceToChan([]int{0, 1, 2, 0, 0, 0, 1})
	actual := MinSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMinSinceAtEnd(t *testing.T) {
	input := SliceToChan([]int{48, 50, 50, 50, 49, 49, 49})
	expected := SliceToChan([]int{0, 1, 2, 0, 0, 0, 0})
	actual := MinSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMinSinceFromStart(t *testing.T) {
	input := SliceToChan([]int{1, 1, 3})
	expected := SliceToChan([]int{0, 0, 1})
	actual := MinSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMinSinceNaN(t *testing.T) {
	nan := math.NaN()
	input := SliceToChan([]float64{3, nan, 2, 4})
	actual := ChanToSlice(MinSince(input, 3))
	expected := []float64{0, 1, 0, 1}

	err := CheckEquals(SliceToChan(actual), SliceToChan(expected))
	if err != nil {
		t.Fatal(err)
	}

}
