// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import "testing"

func TestMoneyFlowIndex(t *testing.T) {
	high := []float64{10, 9, 12, 14, 12}
	low := []float64{6, 7, 9, 12, 10}
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{100, 100, 406.85, 207.69, 266.67}
	period := 2

	result := MoneyFlowIndex(period, high, low, closing, volume)
	if len(result) != len(expected) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		actual := roundDigits(result[i], 2)

		if actual != expected[i] {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected[i])
		}
	}
}

func TestForceIndex(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{900, 220, -320, 360, -180}
	period := 1

	result := ForceIndex(period, closing, volume)
	if len(result) != len(expected) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		actual := roundDigits(result[i], 2)

		if actual != expected[i] {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected[i])
		}
	}
}

func TestDefaultForceIndex(t *testing.T) {
	closing := []float64{9, 11, 7, 10, 8}
	volume := []int64{100, 110, 80, 120, 90}
	expected := []float64{900, 802.86, 642.45, 602.1, 490.37}

	result := DefaultForceIndex(closing, volume)
	if len(result) != len(expected) {
		t.Fatal("result not same size")
	}

	for i := 0; i < len(result); i++ {
		actual := roundDigits(result[i], 2)

		if actual != expected[i] {
			t.Fatalf("result %d actual %f expected %f", i, actual, expected[i])
		}
	}
}
