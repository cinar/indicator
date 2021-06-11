package indicator

// Multiply values with multipler.
func Multiply(values []float64, multiplier float64) []float64 {
	result := make([]float64, len(values))

	for i, value := range values {
		result[i] = value * multiplier
	}

	return result
}

// Divide values with divider.
func Divide(values []float64, divider float64) []float64 {
	return Multiply(values, float64(1)/divider)
}

// Add values1 and values2.
func Add(values1, values2 []float64) []float64 {
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
func Substract(values1, values2 []float64) []float64 {
	substract := Multiply(values2, float64(-1))
	return Add(values1, substract)
}
