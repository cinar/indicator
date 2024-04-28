// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "fmt"

// numericReportColumn is the number report column struct.
type numericReportColumn[T Number] struct {
	ReportColumn
	name   string
	values <-chan T
}

// NewNumericReportColumn returns a new instance of a numeric data column for a report.
func NewNumericReportColumn[T Number](name string, values <-chan T) ReportColumn {
	return &numericReportColumn[T]{
		name:   name,
		values: values,
	}
}

// Name returns the name of the report column.
func (c *numericReportColumn[T]) Name() string {
	return c.name
}

// Type returns number as the data type.
func (*numericReportColumn[T]) Type() string {
	return "number"
}

// Role returns the role of the report column.
func (*numericReportColumn[T]) Role() string {
	return "data"
}

// Value returns the next data value for the report column.
func (c *numericReportColumn[T]) Value() string {
	return fmt.Sprintf("%v", <-c.values)
}
