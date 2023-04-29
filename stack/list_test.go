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
		id   string
		ini  *data[int]
		a    []int
		want []int
	}{
		{"nilone", nil, []int{5}, []int{5}},
		{"nilsome", nil, []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"someone", &data[int]{2, &data[int]{1, nil}}, []int{3}, []int{3, 2, 1}},
		{"somesome", &data[int]{2, &data[int]{1, nil}}, []int{3, 4, 5, 6}, []int{6, 5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		s := &Stack[int]{tt.ini}
		for _, a := range tt.a {
			s.Push(a)
		}
		cur := s.top
		for _, w := range tt.want {
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
		want []retT
	}{
		{"nil", nil, []retT{{0, false}}},
		{"one", &data[int]{5, nil}, []retT{{5, true}}},
		{"some", &data[int]{5, &data[int]{4, &data[int]{3, nil}}}, []retT{{5, true}, {4, true}, {3, true}}},
	}
	for _, tt := range tests {
		s := &Stack[int]{tt.ini}
		for _, w := range tt.want {
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
