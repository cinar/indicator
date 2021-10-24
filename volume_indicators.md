### Volume Indicators

Volumne indicators measure the strength of a trend based the volume.

- [Accumulation/Distribution (A/D)](#accumulationdistribution-ad)
- [On-Balance Volume (OBV)](#on-balance-volume-obv)

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

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
