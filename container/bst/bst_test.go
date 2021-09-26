package bst

import (
	"testing"
)

func TestInsertAndRemove(t *testing.T) {
	values := []float64{2, 1, 3, 4, 0, 6, 6, 10, -1, 9}

	tree := New()

	for _, value := range values {
		tree.Insert(value)
	}

	for _, value := range values {
		if !tree.Remove(value) {
			t.Fatalf("unable to remove %f", value)
		}
	}
}

func TestMinAndMax(t *testing.T) {
	values := []float64{2, 1, 3, 4, 0, 6, 6, 10, -1, 9}
	mins := []float64{2, 1, 1, 1, 0, 0, 0, 0, -1, -1}
	maxs := []float64{2, 2, 3, 4, 4, 6, 6, 10, 10, 10}

	tree := New()

	for i := 0; i < len(values); i++ {
		tree.Insert(values[i])

		min := tree.Min()
		if min != mins[i] {
			t.Fatalf("at %d actual %f expected %f", i, min, mins[i])
		}

		max := tree.Max()
		if max != maxs[i] {
			t.Fatalf("at %d actual %f expected %f", i, max, maxs[i])
		}
	}

	for i := len(values) - 1; i > 0; i-- {
		tree.Remove(values[i])

		min := tree.Min()
		if min != mins[i-1] {
			t.Fatalf("at %d actual %f expected %f", i, min, mins[i-1])
		}

		max := tree.Max()
		if max != maxs[i-1] {
			t.Fatalf("at %d actual %f expected %f", i, max, maxs[i-1])
		}
	}
}
