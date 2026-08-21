// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package registry

import (
	"context"

	"github.com/cinar/indicator/v2/momentum"
)

// init registers the momentum indicators.
//
// momentum.Streak, ConnorsRsi's internal building block, is not registered
// for the same reason as trend's MovingSum: it's a component other
// indicators compose, not one traded on its own.
func init() {
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

	Register("awesome_oscillator", Definition{
		Fields: []Field{FieldHigh, FieldLow},
		New:    func() any { return momentum.NewAwesomeOscillator[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			awesomeOscillator := indicator.(*momentum.AwesomeOscillator[float64])
			return []Series{
				{Name: "awesome_oscillator", Values: awesomeOscillator.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})

	Register("chaikin_oscillator", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose, FieldVolume},
		New:    func() any { return momentum.NewChaikinOscillator[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			chaikinOscillator := indicator.(*momentum.ChaikinOscillator[float64])
			co, ad := chaikinOscillator.ComputeWithContext(ctx, in[0], in[1], in[2], in[3])
			return []Series{
				{Name: "chaikin_oscillator", Values: co},
				{Name: "ad", Values: ad},
			}, nil
		},
	})

	Register("connors_rsi", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return momentum.NewConnorsRsi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			connorsRsi := indicator.(*momentum.ConnorsRsi[float64])
			return []Series{
				{Name: "connors_rsi", Values: connorsRsi.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("coppock_curve", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return momentum.NewCoppockCurve[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			coppockCurve := indicator.(*momentum.CoppockCurve[float64])
			return []Series{
				{Name: "coppock_curve", Values: coppockCurve.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("elder_ray", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return momentum.NewElderRay[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			elderRay := indicator.(*momentum.ElderRay[float64])
			bullPower, bearPower := elderRay.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "bull_power", Values: bullPower},
				{Name: "bear_power", Values: bearPower},
			}, nil
		},
	})

	// Fisher is not registered: its ComputeWithContext produces far fewer
	// values than IdlePeriod() promises (verified independently of this
	// registry -- 400 closings in, IdlePeriod says 382 values should come
	// out, only 195 actually do), so Dates and its Series can't be
	// aligned. Worth its own issue and fix in momentum.Fisher before
	// registering it here.

	Register("ibs", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return momentum.NewInternalBarStrength[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			ibs := indicator.(*momentum.InternalBarStrength[float64])
			return []Series{
				{Name: "ibs", Values: ibs.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("ichimoku_cloud", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return momentum.NewIchimokuCloud[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			ichimokuCloud := indicator.(*momentum.IchimokuCloud[float64])
			conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingSpan := ichimokuCloud.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "conversion_line", Values: conversionLine},
				{Name: "base_line", Values: baseLine},
				{Name: "leading_span_a", Values: leadingSpanA},
				{Name: "leading_span_b", Values: leadingSpanB},
				{Name: "lagging_span", Values: laggingSpan},
			}, nil
		},
	})

	Register("ppo", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return momentum.NewPpo[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			ppo := indicator.(*momentum.Ppo[float64])
			line, signal, histogram := ppo.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "ppo", Values: line},
				{Name: "signal", Values: signal},
				{Name: "histogram", Values: histogram},
			}, nil
		},
	})

	// PringsSpecialK is not registered: unlike its peers it has no
	// IdlePeriod() method, so registry.Run can't skip the leading dates
	// its (very long, 500+ period) internal ROC/SMA chain needs before it
	// yields a value -- Run would treat every date as immediately valid
	// and the series/dates would come out misaligned. Its own test works
	// around this the same way, computing the required length by hand
	// instead of calling an IdlePeriod(). Worth an IdlePeriod() method on
	// PringsSpecialK itself before registering it here.

	Register("pvo", Definition{
		Fields: []Field{FieldVolume},
		New:    func() any { return momentum.NewPvo[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			pvo := indicator.(*momentum.Pvo[float64])
			line, signal, histogram := pvo.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "pvo", Values: line},
				{Name: "signal", Values: signal},
				{Name: "histogram", Values: histogram},
			}, nil
		},
	})

	Register("qstick", Definition{
		Fields: []Field{FieldOpen, FieldClose},
		New:    func() any { return momentum.NewQstick[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			qstick := indicator.(*momentum.Qstick[float64])
			return []Series{
				{Name: "qstick", Values: qstick.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})

	Register("rvi", Definition{
		Fields: []Field{FieldOpen, FieldHigh, FieldLow, FieldClose},
		New:    func() any { return momentum.NewRvi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			rvi := indicator.(*momentum.Rvi[float64])
			line, signal := rvi.ComputeWithContext(ctx, in[0], in[1], in[2], in[3])
			return []Series{
				{Name: "rvi", Values: line},
				{Name: "signal", Values: signal},
			}, nil
		},
	})

	Register("stochastic_oscillator", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return momentum.NewStochasticOscillator[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			stochasticOscillator := indicator.(*momentum.StochasticOscillator[float64])
			k, d := stochasticOscillator.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "k", Values: k},
				{Name: "d", Values: d},
			}, nil
		},
	})

	Register("stochastic_rsi", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return momentum.NewStochasticRsi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			stochasticRsi := indicator.(*momentum.StochasticRsi[float64])
			return []Series{
				{Name: "stochastic_rsi", Values: stochasticRsi.ComputeWithContext(ctx, in[0])},
			}, nil
		},
	})

	Register("td_sequential", Definition{
		Fields: []Field{FieldClose},
		New:    func() any { return momentum.NewTdSequential[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			tdSequential := indicator.(*momentum.TdSequential[float64])
			buySetup, sellSetup, buyCountdown, sellCountdown := tdSequential.ComputeWithContext(ctx, in[0])
			return []Series{
				{Name: "buy_setup", Values: buySetup},
				{Name: "sell_setup", Values: sellSetup},
				{Name: "buy_countdown", Values: buyCountdown},
				{Name: "sell_countdown", Values: sellCountdown},
			}, nil
		},
	})

	Register("ultimate_oscillator", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return momentum.NewUltimateOscillator[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			ultimateOscillator := indicator.(*momentum.UltimateOscillator[float64])
			return []Series{
				{Name: "ultimate_oscillator", Values: ultimateOscillator.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("williams_r", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return momentum.NewWilliamsR[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			williamsR := indicator.(*momentum.WilliamsR[float64])
			return []Series{
				{Name: "williams_r", Values: williamsR.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})
}
