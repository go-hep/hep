// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
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

func ObjectFrom(si rbytes.StreamerInfo, sictx rbytes.StreamerInfoContext) *Object {
	return newObjectFrom(si, sictx)
}

type rfunc func(recv interface{}, r *rbytes.RBuffer) error
type wfunc func(recv interface{}, w *rbytes.WBuffer) (int, error)

// Object wraps a type created from a Streamer and implements the
// following interfaces:
//  - root.Object
//  - rbytes.RVersioner
//  - rbytes.Marshaler
//  - rbytes.Unmarshaler
type Object struct {
	v interface{}

	si    rbytes.StreamerInfo
	rvers int16
	class string

	rfuncs  []rfunc
	marshal wfunc
}

func (obj *Object) Class() string {
	return obj.class
}

func (obj *Object) SetClass(name string) {
	si, ok := Streamers.Get(name, -1)
	if !ok {
		panic(errors.Errorf("rdict: no streamer for %q", name))
	}
	*obj = *newObjectFrom(si, nil)
}

func (obj *Object) String() string {
	return fmt.Sprintf("%v", obj.v)
}

func (obj *Object) RVersion() int16 {
	return obj.rvers
}

func (obj *Object) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers != obj.rvers {
		r.SetErr(errors.Errorf("rdict: inconsistent ROOT version (got=%d, want=%d)", vers, obj.rvers))
		return r.Err()
	}

	rv := reflect.Indirect(reflect.ValueOf(obj.v))
	for i, rfunc := range obj.rfuncs {
		recv := rv.Field(i).Addr().Interface()
		err := rfunc(recv, r)
		if err != nil {
			return err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, obj.Class())
	return r.Err()
}

func (obj *Object) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	return obj.marshal(obj.v, w)
}

func newObjectFrom(si rbytes.StreamerInfo, sictx rbytes.StreamerInfoContext) *Object {
	rt := genTypeFromSI(sictx, si)
	recv := reflect.New(rt)
	obj := &Object{
		v:     recv.Interface(),
		si:    si,
		rvers: int16(si.ClassVersion()),
		class: si.Name(),
	}
	obj.rfuncs = genRStreamerFromSI(sictx, si, recv)
	return obj
}

type counter interface {
	ivalue() int
	imax() int
}

func genTypeFromSI(sictx rbytes.StreamerInfoContext, si rbytes.StreamerInfo) reflect.Type {
	var fields = make([]reflect.StructField, 0, len(si.Elements()))
	for _, se := range si.Elements() {
		rt := genTypeFromSE(sictx, se)
		ft := reflect.StructField{
			Name: "ROOT_" + cxxNameSanitizer.Replace(se.Name()),
			Type: rt,
			Tag:  reflect.StructTag(fmt.Sprintf("groot:%q", se.Name())),
		}
		fields = append(fields, ft)
	}
	return reflect.StructOf(fields)
}

func genTypeFromSE(sictx rbytes.StreamerInfoContext, se rbytes.StreamerElement) reflect.Type {
	switch se := se.(type) {
	default:
		panic(errors.Errorf("rdict: unknown streamer element: %#v (%T)", se, se))
	case *StreamerBasicType:
		return genType(sictx, se.Type(), se.ArrayLen())
	case *StreamerString:
		return gotypes[reflect.String]
	case *StreamerBasicPointer:
		return genType(sictx, se.Type(), -1)
	case *StreamerSTLstring:
		return gotypes[reflect.String]
	case *StreamerSTL:
		switch se.STLVectorType() {
		case rmeta.STLvector:
			return genType(sictx, se.ContainedType(), -1)
		}
		panic("STL container not implemented")
	}
	return nil
}

func genRStreamerFromSI(sictx rbytes.StreamerInfoContext, si rbytes.StreamerInfo, recv reflect.Value) []rfunc {
	var funcs = make([]rfunc, 0, len(si.Elements()))

	for i, se := range si.Elements() {
		sub := reflect.Indirect(recv).Field(i).Addr()
		rfunc := genRStreamerFromSE(sictx, se, sub)
		funcs = append(funcs, rfunc)
	}
	return funcs
}

