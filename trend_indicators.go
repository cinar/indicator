// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"math"

	"github.com/cinar/indicator/container/bst"
)

// Trend indicator.
type Trend int

const (
	Falling Trend = -1
	Rising  Trend = 1
)

const (
	psarAfStep = 0.02
	psarAfMax  = 0.20
)

// The AbsolutePriceOscillator function calculates a technical indicator that is used
// to follow trends. APO crossing above zero indicates bullish, while crossing below
// zero indicates bearish. Positive value is upward trend, while negative value is
// downward trend.
//
// Fast = Ema(fastPeriod, values)
// Slow = Ema(slowPeriod, values)
// APO = Fast - Slow
//
// Returns apo.
func AbsolutePriceOscillator(fastPeriod, slowPeriod int, values []float64) []float64 {
	fast := Ema(fastPeriod, values)
	slow := Ema(slowPeriod, values)
	apo := subtract(fast, slow)

	return apo
}

// The DefaultAbsolutePriceOscillator function calculates APO with the most
// frequently used fast and short periods are 14 and 30.
//
// Returns apo.
func DefaultAbsolutePriceOscillator(values []float64) []float64 {
	return AbsolutePriceOscillator(14, 30, values)
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

// The BalanceOfPower function calculates the strength of buying and selling
// pressure. Positive value indicates an upward trend, and negative value
// indicates a downward trend. Zero indicates a balance between the two.
//
// BOP = (Closing - Opening) / (High - Low)
//
// Returns bop.
func BalanceOfPower(opening, high, low, closing []float64) []float64 {
	bop := divide(subtract(closing, opening), subtract(high, low))
	return bop
}

// The Chande Forecast Oscillator developed by Tushar Chande The Forecast
// Oscillator plots the percentage difference between the closing price and
// the n-period linear regression forecasted price. The oscillator is above
// zero when the forecast price is greater than the closing price and less
// than zero if it is below.
//
// R = Linreg(Closing)
// CFO = ((Closing - R) / Closing) * 100
//
// Returns cfo.
func ChandeForecastOscillator(closing []float64) []float64 {
	x := generateNumbers(0, float64(len(closing)), 1)
	r := LinearRegressionUsingLeastSquare(x, closing)

	cfo := multiplyBy(divide(subtract(closing, r), closing), 100)

	return cfo
}

// The Community Channel Index (CMI) is a momentum-based oscillator
// used to help determine when an investment vehicle is reaching a
// condition of being overbought or oversold.
//
// Moving Average = Sma(Period, Typical Price)
// Mean Deviation = Sma(Period, Abs(Typical Price - Moving Average))
// CMI = (Typical Price - Moving Average) / (0.015 * Mean Deviation)
//
// Returns cmi.
func CommunityChannelIndex(period int, high, low, closing []float64) []float64 {
	tp, _ := TypicalPrice(low, high, closing)
	ma := Sma(period, tp)
	md := Sma(period, abs(subtract(tp, ma)))
	cci := divide(subtract(tp, ma), multiplyBy(md, 0.015))
	cci[0] = 0

	return cci
}

// The default community channel index with the period of 20.
func DefaultCommunityChannelIndex(high, low, closing []float64) []float64 {
	return CommunityChannelIndex(20, high, low, closing)
}

// Dema calculates the Double Exponential Moving Average (DEMA).
//
// DEMA = (2 * EMA(values)) - EMA(EMA(values))
//
// Returns dema.
func Dema(period int, values []float64) []float64 {
	ema1 := Ema(period, values)
	ema2 := Ema(period, ema1)

	dema := subtract(multiplyBy(ema1, 2), ema2)

	return dema
}

// Exponential Moving Average (EMA).
func Ema(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	k := float64(2) / float64(1+period)

	for i, value := range values {
		if i > 0 {
			result[i] = (value * k) + (result[i-1] * float64(1-k))
		} else {
			result[i] = value
		}
	}

	return result
}

// Moving Average Convergence Divergence (MACD).
//
// MACD = 12-Period EMA - 26-Period EMA.
// Signal = 9-Period EMA of MACD.
//
// Returns MACD, signal.
func Macd(closing []float64) ([]float64, []float64) {
	ema12 := Ema(12, closing)
	ema26 := Ema(26, closing)
	macd := subtract(ema12, ema26)
	signal := Ema(9, macd)

	return macd, signal
}

// The Mass Index (MI) uses the high-low range to identify trend reversals
// based on range expansions.
//
// Singe EMA = EMA(9, Highs - Lows)
// Double EMA = EMA(9, Single EMA)
// Ratio = Single EMA / Double EMA
// MI = Sum(25, Ratio)
//
// Returns mi.
func MassIndex(high, low []float64) []float64 {
	ema1 := Ema(9, subtract(high, low))
	ema2 := Ema(9, ema1)
	ratio := divide(ema1, ema2)
	mi := Sum(25, ratio)

	return mi
}

// Moving Chande Forecast Oscillator calculates based on
// the given period.
//
// The Chande Forecast Oscillator developed by Tushar Chande The Forecast
// Oscillator plots the percentage difference between the closing price and
// the n-period linear regression forecasted price. The oscillator is above
// zero when the forecast price is greater than the closing price and less
// than zero if it is below.
//
// R = Linreg(Closing)
// CFO = ((Closing - R) / Closing) * 100
//
// Returns cfo.
func MovingChandeForecastOscillator(period int, closing []float64) []float64 {
	x := generateNumbers(0, float64(len(closing)), 1)
	r := MovingLinearRegressionUsingLeastSquare(period, x, closing)

	cfo := multiplyBy(divide(subtract(closing, r), closing), 100)

	return cfo
}

// Moving max for the given period.
func Max(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)
	bst := bst.New()

	for i := 0; i < len(values); i++ {
		bst.Insert(values[i])

		if i >= period {
			bst.Remove(buffer[i%period])
		}

		buffer[i%period] = values[i]
		result[i] = bst.Max().(float64)
	}

	return result
}

