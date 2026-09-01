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
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/backtest"
	"github.com/cinar/indicator/v2/examples/compound"
	"github.com/cinar/indicator/v2/examples/momentum"
	"github.com/cinar/indicator/v2/examples/trend"
	"github.com/cinar/indicator/v2/examples/volatility"
	"github.com/cinar/indicator/v2/examples/volume"
	"github.com/cinar/indicator/v2/strategy"
)

var parensRegex = regexp.MustCompile(`\(([^)]+)\)`)

// camelToSnake converts CamelCase string to snake_case.
func camelToSnake(s string) string {
	var res strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(s[i-1])
				if unicode.IsLower(prev) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) {
					res.WriteRune('_')
				}
			}
			res.WriteRune(unicode.ToLower(r))
		} else {
			res.WriteRune(r)
		}
	}
	return res.String()
}

// normalize cleans and standardizes a name into snake_case format.
func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "-", "_")

	var res strings.Builder
	prevUnderscore := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			res.WriteRune(r)
			prevUnderscore = false
		} else if !prevUnderscore {
			res.WriteRune('_')
			prevUnderscore = true
		}
	}
	return strings.Trim(res.String(), "_")
}

// strategyKeys extracts all candidate lookup keys for a given strategy.
func strategyKeys(s strategy.Strategy) []string {
	keys := make(map[string]bool)

	// 1. From struct type name
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	typeName := t.Name()
	snakeType := camelToSnake(typeName)
	trimmedType := strings.TrimSuffix(snakeType, "_strategy")

	keys[snakeType] = true
	keys[trimmedType] = true

	// Custom aliases for common abbreviations
	if trimmedType == "triple_moving_average_crossover" {
		keys["triple_ma_crossover"] = true
	}
	if trimmedType == "percent_band_mfi" {
		keys["percent_b_and_mfi"] = true
		keys["percent_band_and_mfi"] = true
		keys["percent_b_mfi"] = true
	}

	// 2. From strategy Name()
	displayName := s.Name()
	normDisplay := normalize(displayName)
	keys[normDisplay] = true
	keys[strings.TrimSuffix(normDisplay, "_strategy")] = true

	// Extract acronyms or text inside parentheses, e.g. "Balance of Power (BoP) Strategy" -> "bop"
	matches := parensRegex.FindAllStringSubmatch(displayName, -1)
	for _, match := range matches {
		if len(match) > 1 {
			normParen := normalize(match[1])
			if normParen != "" && !unicode.IsDigit(rune(normParen[0])) {
				keys[normParen] = true
			}
		}
	}

	var result []string
	for k := range keys {
		if k != "" {
			result = append(result, k)
		}
	}
	return result
}

type strategyCatalog struct {
	categories map[string][]strategy.Strategy
	strategies map[string][]strategy.Strategy
}

func newStrategyCatalog() *strategyCatalog {
	cat := &strategyCatalog{
		categories: map[string][]strategy.Strategy{
			"compound":   compound.AllStrategies(),
			"momentum":   momentum.AllStrategies(),
			"strategy":   strategy.AllStrategies(),
			"trend":      trend.AllStrategies(),
			"volatility": volatility.AllStrategies(),
			"volume":     volume.AllStrategies(),
		},
		strategies: make(map[string][]strategy.Strategy),
	}

	for _, strats := range cat.categories {
		for _, s := range strats {
			keys := strategyKeys(s)
			for _, k := range keys {
				cat.strategies[k] = append(cat.strategies[k], s)
			}
		}
	}

	return cat
}

func (c *strategyCatalog) Resolve(token string) ([]strategy.Strategy, error) {
	norm := normalize(token)
	norm = strings.TrimSuffix(norm, "_strategy")

	if norm == "all" {
		var all []strategy.Strategy
		categoryOrder := []string{"strategy", "compound", "momentum", "trend", "volatility", "volume"}
		for _, catName := range categoryOrder {
			all = append(all, c.categories[catName]...)
		}
		return all, nil
	}

	if cat, ok := c.categories[norm]; ok {
		return cat, nil
	}
	if norm == "base" {
		return c.categories["strategy"], nil
	}

	if strats, ok := c.strategies[norm]; ok && len(strats) > 0 {
		return strats, nil
	}

	return nil, fmt.Errorf("unknown strategy: %s", token)
}

func parseStrategies(input string) ([]strategy.Strategy, error) {
	tokens := strings.Split(input, ",")
	catalog := newStrategyCatalog()

	var result []strategy.Strategy
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		strats, err := catalog.Resolve(token)
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
