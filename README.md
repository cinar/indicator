[![GoDoc](https://godoc.org/github.com/cinar/indicator?status.svg)](https://godoc.org/github.com/cinar/indicator)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.com/cinar/indicator.svg?branch=master)](https://travis-ci.com/cinar/indicator)

# Indicator Go

Indicator is a Golang module providing various stock technical analysis indicators, strategies, and a backtest framework for trading.

## Indicators Provided

The following list of indicators are currently supported by this package:

### Trend Indicators

- [Absolute Price Oscillator (APO)](trend_indicators.md#absolute-price-oscillator-apo)
- [Aroon Indicator](trend_indicators.md#aroon-indicator)
- [Balance of Power (BOP)](trend_indicators.md#balance-of-power-bop)
- [Chande Forecast Oscillator (CFO)](trend_indicators.md#chande-forecast-oscillator-cfo)
- [Double Exponential Moving Average (DEMA)](trend_indicators.md#double-exponential-moving-average-dema)
- [Exponential Moving Average (EMA)](trend_indicators.md#exponential-moving-average-ema)
- [Moving Average Convergence Divergence (MACD)](trend_indicators.md#moving-average-convergence-divergence-macd)
- [Moving Max](trend_indicators.md#moving-max)
- [Moving Min](trend_indicators.md#moving-min)
- [Moving Sum](trend_indicators.md#moving-sum)
- [Parabolic SAR](trend_indicators.md#parabolic-sar)
- [Qstick](trend_indicators.md#qstick)
- [Random Index (KDJ)](trend_indicators.md#random-index-kdj)
- [Simple Moving Average (SMA)](trend_indicators.md#simple-moving-average-sma)
- [Since Change](trend_indicators.md#since-change)
- [Triangular Moving Average (TRIMA)](trend_indicators.md#triangular-moving-average-trima)
- [Triple Exponential Moving Average (TEMA)](trend_indicators.md#triple-exponential-moving-average-tema)
- [Typical Price](trend_indicators.md#typical-price)
- [Vortex Indicator](trend_indicators.md#vortex-indicator)

### Momentum Indicators

- [Awesome Oscillator](momentum_indicators.md#awesome-oscillator)
- [Chaikin Oscillator](momentum_indicators.md#chaikin-oscillator)
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

- [Asset](strategy.md#asset)
- [Action](strategy.md#action)
- [Strategy Function](strategy.md#strategy-function)
- [Buy and Hold Strategy](strategy.md#buy-and-hold-strategy)

### Trend Strategies

- [Chande Forecast Oscillator Strategy](trend_strategies.md#chande-forecast-oscillator-strategy)
- [KDJ Strategy](trend_strategies.md#kdj-strategy)
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

- No Strategies.

### Compound Strategies

- [All Strategies](compound_strategies.md#all-strategies)
- [Run Strategies](compound_strategies.md#run-strategies)
- [Separate Strategies](compound_strategies.md#separate-strategies)
- [MACD and RSI Strategy](compound_strategies.md#macd-and-rsi-strategy)

## Backtest

Backtesting is the method for seeing how well a strategy would have done. The following backtesting functions are provided for evaluating strategies.

- [Apply Actions](backtest.md#apply-actions)
- [Count Transactions](backtest.md#count-transactions)
- [Normalize Actions](backtest.md#normalize-actions)
- [Normalize Gains](backtest.md#normalize-gains)

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

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
