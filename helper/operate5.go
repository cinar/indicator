// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Operate5 wraps Operate5WithContext for backwards compatibility.
//
// Deprecated: Use Operate5WithContext instead.
func Operate5[A any, B any, C any, D any, E any, R any](
	ac <-chan A,
	bc <-chan B,
	cc <-chan C,
	dc <-chan D,
	ec <-chan E,
	o func(A, B, C, D, E) R,
) <-chan R {
	return Operate5WithContext(context.Background(), ac, bc, cc, dc, ec, o)
}

// Operate5WithContext applies the provided operate function to corresponding values from five
// numeric input channels and sends the resulting values to an output channel, supporting context cancellation.
func Operate5WithContext[A any, B any, C any, D any, E any, R any](
	ctx context.Context,
	ac <-chan A,
	bc <-chan B,
	cc <-chan C,
	dc <-chan D,
	ec <-chan E,
	o func(A, B, C, D, E) R,
) <-chan R {
	rc := make(chan R)

	go func() {
		defer func() {
			close(rc)
			DrainWithContext(ctx, ac)
			DrainWithContext(ctx, bc)
			DrainWithContext(ctx, cc)
			DrainWithContext(ctx, dc)
			DrainWithContext(ctx, ec)
		}()

		for {
			var an A
			var bn B
			var cn C
			var dn D
			var en E
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
			case en, ok = <-ec:
				if !ok {
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case rc <- o(an, bn, cn, dn, en):
			}
		}
	}()

	return rc
}
