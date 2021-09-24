package indicator

// Least square.
//
// y = mx + b
// b = y-intercept
// y = slope
//
// m = (n * sumXY - sumX * sumY) / (n * sumX2 - sumX * sumX)
// b = (sumY - m * sumX) / n
func LeastSquare(x, y []float64) (float64, float64) {
	checkSameSize(x, y)

	var sumX, sumX2, sumY, sumXY float64

	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumX2 += x[i] * x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
	}

	n := float64(len(x))
	m := ((n * sumXY) - (sumX * sumY)) / ((n * sumX2) - (sumX * sumX))
	b := (sumY - (m * sumX)) / n

	return m, b
}

// Linear regression using least square method.
//
// y = mx + b
func LinearRegressionUsingLeastSquare(x, y []float64) []float64 {
	m, b := LeastSquare(x, y)

	r := make([]float64, len(y))

	for i := 0; i < len(r); i++ {
		r[i] = (m * x[i]) + b
	}

	return r
}
