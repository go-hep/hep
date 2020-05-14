// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
)

func TestMyReader(t *testing.T) {
	//const fname = "_data-perf/f64s-2.root"
	const fname = "../testdata/small-flat-tree.root"
	f, err := riofs.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := o.(Tree)

	var data struct {
		N   int32
		Sli []float64
	}

	rvars := []ReadVar{
		{Name: "N", Value: &data.N},
		{Name: "SliceFloat64", Value: &data.Sli},
	}

	pr, err := NewReader(tree, rvars, WithRange(0, 10))
	if err != nil {
		t.Fatal(err)
	}
	err = pr.Read(func(rctx RCtx) error {
		log.Printf("[%03d][Sli]: %v\n", rctx.Entry, data.Sli)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPReader(t *testing.T) {
	const fname = "_data-perf/f64s-1024.root"
	//const fname = "../testdata/small-flat-tree.root"
	f, err := riofs.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := o.(Tree)

	var data struct {
		N   int32
		Sli []float64
	}

	rvars := []ReadVar{
		{Name: "N", Value: &data.N},
		{Name: "Sli", Value: &data.Sli},
	}

	//	oo, err := os.Create("new.dump.txt")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	defer oo.Close()

	var sum int
	pr, err := NewPReader(tree, rvars)
	if err != nil {
		t.Fatal(err)
	}
	//pr.nevt = 10
	err = pr.Read(func(rctx RCtx) error {
		sum += len(data.Sli)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

var sumBenchPReadTreeSliF64 = 0

func BenchmarkPReadTreeSliF64(b *testing.B) {
	tmp, err := ioutil.TempDir("", "groot-rtree-read-tree-sli-f64s-")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	const nevts = 1000

	for _, sz := range []int{0, 1, 2, 4, 8, 16, 64, 128, 512, 1024, 1024 * 1024} {
		fname := path.Join(tmp, fmt.Sprintf("f64s-%d.root", sz))
		func() {
			b.StopTimer()
			defer b.StartTimer()

			f, err := riofs.Create(fname, riofs.WithoutCompression())
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()

			var data struct {
				N   int32
				Sli []float64
			}
			wvars := []WriteVar{
				{Name: "N", Value: &data.N},
				{Name: "Sli", Value: &data.Sli, Count: "N"},
			}
			tree, err := NewWriter(f, "tree", wvars, WithoutCompression())
			if err != nil {
				b.Fatal(err)
			}
			defer tree.Close()

			rnd := rand.New(rand.NewSource(1234))
			for i := 0; i < nevts; i++ {
				data.N = int32(rnd.Float64() * 100)
				data.Sli = make([]float64, int(data.N))
				for j := range data.Sli {
					data.Sli[j] = rnd.Float64() * 10
				}

				_, err = tree.Write()
				if err != nil {
					b.Fatal(err)
				}
			}

			err = tree.Close()
			if err != nil {
				b.Fatal(err)
			}

			err = f.Close()
			if err != nil {
				b.Fatal(err)
			}
		}()

		b.Run(fmt.Sprintf("%d", sz), func(b *testing.B) {
			b.StopTimer()
			f, err := riofs.Open(fname)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()

			o, err := f.Get("tree")
			if err != nil {
				b.Fatal(err)
			}

			tree := o.(Tree)

			var data struct {
				N   int32
				Sli []float64
			}

			rvars := ReadVarsFromStruct(&data)
			r, err := NewPReader(tree, rvars)
			if err != nil {
				b.Fatal(err)
			}
			//defer r.close()

			b.StartTimer()
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r.reset()
				data.N = 0
				data.Sli = data.Sli[:0]
				b.StartTimer()

				err = r.Read(func(RCtx) error {
					sumBenchPReadTreeSliF64 += len(data.Sli)
					return nil
				})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
