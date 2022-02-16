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
	"go-hep.org/x/hep/groot/rvers"
)

// WStreamerOf returns a write-streamer for the i-th element of the provided
// streamer info and stream kind.
func WStreamerOf(sinfo rbytes.StreamerInfo, i int, kind rbytes.StreamKind) (rbytes.WStreamer, error) {
	si, ok := sinfo.(*StreamerInfo)
	if !ok {
		return nil, fmt.Errorf("rdict: not a rdict.StreamerInfo (got=%T)", sinfo)
	}

	err := si.BuildStreamers()
	if err != nil {
		return nil, fmt.Errorf("rdict: could not build streamers: %w", err)
	}

	switch kind {
	case rbytes.ObjectWise:
		return newWStreamer(i, si, kind, si.woops)
	case rbytes.MemberWise:
		return newWStreamer(i, si, kind, si.wmops)
	default:
		return nil, fmt.Errorf("rdict: invalid stream kind %v", kind)
	}
}

type wstreamerElem struct {
	recv interface{}
	wop  *wstreamer
	i    int // streamer-element index (or -1 for the whole StreamerInfo)
	kind rbytes.StreamKind
	si   *StreamerInfo
	se   rbytes.StreamerElement
}

func newWStreamer(i int, si *StreamerInfo, kind rbytes.StreamKind, wops []wstreamer) (*wstreamerElem, error) {
	return &wstreamerElem{
		recv: nil,
		wop:  &wops[i],
		i:    i,
		kind: kind,
		si:   si,
		se:   si.elems[i],
	}, nil
}

func (ww *wstreamerElem) Bind(recv interface{}) error {
	rv := reflect.ValueOf(recv)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rdict: invalid kind (got=%T, want=pointer)", recv)
	}
	ww.recv = recv
	ww.wop.cfg.offset = -1 // binding directly to 'recv'. assume no offset is to be applied
	return nil
}

func (ww *wstreamerElem) WStreamROOT(w *rbytes.WBuffer) error {
	_, err := ww.wop.wstream(w, ww.recv)
	if err != nil {
		var (
			name  = ww.si.Name()
			ename = ww.se.TypeName()
		)
		return fmt.Errorf("rdict: could not write element %d (type=%q) from %q: %w",
			ww.i, ename, name, err,
		)
	}

	return nil
}

var (
	_ rbytes.WStreamer = (*wstreamerElem)(nil)
	_ rbytes.Binder    = (*wstreamerElem)(nil)
)

type wstreamOp interface {
	wstream(w *rbytes.WBuffer, recv interface{}) (int, error)
}

// type wstreamBufOp interface {
// 	wstreamBuf(w *rbytes.WBuffer, recv reflect.Value, descr *elemDescr, beg, end int, n int, offset int, arrmode arrayMode) (int, error)
// }

type wopFunc func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error)

type wstreamer struct {
	op  wopFunc
	cfg *streamerConfig
}

func (ww wstreamer) wstream(w *rbytes.WBuffer, recv interface{}) (int, error) {
	return ww.op(w, recv, ww.cfg)
}

var (
	_ wstreamOp = (*wstreamer)(nil)
)

