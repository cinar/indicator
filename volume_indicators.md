### Volume Indicators

Volumne indicators measure the strength of a trend based the volume.

- [Accumulation/Distribution (A/D)](#accumulationdistribution-ad)
- [Chaikin Money Flow (CMF)](#chaikin-money-flow-cmf)
- [Ease of Movement (EMV)](#ease-of-movement-emv)
- [Force Index (FI)](#force-index-fi)
- [Money Flow Index (MFI)](#money-flow-index-mfi)
- [Negative Volume Index (NVI)](#negative-volume-index-nvi)
- [On-Balance Volume (OBV)](#on-balance-volume-obv)
- [Volume Price Trend (VPT)](#volume-price-trend-vpt)
- [Volume Weighted Average Price (VWAP)](#volume-weighted-average-price-vwap)

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
ad := indicator.AccumulationDistribution(high, low, closing, volume)
```

#### Chaikin Money Flow (CMF)

The [ChaikinMoneyFlow](https://pkg.go.dev/github.com/cinar/indicator#ChaikinMoneyFlow) measures the amount of money flow volume over a given period.

```
Money Flow Multiplier = ((Closing - Low) - (High - Closing)) / (High - Low)
Money Flow Volume = Money Flow Multiplier * Volume
Chaikin Money Flow = Sum(20, Money Flow Volume) / Sum(20, Volume)
```

```Golang
cmf := indicator.ChaikinMoneyFlow(high, low, closing, volume)
```

#### Ease of Movement (EMV)

The [EaseOfMovement](https://pkg.go.dev/github.com/cinar/indicator#EaseOfMovement) is a volume based oscillator measuring the ease of price movement.

```
Distance Moved = ((High + Low) / 2) - ((Priod High + Prior Low) /2)
Box Ratio = ((Volume / 100000000) / (High - Low))
EMV(1) = Distance Moved / Box Ratio
EMV(14) = SMA(14, EMV(1))
```

```Golang
emv := indicator.EaseOfMovement(period, high, low, volume)
```

The [DefaultEaseOfMovement](https://pkg.go.dev/github.com/cinar/indicator#DefaultEaseOfMovement) functio uses the default period of 14.

```Golang
emv := indicator.DefaultEaseOfMovement(high, low, volume)
```

#### Force Index (FI)

The [ForceIndex](https://pkg.go.dev/github.com/cinar/indicator#ForceIndex) uses the closing price and the volume to assess the power behind a move and identify turning points.

```
Force Index = EMA(period, (Current - Previous) * Volume)
```

```Golang
fi := indicator.ForceIndex(period, closing, volume)
```

The [DefaultForceIndex](https://pkg.go.dev/github.com/cinar/indicator#DefaultForceIndex) function uses the default period of 13.

```Golang
fi := DefaultForceIndex(closing, volume)
```

#### Money Flow Index (MFI)

The [MoneyFlowIndex](https://pkg.go.dev/github.com/cinar/indicator#MoneyFlowIndex) function analyzes both the closing price and the volume to measure to identify overbought and oversold states. It is similar to the Relative Strength Index (RSI), but it also uses the volume.

```
Raw Money Flow = Typical Price * Volume
Money Ratio = Positive Money Flow / Negative Money Flow
Money Flow Index = 100 - (100 / (1 + Money Ratio))
```

```Golang
result := indicator.MoneyFlowIndex(period, high, low, closing, volume)
```

The [DefaultMoneyFlowIndex](https://pkg.go.dev/github.com/cinar/indicator#DefaultMoneyFlowIndex) function uses the default period of 14.

#### Negative Volume Index (NVI)

The [NegativeVolumeIndex](https://pkg.go.dev/github.com/cinar/indicator#NegativeVolumeIndex) function calculates a cumulative indicator using the change in volume to decide when the smart money is active.

```
If Volume is greather than Previous Volume:

    NVI = Previous NVI

Otherwise:

    NVI = Previous NVI + (((Closing - Previous Closing) / Previous Closing) * Previous NVI)
```

```Golang
result := indicator.NegativeVolumeIndex(closing, volume)
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

#### Volume Price Trend (VPT)

The [VolumePriceTrend](https://pkg.go.dev/github.com/cinar/indicator#VolumePriceTrend) provides a correlation between the volume and the price.

```
VPT = Previous VPT + (Volume * (Current Closing - Previous Closing) / Previous Closing)
```

```Golang
result := indicator.VolumePriceTrend(closing, volume)
```

#### Volume Weighted Average Price (VWAP)

The [VolumeWeightedAveragePrice](https://pkg.go.dev/github.com/cinar/indicator#VolumeWeightedAveragePrice) provides the average price the asset has traded.

```
VWAP = Sum(Closing * Volume) / Sum(Volume)
```

```Golang
result := indicator.VolumeWeightedAveragePrice(period, closing, volume)
```

The [DefaultVolumeWeightedAveragePrice](https://pkg.go.dev/github.com/cinar/indicator#DefaultVolumeWeightedAveragePrice) function uses the default period of 14.

```Golang
result := indicator.DefaultVolumeWeightedAveragePrice(closing, volume)
```

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
