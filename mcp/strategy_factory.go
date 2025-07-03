package main

import (
	"fmt"

	"github.com/cinar/indicator/v2/strategy"
	"github.com/cinar/indicator/v2/strategy/momentum"
	"github.com/cinar/indicator/v2/strategy/trend"
	"github.com/cinar/indicator/v2/strategy/volume"
)

// StrategyType defines the type of trading strategy to be used in a backtest.
// It is represented as a string to allow for easy identification and selection.
type StrategyType string

// Constants for all supported strategy types.
// This list includes a variety of strategies from different categories, such as
// trend, momentum, and volume-based approaches.
const (
	// Base strategies
	StrategyBuyAndHold StrategyType = "buy_and_hold"

	// Momentum strategies
	StrategyAwesomeOscillator StrategyType = "awesome_oscillator"
	StrategyRsi               StrategyType = "rsi"
	StrategyStochasticRsi     StrategyType = "stochastic_rsi"
	StrategyTripleRsi         StrategyType = "triple_rsi"

	// Volume strategies
	StrategyChaikinMoneyFlow     StrategyType = "chaikin_money_flow"
	StrategyEaseOfMovement       StrategyType = "ease_of_movement"
	StrategyForceIndex           StrategyType = "force_index"
	StrategyMoneyFlowIndex       StrategyType = "money_flow_index"
	StrategyNegativeVolumeIndex  StrategyType = "negative_volume_index"
	StrategyWeightedAveragePrice StrategyType = "weighted_average_price"

	// Trend strategies
	StrategyMACD              StrategyType = "macd"
	StrategyAlligator         StrategyType = "alligator"
	StrategyAroon             StrategyType = "aroon"
	StrategyApo               StrategyType = "apo"
	StrategyBop               StrategyType = "bop"
	StrategyCci               StrategyType = "cci"
	StrategyDema              StrategyType = "dema"
	StrategyGoldenCross       StrategyType = "golden_cross"
	StrategyKama              StrategyType = "kama"
	StrategyKdj               StrategyType = "kdj"
	StrategyQstick            StrategyType = "qstick"
	StrategySmma              StrategyType = "smma"
	StrategyTrima             StrategyType = "trima"
	StrategyTripleMaCrossover StrategyType = "triple_ma_crossover"
	StrategyTsi               StrategyType = "tsi"
	StrategyVwma              StrategyType = "vwma"
	StrategyWeightedClose     StrategyType = "weighted_close"
)

// CreateStrategy creates a new strategy instance based on the specified type.
// It acts as a factory function, mapping a StrategyType to a concrete
// implementation of the strategy.Strategy interface.
//
// This function is essential for dynamically selecting and initializing the
// desired trading strategy at runtime. If an unsupported strategy type is
// provided, it returns an error.
func CreateStrategy(strategyType StrategyType) (strategy.Strategy, error) {
	switch strategyType {
	// Base strategies
	case StrategyBuyAndHold:
		return strategy.NewBuyAndHoldStrategy(), nil

	// Trend strategies
	case StrategyAlligator:
		return trend.NewAlligatorStrategy(), nil
	case StrategyAroon:
		return trend.NewAroonStrategy(), nil
	case StrategyApo:
		return trend.NewApoStrategy(), nil
	case StrategyBop:
		return trend.NewBopStrategy(), nil
	case StrategyCci:
		return trend.NewCciStrategy(), nil
	case StrategyDema:
		return trend.NewDemaStrategy(), nil
	case StrategyGoldenCross:
		return trend.NewGoldenCrossStrategy(), nil
	case StrategyKama:
		return trend.NewKamaStrategy(), nil
	case StrategyKdj:
		return trend.NewKdjStrategy(), nil
	case StrategyMACD:
		return trend.NewMacdStrategy(), nil
	case StrategyQstick:
		return trend.NewQstickStrategy(), nil
	case StrategySmma:
		return trend.NewSmmaStrategy(), nil
	case StrategyTrima:
		return trend.NewTrimaStrategy(), nil
	case StrategyTripleMaCrossover:
		return trend.NewTripleMovingAverageCrossoverStrategy(), nil
	case StrategyTsi:
		return trend.NewTsiStrategy(), nil
	case StrategyVwma:
		return trend.NewVwmaStrategy(), nil
	case StrategyWeightedClose:
		return trend.NewWeightedCloseStrategy(), nil

	// Momentum strategies
	case StrategyAwesomeOscillator:
		return momentum.NewAwesomeOscillatorStrategy(), nil
	case StrategyRsi:
		return momentum.NewRsiStrategy(), nil
	case StrategyStochasticRsi:
		return momentum.NewStochasticRsiStrategy(), nil
	case StrategyTripleRsi:
		return momentum.NewTripleRsiStrategy(), nil

	// Volume strategies
	case StrategyChaikinMoneyFlow:
		return volume.NewChaikinMoneyFlowStrategy(), nil
	case StrategyEaseOfMovement:
		return volume.NewEaseOfMovementStrategy(), nil
	case StrategyForceIndex:
		return volume.NewForceIndexStrategy(), nil
	case StrategyMoneyFlowIndex:
		return volume.NewMoneyFlowIndexStrategy(), nil
	case StrategyNegativeVolumeIndex:
		return volume.NewNegativeVolumeIndexStrategy(), nil
	case StrategyWeightedAveragePrice:
		return volume.NewWeightedAveragePriceStrategy(), nil

	default:
		return nil, fmt.Errorf("unsupported strategy: %s", strategyType)
	}
}
