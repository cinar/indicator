### Compound Strategies

Compound strategies combine multiple strategies together to generate signals.

- [All Strategies](#all-strategies)
- [Run Strategies](#run-strategies)
- [Separate Strategies](#separate-strategies)
- [MACD and RSI Strategy](#macd-and-rsi-strategy)

#### All Strategies

The [AllStrategies](https://pkg.go.dev/github.com/cinar/indicator#AllStrategies) function takes one or more [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) and provides a [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) that will return a _BUY_ or _SELL_ action if all strategies are returning the same action, otherwise it will return a _HOLD_ action.

```golang
strategy := indicator.AllStrategies(indicator.MacdStrategy, indicator.DefaultRsiStrategy)
actions := strategy(asset)
```

#### Run Strategies

The [RunStrategies](https://pkg.go.dev/github.com/cinar/indicator#RunStrategies) function takes one or more [StrategyFunction](https://pkg.go.dev/github.com/cinar/indicator#StrategyFunction) and returns the actions for each.

```golang
actions := indicator.RunStrategies(asset, indicator.MacdStrategy, indicator.DefaultRsiStrategy)
```

#### Separate Strategies

The [SeparateStrategies](https://pkg.go.dev/github.com/cinar/indicator#SeparateStrategies) function takes a buy strategy and a sell strategy.

It returns a _BUY_ action if the buy strategy returns a _BUY_ action and the the sell strategy returns a _HOLD_ action.

It returns a _SELL_ action if the sell strategy returns a _SELL_ action and the buy strategy returns a _HOLD_ action.

It returns _HOLD_ otherwise.

```golang
strategy := indicator.SeparateStrategies(indicator.MacdStrategy, indicator.DefaultRsiStrategy)
actions := strategy(asset)
```

#### MACD and RSI Strategy

The [MacdAndRsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#MacdAndRsiStrategy) function is a compound strategy that combines the [MacdStrategy](https://pkg.go.dev/github.com/cinar/indicator#MacdStrategy) and the [DefaultRsiStrategy](https://pkg.go.dev/github.com/cinar/indicator#DefaultRsiStrategy). It will return a _BUY_ or _SELL_ action if both strategies are turning the same action, otherwise, it will return a _HOLD_ action.

```golang
actions := indicator.MacdAndRsiStrategy(asset)
```

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
