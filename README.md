[![GoDoc](https://godoc.org/github.com/cinar/indicator?status.svg)](https://godoc.org/github.com/cinar/indicator) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Build Status](https://github.com/cinar/indicator/actions/workflows/ci.yml/badge.svg)](https://github.com/cinar/indicator/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/cinar/indicator)](https://goreportcard.com/report/github.com/cinar/indicator) [![codecov](https://codecov.io/gh/cinar/indicator/branch/master/graph/badge.svg?token=MB7L69UAWM)](https://codecov.io/gh/cinar/indicator)

Indicator Go
============

Indicator is a Golang module providing various stock technical analysis indicators, strategies, and a backtest framework for trading.

> [!IMPORTANT] I have recently started working on the [v2 version](https://github.com/cinar/indicator/tree/v2). It is a complete rewrite of the library with the following goals:
>
> -	Achieving and maintaining minimum of 90% code coverage.
> -	Having test data in CSV format for each indicator and strategy for easyt validation.
> -	Operating on data streams (Go channels) for both inputs and outputs.
> -	Having each indicator and strategy fully configurable with no preset values.
> -	Supporting all numeric formats using Golang generics.
>
> [!NOTE] I also have a TypeScript version of this module now at [Indicator TS](https://github.com/cinar/indicatorts).

👆 Indicators Provided
----------------------

The following list of indicators are currently supported by this package:

### 📈 Trend Indicators

-	[Absolute Price Oscillator (APO)](trend_indicators.md#absolute-price-oscillator-apo)
-	[Aroon Indicator](trend_indicators.md#aroon-indicator)
-	[Balance of Power (BOP)](trend_indicators.md#balance-of-power-bop)
-	[Chande Forecast Oscillator (CFO)](trend_indicators.md#chande-forecast-oscillator-cfo)
-	[Community Channel Index (CMI)](trend_indicators.md#community-channel-index-cmi)
-	[Double Exponential Moving Average (DEMA)](trend_indicators.md#double-exponential-moving-average-dema)
-	[Exponential Moving Average (EMA)](trend_indicators.md#exponential-moving-average-ema)
-	[Mass Index (MI)](trend_indicators.md#mass-index-mi)
-	[Moving Average Convergence Divergence (MACD)](trend_indicators.md#moving-average-convergence-divergence-macd)
-	[Moving Max](trend_indicators.md#moving-max)
-	[Moving Min](trend_indicators.md#moving-min)
-	[Moving Sum](trend_indicators.md#moving-sum)
-	[Parabolic SAR](trend_indicators.md#parabolic-sar)
-	[Qstick](trend_indicators.md#qstick)
-	[Random Index (KDJ)](trend_indicators.md#random-index-kdj)
-	[Rolling Moving Average (RMA)](trend_indicators.md#rolling-moving-average-rma)
-	[Simple Moving Average (SMA)](trend_indicators.md#simple-moving-average-sma)
-	[Since Change](trend_indicators.md#since-change)
-	[Triple Exponential Moving Average (TEMA)](trend_indicators.md#triple-exponential-moving-average-tema)
-	[Triangular Moving Average (TRIMA)](trend_indicators.md#triangular-moving-average-trima)
-	[Triple Exponential Average (TRIX)](trend_indicators.md#triple-exponential-average-trix)
-	[Typical Price](trend_indicators.md#typical-price)
-	[Volume Weighted Moving Average (VWMA)](trend_indicators.md#volume-weighted-moving-average-vwma)
-	[Vortex Indicator](trend_indicators.md#vortex-indicator)

### 🚀 Momentum Indicators

-	[Awesome Oscillator](momentum_indicators.md#awesome-oscillator)
-	[Chaikin Oscillator](momentum_indicators.md#chaikin-oscillator)
-	[Ichimoku Cloud](momentum_indicators.md#ichimoku-cloud)
-	[Percentage Price Oscillator (PPO)](momentum_indicators.md#percentage-price-oscillator-ppo)
-	[Percentage Volume Oscillator (PVO)](momentum_indicators.md#percentage-volume-oscillator-pvo)
-	[Relative Strength Index (RSI)](momentum_indicators.md#relative-strength-index-rsi)
-	[RSI 2](momentum_indicators.md#rsi-2)
-	[RSI Period](momentum_indicators.md#rsi-period)
-	[Stochastic Oscillator](momentum_indicators.md#stochastic-oscillator)
-	[Williams R](momentum_indicators.md#williams-r)

### 🎢 Volatility Indicators

-	[Acceleration Bands](volatility_indicators.md#acceleration-bands)
-	[Actual True Range (ATR)](volatility_indicators.md#actual-true-range-atr)
-	[Bollinger Band Width](volatility_indicators.md#bollinger-band-width)
-	[Bollinger Bands](volatility_indicators.md#bollinger-bands)
-	[Chandelier Exit](volatility_indicators.md#chandelier-exit)
-	[Donchian Channel (DC)](volatility_indicators.md#donchian-channel-dc)
-	[Keltner Channel (KC)](volatility_indicators.md#keltner-channel-kc)
-	[Moving Standard Deviation (Std)](volatility_indicators.md#moving-standard-deviation-std)
-	[Projection Oscillator (PO)](volatility_indicators.md#projection-oscillator-po)
-	[Ulcer Index (UI)](volatility_indicators.md#ulcer-index-ui)

### 📢 Volume Indicators

-	[Accumulation/Distribution (A/D)](volume_indicators.md#accumulationdistribution-ad)
-	[Chaikin Money Flow (CMF)](volume_indicators.md#chaikin-money-flow-cmf)
-	[Ease of Movement (EMV)](volume_indicators.md#ease-of-movement-emv)
-	[Force Index (FI)](volume_indicators.md#force-index-fi)
-	[Money Flow Index (MFI)](volume_indicators.md#money-flow-index-mfi)
-	[Negative Volume Index (NVI)](volume_indicators.md#negative-volume-index-nvi)
-	[On-Balance Volume (OBV)](volume_indicators.md#on-balance-volume-obv)
-	[Volume Price Trend (VPT)](volume_indicators.md#volume-price-trend-vpt)
-	[Volume Weighted Average Price (VWAP)](volume_indicators.md#volume-weighted-average-price-vwap)

🧠 Strategies Provided
----------------------

Strategies relies on the following:

-	[Asset](strategy.md#asset)
-	[Action](strategy.md#action)
-	[Strategy Function](strategy.md#strategy-function)
-	[Buy and Hold Strategy](strategy.md#buy-and-hold-strategy)

The following list of strategies are currently supported by this package:

### 📈 Trend Strategies

-	[Chande Forecast Oscillator Strategy](trend_strategies.md#chande-forecast-oscillator-strategy)
-	[KDJ Strategy](trend_strategies.md#kdj-strategy)
-	[MACD Strategy](trend_strategies.md#macd-strategy)
-	[Trend Strategy](trend_strategies.md#trend-strategy)
-	[Volume Weighted Moving Average (VWMA) Strategy](trend_strategies.md#volume-weighted-moving-average-vwma-strategy)

### 🚀 Momentum Strategies

-	[Awesome Oscillator Strategy](momentum_strategies.md#awesome-oscillator-strategy)
-	[RSI Strategy](momentum_strategies.md#rsi-strategy)
-	[RSI 2 Strategy](momentum_strategies.md#rsi-2-strategy)
-	[Williams R Strategy](momentum_strategies.md#williams-r-strategy)

### 🎢 Volatility Strategies

-	[Bollinger Bands Strategy](volatility_strategies.md#bollinger-bands-strategy)
-	[Projection Oscillator Strategy](volatility_strategies.md#projection-oscillator-strategy)

### 📢 Volume Strategies

-	[Chaikin Money Flow Strategy](volume_strategies.md#chaikin-money-flow-strategy)
-	[Ease of Movement Strategy](volume_strategies.md#ease-of-movement-strategy)
-	[Force Index Strategy](volume_strategies.md#force-index-strategy)
-	[Money Flow Index Strategy](volume_strategies.md#money-flow-index-strategy)
-	[Negative Volume Index Strategy](volume_strategies.md#negative-volume-index-strategy)
-	[Volume Weighted Average Price Strategy](volume_strategies.md#volume-weighted-average-price-strategy)

### 🧪 Compound Strategies

-	[All Strategies](compound_strategies.md#all-strategies)
-	[Run Strategies](compound_strategies.md#run-strategies)
-	[Separate Strategies](compound_strategies.md#separate-strategies)
-	[MACD and RSI Strategy](compound_strategies.md#macd-and-rsi-strategy)

⏳ Backtesting
--------------

Backtesting is the method for seeing how well a strategy would have done. The following backtesting functions are provided for evaluating strategies.

-	[Apply Actions](backtest.md#apply-actions)
-	[Count Transactions](backtest.md#count-transactions)
-	[Normalize Actions](backtest.md#normalize-actions)
-	[Normalize Gains](backtest.md#normalize-gains)

Usage
-----

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

Disclaimer
----------

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

License
-------

The `v2.x.x` and above are provided under GNU AGPLv3 License.

```
Copyright (c) 2021-2024 Onur Cinar.    
The source code is provided under GNU AGPLv3 License.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```

The version `v1.x.x` is provided under MIT License.

```
Copyright (c) 2021-2024 Onur Cinar.
The source code is provided under MIT License.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
