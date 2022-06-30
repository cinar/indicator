### Trend Indicators

Trend indicators measure the direction and strength of a trend.

- [Absolute Price Oscillator (APO)](#absolute-price-oscillator-apo)
- [Aroon Indicator](#aroon-indicator)
- [Balance of Power (BOP)](trend_indicators.md#balance-of-power-bop)
- [Chande Forecast Oscillator (CFO)](#chande-forecast-oscillator-cfo)
- [Community Channel Index (CMI)](#community-channel-index-cmi)
- [Double Exponential Moving Average (DEMA)](#double-exponential-moving-average-dema)
- [Exponential Moving Average (EMA)](#exponential-moving-average-ema)
- [Mass Index (MI)](#mass-index-mi)
- [Moving Average Convergence Divergence (MACD)](#moving-average-convergence-divergence-macd)
- [Moving Max](#moving-max)
- [Moving Min](#moving-min)
- [Moving Sum](#moving-sum)
- [Parabolic SAR](#parabolic-sar)
- [Qstick](trend_indicator.md#qstick)
- [Random Index (KDJ)](#random-index-kdj)
- [Rolling Moving Average (RMA)](#rolling-moving-average-rma)
- [Simple Moving Average (SMA)](#simple-moving-average-sma)
- [Since Change](#since-change)
- [Triple Exponential Moving Average (TEMA)](#triple-exponential-moving-average-tema)
- [Triangular Moving Average (TRIMA)](#triangular-moving-average-trima)
- [Triple Exponential Average (TRIX)](#triple-exponential-average-trix)
- [Typical Price](#typical-price)
- [Volume Weighted Moving Average (VWMA)](#volume-weighted-moving-average-vwma)
- [Vortex Indicator](#vortex-indicator)

#### Absolute Price Oscillator (APO)

The [AbsolutePriceOscillator](https://pkg.go.dev/github.com/cinar/indicator#AbsolutePriceOscillator) function calculates a technical indicator that is used to follow trends. APO crossing above zero indicates bullish, while crossing below zero indicates bearish. Positive value is upward trend, while negative value is downward trend.

```
Fast = Ema(fastPeriod, values)
Slow = Ema(slowPeriod, values)
APO = Fast - Slow
```

```Golang
apo := indicator.AbsolutePriceOscillator(fastPeriod, slowPeriod, values)
```

Most frequently used fast and short periods are 14 and 30. The [DefaultAbsoluePriceOscillator](https://pkg.go.dev/github.com/cinar/indicator#DefaultAbsolutePriceOscillator) function calculates APO with those periods.

#### Aroon Indicator

The [Aroon](https://pkg.go.dev/github.com/cinar/indicator#Aroon) function calculates a technical indicator that is used to identify trend changes in the price of a stock, as well as the strength of that trend. It consists of two lines, Aroon Up, and Aroon Down. The Aroon Up line measures measures the strength of the uptrend, and the Aroon Down measures the strength of the downtrend. When Aroon Up is above Aroon Down, it indicates bullish price, and when Aroon Down is above Aroon Up, it indicates bearish price.

```
Aroon Up = ((25 - Period Since Last 25 Period High) / 25) * 100
Aroon Down = ((25 - Period Since Last 25 Period Low) / 25) * 100
```

```Golang
aroonUp, aroonDown := indicator.Aroon(high, low)
```

#### Balance of Power (BOP)

The [BalanceOfPower](https://pkg.go.dev/github.com/cinar/indicator#BalanceOfPower) function calculates the strength of buying and selling pressure. Positive value indicates an upward trend, and negative value indicates a downward trend. Zero indicates a balance between the two.

```
BOP = (Closing - Opening) / (High - Low)
```

```Golang
bop := indicator.BalanceOfPower(opening, high, low, closing)
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
#### Community Channel Index (CMI)

The [CommunityChannelIndex](https://pkg.go.dev/github.com/cinar/indicator#CommunityChannelIndex) is a momentum-based oscillator used to help determine when an investment vehicle is reaching a condition of being overbought or oversold.

```
Moving Average = Sma(Period, Typical Price)
Mean Deviation = Sma(Period, Abs(Typical Price - Moving Average))
CMI = (Typical Price - Moving Average) / (0.015 * Mean Deviation)
```

```golang
result := indicator.CommunityChannelIndex(period, high, low, closing)
```

The [DefaultCommunityChannelIndex](https://pkg.go.dev/github.com/cinar/indicator#DefaultCommunityChannelIndex) calculates with the period of 20.

```golang
result := indicator.DefaultCommunityChannelIndex(high, low, closing)
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

#### Mass Index (MI)

The [MassIndex](https://pkg.go.dev/github.com/cinar/indicator#MassIndex) uses the high-low range to identify trend reversals based on range expansions.

```
Singe EMA = EMA(9, Highs - Lows)
Double EMA = EMA(9, Single EMA)
Ratio = Single EMA / Double EMA
MI = Sum(25, Ratio)
```

```Golang
result := indicator.MassIndex(high, low)
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

#### Moving Max

The [Max](https://pkg.go.dev/github.com/cinar/indicator#Max) function gives the maximum value within the given moving period. It can be used to get the moving maximum closing price and other values.

```Golang
max := indicator.Max(period, values)
```

#### Moving Min

The [Min](https://pkg.go.dev/github.com/cinar/indicator#Min) function gives the minimum value within the given moving period. It can be used to get the moving minimum closing price and other values.

```Golang
max := indicator.Min(period, values)
```

#### Moving Sum

The [Sum](https://pkg.go.dev/github.com/cinar/indicator#Sum) function gives the sum value within the given moving period.

```Golang
sum := indicator.Sum(period, values)
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

#### Qstick

The [Qstick](https://pkg.go.dev/github.com/cinar/indicator#Qstick) function calculates the ratio of recent up and down bars.

```
QS = Sma(Closing - Opening)
```

```Golang
qs := indicator.Qstick(period, closing, opening)
```

#### Random Index (KDJ)

The [Kdj](https://pkg.go.dev/github.com/cinar/indicator#Kdj) function calculates the KDJ  indicator, also known as the Random Index. KDJ is calculated similar to the Stochastic Oscillator with the difference of having the J line. It is used to analyze the trend and entry points.

The K and D lines show if the asset is overbought when they crosses above 80%, and oversold when they crosses below 20%. The J line represents the divergence.

```
RSV = ((Closing - Min(Low, rPeriod)) / (Max(High, rPeriod) - Min(Low, rPeriod))) * 100
K = Sma(RSV, kPeriod)
D = Sma(K, dPeriod)
J = (3 * K) - (2 * D)
```

```Golang
k, d, j := indicator.Kdj(rPeriod, kPeriod, dPeriod, high, low, closing)
```

By default, _rPeriod_ of 9, _kPeriod_ of 3, and _dPeriod_ of 3 are used. The [DefaultKdj](https://pkg.go.dev/github.com/cinar/indicator#DefaultKdj) function can be used with those periods.

```Golang
k, d, j := indicator.DefaultKdj(high, low, closing)
```

#### Rolling Moving Average (RMA)

The [Rma](https://pkg.go.dev/github.com/cinar/indicator#Rma) function calculates the rolling moving average for a given period.

```
R[0] to R[p-1] is SMA(values)
R[p] and after is R[i] = ((R[i-1]*(p-1)) + v[i]) / p
```

```Golang
result := indicator.Rma(2, []float64{2, 4, 6, 8, 10, 12})
```

#### Simple Moving Average (SMA)

The [Sma](https://pkg.go.dev/github.com/cinar/indicator#Sma) function calculates the simple moving average for a given period.

```Golang
result := indicator.Sma(2, []float64{2, 4, 6, 8, 10})
```

#### Since Change

The [Since](https://pkg.go.dev/github.com/cinar/indicator#Since) function provides the number values since the last change.

```Golang
changes := indicator.Since(values)
```

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

#### Triple Exponential Average (TRIX)

The [Trix](https://pkg.go.dev/github.com/cinar/indicator#Trix) indicator is an oscillator used to identify oversold and overbought markets, and it can also be used as a momentum indicator. Like many oscillators, TRIX oscillates around a zero line.

```
EMA1 = EMA(period, values)
EMA2 = EMA(period, EMA1)
EMA3 = EMA(period, EMA2)
TRIX = (EMA3 - Previous EMA3) / Previous EMA3
```

```Golang
trix := indicator.Trix(period, values)
```

#### Typical Price

The [TypicalPrice](https://pkg.go.dev/github.com/cinar/indicator#TypicalPrice) function calculates another approximation of average price for each period and can be used as a filter for moving average systems.

```
Typical Price = (High + Low + Closing) / 3
```

```Golang
ta, sma20 := indicator.TypicalPrice(high, low, closing)
```
#### Volume Weighted Moving Average (VWMA)

The [Vwma](https://pkg.go.dev/github.com/cinar/indicator#Vwma) function calculates the Volume Weighted Moving Average (VWMA) averaging the price data with an emphasis on volume, meaning areas with higher volume will have a greater weight.

```
VWMA = Sum(Price * Volume) / Sum(Volume) for a given Period.
```

```Golang
vwma := indicator.Vwma(period, closing, volume)
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

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
