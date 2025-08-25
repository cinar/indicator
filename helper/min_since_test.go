package helper

import "testing"

func TestMinSince(t *testing.T) {
	input := SliceToChan([]int{48, 50, 50, 50, 49, 49, 51})
	expected := SliceToChan([]int{0, 1, 2, 2, 0, 1, 2})
	actual := MinSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMinSinceAtEnd(t *testing.T) {
	input := SliceToChan([]int{48, 50, 50, 50, 49, 49, 49})
	expected := SliceToChan([]int{0, 1, 2, 2, 0, 1, 2})
	actual := MinSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMinSinceFromStart(t *testing.T) {
	input := SliceToChan([]int{1, 1, 3})
	expected := SliceToChan([]int{0, 1, 2})
	actual := MinSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