func (si *StreamerInfo) makeWOp(sictx rbytes.StreamerInfoContext, i int, descr elemDescr) wstreamer {
	cfg := &streamerConfig{si, i, &descr, descr.offset, 0, nil}
	switch descr.otype {
	case rmeta.Base:
		var (
			se       = descr.elem.(*StreamerBase)
			typename = se.Name()
			typevers = se.vbase
			wop, _   = wopFrom(sictx, typename, int16(typevers), 0, nil)
		)
		return wstreamer{wop, cfg}

	case rmeta.Bool:
		return wstreamer{wstreamBool, cfg}
	case rmeta.Char:
		return wstreamer{wstreamI8, cfg}
	case rmeta.Short:
		return wstreamer{wstreamI16, cfg}
	case rmeta.Int:
		return wstreamer{wstreamI32, cfg}
	case rmeta.Long, rmeta.Long64:
		return wstreamer{wstreamI64, cfg}
	case rmeta.UChar:
		return wstreamer{wstreamU8, cfg}
	case rmeta.UShort:
		return wstreamer{wstreamU16, cfg}
	case rmeta.UInt:
		return wstreamer{wstreamU32, cfg}
	case rmeta.ULong, rmeta.ULong64:
		return wstreamer{wstreamU64, cfg}
	case rmeta.Float32:
		return wstreamer{wstreamF32, cfg}
	case rmeta.Float64:
		return wstreamer{wstreamF64, cfg}
	case rmeta.Bits:
		return wstreamer{wstreamBits, cfg}
	case rmeta.Float16:
		return wstreamer{wstreamF16(descr.elem), cfg}
	case rmeta.Double32:
		return wstreamer{wstreamD32(descr.elem), cfg}

	case rmeta.Counter:
		se := descr.elem.(*StreamerBasicType)
		switch se.esize {
		case 4:
			return wstreamer{wstreamI32, cfg}
		case 8:
			return wstreamer{wstreamI64, cfg}
		default:
			panic(fmt.Errorf("rdict: invalid counter size (%d) in %#v", se.esize, se))
		}

	case rmeta.CharStar:
		return wstreamer{wstreamTString, cfg}

	case rmeta.TNamed:
		return wstreamer{wstreamTNamed, cfg}
	case rmeta.TObject:
		return wstreamer{wstreamTObject, cfg}
	case rmeta.TString, rmeta.STLstring:
		return wstreamer{wstreamTString, cfg}

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
					return wstreamer{
						wstreamType("string", wstreamStdString),
						cfg,
					}

				case *StreamerSTL:
					switch se.STLType() {
					case rmeta.STLvector, rmeta.STLlist, rmeta.STLdeque:
						var (
							ct       = se.ContainedType()
							typename = se.TypeName()
							enames   = rmeta.CxxTemplateFrom(typename).Args
							wop, _   = wopFrom(sictx, enames[0], -1, ct, &descr)
						)
						return wstreamer{
							wstreamType(typename, wstreamStdSlice(typename, wop)),
							cfg,
						}

					case rmeta.STLset, rmeta.STLmultiset, rmeta.STLunorderedset, rmeta.STLunorderedmultiset:
						var (
							ct       = se.ContainedType()
							typename = se.TypeName()
							enames   = rmeta.CxxTemplateFrom(typename).Args
							wop, _   = wopFrom(sictx, enames[0], -1, ct, &descr)
						)
						return wstreamer{
							wstreamType(typename, wstreamStdSet(typename, wop)),
							cfg,
						}

					case rmeta.STLmap, rmeta.STLmultimap, rmeta.STLunorderedmap, rmeta.STLunorderedmultimap:
						var (
							ct     = se.ContainedType()
							enames = rmeta.CxxTemplateFrom(se.TypeName()).Args
							kname  = enames[0]
							vname  = enames[1]
						)

						kwop, kvers := wopFrom(sictx, kname, -1, ct, &descr)
						vwop, vvers := wopFrom(sictx, vname, -1, ct, &descr)
						return wstreamer{
							wstreamStdMap(
								kname, vname,
								kwop, vwop,
								kvers, vvers,
							),
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
						return wstreamer{
							wstreamType(typename, wstreamStdBitset(typename, n)),
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
			//				return wstreamer{
			//					wstreamSTL(wstreamSTLArrayMbrWise, wstreamSTLObjWise, descr.oclass),
			//					&wtreamerConfig{si, i, &descr, descr.offset, se.ArrayLen()},
			//				}
			//			}
		}

	// FIXME(sbinet): add rmeta.Conv handling.

	// fixed-size arrays of basic types: [32]int

	case rmeta.OffsetL + rmeta.Bool:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamBool), cfg}

	case rmeta.OffsetL + rmeta.Char:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamI8), cfg}

	case rmeta.OffsetL + rmeta.Short:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamI16), cfg}

	case rmeta.OffsetL + rmeta.Int:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamI32), cfg}

	case rmeta.OffsetL + rmeta.Long, rmeta.OffsetL + rmeta.Long64:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamI64), cfg}

	case rmeta.OffsetL + rmeta.UChar:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamU8), cfg}

	case rmeta.OffsetL + rmeta.UShort:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamU16), cfg}

	case rmeta.OffsetL + rmeta.UInt:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamU32), cfg}

	case rmeta.OffsetL + rmeta.ULong, rmeta.OffsetL + rmeta.ULong64:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamU64), cfg}

	case rmeta.OffsetL + rmeta.Float32:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamF32), cfg}

	case rmeta.OffsetL + rmeta.Float64:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamF64), cfg}

	case rmeta.OffsetL + rmeta.Float16:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{
			wstreamBasicArray(cfg.length, wstreamF16(descr.elem)), // FIXME(sbinet): ROOT uses wstreamCnv here.
			cfg,
		}

	case rmeta.OffsetL + rmeta.Double32:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{
			wstreamBasicArray(cfg.length, wstreamD32(descr.elem)), // FIXME(sbinet): ROOT uses wstreamCnv here.
			cfg,
		}

	case rmeta.OffsetL + rmeta.CharStar:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamTString), cfg}

	case rmeta.OffsetL + rmeta.TString:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamTString), cfg}

	case rmeta.OffsetL + rmeta.TObject:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamTObject), cfg}

	case rmeta.OffsetL + rmeta.TNamed:
		cfg.length = descr.elem.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wstreamTNamed), cfg}

	case rmeta.OffsetL + rmeta.Any,
		rmeta.OffsetL + rmeta.Object:
		var (
			se       = descr.elem
			typename = se.TypeName()
			wop, _   = wopFrom(sictx, typename, -1, 0, nil)
		)
		cfg.length = se.ArrayLen()
		return wstreamer{wstreamBasicArray(cfg.length, wop), cfg}

	// var-size arrays of basic types: [n]int

	case rmeta.OffsetP + rmeta.Bool:
		return wstreamer{wstreamBools, cfg}

	case rmeta.OffsetP + rmeta.Char:
		return wstreamer{wstreamI8s, cfg}

	case rmeta.OffsetP + rmeta.Short:
		return wstreamer{wstreamI16s, cfg}

	case rmeta.OffsetP + rmeta.Int:
		return wstreamer{wstreamI32s, cfg}

	case rmeta.OffsetP + rmeta.Long, rmeta.OffsetP + rmeta.Long64:
		return wstreamer{wstreamI64s, cfg}

	case rmeta.OffsetP + rmeta.UChar:
		return wstreamer{wstreamU8s, cfg}

	case rmeta.OffsetP + rmeta.UShort:
		return wstreamer{wstreamU16s, cfg}

	case rmeta.OffsetP + rmeta.UInt:
		return wstreamer{wstreamU32s, cfg}

	case rmeta.OffsetP + rmeta.ULong, rmeta.OffsetP + rmeta.ULong64:
		return wstreamer{wstreamU64s, cfg}

	case rmeta.OffsetP + rmeta.Float32:
		return wstreamer{wstreamF32s, cfg}

	case rmeta.OffsetP + rmeta.Float64:
		return wstreamer{wstreamF64s, cfg}

	case rmeta.OffsetP + rmeta.Float16:
		return wstreamer{wstreamF16s, cfg}

	case rmeta.OffsetP + rmeta.Double32:
		return wstreamer{wstreamD32s, cfg}

	case rmeta.OffsetP + rmeta.CharStar:
		return wstreamer{wstreamStrs, cfg}

	case rmeta.Streamer:
		switch se := descr.elem.(type) {
		default:
			panic(fmt.Errorf("rdict: invalid streamer element: %#v", se))

		case *StreamerSTLstring:
			return wstreamer{
				wstreamType("string", wstreamStdString),
				cfg,
			}

		case *StreamerSTL:
			switch se.STLType() {
			case rmeta.STLvector, rmeta.STLlist, rmeta.STLdeque:
				var (
					ct       = se.ContainedType()
					typename = se.TypeName()
					enames   = rmeta.CxxTemplateFrom(typename).Args
					vname    = enames[0]
					wop, _   = wopFrom(sictx, vname, -1, ct, &descr)
				)
				return wstreamer{
					wstreamType(typename, wstreamStdSlice(typename, wop)),
					cfg,
				}

			case rmeta.STLset, rmeta.STLmultiset, rmeta.STLunorderedset, rmeta.STLunorderedmultiset:
				var (
					ct       = se.ContainedType()
					typename = se.TypeName()
					enames   = rmeta.CxxTemplateFrom(typename).Args
					vname    = enames[0]
					wop, _   = wopFrom(sictx, vname, -1, ct, &descr)
				)
				return wstreamer{
					wstreamType(typename, wstreamStdSet(typename, wop)),
					cfg,
				}

			case rmeta.STLmap, rmeta.STLmultimap, rmeta.STLunorderedmap, rmeta.STLunorderedmultimap:
				var (
					ct     = se.ContainedType()
					enames = rmeta.CxxTemplateFrom(se.TypeName()).Args
					kname  = enames[0]
					vname  = enames[1]
				)

				kwop, kvers := wopFrom(sictx, kname, -1, ct, &descr)
				vwop, vvers := wopFrom(sictx, vname, -1, ct, &descr)
				return wstreamer{
					wstreamStdMap(
						kname, vname,
						kwop,
						vwop,
						kvers,
						vvers,
					),
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
				return wstreamer{
					wstreamType(typename, wstreamStdBitset(typename, n)),
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
			wop, _   = wopFrom(sictx, typename, -1, 0, nil)
		)
		return wstreamer{wop, cfg}

	case rmeta.AnyP, rmeta.Anyp:
		var (
			se       = descr.elem
			typename = strings.TrimRight(se.TypeName(), "*") // FIXME(sbinet): handle T** ?
			wop, _   = wopFrom(sictx, typename, -1, 0, nil)
		)
		return wstreamer{
			wstreamAnyPtr(wop),
			cfg,
		}

	case rmeta.ObjectP, rmeta.Objectp:
		var (
			se       = descr.elem
			typename = strings.TrimRight(se.TypeName(), "*") // FIXME(sbinet): handle T** ?
			wop, _   = wopFrom(sictx, typename, -1, 0, nil)
		)
		return wstreamer{
			wstreamObjPtr(wop),
			cfg,
		}

	case rmeta.StreamLoop:
		var (
			se       = descr.elem.(*StreamerLoop)
			typename = strings.TrimRight(se.TypeName(), "*") // FIXME(sbinet): handle T** ?
			wop, _   = wopFrom(sictx, typename, -1, 0, nil)
		)
		return wstreamer{
			wstreamBasicSlice(wop),
			cfg,
		}

	default:
		panic(fmt.Errorf("not implemented k=%d (%v)", descr.otype, descr.otype))
		// return wstreamer{wstreamGeneric, &streamerConfig{si, i, &descr, descr.offset, 0}}
	}
}

func wstreamSI(si *StreamerInfo) wopFunc {
	typename := si.Name()
	switch {
	case typename == "TObject":
		return wstreamTObject
	case typename == "TNamed":
		return wstreamTNamed
	case typename == "TString":
		return wstreamTString
	case rtypes.Factory.HasKey(typename):
		obj := rtypes.Factory.Get(typename)().Interface()
		_, ok := obj.(rbytes.Marshaler)
		if ok {
			return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
				obj := cfg.adjust(recv).(rbytes.Marshaler)
				return obj.MarshalROOT(w)
			}
		}
	}
	return wstreamCat(typename, int16(si.ClassVersion()), si.woops)
}

