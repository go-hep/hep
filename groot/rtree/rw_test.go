// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
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
