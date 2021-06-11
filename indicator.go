package indicator

import "sort"

// Moving Average Convergence Divergence (MACD).
//
// MACD = 12-Period EMA - 26-Period EMA.
// Signal = 9-Period EMA of MACD.
//
// Returns MACD, signal.
func Macd(close []float64) ([]float64, []float64) {
	ema12 := Ema(12, close)
	ema26 := Ema(26, close)
	macd := substract(ema12, ema26)
	signal := Ema(9, macd)

	return macd, signal
}

// Bollinger Bands.
//
// Middle Band = 20-Period SMA.
// Upper Band = 20-Period SMA + 2 (20-Period Std)
// Lower Band = 20-Period SMA - 2 (20-Period Std)
//
// Returns middle band, upper band, lower band.
func BollingerBands(close []float64) ([]float64, []float64, []float64) {
	std := Std(20, close)
	std2 := multiply(std, 2)

	middleBand := Sma(20, close)
	upperBand := add(middleBand, std2)
	lowerBand := substract(middleBand, std2)

	return middleBand, upperBand, lowerBand
}

// Awesome Oscillator.
//
// Median Price = ((Low + High) / 2).
// AO = 5-Period SMA - 34-Period SMA.
//
// Returns ao.
func AwesomeOscillator(low, high []float64) []float64 {
	medianPrice := divide(add(low, high), float64(2))
	sma5 := Sma(5, medianPrice)
	sma34 := Sma(34, medianPrice)
	ao := substract(sma5, sma34)

	return ao
}

// Williams R. Determine overbought and oversold.
//
// WR = (Highest High - Close) / (Highest High - Lowest Low) * -100.
//
// Buy when -80 and below. Sell when -20 and above.
//
// Returns wr.
func WilliamsR(low, high, close []float64) []float64 {
	period := 14
	result := make([]float64, len(close))
	lowPeriod := make([]float64, period)
	highPeriod := make([]float64, period)

	for i := 0; i < len(close); i++ {
		lowPeriod[i%period] = low[i]
		sort.Float64s(lowPeriod)

		highPeriod[i%period] = high[i]
		sort.Float64s(highPeriod)

		lowestLowIndex := 0
		if i < period {
			lowestLowIndex = period - i - 1
		}

		highestHighIndex := period - 1

		result[i] = (highPeriod[highestHighIndex] - close[i]) / (highPeriod[highestHighIndex] - lowPeriod[lowestLowIndex]) * float64(-100)
	}

	return result
}