func genRStreamerFromSE(sictx rbytes.StreamerInfoContext, se rbytes.StreamerElement, recv reflect.Value) rfunc {
	switch se := se.(type) {
	default:
		panic(errors.Errorf("rdict: unknown read-streamer element: %#v (%T)", se, se))
	case *StreamerBasicType:
		return genRStreamer(sictx, se.Type(), se.ArrayLen(), recv)
	case *StreamerString:
		return func(recv interface{}, r *rbytes.RBuffer) error {
			*(recv.(*string)) = r.ReadString()
			return r.Err()
		}
	}
	return nil
}

func genType(sictx rbytes.StreamerInfoContext, enum rmeta.Enum, n int) reflect.Type {
	switch enum {
	case rmeta.Bool:
		return gotypes[reflect.Bool]
	case rmeta.Uint8:
		return gotypes[reflect.Uint8]
	case rmeta.Uint16:
		return gotypes[reflect.Uint16]
	case rmeta.Uint32:
		return gotypes[reflect.Uint32]
	case rmeta.Uint64:
		return gotypes[reflect.Uint64]
	case rmeta.Int8:
		return gotypes[reflect.Int8]
	case rmeta.Int16:
		return gotypes[reflect.Int16]
	case rmeta.Int32:
		return gotypes[reflect.Int32]
	case rmeta.Int64:
		return gotypes[reflect.Int64]
	case rmeta.Float32:
		return gotypes[reflect.Float32]
	case rmeta.Float64:
		return gotypes[reflect.Float64]

	case rmeta.Counter:
		return gotypes[reflect.Int]

	case rmeta.OffsetL + rmeta.Bool:
		return reflect.ArrayOf(n, gotypes[reflect.Bool])
	case rmeta.OffsetL + rmeta.Uint8:
		return reflect.ArrayOf(n, gotypes[reflect.Uint8])
	case rmeta.OffsetL + rmeta.Uint16:
		return reflect.ArrayOf(n, gotypes[reflect.Uint16])
	case rmeta.OffsetL + rmeta.Uint32:
		return reflect.ArrayOf(n, gotypes[reflect.Uint32])
	case rmeta.OffsetL + rmeta.Uint64:
		return reflect.ArrayOf(n, gotypes[reflect.Uint64])
	case rmeta.OffsetL + rmeta.Int8:
		return reflect.ArrayOf(n, gotypes[reflect.Int8])
	case rmeta.OffsetL + rmeta.Int16:
		return reflect.ArrayOf(n, gotypes[reflect.Int16])
	case rmeta.OffsetL + rmeta.Int32:
		return reflect.ArrayOf(n, gotypes[reflect.Int32])
	case rmeta.OffsetL + rmeta.Int64:
		return reflect.ArrayOf(n, gotypes[reflect.Int64])
	case rmeta.OffsetL + rmeta.Float32:
		return reflect.ArrayOf(n, gotypes[reflect.Float32])
	case rmeta.OffsetL + rmeta.Float64:
		return reflect.ArrayOf(n, gotypes[reflect.Float64])

	case rmeta.OffsetP + rmeta.Bool:
		return reflect.SliceOf(gotypes[reflect.Bool])
	case rmeta.OffsetP + rmeta.Uint8:
		return reflect.SliceOf(gotypes[reflect.Uint8])
	case rmeta.OffsetP + rmeta.Uint16:
		return reflect.SliceOf(gotypes[reflect.Uint16])
	case rmeta.OffsetP + rmeta.Uint32:
		return reflect.SliceOf(gotypes[reflect.Uint32])
	case rmeta.OffsetP + rmeta.Uint64:
		return reflect.SliceOf(gotypes[reflect.Uint64])
	case rmeta.OffsetP + rmeta.Int8:
		return reflect.SliceOf(gotypes[reflect.Int8])
	case rmeta.OffsetP + rmeta.Int16:
		return reflect.SliceOf(gotypes[reflect.Int16])
	case rmeta.OffsetP + rmeta.Int32:
		return reflect.SliceOf(gotypes[reflect.Int32])
	case rmeta.OffsetP + rmeta.Int64:
		return reflect.SliceOf(gotypes[reflect.Int64])
	case rmeta.OffsetP + rmeta.Float32:
		return reflect.SliceOf(gotypes[reflect.Float32])
	case rmeta.OffsetP + rmeta.Float64:
		return reflect.SliceOf(gotypes[reflect.Float64])

	}
	panic("not implemented")
}

