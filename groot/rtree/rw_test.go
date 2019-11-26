// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"compress/flate"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtypes"
)

func TestBasketRW(t *testing.T) {
	tmp, err := ioutil.TempDir("", "groot-rtree-")
	if err != nil {
		t.Fatalf("could not create temporary directory: %v", err)
	}
	f, err := riofs.Create(filepath.Join(tmp, "basket.root"))
	if err != nil {
		t.Fatalf("could not create temporary file: %v", err)
	}
	defer f.Close()
	defer os.RemoveAll(tmp)

	dir, err := f.Mkdir("data")
	if err != nil {
		t.Fatalf("could not create TDirectory: %v", err)
	}

	var (
		signed = false
		branch = &tbranch{
			named: *rbase.NewNamed("b1", "branch1"),
		}
		leaf = newLeafI(branch, "I32", 1, signed, nil)
	)
	branch.leaves = append(branch.leaves, leaf)

	for _, tc := range []struct {
		basket Basket
	}{
		{
			basket: Basket{
				key:    riofs.KeyFromDir(dir, "empty", "title", "TBasket"),
				branch: branch,
			},
		},
		{
			basket: Basket{
				key:     riofs.KeyFromDir(dir, "simple", "title", "TBasket"),
				bufsize: 5,
				nevsize: 4,
				nevbuf:  3,
				last:    0,
				branch:  branch,
			},
		},
		{
			basket: Basket{
				key:     riofs.KeyFromDir(dir, "with-iobits", "title", "TBasket"),
				bufsize: 5,
				nevsize: 4,
				nevbuf:  3,
				last:    0,
				iobits:  1,
				branch:  branch,
			},
		},
		{
			basket: Basket{
				key:     riofs.KeyFromDir(dir, "with-offsets", "title", "TBasket"),
				bufsize: 5,
				nevsize: 4,
				nevbuf:  4,
				last:    0,
				iobits:  1,
				offsets: []int32{1, 2, 3, 4},
				branch:  branch,
			},
		},
		{
			basket: Basket{
				key:     riofs.KeyFromDir(dir, "with-offsets-displ", "title", "TBasket"),
				bufsize: 5,
				nevsize: 4,
				nevbuf:  4,
				last:    0,
				iobits:  1,
				displ:   []int32{0x11, 0x12, 0x13, 0x14},
				offsets: []int32{1, 2, 3, 4},
				branch:  branch,
			},
		},
		{
			basket: Basket{
				key:     riofs.KeyFromDir(dir, "with-buffer-ref", "title", "TBasket"),
				bufsize: 5,
				nevsize: 4,
				nevbuf:  4,
				last:    1,
				iobits:  1,
				branch:  branch,
				wbuf:    rbytes.NewWBuffer([]byte{42}, nil, 0, nil),
			},
		},
	} {
		t.Run(tc.basket.Name(), func(t *testing.T) {
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)

			var wantBufRef []byte
			if tc.basket.wbuf != nil {
				tc.basket.wbuf.SetPos(int64(tc.basket.last))
				want := tc.basket.wbuf.Bytes()
				wantBufRef = make([]byte, len(want))
				copy(wantBufRef, want)
			}

			n, err := tc.basket.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal basket: n=%d err=%v", n, err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			var b Basket
			err = b.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal basket: %v", err)
			}
			b.branch = branch

			for _, tt := range []struct {
				name      string
				got, want interface{}
			}{
				{"bufsize", b.bufsize, tc.basket.bufsize},
				{"nevsize", b.nevsize, tc.basket.nevsize},
				{"nevbuf", b.nevbuf, tc.basket.nevbuf},
				{"last", b.last, tc.basket.last},
				{"header", b.header, tc.basket.header},
				{"iobits", b.iobits, tc.basket.iobits},
				{"displ", b.displ, tc.basket.displ},
				{"offsets", b.offsets, tc.basket.offsets},
			} {
				if !reflect.DeepEqual(tt.got, tt.want) {
					t.Fatalf("invalid round-trip for %s:\ngot= %#v\nwant=%#v", tt.name, tt.got, tt.want)
				}
			}

			if wantBufRef != nil {
				raw, err := b.key.Bytes()
				if err != nil {
					t.Fatalf("could not unpack key payload: %v", err)
				}
				if got, want := raw, wantBufRef; !reflect.DeepEqual(got, want) {
					t.Fatalf("invalid-roundtrip for wbuf:\ngot= %v\nwant=%v", got, want)
				}
			}
		})
	}

}

