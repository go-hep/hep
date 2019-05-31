// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rarrow handles conversion between ROOT and ARROW data models.
package rarrow // import "go-hep.org/x/hep/groot/rarrow"

import (
	"reflect"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rtree"
)

// SchemaFrom returns an Arrow schema from the provided ROOT tree.
func SchemaFrom(t rtree.Tree) *arrow.Schema {
	fields := make([]arrow.Field, len(t.Branches()))
	for i, b := range t.Branches() {
		fields[i] = fieldFromBranch(b)
	}

	return arrow.NewSchema(fields, nil) // FIXME(sbinet): add metadata.
}

func fieldFromBranch(b rtree.Branch) arrow.Field {
	fields := make([]arrow.Field, len(b.Leaves()))
	for i, leaf := range b.Leaves() {
		fields[i] = arrow.Field{
			Name: leaf.Name(),
			Type: dataTypeFromLeaf(leaf),
		}
	}

	if len(fields) == 1 {
		fields[0].Name = b.Name()
		return fields[0]
	}

	return arrow.Field{
		Name: b.Name(),
		Type: arrow.StructOf(fields...),
	}
}

func dataTypeFromLeaf(leaf rtree.Leaf) arrow.DataType {
	var (
		unsigned = leaf.IsUnsigned()
		kind     = leaf.Kind()
		dt       arrow.DataType
	)

	switch kind {
	case reflect.Bool:
		dt = arrow.FixedWidthTypes.Boolean
	case reflect.Int8:
		switch {
		case unsigned:
			dt = arrow.PrimitiveTypes.Uint8
		default:
			dt = arrow.PrimitiveTypes.Int8
		}
	case reflect.Int16:
		switch {
		case unsigned:
			dt = arrow.PrimitiveTypes.Uint16
		default:
			dt = arrow.PrimitiveTypes.Int16
		}
	case reflect.Int32:
		switch {
		case unsigned:
			dt = arrow.PrimitiveTypes.Uint32
		default:
			dt = arrow.PrimitiveTypes.Int32
		}
	case reflect.Int64:
		switch {
		case unsigned:
			dt = arrow.PrimitiveTypes.Uint64
		default:
			dt = arrow.PrimitiveTypes.Int64
		}
	case reflect.Float32:
		dt = arrow.PrimitiveTypes.Float32
	case reflect.Float64:
		dt = arrow.PrimitiveTypes.Float64
	case reflect.String:
		dt = arrow.BinaryTypes.String

	case reflect.Struct:
		dt = dataTypeFromGo(leaf.Type())

	default:
		panic(errors.Errorf("not implemented %#v", leaf))
	}

	switch {
	case leaf.LeafCount() != nil:
		dt = arrow.ListOf(dt)
	case leaf.Len() > 1:
		switch leaf.Kind() {
		case reflect.String:
			switch dims := leaf.ArrayDim(); dims {
			case 0, 1:
				// interpret as a single string
			default:
				// FIXME(sbinet): properly handle [N]string (but ROOT doesn't support that.)
				// see: https://root-forum.cern.ch/t/char-t-in-a-branch/5591/2
				// etype = reflect.ArrayOf(leaf.Len(), etype)
				panic(errors.Errorf("groot/rtree: invalid number of dimensions (%d)", dims))
			}
		default:
			dt = arrow.FixedSizeListOf(int32(leaf.Len()), dt)
		}
	}

	return dt
}

func dataTypeFromGo(typ reflect.Type) arrow.DataType {
	switch typ.Kind() {
	case reflect.Bool:
		return arrow.FixedWidthTypes.Boolean
	case reflect.Int8:
		return arrow.PrimitiveTypes.Int8
	case reflect.Int16:
		return arrow.PrimitiveTypes.Int16
	case reflect.Int32:
		return arrow.PrimitiveTypes.Int32
	case reflect.Int64:
		return arrow.PrimitiveTypes.Int64
	case reflect.Uint8:
		return arrow.PrimitiveTypes.Uint8
	case reflect.Uint16:
		return arrow.PrimitiveTypes.Uint16
	case reflect.Uint32:
		return arrow.PrimitiveTypes.Uint32
	case reflect.Uint64:
		return arrow.PrimitiveTypes.Uint64
	case reflect.Float32:
		return arrow.PrimitiveTypes.Float32
	case reflect.Float64:
		return arrow.PrimitiveTypes.Float64
	case reflect.Slice:
		// special case []byte
		if typ.Elem().Kind() == reflect.Uint8 {
			return arrow.BinaryTypes.Binary
		}
		return arrow.ListOf(dataTypeFromGo(typ.Elem()))
	case reflect.Array:
		return arrow.FixedSizeListOf(int32(typ.Len()), dataTypeFromGo(typ.Elem()))
	case reflect.String:
		return arrow.BinaryTypes.String

	case reflect.Struct:
		fields := make([]arrow.Field, typ.NumField())
		for i := range fields {
			f := typ.Field(i)
			name := f.Name
			if v, ok := f.Tag.Lookup("groot"); ok {
				name = v
			}
			fields[i] = arrow.Field{
				Name: name,
				Type: dataTypeFromGo(f.Type),
			}
		}
		return arrow.StructOf(fields...)

	default:
		panic(errors.Errorf("rarrow: unsupported Go type %v", typ))
	}
}

