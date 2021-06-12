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

func TestDiffWithLargerBefore(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("did not check size")
		}
	}()

	diff([]float64{}, 10)
}

func TestDiff(t *testing.T) {
	values := []float64{1, 2, 1, 4, 2, 2, 6, 8, 2, 10}
	expected := []float64{0, 1, -1, 3, -2, 0, 4, 2, -6, 8}
	before := 1

	actual := diff(values, before)

	if len(actual) != len(expected) {
		t.Fatalf("actual %d expected %d", len(actual), len(expected))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestGroupPositivesAndNegatives(t *testing.T) {
	values := []float64{1, 0, -2, -4, 6, 0, 0, 4, 2, -20}
	expectedPositives := []float64{1, 6, 4, 2}
	expectedNegatives := []float64{-2, -4, -20}

	actualPositives, actualNegatives := groupPositivesAndNegatives(values)

	if len(actualPositives) != len(expectedPositives) {
		t.Fatalf("actual positives %d expected positives %d", len(actualPositives), len(expectedPositives))
	}

	for i := 0; i < len(actualPositives); i++ {
		if actualPositives[i] != expectedPositives[i] {
			t.Fatalf("at %d actual positive %f expected positive %f", i, actualPositives[i], expectedPositives[i])
		}
	}

	if len(actualNegatives) != len(expectedNegatives) {
		t.Fatalf("actual positives %d expected positives %d", len(actualNegatives), len(expectedNegatives))
	}

	for i := 0; i < len(actualNegatives); i++ {
		if actualNegatives[i] != expectedNegatives[i] {
			t.Fatalf("at %d actual negative %f expected negative %f", i, actualNegatives[i], expectedNegatives[i])
		}
	}
}
