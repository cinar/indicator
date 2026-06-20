// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"context"
	"sort"
)

// PercentRank wraps PercentRankWithContext for backwards compatibility.
//
// Deprecated: Use PercentRankWithContext instead.
func PercentRank[T Number](c <-chan T, period int) <-chan T {
	return PercentRankWithContext(context.Background(), c, period)
}

// PercentRankWithContext returns a channel that emits the percentile rank
// of each value compared to the previous period-1 values, supporting context cancellation.
func PercentRankWithContext[T Number](ctx context.Context, c <-chan T, period int) <-chan T {
	r := make(chan T)

	if period <= 1 {
		close(r)
		return r
	}

	go func() {
		defer close(r)

		values := make([]T, 0, period)
		count := 0

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-c:
				if !ok {
					return
				}
				if count < period {
					values = append(values, value)
					count++
					continue
				}

				// Shift: remove oldest, add new
				copy(values[0:period-1], values[1:period])
				values[period-1] = value

				// Count how many values are less than current
				lessCount := 0
				for i := 0; i < period-1; i++ {
					if values[i] < value {
						lessCount++
					}
				}

				rank := float64(lessCount) * 100.0 / float64(period-1)
				select {
				case <-ctx.Done():
					return
				case r <- T(rank):
				}
			}
		}
	}()

	return r
}

// SortedPercentRank wraps SortedPercentRankWithContext for backwards compatibility.
//
// Deprecated: Use SortedPercentRankWithContext instead.
func SortedPercentRank[T Number](c <-chan T, period int) <-chan T {
	return SortedPercentRankWithContext(context.Background(), c, period)
}

// SortedPercentRankWithContext returns a channel that emits the percentile rank
// by sorting the window values, supporting context cancellation.
func SortedPercentRankWithContext[T Number](ctx context.Context, c <-chan T, period int) <-chan T {
	r := make(chan T)

	if period <= 1 {
		close(r)
		return r
	}

	go func() {
		defer close(r)

		values := make([]T, 0, period)
		count := 0

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-c:
				if !ok {
					return
				}
				if count < period {
					values = append(values, value)
					count++
					continue
				}

				// Shift: remove oldest, add new
				copy(values[0:period-1], values[1:period])
				values[period-1] = value

				// Sort copy for ranking
				sorted := make([]T, period-1)
				copy(sorted, values[:period-1])
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i] < sorted[j]
				})

				// Binary search for rank
				rankIdx := sort.Search(period-1, func(i int) bool {
					return sorted[i] >= value
				})

				rank := float64(rankIdx) * 100.0 / float64(period-1)
				select {
				case <-ctx.Done():
					return
				case r <- T(rank):
				}
			}
		}
	}()

	return r
}
