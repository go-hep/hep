// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

type tchain struct {
	trees []itree
	cur   *itree
	icur  int
}

// Class returns the ROOT class of the argument.
func (tchain) Class() string {
	return "TChain"
}

// Name returns the name of the ROOT objet in the argument.
func (t tchain) Name() string {
	if len(t.trees) == 0 {
		return ""
	}
	return t.trees[0].tree.Name()
}

// Title returns the title of the ROOT object in the argument.
func (t tchain) Title() string {
	if len(t.trees) == 0 {
		return ""
	}
	return t.trees[0].tree.Title()
}

// Chain returns a tchain that is the concatenation of all the input Trees.
func Chain(trees ...Tree) (*tchain, error) {
	ch := &tchain{
		trees: make([]itree, len(trees)),
		icur:  -1,
	}
	var sum int64
	var offset int64
	for i, ttree := range trees {
		if len(ttree.Branches()) != len(trees[0].Branches()) {
			return nil, errorf("Trees with different layouts")
		}
		t := trees[i]
		n := t.Entries()
		sum += n
		ch.trees[i] = itree{tree: t, offset: offset, total: sum}
		offset += n
	}
	if len(trees) > 0 {
		ch.icur = 0
		ch.cur = &ch.trees[ch.icur]
	}
	return ch, nil
}

// Entries returns the total number of entries.
func (t tchain) Entries() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.Entries()
	}
	return v
}

// TotBytes return the total number of bytes before compression.
func (t tchain) TotBytes() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.tree.TotBytes()
	}
	return v
}

// ZipBytes returns the total number of bytes after compression.
func (t tchain) ZipBytes() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.tree.ZipBytes()
	}
	return v

}

// Branches returns the list of branches.
func (t tchain) Branches() []Branch {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.Branches()
}

// Branch returns the branch whose name is the argument.
func (t tchain) Branch(name string) Branch {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.Branch(name)
}

// Leaves returns direct pointers to individual branch leaves.
func (t tchain) Leaves() []Leaf {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.Leaves()
}

// getFile returns the underlying file.
func (t tchain) getFile() *File {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.getFile()
}

// loadEntry returns an error if there is a problem during the loading.
func (t tchain) loadEntry(i int64) error {
	if len(t.trees) == 0 {
		return nil
	}
	//return t.trees[0].loadEntry(i)
	if i >= t.cur.Entries() {
		t.icur++
		t.cur = &t.trees[t.icur]
	}
	return t.cur.loadEntry(i)
}

type itree struct {
	tree   Tree
	offset int64
	total  int64
}

func (t itree) Entries() int64 {
	return t.tree.Entries()
}

func (t *itree) loadEntry(i int64) error {
	j := i - t.offset
	return t.tree.loadEntry(j)
}

var (
	_ Object = (*tchain)(nil)
	_ Named  = (*tchain)(nil)
	_ Tree   = (*tchain)(nil)
)
