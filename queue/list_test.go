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
	type (
		state struct {
			h *data[int]
			t *data[int]
		}
	)
	t1 := &data[int]{4, nil}
	h1 := &data[int]{1, &data[int]{2, &data[int]{3, t1}}}
	tests := []struct {
		id    string
		ini   state
		arg   int
		after []int // head to tail
	}{
		{"nil", state{nil, nil}, 1, []int{1}},
		{"some-one", state{h1, t1}, 5, []int{1, 2, 3, 4, 5}},
	}
	for _, tt := range tests {
		q := Queue[int]{tt.ini.h, tt.ini.t}
		q.Enqueue(tt.arg)
		cur := q.head
		var tail *data[int]
		for _, w := range tt.after {
			if cur.v != w {
				t.Fatalf("id:%v, v:%v, w:%v", tt.id, cur.v, w)
			}
			if cur.next == nil {
				tail = cur
			}
			cur = cur.next
		}
		if cur != nil {
			t.Fatalf("id:%v, extradata:%v", tt.id, cur.v)
		}
		if q.tail != tail {
			t.Fatalf("id:%v, unlinkedtail", tt.id)
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
		id    string
		ini   state
		ret   retT
		after []int
	}{
		{"nil", state{nil, nil}, retT{0, false}, []int{}},
		{"one", state{d1, d1}, retT{5, true}, []int{}},
		{"some", state{h1, t1}, retT{1, true}, []int{700, 10}},
	}
	for _, tt := range tests {
		q := Queue[int]{tt.ini.h, tt.ini.t}
		v, ok := q.Dequeue()
		if v != tt.ret.v || ok != tt.ret.ok {
			t.Fatalf("id:%v, got:(%v,%v), w:(%v,%v)", tt.id, v, ok, tt.ret.v, tt.ret.ok)
		}
		cur := q.head
		var tail *data[int]
		for _, w := range tt.after {
			if cur.v != w {
				t.Fatalf("id:%v, v:%v, w:%v", tt.id, cur.v, w)
			}
			if cur.next == nil {
				tail = cur
			}
			cur = cur.next
		}
		if q.tail != tail {
			t.Fatalf("id:%v, unlinkedtail", tt.id)
		}
	}
}

func TestQueue(t *testing.T) {
	type retT struct {
		v  int
		ok bool
	}
	tests := []struct {
		id    string
		isEns []bool
		en    []int
		rets  []retT
	}{
		{"one", []bool{true, false}, []int{100}, []retT{{100, true}}},
		{"allen-allde", []bool{true, true, true, true, true, false, false, false, false, false}, []int{1, 2, 3, 4, 5}, []retT{{1, true}, {2, true}, {3, true}, {4, true}, {5, true}}},
		{"mixed", []bool{true, false, true, true, false, false}, []int{1, 2, 3}, []retT{{1, true}, {2, true}, {3, true}}},
	}
	for _, tt := range tests {
		q := NewQueue[int]()

		enCur, deCur := 0, 0
		for _, isEn := range tt.isEns {
			if isEn {
				q.Enqueue(tt.en[enCur])
				enCur++
			} else {
				v, ok := q.Dequeue()
				w := tt.rets[deCur]
				if v != w.v || ok != w.ok {
					t.Fatalf("id:%v, got:(%v,%v), w:(%v, %v)", tt.id, v, ok, w.v, w.ok)
				}
				deCur++
			}
		}
		if v, ok := q.Dequeue(); ok {
			t.Fatalf("id:%v, extradata:%v", tt.id, v)
		}
	}
}
