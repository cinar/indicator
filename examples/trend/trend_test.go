// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"testing"

	"github.com/cinar/indicator/v2/examples/trend"
)

func TestAllStrategies(t *testing.T) {
	strategies := trend.AllStrategies()
	if len(strategies) != 21 {
		t.Fatalf("expected 21 strategies, got %d", len(strategies))
	}
}
