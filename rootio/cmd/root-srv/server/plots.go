// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/rootio"
)

func walk(f rootio.Directory, path []string) (rootio.Object, error) {
	o, err := f.Get(path[0])
	if err != nil {
		return nil, err
	}
	if dir, ok := o.(rootio.Directory); ok {
		return walk(dir, path[1:])
	}
	return o, nil
}

func (srv *server) plotH1Handle(w http.ResponseWriter, r *http.Request) error {
	uri := r.URL.Path[len("/plot-h1/"):]
	var err error
	uri, err = url.PathUnescape(uri)
	if err != nil {
		return err
	}
	toks := strings.Split(uri, "/")
	fname := toks[0]

	db, err := srv.db(r)
	if err != nil {
		return err
	}
	db.RLock()
	defer db.RUnlock()

	f := db.get(fname)
	obj, err := walk(f, toks[1:])
	if err != nil {
		return fmt.Errorf("could not find %q in file %q: %v", filepath.Join(toks[1:]...), fname, err)
	}

	robj, ok := obj.(rootio.H1)
	if !ok {
		return fmt.Errorf("object %q could not be converted to hbook.H1D", toks[1])
	}
	h1d, err := rootcnv.H1D(robj)
	if err != nil {
		return err
	}

	plot := hplot.New()
	plot.Title.Text = robj.Title()

	h := hplot.NewH1D(h1d)
	h.Infos.Style = hplot.HInfoSummary
	h.Color = color.RGBA{255, 0, 0, 255}

	plot.Add(h, hplot.NewGrid())

	svg, err := renderSVG(plot)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(string(svg))
}

func (srv *server) plotH2Handle(w http.ResponseWriter, r *http.Request) error {
	url := r.URL.Path[len("/plot-h2/"):]
	toks := strings.Split(url, "/")
	fname := toks[0]

	db, err := srv.db(r)
	if err != nil {
		return err
	}
	db.RLock()
	defer db.RUnlock()

	f := db.get(fname)
	obj, err := walk(f, toks[1:])
	if err != nil {
		return fmt.Errorf("could not find %q in file %q: %v", filepath.Join(toks[1:]...), fname, err)
	}

	robj, ok := obj.(rootio.H2)
	if !ok {
		return fmt.Errorf("object %q could not be converted to hbook.H1D", toks[1])
	}
	h2d, err := rootcnv.H2D(robj)
	if err != nil {
		return err
	}

	plot := hplot.New()
	plot.Title.Text = robj.Title()

	h := hplot.NewH2D(h2d, nil)
	h.Infos.Style = hplot.HInfoSummary

	plot.Add(h, hplot.NewGrid())

	svg, err := renderSVG(plot)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(string(svg))
}

func (srv *server) plotS2Handle(w http.ResponseWriter, r *http.Request) error {
	url := r.URL.Path[len("/plot-s2/"):]
	toks := strings.Split(url, "/")
	fname := toks[0]

	db, err := srv.db(r)
	if err != nil {
		return err
	}
	db.RLock()
	defer db.RUnlock()

	f := db.get(fname)
	obj, err := walk(f, toks[1:])
	if err != nil {
		return fmt.Errorf("could not find %q in file %q: %v", filepath.Join(toks[1:]...), fname, err)
	}

	robj, ok := obj.(rootio.Graph)
	if !ok {
		return fmt.Errorf("object %q could not be converted to rootio.Graph", toks[1])
	}
	s2d, err := rootcnv.S2D(robj)
	if err != nil {
		return err
	}

	plot := hplot.New()
	plot.Title.Text = robj.Title()

	var opts hplot.Options
	if _, ok := obj.(rootio.GraphErrors); ok {
		opts = hplot.WithXErrBars | hplot.WithYErrBars
	}
	h := hplot.NewS2D(s2d, opts)
	if err != nil {
		return err
	}
	h.Color = color.RGBA{255, 0, 0, 255}

	plot.Add(h, hplot.NewGrid())

	svg, err := renderSVG(plot)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(string(svg))
}

func (srv *server) plotBranchHandle(w http.ResponseWriter, r *http.Request) error {
	url := r.URL.Path[len("/plot-branch/"):]
	toks := strings.Split(url, "/")
	fname := toks[0]

	db, err := srv.db(r)
	if err != nil {
		return err
	}
	db.RLock()
	defer db.RUnlock()

	f := db.get(fname)
	obj, err := walk(f, toks[1:])
	if err != nil {
		return fmt.Errorf("could not find %q in file %q: %v", filepath.Join(toks[1:]...), fname, err)
	}

	tree := obj.(rootio.Tree)

	bname := toks[len(toks)-1]
	b := tree.Branch(bname) // FIXME(sbinet): handle sub-branches
	if b == nil {
		return fmt.Errorf("could not find branch %q in tree %q of file %q", bname, tree.Name(), fname)
	}

	leaves := b.Leaves()
	leaf := leaves[0]
	fv, err := newFloats(leaf)
	if err != nil {
		log.Printf("error creating float-val: %v\n", err)
		return err
	}

	min := +math.MaxFloat64
	max := -math.MaxFloat64
	vals := make([]float64, 0, int(tree.Entries()))
	sc, err := rootio.NewTreeScannerVars(tree, rootio.ScanVar{Name: bname, Leaf: leaf.Name()})
	if err != nil {
		return fmt.Errorf("error creating scanner for branch %q in tree %q of file %q: %v", bname, tree.Name(), fname, err)
	}
	defer sc.Close()

	for sc.Next() {
		err = sc.Scan(fv.ptr)
		if err != nil {
			log.Printf("error scan: %v\n", err)
			return err
		}
		for _, v := range fv.vals() {
			max = math.Max(max, v)
			min = math.Min(min, v)
			vals = append(vals, v)
		}
	}

	err = sc.Err()
	if err != nil {
		log.Printf("error finding min/max: %v\n", err)
		return err
	}

	err = sc.Close()
	if err != nil {
		log.Printf("error closing min/max-scanner: %v\n", err)
		return err
	}

	min = math.Nextafter(min, min-1)
	max = math.Nextafter(max, max+1)
	h := hbook.NewH1D(100, min, max)
	for _, v := range vals {
		h.Fill(v, 1)
	}

	plot := hplot.New()
	plot.Title.Text = leaf.Name()

	hh := hplot.NewH1D(h)
	hh.Infos.Style = hplot.HInfoSummary
	hh.Color = color.RGBA{255, 0, 0, 255}

	plot.Add(hh, hplot.NewGrid())

	svg, err := renderSVG(plot)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(string(svg))
}

