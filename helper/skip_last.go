package helper

// SkipLast skips the specified number of elements
// from the end of the given channel.
//
// Example:
//
//	c := helper.SliceToChan([]int{2, 4, 6, 8})
//	actual := helper.SkipLast(c, 2)
//	fmt.Println(helper.ChanToSlice(actual)) // [2, 4]
func SkipLast[T any](c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		defer close(result)

		// Buffer to hold the last "count" elements
		buf := make([]T, 0, count)

		for v := range c {
			buf = append(buf, v)
			if len(buf) > count {
				// send the oldest value
				result <- buf[0]
				buf = buf[1:]
			}
		}
		// drop the last `count` elements automatically
	}()

	return result
}
