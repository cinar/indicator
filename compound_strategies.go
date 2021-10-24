// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// MACD and RSI strategy.
func MacdAndRsiStrategy(asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	macdActions := MacdStrategy(asset)
	rsiActions := DefaultRsiStrategy(asset)

	for i := 0; i < len(actions); i++ {
		if macdActions[i] == rsiActions[i] {
			actions[i] = macdActions[i]
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}
