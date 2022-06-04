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
	return func(asset *Asset) []Action {
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

// The SeparateStrategies function takes a buy strategy and a sell strategy.
//
// It returns a BUY action if the buy strategy returns a BUY action and
// the the sell strategy returns a HOLD action.
//
// It returns a SELL action if the sell strategy returns a SELL action
// and the buy strategy returns a HOLD action.
//
// It returns HOLD otherwise.
func SeparateStategies(buyStrategy, sellStrategy StrategyFunction) StrategyFunction {
	return func(asset *Asset) []Action {
		actions := make([]Action, len(asset.Date))

		buyActions := buyStrategy(asset)
		sellActions := sellStrategy(asset)

		for i := 0; i < len(actions); i++ {
			if buyActions[i] == BUY && sellActions[i] == HOLD {
				actions[i] = BUY
			} else if sellActions[i] == SELL && buyActions[i] == HOLD {
				actions[i] = SELL
			} else {
				actions[i] = HOLD
			}
		}

		return actions
	}
}

// The RunStrategies takes one or more StrategyFunction and returns
// the acitions for each.
func RunStrategies(asset *Asset, strategies ...StrategyFunction) [][]Action {
	actions := make([][]Action, len(strategies))

	for i := 0; i < len(strategies); i++ {
		actions[i] = strategies[i](asset)
	}

	return actions
}

// MACD and RSI strategy.
func MacdAndRsiStrategy(asset *Asset) []Action {
	strategy := AllStrategies(MacdStrategy, DefaultRsiStrategy)
	return strategy(asset)
}
