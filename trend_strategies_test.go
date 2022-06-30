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

	actual := TrendStrategy(asset, 2)
	testEqualsAction(t, actual, expected)
}

func TestVwmaStrategy(t *testing.T) {
	asset := &Asset{
		Date: []time.Time{
			time.Now(), time.Now(), time.Now(), time.Now(), time.Now(),
		},
		Opening: []float64{
			0, 0, 0, 0, 0,
		},
		Closing: []float64{
			20, 21, 21, 19, 16,
		},
		High: []float64{
			0, 0, 0, 0, 0,
		},
		Low: []float64{
			0, 0, 0, 0, 0,
		},
		Volume: []int64{
			100, 50, 40, 50, 100,
		},
	}

	expected := []Action{
		HOLD, SELL, SELL, SELL, SELL,
	}

	period := 3

	strategy := MakeVwmaStrategy(period)
	actual := strategy(asset)
	testEqualsAction(t, actual, expected)
}

func TestDefaultVwmaStrategy(t *testing.T) {
	asset := &Asset{
		Date: []time.Time{
			time.Now(), time.Now(), time.Now(), time.Now(), time.Now(),
		},
		Opening: []float64{
			0, 0, 0, 0, 0,
		},
		Closing: []float64{
			20, 21, 21, 19, 16,
		},
		High: []float64{
			0, 0, 0, 0, 0,
		},
		Low: []float64{
			0, 0, 0, 0, 0,
		},
		Volume: []int64{
			100, 50, 40, 50, 100,
		},
	}

	expected := []Action{
		HOLD, SELL, SELL, SELL, SELL,
	}

	actual := DefaultVwmaStrategy(asset)
	testEqualsAction(t, actual, expected)
}
