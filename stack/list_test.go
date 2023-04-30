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
		arg   int
		after []int
	}{
		{"nil", nil, 5, []int{5}},
		{"some", &data[int]{2, &data[int]{1, nil}}, 3, []int{3, 2, 1}},
	}

	for _, tt := range tests {
		s := &Stack[int]{tt.ini}
		s.Push(tt.arg)
		cur := s.top
		for _, w := range tt.after {
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
		id    string
		ini   *data[int]
		ret   retT
		after []int
	}{
		{"nil", nil, retT{0, false}, []int{}},
		{"one", &data[int]{5, nil}, retT{5, true}, []int{}},
		{"some", &data[int]{5, &data[int]{4, &data[int]{3, nil}}}, retT{5, true}, []int{4, 3}},
	}
	for _, tt := range tests {
		s := &Stack[int]{tt.ini}
		v, ok := s.Pop()
		if v != tt.ret.v || ok != tt.ret.ok {
			t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, tt.ret.v, tt.ret.ok)
		}
		cur := s.top
		for _, w := range tt.after {
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

func TestStack(t *testing.T) {
	type retT struct {
		v  int
		ok bool
	}
	tests := []struct {
		id      string
		isPushs []bool // true:Push
		pu      []int
		po      []retT
	}{
		{"one", []bool{true, false}, []int{100}, []retT{{100, true}}},
		{"allpush-allpop", []bool{true, true, true, true, false, false, false, false}, []int{5, 10, 15, 20}, []retT{{20, true}, {15, true}, {10, true}, {5, true}}},
		{"mixed", []bool{true, false, true, true, false, true, false, false}, []int{1, 2, 3, 4}, []retT{{1, true}, {3, true}, {4, true}, {2, true}}},
	}
	for _, tt := range tests {
		s := NewStack[int]()
		puCur, poCur := 0, 0
		for _, isPush := range tt.isPushs {
			if isPush {
				s.Push(tt.pu[puCur])
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
