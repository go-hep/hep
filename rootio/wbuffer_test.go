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
	wbuf := NewWBuffer(nil, nil, 0, nil)
	want := true
	wbuf.WriteBool(want)
	rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
	got := rbuf.ReadBool()
	if got != want {
		t.Fatalf("Invalid value. got:%v, want:%v", got, want)
	}
}

func TestWBuffer_WriteString(t *testing.T) {
	for i := 0; i < 512; i++ {
		t.Run("", func(t *testing.T) {
			wbuf := NewWBuffer(nil, nil, 0, nil)
			want := strings.Repeat("=", i)
			wbuf.WriteString(want)
			rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
			got := rbuf.ReadString()
			if got != want {
				t.Fatalf("Invalid value for len=%d.\ngot: %q\nwant:%q", i, got, want)
			}
		})
	}
}

func TestWBuffer_WriteCString(t *testing.T) {
	wbuf := NewWBuffer(nil, nil, 0, nil)
	want := "hello"
	cstr := string(append([]byte(want), 0))
	wbuf.WriteCString(cstr)
	rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)

	got := rbuf.ReadCString(len(cstr))
	if want != got {
		t.Fatalf("got=%q, want=%q", got, want)
	}
}

func TestWBufferEmpty(t *testing.T) {
	wbuf := new(WBuffer)
	wbuf.WriteString(string([]byte{1, 2, 3, 4, 5}))
	if wbuf.err != nil {
		t.Fatalf("err: %v, buf=%v", wbuf.err, wbuf.w.p)
	}
}

