package helper

import (
	"testing"
)

func TestWindowSimple(t *testing.T) {
	input := SliceToChan([]int{1, 2, 3, 4, 5})
	expected := SliceToChan([]int{1, 2, 3, 4, 5})
	actual := Window(input, func(v []int, i int) int {
		if len(v) != 1 {
			panic("size != 1")
		}
		return v[0]
	}, 1)
	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWindowPairSum(t *testing.T) {
	input := SliceToChan([]int{1, 2, 3, 4, 5})
	expected := SliceToChan([]int{1, 3, 5, 7, 9})
	actual := Window(input, func(v []int, i int) int {
		if len(v) < 1 {
			panic("size < 1")
		}
		sum := 0
		for _, n := range v {
			sum += n
		}
		return sum
	}, 2)
	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWindowIndex(t *testing.T) {
	input := SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	expected := SliceToChan([]int{0, 0, 0, 1, 2, 0, 1, 2, 0, 1})
	actual := Window(input, func(v []int, i int) int {
		return i
	}, 3)
	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWindowIndexElement(t *testing.T) {
	input := SliceToChan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	expected := SliceToChan([]int{1, 1, 1, 2, 3, 4, 5, 6, 7, 8})
	actual := Window(input, func(v []int, i int) int {
		return v[i]
	}, 3)
	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
