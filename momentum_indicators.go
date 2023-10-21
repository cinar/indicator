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
	ao := subtract(sma5, sma34)

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
func ChaikinOscillator(fastPeriod, slowPeriod int, low, high, closing, volume []float64) ([]float64, []float64) {
	ad := AccumulationDistribution(high, low, closing, volume)
	co := subtract(Ema(fastPeriod, ad), Ema(slowPeriod, ad))

	return co, ad
}

// The DefaultChaikinOscillator function calculates Chaikin
// Oscillator with the most frequently used fast and short
// periods, 3 and 10.
//
// Returns co, ad.
func DefaultChaikinOscillator(low, high, closing, volume []float64) ([]float64, []float64) {
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

// Percentage Price Oscillator (PPO). It is a momentum oscillator for the price.
// It is used to indicate the ups and downs based on the price. A breakout is
// confirmed when PPO is positive.
//
// PPO = ((EMA(fastPeriod, prices) - EMA(slowPeriod, prices)) / EMA(longPeriod, prices)) * 100
// Signal = EMA(9, PPO)
// Histogram = PPO - Signal
//
// Returns ppo, signal, histogram
func PercentagePriceOscillator(fastPeriod, slowPeriod, signalPeriod int, price []float64) ([]float64, []float64, []float64) {
	fastEma := Ema(fastPeriod, price)
	slowEma := Ema(slowPeriod, price)
	ppo := multiplyBy(divide(subtract(fastEma, slowEma), slowEma), 100)
	signal := Ema(signalPeriod, ppo)
	histogram := subtract(ppo, signal)

	return ppo, signal, histogram
}

// Default Percentage Price Oscillator calculates it with the default periods of 12, 26, 9.
//
// Returns ppo, signal, histogram
func DefaultPercentagePriceOscillator(price []float64) ([]float64, []float64, []float64) {
	return PercentagePriceOscillator(12, 26, 9, price)
}

// Percentage Volume Oscillator (PVO). It is a momentum oscillator for the volume.
// It is used to indicate the ups and downs based on the volume. A breakout is
// confirmed when PVO is positive.
//
// PVO = ((EMA(fastPeriod, volumes) - EMA(slowPeriod, volumes)) / EMA(longPeriod, volumes)) * 100
// Signal = EMA(9, PVO)
// Histogram = PVO - Signal
//
// Returns pvo, signal, histogram
func PercentageVolumeOscillator(fastPeriod, slowPeriod, signalPeriod int, volume []float64) ([]float64, []float64, []float64) {
	fastEma := Ema(fastPeriod, volume)
	slowEma := Ema(slowPeriod, volume)
	pvo := multiplyBy(divide(subtract(fastEma, slowEma), slowEma), 100)
	signal := Ema(signalPeriod, pvo)
	histogram := subtract(pvo, signal)

	return pvo, signal, histogram
}

// Default Percentage Volume Oscillator calculates it with the default periods of 12, 26, 9.
//
// Returns pvo, signal, histogram
func DefaultPercentageVolumeOscillator(volume []float64) ([]float64, []float64, []float64) {
	return PercentageVolumeOscillator(12, 26, 9, volume)
}

// Relative Strength Index (RSI). It is a momentum indicator that measures the magnitude
// of recent price changes to evaluate overbought and oversold conditions.
//
// RS = Average Gain / Average Loss
// RSI = 100 - (100 / (1 + RS))
//
// Returns rs, rsi
func Rsi(closing []float64) ([]float64, []float64) {
	return RsiPeriod(14, closing)
}

// RSI with 2 period, a mean-reversion trading strategy
// developed by Larry Connors.
//
// REturns rs, rsi
func Rsi2(closing []float64) ([]float64, []float64) {
	return RsiPeriod(2, closing)
}

// RsiPeriod allows to calculate the RSI indicator with a non-standard period.
func RsiPeriod(period int, closing []float64) ([]float64, []float64) {
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

	meanGains := Rma(period, gains)
	meanLosses := Rma(period, losses)

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
	lowestLow14 := Min(14, low)

	k := multiplyBy(divide(subtract(closing, lowestLow14), subtract(highestHigh14, lowestLow14)), float64(100))
	d := Sma(3, fillNaNWith(k, 0))

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