func wstreamObjPtr(wop wopFunc) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		var (
			pos = w.Pos()
			rv  = reflect.ValueOf(cfg.adjust(recv)).Elem()
			ptr root.Object
		)
		if !((rv == reflect.Value{}) || rv.IsNil()) {
			ptr = rv.Interface().(root.Object)
		}

		w.WriteObjectAny(ptr)
		return int(w.Pos() - pos), w.Err()
	}
}

func wstreamAnyPtr(wop wopFunc) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		var (
			pos = w.Pos()
			rv  = reflect.ValueOf(cfg.adjust(recv)).Elem()
			ptr root.Object
		)
		if !((rv == reflect.Value{}) || rv.IsNil()) {
			ptr = rv.Interface().(root.Object)
		}

		w.WriteObjectAny(ptr)
		return int(w.Pos() - pos), w.Err()
	}
}

func wstreamBool(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteBool(*cfg.adjust(recv).(*bool))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 1, nil
}

func wstreamU8(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteU8(*cfg.adjust(recv).(*uint8))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 1, nil
}

func wstreamU16(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteU16(*cfg.adjust(recv).(*uint16))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 2, nil
}

func wstreamU32(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteU32(*cfg.adjust(recv).(*uint32))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 4, nil
}

func wstreamU64(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteU64(*cfg.adjust(recv).(*uint64))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 8, nil
}

