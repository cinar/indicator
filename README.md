<p align="center">
    <img src="logo.png" />
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/cinar/indicator/v2"><img src="https://img.shields.io/badge/Go_Reference-007D9C?style=for-the-badge&logo=go&logoColor=white" alt="Go Reference" /></a>
    <a href="LICENSE.md"><img src="https://img.shields.io/github/license/cinar/indicator?style=for-the-badge" alt="License" /></a>
    <a href="https://github.com/cinar/indicator/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/cinar/indicator/ci.yml?branch=main&style=for-the-badge&logo=github&label=CI" alt="Go CI" /></a>
    <a href="https://codecov.io/gh/cinar/indicator"><img src="https://img.shields.io/codecov/c/github/cinar/indicator?style=for-the-badge&logo=codecov&logoColor=white" alt="Codecov" /></a>
    <a href="https://github.com/cinar/indicator/pkgs/container/indicator"><img src="https://img.shields.io/github/v/tag/cinar/indicator?label=ghcr&sort=semver&logo=github&style=for-the-badge" alt="GHCR Version" /></a>
    <a href="https://github.com/cinar/indicator/stargazers"><img src="https://img.shields.io/github/stars/cinar/indicator?style=for-the-badge&logo=github&logoColor=white" alt="GitHub Stars" /></a>
</p>
<p align="center">
    <a href="https://trendshift.io/repositories/23333?utm_source=repository-badge&amp;utm_medium=badge&amp;utm_campaign=badge-repository-23333" target="_blank" rel="noopener noreferrer"><img src="https://trendshift.io/api/badge/repositories/23333" alt="cinar%2Findicator | Trendshift" width="250" height="55"/></a>
</p>


Indicator Go
============

Indicator is a Golang module that provides an extensive set of technical analysis indicators, strategies, and a framework for backtesting.

> An extensive technical analysis library for algorithmic trading - 80+ indicators and backtesting framework.

### Major improvements in v2:

