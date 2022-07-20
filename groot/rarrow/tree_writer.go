// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow

import (
	"fmt"
	"reflect"

	"git.sr.ht/~sbinet/go-arrow"
	"git.sr.ht/~sbinet/go-arrow/array"
	"git.sr.ht/~sbinet/go-arrow/arrio"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

// flatTreeWriter writes ARROW data as a ROOT flat-tree.
type flatTreeWriter struct {
	w      rtree.Writer
	schema *arrow.Schema
	ctx    contextWriter
}

// NewFlatTreeWriter creates an arrio.Writer that writes ARROW data as a ROOT
// flat-tree under the provided dir directory.
func NewFlatTreeWriter(dir riofs.Directory, name string, schema *arrow.Schema, opts ...rtree.WriteOption) (*flatTreeWriter, error) {
	var (
		ctx   = newContextWriter(schema)
		wvars = make([]rtree.WriteVar, 0, len(ctx.wvars)+len(ctx.count))
	)

	for _, wvar := range ctx.count {
		wvars = append(wvars, wvar)
	}
	wvars = append(wvars, ctx.wvars...)

	tree, err := rtree.NewWriter(dir, name, wvars, opts...)
	if err != nil {
		return nil, fmt.Errorf("rarrow: could not create flat-tree writer %q: %w", name, err)
	}
	return &flatTreeWriter{w: tree, schema: schema, ctx: ctx}, nil
}

// Close closes the underlying ROOT tree writer.
func (fw *flatTreeWriter) Close() error {
	return fw.w.Close()
}

// Write writes the provided ARROW record to the underlying ROOT flat-tree.
// Write implements arrio.Writer.
func (fw *flatTreeWriter) Write(rec array.Record) error {
	if src := rec.Schema(); !fw.schema.Equal(src) {
		return fmt.Errorf("rarrow: invalid input record schema:\n - got= %v\n - want=%v", src, fw.schema)
	}

	nrows := rec.Column(0).Len()
	for icol, col := range rec.Columns() {
		if col.Len() != nrows {
			return fmt.Errorf(
				"rarrow: column %q (index=%d) has not the same number of rows than others (got=%d, want=%d)",
				rec.ColumnName(icol), icol, col.Len(), nrows,
			)
		}
	}

	for irow := 0; irow < nrows; irow++ {
		for icol, col := range rec.Columns() {
			wvar := &fw.ctx.wvars[icol]
			err := fw.ctx.readFrom(wvar, irow, col)
			if err != nil {
				return fmt.Errorf(
					"rarrow: could not read row=%d from column[%d](name=%s): %w",
					irow, icol, rec.ColumnName(icol), err,
				)
			}
		}
		_, err := fw.w.Write()
		if err != nil {
			return fmt.Errorf("rarrow: could not write row=%d to tree: %w", irow, err)
		}
	}

	return nil
}

type contextWriter struct {
	wvars []rtree.WriteVar
	count map[string]rtree.WriteVar
}

func newContextWriter(schema *arrow.Schema) contextWriter {
	ctx := contextWriter{
		wvars: make([]rtree.WriteVar, len(schema.Fields())),
		count: make(map[string]rtree.WriteVar),
	}
	for i, field := range schema.Fields() {
		ctx.wvars[i] = ctx.writeVarFrom(field)
	}
	return ctx
}

func (ctx *contextWriter) writeVarFrom(field arrow.Field) rtree.WriteVar {
	switch dt := field.Type.(type) {
	case *arrow.BooleanType:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(bool),
		}

	case *arrow.Int8Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(int8),
		}

	case *arrow.Int16Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(int16),
		}

	case *arrow.Int32Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(int32),
		}

	case *arrow.Int64Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(int64),
		}

	case *arrow.Uint8Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(uint8),
		}

	case *arrow.Uint16Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(uint16),
		}

	case *arrow.Uint32Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(uint32),
		}

	case *arrow.Uint64Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(uint64),
		}

	case *arrow.Float32Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(float32),
		}

	case *arrow.Float64Type:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(float64),
		}

	case *arrow.StringType:
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(string),
		}
	case *arrow.BinaryType:
		// FIXME(sbinet): differentiate the 2 (Binary/String) ?
		return rtree.WriteVar{
			Name:  field.Name,
			Value: new(string),
		}

	case *arrow.FixedSizeListType:
		wv := ctx.writeVarFrom(arrow.Field{Type: dt.Elem(), Name: "elem"})
		rt := reflect.ArrayOf(int(dt.Len()), reflect.TypeOf(wv.Value).Elem())
		return rtree.WriteVar{
			Name:  field.Name,
			Value: reflect.New(rt).Interface(),
		}

	case *arrow.FixedSizeBinaryType:
		rt := reflect.ArrayOf(dt.ByteWidth, reflect.TypeOf(byte(0)))
		return rtree.WriteVar{
			Name:  field.Name,
			Value: reflect.New(rt).Interface(),
		}

	case *arrow.ListType:
		wv := ctx.writeVarFrom(arrow.Field{Type: dt.Elem(), Name: "elem"})
		rt := reflect.SliceOf(reflect.TypeOf(wv.Value).Elem())
		nn := "rarrow_n_" + field.Name
		ctx.count[field.Name] = rtree.WriteVar{
			Name:  nn,
			Value: new(int32),
		}
		return rtree.WriteVar{
			Name:  field.Name,
			Value: reflect.New(rt).Interface(),
			Count: nn,
		}

		//	case *arrow.StructType:
		//		fields := make([]reflect.StructField, len(dt.Fields()))
		//		for i, ft := range dt.Fields() {
		//			wv := writeVarFrom(ft)
		//			fields[i] = reflect.StructField{
		//				Name: "ROOT_" + ft.Name,
		//				Type: reflect.TypeOf(wv.Value).Elem(),
		//				Tag:  reflect.StructTag(fmt.Sprintf("groot:%q", ft.Name)),
		//			}
		//		}
		//		rt := reflect.StructOf(fields)
		//		return rtree.WriteVar{
		//			Name:  field.Name,
		//			Value: reflect.New(rt).Interface(),
		//		}

	default:
		panic(fmt.Errorf("invalid ARROW data-type: %T", dt))
	}
}