func TestWBuffer_Write(t *testing.T) {
	for _, tc := range []struct {
		name string
		want interface{}
		wfct func(*WBuffer, interface{})
		rfct func(*RBuffer) interface{}
		cmp  func(a, b interface{}) bool
	}{
		{
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
			wbuf := NewWBuffer(nil, nil, 0, nil)
			tc.wfct(wbuf, tc.want)
			if wbuf.err != nil {
				t.Fatalf("error writing to buffer: %v", wbuf.err)
			}
			rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
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
	for _, test := range rwBufferCases {
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

	w := NewWBuffer(nil, nil, 0, nil)
	_, err = want.(ROOTMarshaler).MarshalROOT(w)
	if err != nil {
		t.Fatal(err)
	}
	wdata := w.w.p

	r := NewRBuffer(wdata, nil, 0, nil)
	obj := Factory.get(name)().Interface().(ROOTUnmarshaler)
	err = obj.UnmarshalROOT(r)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(file+".new", wdata, 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(obj, want) {
		t.Fatalf("error: %q\ngot= %+v\nwant=%+v\ngot= %+v\nwant=%+v", file, wdata, rdata, obj, want)
	}

	os.Remove(file + ".new")
}

func TestWRBuffer(t *testing.T) {
	for _, tc := range []struct {
		name string
		want interface {
			ROOTMarshaler
			ROOTUnmarshaler
		}
	}{
		{
			name: "TObject",
			want: &tobject{id: 0x0, bits: 0x3000000},
		},
		{
			name: "TObject",
			want: &tobject{id: 0x1, bits: 0x3000001},
		},
		{
			name: "TNamed",
			want: &tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "my-name", title: "my-title"},
		},
		{
			name: "TNamed",
			want: &tnamed{
				rvers: 1,
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			want: &tnamed{
				rvers: 1,
				obj:   tobject{id: 0x0, bits: 0x3500000},
				name:  "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			want: &tnamed{
				rvers: 1,
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  strings.Repeat("*", 256),
				title: "my-title",
			},
		},
		{
			name: "TList",
			want: &tlist{
				rvers: 5,
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  "list-name",
				objs: []Object{
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n0", title: "t0"},
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n1", title: "t1"},
				},
			},
		},
		{
			name: "TObjString",
			want: &tobjstring{
				rvers: 1,
				obj:   tobject{id: 0x0, bits: 0x3000008},
				str:   "tobjstring-string",
			},
		},
		{
			name: "TObjArray",
			want: &tobjarray{
				rvers: 3,
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  "my-objs",
				arr: []Object{
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n0", title: "t0"},
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n1", title: "t1"},
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n2", title: "t2"},
				},
				last: 2,
			},
		},
		{
			name: "TStreamerBase",
			want: &tstreamerBase{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x0, bits: 0x3000000},
						name:  "TAttLine",
						title: "Line attributes",
					},
					etype:  0,
					esize:  0,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "BASE",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
				rvers: 3,
				vbase: 1,
			},
		},
		{
			name: "TStreamerBasicType",
			want: &tstreamerBasicType{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x0, bits: 0x3000000},
						name:  "fEntries",
						title: "Number of entries",
					},
					etype:  16,
					esize:  8,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "Long64_t",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
				rvers: 2,
			},
		},
		{
			name: "TStreamerBasicType",
			want: &tstreamerBasicType{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x1, bits: 0x3000001},
						name:  "fEntries",
						title: "Array of entries",
					},
					etype:  kOffsetL + kULong,
					esize:  40,
					arrlen: 5,
					arrdim: 1,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "ULong_t",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
				rvers: 2,
			},
		},
		{
			name: "TStreamerBasicType",
			want: &tstreamerBasicType{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x1, bits: 0x3000001},
						name:  "fEntries",
						title: "DynArray of entries",
					},
					etype:  kOffsetP + kULong,
					esize:  8,
					arrlen: 0,
					arrdim: 1,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "ULong_t",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
				rvers: 2,
			},
		},
		{
			name: "TStreamerLoop",
			want: &tstreamerLoop{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x1, bits: 0x3000001},
						name:  "fLoop",
						title: "A streamer loop",
					},
				},
				rvers:  2,
				cvers:  1,
				cname:  "fArrayCount",
				cclass: "MyArrayCount",
			},
		},
		{
			name: "TStreamerObjectAnyPointer",
			want: &tstreamerObjectAnyPointer{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x1, bits: 0x3000001},
						name:  "fObjAnyPtr",
						title: "A pointer to any object",
					},
				},
				rvers: 2,
			},
		},
		{
			name: "TStreamerSTL",
			want: &tstreamerSTL{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x1, bits: 0x3000001},
						name:  "fStdSet",
						title: "A std::set<int>",
					},
					etype: kSTL,
					ename: "std::set<int>",
				},
				rvers: 2,
				vtype: kSTLset,
				ctype: kSTLset,
			},
		},
		{
			name: "TStreamerSTL",
			want: &tstreamerSTL{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x1, bits: 0x3000001},
						name:  "fStdMultimap",
						title: "A std::multimap<int,int>",
					},
					etype: kSTL,
					ename: "std::multimap<int,int>",
				},
				rvers: 2,
				vtype: kSTLmultimap,
				ctype: kSTLmultimap,
			},
		},
		{
			name: "TStreamerSTLstring",
			want: &tstreamerSTLstring{
				tstreamerSTL: tstreamerSTL{
					tstreamerElement: tstreamerElement{
						rvers: 4,
						named: tnamed{
							rvers: 1,
							obj:   tobject{id: 0x1, bits: 0x3000001},
							name:  "fStdString",
							title: "A std::string",
						},
						etype: kSTL,
						ename: "std::string",
					},
					rvers: 2,
					vtype: kSTLstring,
					ctype: kSTLstring,
				},
				rvers: 2,
			},
		},
		{
			name: "TStreamerArtificial",
			want: &tstreamerArtificial{
				tstreamerElement: tstreamerElement{
					rvers: 4,
					named: tnamed{
						rvers: 1,
						obj:   tobject{id: 0x1, bits: 0x3000001},
						name:  "fArtificial",
						title: "An artificial streamer",
					},
					ename: "std::artificial",
				},
				rvers: 2,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wbuf := NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal ROOT: %v", err)
			}

			rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
			obj := Factory.get(tc.name)().Interface().(ROOTUnmarshaler)
			err = obj.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal ROOT: %v", err)
			}

			if !reflect.DeepEqual(obj, tc.want) {
				t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
			}
		})
	}
}
