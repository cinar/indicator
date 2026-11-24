// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// extremeSinceWithContext returns a channel of T indicating how many values
// ago the extreme value within the sliding window of size w was, where better
// reports whether a is strictly more extreme than b. On ties the most recent
// value counts, as in TA-Lib. NaN values are never taken as the extreme; a
// window that holds only NaN values yields NaN. While the window fills, the
// values seen so far are used.
//
// The position of the extreme is kept between values, and the window is
// scanned again only when that position leaves it.
func extremeSinceWithContext[T Number](ctx context.Context, c <-chan T, w int, better func(a, b T) bool) <-chan T {
	r := make(chan T)

	if w <= 0 {
		close(r)
		return r
	}

	go func() {
		defer close(r)

		// h is a ring buffer where the value at index t is h[t%w].
		h := make([]T, w)

		// best is the index of the extreme value, or -1 if there is none.
		best := -1

		for t := 0; ; t++ {
			var n T
			var ok bool

			select {
			case <-ctx.Done():
				return
			case n, ok = <-c:
				if !ok {
					return
				}
			}

			h[t%w] = n
			oldest := max(t-w+1, 0)

			switch {
			case n != n: // NaN
				if best < oldest {
					best = -1
				}
			case best < oldest:
				// The extreme left the window. Scan it from the newest value
				// to the oldest, so that on ties the newest one is kept.
				best = -1
				for k := t; k >= oldest; k-- {
					v := h[k%w]
					if v == v && (best == -1 || better(v, h[best%w])) {
						best = k
					}
				}
			case !better(h[best%w], n):
				best = t
			}

			out := n
			if best != -1 {
				out = T(t - best)
			}

			select {
			case <-ctx.Done():
				return
			case r <- out:
			}
		}
	}()

	return r
}
