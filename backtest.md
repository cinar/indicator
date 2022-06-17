### Backtest

Backtesting is the method for seeing how well a strategy would have done. The following backtesting functions are provided for evaluating strategies.

- [Apply Actions](#apply-actions)
- [Count Transactions](#count-transactions)
- [Normalize Actions](#normalize-actions)
- [Normalize Gains](#normalize-gains)

#### Apply Actions

The [ApplyActions](https://pkg.go.dev/github.com/cinar/indicator#ApplyActions) takes the given list of prices, applies the given list of normalized actions, and returns the gains.

```golang
gains := indicator.ApplyActions(prices, actions)
```

#### Count Transactions

The [CountTransactions](https://pkg.go.dev/github.com/cinar/indicator#CountTransactions) takes a list of normalized actions, and counts the _BUY_ and _SELL_ actions.

```golang
count := indicator.CountTransactions(actions)
```

#### Normalize Actions

The [NormalizeActions](https://pkg.go.dev/github.com/cinar/indicator#NormalizeActions) takes a list of independenc actions, such as _SELL_, _SELL_, _BUY_, _SELL_, _HOLD_, _SELL_, and produces a normalized list where the actions are following the proper _BUY_, _HOLD_, _SELL_, _HOLD_ order.

```golang
normalized := indicator.NormalizeActions(actions)
```

#### Normalize Gains

The [NormalizeGains](https://pkg.go.dev/github.com/cinar/indicator#NormalizeGains) takes the given list of prices, calculates the price gains, subtracts it from the given list of gains.

```golang
normalizedGains := indicator.NormalizeGains(prices, gains)
```

## Disclaimer

The information provided on this project is strictly for informational purposes and is not to be construed as advice or solicitation to buy or sell any security.

## License

Copyright (c) 2021 Onur Cinar. All Rights Reserved.

The source code is provided under MIT License.
