package indicator

// Multiply values with multipler.
func multiply(values []float64, multiplier float64) []float64 {
	result := make([]float64, len(values))

	for i, value := range values {
		result[i] = value * multiplier
	}

	return result
}

// Divide values with divider.
func divide(values []float64, divider float64) []float64 {
	return multiply(values, float64(1)/divider)
}

// Add values1 and values2.
func add(values1, values2 []float64) []float64 {
	if len(values1) != len(values2) {
		panic("not the same length")
	}

	result := make([]float64, len(values1))
	for i := 0; i < len(result); i++ {
		result[i] = values1[i] + values2[i]
	}

	return result
}

// Substract values2 from values1.
func substract(values1, values2 []float64) []float64 {
	substract := multiply(values2, float64(-1))
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
