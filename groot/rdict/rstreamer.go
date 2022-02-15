// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

// RStreamerOf returns a read-streamer for the i-th element of the provided
// streamer info and stream kind.
func RStreamerOf(sinfo rbytes.StreamerInfo, i int, kind rbytes.StreamKind) (rbytes.RStreamer, error) {
	si, ok := sinfo.(*StreamerInfo)
	if !ok {
		return nil, fmt.Errorf("rdict: not a rdict.StreamerInfo (got=%T)", sinfo)
	}

	err := si.BuildStreamers()
	if err != nil {
		return nil, fmt.Errorf("rdict: could not build streamers: %w", err)
	}

	rops := make([]rstreamer, len(si.descr))
	switch kind {
	case rbytes.ObjectWise:
		copy(rops, si.roops)
	case rbytes.MemberWise:
		copy(rops, si.rmops)
	default:
		return nil, fmt.Errorf("rdict: invalid stream kind %v", kind)
	}

	return newRStreamerElem(i, si, kind, rops)
}

type rstreamerElem struct {
	recv interface{}
	rop  *rstreamer
	i    int // streamer-element index (or -1 for the whole StreamerInfo)
	kind rbytes.StreamKind
	si   *StreamerInfo
	se   rbytes.StreamerElement
}

func newRStreamerElem(i int, si *StreamerInfo, kind rbytes.StreamKind, rops []rstreamer) (*rstreamerElem, error) {
	return &rstreamerElem{
		recv: nil,
		rop:  &rops[i],
		i:    i,
		kind: kind,
		si:   si,
		se:   si.elems[i],
	}, nil
}

func (rr *rstreamerElem) Bind(recv interface{}) error {
	rv := reflect.ValueOf(recv)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rdict: invalid kind (got=%T, want=pointer)", recv)
	}
	rr.recv = recv
	rr.rop.cfg.offset = -1 // binding directly to 'recv'. assume no offset is to be applied
	return nil
}

func (rr *rstreamerElem) RStreamROOT(r *rbytes.RBuffer) error {
	err := rr.rop.rstream(r, rr.recv)
	if err != nil {
		var (
			name  = rr.si.Name()
			ename = rr.se.TypeName()
		)
		return fmt.Errorf(
			"rdict: could not read element %d (type=%q) from %q: %w",
			rr.i, ename, name, err,
		)
	}
	return nil
}

func (rr *rstreamerElem) Count(f func() int) error {
	rr.rop.cfg.count = f
	return nil
}

var (
	_ rbytes.RStreamer = (*rstreamerElem)(nil)
	_ rbytes.Binder    = (*rstreamerElem)(nil)
	_ rbytes.Counter   = (*rstreamerElem)(nil)
)

type rstreamOp interface {
	rstream(r *rbytes.RBuffer, recv interface{}) error
}

// type rstreamBufOp interface {
// 	rstreamBuf(r *rbytes.RBuffer, recv reflect.Value, descr *elemDescr, beg, end int, n int, offset int, arrmode arrayMode) error
// }

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

