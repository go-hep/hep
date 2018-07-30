// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

// pair is an item in the heap. It contains the indices of the two jets to be
// clustered and their kt distance.
type pair struct {
	// jeti and jetj are the indices of the jets to be recombined.
	jeti, jetj int
	// dij is the kt distance for the two particles.
	dij float64
}

// Heap contains a slice of items and the last index.
type Heap struct {
	n     int // n is the last index of the slice
	items []pair
}

// New returns a heap pointer.
func New() *Heap {
	h := &Heap{n: 0}
	h.items = make([]pair, 1)
	return h
}

// Push inserts two new clustering candidates and their kt distance.
func (h *Heap) Push(jeti, jetj int, dij float64) {
	item := pair{jeti: jeti, jetj: jetj, dij: dij}
	h.n++
	if h.n >= len(h.items) {
		h.items = append(h.items, item)
	} else {
		h.items[h.n] = item
	}
	h.moveUp(h.n)
}

// Pop returns the two jets with the smallest distance.
// It returns -1, -1, 0 if the heap is empty.
func (h *Heap) Pop() (jeti, jetj int, dij float64) {
	if h.n == 0 {
		return -1, -1, 0
	}
	item := h.items[1]
	h.n--
	if h.n == 0 {
		return item.jeti, item.jetj, item.dij
	}
	h.swap(1, h.n+1)
	h.moveDown(1)
	return item.jeti, item.jetj, item.dij
}

// IsEmpty returns whether a heap is empty.
func (h *Heap) IsEmpty() bool {
	return h.n == 0
}

// moveUp compares an item with its parent and moves it up if it has
// a smaller distance. It will keep moving up until the parent's distance is smaller
// or if it reaches the top.
func (h *Heap) moveUp(i int) {
	for {
		parent := int(i / 2)
		if parent == 0 {
			return
		}
		if h.less(i, parent) {
			h.swap(i, parent)
			i = parent
			continue
		}
		return
	}
}

// moveDown compares an item to its children and moves it down
// if one of the children has a smaller distance. It will keep
// moving down until it reaches a leaf or both children have
// a bigger distance.
func (h *Heap) moveDown(i int) {
	for {
		left := 2 * i
		if left > h.n {
			return
		}
		right := 2*i + 1
		var smallestChild int
		switch {
		case left == h.n:
			smallestChild = left
		case h.less(left, right):
			smallestChild = left
		default:
			smallestChild = right
		}
		if h.less(smallestChild, i) {
			h.swap(i, smallestChild)
			i = smallestChild
			continue
		}
		return
	}
}

// less returns whether item at index i has a smaller distance
// than the item at index j.
func (h *Heap) less(i, j int) bool {
	return h.items[i].dij < h.items[j].dij
}

// swap swaps the two items at the indices i and j.
func (h *Heap) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}
