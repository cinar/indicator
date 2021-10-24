// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// Chande forecast oscillator strategy.
func ChandeForecastOscillatorStrategy(asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	cfo := ChandeForecastOscillator(asset.Closing)

	for i := 0; i < len(actions); i++ {
		if cfo[i] < 0 {
			actions[i] = BUY
		} else if cfo[i] > 0 {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Moving chande forecast oscillator strategy function.
func MovingChandeForecastOscillatorStrategy(period int, asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	cfo := MovingChandeForecastOscillator(period, asset.Closing)

	for i := 0; i < len(actions); i++ {
		if cfo[i] < 0 {
			actions[i] = BUY
		} else if cfo[i] > 0 {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Make moving chande forecast oscillator strategy.
func MakeMovingChandeForecastOscillatorStrategy(period int) StrategyFunction {
	return func(asset Asset) []Action {
		return MovingChandeForecastOscillatorStrategy(period, asset)
	}
}

// MACD strategy.
func MacdStrategy(asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	macd, signal := Macd(asset.Closing)

	for i := 0; i < len(actions); i++ {
		if macd[i] > signal[i] {
			actions[i] = BUY
		} else if macd[i] < signal[i] {
			actions[i] = SELL
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Trend strategy. Buy when trending up for count times,
// sell when trending down for count times.
func TrendStrategy(asset Asset, count uint) []Action {
	actions := make([]Action, len(asset.Date))

	if len(actions) == 0 {
		return actions
	}

	lastClosing := asset.Closing[0]
	trendCount := uint(1)
	trendUp := false

	actions[0] = HOLD

	for i := 1; i < len(actions); i++ {
		closing := asset.Closing[i]

		if trendUp && (lastClosing <= closing) {
			trendCount++
		} else if !trendUp && (lastClosing >= closing) {
			trendCount++
		} else {
			trendUp = !trendUp
			trendCount = 1
		}

		lastClosing = closing

		if trendCount >= count {
			if trendUp {
				actions[i] = BUY
			} else {
				actions[i] = SELL
			}
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Make trend strategy function.
func MakeTrendStrategy(count uint) StrategyFunction {
	return func(asset Asset) []Action {
		return TrendStrategy(asset, count)
	}
}
