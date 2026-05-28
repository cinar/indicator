// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package valuation

import "math"

// Fv calculates the Future Value (FV) of a Present Value (PV).
//
//  Formula: FV = PV * (1 + rate)^years
func Fv(pv, rate float64, years int) float64 {
	return pv * math.Pow((1+rate), float64(years))
}
