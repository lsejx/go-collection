package stack

import "errors"

var ErrBufferOverflow = errors.New("buffer overflow")

// BufferedStack is stack using slice.
// Its buffer size doesn't grow.
type BufferedStack[T any] struct {
	buf []T
}

// Push pushed a value to stack.
// Push doesn't use built-in append function, so the buffer size doesn't grow.
// If buffer is full before pushing, ErrBufferOverflow is returned.
func (s *BufferedStack[T]) Push(value T) error {
	if len(s.buf) == cap(s.buf) {
		return ErrBufferOverflow
	}
	s.buf = s.buf[:len(s.buf)+1]
	s.buf[len(s.buf)-1] = value
	return nil
}

// Pop returns (value, true) if any value is existing in stack, (default-value, false) if stack is empty.
func (s *BufferedStack[T]) Pop() (T, bool) {
	if len(s.buf) == 0 {
		var v T
		return v, false
	}
	ret := s.buf[len(s.buf)-1]
	s.buf = s.buf[:len(s.buf)-1]
	return ret, true
}

// NewBufferedStack returns empty BufferedStack which has specified capacity.
func NewBufferedStack[T any](capacity uint) *BufferedStack[T] {
	return &BufferedStack[T]{
		buf: make([]T, 0, capacity),
	}
}
