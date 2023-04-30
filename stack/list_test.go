package stack

import "testing"

func TestNewStack(t *testing.T) {
	s := NewStack[any]()
	if s.top != nil {
		t.Fatal(*s.top)
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		id    string
		ini   *data[int]
		args  []int
		wants []int
	}{
		{"nil-one", nil, []int{5}, []int{5}},
		{"nil-some", nil, []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"some-one", &data[int]{2, &data[int]{1, nil}}, []int{3}, []int{3, 2, 1}},
		{"some-some", &data[int]{2, &data[int]{1, nil}}, []int{3, 4, 5, 6}, []int{6, 5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		s := &Stack[int]{tt.ini}
		for _, a := range tt.args {
			s.Push(a)
		}
		cur := s.top
		for _, w := range tt.wants {
			if cur.v != w {
				t.Fatalf("id:%v, v:%v, w:%v", tt.id, cur.v, w)
			}
			cur = cur.prev
		}
		if cur != nil {
			t.Fatalf("id:%v, extradata:%v", tt.id, cur.v)
		}
	}
}

func TestPop(t *testing.T) {
	type retT struct {
		v  int
		ok bool
	}
	tests := []struct {
		id   string
		ini  *data[int]
		rets []retT
	}{
		{"nil", nil, []retT{{0, false}}},
		{"one", &data[int]{5, nil}, []retT{{5, true}}},
		{"some", &data[int]{5, &data[int]{4, &data[int]{3, nil}}}, []retT{{5, true}, {4, true}, {3, true}}},
	}
	for _, tt := range tests {
		s := &Stack[int]{tt.ini}
		for _, ret := range tt.rets {
			v, ok := s.Pop()
			if ok != ret.ok {
				t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, ret.v, ret.ok)
			}
			if v != ret.v {
				t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, ret.v, ret.ok)
			}
		}
		if v, ok := s.Pop(); ok {
			t.Fatalf("id:%v, extradata:%v", tt.id, v)
		}
	}
}

func TestStack(t *testing.T) {
	type retT struct {
		v  int
		ok bool
	}
	tests := []struct {
		id string
		pu []int
		po []retT
	}{
		{"one", []int{100}, []retT{{100, true}}},
		{"some", []int{5, 10, 15, 20}, []retT{{20, true}, {15, true}, {10, true}, {5, true}}},
	}
	for _, tt := range tests {
		s := NewStack[int]()
		for _, v := range tt.pu {
			s.Push(v)
		}
		for _, w := range tt.po {
			v, ok := s.Pop()
			if ok != w.ok {
				t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, w.v, w.ok)
			}
			if v != w.v {
				t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, w.v, w.ok)
			}
		}
		if v, ok := s.Pop(); ok {
			t.Fatalf("id:%v, extradata:%v", tt.id, v)
		}
	}
}
