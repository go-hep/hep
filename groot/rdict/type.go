// Copyright ©2020 The go-hep Authors. All rights reserved.
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

// TypeFromSI returns a Go type corresponding to the provided StreamerInfo.
// TypeFromSI first reaches out to the known groot types (via groot/rtypes) and
// then resorts to building a new type with reflect.
func TypeFromSI(ctx rbytes.StreamerInfoContext, si rbytes.StreamerInfo) (reflect.Type, error) {
	name := si.Name()
	if rtypes.Factory.HasKey(name) {
		fct := rtypes.Factory.Get(name)
		v := fct()
		return v.Type().Elem(), nil
	}

	switch {
	case name == "string", name == "std::string":
		if len(si.Elements()) == 0 {
			// fix for old (v=2) streamer for string
			sinfo := si.(*StreamerInfo)
			sinfo.elems = append(sinfo.elems, &StreamerSTLstring{
				StreamerSTL: StreamerSTL{
					StreamerElement: Element{
						Name:   *rbase.NewNamed("This", ""),
						Type:   rmeta.STLstring,
						Size:   32,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "string",
					}.New(),
					vtype: rmeta.ESTLType(rmeta.STLstring),
					ctype: rmeta.STLstring,
				},
			})
		}
		return gotypes[reflect.String], nil

	case hasStdPrefix(name,
		"vector", "list", "deque",
		"set", "multiset", "unordered_set", "unordered_multiset",
		"map", "multimap", "unordered_map", "unordered_multimap"):
		var (
			se      = si.Elements()[0]
			rt, err = TypeFromSE(ctx, se)
		)
		if err != nil {
			return nil, fmt.Errorf(
				"rdict: could not build element %q type for %q: %w",
				se.Name(), si.Name(), err,
			)
		}
		return rt, nil
	}

	fields := make([]reflect.StructField, 0, len(si.Elements()))
	for _, se := range si.Elements() {
		rt, err := TypeFromSE(ctx, se)
		if err != nil {
			return nil, fmt.Errorf(
				"rdict: could not build element %q type for %q: %w",
				se.Name(), si.Name(), err,
			)
		}
		et := se.Title()
		if !strings.HasPrefix(et, "[") {
			et = ""
		}
		if rt.Kind() == reflect.Array {
			et = ndimsFromType(rt)
		}
		ft := reflect.StructField{
			Name: "ROOT_" + cxxNameSanitizer.Replace(se.Name()),
			Type: rt,
			Tag:  reflect.StructTag(fmt.Sprintf("groot:%q", se.Name()+et)),
		}
		fields = append(fields, ft)
	}

	return reflect.StructOf(fields), nil
}

