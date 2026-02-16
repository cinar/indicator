// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package valuation_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/valuation"
)

func TestPv(t *testing.T) {
	type PvData struct {
		FV    float64
		Rate  float64
		Years int
		PV    float64
	}

	input, err := helper.ReadFromCsvFile[PvData]("testdata/pv.csv")
	if err != nil {
		t.Fatal(err)
	}

	for row := range input {
		pv := helper.RoundDigit(valuation.Pv(row.FV, row.Rate, row.Years), 2)
		if pv != row.PV {
			t.Fatalf("actual %v expected %v", pv, row.PV)
		}
	}
}
