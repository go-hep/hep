// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

func TestLeafDims(t *testing.T) {
	for _, tc := range []struct {
		s    string
		want []int
	}{
		{s: "Leaf", want: nil},
		{s: "Leaf/F", want: nil},
		{s: "Leaf[]", want: []int{-1}},
		{s: "Leaf[]/F", want: []int{-1}},
		{s: "Leaf[20]", want: []int{20}},
		{s: "Leaf[20]/F", want: []int{20}},
		{s: "Leaf[2000]", want: []int{2000}},
		{s: "Leaf[1][2]", want: []int{1, 2}},
		{s: "Leaf[2][]", want: []int{2, -1}},
		{s: "Leaf[1][2][3]", want: []int{1, 2, 3}},
		{s: "Leaf[1][2][3]/F", want: []int{1, 2, 3}},
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

	var (
		br   = new(testBranchImpl)
		scnt = newLeafI(br, "N", nil, signed, nil)
		ucnt = newLeafI(br, "N", nil, unsigned, nil)
	)

	for _, tc := range []struct {
		leaf Leaf
		lcnt Leaf
		data interface{}
	}{
		{
			leaf: newLeafO(br, "BoolTrue", nil, signed, nil),
			data: true,
		},
		{
			leaf: newLeafO(br, "BoolFalse", nil, signed, nil),
			data: false,
		},
		{
			leaf: newLeafB(br, "I8", nil, signed, nil),
			data: int8(-42),
		},
		{
			leaf: newLeafS(br, "I16", nil, signed, nil),
			data: int16(-42),
		},
		{
			leaf: newLeafI(br, "I32", nil, signed, nil),
			data: int32(-42),
		},
		{
			leaf: newLeafL(br, "I64", nil, signed, nil),
			data: int64(-42),
		},
		{
			leaf: newLeafB(br, "U8", nil, unsigned, nil),
			data: uint8(42),
		},
		{
			leaf: newLeafS(br, "U16", nil, unsigned, nil),
			data: uint16(42),
		},
		{
			leaf: newLeafI(br, "U32", nil, unsigned, nil),
			data: uint32(42),
		},
		{
			leaf: newLeafL(br, "U64", nil, unsigned, nil),
			data: uint64(42),
		},
		{
			leaf: newLeafF(br, "F32", nil, signed, nil),
			data: float32(42),
		},
		{
			leaf: newLeafD(br, "F64", nil, signed, nil),
			data: float64(42),
		},
		{
			leaf: newLeafF16(br, "D16", nil, signed, nil, nil),
			data: root.Float16(42),
		},
		{
			leaf: newLeafF16(br, "D16Range", nil, signed, nil, func() rbytes.StreamerElement {
				elm := rdict.Element{
					Name: *rbase.NewNamed("D16Range", "D16Range/f[0, 42]"),
					Type: rmeta.Float16,
				}.New()
				return &elm
			}()),
			data: root.Float16(42),
		},
		{
			leaf: newLeafD32(br, "D32", nil, signed, nil, nil),
			data: root.Double32(42),
		},
		{
			leaf: newLeafD32(br, "D32Range", nil, signed, nil, func() rbytes.StreamerElement {
				elm := rdict.Element{
					Name: *rbase.NewNamed("D32Range", "D32Range/d[0, 42]"),
					Type: rmeta.Double32,
				}.New()
				return &elm
			}()),
			data: root.Double32(42),
		},
		{
			leaf: newLeafO(br, "ArrBools", []int{4}, signed, nil),
			data: [4]bool{true, false, true, false},
		},
		{
			leaf: newLeafB(br, "ArrI8", []int{4}, signed, nil),
			data: [4]int8{1, 2, 3, 4},
		},
		{
			leaf: newLeafS(br, "ArrI16", []int{4}, signed, nil),
			data: [4]int16{1, 2, 3, 4},
		},
		{
			leaf: newLeafI(br, "ArrI32", []int{4}, signed, nil),
			data: [4]int32{1, 2, 3, 4},
		},
		{
			leaf: newLeafL(br, "ArrI64", []int{4}, signed, nil),
			data: [4]int64{1, 2, 3, 4},
		},
		{
			leaf: newLeafB(br, "ArrU8", []int{4}, unsigned, nil),
			data: [4]uint8{1, 2, 3, 4},
		},
		{
			leaf: newLeafS(br, "ArrU16", []int{4}, unsigned, nil),
			data: [4]uint16{1, 2, 3, 4},
		},
		{
			leaf: newLeafI(br, "ArrU32", []int{4}, unsigned, nil),
			data: [4]uint32{1, 2, 3, 4},
		},
		{
			leaf: newLeafL(br, "ArrU64", []int{4}, unsigned, nil),
			data: [4]uint64{1, 2, 3, 4},
		},
		{
			leaf: newLeafF(br, "ArrF32", []int{4}, signed, nil),
			data: [4]float32{1, 2, 3, 4},
		},
		{
			leaf: newLeafD(br, "ArrF64", []int{4}, signed, nil),
			data: [4]float64{1, 2, 3, 4},
		},
		{
			leaf: newLeafD32(br, "ArrD32", []int{4}, signed, nil, nil),
			data: [4]root.Double32{1, 2, 3, 4},
		},
		{
			leaf: newLeafD32(br, "ArrD32Range", []int{4}, signed, nil, func() rbytes.StreamerElement {
				elm := rdict.Element{
					Name: *rbase.NewNamed("ArrD32Range", "ArrD32Range[4]d/[0,4]"),
					Type: rmeta.Double32,
				}.New()
				return &elm
			}()),
			data: [4]root.Double32{0, 4, 0, 4},
		},
		{
			leaf: newLeafO(br, "SliBools", nil, signed, scnt),
			data: []bool{true, false, true, false},
			lcnt: newLeafI(br, "N", nil, signed, nil),
		},
		{
			leaf: newLeafO(br, "SliUBools", nil, unsigned, ucnt),
			data: []bool{true, false, true, false},
			lcnt: newLeafI(br, "N", nil, unsigned, nil),
		},
		{
			leaf: newLeafB(br, "SliI8", nil, signed, scnt),
			data: []int8{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, signed, nil),
		},
		{
			leaf: newLeafS(br, "SliI16", nil, signed, scnt),
			data: []int16{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, signed, nil),
		},
		{
			leaf: newLeafI(br, "SliI32", nil, signed, scnt),
			data: []int32{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, signed, nil),
		},
		{
			leaf: newLeafL(br, "SliI64", nil, signed, scnt),
			data: []int64{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, signed, nil),
		},
		{
			leaf: newLeafB(br, "SliU8", nil, unsigned, ucnt),
			data: []uint8{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, unsigned, nil),
		},
		{
			leaf: newLeafS(br, "SliU16", nil, unsigned, ucnt),
			data: []uint16{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, unsigned, nil),
		},
		{
			leaf: newLeafI(br, "SliU32", nil, unsigned, ucnt),
			data: []uint32{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, unsigned, nil),
		},
		{
			leaf: newLeafL(br, "SliU64", nil, unsigned, ucnt),
			data: []uint64{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, unsigned, nil),
		},
		{
			leaf: newLeafF(br, "SliF32", nil, unsigned, ucnt),
			data: []float32{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, unsigned, nil),
		},
		{
			leaf: newLeafD(br, "SliF64", nil, unsigned, ucnt),
			data: []float64{1, 2, 3, 4},
			lcnt: newLeafI(br, "N", nil, unsigned, nil),
		},
	} {
		t.Run(tc.leaf.Name(), func(t *testing.T) {
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)

			if tc.lcnt != nil {
				tc.leaf.(interface{ setLeafCount(leaf Leaf) }).setLeafCount(tc.lcnt)
				wv := reflect.ValueOf(newValue(tc.lcnt))
				switch {
				case tc.lcnt.IsUnsigned():
					wv.Elem().SetUint(uint64(reflect.ValueOf(tc.data).Len()))
				default:
					wv.Elem().SetInt(int64(reflect.ValueOf(tc.data).Len()))
				}
				err := tc.lcnt.setAddress(wv.Interface())
				if err != nil {
					t.Fatalf("could not setup leaf count: %v", err)
				}

				n, err := tc.lcnt.writeToBuffer(wbuf)
				if err != nil {
					t.Fatalf("could not write count to basket: %v", err)
				}
				if n == 0 {
					t.Fatalf("short write")
				}
			}

			wv := reflect.ValueOf(newValue(tc.leaf))
			wv.Elem().Set(reflect.ValueOf(tc.data))

			if got, want := wv.Elem().Interface(), tc.data; !reflect.DeepEqual(got, want) {
				t.Fatalf("could not setup input data: got=%v, want=%v", got, want)
			}

			err := tc.leaf.setAddress(wv.Interface())
			if err != nil {
				t.Fatalf("could not set write-address: %v", err)
			}

			n, err := tc.leaf.writeToBuffer(wbuf)
			if err != nil {
				t.Fatalf("could not write to basket: %v", err)
			}
			if n == 0 {
				t.Fatalf("short write")
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)

			if tc.lcnt != nil {
				tc.leaf.(interface{ setLeafCount(leaf Leaf) }).setLeafCount(tc.lcnt)
				rv := reflect.ValueOf(newValue(tc.lcnt))
				err := tc.lcnt.setAddress(rv.Interface())
				if err != nil {
					t.Fatalf("could not setup read leaf count: %v", err)
				}

				err = tc.lcnt.readFromBuffer(rbuf)
				if err != nil {
					t.Fatalf("could not write count to basket: %v", err)
				}

				dlen := reflect.ValueOf(tc.data).Len()
				switch {
				case tc.lcnt.IsUnsigned():
					if got, want := int(rv.Elem().Uint()), dlen; got != want {
						t.Fatalf("invalid r/w cycle leaf-count: got=%d, want=%d", got, want)
					}
				default:
					if got, want := int(rv.Elem().Int()), dlen; got != want {
						t.Fatalf("invalid r/w cycle leaf-count: got=%d, want=%d", got, want)
					}
				}
			}

			rv := reflect.ValueOf(newValue(tc.leaf))
			err = tc.leaf.setAddress(rv.Interface())
			if err != nil {
				t.Fatalf("could not set read-address: %v", err)
			}

			err = tc.leaf.readFromBuffer(rbuf)
			if err != nil {
				t.Fatalf("could not read from basket: %v", err)
			}

			if got, want := rv.Elem().Interface(), wv.Elem().Interface(); !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid r/w cycle:\ngot= %v\nwant=%v", got, want)
			}
		})
	}
}

