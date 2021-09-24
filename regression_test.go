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

func TestLinearRegressionUsingLeastSquare(t *testing.T) {
	x := []float64{0, 2, 5, 7}
	y := []float64{-1, 5, 12, 20}

	expected := []float64{-1.1379, 4.6552, 13.3448, 19.1379}

	actual := LinearRegressionUsingLeastSquare(x, y)

	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		a := roundDigits(actual[i], 4)

		if a != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, a, expected[i])
		}
	}
}
