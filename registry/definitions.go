// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package registry

import (
	"context"

	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/trend"
	"github.com/cinar/indicator/v2/volatility"
	"github.com/cinar/indicator/v2/volume"
)

// init registers the built-in indicator catalog. It is a representative
// slice across trend, momentum, volatility, and volume indicators, and
// across single/multi input and single/multi output shapes, meant to be
// extended with more indicators over time following the same pattern.
func init() {
	Register("sma", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewSma[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			sma := indicator.(*trend.Sma[float64])
			return []Series{
				{Name: "sma", Values: sma.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("ema", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewEma[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			ema := indicator.(*trend.Ema[float64])
			return []Series{
				{Name: "ema", Values: ema.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("macd", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewMacd[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			macd := indicator.(*trend.Macd[float64])
			macdLine, signal := macd.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "macd", Values: macdLine},
				{Name: "signal", Values: signal},
			}, nil
		},
	})

	Register("cci", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return trend.NewCci[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			cci := indicator.(*trend.Cci[float64])
			return []Series{
				{Name: "cci", Values: cci.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("rsi", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return momentum.NewRsi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			rsi := indicator.(*momentum.Rsi[float64])
			return []Series{
				{Name: "rsi", Values: rsi.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("bollinger_bands", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewBollingerBands[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			bollingerBands := indicator.(*volatility.BollingerBands[float64])
			upper, middle, lower := bollingerBands.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "upper", Values: upper},
				{Name: "middle", Values: middle},
				{Name: "lower", Values: lower},
			}, nil
		},
	})

	Register("obv", Definition{
		Fields: []Field{FieldClose, FieldVolume},
		New:    func() any { return volume.NewObv[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			obv := indicator.(*volume.Obv[float64])
			return []Series{
				{Name: "obv", Values: obv.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})

	Register("vwap", Definition{
		Fields: []Field{FieldClose, FieldVolume},
		New:    func() any { return volume.NewVwap[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			vwap := indicator.(*volume.Vwap[float64])
			return []Series{
				{Name: "vwap", Values: vwap.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})
}
