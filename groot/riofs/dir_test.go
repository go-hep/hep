// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func TestDirs(t *testing.T) {
	rootdir, err := ioutil.TempDir("", "groot-dir-subdir-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(rootdir)

	fname := filepath.Join(rootdir, "subdirs.root")

	{
		w, err := groot.Create(fname)
		if err != nil {
			t.Fatal(err)
		}
		defer w.Close()

		dir1, err := w.Mkdir("dir1")
		if err != nil {
			t.Fatal(err)
		}

		dir11, err := dir1.Mkdir("dir11")
		if err != nil {
			t.Fatal(err)
		}

		err = dir11.Put("obj1", rbase.NewObjString("data-obj1"))
		if err != nil {
			t.Fatal(err)
		}

		dir2, err := w.Mkdir("dir2")
		if err != nil {
			t.Fatal(err)
		}

		err = dir2.Put("obj2", rbase.NewObjString("data-obj2"))
		if err != nil {
			t.Fatal(err)
		}

		err = w.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		r, err := groot.Open(fname)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		var obj root.Object

		obj, err = r.Get("dir1")
		if err != nil {
			t.Fatal(err)
		}

		dir1 := obj.(riofs.Directory)

		obj, err = r.Get("dir2")
		if err != nil {
			t.Fatal(err)
		}

		dir2 := obj.(riofs.Directory)

		obj, err = dir1.Get("dir11")
		if err != nil {
			t.Fatal(err)
		}

		dir11 := obj.(riofs.Directory)

		obj, err = dir11.Get("obj1")
		if err != nil {
			t.Fatal(err)
		}

		str1 := obj.(*rbase.ObjString)
		if got, want := str1.String(), "data-obj1"; got != want {
			t.Fatalf("got=%q, want=%q", got, want)
		}

		obj, err = dir2.Get("obj2")
		if err != nil {
			t.Fatal(err)
		}

		str2 := obj.(*rbase.ObjString)
		if got, want := str2.String(), "data-obj2"; got != want {
			t.Fatalf("got=%q, want=%q", got, want)
		}
	}
}
