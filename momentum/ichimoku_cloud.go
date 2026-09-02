// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum

import (
	"context"
	"fmt"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

const (
	// DefaultIchimokuCloudConversionPeriod is the default conversion period for the Ichimoku Cloud.
	DefaultIchimokuCloudConversionPeriod = 9

	// DefaultIchimokuCloudBasePeriod is the default base period for the Ichimoku Cloud.
	DefaultIchimokuCloudBasePeriod = 26

	// DefaultIchimokuCloudLeadingPeriod is the default leading period for the Ichimoku Cloud.
	DefaultIchimokuCloudLeadingPeriod = 52

	// DefaultIchimokuCloudLaggingPeriod is the default lagging period for the Ichimoku Cloud.
	DefaultIchimokuCloudLaggingPeriod = 26
)

// IchimokuCloud represents the configuration parameter for calculating the Ichimoku Cloud. It is also known as the
// Ichimoku Kinko Hyo, is a versatile indicator that defines support and resistance, identifies trend direction,
// gauges momentum, and provides trading signals.
//
//	Tenkan-sen (Conversion Line) = (9-Period High + 9-Period Low) / 2
//	Kijun-sen (Base Line) = (26-Period High + 26-Period Low) / 2
//	Senkou Span A (Leading Span A) = (Conversion Line + Base Line) / 2, plotted LaggingPeriod periods ahead.
//	Senkou Span B (Leading Span B) = (52-Period High + 52-Period Low) / 2, plotted LaggingPeriod periods ahead.
//	Chikou Span (Lagging Span) = Today's closing, plotted LaggingPeriod periods back on the chart.
//
// Senkou Span A/B and Chikou Span are all displaced on the chart relative to the day they are computed
// from. Since real-world implementations almost universally default both displacements to the same value,
// LaggingPeriod is reused for both: it is the number of periods Senkou Span A/B are shifted forward, and
// also the number of periods Chikou Span is shifted back.
//
// Example:
//
//	ic := momentum.NewIchimokuCloud[float64]()
//	conversionLine, baseLine, leadingSpanA, leasingSpanB, laggingSpan := ic.Compute(highs, lows, closings)
type IchimokuCloud[T helper.Number] struct {
	// ConversionMax is the conversion Moving Max instance.
	ConversionMax *trend.MovingMax[T]

	// ConversionMin is the conversion Moving Min instance.
	ConversionMin *trend.MovingMin[T]

	// BaseMax is the base Moving Max instance.
	BaseMax *trend.MovingMax[T]

	// BaseMin is the base Moving Min instance.
	BaseMin *trend.MovingMin[T]

	// LeadingMax is the leading Moving Max instance.
	LeadingMax *trend.MovingMax[T]

	// LeadingMin is the leading Moving Min instance.
	LeadingMin *trend.MovingMin[T]

	// LaggingPeriod is the lagging period.
	LaggingPeriod int
}

// NewIchimokuCloud function initializes a new Ichimoku Cloud instance.
func NewIchimokuCloud[T helper.Number]() *IchimokuCloud[T] {
	return &IchimokuCloud[T]{
		ConversionMax: trend.NewMovingMaxWithPeriod[T](DefaultIchimokuCloudConversionPeriod),
		ConversionMin: trend.NewMovingMinWithPeriod[T](DefaultIchimokuCloudConversionPeriod),
		BaseMax:       trend.NewMovingMaxWithPeriod[T](DefaultIchimokuCloudBasePeriod),
		BaseMin:       trend.NewMovingMinWithPeriod[T](DefaultIchimokuCloudBasePeriod),
		LeadingMax:    trend.NewMovingMaxWithPeriod[T](DefaultIchimokuCloudLeadingPeriod),
		LeadingMin:    trend.NewMovingMinWithPeriod[T](DefaultIchimokuCloudLeadingPeriod),
		LaggingPeriod: DefaultIchimokuCloudLaggingPeriod,
	}
}

// ComputeWithContext function takes a channel of numbers and computes the Ichimoku Cloud.
// Returns conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingSpan
func (i *IchimokuCloud[T]) ComputeWithContext(ctx context.Context, highs, lows, closings <-chan T) (<-chan T, <-chan T, <-chan T, <-chan T, <-chan T) {
	highsSplice := helper.DuplicateWithContext(ctx, highs, 3)
	lowsSplice := helper.DuplicateWithContext(ctx, lows, 3)

	//	Tenkan-sen (Conversion Line) = (9-Period High + 9-Period Low) / 2
	conversionLineSplice := helper.DuplicateWithContext(ctx, helper.DivideByWithContext(ctx, helper.AddWithContext(ctx, i.ConversionMax.ComputeWithContext(ctx, highsSplice[0]),
		i.ConversionMin.ComputeWithContext(ctx, lowsSplice[0]),
	),
		2,
	),
		2,
	)

	//	Kijun-sen (Base Line) = (26-Period High + 26-Period Low) / 2
	baseLineSplice := helper.DuplicateWithContext(ctx, helper.DivideByWithContext(ctx, helper.AddWithContext(ctx, i.BaseMax.ComputeWithContext(ctx, highsSplice[1]),
		i.BaseMin.ComputeWithContext(ctx, lowsSplice[1]),
	),
		2,
	),
		2,
	)

	conversionLineSplice[0] = helper.SkipWithContext(ctx, conversionLineSplice[0], i.BaseMax.IdlePeriod()-i.ConversionMax.IdlePeriod())
	conversionLineSplice[1] = helper.SkipWithContext(ctx, conversionLineSplice[1], i.BaseMax.IdlePeriod()-i.ConversionMax.IdlePeriod())

	//	Senkou Span A (Leading Span A) = (Conversion Line + Base Line) / 2
	leadingSpanA := helper.DivideByWithContext(ctx, helper.AddWithContext(ctx, conversionLineSplice[0],
		baseLineSplice[0],
	),
		2,
	)

	//	Senkou Span B (Leading Span B) = (52-Period High + 52-Period Low) / 2
	leadingSpanB := helper.DivideByWithContext(ctx, helper.AddWithContext(ctx, i.LeadingMax.ComputeWithContext(ctx, highsSplice[2]),
		i.LeadingMin.ComputeWithContext(ctx, lowsSplice[2]),
	),
		2,
	)

	leadingSpanA = helper.SkipWithContext(ctx, leadingSpanA, i.LeadingMax.IdlePeriod()-i.BaseMax.IdlePeriod())
	conversionLineSplice[1] = helper.SkipWithContext(ctx, conversionLineSplice[1], i.LeadingMax.IdlePeriod()-i.BaseMax.IdlePeriod())
	baseLineSplice[1] = helper.SkipWithContext(ctx, baseLineSplice[1], i.LeadingMax.IdlePeriod()-i.BaseMax.IdlePeriod())

	// At this point conversionLineSplice[1], baseLineSplice[1], leadingSpanA, and leadingSpanB all start
	// emitting at the same absolute input row, idle := i.LeadingMax.IdlePeriod(). Senkou Span A/B (leadingSpanA
	// and leadingSpanB) are supposed to be plotted disp periods ahead on the chart, i.e. the span value
	// computed from data through day idle+k belongs at chart position idle+disp+k. Rather than delaying the
	// span streams themselves, we delay the returned Conversion/Base Line streams by disp instead: this makes
	// the whole output start disp rows later, so output row p corresponds to absolute chart row idle+disp+p.
	// leadingSpanA/leadingSpanB's row p is then still the span computed as of day idle+p (unchanged), which is
	// exactly the value that belongs at chart row idle+disp+p. No front skip on the spans is needed at all.
	disp := i.LaggingPeriod

	conversionLineSplice[1] = helper.SkipWithContext(ctx, conversionLineSplice[1], disp)
	baseLineSplice[1] = helper.SkipWithContext(ctx, baseLineSplice[1], disp)

	//	Chikou Span (Lagging Span) = Today's closing, plotted disp periods back on the chart, i.e. chart row p
	//	shows the closing from row p+disp. Combined with the output starting at absolute row idle+disp (per
	//	above), chart row p needs closings[idle+2*disp+p], so simply skip idle+2*disp raw closings.
	laggingLine := helper.SkipWithContext(ctx, closings, i.LeadingMax.IdlePeriod()+2*disp)

	// Delaying conversionLine/baseLine by disp above leaves them disp elements shorter than the still
	// front-aligned leadingSpanA/leadingSpanB, and laggingLine (skipped by an extra disp again) ends up a
	// further disp elements shorter still. Trim the tails so all five returned channels remain the same
	// length and row-aligned, as required by callers that zip them together (e.g. helper.Operate5WithContext).
	leadingSpanA = helper.SkipLastWithContext(ctx, leadingSpanA, 2*disp)
	leadingSpanB = helper.SkipLastWithContext(ctx, leadingSpanB, 2*disp)
	conversionLineSplice[1] = helper.SkipLastWithContext(ctx, conversionLineSplice[1], disp)
	baseLineSplice[1] = helper.SkipLastWithContext(ctx, baseLineSplice[1], disp)

	return conversionLineSplice[1], baseLineSplice[1], leadingSpanA, leadingSpanB, laggingLine
}

// IdlePeriod is the initial period that Ichimoku Cloud won't yield any results.
func (i *IchimokuCloud[T]) IdlePeriod() int {
	return i.LeadingMax.IdlePeriod() + i.LaggingPeriod
}

// String is the string representation of the Ichimoku Cloud.
func (i *IchimokuCloud[T]) String() string {
	return fmt.Sprintf("ICHIMOKU(%d,%d,%d,%d)", i.ConversionMax.Period, i.BaseMax.Period, i.LeadingMax.Period, i.LaggingPeriod)
}

// Compute wraps ComputeWithContext for backwards compatibility.
//
// Deprecated: Use ComputeWithContext instead.
func (i *IchimokuCloud[T]) Compute(highs, lows, closings <-chan T) (<-chan T, <-chan T, <-chan T, <-chan T, <-chan T) {
	return i.ComputeWithContext(context.Background(), highs, lows, closings)
}
