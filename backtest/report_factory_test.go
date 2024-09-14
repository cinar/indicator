// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest_test

import (
	"testing"

	"github.com/cinar/indicator/v2/backtest"
)

func TestNewReportUnknown(t *testing.T) {
	report, err := backtest.NewReport("unknown", "")
	if err == nil {
		t.Fatalf("unknown report: %T", report)
	}
}

func TestRegisterReportBuilder(t *testing.T) {
	builderName := "testbuilder"

	report, err := backtest.NewReport(builderName, "")
	if err == nil {
		t.Fatalf("testbuilder is: %T", report)
	}

	backtest.RegisterReportBuilder(builderName, func(config string) (backtest.Report, error) {
		return backtest.NewHTMLReport(config), nil
	})

	report, err = backtest.NewReport(builderName, "")
	if err != nil {
		t.Fatalf("testbuilder is not found: %v", err)
	}

	_, ok := report.(*backtest.HTMLReport)
	if !ok {
		t.Fatalf("testbuilder is: %T", report)
	}
}

func TestNewReportMemory(t *testing.T) {
	report, err := backtest.NewReport(backtest.HTMLReportBuilderName, "")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := report.(*backtest.HTMLReport)
	if !ok {
		t.Fatalf("report not correct type: %T", report)
	}
}
