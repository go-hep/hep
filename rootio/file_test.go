// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"testing"
)

func TestFileDirectory(t *testing.T) {
	f, err := Open("testdata/small-flat-tree.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	for _, table := range []struct {
		test  string
		value string
		want  string
	}{
		{"Name", f.Name(), "test-small.root"}, // name when created
		{"Title", f.Title(), "small event file"},
		{"Class", f.Class(), "TFile"},
	} {
		if table.value != table.want {
			t.Fatalf("%v: got=%q, want=%q", table.test, table.value, table.want)
		}
	}

	for _, table := range []struct {
		name string
		want bool
	}{
		{"tree", true},
		{"tree;0", false},
		{"tree;1", true},
		{"tree;9999", true},
		{"tree_nope", false},
		{"tree_nope;0", false},
		{"tree_nope;1", false},
		{"tree_nope;9999", false},
	} {
		_, err := f.Get(table.name)
		if (err == nil) != table.want {
			t.Fatalf("%s: got key (err=%v). want=%v", table.name, err, table.want)
		}
	}

	for _, table := range []struct {
		name string
		want string
	}{
		{"tree", "TTree"},
		{"tree;1", "TTree"},
	} {
		k, err := f.Get(table.name)
		if err != nil {
			t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
		}

		if k.Class() != table.want {
			t.Fatalf("%s: got key with class=%s (want=%s)", table.name, k.Class(), table.want)
		}
	}

	for _, table := range []struct {
		name string
		want string
	}{
		{"tree", "tree"},
		{"tree;1", "tree"},
	} {
		o, err := f.Get(table.name)
		if err != nil {
			t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
		}

		k := o.(Named)
		if k.Name() != table.want {
			t.Fatalf("%s: got key with name=%s (want=%v)", table.name, k.Name(), table.want)
		}
	}

	for _, table := range []struct {
		name string
		want string
	}{
		{"tree", "my tree title"},
		{"tree;1", "my tree title"},
	} {
		o, err := f.Get(table.name)
		if err != nil {
			t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
		}

		k := o.(Named)
		if k.Title() != table.want {
			t.Fatalf("%s: got key with title=%s (want=%v)", table.name, k.Title(), table.want)
		}
	}
}

func TestFileOpenStreamerInfo(t *testing.T) {
	for _, fname := range []string{
		"testdata/small-flat-tree.root",
		"testdata/simple.root",
	} {
		f, err := Open(fname)
		if err != nil {
			t.Errorf("error opening %q: %v\n", fname, err)
			continue
		}
		defer f.Close()

		_ = f.StreamerInfo()
	}
}
