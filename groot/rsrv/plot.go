// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsrv

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgeps"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgpdf"
	"gonum.org/v1/plot/vg/vgsvg"
	"gonum.org/v1/plot/vg/vgtex"
)

func (srv *Server) render(p *hplot.Plot, opt PlotOptions) ([]byte, error) {
	var canvas vg.CanvasWriterTo

	switch opt.Type {
	case "eps":
		canvas = vgeps.NewTitle(opt.Width, opt.Height, p.Title.Text)
	case "jpg", "jpeg":
		canvas = vgimg.JpegCanvas{vgimg.New(opt.Width, opt.Height)}
	case "pdf":
		canvas = vgpdf.New(opt.Width, opt.Height)
	case "png":
		canvas = vgimg.PngCanvas{vgimg.New(opt.Width, opt.Height)}
	case "svg":
		canvas = vgsvg.New(opt.Width, opt.Height)
	case "tex":
		canvas = vgtex.New(opt.Width, opt.Height)
	case "tiff":
		canvas = vgimg.TiffCanvas{vgimg.New(opt.Width, opt.Height)}
	}

	p.Draw(draw.New(canvas))

	out := new(bytes.Buffer)
	_, err := canvas.WriteTo(out)
	if err != nil {
		return nil, errors.Wrap(err, "could not write canvas")
	}

	return out.Bytes(), nil
}

type floats struct {
	leaf rtree.Leaf
	ptr  interface{}
	vals func() []float64
}

func newFloats(leaf rtree.Leaf) (floats, error) {
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
