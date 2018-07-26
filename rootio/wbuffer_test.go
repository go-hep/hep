// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestWBuffer_WriteBool(t *testing.T) {
	data := make([]byte, 20)
	wbuf := NewWBuffer(data, nil, 0, nil)
	want := true
	wbuf.WriteBool(want)
	rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
	got := rbuf.ReadBool()
	if got != want {
		t.Fatalf("Invalid value. got:%v, want:%v", got, want)
	}
}

func TestWBuffer_WriteString(t *testing.T) {
	data := make([]byte, 520)
	for i := 0; i < 512; i++ {
		wbuf := NewWBuffer(data, nil, 0, nil)
		want := strings.Repeat("=", i)
		wbuf.WriteString(want)
		rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
		got := rbuf.ReadString()
		if got != want {
			t.Fatalf("Invalid value for len=%d.\ngot: %q\nwant:%q", i, got, want)
		}
	}
}

func TestWBuffer_Write(t *testing.T) {
	for _, tc := range []struct {
		buf  []byte
		name string
		want interface{}
		wfct func(*WBuffer, interface{})
		rfct func(*RBuffer) interface{}
		cmp  func(a, b interface{}) bool
	}{
		{
			buf:  make([]byte, 1),
			name: "bool-true",
			want: true,
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteBool(v.(bool))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadBool()
			},
		},
		{
			buf:  make([]byte, 1),
			name: "bool-false",
			want: false,
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteBool(v.(bool))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadBool()
			},
		},
		{
			buf:  make([]byte, 1),
			name: "int8",
			want: int8(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteI8(v.(int8))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadI8()
			},
		},
		{
			buf:  make([]byte, 2),
			name: "int16",
			want: int16(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteI16(v.(int16))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadI16()
			},
		},
		{
			buf:  make([]byte, 4),
			name: "int32",
			want: int32(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteI32(v.(int32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadI32()
			},
		},
		{
			buf:  make([]byte, 8),
			name: "int64",
			want: int64(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteI64(v.(int64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadI64()
			},
		},
		{
			buf:  make([]byte, 1),
			name: "uint8",
			want: uint8(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteU8(v.(uint8))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadU8()
			},
		},
		{
			buf:  make([]byte, 2),
			name: "uint16",
			want: uint16(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteU16(v.(uint16))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadU16()
			},
		},
		{
			buf:  make([]byte, 4),
			name: "uint32",
			want: uint32(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteU32(v.(uint32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadU32()
			},
		},
		{
			buf:  make([]byte, 8),
			name: "uint64",
			want: uint64(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteU64(v.(uint64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadU64()
			},
		},
		{
			buf:  make([]byte, 4),
			name: "float32",
			want: float32(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF32()
			},
		},
		{
			buf:  make([]byte, 4),
			name: "float32-nan",
			want: float32(math.NaN()),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF32()
			},
			cmp: func(a, b interface{}) bool {
				return math.IsNaN(float64(a.(float32))) && math.IsNaN(float64(b.(float32)))
			},
		},
		{
			buf:  make([]byte, 4),
			name: "float32-inf",
			want: float32(math.Inf(-1)),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF32()
			},
			cmp: func(a, b interface{}) bool {
				return math.IsInf(float64(a.(float32)), -1) && math.IsInf(float64(b.(float32)), -1)
			},
		},
		{
			buf:  make([]byte, 4),
			name: "float32+inf",
			want: float32(math.Inf(+1)),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF32()
			},
			cmp: func(a, b interface{}) bool {
				return math.IsInf(float64(a.(float32)), +1) && math.IsInf(float64(b.(float32)), +1)
			},
		},
		{
			buf:  make([]byte, 8),
			name: "float64",
			want: float64(42),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF64()
			},
		},
		{
			buf:  make([]byte, 8),
			name: "float64-nan",
			want: math.NaN(),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF64()
			},
			cmp: func(a, b interface{}) bool {
				return math.IsNaN(a.(float64)) && math.IsNaN(b.(float64))
			},
		},
		{
			buf:  make([]byte, 8),
			name: "float64-inf",
			want: math.Inf(-1),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF64()
			},
			cmp: func(a, b interface{}) bool {
				return math.IsInf(a.(float64), -1) && math.IsInf(b.(float64), -1)
			},
		},
		{
			buf:  make([]byte, 8),
			name: "float64+inf",
			want: math.Inf(+1),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadF64()
			},
			cmp: func(a, b interface{}) bool {
				return math.IsInf(a.(float64), +1) && math.IsInf(b.(float64), +1)
			},
		},
		{
			buf:  make([]byte, len("hello world")+1),
			name: "cstring-1",
			want: "hello world",
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadCString(len("hello world"))
			},
		},
		{
			buf:  make([]byte, len("hello world")+1),
			name: "cstring-2",
			want: "hello world",
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadCString(len("hello world") + 1)
			},
		},
		{
			buf:  make([]byte, len("hello world")+1),
			name: "cstring-3",
			want: "hello world",
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadCString(len("hello world") + 10)
			},
		},
		{
			buf:  make([]byte, len("hello world")+1),
			name: "cstring-4",
			want: "hello",
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadCString(len("hello"))
			},
		},
		{
			buf:  make([]byte, len([]byte{1, 2, 3, 4, 0, 1})+1),
			name: "cstring-5",
			want: string([]byte{1, 2, 3, 4}),
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadCString(len([]byte{1, 2, 3, 4, 0, 1}))
			},
		},
		{
			buf:  make([]byte, 4+5*4),
			name: "static-arr-i32",
			want: []int32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteStaticArrayI32(v.([]int32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadStaticArrayI32()
			},
		},
		{
			buf:  make([]byte, 5*1),
			name: "fast-arr-bool",
			want: []bool{true, false, false, true, false},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayBool(v.([]bool))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayBool(5)
			},
		},
		{
			buf:  make([]byte, 5*1),
			name: "fast-arr-i8",
			want: []int8{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayI8(v.([]int8))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayI8(5)
			},
		},
		{
			buf:  make([]byte, 5*2),
			name: "fast-arr-i16",
			want: []int16{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayI16(v.([]int16))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayI16(5)
			},
		},
		{
			buf:  make([]byte, 5*4),
			name: "fast-arr-i32",
			want: []int32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayI32(v.([]int32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayI32(5)
			},
		},
		{
			buf:  make([]byte, 5*8),
			name: "fast-arr-i64",
			want: []int64{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayI64(v.([]int64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayI64(5)
			},
		},
		{
			buf:  make([]byte, 5*1),
			name: "fast-arr-u8",
			want: []uint8{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayU8(v.([]uint8))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayU8(5)
			},
		},
		{
			buf:  make([]byte, 5*2),
			name: "fast-arr-u16",
			want: []uint16{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayU16(v.([]uint16))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayU16(5)
			},
		},
		{
			buf:  make([]byte, 5*4),
			name: "fast-arr-u32",
			want: []uint32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayU32(v.([]uint32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayU32(5)
			},
		},
		{
			buf:  make([]byte, 5*8),
			name: "fast-arr-u64",
			want: []uint64{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayU64(v.([]uint64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayU64(5)
			},
		},
		{
			buf:  make([]byte, 5*4),
			name: "fast-arr-f32",
			want: []float32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayF32(v.([]float32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayF32(5)
			},
		},
		{
			buf:  make([]byte, 5*4),
			name: "fast-arr-f32-nan+inf-inf",
			want: []float32{1, float32(math.Inf(+1)), 0, float32(math.NaN()), float32(math.Inf(-1))},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayF32(v.([]float32))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayF32(5)
			},
			cmp: func(a, b interface{}) bool {
				aa := a.([]float32)
				bb := b.([]float32)
				if len(aa) != len(bb) {
					return false
				}
				for i := range aa {
					va := float64(aa[i])
					vb := float64(bb[i])
					switch {
					case math.IsNaN(va):
						if !math.IsNaN(vb) {
							return false
						}
					case math.IsNaN(vb):
						if !math.IsNaN(va) {
							return false
						}
					case math.IsInf(va, -1):
						if !math.IsInf(vb, -1) {
							return false
						}
					case math.IsInf(vb, -1):
						if !math.IsInf(va, -1) {
							return false
						}
					case math.IsInf(va, +1):
						if !math.IsInf(vb, +1) {
							return false
						}
					case math.IsInf(vb, +1):
						if !math.IsInf(va, +1) {
							return false
						}
					case va != vb:
						return false
					}
				}
				return true
			},
		},
		{
			buf:  make([]byte, 5*8),
			name: "fast-arr-f64",
			want: []float64{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayF64(v.([]float64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayF64(5)
			},
		},
		{
			buf:  make([]byte, 5*8),
			name: "fast-arr-f64-nan+inf-inf",
			want: []float64{1, math.Inf(+1), 0, math.NaN(), math.Inf(-1)},
			wfct: func(w *WBuffer, v interface{}) {
				w.WriteFastArrayF64(v.([]float64))
			},
			rfct: func(r *RBuffer) interface{} {
				return r.ReadFastArrayF64(5)
			},
			cmp: func(a, b interface{}) bool {
				aa := a.([]float64)
				bb := b.([]float64)
				if len(aa) != len(bb) {
					return false
				}
				for i := range aa {
					va := aa[i]
					vb := bb[i]
					switch {
					case math.IsNaN(va):
						if !math.IsNaN(vb) {
							return false
						}
					case math.IsNaN(vb):
						if !math.IsNaN(va) {
							return false
						}
					case math.IsInf(va, -1):
						if !math.IsInf(vb, -1) {
							return false
						}
					case math.IsInf(vb, -1):
						if !math.IsInf(va, -1) {
							return false
						}
					case math.IsInf(va, +1):
						if !math.IsInf(vb, +1) {
							return false
						}
					case math.IsInf(vb, +1):
						if !math.IsInf(va, +1) {
							return false
						}
					case va != vb:
						return false
					}
				}
				return true
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wbuf := NewWBuffer(tc.buf, nil, 0, nil)
			tc.wfct(wbuf, tc.want)
			if wbuf.err != nil {
				t.Fatalf("error writing to buffer: %v", wbuf.err)
			}
			rbuf := NewRBuffer(tc.buf, nil, 0, nil)
			if rbuf.Err() != nil {
				t.Fatalf("error reading from buffer: %v", rbuf.Err())
			}
			got := tc.rfct(rbuf)
			cmp := reflect.DeepEqual
			if tc.cmp != nil {
				cmp = tc.cmp
			}
			if !cmp(tc.want, got) {
				t.Fatalf("error.\ngot = %v (%T)\nwant= %v (%T)", got, got, tc.want, tc.want)
			}
		})
	}
}

func TestWriteWBuffer(t *testing.T) {
	for _, test := range []struct {
		name string
		file string
		want ROOTMarshaler
	}{
		{
			name: "TObject",
			file: "testdata/tobject.dat",
			want: &tobject{id: 0x0, bits: 0x3000000},
		},
		{
			name: "TNamed",
			file: "testdata/tnamed.dat",
			want: &tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "my-name", title: "my-title"},
		},
		{
			name: "TNamed",
			file: "testdata/tnamed-cmssw.dat",
			want: &tnamed{
				rvers: 1,
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			file: "testdata/tnamed-cmssw-2.dat",
			want: &tnamed{
				rvers: 1,
				obj:   tobject{id: 0x0, bits: 0x3500000},
				name:  "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			file: "testdata/tnamed-long-string.dat",
			want: &tnamed{
				rvers: 1,
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  strings.Repeat("*", 256),
				title: "my-title",
			},
		},
		{
			name: "TArrayI",
			file: "testdata/tarrayi.dat",
			want: &ArrayI{Data: []int32{0, 1, 2, 3, 4}},
		},
		{
			name: "TArrayL64",
			file: "testdata/tarrayl64.dat",
			want: &ArrayL64{Data: []int64{0, 1, 2, 3, 4}},
		},
		{
			name: "TArrayF",
			file: "testdata/tarrayf.dat",
			want: &ArrayF{Data: []float32{0, 1, 2, 3, 4}},
		},
		{
			name: "TArrayD",
			file: "testdata/tarrayd.dat",
			want: &ArrayD{Data: []float64{0, 1, 2, 3, 4}},
		},
	} {
		t.Run("write-buffer="+test.file, func(t *testing.T) {
			testWriteWBuffer(t, test.name, test.file, test.want)
		})
	}
}

func testWriteWBuffer(t *testing.T, name, file string, want interface{}) {
	rdata, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	wdata := make([]byte, len(rdata))
	err = ioutil.WriteFile(file+".new", wdata, 0644)
	if err != nil {
		t.Fatal(err)
	}

	w := NewWBuffer(wdata, nil, 0, nil)
	_, err = want.(ROOTMarshaler).MarshalROOT(w)
	if err != nil {
		t.Fatal(err)
	}

	r := NewRBuffer(wdata, nil, 0, nil)
	obj := Factory.get(name)().Interface().(ROOTUnmarshaler)
	err = obj.UnmarshalROOT(r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(obj, want) {
		t.Fatalf("error: %q\ngot= %+v\nwant=%+v\ngot= %+v\nwant=%+v", file, wdata, rdata, obj, want)
	}

	os.Remove(file + ".new")
}
