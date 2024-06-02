// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"reflect"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

// Testing against the example at URL:
// https://school.stockcharts.com/doku.php?id=technical_indicators:moving_averages

func TestSma(t *testing.T) {
	input := helper.SliceToChan([]float64{
		22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24,
		22.29, 22.15, 22.39, 22.38, 22.61, 23.36, 24.05, 23.75, 23.83,
		23.95, 23.63, 23.82, 23.87, 23.65, 23.19, 23.10, 23.33, 22.68,
		23.10, 22.40, 22.17,
	})

	expected := []float64{
		22.22, 22.21, 22.23, 22.26, 22.30, 22.42, 22.61, 22.77, 22.91,
		23.08, 23.21, 23.38, 23.53, 23.65, 23.71, 23.68, 23.61, 23.50,
		23.43, 23.28, 23.13,
	}

	sma := trend.NewSmaWithPeriod[float64](10)

	actual := helper.ChanToSlice(helper.RoundDigits(sma.Compute(input), 2))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSmaString(t *testing.T) {
	expected := "SMA(10)"
	actual := trend.NewSmaWithPeriod[float64](10).String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
