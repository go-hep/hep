// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"unsafe"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
)

// rleafCtx is the interface that wraps the rcount method.
type rleafCtx interface {
	// rcountFunc returns the function that gives the leaf-count
	// of the provided leaf.
	rcountFunc(leaf string) func() int
	rcountLeaf(leaf string) leafCount
}

// rleaf is the leaf reading interface.
type rleaf interface {
	Leaf() Leaf
	Offset() int64
	readFromBuffer(*rbytes.RBuffer) error
}

// rleafDefaultSliceCap is the default capacity for all
// rleaves that hold slices of data.
const rleafDefaultSliceCap = 8

func rleafFrom(leaf Leaf, rvar ReadVar, rctx rleafCtx) rleaf {
	switch leaf := leaf.(type) {
	case *LeafO:
		return newRLeafBool(leaf, rvar, rctx)
	case *LeafB:
		switch rv := reflect.ValueOf(rvar.Value); rv.Interface().(type) {
		case *int8, *[]int8:
			return newRLeafI8(leaf, rvar, rctx)
		case *uint8, *[]uint8:
			return newRLeafU8(leaf, rvar, rctx)
		default:
			rv := rv.Elem()
			switch rv.Kind() {
			case reflect.Array:
				rt, _ := flattenArrayType(rv.Type())
				switch rt.Kind() {
				case reflect.Int8:
					return newRLeafI8(leaf, rvar, rctx)
				case reflect.Uint8:
					return newRLeafU8(leaf, rvar, rctx)
				}
			case reflect.Slice:
				rt, _ := flattenArrayType(rv.Type().Elem())
				switch rt.Kind() {
				case reflect.Int8:
					return newRLeafI8(leaf, rvar, rctx)
				case reflect.Uint8:
					return newRLeafU8(leaf, rvar, rctx)
				}
			}
		}
		panic(fmt.Errorf("rvar mismatch for %T", leaf))
	case *LeafS:
		switch rv := reflect.ValueOf(rvar.Value); rv.Interface().(type) {
		case *int16, *[]int16:
			return newRLeafI16(leaf, rvar, rctx)
		case *uint16, *[]uint16:
			return newRLeafU16(leaf, rvar, rctx)
		default:
			rv := rv.Elem()
			switch rv.Kind() {
			case reflect.Array:
				rt, _ := flattenArrayType(rv.Type())
				switch rt.Kind() {
				case reflect.Int16:
					return newRLeafI16(leaf, rvar, rctx)
				case reflect.Uint16:
					return newRLeafU16(leaf, rvar, rctx)
				}
			case reflect.Slice:
				rt, _ := flattenArrayType(rv.Type().Elem())
				switch rt.Kind() {
				case reflect.Int16:
					return newRLeafI16(leaf, rvar, rctx)
				case reflect.Uint16:
					return newRLeafU16(leaf, rvar, rctx)
				}
			}
		}
		panic(fmt.Errorf("rvar mismatch for %T", leaf))
	case *LeafI:
		switch rv := reflect.ValueOf(rvar.Value); rv.Interface().(type) {
		case *int32, *[]int32:
			return newRLeafI32(leaf, rvar, rctx)
		case *uint32, *[]uint32:
			return newRLeafU32(leaf, rvar, rctx)
		default:
			rv := rv.Elem()
			switch rv.Kind() {
			case reflect.Array:
				rt, _ := flattenArrayType(rv.Type())
				switch rt.Kind() {
				case reflect.Int32:
					return newRLeafI32(leaf, rvar, rctx)
				case reflect.Uint32:
					return newRLeafU32(leaf, rvar, rctx)
				}
			case reflect.Slice:
				rt, _ := flattenArrayType(rv.Type().Elem())
				switch rt.Kind() {
				case reflect.Int32:
					return newRLeafI32(leaf, rvar, rctx)
				case reflect.Uint32:
					return newRLeafU32(leaf, rvar, rctx)
				}
			}
		}
		panic(fmt.Errorf("rvar mismatch for %T", leaf))
	case *LeafL:
		switch rv := reflect.ValueOf(rvar.Value); rv.Interface().(type) {
		case *int64, *[]int64:
			return newRLeafI64(leaf, rvar, rctx)
		case *uint64, *[]uint64:
			return newRLeafU64(leaf, rvar, rctx)
		default:
			rv := rv.Elem()
			switch rv.Kind() {
			case reflect.Array:
				rt, _ := flattenArrayType(rv.Type())
				switch rt.Kind() {
				case reflect.Int64:
					return newRLeafI64(leaf, rvar, rctx)
				case reflect.Uint64:
					return newRLeafU64(leaf, rvar, rctx)
				}
			case reflect.Slice:
				rt, _ := flattenArrayType(rv.Type().Elem())
				switch rt.Kind() {
				case reflect.Int64:
					return newRLeafI64(leaf, rvar, rctx)
				case reflect.Uint64:
					return newRLeafU64(leaf, rvar, rctx)
				}
			}
			panic(fmt.Errorf("rvar mismatch for %T", leaf))
		}
	case *LeafG:
		// FIXME(sbinet): should we bite the bullet and generate a whole
		// set of types+funcs for LeafG instead of relying on the
		// assumption that LeafG data has the same underlying layout and size
		// than LeafL ? (ie: sizeof(Long_t) == sizeof(Long64_t))
		return rleafFrom((*LeafL)(unsafe.Pointer(leaf)), rvar, rctx)
	case *LeafF:
		return newRLeafF32(leaf, rvar, rctx)
	case *LeafD:
		return newRLeafF64(leaf, rvar, rctx)
	case *LeafF16:
		return newRLeafF16(leaf, rvar, rctx)
	case *LeafD32:
		return newRLeafD32(leaf, rvar, rctx)
	case *LeafC:
		return newRLeafStr(leaf, rvar, rctx)

	case *tleafElement:
		return newRLeafElem(leaf, rvar, rctx)

	case *tleafObject:
		return newRLeafObject(leaf, rvar, rctx)

	default:
		panic(fmt.Errorf("not implemented %T", leaf))
	}
}