// TypeFromSE returns a Go type corresponding to the provided StreamerElement.
// TypeFromSE first reaches out to the known groot types (via groot/rtypes) and
// then resorts to building a new type with reflect.
func TypeFromSE(ctx rbytes.StreamerInfoContext, se rbytes.StreamerElement) (reflect.Type, error) {
	name := se.TypeName()
	name = strings.TrimRight(name, "*")
	if rtypes.Factory.HasKey(name) {
		var (
			fct = rtypes.Factory.Get(name)
			v   = fct()
			typ = v.Elem().Type()
		)

		return typeFromDescr(typ, se.TypeName(), se.ArrayLen(), se.ArrayDims()), nil
	}

	switch se := se.(type) {
	default:
		return nil, fmt.Errorf("rdict: unknown streamer element: %#v (%T)", se, se)

	case *StreamerBase:
		var (
			typename = se.Name()
			typevers = se.vbase
		)

		si, err := ctx.StreamerInfo(se.Name(), int(typevers))
		if err != nil {
			return nil, fmt.Errorf("rdict: could not find streamer info for base %q: %w", typename, err)
		}
		return TypeFromSI(ctx, si)

	case *StreamerBasicType:
		return typeFrom(ctx, se.TypeName(), se.Type(), se.Size(), se.ArrayLen(), se.ArrayDims())

	case *StreamerString:
		return typeFrom(ctx, se.TypeName(), se.Type(), se.Size(), se.ArrayLen(), se.ArrayDims())

	case *StreamerBasicPointer:
		return typeFrom(ctx, se.TypeName(), se.Type(), se.Size(), -1, se.ArrayDims())

	case *StreamerSTLstring:
		return gotypes[reflect.String], nil

	case *StreamerLoop:
		var (
			typename = se.TypeName()
			typevers = int16(-1)
		)
		typename = typename[:len(typename)-1] // drop final '*'
		elt, err := typeFromTypeName(ctx, typename, typevers, se.Type(), se, 1)
		if err != nil {
			return nil, fmt.Errorf(
				"rdict: could not find type of looper %q: %w",
				typename, err,
			)
		}
		return reflect.SliceOf(elt), nil

	case *StreamerObject, *StreamerObjectAny:
		var (
			alen     = se.ArrayLen()
			typename = se.TypeName()
			typevers = -1
			si, err  = ctx.StreamerInfo(typename, typevers)
		)
		if err != nil {
			return nil, fmt.Errorf("rdict: could not find streamer info for type %q: %w", typename, err)
		}

		typ, err := TypeFromSI(ctx, si)
		if err != nil {
			return nil, fmt.Errorf("rdict: could not build type for %q: %w", typename, err)
		}
		return typeFromDescr(typ, typename, alen, se.ArrayDims()), nil

	case *StreamerObjectPointer, *StreamerObjectAnyPointer:
		var (
			alen     = se.ArrayLen()
			typename = se.TypeName()
			typevers = -1
		)
		typename = typename[:len(typename)-1] // drop final '*'

		si, err := ctx.StreamerInfo(typename, typevers)
		if err != nil {
			return nil, fmt.Errorf("rdict: could not find streamer info for ptr-to-object %q: %w", typename, err)
		}

		typ, err := TypeFromSI(ctx, si)
		if err != nil {
			return nil, fmt.Errorf("rdict: could not create type for ptr-to-object %q: %w", typename, err)
		}
		typ = reflect.PtrTo(typ)
		return typeFromDescr(typ, typename, alen, se.ArrayDims()), nil

	case *StreamerSTL:
		switch se.STLType() {
		case rmeta.STLvector, rmeta.STLlist, rmeta.STLdeque:
			var (
				ct       = se.ContainedType()
				ename    = rmeta.CxxTemplateFrom(se.TypeName()).Args[0]
				typevers = int16(-1)
				elt, err = typeFromTypeName(ctx, ename, typevers, ct, se, 1)
			)
			if err != nil {
				return nil, fmt.Errorf("rdict: could not create type for %q: %w", se.TypeName(), err)
			}
			return reflect.SliceOf(elt), nil

		case rmeta.STLset, rmeta.STLunorderedset, rmeta.STLmultiset, rmeta.STLunorderedmultiset:
			var (
				ct       = se.ContainedType()
				typename = se.TypeName()
				typevers = int16(-1)
				kname    = rmeta.CxxTemplateFrom(se.TypeName()).Args[0]
			)

			key, err := typeFromTypeName(ctx, kname, typevers, ct, se, 1)
			if err != nil {
				return nil, fmt.Errorf(
					"could not find key type %q for std::{,multi,unordered_}set %q: %w", kname, typename, err,
				)
			}
			return reflect.MapOf(key, reflect.TypeOf(struct{}{})), nil

		case rmeta.STLmap, rmeta.STLunorderedmap, rmeta.STLmultimap, rmeta.STLunorderedmultimap:
			var (
				ct       = se.ContainedType()
				typename = se.TypeName()
				typevers = int16(-1)
				enames   = rmeta.CxxTemplateFrom(se.TypeName()).Args
				kname    = enames[0]
				vname    = enames[1]
			)

			key, err := typeFromTypeName(ctx, kname, typevers, ct, se, 1)
			if err != nil {
				return nil, fmt.Errorf(
					"could not find key type %q for std::{,multi,unordered_}map %q: %w", kname, typename, err,
				)
			}
			val, err := typeFromTypeName(ctx, vname, typevers, ct, se, 1)
			if err != nil {
				return nil, fmt.Errorf(
					"could not find val type %q for std::{,multi,unordered_}map %q: %w", vname, typename, err,
				)
			}
			return reflect.MapOf(key, val), nil

		case rmeta.STLbitset:
			var (
				typename = se.TypeName()
				enames   = rmeta.CxxTemplateFrom(typename).Args
				_, err   = strconv.Atoi(enames[0])
			)
			if err != nil {
				return nil, fmt.Errorf(
					"could not infer bitset argument (type=%q): %w", typename, err,
				)
			}
			// FIXME(sbinet): use a fixed-sized array rounded-up to n/8-bytes
			//	bits2bytes := func(v int) int {
			//		const len = 8-1
			//		return (v + (8-(sz&len))&len)/8
			//	}
			//	n := bits2bytes(v)
			return reflect.SliceOf(gotypes[reflect.Uint8]), nil

		default:
			return nil, fmt.Errorf("rdict: STL container not implemented: %#v (vtype=%+v)", se, se.STLType())
		}
	}
}

