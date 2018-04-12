// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.



package rootio


type tchain struct {
	trees []Tree
}

//Class returns the ROOT class of the argument.
func (tchain) Class() string {
	return "TChain"
}



//Name returns the name of the ROOT objet in the argument.
func (t tchain) Name() string {
	return t.trees[0].Name()
}


//Title returns the title of the ROOT object in the argument
func (t tchain) Title() string {
	return t.trees[0].Title()
}

// Chain returns a tchain that is the concatenation of all the input Trees.
func Chain(trees ...Tree) tchain {
	var t tchain
	t.trees = append(t.trees, trees...)
	return t
}

//Entries returns the total number of entries 
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
	for _,tree:= range t.trees {
		v += tree.TotBytes()
	}
	return v
}

//ZipBytes returns the total number of bytes after compression.
func (t tchain) ZipBytes() int64 {
	var v int64 
	for _,tree:= range t.trees {
		v += tree.ZipBytes()
	}
	return v

}

//Branches returns the list of branches.
func (t tchain) Branches() []Branch {
	branch := t.trees[0].Branches()
	if branch != nil {
		return branch
	}
	return nil	
}

//Branch returns the branch whose name is the argument.
func (t tchain) Branch(name string) Branch {
	branch := t.trees[0].Branches()
	if branch != nil {	
		for _, br := range t.trees[0].Branches() {
			if br.Name() == name {
				return br
			}
		}
		return nil
	}
	return nil
}

//Leaves returns direct pointers to individual branch leaves.
func (t tchain) Leaves() []Leaf {
	leaf := t.trees[0].Leaves()
	if leaf != nil {
		return leaf
	}
	return nil
}

//getFile returns the underlying file.
func (t tchain) getFile() *File {
	f := t.trees[0].getFile()
	if f != nil {
		return f
	}
	return nil
}

//loadEntry returns an error if there is a problem during the loading
func (t tchain) loadEntry(i int64) error {
	branch := t.trees[0].Branches()
	if branch != nil {
		for _, b := range t.trees[0].Branches() {
			err := b.loadEntry(i)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

var (
	_ Object = (*tchain)(nil)
	_ Named = (*tchain)(nil)
	_ Tree = (*tchain)(nil)
)


