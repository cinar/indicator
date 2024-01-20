// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator/v2

package strategy

import (
	"errors"
	"fmt"
)

// Result is only used inside the test cases to facilitate the
// comparison between the actual and expected strategy results.
type Result struct {
	Action  Action
	Outcome float64
}

// CheckResults checks the actual strategy results against the expected results.
func CheckResults(results <-chan *Result, actions <-chan Action, outcomes <-chan float64) error {
	for result := range results {
		action, ok := <-actions
		if !ok {
			return errors.New("actual actions ended early")
		}

		if action != result.Action {
			return fmt.Errorf("actual %v expected %v", action, result.Action)
		}

		outcome, ok := <-outcomes
		if !ok {
			return errors.New("actual outcomes ended early")
		}

		if outcome != result.Outcome {
			return fmt.Errorf("actual %v expected %v", outcome, result.Outcome)
		}
	}

	return nil
}
