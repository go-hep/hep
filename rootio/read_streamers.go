// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
	"strings"
)

type RStreamer interface {
	RStream(r *RBuffer) error
}

type rstreamerFunc func(r *RBuffer) error

type rstreamerImpl struct {
	funcs []rstreamerFunc
}

func (rs *rstreamerImpl) RStream(r *RBuffer) error {
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
		name := f.Tag.Get("rootio")
		if name == "" {
			name = f.Name
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

func rstreamerFrom(se StreamerElement, ptr interface{}, lcnt leafCount, sictx StreamerInfoContext) rstreamerFunc {
	rt := reflect.TypeOf(ptr).Elem()
	rv := reflect.ValueOf(ptr).Elem()
	rf := rv
	if rt.Kind() == reflect.Struct {
		field := fieldOf(rt, se.Name())
		if field < 0 {
			panic(fmt.Errorf("rootio: no such field %q in type %T", se.Name(), ptr))
		}

		rf = rv.Field(field)
	}

	switch se := se.(type) {
	default:
		panic(fmt.Errorf("rootio: unknown streamer element: %#v", se))

	case *tstreamerBasicType:
		switch se.etype {
		case kCounter:
			switch se.esize {
			case 4:
				fptr := rf.Addr().Interface().(*int32)
				return func(r *RBuffer) error {
					if r.err != nil {
						return r.err
					}
					*fptr = r.ReadI32()
					return r.err
				}
			case 8:
				fptr := rf.Addr().Interface().(*int64)
				return func(r *RBuffer) error {
					if r.err != nil {
						return r.err
					}
					*fptr = r.ReadI64()
					return r.err
				}
			default:
				panic(fmt.Errorf("rootio: invalid kCounter size %d", se.esize))
			}

		case kChar:
			fptr := rf.Addr().Interface().(*int8)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadI8()
				return r.err
			}

		case kShort:
			fptr := rf.Addr().Interface().(*int16)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadI16()
				return r.err
			}

		case kInt:
			fptr := rf.Addr().Interface().(*int32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadI32()
				return r.err
			}

		case kLong, kLong64:
			fptr := rf.Addr().Interface().(*int64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadI64()
				return r.err
			}

		case kFloat:
			fptr := rf.Addr().Interface().(*float32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadF32()
				return r.err
			}

		case kDouble:
			fptr := rf.Addr().Interface().(*float64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadF64()
				return r.err
			}

		case kUChar, kCharStar:
			fptr := rf.Addr().Interface().(*uint8)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadU8()
				return r.err
			}

		case kUShort:
			fptr := rf.Addr().Interface().(*uint16)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadU16()
				return r.err
			}

		case kUInt, kBits:
			fptr := rf.Addr().Interface().(*uint32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadU32()
				return r.err
			}

		case kULong, kULong64:
			fptr := rf.Addr().Interface().(*uint64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadU64()
				return r.err
			}

		case kBool:
			fptr := rf.Addr().Interface().(*bool)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadI8() != 0
				return r.err
			}

		case kOffsetL + kChar:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int8)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayI8(n))
				return r.err
			}

		case kOffsetL + kShort:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int16)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayI16(n))
				return r.err
			}

		case kOffsetL + kInt:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayI32(n))
				return r.err
			}

		case kOffsetL + kLong, kOffsetL + kLong64:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]int64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayI64(n))
				return r.err
			}

		case kOffsetL + kFloat:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]float32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayF32(n))
				return r.err
			}

		case kOffsetL + kDouble:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]float64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayF64(n))
				return r.err
			}

		case kOffsetL + kUChar, kOffsetL + kCharStar:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint8)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayU8(n))
				return r.err
			}

		case kOffsetL + kUShort:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint16)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayU16(n))
				return r.err
			}

		case kOffsetL + kUInt, kOffsetL + kBits:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayU32(n))
				return r.err
			}

		case kOffsetL + kULong, kOffsetL + kULong64:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayU64(n))
				return r.err
			}

		case kOffsetL + kBool:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]bool)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayBool(n))
				return r.err
			}

		default:
			panic(fmt.Errorf("rootio: invalid element type value %d for %#v", se.etype, se))
		}

	case *tstreamerString:
		fptr := rf.Addr().Interface().(*string)
		return func(r *RBuffer) error {
			*fptr = r.ReadString()
			return r.err
		}

	case *tstreamerBasicPointer:
		flen := func() int { return 1 }
		if se.cname != "" {
			switch rv.Kind() {
			case reflect.Struct:
				fln := se.cname
				fptr := rv.FieldByNameFunc(func(n string) bool {
					if n == fln {
						return true
					}
					rf, ok := rt.FieldByName(n)
					if !ok {
						return false
					}
					if rf.Tag.Get("rootio") == fln {
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
		switch se.etype {
		case kOffsetP + kChar:
			fptr := rf.Addr().Interface().(*[]int8)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayI8(n)
				} else {
					*fptr = []int8{}
				}
				return r.err
			}

		case kOffsetP + kShort:
			fptr := rf.Addr().Interface().(*[]int16)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayI16(n)
				} else {
					*fptr = []int16{}
				}
				return r.err
			}

		case kOffsetP + kInt:
			fptr := rf.Addr().Interface().(*[]int32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayI32(n)
				} else {
					*fptr = []int32{}
				}
				return r.err
			}

		case kOffsetP + kLong, kOffsetP + kLong64:
			fptr := rf.Addr().Interface().(*[]int64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayI64(n)
				} else {
					*fptr = []int64{}
				}
				return r.err
			}

		case kOffsetP + kFloat:
			fptr := rf.Addr().Interface().(*[]float32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayF32(n)
				} else {
					*fptr = []float32{}
				}
				return r.err
			}

		case kOffsetP + kDouble:
			fptr := rf.Addr().Interface().(*[]float64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayF64(n)
				} else {
					*fptr = []float64{}
				}
				return r.err
			}

		case kOffsetP + kUChar, kOffsetP + kCharStar:
			fptr := rf.Addr().Interface().(*[]uint8)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayU8(n)
				} else {
					*fptr = []uint8{}
				}
				return r.err
			}

		case kOffsetP + kUShort:
			fptr := rf.Addr().Interface().(*[]uint16)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayU16(n)
				} else {
					*fptr = []uint16{}
				}
				return r.err
			}

		case kOffsetP + kUInt, kOffsetP + kBits:
			fptr := rf.Addr().Interface().(*[]uint32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayU32(n)
				} else {
					*fptr = []uint32{}
				}
				return r.err
			}

		case kOffsetP + kULong, kOffsetP + kULong64:
			fptr := rf.Addr().Interface().(*[]uint64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayU64(n)
				} else {
					*fptr = []uint64{}
				}
				return r.err
			}

		case kOffsetP + kBool:
			fptr := rf.Addr().Interface().(*[]bool)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				n := flen()
				_ = r.ReadU8()
				if n > 0 {
					*fptr = r.ReadFastArrayBool(n)
				} else {
					*fptr = []bool{}
				}
				return r.err
			}

		default:
			panic(fmt.Errorf("rootio: invalid element type value %d for %#v", se.etype, se))
		}

	case *tstreamerSTLstring:
		switch se.ctype {
		case kSTLstring:
			fptr := rf.Addr().Interface().(*string)
			return func(r *RBuffer) error {
				start := r.Pos()
				_, pos, bcnt := r.ReadVersion()
				*fptr = r.ReadString()
				r.CheckByteCount(pos, bcnt, start, "std::string")
				return r.err
			}
		default:
			panic(fmt.Errorf("rootio: invalid element type value %d for %#v", se.ctype, se))
		}

	case *tstreamerSTL:
		switch se.vtype {
		case kSTLvector:
			switch se.ctype {
			case kShort:
				fptr := rf.Addr().Interface().(*[]int16)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayI16(n)
					} else {
						*fptr = []int16{}
					}
					return r.err
				}

			case kInt:
				fptr := rf.Addr().Interface().(*[]int32)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayI32(n)
					} else {
						*fptr = []int32{}
					}
					return r.err
				}

			case kLong, kLong64:
				fptr := rf.Addr().Interface().(*[]int64)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayI64(n)
					} else {
						*fptr = []int64{}
					}
					return r.err
				}

			case kFloat:
				fptr := rf.Addr().Interface().(*[]float32)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayF32(n)
					} else {
						*fptr = []float32{}
					}
					return r.err
				}

			case kDouble:
				fptr := rf.Addr().Interface().(*[]float64)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayF64(n)
					} else {
						*fptr = []float64{}
					}
					return r.err
				}

			case kUShort:
				fptr := rf.Addr().Interface().(*[]uint16)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayU16(n)
					} else {
						*fptr = []uint16{}
					}
					return r.err
				}

			case kUInt, kBits:
				fptr := rf.Addr().Interface().(*[]uint32)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayU32(n)
					} else {
						*fptr = []uint32{}
					}
					return r.err
				}

			case kULong, kULong64:
				fptr := rf.Addr().Interface().(*[]uint64)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayU64(n)
					} else {
						*fptr = []uint64{}
					}
					return r.err
				}

			case kBool:
				fptr := rf.Addr().Interface().(*[]bool)
				return func(r *RBuffer) error {
					var hdr [6]byte
					r.read(hdr[:])
					n := int(r.ReadI32())
					if n > 0 {
						*fptr = r.ReadFastArrayBool(n)
					} else {
						*fptr = []bool{}
					}
					return r.err
				}

			case kObject:
				switch se.ename {
				case "vector<string>", "std::vector<std::string>":
					fptr := rf.Addr().Interface().(*[]string)
					return func(r *RBuffer) error {
						start := r.Pos()
						_, pos, bcnt := r.ReadVersion()
						n := int(r.ReadI32())
						*fptr = make([]string, n)
						for i := 0; i < n; i++ {
							(*fptr)[i] = r.ReadString()
						}
						r.CheckByteCount(pos, bcnt, start, "std::vector<std::string>")
						return r.err
					}
				default:
					subsi, err := sictx.StreamerInfo(se.elemTypeName())
					if err != nil {
						panic(fmt.Errorf("rootio: could not retrieve streamer for %q: %v", se.elemTypeName(), err))
					}
					eptr := reflect.New(rf.Type().Elem())
					felt := rstreamerFrom(subsi.Elements()[0], eptr.Interface(), lcnt, sictx)
					fptr := rf.Addr()
					return func(r *RBuffer) error {
						start := r.Pos()
						_, pos, bcnt := r.ReadVersion()
						n := int(r.ReadI32())
						if fptr.Elem().Len() < n {
							fptr.Elem().Set(reflect.MakeSlice(rf.Type(), n, n))
						}
						sli := fptr.Elem()
						for i := 0; i < n; i++ {
							felt(r)
							sli.Index(i).Set(eptr.Elem())
						}

						r.CheckByteCount(pos, bcnt, start, se.TypeName())
						return r.err
					}
				}
			}
		default:
			panic(fmt.Errorf("rootio: invalid STL type %d for %#v", se.vtype, se))
		}

	case *tstreamerObjectAny:
		sinfo, ok := streamers.getAny(se.ename)
		if !ok {
			panic(fmt.Errorf("no streamer-info for %q", se.ename))
		}
		var funcs []func(r *RBuffer) error
		for i, elt := range sinfo.Elements() {
			fptr := rf.Field(i).Addr().Interface()
			funcs = append(funcs, rstreamerFrom(elt, fptr, lcnt, sictx))
		}
		return func(r *RBuffer) error {
			start := r.Pos()
			_, pos, bcnt := r.ReadVersion()
			chksum := int(r.ReadI32())
			if sinfo.CheckSum() != chksum {
				return fmt.Errorf("rootio: on-disk checksum=%d, streamer=%d (type=%q)", chksum, sinfo.CheckSum(), se.ename)
			}
			for _, fct := range funcs {
				err := fct(r)
				if err != nil {
					return err
				}
			}
			r.CheckByteCount(pos, bcnt, start, se.ename)
			return nil
		}

	}
	panic(fmt.Errorf("rootio: unknown streamer element: %#v", se))
}

