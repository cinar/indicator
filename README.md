[![GoDoc](https://godoc.org/github.com/cinar/indicator?status.svg)](https://godoc.org/github.com/cinar/indicator)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.com/cinar/indicator.svg?branch=master)](https://travis-ci.com/cinar/indicator)

# Indicator Go

Indicator is a Golang module providing various stock technical analysis indicators for trading. The following list of indicators are currently supported by this package:

- [Simple Moving Average (SMA)](#simple-moving-average-sma)
- [Moving Standard Deviation (Std)](#moving-standard-deviation-std)
- [Exponential Moving Average (EMA)](#exponential-moving-average-ema)
- [Moving Average Convergence Divergence (MACD)](#moving-average-convergence-divergence-macd)
- [Bollinger Bands](#bollinger-bands)
- [Bollinger Band Width](#bollinger-band-width)
- [Awesome Oscillator](#awesome-oscillator)
- [Williams R](#williams-r)
- [Typical Price](#typical-price)
- [Relative Strength Index (RSI)](#relative-strength-index-rsi)
- [On-Balance Volume (OBV)](#on-balance-volume-obv)
- [Actual True Range (ATR)](#actual-true-range-atr)
- [Chandelier Exit](#chandelier-exit)
- [Ichimoku Cloud](#ichimoku-cloud)
- [Stochastic Oscillator](#stochastic-oscillator)
- [Aroon Indicator](#aroon-indicator)
- [Parabolic SAR](#parabolic-sar)
- [Vortex Indicator](#vortex-indicator)

## Usage

Install package.

```bash
go get github.com/cinar/indicator
```

Import indicator.

```Golang
import (
    "indicator"
)
```

### Moving Averages

#### Simple Moving Average (SMA)

The [Sma](https://pkg.go.dev/github.com/cinar/indicator#Sma) function calculates the simple moving average for a given period.

```Golang
result := indicator.Sma(2, []float64{2, 4, 6, 8, 10})
```

#### Moving Standard Deviation (Std)

The [Std](https://pkg.go.dev/github.com/cinar/indicator#Std) function calculates the moving standard deviation for a given period.

```Golang
result := indicator.Std(2, []float64{2, 4, 6, 8, 12, 14, 16, 18, 20})
```

#### Exponential Moving Average (EMA)

The [Ema](https://pkg.go.dev/github.com/cinar/indicator#Ema) function calculates the exponential moving average for a given period.

```Golang
result := indicator.Ema(2, []float64{2, 4, 6, 8, 12, 14, 16, 18, 20})
```

### Indicators

#### Moving Average Convergence Divergence (MACD)

The [Macd](https://pkg.go.dev/github.com/cinar/indicator#Macd) function calculates a trend-following momentum indicator that shows the relationship between two moving averages of price.

```
MACD = 12-Period EMA - 26-Period EMA.
Signal = 9-Period EMA of MACD.
```

```Golang
macd, signal := indicator.Macd(close)
```

#### Bollinger Bands

The [BollingerBands](https://pkg.go.dev/github.com/cinar/indicator#BollingerBands) function calculates the bollinger bands, middle band, upper band, lower band, provides identification of when a stock is oversold or overbought.

```
Middle Band = 20-Period SMA.
Upper Band = 20-Period SMA + 2 (20-Period Std)
Lower Band = 20-Period SMA - 2 (20-Period Std)
```

```Golang
middleBand, upperBand, lowerBand := indicator.BollingerBands(close)
```

#### Bollinger Band Width

The [BollingerBandWidth](https://pkg.go.dev/github.com/cinar/indicator#BollingerBandWidth) function  measures the percentage difference between the upper band and the lower band. It decreases as Bollinger Bands narrows and increases as Bollinger Bands widens.

During a period of rising price volatility the band width widens, and during a period of low market volatility band width contracts.

```
Band Width = (Upper Band - Lower Band) / Middle Band
```

```Golang
bandWidth, bandWidthEma90 := indicator.BollingerBandWidth(middleBand, upperBand, lowerBand)
```

#### Awesome Oscillator

The [AwesomeOscillator](https://pkg.go.dev/github.com/cinar/indicator#AwesomeOscillator) function calculates the awesome oscillator based on low and high daily prices for a given stock. It is an indicator used to measure market momentum. 

```
Median Price = ((Low + High) / 2)
AO = 5-Period SMA - 34-Period SMA.
```

```Golang
result := indicator.AwesomeOscillator(low, high)
```

#### Williams R

The [WilliamsR](https://pkg.go.dev/github.com/cinar/indicator#WilliamsR) function calculates the Williams R based on low, high, and close prices. It is a type of momentum indicator that moves between 0 and -100 and measures overbought and oversold levels.

```
WR = (Highest High - Close) / (Highest High - Lowest Low)
```

```Golang
result := indicator.WilliamsR(low, high, close)
```

#### Typical Price

The [TypicalPrice](https://pkg.go.dev/github.com/cinar/indicator#TypicalPrice) function calculates another approximation of average price for each period and can be used as a filter for moving average systems.

```
Typical Price = (High + Low + Close) / 3
```

```Golang
ta, sma20 := indicator.TypicalPrice(high, low, close)
```

#### Relative Strength Index (RSI)

The [Rsi](https://pkg.go.dev/github.com/cinar/indicator#Rsi) function calculates a momentum indicator that measures the magnitude of recent price changes to evaluate overbought and oversold conditions.

```
RS = Average Gain / Average Loss
RSI = 100 - (100 / (1 + RS))
```

```Golang
rs, rsi := indicator.Rsi(close)
```

#### On-Balance Volume (OBV)

The [Obv](https://pkg.go.dev/github.com/cinar/indicator#Obv) function calculates a technical trading momentum indicator that uses volume flow to predict changes in stock price.

```
                  volume, if Close > Close-Prev
OBV = OBV-Prev +       0, if Close = Close-Prev
                 -volume, if Close < Close-Prev
```

```Golang
result := indicator.Obv(close, volume)
```

#### Actual True Range (ATR)

The [Atr](https://pkg.go.dev/github.com/cinar/indicator#Atr) function calculates a technical analysis indicator that measures market volatility by decomposing the entire range of stock prices for that period.

```
TR = Max((High - Low), (High - Close), (Close - Low))
ATR = 14-Period SMA TR
```

```Golang
tr, atr := indicator.Atr(14, high, low, close)
```

#### Chandelier Exit

The [ChandelierExit](https://pkg.go.dev/github.com/cinar/indicator#ChandelierExit) function sets a trailing stop-loss based on the Average True Value (ATR).

```
Chandelier Exit Long = 22-Period SMA High - ATR(22) * 3
Chandelier Exit Short = 22-Period SMA Low + ATR(22) * 3
```

```Golang
chandelierExitLong, chandelierExitShort := indicator.ChandelierExit(high, low, close)
```

#### Ichimoku Cloud

The [IchimokuCloud](https://pkg.go.dev/github.com/cinar/indicator#IchimokuCloud), also known as Ichimoku Kinko Hyo, calculates a versatile indicator that defines support and resistence, identifies tred direction, gauges momentum, and provides trading signals.

```
Tenkan-sen (Conversion Line) = (9-Period High + 9-Period Low) / 2
Kijun-sen (Base Line) = (26-Period High + 26-Period Low) / 2
Senkou Span A (Leading Span A) = (Conversion Line + Base Line) / 2
Senkou Span B (Leading Span B) = (52-Period High + 52-Period Low) / 2
Chikou Span (Lagging Span) = Close plotted 26 days in the past.
```

```Golang
conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingLine := indicator.IchimokuCloud(high, low, close)
```

#### Stochastic Oscillator

The [StochasticOscillator](https://pkg.go.dev/github.com/cinar/indicator#StochasticOscillator) function calculates a momentum indicator that shows the location of the close relative to high-low range over a set number of periods.

```
K = (Close - Lowest Low) / (Highest High - Lowest Low) * 100
D = 3-Period SMA of K
```

```Golang
k, d := indicator.StochasticOscillator(high, low, close)
```

#### Aroon Indicator

The [Aroon](https://pkg.go.dev/github.com/cinar/indicator#Aroon) function calculates a technical indicator that is used to identify trend changes in the price of a stock, as well as the strength of that trend. It consists of two lines, Aroon Up, and Aroon Down. The Aroon Up line measures measures the strength of the uptrend, and the Aroon Down measures the strength of the downtrend. When Aroon Up is above Aroon Down, it indicates bullish price, and when Aroon Down is above Aroon Up, it indicates bearish price.

```
Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
```

```Golang
aroonUp, aroonDown := indicator.Aroon(high, low)
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
psar, trend := indicator.ParabolicSar(high, low, close)
```

#### Vortex Indicator

The [Vortex](https://pkg.go.dev/github.com/cinar/indicator#Vortex) function provides two oscillators that capture positive and negative trend movement. A bullish signal triggers when the positive trend indicator crosses above the negative trend indicator or a key level. A bearish signal triggers when the negative trend indicator crosses above the positive trend indicator or a key level.

```
+VM = Abs(Current High - Prior Low)
-VM = Abs(Current Low - Prior High)

+VM14 = 14-Period Sum of +VM
-VM14 = 14-Period Sum of -VM

TR = Max((High[i]-Low[i]), Abs(High[i]-Close[i-1]), Abs(Low[i]-Close[i-1]))
TR14 = 14-Period Sum of TR

+VI14 = +VM14 / TR14
-VI14 = -VM14 / TR14
```

Based on [Vortex Indicator](https://school.stockcharts.com/doku.php?id=technical_indicators:vortex_indicator)

```Golang
plusVi, minusVi := indicator.Vortex(high, low, close)
```

## License

The source code is provided under MIT License.
