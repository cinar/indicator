// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package backtest

import (
	"fmt"
)

const (
	// HTMLReportBuilderName is the name for the HTML report builder.
	HTMLReportBuilderName = "html"
)

// ReportBuilderFunc defines a function to build a new report using the given configuration parameter.
type ReportBuilderFunc func(config string) (Report, error)

// reportBuilders provides mapping for the report builders.
var reportBuilders = map[string]ReportBuilderFunc{
	HTMLReportBuilderName: htmlReportBuilder,
}

// RegisterReportBuilder registers the given builder.
func RegisterReportBuilder(name string, builder ReportBuilderFunc) {
	reportBuilders[name] = builder
}

// NewReport builds a new report by the given name type and the configuration.
func NewReport(name, config string) (Report, error) {
	builder, ok := reportBuilders[name]
	if !ok {
		return nil, fmt.Errorf("unknown report: %s", name)
	}

	return builder(config)
}

// htmlReportBuilder builds a new HTML report instance.
func htmlReportBuilder(config string) (Report, error) {
	return NewHTMLReport(config), nil
}
