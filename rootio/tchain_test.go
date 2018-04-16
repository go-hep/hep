// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio_test

import (
	"testing"

	"go-hep.org/x/hep/rootio"
)

func TestChain(t *testing.T) {
	f, err := rootio.Open("testdata/chain.1.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}
	tree := obj.(rootio.Tree)
	chain := rootio.Chain(tree)

	if entry_tree, entry_chain := tree.Entries(), chain.Entries(); entry_tree != entry_chain {
		t.Fatalf("entries differ: got=%v, want=%v\n", entry_chain, entry_tree)
	}
	if name_tree, name_chain := tree.Name(), chain.Name(); name_tree != name_chain {
		t.Fatalf("names differ : got=%v, want=%v\n", name_chain, name_tree)
	}
	if title_tree, title_chain := tree.Title(), chain.Title(); title_tree != title_chain {
		t.Fatalf("titles differ : got=%v, want=%v\n", title_chain, title_tree)
	}
}

func TestChainEmpty(t *testing.T) {
	chain := rootio.Chain()
	if chain.Entries() != 0 {
		t.Fatalf("unexpected entry : got=%v,want=0\n", chain.Entries())
	}
	if chain.Name() != "" {
		t.Fatalf("unexpected name : got=%v, want empty string\n", chain.Name())
	}
	if chain.Title() != "" {
		t.Fatalf("unexpected Title : got=%v, want empty string\n", chain.Title())
	}
}

func TestTwoChains(t *testing.T) {
	f1, err := rootio.Open("testdata/chain.1.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f1.Close()
	f2, err := rootio.Open("testdata/chain.2.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	obj1, err1 := f1.Get("tree")
	if err != nil {
		t.Fatal(err1)
	}
	obj2, err2 := f2.Get("tree")
	if err != nil {
		t.Fatal(err2)
	}

	tree1 := obj1.(rootio.Tree)
	tree2 := obj2.(rootio.Tree)
	chain := rootio.Chain(tree1, tree2)

	if entry_tree1, entry_tree2, entry_chain := tree1.Entries(), tree2.Entries(), chain.Entries(); entry_tree1+entry_tree2 != entry_chain {
		t.Fatalf("entries differ: got=%v, want=%v\n", entry_chain, entry_tree1+entry_tree2)
	}
	if name_tree1, name_chain := tree1.Name(), chain.Name(); name_tree1 != name_chain {
		t.Fatalf("names differ : got=%v, want=%v\n", name_chain, name_tree1)
	}

}