func (si *StreamerInfo) makeROp(sictx rbytes.StreamerInfoContext, i int, descr elemDescr) rstreamer {
	cfg := &streamerConfig{si, i, &descr, descr.offset, 0, nil}
	switch descr.otype {
	case rmeta.Base:
		var (
			se       = descr.elem.(*StreamerBase)
			typename = se.Name()
			typevers = se.vbase
			rop      = ropFrom(sictx, typename, int16(typevers), rmeta.Base, nil)
		)
		return rstreamer{rop, cfg}

	case rmeta.Bool:
		return rstreamer{rstreamBool, cfg}
	case rmeta.Char:
		return rstreamer{rstreamI8, cfg}
	case rmeta.Short:
		return rstreamer{rstreamI16, cfg}
	case rmeta.Int:
		return rstreamer{rstreamI32, cfg}
	case rmeta.Long, rmeta.Long64:
		return rstreamer{rstreamI64, cfg}
	case rmeta.UChar:
		return rstreamer{rstreamU8, cfg}
	case rmeta.UShort:
		return rstreamer{rstreamU16, cfg}
	case rmeta.UInt:
		return rstreamer{rstreamU32, cfg}
	case rmeta.ULong, rmeta.ULong64:
		return rstreamer{rstreamU64, cfg}
	case rmeta.Float32:
		return rstreamer{rstreamF32, cfg}
	case rmeta.Float64:
		return rstreamer{rstreamF64, cfg}
	case rmeta.Bits:
		return rstreamer{rstreamBits, cfg}
	case rmeta.Float16:
		return rstreamer{rstreamF16(descr.elem), cfg}
	case rmeta.Double32:
		return rstreamer{rstreamD32(descr.elem), cfg}

	case rmeta.Counter:
		se := descr.elem.(*StreamerBasicType)
		switch se.esize {
		case 4:
			return rstreamer{rstreamI32, cfg}
		case 8:
			return rstreamer{rstreamI64, cfg}
		default:
			panic(fmt.Errorf("rdict: invalid counter size (%d) in %#v", se.esize, se))
		}

	case rmeta.CharStar:
		return rstreamer{rstreamTString, cfg}

	case rmeta.TNamed:
		return rstreamer{rstreamTNamed, cfg}
	case rmeta.TObject:
		return rstreamer{rstreamTObject, cfg}
	case rmeta.TString, rmeta.STLstring:
		return rstreamer{rstreamTString, cfg}

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
				panic("rdict: rmeta.STL (w/ old-class != new-class) not implemented")
			default:
				switch se := se.(type) {
				default:
					panic(fmt.Errorf("rdict: invalid streamer element: %#v", se))

				case *StreamerSTLstring:
					return rstreamer{
						rstreamType("string", rstreamStdString),
						cfg,
					}

				case *StreamerSTL:
					switch se.STLType() {
					case rmeta.STLvector, rmeta.STLlist, rmeta.STLdeque:
						var (
							ct       = se.ContainedType()
							typename = se.TypeName()
							enames   = rmeta.CxxTemplateFrom(typename).Args
							rop      = ropFrom(sictx, enames[0], -1, ct, &descr)
						)
						return rstreamer{
							rstreamType(typename, rstreamStdSlice(typename, rop)),
							cfg,
						}

					case rmeta.STLset, rmeta.STLmultiset, rmeta.STLunorderedset, rmeta.STLunorderedmultiset:
						var (
							ct       = se.ContainedType()
							typename = se.TypeName()
							enames   = rmeta.CxxTemplateFrom(typename).Args
							rop      = ropFrom(sictx, enames[0], -1, ct, &descr)
						)
						return rstreamer{
							rstreamType(typename, rstreamStdSet(typename, rop)),
							cfg,
						}

					case rmeta.STLmap, rmeta.STLunorderedmap, rmeta.STLmultimap, rmeta.STLunorderedmultimap:
						var (
							ct     = se.ContainedType()
							enames = rmeta.CxxTemplateFrom(se.TypeName()).Args
							kname  = enames[0]
							vname  = enames[1]
						)

						krop := ropFrom(sictx, kname, -1, ct, &descr)
						vrop := ropFrom(sictx, vname, -1, ct, &descr)
						return rstreamer{
							rstreamStdMap(kname, vname, krop, vrop),
							cfg,
						}

					case rmeta.STLbitset:
						var (
							typename = se.TypeName()
							enames   = rmeta.CxxTemplateFrom(typename).Args
							n, err   = strconv.Atoi(enames[0])
						)

						if err != nil {
							panic(fmt.Errorf("rdict: invalid STL bitset argument (type=%q): %+v", typename, err))
						}
						return rstreamer{
							rstreamType(typename, rstreamStdBitset(typename, n)),
							cfg,
						}

					default:
						panic(fmt.Errorf("rdict: STL container type=%v not handled", se.STLType()))
					}
				}
			}
		default:
			panic("rdict: rmeta.STL (w/ array-len > 1) not implemented")
			//			switch {
			//			case newClass != oldClass:
			//				panic("not implemented")
			//			default:
			//				return rstreamer{
			//					rstreamSTL(rstreamSTLArrayMbrWise, rstreamSTLObjWise, descr.oclass),
			//					&streamerConfig{si, i, &descr, descr.offset, se.ArrayLen()},
			//				}
			//			}
		}

	case rmeta.Conv + rmeta.Bool:
		return rstreamer{rstreamCnv(descr.ntype, rstreamBool), cfg}

	case rmeta.Conv + rmeta.Char:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI8), cfg}

	case rmeta.Conv + rmeta.Short:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI16), cfg}

	case rmeta.Conv + rmeta.Int:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI32), cfg}

	case rmeta.Conv + rmeta.Long, rmeta.Conv + rmeta.Long64:
		return rstreamer{rstreamCnv(descr.ntype, rstreamI64), cfg}

	case rmeta.Conv + rmeta.UChar:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU8), cfg}

	case rmeta.Conv + rmeta.UShort:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU16), cfg}

	case rmeta.Conv + rmeta.UInt:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU32), cfg}

	case rmeta.Conv + rmeta.ULong, rmeta.Conv + rmeta.ULong64:
		return rstreamer{rstreamCnv(descr.ntype, rstreamU64), cfg}

	case rmeta.Conv + rmeta.Float32:
		return rstreamer{rstreamCnv(descr.ntype, rstreamF32), cfg}

	case rmeta.Conv + rmeta.Float64:
		return rstreamer{rstreamCnv(descr.ntype, rstreamF64), cfg}

	case rmeta.Conv + rmeta.Bits:
		return rstreamer{rstreamCnv(descr.ntype, rstreamBits), cfg}

	case rmeta.Conv + rmeta.Float16:
		return rstreamer{rstreamCnv(descr.ntype, rstreamF16(descr.elem)), cfg}

	case rmeta.Conv + rmeta.Double32:
		return rstreamer{rstreamCnv(descr.ntype, rstreamD32(descr.elem)), cfg}

		// fixed-size arrays of basic types: [32]int

	case rmeta.OffsetL + rmeta.Bool:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamBool), cfg}

	case rmeta.OffsetL + rmeta.Char:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamI8), cfg}

	case rmeta.OffsetL + rmeta.Short:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamI16), cfg}

	case rmeta.OffsetL + rmeta.Int:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamI32), cfg}

	case rmeta.OffsetL + rmeta.Long, rmeta.OffsetL + rmeta.Long64:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamI64), cfg}

	case rmeta.OffsetL + rmeta.UChar:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamU8), cfg}

	case rmeta.OffsetL + rmeta.UShort:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamU16), cfg}

	case rmeta.OffsetL + rmeta.UInt:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamU32), cfg}

	case rmeta.OffsetL + rmeta.ULong, rmeta.OffsetL + rmeta.ULong64:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamU64), cfg}

	case rmeta.OffsetL + rmeta.Float32:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamF32), cfg}

	case rmeta.OffsetL + rmeta.Float64:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamF64), cfg}

	case rmeta.OffsetL + rmeta.Float16:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{
			rstreamBasicArray(cfg.length, rstreamF16(descr.elem)), // FIXME(sbinet): ROOT uses rstreamCnv here.
			cfg,
		}

	case rmeta.OffsetL + rmeta.Double32:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{
			rstreamBasicArray(cfg.length, rstreamD32(descr.elem)), // FIXME(sbinet): ROOT uses rstreamCnv here.
			cfg,
		}

	case rmeta.OffsetL + rmeta.CharStar:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamTString), cfg}

	case rmeta.OffsetL + rmeta.TString:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamTString), cfg}

	case rmeta.OffsetL + rmeta.TObject:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamTObject), cfg}

	case rmeta.OffsetL + rmeta.TNamed:
		cfg.length = descr.elem.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rstreamTNamed), cfg}

	case rmeta.OffsetL + rmeta.Any,
		rmeta.OffsetL + rmeta.Object:
		var (
			se       = descr.elem
			typename = se.TypeName()
			rop      = ropFrom(sictx, typename, -1, 0, nil)
		)
		cfg.length = se.ArrayLen()
		return rstreamer{rstreamBasicArray(cfg.length, rop), cfg}

		// var-size arrays of basic types: [n]int

	case rmeta.OffsetP + rmeta.Bool:
		return rstreamer{rstreamBools, cfg}

	case rmeta.OffsetP + rmeta.Char:
		return rstreamer{rstreamI8s, cfg}

	case rmeta.OffsetP + rmeta.Short:
		return rstreamer{rstreamI16s, cfg}

	case rmeta.OffsetP + rmeta.Int:
		return rstreamer{rstreamI32s, cfg}

	case rmeta.OffsetP + rmeta.Long, rmeta.OffsetP + rmeta.Long64:
		return rstreamer{rstreamI64s, cfg}

	case rmeta.OffsetP + rmeta.UChar:
		return rstreamer{rstreamU8s, cfg}

	case rmeta.OffsetP + rmeta.UShort:
		return rstreamer{rstreamU16s, cfg}

	case rmeta.OffsetP + rmeta.UInt:
		return rstreamer{rstreamU32s, cfg}

	case rmeta.OffsetP + rmeta.ULong, rmeta.OffsetP + rmeta.ULong64:
		return rstreamer{rstreamU64s, cfg}

	case rmeta.OffsetP + rmeta.Float32:
		return rstreamer{rstreamF32s, cfg}

	case rmeta.OffsetP + rmeta.Float64:
		return rstreamer{rstreamF64s, cfg}

	case rmeta.OffsetP + rmeta.Float16:
		return rstreamer{rstreamF16s, cfg}

	case rmeta.OffsetP + rmeta.Double32:
		return rstreamer{rstreamD32s, cfg}

	case rmeta.OffsetP + rmeta.CharStar:
		return rstreamer{rstreamStrs, cfg}

	case rmeta.Streamer:
		switch se := descr.elem.(type) {
		default:
			panic(fmt.Errorf("rdict: invalid streamer element: %#v", se))

		case *StreamerSTLstring:
			return rstreamer{
				rstreamType("string", rstreamStdString),
				cfg,
			}

		case *StreamerSTL:
			switch se.STLType() {
			case rmeta.STLvector, rmeta.STLlist, rmeta.STLdeque:
				var (
					ct       = se.ContainedType()
					typename = se.TypeName()
					enames   = rmeta.CxxTemplateFrom(typename).Args
					rop      = ropFrom(sictx, enames[0], -1, ct, &descr)
				)
				return rstreamer{
					rstreamType(typename, rstreamStdSlice(typename, rop)),
					cfg,
				}

			case rmeta.STLset, rmeta.STLmultiset, rmeta.STLunorderedset, rmeta.STLunorderedmultiset:
				var (
					ct       = se.ContainedType()
					typename = se.TypeName()
					enames   = rmeta.CxxTemplateFrom(typename).Args
					rop      = ropFrom(sictx, enames[0], -1, ct, &descr)
				)
				return rstreamer{
					rstreamType(typename, rstreamStdSet(typename, rop)),
					cfg,
				}

			case rmeta.STLmap, rmeta.STLmultimap, rmeta.STLunorderedmap, rmeta.STLunorderedmultimap:
				var (
					ct     = se.ContainedType()
					enames = rmeta.CxxTemplateFrom(se.TypeName()).Args
					kname  = enames[0]
					vname  = enames[1]
				)

				krop := ropFrom(sictx, kname, -1, ct, &descr)
				vrop := ropFrom(sictx, vname, -1, ct, &descr)
				return rstreamer{
					rstreamStdMap(kname, vname, krop, vrop),
					cfg,
				}

			case rmeta.STLbitset:
				var (
					typename = se.TypeName()
					enames   = rmeta.CxxTemplateFrom(typename).Args
					n, err   = strconv.Atoi(enames[0])
				)

				if err != nil {
					panic(fmt.Errorf("rdict: invalid STL bitset argument (type=%q): %+v", typename, err))
				}
				return rstreamer{
					rstreamType(typename, rstreamStdBitset(typename, n)),
					cfg,
				}

			default:
				panic(fmt.Errorf("rdict: STL container type=%v not handled", se.STLType()))
			}
		}

	case rmeta.Any, rmeta.Object:
		var (
			se       = descr.elem
			typename = se.TypeName()
			rop      = ropFrom(sictx, typename, -1, 0, nil)
		)
		return rstreamer{rop, cfg}

	case rmeta.AnyP, rmeta.Anyp:
		var (
			se       = descr.elem
			typename = strings.TrimRight(se.TypeName(), "*") // FIXME(sbinet): handle T** ?
			rop      = ropFrom(sictx, typename, -1, 0, nil)
		)
		return rstreamer{
			rstreamAnyPtr(rop),
			cfg,
		}

	case rmeta.ObjectP, rmeta.Objectp:
		var (
			se       = descr.elem
			typename = strings.TrimRight(se.TypeName(), "*") // FIXME(sbinet): handle T** ?
			rop      = ropFrom(sictx, typename, -1, 0, nil)
		)
		return rstreamer{
			rstreamObjPtr(rop),
			cfg,
		}

	case rmeta.StreamLoop:
		var (
			se       = descr.elem.(*StreamerLoop)
			typename = strings.TrimRight(se.TypeName(), "*") // FIXME(sbinet): handle T** ?
			rop      = ropFrom(sictx, typename, -1, 0, nil)
		)
		return rstreamer{
			rstreamBasicSlice(rop),
			cfg,
		}

	default:
		panic(fmt.Errorf("not implemented k=%d (%v)", descr.otype, descr.otype))
		// return rstreamer{rstreamGeneric, &streamerConfig{si, i, &descr, descr.offset, 0}}
	}
}

