// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package talibcompat_test

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/trend"
	"github.com/cinar/indicator/v2/volatility"
	"github.com/cinar/indicator/v2/volume"
)

// datasets are the generated reference files in testdata/.
var datasets = []string{"walk", "ticks"}

// tolerance is the allowed relative difference, with an absolute floor of
// the same size for values near zero. The reference values are written with
// 10 significant digits.
const tolerance = 1e-8

// bars holds one dataset: the OHLCV inputs and the TA-Lib reference columns.
type bars struct {
	open, high, low, close, volume []float64
	columns                        map[string][]float64
}

// in returns a fresh input channel for a series.
func in(s []float64) <-chan float64 {
	return helper.SliceToChan(s)
}

// outs collects an indicator's output channels.
func outs(c ...<-chan float64) []<-chan float64 {
	return c
}

// testCase maps a Go indicator to the TA-Lib columns it is compared against.
type testCase struct {
	// name identifies the case in test output.
	name string

	// columns are the reference columns, one per Go output, in output order.
	columns []string

	// compute runs the Go indicator on the dataset.
	compute func(b *bars) []<-chan float64

	// scale multiplies the Go output before comparison. It records a
	// documented unit difference, e.g. 100 for a ratio that TA-Lib reports
	// as a percentage. Zero means 1.
	scale float64

	// fromBar skips the comparison before this bar index. It records a
	// documented EMA seeding difference that converges after a warm-up.
	fromBar int

	// skip is the link to the issue tracking a known difference.
	skip string
}

