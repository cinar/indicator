// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"math"
)

// Check values same size.
func checkSameSize(values ...[]float64) {
	if len(values) < 2 {
		return
	}

	n := len(values[0])

	for i := 1; i < len(values); i++ {
		if len(values[i]) != n {
			panic("not all same size")
		}
	}
}

// Multiply values by multipler.
func multiplyBy(values []float64, multiplier float64) []float64 {
	result := make([]float64, len(values))

	for i, value := range values {
		result[i] = value * multiplier
	}

	return result
}

// Multiply values1 and values2.
func multiply(values1, values2 []float64) []float64 {
	checkSameSize(values1, values2)

	result := make([]float64, len(values1))

	for i := 0; i < len(result); i++ {
		result[i] = values1[i] * values2[i]
	}

	return result
}

// Divide values by divider.
func divideBy(values []float64, divider float64) []float64 {
	return multiplyBy(values, float64(1)/divider)
}

// Divide values1 by values2.
func divide(values1, values2 []float64) []float64 {
	checkSameSize(values1, values2)

	result := make([]float64, len(values1))

	for i := 0; i < len(result); i++ {
		result[i] = values1[i] / values2[i]
	}

	return result
}

// Add values1 and values2.
func add(values1, values2 []float64) []float64 {
	checkSameSize(values1, values2)

	result := make([]float64, len(values1))
	for i := 0; i < len(result); i++ {
		result[i] = values1[i] + values2[i]
	}

	return result
}

// Add addition to values.
func addBy(values []float64, addition float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(result); i++ {
		result[i] = values[i] + addition
	}

	return result
}

// Substract values2 from values1.
func substract(values1, values2 []float64) []float64 {
	substract := multiplyBy(values2, float64(-1))
	return add(values1, substract)
}

// Difference between current and before values.
func diff(values []float64, before int) []float64 {
	return substract(values, shiftRight(1, values))
}

// Percent difference between current and before values.
func percentDiff(values []float64, before int) []float64 {
	result := make([]float64, len(values))

	for i := before; i < len(values); i++ {
		result[i] = (values[i] - values[i-before]) / values[i-before]
	}

	return result
}

// Shift right for period.
func shiftRight(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	for i := period; i < len(result); i++ {
		result[i] = values[i-period]
	}

	return result
}

// Round value to digits.
func roundDigits(value float64, digits int) float64 {
	n := math.Pow(10, float64(digits))

	return math.Round(value*n) / n
}

// Generate numbers.
func generateNumbers(begin, end, step float64) []float64 {
	n := int(math.Round((end - begin) / step))

	numbers := make([]float64, n)

	for i := 0; i < n; i++ {
		numbers[i] = begin + (step * float64(i))
	}

	return numbers
}

// Convets the []int64 to []float64.
func asFloat64(values []int64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		result[i] = float64(values[i])
	}

	return result
}

// Calculate power of base with exponent.
func pow(base []float64, exponent float64) []float64 {
	result := make([]float64, len(base))

	for i := 0; i < len(result); i++ {
		result[i] = math.Pow(base[i], exponent)
	}

	return result
}

// Extact sign.
func extractSign(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(result); i++ {
		if values[i] >= 0 {
			result[i] = 1
		} else {
			result[i] = -1
		}
	}

	return result
}

// Keep positives.
func keepPositives(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		if values[i] > 0 {
			result[i] = values[i]
		} else {
			result[i] = 0
		}
	}

	return result
}

// Keep negatives.
func keepNegatives(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		if values[i] < 0 {
			result[i] = values[i]
		} else {
			result[i] = 0
		}
	}

	return result
}
