// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import (
	"container/heap"
	"math/rand"
	"testing"
)

// FIXME B/op is higher here, than for PQ
func BenchmarkHeap(b *testing.B) {
	var items []*pair
	for range 5000 {
		jeti := rand.Int()
		jetj := rand.Int()
		dij := rand.Float64()
		items = append(items, &pair{jeti: jeti, jetj: jetj, dij: dij})
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h := New()
		for _, item := range items {
			h.Push(item.jeti, item.jetj, item.dij)
		}
		for !h.IsEmpty() {
			h.Pop()
		}
	}
}

func BenchmarkPQ(b *testing.B) {
	var items []*pair
	for range 5000 {
		jeti := rand.Int()
		jetj := rand.Int()
		dij := rand.Float64()
		items = append(items, &pair{jeti: jeti, jetj: jetj, dij: dij})
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var pq PQ
		heap.Init(&pq)
		for _, item := range items {
			heap.Push(&pq, item)
		}
		for pq.Len() > 0 {
			heap.Pop(&pq)
		}
	}
}

func TestHeapPush(t *testing.T) {
	h := New()
	h.Push(0, 1, 10)
	h.Push(2, 3, 5)
	h.Push(4, 5, 8)
	h.Push(6, 7, 2)
	want := []struct {
		jeti, jetj int
		dij        float64
	}{
		{-1, -1, 0},
		{6, 7, 2},
		{2, 3, 5},
		{4, 5, 8},
		{0, 1, 10},
	}
	got := h.items
	if len(got) != len(want) {
		t.Errorf("got = %d items, want = %d", len(got), len(want))
	}
	for i, pair := range got {
		if i == 0 {
			continue
		}
		if pair.jeti != want[i].jeti {
			t.Errorf("h.items[%d].jeti: got = %d, want = %d", i, pair.jeti, want[i].jeti)
		}
		if pair.jetj != want[i].jetj {
			t.Errorf("h.items[%d].jetj: got = %d, want = %d", i, pair.jetj, want[i].jetj)
		}
		if pair.dij != want[i].dij {
			t.Errorf("h.items[%d].dij: got = %f, want = %f", i, pair.dij, want[i].dij)
		}
	}
}

func TestHeapPop(t *testing.T) {
	pairs := []pair{{}, {jeti: 6, jetj: 7, dij: 2}, {jeti: 2, jetj: 3, dij: 5}, {jeti: 4, jetj: 5, dij: 8},
		{jeti: 0, jetj: 1, dij: 10}, {jeti: 8, jetj: 9, dij: 6},
	}
	h := &Heap{items: pairs, n: 5}
	want := []struct {
		jeti, jetj int
		dij        float64
	}{
		{6, 7, 2},
		{2, 3, 5},
		{8, 9, 6},
		{4, 5, 8},
		{0, 1, 10},
	}
	got := make([]struct {
		jeti, jetj int
		dij        float64
	}, len(pairs)-1)
	for i := 0; !h.IsEmpty(); i++ {
		if i >= len(got) {
			t.Fatalf("Heap with n = %d should be empty", h.n)
		}
		jeti, jetj, dij := h.Pop()
		got[i].jeti = jeti
		got[i].jetj = jetj
		got[i].dij = dij
	}
	if len(got) != len(want) {
		t.Errorf("got = %d items, want = %d", len(got), len(want))
	}
	for i := range got {
		if want[i].jeti != got[i].jeti {
			t.Errorf("got[%d].jeti = %d, want[%d].jeti = %d", i, got[i].jeti, i, want[i].jeti)
		}
		if want[i].jetj != got[i].jetj {
			t.Errorf("got[%d].jetj = %d, want[%d].jetj = %d", i, got[i].jetj, i, want[i].jetj)
		}
		if want[i].dij != got[i].dij {
			t.Errorf("got[%d] = %f, want[%d] = %f", i, got[i].dij, i, want[i].dij)
		}
	}
}

// PQ is a priority queue of pairs. It is implemented using the
// container/heap interface.
type PQ []*pair

func (pq PQ) Len() int { return len(pq) }

func (pq PQ) Less(i, j int) bool {
	return pq[i].dij < pq[j].dij
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PQ) Push(x any) {
	item := x.(*pair)
	*pq = append(*pq, item)
}

func (pq *PQ) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