var testCases = []testCase{
	// trend
	{name: "Apo", columns: []string{"APO"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewApo[float64]().Compute(in(b.close)))
	}},
	{name: "Aroon", columns: []string{"AROON_UP", "AROON_DOWN"}, skip: "https://github.com/cinar/indicator/issues/520", compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewAroon[float64]().Compute(in(b.high), in(b.low)))
	}},
	{name: "Bop", columns: []string{"BOP"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewBop[float64]().Compute(in(b.open), in(b.high), in(b.low), in(b.close)))
	}},
	{name: "Cci", columns: []string{"CCI"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewCci[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "Cfo", columns: []string{"CFO"}, skip: "https://github.com/cinar/indicator/issues/530", compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewCfo[float64]().Compute(in(b.close)))
	}},
	{name: "Dema", columns: []string{"DEMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewDema[float64]().Compute(in(b.close)))
	}},
	{name: "Dpo", columns: []string{"DPO"}, skip: "https://github.com/cinar/indicator/issues/522", compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewDpo[float64]().Compute(in(b.close)))
	}},
	{name: "Ema", columns: []string{"EMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewEma[float64]().Compute(in(b.close)))
	}},
	{name: "Hma", columns: []string{"HMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewHma[float64]().Compute(in(b.close)))
	}},
	{name: "Kama", columns: []string{"KAMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewKama[float64]().Compute(in(b.close)))
	}},
	{name: "Kdj", columns: []string{"KDJ_K", "KDJ_D", "KDJ_J"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewKdj[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	// Kst is built on Roc, which returns a ratio (#533).
	{name: "Kst", columns: []string{"KST", "KST_SIGNAL"}, scale: 100, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewKst[float64]().Compute(in(b.close)))
	}},
	// Macd seeds the fast EMA on its own first values (#535).
	{name: "Macd", columns: []string{"MACD", "MACD_SIGNAL"}, fromBar: 170, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewMacd[float64]().Compute(in(b.close)))
	}},
	{name: "MassIndex", columns: []string{"MASSI"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewMassIndex[float64]().Compute(in(b.high), in(b.low)))
	}},
	{name: "McGinleyDynamic", columns: []string{"MCGD"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewMcGinleyDynamic[float64]().Compute(in(b.close)))
	}},
	{name: "MovingMax", columns: []string{"MAX"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewMovingMaxWithPeriod[float64](14).Compute(in(b.close)))
	}},
	{name: "MovingMin", columns: []string{"MIN"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewMovingMinWithPeriod[float64](14).Compute(in(b.close)))
	}},
	{name: "MovingSum", columns: []string{"SUM"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewMovingSumWithPeriod[float64](14).Compute(in(b.close)))
	}},
	{name: "Rma", columns: []string{"RMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewRma[float64]().Compute(in(b.close)))
	}},
	// Roc returns a ratio (#533).
	{name: "Roc", columns: []string{"ROC"}, scale: 100, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewRoc[float64]().Compute(in(b.close)))
	}},
	{name: "Slope", columns: []string{"SLOPE"}, skip: "https://github.com/cinar/indicator/issues/531", compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewSlope[float64]().Compute(in(b.close)))
	}},
	{name: "SlowStochastic", columns: []string{"SLOWSTOCH_K", "SLOWSTOCH_D"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewSlowStochastic[float64]().Compute(in(b.close)))
	}},
	{name: "Sma", columns: []string{"SMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewSma[float64]().Compute(in(b.close)))
	}},
	{name: "Smma", columns: []string{"SMMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewSmma[float64]().Compute(in(b.close)))
	}},
	{name: "Stc", columns: []string{"STC"}, skip: "https://github.com/cinar/indicator/issues/525", compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewStc[float64]().Compute(in(b.close)))
	}},
	{name: "Stochastic", columns: []string{"STOCH1_K", "STOCH1_D"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewStochastic[float64]().Compute(in(b.close)))
	}},
	{name: "T3", columns: []string{"T3"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewT3[float64]().Compute(in(b.close)))
	}},
	{name: "Tema", columns: []string{"TEMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewTema[float64]().Compute(in(b.close)))
	}},
	{name: "Trima", columns: []string{"TRIMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewTrima[float64]().Compute(in(b.close)))
	}},
	// Trix returns a ratio (#533).
	{name: "Trix", columns: []string{"TRIX"}, scale: 100, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewTrix[float64]().Compute(in(b.close)))
	}},
	{name: "Tsi", columns: []string{"TSI"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewTsi[float64]().Compute(in(b.close)))
	}},
	{name: "TypicalPrice", columns: []string{"TYPPRICE"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewTypicalPrice[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "Vwma", columns: []string{"VWMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewVwma[float64]().Compute(in(b.close), in(b.volume)))
	}},
	{name: "WeightedClose", columns: []string{"WCLPRICE"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewWeightedClose[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "Wma", columns: []string{"WMA"}, compute: func(b *bars) []<-chan float64 {
		return outs(trend.NewWma[float64]().Compute(in(b.close)))
	}},

	// momentum
	{name: "AwesomeOscillator", columns: []string{"AO"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewAwesomeOscillator[float64]().Compute(in(b.high), in(b.low)))
	}},
	// ChaikinOscillator seeds each EMA on its own first values (#535).
	{name: "ChaikinOscillator", columns: []string{"ADOSC", "AD"}, fromBar: 150, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewChaikinOscillator[float64]().Compute(in(b.high), in(b.low), in(b.close), in(b.volume)))
	}},
	{name: "ConnorsRsi", columns: []string{"CRSI"}, skip: "https://github.com/cinar/indicator/issues/524", compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewConnorsRsi[float64]().Compute(in(b.close)))
	}},
	{name: "CoppockCurve", columns: []string{"COPPOCK"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewCoppockCurve[float64]().Compute(in(b.close)))
	}},
	{name: "EhlersFisher", columns: []string{"EFISHER"}, skip: "https://github.com/cinar/indicator/issues/523", compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewEhlersFisher[float64]().Compute(in(b.high), in(b.low)))
	}},
	{name: "ElderRay", columns: []string{"ELDER_BULL", "ELDER_BEAR"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewElderRay[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "IchimokuCloud", columns: []string{"ICH_CONV", "ICH_BASE", "ICH_SPANA", "ICH_SPANB"}, skip: "https://github.com/cinar/indicator/issues/534", compute: func(b *bars) []<-chan float64 {
		conversion, base, spanA, spanB, lagging := momentum.NewIchimokuCloud[float64]().Compute(in(b.high), in(b.low), in(b.close))
		go helper.Drain(lagging)
		return outs(conversion, base, spanA, spanB)
	}},
	{name: "InternalBarStrength", columns: []string{"IBS"}, skip: "https://github.com/cinar/indicator/issues/532", compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewInternalBarStrength[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "Ppo", columns: []string{"PPO", "PPO_SIGNAL", "PPO_HIST"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewPpo[float64]().Compute(in(b.close)))
	}},
	{name: "Pvo", columns: []string{"PVO", "PVO_SIGNAL", "PVO_HIST"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewPvo[float64]().Compute(in(b.volume)))
	}},
	{name: "Qstick", columns: []string{"QSTICK"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewQstick[float64]().Compute(in(b.open), in(b.close)))
	}},
	{name: "Rsi", columns: []string{"RSI"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewRsi[float64]().Compute(in(b.close)))
	}},
	{name: "StochasticOscillator", columns: []string{"STOCHF_K", "STOCHF_D"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewStochasticOscillator[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	// StochasticRsi returns a 0-1 ratio (#533).
	{name: "StochasticRsi", columns: []string{"STOCHRSI"}, scale: 100, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewStochasticRsi[float64]().Compute(in(b.close)))
	}},
	{name: "UltimateOscillator", columns: []string{"ULTOSC"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewUltimateOscillator[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "WilliamsR", columns: []string{"WILLR"}, compute: func(b *bars) []<-chan float64 {
		return outs(momentum.NewWilliamsR[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},

	// volatility
	{name: "AccelerationBands", columns: []string{"ACC_UP", "ACC_MID", "ACC_LOW"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewAccelerationBands[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "Atr", columns: []string{"ATR"}, skip: "https://github.com/cinar/indicator/issues/526", compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewAtr[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	// Wilder's ATR through the existing extension point (#526).
	{name: "AtrWithRma", columns: []string{"ATR"}, compute: func(b *bars) []<-chan float64 {
		atr := volatility.NewAtrWithMa(trend.NewRmaWithPeriod[float64](volatility.DefaultAtrPeriod))
		return outs(atr.Compute(in(b.high), in(b.low), in(b.close)))
	}},
	// BollingerBandWidth returns a ratio (#533).
	{name: "BollingerBandWidth", columns: []string{"BBW"}, scale: 100, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewBollingerBandWidth[float64]().Compute(in(b.close)))
	}},
	{name: "BollingerBands", columns: []string{"BB_UP", "BB_MID", "BB_LOW"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewBollingerBands[float64]().Compute(in(b.close)))
	}},
	{name: "ChandelierExit", columns: []string{"CE_LONG", "CE_SHORT"}, skip: "https://github.com/cinar/indicator/issues/526", compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewChandelierExit[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "Chop", columns: []string{"CHOP"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewChop[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "DonchianChannel", columns: []string{"DON_UP", "DON_MID", "DON_LOW"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewDonchianChannel[float64]().Compute(in(b.high), in(b.low)))
	}},
	{name: "KeltnerChannel", columns: []string{"KC_UP", "KC_MID", "KC_LOW"}, skip: "https://github.com/cinar/indicator/issues/528", compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewKeltnerChannel[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "MovingStd", columns: []string{"STDDEV"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewMovingStdWithPeriod[float64](20).Compute(in(b.close)))
	}},
	{name: "PercentB", columns: []string{"PERCENTB"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewPercentB[float64]().Compute(in(b.close)))
	}},
	{name: "SuperTrend", columns: []string{"SUPERTREND"}, skip: "https://github.com/cinar/indicator/issues/527", compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewSuperTrend[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	// SuperTrend on Wilder's ATR through the existing extension point (#527).
	{name: "SuperTrendWithRma", columns: []string{"SUPERTREND"}, compute: func(b *bars) []<-chan float64 {
		superTrend := volatility.NewSuperTrendWithMa(
			trend.NewRmaWithPeriod[float64](volatility.DefaultSuperTrendPeriod),
			volatility.DefaultSuperTrendMultiplier,
		)
		return outs(superTrend.Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "TrueRange", columns: []string{"TRANGE"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewTrueRange[float64]().Compute(in(b.high), in(b.low), in(b.close)))
	}},
	{name: "ZScore", columns: []string{"ZSCORE"}, compute: func(b *bars) []<-chan float64 {
		return outs(volatility.NewZScore[float64]().Compute(in(b.close)))
	}},

	// volume
	{name: "Ad", columns: []string{"AD"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewAd[float64]().Compute(in(b.high), in(b.low), in(b.close), in(b.volume)))
	}},
	{name: "Cmf", columns: []string{"CMF"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewCmf[float64]().Compute(in(b.high), in(b.low), in(b.close), in(b.volume)))
	}},
	{name: "Emv", columns: []string{"EMV"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewEmv[float64]().Compute(in(b.high), in(b.low), in(b.volume)))
	}},
	{name: "Fi", columns: []string{"EFI"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewFi[float64]().Compute(in(b.close), in(b.volume)))
	}},
	{name: "Kvo", columns: []string{"KVO", "KVO_SIGNAL"}, skip: "https://github.com/cinar/indicator/issues/529", compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewKvo[float64]().Compute(in(b.high), in(b.low), in(b.volume)))
	}},
	{name: "Mfi", columns: []string{"MFI"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewMfi[float64]().Compute(in(b.high), in(b.low), in(b.close), in(b.volume)))
	}},
	{name: "Nvi", columns: []string{"NVI"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewNvi[float64]().Compute(in(b.close), in(b.volume)))
	}},
	{name: "Obv", columns: []string{"OBV"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewObv[float64]().Compute(in(b.close), in(b.volume)))
	}},
	{name: "Vpt", columns: []string{"PVT"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewVpt[float64]().Compute(in(b.close), in(b.volume)))
	}},
	// Vwap is a rolling window, which TA-Lib calls VWMA (#535).
	{name: "Vwap", columns: []string{"VWAP"}, compute: func(b *bars) []<-chan float64 {
		return outs(volume.NewVwap[float64]().Compute(in(b.close), in(b.volume)))
	}},
}

