[![GoDoc](https://godoc.org/github.com/cinar/indicator?status.svg)](https://godoc.org/github.com/cinar/indicator)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.com/cinar/indicator.svg?branch=master)](https://travis-ci.com/cinar/indicator)

# Indicator Go

Indicator is a Golang module providing various stock technical analysis indicators, and strategies for trading.

## Indicators Provided

The following list of indicators are currently supported by this package:

### Trend Indicators

- [Aroon Indicator](#aroon-indicator)
- [Chande Forecast Oscillator (CFO)](#chande-forecast-oscillator-cfo)
- [Double Exponential Moving Average (DEMA)](#double-exponential-moving-average-dema)
- [Exponential Moving Average (EMA)](#exponential-moving-average-ema)
- [Moving Average Convergence Divergence (MACD)](#moving-average-convergence-divergence-macd)
- [Parabolic SAR](#parabolic-sar)
- [Simple Moving Average (SMA)](#simple-moving-average-sma)
- [Triangular Moving Average (TRIMA)](#triangular-moving-average-trima)
- [Triple Exponential Moving Average (TEMA)](#triple-exponential-moving-average-tema)
- [Typical Price](#typical-price)
- [Vortex Indicator](#vortex-indicator)

### Momentum Indicators

- [Awesome Oscillator](#awesome-oscillator)
- [Ichimoku Cloud](#ichimoku-cloud)
- [Relative Strength Index (RSI)](#relative-strength-index-rsi)
- [Stochastic Oscillator](#stochastic-oscillator)
- [Williams R](#williams-r)

### Volatility Indicators

- [Acceleration Bands](#acceleration-bands)
- [Actual True Range (ATR)](#actual-true-range-atr)
- [Bollinger Band Width](#bollinger-band-width)
- [Bollinger Bands](#bollinger-bands)
- [Chandelier Exit](#chandelier-exit)
- [Moving Standard Deviation (Std)](#moving-standard-deviation-std)
- [Projection Oscillator (PO)](#projection-oscillator-po)

### Volume Indicators

- [Accumulation/Distribution (A/D)](#accumulationdistribution-ad)
- [On-Balance Volume (OBV)](#on-balance-volume-obv)

## Strategies Provided

The following list of strategies are currently supported by this package:

### Reference Strategies

- [Buy and Hold Strategy](#buy-and-hold-strategy)

### Trend Strategies

- [Chande Forecast Oscillator Strategy](#chande-forecast-oscillator-strategy)
- [MACD Strategy](#macd-strategy)
- [Trend Strategy](#trend-strategy)

### Momentum Strategies

- [Awesome Oscillator Strategy](#awesome-oscillator-strategy)
- [RSI Strategy](#rsi-strategy)
- [Williams R Strategy](#williams-r-strategy)

### Volatility Strategies

- [Bollinger Bands Strategy](#bollinger-bands-strategy)
- [Projection Oscillator Strategy](#projection-oscillator-strategy)

### Volume Strategies

### Compound Strategies

- [MACD and RSI Strategy](#macd-and-rsi-strategy)

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

## Indicators

### Trend Indicators

Trend indicators measure the direction and strength of a trend.

#### Aroon Indicator

The [Aroon](https://pkg.go.dev/github.com/cinar/indicator#Aroon) function calculates a technical indicator that is used to identify trend changes in the price of a stock, as well as the strength of that trend. It consists of two lines, Aroon Up, and Aroon Down. The Aroon Up line measures measures the strength of the uptrend, and the Aroon Down measures the strength of the downtrend. When Aroon Up is above Aroon Down, it indicates bullish price, and when Aroon Down is above Aroon Up, it indicates bearish price.

```
Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
```

```Golang
aroonUp, aroonDown := indicator.Aroon(high, low)
```

#### Chande Forecast Oscillator (CFO)

The [ChandeForecastOscillator](https://pkg.go.dev/github.com/cinar/indicator#ChandeForecastOscillator) developed by Tushar Chande The Forecast Oscillator plots the percentage difference between the closing price and the n-period linear regression forecasted price. The oscillator is above zero when the forecast price is greater than the closing price and less than zero if it is below.

```
R = Linreg(Closing)
CFO = ((Closing - R) / Closing) * 100
```

Based on [Chande Forecast Oscillator Formula, Strategy](https://www.stockmaniacs.net/chande-forecast-oscillator/), [Forecast Oscillator
](https://www.fmlabs.com/reference/default.htm?url=ForecastOscillator.htm), and [Least Squares Regression](https://www.mathsisfun.com/data/least-squares-regression.html).

```golang
cfo := indicator.ChandeForecastOscillator(closing)
```

#### Double Exponential Moving Average (DEMA)

The [Dema](https://pkg.go.dev/github.com/cinar/indicator#Dema) function calculates the Double Exponential Moving Average (DEMA) for a given period.

The double exponential moving average (DEMA) is a technical indicator introduced by Patrick Mulloy. The purpose is to reduce the amount of noise present in price charts used by technical traders. The DEMA uses two exponential moving averages (EMAs) to eliminate lag. It helps confirm uptrends when the price is above the average, and helps confirm downtrends when the price is below the average. When the price crosses the average that may signal a trend change.

```
DEMA = (2 * EMA(values)) - EMA(EMA(values))
```

```Golang
dema := indicator.Dema(period, values)
```

Based on [Double Exponential Moving Average (DEMA)](https://www.investopedia.com/terms/d/double-exponential-moving-average.asp).

#### Exponential Moving Average (EMA)

The [Ema](https://pkg.go.dev/github.com/cinar/indicator#Ema) function calculates the exponential moving average for a given period.

```Golang
result := indicator.Ema(2, []float64{2, 4, 6, 8, 12, 14, 16, 18, 20})
```

#### Moving Average Convergence Divergence (MACD)

The [Macd](https://pkg.go.dev/github.com/cinar/indicator#Macd) function calculates a trend-following momentum indicator that shows the relationship between two moving averages of price.

```
MACD = 12-Period EMA - 26-Period EMA.
Signal = 9-Period EMA of MACD.
```

```Golang
macd, signal := indicator.Macd(closing)
```

#### Parabolic SAR

The [ParabolicSar](https://pkg.go.dev/github.com/cinar/indicator#ParabolicSar) function calculates an identifier for the trend and the trailing stop.

```
PSAR = PSAR[i - 1] - ((PSAR[i - 1] - EP) * AF)
```

If the trend is Falling:

- PSAR is the maximum of PSAR or the previous two high values.
- If the current high is greather than or equals to PSAR, use EP.

If the trend is Rising:

- PSAR is the minimum of PSAR or the previous two low values.
- If the current low is less than or equials to PSAR, use EP.

If PSAR is greather than the closing, trend is falling, and the EP is set to the minimum of EP or the low.

If PSAR is lower than or equals to the closing, trend is rising, and the EP is set to the maximum of EP or the high.

If the trend is the same, and AF is less than 0.20, increment it by 0.02. If the trend is not the same, set AF to 0.02.

Based on video [How to Calculate the PSAR Using Excel - Revised Version](https://www.youtube.com/watch?v=MuEpGBAH7pw&t=0s).

```Golang
psar, trend := indicator.ParabolicSar(high, low, closing)
```

#### Simple Moving Average (SMA)

The [Sma](https://pkg.go.dev/github.com/cinar/indicator#Sma) function calculates the simple moving average for a given period.

```Golang
result := indicator.Sma(2, []float64{2, 4, 6, 8, 10})
```

#### Triangular Moving Average (TRIMA)

The [Trima](https://pkg.go.dev/github.com/cinar/indicator#Trima) function calculates the Triangular Moving Average (TRIMA) for a given period.

The Triangular Moving Average (TRIMA) is a weighted moving average putting more weight to the middle values.

```
If period is even:
   TRIMA = SMA(period / 2, SMA((period / 2) + 1, values))
If period is odd:
   TRIMA = SMA((period + 1) / 2, SMA((period + 1) / 2, values))
```

```Golang
trima := indicator.Trima(period, values)
```

Based on [Triangular Moving Average](https://tulipindicators.org/trima).

#### Triple Exponential Moving Average (TEMA)

The [Tema](https://pkg.go.dev/github.com/cinar/indicator#Tema) function calculates the Triple Exponential Moving Average (TEMA) for a given period.

The triple exponential moving average (TEMA) was designed to smooth value fluctuations, thereby making it easier to identify trends without the lag associated with traditional moving averages. It does this by taking multiple exponential moving averages (EMA) of the original EMA and subtracting out some of the lag.

```
TEMA = (3 * EMA1) - (3 * EMA2) + EMA3
EMA1 = EMA(values)
EMA2 = EMA(EMA1)
EMA3 = EMA(EMA2)
```

```Golang
tema := indicator.Tema(period, values)
```

Based on [Triple Exponential Moving Average (TEMA)](https://www.investopedia.com/terms/t/triple-exponential-moving-average.asp).

#### Typical Price

The [TypicalPrice](https://pkg.go.dev/github.com/cinar/indicator#TypicalPrice) function calculates another approximation of average price for each period and can be used as a filter for moving average systems.

```
Typical Price = (High + Low + Closing) / 3
```

```Golang
ta, sma20 := indicator.TypicalPrice(high, low, closing)
```

#### Vortex Indicator

The [Vortex](https://pkg.go.dev/github.com/cinar/indicator#Vortex) function provides two oscillators that capture positive and negative trend movement. A bullish signal triggers when the positive trend indicator crosses above the negative trend indicator or a key level. A bearish signal triggers when the negative trend indicator crosses above the positive trend indicator or a key level.

```
+VM = Abs(Current High - Prior Low)
-VM = Abs(Current Low - Prior High)

+VM14 = 14-Period Sum of +VM
-VM14 = 14-Period Sum of -VM

TR = Max((High[i]-Low[i]), Abs(High[i]-Closing[i-1]), Abs(Low[i]-Closing[i-1]))
TR14 = 14-Period Sum of TR

+VI14 = +VM14 / TR14
-VI14 = -VM14 / TR14
```

Based on [Vortex Indicator](https://school.stockcharts.com/doku.php?id=technical_indicators:vortex_indicator)

```Golang
plusVi, minusVi := indicator.Vortex(high, low, closing)
```

### Momentum Indicators

Momentum indicators measure the speed of movement.

#### Awesome Oscillator

The [AwesomeOscillator](https://pkg.go.dev/github.com/cinar/indicator#AwesomeOscillator) function calculates the awesome oscillator based on low and high daily prices for a given stock. It is an indicator used to measure market momentum.

```
Median Price = ((Low + High) / 2)
AO = 5-Period SMA - 34-Period SMA.
```

```Golang
result := indicator.AwesomeOscillator(low, high)
```

#### Ichimoku Cloud

The [IchimokuCloud](https://pkg.go.dev/github.com/cinar/indicator#IchimokuCloud), also known as Ichimoku Kinko Hyo, calculates a versatile indicator that defines support and resistence, identifies tred direction, gauges momentum, and provides trading signals.

```
Tenkan-sen (Conversion Line) = (9-Period High + 9-Period Low) / 2
Kijun-sen (Base Line) = (26-Period High + 26-Period Low) / 2
Senkou Span A (Leading Span A) = (Conversion Line + Base Line) / 2
Senkou Span B (Leading Span B) = (52-Period High + 52-Period Low) / 2
Chikou Span (Lagging Span) = Closing plotted 26 days in the past.
```

```Golang
conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingLine := indicator.IchimokuCloud(high, low, closing)
```

#### Relative Strength Index (RSI)

The [Rsi](https://pkg.go.dev/github.com/cinar/indicator#Rsi) function calculates a momentum indicator that measures the magnitude of recent price changes to evaluate overbought and oversold conditions.

```
RS = Average Gain / Average Loss
RSI = 100 - (100 / (1 + RS))
```

```Golang
rs, rsi := indicator.Rsi(closing)
```

#### Stochastic Oscillator

The [StochasticOscillator](https://pkg.go.dev/github.com/cinar/indicator#StochasticOscillator) function calculates a momentum indicator that shows the location of the closing relative to high-low range over a set number of periods.

```
K = (Closing - Lowest Low) / (Highest High - Lowest Low) * 100
D = 3-Period SMA of K
```

```Golang
k, d := indicator.StochasticOscillator(high, low, closing)
```

#### Williams R

The [WilliamsR](https://pkg.go.dev/github.com/cinar/indicator#WilliamsR) function calculates the Williams R based on low, high, and closing prices. It is a type of momentum indicator that moves between 0 and -100 and measures overbought and oversold levels.

```
WR = (Highest High - Closing) / (Highest High - Lowest Low)
```

```Golang
result := indicator.WilliamsR(low, high, closing)
```

### Volatility Indicators

Volatility indicators measure the rate of movement regardless of its direction.

#### Acceleration Bands

The [AccelerationBands](https://pkg.go.dev/github.com/cinar/indicator#AccelerationBands) plots upper and lower envelope bands around a simple moving average.

```
Upper Band = SMA(High * (1 + 4 * (High - Low) / (High + Low)))
Middle Band = SMA(Closing)
Lower Band = SMA(Low * (1 + 4 * (High - Low) / (High + Low)))
```

```golang
upperBand, middleBand, lowerBand := indicator.AccelerationBands(high, low, closing)
```

#### Actual True Range (ATR)

The [Atr](https://pkg.go.dev/github.com/cinar/indicator#Atr) function calculates a technical analysis indicator that measures market volatility by decomposing the entire range of stock prices for that period.

```
TR = Max((High - Low), (High - Closing), (Closing - Low))
ATR = 14-Period SMA TR
```

```Golang
tr, atr := indicator.Atr(14, high, low, closing)
```

#### Bollinger Bands

The [BollingerBands](https://pkg.go.dev/github.com/cinar/indicator#BollingerBands) function calculates the bollinger bands, middle band, upper band, lower band, provides identification of when a stock is oversold or overbought.

```
Middle Band = 20-Period SMA.
Upper Band = 20-Period SMA + 2 (20-Period Std)
Lower Band = 20-Period SMA - 2 (20-Period Std)
```

```Golang
middleBand, upperBand, lowerBand := indicator.BollingerBands(closing)
```

#### Bollinger Band Width

The [BollingerBandWidth](https://pkg.go.dev/github.com/cinar/indicator#BollingerBandWidth) function measures the percentage difference between the upper band and the lower band. It decreases as Bollinger Bands narrows and increases as Bollinger Bands widens.

During a period of rising price volatility the band width widens, and during a period of low market volatility band width contracts.

```
Band Width = (Upper Band - Lower Band) / Middle Band
```

```Golang
bandWidth, bandWidthEma90 := indicator.BollingerBandWidth(middleBand, upperBand, lowerBand)
```

#### Chandelier Exit

The [ChandelierExit](https://pkg.go.dev/github.com/cinar/indicator#ChandelierExit) function sets a trailing stop-loss based on the Average True Value (ATR).

```
Chandelier Exit Long = 22-Period SMA High - ATR(22) * 3
Chandelier Exit Short = 22-Period SMA Low + ATR(22) * 3
```

```Golang
chandelierExitLong, chandelierExitShort := indicator.ChandelierExit(high, low, closing)
```

#### Moving Standard Deviation (Std)

The [Std](https://pkg.go.dev/github.com/cinar/indicator#Std) function calculates the moving standard deviation for a given period.

```Golang
result := indicator.Std(2, []float64{2, 4, 6, 8, 12, 14, 16, 18, 20})
```

#### Projection Oscillator (PO)

The [ProjectionOscillator](https://pkg.go.dev/github.com/cinar/indicator#ProjectionOscillator) calculates the Projection Oscillator (PO). The PO uses the linear regression slope, along with highs and lows.

Period defines the moving window to calculates the PO, and the smooth period defines the moving windows to take EMA of PO.

```
PL = Min(period, (high + MLS(period, x, high)))
PU = Max(period, (low + MLS(period, x, low)))
PO = 100 * (Closing - PL) / (PU - PL)
SPO = EMA(smooth, PO)
```

```golang
po, spo := indicator.ProjectionOscillator(12, 4, high, low, closing)
```

### Volume Indicators

Volumne indicators measure the strength of a trend based the volume.

#### Accumulation/Distribution (A/D)

The [AccumulationDistribution](https://pkg.go.dev/github.com/cinar/indicator#AccumulationDistribution) is a cumulative indicator that uses volume and price to assess whether a stock is being accumulated or distributed.

The Accumulation/Distribution seeks to identify divergences between the stock price and the volume flow.

```
MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
MFV = MFM * Period Volume
AD = Previous AD + CMFV
```

Based on [Accumulation/Distribution Indicator (A/D)](https://www.investopedia.com/terms/a/accumulationdistribution.asp).

```golang
ad := indicator.AccumulationDistribution(high, low, closing)
```

#### On-Balance Volume (OBV)

The [Obv](https://pkg.go.dev/github.com/cinar/indicator#Obv) function calculates a technical trading momentum indicator that uses volume flow to predict changes in stock price.

```
                  volume, if Closing > Closing-Prev
OBV = OBV-Prev +       0, if Closing = Closing-Prev
                 -volume, if Closing < Closing-Prev
```

```Golang
result := indicator.Obv(closing, volume)
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
	HOLD Action = iota
	BUY
	SELL
)
```

### Reference Strategies

Reference strategies provide a base for comparing the actual outcome of a given strategy.

#### Buy and Hold Strategy

The [BuyAndHoldStrategy](https://pkg.go.dev/github.com/cinar/indicator#BuyAndHoldStrategy) provides a simple strategy to buy the given asset and hold it. It provides a good indicator for the change of asset's value without any other strategy is used.

```golang
actions := indicator.BuyAndHoldStrategy(asset)
```

### Trend Strategies

Trend strategies generate signals based on a trend indicator.

#### Chande Forecast Oscillator Strategy

The [ChandeForecastOscillatorStrategy](https://pkg.go.dev/github.com/cinar/indicator#ChandeForecastOscillatorStrategy) uses _cfo_ values that are generated by the [ChandeForecastOscillator](https://pkg.go.dev/github.com/cinar/indicator#ChandeForecastOscillator) indicator function to provide a BUY action when _cfo_ is below zero, and SELL action when _cfo_ is above zero.

```golang
actions := indicator.ChandeForecastOscillatorStrategy(asset)
```

#### MACD Strategy

The [MacdStrategy](https://pkg.go.dev/github.com/cinar/indicator#MacdStrategy) uses the _macd_, and _signal_ values that are generated by the [Macd](https://pkg.go.dev/github.com/cinar/indicator#Macd) indicator function to provide a BUY action when _macd_ crosses above _signal_, and SELL action when _macd_ crosses below _signal_.

```golang
actions := indicator.MacdStrategy(asset)
```

#### Trend Strategy

The [TrendStrategy](https://pkg.go.dev/github.com/cinar/indicator#TrendStrategy) provides a simply strategy to buy the given asset following the asset's closing value increases in _count_ subsequent rows. Produces the sell action following the asset's closing value decreases in _count_ subsequent rows.

```golang
actions := indicator.TrendStrategy(asset, 4)
```

The function signature of [TrendStrategy](https://pkg.go.dev/github.com/cinar/indicator#TrendStrategy) does not match the [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) type, as it requires an additional _count_ parameter. The [MakeTrendStrategy](https://pkg.go.dev/github.com/cinar/indicator#MakeTrendStrategy) function can be used to return a [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) instance based on the given _count_ value.

```golang
strategy := indicator.MakeTrendStrategy(4)
actions := strategy(asset)
```

### Momentum Strategies

Momentum strategies generate signals based on a momentum indicator.

#### Awesome Oscillator Strategy

The [AwesomeOscillatorStrategy](https://pkg.go.dev/github.com/cinar/indicator#AwesomeOscillatorStrategy) uses the _ao_ values that are generated by the [AwesomeOscillator](https://pkg.go.dev/github.com/cinar/indicator#AwesomeOscillator) indicator function to provide a SELL action when the _ao_ is below zero, and a BUY action when _ao_ is above zero.

```golang
actions := indicator.AwesomeOscillatorStrategy(asset)
```

#### RSI Strategy

The [RsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#RsiStrategy) uses the _rsi_ values that are generated by the [Rsi](https://pkg.go.dev/github.com/cinar/indicator#Rsi) indicator function to provide a BUY action when _rsi_ is below the _buyAt_ parameter, and a SELL action when _rsi_ is above the _sellAt_ parameter.

```golang
actions := indicator.RsiStrategy(asset, 70, 30)
```

The RSI strategy is usually used with 70-30, or 80-20 values. The [DefaultRsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#DefaultRsiStrategy) function uses the 70-30 values.

```golang
actions := indicator.DefaultRsiStrategy(asset)
```

The function signature of [RsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#RsiStrategy) does not match the [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) type, as it requires an additional _sellAt_, and _buyAt_ parameters. The [MakeRsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#MakeRsiStrategy) function can be used to return a [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) instance based on the given _sellAt_, and _buyAt_ values.

```golang
strategy := indicator.MakeRsiStrategy(80, 20)
actions := strategy(asset)
```

#### Williams R Strategy

The [WilliamsRStrategy](https://pkg.go.dev/github.com/cinar/indicator#WilliamsRStrategy) uses the _wr_ values that are generated by the [WilliamsR](https://pkg.go.dev/github.com/cinar/indicator#WilliamsR) indicator function to provide a SELL action when the _wr_ is below -20, and a BUY action when _wr_ is above -80.

```golang
actions := indicator.WilliamsRStrategy(asset)
```

### Volatility Strategies

Volatility strategies generate signals based on a volatility indicator.

#### Bollinger Bands Strategy

The [BollingerBandsStrategy](https://pkg.go.dev/github.com/cinar/indicator#BollingerBandsStrategy) uses the _upperBand_, and _lowerBand_ values that are generated by the [BollingerBands](https://pkg.go.dev/github.com/cinar/indicator#BollingerBands) indicator function to provide a SELL action when the asset's closing is above the _upperBand_, and a BUY action when the asset's closing is below the _lowerBand_ values.

```golang
actions := indicator.BollingerBandsStrategy(asset)
```

#### Projection Oscillator Strategy

The [ProjectionOscillatorStrategy](https://pkg.go.dev/github.com/cinar/indicator#ProjectionOscillatorStrategy) uses _po_ and _spo_ values that are generated by the [ProjectionOscillator](https://pkg.go.dev/github.com/cinar/indicator#ProjectionOscillator) indicator function to provide a BUY action when _po_ is above _spo_, and SELL action when _po_ is below _spo_.

```golang
actions := indicator.ProjectionOscillatorStrategy(period, smooth, asset)
```

### Volume Strategies

Volume strategies generate signals based on a volume indicator.

### Compound Strategies

Compound strategies combine multiple strategies together to generate signals.

#### MACD and RSI Strategy

The [MacdAndRsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#MacdAndRsiStrategy) function uses the actions generated by the [MacdStrategy](https://pkg.go.dev/github.com/cinar/indicator#MacdStrategy) and the [DefaultRsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#DefaultRsiStrategy) to provide BUY and SELL actions.

```golang
actions := indicator.MacdAndRsiStrategy(asset)
```

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.
The source code is provided under MIT License.
