// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

func TestReportWriteToFile(t *testing.T) {
	type Row struct {
		Date       time.Time `format:"2006-01-02"`
		High       float64
		Low        float64
		Close      float64
		Annotation string
	}

	input, err := helper.ReadFromCsvFile[Row]("testdata/report.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 5)
	dates := helper.Map(inputs[0], func(row *Row) time.Time { return row.Date })
	highs := helper.Map(inputs[1], func(row *Row) float64 { return row.High })
	lows := helper.Map(inputs[2], func(row *Row) float64 { return row.Low })
	closes := helper.Map(inputs[3], func(row *Row) float64 { return row.Close })
	annotations := helper.Map(inputs[4], func(row *Row) string { return row.Annotation })

	report := helper.NewReport("Test Report", dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("High", highs))
	report.AddColumn(helper.NewNumericReportColumn("Low", lows))
	report.AddColumn(helper.NewNumericReportColumn("Close", closes), 0, 1)
	report.AddColumn(helper.NewAnnotationReportColumn(annotations), 0, 1)

	fileName := "report.html"
	defer helper.Remove(t, fileName)

	err = report.WriteToFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReportWriteToFileFailed(t *testing.T) {
	type Row struct {
		Date       time.Time `format:"2006-01-02"`
		High       float64
		Low        float64
		Close      float64
		Annotation string
	}

	input, err := helper.ReadFromCsvFile[Row]("testdata/report.csv")
	if err != nil {
		t.Fatal(err)
	}

	inputs := helper.Duplicate(input, 2)
	dates := helper.Map(inputs[0], func(row *Row) time.Time { return row.Date })
	closes := helper.Map(inputs[1], func(row *Row) float64 { return row.Close })

	report := helper.NewReport("Test Report", dates)
	report.AddColumn(helper.NewNumericReportColumn("Close", closes))

	err = report.WriteToFile("")
	if err == nil {
		t.Fatal("expected error")
	}
}
