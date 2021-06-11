[![GoDoc](https://godoc.org/github.com/cinar/indicator?status.svg)](https://godoc.org/github.com/cinar/indicator)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.com/cinar/indicator.svg?branch=main)](https://travis-ci.com/cinar/indicator)

# Indicator Go

Indicator is a Golang module providing various stock technical analysis indicators for trading.

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

The [Macd](https://pkg.go.dev/github.com/cinar/indicator#Macd) function calculates the MACD and the Signal.

```
MACD = 12-Period EMA - 26-Period EMA.
Signal = 9-Period EMA of MACD.
```

```Golang
macd, signal := indicator.Macd(close)
```

#### Bollinger Bands

The [BollingerBands](https://pkg.go.dev/github.com/cinar/indicator#BollingerBands) function calculates the bollinger bands, middle band, upper band, lower band.

```
Middle Band = 20-Period SMA.
Upper Band = 20-Period SMA + 2 (20-Period Std)
Lower Band = 20-Period SMA - 2 (20-Period Std)
```

```Golang
middleBand, upperBand, lowerBand := indicator.BollingerBands(close)
```

#### Awesome Oscillator

The [AwesomeOscillator](https://pkg.go.dev/github.com/cinar/indicator#AwesomeOscillator) function calculates the awesome oscillator based on low and high daily prices for a given stock.

```
Median Price = ((Low + High) / 2)
AO = 5-Period SMA - 34-Period SMA.
```

```Golang
result := indicator.AwesomeOscillator(low, high)
```

#### Williams R

The [WilliamsR](https://pkg.go.dev/github.com/cinar/indicator#WilliamsR) function calculates the Williams R based on low, high, and close prices.

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

## License

The source code is provided under MIT License.
