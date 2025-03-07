// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"compress/flate"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/internal/diff"
)

func TestBasketRW(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-rtree-")
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
		leaf = newLeafI(branch, "I32", nil, signed, nil)
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
				got, want any
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

	tmp, err := os.MkdirTemp("", "groot-rtree-")
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
				ctx: basketCtx{
					entry: -1,
					first: -1,
					next:  -1,
				},
			},
		},
		{
			name: "TBranch-with-leaves",
			want: &tbranch{
				named:          *rbase.NewNamed("branch", "leaf1/I:leaf2/L:leaf3/G"),
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
					newLeafI(nil, "leaf1", nil, signed, nil),
					newLeafL(nil, "leaf2", nil, signed, nil),
					newLeafG(nil, "leaf3", nil, signed, nil),
				},
				baskets:     []Basket{},
				basketBytes: []int32{86},
				basketEntry: []int64{0, 4},
				basketSeek:  []int64{304},
				fname:       "foo.root",

				//
				ctx: basketCtx{
					entry: -1,
					first: -1,
					next:  -1,
				},
			},
		},
		{
			name: "TBranch-with-baskets",
			want: &tbranch{
				named:          *rbase.NewNamed("branch", "leaf1/I:leaf2/L:leaf3/G"),
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
					newLeafI(nil, "leaf1", nil, signed, nil),
					newLeafL(nil, "leaf2", nil, signed, nil),
					newLeafG(nil, "leaf3", nil, signed, nil),
				},
				baskets: []Basket{
					{
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
				ctx: basketCtx{
					entry: -1,
					first: -1,
					next:  -1,
				},
			},
		},
		{
			name: "TBranchElement",
			want: &tbranchElement{
				tbranch: tbranch{
					named:          *rbase.NewNamed("branch", "leaf1/I:leaf2/L:leaf3/G"),
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
						newLeafI(nil, "leaf1", nil, signed, nil),
						newLeafL(nil, "leaf2", nil, signed, nil),
						newLeafG(nil, "leaf3", nil, signed, nil),
					},
					baskets: []Basket{
						{
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
					ctx: basketCtx{
						entry: -1,
						first: -1,
						next:  -1,
					},
				},
				class:  "myclass",
				parent: "parentclass",
				clones: "clones",
				chksum: 123456789,
				clsver: 42,
				id:     3,
				btype:  4,
				stype:  5,
				max:    42,
				//stltyp: 45,
			},
		},
		{
			name: "TBranchElement-with-bcount1",
			want: &tbranchElement{
				tbranch: tbranch{
					named:          *rbase.NewNamed("branch", "leaf1/I:leaf2/L:leaf3/G"),
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
						newLeafI(nil, "leaf1", nil, signed, nil),
						newLeafL(nil, "leaf2", nil, signed, nil),
						newLeafG(nil, "leaf3", nil, signed, nil),
					},
					baskets: []Basket{
						{
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
					ctx: basketCtx{
						entry: -1,
						first: -1,
						next:  -1,
					},
				},
				class:  "myclass",
				parent: "parentclass",
				clones: "clones",
				chksum: 123456789,
				clsver: 42,
				id:     3,
				btype:  4,
				stype:  5,
				max:    42,
				//stltyp: 45,
				bcount1: &tbranchElement{
					tbranch: tbranch{
						named:          *rbase.NewNamed("count", "leaf1/I"),
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
							newLeafI(nil, "leaf1", nil, signed, nil),
						},
						baskets: []Basket{
							{
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
						ctx: basketCtx{
							entry: -1,
							first: -1,
							next:  -1,
						},
					},
					class:  "myotherclass",
					parent: "parentclass",
					clones: "clones",
					chksum: 123456789,
					clsver: 42,
					id:     3,
					btype:  4,
					stype:  5,
					max:    42,
					// stltyp: 45,
				},
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

			asTBranch := func(b any) *tbranch {
				switch b := b.(type) {
				case *tbranch:
					return b
				case *tbranchElement:
					return &b.tbranch
				}
				panic("impossible")
			}

			setupInput := func(b any) {
				if b := asTBranch(b); len(b.leaves) != 0 {
					for i := range b.leaves {
						b.leaves[i].setBranch(b)
					}
				}
				if b := asTBranch(tc.want); len(b.baskets) != 0 {
					for i := range b.baskets {
						b.baskets[i].branch = b
					}
				}

			}

			setupInput(tc.want)

			if b, ok := tc.want.(*tbranchElement); ok {
				if b.bcount1 != nil {
					setupInput(b.bcount1)
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

			if b := asTBranch(obj); len(b.baskets) != 0 {
				for i := range b.baskets {
					b.baskets[i].branch = b
					b.baskets[i].key = asTBranch(tc.want).baskets[i].key
				}
			}

			if b, ok := obj.(*tbranchElement); ok {

				want := tc.want.(*tbranchElement)
				if want.bcount1 != nil {
					for i := range b.bcount1.baskets {
						b.bcount1.baskets[i].branch = want.bcount1.baskets[i].branch
						b.bcount1.baskets[i].key = want.bcount1.baskets[i].key
					}
				}

				var cmpTBE func(n string, got, want *tbranchElement)
				cmpTBE = func(n string, got, want *tbranchElement) {
					if got == nil && want == nil {
						return
					}
					if got == nil && want != nil {
						t.Fatalf("got=%v, want=%v", got, want)
					}
					if got != nil && want == nil {
						t.Fatalf("got=%v, want=%v", got, want)
					}

					for i, v := range []struct {
						got, want any
					}{
						{got.tbranch, want.tbranch},
						{got.class, want.class},
						{got.parent, want.parent},
						{got.clones, want.clones},
						{got.chksum, want.chksum},
						{got.clsver, want.clsver},
						{got.id, want.id},
						{got.btype, want.btype},
						{got.stype, want.stype},
						{got.max, want.max},
						{got.stltyp, want.stltyp},
						{got.streamer, want.streamer},
						{got.estreamer, want.estreamer},
					} {
						if !reflect.DeepEqual(v.got, v.want) {
							t.Fatalf("error[%s-%d]\ngot= %+v\nwant=%+v\n", n, i, v.got, v.want)
						}
					}
					cmpTBE("bcount1", got.bcount1, want.bcount1)
					cmpTBE("bcount2", got.bcount2, want.bcount2)
				}
				cmpTBE("master", b, want)
				return
			}

			if !reflect.DeepEqual(obj, tc.want) {
				t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
			}
		})
	}
}

func TestTreeRW(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-rtree-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer os.RemoveAll(tmp)

	const (
		treeName = "mytree"
	)

	for _, tc := range []struct {
		name    string
		skip    bool
		wopts   []WriteOption
		nevts   int64
		wvars   []WriteVar
		btitles []string
		ltitles []string
		total   int
		want    func(i int) any
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
			want:    func(i int) any { return nil },
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
			want: func(i int) any {
				return struct {
					I32 int32
					F64 float64
				}{
					I32: int32(i),
					F64: float64(i),
				}
			},
			cxx: `************************************
*    Row   *   i32.i32 *   f64.f64 *
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
				{Name: "D16", Value: new(root.Float16)},
				{Name: "D32", Value: new(root.Double32)},
			},
			btitles: []string{
				"B/O",
				"I8/B", "I16/S", "I32/I", "I64/L",
				"U8/b", "U16/s", "U32/i", "U64/l",
				"F32/F", "F64/D", "D16/f", "D32/d",
			},
			ltitles: []string{
				"B",
				"I8", "I16", "I32", "I64",
				"U8", "U16", "U32", "U64",
				"F32", "F64", "D16", "D32",
			},
			total: 5 * 50,
			want: func(i int) any {
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
					D16 root.Float16
					D32 root.Double32
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
					D16: root.Float16(i),
					D32: root.Double32(i),
				}
			},
			cxx: `************************************************************************************************************************************************************************
*    Row   *       B.B *     I8.I8 *   I16.I16 *   I32.I32 *   I64.I64 *     U8.U8 *   U16.U16 *   U32.U32 *   U64.U64 *   F32.F32 *   F64.F64 *   D16.D16 *   D32.D32 *
************************************************************************************************************************************************************************
*        0 *         1 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *
*        1 *         0 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        2 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        3 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        4 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
************************************************************************************************************************************************************************
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
			want: func(i int) any {
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
*    Row   *   i32.i32 *   f64.f64 *   str.str *
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
			name:  "strings-empty",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "s1", Value: new(string)},
				{Name: "s2", Value: new(string)},
			},
			btitles: []string{"s1/C", "s2/C"},
			ltitles: []string{"s1", "s2"},
			total:   30,
			want: func(i int) any {
				return struct {
					S1 string
					S2 string
				}{
					S1: strings.Repeat("x", 4-i),
					S2: strings.Repeat("x", i),
				}
			},
			cxx: `************************************
*    Row   *     s1.s1 *     s2.s2 *
************************************
*        0 *      xxxx *           *
*        1 *       xxx *         x *
*        2 *        xx *        xx *
*        3 *         x *       xxx *
*        4 *           *      xxxx *
************************************
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
			want: func(i int) any {
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
					ArrI8:   [5]int8{'a' + int8(i), int8('a' + i + 1), int8('a' + i + 2), int8('a' + i + 3), int8(0)},
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
				"ArrI8", "ArrI16", "ArrI32", "ArrI64",
				"ArrU8", "ArrU16", "ArrU32", "ArrU64",
				"ArrF32", "ArrF64",
			},
			cxx: `***********************************************************************************************************************************************************
*    Row   * Instance *      ArrB *     ArrI8 *    ArrI16 *    ArrI32 *    ArrI64 *     ArrU8 *    ArrU16 *    ArrU32 *    ArrU64 *    ArrF32 *    ArrF64 *
***********************************************************************************************************************************************************
*        0 *        0 *         1 *      abcd *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *
*        0 *        1 *         0 *      abcd *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        0 *        2 *         1 *      abcd *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        0 *        3 *         0 *      abcd *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        0 *        4 *         1 *      abcd *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        1 *        0 *         0 *      bcde *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        1 *        1 *         1 *      bcde *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        1 *        2 *         0 *      bcde *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        1 *        3 *         1 *      bcde *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        1 *        4 *         0 *      bcde *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        2 *        0 *         1 *      cdef *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        2 *        1 *         0 *      cdef *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        2 *        2 *         1 *      cdef *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        2 *        3 *         0 *      cdef *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        2 *        4 *         1 *      cdef *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        3 *        0 *         0 *      defg *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        3 *        1 *         1 *      defg *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        3 *        2 *         0 *      defg *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        3 *        3 *         1 *      defg *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        3 *        4 *         0 *      defg *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
*        4 *        0 *         1 *      efgh *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        4 *        1 *         0 *      efgh *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        4 *        2 *         1 *      efgh *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        4 *        3 *         0 *      efgh *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
*        4 *        4 *         1 *      efgh *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *
***********************************************************************************************************************************************************
`,
		},
		{
			name:  "arrays-2d",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "ArrB", Value: new([2][3]bool)},
				{Name: "ArrI8", Value: new([2][3]int8)},
				{Name: "ArrI16", Value: new([2][3]int16)},
				{Name: "ArrI32", Value: new([2][3]int32)},
				{Name: "ArrI64", Value: new([2][3]int64)},
				{Name: "ArrU8", Value: new([2][3]uint8)},
				{Name: "ArrU16", Value: new([2][3]uint16)},
				{Name: "ArrU32", Value: new([2][3]uint32)},
				{Name: "ArrU64", Value: new([2][3]uint64)},
				{Name: "ArrF32", Value: new([2][3]float32)},
				{Name: "ArrF64", Value: new([2][3]float64)},
			},
			btitles: []string{
				"ArrB[2][3]/O",
				"ArrI8[2][3]/B", "ArrI16[2][3]/S", "ArrI32[2][3]/I", "ArrI64[2][3]/L",
				"ArrU8[2][3]/b", "ArrU16[2][3]/s", "ArrU32[2][3]/i", "ArrU64[2][3]/l",
				"ArrF32[2][3]/F", "ArrF64[2][3]/D",
			},
			ltitles: []string{
				"ArrB[2][3]",
				"ArrI8[2][3]", "ArrI16[2][3]", "ArrI32[2][3]", "ArrI64[2][3]",
				"ArrU8[2][3]", "ArrU16[2][3]", "ArrU32[2][3]", "ArrU64[2][3]",
				"ArrF32[2][3]", "ArrF64[2][3]",
			},
			total: 5 * 258,
			want: func(i int) any {
				return struct {
					ArrBool [2][3]bool
					ArrI8   [2][3]int8
					ArrI16  [2][3]int16
					ArrI32  [2][3]int32
					ArrI64  [2][3]int64
					ArrU8   [2][3]uint8
					ArrU16  [2][3]uint16
					ArrU32  [2][3]uint32
					ArrU64  [2][3]uint64
					ArrF32  [2][3]float32
					ArrF64  [2][3]float64
				}{
					ArrBool: [2][3]bool{
						{bool(i%2 == 0), bool((i+1)%2 == 0), bool((i+2)%2 == 0)},
						{bool((i+3)%2 == 0), bool((i+4)%2 == 0), bool((i+4)%2 == 0)},
					},
					ArrI8: [2][3]int8{
						{int8(i + 0), int8(i + 1), int8(i + 2)},
						{int8(i + 3), int8(i + 4), int8(i + 5)},
					},
					ArrI16: [2][3]int16{
						{int16(i + 0), int16(i + 1), int16(i + 2)},
						{int16(i + 3), int16(i + 4), int16(i + 5)},
					},
					ArrI32: [2][3]int32{
						{int32(i + 0), int32(i + 1), int32(i + 2)},
						{int32(i + 3), int32(i + 4), int32(i + 5)},
					},
					ArrI64: [2][3]int64{
						{int64(i + 0), int64(i + 1), int64(i + 2)},
						{int64(i + 3), int64(i + 4), int64(i + 5)},
					},
					ArrU8: [2][3]uint8{
						{uint8(i + 0), uint8(i + 1), uint8(i + 2)},
						{uint8(i + 3), uint8(i + 4), uint8(i + 5)},
					},
					ArrU16: [2][3]uint16{
						{uint16(i), uint16(i + 1), uint16(i + 2)},
						{uint16(i + 3), uint16(i + 4), uint16(i + 5)},
					},
					ArrU32: [2][3]uint32{
						{uint32(i + 0), uint32(i + 1), uint32(i + 2)},
						{uint32(i + 3), uint32(i + 4), uint32(i + 5)},
					},
					ArrU64: [2][3]uint64{
						{uint64(i + 0), uint64(i + 1), uint64(i + 2)},
						{uint64(i + 3), uint64(i + 4), uint64(i + 5)},
					},
					ArrF32: [2][3]float32{
						{float32(i + 0), float32(i + 1), float32(i + 2)},
						{float32(i + 3), float32(i + 4), float32(i + 5)},
					},
					ArrF64: [2][3]float64{
						{float64(i + 0), float64(i + 1), float64(i + 2)},
						{float64(i + 3), float64(i + 4), float64(i + 5)},
					},
				}
			},
			scan: []string{
				"ArrB",
				"ArrI8+0", "ArrI16", "ArrI32", "ArrI64",
				"ArrU8", "ArrU16", "ArrU32", "ArrU64",
				"ArrF32", "ArrF64",
			},
			cxx: `***********************************************************************************************************************************************************
*    Row   * Instance *      ArrB *   ArrI8+0 *    ArrI16 *    ArrI32 *    ArrI64 *     ArrU8 *    ArrU16 *    ArrU32 *    ArrU64 *    ArrF32 *    ArrF64 *
***********************************************************************************************************************************************************
*        0 *        0 *         1 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *         0 *
*        0 *        1 *         0 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        0 *        2 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        0 *        3 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        0 *        4 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        0 *        5 *         1 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        1 *        0 *         0 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *         1 *
*        1 *        1 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        1 *        2 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        1 *        3 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        1 *        4 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        1 *        5 *         0 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        2 *        0 *         1 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *         2 *
*        2 *        1 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        2 *        2 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        2 *        3 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        2 *        4 *         1 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        2 *        5 *         1 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
*        3 *        0 *         0 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *         3 *
*        3 *        1 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        3 *        2 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        3 *        3 *         1 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        3 *        4 *         0 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
*        3 *        5 *         0 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *
*        4 *        0 *         1 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *         4 *
*        4 *        1 *         0 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *         5 *
*        4 *        2 *         1 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *         6 *
*        4 *        3 *         0 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *         7 *
*        4 *        4 *         1 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *         8 *
*        4 *        5 *         1 *         9 *         9 *         9 *         9 *         9 *         9 *         9 *         9 *         9 *         9 *
***********************************************************************************************************************************************************
`,
		},
		{
			name:  "slices-uint",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliU8", Value: new([]uint8), Count: "N"},
				{Name: "SliU16", Value: new([]uint16), Count: "N"},
				{Name: "SliU32", Value: new([]uint32), Count: "N"},
				{Name: "SliU64", Value: new([]uint64), Count: "N"},
			},
			btitles: []string{
				"N/I",
				"SliU8[N]/b", "SliU16[N]/s", "SliU32[N]/i", "SliU64[N]/l",
			},
			ltitles: []string{
				"N",
				"SliU8[N]", "SliU16[N]", "SliU32[N]", "SliU64[N]",
			},
			total: 170,
			want: func(i int) any {
				type Data struct {
					N      int32
					SliU8  []uint8
					SliU16 []uint16
					SliU32 []uint32
					SliU64 []uint64
				}
				return Data{
					N:      int32(i),
					SliU8:  []uint8{uint8(i), uint8(i + 1), uint8(i + 2), uint8(i + 3), uint8(i + 4)}[:i],
					SliU16: []uint16{uint16(i), uint16(i + 1), uint16(i + 2), uint16(i + 3), uint16(i + 4)}[:i],
					SliU32: []uint32{uint32(i), uint32(i + 1), uint32(i + 2), uint32(i + 3), uint32(i + 4)}[:i],
					SliU64: []uint64{uint64(i), uint64(i + 1), uint64(i + 2), uint64(i + 3), uint64(i + 4)}[:i],
				}
			},
			scan: []string{
				"N",
				"SliU8", "SliU16", "SliU32", "SliU64",
			},
			cxx: `***********************************************************************************
*    Row   * Instance *         N *     SliU8 *    SliU16 *    SliU32 *    SliU64 *
***********************************************************************************
*        0 *        0 *         0 *           *           *           *           *
*        1 *        0 *         1 *         1 *         1 *         1 *         1 *
*        2 *        0 *         2 *         2 *         2 *         2 *         2 *
*        2 *        1 *         2 *         3 *         3 *         3 *         3 *
*        3 *        0 *         3 *         3 *         3 *         3 *         3 *
*        3 *        1 *         3 *         4 *         4 *         4 *         4 *
*        3 *        2 *         3 *         5 *         5 *         5 *         5 *
*        4 *        0 *         4 *         4 *         4 *         4 *         4 *
*        4 *        1 *         4 *         5 *         5 *         5 *         5 *
*        4 *        2 *         4 *         6 *         6 *         6 *         6 *
*        4 *        3 *         4 *         7 *         7 *         7 *         7 *
***********************************************************************************
`,
		},
		{
			name:  "slices-int",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliI16", Value: new([]int16), Count: "N"},
				{Name: "SliI32", Value: new([]int32), Count: "N"},
				{Name: "SliI64", Value: new([]int64), Count: "N"},
			},
			btitles: []string{
				"N/I",
				"SliI16[N]/S", "SliI32[N]/I", "SliI64[N]/L",
			},
			ltitles: []string{
				"N",
				"SliI16[N]", "SliI32[N]", "SliI64[N]",
			},
			total: 160,
			want: func(i int) any {
				type Data struct {
					N      int32
					SliI16 []int16
					SliI32 []int32
					SliI64 []int64
				}
				return Data{
					N:      int32(i),
					SliI16: []int16{int16(i), int16(i + 1), int16(i + 2), int16(i + 3), int16(i + 4)}[:i],
					SliI32: []int32{int32(i), int32(i + 1), int32(i + 2), int32(i + 3), int32(i + 4)}[:i],
					SliI64: []int64{int64(i), int64(i + 1), int64(i + 2), int64(i + 3), int64(i + 4)}[:i],
				}
			},
			scan: []string{
				"N",
				"SliI16", "SliI32", "SliI64",
			},
			cxx: `***********************************************************************
*    Row   * Instance *         N *    SliI16 *    SliI32 *    SliI64 *
***********************************************************************
*        0 *        0 *         0 *           *           *           *
*        1 *        0 *         1 *         1 *         1 *         1 *
*        2 *        0 *         2 *         2 *         2 *         2 *
*        2 *        1 *         2 *         3 *         3 *         3 *
*        3 *        0 *         3 *         3 *         3 *         3 *
*        3 *        1 *         3 *         4 *         4 *         4 *
*        3 *        2 *         3 *         5 *         5 *         5 *
*        4 *        0 *         4 *         4 *         4 *         4 *
*        4 *        1 *         4 *         5 *         5 *         5 *
*        4 *        2 *         4 *         6 *         6 *         6 *
*        4 *        3 *         4 *         7 *         7 *         7 *
***********************************************************************
`,
		},
		{
			name:  "slices-int8",
			skip:  true,
			nevts: 5,
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliI8", Value: new([]int8), Count: "N"},
			},
			btitles: []string{
				"N/I", "SliI8[N]/B",
			},
			ltitles: []string{
				"N", "SliI8[N]",
			},
			total: 30,
			want: func(i int) any {
				type Data struct {
					N     int32
					SliI8 []int8
				}
				return Data{
					N:     int32(i),
					SliI8: []int8{int8('a' + i), int8('a' + i + 1), int8('a' + i + 2), int8('a' + i + 3), int8(0)}[:i],
				}
			},
			scan: []string{
				"N", "SliI8",
			},
			cxx: `***********************************************
*    Row   * Instance *         N *     SliI8 *
***********************************************
*        0 *        0 *         0 *           *
*        1 *        0 *         1 *         b *
*        2 *        0 *         2 *        cd *
*        2 *        1 *         2 *        cd *
*        3 *        0 *         3 *       def *
*        3 *        1 *         3 *       def *
*        3 *        2 *         3 *       def *
*        4 *        0 *         4 *      efgh *
*        4 *        1 *         4 *      efgh *
*        4 *        2 *         4 *      efgh *
*        4 *        3 *         4 *      efgh *
***********************************************
`,
		},
		{
			name:  "slices-bool-floats",
			nevts: 5,
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "SliB", Value: new([]bool), Count: "N"},
				{Name: "SliF32", Value: new([]float32), Count: "N"},
				{Name: "SliF64", Value: new([]float64), Count: "N"},
			},
			btitles: []string{
				"N/I",
				"SliB[N]/O",
				"SliF32[N]/F", "SliF64[N]/D",
			},
			ltitles: []string{
				"N",
				"SliB[N]",
				"SliF32[N]", "SliF64[N]",
			},
			total: 150,
			want: func(i int) any {
				type Data struct {
					N       int32
					SliBool []bool
					SliF32  []float32
					SliF64  []float64
				}
				return Data{
					N:       int32(i),
					SliBool: []bool{bool(i%2 == 0), bool((i+1)%2 == 0), bool((i+2)%2 == 0), bool((i+3)%2 == 0), bool((i+4)%2 == 0)}[:i],
					SliF32:  []float32{float32(i), float32(i + 1), float32(i + 2), float32(i + 3), float32(i + 4)}[:i],
					SliF64:  []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i],
				}
			},
			scan: []string{
				"N",
				"SliB",
				"SliF32", "SliF64",
			},
			cxx: `***********************************************************************
*    Row   * Instance *         N *      SliB *    SliF32 *    SliF64 *
***********************************************************************
*        0 *        0 *         0 *           *           *           *
*        1 *        0 *         1 *         0 *         1 *         1 *
*        2 *        0 *         2 *         1 *         2 *         2 *
*        2 *        1 *         2 *         0 *         3 *         3 *
*        3 *        0 *         3 *         0 *         3 *         3 *
*        3 *        1 *         3 *         1 *         4 *         4 *
*        3 *        2 *         3 *         0 *         5 *         5 *
*        4 *        0 *         4 *         1 *         4 *         4 *
*        4 *        1 *         4 *         0 *         5 *         5 *
*        4 *        2 *         4 *         1 *         6 *         6 *
*        4 *        3 *         4 *         0 *         7 *         7 *
***********************************************************************
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
			want: func(i int) any {
				type Data struct {
					N      int32
					SliI64 []int64
				}
				n := i % 10
				d := Data{
					N:      int32(n),
					SliI64: make([]int64, n),
				}
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
			want: func(i int) any {
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
			want: func(i int) any {
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
			want: func(i int) any {
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
			want: func(i int) any {
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
			want: func(i int) any {
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

			if tc.skip {
				t.Skipf("skipping %s...", tc.name)
			}

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
				defer tw.Close()

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
				for i := range int(tc.nevts) {
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
					t.Errorf("invalid number of bytes written: got=%d, want=%d", got, want)
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

				rvars := NewReadVars(tree)
				if len(rvars) != len(tc.wvars) {
					t.Fatalf("invalid number of read-vars: got=%d, want=%d", len(rvars), len(tc.wvars))
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

				r, err := NewReader(tree, rvars)
				if err != nil {
					t.Fatalf("could not create reader: %+v", err)
				}
				defer r.Close()

				nn := 0
				err = r.Read(func(ctx RCtx) error {
					i := int(ctx.Entry)
					want := tc.want(i)
					for i, rvar := range rvars {
						var (
							want = reflect.ValueOf(want).Field(i).Interface()
							got  = reflect.ValueOf(rvar.Value).Elem().Interface()
						)
						if !reflect.DeepEqual(got, want) {
							return fmt.Errorf(
								"entry[%d]: invalid scan-value[%s]: got=%v, want=%v",
								ctx.Entry, tc.wvars[i].Name, got, want,
							)
						}
					}
					nn++
					return nil
				})
				if err != nil {
					t.Fatalf("could not read tree: %+v", err)
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

				got, err := os.ReadFile(ofile)
				if err != nil {
					t.Fatalf("could not read C++ ROOT scan file %q: %+v\noutput:\n%s", ofile, err, out)
				}

				if got, want := string(got), tc.cxx; got != want {
					t.Fatalf("invalid ROOT scan:\ngot:\n%v\nwant:\n%v\noutput:\n%s\n%s", got, want, out, diff.Format(got, want))
				}
			}
		})
	}
}

func TestNestedTreeRW(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-rtree-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer os.RemoveAll(tmp)

	const (
		treeName = "mytree"
	)

	sictx := rdict.StreamerInfos

	for _, tc := range []struct {
		name    string
		skip    bool
		wopts   []WriteOption
		nevts   int64
		wvars   []WriteVar
		rvars   []ReadVar
		btitles []string
		ltitles []string
		total   int
		want    func(i int) any
		macro   string // ROOT macro to execute to read back ROOT file
		cxx     string // expected ROOT-TTree::Scan
		sinfos  []rbytes.StreamerInfo
	}{
		{
			name: "struct-with-struct",
			wopts: []WriteOption{
				WithZlib(flate.DefaultCompression),
				WithSplitLevel(0),
			},
			nevts: 10,
			wvars: []WriteVar{
				{Name: "evt", Value: new(TNestedStruct1)},
			},
			rvars: []ReadVar{
				{Name: "evt", Value: new(TNestedStruct1)},
			},
			btitles: []string{"evt"},
			ltitles: []string{"evt"},
			total:   460,
			want: func(i int) any {
				var evt struct {
					Data TNestedStruct1
				}
				evt.Data.RunNbr = 10 + int64(i)
				evt.Data.EvtNbr = int64(i)
				evt.Data.P3.Px = float64(i + 10)
				evt.Data.P3.Py = float64(i + 20)
				evt.Data.P3.Pz = float64(i + 30)
				return evt
			},
			macro: `
#include "TFile.h"
#include "TTree.h"

#include <vector>
#include <fstream>

struct TNestedStruct1P3 {
	double px,py,pz;
};

struct TNestedStruct1 {
	Long64_t runnbr;
	Long64_t evtnbr;
	TNestedStruct1P3 p3;
};


void scan(const char *fname, const char *tname, const char *oname) {
 auto o = std::fstream(oname, std::ofstream::out);
 auto f = TFile::Open(fname, "READ");
 auto t = (TTree*)f->Get(tname);
 t->Print();

 TNestedStruct1 *evt = nullptr;
 t->SetBranchAddress("evt", &evt);

 auto n = t->GetEntries();
 o << "entries: " << n << "\n";
 for (int i = 0; i < n; i++) {
	t->GetEntry(i);
	o << "evt[" << i << "]:"
	  << " run=" << evt->runnbr << ","
	  << " evt=" << evt->evtnbr << ","
	  << " p3(" << evt->p3.px << ", " << evt->p3.py << ", " << evt->p3.pz << ")"
	  << "\n";
 }
 o.flush();
}
			`,
			cxx: `entries: 10
evt[0]: run=10, evt=0, p3(10, 20, 30)
evt[1]: run=11, evt=1, p3(11, 21, 31)
evt[2]: run=12, evt=2, p3(12, 22, 32)
evt[3]: run=13, evt=3, p3(13, 23, 33)
evt[4]: run=14, evt=4, p3(14, 24, 34)
evt[5]: run=15, evt=5, p3(15, 25, 35)
evt[6]: run=16, evt=6, p3(16, 26, 36)
evt[7]: run=17, evt=7, p3(17, 27, 37)
evt[8]: run=18, evt=8, p3(18, 28, 38)
evt[9]: run=19, evt=9, p3(19, 29, 39)
`,
			sinfos: []rbytes.StreamerInfo{
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct1P3{})),
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct1{})),
			},
		},
		{
			name: "large-struct-with-struct",
			wopts: []WriteOption{
				WithZlib(flate.DefaultCompression),
				WithSplitLevel(0),
			},
			nevts: 1000,
			wvars: []WriteVar{
				{Name: "evt", Value: new(TNestedStruct1)},
			},
			rvars: []ReadVar{
				{Name: "evt", Value: new(TNestedStruct1)},
			},
			btitles: []string{"evt"},
			ltitles: []string{"evt"},
			total:   46 * 1000,
			want: func(i int) any {
				var evt struct {
					Data TNestedStruct1
				}
				evt.Data.RunNbr = 10 + int64(i)
				evt.Data.EvtNbr = int64(i)
				evt.Data.P3.Px = float64(i + 10)
				evt.Data.P3.Py = float64(i + 20)
				evt.Data.P3.Pz = float64(i + 30)
				return evt
			},
			sinfos: []rbytes.StreamerInfo{
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct1P3{})),
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct1{})),
			},
		},
		{
			name: "struct-with-struct+slice",
			wopts: []WriteOption{
				WithZlib(flate.DefaultCompression),
				WithSplitLevel(0),
			},
			nevts: 10,
			wvars: []WriteVar{
				{Name: "evt", Value: new(TNestedStruct2)},
			},
			rvars: []ReadVar{
				{Name: "evt", Value: new(TNestedStruct2)},
			},
			btitles: []string{"evt"},
			ltitles: []string{"evt"},
			total:   740,
			want: func(i int) any {
				var evt struct {
					Data TNestedStruct2
				}
				evt.Data.RunNbr = 10 + int64(i)
				evt.Data.EvtNbr = int64(i)
				evt.Data.P3.Px = float64(i + 10)
				evt.Data.P3.Py = float64(i + 20)
				evt.Data.P3.Pz = float64(i + 30)
				switch i {
				case 0:
					evt.Data.F32s = nil
				default:
					evt.Data.F32s = make([]float32, 0, i)
				}
				for j := range i {
					evt.Data.F32s = append(evt.Data.F32s, float32((i+1)*10+j))
				}

				return evt
			},
			macro: `
#include "TFile.h"
#include "TTree.h"

#include <vector>
#include <fstream>

struct TNestedStruct1P3 {
	double px,py,pz;
};

struct TNestedStruct2 {
	Long64_t runnbr;
	Long64_t evtnbr;
	TNestedStruct1P3 p3;
	std::vector<float> f32s;
};

template<class T>
std::string printVec(const std::vector<T>& v) {
	std::stringstream o;
	int i = 0;
	o << "[";
	for (auto e : v) {
		if (i > 0) {
			o << ", ";
		}
		o << e;
		i++;
	}
	o << "]";
	return o.str();
}

void scan(const char *fname, const char *tname, const char *oname) {
 auto o = std::fstream(oname, std::ofstream::out);
 auto f = TFile::Open(fname, "READ");
 auto t = (TTree*)f->Get(tname);
 t->Print();

 TNestedStruct2 *evt = nullptr;
 t->SetBranchAddress("evt", &evt);

 auto n = t->GetEntries();
 o << "entries: " << n << "\n";
 for (int i = 0; i < n; i++) {
	t->GetEntry(i);
	o << "evt[" << i << "]:"
	  << " run=" << evt->runnbr << ","
	  << " evt=" << evt->evtnbr << ","
	  << " p3(" << evt->p3.px << ", " << evt->p3.py << ", " << evt->p3.pz << ")"
	  << " f32s(" << printVec(evt->f32s) << ")"
	  << "\n";
 }
 o.flush();
}
			`,
			cxx: `entries: 10
evt[0]: run=10, evt=0, p3(10, 20, 30) f32s([])
evt[1]: run=11, evt=1, p3(11, 21, 31) f32s([20])
evt[2]: run=12, evt=2, p3(12, 22, 32) f32s([30, 31])
evt[3]: run=13, evt=3, p3(13, 23, 33) f32s([40, 41, 42])
evt[4]: run=14, evt=4, p3(14, 24, 34) f32s([50, 51, 52, 53])
evt[5]: run=15, evt=5, p3(15, 25, 35) f32s([60, 61, 62, 63, 64])
evt[6]: run=16, evt=6, p3(16, 26, 36) f32s([70, 71, 72, 73, 74, 75])
evt[7]: run=17, evt=7, p3(17, 27, 37) f32s([80, 81, 82, 83, 84, 85, 86])
evt[8]: run=18, evt=8, p3(18, 28, 38) f32s([90, 91, 92, 93, 94, 95, 96, 97])
evt[9]: run=19, evt=9, p3(19, 29, 39) f32s([100, 101, 102, 103, 104, 105, 106, 107, 108])
`,
			sinfos: []rbytes.StreamerInfo{
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct1P3{})),
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct2{})),
			},
		},
		{
			name: "struct+slice",
			wopts: []WriteOption{
				WithZlib(flate.DefaultCompression),
				WithSplitLevel(0),
			},
			nevts: 10,
			wvars: []WriteVar{
				{Name: "runnbr", Value: new(int64)},
				{Name: "evtnbr", Value: new(int64)},
				{Name: "p3", Value: new(TNestedStruct1P3)},
				{Name: "f32s", Value: new([]float32)},
			},
			rvars: []ReadVar{
				{Name: "runnbr", Value: new(int64)},
				{Name: "evtnbr", Value: new(int64)},
				{Name: "p3", Value: new(TNestedStruct1P3)},
				{Name: "f32s", Value: new([]float32)},
			},
			btitles: []string{"runnbr/L", "evtnbr/L", "p3", "f32s"},
			ltitles: []string{"runnbr", "evtnbr", "p3", "f32s"},
			total:   680,
			want: func(i int) any {
				var evt struct {
					Data TNestedStruct3
				}
				evt.Data.RunNbr = 10 + int64(i)
				evt.Data.EvtNbr = int64(i)
				evt.Data.P3.Px = float64(i + 10)
				evt.Data.P3.Py = float64(i + 20)
				evt.Data.P3.Pz = float64(i + 30)
				switch i {
				case 0:
					evt.Data.F32s = nil
				default:
					evt.Data.F32s = make([]float32, 0, i)
				}
				for j := range i {
					evt.Data.F32s = append(evt.Data.F32s, float32((i+1)*10+j))
				}

				return evt.Data
			},
			macro: `
#include "TFile.h"
#include "TTree.h"

#include <vector>
#include <fstream>

struct TNestedStruct1P3 {
	double px,py,pz;
};

struct TNestedStruct3 {
	Long64_t runnbr;
	Long64_t evtnbr;
	TNestedStruct1P3 p3;
	std::vector<float> f32s;
};

template<class T>
std::string printVec(const std::vector<T>& v) {
	std::stringstream o;
	int i = 0;
	o << "[";
	for (auto e : v) {
		if (i > 0) {
			o << ", ";
		}
		o << e;
		i++;
	}
	o << "]";
	return o.str();
}

void scan(const char *fname, const char *tname, const char *oname) {
 auto o = std::fstream(oname, std::ofstream::out);
 auto f = TFile::Open(fname, "READ");
 auto t = (TTree*)f->Get(tname);
 t->Print();

 Long64_t runnbr;
 t->SetBranchAddress("runnbr", &runnbr);

 Long64_t evtnbr;
 t->SetBranchAddress("evtnbr", &evtnbr);

 TNestedStruct1P3 *p3 = nullptr;
 t->SetBranchAddress("p3", &p3);

 std::vector<float> *f32s = nullptr;
 t->SetBranchAddress("f32s", &f32s);

 auto n = t->GetEntries();
 o << "entries: " << n << "\n";
 for (int i = 0; i < n; i++) {
	t->GetEntry(i);
	o << "evt[" << i << "]:"
	  << " run=" << runnbr << ","
	  << " evt=" << evtnbr << ","
	  << " p3(" << p3->px << ", " << p3->py << ", " << p3->pz << ")"
	  << " f32s(" << printVec(*f32s) << ")"
	  << "\n";
 }
 o.flush();
}
			`,
			cxx: `entries: 10
evt[0]: run=10, evt=0, p3(10, 20, 30) f32s([])
evt[1]: run=11, evt=1, p3(11, 21, 31) f32s([20])
evt[2]: run=12, evt=2, p3(12, 22, 32) f32s([30, 31])
evt[3]: run=13, evt=3, p3(13, 23, 33) f32s([40, 41, 42])
evt[4]: run=14, evt=4, p3(14, 24, 34) f32s([50, 51, 52, 53])
evt[5]: run=15, evt=5, p3(15, 25, 35) f32s([60, 61, 62, 63, 64])
evt[6]: run=16, evt=6, p3(16, 26, 36) f32s([70, 71, 72, 73, 74, 75])
evt[7]: run=17, evt=7, p3(17, 27, 37) f32s([80, 81, 82, 83, 84, 85, 86])
evt[8]: run=18, evt=8, p3(18, 28, 38) f32s([90, 91, 92, 93, 94, 95, 96, 97])
evt[9]: run=19, evt=9, p3(19, 29, 39) f32s([100, 101, 102, 103, 104, 105, 106, 107, 108])
`,
			sinfos: []rbytes.StreamerInfo{
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct1P3{})),
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedStruct3{})),
			},
		},
		{
			name: "vector+slice",
			wopts: []WriteOption{
				WithZlib(flate.DefaultCompression),
				WithSplitLevel(0),
			},
			nevts: 10,
			wvars: []WriteVar{
				{Name: "N", Value: new(int32)},
				{Name: "vec", Value: new([]float32)},
				//				{Name: "sli", Value: new([]float32), Count: "N"},
			},
			rvars: []ReadVar{
				{Name: "N", Value: new(int32)},
				{Name: "vec", Value: new([]float32)},
				//				{Name: "sli", Value: new([]float32)},
			},
			btitles: []string{"N/I", "vec"}, // "sli[N]/F"},
			ltitles: []string{"N", "vec"},   // "sli[N]"},
			total:   360,
			want: func(i int) any {
				var evt struct {
					N   int32
					Vec []float32
					//					Sli []float32 `groot:"sli[N]"`
				}

				evt.N = int32(i) + 1
				evt.Vec = make([]float32, evt.N)
				//				evt.Sli = make([]float32, evt.N)
				for j := range int(evt.N) {
					evt.Vec[j] = -float32((i+1)*10 + j)
					//					evt.Sli[j] = +float32((i+1)*10 + j)
				}

				return evt
			},
			macro: `
#include "TFile.h"
#include "TTree.h"

#include <vector>
#include <fstream>

template<class T>
std::string printVec(const std::vector<T>& v) {
	std::stringstream o;
	int i = 0;
	o << "[";
	for (auto e : v) {
		if (i > 0) {
			o << " ";
		}
		o << e;
		i++;
	}
	o << "]";
	return o.str();
}

template<class T>
std::string printSli(int32_t n, const T *v) {
	std::stringstream o;
	o << "[";
	for (int i = 0; i < n; i++) {
		auto e = v[i];
		if (i > 0) {
			o << " ";
		}
		o << e;
	}
	o << "]";
	return o.str();
}

void scan(const char *fname, const char *tname, const char *oname) {
 auto o = std::fstream(oname, std::ofstream::out);
 auto f = TFile::Open(fname, "READ");
 auto t = (TTree*)f->Get(tname);
 t->Print();

 int32_t n;
 t->SetBranchAddress("N", &n);

 std::vector<float> *vec = nullptr;
 t->SetBranchAddress("vec", &vec);

// float *sli = nullptr;
// t->SetBranchAddress("sli", &sli);

 auto nevts = t->GetEntries();
 o << "entries: " << nevts << "\n";
 for (int i = 0; i < nevts; i++) {
	t->GetEntry(i);
	o << "evt[" << i << "]:"
	  << " " << n
	  << " " << printVec(*vec)
//	  << " " << printSli(n, sli)
	  << "\n";
 }
 o.flush();
}
			`,
			cxx: `entries: 10
evt[0]: 1 [-10]
evt[1]: 2 [-20 -21]
evt[2]: 3 [-30 -31 -32]
evt[3]: 4 [-40 -41 -42 -43]
evt[4]: 5 [-50 -51 -52 -53 -54]
evt[5]: 6 [-60 -61 -62 -63 -64 -65]
evt[6]: 7 [-70 -71 -72 -73 -74 -75 -76]
evt[7]: 8 [-80 -81 -82 -83 -84 -85 -86 -87]
evt[8]: 9 [-90 -91 -92 -93 -94 -95 -96 -97 -98]
evt[9]: 10 [-100 -101 -102 -103 -104 -105 -106 -107 -108 -109]
`,
			sinfos: []rbytes.StreamerInfo{
				rdict.StreamerOf(sictx, reflect.TypeOf([]float32{})),
			},
		},
		{
			name: "event-nosplit",
			wopts: []WriteOption{
				WithZlib(flate.DefaultCompression),
				WithSplitLevel(0),
			},
			nevts: 10,
			wvars: []WriteVar{
				{Name: "evt", Value: new(TNestedEvent1)},
			},
			rvars: []ReadVar{
				{Name: "evt", Value: new(TNestedEvent1)},
			},
			btitles: []string{"evt"},
			ltitles: []string{"evt"},
			total:   25520,
			want: func(i int) any {
				var evt struct {
					Event TNestedEvent1
				}

				evt.Event = TNestedEvent1{}.want(int64(i))

				return evt
			},
			macro: `
#include "TFile.h"
#include "TTree.h"
#include "TString.h"
#include "TObjString.h"

#include <vector>
#include <fstream>

const int ARRAYSZ  = 10;
const int MAXSLICE = 20;
const int MAXSTR   = 32;

#define OFFSET 0

struct TNestedP2 {
	double px;
	float  py;
};


template<class T>
std::string printV(T v) {
	std::stringstream o;
	o << v;
	return o.str();
}

template<>
std::string printV(TNestedP2 v) {
	std::stringstream o;
	o << "{" << v.px << " " << v.py << "}";
	return o.str();
}

template<>
std::string printV(bool v) {
	if (v) {
		return "true";
	}
	return "false";
}

template<>
std::string printV(int8_t v) {
	std::stringstream o;
	o << int(v);
	return o.str();
}

template<>
std::string printV(uint8_t v) {
	std::stringstream o;
	o << int(v);
	return o.str();
}

template<>
std::string printV(TString v) {
	std::stringstream o;
	o << v.Data();
	return o.str();
}

template<>
std::string printV(TObjString v) {
	std::stringstream o;
	o << v.GetString().Data();
	return o.str();
}

template<>
std::string printV(std::string v) {
	return v;
}

template<class T>
std::string printArr(const T *v) {
	std::stringstream o;
	o << "[";
	for (int i = 0; i < ARRAYSZ; i++) {
		auto e = v[i];
		if (i > 0) {
			o << " ";
		}
		o << printV(e);
	}
	o << "]";
	return o.str();
}

//template<>
std::string printArrCStr(char *v[ARRAYSZ]) {
	std::stringstream o;
	o << "[";
	for (int i = 0; i < ARRAYSZ; i++) {
		auto e = v[i];
		if (i > 0) {
			o << " ";
		}
		o << e;
	}
	o << "]";
	return o.str();
}

template<class T>
std::string printSli(int32_t n, const T *v) {
	std::stringstream o;
	o << "[";
	for (int i = 0; i < n; i++) {
		auto e = v[i];
		if (i > 0) {
			o << " ";
		}
		o << printV(e);
	}
	o << "]";
	return o.str();
}

std::string printSliStr(int32_t n, char *v[ARRAYSZ]) {
	std::stringstream o;
	o << "[";
	for (int i = 0; i < n; i++) {
		auto e = v[i];
		if (i > 0) {
			o << " ";
		}
		o << e;
	}
	o << "]";
	return o.str();
}

template<class T>
std::string printVec(const std::vector<T>& v) {
	std::stringstream o;
	int i = 0;
	o << "[";
	for (auto e : v) {
		if (i > 0) {
			o << " ";
		}
		o << printV(e);
		i++;
	}
	o << "]";
	return o.str();
}

template<class T>
std::string printVecVec(const std::vector<std::vector<T> >& v) {
	std::stringstream o;
	int i = 0;
	o << "[";
	for (auto e : v) {
		if (i > 0) {
			o << " ";
		}
		o << printVec(e);
		i++;
	}
	o << "]";
	return o.str();
}

struct TNestedEvent1 {
	bool     Bool;
	//char     Str[MAXSTR];
	//char    *Str;
	std::string Str;
	int8_t   I8;
	int16_t  I16;
	int32_t  I32;
	int64_t  I64;
	uint8_t  U8;
	uint16_t U16;
	uint32_t U32;
	uint64_t U64;
	float    F32;
	double   F64;

	Float16_t  D16;
	Double32_t D32;

	TNestedP2  P2;
	TObjString Obj;

	bool     ArrBs[ARRAYSZ];
//	TString  ArrStr[ARRAYSZ];
	int8_t   ArrI8[ARRAYSZ];
	int16_t  ArrI16[ARRAYSZ];
	int32_t  ArrI32[ARRAYSZ];
	int64_t  ArrI64[ARRAYSZ];
	uint8_t  ArrU8[ARRAYSZ];
	uint16_t ArrU16[ARRAYSZ];
	uint32_t ArrU32[ARRAYSZ];
	uint64_t ArrU64[ARRAYSZ];
	float    ArrF32[ARRAYSZ];
	double   ArrF64[ARRAYSZ];

	Float16_t    ArrD16[ARRAYSZ];
	Double32_t   ArrD32[ARRAYSZ];

	TNestedP2  ArrP2[ARRAYSZ];
	TObjString ArrObj[ARRAYSZ];

	int32_t  N;
	bool     *SliBs;   //[N]
//	char*    *SliStr;  //[N]
	int8_t   *SliI8;   //[N]
	int16_t  *SliI16;  //[N]
	int32_t  *SliI32;  //[N]
	int64_t  *SliI64;  //[N]
	uint8_t  *SliU8;   //[N]
	uint16_t *SliU16;  //[N]
	uint32_t *SliU32;  //[N]
	uint64_t *SliU64;  //[N]
	float    *SliF32;  //[N]
	double   *SliF64;  //[N]

	Float16_t  *SliD16; //[N]
	Double32_t *SliD32; //[N]
//	TNestedP2  *SliP2;  //[N]
//	TObjString *SliObj; //[N]

	std::vector<bool> StdVecBs;
	std::vector<std::string> StdVecStr;
	std::vector<int8_t>  StdVecI8;
	std::vector<int16_t> StdVecI16;
	std::vector<int32_t> StdVecI32;
	std::vector<int64_t> StdVecI64;
	std::vector<uint8_t>  StdVecU8;
	std::vector<uint16_t> StdVecU16;
	std::vector<uint32_t> StdVecU32;
	std::vector<uint64_t> StdVecU64;
	std::vector<float>    StdVecF32;
	std::vector<double>   StdVecF64;

	std::vector<Float16_t>  StdVecD16;
	std::vector<Double32_t> StdVecD32;
	std::vector<TNestedP2>  StdVecP2;
	std::vector<TObjString> StdVecObj;

	std::vector<std::vector<double> >      StdVecVecF64;
	std::vector<std::vector<std::string> > StdVecVecStr;
	std::vector<std::vector<TNestedP2> >   StdVecVecP2;
};

TNestedEvent1 *newEvent() {
	auto *evt = new TNestedEvent1;
	evt->SliBs = (bool*)malloc(sizeof(bool)*0);
//	evt->SliStr = (char**)malloc(sizeof(char*)*0);
	evt->SliI8  = (int8_t*)malloc(sizeof(int8_t)*0);
	evt->SliI16 = (int16_t*)malloc(sizeof(int16_t)*0);
	evt->SliI32 = (int32_t*)malloc(sizeof(int32_t)*0);
	evt->SliI64 = (int64_t*)malloc(sizeof(int64_t)*0);
	evt->SliU8  = (uint8_t*)malloc(sizeof(uint8_t)*0);
	evt->SliU16 = (uint16_t*)malloc(sizeof(uint16_t)*0);
	evt->SliU32 = (uint32_t*)malloc(sizeof(uint32_t)*0);
	evt->SliU64 = (uint64_t*)malloc(sizeof(uint64_t)*0);
	evt->SliF32 = (float*)malloc(sizeof(float)*0);
	evt->SliF64 = (double*)malloc(sizeof(double)*0);

	evt->SliD16 = (Float16_t*)malloc(sizeof(Float16_t)*0);
	evt->SliD32 = (Double32_t*)malloc(sizeof(Double32_t)*0);
//	evt->SliP2 =  (TNestedP2*)malloc(sizeof(TNestedP2)*0);
//	evt->SliObj = (TObjString*)malloc(sizeof(TObjString)*0);

	return evt;
}


void scan(const char *fname, const char *tname, const char *oname) {
 auto o = std::fstream(oname, std::ofstream::out);
 auto f = TFile::Open(fname, "READ");
 auto t = (TTree*)f->Get(tname);
 t->Print();

 TNestedEvent1 *evt = newEvent();
 t->SetBranchAddress("evt", &evt);

 o << "key[000]: " << tname << ";1 \"\" (TTree)\n";

 auto n = t->GetEntries();
 for (int i = 0; i < n; i++) {
	t->GetEntry(i);
	o << "[00" << i << "][evt]: "
	  << "{" << printV(evt->Bool)
	  << " " << printV(evt->Str)
	  << " " << printV(int(evt->I8))
	  << " " << printV(evt->I16)
	  << " " << printV(evt->I32)
	  << " " << printV(evt->I64)
	  << " " << printV(int(evt->U8))
	  << " " << printV(evt->U16)
	  << " " << printV(evt->U32)
	  << " " << printV(evt->U64)
	  << " " << printV(evt->F32)
	  << " " << printV(evt->F64)
	  << " " << printV(evt->D16)
	  << " " << printV(evt->D32)
	  << " " << printV(evt->P2)
	  << " " << printV(evt->Obj)

	  << " " << printArr(evt->ArrBs)
//	  << " " << printArr(evt->ArrStr)
	  << " " << printArr(evt->ArrI8)
	  << " " << printArr(evt->ArrI16)
	  << " " << printArr(evt->ArrI32)
	  << " " << printArr(evt->ArrI64)
	  << " " << printArr(evt->ArrU8)
	  << " " << printArr(evt->ArrU16)
	  << " " << printArr(evt->ArrU32)
	  << " " << printArr(evt->ArrU64)
	  << " " << printArr(evt->ArrF32)
	  << " " << printArr(evt->ArrF64)
	  << " " << printArr(evt->ArrD16)
	  << " " << printArr(evt->ArrD32)
	  << " " << printArr(evt->ArrP2)
	  << " " << printArr(evt->ArrObj)

	  << " " << evt->N
	  << " " << printSli(evt->N, evt->SliBs)
//	  << " " << printSli(evt->N, evt->SliStr)
	  << " " << printSli(evt->N, evt->SliI8)
	  << " " << printSli(evt->N, evt->SliI16)
	  << " " << printSli(evt->N, evt->SliI32)
	  << " " << printSli(evt->N, evt->SliI64)
	  << " " << printSli(evt->N, evt->SliU8)
	  << " " << printSli(evt->N, evt->SliU16)
	  << " " << printSli(evt->N, evt->SliU32)
	  << " " << printSli(evt->N, evt->SliU64)
	  << " " << printSli(evt->N, evt->SliF32)
	  << " " << printSli(evt->N, evt->SliF64)
	  << " " << printSli(evt->N, evt->SliD16)
	  << " " << printSli(evt->N, evt->SliD32)
//	  << " " << printSli(evt->N, evt->SliP2)
//	  << " " << printSli(evt->N, evt->SliObj)

	  << " " << printVec(evt->StdVecBs)
	  << " " << printVec(evt->StdVecStr)
	  << " " << printVec(evt->StdVecI8)
	  << " " << printVec(evt->StdVecI16)
	  << " " << printVec(evt->StdVecI32)
	  << " " << printVec(evt->StdVecI64)
	  << " " << printVec(evt->StdVecU8)
	  << " " << printVec(evt->StdVecU16)
	  << " " << printVec(evt->StdVecU32)
	  << " " << printVec(evt->StdVecU64)
	  << " " << printVec(evt->StdVecF32)
	  << " " << printVec(evt->StdVecF64)
	  << " " << printVec(evt->StdVecD16)
	  << " " << printVec(evt->StdVecD32)
	  << " " << printVec(evt->StdVecP2)
	  << " " << printVec(evt->StdVecObj)

	  << " " << printVecVec(evt->StdVecVecF64)
	  << " " << printVecVec(evt->StdVecVecStr)
	  << " " << printVecVec(evt->StdVecVecP2)
	  << "}\n";
 }
 o.flush();
}
			`,
			cxx: `key[000]: mytree;1 "" (TTree)
[000][evt]: {true str-000 0 0 0 0 0 0 0 0 0 0 0 0 {0 0} obj-0 [true false false false false false false false false false] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [{0 0} {0 0} {0 0} {0 0} {0 0} {0 0} {0 0} {0 0} {0 0} {0 0}] [obj-0 obj-0 obj-0 obj-0 obj-0 obj-0 obj-0 obj-0 obj-0 obj-0] 0 [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] [] []}
[001][evt]: {false str-001 -1 -1 -1 -1 1 1 1 1 1 1 1 1 {1 1} obj-1 [false true false false false false false false false false] [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1] [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1] [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1] [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1] [1 1 1 1 1 1 1 1 1 1] [1 1 1 1 1 1 1 1 1 1] [1 1 1 1 1 1 1 1 1 1] [1 1 1 1 1 1 1 1 1 1] [1 1 1 1 1 1 1 1 1 1] [1 1 1 1 1 1 1 1 1 1] [1 1 1 1 1 1 1 1 1 1] [1 1 1 1 1 1 1 1 1 1] [{1 1} {1 1} {1 1} {1 1} {1 1} {1 1} {1 1} {1 1} {1 1} {1 1}] [obj-1 obj-1 obj-1 obj-1 obj-1 obj-1 obj-1 obj-1 obj-1 obj-1] 1 [true] [-1] [-1] [-1] [-1] [1] [1] [1] [1] [1] [1] [1] [1] [true] [std-001] [-1] [-1] [-1] [-1] [1] [1] [1] [1] [1] [1] [1] [1] [{1 1}] [obj-001] [[0 1 2 3]] [[vec-001 vec-002 vec-003 vec-004]] [[{1 1} {2 2} {3 3} {4 4}]]}
[002][evt]: {true str-002 -2 -2 -2 -2 2 2 2 2 2 2 2 2 {2 2} obj-2 [false false true false false false false false false false] [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2] [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2] [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2] [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2] [2 2 2 2 2 2 2 2 2 2] [2 2 2 2 2 2 2 2 2 2] [2 2 2 2 2 2 2 2 2 2] [2 2 2 2 2 2 2 2 2 2] [2 2 2 2 2 2 2 2 2 2] [2 2 2 2 2 2 2 2 2 2] [2 2 2 2 2 2 2 2 2 2] [2 2 2 2 2 2 2 2 2 2] [{2 2} {2 2} {2 2} {2 2} {2 2} {2 2} {2 2} {2 2} {2 2} {2 2}] [obj-2 obj-2 obj-2 obj-2 obj-2 obj-2 obj-2 obj-2 obj-2 obj-2] 2 [false true] [-2 -2] [-2 -2] [-2 -2] [-2 -2] [2 2] [2 2] [2 2] [2 2] [2 2] [2 2] [2 2] [2 2] [false true] [std-002 std-002] [-2 -2] [-2 -2] [-2 -2] [-2 -2] [2 2] [2 2] [2 2] [2 2] [2 2] [2 2] [2 2] [2 2] [{2 2} {2 2}] [obj-002 obj-002] [[0 1 2 3] [1 2 3 4]] [[vec-002 vec-003 vec-004 vec-005] [vec-002 vec-003 vec-004 vec-005]] [[{2 2} {3 3} {4 4} {5 5}] [{2 2} {3 3} {4 4} {5 5}]]}
[003][evt]: {false str-003 -3 -3 -3 -3 3 3 3 3 3 3 3 3 {3 3} obj-3 [false false false true false false false false false false] [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3] [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3] [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3] [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3] [3 3 3 3 3 3 3 3 3 3] [3 3 3 3 3 3 3 3 3 3] [3 3 3 3 3 3 3 3 3 3] [3 3 3 3 3 3 3 3 3 3] [3 3 3 3 3 3 3 3 3 3] [3 3 3 3 3 3 3 3 3 3] [3 3 3 3 3 3 3 3 3 3] [3 3 3 3 3 3 3 3 3 3] [{3 3} {3 3} {3 3} {3 3} {3 3} {3 3} {3 3} {3 3} {3 3} {3 3}] [obj-3 obj-3 obj-3 obj-3 obj-3 obj-3 obj-3 obj-3 obj-3 obj-3] 3 [false false true] [-3 -3 -3] [-3 -3 -3] [-3 -3 -3] [-3 -3 -3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [false false true] [std-003 std-003 std-003] [-3 -3 -3] [-3 -3 -3] [-3 -3 -3] [-3 -3 -3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [3 3 3] [{3 3} {3 3} {3 3}] [obj-003 obj-003 obj-003] [[0 1 2 3] [1 2 3 4] [2 3 4 5]] [[vec-003 vec-004 vec-005 vec-006] [vec-003 vec-004 vec-005 vec-006] [vec-003 vec-004 vec-005 vec-006]] [[{3 3} {4 4} {5 5} {6 6}] [{3 3} {4 4} {5 5} {6 6}] [{3 3} {4 4} {5 5} {6 6}]]}
[004][evt]: {true str-004 -4 -4 -4 -4 4 4 4 4 4 4 4 4 {4 4} obj-4 [false false false false true false false false false false] [-4 -4 -4 -4 -4 -4 -4 -4 -4 -4] [-4 -4 -4 -4 -4 -4 -4 -4 -4 -4] [-4 -4 -4 -4 -4 -4 -4 -4 -4 -4] [-4 -4 -4 -4 -4 -4 -4 -4 -4 -4] [4 4 4 4 4 4 4 4 4 4] [4 4 4 4 4 4 4 4 4 4] [4 4 4 4 4 4 4 4 4 4] [4 4 4 4 4 4 4 4 4 4] [4 4 4 4 4 4 4 4 4 4] [4 4 4 4 4 4 4 4 4 4] [4 4 4 4 4 4 4 4 4 4] [4 4 4 4 4 4 4 4 4 4] [{4 4} {4 4} {4 4} {4 4} {4 4} {4 4} {4 4} {4 4} {4 4} {4 4}] [obj-4 obj-4 obj-4 obj-4 obj-4 obj-4 obj-4 obj-4 obj-4 obj-4] 4 [false false false true] [-4 -4 -4 -4] [-4 -4 -4 -4] [-4 -4 -4 -4] [-4 -4 -4 -4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [false false false true] [std-004 std-004 std-004 std-004] [-4 -4 -4 -4] [-4 -4 -4 -4] [-4 -4 -4 -4] [-4 -4 -4 -4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [4 4 4 4] [{4 4} {4 4} {4 4} {4 4}] [obj-004 obj-004 obj-004 obj-004] [[0 1 2 3] [1 2 3 4] [2 3 4 5] [3 4 5 6]] [[vec-004 vec-005 vec-006 vec-007] [vec-004 vec-005 vec-006 vec-007] [vec-004 vec-005 vec-006 vec-007] [vec-004 vec-005 vec-006 vec-007]] [[{4 4} {5 5} {6 6} {7 7}] [{4 4} {5 5} {6 6} {7 7}] [{4 4} {5 5} {6 6} {7 7}] [{4 4} {5 5} {6 6} {7 7}]]}
[005][evt]: {false str-005 -5 -5 -5 -5 5 5 5 5 5 5 5 5 {5 5} obj-5 [false false false false false true false false false false] [-5 -5 -5 -5 -5 -5 -5 -5 -5 -5] [-5 -5 -5 -5 -5 -5 -5 -5 -5 -5] [-5 -5 -5 -5 -5 -5 -5 -5 -5 -5] [-5 -5 -5 -5 -5 -5 -5 -5 -5 -5] [5 5 5 5 5 5 5 5 5 5] [5 5 5 5 5 5 5 5 5 5] [5 5 5 5 5 5 5 5 5 5] [5 5 5 5 5 5 5 5 5 5] [5 5 5 5 5 5 5 5 5 5] [5 5 5 5 5 5 5 5 5 5] [5 5 5 5 5 5 5 5 5 5] [5 5 5 5 5 5 5 5 5 5] [{5 5} {5 5} {5 5} {5 5} {5 5} {5 5} {5 5} {5 5} {5 5} {5 5}] [obj-5 obj-5 obj-5 obj-5 obj-5 obj-5 obj-5 obj-5 obj-5 obj-5] 5 [false false false false true] [-5 -5 -5 -5 -5] [-5 -5 -5 -5 -5] [-5 -5 -5 -5 -5] [-5 -5 -5 -5 -5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [false false false false true] [std-005 std-005 std-005 std-005 std-005] [-5 -5 -5 -5 -5] [-5 -5 -5 -5 -5] [-5 -5 -5 -5 -5] [-5 -5 -5 -5 -5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [5 5 5 5 5] [{5 5} {5 5} {5 5} {5 5} {5 5}] [obj-005 obj-005 obj-005 obj-005 obj-005] [[0 1 2 3] [1 2 3 4] [2 3 4 5] [3 4 5 6] [4 5 6 7]] [[vec-005 vec-006 vec-007 vec-008] [vec-005 vec-006 vec-007 vec-008] [vec-005 vec-006 vec-007 vec-008] [vec-005 vec-006 vec-007 vec-008] [vec-005 vec-006 vec-007 vec-008]] [[{5 5} {6 6} {7 7} {8 8}] [{5 5} {6 6} {7 7} {8 8}] [{5 5} {6 6} {7 7} {8 8}] [{5 5} {6 6} {7 7} {8 8}] [{5 5} {6 6} {7 7} {8 8}]]}
[006][evt]: {true str-006 -6 -6 -6 -6 6 6 6 6 6 6 6 6 {6 6} obj-6 [false false false false false false true false false false] [-6 -6 -6 -6 -6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6 -6 -6 -6 -6] [6 6 6 6 6 6 6 6 6 6] [6 6 6 6 6 6 6 6 6 6] [6 6 6 6 6 6 6 6 6 6] [6 6 6 6 6 6 6 6 6 6] [6 6 6 6 6 6 6 6 6 6] [6 6 6 6 6 6 6 6 6 6] [6 6 6 6 6 6 6 6 6 6] [6 6 6 6 6 6 6 6 6 6] [{6 6} {6 6} {6 6} {6 6} {6 6} {6 6} {6 6} {6 6} {6 6} {6 6}] [obj-6 obj-6 obj-6 obj-6 obj-6 obj-6 obj-6 obj-6 obj-6 obj-6] 6 [false false false false false true] [-6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [false false false false false true] [std-006 std-006 std-006 std-006 std-006 std-006] [-6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6] [-6 -6 -6 -6 -6 -6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [6 6 6 6 6 6] [{6 6} {6 6} {6 6} {6 6} {6 6} {6 6}] [obj-006 obj-006 obj-006 obj-006 obj-006 obj-006] [[0 1 2 3] [1 2 3 4] [2 3 4 5] [3 4 5 6] [4 5 6 7] [5 6 7 8]] [[vec-006 vec-007 vec-008 vec-009] [vec-006 vec-007 vec-008 vec-009] [vec-006 vec-007 vec-008 vec-009] [vec-006 vec-007 vec-008 vec-009] [vec-006 vec-007 vec-008 vec-009] [vec-006 vec-007 vec-008 vec-009]] [[{6 6} {7 7} {8 8} {9 9}] [{6 6} {7 7} {8 8} {9 9}] [{6 6} {7 7} {8 8} {9 9}] [{6 6} {7 7} {8 8} {9 9}] [{6 6} {7 7} {8 8} {9 9}] [{6 6} {7 7} {8 8} {9 9}]]}
[007][evt]: {false str-007 -7 -7 -7 -7 7 7 7 7 7 7 7 7 {7 7} obj-7 [false false false false false false false true false false] [-7 -7 -7 -7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7 -7 -7 -7] [7 7 7 7 7 7 7 7 7 7] [7 7 7 7 7 7 7 7 7 7] [7 7 7 7 7 7 7 7 7 7] [7 7 7 7 7 7 7 7 7 7] [7 7 7 7 7 7 7 7 7 7] [7 7 7 7 7 7 7 7 7 7] [7 7 7 7 7 7 7 7 7 7] [7 7 7 7 7 7 7 7 7 7] [{7 7} {7 7} {7 7} {7 7} {7 7} {7 7} {7 7} {7 7} {7 7} {7 7}] [obj-7 obj-7 obj-7 obj-7 obj-7 obj-7 obj-7 obj-7 obj-7 obj-7] 7 [false false false false false false true] [-7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [false false false false false false true] [std-007 std-007 std-007 std-007 std-007 std-007 std-007] [-7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7] [-7 -7 -7 -7 -7 -7 -7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [7 7 7 7 7 7 7] [{7 7} {7 7} {7 7} {7 7} {7 7} {7 7} {7 7}] [obj-007 obj-007 obj-007 obj-007 obj-007 obj-007 obj-007] [[0 1 2 3] [1 2 3 4] [2 3 4 5] [3 4 5 6] [4 5 6 7] [5 6 7 8] [6 7 8 9]] [[vec-007 vec-008 vec-009 vec-010] [vec-007 vec-008 vec-009 vec-010] [vec-007 vec-008 vec-009 vec-010] [vec-007 vec-008 vec-009 vec-010] [vec-007 vec-008 vec-009 vec-010] [vec-007 vec-008 vec-009 vec-010] [vec-007 vec-008 vec-009 vec-010]] [[{7 7} {8 8} {9 9} {10 10}] [{7 7} {8 8} {9 9} {10 10}] [{7 7} {8 8} {9 9} {10 10}] [{7 7} {8 8} {9 9} {10 10}] [{7 7} {8 8} {9 9} {10 10}] [{7 7} {8 8} {9 9} {10 10}] [{7 7} {8 8} {9 9} {10 10}]]}
[008][evt]: {true str-008 -8 -8 -8 -8 8 8 8 8 8 8 8 8 {8 8} obj-8 [false false false false false false false false true false] [-8 -8 -8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8 -8 -8] [8 8 8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8 8 8] [{8 8} {8 8} {8 8} {8 8} {8 8} {8 8} {8 8} {8 8} {8 8} {8 8}] [obj-8 obj-8 obj-8 obj-8 obj-8 obj-8 obj-8 obj-8 obj-8 obj-8] 8 [false false false false false false false true] [-8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [false false false false false false false true] [std-008 std-008 std-008 std-008 std-008 std-008 std-008 std-008] [-8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8] [-8 -8 -8 -8 -8 -8 -8 -8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [8 8 8 8 8 8 8 8] [{8 8} {8 8} {8 8} {8 8} {8 8} {8 8} {8 8} {8 8}] [obj-008 obj-008 obj-008 obj-008 obj-008 obj-008 obj-008 obj-008] [[0 1 2 3] [1 2 3 4] [2 3 4 5] [3 4 5 6] [4 5 6 7] [5 6 7 8] [6 7 8 9] [7 8 9 10]] [[vec-008 vec-009 vec-010 vec-011] [vec-008 vec-009 vec-010 vec-011] [vec-008 vec-009 vec-010 vec-011] [vec-008 vec-009 vec-010 vec-011] [vec-008 vec-009 vec-010 vec-011] [vec-008 vec-009 vec-010 vec-011] [vec-008 vec-009 vec-010 vec-011] [vec-008 vec-009 vec-010 vec-011]] [[{8 8} {9 9} {10 10} {11 11}] [{8 8} {9 9} {10 10} {11 11}] [{8 8} {9 9} {10 10} {11 11}] [{8 8} {9 9} {10 10} {11 11}] [{8 8} {9 9} {10 10} {11 11}] [{8 8} {9 9} {10 10} {11 11}] [{8 8} {9 9} {10 10} {11 11}] [{8 8} {9 9} {10 10} {11 11}]]}
[009][evt]: {false str-009 -9 -9 -9 -9 9 9 9 9 9 9 9 9 {9 9} obj-9 [false false false false false false false false false true] [-9 -9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9 -9] [9 9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9 9] [{9 9} {9 9} {9 9} {9 9} {9 9} {9 9} {9 9} {9 9} {9 9} {9 9}] [obj-9 obj-9 obj-9 obj-9 obj-9 obj-9 obj-9 obj-9 obj-9 obj-9] 9 [false false false false false false false false true] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [false false false false false false false false true] [std-009 std-009 std-009 std-009 std-009 std-009 std-009 std-009 std-009] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [-9 -9 -9 -9 -9 -9 -9 -9 -9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [9 9 9 9 9 9 9 9 9] [{9 9} {9 9} {9 9} {9 9} {9 9} {9 9} {9 9} {9 9} {9 9}] [obj-009 obj-009 obj-009 obj-009 obj-009 obj-009 obj-009 obj-009 obj-009] [[0 1 2 3] [1 2 3 4] [2 3 4 5] [3 4 5 6] [4 5 6 7] [5 6 7 8] [6 7 8 9] [7 8 9 10] [8 9 10 11]] [[vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012] [vec-009 vec-010 vec-011 vec-012]] [[{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}] [{9 9} {10 10} {11 11} {12 12}]]}
`,
			sinfos: []rbytes.StreamerInfo{
				rdict.StreamerOf(sictx, reflect.TypeOf([]float64{})),
				rdict.StreamerOf(sictx, reflect.TypeOf([]string{})),
				rdict.StreamerOf(sictx, reflect.TypeOf([][]float64{})),
				rdict.StreamerOf(sictx, reflect.TypeOf([][]string{})),
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedP2{})),
				rdict.StreamerOf(sictx, reflect.TypeOf([]TNestedP2{})),
				rdict.StreamerOf(sictx, reflect.TypeOf([][]TNestedP2{})),
				rdict.StreamerOf(sictx, reflect.TypeOf(TNestedEvent1{})),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fname := filepath.Join(tmp, tc.name+".root")

			if tc.skip {
				t.Skipf("skipping %s...", tc.name)
			}

			for i := range tc.sinfos {
				rdict.StreamerInfos.Add(tc.sinfos[i])
			}

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
				defer tw.Close()

				if got, want := len(tw.Branches()), len(tc.wvars); got != want {
					t.Fatalf("invalid number of branches: got=%d, want=%d", got, want)
				}

				for i, b := range tw.Branches() {
					if got, want := b.Name(), tc.wvars[i].Name; got != want {
						t.Fatalf("branch[%d]: got=%q, want=%q", i, got, want)
					}
					if got, want := b.Title(), tc.btitles[i]; got != want {
						t.Fatalf("branch[%d]: got=%q, want=%q", i, got, want)
					}
				}

				if got, want := len(tw.Leaves()), len(tc.wvars); got != want {
					leaves := make([]string, got)
					for i, l := range tw.Leaves() {
						leaves[i] = l.Name()
					}
					t.Fatalf("invalid number of leaves: got=%d, want=%d (leaves: %q)", got, want, leaves)
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
				for i := range int(tc.nevts) {
					want := tc.want(i)
					for j, wvar := range tc.wvars {
						v := reflect.ValueOf(wvar.Value).Elem()
						o := reflect.ValueOf(want).Field(j)
						v.Set(o)
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
					t.Errorf("invalid number of bytes written: got=%d, want=%d", got, want)
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

				rvars := tc.rvars
				if len(rvars) != len(tc.wvars) {
					t.Fatalf("invalid number of read-vars: got=%d, want=%d", len(rvars), len(tc.wvars))
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

				r, err := NewReader(tree, rvars)
				if err != nil {
					t.Fatalf("could not create reader: %+v", err)
				}
				defer r.Close()

				nn := 0
				err = r.Read(func(ctx RCtx) error {
					i := int(ctx.Entry)
					want := tc.want(i)
					for i, rvar := range rvars {
						var (
							want = reflect.ValueOf(want).Field(i).Interface()
							got  = reflect.ValueOf(rvar.Value).Elem().Interface()
						)
						if !reflect.DeepEqual(got, want) {
							return fmt.Errorf(
								"entry[%d]: invalid scan-value[%s]:\ngot= %v\nwant=%v",
								ctx.Entry, tc.wvars[i].Name, got, want,
							)
						}
					}
					nn++
					return nil
				})
				if err != nil {
					t.Fatalf("could not read tree: %+v", err)
				}

				if got, want := nn, int(tc.nevts); got != want {
					t.Fatalf("invalid number of events: got=%d, want=%d", got, want)
				}
			}()

			if rtests.HasROOT && len(tc.macro) != 0 {
				ofile := filepath.Join(tmp, tc.name+".txt")
				out, err := rtests.RunCxxROOT("scan", []byte(tc.macro), fname, treeName, ofile)
				if err != nil {
					t.Fatalf("could not run C++ ROOT: %+v\noutput:\n%s", err, out)
				}

				got, err := os.ReadFile(ofile)
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

func TestTreeWriteSubdir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-rtree-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, "tree-subdir.root")

	f, err := riofs.Create(fname)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer f.Close()

	dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11/dir-111")
	if err != nil {
		t.Fatalf("could not create sub-dir hierarchy: %+v", err)
	}

	var data struct {
		I32 int32   `groot:"i32"`
		F64 float64 `groot:"f64"`
	}

	ntup, err := NewWriter(dir, "ntup", WriteVarsFromStruct(&data))
	if err != nil {
		t.Fatalf("could not create tree: %+v", err)
	}
	defer ntup.Close()

	for i := range 5 {
		data.I32 = int32(i)
		data.F64 = float64(i)
		_, err = ntup.Write()
		if err != nil {
			t.Fatalf("could not write event %d: %+v", i, err)
		}
	}

	err = ntup.Close()
	if err != nil {
		t.Fatalf("could not close tree: %+v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("could not close file: %+v", err)
	}

	if !rtests.HasROOT {
		return
	}

	code := `#include <iostream>
#include "TDirectory.h"
#include "TFile.h"
#include "TTree.h"
#include "TTreePlayer.h"

void scan(const char *fname, const char *tree, const char *oname) {
	auto f = TFile::Open(fname);
	gDirectory->cd("dir-1");
	gDirectory->cd("dir-11");
	gDirectory->cd("dir-111");

	auto t = (TTree*)gDirectory->Get(tree);
	if (!t) {
		std::cerr << "could not fetch TTree [" << tree << "] from file [" << fname << "]\n";
		exit(1);
	}
	auto player = dynamic_cast<TTreePlayer*>(t->GetPlayer());
	player->SetScanRedirect(kTRUE);
	player->SetScanFileName(oname);
	t->SetScanField(0);
	t->Scan("i32:f64");
}
`
	ofile := filepath.Join(tmp, "tree-subdir.txt")
	out, err := rtests.RunCxxROOT("scan", []byte(code), fname, "ntup", ofile)
	if err != nil {
		t.Fatalf("could not run C++ ROOT: %+v\noutput:\n%s", err, out)
	}

	got, err := os.ReadFile(ofile)
	if err != nil {
		t.Fatalf("could not read C++ ROOT scan file %q: %+v\noutput:\n%s", ofile, err, out)
	}

	want := `************************************
*    Row   *       i32 *       f64 *
************************************
*        0 *         0 *         0 *
*        1 *         1 *         1 *
*        2 *         2 *         2 *
*        3 *         3 *         3 *
*        4 *         4 *         4 *
************************************
`
	if got, want := string(got), want; got != want {
		t.Fatalf("invalid ROOT scan:\ngot:\n%v\nwant:\n%v\noutput:\n%s", got, want, out)
	}

}

var sumBenchReadTreeF64 = 0.0

func BenchmarkReadTreeF64(b *testing.B) {
	tmp, err := os.MkdirTemp("", "groot-rtree-read-tree-f64-")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	const nevts = 10000

	fname := path.Join(tmp, "f64.root")
	func() {
		b.StopTimer()
		defer b.StartTimer()

		f, err := riofs.Create(fname, riofs.WithoutCompression())
		if err != nil {
			b.Fatal(err)
		}
		defer f.Close()

		var data struct {
			F64 float64
		}
		wvars := []WriteVar{
			{Name: "F64", Value: &data.F64},
		}
		tree, err := NewWriter(f, "tree", wvars, WithoutCompression())
		if err != nil {
			b.Fatal(err)
		}
		defer tree.Close()

		rnd := rand.New(rand.NewSource(1234))
		for range nevts {
			data.F64 = rnd.Float64() * 10

			_, err = tree.Write()
			if err != nil {
				b.Fatal(err)
			}
		}

		err = tree.Close()
		if err != nil {
			b.Fatal(err)
		}

		err = f.Close()
		if err != nil {
			b.Fatal(err)
		}
	}()

	b.StopTimer()
	f, err := riofs.Open(fname)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		b.Fatal(err)
	}

	tree := o.(Tree)

	var data struct {
		F64 float64
	}

	rvars := ReadVarsFromStruct(&data)
	r, err := NewReader(tree, rvars)
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()

	b.StartTimer()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r.r.reset()
		b.StartTimer()

		err = r.Read(func(RCtx) error {
			sumBenchReadTreeF64 += data.F64
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

var sumBenchReadTreeSliF64 = 0

func BenchmarkReadTreeSliF64(b *testing.B) {
	tmp, err := os.MkdirTemp("", "groot-rtree-read-tree-sli-f64s-")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	const nevts = 1000

	for _, sz := range []int{0, 1, 2, 4, 8, 16, 64, 128, 512, 1024, 1024 * 1024} {
		fname := path.Join(tmp, fmt.Sprintf("f64s-%d.root", sz))
		func() {
			b.StopTimer()
			defer b.StartTimer()

			f, err := riofs.Create(fname, riofs.WithoutCompression())
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()

			var data struct {
				N   int32
				Sli []float64
			}
			wvars := []WriteVar{
				{Name: "N", Value: &data.N},
				{Name: "Sli", Value: &data.Sli, Count: "N"},
			}
			tree, err := NewWriter(f, "tree", wvars, WithoutCompression())
			if err != nil {
				b.Fatal(err)
			}
			defer tree.Close()

			rnd := rand.New(rand.NewSource(1234))
			for range nevts {
				data.N = int32(rnd.Float64() * 100)
				data.Sli = make([]float64, int(data.N))
				for j := range data.Sli {
					data.Sli[j] = rnd.Float64() * 10
				}

				_, err = tree.Write()
				if err != nil {
					b.Fatal(err)
				}
			}

			err = tree.Close()
			if err != nil {
				b.Fatal(err)
			}

			err = f.Close()
			if err != nil {
				b.Fatal(err)
			}
		}()

		b.Run(fmt.Sprintf("%d", sz), func(b *testing.B) {
			b.StopTimer()
			f, err := riofs.Open(fname)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()

			o, err := f.Get("tree")
			if err != nil {
				b.Fatal(err)
			}

			tree := o.(Tree)

			var data struct {
				N   int32
				Sli []float64
			}

			rvars := ReadVarsFromStruct(&data)
			r, err := NewReader(tree, rvars)
			if err != nil {
				b.Fatal(err)
			}
			defer r.Close()

			b.StartTimer()
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r.r.reset()
				data.N = 0
				data.Sli = data.Sli[:0]
				b.StartTimer()

				err = r.Read(func(RCtx) error {
					sumBenchReadTreeSliF64 += len(data.Sli)
					return nil
				})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type TNestedStruct1 struct {
	RunNbr int64            `groot:"runnbr"`
	EvtNbr int64            `groot:"evtnbr"`
	P3     TNestedStruct1P3 `groot:"p3"`
}

type TNestedStruct1P3 struct {
	Px float64 `groot:"px"`
	Py float64 `groot:"py"`
	Pz float64 `groot:"pz"`
}

type TNestedStruct2 struct {
	RunNbr int64            `groot:"runnbr"`
	EvtNbr int64            `groot:"evtnbr"`
	P3     TNestedStruct1P3 `groot:"p3"`
	F32s   []float32        `groot:"f32s"`
}

type TNestedStruct3 struct {
	RunNbr int64            `groot:"runnbr"`
	EvtNbr int64            `groot:"evtnbr"`
	P3     TNestedStruct1P3 `groot:"p3"`
	F32s   []float32        `groot:"f32s"`
}

type TNestedP2 struct {
	Px float64 `groot:"px"`
	Py float32 `groot:"py"`
}

type TNestedEvent1 struct {
	B   bool            `groot:"Bool"`
	Str string          `groot:"Str"`
	I8  int8            `groot:"I8"`
	I16 int16           `groot:"I16"`
	I32 int32           `groot:"I32"`
	I64 int64           `groot:"I64"`
	U8  uint8           `groot:"U8"`
	U16 uint16          `groot:"U16"`
	U32 uint32          `groot:"U32"`
	U64 uint64          `groot:"U64"`
	F32 float32         `groot:"F32"`
	F64 float64         `groot:"F64"`
	D16 root.Float16    `groot:"D16"`
	D32 root.Double32   `groot:"D32"`
	P2  TNestedP2       `groot:"P2"`
	Obj rbase.ObjString `groot:"Obj"`

	ArrBs [10]bool `groot:"ArrBs[10]"`
	//ArrStr [10]string `groot:"ArrStr[10]"`
	ArrI8  [10]int8            `groot:"ArrI8[10]"`
	ArrI16 [10]int16           `groot:"ArrI16[10]"`
	ArrI32 [10]int32           `groot:"ArrI32[10]"`
	ArrI64 [10]int64           `groot:"ArrI64[10]"`
	ArrU8  [10]uint8           `groot:"ArrU8[10]"`
	ArrU16 [10]uint16          `groot:"ArrU16[10]"`
	ArrU32 [10]uint32          `groot:"ArrU32[10]"`
	ArrU64 [10]uint64          `groot:"ArrU64[10]"`
	ArrF32 [10]float32         `groot:"ArrF32[10]"`
	ArrF64 [10]float64         `groot:"ArrF64[10]"`
	ArrD16 [10]root.Float16    `groot:"ArrD16[10]"`
	ArrD32 [10]root.Double32   `groot:"ArrD32[10]"`
	ArrP2  [10]TNestedP2       `groot:"ArrP2[10]"`
	ArrObj [10]rbase.ObjString `groot:"ArrObj[10]"`

	N     int32  `groot:"N"`
	SliBs []bool `groot:"SliBs[N]"`
	//	SliStr []string        `groot:"SliStr[N]"`
	SliI8  []int8          `groot:"SliI8[N]"`
	SliI16 []int16         `groot:"SliI16[N]"`
	SliI32 []int32         `groot:"SliI32[N]"`
	SliI64 []int64         `groot:"SliI64[N]"`
	SliU8  []uint8         `groot:"SliU8[N]"`
	SliU16 []uint16        `groot:"SliU16[N]"`
	SliU32 []uint32        `groot:"SliU32[N]"`
	SliU64 []uint64        `groot:"SliU64[N]"`
	SliF32 []float32       `groot:"SliF32[N]"`
	SliF64 []float64       `groot:"SliF64[N]"`
	SliD16 []root.Float16  `groot:"SliD16[N]"`
	SliD32 []root.Double32 `groot:"SliD32[N]"`
	//	SliP2  []TNestedP2     `groot:"SliP2[N]"` // FIXME(sbinet): var-len-array of non-builtins has extra bytes in front
	//	SliObj  []rbase.ObjString     `groot:"SliObj[N]"` // FIXME(sbinet): var-len-array of non-builtins has extra bytes in front

	StdVecBs  []bool            `groot:"StdVecBs"`
	StdVecStr []string          `groot:"StdVecStr"`
	StdVecI8  []int8            `groot:"StdVecI8"`
	StdVecI16 []int16           `groot:"StdVecI16"`
	StdVecI32 []int32           `groot:"StdVecI32"`
	StdVecI64 []int64           `groot:"StdVecI64"`
	StdVecU8  []uint8           `groot:"StdVecU8"`
	StdVecU16 []uint16          `groot:"StdVecU16"`
	StdVecU32 []uint32          `groot:"StdVecU32"`
	StdVecU64 []uint64          `groot:"StdVecU64"`
	StdVecF32 []float32         `groot:"StdVecF32"`
	StdVecF64 []float64         `groot:"StdVecF64"`
	StdVecD16 []root.Float16    `groot:"StdVecD16"`
	StdVecD32 []root.Double32   `groot:"StdVecD32"`
	StdVecP2  []TNestedP2       `groot:"StdVecP2"`
	StdVecObj []rbase.ObjString `groot:"StdVecObj"`

	StdVecVecF64 [][]float64   `groot:"StdVecVecF64"`
	StdVecVecStr [][]string    `groot:"StdVecVecStr"`
	StdVecVecP2  [][]TNestedP2 `groot:"StdVecVecP2"`
}

func (TNestedEvent1) want(i int64) (data TNestedEvent1) {
	data.B = i%2 == 0
	data.Str = fmt.Sprintf("str-%03d", i)
	data.I8 = int8(-i)
	data.I16 = int16(-i)
	data.I32 = int32(-i)
	data.I64 = int64(-i)
	data.U8 = uint8(i)
	data.U16 = uint16(i)
	data.U32 = uint32(i)
	data.U64 = uint64(i)
	data.F32 = float32(i)
	data.F64 = float64(i)
	data.D16 = root.Float16(i)
	data.D32 = root.Double32(i)
	data.P2 = TNestedP2{Px: float64(i), Py: float32(i)}
	data.Obj = *rbase.NewObjString(fmt.Sprintf("obj-%d", i))

	for ii := range data.ArrI32 {
		data.ArrBs[ii] = ii == int(i)
		//data.ArrStr[ii] = fmt.Sprintf("arr-%03d", i)
		data.ArrI8[ii] = int8(-i)
		data.ArrI16[ii] = int16(-i)
		data.ArrI32[ii] = int32(-i)
		data.ArrI64[ii] = int64(-i)
		data.ArrU8[ii] = uint8(i)
		data.ArrU16[ii] = uint16(i)
		data.ArrU32[ii] = uint32(i)
		data.ArrU64[ii] = uint64(i)
		data.ArrF32[ii] = float32(i)
		data.ArrF64[ii] = float64(i)
		data.ArrD16[ii] = root.Float16(i)
		data.ArrD32[ii] = root.Double32(i)
		data.ArrP2[ii] = TNestedP2{
			Px: float64(i),
			Py: float32(i),
		}
		data.ArrObj[ii] = *rbase.NewObjString(fmt.Sprintf("obj-%d", i))
	}
	data.N = int32(i) % 10

	switch data.N {
	case 0:
		data.SliBs = nil
		//		data.SliStr = nil
		data.SliI8 = nil
		data.SliI16 = nil
		data.SliI32 = nil
		data.SliI64 = nil
		data.SliU8 = nil
		data.SliU16 = nil
		data.SliU32 = nil
		data.SliU64 = nil
		data.SliF32 = nil
		data.SliF64 = nil
		data.SliD16 = nil
		data.SliD32 = nil
		//		data.SliP2 = nil
		//		data.SliObj = nil
	default:
		data.SliBs = make([]bool, int(data.N))
		//		data.SliStr = make([]string, int(data.N))
		data.SliI8 = make([]int8, int(data.N))
		data.SliI16 = make([]int16, int(data.N))
		data.SliI32 = make([]int32, int(data.N))
		data.SliI64 = make([]int64, int(data.N))
		data.SliU8 = make([]uint8, int(data.N))
		data.SliU16 = make([]uint16, int(data.N))
		data.SliU32 = make([]uint32, int(data.N))
		data.SliU64 = make([]uint64, int(data.N))
		data.SliF32 = make([]float32, int(data.N))
		data.SliF64 = make([]float64, int(data.N))
		data.SliD16 = make([]root.Float16, int(data.N))
		data.SliD32 = make([]root.Double32, int(data.N))
		//		data.SliP2 = make([]TNestedP2, int(data.N))
		//		data.SliObj = make([]rbase.ObjString, int(data.N))
	}
	for ii := range int(data.N) {
		data.SliBs[ii] = (ii + 1) == int(i)
		//		data.SliStr[ii] = fmt.Sprintf("sli-%03d", i)
		data.SliI8[ii] = int8(-i)
		data.SliI16[ii] = int16(-i)
		data.SliI32[ii] = int32(-i)
		data.SliI64[ii] = int64(-i)
		data.SliU8[ii] = uint8(i)
		data.SliU16[ii] = uint16(i)
		data.SliU32[ii] = uint32(i)
		data.SliU64[ii] = uint64(i)
		data.SliF32[ii] = float32(i)
		data.SliF64[ii] = float64(i)
		data.SliD16[ii] = root.Float16(i)
		data.SliD32[ii] = root.Double32(i)
		//		data.SliP2[ii] = TNestedP2{
		//			Px: float64(i),
		//			Py: float32(i),
		//		}
		//		data.SliObj[ii] = *rbase.NewObjString(fmt.Sprintf("obj-%03d", i))
	}

	switch data.N {
	case 0:
		data.StdVecBs = nil
		data.StdVecStr = nil
		data.StdVecI8 = nil
		data.StdVecI16 = nil
		data.StdVecI32 = nil
		data.StdVecI64 = nil
		data.StdVecU8 = nil
		data.StdVecU16 = nil
		data.StdVecU32 = nil
		data.StdVecU64 = nil
		data.StdVecF32 = nil
		data.StdVecF64 = nil
		data.StdVecD16 = nil
		data.StdVecD32 = nil
		data.StdVecP2 = nil
		data.StdVecObj = nil
	default:
		data.StdVecBs = make([]bool, int(data.N))
		data.StdVecStr = make([]string, int(data.N))
		data.StdVecI8 = make([]int8, int(data.N))
		data.StdVecI16 = make([]int16, int(data.N))
		data.StdVecI32 = make([]int32, int(data.N))
		data.StdVecI64 = make([]int64, int(data.N))
		data.StdVecU8 = make([]uint8, int(data.N))
		data.StdVecU16 = make([]uint16, int(data.N))
		data.StdVecU32 = make([]uint32, int(data.N))
		data.StdVecU64 = make([]uint64, int(data.N))
		data.StdVecF32 = make([]float32, int(data.N))
		data.StdVecF64 = make([]float64, int(data.N))
		data.StdVecD16 = make([]root.Float16, int(data.N))
		data.StdVecD32 = make([]root.Double32, int(data.N))
		data.StdVecP2 = make([]TNestedP2, int(data.N))
		data.StdVecObj = make([]rbase.ObjString, int(data.N))
	}
	for ii := range int(data.N) {
		data.StdVecBs[ii] = (ii + 1) == int(i)
		data.StdVecStr[ii] = fmt.Sprintf("std-%03d", i)
		data.StdVecI8[ii] = int8(-i)
		data.StdVecI16[ii] = int16(-i)
		data.StdVecI32[ii] = int32(-i)
		data.StdVecI64[ii] = int64(-i)
		data.StdVecU8[ii] = uint8(i)
		data.StdVecU16[ii] = uint16(i)
		data.StdVecU32[ii] = uint32(i)
		data.StdVecU64[ii] = uint64(i)
		data.StdVecF32[ii] = float32(i)
		data.StdVecF64[ii] = float64(i)
		data.StdVecD16[ii] = root.Float16(i)
		data.StdVecD32[ii] = root.Double32(i)
		data.StdVecP2[ii] = TNestedP2{
			Px: float64(i),
			Py: float32(i),
		}
		data.StdVecObj[ii] = *rbase.NewObjString(fmt.Sprintf("obj-%03d", i))
	}

	switch data.N {
	case 0:
		data.StdVecVecF64 = nil
		data.StdVecVecStr = nil
		data.StdVecVecP2 = nil
	default:
		data.StdVecVecF64 = make([][]float64, data.N)
		data.StdVecVecStr = make([][]string, data.N)
		data.StdVecVecP2 = make([][]TNestedP2, data.N)
	}
	for ii := range int(data.N) {
		data.StdVecVecF64[ii] = []float64{
			float64(ii),
			float64(ii + 1),
			float64(ii + 2),
			float64(ii + 3),
		}
		data.StdVecVecStr[ii] = []string{
			fmt.Sprintf("vec-%03d", i),
			fmt.Sprintf("vec-%03d", i+1),
			fmt.Sprintf("vec-%03d", i+2),
			fmt.Sprintf("vec-%03d", i+3),
		}
		data.StdVecVecP2[ii] = []TNestedP2{
			{Px: float64(i), Py: float32(i)},
			{Px: float64(i + 1), Py: float32(i + 1)},
			{Px: float64(i + 2), Py: float32(i + 2)},
			{Px: float64(i + 3), Py: float32(i + 3)},
		}
	}

	return data
}