func genRStreamer(sictx rbytes.StreamerInfoContext, enum rmeta.Enum, n int, recv reflect.Value) rfunc {
	switch enum {
	case rmeta.Bool:
		return readBool
	case rmeta.Uint8:
		return readU8
	case rmeta.Uint16:
		return readU16
	case rmeta.Uint32:
		return readU32
	case rmeta.Uint64:
		return readU64
	case rmeta.Int8:
		return readI8
	case rmeta.Int16:
		return readI16
	case rmeta.Int32:
		return readI32
	case rmeta.Int64:
		return readI64
	case rmeta.Float32:
		return readF32
	case rmeta.Float64:
		return readF64

	case rmeta.Counter:
		return readInt
	}
	panic("not implemented")
}

func readBool(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*bool)) = r.ReadBool()
	return r.Err()
}

func readU8(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*uint8)) = r.ReadU8()
	return r.Err()
}

func readU16(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*uint16)) = r.ReadU16()
	return r.Err()
}

func readU32(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*uint32)) = r.ReadU32()
	return r.Err()
}

func readU64(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*uint64)) = r.ReadU64()
	return r.Err()
}

func readI8(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*int8)) = r.ReadI8()
	return r.Err()
}

func readI16(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*int16)) = r.ReadI16()
	return r.Err()
}

func readI32(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*int32)) = r.ReadI32()
	return r.Err()
}

func readI64(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*int64)) = r.ReadI64()
	return r.Err()
}

func readF32(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*float32)) = r.ReadF32()
	return r.Err()
}

func readF64(recv interface{}, r *rbytes.RBuffer) error {
	*(recv.(*float64)) = r.ReadF64()
	return r.Err()
}

func readInt(recv interface{}, r *rbytes.RBuffer) error {
	panic("not implemented")
	*(recv.(*int)) = int(r.ReadI64())
	return r.Err()
}

var (
	gotypes = map[reflect.Kind]reflect.Type{
		reflect.Bool:    reflect.TypeOf(false),
		reflect.Uint8:   reflect.TypeOf(uint8(0)),
		reflect.Uint16:  reflect.TypeOf(uint16(0)),
		reflect.Uint32:  reflect.TypeOf(uint32(0)),
		reflect.Uint64:  reflect.TypeOf(uint64(0)),
		reflect.Int8:    reflect.TypeOf(int8(0)),
		reflect.Int16:   reflect.TypeOf(int16(0)),
		reflect.Int32:   reflect.TypeOf(int32(0)),
		reflect.Int64:   reflect.TypeOf(int64(0)),
		reflect.Uint:    reflect.TypeOf(uint(0)),
		reflect.Int:     reflect.TypeOf(int(0)),
		reflect.Float32: reflect.TypeOf(float32(0)),
		reflect.Float64: reflect.TypeOf(float64(0)),
		reflect.String:  reflect.TypeOf(""),
	}
)

var (
	_ root.Object        = (*Object)(nil)
	_ rbytes.RVersioner  = (*Object)(nil)
	_ rbytes.Marshaler   = (*Object)(nil)
	_ rbytes.Unmarshaler = (*Object)(nil)
)

func init() {
	{
		f := func() reflect.Value {
			o := &Object{class: "*rdict.Object"}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("*rdict.Object", f)
	}
}
