// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package valuation_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/valuation"
)

func TestNpv(t *testing.T) {
	type NpvData struct {
		Rate      float64
		CashFlows string
		NPV       float64
	}

	parseCashFlows := func(s string) []float64 {
		var cashFlows []float64
		for _, cfStr := range strings.Split(s, " ") {
			cf, _ := strconv.ParseFloat(cfStr, 64)
			cashFlows = append(cashFlows, cf)
		}
		return cashFlows
	}

	input, err := helper.ReadFromCsvFile[NpvData]("testdata/npv.csv")
	if err != nil {
		t.Fatal(err)
	}

	for row := range input {
		cashFlows := parseCashFlows(row.CashFlows)
		npv := helper.RoundDigit(valuation.Npv(row.Rate, cashFlows), 2)
		if npv != row.NPV {
			t.Fatalf("actual %v expected %v", npv, row.NPV)
		}
	}
}
