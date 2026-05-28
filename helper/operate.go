// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import "context"

// Operate wraps OperateWithContext for backwards compatibility.
//
// Deprecated: Use OperateWithContext instead.
func Operate[A any, B any, R any](ac <-chan A, bc <-chan B, o func(A, B) R) <-chan R {
	return OperateWithContext(context.Background(), ac, bc, o)
}

// OperateWithContext applies the provided operate function to corresponding values from two
// numeric input channels and sends the resulting values to an output channel, supporting context cancellation.
func OperateWithContext[A any, B any, R any](ctx context.Context, ac <-chan A, bc <-chan B, o func(A, B) R) <-chan R {
	oc := make(chan R)

	go func() {
		defer close(oc)

		for {
			var an A
			var bn B
			var ok bool

			select {
			case <-ctx.Done():
				return
			case an, ok = <-ac:
				if !ok {
					DrainWithContext(ctx, bc)
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case bn, ok = <-bc:
				if !ok {
					DrainWithContext(ctx, ac)
					return
				}
			}

			select {
			case <-ctx.Done():
				return
			case oc <- o(an, bn):
			}
		}
	}()

	return oc
}