type floats struct {
	leaf rootio.Leaf
	ptr  interface{}
	vals func() []float64
}

func newFloats(leaf rootio.Leaf) (floats, error) {
	fv := floats{leaf: leaf}
	n := 1 // scalar
	switch {
	case leaf.LeafCount() != nil:
		n = -1
	case leaf.Len() > 1:
		n = leaf.Len()
	}

	switch leaf.TypeName() {
	case "bool":
		switch n {
		case 1:
			var vv bool
			fv.ptr = &vv
			fv.vals = func() []float64 {
				b := *fv.ptr.(*bool)
				if b {
					return []float64{1}
				}
				return []float64{0}
			}
		case -1:
			var vv []bool
			fv.ptr = &vv
			fv.vals = func() []float64 {
				bs := *fv.ptr.(*[]bool)
				vs := make([]float64, len(bs))
				for i, b := range bs {
					if b {
						vs[i] = 1
					} else {
						vs[i] = 0
					}
				}
				return vs
			}
		default:
			vv := newArrB(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					b := vv.Elem().Index(i).Bool()
					if b {
						vs[i] = 1
					} else {
						vs[i] = 0
					}
				}
				return vs
			}
		}
	case "uint8", "byte":
		switch n {
		case 1:
			var vv uint8
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*uint8))} }
		case -1:
			var vv []uint8
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]uint8)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrU8(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "uint16":
		switch n {
		case 1:
			var vv uint16
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*uint16))} }
		case -1:
			var vv []uint16
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]uint16)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrU16(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "uint32":
		switch n {
		case 1:
			var vv uint32
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*uint32))} }
		case -1:
			var vv []uint32
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]uint32)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrU32(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "uint64":
		switch n {
		case 1:
			var vv uint64
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*uint64))} }
		case -1:
			var vv []uint64
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]uint64)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrU64(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "int8":
		switch n {
		case 1:
			var vv int8
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*int8))} }
		case -1:
			var vv []int8
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]int8)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrI8(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "int16":
		switch n {
		case 1:
			var vv int16
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*int16))} }
		case -1:
			var vv []int16
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]int16)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrI16(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "int32":
		switch n {
		case 1:
			var vv int32
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*int32))} }
		case -1:
			var vv []int32
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]int32)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrI32(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "int64":
		switch n {
		case 1:
			var vv int64
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*int64))} }
		case -1:
			var vv []int64
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]int64)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrI64(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = float64(vv.Elem().Index(i).Int())
				}
				return vs
			}
		}
	case "float32":
		switch n {
		case 1:
			var vv float32
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{float64(*fv.ptr.(*float32))} }
		case -1:
			var vv []float32
			fv.ptr = &vv
			fv.vals = func() []float64 {
				vv := *fv.ptr.(*[]float32)
				vs := make([]float64, len(vv))
				for i, v := range vv {
					vs[i] = float64(v)
				}
				return vs
			}
		default:
			vv := newArrF32(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = vv.Elem().Index(i).Float()
				}
				return vs
			}
		}
	case "float64":
		switch n {
		case 1:
			var vv float64
			fv.ptr = &vv
			fv.vals = func() []float64 { return []float64{*fv.ptr.(*float64)} }
		case -1:
			var vv []float64
			fv.ptr = &vv
			fv.vals = func() []float64 { return *fv.ptr.(*[]float64) }
		default:
			vv := newArrF64(n)
			fv.ptr = vv.Interface()
			fv.vals = func() []float64 {
				vs := make([]float64, n)
				for i := range vs {
					vs[i] = vv.Elem().Index(i).Float()
				}
				return vs
			}
		}
	default:
		return fv, fmt.Errorf("unhandled value of type %q", leaf.TypeName())
	}
	return fv, nil
}

func newArrB(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*bool)(nil)).Elem())
	return reflect.New(typ)
}

func newArrU8(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*uint8)(nil)).Elem())
	return reflect.New(typ)
}

func newArrU16(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*uint16)(nil)).Elem())
	return reflect.New(typ)
}

func newArrU32(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*uint32)(nil)).Elem())
	return reflect.New(typ)
}

func newArrU64(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*uint64)(nil)).Elem())
	return reflect.New(typ)
}

func newArrI8(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*int8)(nil)).Elem())
	return reflect.New(typ)
}

func newArrI16(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*int16)(nil)).Elem())
	return reflect.New(typ)
}

func newArrI32(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*int32)(nil)).Elem())
	return reflect.New(typ)
}

func newArrI64(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*int64)(nil)).Elem())
	return reflect.New(typ)
}

func newArrF32(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*float32)(nil)).Elem())
	return reflect.New(typ)
}

func newArrF64(n int) reflect.Value {
	typ := reflect.ArrayOf(n, reflect.TypeOf((*float64)(nil)).Elem())
	return reflect.New(typ)
}
