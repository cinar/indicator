// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package helper_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestField(t *testing.T) {
	type Row struct {
		Open  float64
		Close float64
	}

	row := &Row{
		Open:  1,
		Close: 2,
	}

	rows := make(chan *Row, 1)
	rows <- row
	close(rows)

	closeField, err := helper.Field[float64](rows, "Close")
	if err != nil {
		t.Fatal(err)
	}

	closeValue := <-closeField

	if closeValue != row.Close {
		t.Fatalf("actual %v expected %v", closeValue, row.Close)
	}
}

func TestFieldNotStruct(t *testing.T) {
	c := make(chan *float64)
	close(c)

	_, err := helper.Field[float64](c, "Name")
	if err == nil {
		t.Fatal("expecting error")
	}
}

func TestFieldUnknownName(t *testing.T) {
	type Row struct {
		Open  float64
		Close float64
	}

	c := make(chan *Row)
	close(c)

	_, err := helper.Field[float64](c, "High")
	if err == nil {
		t.Fatal("expecting error")
	}
}
