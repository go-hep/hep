// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// arrow2root converts the content of an ARROW file to a ROOT TTree.
package main // import "go-hep.org/x/hep/cmd/arrow2root"

import (
	"flag"
	"log"
	"os"
	"reflect"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"golang.org/x/xerrors"
)

func main() {
	log.SetPrefix("arrow2root: ")
	log.SetFlags(0)

	oname := flag.String("o", "output.root", "path to output ROOT file name")
	tname := flag.String("t", "tree", "name of the output tree")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing input ARROW filename argument")
	}
	fname := flag.Arg(0)

	err := process(*oname, *tname, fname)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func process(oname, tname, fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return xerrors.Errorf("could not open ARROW file %q: %w", fname, err)
	}
	defer f.Close()

	mem := memory.NewGoAllocator()
	r, err := ipc.NewFileReader(f, ipc.WithAllocator(mem))
	if err != nil {
		return xerrors.Errorf("could not create ARROW IPC reader from %q: %w", fname, err)
	}
	defer r.Close()

	o, err := groot.Create(oname)
	if err != nil {
		return xerrors.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer o.Close()

	schema := r.Schema()
	wvars := writeVarsFrom(schema)
	tree, err := rtree.NewWriter(o, tname, wvars, rtree.WithTitle(tname))
	if err != nil {
		return xerrors.Errorf("could not create output ROOT tree %q: %w", tname, err)
	}

	err = convert(tree, wvars, r)
	if err != nil {
		return xerrors.Errorf("could not convert ARROW file to ROOT tree: %w", err)
	}

	err = tree.Close()
	if err != nil {
		return xerrors.Errorf("could not close ROOT tree writer: %w", err)
	}

	err = o.Close()
	if err != nil {
		return xerrors.Errorf("could not close output ROOT file %q: %w", oname, err)
	}

	return nil
}

func writeVarsFrom(schema *arrow.Schema) []rtree.WriteVar {
	wvars := make([]rtree.WriteVar, len(schema.Fields()))
	for i, field := range schema.Fields() {
		wvars[i] = writeVarFrom(field)
	}
	return wvars
}

func writeVarFrom(field arrow.Field) rtree.WriteVar {
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
		wv := writeVarFrom(arrow.Field{Type: dt.Elem(), Name: "elem"})
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

		//	case *arrow.ListType:
		//		wv := writeVarFrom(arrow.Field{Type: dt.Elem(), Name: "elem"})
		//		rt := reflect.SliceOf(reflect.TypeOf(wv.Value).Elem())
		//		return rtree.WriteVar{
		//			Name:  field.Name,
		//			Value: reflect.New(rt).Interface(),
		//		}
		//
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
		panic(xerrors.Errorf("invalid ARROW data-type: %T", dt))
	}
}

func convert(tree rtree.Writer, wvars []rtree.WriteVar, r *ipc.FileReader) error {

	for irec := 0; irec < r.NumRecords(); irec++ {
		err := convertRecord(tree, wvars, irec, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertRecord(tree rtree.Writer, wvars []rtree.WriteVar, irec int, r *ipc.FileReader) error {
	rec, err := r.Record(irec)
	if err != nil {
		return xerrors.Errorf("could not read record %d: %w", irec, err)
	}
	defer rec.Release()

	ncols := len(rec.Columns())
	if ncols != len(wvars) {
		return xerrors.Errorf("record %d has not the same number of columns than reference (got=%d, want=%d)",
			irec, ncols, len(wvars),
		)
	}

	nrows := rec.Column(0).Len()
	for icol, col := range rec.Columns() {
		if col.Len() != nrows {
			return xerrors.Errorf(
				"column %q (index=%d) from record %d has not the same number of rows than others (got=%d, want=%d)",
				rec.ColumnName(icol), icol, irec, col.Len(), nrows,
			)
		}
	}

	for irow := 0; irow < nrows; irow++ {
		for icol, col := range rec.Columns() {
			err = readFrom(wvars[icol].Value, irow, col)
			if err != nil {
				return xerrors.Errorf(
					"record[%d]: could not read row=%d from column[%d](name=%s): %w",
					irec, irow, icol, rec.ColumnName(icol), err,
				)
			}
		}
		_, err = tree.Write()
		if err != nil {
			return xerrors.Errorf("could not write row=%d to tree: %w", irow, err)
		}
	}

	return nil
}

func readFrom(ptr interface{}, irow int, arr array.Interface) error {
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
		for i := 0; i < rv.Len(); i++ {
			err := readFrom(rv.Index(i).Addr().Interface(), i, ra)
			if err != nil {
				return err
			}
		}

	case *array.FixedSizeBinary:
		rv := reflect.ValueOf(ptr).Elem()
		sli := rv.Slice(0, rv.Len()).Interface().([]byte)
		copy(sli, arr.Value(irow))

		//	case *array.List:
		//		panic("slice")
	default:
		panic(xerrors.Errorf("invalid array type %T", arr))
	}
	return nil
}