func wstreamI8(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteI8(*cfg.adjust(recv).(*int8))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 1, nil
}

func wstreamI16(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteI16(*cfg.adjust(recv).(*int16))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 2, nil
}

func wstreamI32(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteI32(*cfg.adjust(recv).(*int32))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 4, nil
}

func wstreamI64(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteI64(*cfg.adjust(recv).(*int64))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 8, nil
}

func wstreamF32(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteF32(*cfg.adjust(recv).(*float32))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 4, nil
}

func wstreamF64(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	w.WriteF64(*cfg.adjust(recv).(*float64))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 8, nil
}

func wstreamBits(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	// FIXME(sbinet) handle TObject reference
	// if (bits&kIsReferenced) != 0 { ... }
	w.WriteU32(*cfg.adjust(recv).(*uint32))
	if err := w.Err(); err != nil {
		return 0, err
	}
	return 4, nil
}

func wstreamF16(se rbytes.StreamerElement) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		beg := w.Pos()
		w.WriteF16(*cfg.adjust(recv).(*root.Float16), se)
		if err := w.Err(); err != nil {
			return 0, err
		}
		return int(w.Pos() - beg), w.Err()
	}
}

func wstreamD32(se rbytes.StreamerElement) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		beg := w.Pos()
		w.WriteD32(*cfg.adjust(recv).(*root.Double32), se)
		if err := w.Err(); err != nil {
			return 0, err
		}
		return int(w.Pos() - beg), w.Err()
	}
}

