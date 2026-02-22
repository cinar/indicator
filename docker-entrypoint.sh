#!/bin/sh
# Copyright (c) 2021-2026 Onur Cinar.
# The source code is provided under GNU AGPLv3 License.
# https://github.com/cinar/indicator

set -e

API_KEY=""
DAYS=365
LAST=365
ASSETS=""
OUTPUT="/app/output"
DATA_DIR="/app/data"

show_usage() {
    echo "Usage: indicator [OPTIONS]"
    echo ""
    echo "Sync market data from Tiingo and run backtests."
    echo ""
    echo "Options:"
    echo "  --api-key KEY   Tiingo API key (required)"
    echo "  --days N        Days of historical data to fetch (default: 365)"
    echo "  --last N        Days to backtest (default: 365)"
    echo "  --assets TICKERS Space-separated list of ticker symbols (default: all available)"
    echo "  --output DIR    Output directory for reports (default: /app/output)"
    echo "  --help          Show this help message"
    echo ""
    echo "Example:"
    echo "  indicator --api-key YOUR_KEY --days 365 --assets aapl msft googl"
    echo ""
    echo "Get your free Tiingo API key at: https://www.tiingo.com/"
}

while [ $# -gt 0 ]; do
    case "$1" in
        --api-key)
            API_KEY="$2"
            shift 2
            ;;
        --days)
            DAYS="$2"
            shift 2
            ;;
        --last)
            LAST="$2"
            shift 2
            ;;
        --assets)
            ASSETS="$2"
            shift 2
            ;;
        --output)
            OUTPUT="$2"
            shift 2
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

if [ -z "$API_KEY" ]; then
    echo "Error: --api-key is required"
    echo ""
    show_usage
    exit 1
fi

echo "=========================================="
echo "Indicator Docker - Sync & Backtest"
echo "=========================================="
echo ""
echo "Configuration:"
echo "  API Key: ***${API_KEY:-none}"
echo "  Days: $DAYS"
echo "  Backtest Period: $LAST days"
echo "  Assets: ${ASSETS:-all}"
echo "  Output: $OUTPUT"
echo ""

mkdir -p "$OUTPUT"

echo "Step 1: Syncing market data from Tiingo..."
echo "-------------------------------------------"

if [ -z "$ASSETS" ]; then
    ./indicator-sync \
        -source-name tiingo \
        -source-config "$API_KEY" \
        -target-name filesystem \
        -target-config "$DATA_DIR" \
        -days "$DAYS"
else
    ./indicator-sync \
        -source-name tiingo \
        -source-config "$API_KEY" \
        -target-name filesystem \
        -target-config "$DATA_DIR" \
        -days "$DAYS" \
        $ASSETS
fi

echo ""
echo "Step 2: Running backtests..."
echo "-------------------------------------------"

if [ -z "$ASSETS" ]; then
    ./indicator-backtest \
        -repository-name filesystem \
        -repository-config "$DATA_DIR" \
        -report-name html \
        -report-config "$OUTPUT" \
        -last "$LAST"
else
    ./indicator-backtest \
        -repository-name filesystem \
        -repository-config "$DATA_DIR" \
        -report-name html \
        -report-config "$OUTPUT" \
        -last "$LAST" \
        $ASSETS
fi

echo ""
echo "=========================================="
echo "Done! Reports generated in: $OUTPUT"
echo "=========================================="
echo ""
echo "Open $OUTPUT/index.html in your browser to view the results."
