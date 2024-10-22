// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rvers"
)

// StreamerOf generates a StreamerInfo from a reflect.Type.
//
// StreamerOf panics if the provided type contains non-ROOT compatible types
// such as chan, int, uint or func.
func StreamerOf(ctx rbytes.StreamerInfoContext, typ reflect.Type) rbytes.StreamerInfo {
	if isTObject(typ) {
		name := reflect.New(typ).Elem().Interface().(root.Object).Class()
		si, err := ctx.StreamerInfo(name, -1)
		if err == nil {
			return si
		}
	}

	bldr := newStreamerBuilder(ctx, typ)
	return bldr.genStreamer(typ)
}

type streamerStore interface {
	rbytes.StreamerInfoContext
	addStreamer(si rbytes.StreamerInfo)
}

type streamerBuilder struct {
	ctx streamerStore
	typ reflect.Type
}

func newStreamerBuilder(ctx rbytes.StreamerInfoContext, typ reflect.Type) *streamerBuilder {
	return &streamerBuilder{ctx: newStreamerStore(ctx), typ: typ}
}

func (bld *streamerBuilder) genStreamer(typ reflect.Type) rbytes.StreamerInfo {
	name := typenameOf(typ)
	si := &StreamerInfo{
		named:  *rbase.NewNamed(name, name),
		objarr: rcont.NewObjArray(),
		clsver: 1,
	}
	switch typ.Kind() {
	case reflect.Struct:
		si.elems = make([]rbytes.StreamerElement, 0, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			ft := typ.Field(i)
			si.elems = append(si.elems, bld.genField(typ, ft))
		}
	case reflect.Slice:
		si.clsver = rvers.StreamerInfo
		si.elems = []rbytes.StreamerElement{
			bld.genStdVectorOf(typ.Elem(), "This", 0),
		}
	}
	return si
}

func (bld *streamerBuilder) genStdVectorOf(typ reflect.Type, name string, offset int32) rbytes.StreamerElement {
	const esize = 3 * diskPtrSize
	var (
		ename = ""
		etype rmeta.Enum
	)
	switch typ.Kind() {
	case reflect.Bool:
		ename = "vector<bool>"
		etype = rmeta.Bool
	case reflect.Int8:
		ename = "vector<int8_t>"
		etype = rmeta.Int8
	case reflect.Int16:
		ename = "vector<int16_t>"
		etype = rmeta.Int16
	case reflect.Int32:
		ename = "vector<int32_t>"
		etype = rmeta.Int32
	case reflect.Int64:
		ename = "vector<int64_t>"
		etype = rmeta.Int64
	case reflect.Uint8:
		ename = "vector<uint8_t>"
		etype = rmeta.Uint8
	case reflect.Uint16:
		ename = "vector<uint16_t>"
		etype = rmeta.Uint16
	case reflect.Uint32:
		ename = "vector<uint32_t>"
		etype = rmeta.Uint32
	case reflect.Uint64:
		ename = "vector<uint64_t>"
		etype = rmeta.Uint64
	case reflect.Float32:
		switch typ {
		case reflect.TypeOf(root.Float16(0)):
			ename = "vector<Float16_t>"
			etype = rmeta.Float16
		default:
			ename = "vector<float>"
			etype = rmeta.Float32
		}
	case reflect.Float64:
		switch typ {
		case reflect.TypeOf(root.Double32(0)):
			ename = "vector<Double32_t>"
			etype = rmeta.Double32
		default:
			ename = "vector<double>"
			etype = rmeta.Float64
		}
	case reflect.String:
		ename = "vector<string>"
		etype = rmeta.STLstring
	case reflect.Struct:
		ename = fmt.Sprintf("vector<%s>", typenameOf(typ))
		etype = rmeta.Any
		if isTObject(typ) || isTObject(reflect.PointerTo(typ)) {
			etype = rmeta.Object
		}
	case reflect.Slice:
		ename = typenameOf(typ)
		if strings.HasSuffix(ename, ">") {
			ename += " "
		}
		ename = fmt.Sprintf("vector<%s>", ename)
		etype = rmeta.Any
	default:
		panic(fmt.Errorf("rdict: invalid slice type %v", typ))
	}

	return NewCxxStreamerSTL(
		StreamerElement{
			named:  *rbase.NewNamed(name, ""),
			etype:  rmeta.Streamer,
			esize:  esize,
			offset: offset,
			ename:  ename,
		}, rmeta.STLvector, etype,
	)
}

