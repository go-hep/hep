// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestWBufferGrow(t *testing.T) {
	buf := new(WBuffer)
	if len(buf.w.p) != 0 {
		t.Fatalf("got=%d, want 0-size buffer", len(buf.w.p))
	}

	buf.w.grow(8)
	if got, want := len(buf.w.p), 8; got != want {
		t.Fatalf("got=%d, want=%d buffer len-size", got, want)
	}

	buf.w.grow(8)
	if got, want := len(buf.w.p), 2*8; got != want {
		t.Fatalf("got=%d, want=%d buffer len-size", got, want)
	}

	buf.w.grow(1)
	if got, want := len(buf.w.p), 2*8+1; got != want {
		t.Fatalf("got=%d, want=%d buffer len-size", got, want)
	}

	buf.w.grow(0)
	if got, want := len(buf.w.p), 2*8+1; got != want {
		t.Fatalf("got=%d, want=%d buffer len-size", got, want)
	}
	if got, want := cap(buf.w.p), 8*1024; got != want {
		t.Fatalf("got=%d, want=%d buffer cap-size", got, want)
	}

	const n = 10 * 1024
	buf.w.grow(n)
	if got, want := len(buf.w.p), 2*8+1+n; got != want {
		t.Fatalf("got=%d, want=%d buffer len-size", got, want)
	}
	if got, want := cap(buf.w.p), 8*4*1024; got != want {
		t.Fatalf("got=%d, want=%d buffer cap-size", got, want)
	}

	defer func() {
		e := recover()
		if e == nil {
			t.Fatalf("expected a panic")
		}
	}()

	buf.w.grow(-1)
}

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
	for _, i := range []int{0, 1, 2, 8, 16, 32, 64, 128, 253, 254, 255, 256, 512} {
		t.Run(fmt.Sprintf("str-%03d", i), func(t *testing.T) {
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
	if wbuf.Err() != nil {
		t.Fatalf("err: %v, buf=%v", wbuf.Err(), wbuf.w.p)
	}
}

func TestWBuffer_Write(t *testing.T) {
	for _, tc := range []struct {
		name string
		want any
		wfct func(*WBuffer, any)
		rfct func(*RBuffer) any
		cmp  func(a, b any) bool
	}{
		{
			name: "bool-true",
			want: true,
			wfct: func(w *WBuffer, v any) {
				w.WriteBool(v.(bool))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadBool()
			},
		},
		{
			name: "bool-false",
			want: false,
			wfct: func(w *WBuffer, v any) {
				w.WriteBool(v.(bool))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadBool()
			},
		},
		{
			name: "int8",
			want: int8(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteI8(v.(int8))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadI8()
			},
		},
		{
			name: "int16",
			want: int16(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteI16(v.(int16))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadI16()
			},
		},
		{
			name: "int32",
			want: int32(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteI32(v.(int32))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadI32()
			},
		},
		{
			name: "int64",
			want: int64(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteI64(v.(int64))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadI64()
			},
		},
		{
			name: "uint8",
			want: uint8(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteU8(v.(uint8))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadU8()
			},
		},
		{
			name: "uint16",
			want: uint16(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteU16(v.(uint16))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadU16()
			},
		},
		{
			name: "uint32",
			want: uint32(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteU32(v.(uint32))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadU32()
			},
		},
		{
			name: "uint64",
			want: uint64(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteU64(v.(uint64))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadU64()
			},
		},
		{
			name: "float32",
			want: float32(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF32()
			},
		},
		{
			name: "float32-nan",
			want: float32(math.NaN()),
			wfct: func(w *WBuffer, v any) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF32()
			},
			cmp: func(a, b any) bool {
				return math.IsNaN(float64(a.(float32))) && math.IsNaN(float64(b.(float32)))
			},
		},
		{
			name: "float32-inf",
			want: float32(math.Inf(-1)),
			wfct: func(w *WBuffer, v any) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF32()
			},
			cmp: func(a, b any) bool {
				return math.IsInf(float64(a.(float32)), -1) && math.IsInf(float64(b.(float32)), -1)
			},
		},
		{
			name: "float32+inf",
			want: float32(math.Inf(+1)),
			wfct: func(w *WBuffer, v any) {
				w.WriteF32(v.(float32))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF32()
			},
			cmp: func(a, b any) bool {
				return math.IsInf(float64(a.(float32)), +1) && math.IsInf(float64(b.(float32)), +1)
			},
		},
		{
			name: "float64",
			want: float64(42),
			wfct: func(w *WBuffer, v any) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF64()
			},
		},
		{
			name: "float64-nan",
			want: math.NaN(),
			wfct: func(w *WBuffer, v any) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF64()
			},
			cmp: func(a, b any) bool {
				return math.IsNaN(a.(float64)) && math.IsNaN(b.(float64))
			},
		},
		{
			name: "float64-inf",
			want: math.Inf(-1),
			wfct: func(w *WBuffer, v any) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF64()
			},
			cmp: func(a, b any) bool {
				return math.IsInf(a.(float64), -1) && math.IsInf(b.(float64), -1)
			},
		},
		{
			name: "float64+inf",
			want: math.Inf(+1),
			wfct: func(w *WBuffer, v any) {
				w.WriteF64(v.(float64))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadF64()
			},
			cmp: func(a, b any) bool {
				return math.IsInf(a.(float64), +1) && math.IsInf(b.(float64), +1)
			},
		},
		{
			name: "cstring-1",
			want: "hello world",
			wfct: func(w *WBuffer, v any) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadCString(len("hello world"))
			},
		},
		{
			name: "cstring-2",
			want: "hello world",
			wfct: func(w *WBuffer, v any) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadCString(len("hello world") + 1)
			},
		},
		{
			name: "cstring-3",
			want: "hello world",
			wfct: func(w *WBuffer, v any) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadCString(len("hello world") + 10)
			},
		},
		{
			name: "cstring-4",
			want: "hello",
			wfct: func(w *WBuffer, v any) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadCString(len("hello"))
			},
		},
		{
			name: "cstring-5",
			want: string([]byte{1, 2, 3, 4}),
			wfct: func(w *WBuffer, v any) {
				w.WriteCString(v.(string))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadCString(len([]byte{1, 2, 3, 4, 0, 1}))
			},
		},
		{
			name: "std::string-1",
			want: "hello",
			wfct: func(w *WBuffer, v any) {
				w.WriteStdString(v.(string))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadStdString()
			},
		},
		{
			name: "std::string-2",
			want: strings.Repeat("hello", 256),
			wfct: func(w *WBuffer, v any) {
				w.WriteStdString(v.(string))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadStdString()
			},
		},
		{
			name: "static-arr-i32",
			want: []int32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteStaticArrayI32(v.([]int32))
			},
			rfct: func(r *RBuffer) any {
				return r.ReadStaticArrayI32()
			},
		},
		{
			name: "fast-arr-bool",
			want: []bool{true, false, false, true, false},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayBool(v.([]bool))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]bool, 5)
				r.ReadArrayBool(sli)
				return sli
			},
		},
		{
			name: "fast-arr-i8",
			want: []int8{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayI8(v.([]int8))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]int8, 5)
				r.ReadArrayI8(sli)
				return sli
			},
		},
		{
			name: "fast-arr-i16",
			want: []int16{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayI16(v.([]int16))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]int16, 5)
				r.ReadArrayI16(sli)
				return sli
			},
		},
		{
			name: "fast-arr-i32",
			want: []int32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayI32(v.([]int32))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]int32, 5)
				r.ReadArrayI32(sli)
				return sli
			},
		},
		{
			name: "fast-arr-i64",
			want: []int64{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayI64(v.([]int64))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]int64, 5)
				r.ReadArrayI64(sli)
				return sli
			},
		},
		{
			name: "fast-arr-u8",
			want: []uint8{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayU8(v.([]uint8))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]uint8, 5)
				r.ReadArrayU8(sli)
				return sli
			},
		},
		{
			name: "fast-arr-u16",
			want: []uint16{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayU16(v.([]uint16))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]uint16, 5)
				r.ReadArrayU16(sli)
				return sli
			},
		},
		{
			name: "fast-arr-u32",
			want: []uint32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayU32(v.([]uint32))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]uint32, 5)
				r.ReadArrayU32(sli)
				return sli
			},
		},
		{
			name: "fast-arr-u64",
			want: []uint64{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayU64(v.([]uint64))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]uint64, 5)
				r.ReadArrayU64(sli)
				return sli
			},
		},
		{
			name: "fast-arr-f32",
			want: []float32{1, 2, 0, 2, 1},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayF32(v.([]float32))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]float32, 5)
				r.ReadArrayF32(sli)
				return sli
			},
		},
		{
			name: "fast-arr-f32-nan+inf-inf",
			want: []float32{1, float32(math.Inf(+1)), 0, float32(math.NaN()), float32(math.Inf(-1))},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayF32(v.([]float32))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]float32, 5)
				r.ReadArrayF32(sli)
				return sli
			},
			cmp: func(a, b any) bool {
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
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayF64(v.([]float64))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]float64, 5)
				r.ReadArrayF64(sli)
				return sli
			},
		},
		{
			name: "fast-arr-f64-nan+inf-inf",
			want: []float64{1, math.Inf(+1), 0, math.NaN(), math.Inf(-1)},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayF64(v.([]float64))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]float64, 5)
				r.ReadArrayF64(sli)
				return sli
			},
			cmp: func(a, b any) bool {
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
		{
			name: "fast-arr-str",
			want: []string{"hello", "world"},
			wfct: func(w *WBuffer, v any) {
				w.WriteArrayString(v.([]string))
			},
			rfct: func(r *RBuffer) any {
				sli := make([]string, 2)
				r.ReadArrayString(sli)
				return sli
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wbuf := NewWBuffer(nil, nil, 0, nil)
			tc.wfct(wbuf, tc.want)
			if wbuf.Err() != nil {
				t.Fatalf("error writing to buffer: %v", wbuf.Err())
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
