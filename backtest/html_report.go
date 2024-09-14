// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest

import (
	// Go embed report template.
	_ "embed"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"text/template"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

const (
	// DefaultWriteStrategyReports is the default state of writing individual strategy reports.
	DefaultWriteStrategyReports = true
)

//go:embed "html_report.tmpl"
var htmlReportTmpl string

//go:embed "html_asset_report.tmpl"
var htmlAssetReportTmpl string

// HTMLReport is the backtest HTML report.
type HTMLReport struct {
	Report

	// outputDir is the output directory for the generated reports.
	outputDir string

	// assetResults is the mapping from the asset name to strategy results.
	assetResults map[string][]*htmlReportResult

	// bestResults is the best results for each asset.
	bestResults []*htmlReportResult

	// WriteStrategyReports indicates whether the individual strategy reports should be generated.
	WriteStrategyReports bool

	// DateFormat is the date format that is used in the reports.
	DateFormat string
}

// htmlReportResult encapsulates the outcome of running a strategy.
type htmlReportResult struct {
	// AssetName is the name of the asset.
	AssetName string

	// StrategyName is the name of the strategy.
	StrategyName string

	// Action is the last recommended action by the strategy.
	Action strategy.Action

	// Since indicates how long the current action recommendation has been in effect.
	Since int

	// Outcome is the effectiveness of applying the recommended actions.
	Outcome float64

	// Transactions is the number of transactions made by the strategy.
	Transactions int
}

// NewHTMLReport initializes a new HTML report instance.
func NewHTMLReport(outputDir string) *HTMLReport {
	return &HTMLReport{
		outputDir:            outputDir,
		assetResults:         make(map[string][]*htmlReportResult),
		WriteStrategyReports: DefaultWriteStrategyReports,
		DateFormat:           helper.DefaultReportDateFormat,
	}
}

// Begin is called when the backtest starts.
func (h *HTMLReport) Begin(assetNames []string, _ []strategy.Strategy) error {
	// Make sure that output directory exists.
	err := os.MkdirAll(h.outputDir, 0o700)
	if err != nil {
		return fmt.Errorf("unable to make the output directory: %w", err)
	}

	h.bestResults = make([]*htmlReportResult, 0, len(assetNames))

	return nil
}

// AssetBegin is called when backtesting for the given asset begins.
func (h *HTMLReport) AssetBegin(name string, strategies []strategy.Strategy) error {
	_, ok := h.assetResults[name]
	if ok {
		return fmt.Errorf("asset has already begun: %s", name)
	}

	h.assetResults[name] = make([]*htmlReportResult, 0, len(strategies))

	return nil
}

// Write writes the given strategy actions and outomes to the report.
func (h *HTMLReport) Write(assetName string, currentStrategy strategy.Strategy, snapshots <-chan *asset.Snapshot, actions <-chan strategy.Action, outcomes <-chan float64) error {
	actionsSplice := helper.Duplicate(actions, 3)

	actions = helper.Last(actionsSplice[0], 1)
	sinces := helper.Last(helper.Since[strategy.Action, int](actionsSplice[1]), 1)
	outcomes = helper.Last(outcomes, 1)
	transactions := helper.Last(strategy.CountTransactions(actionsSplice[2]), 1)

	// Generate inidividual strategy report.
	if h.WriteStrategyReports {
		report := currentStrategy.Report(snapshots)
		report.DateFormat = h.DateFormat

		reportFile := h.strategyReportFileName(assetName, currentStrategy.Name())

		err := report.WriteToFile(path.Join(h.outputDir, reportFile))
		if err != nil {
			return fmt.Errorf("unable to write report for %s (%v)", assetName, err)
		}
	} else {
		go helper.Drain(snapshots)
	}

	// Get asset strategy results.
	results, ok := h.assetResults[assetName]
	if !ok {
		return fmt.Errorf("asset has not begun: %s", assetName)
	}

	// Append current strategy result for the asset.
	h.assetResults[assetName] = append(results, &htmlReportResult{
		AssetName:    assetName,
		StrategyName: currentStrategy.Name(),
		Action:       <-actions,
		Since:        <-sinces,
		Outcome:      <-outcomes * 100,
		Transactions: <-transactions,
	})

	return nil
}

// AssetEnd is called when backtesting for the given asset ends.
func (h *HTMLReport) AssetEnd(name string) error {
	results, ok := h.assetResults[name]
	if !ok {
		return fmt.Errorf("asset has not begun: %s", name)
	}

	delete(h.assetResults, name)

	// Sort the backtest results by the outcomes.
	slices.SortFunc(results, func(a, b *htmlReportResult) int {
		return int(b.Outcome - a.Outcome)
	})

	bestResult := results[0]

	// Report the best result for the current asset.
	log.Printf("Best outcome for %s is %.2f%% with %s.", name, bestResult.Outcome, bestResult.StrategyName)
	h.bestResults = append(h.bestResults, bestResult)

	// Write the asset report.
	err := h.writeAssetReport(name, results)
	if err != nil {
		return fmt.Errorf("unable to write report for %s: %w", name, err)
	}

	return nil
}

// End is called when the backtest ends.
func (h *HTMLReport) End() error {
	// Sort the best results by the outcomes.
	slices.SortFunc(h.bestResults, func(a, b *htmlReportResult) int {
		return int(b.Outcome - a.Outcome)
	})

	return h.writeReport()
}

// strategyReportFileName defines the HTML report file name for the given asset and strategy.
func (*HTMLReport) strategyReportFileName(assetName, strategyName string) string {
	return fmt.Sprintf("%s - %s.html", assetName, strategyName)
}

// writeAssetReport generates a detailed report for the asset, summarizing the backtest results.
func (h *HTMLReport) writeAssetReport(name string, results []*htmlReportResult) error {
	type Model struct {
		AssetName   string
		Results     []*htmlReportResult
		GeneratedOn string
	}

	model := Model{
		AssetName:   name,
		Results:     results,
		GeneratedOn: time.Now().String(),
	}

	file, err := os.Create(filepath.Join(h.outputDir, fmt.Sprintf("%s.html", name)))
	if err != nil {
		return fmt.Errorf("unable to open asset report file for %s: %w", name, err)
	}

	defer helper.CloseAndLogError(file, "unable to close asset report file")

	tmpl := template.Must(template.New("report").Parse(htmlAssetReportTmpl))

	err = tmpl.Execute(file, model)
	if err != nil {
		return fmt.Errorf("unable to execute asset report template: %w", err)
	}

	return nil
}

// writeReport generates a detailed report for the best results for all the assets.
func (h *HTMLReport) writeReport() error {
	type Model struct {
		Results     []*htmlReportResult
		GeneratedOn string
	}

	model := Model{
		Results:     h.bestResults,
		GeneratedOn: time.Now().String(),
	}

	file, err := os.Create(filepath.Join(h.outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("unable to open main report file: %w", err)
	}

	defer helper.CloseAndLogError(file, "unable to close main report file")

	tmpl := template.Must(template.New("report").Parse(htmlReportTmpl))

	err = tmpl.Execute(file, model)
	if err != nil {
		return fmt.Errorf("unable to execute main report template: %w", err)
	}

	return nil
}
