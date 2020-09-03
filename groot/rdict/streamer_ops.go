// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"log"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

type rstreamOp interface {
	rstream(r *rbytes.RBuffer, recv interface{}) error
}

type wstreamOp interface {
	wstream(w *rbytes.WBuffer, recv interface{}) (int, error)
}

type rstreamBufOp interface {
	rstreamBuf(r *rbytes.RBuffer, recv reflect.Value, descr *elemDescr, beg, end int, n int, offset int, arrmode arrayMode) error
}

type wstreamBufOp interface {
	wstreamBuf(w *rbytes.WBuffer, recv reflect.Value, descr *elemDescr, beg, end int, n int, offset int, arrmode arrayMode) (int, error)
}

type arrayMode int32

type elemDescr struct {
	otype  rmeta.Enum
	ntype  rmeta.Enum
	offset int // actually an index to the struct's field or to array's element
	length int
	elem   rbytes.StreamerElement
	method int
	oclass string
	nclass string
	mbr    interface{} // member streamer
}

type streamerConfig struct {
	si     *StreamerInfo
	eid    int // element ID
	descr  *elemDescr
	offset int // offset/index within object
	length int // number of elements for fixed-length arrays
}

func (cfg *streamerConfig) counter(recv interface{}) int {
	return int(reflect.ValueOf(recv).Elem().Field(cfg.descr.method).Int())
}

func (cfg *streamerConfig) adjust(recv interface{}) interface{} {
	if cfg == nil {
		return recv
	}
	rv := reflect.ValueOf(recv).Elem()
	switch rv.Kind() {
	case reflect.Struct:
		return rv.Field(cfg.offset).Addr().Interface()
	case reflect.Array, reflect.Slice:
		return rv.Index(cfg.offset).Addr().Interface()
	default:
		return recv
	}
}

