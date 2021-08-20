package indicator

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
	if before >= len(values) {
		panic("before greather or equals to size")
	}

	result := make([]float64, len(values))

	for i := before; i < len(values); i++ {
		result[i] = values[i] - values[i-before]
	}

	return result
}

// Groups positives and negatives. Ignores zeros.
func groupPositivesAndNegatives(values []float64) ([]float64, []float64) {
	var positives, negatives []float64

	for _, value := range values {
		if value > 0 {
			positives = append(positives, value)
		} else if value < 0 {
			negatives = append(negatives, value)
		}
	}

	return positives, negatives
}

// Shift right for period.
func shiftRight(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	for i := period; i < len(result); i++ {
		result[i] = values[i-period]
	}

	return result
}