// Moving min for the given period.
func Min(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)
	bst := bst.New()

	for i := 0; i < len(values); i++ {
		bst.Insert(values[i])

		if i >= period {
			bst.Remove(buffer[i%period])
		}

		buffer[i%period] = values[i]
		result[i] = bst.Min().(float64)
	}

	return result
}

// Parabolic SAR. It is a popular technical indicator for identifying the trend
// and as a trailing stop.
//
// PSAR = PSAR[i - 1] - ((PSAR[i - 1] - EP) * AF)
//
// If the trend is Falling:
//   - PSAR is the maximum of PSAR or the previous two high values.
//   - If the current high is greather than or equals to PSAR, use EP.
//
// If the trend is Rising:
//   - PSAR is the minimum of PSAR or the previous two low values.
//   - If the current low is less than or equals to PSAR, use EP.
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
func ParabolicSar(high, low, closing []float64) ([]float64, []Trend) {
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

		if psar[i] > closing[i] {
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

// The Qstick function calculates the ratio of recent up and down bars.
//
// QS = Sma(Closing - Opening)
//
// Returns qs.
func Qstick(period int, opening, closing []float64) []float64 {
	qs := Sma(period, subtract(closing, opening))
	return qs
}

// The Kdj function calculates the KDJ  indicator, also known as
// the Random Index. KDJ is calculated similar to the Stochastic
// Oscillator with the difference of having the J line. It is
// used to analyze the trend and entry points.
//
// The K and D lines show if the asset is overbought when they
// crosses above 80%, and oversold when they crosses below
// 20%. The J line represents the divergence.
//
// RSV = ((Closing - Min(Low, rPeriod))
//
//	/ (Max(High, rPeriod) - Min(Low, rPeriod))) * 100
//
// K = Sma(RSV, kPeriod)
// D = Sma(K, dPeriod)
// J = (3 * K) - (2 * D)
//
// Returns k, d, j.
func Kdj(rPeriod, kPeriod, dPeriod int, high, low, closing []float64) ([]float64, []float64, []float64) {
	highest := Max(rPeriod, high)
	lowest := Min(rPeriod, low)

	rsv := multiplyBy(divide(subtract(closing, lowest), subtract(highest, lowest)), 100)

	k := Sma(kPeriod, rsv)
	d := Sma(dPeriod, k)
	j := subtract(multiplyBy(k, 3), multiplyBy(d, 2))

	return k, d, j
}

// The DefaultKdj function calculates KDJ based on default periods
// consisting of rPeriod of 9, kPeriod of 3, and dPeriod of 3.
//
// Returns k, d, j.
func DefaultKdj(high, low, closing []float64) ([]float64, []float64, []float64) {
	return Kdj(9, 3, 3, high, low, closing)
}

// Rolling Moving Average (RMA).
//
// R[0] to R[p-1] is SMA(values)
// R[p] and after is R[i] = ((R[i-1]*(p-1)) + v[i]) / p
//
// Returns r.
func Rma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1

		if i < period {
			sum += value
		} else {
			sum = (result[i-1] * float64(period-1)) + value
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}

// Simple Moving Average (SMA).
func Sma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1
		sum += value

		if i >= period {
			sum -= values[i-period]
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}

// Since last values change.
func Since(values []float64) []int {
	result := make([]int, len(values))

	lastValue := math.NaN()
	sinceLast := 0

	for i := 0; i < len(values); i++ {
		value := values[i]

		if value != lastValue {
			lastValue = value
			sinceLast = 0
		} else {
			sinceLast++
		}

		result[i] = sinceLast
	}

	return result
}

// Moving sum for the given period.
func Sum(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := 0.0

	for i := 0; i < len(values); i++ {
		sum += values[i]

		if i >= period {
			sum -= values[i-period]
		}

		result[i] = sum
	}

	return result
}

// Tema calculates the Triple Exponential Moving Average (TEMA).
//
// TEMA = (3 * EMA1) - (3 * EMA2) + EMA3
// EMA1 = EMA(values)
// EMA2 = EMA(EMA1)
// EMA3 = EMA(EMA2)
//
// Returns tema.
func Tema(period int, values []float64) []float64 {
	ema1 := Ema(period, values)
	ema2 := Ema(period, ema1)
	ema3 := Ema(period, ema2)

	tema := add(subtract(multiplyBy(ema1, 3), multiplyBy(ema2, 3)), ema3)

	return tema
}

// Trima function calculates the Triangular Moving Average (TRIMA).
//
// If period is even:
//
//	TRIMA = SMA(period / 2, SMA((period / 2) + 1, values))
//
// If period is odd:
//
//	TRIMA = SMA((period + 1) / 2, SMA((period + 1) / 2, values))
//
// Returns trima.
func Trima(period int, values []float64) []float64 {
	var n1, n2 int

	if period%2 == 0 {
		n1 = period / 2
		n2 = n1 + 1
	} else {
		n1 = (period + 1) / 2
		n2 = n1
	}

	trima := Sma(n1, Sma(n2, values))

	return trima
}

// Triple Exponential Average (TRIX) indicator is an oscillator used to
// identify oversold and overbought markets, and it can also be used
// as a momentum indicator. Like many oscillators, TRIX oscillates
// around a zero line.
//
// EMA1 = EMA(period, values)
// EMA2 = EMA(period, EMA1)
// EMA3 = EMA(period, EMA2)
// TRIX = (EMA3 - Previous EMA3) / Previous EMA3
//
// Returns trix.
func Trix(period int, values []float64) []float64 {
	ema1 := Ema(period, values)
	ema2 := Ema(period, ema1)
	ema3 := Ema(period, ema2)
	previous := shiftRightAndFillBy(1, ema3[0], ema3)
	trix := divide(subtract(ema3, previous), previous)

	return trix
}

// Typical Price. It is another approximation of average price for each
// period and can be used as a filter for moving average systems.
//
// Typical Price = (High + Low + Closing) / 3
//
// Returns typical price, 20-Period SMA.
func TypicalPrice(low, high, closing []float64) ([]float64, []float64) {
	checkSameSize(high, low, closing)

	sma20 := Sma(20, closing)

	ta := make([]float64, len(closing))
	for i := 0; i < len(ta); i++ {
		ta[i] = (high[i] + low[i] + closing[i]) / float64(3)
	}

	return ta, sma20
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
// TR = Max((High[i]-Low[i]), Abs(High[i]-Closing[i-1]), Abs(Low[i]-Closing[i-1]))
// TR14 = 14-Period Sum of TR
//
// +VI14 = +VM14 / TR14
// -VI14 = -VM14 / TR14
//
// Based on https://school.stockcharts.com/doku.php?id=technical_indicators:vortex_indicator
//
// Returns plusVi, minusVi
func Vortex(high, low, closing []float64) ([]float64, []float64) {
	checkSameSize(high, low, closing)

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
		highPrevClosing := math.Abs(high[i] - closing[i-1])
		lowPrevClosing := math.Abs(low[i] - closing[i-1])

		trSum -= tr[j]
		tr[j] = math.Max(highLow, math.Max(highPrevClosing, lowPrevClosing))
		trSum += tr[j]

		plusVi[i] = plusVmSum / trSum
		minusVi[i] = minusVmSum / trSum
	}

	return plusVi, minusVi
}

// The Vwma function calculates the Volume Weighted Moving Average (VWMA)
// averaging the price data with an emphasis on volume, meaning areas
// with higher volume will have a greater weight.
//
// VWMA = Sum(Price * Volume) / Sum(Volume) for a given Period.
//
// Returns vwma
func Vwma(period int, closing, volume []float64) []float64 {
	vwma := divide(Sum(period, multiply(closing, volume)), Sum(period, volume))

	return vwma
}

// The DefaultVwma function calculates VWMA with a period of 20.
func DefaultVwma(closing, volume []float64) []float64 {
	return Vwma(20, closing, volume)
}