func wstreamTString(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	beg := w.Pos()
	w.WriteString(*cfg.adjust(recv).(*string))
	return int(w.Pos() - beg), w.Err()
}

func wstreamTObject(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	obj := cfg.adjust(recv).(*rbase.Object)
	return obj.MarshalROOT(w)
}

func wstreamTNamed(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	named := cfg.adjust(recv).(*rbase.Named)
	return named.MarshalROOT(w)
}

func wstreamBasicArray(n int, arr wopFunc) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		var (
			nn = 0
			rv = reflect.ValueOf(cfg.adjust(recv)).Elem()
		)
		for i := 0; i < n; i++ {
			nb, err := arr(w, rv.Index(i).Addr().Interface(), nil)
			if err != nil {
				return 0, fmt.Errorf(
					"rdict: could not wstream array element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, cfg.si.Name(), err,
				)
			}
			nn += nb
		}
		return nn, nil
	}
}

func wstreamBasicSlice(sli wopFunc) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		w.WriteI8(1) // is-array
		var (
			nn = 1
			n  = int(reflect.ValueOf(recv).Elem().FieldByIndex(cfg.descr.method).Int())
			rv = reflect.ValueOf(cfg.adjust(recv)).Elem()
		)
		for i := 0; i < n; i++ {
			nb, err := sli(w, rv.Index(i).Addr().Interface(), nil)
			if err != nil {
				return nn, fmt.Errorf(
					"rdict: could not wstream slice element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, cfg.si.Name(), err,
				)
			}
			nn += nb
		}
		return nn, nil
	}
}

func wstreamHeader(w *rbytes.WBuffer, typename string, typevers int16) rbytes.Header {
	if _, ok := rmeta.CxxBuiltins[typename]; ok && typename != "string" {
		return rbytes.Header{Pos: -1}
	}
	if typename == "TString" {
		return rbytes.Header{Pos: -1}
	}
	return w.WriteHeader(typename, typevers)
}

func wsetHeader(w *rbytes.WBuffer, hdr rbytes.Header) (int, error) {
	if hdr.Pos < 0 {
		return 0, nil
	}
	return w.SetHeader(hdr)
}

func wstreamType(typename string, wop wopFunc) wopFunc {
	const typevers = rvers.StreamerInfo
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		hdr := w.WriteHeader(typename, int16(typevers))
		n, err := wop(w, recv, cfg)
		if err != nil {
			return n, err
		}
		return w.SetHeader(hdr)
	}
}

func wstreamStdSlice(typename string, wop wopFunc) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		var (
			rv = reflect.ValueOf(cfg.adjust(recv)).Elem()
			n  = rv.Len()
			nn = 0
		)
		w.WriteI32(int32(n))
		for i := 0; i < n; i++ {
			nb, err := wop(w, rv.Index(i).Addr().Interface(), nil)
			if err != nil {
				return nn, fmt.Errorf(
					"rdict: could not wstream element %s[%d] of %s: %w",
					cfg.descr.elem.Name(), i, typename, err,
				)
			}
			nn += nb
		}
		return nn, w.Err()
	}
}

