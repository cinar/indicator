// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// Awesome Oscillator.
//
// Median Price = ((Low + High) / 2).
// AO = 5-Period SMA - 34-Period SMA.
//
// Returns ao.
func AwesomeOscillator(low, high []float64) []float64 {
	medianPrice := divideBy(add(low, high), float64(2))
	sma5 := Sma(5, medianPrice)
	sma34 := Sma(34, medianPrice)
	ao := substract(sma5, sma34)

	return ao
}

// The ChaikinOscillator function measures the momentum of the
// Accumulation/Distribution (A/D) using the Moving Average
// Convergence Divergence (MACD) formula. It takes the
// difference between fast and slow periods EMA of the A/D.
// Cross above the A/D line indicates bullish.
//
// CO = Ema(fastPeriod, AD) - Ema(slowPeriod, AD)
//
// Returns co, ad.
func ChaikinOscillator(fastPeriod, slowPeriod int, low, high, closing []float64, volume []int64) ([]float64, []float64) {
	ad := AccumulationDistribution(high, low, closing, volume)
	co := substract(Ema(fastPeriod, ad), Ema(slowPeriod, ad))

	return co, ad
}

// The DefaultChaikinOscillator function calculates Chaikin
// Oscillator with the most frequently used fast and short
// periods, 3 and 10.
//
// Returns co, ad.
func DefaultChaikinOscillator(low, high, closing []float64, volume []int64) ([]float64, []float64) {
	return ChaikinOscillator(3, 10, low, high, closing, volume)
}

// Ichimoku Cloud. Also known as Ichimoku Kinko Hyo, is a versatile indicator that defines support and
// resistence, identifies trend direction, gauges momentum, and provides trading signals.
//
// Tenkan-sen (Conversion Line) = (9-Period High + 9-Period Low) / 2
// Kijun-sen (Base Line) = (26-Period High + 26-Period Low) / 2
// Senkou Span A (Leading Span A) = (Conversion Line + Base Line) / 2
// Senkou Span B (Leading Span B) = (52-Period High + 52-Period Low) / 2
// Chikou Span (Lagging Span) = Closing plotted 26 days in the past.
//
// Returns conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingSpan
func IchimokuCloud(high, low, closing []float64) ([]float64, []float64, []float64, []float64, []float64) {
	checkSameSize(high, low, closing)

	conversionLine := divideBy(add(Max(9, high), Min(9, low)), float64(2))
	baseLine := divideBy(add(Max(26, high), Min(26, low)), float64(2))
	leadingSpanA := divideBy(add(conversionLine, baseLine), float64(2))
	leadingSpanB := divideBy(add(Max(52, high), Min(52, low)), float64(2))
	laggingLine := shiftRight(26, closing)

	return conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingLine
}

// Relative Strength Index (RSI). It is a momentum indicator that measures the magnitude
// of recent price changes to evaluate overbought and oversold conditions.
//
// RS = Average Gain / Average Loss
// RSI = 100 - (100 / (1 + RS))
//
// Returns rs, rsi
func Rsi(closing []float64) ([]float64, []float64) {
	gains := make([]float64, len(closing))
	losses := make([]float64, len(closing))

	for i := 1; i < len(closing); i++ {
		difference := closing[i] - closing[i-1]

		if difference > 0 {
			gains[i] = difference
			losses[i] = 0
		} else {
			losses[i] = -difference
			gains[i] = 0
		}
	}

	meanGains := Sma(14, gains)
	meanLosses := Sma(14, losses)

	rsi := make([]float64, len(closing))
	rs := make([]float64, len(closing))

	for i := 0; i < len(rsi); i++ {
		rs[i] = meanGains[i] / meanLosses[i]
		rsi[i] = 100 - (100 / (1 + rs[i]))
	}

	return rs, rsi
}

// Stochastic Oscillator. It is a momentum indicator that shows the location of the closing
// relative to high-low range over a set number of periods.
//
// K = (Closing - Lowest Low) / (Highest High - Lowest Low) * 100
// D = 3-Period SMA of K
//
// Returns k, d
func StochasticOscillator(high, low, closing []float64) ([]float64, []float64) {
	checkSameSize(high, low, closing)

	highestHigh14 := Max(14, high)
	lowestLow14 := Min(15, low)

	k := divide(substract(closing, lowestLow14), multiplyBy(substract(highestHigh14, lowestLow14), float64(100)))
	d := Sma(3, k)

	return k, d
}

// Williams R. Determine overbought and oversold.
//
// WR = (Highest High - Closing) / (Highest High - Lowest Low) * -100.
//
// Buy when -80 and below. Sell when -20 and above.
//
// Returns wr.
func WilliamsR(low, high, closing []float64) []float64 {
	period := 14

	highestHigh := Max(period, high)
	lowestLow := Min(period, low)

	result := make([]float64, len(closing))

	for i := 0; i < len(closing); i++ {
		result[i] = (highestHigh[i] - closing[i]) / (highestHigh[i] - lowestLow[i]) * float64(-100)
	}

	return result
}
