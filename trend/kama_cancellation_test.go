//go:build go1.25

// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package trend_test

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

func TestKamaCancellationWithPendingSmoothing(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		closings := make([]float64, 40)
		for i := range closings {
			closings[i] = float64(i%7 + 1)
		}
		kama := trend.NewKamaWith[float64](2, 2, 4)
		actual := kama.ComputeWithContext(ctx, helper.SliceToChanWithContext(ctx, closings))
		synctest.Wait()
		cancel()
		for range actual {
		}
		synctest.Wait()
	})
}
