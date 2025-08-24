// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

func TestCsv(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	reader := strings.NewReader("Date,Asset,Open,Close\n\"2023-11-26 00:00:00\",\"SP500\",10.2,30.4\n")

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	row := <-csv.ReadFromReader(reader)

	if row.Close != 30.4 {
		t.Fatalf("actual %v expected 30.4", row.Close)
	}
}

func TestCsvNoHeader(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	reader := strings.NewReader("10.2,30.4\n")

	csv, err := helper.NewCsv[Row](helper.WithoutCsvHeader[Row]())
	if err != nil {
		t.Fatal(err)
	}

	row := <-csv.ReadFromReader(reader)

	if row.Close != 10.2 {
		t.Fatalf("actual %v expected 10.2", row.Close)
	}

	fmt.Printf("%v", row)
}

func TestCsvInvalidColumns(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	reader := strings.NewReader("1,2\n1\n")

	csv, err := helper.NewCsv[Row](helper.WithoutCsvHeader[Row]())
	if err != nil {
		t.Fatal(err)
	}

	_, ok := <-csv.ReadFromReader(reader)
	if !ok {
		t.Fatalf("actual closed expected 1,2")
	}

	row, ok := <-csv.ReadFromReader(reader)
	if ok {
		t.Fatalf("actual %v expected closed", row)
	}
}

func TestCsvMissingHeader(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	reader := strings.NewReader("")

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	row, ok := <-csv.ReadFromReader(reader)
	if ok {
		t.Fatalf("actual %v expected closed", row)
	}
}

func TestCsvInvalidField(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	reader := strings.NewReader("\"ABCD\",\"EFGH\"\n")

	csv, err := helper.NewCsv[Row](helper.WithoutCsvHeader[Row]())
	if err != nil {
		t.Fatal(err)
	}

	row, ok := <-csv.ReadFromReader(reader)
	if ok {
		t.Fatalf("actual %v expected closed", row)
	}
}

func TestCsvNoStruct(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	_, err := helper.NewCsv[*Row]()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCsvReadFromFile(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := csv.ReadFromFile("testdata/with_header.csv")
	if err != nil {
		t.Fatal(err)
	}

	row := <-rows

	if row.Close != 30.4 {
		t.Fatalf("actual %v expected 30.4", row.Close)
	}
}

func TestCsvReadFromMissingFile(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	_, err = csv.ReadFromFile("testdata/missing.csv")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestReadFromCsvFile(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	rows, err := helper.ReadFromCsvFile[Row]("testdata/with_header.csv")
	if err != nil {
		t.Fatal(err)
	}

	row := <-rows

	if row.Close != 30.4 {
		t.Fatalf("actual %v expected 30.4", row.Close)
	}
}

func TestReadFromCsvFileNoStruct(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	_, err := helper.ReadFromCsvFile[*Row]("testdata/with_header.csv")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCsvWriteToFile(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	input := []*Row{
		{Close: 10, High: 20},
		{Close: 30, High: 40},
	}

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	fileName := "test_csv_write_to_file.csv"
	defer helper.Remove(t, fileName)

	err = csv.WriteToFile(fileName, helper.SliceToChan(input))
	if err != nil {
		t.Fatal(err)
	}

	csv, err = helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := csv.ReadFromFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.SliceToChan(input)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCsvAppendToFile(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	input := []*Row{
		{Close: 10, High: 20},
		{Close: 30, High: 40},
	}

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	fileName := "test_csv_append_to_file.csv"
	defer helper.Remove(t, fileName)

	err = csv.WriteToFile(fileName, helper.SliceToChan(input[:1]))
	if err != nil {
		t.Fatal(err)
	}

	csv, err = helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	err = csv.AppendToFile(fileName, helper.SliceToChan(input[1:]))
	if err != nil {
		t.Fatal(err)
	}

	csv, err = helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := csv.ReadFromFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.SliceToChan(input)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCsvWriteToInvalidFile(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	rows := helper.SliceToChan([]*Row{
		{Close: 10, High: 20},
		{Close: 30, High: 40},
	})

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	err = csv.WriteToFile("testdata/invalid/data.csv", rows)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCsvWriteToFileInvalidField(t *testing.T) {
	type Invalid interface{}

	type Row struct {
		Close float64
		High  Invalid
	}

	rows := helper.SliceToChan([]*Row{
		{Close: 10, High: nil},
		{Close: 30, High: nil},
	})

	csv, err := helper.NewCsv[Row]()
	if err != nil {
		t.Fatal(err)
	}

	fileName := "test_csv_write_to_file_invalid_field.csv"
	defer helper.Remove(t, fileName)

	err = csv.WriteToFile(fileName, rows)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAppendOrWriteToCsvFile(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	input := []*Row{
		{Close: 10, High: 20},
		{Close: 30, High: 40},
	}

	fileName := "test_append_or_write_to_csv_file.csv"
	defer helper.Remove(t, fileName)

	err := helper.AppendOrWriteToCsvFile(fileName, helper.SliceToChan(input[:1]))
	if err != nil {
		t.Fatal(err)
	}

	err = helper.AppendOrWriteToCsvFile(fileName, helper.SliceToChan(input[1:]))
	if err != nil {
		t.Fatal(err)
	}

	actual, err := helper.ReadFromCsvFile[Row](fileName)
	if err != nil {
		t.Fatal(err)
	}

	expected := helper.SliceToChan(input)

	err = helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAppendOrWriteToCsvFileNoStruct(t *testing.T) {
	type Row struct {
		Close float64
		High  float64
	}

	input := helper.SliceToChan([]**Row{})

	err := helper.AppendOrWriteToCsvFile[*Row]("test_append_or_write_to_csv_file_no_struct.csv", input)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCsvWithDefaultDateFormat(t *testing.T) {
	type Row struct {
		Date  time.Time
		Close float64
		High  float64
	}

	reader := strings.NewReader("Date,Asset,Open,Close\n\"2023-11-26 01:02:03\",\"SP500\",10.2,30.4\n")

	csv, err := helper.NewCsv[Row](
		helper.WithCsvDefaultDateTimeFormat[Row]("2006-01-02 15:04:05"),
	)
	if err != nil {
		t.Fatal(err)
	}

	row := <-csv.ReadFromReader(reader)

	expected := time.Date(2023, 11, 26, 1, 2, 3, 0, time.UTC)

	if !row.Date.Equal(expected) {
		t.Fatalf("actual %v expected %v", row.Date, expected)
	}
}
