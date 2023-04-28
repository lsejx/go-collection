package queue

import "errors"

var ErrBufferOverflow = errors.New("buffer overflow")

// BufferedQueue is queue using slice as "ring buffer".
// Its buffer size doesn't grow.
type BufferedQueue[T any] struct {
	buf  []T
	head int
	tail int
}

// Enqueue enqueues a value to queue.
// If buffer is full before enqueuing, ErrBufferOverflow is returned.
func (q *BufferedQueue[T]) Enqueue(value T) error {
	if (q.tail+2)%len(q.buf) == q.head {
		return ErrBufferOverflow
	}
	q.tail = (q.tail + 1) % len(q.buf)
	q.buf[q.tail] = value
	return nil
}

// Dequeue returns (value, true) if any value is existing in queue, (default-value, false) if queue is empty.
func (q *BufferedQueue[T]) Dequeue() (T, bool) {
	if (q.tail+1)%len(q.buf) == q.head {
		var v T
		return v, false
	}
	v := q.buf[q.head]
	q.head = (q.head + 1) % len(q.buf)
	return v, true
}

// NewBufferedQueue returns empty BufferedQueue which has specified size.
func NewBufferedQueue[T any](size int) *BufferedQueue[T] {
	return &BufferedQueue[T]{
		buf:  make([]T, size+1),
		head: 0,
		tail: size,
	}
}