func rstreamSI(si *StreamerInfo) ropFunc {
	typename := si.Name()
	switch {
	case typename == "TObject":
		return rstreamTObject
	case typename == "TNamed":
		return rstreamTNamed
	case typename == "TString":
		return rstreamTString
	case rtypes.Factory.HasKey(typename):
		obj := rtypes.Factory.Get(typename)().Interface()
		_, ok := obj.(rbytes.Unmarshaler)
		if ok {
			return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
				obj := cfg.adjust(recv).(rbytes.Unmarshaler)
				return obj.UnmarshalROOT(r)
			}
		}
	}
	return rstreamCat(typename, int16(si.ClassVersion()), si.roops)
}

func rstreamObjPtr(rop ropFunc) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		obj := r.ReadObjectAny()
		if r.Err() != nil {
			return r.Err()
		}
		rv := reflect.ValueOf(cfg.adjust(recv))
		if obj == nil {
			if !rv.Elem().IsNil() {
				rv.Elem().Set(reflect.Value{})
			}
			return nil
		}
		rv.Elem().Set(reflect.ValueOf(obj))
		return nil
	}
}

func rstreamAnyPtr(rop ropFunc) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		obj := r.ReadObjectAny()
		if r.Err() != nil {
			return r.Err()
		}
		rv := reflect.ValueOf(cfg.adjust(recv))
		if obj == nil {
			if !rv.Elem().IsNil() {
				rv.Elem().Set(reflect.Value{})
			}
			return nil
		}
		rv.Elem().Set(reflect.ValueOf(obj))
		return nil
	}
}

