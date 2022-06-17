// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

// The NormalizeActions takes a list of independenc actions,
// such as SELL, SELL, BUY, SELL, HOLD, SELL, and produces
// a normalized list where the actions are following the
// proper BUY, HOLD, SELL, HOLD order.
func NormalizeActions(actions []Action) []Action {
	normalized := make([]Action, len(actions))

	last := SELL

	for i, action := range actions {
		if (action != HOLD) && (action != last) {
			normalized[i] = action
			last = action
		} else {
			normalized[i] = HOLD
		}
	}

	return normalized
}

// The CountTransactions takes a list of normalized actions,
// and counts the BUY and SELL actions.
func CountTransactions(actions []Action) int {
	count := 0

	for _, action := range actions {
		if (action == BUY) || (action == SELL) {
			count++
		}
	}

	return count
}

// The ApplyActions takes the given list of prices, applies the
// given list of normalized actions, and returns the gains.
func ApplyActions(prices []float64, actions []Action) []float64 {
	gains := make([]float64, len(actions))

	initialBalance := 1.0
	balance := initialBalance
	shares := 0.0

	for i := 0; i < len(actions); i++ {
		if actions[i] == BUY {
			if balance > 0 {
				shares = balance / prices[i]
				balance = 0
			}
		} else if actions[i] == SELL {
			if shares > 0 {
				balance = shares * prices[i]
				shares = 0
			}
		}

		gains[i] = ((shares * prices[i]) + balance - initialBalance) / initialBalance
	}

	return gains
}

// The NormalizeGains takes the given list of prices, calculates the
// price gains, subtracts it from the given list of gains.
func NormalizeGains(prices, gains []float64) []float64 {
	priceGains := Sum(len(prices), percentDiff(prices, 1))
	normalized := subtract(gains, priceGains)

	return normalized
}