func stdvecSIFrom(name, ename string, ctx StreamerInfoContext) StreamerInfo {
	ename = strings.TrimSpace(ename)
	if etyp, ok := cxxbuiltins[ename]; ok {
		si := &tstreamerInfo{
			named: tnamed{
				name:  name,
				title: name,
			},
			elems: []StreamerElement{
				&tstreamerSTL{
					tstreamerElement: tstreamerElement{
						named: tnamed{
							name:  name,
							title: name,
						},
						ename: name,
					},
					rvers: 0,
					vtype: kSTLvector,
					ctype: gotype2ROOTEnum[etyp],
				},
			},
		}
		return si
	}
	esi, err := ctx.StreamerInfo(ename)
	if esi == nil || err != nil {
		return nil
	}

	si := &tstreamerInfo{
		named: tnamed{
			name:  name,
			title: name,
		},
		elems: []StreamerElement{
			&tstreamerSTL{
				tstreamerElement: tstreamerElement{
					named: tnamed{
						name:  name,
						title: name,
					},
					ename: name,
				},
				rvers: 0,
				vtype: kSTLvector,
				ctype: kObject,
			},
		},
	}
	return si
}

func gotypeFromSI(sinfo StreamerInfo, ctx StreamerInfoContext) reflect.Type {
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

		var lcount Leaf
		if elt.Title() != "" {
			lcount = &tleaf{}
		}
		ft.Type = gotypeFromSE(elt, lcount, ctx)
		ft.Tag = reflect.StructTag(`rootio:"` + elt.Name() + `"`)
	}

	return reflect.StructOf(fields)
}

