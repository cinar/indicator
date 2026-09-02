// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package main

import (
	"fmt"
	"sort"

	"github.com/cinar/indicator/v2/examples/compound"
	"github.com/cinar/indicator/v2/examples/momentum"
	"github.com/cinar/indicator/v2/examples/trend"
	"github.com/cinar/indicator/v2/examples/volatility"
	"github.com/cinar/indicator/v2/examples/volume"
	"github.com/cinar/indicator/v2/strategy"
)

// StrategyBuilderFunc defines a function that builds a new strategy instance.
type StrategyBuilderFunc func() strategy.Strategy

// strategyBuilders provides the name to strategy builder mapping for the
// strategies that the indicator-backtest command line tool knows how to
// build. Strategies are only ever added to a backtest run when the user
// explicitly names them through the -strategies flag.
var strategyBuilders = map[string]StrategyBuilderFunc{
	"buy-and-hold": func() strategy.Strategy { return strategy.NewBuyAndHoldStrategy() },

	"macd-rsi": func() strategy.Strategy { return compound.NewMacdRsiStrategy() },

	"awesome-oscillator":    func() strategy.Strategy { return momentum.NewAwesomeOscillatorStrategy() },
	"coppock-curve":         func() strategy.Strategy { return momentum.NewCoppockCurveStrategy() },
	"elder-ray":             func() strategy.Strategy { return momentum.NewElderRayStrategy() },
	"ichimoku-cloud":        func() strategy.Strategy { return momentum.NewIchimokuCloudStrategy() },
	"rsi":                   func() strategy.Strategy { return momentum.NewRsiStrategy() },
	"stochastic-oscillator": func() strategy.Strategy { return momentum.NewStochasticOscillatorStrategy() },
	"stochastic-rsi":        func() strategy.Strategy { return momentum.NewStochasticRsiStrategy() },
	"triple-rsi":            func() strategy.Strategy { return momentum.NewTripleRsiStrategy() },
	"williams-r":            func() strategy.Strategy { return momentum.NewWilliamsRStrategy() },

	"alligator":                       func() strategy.Strategy { return trend.NewAlligatorStrategy() },
	"apo":                             func() strategy.Strategy { return trend.NewApoStrategy() },
	"aroon":                           func() strategy.Strategy { return trend.NewAroonStrategy() },
	"bop":                             func() strategy.Strategy { return trend.NewBopStrategy() },
	"cci":                             func() strategy.Strategy { return trend.NewCciStrategy() },
	"cfo":                             func() strategy.Strategy { return trend.NewCfoStrategy() },
	"dema":                            func() strategy.Strategy { return trend.NewDemaStrategy() },
	"golden-cross":                    func() strategy.Strategy { return trend.NewGoldenCrossStrategy() },
	"hma":                             func() strategy.Strategy { return trend.NewHmaStrategy() },
	"kama":                            func() strategy.Strategy { return trend.NewKamaStrategy() },
	"kdj":                             func() strategy.Strategy { return trend.NewKdjStrategy() },
	"macd":                            func() strategy.Strategy { return trend.NewMacdStrategy() },
	"qstick":                          func() strategy.Strategy { return trend.NewQstickStrategy() },
	"smma":                            func() strategy.Strategy { return trend.NewSmmaStrategy() },
	"trima":                           func() strategy.Strategy { return trend.NewTrimaStrategy() },
	"triple-moving-average-crossover": func() strategy.Strategy { return trend.NewTripleMovingAverageCrossoverStrategy() },
	"tsi":                             func() strategy.Strategy { return trend.NewTsiStrategy() },
	"vwma":                            func() strategy.Strategy { return trend.NewVwmaStrategy() },
	"weighted-close":                  func() strategy.Strategy { return trend.NewWeightedCloseStrategy() },

	"bollinger-bands":           func() strategy.Strategy { return volatility.NewBollingerBandsStrategy() },
	"donchian-channel-breakout": func() strategy.Strategy { return volatility.NewDonchianChannelBreakoutStrategy() },
	"keltner-channel":           func() strategy.Strategy { return volatility.NewKeltnerChannelStrategy() },
	"super-trend":               func() strategy.Strategy { return volatility.NewSuperTrendStrategy() },

	"chaikin-money-flow":     func() strategy.Strategy { return volume.NewChaikinMoneyFlowStrategy() },
	"ease-of-movement":       func() strategy.Strategy { return volume.NewEaseOfMovementStrategy() },
	"force-index":            func() strategy.Strategy { return volume.NewForceIndexStrategy() },
	"money-flow-index":       func() strategy.Strategy { return volume.NewMoneyFlowIndexStrategy() },
	"negative-volume-index":  func() strategy.Strategy { return volume.NewNegativeVolumeIndexStrategy() },
	"obv":                    func() strategy.Strategy { return volume.NewObvStrategy() },
	"weighted-average-price": func() strategy.Strategy { return volume.NewWeightedAveragePriceStrategy() },
}

// StrategyNames returns the sorted list of the registered strategy names.
func StrategyNames() []string {
	names := make([]string, 0, len(strategyBuilders))
	for name := range strategyBuilders {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// NewStrategy builds a new strategy instance for the given registered name.
func NewStrategy(name string) (strategy.Strategy, error) {
	builder, ok := strategyBuilders[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}

	return builder(), nil
}
