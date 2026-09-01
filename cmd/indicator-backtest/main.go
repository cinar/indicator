// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// main is the indicator backtest command line program.
package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/examples/compound"
	"github.com/cinar/indicator/v2/examples/momentum"
	"github.com/cinar/indicator/v2/examples/trend"
	"github.com/cinar/indicator/v2/examples/volatility"
	"github.com/cinar/indicator/v2/examples/volume"
	"github.com/cinar/indicator/v2/strategy"
)

func parseStrategies(input string) ([]strategy.Strategy, error) {
	var result []strategy.Strategy
	tokens := strings.Split(input, ",")
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		strats, err := resolveStrategy(token)
		if err != nil {
			return nil, err
		}
		result = append(result, strats...)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no strategies specified")
	}

	return result, nil
}

func resolveStrategy(name string) ([]strategy.Strategy, error) {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "-", "_"))
	normalized = strings.TrimSuffix(normalized, "_strategy")

	switch normalized {
	case "all":
		var all []strategy.Strategy
		all = append(all, strategy.AllStrategies()...)
		all = append(all, compound.AllStrategies()...)
		all = append(all, momentum.AllStrategies()...)
		all = append(all, trend.AllStrategies()...)
		all = append(all, volatility.AllStrategies()...)
		all = append(all, volume.AllStrategies()...)
		return all, nil

	case "compound":
		return compound.AllStrategies(), nil

	case "momentum":
		return momentum.AllStrategies(), nil

	case "trend":
		return trend.AllStrategies(), nil

	case "volatility":
		return volatility.AllStrategies(), nil

	case "volume":
		return volume.AllStrategies(), nil

	// Base strategies
	case "buy_and_hold":
		return []strategy.Strategy{strategy.NewBuyAndHoldStrategy()}, nil

	// Trend strategies
	case "alligator":
		return []strategy.Strategy{trend.NewAlligatorStrategy()}, nil
	case "apo":
		return []strategy.Strategy{trend.NewApoStrategy()}, nil
	case "aroon":
		return []strategy.Strategy{trend.NewAroonStrategy()}, nil
	case "bop":
		return []strategy.Strategy{trend.NewBopStrategy()}, nil
	case "cci":
		return []strategy.Strategy{trend.NewCciStrategy()}, nil
	case "cfo":
		return []strategy.Strategy{trend.NewCfoStrategy()}, nil
	case "dema":
		return []strategy.Strategy{trend.NewDemaStrategy()}, nil
	case "envelope":
		return []strategy.Strategy{trend.NewEnvelopeStrategy()}, nil
	case "golden_cross":
		return []strategy.Strategy{trend.NewGoldenCrossStrategy()}, nil
	case "hma":
		return []strategy.Strategy{trend.NewHmaStrategy()}, nil
	case "kama":
		return []strategy.Strategy{trend.NewKamaStrategy()}, nil
	case "kdj":
		return []strategy.Strategy{trend.NewKdjStrategy()}, nil
	case "macd":
		return []strategy.Strategy{trend.NewMacdStrategy()}, nil
	case "qstick":
		return []strategy.Strategy{trend.NewQstickStrategy()}, nil
	case "smma":
		return []strategy.Strategy{trend.NewSmmaStrategy()}, nil
	case "trima":
		return []strategy.Strategy{trend.NewTrimaStrategy()}, nil
	case "triple_moving_average_crossover", "triple_ma_crossover":
		return []strategy.Strategy{trend.NewTripleMovingAverageCrossoverStrategy()}, nil
	case "trix":
		return []strategy.Strategy{trend.NewTrixStrategy()}, nil
	case "tsi":
		return []strategy.Strategy{trend.NewTsiStrategy()}, nil
	case "vwma":
		return []strategy.Strategy{trend.NewVwmaStrategy()}, nil
	case "weighted_close":
		return []strategy.Strategy{trend.NewWeightedCloseStrategy()}, nil

	// Momentum strategies
	case "awesome_oscillator":
		return []strategy.Strategy{momentum.NewAwesomeOscillatorStrategy()}, nil
	case "coppock_curve":
		return []strategy.Strategy{momentum.NewCoppockCurveStrategy()}, nil
	case "elder_ray":
		return []strategy.Strategy{momentum.NewElderRayStrategy()}, nil
	case "ichimoku_cloud":
		return []strategy.Strategy{momentum.NewIchimokuCloudStrategy()}, nil
	case "rsi":
		return []strategy.Strategy{momentum.NewRsiStrategy()}, nil
	case "stochastic_oscillator":
		return []strategy.Strategy{momentum.NewStochasticOscillatorStrategy()}, nil
	case "stochastic_rsi":
		return []strategy.Strategy{momentum.NewStochasticRsiStrategy()}, nil
	case "triple_rsi":
		return []strategy.Strategy{momentum.NewTripleRsiStrategy()}, nil
	case "williams_r":
		return []strategy.Strategy{momentum.NewWilliamsRStrategy()}, nil

	// Volatility strategies
	case "bollinger_bands":
		return []strategy.Strategy{volatility.NewBollingerBandsStrategy()}, nil
	case "donchian_channel_breakout":
		return []strategy.Strategy{volatility.NewDonchianChannelBreakoutStrategy()}, nil
	case "keltner_channel":
		return []strategy.Strategy{volatility.NewKeltnerChannelStrategy()}, nil
	case "super_trend":
		return []strategy.Strategy{volatility.NewSuperTrendStrategy()}, nil

	// Volume strategies
	case "chaikin_money_flow":
		return []strategy.Strategy{volume.NewChaikinMoneyFlowStrategy()}, nil
	case "ease_of_movement":
		return []strategy.Strategy{volume.NewEaseOfMovementStrategy()}, nil
	case "force_index":
		return []strategy.Strategy{volume.NewForceIndexStrategy()}, nil
	case "money_flow_index":
		return []strategy.Strategy{volume.NewMoneyFlowIndexStrategy()}, nil
	case "negative_volume_index":
		return []strategy.Strategy{volume.NewNegativeVolumeIndexStrategy()}, nil
	case "obv":
		return []strategy.Strategy{volume.NewObvStrategy()}, nil
	case "percent_b_and_mfi", "percent_band_and_mfi", "percent_b_mfi":
		return []strategy.Strategy{volume.NewPercentBandMFIStrategy()}, nil
	case "weighted_average_price":
		return []strategy.Strategy{volume.NewWeightedAveragePriceStrategy()}, nil

	// Compound strategies
	case "macd_rsi":
		return []strategy.Strategy{compound.NewMacdRsiStrategy()}, nil

	default:
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
}

