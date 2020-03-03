// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
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

	case strings.HasPrefix(name, "vector<"),
		strings.HasPrefix(name, "map<"): // FIXME(sbinet): handle other std::containers?
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
		if rt.Kind() == reflect.Array {
			et = fmt.Sprintf("[%d]", rt.Len())
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

		switch alen := se.ArrayLen(); alen {
		case 0:
			return typ, nil
		default:
			return reflect.ArrayOf(alen, typ), nil
		}
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
		return typeFrom(ctx, se.Type(), se.Size(), se.ArrayLen())

	case *StreamerString:
		return typeFrom(ctx, se.Type(), se.Size(), se.ArrayLen())

	case *StreamerBasicPointer:
		return typeFrom(ctx, se.Type(), se.Size(), -1)

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

		if alen > 0 {
			return reflect.ArrayOf(se.ArrayLen(), typ), nil
		}
		return typ, nil

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
		ptr := reflect.PtrTo(typ)

		if alen > 0 {
			return reflect.ArrayOf(alen, ptr), nil
		}
		return ptr, nil

	case *StreamerSTL:
		switch se.STLType() {
		case rmeta.STLvector:
			var (
				ct       = se.ContainedType()
				typename = rmeta.CxxTemplateArgsOf(se.TypeName())[0]
				typevers = int16(-1)
				elt, err = typeFromTypeName(ctx, typename, typevers, ct, se, 1)
			)
			if err != nil {
				return nil, fmt.Errorf("rdict: could not create type for std::vector<T> T=%q: %w", typename, err)
			}
			return reflect.SliceOf(elt), nil

		case rmeta.STLmap, rmeta.STLunorderedmap:
			var (
				ct       = se.ContainedType()
				typename = se.TypeName()
				typevers = int16(-1)
				enames   = rmeta.CxxTemplateArgsOf(se.TypeName())
				kname    = enames[0]
				vname    = enames[1]
			)

			key, err := typeFromTypeName(ctx, kname, typevers, ct, se, 1)
			if err != nil {
				return nil, fmt.Errorf(
					"could not find key type %q for std::map %q: %w", kname, typename, err,
				)
			}
			val, err := typeFromTypeName(ctx, vname, typevers, ct, se, 1)
			if err != nil {
				return nil, fmt.Errorf(
					"could not find val type %q for std::map %q: %w", vname, typename, err,
				)
			}
			return reflect.MapOf(key, val), nil
		default:
			return nil, fmt.Errorf("rdict: STL container not implemented: %#v", se)
		}
	}
}

