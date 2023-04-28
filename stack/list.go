package stack

type data[T any] struct {
	v    T
	prev *data[T]
}

// Stack is stack using list-structure.
type Stack[T any] struct {
	top *data[T]
}

// Push pushes a value to stack.
func (s *Stack[T]) Push(value T) {
	d := &data[T]{
		v:    value,
		prev: s.top,
	}
	s.top = d
}

// Pop returns true if any value is existing in stack, false if stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	if s.top == nil {
		var v T
		return v, false
	}
	ret := s.top.v
	s.top = s.top.prev
	return ret, true
}

// NewLStack returns empty Stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		top: nil,
	}
}
