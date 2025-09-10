package momentum

import (
	"math"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestNewPringsSpecialKInitialization(t *testing.T) {
	psk := NewPringsSpecialK[float64]()

	if psk.Roc10 == nil || psk.Roc15 == nil || psk.Roc50 == nil || psk.Roc65 == nil ||
		psk.Roc75 == nil || psk.Roc100 == nil || psk.Roc130 == nil || psk.Roc195 == nil {
		t.Error("ROC pointers should be initialized")
	}
	if psk.Sma10 == nil || psk.Sma15 == nil || psk.Sma20 == nil || psk.Sma30 == nil ||
		psk.Sma40 == nil || psk.Sma65 == nil || psk.Sma75 == nil || psk.Sma100 == nil ||
		psk.Sma195 == nil || psk.Sma265 == nil || psk.Sma390 == nil || psk.Sma530 == nil {
		t.Error("SMA pointers should be initialized")
	}
}

func TestPringsSpecialKIsDecreasing(t *testing.T) {
	f := 0.1
	r := make([]float64, 0)
	for i := 0; i < 800; i++ {
		r = append(r, f)
		f += 0.01
	}
	cl := helper.SliceToChan(r)
	sp := NewPringsSpecialK[float64]()
	rs := sp.Compute(cl)
	p := math.NaN()
	upwardCount := 0
	for v := range rs {
		if !math.IsNaN(p) && p <= v {
			upwardCount++
		}
		p = v
	}
	if upwardCount != 0 {
		t.Error("Prings Special K should be decreasing on increasing values")
	}
}

func pringsSpecialKComputeOutputLengthOnInputLength(inputLen int) int {
	input := make(chan float64, inputLen)

	go func() {
		for i := 1; i <= inputLen; i++ {
			input <- float64(i)
		}
		close(input)
	}()

	psk := NewPringsSpecialK[float64]()
	out := psk.Compute(input)

	var got []float64
	for v := range out {
		got = append(got, v)
	}

	return len(got)
}

// Test sufficient number of samples for output
func TestPringsSpecialKComputeBasicOutput(t *testing.T) {
	sma530 := trend.NewSmaWithPeriod[float64](530)
	roc195 := trend.NewRocWithPeriod[float64](195)
	minimumRequiredInputLength := sma530.IdlePeriod() + roc195.IdlePeriod() + 1

	expectedOutputLen := 0
	actualOutputLen := pringsSpecialKComputeOutputLengthOnInputLength(minimumRequiredInputLength - 1)
	if actualOutputLen != expectedOutputLen {
		t.Errorf("Expected %d output values, got %d", expectedOutputLen, actualOutputLen)
	}
	expectedOutputLen = 1
	actualOutputLen = pringsSpecialKComputeOutputLengthOnInputLength(minimumRequiredInputLength)
	if actualOutputLen != expectedOutputLen {
		t.Errorf("Expected %d output values, got %d", expectedOutputLen, actualOutputLen)
	}
}

func TestPringsSpecialKComputeConstantInput(t *testing.T) {
	inputLen := 800
	val := 100.0
	input := make(chan float64, inputLen)
	go func() {
		for i := 0; i < inputLen; i++ {
			input <- val
		}
		close(input)
	}()
	psk := NewPringsSpecialK[float64]()
	out := psk.Compute(input)
	for v := range out {
		if v != 0.0 {
			t.Errorf("Expected output to be 0 for constant input, got %v", v)
			break
		}
	}
}
