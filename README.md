[![GoDoc](https://godoc.org/github.com/cinar/indicator/v2?status.svg)](https://godoc.org/github.com/cinar/indicator/v2) [![License](https://img.shields.io/badge/License-AGPLv3-blue.svg)](https://opensource.org/licenses/AGPLv3) [![Go Report Card](https://goreportcard.com/badge/github.com/cinar/indicator/v2)](https://goreportcard.com/report/github.com/cinar/indicator/v2) ![Go CI](https://github.com/cinar/indicator/actions/workflows/ci.yml/badge.svg) [![codecov](https://codecov.io/gh/cinar/indicator/graph/badge.svg?token=MB7L69UAWM)](https://codecov.io/gh/cinar/indicator) [![GHCR Version](https://img.shields.io/github/v/tag/cinar/indicator?label=ghcr&sort=semver&logo=github)](https://github.com/cinar/indicator/pkgs/container/indicator)

<p align="center">
    <img src="logo.png" />
</p>

Indicator Go
============

Indicator is a Golang module that provides a rich set of technical analysis indicators, strategies, and a framework for backtesting.

### Major improvements in v2:

-	**Enhanced Code Quality:** A complete rewrite was undertaken to achieve and maintain at least 90% code coverage.
-	**Improved Testability:** Each indicator and strategy have dedicated test data in CSV format for easier validation.
-	**Streamlined Data Handling:** The library was rewritten to operate on data streams (Go channels) for both inputs and outputs. If you prefer using slices, helper functions like [helper.SliceToChan](helper/README.md#SliceToChan) and [helper.ChanToSlice](helper/README.md#ChanToSlice) are available. Alternatively, you can still use the [v1 version](https://github.com/cinar/indicator/tree/v1).
-	**Configurable Indicators and Strategies:** All indicators and strategies were designed to be fully configurable with no preset values.
-	**Generics Support:** The library leverages Golang generics to support various numeric data formats.
-   **MCP Support:** MCP (Multi-Client Protocol Server) support is integrated into the library, facilitating its use with various AI tools.

I also have a TypeScript version of this module now at [Indicator TS](https://github.com/cinar/indicatorts).

👆 Indicators Provided
----------------------

The following list of indicators are currently supported by this package:

### 📈 Trend Indicators

-	[Absolute Price Oscillator (APO)](trend/README.md#Apo)
-   [Alligator Indicator](trend/README.md#Alligator)
-	[Aroon Indicator](trend/README.md#Aroon)
-	[Balance of Power (BoP)](trend/README.md#Bop)
-	[Chande Forecast Oscillator (CFO)](trend/README.md#Cfo)
-	[Commodity Channel Index (CCI)](trend/README.md#Cci)
-   [Envelope](trend/README.md#Envelope)
-	[Hull Moving Average (HMA)](trend/README.md#Hma)
-   [Detrended Price Oscillator (DPO)](trend/README.md#Dpo)
-	[Double Exponential Moving Average (DEMA)](trend/README.md#Dema)
-	[Exponential Moving Average (EMA)](trend/README.md#Ema)
-	[Kaufman's Adaptive Moving Average (KAMA)](trend/README.md#Kama)
-	[Know Sure Thing (KST)](trend/README.md#Kst)
-	[Mass Index (MI)](trend/README.md#MassIndex)
-	[McGinley Dynamic](trend/README.md#McGinleyDynamic)
-	[Moving Average Convergence Divergence (MACD)](trend/README.md#Macd)
-	[Moving Least Square (MLS)](trend/README.md#Mls)
-	[Moving Linear Regression (MLR)](trend/README.md#Mlr)
-	[Moving Max](trend/README.md#MovingMax)
-	[Moving Min](trend/README.md#MovingMin)
-	[Moving Sum](trend/README.md#MovingSum)
-	[Pivot Point](trend/README.md#PivotPoint)
-	[Random Index (KDJ)](trend/README.md#Kdj)
-	[Stochastic](trend/README.md#Stochastic)
-	[Slow Stochastic](trend/README.md#SlowStochastic)
-	[Schaff Trend Cycle (STC)](trend/README.md#Stc)
-	[Rolling Moving Average (RMA)](trend/README.md#Rma)
-	[Simple Moving Average (SMA)](trend/README.md#Sma)
-	[Since Change](helper/README.md#Since)
-   [Smoothed Moving Average (SMMA)](trend/README.md#Smma)
-	[Triple Exponential Moving Average (TEMA)](trend/README.md#Tema)
-	[Triangular Moving Average (TRIMA)](trend/README.md#Trima)
-	[Triple Exponential Average (TRIX)](trend/README.md#Trix)
-	[True Strength Index (TSI)](trend/README.md#Tsi)
-	[Tillson T3](trend/README.md#T3)
-	[Typical Price](trend/README.md#TypicalPrice)
-	[Volume Weighted Moving Average (VWMA)](trend/README.md#Vwma)
-   [Weighted Close](trend/README.md#WeightedClose)
-	[Weighted Moving Average (WMA)](trend/README.md#Wma)

### 🚀 Momentum Indicators

-	[Awesome Oscillator](momentum/README.md#AwesomeOscillator)
-	[Chaikin Oscillator](momentum/README.md#ChaikinOscillator)
-	[Connors RSI](momentum/README.md#ConnorsRsi)
-	[Coppock Curve](momentum/README.md#CoppockCurve)
-	[Elder-Ray Index](momentum/README.md#ElderRay)
-	[Fisher Transform](momentum/README.md#Fisher)
-	[Ichimoku Cloud](momentum/README.md#IchimokuCloud)
-	[Percentage Price Oscillator (PPO)](momentum/README.md#Ppo)
-	[Percentage Volume Oscillator (PVO)](momentum/README.md#Pvo)
-   [Martin Pring's Special K](momentum/README.md#PringsSpecialK)
-	[Relative Strength Index (RSI)](momentum/README.md#Rsi)
-	[Relative Vigor Index (RVI)](momentum/README.md#Rvi)
-	[Qstick](momentum/README.md#Qstick)
-	[Stochastic Oscillator](momentum/README.md#StochasticOscillator)
-	[Stochastic RSI](momentum/README.md#StochasticRsi)
-	[TD Sequential](momentum/README.md#TdSequential)

### 🎢 Volatility Indicators

-   [Percent B](volatility/README.md#PercentB)
-	[Acceleration Bands](volatility/README.md#AccelerationBands)
-	[Average True Range (ATR)](volatility/README.md#Atr)
-	[Bollinger Band Width](volatility/README.md#BollingerBandWidth)
-	[Bollinger Bands](volatility/README.md#BollingerBands)
-	[Chandelier Exit](volatility/README.md#ChandelierExit)
-	[Choppiness Index (CHOP)](volatility/README.md#Chop)
-	[Donchian Channel (DC)](volatility/README.md#DonchianChannel)
-	[Keltner Channel (KC)](volatility/README.md#KeltnerChannel)
-	[Moving Standard Deviation (Std)](volatility/README.md#MovingStd)
-	[Projection Oscillator (PO)](volatility/README.md#Po)
-	[Super Trend](volatility/README.md#SuperTrend)
-	[Ulcer Index (UI)](volatility/README.md#UlcerIndex)

### 📢 Volume Indicators

-	[Accumulation/Distribution (A/D)](volume/README.md#Ad)
-	[Chaikin Money Flow (CMF)](volume/README.md#Cmf)
-	[Ease of Movement (EMV)](volume/README.md#Emv)
-	[Force Index (FI)](volume/README.md#Fi)
-	[Klinger Volume Oscillator (KVO)](volume/README.md#Kvo)
-	[Money Flow Index (MFI)](volume/README.md#Mfi)
-	[Money Flow Multiplier (MFM)](volume/README.md#Mfm)
-	[Money Flow Volume (MFV)](volume/README.md#Mfv)
-	[Negative Volume Index (NVI)](volume/README.md#Nvi)
-	[On-Balance Volume (OBV)](volume/README.md#Obv)
-	[Volume Price Trend (VPT)](volume/README.md#Vpt)
-	[Volume Weighted Average Price (VWAP)](volume/README.md#Vwap)

### 💰 Asset Valuation
-   [Future Value (FV)](valuation/README.md#Fv)
-   [Net Present Value (NPV)](valuation/README.md#Npv)
-   [Present Value (PV)](valuation/README.md#Pv)

🧠 Strategies Provided
----------------------

The following list of strategies are currently supported by this package:

### ⚖ Base Strategies

-	[Buy and Hold Strategy](strategy/README.md#type-buyandholdstrategy)

### 📈 Trend Strategies

-   [Alligator Strategy](strategy/trend/README.md#AlligatorStrategy)
-	[Absolute Price Oscillator (APO) Strategy](strategy/trend/README.md#ApoStrategy)
-	[Aroon Strategy](strategy/trend/README.md#AroonStrategy)
-	[Balance of Power (BoP) Strategy](strategy/trend/README.md#BopStrategy)
-	[Chande Forecast Oscillator Strategy](strategy/trend/README.md#CfoStrategy)
-	[Commodity Channel Index (CCI) Strategy](strategy/trend/README.md#CciStrategy)
-	[Double Exponential Moving Average (DEMA) Strategy](strategy/trend/README.md#DemaStrategy)
-   [Envelope Strategy](strategy/trend/README.md#EnvelopeStrategy)
-	[Golden Cross Strategy](strategy/trend/README.md#GoldenCrossStrategy)
-	[Kaufman's Adaptive Moving Average (KAMA) Strategy](strategy/trend/README.md#KamaStrategy)
-	[Moving Average Convergence Divergence (MACD) Strategy](strategy/trend/README.md#MacdStrategy)
-	[Qstick Strategy](strategy/trend/README.md#QstickStrategy)
-	[Random Index (KDJ) Strategy](strategy/trend/README.md#KdjStrategy)
-   [Smoothed Moving Average (SMMA) Strategy](strategy/trend/README.md#SmmaStrategy)
-	[Triangular Moving Average (TRIMA) Strategy](strategy/trend/README.md#TrimaStrategy)
-	[Triple Exponential Average (TRIX) Strategy](strategy/trend/README.md#TrixStrategy)
-	[Triple Moving Average Crossover Strategy](strategy/trend/README.md#TripleMovingAverageCrossoverStrategy)
-	[True Strength Index (TSI) Strategy](strategy/trend/README.md#TsiStrategy)
-	[Volume Weighted Moving Average (VWMA) Strategy](strategy/trend/README.md#VwmaStrategy)
-   [Weighted Close Strategy](strategy/trend/README.md#WeightedCloseStrategy)

### 🚀 Momentum Strategies

-	[Awesome Oscillator Strategy](strategy/momentum/README.md#AwesomeOscillatorStrategy)
-	[Ichimoku Cloud Strategy](strategy/momentum/README.md#IchimokuCloudStrategy)
-	[RSI Strategy](strategy/momentum/README.md#RsiStrategy)
-	[Stochastic RSI Strategy](strategy/momentum/README.md#StochasticRsiStrategy)
-	[Triple RSI Strategy](strategy/momentum/README.md#TripleRsiStrategy)

### 🎢 Volatility Strategies

-	[Bollinger Bands Strategy](strategy/volatility/README.md#BollingerBandsStrategy)
-	[Super Trend Strategy](strategy/volatility/README.md#SuperTrendStrategy)

### 📢 Volume Strategies

-	[Chaikin Money Flow Strategy](strategy/volume/README.md#ChaikinMoneyFlowStrategy)
-	[Ease of Movement Strategy](strategy/volume/README.md#EaseOfMovementStrategy)
-	[Force Index Strategy](strategy/volume/README.md#ForceIndexStrategy)
-	[Money Flow Index Strategy](strategy/volume/README.md#MoneyFlowIndexStrategy)
-	[Negative Volume Index Strategy](strategy/volume/README.md#NegativeVolumeIndexStrategy)
-	[Percent Band and MFI Strategy](strategy/volume/README.md#PercentBandMFIStrategy)
-	[Weighted Average Price Strategy](strategy/volume/README.md#WeightedAveragePriceStrategy)

### 🧪 Compound Strategies

Compound strategies merge multiple strategies to produce integrated recommendations. They combine individual strategies' recommendations using various decision-making logic.

-	[And Strategy](strategy/README.md#AndStrategy)
-	[Majority Strategy](strategy/README.md#MajorityStrategy)
-	[MACD-RSI Strategy](strategy/compound/README.md#MacdRsiStrategy)
-	[Or Strategy](strategy/README.md#OrStrategy)
-	[Split Strategy](strategy/README.md#SplitStrategy)

### 🎁 Decorator Strategies

Decorator strategies offer a way to alter the recommendations of other strategies.

-   [Inverse Strategy](strategy/decorator/README.md#InverseStrategy)
-   [No Loss Strategy](strategy/decorator/README.md#NoLossStrategy)
-   [Stop Loss Strategy](strategy/decorator/README.md#StopLossStrategy)

🗃 Repositories
--------------

Repository serves as a centralized storage and retrieval location for [asset snapshots](asset/README.md#Snapshot).

The following [repository implementations](asset/README.md#Repository) are provided.

-	[File System Repository](asset/README.md#FileSystemRepository)
-	[In Memory Repository](asset/README.md#InMemoryRepository)
-	[Tiingo Repository](asset/README.md#TiingoRepository)
-	[Alpaca Markets Repository](https://github.com/cinar/indicatoralpaca)

The [Sync function](asset/README.md#Sync) facilitates the synchronization of assets between designated source and target repositories by employing multi-worker concurrency for enhanced efficiency. This function serves the purpose of procuring the most recent snapshots from remote repositories and seamlessly transferring them to local repositories, such as file system repositories.

The `indicator-sync` command line tool also offers the capability of synchronizing data between the Tiingo Repository and the File System Repository. To illustrate its usage, consider the following example command:

```bash
$ indicator-sync \
    -source-name tiingo \
    -source-config $TIINGO_KEY \
    -target-name filesystem \
    -target-config /home/user/assets \
    -days 30
```

This command effectively retrieves the most recent snapshots for assets residing within the `/home/user/assets` directory from the Tiingo Repository. In the event that the local asset file is devoid of content, it automatically extends its reach to synchronize 30 days' worth of snapshots, ensuring a comprehensive and up-to-date repository.

⏳ Backtesting
--------------

The [Backtest functionality](backtest/README.md#Backtest), using the [Outcome](strategy/README.md#Outcome), rigorously evaluates the potential performance of the specified strategies applied to a defined set of assets. It generates comprehensive visual representations for each strategy-asset pairing.

```go
report := backtest.NewHTMLReport(outputDir)
bt := backtest.NewBacktest(repository, report)
bt.Names = append(bt.Names, "brk-b")
bt.Strategies = append(bt.Strategies, trend.NewApoStrategy())

err = bt.Run()
if err != nil {
	t.Fatal(err)
}
```

The `indicator-backtest` command line tool empowers users to conduct comprehensive backtesting of assets residing within a specified repository. This capability encompasses the application of all currently recognized strategies, culminating in the generation of detailed reports within a designated output directory.

```bash
$ indicator-backtest \
    -source-name filesystem \
    -source-config /home/user/assets \
    -output /home/user/reports \
    -workers 1
```

☁️  MCP Server
--------------

The [MCP Server](mcp/README.md) (Multi-Client Protocol Server) provides a robust and scalable solution for serving trading strategies to multiple clients. It enables real-time strategy execution and data processing, making it ideal for applications requiring high-throughput and low-latency interactions with trading algorithms.

🐳 Docker
---------

The easiest way to get started is using our Docker image. It handles everything - syncing market data from Tiingo and generating backtest reports - in a single command.

### Quick Start

```bash
# Get your free Tiingo API key at: https://www.tiingo.com/

# Run backtest for specific assets
docker run -it --rm \
  -v $(pwd)/output:/app/output \
  ghcr.io/cinar/indicator:latest \
  --api-key YOUR_TIINGO_API_KEY \
  --days 365 \
  --assets aapl msft googl

# View results (macOS)
open output/index.html

# View results (Linux)
xdg-open output/index.html
```

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--api-key` | Tiingo API key (required) | - |
| `--days` | Days of historical data to fetch | 365 |
| `--last` | Days to backtest | 365 |
| `--assets` | Space-separated ticker symbols (default: all) | all |
| `--output` | Output directory for reports | /app/output |

### Examples

```bash
# Backtest all available assets for 1 year
docker run -it --rm \
  -v $(pwd)/reports:/app/output \
  ghcr.io/cinar/indicator:latest \
  --api-key YOUR_TIINGO_API_KEY

# Backtest specific stocks for last 6 months, test last 30 days
docker run -it --rm \
  -v $(pwd)/reports:/app/output \
  ghcr.io/cinar/indicator:latest \
  --api-key YOUR_TIINGO_API_KEY \
  --days 180 \
  --last 30 \
  --assets aapl msft googl amzn

# Custom output directory
docker run -it --rm \
  -v /path/to/my/reports:/app/output \
  ghcr.io/cinar/indicator:latest \
  --api-key YOUR_TIINGO_API_KEY \
  --output /app/output
```

### Build Locally

```bash
docker build -t indicator .
docker run -it --rm -v $(pwd)/output:/app/output indicator --api-key YOUR_KEY
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

🌐 Ecosystem
------------

Indicator Go is part of a broader ecosystem of technical analysis tools:

- [Indicator TS](https://github.com/cinar/indicatorts) - TypeScript/JavaScript implementation of the same indicators and strategies
- [Indicator Alpaca](https://github.com/cinar/indicatoralpaca) - Alpaca Markets integration for live trading
- [MCP Server](mcp/README.md) - Model Context Protocol server for AI integration

💖 Our Sponsors
---------------

Indicator is a community-supported project. The following companies, organizations, and individuals help make our work possible.  Become [a sponsor](https://github.com/sponsors/cinar) and help us continue to grow!

![Our Sponsors](./sponsors.svg)

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
Copyright (c) 2021-2026 Onur Cinar.    
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
Copyright (c) 2021-2026 Onur Cinar.
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