func wstreamStdSet(typename string, wop wopFunc) wopFunc {
	// FIXME(sbinet): add special handling for std::set-like types
	// the correct equivalent Go-type of std::set<T> is map[T]struct{}
	// (or, when availaible, std.Set[T])
	return wstreamStdSlice(typename, wop)
}

func wstreamStdMap(kname, vname string, kwop, vwop wopFunc, kvers, vvers int16) wopFunc {
	typename := fmt.Sprintf("map<%s,%s>", kname, vname)
	if strings.HasSuffix(vname, ">") {
		typename = fmt.Sprintf("map<%s,%s >", kname, vname)
	}
	const typevers = rvers.StreamerInfo
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		var (
			rv = reflect.ValueOf(cfg.adjust(recv)).Elem()
			n  = rv.Len()
			nn = 0
		)
		hdr := w.WriteHeader(typename, int16(typevers))
		w.WriteI32(int32(n))
		keyT := reflect.SliceOf(rv.Type().Key())
		valT := reflect.SliceOf(rv.Type().Elem())
		keys := reflect.New(keyT).Elem()
		vals := reflect.New(valT).Elem()
		keys.Set(reflect.AppendSlice(keys, reflect.MakeSlice(keyT, n, n)))
		vals.Set(reflect.AppendSlice(vals, reflect.MakeSlice(valT, n, n)))

		iter := rv.MapRange()
		for i := 0; iter.Next(); i++ {
			key := iter.Key()
			val := iter.Value()
			keys.Index(i).Set(key)
			vals.Index(i).Set(val)
		}
		if n > 0 {
			hdr := wstreamHeader(w, kname, kvers)
			for i := 0; i < n; i++ {
				nb, err := kwop(w, keys.Index(i).Addr().Interface(), nil)
				if err != nil {
					return nn, fmt.Errorf(
						"rdict: could not wstream key-element %s[%d] of %s: %w",
						kname, i, cfg.si.Name(), err,
					)
				}
				nn += nb
			}
			nb, err := wsetHeader(w, hdr)
			if err != nil {
				return nn, err
			}
			nn += nb
		}

		if n > 0 {
			hdr := wstreamHeader(w, vname, vvers)
			for i := 0; i < n; i++ {
				nb, err := vwop(w, vals.Index(i).Addr().Interface(), nil)
				if err != nil {
					return nn, fmt.Errorf(
						"rdict: could not rstream val-element %s[%d] of %s: %w",
						vname, i, cfg.si.Name(), err,
					)
				}
				nn += nb
			}
			_, err := wsetHeader(w, hdr)
			if err != nil {
				return nn, err
			}
		}

		return w.SetHeader(hdr)
	}
}

func wstreamStdBitset(typename string, n int) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		sli := *cfg.adjust(recv).(*[]uint8)
		sli = sli[:n]

		w.WriteI32(int32(n))
		w.WriteStdBitset(sli)
		return n + 4, w.Err()
	}
}

func wstreamBools(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]bool)
	)
	sli = sli[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayBool(sli)
	return 1 + n, w.Err()
}

func wstreamI8s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]int8)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayI8(sli)
	return 1 + n, w.Err()
}

func wstreamI16s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]int16)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayI16(sli)
	return 1 + n*2, w.Err()
}

func wstreamI32s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]int32)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayI32(sli)
	return 1 + n*4, w.Err()
}

func wstreamI64s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]int64)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayI64(sli)
	return 1 + n*8, w.Err()
}

func wstreamU8s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]uint8)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayU8(sli)
	return 1 + n, w.Err()
}

func wstreamU16s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]uint16)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayU16(sli)
	return 1 + n*2, w.Err()
}

func wstreamU32s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]uint32)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayU32(sli)
	return 1 + n*4, w.Err()
}

func wstreamU64s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]uint64)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayU64(sli)
	return 1 + n*8, w.Err()
}

func wstreamF32s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]float32)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayF32(sli)
	return 1 + n*4, w.Err()
}

func wstreamF64s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]float64)
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayF64(sli)
	return 1 + n*8, w.Err()
}

func wstreamF16s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]root.Float16)
		beg = w.Pos()
	)
	sli = sli[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayF16(sli, cfg.descr.elem)
	return int(w.Pos() - beg), w.Err()
}

