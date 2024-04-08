// Package volatility contains the volatility strategy functions.
//
// This package belongs to the Indicator project. Indicator is
// a Golang module that supplies a variety of technical
// indicators, strategies, and a backtesting framework
// for analysis.
//
// # License
//
//	Copyright (c) 2021-2024 Onur Cinar.
//	The source code is provided under GNU AGPLv3 License.
//	https://github.com/cinar/indicator/v2
//
// # Disclaimer
//
// The information provided on this project is strictly for
// informational purposes and is not to be construed as
// advice or solicitation to buy or sell any security.
package volatility

import "github.com/cinar/indicator/v2/strategy"

// AllStrategies returns a slice containing references to all available volatility strategies.
func AllStrategies() []strategy.Strategy {
	return []strategy.Strategy{
		NewBollingerBandsStrategy(),
		NewSuperTrendStrategy(),
		NewSuperTrendStrategyWith(14, 3),
		NewSuperTrendStrategyWith(10, 3),
		NewSuperTrendStrategyWith(7, 3),
	}
}
