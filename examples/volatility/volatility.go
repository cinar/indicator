// Package volatility provides illustrative examples demonstrating how to compose
// volatility indicators into trading strategies for educational and research purposes.
//
// This package belongs to the Indicator project. Indicator is
// a Golang module that supplies a variety of technical
// indicators, strategies, and a backtesting framework
// for analysis.
//
// # License
//
//	Copyright (c) 2021-2026 The Indicator Authors.
//	The source code is provided under GNU AGPLv3 License.
//	https://github.com/cinar/indicator
//
// # Disclaimer
//
// The information provided on this project is strictly for
// informational and educational purposes and is not to be construed as
// investment, financial, or trading advice.
package volatility

import (
	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/trend"
	"github.com/cinar/indicator/v2/volatility"
)

// AllStrategies returns a slice containing references to all available example volatility strategies.
func AllStrategies() []strategy.Strategy {
	return []strategy.Strategy{
		NewBollingerBandsStrategy(),
		NewKeltnerChannelStrategy(),
		NewDonchianChannelBreakoutStrategy(),
		NewSuperTrendStrategy(),
		NewSuperTrendStrategyWith(
			volatility.NewSuperTrendWithMa[float64](
				trend.NewSmaWithPeriod[float64](volatility.DefaultSuperTrendPeriod),
				volatility.DefaultSuperTrendMultiplier,
			),
		),
		NewSuperTrendStrategyWith(
			volatility.NewSuperTrendWithMa[float64](
				trend.NewEmaWithPeriod[float64](volatility.DefaultSuperTrendPeriod),
				volatility.DefaultSuperTrendMultiplier,
			),
		),
		NewSuperTrendStrategyWith(volatility.NewSuperTrendWithPeriod[float64](14, 3)),
		NewSuperTrendStrategyWith(volatility.NewSuperTrendWithPeriod[float64](10, 3)),
		NewSuperTrendStrategyWith(volatility.NewSuperTrendWithPeriod[float64](7, 3)),
	}
}
