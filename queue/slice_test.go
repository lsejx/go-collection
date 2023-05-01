package queue

import (
	"errors"
	"testing"

	"golang.org/x/exp/slices"
)

func TestNewBufferedQueue(t *testing.T) {
	tests := []struct {
		arg uint
	}{
		{0},
		{5},
	}
	for _, tt := range tests {
		q := NewBufferedQueue[any](tt.arg)
		if uint(len(q.buf)) != tt.arg+1 {
			t.Fatalf("a:%v, len:%v", tt.arg, len(q.buf))
		}
		if q.head != 0 {
			t.Fatalf("head:%v", q.head)
		}
		if q.tail != tt.arg {
			t.Fatalf("a:%v, tail:%v", tt.arg, q.tail)
		}
	}
}

func TestLocalCap(t *testing.T) {
	type state struct {
		b []int
		h uint
		t uint
	}
	tests := []struct {
		ini  state
		want uint
	}{
		{state{make([]int, 1), 0, 0}, 1},
		{state{make([]int, 5), 0, 4}, 5},
	}
	for _, tt := range tests {
		q := &BufferedQueue[int]{tt.ini.b, tt.ini.h, tt.ini.t}
		if q.cap() != tt.want {
			t.Fatalf("cap:%v, want:%v", q.cap(), tt.want)
		}
	}
}

func TestGlobalCap(t *testing.T) {
	type state struct {
		b []int
		h uint
		t uint
	}
	tests := []struct {
		ini  state
		want uint
	}{
		{state{make([]int, 1), 0, 0}, 0},
		{state{make([]int, 5), 0, 4}, 4},
	}
	for _, tt := range tests {
		q := &BufferedQueue[int]{tt.ini.b, tt.ini.h, tt.ini.t}
		if q.Cap() != tt.want {
			t.Fatalf("cap:%v, want:%v", q.Cap(), tt.want)
		}
	}
}

func TestLen(t *testing.T) {
	type state struct {
		b []int
		h uint
		t uint
	}
	tests := []struct {
		ini  state
		want uint
	}{
		{state{make([]int, 1), 0, 0}, 0},
		{state{make([]int, 2), 0, 0}, 1},
		{state{make([]int, 2), 0, 1}, 0},
		{state{make([]int, 2), 1, 0}, 0},
		{state{make([]int, 2), 1, 1}, 1},
		{state{make([]int, 3), 0, 0}, 1},
		{state{make([]int, 3), 0, 1}, 2},
		{state{make([]int, 3), 0, 2}, 0},
		{state{make([]int, 3), 1, 0}, 0},
		{state{make([]int, 3), 1, 1}, 1},
		{state{make([]int, 3), 1, 2}, 2},
		{state{make([]int, 3), 2, 0}, 2},
		{state{make([]int, 3), 2, 1}, 0},
		{state{make([]int, 3), 2, 2}, 1},
		{state{make([]int, 4), 0, 0}, 1},
		{state{make([]int, 4), 0, 1}, 2},
		{state{make([]int, 4), 0, 2}, 3},
		{state{make([]int, 4), 0, 3}, 0},
		{state{make([]int, 4), 1, 0}, 0},
		{state{make([]int, 4), 1, 1}, 1},
		{state{make([]int, 4), 1, 2}, 2},
		{state{make([]int, 4), 1, 3}, 3},
		{state{make([]int, 4), 2, 0}, 3},
		{state{make([]int, 4), 2, 1}, 0},
		{state{make([]int, 4), 2, 2}, 1},
		{state{make([]int, 4), 2, 3}, 2},
		{state{make([]int, 4), 3, 0}, 2},
		{state{make([]int, 4), 3, 1}, 3},
		{state{make([]int, 4), 3, 2}, 0},
		{state{make([]int, 4), 3, 3}, 1},
	}
	for _, tt := range tests {
		q := &BufferedQueue[int]{tt.ini.b, tt.ini.h, tt.ini.t}
		if q.Len() != tt.want {
			t.Fatalf("ini:%v, len:%v, want:%v", tt.ini, q.Len(), tt.want)
		}
	}
}

func TestIsFull(t *testing.T) {
	type state struct {
		b []int
		h uint
		t uint
	}
	tests := []struct {
		ini  state
		want bool
	}{
		{state{make([]int, 1), 0, 0}, true},
		{state{make([]int, 2), 0, 0}, true},
		{state{make([]int, 2), 0, 1}, false},
		{state{make([]int, 2), 1, 0}, false},
		{state{make([]int, 2), 1, 1}, true},
		{state{make([]int, 3), 0, 0}, false},
		{state{make([]int, 3), 0, 1}, true},
		{state{make([]int, 3), 0, 2}, false},
		{state{make([]int, 3), 1, 0}, false},
		{state{make([]int, 3), 1, 1}, false},
		{state{make([]int, 3), 1, 2}, true},
		{state{make([]int, 3), 2, 0}, true},
		{state{make([]int, 3), 2, 1}, false},
		{state{make([]int, 3), 2, 2}, false},
	}
	for _, tt := range tests {
		q := &BufferedQueue[int]{tt.ini.b, tt.ini.h, tt.ini.t}
		if q.IsFull() != tt.want {
			t.Fatalf("ini:%v, f:%v, w:%v", tt.ini, q.IsFull(), tt.want)
		}
	}
}

