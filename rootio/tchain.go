// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

type tchain struct {
	trees []Tree
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
	return t.trees[0].Name()
}

// Title returns the title of the ROOT object in the argument.
func (t tchain) Title() string {
	if len(t.trees) == 0 {
		return ""
	}
	return t.trees[0].Title()
}

// Chain returns a tchain that is the concatenation of all the input Trees.
func Chain(trees ...Tree) tchain {
	var t tchain
	t.trees = append(t.trees, trees...)
	return t
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
		v += tree.TotBytes()
	}
	return v
}

// ZipBytes returns the total number of bytes after compression.
func (t tchain) ZipBytes() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.ZipBytes()
	}
	return v

}

// Branches returns the list of branches.
func (t tchain) Branches() []Branch {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].Branches()
}

// Branch returns the branch whose name is the argument.
func (t tchain) Branch(name string) Branch {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].Branch(name)
}

// Leaves returns direct pointers to individual branch leaves.
func (t tchain) Leaves() []Leaf {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].Leaves()
}

// getFile returns the underlying file.
func (t tchain) getFile() *File {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].getFile()
}

// loadEntry returns an error if there is a problem during the loading.
func (t tchain) loadEntry(i int64) error {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].loadEntry(i)
}

var (
	_ Object = (*tchain)(nil)
	_ Named  = (*tchain)(nil)
	_ Tree   = (*tchain)(nil)
)
