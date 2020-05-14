// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
)

type rleafCtx interface {
	rcount(leaf string) func() int
}

type rleaf interface {
	Leaf() Leaf
	Offset() int64
	readFromBuffer(*rbytes.RBuffer) error
}

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
			if rv.Kind() == reflect.Array {
				switch rv.Type().Elem().Kind() {
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
			if rv.Kind() == reflect.Array {
				switch rv.Type().Elem().Kind() {
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
			if rv.Kind() == reflect.Array {
				switch rv.Type().Elem().Kind() {
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
			if rv.Kind() == reflect.Array {
				switch rv.Type().Elem().Kind() {
				case reflect.Int64:
					return newRLeafI64(leaf, rvar, rctx)
				case reflect.Uint64:
					return newRLeafU64(leaf, rvar, rctx)
				}
			}
			panic(fmt.Errorf("rvar mismatch for %T", leaf))
		}
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

	default:
		panic(fmt.Errorf("not implemented %T", leaf))
	}
}
