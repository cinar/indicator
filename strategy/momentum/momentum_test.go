// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package momentum_test

import (
	"testing"

	"github.com/cinar/indicator/v2/strategy/momentum"
)

func TestAllStrategies(t *testing.T) {
	strategies := momentum.AllStrategies()
	if len(strategies) != 8 {
		t.Fatalf("expected 8 strategies, got %d", len(strategies))
	}
}
