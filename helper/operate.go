// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Operate applies the provided operate function to corresponding values from two
// numeric input channels and sends the resulting values to an output channel.
//
// Example:
//
//	add := helper.Operate(ac, bc, func(a, b int) int {
//	  return a + b
//	})
func Operate[A any, B any, R any](ac <-chan A, bc <-chan B, o func(A, B) R) <-chan R {
	oc := make(chan R)

	go func() {
		defer close(oc)

		for {
			an, ok := <-ac
			if !ok {
				Drain(bc)
				break
			}

			bn, ok := <-bc
			if !ok {
				Drain(ac)
				break
			}

			oc <- o(an, bn)
		}
	}()

	return oc
}
