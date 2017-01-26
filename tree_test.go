// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "testing"

func TestFlatTree(t *testing.T) {
	f, err := Open("testdata/small.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	myprintf(">>> f.Get(tree)...\n")
	obj, ok := f.Get("tree")
	if !ok {
		t.Fatalf("could not retrieve tree [tree]")
	}

	key := obj.(*Key)
	if key.Name() != "tree" {
		t.Fatalf("key.Name: expected [tree] (got=%v)", key.Name())
	}

	tree := key.Value().(*Tree)
	if tree.Name() != "tree" {
		t.Fatalf("tree.Name: expected [tree] (got=%v)", tree.Name())
	}
	myprintf(">>> f.Get(tree)... [done]\n")

	for _, table := range []struct {
		test     string
		value    string
		expected string
	}{
		{"Name", tree.Name(), "tree"}, // name when created
		{"Title", tree.Title(), "my tree title"},
		{"Class", tree.Class(), "TTree"},
	} {
		if table.value != table.expected {
			t.Fatalf("%v: expected [%v] got [%v]", table.test, table.expected, table.value)
		}
	}

	entries := tree.Entries()
	if entries != 100 {
		t.Fatalf("tree.Entries: expected [100] (got=%v)", entries)
	}

	if tree.totbytes != 40506 {
		t.Fatalf("tree.totbytes: expected [40506] (got=%v)", tree.totbytes)
	}

	if tree.zipbytes != 4184 {
		t.Fatalf("tree.zipbytes: expected [4184] (got=%v)", tree.zipbytes)
	}
}
