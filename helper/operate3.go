// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Operate3 applies the provided operate function to corresponding values from three
// numeric input channels and sends the resulting values to an output channel.
//
// Example:
//
//	add := helper.Operate3(ac, bc, cc, func(a, b, c int) int {
//	  return a + b + c
//	})
func Operate3[A any, B any, C any, R any](ac <-chan A, bc <-chan B, cc <-chan C, o func(A, B, C) R) <-chan R {
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

			rc <- o(an, bn, cn)
		}

		Drain(ac)
		Drain(bc)
		Drain(cc)
	}()

	return rc
}
