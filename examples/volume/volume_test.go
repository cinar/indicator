// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package volume_test

import (
	"testing"

	"github.com/cinar/indicator/v2/examples/volume"
)

func TestAllStrategies(t *testing.T) {
	strategies := volume.AllStrategies()
	if len(strategies) != 7 {
		t.Fatalf("expected 7 strategies, got %d", len(strategies))
	}
}
