// Copyright (c) 2021-2024 Onur Cinar.
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
}

// NewRing creates a new ring instance with the given size.
func NewRing[T any](size int) *Ring[T] {
	return &Ring[T]{
		buffer: make([]T, size),
		begin:  0,
		end:    0,
		empty:  true,
	}
}

// Put inserts the specified value into the ring and returns the
// value that was previously stored at that index.
func (r *Ring[T]) Put(t T) T {
	if r.IsFull() {
		r.begin = r.nextIndex(r.begin)
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

	if r.begin == r.end {
		r.empty = true
	}

	return t, true
}

// At returns the value at the given index.
func (r *Ring[T]) At(index int) T {
	return r.buffer[(r.begin+index)%len(r.buffer)]
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