func builderFrom(mem memory.Allocator, dt arrow.DataType, size int64) array.Builder {
	var bldr array.Builder
	switch dt := dt.(type) {
	case *arrow.BooleanType:
		bldr = array.NewBooleanBuilder(mem)
	case *arrow.Int8Type:
		bldr = array.NewInt8Builder(mem)
	case *arrow.Int16Type:
		bldr = array.NewInt16Builder(mem)
	case *arrow.Int32Type:
		bldr = array.NewInt32Builder(mem)
	case *arrow.Int64Type:
		bldr = array.NewInt64Builder(mem)
	case *arrow.Uint8Type:
		bldr = array.NewUint8Builder(mem)
	case *arrow.Uint16Type:
		bldr = array.NewUint16Builder(mem)
	case *arrow.Uint32Type:
		bldr = array.NewUint32Builder(mem)
	case *arrow.Uint64Type:
		bldr = array.NewUint64Builder(mem)
	case *arrow.Float32Type:
		bldr = array.NewFloat32Builder(mem)
	case *arrow.Float64Type:
		bldr = array.NewFloat64Builder(mem)
	case *arrow.BinaryType:
		bldr = array.NewBinaryBuilder(mem, dt)
	case *arrow.StringType:
		bldr = array.NewStringBuilder(mem)
	case *arrow.ListType:
		bldr = array.NewListBuilder(mem, dt.Elem())
	case *arrow.FixedSizeListType:
		bldr = array.NewFixedSizeListBuilder(mem, dt.Len(), dt.Elem())
	case *arrow.StructType:
		bldr = array.NewStructBuilder(mem, dt)
	default:
		panic(errors.Errorf("groot/rarrow: invalid Arrow type %v", dt))
	}
	bldr.Reserve(int(size))
	return bldr
}

func appendData(bldr array.Builder, v rtree.ScanVar, dt arrow.DataType) {
	switch bldr := bldr.(type) {
	case *array.BooleanBuilder:
		bldr.Append(*v.Value.(*bool))
	case *array.Int8Builder:
		bldr.Append(*v.Value.(*int8))
	case *array.Int16Builder:
		bldr.Append(*v.Value.(*int16))
	case *array.Int32Builder:
		bldr.Append(*v.Value.(*int32))
	case *array.Int64Builder:
		bldr.Append(*v.Value.(*int64))
	case *array.Uint8Builder:
		bldr.Append(*v.Value.(*uint8))
	case *array.Uint16Builder:
		bldr.Append(*v.Value.(*uint16))
	case *array.Uint32Builder:
		bldr.Append(*v.Value.(*uint32))
	case *array.Uint64Builder:
		bldr.Append(*v.Value.(*uint64))
	case *array.Float32Builder:
		bldr.Append(*v.Value.(*float32))
	case *array.Float64Builder:
		bldr.Append(*v.Value.(*float64))
	case *array.StringBuilder:
		bldr.Append(*v.Value.(*string))

	case *array.ListBuilder:
		sub := bldr.ValueBuilder()
		v := reflect.ValueOf(v.Value).Elem()
		sub.Reserve(v.Len())
		bldr.Append(true)
		for i := 0; i < v.Len(); i++ {
			appendValue(sub, v.Index(i).Interface())
		}

	case *array.FixedSizeListBuilder:
		sub := bldr.ValueBuilder()
		v := reflect.ValueOf(v.Value).Elem()
		sub.Reserve(v.Len())
		bldr.Append(true)
		for i := 0; i < v.Len(); i++ {
			appendValue(sub, v.Index(i).Interface())
		}

	case *array.StructBuilder:
		bldr.Append(true)
		v := reflect.ValueOf(v.Value).Elem()
		for i := 0; i < bldr.NumField(); i++ {
			f := bldr.FieldBuilder(i)
			appendValue(f, v.Field(i).Interface())
		}

	default:
		panic(errors.Errorf("groot/rarrow: invalid Arrow builder type %T", bldr))
	}
}

func appendValue(bldr array.Builder, v interface{}) {
	switch b := bldr.(type) {
	case *array.BooleanBuilder:
		b.Append(v.(bool))
	case *array.Int8Builder:
		b.Append(v.(int8))
	case *array.Int16Builder:
		b.Append(v.(int16))
	case *array.Int32Builder:
		b.Append(v.(int32))
	case *array.Int64Builder:
		b.Append(v.(int64))
	case *array.Uint8Builder:
		b.Append(v.(uint8))
	case *array.Uint16Builder:
		b.Append(v.(uint16))
	case *array.Uint32Builder:
		b.Append(v.(uint32))
	case *array.Uint64Builder:
		b.Append(v.(uint64))
	case *array.Float32Builder:
		b.Append(v.(float32))
	case *array.Float64Builder:
		b.Append(v.(float64))
	case *array.StringBuilder:
		b.Append(v.(string))

	case *array.ListBuilder:
		b.Append(true)
		sub := b.ValueBuilder()
		v := reflect.ValueOf(v)
		for i := 0; i < v.Len(); i++ {
			appendValue(sub, v.Index(i).Interface())
		}

	case *array.FixedSizeListBuilder:
		b.Append(true)
		sub := b.ValueBuilder()
		v := reflect.ValueOf(v)
		for i := 0; i < v.Len(); i++ {
			appendValue(sub, v.Index(i).Interface())
		}

	case *array.StructBuilder:
		v := reflect.ValueOf(v)
		for i := 0; i < b.NumField(); i++ {
			f := b.FieldBuilder(i)
			appendValue(f, v.Field(i).Interface())
		}

	default:
		panic(errors.Errorf("groot/rarrow: invalid Arrow builder type %T", b))
	}
}