func main() {
	var repositoryName string
	var repositoryConfig string
	var reportName string
	var reportConfig string
	var strategyNames string
	var workers int
	var lastDays int
	var addSplits bool
	var addAnds bool

	stdErr := log.New(os.Stderr, "", 0)
	stdErr.Println("Indicator Backtest")
	stdErr.Println("Copyright (c) 2021-2026 The Indicator Authors.")
	stdErr.Println("The source code is provided under GNU AGPLv3 License.")
	stdErr.Println("https://github.com/cinar/indicator")
	stdErr.Println()
	stdErr.Println("DISCLAIMER: For educational and research purposes only. Not investment or financial advice.")
	stdErr.Println("Backtested performance is hypothetical and does not guarantee future results.")
	stdErr.Println()

	flag.StringVar(&repositoryName, "repository-name", "filesystem", "repository name")
	flag.StringVar(&repositoryConfig, "repository-config", "", "repository config")
	flag.StringVar(&reportName, "report-name", "html", "report name")
	flag.StringVar(&reportConfig, "report-config", ".", "report type")
	flag.StringVar(&strategyNames, "strategies", "", "comma-separated list of strategies to backtest (e.g. macd, rsi, buy_and_hold)")
	flag.StringVar(&strategyNames, "strategy", "", "comma-separated list of strategies to backtest (alias for -strategies)")
	flag.IntVar(&workers, "workers", backtest.DefaultBacktestWorkers, "number of concurrent workers")
	flag.IntVar(&lastDays, "last", backtest.DefaultLastDays, "number of days to do backtest")
	flag.BoolVar(&addSplits, "splits", false, "add the split strategies")
	flag.BoolVar(&addAnds, "ands", false, "add the and strategies")
	flag.Parse()

	logger := slog.Default()

	if strings.TrimSpace(strategyNames) == "" {
		logger.Error("No strategies specified. Please specify strategies using the -strategies flag (e.g. -strategies macd,rsi or -strategies all).")
		os.Exit(1)
	}

	selectedStrategies, err := parseStrategies(strategyNames)
	if err != nil {
		logger.Error("Unable to parse strategies.", "error", err)
		os.Exit(1)
	}

	source, err := asset.NewRepository(repositoryName, repositoryConfig)
	if err != nil {
		logger.Error("Unable to initialize source.", "error", err)
		os.Exit(1)
	}

	report, err := backtest.NewReport(reportName, reportConfig)
	if err != nil {
		logger.Error("Unable to initialize report.", "error", err)
		os.Exit(1)
	}

	backtester := backtest.NewBacktest(source, report)
	backtester.Workers = workers
	backtester.LastDays = lastDays
	backtester.Logger = logger
	backtester.Names = append(backtester.Names, flag.Args()...)
	backtester.Strategies = append(backtester.Strategies, selectedStrategies...)

	if addSplits {
		backtester.Strategies = append(backtester.Strategies, strategy.AllSplitStrategies(backtester.Strategies)...)
	}

	if addAnds {
		backtester.Strategies = append(backtester.Strategies, strategy.AllAndStrategies(backtester.Strategies)...)
	}

	err = backtester.Run()
	if err != nil {
		logger.Error("Unable to run backtest.", "error", err)
		os.Exit(1)
	}
}