func (bld *streamerBuilder) genPtr(typ reflect.Type, name string, offset int32) rbytes.StreamerElement {
	// FIXME(sbinet): is typ always a struct?
	//	switch typ.Kind() {
	//	case reflect.Struct:
	//	default:
	//		panic(fmt.Errorf("rdict: invalid ptr-to type %v", typ))
	//	}

	ptr := reflect.PointerTo(typ)
	se := StreamerElement{
		named:  *rbase.NewNamed(name, ""),
		etype:  rmeta.AnyP,
		esize:  int32(ptrSize),
		offset: offset,
		ename:  typenameOf(ptr),
	}

	if isTObject(ptr) {
		se.etype = rmeta.ObjectP
		return &StreamerObjectPointer{se}
	}

	return &StreamerObjectAnyPointer{se}
}

func (bld *streamerBuilder) genArrayOf(n int, typ reflect.Type, name string, offset int32) rbytes.StreamerElement {
	var (
		arrlen = int32(n)
		dims   = []int32{arrlen}
		maxidx [5]int32
		esize  = arrlen
		etype  rmeta.Enum
		ename  = ""
	)

	for typ.Kind() == reflect.Array && len(dims) < 5 {
		dim := int32(typ.Len())
		dims = append(dims, dim)
		typ = typ.Elem()
		esize *= dim
	}
	copy(maxidx[:], dims)
	arrdim := int32(len(dims))

	switch typ.Kind() {
	case reflect.Bool:
		esize *= 1
		ename = "bool"
		etype = rmeta.Bool

	case reflect.Int8:
		esize *= 1
		ename = "int8_t"
		etype = rmeta.Int8

	case reflect.Int16:
		esize *= 2
		ename = "int16_t"
		etype = rmeta.Int16

	case reflect.Int32:
		esize *= 4
		ename = "int32_t"
		etype = rmeta.Int32

	case reflect.Int64:
		esize *= 8
		ename = "int64_t"
		etype = rmeta.Int64

	case reflect.Uint8:
		esize *= 1
		ename = "uint8_t"
		etype = rmeta.Uint8

	case reflect.Uint16:
		esize *= 2
		ename = "uint16_t"
		etype = rmeta.Uint16

	case reflect.Uint32:
		esize *= 4
		ename = "uint32_t"
		etype = rmeta.Uint32

	case reflect.Uint64:
		esize *= 8
		ename = "uint64_t"
		etype = rmeta.Uint64

	case reflect.Float32:
		switch typ {
		case reflect.TypeOf(root.Float16(0)):
			esize *= 4
			ename = "Float16_t"
			etype = rmeta.Float16
		default:
			esize *= 4
			ename = "float"
			etype = rmeta.Float32
		}

	case reflect.Float64:
		switch typ {
		case reflect.TypeOf(root.Double32(0)):
			esize *= 8
			ename = "Double32_t"
			etype = rmeta.Double32
		default:
			esize *= 8
			ename = "double"
			etype = rmeta.Float64
		}

	case reflect.String:
		return &StreamerString{
			StreamerElement{
				named:  *rbase.NewNamed(name, ""),
				etype:  rmeta.OffsetL + rmeta.TString,
				esize:  esize * sizeOfTString,
				arrlen: arrlen,
				arrdim: arrdim,
				maxidx: maxidx,
				offset: offset,
				ename:  "TString",
			},
		}
	case reflect.Struct:
		if isTObject(typ) || isTObject(reflect.PointerTo(typ)) {
			return &StreamerObject{
				StreamerElement{
					named:  *rbase.NewNamed(name, ""),
					etype:  rmeta.OffsetL + rmeta.Object,
					esize:  esize * bld.sizeOf(typ),
					arrlen: arrlen,
					arrdim: arrdim,
					maxidx: maxidx,
					offset: offset,
					ename:  typenameOf(typ),
				},
			}
		}
		return &StreamerObjectAny{
			StreamerElement{
				named:  *rbase.NewNamed(name, ""),
				etype:  rmeta.OffsetL + rmeta.Any,
				esize:  esize * bld.sizeOf(typ),
				arrlen: arrlen,
				arrdim: arrdim,
				maxidx: maxidx,
				offset: offset,
				ename:  typenameOf(typ),
			},
		}
	default:
		panic(fmt.Errorf("rdict: invalid array element type %v", typ))
	}

	return &StreamerBasicType{
		StreamerElement{
			named:  *rbase.NewNamed(name, ""),
			etype:  rmeta.OffsetL + etype,
			esize:  esize,
			offset: offset,
			arrlen: arrlen,
			arrdim: arrdim,
			maxidx: maxidx,
			ename:  ename,
		},
	}
}

