// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestSlopeSimple(t *testing.T) {
	closing := helper.SliceToChan([]float64{10, 13, 17, 16, 20, 29})
	expected := helper.SliceToChan([]float64{(16 - 10) / 3.0, (20 - 13) / 3.0, (29 - 17) / 3.0})

	slope := NewSlopeWithPeriod[float64](3)
	actual := slope.Compute(closing)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSlopeFallbackPeriod(t *testing.T) {
	slope := NewSlopeWithPeriod[float64](-17)

	if slope.Period != DefaultSlopePeriod {
		t.Fatal("expected period to be fallback to default value")
	}
}

func TestSlopeToStringAndIdlePeriod(t *testing.T) {
	slope := NewSlopeWithPeriod[float64](0)
	if slope.IdlePeriod() != DefaultSlopePeriod {
		t.Fatalf("unexpected IdlePeriod: %d", slope.IdlePeriod())
	}

	slope.Period = 3
	if s := slope.String(); s != "SLOPE(3)" {
		t.Fatalf("unexpected String(): %s", s)
	}
}
