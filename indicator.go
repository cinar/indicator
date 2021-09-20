package indicator

import (
	"math"
)

// Trend indicator.
type Trend int

const (
	// Falling trend.
	Falling Trend = iota

	// Rising trend.
	Rising
)

const (
	psarAfStep = 0.02
	psarAfMax  = 0.20
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
	gains := make([]float64, len(close))
	losses := make([]float64, len(close))

	for i := 1; i < len(close); i++ {
		difference := close[i] - close[i-1]

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

	rsi := make([]float64, len(close))
	rs := make([]float64, len(close))

	for i := 0; i < len(rsi); i++ {
		rs[i] = meanGains[i] / meanLosses[i]
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

// Aroon Indicator. It is a technical indicator that is used to identify trend changes
// in the price of a stock, as well as the strength of that trend. It consists of two
// lines, Arron Up, and Aroon Down. The Aroon Up line measures the strength of the
// uptrend, and the Aroon Down measures the strength of the downtrend. When Aroon Up
// is above Aroon Down, it indicates bullish price, and when Aroon Down is above
// Aroon Up, it indicates bearish price.
//
// Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
// Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
//
// Returns aroonUp, aroonDown
func Aroon(high, low []float64) ([]float64, []float64) {
	checkSameSize(high, low)

	sinceLastHigh25 := Since(Max(25, high))
	sinceLastLow25 := Since(Min(25, low))

	aroonUp := make([]float64, len(high))
	aroonDown := make([]float64, len(high))

	for i := 0; i < len(aroonUp); i++ {
		aroonUp[i] = (float64(25-sinceLastHigh25[i]) / 25) * 100
		aroonDown[i] = (float64(25-sinceLastLow25[i]) / 25) * 100
	}

	return aroonUp, aroonDown
}

// Parabolic SAR. It is a popular technical indicator for identifying the trend
// and as a trailing stop.
//
// PSAR = PSAR[i - 1] - ((PSAR[i - 1] - EP) * AF)
//
// If the trend is Falling:
//  - PSAR is the maximum of PSAR or the previous two high values.
//  - If the current high is greather than or equals to PSAR, use EP.
//
// If the trend is Rising:
//  - PSAR is the minimum of PSAR or the previous two low values.
//  - If the current low is less than or equals to PSAR, use EP.
//
// If PSAR is greater than the closing, trend is falling, and the EP
// is set to the minimum of EP or the low.
//
// If PSAR is lower than or equals to the closing, trend is rising, and the EP
// is set to the maximum of EP or the high.
//
// If the trend is the same, and AF is less than 0.20, increment it by 0.02.
// If the trend is not the same, set AF to 0.02.
//
// Based on video https://www.youtube.com/watch?v=MuEpGBAH7pw&t=0s.
//
// Returns psar, trend
func ParabolicSar(high, low, close []float64) ([]float64, []Trend) {
	checkSameSize(high, low)

	trend := make([]Trend, len(high))
	psar := make([]float64, len(high))

	var af, ep float64

	trend[0] = Falling
	psar[0] = high[0]
	af = psarAfStep
	ep = low[0]

	for i := 1; i < len(psar); i++ {
		psar[i] = psar[i-1] - ((psar[i-1] - ep) * af)

		if trend[i-1] == Falling {
			psar[i] = math.Max(psar[i], high[i-1])
			if i > 1 {
				psar[i] = math.Max(psar[i], high[i-2])
			}

			if high[i] >= psar[i] {
				psar[i] = ep
			}
		} else {
			psar[i] = math.Min(psar[i], low[i-1])
			if i > 1 {
				psar[i] = math.Min(psar[i], low[i-2])
			}

			if low[i] <= psar[i] {
				psar[i] = ep
			}
		}

		prevEp := ep

		if psar[i] > close[i] {
			trend[i] = Falling
			ep = math.Min(ep, low[i])
		} else {
			trend[i] = Rising
			ep = math.Max(ep, high[i])
		}

		if trend[i] != trend[i-1] {
			af = psarAfStep
		} else if prevEp != ep && af < psarAfMax {
			af += psarAfStep
		}
	}

	return psar, trend
}

// Vortex Indicator. It provides two oscillators that capture positive and
// negative trend movement. A bullish signal triggers when the positive
// trend indicator crosses above the negative trend indicator or a key
// level. A bearish signal triggers when the negative trend indicator
// crosses above the positive trend indicator or a key level.
//
// +VM = Abs(Current High - Prior Low)
// -VM = Abd(Current Low - Prior High)
//
// +VM14 = 14-Period Sum of +VM
// -VM14 = 14-Period Sum of -VM
//
// TR = Max((High[i]-Low[i]), Abs(High[i]-Close[i-1]), Abs(Low[i]-Close[i-1]))
// TR14 = 14-Period Sum of TR
//
// +VI14 = +VM14 / TR14
// -VI14 = -VM14 / TR14
//
// Based on https://school.stockcharts.com/doku.php?id=technical_indicators:vortex_indicator
//
// Returns plusVi, minusVi
func Vortex(high, low, close []float64) ([]float64, []float64) {
	checkSameSize(high, low, close)

	period := 14

	plusVi := make([]float64, len(high))
	minusVi := make([]float64, len(high))

	plusVm := make([]float64, period)
	minusVm := make([]float64, period)
	tr := make([]float64, period)

	var plusVmSum, minusVmSum, trSum float64

	for i := 1; i < len(high); i++ {
		j := i % period

		plusVmSum -= plusVm[j]
		plusVm[j] = math.Abs(high[i] - low[i-1])
		plusVmSum += plusVm[j]

		minusVmSum -= minusVm[j]
		minusVm[j] = math.Abs(low[i] - high[i-1])
		minusVmSum += minusVm[j]

		highLow := high[i] - low[i]
		highPrevClose := math.Abs(high[i] - close[i-1])
		lowPrevClose := math.Abs(low[i] - close[i-1])

		trSum -= tr[j]
		tr[j] = math.Max(highLow, math.Max(highPrevClose, lowPrevClose))
		trSum += tr[j]

		plusVi[i] = plusVmSum / trSum
		minusVi[i] = minusVmSum / trSum
	}

	return plusVi, minusVi
}

// Acceleration Bands. Plots upper and lower envelope bands
// around a simple moving average.
//
// Upper Band = SMA(High * (1 + 4 * (High - Low) / (High + Low)))
// Middle Band = SMA(Close)
// Lower Band = SMA(Low * (1 - 4 * (High - Low) / (High + Low)))
//
// Returns upper band, middle band, lower band.
func AccelerationBands(high, low, close []float64) ([]float64, []float64, []float64) {
	checkSameSize(high, low, close)

	k := divide(substract(high, low), add(high, low))

	upperBand := Sma(20, multiply(high, addBy(multiplyBy(k, 4), 1)))
	middleBand := Sma(20, close)
	lowerBand := Sma(20, multiply(low, addBy(multiplyBy(k, -4), 1)))

	return upperBand, middleBand, lowerBand
}

// Accumulation/Distribution Indicator (A/D). Cumulative indicator
// that uses volume and price to assess whether a stock is
// being accumulated or distributed.
//
// MFM = ((Close - Low) - (High - Close)) / (High - Low)
// MFV = MFM * Period Volume
// AD = Previous AD + CMFV
//
// Returns ad.
func AccumulationDistribution(high, low, close []float64, volume []int64) []float64 {
	checkSameSize(high, low, close)

	ad := make([]float64, len(close))

	for i := 0; i < len(ad); i++ {
		if i > 0 {
			ad[i] = ad[i-1]
		}

		ad[i] += float64(volume[i]) * (((close[i] - low[i]) - (high[i] - close[i])) / (high[i] - low[i]))
	}

	return ad
}
