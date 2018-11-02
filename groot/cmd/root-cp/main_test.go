// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/root"
)

func TestROOTCp(t *testing.T) {
	dir, err := ioutil.TempDir("", "groot-root-cp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	refname := filepath.Join(dir, "ref.root")
	ref, err := groot.Create(refname)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Close()

	refs := []root.Object{
		rbase.NewObjString("string1"),
		rbase.NewObjString("string2"),
		rbase.NewObjString("string3"),
	}
	keys := []string{
		"key", "key-1", "str-3",
	}

	for i := range refs {
		err := ref.Put(keys[i], refs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	err = ref.Close()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		oname string
		fname string
		keys  []int
	}{
		{
			oname: "out-all.root",
			fname: refname,
			keys:  []int{0, 1, 2},
		},
		{
			oname: "out-key.root",
			fname: refname + ":key",
			keys:  []int{0, 1},
		},
		{
			oname: "out-key-star.root",
			fname: refname + ":key.*",
			keys:  []int{0, 1},
		},
		{
			oname: "out-key-star2.root",
			fname: refname + ":key-.*",
			keys:  []int{1},
		},
		{
			oname: "out-str.root",
			fname: refname + ":str",
			keys:  []int{2},
		},
		{
			oname: "out-str.root",
			fname: refname + ":str.*",
			keys:  []int{2},
		},
		{
			oname: "empty.root",
			fname: refname + ":NONE.*",
			keys:  []int{},
		},
	} {
		t.Run(tc.oname, func(t *testing.T) {
			oname := filepath.Join(dir, tc.oname)
			err := rootcp(oname, []string{tc.fname})
			if err != nil {
				t.Fatal(err)
			}

			f, err := groot.Open(oname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			if got, want := len(f.Keys()), len(tc.keys); got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			for _, i := range tc.keys {
				v, err := f.Get(keys[i])
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(v, refs[i]) {
					t.Fatalf(
						"invalid value for %q:\ngot=%v\nwant=%v\n",
						keys[i],
						v, refs[i],
					)
				}
			}

		})
	}
}

func TestSplitArg(t *testing.T) {
	for _, tc := range []struct {
		cmd   string
		fname string
		sel   string
		err   error
	}{
		{
			cmd:   "file.root",
			fname: "file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root",
			fname: "dir/sub/file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "/dir/sub/file.root",
			fname: "/dir/sub/file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "../dir/sub/file.root",
			fname: "../dir/sub/file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root:hist",
			fname: "dir/sub/file.root",
			sel:   "hist",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root:hist*",
			fname: "dir/sub/file.root",
			sel:   "hist*",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root:",
			fname: "dir/sub/file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "file://dir/sub/file.root:",
			fname: "file://dir/sub/file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "https://dir/sub/file.root",
			fname: "https://dir/sub/file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "http://dir/sub/file.root",
			fname: "http://dir/sub/file.root",
			sel:   ".*",
			err:   nil,
		},
		{
			cmd:   "https://dir/sub/file.root:hist*",
			fname: "https://dir/sub/file.root",
			sel:   "hist*",
			err:   nil,
		},
		{
			cmd:   "root://dir/sub/file.root:hist*",
			fname: "root://dir/sub/file.root",
			sel:   "hist*",
			err:   nil,
		},
		{
			cmd: "dir/sub/file.root:h:h",
			err: errors.Errorf("root-cp: too many ':' in %q", "dir/sub/file.root:h:h"),
		},
		{
			cmd: "root://dir/sub/file.root:h:h",
			err: errors.Errorf("root-cp: too many ':' in %q", "root://dir/sub/file.root:h:h"),
		},
		{
			cmd: "root://dir/sub/file.root::h:",
			err: errors.Errorf("root-cp: too many ':' in %q", "root://dir/sub/file.root::h:"),
		},
	} {
		t.Run(tc.cmd, func(t *testing.T) {
			fname, sel, err := splitArg(tc.cmd)
			switch {
			case err != nil && tc.err != nil:
				if !reflect.DeepEqual(err.Error(), tc.err.Error()) {
					t.Fatalf("got err=%v, want=%v", err, tc.err)
				}
				return
			case err != nil && tc.err == nil:
				t.Fatalf("got err=%v, want=%v", err, tc.err)
			case err == nil && tc.err != nil:
				t.Fatalf("got err=%v, want=%v", err, tc.err)
			}

			if got, want := fname, tc.fname; got != want {
				t.Fatalf("fname=%q, want=%q", got, want)
			}

			if got, want := sel, tc.sel; got != want {
				t.Fatalf("selection=%q, want=%q", got, want)
			}
		})
	}
}