func (bld *streamerBuilder) genVarLenArrayOf(typ reflect.Type, class, count, name string, offset int32) rbytes.StreamerElement {
	var (
		esize = 0
		ename = ""
		etype rmeta.Enum
	)
	switch typ.Kind() {
	case reflect.Bool:
		esize = 1
		ename = "bool"
		etype = rmeta.Bool
	case reflect.Int8:
		esize = 1
		ename = "int8_t"
		etype = rmeta.Int8
	case reflect.Int16:
		esize = 2
		ename = "int16_t"
		etype = rmeta.Int16
	case reflect.Int32:
		esize = 4
		ename = "int32_t"
		etype = rmeta.Int32
	case reflect.Int64:
		esize = 8
		ename = "int64_t"
		etype = rmeta.Int64
	case reflect.Uint8:
		esize = 1
		ename = "uint8_t"
		etype = rmeta.Uint8
	case reflect.Uint16:
		esize = 2
		ename = "uint16_t"
		etype = rmeta.Uint16
	case reflect.Uint32:
		esize = 4
		ename = "uint32_t"
		etype = rmeta.Uint32
	case reflect.Uint64:
		esize = 8
		ename = "uint64_t"
		etype = rmeta.Uint64
	case reflect.Float32:
		esize = 4
		switch typ {
		case reflect.TypeOf(root.Float16(0)):
			ename = "Float16_t"
			etype = rmeta.Float16
		default:
			ename = "float"
			etype = rmeta.Float32
		}
	case reflect.Float64:
		esize = 8
		switch typ {
		case reflect.TypeOf(root.Double32(0)):
			ename = "Double32_t"
			etype = rmeta.Double32
		default:
			ename = "double"
			etype = rmeta.Float64
		}
	case reflect.String:
		return NewStreamerLoop(
			StreamerElement{
				named:  *rbase.NewNamed(name, "["+count+"]"),
				esize:  4,
				offset: offset,
				ename:  "TString*",
			},
			1, count, class,
		)
	case reflect.Struct:
		ename = typenameOf(typ)
		return NewStreamerLoop(
			StreamerElement{
				named:  *rbase.NewNamed(name, "["+count+"]"),
				esize:  diskPtrSize,
				offset: offset,
				ename:  ename + "*",
			},
			1, count, class,
		)
	default:
		panic(fmt.Errorf("rdict: invalid c-var-len-array type %v", typ))
	}

	return NewStreamerBasicPointer(
		StreamerElement{
			named:  *rbase.NewNamed(name, "["+count+"]"),
			etype:  rmeta.OffsetP + etype,
			esize:  int32(esize),
			offset: offset,
			ename:  ename + "*",
		}, 1, count, class,
	)
}

