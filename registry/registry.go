// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// Package registry provides a name-based catalog of technical indicators,
// so that a caller who only has an indicator's name, a set of string
// parameters, and a stream of asset snapshots can compute it without
// knowing the underlying Go type. It exists to give the CLI and the MCP
// server a single, shared way to look up and run an indicator by name.
package registry

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cinar/indicator/v2/asset"
	"github.com/cinar/indicator/v2/helper"
)

// Field identifies which OHLCV series of an asset snapshot an indicator
// input is drawn from.
type Field int

const (
	// FieldOpen is the opening price.
	FieldOpen Field = iota

	// FieldHigh is the high price.
	FieldHigh

	// FieldLow is the low price.
	FieldLow

	// FieldClose is the closing price.
	FieldClose

	// FieldVolume is the traded volume.
	FieldVolume
)

// Series is one named output stream produced by an indicator, such as the
// "macd" and "signal" lines of a MACD, or the single "sma" line of an SMA.
type Series struct {
	// Name identifies the series, e.g. "sma", "upper", "signal".
	Name string

	// Values is the computed series, aligned with Result.Dates.
	Values <-chan float64
}

// Result is the outcome of running an indicator: the dates aligned with
// each Series, and the series themselves.
type Result struct {
	// Dates is aligned with every channel in Series: the Nth date
	// corresponds to the Nth value read off each series.
	Dates <-chan time.Time

	// Series are the indicator's named output streams.
	Series []Series
}

// idlePerioder is implemented by indicators that need an initial run of
// snapshots before they can yield a value. Its result is used to skip the
// same number of leading dates, so Result.Dates stays aligned with Series.
type idlePerioder interface {
	IdlePeriod() int
}

// Definition describes how to build and run one named indicator.
type Definition struct {
	// Fields lists the OHLCV inputs the indicator needs, in the order
	// its Compute function expects them.
	Fields []Field

	// New returns a fresh indicator instance configured with its
	// default parameters.
	New func() any

	// Compute runs the indicator against inputs (one channel per
	// Fields entry, in the same order) and returns its named output
	// series.
	Compute func(ctx context.Context, indicator any, inputs []<-chan float64) ([]Series, error)
}

// definitions is the name-to-indicator catalog. Entries are added in
// definitions.go, grouped by the indicator category they come from.
var definitions = map[string]Definition{}

// Register adds or replaces the definition for the given indicator name.
// It is exported so that callers outside this package can extend the
// catalog without needing to fork it.
func Register(name string, definition Definition) {
	definitions[name] = definition
}

// Names returns the registered indicator names in sorted order.
func Names() []string {
	names := make([]string, 0, len(definitions))
	for name := range definitions {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Get returns the definition registered for the given indicator name.
func Get(name string) (Definition, bool) {
	definition, ok := definitions[name]
	return definition, ok
}

// Run looks up the named indicator, applies the given parameters to a
// fresh instance, and computes it over the given snapshots. Parameter
// names may use a dotted path to reach a nested indicator, e.g. "Ema1.Period"
// on a MACD.
func Run(ctx context.Context, name string, params map[string]string, snapshots <-chan *asset.Snapshot) (*Result, error) {
	definition, ok := Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown indicator: %s", name)
	}

	indicator := definition.New()

	for param, value := range params {
		if err := SetParam(indicator, param, value); err != nil {
			return nil, fmt.Errorf("param %s: %w", param, err)
		}
	}

	// Materialize the snapshots once, then give every consumer (the dates
	// and each Field) its own independent producer reading from the slice,
	// rather than a single lock-step fan-out shared between them. Most
	// indicators are lazy and only pull from their input as their own
	// output is consumed, but a few (e.g. momentum.Fisher) read their
	// entire input channel before returning anything; sharing one
	// synchronized fan-out with those would deadlock, since nothing
	// consumes the dates channel until Run itself returns.
	data := helper.ChanToSlice(snapshots)

	dates := asset.SnapshotsAsDatesWithContext(ctx, helper.SliceToChanWithContext(ctx, data))
	if idler, ok := indicator.(idlePerioder); ok {
		dates = helper.SkipWithContext(ctx, dates, idler.IdlePeriod())
	}

	inputs := make([]<-chan float64, len(definition.Fields))
	for i, field := range definition.Fields {
		inputs[i] = fieldChannel(ctx, helper.SliceToChanWithContext(ctx, data), field)
	}

	series, err := definition.Compute(ctx, indicator, inputs)
	if err != nil {
		return nil, err
	}

	return &Result{Dates: dates, Series: series}, nil
}

// fieldChannel extracts the requested OHLCV field from a channel of
// snapshots.
func fieldChannel(ctx context.Context, snapshots <-chan *asset.Snapshot, field Field) <-chan float64 {
	switch field {
	case FieldOpen:
		return asset.SnapshotsAsOpeningsWithContext(ctx, snapshots)
	case FieldHigh:
		return asset.SnapshotsAsHighsWithContext(ctx, snapshots)
	case FieldLow:
		return asset.SnapshotsAsLowsWithContext(ctx, snapshots)
	case FieldClose:
		return asset.SnapshotsAsClosingsWithContext(ctx, snapshots)
	case FieldVolume:
		return asset.SnapshotsAsVolumesWithContext(ctx, snapshots)
	default:
		return asset.SnapshotsAsClosingsWithContext(ctx, snapshots)
	}
}

// SetParam assigns the named parameter its string value on the given
// indicator, using reflection. name may be a dotted path (e.g.
// "Ema1.Period") to reach a field on a nested indicator pointer.
func SetParam(indicator any, name, value string) error {
	v := reflect.ValueOf(indicator)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("indicator must be a non-nil pointer, got %T", indicator)
	}

	v = v.Elem()

	parts := strings.Split(name, ".")

	for i, part := range parts {
		if v.Kind() != reflect.Struct {
			return fmt.Errorf("param %q: %s is not a struct field", name, part)
		}

		field := v.FieldByName(part)
		if !field.IsValid() {
			return fmt.Errorf("unknown param: %s", name)
		}

		if i == len(parts)-1 {
			return setScalar(field, value)
		}

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return fmt.Errorf("param %q: %s is nil", name, part)
			}

			v = field.Elem()
		} else {
			v = field
		}
	}

	return nil
}

// setScalar parses value and assigns it to the given settable field.
func setScalar(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("field is not settable")
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		field.SetInt(n)

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		field.SetFloat(n)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		field.SetBool(b)

	case reflect.String:
		field.SetString(value)

	default:
		return fmt.Errorf("unsupported param kind: %s", field.Kind())
	}

	return nil
}
