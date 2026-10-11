// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"errors"
	"fmt"
)

// annotationReportColumn is the annotation report column struct.
type annotationReportColumn struct {
	ReportColumn
	values <-chan string
}

// NewAnnotationReportColumn returns a new instance of an annotation column for a report.
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

// Value returns the next data value for the report column. It returns
// an error when the backing channel is exhausted before the report's
// time axis is, which would otherwise silently produce a "null" value.
func (c *annotationReportColumn) Value() (string, error) {
	value, ok := <-c.values

	if !ok {
		return "", errors.New("annotation report column: no more data available")
	}

	if value != "" {
		return fmt.Sprintf("%q", value), nil
	}

	return "null", nil
}