func (bld *streamerBuilder) genField(typ reflect.Type, field reflect.StructField) rbytes.StreamerElement {

	offset := offsetOf(field)

	switch field.Type.Kind() {
	case reflect.Bool:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  1,
				offset: offset,
				ename:  "bool",
			},
		}
	case reflect.Int8:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  1,
				offset: offset,
				ename:  "int8_t",
			},
		}
	case reflect.Int16:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  2,
				offset: offset,
				ename:  "int16_t",
			},
		}
	case reflect.Int32:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  4,
				offset: offset,
				ename:  "int32_t",
			},
		}
	case reflect.Int64:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  8,
				offset: offset,
				ename:  "int64_t",
			},
		}
	case reflect.Uint8:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  1,
				offset: offset,
				ename:  "uint8_t",
			},
		}
	case reflect.Uint16:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  2,
				offset: offset,
				ename:  "uint16_t",
			},
		}
	case reflect.Uint32:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  4,
				offset: offset,
				ename:  "uint32_t",
			},
		}
	case reflect.Uint64:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  8,
				offset: offset,
				ename:  "uint64_t",
			},
		}
	case reflect.Float32:
		switch field.Type {
		case reflect.TypeOf(root.Float16(0)):
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Float16,
					esize:  4,
					offset: offset,
					ename:  "Float16_t",
				},
			}
		default:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.GoType2ROOTEnum[field.Type],
					esize:  4,
					offset: offset,
					ename:  "float",
				},
			}
		}
	case reflect.Float64:
		switch field.Type {
		case reflect.TypeOf(root.Double32(0)):
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Double32,
					esize:  8,
					offset: offset,
					ename:  "Double32_t",
				},
			}
		default:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.GoType2ROOTEnum[field.Type],
					esize:  8,
					offset: offset,
					ename:  "double",
				},
			}
		}
	case reflect.String:
		return &StreamerSTLstring{
			StreamerSTL: StreamerSTL{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Streamer,
					esize:  sizeOfStdString,
					offset: offset,
					ename:  "string",
				},
			},
		}
	case reflect.Struct:
		return &StreamerObjectAny{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.Any,
				esize:  bld.sizeOf(field.Type),
				offset: offset,
				ename:  typenameOf(field.Type),
			},
		}

	case reflect.Array:
		var (
			et   = field.Type.Elem()
			n    = field.Type.Len()
			name = nameOf(field)
		)
		return bld.genArrayOf(n, et, name, offsetOf(field))

	case reflect.Slice:
		et := field.Type.Elem()
		count, ok := hasCount(field)
		if ok {
			class := typenameOf(typ)
			return bld.genVarLenArrayOf(et, class, count, nameOf(field), offsetOf(field))
		}
		return bld.genStdVectorOf(et, nameOf(field), offsetOf(field))

	case reflect.Ptr:
		et := field.Type.Elem()
		return bld.genPtr(et, nameOf(field), offsetOf(field))

	default:
		panic(fmt.Errorf(
			"rdict: invalid struct field (name=%v, type=%v, kind=%v)",
			field.Name, field.Type, field.Type.Kind(),
		))
	}
}

type streamerStoreImpl struct {
	ctx rbytes.StreamerInfoContext
	db  map[string]rbytes.StreamerInfo
}

func newStreamerStore(ctx rbytes.StreamerInfoContext) streamerStore {
	if ctx, ok := ctx.(streamerStore); ok {
		return ctx
	}

	return &streamerStoreImpl{
		ctx: ctx,
		db:  make(map[string]rbytes.StreamerInfo),
	}
}

// StreamerInfo returns the named StreamerInfo.
// If version is negative, the latest version should be returned.
func (store *streamerStoreImpl) StreamerInfo(name string, version int) (rbytes.StreamerInfo, error) {
	return store.ctx.StreamerInfo(name, version)
}

func (store *streamerStoreImpl) addStreamer(si rbytes.StreamerInfo) {
	store.db[si.Name()] = si
}

func nameOf(field reflect.StructField) string {
	tag, ok := field.Tag.Lookup("groot")
	if ok {
		i := strings.Index(tag, "[")
		if i < 0 {
			return tag
		}
		return tag[:i]
	}
	return field.Name
}

func hasCount(field reflect.StructField) (string, bool) {
	tag, ok := field.Tag.Lookup("groot")
	if !ok || !strings.Contains(tag, "[") {
		return "", false
	}
	var (
		count = new(strings.Builder)
		brack bool
	)
	for _, v := range tag {
		switch v {
		case '[':
			brack = true
		case ']':
			name := strings.TrimSpace(count.String())
			if !isIdent(name) {
				return "", false
			}
			return name, true
		default:
			if brack {
				_, _ = count.WriteString(string(v))
			}
		}
	}
	return "", false
}

