// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume

import (
	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultNviInitial is the default initial for the NVI.
	DefaultNviInitial = 1000
)

// Nvi holds configuration parameters for calculating the Negative Volume Index (NVI). It is a cumulative
// indicator using the change in volume to decide when the smart money is active.
//
// If Volume is greather than Previous Volume:
//
//	NVI = Previous NVI
//
// Otherwise:
//
//	NVI = Previous NVI + (((Closing - Previous Closing) / Previous Closing) * Previous NVI)
//
// Example:
//
//	nvi := volume.NewNvi[float64]()
//	result := nvi.Compute(closings, volumes)
type Nvi[T helper.Number] struct {
	// Initial is the initial NVI value.
	Initial T
}

// NewNvi function initializes a new NVI instance with the default parameters.
func NewNvi[T helper.Number]() *Nvi[T] {
	initial := DefaultNviInitial

	return &Nvi[T]{
		Initial: T(initial),
	}
}

// Compute function takes a channel of numbers and computes the NVI.
func (n *Nvi[T]) Compute(closings, volumes <-chan T) <-chan T {
	closingRatios := helper.ChangeRatio(closings, 1)
	volumeChanges := helper.Change(volumes, 1)

	previous := n.Initial

	return helper.Operate(closingRatios, volumeChanges, func(closingRatio, volumeChange T) T {
		// If Volume is greather than Previous Volume:
		//	NVI = Previous NVI
		current := previous

		// Otherwise:
		//	NVI = Previous NVI + (((Closing - Previous Closing) / Previous Closing) * Previous NVI)
		if volumeChange <= 0 {
			current += closingRatio * previous
		}

		previous = current
		return current
	})
}

// IdlePeriod is the initial period that NVI won't yield any results.
func (*Nvi[T]) IdlePeriod() int {
	return 1
}