func wstreamD32s(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]root.Double32)
		beg = w.Pos()
	)
	sli = sli[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayD32(sli, cfg.descr.elem)
	return int(w.Pos() - beg), w.Err()
}

func wstreamStrs(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	var (
		n   = cfg.counter(recv)
		sli = *cfg.adjust(recv).(*[]string)
		beg = w.Pos()
	)
	sli = (sli)[:n]
	w.WriteI8(1) // is-array
	w.WriteArrayString(sli)
	return int(w.Pos() - beg), w.Err()
}

func wstreamCat(typename string, typevers int16, wops []wstreamer) wopFunc {
	return func(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
		hdr := w.WriteHeader(typename, typevers)
		recv = cfg.adjust(recv)
		for i, wop := range wops {
			_, err := wop.wstream(w, recv)
			if err != nil {
				return 0, fmt.Errorf(
					"rdict: could not wstream element %d (%s) of %s: %w",
					i, wop.cfg.descr.elem.Name(), cfg.si.Name(), err,
				)
			}
		}
		return w.SetHeader(hdr)
	}
}

func wstreamStdString(w *rbytes.WBuffer, recv interface{}, cfg *streamerConfig) (int, error) {
	beg := w.Pos()
	w.WriteString(*cfg.adjust(recv).(*string))
	return int(w.Pos() - beg), w.Err()

}

func wopFuncFor(e rmeta.Enum, descr *elemDescr) wopFunc {
	switch e {
	case rmeta.Bool:
		return wstreamBool
	case rmeta.Bits:
		return wstreamBits
	case rmeta.Int8:
		return wstreamI8
	case rmeta.Int16:
		return wstreamI16
	case rmeta.Int32:
		return wstreamI32
	case rmeta.Int64, rmeta.Long64:
		return wstreamI64
	case rmeta.Uint8:
		return wstreamU8
	case rmeta.Uint16:
		return wstreamU16
	case rmeta.Uint32:
		return wstreamU32
	case rmeta.Uint64, rmeta.ULong64:
		return wstreamU64
	case rmeta.Float32:
		return wstreamF32
	case rmeta.Float64:
		return wstreamF64
	case rmeta.Float16:
		return wstreamF16(descr.elem)
	case rmeta.Double32:
		return wstreamD32(descr.elem)
	case rmeta.TString, rmeta.CharStar:
		return wstreamTString
	case rmeta.STLstring:
		return wstreamStdString
	case rmeta.TObject:
		return wstreamTObject
	case rmeta.TNamed:
		return wstreamTNamed
	default:
		return nil
	}
}

func wopFrom(sictx rbytes.StreamerInfoContext, typename string, typevers int16, enum rmeta.Enum, descr *elemDescr) (wopFunc, int16) {
	e, ok := rmeta.TypeName2Enum(typename)
	if ok {
		wop := wopFuncFor(e, descr)
		if wop != nil {
			return wop, -1
		}
	}

	wop := wopFuncFor(enum, descr)
	if wop != nil {
		return wop, -1
	}

	switch {
	case hasStdPrefix(typename, "vector", "list", "deque"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		wop, _ := wopFrom(sictx, enames[0], -1, 0, nil)
		return wstreamStdSlice(typename, wop), rvers.StreamerInfo

	case hasStdPrefix(typename, "set", "multiset", "unordered_set", "unordered_multiset"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		wop, _ := wopFrom(sictx, enames[0], -1, 0, nil)
		return wstreamStdSet(typename, wop), rvers.StreamerInfo

	case hasStdPrefix(typename, "map", "multimap", "unordered_map", "unordered_multimap"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		kname := enames[0]
		vname := enames[1]

		kwop, kvers := wopFrom(sictx, kname, -1, 0, nil)
		vwop, vvers := wopFrom(sictx, vname, -1, 0, nil)
		return wstreamStdMap(kname, vname, kwop, vwop, kvers, vvers), rvers.StreamerInfo

	case hasStdPrefix(typename, "bitset"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		n, err := strconv.Atoi(enames[0])
		if err != nil {
			panic(fmt.Errorf("rdict: invalid STL bitset argument (type=%q): %+v", typename, err))
		}
		return wstreamStdBitset(typename, n), rvers.StreamerInfo
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

	wop = wstreamSI(esi)
	return wop, int16(esi.ClassVersion())
}
