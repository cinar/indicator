// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package registry

import (
	"context"

	"github.com/cinar/indicator/v2/volatility"
)

// init registers the volatility indicators.
//
// volatility.MovingStd is not registered for the same reason as trend's
// MovingSum: BollingerBands and ZScore compose it, it isn't traded on its
// own.
func init() {
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

	Register("acceleration_bands", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewAccelerationBands[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			accelerationBands := indicator.(*volatility.AccelerationBands[float64])
			upper, middle, lower := accelerationBands.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "upper", Values: upper},
				{Name: "middle", Values: middle},
				{Name: "lower", Values: lower},
			}, nil
		},
	})

	Register("annualized_historical_volatility", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewAnnualizedHistoricalVolatility[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			annualizedHistoricalVolatility := indicator.(*volatility.AnnualizedHistoricalVolatility[float64])
			return []Series{
				{Name: "annualized_historical_volatility", Values: annualizedHistoricalVolatility.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("atr", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewAtr[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			atr := indicator.(*volatility.Atr[float64])
			return []Series{
				{Name: "atr", Values: atr.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("bollinger_band_width", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewBollingerBandWidth[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			bollingerBandWidth := indicator.(*volatility.BollingerBandWidth[float64])
			return []Series{
				{Name: "bollinger_band_width", Values: bollingerBandWidth.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("chandelier_exit", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewChandelierExit[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			chandelierExit := indicator.(*volatility.ChandelierExit[float64])
			long, short := chandelierExit.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "long", Values: long},
				{Name: "short", Values: short},
			}, nil
		},
	})

	Register("chop", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewChop[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			chop := indicator.(*volatility.Chop[float64])
			return []Series{
				{Name: "chop", Values: chop.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("donchian_channel", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewDonchianChannel[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			donchianChannel := indicator.(*volatility.DonchianChannel[float64])
			upper, middle, lower := donchianChannel.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "upper", Values: upper},
				{Name: "middle", Values: middle},
				{Name: "lower", Values: lower},
			}, nil
		},
	})

	Register("historical_volatility", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewHistoricalVolatility[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			historicalVolatility := indicator.(*volatility.HistoricalVolatility[float64])
			return []Series{
				{Name: "historical_volatility", Values: historicalVolatility.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("keltner_channel", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewKeltnerChannel[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			keltnerChannel := indicator.(*volatility.KeltnerChannel[float64])
			upper, middle, lower := keltnerChannel.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "upper", Values: upper},
				{Name: "middle", Values: middle},
				{Name: "lower", Values: lower},
			}, nil
		},
	})

	Register("percent_b", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewPercentB[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			percentB := indicator.(*volatility.PercentB[float64])
			return []Series{
				{Name: "percent_b", Values: percentB.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("po", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewPo[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			po := indicator.(*volatility.Po[float64])
			return []Series{
				{Name: "po", Values: po.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("super_trend", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewSuperTrend[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			superTrend := indicator.(*volatility.SuperTrend[float64])
			return []Series{
				{Name: "super_trend", Values: superTrend.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("true_range", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volatility.NewTrueRange[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			trueRange := indicator.(*volatility.TrueRange[float64])
			return []Series{
				{Name: "true_range", Values: trueRange.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("ulcer_index", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewUlcerIndex[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			ulcerIndex := indicator.(*volatility.UlcerIndex[float64])
			return []Series{
				{Name: "ulcer_index", Values: ulcerIndex.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("z_score", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return volatility.NewZScore[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			zScore := indicator.(*volatility.ZScore[float64])
			return []Series{
				{Name: "z_score", Values: zScore.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})
}
