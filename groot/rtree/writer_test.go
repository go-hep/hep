// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

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
	tmp, err := ioutil.TempDir("", "groot-rtree-")
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
