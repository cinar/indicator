// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestKstSimple(t *testing.T) {
	// Use a simple linear price progression to verify correctness
	// With steadily increasing prices, all ROCs should be positive
	prices := make([]float64, 100)
	for i := range prices {
		prices[i] = 100 + float64(i) // 100, 101, 102, ...
	}

	input := helper.SliceToChan(prices)
	kst := trend.NewKst[float64]()
	actualKst, actualSignal := kst.Compute(input)

	// Verify we get some output
	kstSlice := helper.ChanToSlice(actualKst)
	signalSlice := helper.ChanToSlice(actualSignal)

	if len(kstSlice) == 0 {
		t.Fatal("KST produced no output")
	}

	// With steadily increasing prices, KST should be positive
	for i, v := range kstSlice {
		if v < 0 {
			t.Fatalf("KST at index %d is negative: %v", i, v)
		}
	}

	// Signal should also be positive
	for i, v := range signalSlice {
		if v < 0 {
			t.Fatalf("Signal at index %d is negative: %v", i, v)
		}
	}

	t.Logf("KST first value: %v", kstSlice[0])
	if len(signalSlice) > 0 {
		t.Logf("Signal first value: %v", signalSlice[0])
	}
}

func TestKstString(t *testing.T) {
	kst := trend.NewKst[float64]()
	expected := "KST(10,15,20,30,10,10,10,15,9)"
	actual := kst.String()
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func TestKstIdlePeriod(t *testing.T) {
	kst := trend.NewKst[float64]()
	expected := 52
	actual := kst.IdlePeriod()
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func TestKstWithNaNCheck(t *testing.T) {
	// Test with some random prices to ensure no NaN values
	prices := []float64{
		100, 101, 102, 103, 104, 105, 106, 107, 108, 109,
		110, 111, 112, 113, 114, 115, 116, 117, 118, 119,
		120, 121, 122, 123, 124, 125, 126, 127, 128, 129,
		130, 131, 132, 133, 134, 135, 136, 137, 138, 139,
		140, 141, 142, 143, 144, 145, 146, 147, 148, 149,
		150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169,
		170, 171, 172, 173, 174, 175, 176, 177, 178, 179,
		180, 181, 182, 183, 184, 185, 186, 187, 188, 189,
		190, 191, 192, 193, 194, 195, 196, 197, 198, 199,
	}

	input := helper.SliceToChan(prices)
	kst := trend.NewKst[float64]()
	actualKst, actualSignal := kst.Compute(input)

	kstSlice := helper.ChanToSlice(actualKst)
	signalSlice := helper.ChanToSlice(actualSignal)

	// Check for NaN values
	for i, v := range kstSlice {
		if math.IsNaN(v) {
			t.Fatalf("KST at index %d is NaN", i)
		}
	}

	for i, v := range signalSlice {
		if math.IsNaN(v) {
			t.Fatalf("Signal at index %d is NaN", i)
		}
	}

	t.Logf("KST: %d values, Signal: %d values", len(kstSlice), len(signalSlice))
}

func TestKst(t *testing.T) {
	type Data struct {
		Close  float64
		Kst    float64
		Signal float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/kst.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

	kst := trend.NewKst[float64]()
	actualKst, actualSignal := kst.Compute(closing)

	actualKst = helper.RoundDigits(actualKst, 2)
	actualSignal = helper.RoundDigits(actualSignal, 2)

	inputs[1] = helper.Skip(inputs[1], kst.IdlePeriod())

	for data := range inputs[1] {
		actualK := <-actualKst
		actualS := <-actualSignal

		if actualK != data.Kst {
			t.Fatalf("actual %v expected %v", actualK, data.Kst)
		}

		if actualS != data.Signal {
			t.Fatalf("actual %v expected %v", actualS, data.Signal)
		}
	}
}
