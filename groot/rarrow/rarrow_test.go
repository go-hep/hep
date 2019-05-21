// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow // import "go-hep.org/x/hep/groot/rarrow"

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/apache/arrow/go/arrow"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func mdFrom(vs ...string) arrow.Metadata {
	keys := make([]string, 0, len(vs)/2)
	vals := make([]string, 0, len(vs)/2)
	for i, v := range vs {
		switch {
		case i%2 == 0:
			keys = append(keys, v)
		default:
			vals = append(vals, v)
		}
	}
	return arrow.NewMetadata(keys, vals)
}

func TestSchemaFrom(t *testing.T) {
	for _, tc := range []struct {
		file string
		tree string
		want *arrow.Schema
	}{
		{
			file: "../testdata/simple.root",
			tree: "tree",
			want: arrow.NewSchema([]arrow.Field{
				{Name: "one", Type: arrow.PrimitiveTypes.Int32},
				{Name: "two", Type: arrow.PrimitiveTypes.Float32},
				{Name: "three", Type: arrow.BinaryTypes.String},
			}, nil),
		},
		{
			file: "../testdata/small-flat-tree.root",
			tree: "tree",
			want: arrow.NewSchema([]arrow.Field{
				{Name: "Int32", Type: arrow.PrimitiveTypes.Int32},
				{Name: "Int64", Type: arrow.PrimitiveTypes.Int64},
				{Name: "UInt32", Type: arrow.PrimitiveTypes.Uint32},
				{Name: "UInt64", Type: arrow.PrimitiveTypes.Uint64},
				{Name: "Float32", Type: arrow.PrimitiveTypes.Float32},
				{Name: "Float64", Type: arrow.PrimitiveTypes.Float64},
				{Name: "Str", Type: arrow.BinaryTypes.String},
				{Name: "ArrayInt32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int32)},
				{Name: "ArrayInt64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int64)},
				{Name: "ArrayUInt32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint32)},
				{Name: "ArrayUInt64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint64)},
				{Name: "ArrayFloat32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Float32)},
				{Name: "ArrayFloat64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Float64)},
				{Name: "N", Type: arrow.PrimitiveTypes.Int32},
				{Name: "SliceInt32", Type: arrow.ListOf(arrow.PrimitiveTypes.Int32) /*, Metadata: mdFrom("ROOT:LeafCount", "N")*/},
				{Name: "SliceInt64", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64) /*, Metadata: mdFrom("ROOT:LeafCount", "N")*/},
				{Name: "SliceUInt32", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint32) /*, Metadata: mdFrom("ROOT:LeafCount", "N")*/},
				{Name: "SliceUInt64", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint64) /*, Metadata: mdFrom("ROOT:LeafCount", "N")*/},
				{Name: "SliceFloat32", Type: arrow.ListOf(arrow.PrimitiveTypes.Float32) /*, Metadata: mdFrom("ROOT:LeafCount", "N")*/},
				{Name: "SliceFloat64", Type: arrow.ListOf(arrow.PrimitiveTypes.Float64) /*, Metadata: mdFrom("ROOT:LeafCount", "N")*/},
			}, nil),
		},
		{
			file: "../testdata/small-evnt-tree-fullsplit.root",
			tree: "tree",
			want: arrow.NewSchema([]arrow.Field{
				{Name: "evt", Type: arrow.StructOf([]arrow.Field{
					{Name: "Beg", Type: arrow.BinaryTypes.String},
					{Name: "I16", Type: arrow.PrimitiveTypes.Int16},
					{Name: "I32", Type: arrow.PrimitiveTypes.Int32},
					{Name: "I64", Type: arrow.PrimitiveTypes.Int64},
					{Name: "U16", Type: arrow.PrimitiveTypes.Uint16},
					{Name: "U32", Type: arrow.PrimitiveTypes.Uint32},
					{Name: "U64", Type: arrow.PrimitiveTypes.Uint64},
					{Name: "F32", Type: arrow.PrimitiveTypes.Float32},
					{Name: "F64", Type: arrow.PrimitiveTypes.Float64},
					{Name: "Str", Type: arrow.BinaryTypes.String},
					{Name: "P3", Type: arrow.StructOf([]arrow.Field{
						{Name: "Px", Type: arrow.PrimitiveTypes.Int32},
						{Name: "Py", Type: arrow.PrimitiveTypes.Float64},
						{Name: "Pz", Type: arrow.PrimitiveTypes.Int32},
					}...)},
					{Name: "ArrayI16", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int16)},
					{Name: "ArrayI32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int32)},
					{Name: "ArrayI64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int64)},
					{Name: "ArrayU16", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint16)},
					{Name: "ArrayU32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint32)},
					{Name: "ArrayU64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint64)},
					{Name: "ArrayF32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Float32)},
					{Name: "ArrayF64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Float64)},
					{Name: "N", Type: arrow.PrimitiveTypes.Int32},
					{Name: "SliceI16", Type: arrow.ListOf(arrow.PrimitiveTypes.Int16)},
					{Name: "SliceI32", Type: arrow.ListOf(arrow.PrimitiveTypes.Int32)},
					{Name: "SliceI64", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)},
					{Name: "SliceU16", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint16)},
					{Name: "SliceU32", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint32)},
					{Name: "SliceU64", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint64)},
					{Name: "SliceF32", Type: arrow.ListOf(arrow.PrimitiveTypes.Float32)},
					{Name: "SliceF64", Type: arrow.ListOf(arrow.PrimitiveTypes.Float64)},
					{Name: "StdStr", Type: arrow.BinaryTypes.String},
					{Name: "StlVecI16", Type: arrow.ListOf(arrow.PrimitiveTypes.Int16)},
					{Name: "StlVecI32", Type: arrow.ListOf(arrow.PrimitiveTypes.Int32)},
					{Name: "StlVecI64", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)},
					{Name: "StlVecU16", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint16)},
					{Name: "StlVecU32", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint32)},
					{Name: "StlVecU64", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint64)},
					{Name: "StlVecF32", Type: arrow.ListOf(arrow.PrimitiveTypes.Float32)},
					{Name: "StlVecF64", Type: arrow.ListOf(arrow.PrimitiveTypes.Float64)},
					{Name: "StlVecStr", Type: arrow.ListOf(arrow.BinaryTypes.String)},
					{Name: "End", Type: arrow.BinaryTypes.String},
				}...)},
			}, nil),
		},
		{
			file: "../testdata/small-evnt-tree-nosplit.root",
			tree: "tree",
			want: arrow.NewSchema([]arrow.Field{
				{Name: "evt", Type: arrow.StructOf([]arrow.Field{
					{Name: "Beg", Type: arrow.BinaryTypes.String},
					{Name: "I16", Type: arrow.PrimitiveTypes.Int16},
					{Name: "I32", Type: arrow.PrimitiveTypes.Int32},
					{Name: "I64", Type: arrow.PrimitiveTypes.Int64},
					{Name: "U16", Type: arrow.PrimitiveTypes.Uint16},
					{Name: "U32", Type: arrow.PrimitiveTypes.Uint32},
					{Name: "U64", Type: arrow.PrimitiveTypes.Uint64},
					{Name: "F32", Type: arrow.PrimitiveTypes.Float32},
					{Name: "F64", Type: arrow.PrimitiveTypes.Float64},
					{Name: "Str", Type: arrow.BinaryTypes.String},
					{Name: "P3", Type: arrow.StructOf([]arrow.Field{
						{Name: "Px", Type: arrow.PrimitiveTypes.Int32},
						{Name: "Py", Type: arrow.PrimitiveTypes.Float64},
						{Name: "Pz", Type: arrow.PrimitiveTypes.Int32},
					}...)},
					{Name: "ArrayI16", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int16)},
					{Name: "ArrayI32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int32)},
					{Name: "ArrayI64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Int64)},
					{Name: "ArrayU16", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint16)},
					{Name: "ArrayU32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint32)},
					{Name: "ArrayU64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Uint64)},
					{Name: "ArrayF32", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Float32)},
					{Name: "ArrayF64", Type: fixedSizeListOf(10, arrow.PrimitiveTypes.Float64)},
					{Name: "N", Type: arrow.PrimitiveTypes.Int32},
					{Name: "SliceI16", Type: arrow.ListOf(arrow.PrimitiveTypes.Int16)},
					{Name: "SliceI32", Type: arrow.ListOf(arrow.PrimitiveTypes.Int32)},
					{Name: "SliceI64", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)},
					{Name: "SliceU16", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint16)},
					{Name: "SliceU32", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint32)},
					{Name: "SliceU64", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint64)},
					{Name: "SliceF32", Type: arrow.ListOf(arrow.PrimitiveTypes.Float32)},
					{Name: "SliceF64", Type: arrow.ListOf(arrow.PrimitiveTypes.Float64)},
					{Name: "StdStr", Type: arrow.BinaryTypes.String},
					{Name: "StlVecI16", Type: arrow.ListOf(arrow.PrimitiveTypes.Int16)},
					{Name: "StlVecI32", Type: arrow.ListOf(arrow.PrimitiveTypes.Int32)},
					{Name: "StlVecI64", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)},
					{Name: "StlVecU16", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint16)},
					{Name: "StlVecU32", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint32)},
					{Name: "StlVecU64", Type: arrow.ListOf(arrow.PrimitiveTypes.Uint64)},
					{Name: "StlVecF32", Type: arrow.ListOf(arrow.PrimitiveTypes.Float32)},
					{Name: "StlVecF64", Type: arrow.ListOf(arrow.PrimitiveTypes.Float64)},
					{Name: "StlVecStr", Type: arrow.ListOf(arrow.BinaryTypes.String)},
					{Name: "End", Type: arrow.BinaryTypes.String},
				}...)},
			}, nil),
		},
		{
			file: "../testdata/root_numpy_struct.root",
			tree: "test",
			want: arrow.NewSchema([]arrow.Field{
				{Name: "branch1", Type: arrow.StructOf([]arrow.Field{
					{Name: "intleaf", Type: arrow.PrimitiveTypes.Int32},
					{Name: "floatleaf", Type: arrow.PrimitiveTypes.Float32},
				}...)},
				{Name: "branch2", Type: arrow.StructOf([]arrow.Field{
					{Name: "intleaf", Type: arrow.PrimitiveTypes.Int32},
					{Name: "floatleaf", Type: arrow.PrimitiveTypes.Float32},
				}...)},
			}, nil),
		},
	} {
		t.Run(tc.file, func(t *testing.T) {
			f, err := groot.Open(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			o, err := riofs.Dir(f).Get(tc.tree)
			if err != nil {
				t.Fatal(err)
			}

			tree := o.(rtree.Tree)

			got := SchemaFrom(tree)

			if !got.Equal(tc.want) {
				t.Fatalf("invalid schema.\ngot:\n%s\nwant:\n%s\n", displaySchema(got), displaySchema(tc.want))
			}
		})
	}
}

func displaySchema(s *arrow.Schema) string {
	o := new(strings.Builder)
	fmt.Fprintf(o, "%*.sfields: %d\n", 2, "", len(s.Fields()))
	for _, f := range s.Fields() {
		displayField(o, f, 4)
	}
	if meta := s.Metadata(); meta.Len() > 0 {
		fmt.Fprintf(o, "metadata: %v\n", meta)
	}
	return o.String()
}

func displayField(o io.Writer, field arrow.Field, inc int) {
	nullable := ""
	if field.Nullable {
		nullable = ", nullable"
	}
	fmt.Fprintf(o, "%*.s- %s: type=%v%v\n", inc, "", field.Name, field.Type, nullable)
	if field.HasMetadata() {
		fmt.Fprintf(o, "%*.smetadata: %v\n", inc, "", field.Metadata)
	}
}
