package helper

// SlicesReverse loops through a slice in reverse order starting from the given index.
// The given function is called for each element in the slice. If the function returns false,
// the loop is terminated.
func SlicesReverse[T any](r []T, i int, f func(T) bool) {
	l := len(r)
	if l == 0 || i < 0 || i >= l {
		return
	}
	for m := i - 1; m >= 0; m-- {
		if !f(r[m]) {
			return
		}
	}
	for m := l - 1; m >= i; m-- {
		if !f(r[m]) {
			return
		}
	}
}