func (si *StreamerInfo) makeReadOp(sictx rbytes.StreamerInfoContext, i int, descr elemDescr) rstreamer {
	log.Printf("--- makeReadOp(%d, %q, %#v)...", i, si.Name()+"."+si.elems[i].Name(), descr)
	switch descr.otype {
	case rmeta.Bool:
		return rstreamer{rstreamBool, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Char:
		return rstreamer{rstreamI8, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Short:
		return rstreamer{rstreamI16, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Int:
		return rstreamer{rstreamI32, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Long, rmeta.Long64:
		return rstreamer{rstreamI64, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.UChar:
		return rstreamer{rstreamU8, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.UShort:
		return rstreamer{rstreamU16, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.UInt:
		return rstreamer{rstreamU32, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.ULong, rmeta.ULong64:
		return rstreamer{rstreamU64, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Float32:
		return rstreamer{rstreamF32, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Float64:
		return rstreamer{rstreamF64, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Bits:
		return rstreamer{rstreamBits, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Float16:
		return rstreamer{rstreamF16(descr.elem), &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.Double32:
		return rstreamer{rstreamD32(descr.elem), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Counter:
		se := descr.elem.(*StreamerBasicType)
		switch se.esize {
		case 4:
			return rstreamer{rstreamI32, &streamerConfig{si, i, &descr, descr.offset, 0}}
		case 8:
			return rstreamer{rstreamI64, &streamerConfig{si, i, &descr, descr.offset, 0}}
		default:
			panic(fmt.Errorf("rdict: invalid counter size (%d) in %#v", se.esize, se))
		}

	case rmeta.TNamed:
		return rstreamer{rstreamTNamed, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.TObject:
		return rstreamer{rstreamTObject, &streamerConfig{si, i, &descr, descr.offset, 0}}
	case rmeta.TString:
		return rstreamer{rstreamTString, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.STL:
		var (
			se       = descr.elem
			newClass = descr.nclass
			oldClass = descr.oclass
			// _, isSTLbase = se.(*StreamerBase) // FIXME(sbinet)
		)

		switch {
		case se.ArrayLen() <= 1:
			switch {
			case newClass != oldClass:
				panic("not implemented")
			default:
				switch se := se.(type) {
				case *StreamerSTL:
					switch se.STLType() {
					case rmeta.STLvector:
						panic("not implemented")
					case rmeta.STLmap, rmeta.STLmultimap,
						rmeta.STLset, rmeta.STLmultiset,
						rmeta.STLunorderedmap, rmeta.STLunorderedmultimap,
						rmeta.STLunorderedset, rmeta.STLunorderedmultiset:
						panic("not implemented")
					default:
						panic("not implemented")
					}
				case *StreamerSTLstring:
					panic("not implemented")
				default:
					panic("not implemented")
				}
			}
		default:
			switch {
			case newClass != oldClass:
				panic("not implemented")
			default:
				return rstreamer{
					rstreamSTL(rstreamSTLArrayMbrWise, rstreamSTLObjWise, descr.oclass),
					&streamerConfig{si, i, &descr, descr.offset, se.ArrayLen()},
				}
			}
		}

	case rmeta.Conv + rmeta.Bool:
		return rstreamer{rstreamCnv(descr.ntype, rstreamBool), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Char:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI8), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Short:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI16), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Int:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI32), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Long, rmeta.Conv + rmeta.Long64:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI64), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.UChar:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU8), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.UShort:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU16), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.UInt:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU32), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.ULong, rmeta.Conv + rmeta.ULong64:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU64), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Float32:
		return rstreamer{rstreamCnv(descr.ntype, rstreamF32), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Float64:
		return rstreamer{rstreamCnv(descr.ntype, rstreamF64), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Bits:
		return rstreamer{rstreamCnv(descr.ntype, rstreamBits), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Float16:
		return rstreamer{rstreamCnv(descr.ntype, rstreamF16(descr.elem)), &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Conv + rmeta.Double32:
		return rstreamer{rstreamCnv(descr.ntype, rstreamD32(descr.elem)), &streamerConfig{si, i, &descr, descr.offset, 0}}

		// fixed-size arrays of basic types: [32]int

	case rmeta.OffsetL + rmeta.Bool:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamBool), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.Char:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamI8), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.Short:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamI16), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.Int:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamI32), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.Long, rmeta.OffsetL + rmeta.Long64:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamI64), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.UChar:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamU8), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.UShort:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamU16), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.UInt:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamI32), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.ULong, rmeta.OffsetL + rmeta.ULong64:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamU64), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.Float32:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamF32), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.Float64:
		alen := descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(alen, rstreamF64), &streamerConfig{si, i, &descr, descr.offset, alen}}

	case rmeta.OffsetL + rmeta.Float16:
		alen := descr.elem.ArrayLen()
		return rstreamer{
			rstreamBasicArray(alen, rstreamCnv(descr.ntype, rstreamF16(descr.elem))),
			&streamerConfig{si, i, &descr, descr.offset, alen},
		}

	case rmeta.OffsetL + rmeta.Double32:
		alen := descr.elem.ArrayLen()
		return rstreamer{
			rstreamBasicArray(alen, rstreamCnv(descr.ntype, rstreamD32(descr.elem))),
			&streamerConfig{si, i, &descr, descr.offset, alen},
		}

		// var-size arrays of basic types: [n]int

	case rmeta.OffsetP + rmeta.Bool:
		return rstreamer{rstreamBools, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Char:
		return rstreamer{rstreamI8s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Short:
		return rstreamer{rstreamI16s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Int:
		return rstreamer{rstreamI32s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Long, rmeta.OffsetP + rmeta.Long64:
		return rstreamer{rstreamI64s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.UChar:
		return rstreamer{rstreamU8s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.UShort:
		return rstreamer{rstreamU16s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.UInt:
		return rstreamer{rstreamI32s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.ULong, rmeta.OffsetP + rmeta.ULong64:
		return rstreamer{rstreamU64s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Float32:
		return rstreamer{rstreamF32s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Float64:
		return rstreamer{rstreamF64s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Float16:
		return rstreamer{rstreamF16s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.OffsetP + rmeta.Double32:
		return rstreamer{rstreamD32s, &streamerConfig{si, i, &descr, descr.offset, 0}}

	case rmeta.Streamer:
		switch se := descr.elem.(type) {
		case *StreamerSTLstring:
			return rstreamer{rstreamStdString, &streamerConfig{si, i, &descr, descr.offset, 0}}
		case *StreamerSTL:
			switch se.STLType() {
			case rmeta.STLvector:
			default:
				panic(fmt.Errorf("rdict: STL container type=%v not handled", se.STLType()))
			}
		}

		se := descr.elem
		osi, err := sictx.StreamerInfo(se.TypeName(), -1)
		if err != nil {
			panic(fmt.Errorf("rdict: could not find streamer info for element %q (type=%q) of %q: %w",
				se.Name(), se.TypeName(), si.Name(),
				err,
			))
		}
		err = osi.BuildStreamers()
		if err != nil {
			panic(fmt.Errorf("rdict: could not build streamers for %q (element %q of streamer %q): %w",
				osi.Name(), se.Name(), si.Name(), err,
			))
		}

		esi := osi.(*StreamerInfo)
		return rstreamer{
			rstreamCat(se.TypeName(), int16(esi.ClassVersion()), esi.roops),
			&streamerConfig{si, i, &descr, descr.offset, 0},
		}

	case rmeta.Any:
		se := descr.elem
		osi, err := sictx.StreamerInfo(se.TypeName(), -1)
		if err != nil {
			panic(fmt.Errorf("rdict: could not find streamer info for element %q (type=%q) of %q: %w",
				se.Name(), se.TypeName(), si.Name(),
				err,
			))
		}
		err = osi.BuildStreamers()
		if err != nil {
			panic(fmt.Errorf("rdict: could not build streamers for %q (element %q of streamer %q): %w",
				osi.Name(), se.Name(), si.Name(), err,
			))
		}

		esi := osi.(*StreamerInfo)
		return rstreamer{
			rstreamCat(se.TypeName(), int16(esi.ClassVersion()), esi.roops),
			&streamerConfig{si, i, &descr, descr.offset, 0},
		}

	default:
		panic(fmt.Errorf("not implemented k=%d (%v)", descr.otype, descr.otype))
		return rstreamer{rstreamGeneric, &streamerConfig{si, i, &descr, descr.offset, 0}}
	}

	panic("impossible")
}

func rstreamBool(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readBool(cfg.adjust(recv), r)
}

func rstreamU8(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readU8(cfg.adjust(recv), r)
}

func rstreamU16(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readU16(cfg.adjust(recv), r)
}

func rstreamU32(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readU32(cfg.adjust(recv), r)
}

func rstreamU64(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readU64(cfg.adjust(recv), r)
}

func rstreamI8(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readI8(cfg.adjust(recv), r)
}

func rstreamI16(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readI16(cfg.adjust(recv), r)
}

func rstreamI32(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readI32(cfg.adjust(recv), r)
}

func rstreamI64(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readI64(cfg.adjust(recv), r)
}

func rstreamF32(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readF32(cfg.adjust(recv), r)
}

func rstreamF64(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readF64(cfg.adjust(recv), r)
}

func rstreamBits(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	bits := r.ReadU32()
	recv = cfg.adjust(recv)
	*(recv.(*uint32)) = bits
	// FIXME(sbinet) handle TObject reference
	// if (bits&kIsReferenced) != 0 { ... }
	return r.Err()
}

func rstreamF16(se rbytes.StreamerElement) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		recv = cfg.adjust(recv)
		*(recv.(*root.Float16)) = r.ReadF16(se)
		return r.Err()
	}
}

func rstreamD32(se rbytes.StreamerElement) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		recv = cfg.adjust(recv)
		*(recv.(*root.Double32)) = r.ReadD32(se)
		return r.Err()
	}
}

func rstreamTString(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	return readStr(cfg.adjust(recv), r)
}

func rstreamTObject(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	obj := cfg.adjust(recv).(*rbase.Object)
	return obj.UnmarshalROOT(r)
}

func rstreamTNamed(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	named := cfg.adjust(recv).(*rbase.Named)
	return named.UnmarshalROOT(r)
}

func rstreamCnv(to rmeta.Enum, from ropFunc) ropFunc {
	switch to {
	case rmeta.Bool:
	case rmeta.Char:
	case rmeta.Short:
	case rmeta.Int:
	case rmeta.Long, rmeta.Long64:
	case rmeta.UChar:
	case rmeta.UShort:
	case rmeta.UInt:
	case rmeta.ULong, rmeta.ULong64:
	case rmeta.Float32:
	case rmeta.Float64:
	case rmeta.Float16:
	case rmeta.Double32:
	case rmeta.Bits:
	}

	panic("not implemented")
}

func rstreamBasicArray(n int, arr ropFunc) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		rv := reflect.ValueOf(cfg.adjust(recv)).Elem()
		for i := 0; i < n; i++ {
			err := arr(r, rv.Index(i).Addr().Interface(), nil)
			if err != nil {
				return fmt.Errorf(
					"rdict: could not rstream element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, cfg.si.Name(), err,
				)
			}
		}
		return nil
	}
}

func rstreamBasicSlice(arr ropFunc) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		_ = r.ReadI8() // is-array
		n := int(reflect.ValueOf(recv).Elem().Field(cfg.descr.method).Int())
		rv := reflect.ValueOf(cfg.adjust(recv)).Elem()
		for i := 0; i < n; i++ {
			err := arr(r, rv.Index(i).Addr().Interface(), nil)
			if err != nil {
				return fmt.Errorf(
					"rdict: could not rstream element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, cfg.si.Name(), err,
				)
			}
		}
		return nil
	}
}

func rstreamBools(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]bool)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]bool, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayBool(*sli)
	return r.Err()
}

func rstreamU8s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint8)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]uint8, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayU8(*sli)
	return r.Err()
}

