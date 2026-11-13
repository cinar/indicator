// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"context"
	"math"
	"runtime"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestAroon(t *testing.T) {
	type AroonData struct {
		High float64
		Low  float64
		Up   float64
		Down float64
	}

	aroon := trend.NewAroon[float64]()

	input, err := helper.ReadFromCsvFile[AroonData]("testdata/aroon.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(row *AroonData) float64 { return row.High })
	low := helper.Map(inputs[1], func(row *AroonData) float64 { return row.Low })
	expectedUp := helper.Map(inputs[2], func(row *AroonData) float64 { return row.Up })
	expectedDown := helper.Map(inputs[3], func(row *AroonData) float64 { return row.Down })

	expectedUp = helper.Skip(expectedUp, aroon.IdlePeriod())
	expectedDown = helper.Skip(expectedDown, aroon.IdlePeriod())

	actualUp, actualDown := aroon.Compute(high, low)

	err = helper.CheckEquals(actualUp, expectedUp, actualDown, expectedDown)
	if err != nil {
		t.Fatal(err)
	}
}

// TestAroonCompareTaLib compares a short period with TA-Lib 0.8.3
// AROON(high, low, 5) on the same data.
func TestAroonCompareTaLib(t *testing.T) {
	type AroonData struct {
		High float64
		Low  float64
		Up   float64
		Down float64
	}
	var testData = []AroonData{
		{High: 82.15, Low: 81.29}, // None
		{High: 81.89, Low: 80.64}, // None
		{High: 83.03, Low: 81.31}, // None
		{High: 83.30, Low: 82.65}, // None
		{High: 83.85, Low: 83.07}, // None
		{High: 83.90, Low: 83.11, Up: 100, Down: 20},
		{High: 83.33, Low: 82.49, Up: 80, Down: 0},
		{High: 84.30, Low: 82.30, Up: 100, Down: 0},
		{High: 84.84, Low: 84.15, Up: 100, Down: 80},
		{High: 85.00, Low: 84.11, Up: 100, Down: 60},
		{High: 85.90, Low: 84.03, Up: 100, Down: 40},
		{High: 86.58, Low: 85.39, Up: 100, Down: 20},
		{High: 86.98, Low: 85.76, Up: 100, Down: 0},
		{High: 88.00, Low: 87.17, Up: 100, Down: 40},
		{High: 87.87, Low: 87.01, Up: 80, Down: 20},
	}

	input := helper.SliceToChan[AroonData](testData)
	aroon := &trend.Aroon[float64]{
		Period: 5,
	}

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(row AroonData) float64 { return row.High })
	low := helper.Map(inputs[1], func(row AroonData) float64 { return row.Low })
	expectedUp := helper.Map(inputs[2], func(row AroonData) float64 { return row.Up })
	expectedDown := helper.Map(inputs[3], func(row AroonData) float64 { return row.Down })

	expectedUp = helper.Skip(expectedUp, aroon.IdlePeriod())
	expectedDown = helper.Skip(expectedDown, aroon.IdlePeriod())

	actualUp, actualDown := aroon.Compute(high, low)

	err := helper.CheckEquals(actualUp, expectedUp, actualDown, expectedDown)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAroonCancellation(t *testing.T) {
	runtime.GC()
	baseline := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	high := make(chan float64)
	low := make(chan float64)

	aroon := trend.NewAroon[float64]()
	up, down := aroon.ComputeWithContext(ctx, high, low)

	cancel()

	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	current := runtime.NumGoroutine()
	if current > baseline+2 {
		t.Fatalf("Goroutine leak detected. Baseline: %d, Current: %d", baseline, current)
	}

	if _, ok := <-up; ok {
		t.Fatal("Up channel should be closed after cancellation")
	}
	if _, ok := <-down; ok {
		t.Fatal("Down channel should be closed after cancellation")
	}
}

/* https://github.com/francescosisini/TTR/blob/27ea28698295f56447fcc87ad515140bdb35c8a5/R/aroon.R#L27
'Aroon up (down) is the elapsed time, expressed as a percentage, between today
#'and the highest (lowest) price in the last \code{n} periods.  If today's
#'price is a new high (low) Aroon up (down) will be 100. Each subsequent period
#'without another new high (low) causes Aroon up (down) to decrease by (1 /
#'\code{n}) x 100.
*/

// TestAroonValidateDefinition checks the definition above on hand-made data
// with repeated highs and lows, where the most recent one counts, as in
// TA-Lib. At period 3 the values are not whole numbers.
func TestAroonValidateDefinition(t *testing.T) {
	type AroonData struct {
		High float64
		Low  float64
		Up   float64
		Down float64
	}
	var testData = []AroonData{
		{High: 10, Low: 5}, // None
		{High: 12, Low: 4}, // None
		{High: 11, Low: 6}, // None
		// The high of 12 and the low of 4 repeat; the most recent one counts.
		{High: 12, Low: 4, Up: 100, Down: 100},
		{High: 9, Low: 7, Up: 66.67, Down: 66.67},
		{High: 8, Low: 8, Up: 33.33, Down: 33.33},
		{High: 7, Low: 9, Up: 0, Down: 0},
	}

	input := helper.SliceToChan[AroonData](testData)
	aroon := &trend.Aroon[float64]{
		Period: 3,
	}

	inputs := helper.Duplicate(input, 4)
	high := helper.Map(inputs[0], func(row AroonData) float64 { return row.High })
	low := helper.Map(inputs[1], func(row AroonData) float64 { return row.Low })
	expectedUp := helper.Map(inputs[2], func(row AroonData) float64 { return row.Up })
	expectedDown := helper.Map(inputs[3], func(row AroonData) float64 { return row.Down })

	expectedUp = helper.Skip(expectedUp, aroon.IdlePeriod())
	expectedDown = helper.Skip(expectedDown, aroon.IdlePeriod())

	actualUp, actualDown := aroon.Compute(high, low)
	actualUp = helper.RoundDigits(actualUp, 2)
	actualDown = helper.RoundDigits(actualDown, 2)

	err := helper.CheckEquals(actualUp, expectedUp, actualDown, expectedDown)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAroonInvalidPeriod(t *testing.T) {
	for _, period := range []int{0, -1} {
		aroon := &trend.Aroon[float64]{
			Period: period,
		}

		up, down := aroon.Compute(helper.SliceToChan([]float64{1, 2, 3}), helper.SliceToChan([]float64{1, 2, 3}))

		if actual := helper.ChanToSlice(up); len(actual) != 0 {
			t.Fatalf("period %d up %v expected none", period, actual)
		}
		if actual := helper.ChanToSlice(down); len(actual) != 0 {
			t.Fatalf("period %d down %v expected none", period, actual)
		}
	}
}

func TestAroonNaN(t *testing.T) {
	// A missing high or low is skipped; the extremes around it still count.
	nan := math.NaN()
	high := helper.SliceToChan([]float64{1, 3, nan, 2, 4})
	low := helper.SliceToChan([]float64{4, 2, nan, 3, 1})
	expected := []float64{50, 0, 100}

	aroon := &trend.Aroon[float64]{
		Period: 2,
	}

	up, down := aroon.Compute(high, low)

	err := helper.CheckEquals(up, helper.SliceToChan(expected), down, helper.SliceToChan(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAroonString(t *testing.T) {
	expected := "AROON(25)"
	actual := trend.NewAroon[float64]().String()

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