func typeFrom(ctx rbytes.StreamerInfoContext, enum rmeta.Enum, size uintptr, n int) (reflect.Type, error) {
	switch enum {
	case rmeta.Bool:
		return gotypes[reflect.Bool], nil
	case rmeta.Uint8:
		return gotypes[reflect.Uint8], nil
	case rmeta.Uint16:
		return gotypes[reflect.Uint16], nil
	case rmeta.Uint32, rmeta.Bits:
		return gotypes[reflect.Uint32], nil
	case rmeta.Uint64:
		return gotypes[reflect.Uint64], nil
	case rmeta.Int8:
		return gotypes[reflect.Int8], nil
	case rmeta.Int16:
		return gotypes[reflect.Int16], nil
	case rmeta.Int32:
		return gotypes[reflect.Int32], nil
	case rmeta.Int64, rmeta.Long64:
		return gotypes[reflect.Int64], nil
	case rmeta.Float32:
		return gotypes[reflect.Float32], nil
	case rmeta.Float64:
		return gotypes[reflect.Float64], nil
	case rmeta.Float16:
		return reflect.TypeOf((*root.Float16)(nil)).Elem(), nil
	case rmeta.Double32:
		return reflect.TypeOf((*root.Double32)(nil)).Elem(), nil
	case rmeta.TString, rmeta.STLstring:
		return gotypes[reflect.String], nil

	case rmeta.CharStar:
		return gotypes[reflect.String], nil

	case rmeta.Counter:
		switch size {
		case 4:
			return gotypes[reflect.Int32], nil
		case 8:
			return gotypes[reflect.Int64], nil
		default:
			return nil, fmt.Errorf("rdict: invalid counter size=%d", size)
		}

	case rmeta.TObject:
		return reflect.TypeOf((*rbase.Object)(nil)).Elem(), nil

	case rmeta.TNamed:
		return reflect.TypeOf((*rbase.Named)(nil)).Elem(), nil

	case rmeta.OffsetL + rmeta.Bool:
		return reflect.ArrayOf(n, gotypes[reflect.Bool]), nil
	case rmeta.OffsetL + rmeta.Uint8:
		return reflect.ArrayOf(n, gotypes[reflect.Uint8]), nil
	case rmeta.OffsetL + rmeta.Uint16:
		return reflect.ArrayOf(n, gotypes[reflect.Uint16]), nil
	case rmeta.OffsetL + rmeta.Uint32:
		return reflect.ArrayOf(n, gotypes[reflect.Uint32]), nil
	case rmeta.OffsetL + rmeta.Uint64:
		return reflect.ArrayOf(n, gotypes[reflect.Uint64]), nil
	case rmeta.OffsetL + rmeta.Int8:
		return reflect.ArrayOf(n, gotypes[reflect.Int8]), nil
	case rmeta.OffsetL + rmeta.Int16:
		return reflect.ArrayOf(n, gotypes[reflect.Int16]), nil
	case rmeta.OffsetL + rmeta.Int32:
		return reflect.ArrayOf(n, gotypes[reflect.Int32]), nil
	case rmeta.OffsetL + rmeta.Int64,
		rmeta.OffsetL + rmeta.Long64:
		return reflect.ArrayOf(n, gotypes[reflect.Int64]), nil
	case rmeta.OffsetL + rmeta.Float32:
		return reflect.ArrayOf(n, gotypes[reflect.Float32]), nil
	case rmeta.OffsetL + rmeta.Float64:
		return reflect.ArrayOf(n, gotypes[reflect.Float64]), nil
	case rmeta.OffsetL + rmeta.Float16:
		return reflect.ArrayOf(n, reflect.TypeOf((*root.Float16)(nil)).Elem()), nil
	case rmeta.OffsetL + rmeta.Double32:
		return reflect.ArrayOf(n, reflect.TypeOf((*root.Double32)(nil)).Elem()), nil
	case rmeta.OffsetL + rmeta.TString,
		rmeta.OffsetL + rmeta.CharStar,
		rmeta.OffsetL + rmeta.STLstring:
		return reflect.ArrayOf(n, gotypes[reflect.String]), nil

	case rmeta.OffsetP + rmeta.Bool:
		return reflect.SliceOf(gotypes[reflect.Bool]), nil
	case rmeta.OffsetP + rmeta.Uint8:
		return reflect.SliceOf(gotypes[reflect.Uint8]), nil
	case rmeta.OffsetP + rmeta.Uint16:
		return reflect.SliceOf(gotypes[reflect.Uint16]), nil
	case rmeta.OffsetP + rmeta.Uint32:
		return reflect.SliceOf(gotypes[reflect.Uint32]), nil
	case rmeta.OffsetP + rmeta.Uint64:
		return reflect.SliceOf(gotypes[reflect.Uint64]), nil
	case rmeta.OffsetP + rmeta.Int8:
		return reflect.SliceOf(gotypes[reflect.Int8]), nil
	case rmeta.OffsetP + rmeta.Int16:
		return reflect.SliceOf(gotypes[reflect.Int16]), nil
	case rmeta.OffsetP + rmeta.Int32:
		return reflect.SliceOf(gotypes[reflect.Int32]), nil
	case rmeta.OffsetP + rmeta.Int64:
		return reflect.SliceOf(gotypes[reflect.Int64]), nil
	case rmeta.OffsetP + rmeta.Float32:
		return reflect.SliceOf(gotypes[reflect.Float32]), nil
	case rmeta.OffsetP + rmeta.Float64:
		return reflect.SliceOf(gotypes[reflect.Float64]), nil
	case rmeta.OffsetP + rmeta.Float16:
		return reflect.SliceOf(reflect.TypeOf((*root.Float16)(nil)).Elem()), nil
	case rmeta.OffsetP + rmeta.Double32:
		return reflect.SliceOf(reflect.TypeOf((*root.Double32)(nil)).Elem()), nil
	case rmeta.OffsetP + rmeta.STLstring,
		rmeta.OffsetP + rmeta.CharStar:
		return reflect.SliceOf(gotypes[reflect.String]), nil

	}
	return nil, fmt.Errorf("rmeta=%d (%v) not implemented (size=%d, n=%v)", enum, enum, size, n)
}

func typeFromTypeName(ctx rbytes.StreamerInfoContext, typename string, typevers int16, enum rmeta.Enum, se rbytes.StreamerElement, n int) (reflect.Type, error) {
	e, ok := rmeta.TypeName2Enum(typename)
	if ok {
		return typeFrom(ctx, e, se.Size(), n)
	}

	switch {
	case strings.HasPrefix(typename, "vector<"), strings.HasPrefix(typename, "std::vector<"):
		enames := rmeta.CxxTemplateArgsOf(typename)
		et, err := typeFromTypeName(ctx, enames[0], -1, -1, se, n)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(et), nil

	case strings.HasPrefix(typename, "map<"), strings.HasPrefix(typename, "std::map<"),
		strings.HasPrefix(typename, "unordered_map<"), strings.HasPrefix(typename, "std::unordered_map<"):
		enames := rmeta.CxxTemplateArgsOf(typename)
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
	}

	osi, err := ctx.StreamerInfo(typename, int(typevers))
	if err != nil {
		return nil, fmt.Errorf("rdict: could not find streamer info for %q (version=%d): %w", typename, typevers, err)
	}

	return TypeFromSI(ctx, osi)
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
