// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// Ring represents a ring structure that can be instantiated
// using the NewRing function.
//
// Example:
//
//	ring := helper.NewRing[int](2)
//
//	fmt.Println(ring.Insert(1)) // 0
//	fmt.Println(ring.Insert(2)) // 0
//	fmt.Println(ring.Insert(3)) // 1
//	fmt.Println(ring.Insert(4)) // 2
type Ring[T any] struct {
	buffer []T
	begin  int
	end    int
	empty  bool
	count  int
}

// NewRing creates a new ring instance with the given size. A ring
// cannot function with zero or negative capacity, so a size less than
// 1 is clamped to 1.
func NewRing[T any](size int) *Ring[T] {
	if size < 1 {
		size = 1
	}

	return &Ring[T]{
		buffer: make([]T, size),
		begin:  0,
		end:    0,
		empty:  true,
		count:  0,
	}
}

// Put inserts the specified value into the ring and returns the
// value that was previously stored at that index.
func (r *Ring[T]) Put(t T) T {
	if r.IsFull() {
		r.begin = r.nextIndex(r.begin)
	} else {
		r.count++
	}

	o := r.buffer[r.end]
	r.buffer[r.end] = t

	r.end = r.nextIndex(r.end)
	r.empty = false

	return o
}

// Get retrieves the available value from the ring buffer. If empty,
// it returns the default value (T) and false.
func (r *Ring[T]) Get() (T, bool) {
	var t T

	if r.empty {
		return t, false
	}

	t = r.buffer[r.begin]
	r.begin = r.nextIndex(r.begin)
	r.count--

	if r.begin == r.end {
		r.empty = true
	}

	return t, true
}

// At returns the value at the given index, relative to the oldest
// element currently in the ring. It returns false if index is out
// of range, or fewer than index+1 elements have ever been Put.
func (r *Ring[T]) At(index int) (T, bool) {
	var t T

	if index < 0 || index >= r.count {
		return t, false
	}

	return r.buffer[(r.begin+index)%len(r.buffer)], true
}

// IsEmpty checks if the current ring buffer is empty.
func (r *Ring[T]) IsEmpty() bool {
	return r.empty
}

// IsFull checks if the current ring buffer is full.
func (r *Ring[T]) IsFull() bool {
	return !r.empty && (r.end == r.begin)
}

// nextIndex returns the next index in a ring buffer, wrapping
// around if it reaches the capacity.
func (r *Ring[T]) nextIndex(i int) int {
	return (i + 1) % len(r.buffer)
}
