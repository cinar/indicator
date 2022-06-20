### Momentum Indicators

Momentum indicators measure the speed of movement.

- [Awesome Oscillator](#awesome-oscillator)
- [Chaikin Oscillator](#chaikin-oscillator)
- [Ichimoku Cloud](#ichimoku-cloud)
- [Percentage Price Oscillator (PPO)](#percentage-price-oscillator-ppo)
- [Percentage Volume Oscillator (PVO)](#percentage-volume-oscillator-pvo)
- [Relative Strength Index (RSI)](#relative-strength-index-rsi)
- [RSI 2](#rsi-2)
- [RSI Period](#rsi-period)
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

#### Chaikin Oscillator

The [ChaikinOscillator](https://pkg.go.dev/github.com/cinar/indicator#ChaikinOscillator) function measures the momentum of the [Accumulation/Distribution (A/D)](volume_indicators.md#accumulationdistribution-ad) using the [Moving Average Convergence Divergence (MACD)](trend_indicators.md#moving-average-convergence-divergence-macd) formula. It takes the difference between fast and slow periods EMA of the A/D. Cross above the A/D line indicates bullish.

```
CO = Ema(fastPeriod, AD) - Ema(slowPeriod, AD)
```

```Golang
co, ad := indicator.ChaikinOscillator(fastPeriod, slowPeriod, high, low, closing)
```

Most frequently used fast and short periods are 3 and 10. The [DefaultChaikinOscillator](https://pkg.go.dev/github.com/cinar/indicator#DefaultChaikinOscillator) function calculates Chaikin Oscillator with those periods.

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

#### Percentage Price Oscillator (PPO)

The [PercentagePriceOscillator](https://pkg.go.dev/github.com/cinar/indicator#PercentagePriceOscillator) function calculates a momentum oscillator for the price It is used to indicate the ups and downs based on the price. A breakout is confirmed when PPO is positive.

```
PPO = ((EMA(fastPeriod, prices) - EMA(slowPeriod, prices)) / EMA(longPeriod, prices)) * 100
Signal = EMA(9, PPO)
Histogram = PPO - Signal
```

```Golang
ppo, signal, histogram := indicator.PercentagePriceOscillator(
    fastPeriod, 
    slowPeriod, 
    signalPeriod, 
    price
)
```

The [DefaultPercentagePriceOscillator](https://pkg.go.dev/github.com/cinar/indicator#DefaultPercentagePriceOscillator) function calculates it with the default periods of 12, 26, 9.

```Golang
ppo, signal, histogram := indicator.DefaultPercentagePriceOscillator(price)
```

#### Percentage Volume Oscillator (PVO)

The [PercentageVolumeOscillator](https://pkg.go.dev/github.com/cinar/indicator#PercentageVolumeOscillator) function calculates a momentum oscillator for the volume It is used to indicate the ups and downs based on the volume. A breakout is confirmed when PVO is positive.

```
PVO = ((EMA(fastPeriod, volumes) - EMA(slowPeriod, volumes)) / EMA(longPeriod, volumes)) * 100
Signal = EMA(9, PVO)
Histogram = PVO - Signal
```

```Golang
pvo, signal, histogram := indicator.PercentageVolumeOscillator(
    fastPeriod, 
    slowPeriod, 
    signalPeriod, 
    volume
)
```

The [DefaultPercentageVolumeOscillator](https://pkg.go.dev/github.com/cinar/indicator#DefaultPercentageVolumeOscillator) function calculates it with the default periods of 12, 26, 9.

```Golang
pvo, signal, histogram := indicator.DefaultPercentageVolumeOscillator(volume)
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

#### RSI 2

The [Rsi2](https://pkg.go.dev/github.com/cinar/indicator#Rsi2) function calculates a calculates a RSI with 2 period that provides a mean-reversion trading strategy. It is developed by Larry Connors.

```Golang
rs, rsi := indicator.Rsi2(closing)
```

#### RSI Period

The [RsiPeriod](https://pkg.go.dev/github.com/cinar/indicator#RsiPeriod) allows to calculate the RSI indicator with a non-standard period.

```Golang
rs, rsi := indicator.RsiPeriod(period, closing)
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
