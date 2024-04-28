// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend

import (
	"github.com/cinar/indicator/v2/helper"
)

// Ma represents the interface for the Moving Average (MA) indicators.
type Ma[T helper.Number] interface {
	// Compute function takes a channel of numbers and computes the MA.
	Compute(<-chan T) <-chan T

	// IdlePeriod is the initial period that MA won't yield any results.
	IdlePeriod() int

	// String is the string representation of the MA instance.
	String() string
}
