// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package valuation

import "math"

// Pv calculates the Present Value (PV) of a Future Value (FV).
//
//  Formula: PV = FV / (1 + rate)^years
func Pv(fv, rate float64, years int) float64 {
	return fv / math.Pow((1+rate), float64(years))
}
