// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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
		nevts    = 5
		treeName = "mytree"
	)

	for _, tc := range []struct {
		name    string
		skip    bool
		wvars   []WriteVar
		btitles []string
		total   int
		want    func(i int) interface{}
	}{
		{
			name:    "empty",
			wvars:   []WriteVar{},
			btitles: []string{},
			total:   nevts * (0),
			want:    func(i int) interface{} { return nil },
		},
		{
			name: "simple",
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
			},
			btitles: []string{"i32/I", "f64/D"},
			total:   nevts * (4 + 8),
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
			name: "builtins",
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
			total: nevts * 43,
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
		},
		{
			name: "strings",
			wvars: []WriteVar{
				{Name: "i32", Value: new(int32)},
				{Name: "f64", Value: new(float64)},
				{Name: "str", Value: new(string)},
			},
			btitles: []string{"i32/I", "f64/D", "str/C"},
			total:   nevts * (4 + 8 + (3 + 1)), // 3: strings are "xxx" + 1:string-size
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
		},
		{
			name: "arrays",
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
			total: nevts * 215,
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
		},
		{
			name: "SliI64",
			skip: true, // FIXME(sbinet): var-len arrays not READY yet
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliI64", Value: &[]int64{}, Count: "N"},
			},
			btitles: []string{"N/I", "SliI64[N]/L"},
			total:   nevts*(4+5*8) - 120,
			want: func(i int) interface{} {
				return struct {
					N      int32
					SliI64 []int64
				}{
					N:      int32(i),
					SliI64: []int64{int64(i), int64(i + 1), int64(i + 2), int64(i + 3), int64(i + 4)}[:i],
				}
			},
		},
		{
			name: "SliF64",
			skip: true, // FIXME(sbinet): var-len arrays not READY yet
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliF64", Value: new([]float64), Count: "N"},
			},
			btitles: []string{"N/I", "SliF64[N]/D"},
			total:   nevts*(4+5*8) - 120,
			want: func(i int) interface{} {
				return struct {
					N      int32
					SliF64 []float64
				}{
					N:      int32(i),
					SliF64: []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i],
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("test %s not ready yet", tc.name)
			}

			fname := filepath.Join(tmp, tc.name+".root")

			func() {
				f, err := riofs.Create(fname)
				if err != nil {
					t.Fatalf("could not create write ROOT file %q: %v", fname, err)
				}
				defer f.Close()

				tw, err := NewWriter(f, treeName, tc.wvars)
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
					if got, want := leaf.Title(), tc.wvars[i].Name; got != want {
						t.Fatalf("leaf[%d]: got=%q, want=%q", i, got, want)
					}
				}

				total := 0
				for i := 0; i < nevts; i++ {
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

				if got, want := tw.Entries(), int64(nevts); got != want {
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

				if got, want := tree.Entries(), int64(nevts); got != want {
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
					if got, want := leaf.Title(), tc.wvars[i].Name; got != want {
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

				if got, want := nn, nevts; got != want {
					t.Fatalf("invalid number of events: got=%d, want=%d", got, want)
				}
			}()
		})
	}
}
