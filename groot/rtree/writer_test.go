// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"compress/flate"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"go-hep.org/x/hep/groot/internal/rcompress"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/riofs"
	"golang.org/x/exp/rand"
)

func TestFlattenArrayType(t *testing.T) {
	for _, tc := range []struct {
		typ   interface{}
		want  interface{}
		shape []int
	}{
		{
			typ:   int32(0),
			want:  int32(0),
			shape: nil,
		},
		{
			typ:   [2]int32{},
			want:  int32(0),
			shape: []int{2},
		},
		{
			typ:   [2][3]int32{},
			want:  int32(0),
			shape: []int{2, 3},
		},
		{
			typ:   [2][3][4]int32{},
			want:  int32(0),
			shape: []int{2, 3, 4},
		},
		{
			typ:   [2][3][4][0]int32{},
			want:  int32(0),
			shape: []int{2, 3, 4, 0},
		},
		{
			typ:   [2][3][4][5]int32{},
			want:  int32(0),
			shape: []int{2, 3, 4, 5},
		},
		{
			typ:   [2][3][4][5][6]int32{},
			want:  int32(0),
			shape: []int{2, 3, 4, 5, 6},
		},
		{
			typ:   [2][3][4][0]struct{}{},
			want:  struct{}{},
			shape: []int{2, 3, 4, 0},
		},
		{
			typ:   [2][3][4][0][]string{},
			want:  []string{},
			shape: []int{2, 3, 4, 0},
		},
		{
			typ:   []string{},
			want:  []string{},
			shape: nil,
		},
	} {
		t.Run(fmt.Sprintf("%T", tc.typ), func(t *testing.T) {
			rt, shape := flattenArrayType(reflect.TypeOf(tc.typ))
			if got, want := rt, reflect.TypeOf(tc.want); got != want {
				t.Fatalf("invalid array element type: got=%v, want=%v", got, want)
			}

			if got, want := shape, tc.shape; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid shape: got=%v, want=%v", got, want)
			}
		})
	}
}

func TestInvalidTreeMerger(t *testing.T) {
	var (
		w   wtree
		src = rbase.NewObjString("foo")
	)

	err := w.ROOTMerge(src)
	if err == nil {
		t.Fatalf("expected an error")
	}

	const want = "rtree: can not merge src=*rbase.ObjString into dst=*rtree.wtree"
	if got, want := err.Error(), want; got != want {
		t.Fatalf("invalid ROOTMerge error. got=%q, want=%q", got, want)
	}
}

func TestConcurrentWrite(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-rtree-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmp)
	}()

	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			f, err := riofs.Create(filepath.Join(tmp, fmt.Sprintf("file-%03d.root", i)))
			if err != nil {
				t.Errorf("could not create root file: %+v", err)
				return
			}
			defer f.Close()

			var (
				evt struct {
					N   int32
					Sli []float64 `groot:"Sli[N]"`
				}
				wvars = WriteVarsFromStruct(&evt)
			)
			w, err := NewWriter(f, "tree", wvars)
			if err != nil {
				t.Errorf("could not create tree writer: %+v", err)
				return
			}
			defer w.Close()

			rng := rand.New(rand.NewSource(1234))
			for i := 0; i < 100; i++ {
				evt.N = rng.Int31n(10) + 1
				evt.Sli = evt.Sli[:0]
				for j := 0; j < int(evt.N); j++ {
					evt.Sli = append(evt.Sli, rng.Float64())
				}
				_, err = w.Write()
				if err != nil {
					t.Errorf("could not write event %d: %+v", i, err)
					return
				}
			}

			err = w.Close()
			if err != nil {
				t.Errorf("could not close tree writer: %+v", err)
				return
			}

			err = f.Close()
			if err != nil {
				t.Errorf("could not close root file: %+v", err)
				return
			}
		}(i)
	}
	wg.Wait()
}

