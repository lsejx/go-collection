package queue

type data[T any] struct {
	v    T
	next *data[T]
}

// Queue is queue using list-structure.
type Queue[T any] struct {
	head *data[T]
	tail *data[T]
}

// Enqueue enqueues a value to queue.
func (q *Queue[T]) Enqueue(value T) {
	d := &data[T]{
		v:    value,
		next: nil,
	}
	if q.tail == nil {
		q.head = d
	} else {
		q.tail.next = d
	}
	q.tail = d
}

// Dequeue returns (value, true) if any value is existing in queue, (default-value, false) if queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.head == nil {
		var v T
		return v, false
	}
	v := q.head.v
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	return v, true
}

// NewQueue returns empty Queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		head: nil,
		tail: nil,
	}
}
