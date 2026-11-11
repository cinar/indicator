# Copyright (c) 2021-2026 The Indicator Authors.
# The source code is provided under GNU AGPLv3 License.
# https://github.com/cinar/indicator

"""Generates the TA-Lib reference data used by talibcompat_test.go.

This script is the only part of the comparison that needs TA-Lib. The Go
tests read the CSV files it writes and use only the standard library, so the
Go module stays free of external dependencies.

Usage (from this directory):

    python3 -m venv /tmp/talib-venv
    /tmp/talib-venv/bin/pip install ta-lib==0.8.3 numpy
    /tmp/talib-venv/bin/python generate.py

Each output file holds the OHLCV inputs followed by one column per reference
output. Bars where TA-Lib produces no value are left empty.
"""

import numpy as np
import talib as ta

BARS = 300

# name -> (seed, tick). Prices are rounded to the tick. The "walk" tick is fine
# enough that ties practically never occur, while the "ticks" tick produces
# equal highs, lows and closes, which exercises tie handling.
DATASETS = {
    "walk": (7, 0.0001),
    "ticks": (11, 0.5),
}


def make_bars(seed, tick):
    """Returns a seeded random-walk OHLCV series."""
    rng = np.random.default_rng(seed)
    close = 100 + np.cumsum(rng.normal(0, 1, BARS))
    open_ = close + rng.normal(0, 0.5, BARS)
    high = np.maximum(open_, close) + rng.uniform(0, 1.5, BARS)
    low = np.minimum(open_, close) - rng.uniform(0, 1.5, BARS)
    volume = rng.integers(1000, 100000, BARS).astype(float)

    open_, high, low, close = (np.round(np.round(x / tick) * tick, 4) for x in (open_, high, low, close))
    high = np.maximum(high, np.maximum(open_, close))
    low = np.minimum(low, np.minimum(open_, close))

    return open_, high, low, close, volume