func rstreamU16s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint16)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]uint16, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayU16(*sli)
	return r.Err()
}

func rstreamU32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint32)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]uint32, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayU32(*sli)
	return r.Err()
}

func rstreamU64s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint64)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]uint64, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayU64(*sli)
	return r.Err()
}

func rstreamI8s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int8)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]int8, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayI8(*sli)
	return r.Err()
}

func rstreamI16s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int16)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]int16, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayI16(*sli)
	return r.Err()
}

func rstreamI32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int32)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]int32, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayI32(*sli)
	return r.Err()
}

func rstreamI64s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int64)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]int64, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayI64(*sli)
	return r.Err()
}

func rstreamF32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]float32)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]float32, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayF32(*sli)
	return r.Err()
}

func rstreamF64s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]float64)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]float64, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayF64(*sli)
	return r.Err()
}

func rstreamF16s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]root.Float16)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]root.Float16, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayF16(*sli, cfg.descr.elem)
	return r.Err()
}

func rstreamD32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]root.Double32)
	)
	if nn := len(*sli); nn < n {
		*sli = append(*sli, make([]root.Double32, n-nn)...)
	}
	*sli = (*sli)[:n]
	r.ReadArrayD32(*sli, cfg.descr.elem)
	return r.Err()
}

