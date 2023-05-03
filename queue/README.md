# go-collection/queue
|type|spec|
|:---|:---|
|Queue|linked data structure|
|BufferedQueue|fixed size ring buffer|
<br><br>

# Import
	import "github.com/lsejx/go-collection/queue

# Example
## Queue
	q := queue.NewQueue[int]()

	q.Enqueue(1)

	n, ok := q.Dequeue()
	// 1, true

	n, ok = q.Dequeue()
	// 0, false

## BufferedQueue
	q := queue.NewBufferedQueue[int](1)

	c := q.Cap()
	// 1

	l := q.Len()
	// 0

	err := q.Enqueue(1)
	// nil

	n, ok := q.Dequeue()
	// 1, true

	n, ok = q.Dequeue()
	// 0, false

	err = q.Enqueue(2)
	// nil

	err = q.Enqueue(3)
	// ErrBufferOverflow