// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"reflect"
	"sync/atomic"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/arrio"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/sbinet/npyio"
	"golang.org/x/xerrors"
)

var (
	boolType    = reflect.TypeOf(true)
	uint8Type   = reflect.TypeOf((*uint8)(nil)).Elem()
	uint16Type  = reflect.TypeOf((*uint16)(nil)).Elem()
	uint32Type  = reflect.TypeOf((*uint32)(nil)).Elem()
	uint64Type  = reflect.TypeOf((*uint64)(nil)).Elem()
	int8Type    = reflect.TypeOf((*int8)(nil)).Elem()
	int16Type   = reflect.TypeOf((*int16)(nil)).Elem()
	int32Type   = reflect.TypeOf((*int32)(nil)).Elem()
	int64Type   = reflect.TypeOf((*int64)(nil)).Elem()
	float32Type = reflect.TypeOf((*float32)(nil)).Elem()
	float64Type = reflect.TypeOf((*float64)(nil)).Elem()

//	complex64Type  = reflect.TypeOf((*complex64)(nil)).Elem()
//	complex128Type = reflect.TypeOf((*complex128)(nil)).Elem()
)

// Record is an in-memory Arrow Record backed by a NumPy data file.
type Record struct {
	refs int64

	mem memory.Allocator

	schema *arrow.Schema
	nrows  int64
	ncols  int64

	cols []array.Interface
}

func NewRecord(npy *npyio.Reader) *Record {
	var (
		mem    = memory.NewGoAllocator()
		schema = schemaFrom(npy)
		shape  = make([]int, len(npy.Header.Descr.Shape))
	)

	copy(shape, npy.Header.Descr.Shape)
	if npy.Header.Descr.Fortran {
		a := shape
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		shape = a
	}
	nrows := int64(shape[0])

	rec := &Record{
		refs:   1,
		mem:    mem,
		schema: schema,
		nrows:  nrows,
		ncols:  1,
	}

	nelem := int64(1)
	for _, v := range shape {
		nelem *= int64(v)
	}

	bldr := builderFrom(mem, schema.Field(0).Type, nrows)
	defer bldr.Release()

	rec.read(npy, nelem, bldr)

	return rec
}

// Retain increases the reference count by 1.
// Retain may be called simultaneously from multiple goroutines.
func (rec *Record) Retain() {
	atomic.AddInt64(&rec.refs, 1)
}

// Release decreases the reference count by 1.
// When the reference count goes to zero, the memory is freed.
// Release may be called simultaneously from multiple goroutines.
func (rec *Record) Release() {
	if atomic.LoadInt64(&rec.refs) <= 0 {
		panic("groot/rarrow: too many releases")
	}

	if atomic.AddInt64(&rec.refs, -1) == 0 {
		for i := range rec.cols {
			rec.cols[i].Release()
		}
		rec.cols = nil
	}
}

func (rec *Record) Schema() *arrow.Schema        { return rec.schema }
func (rec *Record) NumRows() int64               { return rec.nrows }
func (rec *Record) NumCols() int64               { return rec.ncols }
func (rec *Record) Columns() []array.Interface   { return rec.cols }
func (rec *Record) Column(i int) array.Interface { return rec.cols[i] }
func (rec *Record) ColumnName(i int) string      { return rec.schema.Field(i).Name }

// NewSlice constructs a zero-copy slice of the record with the indicated
// indices i and j, corresponding to array[i:j].
// The returned record must be Release()'d after use.
//
// NewSlice panics if the slice is outside the valid range of the record array.
// NewSlice panics if j < i.
func (rec *Record) NewSlice(i, j int64) array.Record {
	panic("not implemented")
}

func (rec *Record) read(r *npyio.Reader, nelem int64, bldr array.Builder) {
	rt := dtypeFrom(rec.schema.Field(0).Type)
	rv := reflect.New(reflect.SliceOf(rt)).Elem()
	rv.Set(reflect.MakeSlice(rv.Type(), int(nelem), int(nelem)))

	err := r.Read(rv.Addr().Interface())
	if err != nil {
		panic(xerrors.Errorf("npy2root: could not read numpy data: %w", err))
	}

	offset := nelem / rec.nrows
	for i := int64(0); i < rec.nrows; i++ {
		appendData(bldr, rv.Interface(), int(offset*i), int(offset*(i+1)), rec.schema.Field(0).Type)
	}

	rec.cols = append(rec.cols, bldr.NewArray())
}

