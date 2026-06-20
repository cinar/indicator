// Copyright (c) 2021-2026 Onur Cinar.
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
					out = f(h[:cnt], 0)
				} else {
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
