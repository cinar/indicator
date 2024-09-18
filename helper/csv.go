// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"encoding/csv"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

const (
	// CsvHeaderTag represents the parameter name for the column header.
	CsvHeaderTag = "header"

	// CsvFormatTag represents the parameter name for the column format.
	CsvFormatTag = "format"

	// DefaultDateTimeFormat denotes the default format of a date and time column.
	DefaultDateTimeFormat = "2006-01-02 15:04:05"
)

// csvColumn represents the mapping between the CSV column and
// the corresponding struct field.
type csvColumn struct {
	Header      string
	ColumnIndex int
	FieldIndex  int
	Format      string
}

// Csv represents the configuration for CSV reader and writer.
type Csv[T any] struct {
	// hasHeader indicates whether the CSV contains a header row.
	hasHeader bool

	// columns are the mappings between the CSV columns and
	// the corresponding struct fields.
	columns []csvColumn

	// Logger is the slog logger instance.
	Logger *slog.Logger
}

// NewCsv function initializes a new CSV instance. The parameter
// hasHeader indicates whether the CSV contains a header row.
func NewCsv[T any](hasHeader bool) (*Csv[T], error) {
	c := &Csv[T]{
		hasHeader: hasHeader,
		Logger:    slog.Default(),
	}

	// Row type must be a pointer to struct.
	structType := reflect.TypeOf((*T)(nil)).Elem()
	if structType.Kind() != reflect.Struct {
		return nil, errors.New("type not a struct")
	}

	// Create a mapping linking CSV columns to corresponding struct fields.
	c.columns = make([]csvColumn, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		header, ok := field.Tag.Lookup(CsvHeaderTag)
		if !ok {
			header = field.Name
		}

		format, ok := field.Tag.Lookup(CsvFormatTag)
		if !ok {
			format = DefaultDateTimeFormat
		}

		c.columns[i] = csvColumn{
			Header:      header,
			ColumnIndex: i,
			FieldIndex:  i,
			Format:      format,
		}
	}

	return c, nil
}

// ReadFromReader parses the CSV data from the provided reader,
// maps the data to corresponding struct fields, and delivers
// the resulting it through the channel.
func (c *Csv[T]) ReadFromReader(reader io.Reader) <-chan *T {
	rows := make(chan *T)

	go func() {
		defer close(rows)

		csvReader := csv.NewReader(reader)

		// If CSV has headers, align column indices to match the
		// order of column headers.
		if c.hasHeader {
			err := c.updateColumnIndexes(csvReader)
			if err != nil {
				c.Logger.Error("Unable to update the column indexes.", "error", err)
				return
			}
		}

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				c.Logger.Error("Unable to read row.", "error", err)
				break
			}

			row := new(T)
			rowValue := reflect.ValueOf(row).Elem()

			for _, column := range c.columns {
				if column.ColumnIndex == -1 {
					continue
				}

				err := setReflectValue(rowValue.Field(column.FieldIndex),
					record[column.ColumnIndex], column.Format)
				if err != nil {
					c.Logger.Error("Unable to set value.", "error", err)
					return
				}
			}

			rows <- row
		}
	}()

	return rows
}

// ReadFromFile parses the CSV data from the provided file name,
// maps the data to corresponding struct fields, and delivers
// the resulting rows through the channel.
func (c *Csv[T]) ReadFromFile(fileName string) (<-chan *T, error) {
	file, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	rows := Waitable(wg, c.ReadFromReader(file))

	go func() {
		wg.Wait()
		err := file.Close()
		if err != nil {
			c.Logger.Error("Unable to close file.", "error", err)
		}
	}()

	return rows, nil
}

// AppendToFile appends the provided rows of data to the end of the specified file, creating
// the file if it doesn't exist.  In append mode, the function assumes that the existing
// file's column order matches the field order of the given row struct to ensure consistent
// data structure.
func (c *Csv[T]) AppendToFile(fileName string, rows <-chan *T) error {
	file, err := os.OpenFile(filepath.Clean(fileName), os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	err = c.writeToWriter(file, false, rows)
	if err != nil {
		return err
	}

	return file.Close()
}

// WriteToFile creates a new file with the given name and writes the provided rows
// of data to it, overwriting any existing content.
func (c *Csv[T]) WriteToFile(fileName string, rows <-chan *T) error {
	file, err := os.OpenFile(filepath.Clean(fileName), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	err = c.writeToWriter(file, true, rows)
	if err != nil {
		return err
	}

	return file.Close()
}

// updateColumnIndexes aligns column indices to match the order of column headers.
func (c *Csv[T]) updateColumnIndexes(csvReader *csv.Reader) error {
	headers, err := csvReader.Read()
	if err != nil {
		return err
	}

	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	for i := range c.columns {
		index, ok := headerMap[c.columns[i].Header]
		if !ok {
			index = -1
		}

		c.columns[i].ColumnIndex = index
	}

	return nil
}

// writeToWriter writes the provided rows of data to the specified writer, with the option
// to include or exclude headers for flexibility in data presentation.
func (c *Csv[T]) writeToWriter(writer io.Writer, writeHeader bool, rows <-chan *T) error {
	csvWriter := csv.NewWriter(writer)

	if writeHeader {
		err := c.writeHeaderToCsvWriter(csvWriter)
		if err != nil {
			return err
		}
	}

	record := make([]string, len(c.columns))

	for row := range rows {
		rowValue := reflect.ValueOf(row).Elem()

		for i, column := range c.columns {
			stringValue, err := getReflectValue(rowValue.Field(column.FieldIndex), column.Format)
			if err != nil {
				return err
			}

			record[i] = stringValue
		}

		err := csvWriter.Write(record)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return csvWriter.Error()
}

// writeHeaderToCsvWriter writes the column headers for the CSV data to the specified CSV writer.
func (c *Csv[T]) writeHeaderToCsvWriter(csvWriter *csv.Writer) error {
	header := make([]string, len(c.columns))

	for i, column := range c.columns {
		header[i] = column.Header
	}

	return csvWriter.Write(header)
}

// ReadFromCsvFile creates a CSV instance, parses CSV data from the provided filename,
// maps the data to corresponding struct fields, and delivers it through the channel.
func ReadFromCsvFile[T any](fileName string, hasHeader bool) (<-chan *T, error) {
	csv, err := NewCsv[T](hasHeader)
	if err != nil {
		return nil, err
	}

	return csv.ReadFromFile(fileName)
}

// AppendOrWriteToCsvFile writes the provided rows of data to the specified file, appending to
// the existing file if it exists or creating a new one if it doesn't. In append mode, the
// function assumes that the existing file's column order matches the field order of the
// given row struct to ensure consistent data structure.
func AppendOrWriteToCsvFile[T any](fileName string, hasHeader bool, rows <-chan *T) error {
	csv, err := NewCsv[T](hasHeader)
	if err != nil {
		return err
	}

	stat, err := os.Stat(filepath.Clean(fileName))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else if stat.Size() > 0 {
		return csv.AppendToFile(fileName, rows)
	}

	return csv.WriteToFile(fileName, rows)
}
