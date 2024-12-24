// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"reflect"
	"testing"
	"time"
)

func TestSetReflectValueFromString(t *testing.T) {
	actual := ""
	value := reflect.ValueOf(&actual).Elem()
	expected := "string value"

	err := setReflectValue(value, expected, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromBool(t *testing.T) {
	actual := false
	value := reflect.ValueOf(&actual).Elem()
	expected := true

	err := setReflectValue(value, "true", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotBool(t *testing.T) {
	actual := false
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromInt(t *testing.T) {
	actual := 0
	value := reflect.ValueOf(&actual).Elem()
	expected := 10

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotInt(t *testing.T) {
	actual := 0
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromInt8(t *testing.T) {
	actual := int8(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := int8(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotInt8(t *testing.T) {
	actual := int8(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromInt16(t *testing.T) {
	actual := int16(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := int16(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotInt16(t *testing.T) {
	actual := int16(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromInt32(t *testing.T) {
	actual := int32(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := int32(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotInt32(t *testing.T) {
	actual := int32(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromInt64(t *testing.T) {
	actual := int64(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := int64(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotInt64(t *testing.T) {
	actual := int64(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromUint(t *testing.T) {
	actual := uint(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := uint(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotUint(t *testing.T) {
	actual := 0
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromUint8(t *testing.T) {
	actual := uint8(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := uint8(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotUint8(t *testing.T) {
	actual := uint8(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromUint16(t *testing.T) {
	actual := uint16(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := uint16(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotUint16(t *testing.T) {
	actual := uint16(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromUint32(t *testing.T) {
	actual := uint32(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := uint32(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotUint32(t *testing.T) {
	actual := uint32(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromUint64(t *testing.T) {
	actual := uint64(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := uint64(10)

	err := setReflectValue(value, "10", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotUint64(t *testing.T) {
	actual := uint64(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromFloat32(t *testing.T) {
	actual := float32(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := float32(10.20)

	err := setReflectValue(value, "10.20", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotFloat32(t *testing.T) {
	actual := float32(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromFloat64(t *testing.T) {
	actual := float64(0)
	value := reflect.ValueOf(&actual).Elem()
	expected := 10.20

	err := setReflectValue(value, "10.20", "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotFloat64(t *testing.T) {
	actual := float64(0)
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromTime(t *testing.T) {
	var actual time.Time
	value := reflect.ValueOf(&actual).Elem()
	expected := time.Date(2023, 11, 28, 19, 14, 0, 0, time.UTC)

	err := setReflectValue(value, "2023-11-28 19:14:00", "2006-01-02 15:04:05")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestSetReflectValueFromNotTime(t *testing.T) {
	var actual time.Time
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "abcd", "2006-01-02 15:04:05")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromUnexpectedStuct(t *testing.T) {
	type Unexpected struct{}
	var actual Unexpected
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "2023-11-28 19:14:00", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestSetReflectValueFromUnexpectedKind(t *testing.T) {
	type Unexpected struct{}
	var actual *Unexpected
	value := reflect.ValueOf(&actual).Elem()

	err := setReflectValue(value, "2023-11-28 19:14:00", "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestGetReflectValueFromString(t *testing.T) {
	expected := "string value"
	value := reflect.ValueOf(&expected).Elem()

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromBool(t *testing.T) {
	input := true
	value := reflect.ValueOf(&input).Elem()
	expected := "true"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromInt(t *testing.T) {
	input := 10
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromInt8(t *testing.T) {
	input := int8(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromInt16(t *testing.T) {
	input := int16(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromInt32(t *testing.T) {
	input := int32(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromInt64(t *testing.T) {
	input := int64(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromUint(t *testing.T) {
	input := uint(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromUint8(t *testing.T) {
	input := uint8(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromUint16(t *testing.T) {
	input := uint16(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromUint32(t *testing.T) {
	input := uint32(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromUint64(t *testing.T) {
	input := uint64(10)
	value := reflect.ValueOf(&input).Elem()
	expected := "10"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromFloat32(t *testing.T) {
	input := float32(10.20)
	value := reflect.ValueOf(&input).Elem()
	expected := "10.2"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromFloat64(t *testing.T) {
	input := 10.20
	value := reflect.ValueOf(&input).Elem()
	expected := "10.2"

	actual, err := getReflectValue(value, "")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromTime(t *testing.T) {
	input := time.Date(2023, 11, 28, 19, 14, 0, 0, time.UTC)
	value := reflect.ValueOf(&input).Elem()
	expected := "2023-11-28 19:14:00"

	actual, err := getReflectValue(value, "2006-01-02 15:04:05")
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

func TestGetReflectValueFromUnexpectedStuct(t *testing.T) {
	type Unexpected struct{}
	var input Unexpected
	value := reflect.ValueOf(&input).Elem()

	actual, err := getReflectValue(value, "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}

func TestGetReflectValueFromUnexpectedKind(t *testing.T) {
	type Unexpected struct{}
	var input *Unexpected
	value := reflect.ValueOf(&input).Elem()

	actual, err := getReflectValue(value, "")
	if err == nil {
		t.Fatalf("actual %v expected error", actual)
	}
}