func rstreamCat(typename string, typevers int16, rops []rstreamer) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		beg := r.Pos()
		vers, pos, bcnt := r.ReadVersion(cfg.descr.oclass)
		if vers != typevers {
			r.SetErr(fmt.Errorf(
				"rdict: inconsistent ROOT version type=%q (got=%d, want=%d)",
				typename, vers, typevers,
			))
			return r.Err()
		}

		recv = cfg.adjust(recv)
		for i, rop := range rops {
			err := rop.rstream(r, rop.cfg.adjust(recv))
			if err != nil {
				return fmt.Errorf(
					"rdict: could not rstream element %d (%s) of %s: %w",
					i, rop.cfg.descr.elem.Name(), cfg.si.Name(), err,
				)
			}
		}
		r.CheckByteCount(pos, bcnt, beg, cfg.descr.oclass)
		return r.Err()
	}
}

func rstreamStdString(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	beg := r.Pos()
	_ /*vers*/, pos, bcnt := r.ReadVersion("string")
	*(cfg.adjust(recv).(*string)) = r.ReadString()
	r.CheckByteCount(pos, bcnt, beg, "string")
	return r.Err()

}

func rstreamGeneric(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	const (
		beg     = 0
		end     = 1
		n       = 1
		arrmode = 2
	)
	return cfg.si.rstream(r, recv, cfg.descr, beg, end, n, cfg.offset, arrmode)
}

