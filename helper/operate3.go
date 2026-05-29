// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Operate3 wraps Operate3WithContext for backwards compatibility.
//
// Deprecated: Use Operate3WithContext instead.
func Operate3[A any, B any, C any, R any](ac <-chan A, bc <-chan B, cc <-chan C, o func(A, B, C) R) <-chan R {
	return Operate3WithContext(context.Background(), ac, bc, cc, o)
}

// Operate3WithContext applies the provided operate function to corresponding values from three
// numeric input channels and sends the resulting values to an output channel, supporting context cancellation.
func Operate3WithContext[A any, B any, C any, R any](ctx context.Context, ac <-chan A, bc <-chan B, cc <-chan C, o func(A, B, C) R) <-chan R {
	rc := make(chan R)

	go func() {
		defer func() {
			close(rc)
			DrainWithContext(ctx, ac)
			DrainWithContext(ctx, bc)
			DrainWithContext(ctx, cc)
		}()

		for {
			var an A
			var bn B
			var cn C
			var ok bool

			select {
			case <-ctx.Done():
				return
			case an, ok = <-ac:
				if !ok {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case bn, ok = <-bc:
				if !ok {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case cn, ok = <-cc:
				if !ok {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case rc <- o(an, bn, cn):
			}
		}
	}()

	return rc
}
