### Strategies

The strategies are where the results from one or more indicators gets combined to produce a recommended action.

**The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.**

- [Asset](#asset)
- [Action](#action)
- [Strategy Function](#strategy-function)
- [Buy and Hold Strategy](#buy-and-hold-strategy)


#### Asset

The stragies operates on an [Asset](https://pkg.go.dev/github.com/cinar/indicator#Asset) with the following members.

```golang
type Asset struct {
	Date    []time.Time
	Opening []float64
	Closing []float64
	High    []float64
	Low     []float64
	Volume  []float64
}
```

#### Strategy Function

The [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) takes an [Asset](https://pkg.go.dev/github.com/cinar/indicator#Asset), and provides an array of [Action](https://pkg.go.dev/github.com/cinar/indicator#Action) for each row.

```golang
// Strategy function.
type StrategyFunction func(*Asset) []Action
```

#### Action

The following [Action](https://pkg.go.dev/github.com/cinar/indicator#Action) values are currently provided.

```golang
type Action int

const (
	SELL Action = -1
	HOLD Action = 0
	BUY  Action = 1
)
```

#### Buy and Hold Strategy

The [BuyAndHoldStrategy](https://pkg.go.dev/github.com/cinar/indicator#BuyAndHoldStrategy) provides a simple strategy to buy the given asset and hold it. It provides a good indicator for the change of asset's value without any other strategy is used.

```golang
actions := indicator.BuyAndHoldStrategy(asset)
```

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