def references(o, h, l, c, v):
    """Returns the reference outputs keyed by column name.

    Parameters match the Go indicators' defaults, so that each Go test case
    can use the default constructor. A few columns with no single TA-Lib
    function are composed from TA-Lib building blocks; each one is noted.
    """
    r = {}

    # trend
    r["APO"] = ta.APO(c, 14, 30, 1)
    r["AROON_DOWN"], r["AROON_UP"] = ta.AROON(h, l, 25)
    r["BOP"] = ta.BOP(o, h, l, c)
    r["CCI"] = ta.CCI(h, l, c, 20)
    r["CFO"] = ta.FOSC(c, 14)
    r["DEMA"] = ta.DEMA(c, 20)
    r["DPO"] = ta.DPO(c, 20)
    r["EMA"] = ta.EMA(c, 20)
    r["HMA"] = ta.HMA(c, 20)
    r["KAMA"] = ta.KAMA(c, 10)
    r["KDJ_K"], r["KDJ_D"], r["KDJ_J"] = ta.KDJ(h, l, c, 9, 3, 0, 3, 0)
    r["KST"], r["KST_SIGNAL"] = ta.KST(c)
    r["MACD"], r["MACD_SIGNAL"], _ = ta.MACD(c, 12, 26, 9)
    r["MASSI"] = ta.MASSI(h, l, 9, 25)
    r["MCGD"] = ta.MCGD(c, 14)
    r["MAX"] = ta.MAX(c, 14)
    r["MIN"] = ta.MIN(c, 14)
    r["SUM"] = ta.SUM(c, 14)
    r["RMA"] = ta.RMA(c, 20)
    r["ROC"] = ta.ROC(c, 9)
    r["SLOPE"] = ta.LINEARREG_SLOPE(c, 14)
    r["SMA"] = ta.SMA(c, 50)
    r["SMMA"] = ta.RMA(c, 7)
    r["STC"] = ta.STC(c, 23, 50, 10)
    r["STOCH1_K"], r["STOCH1_D"] = ta.STOCHF(c, c, c, 10, 3, 0)
    r["SLOWSTOCH_K"], r["SLOWSTOCH_D"] = ta.STOCH(c, c, c, 10, 3, 0, 3, 0)
    r["T3"] = ta.T3(c, 5, 0.7)
    r["TEMA"] = ta.TEMA(c, 20)
    r["TRIMA"] = ta.TRIMA(c, 15)
    r["TRIX"] = ta.TRIX(c, 15)
    r["TSI"] = ta.TSI(c, 25, 13)
    r["TYPPRICE"] = ta.TYPPRICE(h, l, c)
    r["VWMA"] = ta.VWMA(c, v, 20)
    r["WCLPRICE"] = ta.WCLPRICE(h, l, c)
    r["WMA"] = ta.WMA(c, 20)

    # momentum
    r["AO"] = ta.AO(h, l, 5, 34)
    r["ADOSC"] = ta.ADOSC(h, l, c, v, 3, 10)
    r["CRSI"] = ta.CRSI(c, 3, 2, 100)
    r["COPPOCK"] = ta.COPPOCK(c, 10, 11, 14)
    r["EFISHER"] = ta.FISHER(h, l, 10)[0]
    ema13 = ta.EMA(c, 13)  # Elder Ray: composed from TA-Lib EMA.
    r["ELDER_BULL"], r["ELDER_BEAR"] = h - ema13, l - ema13
    r["IBS"] = ta.IBS(h, l, c)
    r["ICH_CONV"], r["ICH_BASE"], r["ICH_SPANA"], r["ICH_SPANB"] = ta.ICHIMOKU(h, l, 9, 26, 52)
    ppo = ta.PPO(c, 12, 26, 1)
    r["PPO"], r["PPO_SIGNAL"] = ppo, ta.EMA(ppo, 9)
    r["PPO_HIST"] = r["PPO"] - r["PPO_SIGNAL"]
    pvo = ta.PVO(v, 12, 26, 1)
    r["PVO"], r["PVO_SIGNAL"] = pvo, ta.EMA(pvo, 9)
    r["PVO_HIST"] = r["PVO"] - r["PVO_SIGNAL"]
    r["QSTICK"] = ta.QSTICK(o, c, 20)
    r["RSI"] = ta.RSI(c, 14)
    r["STOCHF_K"], r["STOCHF_D"] = ta.STOCHF(h, l, c, 14, 3, 0)
    r["STOCHRSI"] = ta.STOCHRSI(c, 14, 14, 1, 0)[0]
    r["ULTOSC"] = ta.ULTOSC(h, l, c, 7, 14, 28)
    r["WILLR"] = ta.WILLR(h, l, c, 14)

    # volatility
    r["ACC_UP"], r["ACC_MID"], r["ACC_LOW"] = ta.ACCBANDS(h, l, c, 20)
    r["ATR"] = ta.ATR(h, l, c, 14)
    r["BBW"] = ta.BBW(c, 20, 2, 2, 0)
    r["BB_UP"], r["BB_MID"], r["BB_LOW"] = ta.BBANDS(c, 20, 2, 2, 0)
    atr22 = ta.ATR(h, l, c, 22)  # Chandelier Exit: composed from TA-Lib MAX/MIN/ATR.
    r["CE_LONG"] = ta.MAX(h, 22) - 3 * atr22
    r["CE_SHORT"] = ta.MIN(l, 22) + 3 * atr22
    r["CHOP"] = ta.CHOP(h, l, c, 14)
    r["DON_UP"], r["DON_MID"], r["DON_LOW"] = ta.DONCHIAN(h, l, 20)
    r["KC_UP"], r["KC_MID"], r["KC_LOW"] = ta.KC(h, l, c, 20, 20, 2.0)
    r["PERCENTB"] = ta.PERCENTB(c, 20, 2, 2, 0)
    r["STDDEV"] = ta.STDDEV(c, 20, 1)
    r["SUPERTREND"] = ta.SUPERTREND(h, l, c, 14, 2.5)[0]
    r["TRANGE"] = ta.TRANGE(h, l, c)
    # Z-Score: composed from TA-Lib SMA/STDDEV.
    r["ZSCORE"] = (c - ta.SMA(c, 20)) / ta.STDDEV(c, 20, 1)

    # volume
    r["AD"] = ta.AD(h, l, c, v)
    r["CMF"] = ta.CMF(h, l, c, v, 20)
    r["EFI"] = ta.EFI(c, v, 13)
    r["EMV"] = ta.EMV(h, l, v, 14, 1e8)  # Go uses a volume divisor of 1e8.
    r["KVO"], r["KVO_SIGNAL"] = ta.KVO(h, l, c, v, 34, 55, 13)
    r["MFI"] = ta.MFI(h, l, c, v, 14)
    r["NVI"] = ta.NVI(c, v)
    r["OBV"] = ta.OBV(c, v)
    r["PVT"] = ta.PVT(c, v)
    # Go Vwap is a rolling window, which is TA-Lib VWMA; TA-Lib VWAP is cumulative.
    r["VWAP"] = ta.VWMA(c, v, 14)

    return r


def fmt(x):
    return "" if np.isnan(x) else f"{x:.10g}"


def main():
    for name, (seed, tick) in DATASETS.items():
        o, h, l, c, v = make_bars(seed, tick)
        refs = references(o, h, l, c, v)
        columns = ["Open", "High", "Low", "Close", "Volume"] + sorted(refs)
        data = [o, h, l, c, v] + [np.asarray(refs[k], dtype=float) for k in sorted(refs)]

        with open(f"{name}.csv", "w", encoding="utf-8") as f:
            f.write(",".join(columns) + "\n")
            for i in range(BARS):
                f.write(",".join(fmt(col[i]) for col in data) + "\n")


if __name__ == "__main__":
    main()
