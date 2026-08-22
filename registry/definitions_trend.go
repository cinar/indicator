// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package registry

import (
	"context"

	"github.com/cinar/indicator/v2/trend"
)

// init registers the trend indicators.
//
// A handful of trend package types are deliberately not registered:
//
//   - MovingMax, MovingMin, MovingSum, and the generic single-series
//     Stochastic are building blocks other indicators compose (e.g. Sma
//     uses MovingSum), not indicators traded on their own.
//   - Envelope takes a caller-supplied Ma implementation instead of having
//     its own default constructor; there's no single obvious default to
//     pick on its behalf.
//   - Mlr and Mls fit arbitrary (x, y) series, not a fixed OHLCV field, so
//     there's no natural default input to bind them to.
//   - PivotPoint returns a PivotPointResult struct per value instead of a
//     plain number, which doesn't fit the Series shape every other
//     indicator here uses.
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

	Register("apo", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewApo[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			apo := indicator.(*trend.Apo[float64])
			return []Series{
				{Name: "apo", Values: apo.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("aroon", Definition{
		Fields: []Field{FieldHigh, FieldLow},
		New:    func() any { return trend.NewAroon[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			aroon := indicator.(*trend.Aroon[float64])
			up, down := aroon.ComputeWithContext(ctx, in[0], in[1])
			return []Series{
				{Name: "up", Values: up},
				{Name: "down", Values: down},
			}, nil
		},
	})

	Register("bop", Definition{
		Fields: []Field{FieldOpen, FieldHigh, FieldLow, FieldClose},
		New:    func() any { return trend.NewBop[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			bop := indicator.(*trend.Bop[float64])
			return []Series{
				{Name: "bop", Values: bop.ComputeWithContext(ctx, in[0], in[1], in[2], in[3])},
			}, nil
		},
	})

	Register("cfo", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewCfo[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			cfo := indicator.(*trend.Cfo[float64])
			return []Series{
				{Name: "cfo", Values: cfo.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("dema", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewDema[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			dema := indicator.(*trend.Dema[float64])
			return []Series{
				{Name: "dema", Values: dema.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("dpo", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewDpo[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			dpo := indicator.(*trend.Dpo[float64])
			return []Series{
				{Name: "dpo", Values: dpo.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("hma", Definition{
		Fields: []Field{FieldClose},
		// Hma, like Wma, has no zero-arg constructor; 14 is a commonly used
		// default period, override with -param Period=.
		New: func() any { return trend.NewHmaWithPeriod[float64](14) },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			hma := indicator.(*trend.Hma[float64])
			return []Series{
				{Name: "hma", Values: hma.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("kama", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewKama[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			kama := indicator.(*trend.Kama[float64])
			return []Series{
				{Name: "kama", Values: kama.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("kdj", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return trend.NewKdj[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			kdj := indicator.(*trend.Kdj[float64])
			k, d, j := kdj.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "k", Values: k},
				{Name: "d", Values: d},
				{Name: "j", Values: j},
			}, nil
		},
	})

	Register("kst", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewKst[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			kst := indicator.(*trend.Kst[float64])
			kstLine, signal := kst.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "kst", Values: kstLine},
				{Name: "signal", Values: signal},
			}, nil
		},
	})

	Register("mass_index", Definition{
		Fields: []Field{FieldHigh, FieldLow},
		New:    func() any { return trend.NewMassIndex[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			massIndex := indicator.(*trend.MassIndex[float64])
			return []Series{
				{Name: "mass_index", Values: massIndex.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})

	Register("mcginley_dynamic", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewMcGinleyDynamic[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			mcGinleyDynamic := indicator.(*trend.McGinleyDynamic[float64])
			return []Series{
				{Name: "mcginley_dynamic", Values: mcGinleyDynamic.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("rma", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewRma[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			rma := indicator.(*trend.Rma[float64])
			return []Series{
				{Name: "rma", Values: rma.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("roc", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewRoc[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			roc := indicator.(*trend.Roc[float64])
			return []Series{
				{Name: "roc", Values: roc.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("slope", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewSlope[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			slope := indicator.(*trend.Slope[float64])
			return []Series{
				{Name: "slope", Values: slope.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("slow_stochastic", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewSlowStochastic[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			slowStochastic := indicator.(*trend.SlowStochastic[float64])
			k, d := slowStochastic.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "k", Values: k},
				{Name: "d", Values: d},
			}, nil
		},
	})

	Register("smma", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewSmma[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			smma := indicator.(*trend.Smma[float64])
			return []Series{
				{Name: "smma", Values: smma.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	// Stc (Schaff Trend Cycle) is still not registered: #421 fixed the
	// output channel never closing, but the computed values themselves
	// are wrong -- STC = 100*(MACD-%K)/(%D-%K) divides by a denominator
	// that's frequently near zero, so real-world input (verified against
	// asset/testdata/repository/brk-b.csv) yields -Inf and swings like
	// 1435 or -710 for an indicator that's documented to oscillate
	// 0-100. This is a distinct, still-open bug in trend.Stc's formula,
	// not the one #421 fixed.

	// T3 is still not registered: #423 fixed ComputeWithContext yielding
	// fewer values than IdlePeriod() promised, but the weighted-sum
	// coefficients themselves don't sum to 1 (c1+c2+c3+c4 = 3a(a-1), not
	// 1 as a moving average requires), so T3 doesn't track price at all
	// -- verified by feeding a constant 100 closings and getting -63
	// back, not 100. This is a distinct, still-open bug in T3's formula
	// (see the standard Tillson T3 coefficients: c2 should include a
	// +3a^3 term and c3 a -6a^2-3a^3 term, both missing here), not the
	// one #423 fixed.

	Register("tema", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewTema[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			tema := indicator.(*trend.Tema[float64])
			return []Series{
				{Name: "tema", Values: tema.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("trima", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewTrima[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			trima := indicator.(*trend.Trima[float64])
			return []Series{
				{Name: "trima", Values: trima.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("trix", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewTrix[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			trix := indicator.(*trend.Trix[float64])
			return []Series{
				{Name: "trix", Values: trix.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("tsi", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return trend.NewTsi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			tsi := indicator.(*trend.Tsi[float64])
			return []Series{
				{Name: "tsi", Values: tsi.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("typical_price", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return trend.NewTypicalPrice[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			typicalPrice := indicator.(*trend.TypicalPrice[float64])
			return []Series{
				{Name: "typical_price", Values: typicalPrice.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("vwma", Definition{
		Fields: []Field{FieldClose, FieldVolume},
		New:    func() any { return trend.NewVwma[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			vwma := indicator.(*trend.Vwma[float64])
			return []Series{
				{Name: "vwma", Values: vwma.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})

	Register("weighted_close", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return trend.NewWeightedClose[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			weightedClose := indicator.(*trend.WeightedClose[float64])
			return []Series{
				{Name: "weighted_close", Values: weightedClose.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("wma", Definition{
		Fields: []Field{FieldClose},
		// Wma has no zero-arg constructor (NewWmaWith panics on period <= 0),
		// so the registry picks the same 10-period default many WMA
		// implementations use; override with -param Period=.
		New: func() any { return trend.NewWmaWith[float64](10) },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			wma := indicator.(*trend.Wma[float64])
			return []Series{
				{Name: "wma", Values: wma.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})
}