func TestIOFeaturesRW(t *testing.T) {
	for _, tc := range []struct {
		name string
		want tioFeatures
	}{
		{"io-0x00", 0x00},
		{"io-0x01", 0x01},
		{"io-0x02", 0x02},
		{"io-0x03", 0x03},
		{"io-0xff", 0xff},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)

			n, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal IOFeatures: n=%d err=%v", n, err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			var got tioFeatures
			err = got.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal IOFeatures: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("invalid round-trip: got=%x, want=%x", got, tc.want)
			}
		})
	}
}

func TestBranchRW(t *testing.T) {
	const (
		unsigned = true
		signed   = false
	)

	tmp, err := ioutil.TempDir("", "groot-rtree-")
	if err != nil {
		t.Fatalf("could not create temporary directory: %v", err)
	}
	f, err := riofs.Create(filepath.Join(tmp, "basket.root"))
	if err != nil {
		t.Fatalf("could not create temporary file: %v", err)
	}
	defer f.Close()
	defer os.RemoveAll(tmp)

	dir, err := f.Mkdir("data")
	if err != nil {
		t.Fatalf("could not create TDirectory: %v", err)
	}

	for _, tc := range []struct {
		name string
		want rtests.ROOTer
	}{
		{
			name: "TBranch",
			want: &tbranch{
				named:          *rbase.NewNamed("branch", "leaf1/I"),
				attfill:        *rbase.NewAttFill(),
				compress:       1,
				basketSize:     defaultBasketSize,
				entryOffsetLen: 0,
				writeBasket:    1,
				entryNumber:    4,
				iobits:         0,
				offset:         0,
				maxBaskets:     10,
				splitLevel:     1,
				entries:        4,
				firstEntry:     0,
				totBytes:       86,
				zipBytes:       86,
				branches:       []Branch{},
				leaves:         []Leaf{},
				baskets:        []Basket{},
				basketBytes:    []int32{86},
				basketEntry:    []int64{0, 4},
				basketSeek:     []int64{304},
				fname:          "foo.root",

				//
				readentry:   -1,
				firstbasket: -1,
				nextbasket:  -1,
			},
		},
		{
			name: "TBranch-with-leaves",
			want: &tbranch{
				named:          *rbase.NewNamed("branch", "leaf1/I:leaf2/L"),
				attfill:        *rbase.NewAttFill(),
				compress:       1,
				basketSize:     defaultBasketSize,
				entryOffsetLen: 0,
				writeBasket:    1,
				entryNumber:    4,
				iobits:         0,
				offset:         0,
				maxBaskets:     10,
				splitLevel:     1,
				entries:        4,
				firstEntry:     0,
				totBytes:       86,
				zipBytes:       86,
				branches:       []Branch{},
				leaves: []Leaf{
					newLeafI(nil, "leaf1", 1, signed, nil),
					newLeafL(nil, "leaf2", 1, signed, nil),
				},
				baskets:     []Basket{},
				basketBytes: []int32{86},
				basketEntry: []int64{0, 4},
				basketSeek:  []int64{304},
				fname:       "foo.root",

				//
				readentry:   -1,
				firstbasket: -1,
				nextbasket:  -1,
			},
		},
		{
			name: "TBranch-with-baskets",
			want: &tbranch{
				named:          *rbase.NewNamed("branch", "leaf1/I:leaf2/L"),
				attfill:        *rbase.NewAttFill(),
				compress:       1,
				basketSize:     defaultBasketSize,
				entryOffsetLen: 0,
				writeBasket:    1,
				entryNumber:    4,
				iobits:         0,
				offset:         0,
				maxBaskets:     10,
				splitLevel:     1,
				entries:        4,
				firstEntry:     0,
				totBytes:       86,
				zipBytes:       86,
				branches:       []Branch{},
				leaves: []Leaf{
					newLeafI(nil, "leaf1", 1, signed, nil),
					newLeafL(nil, "leaf2", 1, signed, nil),
				},
				baskets: []Basket{
					Basket{
						key:     riofs.KeyFromDir(dir, "with-offsets", "title", "TBasket"),
						bufsize: 5,
						nevsize: 4,
						nevbuf:  4,
						last:    0,
						iobits:  1,
						offsets: []int32{1, 2, 3, 4},
						branch:  nil,
					},
				},
				basketBytes: []int32{86},
				basketEntry: []int64{0, 4},
				basketSeek:  []int64{304},
				fname:       "foo.root",

				//
				readentry:   -1,
				firstbasket: -1,
				nextbasket:  -1,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			{
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.SetErr(io.EOF)
				_, err := tc.want.MarshalROOT(wbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
			}

			if b := tc.want.(*tbranch); len(b.leaves) != 0 {
				for i := range b.leaves {
					b.leaves[i].setBranch(b)
				}
			}

			if b := tc.want.(*tbranch); len(b.baskets) != 0 {
				for i := range b.baskets {
					b.baskets[i].branch = b
				}
			}

			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal ROOT: %v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			class := tc.want.Class()
			obj := rtypes.Factory.Get(class)().Interface().(rbytes.Unmarshaler)
			{
				rbuf.SetErr(io.EOF)
				err = obj.UnmarshalROOT(rbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
				rbuf.SetErr(nil)
			}
			err = obj.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal ROOT: %v", err)
			}

			if b := obj.(*tbranch); len(b.baskets) != 0 {
				for i := range b.baskets {
					b.baskets[i].branch = b
					b.baskets[i].key = tc.want.(*tbranch).baskets[i].key
				}
			}

			if !reflect.DeepEqual(obj, tc.want) {
				t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
			}
		})
	}
}

func TestTreeRW(t *testing.T) {
	tmp, err := ioutil.TempDir("", "groot-rtree-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer os.RemoveAll(tmp)

	const (
		treeName = "mytree"
	)

	for _, tc := range []struct {
		name    string
		wopts   []WriteOption
		nevts   int64
		wvars   []WriteVar
		btitles []string
		ltitles []string
		total   int
		want    func(i int) interface{}
		scan    []string // list of branches to use for ROOT TTree::Scan
		cxx     string   // expected ROOT-TTree::Scan
	}{
		{
			name:    "empty",
			nevts:   5,
			wvars:   []WriteVar{},
			btitles: []string{},
			ltitles: []string{},
			total:   5 * (0),
			want:    func(i int) interface{} { return nil },
			cxx: `************
*    Row   *
************
*        0 *
*        1 *
*        2 *
*        3 *
*        4 *
************
`,
		},
		{
			name:  "simple",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
			},
			btitles: []string{"i32/I", "f64/D"},
			ltitles: []string{"i32", "f64"},
			total:   5 * (4 + 8),
			want: func(i int) interface{} {
				return struct {
					I32 int32
					F64 float64
				}{
					I32: int32(i),
					F64: float64(i),
				}
			},
			cxx: `************************************
*    Row   *       i32 *       f64 *
************************************
*        0 *         0 *         0 *
*        1 *         1 *         1 *
*        2 *         2 *         2 *
*        3 *         3 *         3 *
*        4 *         4 *         4 *
************************************
`,
		},
		{
			name:  "builtins",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "B", Value: new(bool)},
				{Name: "I8", Value: new(int8)},
				{Name: "I16", Value: new(int16)},
				{Name: "I32", Value: new(int32)},
				{Name: "I64", Value: new(int64)},
				{Name: "U8", Value: new(uint8)},
				{Name: "U16", Value: new(uint16)},
				{Name: "U32", Value: new(uint32)},
				{Name: "U64", Value: new(uint64)},
				{Name: "F32", Value: new(float32)},
				{Name: "F64", Value: new(float64)},
			},
			btitles: []string{
				"B/O",
				"I8/B", "I16/S", "I32/I", "I64/L",
				"U8/b", "U16/s", "U32/i", "U64/l",
				"F32/F", "F64/D",
			},
			ltitles: []string{
				"B",
				"I8", "I16", "I32", "I64",
				"U8", "U16", "U32", "U64",
				"F32", "F64",
			},
			total: 5 * 43,
			want: func(i int) interface{} {
				return struct {
					B   bool
					I8  int8
					I16 int16
					I32 int32
					I64 int64
					U8  uint8
					U16 uint16
					U32 uint32
					U64 uint64
					F32 float32
					F64 float64
				}{
					B:   bool(i%2 == 0),
					I8:  int8(i),
					I16: int16(i),
					I32: int32(i),
					I64: int64(i),
					U8:  uint8(i),
					U16: uint16(i),
					U32: uint32(i),
					U64: uint64(i),
					F32: float32(i),
					F64: float64(i),
				}
			},
			cxx: `************************************************************************************************************************************************
*    Row   *         B *        I8 *       I16 *       I32 *       I64 *        U8 *       U16 *       U32 *       U64 *       F32 *       F64 *
************************************************************************************************************************************************
*        0 *         1 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *
*        1 *         0 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        2 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        3 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        4 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
************************************************************************************************************************************************
`,
		},
		{
			name:  "strings",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
				{Name: "str", Value: new(string)},
			},
			btitles: []string{"i32/I", "f64/D", "str/C"},
			ltitles: []string{"i32", "f64", "str"},
			total:   5 * (4 + 8 + (3 + 1)), // 3: strings are "xxx" + 1:string-size
			want: func(i int) interface{} {
				return struct {
					I32 int32
					F64 float64
					Str string
				}{
					I32: int32(i),
					F64: float64(i),
					Str: fmt.Sprintf("%03d", i),
				}
			},
			cxx: `************************************************
*    Row   *       i32 *       f64 *       str *
************************************************
*        0 *         0 *         0 *       000 *
*        1 *         1 *         1 *       001 *
*        2 *         2 *         2 *       002 *
*        3 *         3 *         3 *       003 *
*        4 *         4 *         4 *       004 *
************************************************
`,
		},
		{
			name:  "arrays",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "ArrB", Value: new([5]bool)},
				{Name: "ArrI8", Value: new([5]int8)},
				{Name: "ArrI16", Value: new([5]int16)},
				{Name: "ArrI32", Value: new([5]int32)},
				{Name: "ArrI64", Value: new([5]int64)},
				{Name: "ArrU8", Value: new([5]uint8)},
				{Name: "ArrU16", Value: new([5]uint16)},
				{Name: "ArrU32", Value: new([5]uint32)},
				{Name: "ArrU64", Value: new([5]uint64)},
				{Name: "ArrF32", Value: new([5]float32)},
				{Name: "ArrF64", Value: new([5]float64)},
			},
			btitles: []string{
				"ArrB[5]/O",
				"ArrI8[5]/B", "ArrI16[5]/S", "ArrI32[5]/I", "ArrI64[5]/L",
				"ArrU8[5]/b", "ArrU16[5]/s", "ArrU32[5]/i", "ArrU64[5]/l",
				"ArrF32[5]/F", "ArrF64[5]/D",
			},
			ltitles: []string{
				"ArrB[5]",
				"ArrI8[5]", "ArrI16[5]", "ArrI32[5]", "ArrI64[5]",
				"ArrU8[5]", "ArrU16[5]", "ArrU32[5]", "ArrU64[5]",
				"ArrF32[5]", "ArrF64[5]",
			},
			total: 5 * 215,
			want: func(i int) interface{} {
				return struct {
					ArrBool [5]bool
					ArrI8   [5]int8
					ArrI16  [5]int16
					ArrI32  [5]int32
					ArrI64  [5]int64
					ArrU8   [5]uint8
					ArrU16  [5]uint16
					ArrU32  [5]uint32
					ArrU64  [5]uint64
					ArrF32  [5]float32
					ArrF64  [5]float64
				}{
					ArrBool: [5]bool{bool(i%2 == 0), bool((i+1)%2 == 0), bool((i+2)%2 == 0), bool((i+3)%2 == 0), bool((i+4)%2 == 0)},
					ArrI8:   [5]int8{int8(i), int8(i + 1), int8(i + 2), int8(i + 3), int8(i + 4)},
					ArrI16:  [5]int16{int16(i), int16(i + 1), int16(i + 2), int16(i + 3), int16(i + 4)},
					ArrI32:  [5]int32{int32(i), int32(i + 1), int32(i + 2), int32(i + 3), int32(i + 4)},
					ArrI64:  [5]int64{int64(i), int64(i + 1), int64(i + 2), int64(i + 3), int64(i + 4)},
					ArrU8:   [5]uint8{uint8(i), uint8(i + 1), uint8(i + 2), uint8(i + 3), uint8(i + 4)},
					ArrU16:  [5]uint16{uint16(i), uint16(i + 1), uint16(i + 2), uint16(i + 3), uint16(i + 4)},
					ArrU32:  [5]uint32{uint32(i), uint32(i + 1), uint32(i + 2), uint32(i + 3), uint32(i + 4)},
					ArrU64:  [5]uint64{uint64(i), uint64(i + 1), uint64(i + 2), uint64(i + 3), uint64(i + 4)},
					ArrF32:  [5]float32{float32(i), float32(i + 1), float32(i + 2), float32(i + 3), float32(i + 4)},
					ArrF64:  [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)},
				}
			},
			scan: []string{
				"ArrB",
				// "ArrI8", // FIXME(sbinet): ROOT's handling of [X]int8 is sub-par.
				"ArrI16", "ArrI32", "ArrI64",
				"ArrU8", "ArrU16", "ArrU32", "ArrU64",
				"ArrF32", "ArrF64",
			},
			cxx: `***********************************************************************************************************************************************
*    Row   * Instance *      ArrB *    ArrI16 *    ArrI32 *    ArrI64 *     ArrU8 *    ArrU16 *    ArrU32 *    ArrU64 *    ArrF32 *    ArrF64 *
***********************************************************************************************************************************************
*        0 *        0 *         1 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *
*        0 *        1 *         0 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        0 *        2 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        0 *        3 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        0 *        4 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        1 *        0 *         0 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        1 *        1 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        1 *        2 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        1 *        3 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        1 *        4 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        2 *        0 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        2 *        1 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        2 *        2 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        2 *        3 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        2 *        4 *         1 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        3 *        0 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        3 *        1 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        3 *        2 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        3 *        3 *         1 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        3 *        4 *         0 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
*        4 *        0 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        4 *        1 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        4 *        2 *         1 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        4 *        3 *         0 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
*        4 *        4 *         1 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *
***********************************************************************************************************************************************
`,
		},
		{
			name:  "slices",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliB", Value: new([]bool), Count: "N"},
				{Name: "SliI8", Value: new([]int8), Count: "N"},
				{Name: "SliI16", Value: new([]int16), Count: "N"},
				{Name: "SliI32", Value: new([]int32), Count: "N"},
				{Name: "SliI64", Value: new([]int64), Count: "N"},
				{Name: "SliU8", Value: new([]uint8), Count: "N"},
				{Name: "SliU16", Value: new([]uint16), Count: "N"},
				{Name: "SliU32", Value: new([]uint32), Count: "N"},
				{Name: "SliU64", Value: new([]uint64), Count: "N"},
				{Name: "SliF32", Value: new([]float32), Count: "N"},
				{Name: "SliF64", Value: new([]float64), Count: "N"},
			},
			btitles: []string{
				"N/I",
				"SliB[N]/O",
				"SliI8[N]/B", "SliI16[N]/S", "SliI32[N]/I", "SliI64[N]/L",
				"SliU8[N]/b", "SliU16[N]/s", "SliU32[N]/i", "SliU64[N]/l",
				"SliF32[N]/F", "SliF64[N]/D",
			},
			ltitles: []string{
				"N",
				"SliB[N]",
				"SliI8[N]", "SliI16[N]", "SliI32[N]", "SliI64[N]",
				"SliU8[N]", "SliU16[N]", "SliU32[N]", "SliU64[N]",
				"SliF32[N]", "SliF64[N]",
			},
			total: 450,
			want: func(i int) interface{} {
				type Data struct {
					N       int32
					SliBool []bool
					SliI8   []int8
					SliI16  []int16
					SliI32  []int32
					SliI64  []int64
					SliU8   []uint8
					SliU16  []uint16
					SliU32  []uint32
					SliU64  []uint64
					SliF32  []float32
					SliF64  []float64
				}
				if i == 0 {
					return Data{N: 0}
				}
				return Data{
					N:       int32(i),
					SliBool: []bool{bool(i%2 == 0), bool((i+1)%2 == 0), bool((i+2)%2 == 0), bool((i+3)%2 == 0), bool((i+4)%2 == 0)}[:i],
					SliI8:   []int8{int8(i), int8(i + 1), int8(i + 2), int8(i + 3), int8(i + 4)}[:i],
					SliI16:  []int16{int16(i), int16(i + 1), int16(i + 2), int16(i + 3), int16(i + 4)}[:i],
					SliI32:  []int32{int32(i), int32(i + 1), int32(i + 2), int32(i + 3), int32(i + 4)}[:i],
					SliI64:  []int64{int64(i), int64(i + 1), int64(i + 2), int64(i + 3), int64(i + 4)}[:i],
					SliU8:   []uint8{uint8(i), uint8(i + 1), uint8(i + 2), uint8(i + 3), uint8(i + 4)}[:i],
					SliU16:  []uint16{uint16(i), uint16(i + 1), uint16(i + 2), uint16(i + 3), uint16(i + 4)}[:i],
					SliU32:  []uint32{uint32(i), uint32(i + 1), uint32(i + 2), uint32(i + 3), uint32(i + 4)}[:i],
					SliU64:  []uint64{uint64(i), uint64(i + 1), uint64(i + 2), uint64(i + 3), uint64(i + 4)}[:i],
					SliF32:  []float32{float32(i), float32(i + 1), float32(i + 2), float32(i + 3), float32(i + 4)}[:i],
					SliF64:  []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i],
				}
			},
			scan: []string{
				"N",
				"SliB",
				// "SliI8[]", // ROOT's handling of []int8 is sub-par.
				"SliI16", "SliI32", "SliI64",
				"SliU8", "SliU16", "SliU32", "SliU64",
				"SliF32", "SliF64",
			},
			cxx: `***********************************************************************************************************************************************************
*    Row   * Instance *         N *      SliB *    SliI16 *    SliI32 *    SliI64 *     SliU8 *    SliU16 *    SliU32 *    SliU64 *    SliF32 *    SliF64 *
***********************************************************************************************************************************************************
*        0 *        0 *         0 *           *           *           *           *           *           *           *           *           *           *
*        1 *        0 *         1 *         0 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        2 *        0 *         2 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        2 *        1 *         2 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        3 *        0 *         3 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        3 *        1 *         3 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        3 *        2 *         3 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        4 *        0 *         4 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        4 *        1 *         4 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        4 *        2 *         4 *         1 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        4 *        3 *         4 *         0 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
***********************************************************************************************************************************************************
`,
		},
		{
			name:  "slices-multi-baskets",
			nevts: 10000,
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliI64", Value: new([]int64), Count: "N"},
			},
			btitles: []string{
				"N/I",
				"SliI64[N]/L",
			},
			ltitles: []string{
				"N",
				"SliI64[N]",
			},
			total: 400000,
			want: func(i int) interface{} {
				type Data struct {
					N      int32
					SliI64 []int64
				}
				n := i % 10
				var d = Data{N: int32(n)}
				if n == 0 {
					return d
				}
				d.SliI64 = make([]int64, n)
				for j := range d.SliI64 {
					d.SliI64[j] = int64(j + 1)
				}
				return d
			},
		},
		{
			name:  "compr-no-compression",
			wopts: []WriteOption{WithoutCompression()},
			nevts: 500,
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
			},
			btitles: []string{"i32/I", "f64/D"},
			ltitles: []string{"i32", "f64"},
			total:   500 * (4 + 8),
			want: func(i int) interface{} {
				return struct {
					I32 int32
					F64 float64
				}{
					I32: int32(i),
					F64: float64(i),
				}
			},
		},
		{
			name:  "compr-lz4-default",
			wopts: []WriteOption{WithLZ4(flate.DefaultCompression)},
			nevts: 500,
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
			},
			btitles: []string{"i32/I", "f64/D"},
			ltitles: []string{"i32", "f64"},
			total:   500 * (4 + 8),
			want: func(i int) interface{} {
				return struct {
					I32 int32
					F64 float64
				}{
					I32: int32(i),
					F64: float64(i),
				}
			},
		},
		{
			name:  "compr-lzma-default",
			wopts: []WriteOption{WithLZMA(flate.DefaultCompression)},
			nevts: 500,
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
			},
			btitles: []string{"i32/I", "f64/D"},
			ltitles: []string{"i32", "f64"},
			total:   500 * (4 + 8),
			want: func(i int) interface{} {
				return struct {
					I32 int32
					F64 float64
				}{
					I32: int32(i),
					F64: float64(i),
				}
			},
		},
		{
			name:  "compr-zlib-2",
			wopts: []WriteOption{WithZlib(2)},
			nevts: 500,
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
			},
			btitles: []string{"i32/I", "f64/D"},
			ltitles: []string{"i32", "f64"},
			total:   500 * (4 + 8),
			want: func(i int) interface{} {
				return struct {
					I32 int32
					F64 float64
				}{
					I32: int32(i),
					F64: float64(i),
				}
			},
		},
		{
			name:  "compr-zlib-default",
			wopts: []WriteOption{WithZlib(flate.DefaultCompression)},
			nevts: 500,
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
			},
			btitles: []string{"i32/I", "f64/D"},
			ltitles: []string{"i32", "f64"},
			total:   500 * (4 + 8),
			want: func(i int) interface{} {
				return struct {
					I32 int32
					F64 float64
				}{
					I32: int32(i),
					F64: float64(i),
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fname := filepath.Join(tmp, tc.name+".root")

			func() {
				f, err := riofs.Create(fname)
				if err != nil {
					t.Fatalf("could not create write ROOT file %q: %v", fname, err)
				}
				defer f.Close()

				tw, err := NewWriter(f, treeName, tc.wvars, tc.wopts...)
				if err != nil {
					t.Fatalf("could not create tree writer: %v", err)
				}
				for i, b := range tw.Branches() {
					if got, want := b.Name(), tc.wvars[i].Name; got != want {
						t.Fatalf("branch[%d]: got=%q, want=%q", i, got, want)
					}
					if got, want := b.Title(), tc.btitles[i]; got != want {
						t.Fatalf("branch[%d]: got=%q, want=%q", i, got, want)
					}
				}

				for i, leaf := range tw.Leaves() {
					if got, want := leaf.Name(), tc.wvars[i].Name; got != want {
						t.Fatalf("leaf[%d]: got=%q, want=%q", i, got, want)
					}
					if got, want := leaf.Title(), tc.ltitles[i]; got != want {
						t.Fatalf("leaf[%d]: got=%q, want=%q", i, got, want)
					}
				}

				total := 0
				for i := 0; i < int(tc.nevts); i++ {
					want := tc.want(i)
					for j, wvar := range tc.wvars {
						v := reflect.ValueOf(wvar.Value).Elem()
						want := reflect.ValueOf(want).Field(j)
						v.Set(want)
					}
					n, err := tw.Write()
					if err != nil {
						t.Fatalf("could not write event %d: %v", i, err)
					}
					total += n
				}

				if got, want := tw.Entries(), tc.nevts; got != want {
					t.Fatalf("invalid number of entries: got=%d, want=%d", got, want)
				}
				if got, want := total, tc.total; got != want {
					t.Fatalf("invalid number of bytes written: got=%d, want=%d", got, want)
				}

				err = tw.Close()
				if err != nil {
					t.Fatalf("could not close tree writer: %v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("could not close write ROOT file %q: %v", fname, err)
				}
			}()

			func() {
				f, err := riofs.Open(fname)
				if err != nil {
					t.Fatalf("could not opend read ROOT file %q: %+v", fname, err)
				}
				defer f.Close()

				obj, err := f.Get(treeName)
				if err != nil {
					t.Fatalf("could not get ROOT tree %q: %+v", treeName, err)
				}
				tree := obj.(Tree)

				if got, want := tree.Entries(), tc.nevts; got != want {
					t.Fatalf("invalid number of events: got=%v, want=%v", got, want)
				}

				for i, b := range tree.Branches() {
					if got, want := b.Name(), tc.wvars[i].Name; got != want {
						t.Fatalf("branch[%d]: got=%q, want=%q", i, got, want)
					}
					if got, want := b.Title(), tc.btitles[i]; got != want {
						t.Fatalf("branch[%d]: got=%q, want=%q", i, got, want)
					}
				}

				for i, leaf := range tree.Leaves() {
					if got, want := leaf.Name(), tc.wvars[i].Name; got != want {
						t.Fatalf("leaf[%d]: got=%q, want=%q", i, got, want)
					}
					if got, want := leaf.Title(), tc.ltitles[i]; got != want {
						t.Fatalf("leaf[%d]: got=%q, want=%q", i, got, want)
					}
				}

				if len(tc.wvars) == 0 {
					return
				}

				rvars := NewScanVars(tree)
				if len(rvars) != len(tc.wvars) {
					t.Fatalf("invalid number of scan-vars: got=%d, want=%d", len(rvars), len(tc.wvars))
				}

				for i, rvar := range rvars {
					wvar := tc.wvars[i]
					if got, want := rvar.Name, wvar.Name; got != want {
						t.Fatalf("invalid name for rvar[%d]: got=%q, want=%q", i, got, want)
					}
					wtyp := reflect.TypeOf(wvar.Value)
					rtyp := reflect.TypeOf(rvar.Value)
					if got, want := rtyp, wtyp; got != want {
						t.Fatalf("invalid type for rvar[%d]: got=%v, want=%v", i, got, want)
					}
				}

				sc, err := NewScannerVars(tree, rvars...)
				if err != nil {
					t.Fatalf("could not create scanner: %+v", err)
				}
				defer sc.Close()

				nn := 0
				for sc.Next() {
					err := sc.Scan()
					if err != nil {
						t.Fatalf("could not scan entry %d: %+v", sc.Entry(), err)
					}
					want := tc.want(nn)
					for i, rvar := range rvars {
						var (
							want = reflect.ValueOf(want).Field(i).Interface()
							got  = reflect.ValueOf(rvar.Value).Elem().Interface()
						)
						if !reflect.DeepEqual(got, want) {
							t.Errorf("entry[%d]: invalid scan-value[%s]: got=%v, want=%v", nn, tc.wvars[i].Name, got, want)
						}
					}
					nn++
				}

				if sc.Err() != nil {
					t.Fatalf("could not scan tree: %+v", sc.Err())
				}

				if got, want := nn, int(tc.nevts); got != want {
					t.Fatalf("invalid number of events: got=%d, want=%d", got, want)
				}
			}()
			if rtests.HasROOT && tc.cxx != "" {
				code := `#include <iostream>
#include "TFile.h"
#include "TTree.h"
#include "TTreePlayer.h"

void scan(const char* fname, const char* tree, const char *list, const char *oname) {
	auto f = TFile::Open(fname);
	auto t = (TTree*)f->Get(tree);
	if (!t) {
		std::cerr << "could not fetch TTree [" << tree << "] from file [" << fname << "]\n";
		exit(1);
	}
	auto player = dynamic_cast<TTreePlayer*>(t->GetPlayer());
	player->SetScanRedirect(kTRUE);
	player->SetScanFileName(oname);
	t->SetScanField(0);
	t->Scan(list);
}
`

				scan := []string{"*"}
				if len(tc.scan) != 0 {
					scan = tc.scan
				}

				ofile := filepath.Join(tmp, tc.name+".txt")
				out, err := rtests.RunCxxROOT("scan", []byte(code), fname, treeName, strings.Join(scan, ":"), ofile)
				if err != nil {
					t.Fatalf("could not run C++ ROOT: %+v\noutput:\n%s", err, out)
				}

				got, err := ioutil.ReadFile(ofile)
				if err != nil {
					t.Fatalf("could not read C++ ROOT scan file %q: %+v\noutput:\n%s", ofile, err, out)
				}

				if got, want := string(got), tc.cxx; got != want {
					t.Fatalf("invalid ROOT scan:\ngot:\n%v\nwant:\n%v\noutput:\n%s", got, want, out)
				}
			}
		})
	}
}
