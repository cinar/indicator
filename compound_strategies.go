// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// The AllStrategies function takes one or more StrategyFunction and
// provides a StrategyFunction that will return a BUY or SELL action
// if all strategies are returning the same action, otherwise it
// will return a HOLD action.
func AllStrategies(strategies ...StrategyFunction) StrategyFunction {
	return func(asset Asset) []Action {
		actions := RunStrategies(asset, strategies...)

		for i := 1; i < len(actions); i++ {
			for j := 0; j < len(actions[0]); j++ {
				if actions[0][j] != actions[i][j] {
					actions[0][j] = HOLD
				}
			}
		}

		return actions[0]
	}
}

// MACD and RSI strategy.
func MacdAndRsiStrategy(asset Asset) []Action {
	strategy := AllStrategies(MacdStrategy, DefaultRsiStrategy)
	return strategy(asset)
}

// The RunStrategies takes one or more StrategyFunction and returns
// the acitions for each.
func RunStrategies(asset Asset, strategies ...StrategyFunction) [][]Action {
	actions := make([][]Action, len(strategies))

	for i := 0; i < len(strategies); i++ {
		actions[i] = strategies[i](asset)
	}

	return actions
}
