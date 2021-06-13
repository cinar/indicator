package indicator

import (
	"math"
	"sort"
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
// TODO: Not optimal. Needs to be done with a binary tree and a ring buffer.
func Max(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)

	for i := 0; i < len(values); i++ {
		buffer[i%period] = values[i]
		sort.Float64s(buffer)

		result[i] = buffer[period-1]
	}

	return result
}

// Moving min for the given period.
// TODO: Not optimal. Needs to be done with a binary tree and a ring buffer.
func Min(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)

	for i := 0; i < len(values); i++ {
		buffer[i%period] = values[i]
		sort.Float64s(buffer)

		lowest := 0
		if i < period {
			lowest = lowest - i - 1
		}

		result[i] = buffer[lowest]
	}

	return result
}
