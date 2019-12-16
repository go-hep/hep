// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/internal/rcmd"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
	"golang.org/x/xerrors"
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
		t.Fatalf("%+v", err)
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
			t.Fatalf("%+v", err)
		}
	}

	err = ref.Close()
	if err != nil {
		t.Fatalf("%+v", err)
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
				t.Fatalf("%+v", err)
			}
			defer f.Close()

			if got, want := len(f.Keys()), len(tc.keys); got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			for _, i := range tc.keys {
				v, err := f.Get(keys[i])
				if err != nil {
					t.Fatalf("%+v", err)
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
			err: xerrors.Errorf("root-cp: too many ':' in %q", "dir/sub/file.root:h:h"),
		},
		{
			cmd: "root://dir/sub/file.root:h:h",
			err: xerrors.Errorf("root-cp: too many ':' in %q", "root://dir/sub/file.root:h:h"),
		},
		{
			cmd: "root://dir/sub/file.root::h:",
			err: xerrors.Errorf("root-cp: too many ':' in %q", "root://dir/sub/file.root::h:"),
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

func TestROOTCpTree(t *testing.T) {
	dir, err := ioutil.TempDir("", "groot-root-cp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	refname := filepath.Join(dir, "ref.root")
	ref, err := groot.Create(refname)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer ref.Close()

	rdata := struct {
		N    int32
		I32s []int32 `groot:"i32s[N]"`
	}{}

	rsrc, err := rtree.NewWriter(ref, "mytree", rtree.WriteVarsFromStruct(&rdata))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		rdata.N = int32(i)
		rdata.I32s = make([]int32, i)
		for j := range rdata.I32s {
			rdata.I32s[j] = int32(i)
		}

		_, err = rsrc.Write()
		if err != nil {
			t.Fatalf("could not write event %d: %+v", i, err)
		}
	}

	err = rsrc.Close()
	if err != nil {
		t.Fatalf("could not close src tree: %+v", err)
	}

	err = ref.Close()
	if err != nil {
		t.Fatalf("could not close ref file: %+v", err)
	}

	chkname := filepath.Join(dir, "chk.root")
	err = rootcp(chkname, []string{refname + ":mytree"})
	if err != nil {
		t.Fatalf("could not copy tree: %+v", err)
	}

	want := new(bytes.Buffer)
	err = rcmd.Dump(want, refname, true, nil)
	if err != nil {
		t.Fatalf("could not dump ref file %q: %+v", refname, err)
	}

	got := new(bytes.Buffer)
	err = rcmd.Dump(got, chkname, true, nil)
	if err != nil {
		t.Fatalf("could not dump new file %q: %+v", chkname, err)
	}

	if got, want := got.String(), want.String(); got != want {
		t.Fatalf("dumps differ:\ngot:\n%s\n===\nwant:\n%s\n===\n", got, want)
	}
}
