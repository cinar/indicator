// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"github.com/cinar/indicator/container/bst"
	"math"
)

// Simple Moving Average (SMA).
func Sma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1
		sum += value

		if i >= period {
			sum -= values[i-period]
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}

// Standard deviation.
func Std(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sma := Sma(period, values)
	sum := float64(0)

	for i, value := range values {
		d1 := math.Pow(value-sma[i], 2)
		count := i + 1
		sum += d1

		if i >= period {
			first := i - period
			sum -= math.Pow(values[first]-sma[first], 2)
			count = period
		}

		result[i] = math.Sqrt(sum / float64(count))
	}

	return result
}

// Exponential Moving Average (EMA).
func Ema(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	k := float64(2) / float64(1+period)

	for i, value := range values {
		if i > 0 {
			result[i] = (value * k) + (result[i-1] * float64(1-k))
		} else {
			result[i] = value
		}
	}

	return result
}

// Moving max for the given period.
func Max(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)
	bst := bst.New()

	for i := 0; i < len(values); i++ {
		bst.Insert(values[i])

		if i >= period {
			bst.Remove(buffer[i%period])
		}

		buffer[i%period] = values[i]
		result[i] = bst.Max().(float64)
	}

	return result
}

// Moving min for the given period.
func Min(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)
	bst := bst.New()

	for i := 0; i < len(values); i++ {
		bst.Insert(values[i])

		if i >= period {
			bst.Remove(buffer[i%period])
		}

		buffer[i%period] = values[i]
		result[i] = bst.Min().(float64)
	}

	return result
}

// Since last values change.
func Since(values []float64) []int {
	result := make([]int, len(values))

	lastValue := math.NaN()
	sinceLast := 0

	for i := 0; i < len(values); i++ {
		value := values[i]

		if value != lastValue {
			lastValue = value
			sinceLast = 0
		} else {
			sinceLast++
		}

		result[i] = sinceLast
	}

	return result
}

// Tema calculates the Triple Exponential Moving Average (TEMA).
//
// TEMA = (3 * EMA1) - (3 * EMA2) + EMA3
// EMA1 = EMA(values)
// EMA2 = EMA(EMA1)
// EMA3 = EMA(EMA2)
//
// Returns tema.
func Tema(period int, values []float64) []float64 {
	ema1 := Ema(period, values)
	ema2 := Ema(period, ema1)
	ema3 := Ema(period, ema2)

	tema := add(substract(multiplyBy(ema1, 3), multiplyBy(ema2, 3)), ema3)

	return tema
}

// Dema calculates the Double Exponential Moving Average (DEMA).
//
// DEMA = (2 * EMA(values)) - EMA(EMA(values))
//
// Returns dema.
func Dema(period int, values []float64) []float64 {
	ema1 := Ema(period, values)
	ema2 := Ema(period, ema1)

	dema := substract(multiplyBy(ema1, 2), ema2)

	return dema
}

// Trima function calculates the Triangular Moving Average (TRIMA).
//
// If period is even:
//   TRIMA = SMA(period / 2, SMA((period / 2) + 1, values))
// If period is odd:
//   TRIMA = SMA((period + 1) / 2, SMA((period + 1) / 2, values))
//
// Returns trima.
func Trima(period int, values []float64) []float64 {
	var n1, n2 int

	if period%2 == 0 {
		n1 = period / 2
		n2 = n1 + 1
	} else {
		n1 = (period + 1) / 2
		n2 = n1
	}

	trima := Sma(n1, Sma(n2, values))

	return trima
}
