// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package valuation_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/valuation"
)

func TestFv(t *testing.T) {
	type FvData struct {
		PV    float64
		Rate  float64
		Years int
		FV    float64
	}

	input, err := helper.ReadFromCsvFile[FvData]("testdata/fv.csv")
	if err != nil {
		t.Fatal(err)
	}

	for row := range input {
		fv := helper.RoundDigit(valuation.Fv(row.PV, row.Rate, row.Years), 2)
		if fv != row.FV {
			t.Fatalf("actual %v expected %v", fv, row.FV)
		}
	}
}
