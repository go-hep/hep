// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"math"
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
