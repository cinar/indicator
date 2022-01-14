// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"testing"
)

func TestLeastSquare(t *testing.T) {
	x := []float64{2, 3, 5, 7, 9}
	y := []float64{4, 5, 7, 10, 15}

	m, b := LeastSquare(x, y)

	expectedM := 1.5183
	expectedB := 0.3049

	m = roundDigits(m, 4)
	b = roundDigits(b, 4)

	if m != expectedM {
		t.Fatalf("m %f expected %f", m, expectedM)
	}

	if b != expectedB {
		t.Fatalf("b %f expected %f", b, expectedB)
	}
}

func TestMovingLeastSquare(t *testing.T) {
	x := []float64{1, 2, 3, 4, 2, 3, 5, 7, 9}
	y := []float64{1, 2, 3, 4, 4, 5, 7, 10, 15}

	period := 5

	m, b := MovingLeastSquare(period, x, y)

	expectedM := 1.5183
	expectedB := 0.3049

	finalM := roundDigits(m[len(m)-1], 4)
	finalB := roundDigits(b[len(b)-1], 4)

	if finalM != expectedM {
		t.Fatalf("m %f expected %f", finalM, expectedM)
	}

	if finalB != expectedB {
		t.Fatalf("m %f expected %f", finalB, expectedB)
	}
}

func TestLinearRegressionUsingLeastSquare(t *testing.T) {
	x := []float64{0, 2, 5, 7}
	y := []float64{-1, 5, 12, 20}

	expected := []float64{-1.1379, 4.6552, 13.3448, 19.1379}

	actual := LinearRegressionUsingLeastSquare(x, y)
	testEquals(t, roundDigitsAll(actual, 4), expected)
}
