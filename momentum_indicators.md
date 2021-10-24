### Momentum Indicators

Momentum indicators measure the speed of movement.

- [Awesome Oscillator](#awesome-oscillator)
- [Ichimoku Cloud](#ichimoku-cloud)
- [Relative Strength Index (RSI)](#relative-strength-index-rsi)
- [Stochastic Oscillator](#stochastic-oscillator)
- [Williams R](#williams-r)

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

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
