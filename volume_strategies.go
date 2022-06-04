// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// Money flow index strategy.
func MoneyFlowIndexStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	moneyFlowIndex := DefaultMoneyFlowIndex(
		asset.High,
		asset.Low,
		asset.Closing,
		asset.Volume)

	for i := 0; i < len(actions); i++ {
		if moneyFlowIndex[i] >= 80 {
			actions[i] = SELL
		} else {
			actions[i] = BUY
		}
	}

	return actions
}

// Force index strategy function.
func ForceIndexStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	forceIndex := DefaultForceIndex(asset.Closing, asset.Volume)

	for i := 0; i < len(actions); i++ {
		if forceIndex[i] > 0 {
			actions[i] = BUY
		} else if forceIndex[i] < 0 {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Ease of movement strategy.
func EaseOfMovementStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	emv := DefaultEaseOfMovement(asset.High, asset.Low, asset.Volume)

	for i := 0; i < len(actions); i++ {
		if emv[i] > 0 {
			actions[i] = BUY
		} else if emv[i] < 0 {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Volume weighted average price strategy function.
func VolumeWeightedAveragePriceStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	vwap := DefaultVolumeWeightedAveragePrice(asset.Closing, asset.Volume)

	for i := 0; i < len(actions); i++ {
		if vwap[i] > asset.Closing[i] {
			actions[i] = BUY
		} else if vwap[i] < asset.Closing[i] {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Negative volume index strategy.
func NegativeVolumeIndexStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	nvi := NegativeVolumeIndex(asset.Closing, asset.Volume)
	nvi255 := Ema(255, nvi)

	for i := 0; i < len(actions); i++ {
		if nvi[i] < nvi255[i] {
			actions[i] = BUY
		} else if nvi[i] > nvi255[i] {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Chaikin money flow strategy.
func ChaikinMoneyFlowStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	cmf := ChaikinMoneyFlow(
		asset.High,
		asset.Low,
		asset.Closing,
		asset.Volume)

	for i := 0; i < len(actions); i++ {
		if cmf[i] < 0 {
			actions[i] = BUY
		} else if cmf[i] > 0 {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}
