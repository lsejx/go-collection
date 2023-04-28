# go-collection/stack
|type|spec|
|:---|:---|
|Stack|linked data structure|
|BufferedStack|fixed size buffer|
<br><br>

# Import
	import "github.com/lsejx/go-collection/stack"

# Example
## Stack
	s := stack.NewStack[int]()

	s.Push(1)

	n, ok := s.Pop()
	// 1, true

	n, ok = s.Pop()
	// 0, false

## BufferedStack
	s := stack.NewBufferedStack[int](1)

	err := s.Push(1)
	// nil

	n, ok := s.Pop()
	// 1, true

	n, ok = s.Pop()
	// 0, false

	err = s.Push(2)
	// nil

	err = s.Push(3)
	// ErrBufferOverflow