func (leaf *LeafO) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafB) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafS) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafI) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafL) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafF) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafD) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafD32) setLeafCount(lcnt Leaf) { leaf.tleaf.count = lcnt.(leafCount) }
func (leaf *LeafC) setLeafCount(lcnt Leaf)   { leaf.tleaf.count = lcnt.(leafCount) }

type testBranchImpl struct {
	tbranch
}

func (b *testBranchImpl) getReadEntry() int64 { return 1 }

func TestAsLeafBase(t *testing.T) {
	for _, tc := range []struct {
		leaf Leaf
		want rmeta.Enum
	}{
		{
			leaf: new(LeafO),
			want: rmeta.Bool,
		},
		{
			leaf: new(LeafB),
			want: rmeta.Int8,
		},
		{
			leaf: new(LeafS),
			want: rmeta.Int16,
		},
		{
			leaf: new(LeafI),
			want: rmeta.Int32,
		},
		{
			leaf: new(LeafL),
			want: rmeta.Int64,
		},
		{
			leaf: new(LeafF),
			want: rmeta.Float32,
		},
		{
			leaf: new(LeafD),
			want: rmeta.Float64,
		},
		{
			leaf: new(LeafF16),
			want: rmeta.Float16,
		},
		{
			leaf: new(LeafD32),
			want: rmeta.Double32,
		},
		{
			leaf: new(LeafC),
			want: rmeta.CharStar, // FIXME(sbinet): rmeta.Char?
		},
	} {
		t.Run(fmt.Sprintf("%T", tc.leaf), func(t *testing.T) {
			_, got := asLeafBase(tc.leaf)
			if got != tc.want {
				t.Fatalf("got=%v, want=%v", got, tc.want)
			}
		})
	}
}
