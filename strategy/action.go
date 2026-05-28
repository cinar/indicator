// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import (
	"context"

	"github.com/cinar/indicator/v2/helper"
)

// Action represents the different action categories that a
// strategy can recommend.
type Action int

const (
	// Hold suggests maintaining the current position and not
	// taking any actions on the asset.
	Hold Action = 0

	// Sell suggests disposing of the asset and exiting the current position.
	// This recommendation typically indicates that the strategy believes the
	// asset's price has reached its peak or is likely to decline.
	Sell Action = -1

	// Buy suggests acquiring the asset and entering a new position. This
	// recommendation usually implies that the strategy believes the
	// asset's price is undervalued.
	Buy Action = 1
)

// Annotation returns a single character string representing the recommended action.
// It returns "S" for Sell, "B" for Buy, and an empty string for Hold.
func (a Action) Annotation() string {
	switch a {
	case Sell:
		return "S"

	case Buy:
		return "B"

	default:
		return ""
	}
}

// ActionsToAnnotationsWithContext takes a channel of action recommendations and returns a
// new channel containing corresponding annotations for those actions, supporting context cancellation.
func ActionsToAnnotationsWithContext(ctx context.Context, ac <-chan Action) <-chan string {
	return helper.MapWithContext(ctx, NormalizeActionsWithContext(ctx, ac), func(a Action) string {
		return a.Annotation()
	})
}

// ActionsToAnnotations takes a channel of action recommendations and returns a
// new channel containing corresponding annotations for those actions.
//
// Deprecated: Use ActionsToAnnotationsWithContext instead.
func ActionsToAnnotations(ac <-chan Action) <-chan string {
	return ActionsToAnnotationsWithContext(context.Background(), ac)
}

// NormalizeActionsWithContext transforms the given channel of actions to ensure a consistent and
// predictable sequence, supporting context cancellation.
func NormalizeActionsWithContext(ctx context.Context, ac <-chan Action) <-chan Action {
	last := Sell

	return helper.MapWithContext(ctx, ac, func(a Action) Action {
		if a != Hold && a != last {
			last = a
			return a
		}

		return Hold
	})
}

// NormalizeActions transforms the given channel of actions to ensure a consistent and
// predictable sequence.
//
// Deprecated: Use NormalizeActionsWithContext instead.
func NormalizeActions(ac <-chan Action) <-chan Action {
	return NormalizeActionsWithContext(context.Background(), ac)
}

// DenormalizeActionsWithContext simplifies the representation of the action sequence and facilitates subsequent
// processing by transforming the given channel of actions, supporting context cancellation.
func DenormalizeActionsWithContext(ctx context.Context, ac <-chan Action) <-chan Action {
	last := Hold

	return helper.MapWithContext(ctx, ac, func(a Action) Action {
		if a != Hold && a != last {
			last = a
		}

		return last
	})
}

// DenormalizeActions simplifies the representation of the action sequence.
//
// Deprecated: Use DenormalizeActionsWithContext instead.
func DenormalizeActions(ac <-chan Action) <-chan Action {
	return DenormalizeActionsWithContext(context.Background(), ac)
}

// CountActions taken a slice of Action channels, and counts them by their type.
func CountActions(acs []<-chan Action) (int, int, int, bool) {
	var buy, hold, sell int

	for _, ac := range acs {
		action, ok := <-ac
		if !ok {
			return 0, 0, 0, false
		}

		switch action {
		case Sell:
			sell++

		case Buy:
			buy++

		default:
			hold++
		}
	}

	return buy, hold, sell, true
}

// CountTransactions counts the number of recommended Buy and Sell actions.
func CountTransactions(ac <-chan Action) <-chan int {
	var transactions int

	return helper.Map(ac, func(action Action) int {
		if action != Hold {
			transactions++
		}

		return transactions
	})
}
