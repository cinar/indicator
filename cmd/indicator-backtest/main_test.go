// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package main

import (
	"testing"
)

func TestParseStrategies(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		expectErr     bool
	}{
		{
			name:          "single strategy",
			input:         "macd",
			expectedCount: 1,
			expectErr:     false,
		},
		{
			name:          "multiple strategies with spaces and hyphens",
			input:         "macd, rsi, buy-and-hold",
			expectedCount: 3,
			expectErr:     false,
		},
		{
			name:          "strategy suffix",
			input:         "macd_strategy, rsi_strategy",
			expectedCount: 2,
			expectErr:     false,
		},
		{
			name:          "uppercase names",
			input:         "MACD, RSI, BOLLINGER_BANDS",
			expectedCount: 3,
			expectErr:     false,
		},
		{
			name:          "category trend",
			input:         "trend",
			expectedCount: 19,
			expectErr:     false,
		},
		{
			name:          "category momentum",
			input:         "momentum",
			expectedCount: 9,
			expectErr:     false,
		},
		{
			name:          "category volatility",
			input:         "volatility",
			expectedCount: 9,
			expectErr:     false,
		},
		{
			name:          "category volume",
			input:         "volume",
			expectedCount: 7,
			expectErr:     false,
		},
		{
			name:          "category compound",
			input:         "compound",
			expectedCount: 2,
			expectErr:     false,
		},
		{
			name:          "all strategies",
			input:         "all",
			expectedCount: 47,
			expectErr:     false,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "only spaces and commas",
			input:     "  ,  ,  ",
			expectErr: true,
		},
		{
			name:      "unknown strategy",
			input:     "nonexistent_strategy",
			expectErr: true,
		},
		{
			name:      "valid mixed with unknown",
			input:     "macd, nonexistent_strategy",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strats, err := parseStrategies(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}

			if len(strats) != tt.expectedCount {
				t.Fatalf("expected %d strategies, got %d", tt.expectedCount, len(strats))
			}
		})
	}
}

func TestResolveAllKnownStrategies(t *testing.T) {
	knownStrategies := []string{
		"all", "compound", "momentum", "trend", "volatility", "volume",
		"buy_and_hold", "alligator", "apo", "aroon", "bop", "cci", "cfo", "dema",
		"envelope", "golden_cross", "hma", "kama", "kdj", "macd", "qstick", "smma",
		"trima", "triple_moving_average_crossover", "triple_ma_crossover", "trix",
		"tsi", "vwma", "weighted_close", "awesome_oscillator", "coppock_curve",
		"elder_ray", "ichimoku_cloud", "rsi", "stochastic_oscillator", "stochastic_rsi",
		"triple_rsi", "williams_r", "bollinger_bands", "donchian_channel_breakout",
		"keltner_channel", "super_trend", "chaikin_money_flow", "ease_of_movement",
		"force_index", "money_flow_index", "negative_volume_index", "obv",
		"percent_b_and_mfi", "percent_band_and_mfi", "percent_b_mfi",
		"weighted_average_price", "macd_rsi",
	}

	for _, name := range knownStrategies {
		t.Run(name, func(t *testing.T) {
			strats, err := resolveStrategy(name)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", name, err)
			}
			if len(strats) == 0 {
				t.Fatalf("expected non-empty strategies for %q", name)
			}
		})
	}
}
