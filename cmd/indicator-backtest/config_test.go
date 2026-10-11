// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/cinar/indicator/v2/examples/momentum"
)

const testSnapshotsCsv = `Date,Open,High,Low,Close,Adj Close,Volume
2023-11-01,340,345,339,343.75,343.75,2789700
2023-11-02,344,349,343,349.02,349.02,3433700
2023-11-03,349,352,347,350.10,350.10,3011200
2023-11-06,350,353,348,351.20,351.20,2987600
2023-11-07,351,356,349,354.30,354.30,3120400
2023-11-08,354,357,352,355.10,355.10,2890300
2023-11-09,355,358,353,356.40,356.40,3054800
`

func writeTestConfig(t *testing.T, contents string) string {
	t.Helper()

	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "TEST.csv"), []byte(testSnapshotsCsv), 0o600); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(dir, "config.json")

	if err := os.WriteFile(configPath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	return configPath
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected an error for a missing config file")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	if err := os.WriteFile(configPath, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("expected an error for invalid JSON")
	}
}

func TestLoadConfigNoStrategies(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	if err := os.WriteFile(configPath, []byte(`{"repository":{"name":"filesystem","config":"."}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("expected an error for a config file without strategies")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	configPath := writeTestConfig(t, `{"strategies":[{"name":"rsi"}]}`)

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if config.Repository.Name != "filesystem" {
		t.Errorf("expected default repository name filesystem, got %s", config.Repository.Name)
	}

	if config.Report.Name != "html" {
		t.Errorf("expected default report name html, got %s", config.Report.Name)
	}

	if config.Workers == 0 || config.LastDays == 0 {
		t.Errorf("expected non-zero default workers and lastDays, got %+v", config)
	}
}

func TestNewStrategyFromConfigUnknownStrategy(t *testing.T) {
	_, err := NewStrategyFromConfig("not-a-real-strategy", nil)
	if err == nil {
		t.Fatal("expected an error for an unknown strategy name")
	}
}

func TestNewStrategyFromConfigOverlay(t *testing.T) {
	s, err := NewStrategyFromConfig("rsi", []byte(`{"BuyAt": 20, "SellAt": 80, "Rsi": {"Rma": {"Period": 5}}}`))
	if err != nil {
		t.Fatal(err)
	}

	rsi, ok := s.(*momentum.RsiStrategy)
	if !ok {
		t.Fatalf("expected *momentum.RsiStrategy, got %T", s)
	}

	if rsi.BuyAt != 20 || rsi.SellAt != 80 {
		t.Errorf("expected BuyAt=20 SellAt=80, got BuyAt=%v SellAt=%v", rsi.BuyAt, rsi.SellAt)
	}

	if rsi.Rsi.Rma.Period != 5 {
		t.Errorf("expected nested Rma.Period=5, got %v", rsi.Rsi.Rma.Period)
	}
}

func TestNewStrategyFromConfigInvalidJSON(t *testing.T) {
	_, err := NewStrategyFromConfig("rsi", []byte(`{not json`))
	if err == nil {
		t.Fatal("expected an error for invalid strategy config JSON")
	}
}

func TestNewBacktestFromConfigRun(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")

	if err := os.WriteFile(filepath.Join(dir, "TEST.csv"), []byte(testSnapshotsCsv), 0o600); err != nil {
		t.Fatal(err)
	}

	config := &Config{
		Repository: SourceConfig{Name: "filesystem", Config: dir},
		Report:     SourceConfig{Name: "html", Config: outDir},
		Symbols:    []string{"TEST"},
		LastDays:   3650,
		Strategies: []StrategyConfig{
			{Name: "rsi"},
			{Name: "rsi", Config: []byte(`{"BuyAt": 20, "SellAt": 80}`)},
		},
	}

	backtester, err := NewBacktestFromConfig(config, slog.Default())
	if err != nil {
		t.Fatal(err)
	}

	if err := backtester.Run(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
		t.Errorf("expected an index.html report to be written: %v", err)
	}
}

func TestNewBacktestFromConfigUnknownRepository(t *testing.T) {
	config := &Config{
		Repository: SourceConfig{Name: "not-a-real-repository"},
		Report:     SourceConfig{Name: "html", Config: t.TempDir()},
		Strategies: []StrategyConfig{{Name: "rsi"}},
	}

	_, err := NewBacktestFromConfig(config, slog.Default())
	if err == nil {
		t.Fatal("expected an error for an unknown repository name")
	}
}
