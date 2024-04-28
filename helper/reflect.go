// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"fmt"
	"math/bits"
	"reflect"
	"strconv"
	"time"
)

// kindToBits maps numeric kinds (e.g., int, float) to their corresponding sizes in bits.
var kindToBits = map[reflect.Kind]int{
	reflect.Int:     strconv.IntSize,
	reflect.Int8:    8,
	reflect.Int16:   16,
	reflect.Int32:   32,
	reflect.Int64:   64,
	reflect.Uint:    bits.UintSize,
	reflect.Uint16:  16,
	reflect.Uint32:  32,
	reflect.Uint64:  64,
	reflect.Float32: 32,
	reflect.Float64: 64,
}

// setReflectValueFromBool assigns the parsed boolean value to the specified variable.
func setReflectValueFromBool(value reflect.Value, stringValue string) error {
	actualValue, err := strconv.ParseBool(stringValue)
	if err == nil {
		value.SetBool(actualValue)
	}

	return err
}

// setReflectValueFromInt assigns the parsed integer value to the specified variable.
func setReflectValueFromInt(value reflect.Value, stringValue string, bitSize int) error {
	actualValue, err := strconv.ParseInt(stringValue, 10, bitSize)
	if err == nil {
		value.SetInt(actualValue)
	}

	return err
}

// setReflectValueFromUint assigns the parsed unsigned integer value to the specified variable.
func setReflectValueFromUint(value reflect.Value, stringValue string, bitSize int) error {
	actualValue, err := strconv.ParseUint(stringValue, 10, bitSize)
	if err == nil {
		value.SetUint(actualValue)
	}

	return err
}

// setReflectValueFromFloat assigns the parsed float value to the specified variable.
func setReflectValueFromFloat(value reflect.Value, stringValue string, bitSize int) error {
	actualValue, err := strconv.ParseFloat(stringValue, bitSize)
	if err == nil {
		value.SetFloat(actualValue)
	}

	return err
}

// setReflectValueFromTime assigns the parsed unsigned float value to the specified variable.
func setReflectValueFromTime(value reflect.Value, stringValue string, format string) error {
	actualValue, err := time.Parse(format, stringValue)
	if err == nil {
		value.Set(reflect.ValueOf(actualValue))
	}

	return err
}

// setReflectValue assigns the parsed value to the specified variable.
func setReflectValue(value reflect.Value, stringValue string, format string) error {
	kind := value.Kind()

	switch kind {
	case reflect.String:
		value.SetString(stringValue)
		return nil

	case reflect.Bool:
		return setReflectValueFromBool(value, stringValue)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setReflectValueFromInt(value, stringValue, kindToBits[kind])

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setReflectValueFromUint(value, stringValue, kindToBits[kind])

	case reflect.Float32, reflect.Float64:
		return setReflectValueFromFloat(value, stringValue, kindToBits[kind])

	case reflect.Struct:
		typeString := value.Type().String()

		switch typeString {
		case "time.Time":
			return setReflectValueFromTime(value, stringValue, format)

		default:
			return fmt.Errorf("unsupported struct type %s", typeString)
		}

	default:
		return fmt.Errorf("unsupported value kind %s", kind)
	}
}

// getReflectValue returns the string representation of the given value.
func getReflectValue(value reflect.Value, format string) (string, error) {
	kind := value.Kind()

	switch kind {
	case reflect.String:
		return value.String(), nil

	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, kindToBits[kind]), nil

	case reflect.Struct:
		typeString := value.Type().String()

		switch typeString {
		case "time.Time":
			return value.Interface().(time.Time).Format(format), nil

		default:
			return "", fmt.Errorf("unsupported struct type %s", typeString)
		}

	default:
		return "", fmt.Errorf("unsupported value kind %s", kind)
	}
}
