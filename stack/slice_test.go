package stack

import (
	"errors"
	"testing"
)

func TestNewBufferedStack(t *testing.T) {
	tests := []struct {
		a uint
	}{
		{0},
		{1},
		{5},
		{128},
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

func TestBPush(t *testing.T) {
	tests := []struct {
		id   string
		ini  []int
		args []int
		want []int
		rets []error
	}{
		// nil-error case
		{"nil-one", make([]int, 0, 1), []int{1}, []int{1}, []error{nil}},
		{"nil-some", make([]int, 0, 3), []int{0, 5, 10}, []int{0, 5, 10}, []error{nil, nil, nil}},
		{"some-one", []int{1, 2, 3, 0}[:3], []int{5}, []int{1, 2, 3, 5}, []error{nil, nil, nil}},
		{"some-some", []int{1, 2, 3, 0, 0}[:3], []int{5, 100}, []int{1, 2, 3, 5, 100}, []error{nil, nil, nil}},
		// error case
		{"nil-0buf-one-err", make([]int, 0), []int{1}, make([]int, 0), []error{ErrBufferOverflow}},
		{"nil-0buf-some-err", make([]int, 0), []int{1, 4, 9}, make([]int, 0), []error{ErrBufferOverflow, ErrBufferOverflow, ErrBufferOverflow}},
		{"nil-somebuf-some-err", make([]int, 0, 2), []int{5, 25, 125}, []int{5, 25}, []error{nil, nil, ErrBufferOverflow}},
		{"some-0buf-one-err", []int{8, 2000}, []int{1}, []int{8, 2000}, []error{ErrBufferOverflow}},
		{"some-0buf-some-err", []int{1, 6, 100}, []int{1, 2}, []int{1, 6, 100}, []error{ErrBufferOverflow, ErrBufferOverflow}},
		{"some-somebuf-some-err", []int{1, 2, 0}[:2], []int{3, 9, 27}, []int{1, 2, 3}, []error{nil, ErrBufferOverflow, ErrBufferOverflow}},
	}
	for _, tt := range tests {
		s := &BufferedStack[int]{tt.ini}
		for i, a := range tt.args {
			err := s.Push(a)
			if !errors.Is(err, tt.rets[i]) {
				t.Fatalf("id:%v, err:%v, a:%v", tt.id, err, a)
			}
		}
		for i, w := range tt.want {
			if s.buf[i] != w {
				t.Fatalf("id:%v, v:%v, w:%v", tt.id, s.buf[i], w)
			}
		}
	}
}

func TestBPop(t *testing.T) {
	type retT struct {
		v  int
		ok bool
	}
	tests := []struct {
		id   string
		ini  []int
		rets []retT
	}{
		{"nil", make([]int, 0), make([]retT, 0)},
		{"one", []int{4}, []retT{{4, true}}},
		{"some", []int{2, 41, 99}, []retT{{99, true}, {41, true}, {2, true}}},
	}
	for _, tt := range tests {
		s := &BufferedStack[int]{tt.ini}
		for _, r := range tt.rets {
			v, ok := s.Pop()
			if v != r.v || ok != r.ok {
				t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, r.v, r.ok)
			}
		}
		if len(s.buf) != 0 {
			t.Fatalf("id:%v, extradata:%v", tt.id, s.buf)
		}
		if v, ok := s.Pop(); ok {
			t.Fatalf("id:%v, extrapop:%v", tt.id, v)
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
