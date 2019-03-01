// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	regen = flag.Bool("regen", false, "regenerate reference files")
)

func TestGenerate(t *testing.T) {
	dir, err := ioutil.TempDir("", "groot-gen-type-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, tc := range []struct {
		fname     string
		want      string
		types     []string
		verbose   bool
		streamers bool
	}{
		{
			fname:     "../../testdata/small-evnt-tree-fullsplit.root",
			want:      "testdata/small-evnt-tree-fullsplit.txt",
			types:     []string{"Event", "P3"},
			streamers: true,
		},
	} {
		t.Run(tc.fname, func(t *testing.T) {
			oname := filepath.Base(tc.fname) + ".go"
			o, err := os.Create(filepath.Join(dir, oname))
			if err != nil {
				t.Fatal(err)
			}
			defer o.Close()

			err = generate(o, "main", tc.types, tc.fname, tc.verbose, tc.streamers)
			if err != nil {
				t.Fatalf("could not generate types: %v", err)
			}

			err = o.Close()
			if err != nil {
				t.Fatal(err)
			}

			got, err := ioutil.ReadFile(o.Name())
			if err != nil {
				t.Fatalf("could not read generated file: %v", err)
			}

			if *regen {
				ioutil.WriteFile(tc.want, got, 0644)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference file: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("error:\n%v", diff(t, string(got), string(want)))
			}
		})
	}
}

func TestRW(t *testing.T) {
	dir, err := ioutil.TempDir("", "groot-gen-type-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, tc := range []struct {
		fname     string
		want      string
		types     []string
		verbose   bool
		streamers bool
		main      string
	}{
		{
			fname:     "../../testdata/small-evnt-tree-fullsplit.root",
			want:      "testdata/small-evnt-tree-fullsplit.txt",
			types:     []string{"Event", "P3"},
			streamers: true,
			main: `
package main

import (
	"flag"
	"log"
	"reflect"

	"go-hep.org/x/hep/groot"
)

func main() {
	{ // FIXME(sbinet): this shouldn't be necessary => bundle streamerinfo
		fname := flag.String("f", "", "file")
		flag.Parse()
		f, err := groot.Open(*fname)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
	w, err := groot.Create("out.root")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	wevt := &Event{
		Beg:       "beg",
		I16:       -16,
		I32:       -32,
		I64:       -64,
		U16:       +16,
		U32:       +32,
		U64:       +64,
		F32:       +32,
		F64:       +64,
		Str:       "my-string",
		P3:        P3{1, 2, 3},
		ArrayI16:  [10]int16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayI32:  [10]int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayI64:  [10]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayU16:  [10]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayU32:  [10]uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayU64:  [10]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayF32:  [10]float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayF64:  [10]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		N:         10,
		SliceI16:  []int16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceI32:  []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceI64:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceU16:  []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceU32:  []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceU64:  []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceF32:  []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceF64:  []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StdStr:    "std-string",
		StlVecI16: []int16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecI32: []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecI64: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecU16: []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecU32: []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecU64: []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecF32: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecF64: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		StlVecStr: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		End:       "end",
	}

	err = w.Put("evt", wevt)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatalf("error closing out.root file: %v", err)
	}

	r, err := groot.Open("out.root")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	o, err := r.Get("evt")
	if err != nil {
		log.Fatal(err)
	}

	revt := o.(*Event)
	if !reflect.DeepEqual(revt, wevt) {
		log.Fatalf("error:\ngot= %#v\nwant=%#v", revt, wevt)
	}
}
`,
		},
	} {
		t.Run(tc.fname, func(t *testing.T) {
			oname := filepath.Base(tc.fname) + ".go"
			o, err := os.Create(filepath.Join(dir, oname))
			if err != nil {
				t.Fatal(err)
			}
			defer o.Close()

			err = generate(o, "main", tc.types, tc.fname, tc.verbose, tc.streamers)
			if err != nil {
				t.Fatalf("could not generate types: %v", err)
			}

			err = o.Close()
			if err != nil {
				t.Fatal(err)
			}

			got, err := ioutil.ReadFile(o.Name())
			if err != nil {
				t.Fatalf("could not read generated file: %v", err)
			}

			if *regen {
				ioutil.WriteFile(tc.want, got, 0644)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference file: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("error:\n%v", diff(t, string(got), string(want)))
			}

			err = ioutil.WriteFile(filepath.Join(dir, "main.go"), []byte(tc.main), 0644)
			if err != nil {
				t.Fatal(err)
			}

			cwd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			buf := new(bytes.Buffer)
			cmd := exec.Command("go", "build",
				"-o", filepath.Join(dir, "a.out"),
				filepath.Join(dir, "main.go"),
				filepath.Join(dir, oname),
			)
			cmd.Stdout = buf
			cmd.Stderr = buf
			err = cmd.Run()
			if err != nil {
				t.Fatalf("could not run command %v:\n%v\nerr=%v",
					cmd.Args,
					buf.String(), err)
			}
			buf.Reset()

			cmd = exec.Command("./a.out", "-f", filepath.Join(cwd, tc.fname))
			cmd.Dir = dir
			cmd.Stdout = buf
			cmd.Stderr = buf
			err = cmd.Run()
			if err != nil {
				t.Fatalf("could not run command %v:\n%v\nerr=%v",
					cmd.Args,
					buf.String(), err)
			}
		})
	}
}

func diff(t *testing.T, chk, ref string) string {
	t.Helper()

	if !hasDiffCmd {
		return fmt.Sprintf("=== got ===\n%s\n=== want ===\n%s\n", chk, ref)
	}

	tmpdir, err := ioutil.TempDir("", "groot-diff-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	got := filepath.Join(tmpdir, "got.txt")
	err = ioutil.WriteFile(got, []byte(chk), 0644)
	if err != nil {
		t.Fatalf("could not create %s file: %v", got, err)
	}

	want := filepath.Join(tmpdir, "want.txt")
	err = ioutil.WriteFile(want, []byte(ref), 0644)
	if err != nil {
		t.Fatalf("could not create %s file: %v", want, err)
	}

	out := new(bytes.Buffer)
	cmd := exec.Command("diff", "-urN", want, got)
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return out.String() + "\nerror: " + err.Error()
}

var hasDiffCmd = false

func init() {
	_, err := exec.LookPath("diff")
	if err == nil {
		hasDiffCmd = true
	}
}