func schemaFrom(npy *npyio.Reader) *arrow.Schema {
	var (
		hdr   = npy.Header
		dtype arrow.DataType
	)
	switch hdr.Descr.Type {
	case "b1", "<b1", "|b1", "bool":
		dtype = arrow.FixedWidthTypes.Boolean

	case "u1", "<u1", "|u1", "uint8":
		dtype = arrow.PrimitiveTypes.Uint8

	case "u2", "<u2", "|u2", ">u2", "uint16":
		dtype = arrow.PrimitiveTypes.Uint16

	case "u4", "<u4", "|u4", ">u4", "uint32":
		dtype = arrow.PrimitiveTypes.Uint32

	case "u8", "<u8", "|u8", ">u8", "uint64":
		dtype = arrow.PrimitiveTypes.Uint64

	case "i1", "<i1", "|i1", ">i1", "int8":
		dtype = arrow.PrimitiveTypes.Int8

	case "i2", "<i2", "|i2", ">i2", "int16":
		dtype = arrow.PrimitiveTypes.Int16

	case "i4", "<i4", "|i4", ">i4", "int32":
		dtype = arrow.PrimitiveTypes.Int32

	case "i8", "<i8", "|i8", ">i8", "int64":
		dtype = arrow.PrimitiveTypes.Int64

	case "f4", "<f4", "|f4", ">f4", "float32":
		dtype = arrow.PrimitiveTypes.Float32

	case "f8", "<f8", "|f8", ">f8", "float64":
		dtype = arrow.PrimitiveTypes.Float64

		//	case "c8", "<c8", "|c8", ">c8", "complex64":
		//		panic(xerrors.Errorf("npy2root: complex64 not supported"))
		//
		//	case "c16", "<c16", "|c16", ">c16", "complex128":
		//		panic(xerrors.Errorf("npy2root: complex128 not supported"))

	default:
		panic(xerrors.Errorf("npy2root: invalid dtype descriptor %q", hdr.Descr.Type))
	}

	shape := make([]int, len(hdr.Descr.Shape))
	copy(shape, hdr.Descr.Shape)
	if hdr.Descr.Fortran {
		a := shape
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		shape = a
	}

	switch len(shape) {
	case 1:
		// scalar

	case 2:
		// 1d-array
		dtype = arrow.FixedSizeListOf(int32(shape[1]), dtype)

		//	case 3:
		//		panic(xerrors.Errorf("npy2root: 2d-arrays not yet supported by groot"))
		//		// 		// 2d-array
		//		// 		dtype = arrow.FixedSizeListOf(int32(shape[1]), dtype)
		//		// 		dtype = arrow.FixedSizeListOf(int32(shape[2]), dtype)

	default:
		panic(xerrors.Errorf("npy2root: invalid shape descriptor %v", hdr.Descr.Shape))
	}

	schema := arrow.NewSchema([]arrow.Field{{Name: "numpy", Type: dtype}}, nil)
	return schema
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
		//	case *arrow.BinaryType:
		//		bldr = array.NewBinaryBuilder(mem, dt)
		//	case *arrow.StringType:
		//		bldr = array.NewStringBuilder(mem)
	case *arrow.FixedSizeListType:
		bldr = array.NewFixedSizeListBuilder(mem, dt.Len(), dt.Elem())
	default:
		panic(xerrors.Errorf("npy2root: invalid Arrow type %v", dt))
	}
	bldr.Reserve(int(size))
	return bldr
}

func dtypeFrom(dt arrow.DataType) reflect.Type {
	switch dt := dt.(type) {
	case *arrow.BooleanType:
		return boolType
	case *arrow.Int8Type:
		return int8Type
	case *arrow.Int16Type:
		return int16Type
	case *arrow.Int32Type:
		return int32Type
	case *arrow.Int64Type:
		return int64Type
	case *arrow.Uint8Type:
		return uint8Type
	case *arrow.Uint16Type:
		return uint16Type
	case *arrow.Uint32Type:
		return uint32Type
	case *arrow.Uint64Type:
		return uint64Type
	case *arrow.Float32Type:
		return float32Type
	case *arrow.Float64Type:
		return float64Type
		//	case *arrow.BinaryType:
		//		bldr = array.NewBinaryBuilder(mem, dt)
		//	case *arrow.StringType:
		//		bldr = array.NewStringBuilder(mem)
	case *arrow.FixedSizeListType:
		return dtypeFrom(dt.Elem())
	default:
		panic(xerrors.Errorf("npy2root: invalid Arrow type %v", dt))
	}
}

func appendData(bldr array.Builder, vs interface{}, beg, end int, dt arrow.DataType) {
	switch bldr := bldr.(type) {
	case *array.BooleanBuilder:
		bldr.Append(vs.([]bool)[beg])
	case *array.Int8Builder:
		bldr.Append(vs.([]int8)[beg])
	case *array.Int16Builder:
		bldr.Append(vs.([]int16)[beg])
	case *array.Int32Builder:
		bldr.Append(vs.([]int32)[beg])
	case *array.Int64Builder:
		bldr.Append(vs.([]int64)[beg])
	case *array.Uint8Builder:
		bldr.Append(vs.([]uint8)[beg])
	case *array.Uint16Builder:
		bldr.Append(vs.([]uint16)[beg])
	case *array.Uint32Builder:
		bldr.Append(vs.([]uint32)[beg])
	case *array.Uint64Builder:
		bldr.Append(vs.([]uint64)[beg])
	case *array.Float32Builder:
		bldr.Append(vs.([]float32)[beg])
	case *array.Float64Builder:
		bldr.Append(vs.([]float64)[beg])
	case *array.FixedSizeListBuilder:
		elemType := dt.(*arrow.FixedSizeListType)
		sub := bldr.ValueBuilder()
		n := int(elemType.Len())
		sub.Reserve(n)
		bldr.Append(true)
		for i := 0; i < n; i++ {
			appendData(sub, vs, beg+i, end, elemType)
		}
	default:
		panic(xerrors.Errorf("npy2root: invalid Arrow builder type %T", bldr))
	}
}

type RecordReader struct {
	recs []array.Record
	cur  int
}

func NewRecordReader(recs ...array.Record) *RecordReader {
	return &RecordReader{
		recs: recs,
		cur:  0,
	}
}

func (rr *RecordReader) Read() (array.Record, error) {
	if rr.cur >= len(rr.recs) {
		return nil, io.EOF
	}
	rec := rr.recs[rr.cur]
	rr.cur++
	return rec, nil
}

var (
	_ array.Record = (*Record)(nil)
	_ arrio.Reader = (*RecordReader)(nil)
)
