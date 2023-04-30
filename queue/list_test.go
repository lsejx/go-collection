package queue

import "testing"

func TestNewQueue(t *testing.T) {
	q := NewQueue[any]()
	if q.head != nil {
		t.Fatal("head != nil")
	}
	if q.tail != nil {
		t.Fatal("tail != nil")
	}
}

func TestEnqueue(t *testing.T) {
	type state struct {
		h *data[int]
		t *data[int]
	}
	t1 := &data[int]{4, nil}
	h1 := &data[int]{1, &data[int]{2, &data[int]{3, t1}}}
	t2 := &data[int]{100, nil}
	h2 := &data[int]{6, &data[int]{0, t2}}
	tests := []struct {
		id   string
		ini  state
		args []int
		want []int // head to tail
	}{
		{"nil-one", state{nil, nil}, []int{1}, []int{1}},
		{"nil-some", state{nil, nil}, []int{2, 4, 6, 8}, []int{2, 4, 6, 8}},
		{"some-one", state{h1, t1}, []int{1}, []int{1, 2, 3, 4, 1}},
		{"some-some", state{h2, t2}, []int{20, 90}, []int{6, 0, 100, 20, 90}},
	}
	for _, tt := range tests {
		q := Queue[int]{tt.ini.h, tt.ini.t}
		for _, a := range tt.args {
			q.Enqueue(a)
		}
		cur := q.head
		for _, w := range tt.want {
			if cur.v != w {
				t.Fatalf("id:%v, v:%v, w:%v", tt.id, cur.v, w)
			}
			cur = cur.next
		}
		if cur != nil {
			t.Fatalf("id:%v, extradata:%v", tt.id, cur.v)
		}
	}
}

func TestDequeue(t *testing.T) {
	type (
		state struct {
			h *data[int]
			t *data[int]
		}
		retT struct {
			v  int
			ok bool
		}
	)
	d1 := &data[int]{5, nil}
	t1 := &data[int]{10, nil}
	h1 := &data[int]{1, &data[int]{700, t1}}
	tests := []struct {
		id   string
		ini  state
		rets []retT
	}{
		{"nil", state{nil, nil}, []retT{}},
		{"one", state{d1, d1}, []retT{{5, true}}},
		{"some", state{h1, t1}, []retT{{1, true}, {700, true}, {10, true}}},
	}
	for _, tt := range tests {
		q := Queue[int]{tt.ini.h, tt.ini.t}
		for _, r := range tt.rets {
			v, ok := q.Dequeue()
			if v != r.v || ok != r.ok {
				t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, r.v, r.ok)
			}
		}
		if q.head != nil {
			t.Fatalf("id:%v, extrahead:%v", tt.id, q.head.v)
		}
		if q.tail != nil {
			t.Fatalf("id:%v, extratail:%v", tt.id, q.tail.v)
		}
	}
}
