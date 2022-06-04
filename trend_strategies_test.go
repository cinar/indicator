// Copyright (c) 2021 Onur Cinar. All Rights Reserved.
// The source code is provided under MIT License.
//
// https://github.com/cinar/indicator

package indicator

import (
	"testing"
	"time"
)

func TestTrendStrategy(t *testing.T) {
	asset := &Asset{
		Date: []time.Time{
			time.Now(), time.Now(), time.Now(), time.Now(), time.Now(),
		},
		Opening: []float64{
			0, 0, 0, 0, 0,
		},
		Closing: []float64{
			0, 1, 2, 1, 0,
		},
		High: []float64{
			0, 0, 0, 0, 0,
		},
		Low: []float64{
			0, 0, 0, 0, 0,
		},
		Volume: []int64{
			0, 0, 0, 0, 0,
		},
	}

	expected := []Action{
		HOLD, HOLD, BUY, HOLD, SELL,
	}

	actions := TrendStrategy(asset, 2)

	for i := 0; i < len(expected); i++ {
		if actions[i] != expected[i] {
			t.Fatalf("at %d actual %d expected %d", i, actions[i], expected[i])
		}
	}
}
