// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// main is the indicator command line program. It computes a single named
// indicator over an asset's OHLCV data and writes the result as CSV.
package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/registry"
)

// paramFlags collects repeated "-param Name=Value" flags into a map.
type paramFlags map[string]string

// String returns the flag's default text representation.
func (p paramFlags) String() string {
	return ""
}

// Set parses a single "Name=Value" flag occurrence.
func (p paramFlags) Set(param string) error {
	name, value, ok := strings.Cut(param, "=")
	if !ok {
		return fmt.Errorf("invalid param %q, expected Name=Value", param)
	}

	p[name] = value

	return nil
}

func main() {
	var indicatorName string
	var repositoryName string
	var repositoryConfig string
	var lastDays int

	params := make(paramFlags)

	stdErr := log.New(os.Stderr, "", 0)
	stdErr.Println("Indicator CLI")
	stdErr.Println("Copyright (c) 2021-2026 Onur Cinar.")
	stdErr.Println("The source code is provided under GNU AGPLv3 License.")
	stdErr.Println("https://github.com/cinar/indicator")
	stdErr.Println()

	flag.StringVar(&indicatorName, "indicator", "", "indicator name, see -list")
	flag.StringVar(&repositoryName, "source-name", "filesystem", "source repository name")
	flag.StringVar(&repositoryConfig, "source-config", "", "source repository config")
	flag.IntVar(&lastDays, "last", 0, "number of days to go back, 0 for the full history")
	flag.Var(params, "param", "indicator parameter as Name=Value, repeatable")
	list := flag.Bool("list", false, "list the available indicator names and exit")
	flag.Parse()

	logger := slog.Default()

	if *list {
		for _, name := range registry.Names() {
			fmt.Println(name)
		}

		return
	}

	if indicatorName == "" || flag.NArg() != 1 {
		stdErr.Println("Usage: indicator-cli -indicator <name> [-param Name=Value ...] [-source-name <name> -source-config <config>] <asset>")
		os.Exit(1)
	}

	source, err := asset.NewRepository(repositoryName, repositoryConfig)
	if err != nil {
		logger.Error("Unable to initialize source.", "error", err)
		os.Exit(1)
	}

	assetName := flag.Arg(0)

	var snapshots <-chan *asset.Snapshot
	if lastDays > 0 {
		since := time.Now().AddDate(0, 0, -lastDays)
		snapshots, err = source.GetSince(assetName, since)
	} else {
		snapshots, err = source.Get(assetName)
	}

	if err != nil {
		logger.Error("Unable to get asset snapshots.", "error", err)
		os.Exit(1)
	}

	result, err := registry.Run(context.Background(), indicatorName, params, snapshots)
	if err != nil {
		logger.Error("Unable to compute indicator.", "error", err)
		os.Exit(1)
	}

	if err := writeCsv(os.Stdout, result); err != nil {
		logger.Error("Unable to write result.", "error", err)
		os.Exit(1)
	}
}

// writeCsv writes the indicator result as CSV, one row per date, with a
// "date" column followed by one column per output series.
func writeCsv(w *os.File, result *registry.Result) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := make([]string, 0, len(result.Series)+1)
	header = append(header, "date")
	for _, series := range result.Series {
		header = append(header, series.Name)
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for date := range result.Dates {
		row := make([]string, 0, len(result.Series)+1)
		row = append(row, date.Format("2006-01-02"))

		for _, series := range result.Series {
			value, ok := <-series.Values
			if !ok {
				return fmt.Errorf("indicator produced fewer values than dates for series %q", series.Name)
			}

			row = append(row, strconv.FormatFloat(value, 'f', -1, 64))
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
