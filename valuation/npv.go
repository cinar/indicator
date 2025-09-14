// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package valuation

import "math"

// Npv calculates the Net Present Value (NPV) of a series of cash flows.
//
//  Formula: NPV = sum(CF_i / (1 + rate)^i) for i = 1 to n
func Npv(rate float64, cfs []float64) float64 {
	var npv float64
	for i, cf := range cfs {
		npv += cf / math.Pow(1+rate, float64(i+1))
	}

	return npv
}

