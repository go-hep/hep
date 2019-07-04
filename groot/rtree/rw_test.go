// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
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
