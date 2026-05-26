// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/strategy/trend"
)

func TestAllStrategies(t *testing.T) {
	strategies := trend.AllStrategies()
	if len(strategies) != 19 {
		t.Fatalf("expected 19 strategies, got %d", len(strategies))
	}
}
