[![GoDoc](https://godoc.org/github.com/cinar/indicator?status.svg)](https://godoc.org/github.com/cinar/indicator)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.com/cinar/indicator.svg?branch=master)](https://travis-ci.com/cinar/indicator)

# Indicator Go

Indicator is a Golang module providing various stock technical analysis indicators, and strategies for trading.

## Indicators Provided

The following list of indicators are currently supported by this package:

### Trend Indicators

- [Aroon Indicator](trend_indicators.md#aroon-indicator)
- [Chande Forecast Oscillator (CFO)](trend_indicators.md#chande-forecast-oscillator-cfo)
- [Double Exponential Moving Average (DEMA)](trend_indicators.md#double-exponential-moving-average-dema)
- [Exponential Moving Average (EMA)](trend_indicators.md#exponential-moving-average-ema)
- [Moving Average Convergence Divergence (MACD)](trend_indicators.md#moving-average-convergence-divergence-macd)
- [Moving Max](trend_indicators.md#moving-max)
- [Moving Min](trend_indicators.md#moving-min)
- [Parabolic SAR](#parabolic-sar)
- [Simple Moving Average (SMA)](trend_indicators.md#simple-moving-average-sma)
- [Since Change](trend_indicators.md#since-change)
- [Triangular Moving Average (TRIMA)](trend_indicators.md#triangular-moving-average-trima)
- [Triple Exponential Moving Average (TEMA)](trend_indicators.md#triple-exponential-moving-average-tema)
- [Typical Price](trend_indicators.md#typical-price)
- [Vortex Indicator](trend_indicators.md#vortex-indicator)

### Momentum Indicators

- [Awesome Oscillator](momentum_indicators.md#awesome-oscillator)
- [Ichimoku Cloud](momentum_indicators.md#ichimoku-cloud)
- [Relative Strength Index (RSI)](momentum_indicators.md#relative-strength-index-rsi)
- [Stochastic Oscillator](momentum_indicators.md#stochastic-oscillator)
- [Williams R](momentum_indicators.md#williams-r)

### Volatility Indicators

- [Acceleration Bands](volatility_indicators.md#acceleration-bands)
- [Actual True Range (ATR)](volatility_indicators.md#actual-true-range-atr)
- [Bollinger Band Width](volatility_indicators.md#bollinger-band-width)
- [Bollinger Bands](volatility_indicators.md#bollinger-bands)
- [Chandelier Exit](volatility_indicators.md#chandelier-exit)
- [Moving Standard Deviation (Std)](volatility_indicators.md#moving-standard-deviation-std)
- [Projection Oscillator (PO)](volatility_indicators.md#projection-oscillator-po)

### Volume Indicators

- [Accumulation/Distribution (A/D)](volume_indicators.md#accumulationdistribution-ad)
- [On-Balance Volume (OBV)](volume_indicators.md#on-balance-volume-obv)

## Strategies Provided

The following list of strategies are currently supported by this package:

### Reference Strategies

- [Buy and Hold Strategy](#buy-and-hold-strategy)

### Trend Strategies

- [Chande Forecast Oscillator Strategy](trend_strategies.md#chande-forecast-oscillator-strategy)
- [MACD Strategy](trend_strategies.md#macd-strategy)
- [Trend Strategy](trend_strategies.md#trend-strategy)

### Momentum Strategies

- [Awesome Oscillator Strategy](momentum_strategies.md#awesome-oscillator-strategy)
- [RSI Strategy](momentum_strategies.md#rsi-strategy)
- [Williams R Strategy](momentum_strategies.md#williams-r-strategy)

### Volatility Strategies

- [Bollinger Bands Strategy](volatility_strategies.md#bollinger-bands-strategy)
- [Projection Oscillator Strategy](volatility_strategies.md#projection-oscillator-strategy)

### Volume Strategies

### Compound Strategies

- [MACD and RSI Strategy](compound_strategies.md#macd-and-rsi-strategy)

## Usage

Install package.

```bash
go get github.com/cinar/indicator
```

Import indicator.

```Golang
import (
    "github.com/cinar/indicator"
)
```

### Strategies

The strategies are where the results from one or more indicators gets combined to produce a recommended action.

**The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.**

The stragies operates on an [Asset](https://pkg.go.dev/github.com/cinar/indicator#Asset) with the following members.

```golang
type Asset struct {
	Date    []time.Time
	Opening []float64
	Closing []float64
	High    []float64
	Low     []float64
	Volume  []int64
}
```

The [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) takes an [Asset](https://pkg.go.dev/github.com/cinar/indicator#Asset), and provides an array of [Action](https://pkg.go.dev/github.com/cinar/indicator#Action) for each row.

```golang
// Strategy function.
type StrategyFunction func(Asset) []Action
```

The following [Action](https://pkg.go.dev/github.com/cinar/indicator#Action) values are currently provided.

```golang
type Action int

const (
	SELL Action = -1
	HOLD Action = 0
	BUY  Action = 1
)
```

### Reference Strategies

Reference strategies provide a base for comparing the actual outcome of a given strategy.

#### Buy and Hold Strategy

The [BuyAndHoldStrategy](https://pkg.go.dev/github.com/cinar/indicator#BuyAndHoldStrategy) provides a simple strategy to buy the given asset and hold it. It provides a good indicator for the change of asset's value without any other strategy is used.

```golang
actions := indicator.BuyAndHoldStrategy(asset)
```

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
