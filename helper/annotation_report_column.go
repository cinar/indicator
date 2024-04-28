// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "fmt"

// annotationReportColumn is the annotation report column struct.
type annotationReportColumn struct {
	ReportColumn
	values <-chan string
}

// NewAnnotationReportColumn returns a new instance of a annotation column for a report.
func NewAnnotationReportColumn(values <-chan string) ReportColumn {
	return &annotationReportColumn{
		values: values,
	}
}

// Name returns the name of the report column.
func (*annotationReportColumn) Name() string {
	return ""
}

// Type returns number as the data type.
func (*annotationReportColumn) Type() string {
	return "string"
}

// Role returns the role of the report column.
func (*annotationReportColumn) Role() string {
	return "annotation"
}

// Value returns the next data value for the report column.
func (c *annotationReportColumn) Value() string {
	value := <-c.values

	if value != "" {
		return fmt.Sprintf("\"%s\"", value)
	}

	return "null"
}
