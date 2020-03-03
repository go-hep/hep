// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rcmd"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func TestDiff(t *testing.T) {
	tmp, err := ioutil.TempDir("", "groot-rcmd-diff-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, tc := range []struct {
		name string
		keys []string
		err  error
		want string
		fref func(name string) *riofs.File
		fchk func(name string) *riofs.File
	}{
		{
			name: "same",
			fref: func(name string) *riofs.File {
				f, err := groot.Open("../testdata/small-flat-tree.root")
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Open("../testdata/small-flat-tree.root")
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
		},
		{
			name: "empty",
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				return f
			},
			err: fmt.Errorf("could not compute keys to compare: empty key set"),
		},
		{
			name: "empty-with-keys",
			keys: []string{" "},
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				return f
			},
			err: fmt.Errorf("could not compute keys to compare: empty key set"),
		},
		{
			name: "only-dirs",
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				_, err = riofs.Dir(f).Mkdir("dir-1/dir-11/dir-111")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				_, err = riofs.Dir(f).Mkdir("dir-1/dir-11/dir-111")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				return f
			},
		},
		{
			name: "different-key-type",
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = dir.Put("k1", rbase.NewNamed("k1-name", "k1-title"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = dir.Put("k1", rbase.NewObjString("obj-string"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			err: fmt.Errorf("dir-1: values for dir-11 in directory differ: dir-1/dir-11: values for k1 in directory differ: dir-1/dir-11/k1: type of keys differ: ref=*rbase.Named chk=*rbase.ObjString"),
		},
		{
			name: "different-key-set-chk",
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir11, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}
				_ = dir21

				err = dir11.Put("k1", rbase.NewObjString("obj-string"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir11, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir31, err := riofs.Dir(f).Mkdir("dir-3/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}
				_ = dir31

				err = dir11.Put("k1", rbase.NewObjString("obj-string-xxx"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			err:  fmt.Errorf("could not compute keys to compare: key set differ"),
			want: "key[dir-2] -- missing from chk-file\n",
		},
		{
			name: "different-key-set-ref",
			keys: []string{"dir-1", "dir-3"},
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir11, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}
				_ = dir21

				err = dir11.Put("k1", rbase.NewObjString("obj-string"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir11, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir31, err := riofs.Dir(f).Mkdir("dir-3/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}
				_ = dir31

				err = dir11.Put("k1", rbase.NewObjString("obj-string-xxx"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			err:  fmt.Errorf("could not compute keys to compare: key set differ"),
			want: "key[dir-3] -- missing from ref-file\n",
		},
		{
			name: "different-key-value",
			keys: []string{"dir-1"},
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = dir.Put("k1", rbase.NewObjString("obj-string"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = dir.Put("k1", rbase.NewObjString("obj-string-xxx"))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			err: fmt.Errorf("dir-1: values for dir-11 in directory differ: dir-1/dir-11: values for k1 in directory differ: dir-1/dir-11/k1: keys differ"),
			want: `key[dir-1/dir-11/k1] (*rbase.ObjString) -- (-ref +chk)
-obj-string
+obj-string-xxx
`,
		},
		{
			name: "different-trees-entries",
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				var data struct {
					I32 int32
					F64 float64
					Arr [2]float64
				}
				w, err := rtree.NewWriter(dir, "tree", rtree.WriteVarsFromStruct(&data))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				for i := 0; i < 5; i++ {
					data.I32 = int32(i)
					data.F64 = float64(i)
					data.Arr = [2]float64{float64(i + 1), float64(i + 2)}
					_, err = w.Write()
					if err != nil {
						t.Fatalf("could not write event #%d: %+v", i, err)
					}
				}

				err = w.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				var data struct {
					I32 int32
					F64 float64
					Arr [2]float64
				}
				w, err := rtree.NewWriter(dir, "tree", rtree.WriteVarsFromStruct(&data))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				for i := 0; i < 6; i++ {
					data.I32 = int32(i)
					data.F64 = float64(i)
					data.Arr = [2]float64{float64(i + 1), float64(i + 2)}
					_, err = w.Write()
					if err != nil {
						t.Fatalf("could not write event #%d: %+v", i, err)
					}
				}

				err = w.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			err: fmt.Errorf("dir-1: values for dir-11 in directory differ: dir-1/dir-11: values for tree in directory differ: dir-1/dir-11/tree: number of entries differ: ref=5 chk=6"),
		},
		{
			name: "different-trees-values",
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				var data struct {
					I32 int32
					F64 float64
					Arr [2]float64
				}
				w, err := rtree.NewWriter(dir, "tree", rtree.WriteVarsFromStruct(&data))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				for i := 0; i < 5; i++ {
					data.I32 = int32(i)
					data.F64 = float64(i + 1)
					data.Arr = [2]float64{float64(i + 1), float64(i + 2)}
					_, err = w.Write()
					if err != nil {
						t.Fatalf("could not write event #%d: %+v", i, err)
					}
				}

				err = w.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				var data struct {
					I32 int32
					F64 float64
					Arr [2]float64
				}
				w, err := rtree.NewWriter(dir, "tree", rtree.WriteVarsFromStruct(&data))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				for i := 0; i < 5; i++ {
					data.I32 = int32(i)
					data.F64 = float64(i)
					data.Arr = [2]float64{float64(i), float64(i + 2)}
					_, err = w.Write()
					if err != nil {
						t.Fatalf("could not write event #%d: %+v", i, err)
					}
				}

				err = w.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			err: fmt.Errorf("dir-1: values for dir-11 in directory differ: dir-1/dir-11: values for tree in directory differ: dir-1/dir-11/tree: trees differ"),
			want: `key[dir-1/dir-11/tree][0000].F64 -- (-ref +chk)
  float64(
- 	1,
+ 	0,
  )
key[dir-1/dir-11/tree][0000].Arr -- (-ref +chk)
  [2]float64{
- 	1,
+ 	0,
  	2,
  }
key[dir-1/dir-11/tree][0001].F64 -- (-ref +chk)
  float64(
- 	2,
+ 	1,
  )
key[dir-1/dir-11/tree][0001].Arr -- (-ref +chk)
  [2]float64{
- 	2,
+ 	1,
  	3,
  }
key[dir-1/dir-11/tree][0002].F64 -- (-ref +chk)
  float64(
- 	3,
+ 	2,
  )
key[dir-1/dir-11/tree][0002].Arr -- (-ref +chk)
  [2]float64{
- 	3,
+ 	2,
  	4,
  }
key[dir-1/dir-11/tree][0003].F64 -- (-ref +chk)
  float64(
- 	4,
+ 	3,
  )
key[dir-1/dir-11/tree][0003].Arr -- (-ref +chk)
  [2]float64{
- 	4,
+ 	3,
  	5,
  }
key[dir-1/dir-11/tree][0004].F64 -- (-ref +chk)
  float64(
- 	5,
+ 	4,
  )
key[dir-1/dir-11/tree][0004].Arr -- (-ref +chk)
  [2]float64{
- 	5,
+ 	4,
  	6,
  }
`,
		},
		{
			name: "different-trees-types",
			fref: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				var data struct {
					I32 int64
					F64 float64
					Arr [2]float64
				}
				w, err := rtree.NewWriter(dir, "tree", rtree.WriteVarsFromStruct(&data))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				for i := 0; i < 5; i++ {
					data.I32 = int64(i)
					data.F64 = float64(i + 1)
					data.Arr = [2]float64{float64(i + 1), float64(i + 2)}
					_, err = w.Write()
					if err != nil {
						t.Fatalf("could not write event #%d: %+v", i, err)
					}
				}

				err = w.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			fchk: func(name string) *riofs.File {
				f, err := groot.Create(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
				if err != nil {
					t.Fatalf("%+v", err)
				}

				var data struct {
					I32 int32
					F64 float64
					Arr [2]float64
				}
				w, err := rtree.NewWriter(dir, "tree", rtree.WriteVarsFromStruct(&data))
				if err != nil {
					t.Fatalf("%+v", err)
				}

				for i := 0; i < 5; i++ {
					data.I32 = int32(i)
					data.F64 = float64(i + 1)
					data.Arr = [2]float64{float64(i + 1), float64(i + 2)}
					_, err = w.Write()
					if err != nil {
						t.Fatalf("could not write event #%d: %+v", i, err)
					}
				}

				err = w.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				f, err = groot.Open(name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				return f
			},
			err: fmt.Errorf("dir-1: values for dir-11 in directory differ: dir-1/dir-11: values for tree in directory differ: dir-1/dir-11/tree: trees differ"),
			want: `key[dir-1/dir-11/tree][0000].I32 -- (-ref +chk)
  interface{}(
- 	int64(0),
+ 	int32(0),
  )
key[dir-1/dir-11/tree][0001].I32 -- (-ref +chk)
  interface{}(
- 	int64(1),
+ 	int32(1),
  )
key[dir-1/dir-11/tree][0002].I32 -- (-ref +chk)
  interface{}(
- 	int64(2),
+ 	int32(2),
  )
key[dir-1/dir-11/tree][0003].I32 -- (-ref +chk)
  interface{}(
- 	int64(3),
+ 	int32(3),
  )
key[dir-1/dir-11/tree][0004].I32 -- (-ref +chk)
  interface{}(
- 	int64(4),
+ 	int32(4),
  )
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			refname := filepath.Join(tmp, tc.name+"-ref.root")
			fref := tc.fref(refname)
			defer fref.Close()

			chkname := filepath.Join(tmp, tc.name+"-chk.root")
			fchk := tc.fchk(chkname)
			defer fchk.Close()

			out := new(strings.Builder)
			err := rcmd.Diff(out, fref, fchk, tc.keys)
			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error.\ngot= %s\nwant=%s\n", got, want)
				}
			case err != nil && tc.err == nil:
				t.Fatalf("unexpected error: %+v", err)

			case err == nil && tc.err != nil:
				t.Fatalf("expected an error: %+v", tc.err)

			case err == nil && tc.err == nil:
				// ok
				return
			}

			// replace non-breaking spaces (U+00a0) with regular space (U+0020).
			got := strings.Replace(out.String(), " ", " ", -1)

			if got, want := got, tc.want; got != want {
				t.Fatalf("invalid diff.\ngot:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}
