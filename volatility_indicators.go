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

	k := divide(substract(high, low), add(high, low))

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
	std := Std(20, closing)
	std2 := multiplyBy(std, 2)

	middleBand := Sma(20, closing)
	upperBand := add(middleBand, std2)
	lowerBand := substract(middleBand, std2)

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
	result := make([]float64, len(values))
	sma := Sma(period, values)
	sum := float64(0)

	for i, value := range values {
		d1 := math.Pow(value-sma[i], 2)
		count := i + 1
		sum += d1

		if i >= period {
			first := i - period
			sum -= math.Pow(values[first]-sma[first], 2)
			count = period
		}

		result[i] = math.Sqrt(sum / float64(count))
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

	po := divide(multiplyBy(substract(closing, pl), 100), substract(pu, pl))
	spo := Ema(smooth, po)

	return po, spo
}
