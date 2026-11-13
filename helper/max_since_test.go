package helper

import (
	"context"
	"math"
	"runtime"
	"testing"
	"time"
)

func TestMaxSince(t *testing.T) {
	input := SliceToChan([]int{48, 49, 47, 52, 52, 52, 53, 50, 55})
	expected := SliceToChan([]int{0, 0, 1, 0, 0, 0, 0, 1, 0})
	actual := MaxSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMaxSinceAtEnd(t *testing.T) {
	input := SliceToChan([]int{48, 49, 47, 52, 52, 52, 52})
	expected := SliceToChan([]int{0, 0, 1, 0, 0, 0, 0})
	actual := MaxSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMaxSinceLeavesWindow(t *testing.T) {
	// The 9 leaves the window at the fourth value, and the window is scanned
	// again: of the two 5s, the most recent one counts.
	input := SliceToChan([]int{9, 5, 5, 1, 1, 1})
	expected := SliceToChan([]int{0, 1, 2, 1, 2, 0})
	actual := MaxSince(input, 3)

	err := CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMaxSinceNaN(t *testing.T) {
	nan := math.NaN()
	input := SliceToChan([]float64{1, nan, 2, nan, nan, nan, 3})
	actual := ChanToSlice(MaxSince(input, 3))
	expected := []float64{0, 1, 0, 1, 2, nan, 0}

	for i := range expected {
		if math.IsNaN(expected[i]) != math.IsNaN(actual[i]) ||
			(!math.IsNaN(expected[i]) && actual[i] != expected[i]) {
			t.Fatalf("index %d actual %v expected %v", i, actual, expected)
		}
	}
}

func TestMaxSinceInvalidWindow(t *testing.T) {
	actual := ChanToSlice(MaxSince(SliceToChan([]int{1, 2, 3}), 0))
	if len(actual) != 0 {
		t.Fatalf("actual %v expected none", actual)
	}
}

func TestMaxSinceCancelWhileSending(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int, 1)
	input <- 1

	actual := MaxSinceWithContext(ctx, input, 3)

	// Wait until the value is taken, so that it is pending on the output.
	for len(input) > 0 {
		runtime.Gosched()
	}
	cancel()
	time.Sleep(50 * time.Millisecond)

	// Nothing read the output, so the cancellation was seen while sending.
	// The pending value may still be received before the channel closes.
	for range actual {
		continue
	}
}
