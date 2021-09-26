package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	values := []float64{1, 2, 3, 4}

	queue := New()

	for _, value := range values {
		queue.Enqueue(value)
	}

	for i, value := range values {
		if queue.Empty() {
			t.Fatal("queue empty")
		}

		actual := queue.Dequeue()

		if actual != value {
			t.Fatalf("at %d actual %f expected %f", i, actual, value)
		}
	}
}
