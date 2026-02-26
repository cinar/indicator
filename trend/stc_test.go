// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestStcSlowStochastic(t *testing.T) {
	resultChan := make(chan float64, 10)
	go func() {
		values := []float64{
			100, 102, 101, 103, 105, 104, 106, 108, 107, 109,
			110, 111, 109, 112, 113, 115, 114, 116, 118, 117,
			119, 120, 118, 121, 122, 120, 123, 125, 124, 126,
		}

		input := helper.SliceToChan(values)
		s := trend.NewSlowStochastic[float64]()
		actualK, actualD := s.Compute(input)

		for {
			select {
			case v, ok := <-actualK:
				if !ok {
					return
				}
				resultChan <- v
			case v, ok := <-actualD:
				if !ok {
					return
				}
				resultChan <- v
			case <-time.After(2 * time.Second):
				return
			}
		}
	}()

	select {
	case v := <-resultChan:
		t.Logf("got value: %v", v)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout - no values produced")
	}
}

func TestStcFull(t *testing.T) {
	type Data struct {
		Value float64
		K     float64
		D     float64
	}

	input, err := helper.ReadFromCsvFile[Data]("testdata/stochastic.csv")
	if err != nil {
		t.Fatal(err)
	}

	values := helper.Map(input, func(d *Data) float64 { return d.Value })

	stc := trend.NewStcWithPeriod[float64](5, 10, 5, 3)
	result := stc.Compute(values)

	slice := helper.ChanToSlice(result)

	t.Logf("STC generated %d values", len(slice))

	if len(slice) == 0 {
		t.Fatal("no results generated")
	}
}
