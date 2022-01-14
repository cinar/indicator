// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// Accumulation/Distribution Indicator (A/D). Cumulative indicator
// that uses volume and price to assess whether a stock is
// being accumulated or distributed.
//
// MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
// MFV = MFM * Period Volume
// AD = Previous AD + CMFV
//
// Returns ad.
func AccumulationDistribution(high, low, closing []float64, volume []int64) []float64 {
	checkSameSize(high, low, closing)

	ad := make([]float64, len(closing))

	for i := 0; i < len(ad); i++ {
		if i > 0 {
			ad[i] = ad[i-1]
		}

		ad[i] += float64(volume[i]) * (((closing[i] - low[i]) - (high[i] - closing[i])) / (high[i] - low[i]))
	}

	return ad
}

// On-Balance Volume (OBV). It is a technical trading momentum indicator that
// uses volume flow to predict changes in stock price.
//
//                   volume, if Closing > Closing-Prev
// OBV = OBV-Prev +       0, if Closing = Closing-Prev
//                  -volume, if Closing < Closing-Prev
//
// Returns obv
func Obv(closing []float64, volume []int64) []int64 {
	if len(closing) != len(volume) {
		panic("not all same size")
	}

	obv := make([]int64, len(volume))

	for i := 1; i < len(obv); i++ {
		obv[i] = obv[i-1]

		if closing[i] > closing[i-1] {
			obv[i] += volume[i]
		} else if closing[i] < closing[i-1] {
			obv[i] -= volume[i]
		}
	}

	return obv
}

// The Money Flow Index (MFI) analyzes both the closing price and the volume
// to measure to identify overbought and oversold states. It is similar to
// the Relative Strength Index (RSI), but it also uses the volume.
//
// Raw Money Flow = Typical Price * Volume
// Money Ratio = Positive Money Flow / Negative Money Flow
// Money Flow Index = 100 - (100 / (1 + Money Ratio))
//
// Retruns money flow index values.
func MoneyFlowIndex(period int, high, low, closing []float64, volume []int64) []float64 {
	typicalPrice, _ := TypicalPrice(low, high, closing)
	rawMoneyFlow := multiply(typicalPrice, asFloat64(volume))

	signs := extractSign(diff(rawMoneyFlow, 1))
	moneyFlow := multiply(signs, rawMoneyFlow)

	positiveMoneyFlow := keepPositives(moneyFlow)
	negativeMoneyFlow := keepNegatives(moneyFlow)

	moneyRatio := divide(
		Sum(period, positiveMoneyFlow),
		Sum(period, negativeMoneyFlow))

	moneyFlowIndex := addBy(multiplyBy(pow(addBy(moneyRatio, 1), -1), -100), 100)

	return moneyFlowIndex
}

// Default money flow index with period 14.
func DefaultMoneyFlowIndex(high, low, closing []float64, volume []int64) []float64 {
	return MoneyFlowIndex(14, high, low, closing, volume)
}

// The Force Index (FI) uses the closing price and the volume to assess
// the power behind a move and identify turning points.
//
// Force Index = EMA(period, (Current - Previous) * Volume)
//
// Returns force index.
func ForceIndex(period int, closing []float64, volume []int64) []float64 {
	return Ema(period, multiply(diff(closing, 1), asFloat64(volume)))
}

// The default Force Index (FI) with window size of 13.
func DefaultForceIndex(closing []float64, volume []int64) []float64 {
	return ForceIndex(13, closing, volume)
}

// The Ease of Movement (EMV) is a volume based oscillator measuring
// the ease of price movement.
//
// Distance Moved = ((High + Low) / 2) - ((Priod High + Prior Low) /2)
// Box Ratio = ((Volume / 100000000) / (High - Low))
// EMV(1) = Distance Moved / Box Ratio
// EMV(14) = SMA(14, EMV(1))
//
// Returns ease of movement values.
func EaseOfMovement(period int, high, low []float64, volume []int64) []float64 {
	distanceMoved := diff(divideBy(add(high, low), 2), 1)
	boxRatio := divide(divideBy(asFloat64(volume), float64(100000000)), substract(high, low))
	emv := Sma(period, divide(distanceMoved, boxRatio))
	return emv
}

// The default Ease of Movement with the default period of 14.
func DefaultEaseOfMovement(high, low []float64, volume []int64) []float64 {
	return EaseOfMovement(14, high, low, volume)
}
