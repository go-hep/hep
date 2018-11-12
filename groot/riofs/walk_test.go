// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"reflect"
	"testing"
)

func TestRecDir(t *testing.T) {
	f, err := Open("../testdata/dirs-6.14.00.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rd := Dir(f)
	for _, tc := range []struct {
		path  string
		class string
	}{
		{"dir1/dir11/h1", "TH1F"},
		{"dir1/dir11/h1;1", "TH1F"},
		{"dir1/dir11/h1;9999", "TH1F"},
		{"dir1/dir11", "TDirectoryFile"},
		{"dir1/dir11;1", "TDirectoryFile"},
		{"dir1/dir11;9999", "TDirectoryFile"},
		{"dir1", "TDirectoryFile"},
		{"dir2", "TDirectoryFile"},
		{"dir3", "TDirectoryFile"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			o, err := rd.Get(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := o.Class(), tc.class; got != want {
				t.Fatalf("got=%q, want=%q", got, want)
			}
		})
	}

	keys := make([]string, len(rd.Keys()))
	for i, k := range rd.Keys() {
		keys[i] = k.Name()
	}

	if got, want := keys, []string{"dir1", "dir2", "dir3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid keys:\ngot = %v\nwant=%v\n", got, want)
	}
}
