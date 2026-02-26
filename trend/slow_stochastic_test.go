// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestSlowStochastic(t *testing.T) {
	type Data struct {
		Close float64 `header:"Close"`
		K     float64 `header:"K"`
		D     float64 `header:"D"`
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/slow_stochastic.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	closing := helper.Map(inputs[0], func(d *Data) float64 { return d.Close })

	s := trend.NewSlowStochastic[float64]()
	actualK, actualD := s.Compute(closing)

	actualK = helper.RoundDigits(actualK, 2)
	actualD = helper.RoundDigits(actualD, 2)

	inputs[1] = helper.Skip(inputs[1], s.IdlePeriod())

	for data := range inputs[1] {
		ak := <-actualK
		ad := <-actualD

		if ak != data.K {
			t.Fatalf("K: actual %v expected %v", ak, data.K)
		}

		if ad != data.D {
			t.Fatalf("D: actual %v expected %v", ad, data.D)
		}
	}
}

func TestNewSlowStochasticWithPeriod(t *testing.T) {
	s := trend.NewSlowStochasticWithPeriod[float64](14, 3, 3)
	if s.Period != 14 {
		t.Fatalf("expected period 14, got %d", s.Period)
	}
}
