package container

// Comparable interface.
type Comparable interface {
	// Compare with other value. Returns -1 if less than, 0 if
	// equals, and 1 if greather than the other value.
	Compare(other Comparable) int
}

// Compares first and second values. The given values must be
// numeric or must implement Comparable interface.
//
// Returns -1 if less than, 0 if equals, 1 if greather than.
func Compare(first, second interface{}) int {
	if _, ok := first.(Comparable); ok {
		return first.(Comparable).Compare(second.(Comparable))
	}

	switch first.(type) {
	case float64:
		return compareFloat64(first.(float64), second.(float64))

	case int64:
		return compareInt64(first.(int64), second.(int64))
	}

	panic("not comparable")
}

func compareFloat64(first, second float64) int {
	if first < second {
		return -1
	}

	if first > second {
		return 1
	}

	return 0
}

func compareInt64(first, second int64) int {
	if first < second {
		return -1
	}

	if first > second {
		return 1
	}

	return 0
}
