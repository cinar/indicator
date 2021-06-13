package indicator

import (
	"math"
)

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
	std2 := multiplyBy(std, 2)

	middleBand := Sma(20, close)
	upperBand := add(middleBand, std2)
	lowerBand := substract(middleBand, std2)

	return middleBand, upperBand, lowerBand
}

// Bollinger Band Width. It measures the percentage difference between the
// upper band and the lower band. It decreases as Bollinger Bands narrows
// and increases as Bollinger Bands widens
//
// During a period of rising price volatity the band width widens, and
// during a period of low market volatity band width contracts.
//
// Band Width = (Upper Band - Lower Band) / Middle Band
//
// Returns bandWidth, bandWidthEma90
func BollingerBandWidth(middleBand, upperBand, lowerBand []float64) ([]float64, []float64) {
	checkSameSize(middleBand, upperBand, lowerBand)

	bandWidth := make([]float64, len(middleBand))
	for i := 0; i < len(bandWidth); i++ {
		bandWidth[i] = (upperBand[i] - lowerBand[i]) / middleBand[i]
	}

	bandWidthEma90 := Ema(90, bandWidth)

	return bandWidth, bandWidthEma90
}

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

// Williams R. Determine overbought and oversold.
//
// WR = (Highest High - Close) / (Highest High - Lowest Low) * -100.
//
// Buy when -80 and below. Sell when -20 and above.
//
// Returns wr.
func WilliamsR(low, high, close []float64) []float64 {
	period := 14

	highestHigh := Max(period, high)
	lowestLow := Min(period, low)

	result := make([]float64, len(close))

	for i := 0; i < len(close); i++ {
		result[i] = (highestHigh[i] - close[i]) / (highestHigh[i] - lowestLow[i]) * float64(-100)
	}

	return result
}

// Typical Price. It is another approximation of average price for each
// period and can be used as a filter for moving average systems.
//
// Typical Price = (High + Low + Close) / 3
//
// Returns typical price, 20-Period SMA.
func TypicalPrice(low, high, close []float64) ([]float64, []float64) {
	checkSameSize(high, low, close)

	sma20 := Sma(20, close)

	ta := make([]float64, len(close))
	for i := 0; i < len(ta); i++ {
		ta[i] = (high[i] + low[i] + close[i]) / float64(3)
	}

	return ta, sma20
}

// Relative Strength Index (RSI). It is a momentum indicator that measures the magnitude
// of recent price changes to evaluate overbought and oversold conditions.
//
// RS = Average Gain / Average Loss
// RSI = 100 - (100 / (1 + RS))
//
// Returns rs, rsi
func Rsi(close []float64) ([]float64, []float64) {
	gains, losses := groupPositivesAndNegatives(diff(close, 1))

	meanGains := Sma(14, gains)
	meanLosses := Sma(14, losses)

	rsi := make([]float64, len(close))
	rs := make([]float64, len(close))

	for i := 0; i < len(rsi); i++ {
		rs[i] = meanGains[i] / (float64(-1) * meanLosses[i])
		rsi[i] = 100 - (100 / (1 + rs[i]))
	}

	return rs, rsi
}

// On-Balance Volume (OBV). It is a technical trading momentum indicator that
// uses volume flow to predict changes in stock price.
//
//                   volume, if Close > Close-Prev
// OBV = OBV-Prev +       0, if Close = Close-Prev
//                  -volume, if Close < Close-Prev
//
// Returns obv
func Obv(close []float64, volume []int64) []int64 {
	if len(close) != len(volume) {
		panic("not all same size")
	}

	obv := make([]int64, len(volume))

	for i := 1; i < len(obv); i++ {
		obv[i] = obv[i-1]

		if close[i] > close[i-1] {
			obv[i] += volume[i]
		} else if close[i] < close[i-1] {
			obv[i] -= volume[i]
		}
	}

	return obv
}

// Average True Range (ATR). It is a technical analysis indicator that measures market
// volatility by decomposing the entire range of stock prices for that period.
//
// TR = Max((High - Low), (High - Close), (Close - Low))
// ATR = SMA TR
//
// Returns tr, atr
func Atr(period int, high, low, close []float64) ([]float64, []float64) {
	checkSameSize(high, low, close)

	tr := make([]float64, len(close))

	for i := 0; i < len(tr); i++ {
		tr[i] = math.Max(high[i]-low[i], math.Max(high[i]-close[i], close[i]-low[i]))
	}

	atr := Sma(period, tr)

	return tr, atr
}

// Chandelier Exit. It sets a trailing stop-loss based on the Average True Value (ATR).
//
// Chandelier Exit Long = 22-Period SMA High - ATR(22) * 3
// Chandelier Exit Short = 22-Period SMA Low + ATR(22) * 3
//
// Returns chandelierExitLong, chandelierExitShort
func ChandelierExit(high, low, close []float64) ([]float64, []float64) {
	_, atr22 := Atr(22, high, low, close)
	highestHigh22 := Max(22, high)
	lowestLow22 := Min(22, low)

	chandelierExitLong := make([]float64, len(close))
	chandelierExitShort := make([]float64, len(close))

	for i := 0; i < len(chandelierExitLong); i++ {
		chandelierExitLong[i] = highestHigh22[i] - (atr22[i] * float64(3))
		chandelierExitShort[i] = lowestLow22[i] + (atr22[i] * float64(3))
	}

	return chandelierExitLong, chandelierExitShort
}

// Ichimoku Cloud. Also known as Ichimoku Kinko Hyo, is a versatile indicator that defines support and
// resistence, identifies trend direction, gauges momentum, and provides trading signals.
//
// Tenkan-sen (Conversion Line) = (9-Period High + 9-Period Low) / 2
// Kijun-sen (Base Line) = (26-Period High + 26-Period Low) / 2
// Senkou Span A (Leading Span A) = (Conversion Line + Base Line) / 2
// Senkou Span B (Leading Span B) = (52-Period High + 52-Period Low) / 2
// Chikou Span (Lagging Span) = Close plotted 26 days in the past.
//
// Returns conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingSpan
func IchimokuCloud(high, low, close []float64) ([]float64, []float64, []float64, []float64, []float64) {
	checkSameSize(high, low, close)

	conversionLine := divideBy(add(Max(9, high), Min(9, low)), float64(2))
	baseLine := divideBy(add(Max(26, high), Min(26, low)), float64(2))
	leadingSpanA := divideBy(add(conversionLine, baseLine), float64(2))
	leadingSpanB := divideBy(add(Max(52, high), Min(52, low)), float64(2))
	laggingLine := shiftRight(26, close)

	return conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingLine
}

// Stochastic Oscillator. It is a momentum indicator that shows the location of the close
// relative to high-low range over a set number of periods.
//
// K = (Close - Lowest Low) / (Highest High - Lowest Low) * 100
// D = 3-Period SMA of K
//
// Returns k, d
func StochasticOscillator(high, low, close []float64) ([]float64, []float64) {
	checkSameSize(high, low, close)

	highestHigh14 := Max(14, high)
	lowestLow14 := Min(15, low)

	k := divide(substract(close, lowestLow14), multiplyBy(substract(highestHigh14, lowestLow14), float64(100)))
	d := Sma(3, k)

	return k, d
}
