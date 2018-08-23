// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRMemFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "rootio-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fname := filepath.Join(dir, "objstring.root")

	w, err := Create(fname)
	if err != nil {
		t.Fatal(err)
	}

	var (
		kname = "my-key"
		want  = NewObjString("Hello World from Go-HEP!")
	)

	err = w.Put(kname, want)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(w.Keys()), 1; got != want {
		t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("error closing file: %v", err)
	}

	raw, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	r, err := NewReader(&memFile{bytes.NewReader(raw)})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	rgot, err := r.Get(kname)
	if err != nil {
		t.Fatal(err)
	}

	if got := rgot.(ObjString); !reflect.DeepEqual(got, want) {
		t.Fatalf("error reading back objstring.\ngot = %#v\nwant = %#v", got, want)
	}

	err = r.Close()
	if err != nil {
		t.Fatalf("error closing file: %v", err)
	}
}
