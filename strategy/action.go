// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy

import "github.com/cinar/indicator/v2/helper"

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

// ActionsToAnnotations takes a channel of action recommendations and returns a
// new channel containing corresponding annotations for those actions.
func ActionsToAnnotations(ac <-chan Action) <-chan string {
	return helper.Map(NormalizeActions(ac), func(a Action) string {
		return a.Annotation()
	})
}

// NormalizeActions transforms the given channel of actions to ensure a consistent and
// predictable sequence. It eliminates consecutive occurrences of the same action
// (Buy/Sell), ensuring the order follows a pattern of Hold, Buy, Hold, Sell.
func NormalizeActions(ac <-chan Action) <-chan Action {
	last := Sell

	return helper.Map(ac, func(a Action) Action {
		if a != Hold && a != last {
			last = a
			return a
		}

		return Hold
	})
}

// DenormalizeActions simplifies the representation of the action sequence and facilitates subsequent
// processing by transforming the given channel of actions. It retains Hold actions until the
// first Buy or Sell action appears. Subsequently, it replaces all remaining Hold actions with
// the preceding Buy or Sell action, effectively merging consecutive actions.
func DenormalizeActions(ac <-chan Action) <-chan Action {
	last := Hold

	return helper.Map(ac, func(a Action) Action {
		if a != Hold && a != last {
			last = a
		}

		return last
	})
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
