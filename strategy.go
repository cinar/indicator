package indicator

import (
	"log"
	"time"
)

// Strategy action.
type Action int

const (
	HOLD Action = iota
	BUY
	SELL
)

// Asset values.
type Asset struct {
	Date   []time.Time
	Open   []float64
	Close  []float64
	High   []float64
	Low    []float64
	Volume []int64
}

// Strategy function. It takes an Asset and returns
// actions for each row.
type StrategyFunction func(Asset) []Action

// Buy and hold strategy. Buys at the beginning and holds.
func BuyAndHoldStrategy(asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	for i := 0; i < len(actions); i++ {
		actions[i] = BUY
	}

	return actions
}

// Trend strategy. Buy when trending up for count times,
// sell when trending down for count times.
func TrendStrategy(asset Asset, count int) []Action {
	if count < 1 {
		log.Fatal("count cannot be less than 1")
	}

	actions := make([]Action, len(asset.Date))

	if len(actions) == 0 {
		return actions
	}

	var lastClose float64 = asset.Close[0]
	var trendCount int
	var trendUp bool

	actions[0] = HOLD

	for i := 1; i < len(actions); i++ {
		close := asset.Close[i]

		if trendUp && (lastClose <= close) {
			trendCount++
		} else if !trendUp && (lastClose >= close) {
			trendCount++
		} else {
			trendUp = !trendUp
			trendCount = 1
		}

		lastClose = close

		if trendCount == count {
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
func MakeTrendStrategy(count int) StrategyFunction {
	return func(asset Asset) []Action {
		return TrendStrategy(asset, count)
	}
}

// MACD strategy.
func MacdStrategy(asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	macd, signal := Macd(asset.Close)

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

// RSI strategy. Sells above sell at, buys below buy at.
func RsiStrategy(asset Asset, sellAt, buyAt float64) []Action {
	actions := make([]Action, len(asset.Date))

	_, rsi := Rsi(asset.Close)

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
func DefaultRsiStrategy(asset Asset) []Action {
	return RsiStrategy(asset, 70, 30)
}

// Make RSI strategy function.
func MakeRsiStrategy(sellAt, buyAt float64) StrategyFunction {
	return func(asset Asset) []Action {
		return RsiStrategy(asset, sellAt, buyAt)
	}
}

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

// Bollinger bands strategy function.
func BollingerBandsStrategy(asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	_, upperBand, lowerBand := BollingerBands(asset.Close)

	for i := 0; i < len(actions); i++ {
		if asset.Close[i] > upperBand[i] {
			actions[i] = SELL
		} else if asset.Close[i] < lowerBand[i] {
			actions[i] = BUY
		} else {
			actions[i] = HOLD
		}
	}

	return actions
}

// Awesome oscillator strategy function.
func AwesomeOscillatorStrategy(asset Asset) []Action {
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

// Williams R strategy function.
func WilliamsRStrategy(asset Asset) []Action {
	actions := make([]Action, len(asset.Date))

	wr := WilliamsR(asset.Low, asset.High, asset.Close)

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
