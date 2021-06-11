package indicator

import "math"

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
