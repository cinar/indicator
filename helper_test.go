package indicator

import (
	"testing"
)

func TestMultiply(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	multiplier := float64(2)

	result := multiply(values, multiplier)
	if len(result) != len(values) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		expected := values[i] * multiplier
		actual := result[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestDivide(t *testing.T) {
	values := []float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	divider := float64(2)

	result := divide(values, divider)
	if len(result) != len(values) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		expected := values[i] / divider
		actual := result[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestAddWithDifferentSizes(t *testing.T) {
	values1 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	values2 := []float64{1, 2, 3, 4, 5}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("did not check size")
		}
	}()

	add(values1, values2)
}

func TestAdd(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := add(values, values)
	if len(result) != len(values) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		expected := values[i] + values[i]
		actual := result[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestSubstract(t *testing.T) {
	values1 := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	values2 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := substract(values1, values2)
	if len(result) != len(values1) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		expected := values1[i] - values2[i]
		actual := result[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}