func (ctx *contextWriter) readFrom(wvar *rtree.WriteVar, irow int, arr array.Interface) error {
	ptr := wvar.Value
	switch arr := arr.(type) {
	case *array.Boolean:
		*ptr.(*bool) = arr.Value(irow)
	case *array.Int8:
		*ptr.(*int8) = arr.Value(irow)
	case *array.Int16:
		*ptr.(*int16) = arr.Value(irow)
	case *array.Int32:
		*ptr.(*int32) = arr.Value(irow)
	case *array.Int64:
		*ptr.(*int64) = arr.Value(irow)
	case *array.Uint8:
		*ptr.(*uint8) = arr.Value(irow)
	case *array.Uint16:
		*ptr.(*uint16) = arr.Value(irow)
	case *array.Uint32:
		*ptr.(*uint32) = arr.Value(irow)
	case *array.Uint64:
		*ptr.(*uint64) = arr.Value(irow)
	case *array.Float32:
		*ptr.(*float32) = arr.Value(irow)
	case *array.Float64:
		*ptr.(*float64) = arr.Value(irow)
	case *array.String:
		*ptr.(*string) = arr.Value(irow)
	case *array.Binary:
		*ptr.(*string) = string(arr.Value(irow))

	case *array.FixedSizeList:
		rv := reflect.ValueOf(ptr).Elem()
		n := int64(rv.Len())
		off := int64(arr.Offset())
		beg := (off + int64(irow)) * n
		end := (off + int64(irow+1)) * n
		ra := array.NewSlice(arr.ListValues(), beg, end)
		defer ra.Release()
		ptr := &rtree.WriteVar{
			Name: "_rarrow_elem_" + wvar.Name,
		}
		for i := 0; i < rv.Len(); i++ {
			ptr.Value = rv.Index(i).Addr().Interface()
			err := ctx.readFrom(ptr, i, ra)
			if err != nil {
				return err
			}
		}

	case *array.FixedSizeBinary:
		rv := reflect.ValueOf(ptr).Elem()
		sli := rv.Slice(0, rv.Len()).Interface().([]byte)
		copy(sli, arr.Value(irow))

	case *array.List:
		rv := reflect.ValueOf(ptr).Elem()
		rc := reflect.ValueOf(ctx.count[wvar.Name].Value).Elem()
		if !arr.IsValid(irow) {
			rc.SetInt(0)
			rv.SetLen(0)
			return nil
		}

		j := irow + arr.Data().Offset()
		beg := int64(arr.Offsets()[j])
		end := int64(arr.Offsets()[j+1])
		sli := array.NewSlice(arr.ListValues(), beg, end)
		defer sli.Release()

		sz := sli.Len()
		rc.SetInt(int64(sz))

		if src, dst := sz, rv.Len(); src > dst {
			rv.Set(reflect.MakeSlice(rv.Type(), src, src))
		}
		rv.SetLen(sz)

		ptr := &rtree.WriteVar{
			Name: "_rarrow_elem_" + wvar.Name,
		}
		for i := 0; i < sli.Len(); i++ {
			ptr.Value = rv.Index(i).Addr().Interface()
			err := ctx.readFrom(ptr, i, sli)
			if err != nil {
				return err
			}
		}

	default:
		panic(fmt.Errorf("invalid array type %T", arr))
	}
	return nil
}

var (
	_ arrio.Writer = (*flatTreeWriter)(nil)
)
