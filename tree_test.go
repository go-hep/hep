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
	if got, want := key.Name(), "tree"; got != want {
		t.Fatalf("key.Name: got=%q. want=%q", got, want)
	}

	tree := key.Value().(*Tree)
	if got, want := tree.Name(), "tree"; got != want {
		t.Fatalf("tree.Name: got=%q. want=%q", got, want)
	}
	myprintf(">>> f.Get(tree)... [done]\n")

	for _, table := range []struct {
		test  string
		value string
		want  string
	}{
		{"Name", tree.Name(), "tree"}, // name when created
		{"Title", tree.Title(), "my tree title"},
		{"Class", tree.Class(), "TTree"},
	} {
		if table.value != table.want {
			t.Fatalf("%v: got=[%v]. want=[%v]", table.test, table.value, table.want)
		}
	}

	entries := tree.Entries()
	if got, want := entries, int64(100); got != want {
		t.Fatalf("tree.Entries: got=%v. want=%v", got, want)
	}

	if got, want := tree.totbytes, int64(40506); got != want {
		t.Fatalf("tree.totbytes: got=%v. want=%v", got, want)
	}

	if got, want := tree.zipbytes, int64(4184); got != want {
		t.Fatalf("tree.zipbytes: got=%v. want=%v", got, want)
	}
}
