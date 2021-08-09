// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rcmd"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
)

func TestROOTCp(t *testing.T) {
	dir, err := os.MkdirTemp("", "groot-root-cp-")
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
	{
		subdir, err := riofs.Dir(ref).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		obj := rbase.NewObjString("string111")
		err = subdir.Put("str-111", obj)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		keys = append(keys, "dir-1/dir-11/str-111")
		refs = append(refs, obj)
	}
	{
		subdir, err := riofs.Dir(ref).Mkdir("dir-1/dir-12")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		obj := rbase.NewObjString("string121")
		err = subdir.Put("str-121", obj)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		keys = append(keys, "dir-1/dir-12/str-121")
		refs = append(refs, obj)
	}
	{
		subdir, err := riofs.Dir(ref).Mkdir("dir-2")
		if err != nil {
			t.Fatalf("%+v", err)
		}

		obj := rbase.NewObjString("string21")
		err = subdir.Put("obj-21", obj)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		keys = append(keys, "dir-2/obj-21")
		refs = append(refs, obj)
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
			keys:  []int{0, 1, 2, 3, 4, 5},
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
			oname: "out-str10.root",
			fname: refname + ":^/str",
			keys:  []int{2},
		},
		{
			oname: "out-str11.root",
			fname: refname + ":/str",
			keys:  []int{2, 3, 4},
		},
		{
			oname: "out-str12.root",
			fname: refname + ":str",
			keys:  []int{2, 3, 4},
		},
		{
			oname: "out-str20.root",
			fname: refname + ":^/str.*",
			keys:  []int{2},
		},
		{
			oname: "out-str21.root",
			fname: refname + ":/str.*",
			keys:  []int{2, 3, 4},
		},
		{
			oname: "out-str22.root",
			fname: refname + ":str.*",
			keys:  []int{2, 3, 4},
		},
		{
			oname: "out-dir.root",
			fname: refname + ":dir",
			keys:  []int{3, 4, 5},
		},
		{
			oname: "empty.root",
			fname: refname + ":NONE.*",
			keys:  []int{},
		},
	} {
		t.Run(tc.oname, func(t *testing.T) {
			oname := filepath.Join(dir, tc.oname)
			err := rcmd.Copy(oname, []string{tc.fname})
			if err != nil {
				t.Fatalf("%+v", err)
			}

			f, err := groot.Open(oname)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			defer f.Close()

			gotKeys := 0
			err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
				if err != nil {
					return err
				}
				name := path[len(f.Name()):]
				if name == "" {
					return nil
				}
				if _, isdir := obj.(riofs.Directory); isdir {
					return nil
				}
				gotKeys++
				return nil
			})
			if err != nil {
				t.Fatalf("could not count keys in output ROOT file: %+v", err)
			}

			if got, want := gotKeys, len(tc.keys); got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			for _, i := range tc.keys {
				v, err := riofs.Dir(f).Get(keys[i])
				if err != nil {
					t.Fatalf("%+v", err)
				}

				switch v := v.(type) {
				case riofs.Directory:
					if got, want := v.(root.Named).Name(), refs[i].(root.Named).Name(); got != want {
						t.Fatalf(
							"invalid value for %q:\ngot=%v\nwant=%v\n",
							keys[i],
							got, want,
						)
					}
				default:
					if !reflect.DeepEqual(v, refs[i]) {
						t.Fatalf(
							"invalid value for %q:\ngot=%v\nwant=%v\n",
							keys[i],
							v, refs[i],
						)
					}
				}
			}

		})
	}
}

func TestROOTCpTree(t *testing.T) {
	dir, err := os.MkdirTemp("", "groot-root-cp-")
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

	refdir, err := riofs.Dir(ref).Mkdir("dir1/dir11")
	if err != nil {
		t.Fatalf("could not create dir hierarchy: %+v", err)
	}

	rsrc, err := rtree.NewWriter(refdir, "mytree", rtree.WriteVarsFromStruct(&rdata))
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
	err = rcmd.Copy(chkname, []string{refname + ":dir1/dir11/mytree"})
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
