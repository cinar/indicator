// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/strategy/volume"
)

func TestAllStrategies(t *testing.T) {
	strategies := volume.AllStrategies()
	if len(strategies) != 6 {
		t.Fatalf("expected 6 strategies, got %d", len(strategies))
	}
}
