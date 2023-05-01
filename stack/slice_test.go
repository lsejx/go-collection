package stack

import (
	"errors"
	"testing"

	"golang.org/x/exp/slices"
)

func TestNewBufferedStack(t *testing.T) {
	tests := []struct {
		a uint
	}{
		{0},
		{5},
	}
	for _, tt := range tests {
		s := NewBufferedStack[any](tt.a)
		if len(s.buf) != 0 {
			t.FailNow()
		}
		if uint(cap(s.buf)) != tt.a {
			t.Fatalf("a:%v, cap:%v", tt.a, cap(s.buf))
		}
	}
}

func TestCap(t *testing.T) {
	tests := []struct {
		ini  []int
		want uint
	}{
		{make([]int, 0), 0},
		{make([]int, 0, 1), 1},
		{make([]int, 0, 5), 5},
	}
	for _, tt := range tests {
		s := &BufferedStack[int]{tt.ini}
		if s.Cap() != tt.want {
			t.Fatalf("cap:%v, w:%v", s.Cap(), tt.want)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		ini  []int
		want uint
	}{
		{make([]int, 0), 0},
		{make([]int, 1), 1},
		{make([]int, 5), 5},
	}
	for _, tt := range tests {
		s := &BufferedStack[int]{tt.ini}
		if s.Len() != tt.want {
			t.Fatalf("len:%v, w:%v", s.Len(), tt.want)
		}
	}
}

func TestIsFull(t *testing.T) {
	tests := []struct {
		ini  []int
		want bool
	}{
		{make([]int, 0), true},
		{make([]int, 0, 1), false},
		{make([]int, 1), true},
		{make([]int, 0, 2), false},
		{make([]int, 1, 2), false},
		{make([]int, 2), true},
	}
	for _, tt := range tests {
		s := &BufferedStack[int]{tt.ini}
		if s.IsFull() != tt.want {
			t.Fatalf("ini:%v, f:%v, w:%v", tt.ini, s.IsFull(), tt.want)
		}
	}
}

func TestBPush(t *testing.T) {
	tests := []struct {
		id    string
		ini   []int
		arg   int
		after []int
		ret   error
	}{
		// nil-error case
		{"nil", make([]int, 0, 1), 1, []int{1}, nil},
		{"some", []int{1, 2, 3, 0}[:3], 5, []int{1, 2, 3, 5}, nil},
		// error case
		{"nil-err", make([]int, 0), 1, make([]int, 0), ErrBufferOverflow},
		{"some-err", []int{8, 2000}, 1, []int{8, 2000}, ErrBufferOverflow},
	}
	for _, tt := range tests {
		s := &BufferedStack[int]{tt.ini}
		err := s.Push(tt.arg)
		if !errors.Is(err, tt.ret) {
			t.Fatalf("id:%v, err:%v, a:%v", tt.id, err, tt.arg)
		}
		if !slices.Equal(s.buf, tt.after) {
			t.Fatalf("id:%v, b:%v, w:%v", tt.id, s.buf, tt.after)
		}
	}
}

func TestBPop(t *testing.T) {
	type retT struct {
		v  int
		ok bool
	}
	tests := []struct {
		id    string
		ini   []int
		ret   retT
		after []int
	}{
		{"nil", []int{}, retT{0, false}, []int{}},
		{"one", []int{4}, retT{4, true}, []int{}},
		{"some", []int{2, 41, 99}, retT{99, true}, []int{2, 41}},
	}
	for _, tt := range tests {
		s := &BufferedStack[int]{tt.ini}
		v, ok := s.Pop()
		if v != tt.ret.v || ok != tt.ret.ok {
			t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, tt.ret.v, tt.ret.ok)
		}
		if !slices.Equal(s.buf, tt.after) {
			t.Fatalf("id:%v, b:%v, w:%v", tt.id, s.buf, tt.after)
		}
	}
}

func TestBufferedStack(t *testing.T) {
	type retT struct {
		v  int
		ok bool
	}
	tests := []struct {
		id      string
		cap     uint
		isPushs []bool
		pu      []int
		puRets  []error
		po      []retT
	}{
		// nil-err case
		{"one", 1, []bool{true, false}, []int{5}, []error{nil}, []retT{{5, true}}},
		{"allpush-allpop", 3, []bool{true, true, true, false, false, false}, []int{100, 101, 50}, []error{nil, nil, nil}, []retT{{50, true}, {101, true}, {100, true}}},
		{"mixed", 2, []bool{true, true, false, true, false, false}, []int{1, 2, 3}, []error{nil, nil, nil}, []retT{{2, true}, {3, true}, {1, true}}},
		// err case
		{"0buf-err", 0, []bool{true}, []int{1}, []error{ErrBufferOverflow}, []retT{}},
		{"somebuf-err", 2, []bool{true, true, true, false, false}, []int{1, 2, 3}, []error{nil, nil, ErrBufferOverflow}, []retT{{2, true}, {1, true}}},
	}
	for _, tt := range tests {
		s := NewBufferedStack[int](tt.cap)
		puCur, poCur := 0, 0
		for _, isPush := range tt.isPushs {
			if isPush {
				v := tt.pu[puCur]
				err := s.Push(v)
				if !errors.Is(err, tt.puRets[puCur]) {
					t.Fatalf("id:%v, err:%v, v:%v", tt.id, err, v)
				}
				puCur++
			} else {
				v, ok := s.Pop()
				w := tt.po[poCur]
				if v != w.v || ok != w.ok {
					t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, w.v, w.ok)
				}
				poCur++
			}
		}
		if v, ok := s.Pop(); ok {
			t.Fatalf("id:%v, extradata:%v", tt.id, v)
		}
	}
}