func gotypeFromSE(se StreamerElement, lcount Leaf, ctx StreamerInfoContext) reflect.Type {
	if typ, ok := builtins[se.TypeName()]; ok {
		return typ
	}
	switch se := se.(type) {
	default:
		panic(fmt.Errorf("rootio: unknown streamer element: %#v", se))

	case *tstreamerBasicType:
		switch se.etype {
		case kCounter:
			switch se.esize {
			case 4:
				return reflect.TypeOf(int32(0))
			case 8:
				return reflect.TypeOf(int64(0))
			default:
				panic(fmt.Errorf("rootio: invalid kCounter size %d", se.esize))
			}

		case kChar:
			return reflect.TypeOf(int8(0))
		case kShort:
			return reflect.TypeOf(int16(0))
		case kInt:
			return reflect.TypeOf(int32(0))
		case kLong, kLong64:
			return reflect.TypeOf(int64(0))
		case kFloat:
			return reflect.TypeOf(float32(0))
		case kFloat16:
			return reflect.TypeOf(Float16(0))
		case kDouble32:
			return reflect.TypeOf(Double32(0))
		case kDouble:
			return reflect.TypeOf(float64(0))
		case kUChar, kCharStar:
			return reflect.TypeOf(uint8(0))
		case kUShort:
			return reflect.TypeOf(uint16(0))
		case kUInt, kBits:
			return reflect.TypeOf(uint32(0))
		case kULong, kULong64:
			return reflect.TypeOf(uint64(0))
		case kBool:
			return reflect.TypeOf(false)
		case kOffsetL + kChar:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(int8(0)))
		case kOffsetL + kShort:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(int16(0)))
		case kOffsetL + kInt:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(int32(0)))
		case kOffsetL + kLong, kOffsetL + kLong64:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(int64(0)))
		case kOffsetL + kFloat:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(float32(0)))
		case kOffsetL + kFloat16:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(Float16(0)))
		case kOffsetL + kDouble32:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(Double32(0)))
		case kOffsetL + kDouble:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(float64(0)))
		case kOffsetL + kUChar, kOffsetL + kCharStar:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(uint8(0)))
		case kOffsetL + kUShort:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(uint16(0)))
		case kOffsetL + kUInt, kOffsetL + kBits:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(uint32(0)))
		case kOffsetL + kULong, kOffsetL + kULong64:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(uint64(0)))
		case kOffsetL + kBool:
			return reflect.ArrayOf(int(se.arrlen), reflect.TypeOf(false))
		default:
			panic(fmt.Errorf("rootio: invalid element type value %d for %#v", se.etype, se))
		}

	case *tstreamerString:
		return reflect.TypeOf("")

	case *tstreamerBasicPointer:
		switch se.etype {
		case kOffsetP + kChar:
			tp := reflect.TypeOf(int8(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kShort:
			tp := reflect.TypeOf(int16(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kInt:
			tp := reflect.TypeOf(int32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kLong, kOffsetP + kLong64:
			tp := reflect.TypeOf(int64(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kFloat:
			tp := reflect.TypeOf(float32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kFloat16:
			tp := reflect.TypeOf(Float16(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kDouble32:
			tp := reflect.TypeOf(Double32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kDouble:
			tp := reflect.TypeOf(float64(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kUChar, kOffsetP + kCharStar:
			tp := reflect.TypeOf(uint8(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kUShort:
			tp := reflect.TypeOf(uint16(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kUInt, kOffsetP + kBits:
			tp := reflect.TypeOf(uint32(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kULong, kOffsetP + kULong64:
			tp := reflect.TypeOf(uint64(0))
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		case kOffsetP + kBool:
			tp := reflect.TypeOf(false)
			if lcount != nil {
				return reflect.SliceOf(tp)
			}
			return reflect.PtrTo(tp)
		default:
			panic(fmt.Errorf("rootio: invalid element type value %d for %#v", se.etype, se))
		}

	case *tstreamerSTLstring:
		switch se.ctype {
		case kSTLstring:
			return reflect.TypeOf("")
		default:
			panic(fmt.Errorf("rootio: invalid element type value %d for %#v", se.ctype, se))
		}

	case *tstreamerSTL:
		switch se.vtype {
		case kSTLvector:
			switch se.ctype {
			case kChar:
				return reflect.SliceOf(reflect.TypeOf(int8(0)))
			case kShort:
				return reflect.SliceOf(reflect.TypeOf(int16(0)))
			case kInt:
				return reflect.SliceOf(reflect.TypeOf(int32(0)))
			case kLong:
				return reflect.SliceOf(reflect.TypeOf(int64(0)))
			case kFloat:
				return reflect.SliceOf(reflect.TypeOf(float32(0)))
			case kDouble:
				return reflect.SliceOf(reflect.TypeOf(float64(0)))
			case kUChar:
				return reflect.SliceOf(reflect.TypeOf(uint8(0)))
			case kUShort:
				return reflect.SliceOf(reflect.TypeOf(uint16(0)))
			case kUInt:
				return reflect.SliceOf(reflect.TypeOf(uint32(0)))
			case kULong:
				return reflect.SliceOf(reflect.TypeOf(uint64(0)))
			case kBool:
				return reflect.SliceOf(reflect.TypeOf(false))
			case kObject:
				switch se.ename {
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
					eltname := se.elemTypeName()
					if eltname == "" {
						panic(fmt.Errorf("rootio: could not find element name for %q", se.ename))
					}
					if et, ok := cxxbuiltins[eltname]; ok {
						return reflect.SliceOf(et)
					}
					sielt, err := ctx.StreamerInfo(eltname)
					if err != nil {
						panic(err)
					}
					o := gotypeFromSI(sielt, ctx)
					if o == nil {
						panic(fmt.Errorf("rootio: invalid std::vector<kObject>: ename=%q", se.ename))
					}
					return reflect.SliceOf(o)
				}
			}
		default:
			panic(fmt.Errorf("rootio: invalid STL type %d for %#v", se.vtype, se))
		}

	case *tstreamerObjectAny:
		si, err := ctx.StreamerInfo(se.ename)
		if err != nil {
			panic(err)
		}
		return gotypeFromSI(si, ctx)

	case *tstreamerBase:
		switch se.ename {
		case "BASE":
			si, err := ctx.StreamerInfo(se.Name())
			if err != nil {
				panic(err)
			}
			return gotypeFromSI(si, ctx)

		default:
			panic(fmt.Errorf("rootio: unknown base class %q in StreamerElement %q: %#v", se.ename, se.Name(), se))
		}

	case *tstreamerObject:
		si, err := ctx.StreamerInfo(se.ename)
		if err != nil {
			panic(err)
		}
		return gotypeFromSI(si, ctx)

	case *tstreamerObjectPointer:
		ename := se.ename[:len(se.ename)-1] // drop final '*'
		si, err := ctx.StreamerInfo(ename)
		if err != nil {
			panic(err)
		}
		typ := gotypeFromSI(si, ctx)
		return reflect.PtrTo(typ)

	case *tstreamerObjectAnyPointer:
		ename := se.ename[:len(se.ename)-1] // drop final '*'
		si, err := ctx.StreamerInfo(ename)
		if err != nil {
			panic(err)
		}
		typ := gotypeFromSI(si, ctx)
		return reflect.PtrTo(typ)
	}

	panic(fmt.Errorf("rootio: unknown streamer element: %#v", se))
}