func rstreamBool(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*bool)) = r.ReadBool()
	return r.Err()
}

func rstreamI8(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*int8)) = r.ReadI8()
	return r.Err()
}

func rstreamI16(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*int16)) = r.ReadI16()
	return r.Err()
}

func rstreamI32(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*int32)) = r.ReadI32()
	return r.Err()
}

func rstreamI64(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*int64)) = r.ReadI64()
	return r.Err()
}

func rstreamU8(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*uint8)) = r.ReadU8()
	return r.Err()
}

func rstreamU16(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*uint16)) = r.ReadU16()
	return r.Err()
}

func rstreamU32(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*uint32)) = r.ReadU32()
	return r.Err()
}

func rstreamU64(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*uint64)) = r.ReadU64()
	return r.Err()
}

func rstreamF32(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*float32)) = r.ReadF32()
	return r.Err()
}

func rstreamF64(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*float64)) = r.ReadF64()
	return r.Err()
}

func rstreamBits(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*uint32)) = r.ReadU32()
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
	*(cfg.adjust(recv).(*string)) = r.ReadString()
	return r.Err()
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
					"rdict: could not rstream array element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, cfg.si.Name(), err,
				)
			}
		}
		return nil
	}
}

func rstreamBasicSlice(sli ropFunc) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		_ = r.ReadI8() // is-array
		n := int(reflect.ValueOf(recv).Elem().FieldByIndex(cfg.descr.method).Int())
		rv := reflect.ValueOf(cfg.adjust(recv)).Elem()
		if nn := rv.Len(); nn < n {
			rv.Set(reflect.AppendSlice(rv, reflect.MakeSlice(rv.Type(), n-nn, n-nn)))
		}
		for i := 0; i < n; i++ {
			err := sli(r, rv.Index(i).Addr().Interface(), nil)
			if err != nil {
				return fmt.Errorf(
					"rdict: could not rstream slice element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, cfg.si.Name(), err,
				)
			}
		}
		return nil
	}
}

