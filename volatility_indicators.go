// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"math"
)

// Acceleration Bands. Plots upper and lower envelope bands
// around a simple moving average.
//
// Upper Band = SMA(High * (1 + 4 * (High - Low) / (High + Low)))
// Middle Band = SMA(Closing)
// Lower Band = SMA(Low * (1 - 4 * (High - Low) / (High + Low)))
//
// Returns upper band, middle band, lower band.
func AccelerationBands(high, low, closing []float64) ([]float64, []float64, []float64) {
	checkSameSize(high, low, closing)

	k := divide(subtract(high, low), add(high, low))

	upperBand := Sma(20, multiply(high, addBy(multiplyBy(k, 4), 1)))
	middleBand := Sma(20, closing)
	lowerBand := Sma(20, multiply(low, addBy(multiplyBy(k, -4), 1)))

	return upperBand, middleBand, lowerBand
}

// Average True Range (ATR). It is a technical analysis indicator that measures market
// volatility by decomposing the entire range of stock prices for that period.
//
// TR = Max((High - Low), (High - Closing), (Closing - Low))
// ATR = SMA TR
//
// Returns tr, atr
func Atr(period int, high, low, closing []float64) ([]float64, []float64) {
	checkSameSize(high, low, closing)

	tr := make([]float64, len(closing))

	for i := 0; i < len(tr); i++ {
		tr[i] = math.Max(high[i]-low[i], math.Max(high[i]-closing[i], closing[i]-low[i]))
	}

	atr := Sma(period, tr)

	return tr, atr
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

// Bollinger Bands.
//
// Middle Band = 20-Period SMA.
// Upper Band = 20-Period SMA + 2 (20-Period Std)
// Lower Band = 20-Period SMA - 2 (20-Period Std)
//
// Returns middle band, upper band, lower band.
func BollingerBands(closing []float64) ([]float64, []float64, []float64) {
	middleBand := Sma(20, closing)

	std := StdFromSma(20, closing, middleBand)
	std2 := multiplyBy(std, 2)

	upperBand := add(middleBand, std2)
	lowerBand := subtract(middleBand, std2)

	return middleBand, upperBand, lowerBand
}

// Chandelier Exit. It sets a trailing stop-loss based on the Average True Value (ATR).
//
// Chandelier Exit Long = 22-Period SMA High - ATR(22) * 3
// Chandelier Exit Short = 22-Period SMA Low + ATR(22) * 3
//
// Returns chandelierExitLong, chandelierExitShort
func ChandelierExit(high, low, closing []float64) ([]float64, []float64) {
	_, atr22 := Atr(22, high, low, closing)
	highestHigh22 := Max(22, high)
	lowestLow22 := Min(22, low)

	chandelierExitLong := make([]float64, len(closing))
	chandelierExitShort := make([]float64, len(closing))

	for i := 0; i < len(chandelierExitLong); i++ {
		chandelierExitLong[i] = highestHigh22[i] - (atr22[i] * float64(3))
		chandelierExitShort[i] = lowestLow22[i] + (atr22[i] * float64(3))
	}

	return chandelierExitLong, chandelierExitShort
}

// Standard deviation.
func Std(period int, values []float64) []float64 {
	return StdFromSma(period, values, Sma(period, values))
}

// Standard deviation from the given SMA.
func StdFromSma(period int, values, sma []float64) []float64 {
	result := make([]float64, len(values))

	sum2 := 0.0
	for i, v := range values {
		sum2 += v * v
		if i < period-1 {
			result[i] = 0.0
		} else {
			result[i] = math.Sqrt(sum2/float64(period) - sma[i]*sma[i])
			w := values[i-(period-1)]
			sum2 -= w * w
		}
	}

	return result
}

// ProjectionOscillator calculates the Projection Oscillator (PO). The PO
// uses the linear regression slope, along with highs and lows.
//
// Period defines the moving window to calculates the PO, and the smooth
// period defines the moving windows to take EMA of PO.
//
// PL = Min(period, (high + MLS(period, x, high)))
// PU = Max(period, (low + MLS(period, x, low)))
// PO = 100 * (Closing - PL) / (PU - PL)
// SPO = EMA(smooth, PO)
//
// Returns po, spo.
func ProjectionOscillator(period, smooth int, high, low, closing []float64) ([]float64, []float64) {
	x := generateNumbers(0, float64(len(closing)), 1)
	mHigh, _ := MovingLeastSquare(period, x, high)
	mLow, _ := MovingLeastSquare(period, x, low)

	vHigh := add(high, multiply(mHigh, x))
	vLow := add(low, multiply(mLow, x))

	pu := Max(period, vHigh)
	pl := Min(period, vLow)

	po := divide(multiplyBy(subtract(closing, pl), 100), subtract(pu, pl))
	spo := Ema(smooth, po)

	return po, spo
}

// The Ulcer Index (UI) measures downside risk. The index increases in value
// as the price moves farther away from a recent high and falls as the price
// rises to new highs.
//
// High Closings = Max(period, Closings)
// Percentage Drawdown = 100 * ((Closings - High Closings) / High Closings)
// Squared Average = Sma(period, Percent Drawdown * Percent Drawdown)
// Ulcer Index = Sqrt(Squared Average)
//
// Returns ui.
func UlcerIndex(period int, closing []float64) []float64 {
	highClosing := Max(period, closing)
	percentageDrawdown := multiplyBy(divide(subtract(closing, highClosing), highClosing), 100)
	squaredAverage := Sma(period, multiply(percentageDrawdown, percentageDrawdown))
	ui := sqrt(squaredAverage)

	return ui
}

// The default ulcer index with the default period of 14.
func DefaultUlcerIndex(closing []float64) []float64 {
	return UlcerIndex(14, closing)
}

// The Donchian Channel (DC) calculates three lines generated by moving average
// calculations that comprise an indicator formed by upper and lower bands
// around a midrange or median band. The upper band marks the highest
// price of an asset while the lower band marks the lowest price of
// an asset, and the area between the upper and lower bands
// represents the Donchian Channel.
//
// Upper Channel = Mmax(period, closings)
// Lower Channel = Mmin(period, closings)
// Middle Channel = (Upper Channel + Lower Channel) / 2
//
// Returns upperChannel, middleChannel, lowerChannel.
func DonchianChannel(period int, closing []float64) ([]float64, []float64, []float64) {
	upperChannel := Max(period, closing)
	lowerChannel := Min(period, closing)
	middleChannel := divideBy(add(upperChannel, lowerChannel), 2)

	return upperChannel, middleChannel, lowerChannel
}

// The Keltner Channel (KC) provides volatility-based bands that are placed
// on either side of an asset's price and can aid in determining the
// direction of a trend.
//
// Middle Line = EMA(period, closings)
// Upper Band = EMA(period, closings) + 2 * ATR(period, highs, lows, closings)
// Lower Band = EMA(period, closings) - 2 * ATR(period, highs, lows, closings)
//
// Returns upperBand, middleLine, lowerBand.
func KeltnerChannel(period int, high, low, closing []float64) ([]float64, []float64, []float64) {
	_, atr := Atr(period, high, low, closing)
	atr2 := multiplyBy(atr, 2)

	middleLine := Ema(period, closing)
	upperBand := add(middleLine, atr2)
	lowerBand := subtract(middleLine, atr2)

	return upperBand, middleLine, lowerBand
}

// The default keltner channel with the default period of 20.
func DefaultKeltnerChannel(high, low, closing []float64) ([]float64, []float64, []float64) {
	return KeltnerChannel(20, high, low, closing)
}
