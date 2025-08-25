// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Window returns a channel that emits the passed function result
// within a sliding window of size w from the input channel c.
// Note: the slice is in the same order than in source channel
// but the 1st element may not be 0, use modulo window size if
// order is important.
func Window[T any](c <-chan T, f func([]T, int) T, w int) <-chan T {
	r := make(chan T)

	if w <= 0 {
		close(r)
		return r
	}
	
	go func() {
		defer close(r)
		h := make([]T, w)
		n, cnt := 0, 0

		for ok := true; ok; {
			if h[n], ok = <-c; ok {
				if cnt < w {
					cnt++
					r <- f(h[:cnt], 0)
				} else {
					r <- f(h, (n+1)%w)
				}
			}
			n = (n + 1) % w
		}
	}()

	return r
}
