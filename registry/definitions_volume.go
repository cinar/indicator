// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package registry

import (
	"context"

	"github.com/cinar/indicator/v2/volume"
)

// init registers the volume indicators.
func init() {
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

	Register("ad", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose, FieldVolume},
		New:    func() any { return volume.NewAd[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			ad := indicator.(*volume.Ad[float64])
			return []Series{
				{Name: "ad", Values: ad.ComputeWithContext(ctx, in[0], in[1], in[2], in[3])},
			}, nil
		},
	})

	Register("cmf", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose, FieldVolume},
		New:    func() any { return volume.NewCmf[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			cmf := indicator.(*volume.Cmf[float64])
			return []Series{
				{Name: "cmf", Values: cmf.ComputeWithContext(ctx, in[0], in[1], in[2], in[3])},
			}, nil
		},
	})

	Register("emv", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldVolume},
		New:    func() any { return volume.NewEmv[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			emv := indicator.(*volume.Emv[float64])
			return []Series{
				{Name: "emv", Values: emv.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("fi", Definition{
		Fields: []Field{FieldClose, FieldVolume},
		New:    func() any { return volume.NewFi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			fi := indicator.(*volume.Fi[float64])
			return []Series{
				{Name: "fi", Values: fi.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})

	Register("kvo", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldVolume},
		New:    func() any { return volume.NewKvo[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			kvo := indicator.(*volume.Kvo[float64])
			line, signal := kvo.ComputeWithContext(ctx, in[0], in[1], in[2])
			return []Series{
				{Name: "kvo", Values: line},
				{Name: "signal", Values: signal},
			}, nil
		},
	})

	Register("mfi", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose, FieldVolume},
		New:    func() any { return volume.NewMfi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			mfi := indicator.(*volume.Mfi[float64])
			return []Series{
				{Name: "mfi", Values: mfi.ComputeWithContext(ctx, in[0], in[1], in[2], in[3])},
			}, nil
		},
	})

	Register("mfm", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose},
		New:    func() any { return volume.NewMfm[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			mfm := indicator.(*volume.Mfm[float64])
			return []Series{
				{Name: "mfm", Values: mfm.ComputeWithContext(ctx, in[0], in[1], in[2])},
			}, nil
		},
	})

	Register("mfv", Definition{
		Fields: []Field{FieldHigh, FieldLow, FieldClose, FieldVolume},
		New:    func() any { return volume.NewMfv[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			mfv := indicator.(*volume.Mfv[float64])
			return []Series{
				{Name: "mfv", Values: mfv.ComputeWithContext(ctx, in[0], in[1], in[2], in[3])},
			}, nil
		},
	})

	Register("nvi", Definition{
		Fields: []Field{FieldClose, FieldVolume},
		New:    func() any { return volume.NewNvi[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			nvi := indicator.(*volume.Nvi[float64])
			return []Series{
				{Name: "nvi", Values: nvi.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})

	Register("vpt", Definition{
		Fields: []Field{FieldClose, FieldVolume},
		New:    func() any { return volume.NewVpt[float64]() },
		Compute: func(ctx context.Context, indicator any, in []<-chan float64) ([]Series, error) {
			vpt := indicator.(*volume.Vpt[float64])
			return []Series{
				{Name: "vpt", Values: vpt.ComputeWithContext(ctx, in[0], in[1])},
			}, nil
		},
	})
}
