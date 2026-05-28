// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"fmt"

	"github.com/cinar/indicator/v2/helper"
)

// PivotPointMethod represents the method used for calculating pivot points.
type PivotPointMethod int

const (
	// PivotPointStandard is the standard pivot point calculation.
	PivotPointStandard PivotPointMethod = iota

	// PivotPointWoodie is the Woodie pivot point calculation.
	PivotPointWoodie

	// PivotPointCamarilla is the Camarilla pivot point calculation.
	PivotPointCamarilla

	// PivotPointFibonacci is the Fibonacci pivot point calculation.
	PivotPointFibonacci
)

// PivotPointResult represents the result of the pivot point calculation, including
// the pivot point itself, and its associated resistance (R) and support (S) levels.
type PivotPointResult[T helper.Float] struct {
	P  T
	R1 T
	R2 T
	R3 T
	R4 T
	S1 T
	S2 T
	S3 T
	S4 T
}

// PivotPoint represents the configuration parameters for calculating Pivot Points.
// Pivot points are calculated based on the previous period's high, low, and close,
// and are used to predict support and resistance levels for the current period.
type PivotPoint[T helper.Float] struct {
	// Method is the pivot point calculation method.
	Method PivotPointMethod
}

// NewPivotPoint function initializes a new Pivot Point instance with the standard method.
func NewPivotPoint[T helper.Float]() *PivotPoint[T] {
	return NewPivotPointWithMethod[T](PivotPointStandard)
}

// NewPivotPointWithMethod function initializes a new Pivot Point instance with the given method.
func NewPivotPointWithMethod[T helper.Float](method PivotPointMethod) *PivotPoint[T] {
	return &PivotPoint[T]{
		Method: method,
	}
}

// Compute function takes channels for open, high, low, and closing prices and
// returns a channel of PivotPointResult. It uses the values from the previous
// period to calculate levels for the current period.
func (p *PivotPoint[T]) Compute(opens, highs, lows, closings <-chan T) <-chan PivotPointResult[T] {
	result := make(chan PivotPointResult[T], cap(closings))

	go func() {
		defer close(result)

		var prevH, prevL, prevC T
		first := true

		for {
			o, okO := <-opens
			h, okH := <-highs
			l, okL := <-lows
			c, okC := <-closings

			if !okO || !okH || !okL || !okC {
				break
			}

			if !first {
				result <- p.calculate(prevH, prevL, prevC, o)
			}

			prevH, prevL, prevC = h, l, c
			first = false
		}
	}()

	return result
}

// calculate calculates the pivot points using the specified method.
func (p *PivotPoint[T]) calculate(h, l, c, currO T) PivotPointResult[T] {
	var res PivotPointResult[T]

	switch p.Method {
	case PivotPointStandard:
		res.P = (h + l + c) / 3
		res.R1 = 2*res.P - l
		res.S1 = 2*res.P - h
		res.R2 = res.P + (h - l)
		res.S2 = res.P - (h - l)
		res.R3 = h + 2*(res.P-l)
		res.S3 = l - 2*(h-res.P)
		res.R4 = h + 3*(res.P-l)
		res.S4 = l - 3*(h-res.P)

	case PivotPointWoodie:
		res.P = (h + l + 2*currO) / 4
		res.R1 = 2*res.P - l
		res.S1 = 2*res.P - h
		res.R2 = res.P + (h - l)
		res.S2 = res.P - (h - l)
		res.R3 = h + 2*(res.P-l)
		res.S3 = l - 2*(h-res.P)

	case PivotPointCamarilla:
		diff := h - l
		res.P = (h + l + c) / 3
		res.R1 = c + diff*T(1.1)/12
		res.R2 = c + diff*T(1.1)/6
		res.R3 = c + diff*T(1.1)/4
		res.R4 = c + diff*T(1.1)/2
		res.S1 = c - diff*T(1.1)/12
		res.S2 = c - diff*T(1.1)/6
		res.S3 = c - diff*T(1.1)/4
		res.S4 = c - diff*T(1.1)/2

	case PivotPointFibonacci:
		diff := h - l
		res.P = (h + l + c) / 3
		res.R1 = res.P + diff*T(0.382)
		res.S1 = res.P - diff*T(0.382)
		res.R2 = res.P + diff*T(0.618)
		res.S2 = res.P - diff*T(0.618)
		res.R3 = res.P + diff*T(1.000)
		res.S3 = res.P - diff*T(1.000)
	}

	return res
}

// IdlePeriod is the initial period that Pivot Point won't yield any results.
func (p *PivotPoint[T]) IdlePeriod() int {
	return 1
}

// String is the string representation of the Pivot Point instance.
func (p *PivotPoint[T]) String() string {
	var methodStr string
	switch p.Method {
	case PivotPointStandard:
		methodStr = "Standard"
	case PivotPointWoodie:
		methodStr = "Woodie"
	case PivotPointCamarilla:
		methodStr = "Camarilla"
	case PivotPointFibonacci:
		methodStr = "Fibonacci"
	}
	return fmt.Sprintf("PivotPoint(%s)", methodStr)
}
