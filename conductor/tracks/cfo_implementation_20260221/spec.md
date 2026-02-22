# Specification: Chande Forecast Oscillator (CFO)

## Overview
The Chande Forecast Oscillator (CFO) is a momentum indicator that measures the difference between a security's price and its linear regression forecast. It helps identify overbought or oversold conditions and potential trend reversals.

## Indicator Logic
- **Calculation:** `CFO = ((Price - Forecast) / Price) * 100`
- **Forecast:** Linear regression forecast for the current price based on a specified period.
- **Period:** Number of data points used for the linear regression calculation.

## Strategy Logic
- **Buy Signal:** CFO crosses above a specified threshold (e.g., 0).
- **Sell Signal:** CFO crosses below a specified threshold (e.g., 0).
