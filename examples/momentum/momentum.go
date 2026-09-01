// Package momentum provides illustrative examples demonstrating how to compose
// momentum indicators into trading strategies for educational and research purposes.
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
package momentum

import "github.com/cinar/indicator/v2/strategy"

// AllStrategies returns a slice containing references to all available example momentum strategies.
func AllStrategies() []strategy.Strategy {
	return []strategy.Strategy{
		NewAwesomeOscillatorStrategy(),
		NewElderRayStrategy(),
		NewCoppockCurveStrategy(),
		NewIchimokuCloudStrategy(),
		NewRsiStrategy(),
		NewStochasticOscillatorStrategy(),
		NewStochasticRsiStrategy(),
		NewTripleRsiStrategy(),
		NewWilliamsRStrategy(),
	}
}
