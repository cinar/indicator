package queue

import (
	"container/list"
)

// Queue type.
type Queue struct {
	list *list.List
}

// New queue.
func New() *Queue {
	return &Queue{
		list: list.New(),
	}
}

// Enqueue the given value.
func (q *Queue) Enqueue(v interface{}) {
	q.list.PushBack(v)
}

// Dequeue from the queue.
func (q *Queue) Dequeue() interface{} {
	front := q.list.Front()
	if front == nil {
		panic("queue empty")
	}

	value := front.Value
	q.list.Remove(front)

	return value
}

// Queue empty.
func (q *Queue) Empty() bool {
	return (q.list.Len() == 0)
}