func rstreamHeader(r *rbytes.RBuffer, typename string) rbytes.Header {
	if _, ok := rmeta.CxxBuiltins[typename]; ok && typename != "string" {
		return rbytes.Header{Pos: -1}
	}
	return r.ReadHeader(typename)
}

func rcheckHeader(r *rbytes.RBuffer, hdr rbytes.Header) error {
	if hdr.Pos < 0 {
		return nil
	}
	r.CheckHeader(hdr)
	return r.Err()
}

func rstreamType(typename string, rop ropFunc) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		hdr := r.ReadHeader(typename)
		err := rop(r, recv, cfg)
		if err != nil {
			return fmt.Errorf(
				"rdict: could not read (type=%q, vers=%d): %w",
				hdr.Name, hdr.Vers, err,
			)
		}
		r.CheckHeader(hdr)
		return r.Err()
	}
}

func rstreamStdSlice(typename string, rop ropFunc) ropFunc {
	//	const typevers = 1
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		// FIXME(sbinet): use typevers to infer obj-/mbr-wise reading.
		n := int(r.ReadI32())
		rv := reflect.ValueOf(cfg.adjust(recv)).Elem()
		if nn := rv.Len(); nn < n {
			rv.Set(reflect.AppendSlice(rv, reflect.MakeSlice(rv.Type(), n-nn, n-nn)))
		}
		rv.SetLen(n)
		for i := 0; i < n; i++ {
			err := rop(r, rv.Index(i).Addr().Interface(), nil)
			if err != nil {
				return fmt.Errorf(
					"rdict: could not rstream element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, cfg.si.Name(), err,
				)
			}
		}
		return r.Err()
	}
}

func rstreamStdSet(typename string, rop ropFunc) ropFunc {
	//	const typevers = 1

	// FIXME(sbinet): add special handling for std::set-like types
	// the correct equivalent Go-type of std::set<T> is map[T]struct{}
	// (or, when availaible, std.Set[T])
	return rstreamStdSlice(typename, rop)
}

