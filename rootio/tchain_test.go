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

	if etree, echain := tree.Entries(), chain.Entries(); etree != echain {
		t.Fatalf("entries differ. got=%v, want=%v", echain, etree)
	}
}

/*
func TestChainEmpty(t *testing.T) {
chain := rootio.Chain()
}
*/
