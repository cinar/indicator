// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Window wraps WindowWithContext for backwards compatibility.
//
// Deprecated: Use WindowWithContext instead.
func Window[T any](c <-chan T, f func([]T, int) T, w int) <-chan T {
	return WindowWithContext(context.Background(), c, f, w)
}

// WindowWithContext returns a channel that emits the passed function result
// within a sliding window of size w from the input channel c, supporting context cancellation.
//
// The aggregation function f receives (s, i), where s is the current window's
// backing slice and i is the rotation offset of the oldest (chronologically
// first) element in s. The backing slice is a fixed-size ring buffer that is
// reused and rotated in place as new values arrive, so once the window has
// filled up (len(s) == w), a plain left-to-right scan of s (e.g. range s, or
// s[0], s[1], ...) does NOT visit the values in chronological order — it
// visits them in whatever order they currently sit in the ring.
//
// A custom f that cares about temporal order must use i, not raw position:
//   - Oldest-to-newest: the k-th oldest element (0-based) is s[(i+k)%len(s)].
//     Equivalently, s rotated left by i is the chronological sequence.
//   - Newest-to-oldest: use the SlicesReverse(s, i, ...) helper, which starts
//     just before i (the newest element) and walks backward, wrapping around,
//     stopping at i (the oldest element). See max_since.go and min_since.go
//     for a worked example.
//   - While the window is still filling (len(s) < w), i is always 0 and s has
//     not wrapped yet, so it is already in chronological order as-is.
//
// Order-independent aggregations (e.g. min/max of the set, as in highest.go
// and lowest.go) can safely ignore i, since any permutation of the same
// values yields the same result.
func WindowWithContext[T any](ctx context.Context, c <-chan T, f func([]T, int) T, w int) <-chan T {
	r := make(chan T)

	if w <= 0 {
		close(r)
		return r
	}

	go func() {
		defer close(r)
		h := make([]T, w)
		n, cnt := 0, 0

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-c:
				if !ok {
					return
				}
				h[n] = val
				var out T
				if cnt < w {
					cnt++
					// Window not yet full: h[:cnt] has not wrapped, so it is
					// already in chronological order and the offset is 0.
					out = f(h[:cnt], 0)
				} else {
					// Window full: h is a rotated ring buffer. (n+1)%w is the
					// index of the oldest element; f must use it (rather than
					// scanning h by raw position) to recover chronological
					// order. See the WindowWithContext doc comment above.
					out = f(h, (n+1)%w)
				}
				select {
				case <-ctx.Done():
					return
				case r <- out:
				}
			}
			n = (n + 1) % w
		}
	}()

	return r
}