func rstreamStdMap(kname, vname string, krop, vrop ropFunc) ropFunc {
	typename := fmt.Sprintf("map<%s,%s>", kname, vname)
	if strings.HasSuffix(vname, ">") {
		typename = fmt.Sprintf("map<%s,%s >", kname, vname)
	}
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		//typevers = int16(cfg.si.ClassVersion())
		hdr := r.ReadHeader(typename)
		mbrwise := hdr.Vers&rbytes.StreamedMemberWise != 0
		// if mbrwise {
		// 	vers &= ^rbytes.StreamedMemberWise
		// }

		if mbrwise {
			clvers := r.ReadI16()
			switch {
			case clvers == 1:
				// TODO
			case clvers <= 0:
				/*chksum*/ _ = r.ReadU32()
			}
		}

		n := int(r.ReadI32())
		rv := reflect.ValueOf(cfg.adjust(recv)).Elem()
		keyT := reflect.SliceOf(rv.Type().Key())
		valT := reflect.SliceOf(rv.Type().Elem())
		keys := reflect.New(keyT).Elem()
		keys.Set(reflect.AppendSlice(keys, reflect.MakeSlice(keyT, n, n)))
		if n > 0 {
			hdr := rstreamHeader(r, kname)
			for i := 0; i < n; i++ {
				err := krop(r, keys.Index(i).Addr().Interface(), nil)
				if err != nil {
					return fmt.Errorf(
						"rdict: could not rstream key-element %s[%d] of %s: %w",
						kname, i, cfg.si.Name(), err,
					)
				}
			}
			err := rcheckHeader(r, hdr)
			if err != nil {
				return err
			}
		}

		vals := reflect.New(valT).Elem()
		vals.Set(reflect.AppendSlice(vals, reflect.MakeSlice(valT, n, n)))
		if n > 0 {
			hdr := rstreamHeader(r, vname)
			for i := 0; i < n; i++ {
				err := vrop(r, vals.Index(i).Addr().Interface(), nil)
				if err != nil {
					return fmt.Errorf(
						"rdict: could not rstream val-element %s[%d] of %s: %w",
						vname, i, cfg.si.Name(), err,
					)
				}
			}
			err := rcheckHeader(r, hdr)
			if err != nil {
				return err
			}
		}

		if rv.IsNil() {
			rv.Set(reflect.MakeMapWithSize(rv.Type(), n))
		}
		for i := 0; i < n; i++ {
			rv.SetMapIndex(keys.Index(i), vals.Index(i))
		}

		r.CheckHeader(hdr)
		return r.Err()
	}
}

func rstreamStdBitset(typename string, n int) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		var (
			bits = int(r.ReadI32())
			sli  = cfg.adjust(recv).(*[]uint8)
		)
		*sli = rbytes.ResizeU8(*sli, bits)
		r.ReadStdBitset(*sli)
		return r.Err()
	}
}

func rstreamBools(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]bool)
	)
	*sli = rbytes.ResizeBool(*sli, n)
	r.ReadArrayBool(*sli)
	return r.Err()
}

func rstreamU8s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint8)
	)
	*sli = rbytes.ResizeU8(*sli, n)
	r.ReadArrayU8(*sli)
	return r.Err()
}

func rstreamU16s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint16)
	)
	*sli = rbytes.ResizeU16(*sli, n)
	r.ReadArrayU16(*sli)
	return r.Err()
}

func rstreamU32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint32)
	)
	*sli = rbytes.ResizeU32(*sli, n)
	r.ReadArrayU32(*sli)
	return r.Err()
}

func rstreamU64s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]uint64)
	)
	*sli = rbytes.ResizeU64(*sli, n)
	r.ReadArrayU64(*sli)
	return r.Err()
}

func rstreamI8s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int8)
	)
	*sli = rbytes.ResizeI8(*sli, n)
	r.ReadArrayI8(*sli)
	return r.Err()
}

func rstreamI16s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int16)
	)
	*sli = rbytes.ResizeI16(*sli, n)
	r.ReadArrayI16(*sli)
	return r.Err()
}

func rstreamI32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int32)
	)
	*sli = rbytes.ResizeI32(*sli, n)
	r.ReadArrayI32(*sli)
	return r.Err()
}

func rstreamI64s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]int64)
	)
	*sli = rbytes.ResizeI64(*sli, n)
	r.ReadArrayI64(*sli)
	return r.Err()
}

func rstreamF32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]float32)
	)
	*sli = rbytes.ResizeF32(*sli, n)
	r.ReadArrayF32(*sli)
	return r.Err()
}

