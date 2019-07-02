// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbytes"
)

func TestLeafDims(t *testing.T) {
	for _, tc := range []struct {
		s    string
		want []int
	}{
		{s: "Leaf", want: nil},
		{s: "Leaf[]", want: []int{-1}},
		{s: "Leaf[20]", want: []int{20}},
		{s: "Leaf[2000]", want: []int{2000}},
		{s: "Leaf[1][2]", want: []int{1, 2}},
		{s: "Leaf[2][]", want: []int{2, -1}},
		{s: "Leaf[1][2][3]", want: []int{1, 2, 3}},
	} {
		t.Run(tc.s, func(t *testing.T) {
			dims := leafDims(tc.s)
			if !reflect.DeepEqual(dims, tc.want) {
				t.Fatalf("invalid dims. got=%#v, want=%#v", dims, tc.want)
			}
		})
	}
}

func TestLeafReadWriteBasket(t *testing.T) {
	const (
		unsigned = true
		signed   = false
	)

	var br Branch

	for _, tc := range []struct {
		leaf Leaf
		lcnt Leaf
		data interface{}
	}{
		{
			leaf: newLeafO(br, "BoolTrue", 1, signed, nil),
			data: true,
		},
		{
			leaf: newLeafO(br, "BoolFalse", 1, signed, nil),
			data: false,
		},
		{
			leaf: newLeafB(br, "I8", 1, signed, nil),
			data: int8(-42),
		},
		{
			leaf: newLeafS(br, "I16", 1, signed, nil),
			data: int16(-42),
		},
		{
			leaf: newLeafI(br, "I32", 1, signed, nil),
			data: int32(-42),
		},
		{
			leaf: newLeafL(br, "I64", 1, signed, nil),
			data: int64(-42),
		},
		{
			leaf: newLeafB(br, "U8", 1, unsigned, nil),
			data: uint8(42),
		},
		{
			leaf: newLeafS(br, "U16", 1, unsigned, nil),
			data: uint16(42),
		},
		{
			leaf: newLeafI(br, "U32", 1, unsigned, nil),
			data: uint32(42),
		},
		{
			leaf: newLeafL(br, "U64", 1, unsigned, nil),
			data: uint64(42),
		},
		{
			leaf: newLeafF(br, "F32", 1, signed, nil),
			data: float32(42),
		},
		{
			leaf: newLeafD(br, "F64", 1, signed, nil),
			data: float64(42),
		},
		{
			leaf: newLeafO(br, "ArrBools", 4, signed, nil),
			data: [4]bool{true, false, true, false},
		},
		{
			leaf: newLeafB(br, "ArrI8", 4, signed, nil),
			data: [4]int8{1, 2, 3, 4},
		},
		{
			leaf: newLeafS(br, "ArrI16", 4, signed, nil),
			data: [4]int16{1, 2, 3, 4},
		},
		{
			leaf: newLeafI(br, "ArrI32", 4, signed, nil),
			data: [4]int32{1, 2, 3, 4},
		},
		{
			leaf: newLeafL(br, "ArrI64", 4, signed, nil),
			data: [4]int64{1, 2, 3, 4},
		},
		{
			leaf: newLeafB(br, "ArrU8", 4, unsigned, nil),
			data: [4]uint8{1, 2, 3, 4},
		},
		{
			leaf: newLeafS(br, "ArrU16", 4, unsigned, nil),
			data: [4]uint16{1, 2, 3, 4},
		},
		{
			leaf: newLeafI(br, "ArrU32", 4, unsigned, nil),
			data: [4]uint32{1, 2, 3, 4},
		},
		{
			leaf: newLeafL(br, "ArrU64", 4, unsigned, nil),
			data: [4]uint64{1, 2, 3, 4},
		},
		{
			leaf: newLeafF(br, "ArrF32", 4, signed, nil),
			data: [4]float32{1, 2, 3, 4},
		},
		{
			leaf: newLeafD(br, "ArrF64", 4, signed, nil),
			data: [4]float64{1, 2, 3, 4},
		},
	} {
		t.Run(tc.leaf.Name(), func(t *testing.T) {
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			wv := reflect.ValueOf(newValue(tc.leaf))
			wv.Elem().Set(reflect.ValueOf(tc.data))

			if got, want := wv.Elem().Interface(), tc.data; !reflect.DeepEqual(got, want) {
				t.Fatalf("could not setup input data: got=%v, want=%v", got, want)
			}

			err := tc.leaf.setAddress(wv.Interface())
			if err != nil {
				t.Fatalf("could not set write-address: %v", err)
			}

			err = tc.leaf.writeToBasket(wbuf)
			if err != nil {
				t.Fatalf("could not write to basket: %v", err)
			}

			rv := reflect.ValueOf(newValue(tc.leaf))
			err = tc.leaf.setAddress(rv.Interface())
			if err != nil {
				t.Fatalf("could not set read-address: %v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			err = tc.leaf.readFromBasket(rbuf)
			if err != nil {
				t.Fatalf("could not read from basket: %v", err)
			}

			if got, want := rv.Elem().Interface(), wv.Elem().Interface(); !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid r/w cycle:\ngot= %v\nwant=%v", got, want)
			}
		})
	}
}