type rleafObject struct {
	base *tleafObject
	v    rbytes.Unmarshaler
}

var (
	_ rleaf = (*rleafObject)(nil)
)

func newRLeafObject(leaf *tleafObject, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		panic("not implemented")
	case leaf.len > 1:
		panic("not implemented")
	default:
		return &rleafObject{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(rbytes.Unmarshaler),
		}
	}
}

func (leaf *rleafObject) Leaf() Leaf { return leaf.base }

func (leaf *rleafObject) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafObject) readFromBuffer(r *rbytes.RBuffer) error {
	if leaf.base.virtual {
		var (
			n     = int(r.ReadU8())
			class = r.ReadCString(n + 1)
		)
		if class != leaf.base.Title() {
			// FIXME(sbinet): we should be able to handle (C++) polymorphism.
			// but in Go, this should translate to interfaces.
			panic(fmt.Errorf(
				"rtree: rleaf object with incompatible class names: got=%q, want=%q",
				class, leaf.base.Title(),
			))
		}
	}

	return leaf.v.UnmarshalROOT(r)
}

func newRLeafElem(leaf *tleafElement, rvar ReadVar, rctx rleafCtx) rleaf {
	const kind = rbytes.ObjectWise // FIXME(sbinet): infer from stream?

	var (
		b         = leaf.branch.(*tbranchElement)
		si        = b.streamer
		err       error
		rstreamer rbytes.RStreamer
	)

	switch {
	case b.id < 0:
		rstreamer, err = si.NewRStreamer(kind)
	default:
		rstreamer, err = rdict.RStreamerOf(si, int(b.id), kind)
	}

	if err != nil {
		panic(fmt.Errorf(
			"rtree: could not find read-streamer for leaf=%q (type=%s): %+v",
			leaf.Name(), leaf.TypeName(), err,
		))
	}
	err = rstreamer.(rbytes.Binder).Bind(rvar.Value)
	if err != nil {
		panic(fmt.Errorf("rtree: could not bind read-streamer for leaf=%q (type=%s) to ptr=%T: %w",
			leaf.Name(), leaf.TypeName(), rvar.Value, err,
		))
	}

	if leaf.count != nil {
		r, ok := rstreamer.(rbytes.Counter)
		if !ok {
			panic(fmt.Errorf(
				"rtree: could not set read-streamer counter for leaf=%q (type=%s): %+v",
				leaf.Name(), leaf.TypeName(), err,
			))
		}
		lc := rctx.rcountLeaf(leaf.count.Name())
		err = r.Count(lc.ivalue)
		if err != nil {
			panic(fmt.Errorf(
				"rtree: could not set read-streamer counter for leaf=%q (type=%s): %+v",
				leaf.Name(), leaf.TypeName(), err,
			))
		}
	}

	return &rleafElem{
		base:     leaf,
		v:        rvar.Value,
		streamer: rstreamer,
	}
}

type rleafElem struct {
	base     *tleafElement
	v        interface{}
	n        func() int
	streamer rbytes.RStreamer
}

func (leaf *rleafElem) Leaf() Leaf { return leaf.base }

func (leaf *rleafElem) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafElem) readFromBuffer(r *rbytes.RBuffer) error {
	return leaf.streamer.RStreamROOT(r)
}

func (leaf *rleafElem) bindCount() {
	switch v := reflect.ValueOf(leaf.v).Interface().(type) {
	case *int8:
		leaf.n = func() int { return int(*v) }
	case *int16:
		leaf.n = func() int { return int(*v) }
	case *int32:
		leaf.n = func() int { return int(*v) }
	case *int64:
		leaf.n = func() int { return int(*v) }
	case *uint8:
		leaf.n = func() int { return int(*v) }
	case *uint16:
		leaf.n = func() int { return int(*v) }
	case *uint32:
		leaf.n = func() int { return int(*v) }
	case *uint64:
		leaf.n = func() int { return int(*v) }
	default:
		panic(fmt.Errorf("invalid leaf-elem type: %T", v))
	}
}

func (leaf *rleafElem) ivalue() int {
	return leaf.n()
}

var (
	_ rleaf = (*rleafElem)(nil)
)

type rleafCount struct {
	Leaf
	n    func() int
	leaf rleaf
}

func (l *rleafCount) ivalue() int {
	return l.n()
}

func (l *rleafCount) imax() int {
	panic("not implemented")
}

var (
	_ leafCount = (*rleafCount)(nil)
)