func rstreamF64s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]float64)
	)
	*sli = rbytes.ResizeF64(*sli, n)
	r.ReadArrayF64(*sli)
	return r.Err()
}

func rstreamF16s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]root.Float16)
	)
	*sli = rbytes.ResizeF16(*sli, n)
	r.ReadArrayF16(*sli, cfg.descr.elem)
	return r.Err()
}

func rstreamD32s(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]root.Double32)
	)
	*sli = rbytes.ResizeD32(*sli, n)
	r.ReadArrayD32(*sli, cfg.descr.elem)
	return r.Err()
}

func rstreamStrs(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	var (
		_   = r.ReadI8() // is-array
		n   = cfg.counter(recv)
		sli = cfg.adjust(recv).(*[]string)
	)
	*sli = rbytes.ResizeStr(*sli, n)
	r.ReadArrayString(*sli)
	return r.Err()
}

func rstreamCat(typename string, typevers int16, rops []rstreamer) ropFunc {
	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
		hdr := r.ReadHeader(typename)
		if hdr.Vers != typevers {
			r.SetErr(fmt.Errorf(
				"rdict: inconsistent ROOT version type=%q (got=%d, want=%d)",
				hdr.Name, hdr.Vers, typevers,
			))
			return r.Err()
		}

		recv = cfg.adjust(recv)
		for i, rop := range rops {
			err := rop.rstream(r, recv)
			if err != nil {
				return fmt.Errorf(
					"rdict: could not rstream element %d (%s) of %s: %w",
					i, rop.cfg.descr.elem.Name(), cfg.si.Name(), err,
				)
			}
		}
		r.CheckHeader(hdr)
		return r.Err()
	}
}

func rstreamStdString(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
	*(cfg.adjust(recv).(*string)) = r.ReadString()
	return r.Err()

}

// func rstreamGeneric(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
// 	const (
// 		beg     = 0
// 		end     = 1
// 		n       = 1
// 		arrmode = 2
// 	)
// 	return cfg.si.rstream(r, recv, cfg.descr, beg, end, n, cfg.offset, arrmode)
// }

// type rmbrwiseSTLFunc func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16) error
// type robjwiseSTLFunc func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16, start int32) error
//
// func rstreamSTL(mbrwise rmbrwiseSTLFunc, objwise robjwiseSTLFunc, oclass string) ropFunc {
// 	return func(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig) error {
// 		err := r.Err()
// 		if err != nil {
// 			return err
// 		}
//
// 		hdr := r.ReadHeader(oclass)
// 		switch {
// 		case hdr.Vers&rbytes.StreamedMemberWise != 0:
// 			err = mbrwise(r, cfg.adjust(recv), cfg, hdr.Vers)
// 		default:
// 			err = objwise(r, cfg.adjust(recv), cfg, hdr.Vers, int32(beg))
// 		}
//
// 		r.CheckHeader(hdr)
// 		return err
// 	}
// }
//
// func rstreamSTLArrayMbrWise(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16) error {
// 	panic("not implemented")
// }
//
// func rstreamSTLObjWise(r *rbytes.RBuffer, recv interface{}, cfg *streamerConfig, vers int16, start int32) error {
// 	panic("not implemented")
// }

// func (si *StreamerInfo) rstream(r *rbytes.RBuffer, recv interface{}, descr *elemDescr, beg, end, n, offset int, mode arrayMode) error {
// 	needIncr := mode&2 == 0
// 	mode &= ^2
//
// 	if needIncr {
// 		panic("not implemented")
// 	}
// 	const kHaveLoop = 1024
// 	var typeOffset rmeta.Enum
// 	if mode == 0 {
// 		typeOffset = kHaveLoop
// 	}
//
// 	var (
// 		ioffset = []int{-1, offset}
// 		err     error
// 	)
// 	for i := beg; i < end; i++ {
// 		ioffset[0] = descr.offset
//
// 		var (
// 			recv  = ptrAdjust(recv, ioffset)
// 			etype = descr.otype
// 		)
//
// 		switch etype + typeOffset {
// 		case rmeta.Bool:
// 			*(recv.(*bool)) = r.ReadBool()
// 		case rmeta.Char:
// 			*(recv.(*int8)) = r.ReadI8()
// 		case rmeta.Short:
// 			*(recv.(*int16)) = r.ReadI16()
// 		case rmeta.Int:
// 			*(recv.(*int32)) = r.ReadI32()
// 		case rmeta.Long, rmeta.Long64:
// 			*(recv.(*int64)) = r.ReadI64()
// 		case rmeta.UChar:
// 			*(recv.(*uint8)) = r.ReadU8()
// 		case rmeta.UShort:
// 			*(recv.(*uint16)) = r.ReadU16()
// 		case rmeta.UInt:
// 			*(recv.(*uint32)) = r.ReadU32()
// 		case rmeta.ULong, rmeta.ULong64:
// 			*(recv.(*uint64)) = r.ReadU64()
// 		case rmeta.Float32:
// 			*(recv.(*float32)) = r.ReadF32()
// 		case rmeta.Float64:
// 			*(recv.(*float64)) = r.ReadF64()
// 		case rmeta.Float16:
// 			*(recv.(*root.Float16)) = r.ReadF16(descr.elem)
// 		case rmeta.Double32:
// 			*(recv.(*root.Double32)) = r.ReadD32(descr.elem)
// 		}
// 		err = r.Err()
//
// 		if err != nil {
// 			return fmt.Errorf("rdict: could not rstream data: %w", err)
// 		}
// 	}
// 	panic("not implemented")
// }

