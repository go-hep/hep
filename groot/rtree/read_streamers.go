// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

var (
	cxxNameSanitizer = strings.NewReplacer(
		"<", "_",
		">", "_",
		":", "_",
		",", "_",
		" ", "_",
	)
)

type rstreamerFunc func(r *rbytes.RBuffer) error

type rstreamerImpl struct {
	funcs []rstreamerFunc
}

func (rs *rstreamerImpl) RStreamROOT(r *rbytes.RBuffer) error {
	for _, rfunc := range rs.funcs {
		err := rfunc(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func fieldOf(rt reflect.Type, field string) int {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		name := f.Tag.Get("groot")
		if name == "" {
			name = f.Name
		}
		if idx := strings.Index(name, "["); idx > 0 {
			name = name[:idx]
		}
		if name == field {
			return i
		}
		if f.Name == field {
			return i
		}
	}
	return -1
}

func rstreamerFrom(se rbytes.StreamerElement, ptr interface{}, lcnt leafCount, sictx rbytes.StreamerInfoContext) rstreamerFunc {
	rt := reflect.TypeOf(ptr).Elem()
	rv := reflect.ValueOf(ptr).Elem()
	rf := rv
	if rt.Kind() == reflect.Struct {
		field := fieldOf(rt, se.Name())
		if field < 0 {
			panic(fmt.Errorf("rtree: no such field %q in type %T", se.Name(), ptr))
		}

		rf = rv.Field(field)
	}

	switch se := se.(type) {
	default:
		panic(fmt.Errorf("rtree: unknown streamer element: %#v", se))

	case *rdict.StreamerBasicType:
		switch se.Type() {
		case rmeta.Counter:
			switch se.Size() {
			case 4:
				fptr := rf.Addr().Interface().(*int32)
				return func(r *rbytes.RBuffer) error {
					if r.Err() != nil {
						return r.Err()
					}
					*fptr = r.ReadI32()
					return r.Err()
				}
			case 8:
				fptr := rf.Addr().Interface().(*int64)
				return func(r *rbytes.RBuffer) error {
					if r.Err() != nil {
						return r.Err()
					}
					*fptr = r.ReadI64()
					return r.Err()
				}
			default:
				panic(fmt.Errorf("rtree: invalid kCounter size %d", se.Size()))
			}

		case rmeta.Char:
			fptr := rf.Addr().Interface().(*int8)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadI8()
				return r.Err()
			}

		case rmeta.Short:
			fptr := rf.Addr().Interface().(*int16)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadI16()
				return r.Err()
			}

		case rmeta.Int:
			fptr := rf.Addr().Interface().(*int32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadI32()
				return r.Err()
			}

		case rmeta.Long, rmeta.Long64:
			fptr := rf.Addr().Interface().(*int64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadI64()
				return r.Err()
			}

		case rmeta.Float:
			fptr := rf.Addr().Interface().(*float32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadF32()
				return r.Err()
			}

		case rmeta.Double:
			fptr := rf.Addr().Interface().(*float64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadF64()
				return r.Err()
			}

		case rmeta.UChar, rmeta.CharStar:
			fptr := rf.Addr().Interface().(*uint8)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadU8()
				return r.Err()
			}

		case rmeta.UShort:
			fptr := rf.Addr().Interface().(*uint16)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadU16()
				return r.Err()
			}

		case rmeta.UInt, rmeta.Bits:
			fptr := rf.Addr().Interface().(*uint32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadU32()
				return r.Err()
			}

		case rmeta.ULong, rmeta.ULong64:
			fptr := rf.Addr().Interface().(*uint64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadU64()
				return r.Err()
			}

		case rmeta.Bool:
			fptr := rf.Addr().Interface().(*bool)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				*fptr = r.ReadI8() != 0
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.Char:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int8)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayI8(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.Short:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int16)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayI16(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.Int:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayI32(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.Long, rmeta.OffsetL + rmeta.Long64:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayI64(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.Float:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]float32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayF32(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.Double:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]float64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayF64(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.UChar, rmeta.OffsetL + rmeta.CharStar:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint8)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayU8(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.UShort:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint16)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayU16(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.UInt, rmeta.OffsetL + rmeta.Bits:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayU32(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.ULong, rmeta.OffsetL + rmeta.ULong64:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayU64(fptr)
				return r.Err()
			}

		case rmeta.OffsetL + rmeta.Bool:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]bool)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				r.ReadArrayBool(fptr)
				return r.Err()
			}

		default:
			panic(fmt.Errorf("rtree: invalid element type value %d for %#v", se.Type(), se))
		}

	case *rdict.StreamerString:
		fptr := rf.Addr().Interface().(*string)
		return func(r *rbytes.RBuffer) error {
			*fptr = r.ReadString()
			return r.Err()
		}

	case *rdict.StreamerBasicPointer:
		flen := func() int { return 1 }
		if se.CountName() != "" {
			switch rv.Kind() {
			case reflect.Struct:
				fln := se.CountName()
				fptr := rv.FieldByNameFunc(func(n string) bool {
					if n == fln {
						return true
					}
					rf, ok := rt.FieldByName(n)
					if !ok {
						return false
					}
					if rf.Tag.Get("groot") == fln {
						return true
					}
					return false
				})
				flen = func() int {
					return int(fptr.Int())
				}
			default:
				if lcnt != nil {
					flen = lcnt.ivalue
				}
			}
		}
		switch se.Type() {
		case rmeta.OffsetP + rmeta.Char:
			fptr := rf.Addr().Interface().(*[]int8)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeI8(*fptr, n)
				if n > 0 {
					r.ReadArrayI8(*fptr)
				} else {
					*fptr = []int8{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.Short:
			fptr := rf.Addr().Interface().(*[]int16)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeI16(*fptr, n)
				if n > 0 {
					r.ReadArrayI16(*fptr)
				} else {
					*fptr = []int16{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.Int:
			fptr := rf.Addr().Interface().(*[]int32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeI32(*fptr, n)
				if n > 0 {
					r.ReadArrayI32(*fptr)
				} else {
					*fptr = []int32{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.Long, rmeta.OffsetP + rmeta.Long64:
			fptr := rf.Addr().Interface().(*[]int64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeI64(*fptr, n)
				if n > 0 {
					r.ReadArrayI64(*fptr)
				} else {
					*fptr = []int64{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.Float:
			fptr := rf.Addr().Interface().(*[]float32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeF32(*fptr, n)
				if n > 0 {
					r.ReadArrayF32(*fptr)
				} else {
					*fptr = []float32{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.Double:
			fptr := rf.Addr().Interface().(*[]float64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeF64(*fptr, n)
				if n > 0 {
					r.ReadArrayF64(*fptr)
				} else {
					*fptr = []float64{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.UChar, rmeta.OffsetP + rmeta.CharStar:
			fptr := rf.Addr().Interface().(*[]uint8)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeU8(*fptr, n)
				if n > 0 {
					r.ReadArrayU8(*fptr)
				} else {
					*fptr = []uint8{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.UShort:
			fptr := rf.Addr().Interface().(*[]uint16)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeU16(*fptr, n)
				if n > 0 {
					r.ReadArrayU16(*fptr)
				} else {
					*fptr = []uint16{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.UInt, rmeta.OffsetP + rmeta.Bits:
			fptr := rf.Addr().Interface().(*[]uint32)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeU32(*fptr, n)
				if n > 0 {
					r.ReadArrayU32(*fptr)
				} else {
					*fptr = []uint32{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.ULong, rmeta.OffsetP + rmeta.ULong64:
			fptr := rf.Addr().Interface().(*[]uint64)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeU64(*fptr, n)
				if n > 0 {
					r.ReadArrayU64(*fptr)
				} else {
					*fptr = []uint64{}
				}
				return r.Err()
			}

		case rmeta.OffsetP + rmeta.Bool:
			fptr := rf.Addr().Interface().(*[]bool)
			return func(r *rbytes.RBuffer) error {
				if r.Err() != nil {
					return r.Err()
				}
				n := flen()
				_ = r.ReadU8()
				*fptr = rbytes.ResizeBool(*fptr, n)
				if n > 0 {
					r.ReadArrayBool(*fptr)
				} else {
					*fptr = []bool{}
				}
				return r.Err()
			}

		default:
			panic(fmt.Errorf("rtree: invalid element type value %d for %#v", se.Type(), se))
		}

	case *rdict.StreamerSTLstring:
		switch se.ContainedType() {
		case rmeta.STLstring:
			fptr := rf.Addr().Interface().(*string)
			return func(r *rbytes.RBuffer) error {
				start := r.Pos()
				_, pos, bcnt := r.ReadVersion("string") // ROOT knows std::string as string.
				*fptr = r.ReadString()
				r.CheckByteCount(pos, bcnt, start, "std::string")
				return r.Err()
			}
		default:
			panic(fmt.Errorf("rtree: invalid element type value %d for %#v", se.ContainedType(), se))
		}

	case *rdict.StreamerSTL:
		switch se.STLType() {
		case rmeta.STLvector:
			switch se.ContainedType() {
			case rmeta.Char:
				fptr := rf.Addr().Interface().(*[]int8)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeI8(*fptr, n)
					if n > 0 {
						r.ReadArrayI8(*fptr)
					} else {
						*fptr = []int8{}
					}
					return r.Err()
				}

			case rmeta.Short:
				fptr := rf.Addr().Interface().(*[]int16)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeI16(*fptr, n)
					if n > 0 {
						r.ReadArrayI16(*fptr)
					} else {
						*fptr = []int16{}
					}
					return r.Err()
				}

			case rmeta.Int:
				fptr := rf.Addr().Interface().(*[]int32)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeI32(*fptr, n)
					if n > 0 {
						r.ReadArrayI32(*fptr)
					} else {
						*fptr = []int32{}
					}
					return r.Err()
				}

			case rmeta.Long, rmeta.Long64:
				fptr := rf.Addr().Interface().(*[]int64)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeI64(*fptr, n)
					if n > 0 {
						r.ReadArrayI64(*fptr)
					} else {
						*fptr = []int64{}
					}
					return r.Err()
				}

			case rmeta.Float:
				fptr := rf.Addr().Interface().(*[]float32)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeF32(*fptr, n)
					if n > 0 {
						r.ReadArrayF32(*fptr)
					} else {
						*fptr = []float32{}
					}
					return r.Err()
				}

			case rmeta.Double:
				fptr := rf.Addr().Interface().(*[]float64)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeF64(*fptr, n)
					if n > 0 {
						r.ReadArrayF64(*fptr)
					} else {
						*fptr = []float64{}
					}
					return r.Err()
				}

			case rmeta.UShort:
				fptr := rf.Addr().Interface().(*[]uint16)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeU16(*fptr, n)
					if n > 0 {
						r.ReadArrayU16(*fptr)
					} else {
						*fptr = []uint16{}
					}
					return r.Err()
				}

			case rmeta.UInt, rmeta.Bits:
				fptr := rf.Addr().Interface().(*[]uint32)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeU32(*fptr, n)
					if n > 0 {
						r.ReadArrayU32(*fptr)
					} else {
						*fptr = []uint32{}
					}
					return r.Err()
				}

			case rmeta.ULong, rmeta.ULong64:
				fptr := rf.Addr().Interface().(*[]uint64)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeU64(*fptr, n)
					if n > 0 {
						r.ReadArrayU64(*fptr)
					} else {
						*fptr = []uint64{}
					}
					return r.Err()
				}

			case rmeta.Bool:
				fptr := rf.Addr().Interface().(*[]bool)
				return func(r *rbytes.RBuffer) error {
					var hdr [6]byte
					_, _ = r.Read(hdr[:])
					n := int(r.ReadI32())
					*fptr = rbytes.ResizeBool(*fptr, n)
					if n > 0 {
						r.ReadArrayBool(*fptr)
					} else {
						*fptr = []bool{}
					}
					return r.Err()
				}

			case rmeta.Object:
				switch se.TypeName() {
				case "vector<string>", "std::vector<std::string>":
					fptr := rf.Addr().Interface().(*[]string)
					*fptr = make([]string, 0, 8)
					return func(r *rbytes.RBuffer) error {
						start := r.Pos()
						_, pos, bcnt := r.ReadVersion("vector<string>")
						n := int(r.ReadI32())
						*fptr = rbytes.ResizeStr(*fptr, n)
						for i := 0; i < n; i++ {
							(*fptr)[i] = r.ReadString()
						}
						r.CheckByteCount(pos, bcnt, start, "std::vector<std::string>")
						return r.Err()
					}
				default:
					// FIXME(sbinet): always load latest version?
					etn := se.ElemTypeName()
					subsi, err := sictx.StreamerInfo(etn[0], -1)
					if err != nil {
						panic(fmt.Errorf("rtree: could not retrieve streamer for %q: %w", etn[0], err))
					}
					eptr := reflect.New(rf.Type().Elem())
					if len(subsi.Elements()) <= 0 {
						panic(fmt.Errorf("rtree: invalid streamer info for %q", etn[0]))
					}

					felt := rstreamerFrom(subsi.Elements()[0], eptr.Interface(), lcnt, sictx)
					fptr := rf.Addr()
					typename := se.TypeName()
					return func(r *rbytes.RBuffer) error {
						start := r.Pos()
						_, pos, bcnt := r.ReadVersion(typename)
						n := int(r.ReadI32())
						if fptr.Elem().Len() < n {
							fptr.Elem().Set(reflect.MakeSlice(rf.Type(), n, n))
						}
						sli := fptr.Elem()
						for i := 0; i < n; i++ {
							_ = felt(r)
							sli.Index(i).Set(eptr.Elem())
						}

						r.CheckByteCount(pos, bcnt, start, typename)
						return r.Err()
					}
				}
			}
		default:
			panic(fmt.Errorf("rtree: invalid STL type %d for %#v", se.STLType(), se))
		}

	case *rdict.StreamerObjectAny:
		sinfo, err := sictx.StreamerInfo(se.TypeName(), -1)
		if err != nil {
			panic(fmt.Errorf("no streamer-info for %q", se.TypeName()))
		}
		var funcs []func(r *rbytes.RBuffer) error
		for i, elt := range sinfo.Elements() {
			fptr := rf.Field(i).Addr().Interface()
			funcs = append(funcs, rstreamerFrom(elt, fptr, lcnt, sictx))
		}
		typename := se.TypeName()
		return func(r *rbytes.RBuffer) error {
			start := r.Pos()
			_, pos, bcnt := r.ReadVersion(typename)
			for _, fct := range funcs {
				err := fct(r)
				if err != nil {
					return err
				}
			}
			r.CheckByteCount(pos, bcnt, start, typename)
			return nil
		}

	}
	panic(fmt.Errorf("rtree: unknown streamer element: %#v", se))
}

func gotypeFromSI(sinfo rbytes.StreamerInfo, ctx rbytes.StreamerInfoContext) reflect.Type {
	if typ, ok := builtins[sinfo.Name()]; ok {
		return typ
	}
	elts := sinfo.Elements()
	fields := make([]reflect.StructField, len(elts))
	for i := range fields {
		ft := &fields[i]
		elt := elts[i]
		ename := elt.Name()
		if ename == "" {
			panic(fmt.Errorf("elt[%d]: %q for si=%v", i, elt.Class(), sinfo))
		}
		ft.Name = "ROOT_" + elt.Name()
		ft.Name = cxxNameSanitizer.Replace(ft.Name)

		var (
			lcount Leaf
			ltitle = ""
		)
		if elt.Title() != "" {
			lcount = &tleaf{}
			ltitle = elt.Title()
		}
		ft.Type = gotypeFromSE(elt, lcount, ctx)
		if ft.Type.Kind() == reflect.Array {
			ltitle = fmt.Sprintf("[%d]", ft.Type.Len())
		}
		ft.Tag = reflect.StructTag(`groot:"` + elt.Name() + ltitle + `"`)
	}

	return reflect.StructOf(fields)
}

func gotypeFromSE(se rbytes.StreamerElement, lcount Leaf, ctx rbytes.StreamerInfoContext) reflect.Type {
	if typ, ok := builtins[se.TypeName()]; ok {
		return typ
	}
	switch se := se.(type) {
	default:
		panic(fmt.Errorf("rtree: unknown streamer element: %#v", se))

	case *rdict.StreamerBasicType:
		switch se.Type() {
		case rmeta.Counter:
			switch se.Size() {
			case 4:
				return reflect.TypeOf(int32(0))
			case 8:
				return reflect.TypeOf(int64(0))
			default:
				panic(fmt.Errorf("rtree: invalid rmeta.Counter size %d", se.Size()))
			}

		case rmeta.Char:
			return reflect.TypeOf(int8(0))
		case rmeta.Short:
			return reflect.TypeOf(int16(0))
		case rmeta.Int:
			return reflect.TypeOf(int32(0))
		case rmeta.Long, rmeta.Long64:
			return reflect.TypeOf(int64(0))
		case rmeta.Float:
			return reflect.TypeOf(float32(0))
		case rmeta.Float16:
			return reflect.TypeOf(root.Float16(0))
		case rmeta.Double32:
			return reflect.TypeOf(root.Double32(0))
		case rmeta.Double:
			return reflect.TypeOf(float64(0))
		case rmeta.UChar, rmeta.CharStar:
			return reflect.TypeOf(uint8(0))
		case rmeta.UShort:
			return reflect.TypeOf(uint16(0))
		case rmeta.UInt, rmeta.Bits:
			return reflect.TypeOf(uint32(0))
		case rmeta.ULong, rmeta.ULong64:
			return reflect.TypeOf(uint64(0))
		case rmeta.Bool:
			return reflect.TypeOf(false)
		case rmeta.OffsetL + rmeta.Char:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(int8(0)))
		case rmeta.OffsetL + rmeta.Short:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(int16(0)))
		case rmeta.OffsetL + rmeta.Int:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(int32(0)))
		case rmeta.OffsetL + rmeta.Long, rmeta.OffsetL + rmeta.Long64:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(int64(0)))
		case rmeta.OffsetL + rmeta.Float:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(float32(0)))
		case rmeta.OffsetL + rmeta.Float16:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(root.Float16(0)))
		case rmeta.OffsetL + rmeta.Double32:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(root.Double32(0)))
		case rmeta.OffsetL + rmeta.Double:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(float64(0)))
		case rmeta.OffsetL + rmeta.UChar, rmeta.OffsetL + rmeta.CharStar:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(uint8(0)))
		case rmeta.OffsetL + rmeta.UShort:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(uint16(0)))
		case rmeta.OffsetL + rmeta.UInt, rmeta.OffsetL + rmeta.Bits:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(uint32(0)))
		case rmeta.OffsetL + rmeta.ULong, rmeta.OffsetL + rmeta.ULong64:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(uint64(0)))
		case rmeta.OffsetL + rmeta.Bool:
			return reflect.ArrayOf(se.ArrayLen(), reflect.TypeOf(false))
		default:
			panic(fmt.Errorf("rtree: invalid element type value %d for %#v", se.Type(), se))
		}

	case *rdict.StreamerString:
		return reflect.TypeOf("")

	case *rdict.StreamerBasicPointer:
		switch se.Type() {
		case rmeta.OffsetP + rmeta.Char:
			tp := reflect.TypeOf(int8(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Short:
			tp := reflect.TypeOf(int16(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Int:
			tp := reflect.TypeOf(int32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Long, rmeta.OffsetP + rmeta.Long64:
			tp := reflect.TypeOf(int64(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Float:
			tp := reflect.TypeOf(float32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Float16:
			tp := reflect.TypeOf(root.Float16(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Double32:
			tp := reflect.TypeOf(root.Double32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Double:
			tp := reflect.TypeOf(float64(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.UChar, rmeta.OffsetP + rmeta.CharStar:
			tp := reflect.TypeOf(uint8(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.UShort:
			tp := reflect.TypeOf(uint16(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.UInt, rmeta.OffsetP + rmeta.Bits:
			tp := reflect.TypeOf(uint32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.ULong, rmeta.OffsetP + rmeta.ULong64:
			tp := reflect.TypeOf(uint64(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case rmeta.OffsetP + rmeta.Bool:
			tp := reflect.TypeOf(false)
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		default:
			panic(fmt.Errorf("rtree: invalid element type value %d for %#v", se.Type(), se))
		}

	case *rdict.StreamerSTLstring:
		switch se.ContainedType() {
		case rmeta.STLstring:
			return reflect.TypeOf("")
		default:
			panic(fmt.Errorf("rtree: invalid element type value %d for %#v", se.ContainedType(), se))
		}

	case *rdict.StreamerSTL:
		switch se.STLType() {
		case rmeta.STLvector:
			switch se.ContainedType() {
			case rmeta.Char:
				return reflect.SliceOf(reflect.TypeOf(int8(0)))
			case rmeta.Short:
				return reflect.SliceOf(reflect.TypeOf(int16(0)))
			case rmeta.Int:
				return reflect.SliceOf(reflect.TypeOf(int32(0)))
			case rmeta.Long:
				return reflect.SliceOf(reflect.TypeOf(int64(0)))
			case rmeta.Float:
				return reflect.SliceOf(reflect.TypeOf(float32(0)))
			case rmeta.Double:
				return reflect.SliceOf(reflect.TypeOf(float64(0)))
			case rmeta.UChar:
				return reflect.SliceOf(reflect.TypeOf(uint8(0)))
			case rmeta.UShort:
				return reflect.SliceOf(reflect.TypeOf(uint16(0)))
			case rmeta.UInt:
				return reflect.SliceOf(reflect.TypeOf(uint32(0)))
			case rmeta.ULong:
				return reflect.SliceOf(reflect.TypeOf(uint64(0)))
			case rmeta.Bool:
				return reflect.SliceOf(reflect.TypeOf(false))
			case rmeta.Object:
				switch se.TypeName() {
				case "vector<string>", "std::vector<std::string>":
					return reflect.TypeOf([]string(nil))
				case "vector<vector<char> >":
					return reflect.TypeOf([][]int8(nil))
				case "vector<vector<short> >":
					return reflect.TypeOf([][]uint16(nil))
				case "vector<vector<int> >":
					return reflect.TypeOf([][]int32(nil))
				case "vector<vector<long int> >", "vector<vector<long> >":
					return reflect.TypeOf([][]int64(nil))
				case "vector<vector<float> >":
					return reflect.TypeOf([][]float32(nil))
				case "vector<vector<double> >":
					return reflect.TypeOf([][]float64(nil))
				case "vector<vector<unsigned char> >":
					return reflect.TypeOf([][]uint8(nil))
				case "vector<vector<unsigned short> >":
					return reflect.TypeOf([][]uint16(nil))
				case "vector<vector<unsigned int> >", "vector<vector<unsigned> >":
					return reflect.TypeOf([][]uint32(nil))
				case "vector<vector<unsigned long int> >", "vector<vector<unsigned long> >":
					return reflect.TypeOf([][]uint64(nil))
				case "vector<vector<bool> >":
					return reflect.TypeOf([][]bool(nil))
				case "vector<vector<string> >":
					return reflect.TypeOf([][]string(nil))
				default:
					eltname := se.ElemTypeName()[0]
					if eltname == "" {
						panic(fmt.Errorf("rtree: could not find element name for %q", se.TypeName()))
					}
					if et, ok := rmeta.CxxBuiltins[eltname]; ok {
						return reflect.SliceOf(et)
					}
					// FIXME(sbinet): always load latest version?
					sielt, err := ctx.StreamerInfo(eltname, -1)
					if err != nil {
						panic(err)
					}
					o := gotypeFromSI(sielt, ctx)
					if o == nil {
						panic(fmt.Errorf("rtree: invalid std::vector<kObject>: ename=%q", se.TypeName()))
					}
					return reflect.SliceOf(o)
				}
			default:
				panic(fmt.Errorf("rtree: invalid STL contained-type %v for %#v", se.ContainedType(), se))
			}
		default:
			panic(fmt.Errorf("rtree: invalid STL type %d for %#v", se.STLType(), se))
		}

	case *rdict.StreamerObjectAny:
		// FIXME(sbinet): always load latest version?
		si, err := ctx.StreamerInfo(se.TypeName(), -1)
		if err != nil {
			panic(err)
		}
		return gotypeFromSI(si, ctx)

	case *rdict.StreamerBase:
		switch se.TypeName() {
		case "BASE":
			// FIXME(sbinet): always load latest version?
			si, err := ctx.StreamerInfo(se.Name(), -1)
			if err != nil {
				panic(err)
			}
			return gotypeFromSI(si, ctx)

		default:
			panic(fmt.Errorf("rtree: unknown base class %q in StreamerElement %q: %#v", se.TypeName(), se.Name(), se))
		}

	case *rdict.StreamerObject:
		// FIXME(sbinet): always load latest version?
		si, err := ctx.StreamerInfo(se.TypeName(), -1)
		if err != nil {
			panic(err)
		}
		return gotypeFromSI(si, ctx)

	case *rdict.StreamerObjectPointer:
		ename := se.TypeName()[:len(se.TypeName())-1] // drop final '*'
		// FIXME(sbinet): always load latest version?
		si, err := ctx.StreamerInfo(ename, -1)
		if err != nil {
			panic(err)
		}
		typ := gotypeFromSI(si, ctx)
		return reflect.PtrTo(typ)

	case *rdict.StreamerObjectAnyPointer:
		ename := se.TypeName()[:len(se.TypeName())-1] // drop final '*'
		// FIXME(sbinet): always load latest version?
		si, err := ctx.StreamerInfo(ename, -1)
		if err != nil {
			panic(err)
		}
		typ := gotypeFromSI(si, ctx)
		return reflect.PtrTo(typ)
	}

	panic(fmt.Errorf("rtree: unknown streamer element: %#v", se))
}
