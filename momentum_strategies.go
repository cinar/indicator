// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// Awesome oscillator strategy function.
func AwesomeOscillatorStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	ao := AwesomeOscillator(asset.Low, asset.High)

	for i := 0; i < len(actions); i++ {
		if ao[i] > 0 {
			actions[i] = BUY
		} else if ao[i] < 0 {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// RSI strategy. Sells above sell at, buys below buy at.
func RsiStrategy(asset *Asset, sellAt, buyAt float64) []Action {
	actions := make([]Action, len(asset.Date))

	_, rsi := Rsi(asset.Closing)

	for i := 0; i < len(actions); i++ {
		if rsi[i] <= buyAt {
			actions[i] = BUY
		} else if rsi[i] >= sellAt {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Default RSI strategy function. It buys
// below 30 and sells above 70.
func DefaultRsiStrategy(asset *Asset) []Action {
	return RsiStrategy(asset, 70, 30)
}

// Make RSI strategy function.
func MakeRsiStrategy(sellAt, buyAt float64) StrategyFunction {
	return func(asset *Asset) []Action {
		return RsiStrategy(asset, sellAt, buyAt)
	}
}

// RSI 2 strategy. When 2-period RSI moves below 10, it is considered deeply oversold,
// and the other way around when moves above 90.
func Rsi2Strategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	_, rsi := Rsi2(asset.Closing)

	for i := 0; i < len(actions); i++ {
		if rsi[i] < 10 {
			actions[i] = BUY
		} else if rsi[i] > 90 {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Williams R strategy function.
func WilliamsRStrategy(asset *Asset) []Action {
	actions := make([]Action, len(asset.Date))

	wr := WilliamsR(asset.Low, asset.High, asset.Closing)

	for i := 0; i < len(actions); i++ {
		if wr[i] < -20 {
			actions[i] = SELL
		} else if wr[i] > -80 {
			actions[i] = BUY
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}
