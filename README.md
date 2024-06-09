[![GoDoc](https://godoc.org/github.com/cinar/indicator/v2?status.svg)](https://godoc.org/github.com/cinar/indicator/v2) [![License](https://img.shields.io/badge/License-AGPLv3-blue.svg)](https://opensource.org/licenses/AGPLv3) [![Go Report Card](https://goreportcard.com/badge/github.com/cinar/indicator/v2)](https://goreportcard.com/report/github.com/cinar/indicator/v2) ![Go CI](https://github.com/cinar/indicator/actions/workflows/ci.yml/badge.svg) [![codecov](https://codecov.io/gh/cinar/indicator/graph/badge.svg?token=MB7L69UAWM)](https://codecov.io/gh/cinar/indicator)

Indicator Go
============

Indicator is a Golang module that provides a rich set of technical analysis indicators, strategies, and a framework for backtesting.

### Major improvements in v2:

-	**Enhanced Code Quality:** A complete rewrite was undertaken to achieve and maintain at least 90% code coverage.
-	**Improved Testability:** Each indicator and strategy have dedicated test data in CSV format for easier validation.
-	**Streamlined Data Handling:** The library was rewritten to operate on data streams (Go channels) for both inputs and outputs. If you prefer using slices, helper functions like [helper.SliceToChan](helper/README.md#func-slicetochan) and [helper.ChanToSlice](helper/README.md#func-chantoslice) are available. Alternatively, you can still use the v1 version.
-	**Configurable Indicators and Strategies:** All indicators and strategies were designed to be fully configurable with no preset values.
-	**Generics Support:** The library leverages Golang generics to support various numeric data formats.

### Migrating from v1 to v2:

-	The porting process from v1 to v2 is ongoing. Indicators or strategies without documentation links are not currently implemented in v2.
-	Contributions are welcome! Feel free to contribute to the project and help expand the library's capabilities.

### Find the v1 library:

You can find the previous version (v1) of the library in the [v1 branch](https://github.com/cinar/indicator/tree/v2) on GitHub.

> [!NOTE] I also have a TypeScript version of this module now at [Indicator TS](https://github.com/cinar/indicatorts).

👆 Indicators Provided
----------------------

The following list of indicators are currently supported by this package:

### 📈 Trend Indicators

-	[Absolute Price Oscillator (APO)](trend/README.md#type-apo)
-	[Aroon Indicator](trend/README.md#type-aroon)
-	[Balance of Power (BoP)](trend/README.md#type-bop)
-	Chande Forecast Oscillator (CFO)
-	[Community Channel Index (CCI)](trend/README.md#type-cci)
-	[Hull Moving Average (HMA)](trend/README.md#type-hma)
-	[Double Exponential Moving Average (DEMA)](trend/README.md#type-dema)
-	[Exponential Moving Average (EMA)](trend/README.md#type-ema)
-	[Kaufman's Adaptive Moving Average (KAMA)](trend/README.md#type-kama)
-	[Mass Index (MI)](trend/README.md#type-massindex)
-	[Moving Average Convergence Divergence (MACD)](trend/README.md#type-macd)
-	[Moving Least Square (MLS)](trend/README.md#type-mls)
-	[Moving Linear Regression (MLR)](trend/README.md#type-mlr)
-	[Moving Max](trend/README.md#type-movingmax)
-	[Moving Min](trend/README.md#type-movingmin)
-	[Moving Sum](trend/README.md#type-movingsum)
-	Parabolic SAR
-	[Random Index (KDJ)](trend/README.md#type-kdj)
-	[Rolling Moving Average (RMA)](trend/README.md#type-rma)
-	[Simple Moving Average (SMA)](trend/README.md#type-sma)
-	[Since Change](helper/README.md#func-since)
-	[Triple Exponential Moving Average (TEMA)](trend/README.md#type-tema)
-	[Triangular Moving Average (TRIMA)](trend/README.md#type-trima)
-	[Triple Exponential Average (TRIX)](trend/README.md#type-trix)
-	[Typical Price](trend/README.md#type-typicalprice)
-	[Volume Weighted Moving Average (VWMA)](trend/README.md#type-vwma)
-	Vortex Indicator
-	[Weighted Moving Average (WMA)](trend/README.md#type-wma)

### 🚀 Momentum Indicators

-	[Awesome Oscillator](momentum/README.md#type-awesomeoscillator)
-	[Chaikin Oscillator](momentum/README.md#type-chaikinoscillator)
-	[Ichimoku Cloud](momentum/README.md#type-ichimokucloud)
-	[Percentage Price Oscillator (PPO)](momentum/README.md#type-ppo)
-	[Percentage Volume Oscillator (PVO)](momentum/README.md#type-pvo)
-	[Relative Strength Index (RSI)](momentum/README.md#type-rsi)
-	[Qstick](momentum/README.md#type-qstick)
-	[Stochastic Oscillator](momentum/README.md#type-stochasticoscillator)
-	[Williams R](momentum/README.md#type-williamsr)

### 🎢 Volatility Indicators

-	[Acceleration Bands](volatility/README.md#type-accelerationbands)
-	[Actual True Range (ATR)](volatility/README.md#type-atr)
-	[Bollinger Band Width](volatility/README.md#type-bollingerbandwidth)
-	[Bollinger Bands](volatility/README.md#type-bollingerbands)
-	[Chandelier Exit](volatility/README.md#type-chandelierexit)
-	[Donchian Channel (DC)](volatility/README.md#type-donchianchannel)
-	[Keltner Channel (KC)](volatility/README.md#type-keltnerchannel)
-	[Moving Standard Deviation (Std)](volatility/README.md#type-movingstd)
-	[Projection Oscillator (PO)](volatility/README.md#type-po)
-	[Super Trend](volatility/README.md#type-supertrend)
-	[Ulcer Index (UI)](volatility/README.md#type-ulcerindex)

### 📢 Volume Indicators

-	[Accumulation/Distribution (A/D)](volume/README.md#type-ad)
-	[Chaikin Money Flow (CMF)](volume/README.md#type-cmf)
-	[Ease of Movement (EMV)](volume/README.md#type-emv)
-	[Force Index (FI)](volume/README.md#type-fi)
-	[Money Flow Index (MFI)](volume/README.md#type-mfi)
-	[Money Flow Multiplier (MFM)](volume/README.md#type-mfm)
-	[Money Flow Volume (MFV)](volume/README.md#type-mfv)
-	[Negative Volume Index (NVI)](volume/README.md#type-nvi)
-	[On-Balance Volume (OBV)](volume/README.md#type-obv)
-	[Volume Price Trend (VPT)](volume/README.md#type-vpt)
-	[Volume Weighted Average Price (VWAP)](volume/README.md#type-vwap)

🧠 Strategies Provided
----------------------

The following list of strategies are currently supported by this package:

### ⚖ Base Strategies

-	[Buy and Hold Strategy](strategy/README.md#type-buyandholdstrategy)

### 📈 Trend Strategies

-	[Absolute Price Oscillator (APO) Strategy](strategy/trend/README.md#type-apostrategy)
-	[Aroon Strategy](strategy/trend/README.md#type-aroonstrategy)
-	[Balance of Power (BoP) Strategy](strategy/trend/README.md#type-bopstrategy)
-	[Double Exponential Moving Average (DEMA) Strategy](strategy/trend/README.md#type-demastrategy)
-	Chande Forecast Oscillator Strategy
-	[Community Channel Index (CCI) Strategy](strategy/trend/README.md#type-ccistrategy)
-	[Golden Cross Strategy](strategy/trend/README.md#type-goldencrossstrategy)
-	[Random Index (KDJ) Strategy](strategy/trend/README.md#type-kdjstrategy)
-	[Moving Average Convergence Divergence (MACD) Strategy](strategy/trend/README.md#type-macdstrategy)
-	[Qstick Strategy](strategy/trend/README.md#type-qstickstrategy)
-	[Triangular Moving Average (TRIMA) Strategy](strategy/trend/README.md#type-trimastrategy)
-	[Triple Exponential Average (TRIX) Strategy](strategy/trend/README.md#type-trixstrategy)
-	[Triple Moving Average Crossover Strategy](strategy/trend/README.md#type-triplemovingaveragecrossoverstrategy)
-	[Volume Weighted Moving Average (VWMA) Strategy](strategy/trend/README.md#type-vwmastrategy)

### 🚀 Momentum Strategies

-	[Awesome Oscillator Strategy](strategy/momentum/README.md#type-awesomeoscillatorstrategy)
-	[RSI Strategy](strategy/momentum/README.md#type-rsistrategy)
-	Williams R Strategy

### 🎢 Volatility Strategies

-	[Bollinger Bands Strategy](strategy/volatility/README.md#type-bollingerbandsstrategy)
-	Projection Oscillator Strategy

### 📢 Volume Strategies

-	Chaikin Money Flow Strategy
-	Ease of Movement Strategy
-	Force Index Strategy
-	Money Flow Index Strategy
-	Negative Volume Index Strategy
-	Volume Weighted Average Price Strategy

### 🧪 Compound Strategies

Compound strategies merge multiple strategies to produce integrated recommendations. They combine individual strategies' recommendations using various decision-making logic.

-	[All Strategy](strategy/README.md#type-allstrategy)
-	[Or Strategy](strategy/README.md#type-orstrategy)
-	[Majority Strategy](strategy/README.md#type-majoritystrategy)
-	[MACD-RSI Strategy](strategy/compound/README.md#type-macdrsistrategy)
-	[Split Strategy](strategy/README.md#type-splitstrategy)

🗃 Repositories
--------------

Repository serves as a centralized storage and retrieval location for [asset snapshots](asset/README.md#type-snapshot).

The following [repository implementations](asset/README.md#type-repository) are provided.

-	[File System Repository](asset/README.md#type-filesystemrepository)
-	[In Memory Repository](asset/README.md#type-inmemoryrepository)
-	[Tiingo Repository](asset/README.md#type-tiingorepository)

The [Sync function]() facilitates the synchronization of assets between designated source and target repositories by employing multi-worker concurrency for enhanced efficiency. This function serves the purpose of procuring the most recent snapshots from remote repositories and seamlessly transferring them to local repositories, such as file system repositories.

The `indicator-sync` command line tool also offers the capability of synchronizing data between the Tiingo Repository and the File System Repository. To illustrate its usage, consider the following example command:

```bash
$ indicator-sync -key $TIINGO_KEY -target /home/user/assets -days 30
```

This command effectively retrieves the most recent snapshots for assets residing within the `/home/user/assets` directory from the Tiingo Repository. In the event that the local asset file is devoid of content, it automatically extends its reach to synchronize 30 days' worth of snapshots, ensuring a comprehensive and up-to-date repository.

⏳ Backtesting
--------------

The [Backtest functionality](strategy/README.md#type-backtest), using the [Outcome](strategy/README.md#func-outcome), rigorously evaluates the potential performance of the specified strategies applied to a defined set of assets. It generates comprehensive visual representations for each strategy-asset pairing.

```go
backtest := strategy.NewBacktest(repository, outputDir)
backtest.Names = append(backtest.Names, "brk-b")
backtest.Strategies = append(backtest.Strategies, trend.NewApoStrategy())

err = backtest.Run()
if err != nil {
	t.Fatal(err)
}
```

The `indicator-backtest` command line tool empowers users to conduct comprehensive backtesting of assets residing within a specified repository. This capability encompasses the application of all currently recognized strategies, culminating in the generation of detailed reports within a designated output directory.

```bash
$ indicator-backtest -repository /home/user/assets -output /home/user/reports -workers 1
```

Usage
-----

Install package.

```bash
go get github.com/cinar/indicator/v2
```

Import indicator.

```Golang
import (
    "github.com/cinar/indicator/v2"
)
```

Contributing to the Project
---------------------------

Anyone can contribute to Indicator library. Please make sure to read our [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md) guide first. Follow the [How to Contribute to Indicator](./CONTRIBUTING.md) to contribute. Signining a [Contributor Agreement](./CLA.md) is also required to contribute to the project.

Disclaimer
----------

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

License
-------

The `v2.x.x` and above are dual-licensed under GNU AGPLv3 License and a commercial license. For free use and modifications of the code, you can use the AGPLv3 license. If you require commercial license with different terms, please contact me.

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
