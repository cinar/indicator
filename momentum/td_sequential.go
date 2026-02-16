// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"github.com/cinar/indicator/v2/helper"
)

const (
	// DefaultTdSequentialLookback is the default lookback period for comparing closes.
	DefaultTdSequentialLookback = 4

	// DefaultTdSequentialCountdownLookback is the default lookback period for countdown comparison.
	DefaultTdSequentialCountdownLookback = 2

	// DefaultTdSequentialSetupPeriod is the default setup period (9).
	DefaultTdSequentialSetupPeriod = 9

	// DefaultTdSequentialCountdownPeriod is the default countdown period (13).
	DefaultTdSequentialCountdownPeriod = 13
)

// TdSequential represents the configuration parameters for calculating the
// Tom DeMark's TD Sequential indicator. TD Sequential is a momentum indicator
// that identifies potential trend exhaustion and reversals.
//
// The indicator consists of two phases:
//
//	TD Setup: Counts 9 consecutive closes higher (sell) or lower (buy) than
//	the close 4 bars ago.
//
//	TD Countdown: After a completed setup, counts 13 closes higher (sell) or
//	lower (buy) than the close 2 bars ago.
//
// Example:
//
//	td := momentum.NewTdSequential[float64]()
//	buySetup, sellSetup, buyCountdown, sellCountdown := td.Compute(closings)
type TdSequential[T helper.Number] struct {
	// Lookback is the number of bars to look back for comparison in the setup phase.
	Lookback int

	// CountdownLookback is the number of bars to look back for comparison in the countdown phase.
	CountdownLookback int

	// SetupPeriod is the number of consecutive closes required to complete a setup.
	SetupPeriod int

	// CountdownPeriod is the number of closes required to complete a countdown.
	CountdownPeriod int
}

// NewTdSequential function initializes a new TD Sequential instance with default parameters.
func NewTdSequential[T helper.Number]() *TdSequential[T] {
	return &TdSequential[T]{
		Lookback:          DefaultTdSequentialLookback,
		CountdownLookback: DefaultTdSequentialCountdownLookback,
		SetupPeriod:       DefaultTdSequentialSetupPeriod,
		CountdownPeriod:   DefaultTdSequentialCountdownPeriod,
	}
}

// lessThan compares two generic numbers and returns true if a < b.
func lessThan[T helper.Number](a, b T) bool {
	return float64(a) < float64(b)
}

// greaterThan compares two generic numbers and returns true if a > b.
func greaterThan[T helper.Number](a, b T) bool {
	return float64(a) > float64(b)
}

// lessOrEqual compares two generic numbers and returns true if a <= b.
func lessOrEqual[T helper.Number](a, b T) bool {
	return float64(a) <= float64(b)
}

// greaterOrEqual compares two generic numbers and returns true if a >= b.
func greaterOrEqual[T helper.Number](a, b T) bool {
	return float64(a) >= float64(b)
}

// Compute function takes a channel of numbers and computes the TD Sequential indicator.
// Returns four channels: buySetup, sellSetup, buyCountdown, sellCountdown.
func (t *TdSequential[T]) Compute(closings <-chan T) (<-chan T, <-chan T, <-chan T, <-chan T) {
	closings = helper.Buffered(closings, t.Lookback+t.CountdownLookback)

	buySetup := make(chan T)
	sellSetup := make(chan T)
	buyCountdown := make(chan T)
	sellCountdown := make(chan T)

	go func() {
		defer close(buySetup)
		defer close(sellSetup)
		defer close(buyCountdown)
		defer close(sellCountdown)

		var currentBuySetup, currentSellSetup T
		var buyCountdownCount, sellCountdownCount int
		inBuyCountdown := false
		inSellCountdown := false
		closeHistory := make([]T, 0, t.Lookback+t.CountdownLookback+1)

		for current := range closings {
			closeHistory = append(closeHistory, current)
			if len(closeHistory) <= t.Lookback {
				buySetup <- 0
				sellSetup <- 0
				buyCountdown <- 0
				sellCountdown <- 0
				continue
			}

			prevClose := closeHistory[len(closeHistory)-1-t.Lookback]

			// Setup phase - buy (close < close 4 bars ago)
			if lessThan(current, prevClose) {
				if float64(currentBuySetup) >= 0 {
					currentBuySetup = T(float64(currentBuySetup) + 1)
				} else {
					currentBuySetup = 1
				}
			} else {
				currentBuySetup = 0
			}

			// Setup phase - sell (close > close 4 bars ago)
			if greaterThan(current, prevClose) {
				if float64(currentSellSetup) <= 0 {
					currentSellSetup = T(float64(currentSellSetup) - 1)
				} else {
					currentSellSetup = -1
				}
			} else {
				currentSellSetup = 0
			}

			// Check if setup completed
			if float64(currentBuySetup) >= float64(t.SetupPeriod) {
				inBuyCountdown = true
			}
			if float64(currentSellSetup) <= -float64(t.SetupPeriod) {
				inSellCountdown = true
			}

			// Countdown phase - buy (close <= close 2 bars ago)
			if inBuyCountdown && buyCountdownCount < t.CountdownPeriod {
				if len(closeHistory) > t.CountdownLookback {
					cdPrevClose := closeHistory[len(closeHistory)-1-t.CountdownLookback]
					if lessOrEqual(current, cdPrevClose) {
						buyCountdownCount++
					}
				}
			}

			// Countdown phase - sell (close >= close 2 bars ago)
			if inSellCountdown && sellCountdownCount < t.CountdownPeriod {
				if len(closeHistory) > t.CountdownLookback {
					cdPrevClose := closeHistory[len(closeHistory)-1-t.CountdownLookback]
					if greaterOrEqual(current, cdPrevClose) {
						sellCountdownCount++
					}
				}
			}

			// Reset countdown when completed
			if buyCountdownCount >= t.CountdownPeriod {
				buyCountdownCount = 0
				inBuyCountdown = false
			}
			if sellCountdownCount >= t.CountdownPeriod {
				sellCountdownCount = 0
				inSellCountdown = false
			}

			buySetup <- currentBuySetup
			sellSetup <- currentSellSetup
			buyCountdown <- T(buyCountdownCount)
			sellCountdown <- T(sellCountdownCount)
		}
	}()

	return buySetup, sellSetup, buyCountdown, sellCountdown
}

// IdlePeriod is the initial period that TD Sequential won't yield meaningful results.
func (t *TdSequential[T]) IdlePeriod() int {
	return t.Lookback + t.SetupPeriod + t.CountdownPeriod
}