-	**Enhanced Code Quality:** A complete rewrite was undertaken to achieve and maintain at least 90% code coverage.
-	**Improved Testability:** Each indicator and strategy have dedicated test data in CSV format for easier validation.
-	**Streamlined Data Handling:** The library was rewritten to operate on data streams (Go channels) for both inputs and outputs. If you prefer using slices, helper functions like [helper.SliceToChan](helper/README.md#SliceToChan) and [helper.ChanToSlice](helper/README.md#ChanToSlice) are available. Alternatively, you can still use the [v1 version](https://github.com/cinar/indicator/tree/v1).
-	**Configurable Indicators and Strategies:** All indicators and strategies were designed to be fully configurable with no preset values.
-	**Generics Support:** The library leverages Golang generics to support various numeric data formats.

There is also a TypeScript version of this module at [Indicator TS](https://github.com/cinar/indicatorts).

⚠️ Risk Disclosure & Legal Disclaimer
------------------------------------

1. **Strict Educational, Research, and Informational Purpose Notice:**
This repository, software library, documentation, and all associated tools, mathematical models, indicator calculations, example strategies, and backtesting utilities are provided strictly and solely for educational, academic, and research purposes. Nothing contained herein constitutes, or is intended to constitute, investment, financial, tax, legal, or trading advice. No communication or output from this software shall be construed as a recommendation, endorsement, solicitation, or offer to buy or sell any security, commodity, currency, digital asset, or other financial instrument.

2. **Non-Advisory and Non-Fiduciary Status:**
The Indicator Authors, maintainers, contributors, and copyright holders are not registered investment advisers (RIAs), registered broker-dealers, commodity trading advisors (CTAs), commodity pool operators (CPOs), financial planners, or tax advisors with the U.S. Securities and Exchange Commission (SEC), the U.S. Commodity Futures Trading Commission (CFTC), the Financial Industry Regulatory Authority (FINRA), or any other regulatory agency or authority in any jurisdiction. The maintainers do not act in any fiduciary capacity, and no fiduciary, advisory, or professional-client relationship is established through the use of, contribution to, or access to this software.

3. **Risk Disclosure Regarding Substantial Trading Losses and Leverage:**
Trading and investing in equities, options, futures, foreign exchange (Forex), commodities, and cryptocurrencies carries substantial risk of loss and is not suitable for all individuals or investors. Prices of financial instruments fluctuate significantly and unpredictably. The use of leverage, margin, or derivative contracts can rapidly amplify losses as well as gains, potentially resulting in the complete loss of all invested capital or losses exceeding initial deposits. You should carefully consider whether trading is appropriate for you in light of your personal financial condition, experience level, and risk tolerance, and you should consult with qualified, licensed financial advisors before making any investment or trading decisions.

4. **CFTC-Style Hypothetical and Backtested Performance Disclosure (CFTC Rule 4.41):**
Hypothetical, simulated, or backtested performance results have inherent and significant limitations. Unlike an actual performance record, simulated results do not represent actual trading. Because trades have not actually been executed, results may over- or under-compensate for the impact, if any, of certain market factors, including but not limited to lack of liquidity, order execution delays, market impact, exchange fees, slippage, and spread. Furthermore, hypothetical and simulated trading programs are designed with the benefit of hindsight and retrospective knowledge of historical price movements. No representation or warranty is being made that any account, strategy, or portfolio will or is likely to achieve profits, avoid losses, or experience outcomes similar to those illustrated in any backtest, example, or simulation.

5. **"AS IS" Software Liability Limitation & Disclaimers:**
This software is provided on an "AS IS" and "AS AVAILABLE" basis, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, TITLE, ACCURACY, COMPLETENESS, AND NON-INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS, MAINTAINERS, CONTRIBUTORS, OR COPYRIGHT HOLDERS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, PUNITIVE, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, FINANCIAL LOSSES, TRADING LOSSES, LOSS OF CAPITAL, LOSS OF PROFITS, LOSS OF DATA, SYSTEM FAILURES, INACCURATE CALCULATIONS, OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND UNDER ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF OR INABILITY TO USE THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGES.

👆 Indicators Provided
----------------------

The following list of indicators are currently supported by this package:

### 📈 Trend Indicators

-	[Absolute Price Oscillator (APO)](trend/README.md#Apo)
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
-	[Rate of Change (ROC)](trend/README.md#Roc)
-	[Stochastic](trend/README.md#Stochastic)
-	[Slow Stochastic](trend/README.md#SlowStochastic)
-	[Schaff Trend Cycle (STC)](trend/README.md#Stc)
-	[Rolling Moving Average (RMA)](trend/README.md#Rma)
-	[Simple Moving Average (SMA)](trend/README.md#Sma)
-	[Since Change](helper/README.md#Since)
-	[Slope](trend/README.md#Slope)
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
-	[Internal Bar Strength (IBS)](momentum/README.md#InternalBarStrength)
-   [Martin Pring's Special K](momentum/README.md#PringsSpecialK)
-	[Percentage Price Oscillator (PPO)](momentum/README.md#Ppo)
-	[Percentage Volume Oscillator (PVO)](momentum/README.md#Pvo)
-	[Qstick](momentum/README.md#Qstick)
-	[Relative Strength Index (RSI)](momentum/README.md#Rsi)
-	[Relative Vigor Index (RVI)](momentum/README.md#Rvi)
-	[Stochastic Oscillator](momentum/README.md#StochasticOscillator)
-	[Stochastic RSI](momentum/README.md#StochasticRsi)
-	[TD Sequential](momentum/README.md#TdSequential)
-	[Ultimate Oscillator](momentum/README.md#UltimateOscillator)
-	[Williams R](momentum/README.md#WilliamsR)

### 🎢 Volatility Indicators

-   [Percent B](volatility/README.md#PercentB)
-	[Acceleration Bands](volatility/README.md#AccelerationBands)
-	[Annualized Historical Volatility (AHV)](volatility/README.md#AnnualizedHistoricalVolatility)
-	[Average True Range (ATR)](volatility/README.md#Atr)
-	[Bollinger Band Width](volatility/README.md#BollingerBandWidth)
-	[Bollinger Bands](volatility/README.md#BollingerBands)
-	[Chandelier Exit](volatility/README.md#ChandelierExit)
-	[Choppiness Index (CHOP)](volatility/README.md#Chop)
-	[Donchian Channel (DC)](volatility/README.md#DonchianChannel)
-	[Historical Volatility (HV)](volatility/README.md#HistoricalVolatility)
-	[Keltner Channel (KC)](volatility/README.md#KeltnerChannel)
-	[Moving Standard Deviation (Std)](volatility/README.md#MovingStd)
-	[Projection Oscillator (PO)](volatility/README.md#Po)
-	[Super Trend](volatility/README.md#SuperTrend)
-	[True Range (TR)](volatility/README.md#TrueRange)
-	[Ulcer Index (UI)](volatility/README.md#UlcerIndex)
-	[Z-Score](volatility/README.md#ZScore)


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

💡 Example Strategies (Educational Demonstrations)
---------------------------------------------------

The following concrete strategies are provided in the `examples/` directory purely as illustrative demonstrations on how developers can consume the core indicator mathematics and generic backtesting framework. They are intended strictly for educational and research purposes:

### ⚖ Base Strategies

-	[Buy and Hold Strategy](strategy/README.md#BuyAndHoldStrategy)

### 📈 Trend Strategy Examples

-   [Alligator Strategy](examples/trend/README.md#AlligatorStrategy)
-	[Absolute Price Oscillator (APO) Strategy](examples/trend/README.md#ApoStrategy)
-	[Aroon Strategy](examples/trend/README.md#AroonStrategy)
-	[Balance of Power (BoP) Strategy](examples/trend/README.md#BopStrategy)
-	[Chande Forecast Oscillator Strategy](examples/trend/README.md#CfoStrategy)
-	[Commodity Channel Index (CCI) Strategy](examples/trend/README.md#CciStrategy)
-	[Double Exponential Moving Average (DEMA) Strategy](examples/trend/README.md#DemaStrategy)
-   [Envelope Strategy](examples/trend/README.md#EnvelopeStrategy)
-	[Golden Cross Strategy](examples/trend/README.md#GoldenCrossStrategy)
-	[Hull Moving Average (HMA) Strategy](examples/trend/README.md#HmaStrategy)
-	[Kaufman's Adaptive Moving Average (KAMA) Strategy](examples/trend/README.md#KamaStrategy)
-	[Moving Average Convergence Divergence (MACD) Strategy](examples/trend/README.md#MacdStrategy)
-	[Qstick Strategy](examples/trend/README.md#QstickStrategy)
-	[Random Index (KDJ) Strategy](examples/trend/README.md#KdjStrategy)
-   [Smoothed Moving Average (SMMA) Strategy](examples/trend/README.md#SmmaStrategy)
-	[Triangular Moving Average (TRIMA) Strategy](examples/trend/README.md#TrimaStrategy)
-	[Triple Exponential Average (TRIX) Strategy](examples/trend/README.md#TrixStrategy)
-	[Triple Moving Average Crossover Strategy](examples/trend/README.md#TripleMovingAverageCrossoverStrategy)
-	[True Strength Index (TSI) Strategy](examples/trend/README.md#TsiStrategy)
-	[Volume Weighted Moving Average (VWMA) Strategy](examples/trend/README.md#VwmaStrategy)
-   [Weighted Close Strategy](examples/trend/README.md#WeightedCloseStrategy)

### 🚀 Momentum Strategy Examples

-	[Awesome Oscillator Strategy](examples/momentum/README.md#AwesomeOscillatorStrategy)
-	[Coppock Curve Strategy](examples/momentum/README.md#CoppockCurveStrategy)
-	[Elder Ray Strategy](examples/momentum/README.md#ElderRayStrategy)
-	[Ichimoku Cloud Strategy](examples/momentum/README.md#IchimokuCloudStrategy)
-	[RSI Strategy](examples/momentum/README.md#RsiStrategy)
-	[Stochastic Oscillator Strategy](examples/momentum/README.md#StochasticOscillatorStrategy)
-	[Stochastic RSI Strategy](examples/momentum/README.md#StochasticRsiStrategy)
-	[Triple RSI Strategy](examples/momentum/README.md#TripleRsiStrategy)
-	[Williams R Strategy](examples/momentum/README.md#WilliamsRStrategy)

### 🎢 Volatility Strategy Examples

-	[Bollinger Bands Strategy](examples/volatility/README.md#BollingerBandsStrategy)
-	[Donchian Channel Breakout Strategy](examples/volatility/README.md#DonchianChannelBreakoutStrategy)
-	[Super Trend Strategy](examples/volatility/README.md#SuperTrendStrategy)

### 📢 Volume Strategy Examples

-	[Chaikin Money Flow Strategy](examples/volume/README.md#ChaikinMoneyFlowStrategy)
-	[Ease of Movement Strategy](examples/volume/README.md#EaseOfMovementStrategy)
-	[Force Index Strategy](examples/volume/README.md#ForceIndexStrategy)
-	[Money Flow Index Strategy](examples/volume/README.md#MoneyFlowIndexStrategy)
-	[Negative Volume Index Strategy](examples/volume/README.md#NegativeVolumeIndexStrategy)
-	[On-Balance Volume (OBV) Strategy](examples/volume/README.md#ObvStrategy)
-	[Percent Band and MFI Strategy](examples/volume/README.md#PercentBandMFIStrategy)
-	[Weighted Average Price Strategy](examples/volume/README.md#WeightedAveragePriceStrategy)

### 🧪 Compound Strategy Examples

Compound strategies demonstrate merging multiple indicator signals into unified actions using various logical combinations.

-	[And Strategy](strategy/README.md#AndStrategy)
-	[Majority Strategy](strategy/README.md#MajorityStrategy)
-	[MACD-RSI Strategy](examples/compound/README.md#MacdRsiStrategy)
-	[Or Strategy](strategy/README.md#OrStrategy)
-	[Split Strategy](strategy/README.md#SplitStrategy)

### 🎁 Decorator Strategy Examples

Decorator strategies offer an illustrative way to alter or filter the recommendations of underlying strategies.

-   [Inverse Strategy](examples/decorator/README.md#InverseStrategy)
-   [Cost Basis Exit Strategy](examples/decorator/README.md#CostBasisExitStrategy)
-   [Stop Loss Strategy](examples/decorator/README.md#StopLossStrategy)

🗃 Repositories
--------------

Repository serves as a centralized storage and retrieval location for [asset snapshots](asset/README.md#Snapshot).

The following [repository implementations](asset/README.md#Repository) are provided.

-	[File System Repository](asset/README.md#FileSystemRepository)
-	[In Memory Repository](asset/README.md#InMemoryRepository)
-	[Tiingo Repository](asset/README.md#TiingoRepository)

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

The [Backtest functionality](backtest/README.md#Backtest), using the [Outcome](strategy/README.md#Outcome), evaluates the hypothetical performance of specified strategies applied to a defined set of historical asset data. It generates comprehensive visual representations for each strategy-asset pairing.

```go
import (
	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/examples/trend"
)

report := backtest.NewHTMLReport(outputDir)
bt := backtest.NewBacktest(repository, report)
bt.Names = append(bt.Names, "brk-b")
bt.Strategies = append(bt.Strategies, trend.NewApoStrategy())

err = bt.Run()
if err != nil {
	t.Fatal(err)
}
```

The `indicator-backtest` command line tool empowers users to conduct comprehensive backtesting of assets residing within a specified repository. It does not run any strategy by default — you must explicitly name the strategies to backtest with the `-strategies` flag, culminating in the generation of detailed reports within a designated output directory.

```bash
$ indicator-backtest \
    -repository-name filesystem \
    -repository-config /home/user/assets \
    -report-config /home/user/reports \
    -strategies apo,macd,rsi \
    -workers 1
```

Run `indicator-backtest -list-strategies` to print the full list of strategy names that are available to pass to `-strategies`.

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
  --assets aapl msft googl \
  --strategies apo,macd,rsi

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
| `--strategies` | Comma-separated strategy names to backtest (required) | - |
| `--output` | Output directory for reports | /app/output |

Run `docker run --rm --entrypoint ./indicator-backtest ghcr.io/cinar/indicator:latest -list-strategies` to see the available strategy names.

### Examples

```bash
# Backtest all available assets for 1 year
docker run -it --rm \
  -v $(pwd)/reports:/app/output \
  ghcr.io/cinar/indicator:latest \
  --api-key YOUR_TIINGO_API_KEY \
  --strategies apo,macd,rsi

# Backtest specific stocks for last 6 months, test last 30 days
docker run -it --rm \
  -v $(pwd)/reports:/app/output \
  ghcr.io/cinar/indicator:latest \
  --api-key YOUR_TIINGO_API_KEY \
  --days 180 \
  --last 30 \
  --assets aapl msft googl amzn \
  --strategies apo,macd,rsi

# Custom output directory
docker run -it --rm \
  -v /path/to/my/reports:/app/output \
  ghcr.io/cinar/indicator:latest \
  --api-key YOUR_TIINGO_API_KEY \
  --strategies apo,macd,rsi \
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

💖 Our Sponsors
---------------

Indicator is a community-supported project. The following companies, organizations, and individuals help make our work possible.  Become [a sponsor](https://github.com/sponsors/cinar) and help us continue to grow!

![Our Sponsors](./sponsors.svg)

Contributing to the Project
---------------------------

Anyone can contribute to Indicator library. Please make sure to read our [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md) guide first. Follow the [How to Contribute to Indicator](./CONTRIBUTING.md) to contribute. Signining a [Contributor Agreement](./CLA.md) is also required to contribute to the project.

License
-------

The `v2.x.x` and above are dual-licensed under GNU AGPLv3 License and a commercial license. For free use and modifications of the code, you can use the AGPLv3 license. If you require commercial license with different terms, please contact the maintainers.

```
Copyright (c) 2021-2026 The Indicator Authors.    
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
Copyright (c) 2021-2026 The Indicator Authors.
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

Trademarks
----------

- **TD Sequential®** is a registered trademark of DeMark Analytics, LLC. This library is an independent open-source mathematical implementation and is not affiliated with, sponsored by, or endorsed by DeMark Analytics, LLC.
- **Bollinger Bands®** is a registered trademark of John Bollinger.
- **ConnorsRSI®** is a registered trademark of Connors Group, Inc.
- All other product names, logos, and brands mentioned herein are trademarks or registered trademarks of their respective owners. Mention of third-party products, services, or trademarks is for nominative identification and educational purposes only and does not imply affiliation or endorsement.
