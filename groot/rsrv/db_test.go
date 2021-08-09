// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsrv

import (
	"os"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
)

func TestDB(t *testing.T) {
	dir, err := os.MkdirTemp("", "groot-rsrv-db-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	os.RemoveAll(dir)

	db := NewDB(dir)
	if got, want := len(db.Files()), 0; got != want {
		t.Fatalf("invalid number of files. got=%d, want=%d", got, want)
	}

	f, err := riofs.Open("../testdata/simple.root")
	if err != nil {
		t.Fatalf("could not open ROOT file: %v", err)
	}
	defer f.Close()

	const uri = "upload-store:///simple.root"
	db.set(uri, f)

	if got, want := db.Files(), []string{uri}; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid list of files. got=%v, want=%v", got, want)
	}

	wantFname := f.Name()
	err = db.Tx(uri, func(f *riofs.File) error {
		got := f.Name()
		if got != wantFname {
			t.Fatalf("invalid filename in transaction. got=%q, want=%q", got, wantFname)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}

	err = db.Tx("not-there", func(f *riofs.File) error { return nil })
	if err == nil {
		t.Fatalf("expected an error")
	}

	db.Close()
}