func TestBEnqueue(t *testing.T) {
	type state struct {
		b []int
		h uint
		t uint
	}
	tests := []struct {
		id    string
		ini   state
		arg   int
		ret   error
		after state
	}{
		// nil-err case
		{"nil", state{[]int{0, 0, 0}, 0, 2}, 5, nil, state{[]int{5, 0, 0}, 0, 0}},
		{"some", state{[]int{1, 2, 3, 0, 0, 0}, 0, 2}, 5, nil, state{[]int{1, 2, 3, 5, 0, 0}, 0, 3}},
		{"some-ring", state{[]int{1, 3, 5, 7, 9}, 2, 4}, 10, nil, state{[]int{10, 3, 5, 7, 9}, 2, 0}},
		// err case
		{"nil-err", state{[]int{0}, 0, 0}, 1, ErrBufferOverflow, state{[]int{0}, 0, 0}},
		{"some-err", state{[]int{1, 2, 0}, 0, 1}, 3, ErrBufferOverflow, state{[]int{1, 2, 0}, 0, 1}},
		{"some-ring-err", state{[]int{1, 2, 3}, 2, 0}, 4, ErrBufferOverflow, state{[]int{1, 2, 3}, 2, 0}},
	}
	for _, tt := range tests {
		q := &BufferedQueue[int]{tt.ini.b, tt.ini.h, tt.ini.t}
		err := q.Enqueue(tt.arg)
		if !errors.Is(err, tt.ret) {
			t.Fatalf("id:%v, err:%v, a:%v", tt.id, err, tt.arg)
		}
		if !slices.Equal(q.buf, tt.after.b) {
			t.Fatalf("id:%v, b:%v, w:%v", tt.id, q.buf, tt.after.b)
		}
		if q.head != tt.after.h {
			t.Fatalf("id:%v, h:%v, w:%v", tt.id, q.head, tt.after.h)
		}
		if q.tail != tt.after.t {
			t.Fatalf("id:%v, t:%v, w:%v", tt.id, q.tail, tt.after.t)
		}
	}
}

func TestBDequeue(t *testing.T) {
	type (
		state struct {
			b []int
			h uint
			t uint
		}
		retT struct {
			v  int
			ok bool
		}
	)
	tests := []struct {
		id    string
		ini   state
		ret   retT
		after state
	}{
		{"nil-false", state{[]int{0}, 0, 0}, retT{0, false}, state{[]int{0}, 0, 0}},
		{"false", state{[]int{1, 2, 3}, 0, 2}, retT{0, false}, state{[]int{1, 2, 3}, 0, 2}},
		{"ring-false", state{[]int{1, 2, 3}, 2, 1}, retT{0, false}, state{[]int{1, 2, 3}, 2, 1}},
		{"one", state{[]int{1, 2, 3}, 0, 0}, retT{1, true}, state{[]int{1, 2, 3}, 1, 0}},
		{"one-ring", state{[]int{1, 2, 3}, 2, 2}, retT{3, true}, state{[]int{1, 2, 3}, 0, 2}},
		{"some", state{[]int{1, 2, 3}, 0, 1}, retT{1, true}, state{[]int{1, 2, 3}, 1, 1}},
		{"some-ring", state{[]int{1, 2, 3}, 2, 0}, retT{3, true}, state{[]int{1, 2, 3}, 0, 0}},
	}
	for _, tt := range tests {
		q := &BufferedQueue[int]{tt.ini.b, tt.ini.h, tt.ini.t}
		v, ok := q.Dequeue()
		if v != tt.ret.v || ok != tt.ret.ok {
			t.Fatalf("id:%v, v:(%v,%v), w:(%v,%v)", tt.id, v, ok, tt.ret.v, tt.ret.ok)
		}
		if !slices.Equal(q.buf, tt.after.b) {
			t.Fatalf("id:%v, b:%v, w:%v", tt.id, q.buf, tt.after.b)
		}
		if q.head != tt.after.h {
			t.Fatalf("id:%v, h:%v, w:%v", tt.id, q.head, tt.after.h)
		}
		if q.tail != tt.after.t {
			t.Fatalf("id:%v, t:%v, w:%v", tt.id, q.tail, tt.after.t)
		}
	}
}
