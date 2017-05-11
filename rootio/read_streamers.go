// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
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

func rstreamerFrom(se StreamerElement, ptr interface{}, lcnt leafCount) rstreamerFunc {
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

		case kLong:
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

		case kUChar:
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

		case kUInt:
			fptr := rf.Addr().Interface().(*uint32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadU32()
				return r.err
			}

		case kULong:
			fptr := rf.Addr().Interface().(*uint64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				*fptr = r.ReadU64()
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

		case kOffsetL + kLong:
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

		case kOffsetL + kUChar:
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

		case kOffsetL + kUInt:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint32)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayU32(n))
				return r.err
			}

		case kOffsetL + kULong:
			n := rf.Len()
			fptr := rf.Slice(0, n).Interface().([]uint64)
			return func(r *RBuffer) error {
				if r.err != nil {
					return r.err
				}
				copy(fptr[:], r.ReadFastArrayU64(n))
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

		case kOffsetP + kLong:
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

		case kOffsetP + kUChar:
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

		case kOffsetP + kUInt:
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

		case kOffsetP + kULong:
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

			case kLong:
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

			case kUInt:
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

			case kULong:
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

			case kObject:
				switch se.ename {
				case "vector<string>":
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
			funcs = append(funcs, rstreamerFrom(elt, fptr, lcnt))
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
	return nil
}
