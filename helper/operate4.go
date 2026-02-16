// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Operate4 applies the provided operate function to corresponding values from four
// numeric input channels and sends the resulting values to an output channel.
//
// Example:
//
//	add := helper.Operate4(ac, bc, cc, dc, func(a, b, c, d int) int {
//	  return a + b + c + d
//	})
func Operate4[A any, B any, C any, D any, R any](ac <-chan A, bc <-chan B, cc <-chan C, dc <-chan D, o func(A, B, C, D) R) <-chan R {
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

			rc <- o(an, bn, cn, dn)
		}

		Drain(ac)
		Drain(bc)
		Drain(cc)
		Drain(dc)
	}()

	return rc
}
