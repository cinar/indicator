// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Operate4 wraps Operate4WithContext for backwards compatibility.
//
// Deprecated: Use Operate4WithContext instead.
func Operate4[A any, B any, C any, D any, R any](ac <-chan A, bc <-chan B, cc <-chan C, dc <-chan D, o func(A, B, C, D) R) <-chan R {
	return Operate4WithContext(context.Background(), ac, bc, cc, dc, o)
}

// Operate4WithContext applies the provided operate function to corresponding values from four
// numeric input channels and sends the resulting values to an output channel, supporting context cancellation.
func Operate4WithContext[A any, B any, C any, D any, R any](ctx context.Context, ac <-chan A, bc <-chan B, cc <-chan C, dc <-chan D, o func(A, B, C, D) R) <-chan R {
	rc := make(chan R)

	go func() {
		defer func() {
			close(rc)
			DrainWithContext(ctx, ac)
			DrainWithContext(ctx, bc)
			DrainWithContext(ctx, cc)
			DrainWithContext(ctx, dc)
		}()

		for {
			var an A
			var bn B
			var cn C
			var dn D
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
			case dn, ok = <-dc:
				if !ok {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case rc <- o(an, bn, cn, dn):
			}
		}
	}()

	return rc
}
