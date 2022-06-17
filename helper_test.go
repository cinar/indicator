// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"testing"
)

func TestMultiplyBy(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	multiplier := float64(2)

	result := multiplyBy(values, multiplier)
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

func TestMultiply(t *testing.T) {
	values1 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	values2 := []float64{2, 4, 2, 4, 2, 4, 2, 4, 2, 4}
	expected := []float64{2, 8, 6, 16, 10, 24, 14, 32, 18, 40}

	actual := multiply(values1, values2)

	if len(actual) != len(expected) {
		t.Fatalf("actual %d expected %d", len(actual), len(expected))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestDivideBy(t *testing.T) {
	values := []float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	divider := float64(2)

	result := divideBy(values, divider)
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

func TestDivide(t *testing.T) {
	values1 := []float64{2, 8, 6, 16, 10, 24, 14, 32, 18, 40}
	values2 := []float64{2, 4, 2, 4, 2, 4, 2, 4, 2, 4}
	expected := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	actual := divide(values1, values2)
	checkSameSize(actual, expected)

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
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
	checkSameSize(result, values)

	for i := 0; i < len(result); i++ {
		expected := values[i] + values[i]
		actual := result[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestAddBy(t *testing.T) {
	values := []float64{1, 2, 3, 4}
	expected := []float64{2, 3, 4, 5}

	actual := addBy(values, 1)
	checkSameSize(actual, expected)

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("result %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestSubtract(t *testing.T) {
	values1 := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	values2 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := subtract(values1, values2)
	checkSameSize(result, values1)

	for i := 0; i < len(result); i++ {
		expected := values1[i] - values2[i]
		actual := result[i]

		if actual != expected {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected)
		}
	}
}

func TestDiff(t *testing.T) {
	values := []float64{1, 2, 1, 4, 2, 2, 6, 8, 2, 10}
	expected := []float64{1, 1, -1, 3, -2, 0, 4, 2, -6, 8}
	before := 1

	actual := diff(values, before)
	checkSameSize(actual, expected)

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestPercentDiff(t *testing.T) {
	values := []float64{1, 2, 8, 32, 64}
	expected := []float64{0, 1, 3, 3, 1}
	before := 1

	actual := percentDiff(values, before)
	checkSameSize(actual, expected)

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestShiftRightAndFillBy(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []float64{1, 1, 1, 1, 1, 2, 3, 4, 5, 6}
	period := 4

	actual := shiftRightAndFillBy(period, values[0], values)
	checkSameSize(actual, expected)

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestShiftRight(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []float64{0, 0, 0, 0, 1, 2, 3, 4, 5, 6}
	period := 4

	actual := shiftRight(period, values)
	checkSameSize(actual, expected)

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestRoundDigits(t *testing.T) {
	value := 1.5182345
	expected := 1.5182

	actual := roundDigits(value, 4)

	if actual != expected {
		t.Fatalf("actual %f expected %f", actual, expected)
	}
}

func TestGenerateNunbers(t *testing.T) {
	expected := []float64{2, 4, 6, 8}

	actual := generateNumbers(2, 10, 2)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestAsFloat64(t *testing.T) {
	values := []int64{1, 2, 3, 4}
	expected := []float64{1, 2, 3, 4}

	actual := asFloat64(values)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestPow(t *testing.T) {
	values := []float64{1, 2, 3, 4}
	expected := []float64{1, 4, 9, 16}
	exponent := 2.0

	actual := pow(values, exponent)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestExtractSign(t *testing.T) {
	values := []float64{1, -2, 3, 4, -5, 0, 6, -8, -9, 10}
	expected := []float64{1, -1, 1, 1, -1, 1, 1, -1, -1, 1}

	actual := extractSign(values)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestKeepPositives(t *testing.T) {
	values := []float64{1, -2, -3, 4}
	expected := []float64{1, 0, 0, 4}

	actual := keepPositives(values)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}

func TestKeepNegatives(t *testing.T) {
	values := []float64{1, -2, -3, 4}
	expected := []float64{0, -2, -3, 0}

	actual := keepNegatives(values)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}
