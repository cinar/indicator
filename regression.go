// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

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

// Moving least square over a period.
//
// y = mx + b
// b = y-intercept
// y = slope
//
// m = (n * sumXY - sumX * sumY) / (n * sumX2 - sumX * sumX)
// b = (sumY - m * sumX) / n
func MovingLeastSquare(period int, x, y []float64) ([]float64, []float64) {
	checkSameSize(x, y)

	m := make([]float64, len(x))
	b := make([]float64, len(x))

	var sumX, sumX2, sumY, sumXY float64

	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumX2 += x[i] * x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]

		n := float64(i + 1)

		if i >= period {
			sumX -= x[i-period]
			sumX2 -= x[i-period] * x[i-period]
			sumY -= y[i-period]
			sumXY -= x[i-period] * y[i-period]
			n = float64(period)
		}

		m[i] = ((n * sumXY) - (sumX * sumY)) / ((n * sumX2) - (sumX * sumX))
		b[i] = (sumY - (m[i] * sumX)) / n
	}

	return m, b
}

// Linear regression using least square method.
//
// y = mx + b
func LinearRegressionUsingLeastSquare(x, y []float64) []float64 {
	m, b := LeastSquare(x, y)

	r := make([]float64, len(x))

	for i := 0; i < len(r); i++ {
		r[i] = (m * x[i]) + b
	}

	return r
}

// Moving linear regression using least square.
//
// y = mx + b
func MovingLinearRegressionUsingLeastSquare(period int, x, y []float64) []float64 {
	m, b := MovingLeastSquare(period, x, y)

	r := make([]float64, len(x))

	for i := 0; i < len(r); i++ {
		r[i] = (m[i] * x[i]) + b[i]
	}

	return r
}
