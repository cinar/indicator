// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
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
	DefaultDateTimeFormat = "2006-01-02"
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

	// defaultDateTimeFormat is the default format for date and time columns.
	defaultDateTimeFormat string
}

// CsvOption represents a functional option for configuring the CSV instance.
type CsvOption[T any] func(*Csv[T])

// WithoutCsvHeader disables the header row in the CSV.
func WithoutCsvHeader[T any]() CsvOption[T] {
	return func(c *Csv[T]) {
		c.hasHeader = false
	}
}

// WithCsvLogger sets the logger for the CSV instance.
func WithCsvLogger[T any](logger *slog.Logger) CsvOption[T] {
	return func(c *Csv[T]) {
		c.Logger = logger
	}
}

// WithCsvDefaultDateTimeFormat sets the default date and time format for the CSV instance.
func WithCsvDefaultDateTimeFormat[T any](format string) CsvOption[T] {
	return func(c *Csv[T]) {
		c.defaultDateTimeFormat = format
	}
}

// NewCsv creates a new CSV instance with the provided options.
func NewCsv[T any](options ...CsvOption[T]) (*Csv[T], error) {
	c := &Csv[T]{
		hasHeader:             true,
		Logger:                slog.Default(),
		defaultDateTimeFormat: DefaultDateTimeFormat,
	}

	// Apply options to the CSV instance.
	for _, option := range options {
		option(c)
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
			format = c.defaultDateTimeFormat
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
	return c.ReadFromReaderWithContext(context.Background(), reader)
}

// ReadFromReaderWithContext parses the CSV data from the provided reader,
// maps the data to corresponding struct fields, and delivers
// the resulting it through the channel, supporting context cancellation.
func (c *Csv[T]) ReadFromReaderWithContext(ctx context.Context, reader io.Reader) <-chan *T {
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
			select {
			case <-ctx.Done():
				return
			default:
			}

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

			select {
			case <-ctx.Done():
				return
			case rows <- row:
			}
		}
	}()

	return rows
}

// ReadFromFile parses the CSV data from the provided file name,
// maps the data to corresponding struct fields, and delivers
// the resulting rows through the channel.
func (c *Csv[T]) ReadFromFile(fileName string) (<-chan *T, error) {
	return c.ReadFromFileWithContext(context.Background(), fileName)
}

// ReadFromFileWithContext parses the CSV data from the provided file name,
// maps the data to corresponding struct fields, and delivers
// the resulting rows through the channel, supporting context cancellation.
func (c *Csv[T]) ReadFromFileWithContext(ctx context.Context, fileName string) (<-chan *T, error) {
	file, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	rows := WaitableWithContext(ctx, wg, c.ReadFromReaderWithContext(ctx, file))

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
// the file if it doesn't exist.
func (c *Csv[T]) AppendToFile(fileName string, rows <-chan *T) error {
	return c.AppendToFileWithContext(context.Background(), fileName, rows)
}

// AppendToFileWithContext appends the provided rows of data to the end of the specified file, creating
// the file if it doesn't exist, supporting context cancellation.
func (c *Csv[T]) AppendToFileWithContext(ctx context.Context, fileName string, rows <-chan *T) error {
	file, err := os.OpenFile(filepath.Clean(fileName), os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	err = c.writeToWriterWithContext(ctx, file, false, rows)
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}

// WriteToFile creates a new file with the given name and writes the provided rows
// of data to it, overwriting any existing content.
func (c *Csv[T]) WriteToFile(fileName string, rows <-chan *T) error {
	return c.WriteToFileWithContext(context.Background(), fileName, rows)
}

// WriteToFileWithContext creates a new file with the given name and writes the provided rows
// of data to it, overwriting any existing content, supporting context cancellation.
func (c *Csv[T]) WriteToFileWithContext(ctx context.Context, fileName string, rows <-chan *T) error {
	file, err := os.OpenFile(filepath.Clean(fileName), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	err = c.writeToWriterWithContext(ctx, file, true, rows)
	if err != nil {
		_ = file.Close()
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

// writeToWriterWithContext writes the provided rows of data to the specified writer, with the option
// to include or exclude headers for flexibility in data presentation, supporting context cancellation.
func (c *Csv[T]) writeToWriterWithContext(ctx context.Context, writer io.Writer, writeHeader bool, rows <-chan *T) error {
	csvWriter := csv.NewWriter(writer)

	if writeHeader {
		err := c.writeHeaderToCsvWriter(csvWriter)
		if err != nil {
			return err
		}
	}

	record := make([]string, len(c.columns))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var row *T
		var ok bool
		select {
		case <-ctx.Done():
			return ctx.Err()
		case row, ok = <-rows:
			if !ok {
				goto done
			}
		}

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

done:
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

// ReadFromCsvFile wraps ReadFromCsvFileWithContext for backwards compatibility.
//
// Deprecated: Use ReadFromCsvFileWithContext instead.
func ReadFromCsvFile[T any](fileName string, options ...CsvOption[T]) (<-chan *T, error) {
	return ReadFromCsvFileWithContext[T](context.Background(), fileName, options...)
}

// ReadFromCsvFileWithContext creates a CSV instance, parses CSV data from the provided filename,
// maps the data to corresponding struct fields, and delivers it through the channel, supporting context cancellation.
func ReadFromCsvFileWithContext[T any](ctx context.Context, fileName string, options ...CsvOption[T]) (<-chan *T, error) {
	c, err := NewCsv[T](options...)
	if err != nil {
		return nil, err
	}

	return c.ReadFromFileWithContext(ctx, fileName)
}

// AppendOrWriteToCsvFile wraps AppendOrWriteToCsvFileWithContext for backwards compatibility.
//
// Deprecated: Use AppendOrWriteToCsvFileWithContext instead.
func AppendOrWriteToCsvFile[T any](fileName string, rows <-chan *T, options ...CsvOption[T]) error {
	return AppendOrWriteToCsvFileWithContext[T](context.Background(), fileName, rows, options...)
}

// AppendOrWriteToCsvFileWithContext writes the provided rows of data to the specified file, appending to
// the existing file if it exists or creating a new one if it doesn't, supporting context cancellation.
func AppendOrWriteToCsvFileWithContext[T any](ctx context.Context, fileName string, rows <-chan *T, options ...CsvOption[T]) error {
	c, err := NewCsv[T](options...)
	if err != nil {
		return err
	}

	stat, err := os.Stat(filepath.Clean(fileName))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else if stat.Size() > 0 {
		return c.AppendToFileWithContext(ctx, fileName, rows)
	}

	return c.WriteToFileWithContext(ctx, fileName, rows)
}