// func ptrAdjust(ptr interface{}, offsets []int) interface{} {
// 	recv := ptr
// 	for _, offset := range offsets {
// 		rv := reflect.ValueOf(recv).Elem()
// 		switch rv.Kind() {
// 		case reflect.Struct:
// 			recv = rv.Field(offset).Addr().Interface()
// 		case reflect.Array, reflect.Slice:
// 			recv = rv.Index(offset).Addr().Interface()
// 		default:
// 			continue
// 		}
// 	}
// 	return recv
// }

func ropFuncFor(e rmeta.Enum, descr *elemDescr) ropFunc {
	switch e {
	case rmeta.Bool:
		return rstreamBool
	case rmeta.Bits:
		return rstreamBits
	case rmeta.Int8:
		return rstreamI8
	case rmeta.Int16:
		return rstreamI16
	case rmeta.Int32:
		return rstreamI32
	case rmeta.Int64, rmeta.Long64:
		return rstreamI64
	case rmeta.Uint8:
		return rstreamU8
	case rmeta.Uint16:
		return rstreamU16
	case rmeta.Uint32:
		return rstreamU32
	case rmeta.Uint64, rmeta.ULong64:
		return rstreamU64
	case rmeta.Float32:
		return rstreamF32
	case rmeta.Float64:
		return rstreamF64
	case rmeta.Float16:
		return rstreamF16(descr.elem)
	case rmeta.Double32:
		return rstreamD32(descr.elem)
	case rmeta.TString, rmeta.CharStar:
		return rstreamTString
	case rmeta.STLstring:
		return rstreamStdString
	case rmeta.TObject:
		return rstreamTObject
	case rmeta.TNamed:
		return rstreamTNamed
	default:
		return nil
	}
}

func ropFrom(sictx rbytes.StreamerInfoContext, typename string, typevers int16, enum rmeta.Enum, descr *elemDescr) ropFunc {
	e, ok := rmeta.TypeName2Enum(typename)
	if ok {
		rop := ropFuncFor(e, descr)
		if rop != nil {
			return rop
		}
	}

	rop := ropFuncFor(enum, descr)
	if rop != nil {
		return rop
	}

	switch {
	case hasStdPrefix(typename, "vector", "list", "deque"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		rop := ropFrom(sictx, enames[0], -1, 0, nil)
		return rstreamStdSlice(typename, rop)

	case hasStdPrefix(typename, "set", "multiset", "unordered_set", "unordered_multiset"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		rop := ropFrom(sictx, enames[0], -1, 0, nil)
		return rstreamStdSet(typename, rop)

	case hasStdPrefix(typename, "map", "multimap", "unordered_map", "unordered_multimap"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		kname := enames[0]
		vname := enames[1]

		krop := ropFrom(sictx, kname, -1, 0, nil)
		vrop := ropFrom(sictx, vname, -1, 0, nil)
		return rstreamStdMap(kname, vname, krop, vrop)

	case hasStdPrefix(typename, "bitset"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		n, err := strconv.Atoi(enames[0])
		if err != nil {
			panic(fmt.Errorf("rdict: invalid STL bitset argument (type=%q): %+v", typename, err))
		}
		return rstreamStdBitset(typename, n)
	}

	osi, err := sictx.StreamerInfo(typename, int(typevers))
	if err != nil {
		panic(fmt.Errorf("rdict: could not find streamer info for %q (version=%d): %w", typename, typevers, err))
	}
	esi := osi.(*StreamerInfo)

	err = esi.BuildStreamers()
	if err != nil {
		panic(fmt.Errorf("rdict: could not build streamers for %q (version=%d): %w", typename, typevers, err))
	}

	rop = rstreamSI(esi)
	return rop
}
