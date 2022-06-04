// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// Bollinger bands strategy function.
func BollingerBandsStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	_, upperBand, lowerBand := BollingerBands(asset.Closing)

	for i := 0; i < len(actions); i++ {
		if asset.Closing[i] > upperBand[i] {
			actions[i] = SELL
		} else if asset.Closing[i] < lowerBand[i] {
			actions[i] = BUY
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Projection oscillator strategy function.
func ProjectionOscillatorStrategy(period, smooth int, asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	po, spo := ProjectionOscillator(
		period,
		smooth,
		asset.High,
		asset.Low,
		asset.Closing)

	for i := 0; i < len(actions); i++ {
		if po[i] > spo[i] {
			actions[i] = BUY
		} else if po[i] < spo[i] {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Make projection oscillator strategy.
func MakeProjectionOscillatorStrategy(period, smooth int) StrategyFunction {
	return func(asset *Asset) []Action {
		return ProjectionOscillatorStrategy(period, smooth, asset)
	}
}