func typeFrom(ctx rbytes.StreamerInfoContext, typename string, enum rmeta.Enum, size uintptr, n int, dims []int32) (reflect.Type, error) {
	var rt reflect.Type

	switch enum {
	case rmeta.Bool:
		rt = gotypes[reflect.Bool]
	case rmeta.Uint8:
		rt = gotypes[reflect.Uint8]
	case rmeta.Uint16:
		rt = gotypes[reflect.Uint16]
	case rmeta.Uint32, rmeta.Bits:
		rt = gotypes[reflect.Uint32]
	case rmeta.Uint64, rmeta.ULong64:
		rt = gotypes[reflect.Uint64]
	case rmeta.Int8:
		rt = gotypes[reflect.Int8]
	case rmeta.Int16:
		rt = gotypes[reflect.Int16]
	case rmeta.Int32:
		rt = gotypes[reflect.Int32]
	case rmeta.Int64, rmeta.Long64:
		rt = gotypes[reflect.Int64]
	case rmeta.Float32:
		rt = gotypes[reflect.Float32]
	case rmeta.Float64:
		rt = gotypes[reflect.Float64]
	case rmeta.Float16:
		rt = reflect.TypeOf((*root.Float16)(nil)).Elem()
	case rmeta.Double32:
		rt = reflect.TypeOf((*root.Double32)(nil)).Elem()
	case rmeta.TString, rmeta.STLstring:
		rt = gotypes[reflect.String]

	case rmeta.CharStar:
		rt = gotypes[reflect.String]

	case rmeta.Counter:
		switch size {
		case 4:
			rt = gotypes[reflect.Int32]
		case 8:
			rt = gotypes[reflect.Int64]
		default:
			return nil, fmt.Errorf("rdict: invalid counter size=%d", size)
		}

	case rmeta.TObject:
		rt = reflect.TypeOf((*rbase.Object)(nil)).Elem()

	case rmeta.TNamed:
		rt = reflect.TypeOf((*rbase.Named)(nil)).Elem()

	case rmeta.OffsetL + rmeta.Bool:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Bool]
	case rmeta.OffsetL + rmeta.Uint8:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Uint8]
	case rmeta.OffsetL + rmeta.Uint16:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Uint16]
	case rmeta.OffsetL + rmeta.Uint32:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Uint32]
	case rmeta.OffsetL + rmeta.Uint64, rmeta.OffsetL + rmeta.ULong64:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Uint64]
	case rmeta.OffsetL + rmeta.Int8:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Int8]
	case rmeta.OffsetL + rmeta.Int16:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Int16]
	case rmeta.OffsetL + rmeta.Int32:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Int32]
	case rmeta.OffsetL + rmeta.Int64, rmeta.OffsetL + rmeta.Long64:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Int64]
	case rmeta.OffsetL + rmeta.Float32:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Float32]
	case rmeta.OffsetL + rmeta.Float64:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.Float64]
	case rmeta.OffsetL + rmeta.Float16:
		// dim handled by typeFromDescr.
		rt = reflect.TypeOf((*root.Float16)(nil)).Elem()
	case rmeta.OffsetL + rmeta.Double32:
		// dim handled by typeFromDescr.
		rt = reflect.TypeOf((*root.Double32)(nil)).Elem()
	case rmeta.OffsetL + rmeta.TString,
		rmeta.OffsetL + rmeta.CharStar,
		rmeta.OffsetL + rmeta.STLstring:
		// dim handled by typeFromDescr.
		rt = gotypes[reflect.String]

	case rmeta.OffsetP + rmeta.Bool:
		rt = reflect.SliceOf(gotypes[reflect.Bool])
	case rmeta.OffsetP + rmeta.Uint8:
		rt = reflect.SliceOf(gotypes[reflect.Uint8])
	case rmeta.OffsetP + rmeta.Uint16:
		rt = reflect.SliceOf(gotypes[reflect.Uint16])
	case rmeta.OffsetP + rmeta.Uint32:
		rt = reflect.SliceOf(gotypes[reflect.Uint32])
	case rmeta.OffsetP + rmeta.Uint64, rmeta.OffsetP + rmeta.ULong64:
		rt = reflect.SliceOf(gotypes[reflect.Uint64])
	case rmeta.OffsetP + rmeta.Int8:
		rt = reflect.SliceOf(gotypes[reflect.Int8])
	case rmeta.OffsetP + rmeta.Int16:
		rt = reflect.SliceOf(gotypes[reflect.Int16])
	case rmeta.OffsetP + rmeta.Int32:
		rt = reflect.SliceOf(gotypes[reflect.Int32])
	case rmeta.OffsetP + rmeta.Int64, rmeta.OffsetP + rmeta.Long64:
		rt = reflect.SliceOf(gotypes[reflect.Int64])
	case rmeta.OffsetP + rmeta.Float32:
		rt = reflect.SliceOf(gotypes[reflect.Float32])
	case rmeta.OffsetP + rmeta.Float64:
		rt = reflect.SliceOf(gotypes[reflect.Float64])
	case rmeta.OffsetP + rmeta.Float16:
		rt = reflect.SliceOf(reflect.TypeOf((*root.Float16)(nil)).Elem())
	case rmeta.OffsetP + rmeta.Double32:
		rt = reflect.SliceOf(reflect.TypeOf((*root.Double32)(nil)).Elem())
	case rmeta.OffsetP + rmeta.STLstring,
		rmeta.OffsetP + rmeta.CharStar:
		rt = reflect.SliceOf(gotypes[reflect.String])
	}

	if rt == nil {
		return nil, fmt.Errorf("rmeta=%d (%v) not implemented (size=%d, n=%v)", enum, enum, size, n)
	}

	return typeFromDescr(rt, typename, n, dims), nil
}

