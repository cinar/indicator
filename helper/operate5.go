// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Operate5 applies the provided operate function to corresponding values from five
// numeric input channels and sends the resulting values to an output channel.
//
// Example:
//
//	result := helper.Operate5(ac, bc, cc, dc, ec, func(a, b, c, d, e int) int {
//	  return a + b + c + d + e
//	})
func Operate5[A any, B any, C any, D any, E any, R any](
	ac <-chan A,
	bc <-chan B,
	cc <-chan C,
	dc <-chan D,
	ec <-chan E,
	o func(A, B, C, D, E) R,
) <-chan R {
	rc := make(chan R)

	go func() {
		defer close(rc)

		for {
			an, ok := <-ac
			if !ok {
				break
			}

			bn, ok := <-bc
			if !ok {
				break
			}

			cn, ok := <-cc
			if !ok {
				break
			}

			dn, ok := <-dc
			if !ok {
				break
			}

			en, ok := <-ec
			if !ok {
				break
			}

			rc <- o(an, bn, cn, dn, en)
		}

		Drain(ac)
		Drain(bc)
		Drain(cc)
		Drain(dc)
		Drain(ec)
	}()

	return rc
}