type rmbrwiseSTLFunc func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16) error
type robjwiseSTLFunc func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16, start int32) error

func rstreamSTL(mbrwise rmbrwiseSTLFunc, objwise robjwiseSTLFunc, oclass string) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		err := r.Err()
		if err != nil {
			return err
		}

		beg := r.Pos()
		vers, pos, bcnt := r.ReadVersion(oclass)
		switch {
		case vers&rbytes.StreamedMemberWise != 0:
			err = mbrwise(r, cfg.adjust(recv), cfg, vers)
		default:
			err = objwise(r, cfg.adjust(recv), cfg, vers, int32(beg))
		}

		r.CheckByteCount(pos, bcnt, beg, oclass)
		return err
	}
}

func rstreamSTLArrayMbrWise(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16) error {
	panic("not implemented")
}

func rstreamSTLObjWise(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16, start int32) error {
	panic("not implemented")
}

func (si *StreamerInfo) rstream(r *rbytes.RBuffer, recv interface{}, descr *elemDescr, beg, end, n, offset int, mode arrayMode) error {
	needIncr := mode&2 == 0
	mode &= ^2

	if needIncr {
		panic("not implemented")
	}
	const kHaveLoop = 1024
	var typeOffset rmeta.Enum
	if mode == 0 {
		typeOffset = kHaveLoop
	}

	var (
		ioffset = []int{-1, offset}
		err     error
	)
	for i := beg; i < end; i++ {
		ioffset[0] = descr.offset

		var (
			recv  = ptrAdjust(recv, ioffset)
			etype = descr.otype
		)

		switch etype + typeOffset {
		case rmeta.Bool:
			err = readBool(recv, r)
		case rmeta.Char:
			err = readI8(recv, r)
		case rmeta.Short:
			err = readI16(recv, r)
		case rmeta.Int:
			err = readI32(recv, r)
		case rmeta.Long, rmeta.Long64:
			err = readI64(recv, r)
		case rmeta.UChar:
			err = readU8(recv, r)
		case rmeta.UShort:
			err = readU16(recv, r)
		case rmeta.UInt:
			err = readU32(recv, r)
		case rmeta.ULong, rmeta.ULong64:
			err = readU64(recv, r)
		case rmeta.Float32:
			err = readF32(recv, r)
		case rmeta.Float64:
			err = readF64(recv, r)
		case rmeta.Float16:
			err = readF16(recv, r, descr.elem)
		case rmeta.Double32:
			err = readD32(recv, r, descr.elem)
		}

		if err != nil {
			return fmt.Errorf("rdict: could not rstream data: %w", err)
		}
	}
	panic("not implemented")
}

func ptrAdjust(ptr interface{}, offsets []int) interface{} {
	recv := ptr
	for _, offset := range offsets {
		rv := reflect.ValueOf(recv).Elem()
		switch rv.Kind() {
		case reflect.Struct:
			recv = rv.Field(offset).Addr().Interface()
		case reflect.Array, reflect.Slice:
			recv = rv.Index(offset).Addr().Interface()
		default:
			continue
		}
	}
	return recv
}

type ropFunc func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error

type rstreamer struct {
	op  ropFunc
	cfg *streamerConfig
}

func (rr rstreamer) rstream(r *rbytes.RBuffer, recv interface{}) error {
	return rr.op(r, recv, rr.cfg)
}

var (
	_ rstreamOp = (*rstreamer)(nil)
)
