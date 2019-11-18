// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"io"
	"reflect"

	"golang.org/x/xerrors"
)

type Basket struct {
	key Key

	vers    uint16
	bufsize int // length in bytes
	nevsize int // length in int_t or fixed length of each entry
	nevbuf  int // number of entries in basket
	last    int // pointer to last used byte in basket
	flag    byte

	header  bool    // true when only the basket header must be read/written
	displ   []int32 // displacement of entries in key.buffer
	offsets []int32 // offset of entries in key.buffer

	rbuf *RBuffer
}

func (b *Basket) Name() string {
	return b.key.Name()
}

func (b *Basket) Title() string {
	return b.key.Title()
}

func (b *Basket) Class() string {
	return "TBasket"
}

func (b *Basket) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if err := b.key.UnmarshalROOT(r); err != nil {
		return err
	}

	if b.Class() != "TBasket" {
		return xerrors.Errorf("rootio.Basket: Key is not a Basket")
	}

	b.vers = r.ReadU16()
	b.bufsize = int(r.ReadI32())
	b.nevsize = int(r.ReadI32())

	if b.nevsize < 0 {
		r.err = xerrors.Errorf("rootio.Basket: incorrect event buffer size [%v]", b.nevsize)
		b.nevsize = 0
		return r.err
	}

	b.nevbuf = int(r.ReadI32())
	b.last = int(r.ReadI32())
	b.flag = r.ReadU8()

	if b.last > b.bufsize {
		b.bufsize = b.last
	}

	if b.flag == 0 {
		return r.err
	}

	if b.flag%10 != 2 {
		if b.nevbuf > 0 {
			n := int(r.ReadI32())
			b.offsets = r.ReadFastArrayI32(n)
			if 20 < b.flag && b.flag < 40 {
				const displacement uint32 = 0xFF000000
				for i, v := range b.offsets {
					b.offsets[i] = int32(uint32(v) &^ displacement)
				}
			}
		}
		if b.flag > 40 {
			n := int(r.ReadI32())
			b.offsets = r.ReadFastArrayI32(n)
		}
	}

	if b.flag == 1 || b.flag > 10 {
		// reading raw data
		var sz = int32(b.last)
		if b.vers <= 1 {
			sz = r.ReadI32()
		}
		buf := make([]byte, int(sz))
		_, err := io.ReadFull(r.r, buf)
		if err != nil {
			r.err = err
			return r.err
		}
		b.key.buf = buf
	}

	return r.err
}

func (b *Basket) loadEntry(entry int64) error {
	var err error
	var offset = int64(b.key.keylen)
	if len(b.offsets) > 0 {
		offset = int64(b.offsets[int(entry)])
	}
	pos := entry*int64(b.nevsize) + offset
	err = b.rbuf.setPos(pos)
	return err
}

func (b *Basket) readLeaf(entry int64, leaf Leaf) error {
	var offset int64
	if len(b.offsets) == 0 {
		offset = entry*int64(b.nevsize) + int64(leaf.Offset()) + int64(b.key.keylen)
	} else {
		offset = int64(b.offsets[int(entry)]) + int64(leaf.Offset())
	}
	err := b.rbuf.setPos(offset)
	if err != nil {
		return err
	}
	return leaf.readBasket(b.rbuf)
}

func init() {
	f := func() reflect.Value {
		o := &Basket{}
		return reflect.ValueOf(o)
	}
	Factory.add("TBasket", f)
	Factory.add("*rootio.Basket", f)
}

var _ Object = (*Basket)(nil)
var _ Named = (*Basket)(nil)
var _ ROOTUnmarshaler = (*Basket)(nil)