func isIdent(name string) bool {
	if len(name) == 0 {
		return false
	}

	ok := func(c rune) bool {
		return 'A' <= c && c <= 'Z' ||
			'a' <= c && c <= 'z' ||
			'0' <= c && c <= '9' ||
			c == '_'
	}

	for i, c := range name {
		if i == 0 && ('0' <= c && c <= '9') {
			return false
		}
		if !ok(c) {
			return false
		}
	}
	return true
}

func offsetOf(field reflect.StructField) int32 {
	// return int32(field.Offset)
	// FIXME(sbinet): it seems ROOT expects 0 here...
	return 0
}

func (bld *streamerBuilder) sizeOf(typ reflect.Type) int32 {
	// FIXME(sbinet): compute ROOT-compatible size.
	if ptr := reflect.PointerTo(typ); isTObject(ptr) || isTObject(typ) {
		name := typenameOf(typ)
		switch name {
		case "TObjString":
			return sizeOfTObjString
		}
	}
	return int32(typ.Size())
}

func isTObject(typ reflect.Type) bool {
	return typ.Implements(rootObjectIface)
}

func typenameOf(typ reflect.Type) string {
	if isTObject(typ) {
		switch typ.Kind() {
		case reflect.Ptr:
			name := reflect.New(typ.Elem()).Interface().(root.Object).Class()
			return name + "*"
		default:
			name := reflect.New(typ).Elem().Interface().(root.Object).Class()
			return name
		}
	}
	if ptr := reflect.PointerTo(typ); isTObject(ptr) {
		name := reflect.New(typ).Interface().(root.Object).Class()
		return name
	}

	switch typ.Kind() {
	case reflect.Slice:
		ename := typenameOf(typ.Elem())
		if strings.HasSuffix(ename, ">") {
			ename += " "
		}
		return "vector<" + ename + ">"
	case reflect.Array:
		var (
			dims []int
			et   = typ
		)
		for i := 0; i < 10; i++ {
			dims = append(dims, et.Len())
			et = et.Elem()
			if et.Kind() != reflect.Array {
				break
			}

		}
		o := new(strings.Builder)
		ename := typenameOf(et)
		_, _ = o.WriteString(ename)
		for _, v := range dims {
			fmt.Fprintf(o, "[%d]", v)
		}
		return o.String()
	case reflect.Bool:
		return "bool"
	case reflect.Int8:
		return "int8_t"
	case reflect.Int16:
		return "int16_t"
	case reflect.Int32:
		return "int32_t"
	case reflect.Int64:
		return "int64_t"
	case reflect.Uint8:
		return "uint8_t"
	case reflect.Uint16:
		return "uint16_t"
	case reflect.Uint32:
		return "uint32_t"
	case reflect.Uint64:
		return "uint64_t"
	case reflect.Float32:
		switch typ {
		case reflect.TypeOf(root.Float16(0)):
			return "Float16_t"
		default:
			return "float"
		}
	case reflect.Float64:
		switch typ {
		case reflect.TypeOf(root.Double32(0)):
			return "Double32_t"
		default:
			return "double"
		}
	case reflect.String:
		return "string"
	case reflect.Ptr:
		return typenameOf(typ.Elem()) + "*"

	default:
		name := typ.Name()
		if name == "" {
			panic(fmt.Errorf("rdict: invalid reflect type %v", typ))
		}
		return name
	}
}

const (
	diskPtrSize      = 8
	sizeOfTObjString = 40
	sizeOfTString    = 3 * diskPtrSize
	sizeOfStdString  = 4 * diskPtrSize
)

var (
	rootObjectIface = reflect.TypeOf((*root.Object)(nil)).Elem()
)

var (
	_ streamerStore              = (*streamerStoreImpl)(nil)
	_ rbytes.StreamerInfoContext = (*streamerStoreImpl)(nil)
)