func TestWriteThisStreamers(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-rtree-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, "streamers.root")
	o, err := riofs.Create(fname)
	if err != nil {
		t.Fatalf("could not create output ROOT file: %+v", err)
	}
	defer o.Close()

	var evt struct {
		F1 []float64 `groot:"F1"`
		F2 []float64 `groot:"F2"`
		F3 []int64   `groot:"F3"`
		F4 []int64   `groot:"F4"`
	}

	wvars := WriteVarsFromStruct(&evt)
	tree, err := NewWriter(o, "tree", wvars)
	if err != nil {
		t.Fatalf("could not create output ROOT tree %q: %+v", "tree", err)
	}

	for i := 0; i < 10; i++ {
		evt.F1 = []float64{float64(i)}
		evt.F2 = []float64{float64(i), float64(i)}
		evt.F3 = []int64{int64(i)}
		evt.F4 = []int64{int64(i), int64(i)}

		_, err = tree.Write()
		if err != nil {
			t.Fatalf("could not write event %d: %+v", i, err)
		}
	}

	err = tree.Close()
	if err != nil {
		t.Fatalf("could not close tree: %+v", err)
	}

	err = o.Close()
	if err != nil {
		t.Fatalf("could not close file: %+v", err)
	}

	f, err := riofs.Open(fname)
	if err != nil {
		t.Fatalf("could not re-open ROOT file: %+v", err)
	}
	defer f.Close()

	sinfos := make(map[string]int)
	for _, si := range f.StreamerInfos() {
		sinfos[si.Name()]++
	}

	for _, tc := range []struct {
		name string
		want int
	}{
		{"vector<double>", 1},
		{"vector<int64_t>", 1},
	} {
		got := sinfos[tc.name]
		if got != tc.want {
			t.Errorf("invalid count for %q: got=%d, want=%d", tc.name, got, tc.want)
		}
	}
}

func TestWriterWithCompression(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-rtree-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmp)
	}()

	for _, tc := range []struct {
		wopt WriteOption
		want rcompress.Kind
	}{
		// {WithoutCompression(flate.BestCompression), rcompress.},
		{WithLZ4(flate.BestCompression), rcompress.LZ4},
		{WithLZMA(flate.BestCompression), rcompress.LZMA},
		{WithZlib(flate.BestCompression), rcompress.ZLIB},
		{WithZstd(flate.BestCompression), rcompress.ZSTD},
	} {
		t.Run("alg-"+tc.want.String(), func(t *testing.T) {
			fname := filepath.Join(tmp, "groot-alg-"+tc.want.String())
			f, err := riofs.Create(fname)
			if err != nil {
				t.Fatalf("could not create file %q: %v", fname, err)
			}
			defer f.Close()

			var (
				evt struct {
					N   int32
					Sli []float64 `groot:"Sli[N]"`
				}
				wvars = WriteVarsFromStruct(&evt)
			)
			w, err := NewWriter(f, "tree", wvars, tc.wopt)
			if err != nil {
				t.Errorf("could not create tree writer: %+v", err)
				return
			}
			defer w.Close()

			rng := rand.New(rand.NewSource(1234))
			for i := 0; i < 100; i++ {
				evt.N = rng.Int31n(10) + 1
				evt.Sli = evt.Sli[:0]
				for j := 0; j < int(evt.N); j++ {
					evt.Sli = append(evt.Sli, rng.Float64())
				}
				_, err = w.Write()
				if err != nil {
					t.Errorf("could not write event %d: %+v", i, err)
					return
				}
			}

			err = w.Close()
			if err != nil {
				t.Errorf("could not close tree writer: %+v", err)
				return
			}

			err = f.Close()
			if err != nil {
				t.Errorf("could not close root file: %+v", err)
				return
			}

			{
				f, err := riofs.Open(fname)
				if err != nil {
					t.Fatalf("could not open ROOT file %q: %v", fname, err)
				}
				defer f.Close()

				tree, err := riofs.Get[Tree](f, "tree")
				if err != nil {
					t.Fatalf("could not open tree: %v", err)
				}
				bname := "Sli"
				b := tree.Branch(bname)
				if b == nil {
					t.Fatalf("could not retrieve branch %q", bname)
				}
				bb, ok := b.(*tbranch)
				if !ok {
					t.Fatalf("unexpected type for branch %q: %T", bname, b)
				}
				xset := rcompress.SettingsFrom(int32(bb.compress))
				if got, want := xset.Alg, tc.want; got != want {
					t.Fatalf("invalid compression algorithm for branch %q: got=%v, want=%v", b.Name(), got, want)
				}
			}
		})
	}
}
