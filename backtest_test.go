// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator"
)

func TestNormalizeActions(t *testing.T) {
	actions := []indicator.Action{
		indicator.HOLD,
		indicator.SELL,
		indicator.SELL,
		indicator.BUY,
		indicator.BUY,
		indicator.SELL,
		indicator.BUY,
		indicator.BUY,
		indicator.HOLD,
		indicator.SELL,
	}

	expected := []indicator.Action{
		indicator.HOLD,
		indicator.HOLD,
		indicator.HOLD,
		indicator.BUY,
		indicator.HOLD,
		indicator.SELL,
		indicator.BUY,
		indicator.HOLD,
		indicator.HOLD,
		indicator.SELL,
	}

	normalized := indicator.NormalizeActions(actions)

	if len(normalized) != len(expected) {
		t.Fatal("not the same size")
	}

	for i, actual := range normalized {
		if actual != expected[i] {
			t.Fatalf("at %d actual %v expected %v", i, actual, expected[i])
		}
	}
}

func TestCountTransactions(t *testing.T) {
	actions := []indicator.Action{
		indicator.HOLD,
		indicator.HOLD,
		indicator.HOLD,
		indicator.BUY,
		indicator.HOLD,
		indicator.SELL,
		indicator.BUY,
		indicator.HOLD,
		indicator.HOLD,
		indicator.SELL,
	}

	expected := 4

	actual := indicator.CountTransactions(actions)

	if actual != expected {
		t.Fatalf("actual %d expected %d", actual, expected)
	}
}

func TestApplyActions(t *testing.T) {
	prices := []float64{
		1.00,
		2.00,
		3.00,
		4.00,
		4.00,
		5.00,
		7.00,
		5.00,
		8.00,
		9.00,
	}

	actions := []indicator.Action{
		indicator.HOLD,
		indicator.HOLD,
		indicator.HOLD,
		indicator.BUY,
		indicator.HOLD,
		indicator.SELL,
		indicator.BUY,
		indicator.HOLD,
		indicator.HOLD,
		indicator.SELL,
	}

	expected := []float64{
		0.00,
		0.00,
		0.00,
		0.00,
		0.00,
		0.25,
		0.25,
		-0.11,
		0.43,
		0.61,
	}

	gains := indicator.ApplyActions(prices, actions)

	if len(gains) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(gains); i++ {
		actual := math.Round(gains[i]*100) / 100

		if actual != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual, expected[i])
		}
	}
}

func TestNormalizeGains(t *testing.T) {
	prices := []float64{2, 4, 6, 12, 18}
	gains := []float64{0, 1, 1.5, 2.5, 3}
	expected := []float64{0, 0, 0, 0, 0}

	actual := indicator.NormalizeGains(prices, gains)
	if len(actual) != len(expected) {
		t.Fatal("not the same size")
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("at %d actual %f expected %f", i, actual[i], expected[i])
		}
	}
}
