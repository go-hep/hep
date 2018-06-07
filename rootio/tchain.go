// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

type tchain struct {
	trees []Tree
	offs  []int64 // number of entries before this tree
	tots  []int64 // total number of entries up to this tree

	cur  int   // index of current tree
	tree Tree  // current tree
	off  int64 // current offset
	tot  int64 // current number of entries
}

// Chain returns a tchain that is the concatenation of all the input Trees.
func Chain(trees ...Tree) Tree {
	if len(trees) == 0 {
		return &tchain{}
	}
	n := len(trees)
	ch := &tchain{
		trees: make([]Tree, n),
		offs:  make([]int64, n),
		tots:  make([]int64, n),
		cur:   -1,
	}
	var (
		sum int64
		off int64
	)
	for i := range trees {
		t := trees[i]
		n := t.Entries()
		sum += n
		ch.trees[i] = t
		ch.offs[i] = off
		ch.tots[i] = sum
		off += n
	}

	ch.loadTree(ch.cur + 1)
	return ch
}

func (ch *tchain) loadTree(i int) {
	ch.cur = i
	if ch.cur >= len(ch.trees) {
		ch.tree = nil
		return
	}
	ch.tree = ch.trees[ch.cur]
	ch.off = ch.offs[ch.cur]
	ch.tot = ch.tots[ch.cur]
	return
}

// Class returns the ROOT class of the argument.
func (*tchain) Class() string {
	return "TChain"
}

// Name returns the name of the ROOT objet in the argument.
func (t *tchain) Name() string {
	if t.tree == nil {
		return ""
	}
	return t.tree.Name()
}

// Title returns the title of the ROOT object in the argument.
func (t *tchain) Title() string {
	if t.tree == nil {
		return ""
	}
	return t.tree.Title()
}

// Entries returns the total number of entries.
func (t *tchain) Entries() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.Entries()
	}
	return v
}

// TotBytes return the total number of bytes before compression.
func (t *tchain) TotBytes() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.TotBytes()
	}
	return v
}

// ZipBytes returns the total number of bytes after compression.
func (t *tchain) ZipBytes() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.ZipBytes()
	}
	return v

}

// Branches returns the list of branches.
func (t *tchain) Branches() []Branch {
	if t.tree == nil {
		return nil
	}
	return t.tree.Branches()
}

// Branch returns the branch whose name is the argument.
func (t *tchain) Branch(name string) Branch {
	if t.tree == nil {
		return nil
	}
	return t.tree.Branch(name)
}

// Leaves returns direct pointers to individual branch leaves.
func (t *tchain) Leaves() []Leaf {
	if t.tree == nil {
		return nil
	}
	return t.tree.Leaves()
}

// Leaf returns the leaf whose name is the argument.
func (t *tchain) Leaf(name string) Leaf {
	if t.tree == nil {
		return nil
	}
	return t.tree.Leaf(name)
}

// getFile returns the underlying file.
func (t *tchain) getFile() *File {
	if t.tree == nil {
		return nil
	}
	return t.tree.getFile()
}

// loadEntry returns an error if there is a problem during the loading.
func (t *tchain) loadEntry(i int64) error {
	if t.tree == nil {
		return nil
	}
	j := i - t.off
	return t.tree.loadEntry(j)
}

var (
	_ Object = (*tchain)(nil)
	_ Named  = (*tchain)(nil)
	_ Tree   = (*tchain)(nil)
)