func typeFromTypeName(ctx rbytes.StreamerInfoContext, typename string, typevers int16, enum rmeta.Enum, se rbytes.StreamerElement, n int) (reflect.Type, error) {
	e, ok := rmeta.TypeName2Enum(typename)
	if ok {
		return typeFrom(ctx, typename, e, se.Size(), n, se.ArrayDims())
	}

	switch {
	case hasStdPrefix(typename, "vector", "list", "deque"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		et, err := typeFromTypeName(ctx, enames[0], -1, -1, se, n)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(et), nil

	case hasStdPrefix(typename, "set", "multiset", "unordered_set", "unordered_multiset"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		kname := enames[0]

		kt, err := typeFromTypeName(ctx, kname, -1, -1, se, n)
		if err != nil {
			return nil, err
		}
		return reflect.MapOf(kt, reflect.TypeOf(struct{}{})), nil

	case hasStdPrefix(typename, "map", "multimap", "unordered_map", "unordered_multimap"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		kname := enames[0]
		vname := enames[1]

		kt, err := typeFromTypeName(ctx, kname, -1, -1, se, n)
		if err != nil {
			return nil, err
		}
		vt, err := typeFromTypeName(ctx, vname, -1, -1, se, n)
		if err != nil {
			return nil, err
		}
		return reflect.MapOf(kt, vt), nil

	case hasStdPrefix(typename, "pair"):
		enames := rmeta.CxxTemplateFrom(typename).Args
		p0 := enames[0]
		p1 := enames[1]
		t0, err := typeFromTypeName(ctx, p0, -1, -1, se, n)
		if err != nil {
			return nil, err
		}
		t1, err := typeFromTypeName(ctx, p1, -1, -1, se, n)
		if err != nil {
			return nil, err
		}
		return reflect.StructOf([]reflect.StructField{
			{
				Name: "ROOT_first",
				Type: t0,
				Tag:  reflect.StructTag(`groot:"first"`),
			},
			{
				Name: "ROOT_second",
				Type: t1,
				Tag:  reflect.StructTag(`groot:"second"`),
			},
		}), nil

	case hasStdPrefix(typename, "bitset"):
		var (
			enames = rmeta.CxxTemplateFrom(typename).Args
			_, err = strconv.Atoi(enames[0])
		)

		if err != nil {
			return nil, fmt.Errorf("rdict: invalid STL bitset argument (type=%q): %+v", typename, err)
		}
		return reflect.SliceOf(gotypes[reflect.Uint8]), nil
	}

	osi, err := ctx.StreamerInfo(typename, int(typevers))
	if err != nil {
		return nil, fmt.Errorf("rdict: could not find streamer info for %q (version=%d): %w", typename, typevers, err)
	}

	return TypeFromSI(ctx, osi)
}

func typeFromDescr(typ reflect.Type, typename string, alen int, dims []int32) reflect.Type {
	if alen > 0 {
		// handle [n][m][u][v][w]T
		ndim := len(dims)
		for i := range dims {
			typ = reflect.ArrayOf(int(dims[ndim-1-i]), typ)
		}
		return typ
	}

	if alen < 0 {
		// slice. drop one '*' from typename.
		typename = strings.TrimSuffix(typename, "*")
	}
	if typename == "char*" {
		// slice. drop one '*' from typename.
		typename = strings.TrimSuffix(typename, "*")
	}

	// handle T***
	for i := range typename {
		if typename[len(typename)-1-i] != '*' {
			break
		}
		typ = reflect.PtrTo(typ)
	}

	return typ
}

func ndimsFromType(rt reflect.Type) string {
	var dims []int
	for rt.Kind() == reflect.Array {
		dims = append(dims, rt.Len())
		rt = rt.Elem()
	}
	var o strings.Builder
	for _, v := range dims {
		o.WriteString("[" + strconv.Itoa(v) + "]")
	}
	return o.String()
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

func hasStdPrefix(typename string, ps ...string) bool {
	for _, p := range ps {
		switch {
		case strings.HasPrefix(typename, p+"<"),
			strings.HasPrefix(typename, "std::"+p+"<"):
			return true
		}
	}
	return false
}