func TestTaLibCompat(t *testing.T) {
	for _, dataset := range datasets {
		b := readBars(t, "testdata/"+dataset+".csv")

		for _, tc := range testCases {
			t.Run(dataset+"/"+tc.name, func(t *testing.T) {
				if tc.skip != "" {
					t.Skip("known difference from TA-Lib: " + tc.skip)
				}

				results := collect(tc.compute(b))
				if len(results) != len(tc.columns) {
					t.Fatalf("got %d outputs, expected %d", len(results), len(tc.columns))
				}

				for i, column := range tc.columns {
					expected, ok := b.columns[column]
					if !ok {
						t.Fatalf("column %s not found", column)
					}

					compare(t, column, results[i], expected, tc.scale, tc.fromBar)
				}
			})
		}
	}
}

// collect drains all output channels concurrently, since the outputs of an
// indicator are usually produced in lockstep.
func collect(channels []<-chan float64) [][]float64 {
	results := make([][]float64, len(channels))

	var wg sync.WaitGroup
	for i, c := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = helper.ChanToSlice(c)
		}()
	}
	wg.Wait()

	return results
}

// compare checks a Go output against a reference column. The Go output is
// aligned to the end of the input, since every compared indicator emits one
// value per bar once its idle period has passed.
func compare(t *testing.T, column string, actual, expected []float64, scale float64, fromBar int) {
	t.Helper()

	if scale == 0 {
		scale = 1
	}

	offset := len(expected) - len(actual)
	if offset < 0 {
		t.Fatalf("%s: got %d values for %d bars", column, len(actual), len(expected))
	}

	compared := 0
	for i, value := range actual {
		bar := offset + i
		if bar < fromBar || math.IsNaN(expected[bar]) {
			continue
		}

		value *= scale
		if math.Abs(value-expected[bar]) > tolerance*math.Max(1, math.Abs(expected[bar])) {
			t.Fatalf("%s: bar %d actual %v expected %v", column, bar, value, expected[bar])
		}

		compared++
	}

	if compared < len(expected)/3 {
		t.Fatalf("%s: only %d of %d bars compared", column, compared, len(expected))
	}
}

// readBars reads a generated dataset. Empty cells are bars where TA-Lib
// produces no value, and are read as NaN.
func readBars(t *testing.T, fileName string) *bars {
	t.Helper()

	f, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	header, rows := records[0], records[1:]
	columns := make(map[string][]float64, len(header))

	for j, name := range header {
		column := make([]float64, len(rows))

		for i, row := range rows {
			if row[j] == "" {
				column[i] = math.NaN()
				continue
			}

			column[i], err = strconv.ParseFloat(row[j], 64)
			if err != nil {
				t.Fatalf("%s row %d column %s: %v", fileName, i+1, name, err)
			}
		}

		columns[name] = column
	}

	return &bars{
		open:    columns["Open"],
		high:    columns["High"],
		low:     columns["Low"],
		close:   columns["Close"],
		volume:  columns["Volume"],
		columns: columns,
	}
}
