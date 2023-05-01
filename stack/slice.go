package stack

import "errors"

var ErrBufferOverflow = errors.New("buffer overflow")

// BufferedStack is stack using slice.
// Its buffer size doesn't grow.
type BufferedStack[T any] struct {
	buf []T
}

// Cap returns capacity of stack.
func (s *BufferedStack[T]) Cap() uint {
	return uint(cap(s.buf))
}

// Len returns number of items which stack is holding.
func (s *BufferedStack[T]) Len() uint {
	return uint(len(s.buf))
}

// IsFull returns whether stack is full.
func (s *BufferedStack[T]) IsFull() bool {
	return s.Len() == s.Cap()
}

// Push pushed a value to stack.
// Push doesn't use built-in append function, so the buffer size doesn't grow.
// If buffer is full before pushing, ErrBufferOverflow is returned.
func (s *BufferedStack[T]) Push(value T) error {
	if s.IsFull() {
		return ErrBufferOverflow
	}
	s.buf = s.buf[:s.Len()+1]
	s.buf[s.Len()-1] = value
	return nil
}

// Pop returns (value, true) if any value is existing in stack, (default-value, false) if stack is empty.
func (s *BufferedStack[T]) Pop() (T, bool) {
	if s.Len() == 0 {
		var v T
		return v, false
	}
	ret := s.buf[s.Len()-1]
	s.buf = s.buf[:s.Len()-1]
	return ret, true
}

// NewBufferedStack returns empty BufferedStack which has specified capacity.
func NewBufferedStack[T any](capacity uint) *BufferedStack[T] {
	return &BufferedStack[T]{
		buf: make([]T, 0, capacity),
	}
}
