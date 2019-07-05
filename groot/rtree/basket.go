// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"io"
	"reflect"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type Basket struct {
	key riofs.Key

	bufsize int // length in bytes
	nevsize int // length in int_t or fixed length of each entry
	nevbuf  int // number of entries in basket
	last    int // pointer to last used byte in basket

	header  bool        // true when only the basket header must be read/written
	iobits  tioFeatures // IO feature flags
	displ   []int32     // displacement of entries in key.buffer
	offsets []int32     // offset of entries in key.buffer

	rbuf *rbytes.RBuffer
	wbuf *rbytes.WBuffer

	branch Branch // basket support branch
}

func newBasketFrom(t Tree, b Branch) Basket {
	var dir riofs.Directory
	switch b := b.(type) {
	case *tbranch:
		dir = b.dir
	case *tbranchElement:
		dir = b.tbranch.dir
	default:
		panic(errors.Errorf("rtree: unknown Branch type %T", b))
	}

	var (
		name  = b.Name()
		title = t.Name()
		class = "TBasket"
	)

	bkt := Basket{
		key:  riofs.KeyFromDir(dir, name, title, class),
		wbuf: rbytes.NewWBuffer(nil, nil, 0, nil),
	}

	return bkt
}

func (b *Basket) Name() string {
	return b.key.Name()
}

func (b *Basket) Title() string {
	return b.key.Title()
}

func (*Basket) RVersion() int16 {
	return rvers.Basket
}

func (b *Basket) Class() string {
	return "TBasket"
}

func (b *Basket) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	beg := w.Pos()

	if n, err := b.key.MarshalROOT(w); err != nil {
		return n, err
	}

	w.WriteI16(b.RVersion())
	w.WriteI32(int32(b.bufsize))
	switch {
	case b.iobits != 0:
		w.WriteI32(int32(-b.nevsize))
		if n, err := b.iobits.MarshalROOT(w); err != nil {
			w.SetErr(errors.Wrapf(err, "rtree: could not marshal iobits basket"))
			return n, w.Err()
		}
	default:
		w.WriteI32(int32(b.nevsize))
	}
	w.WriteI32(int32(b.nevbuf))
	w.WriteI32(int32(b.last))

	mustGenOffsets := (len(b.offsets) > 0 && b.nevbuf > 0 &&
		b.canGenerateOffsetArray())

	if mustGenOffsets && len(b.displ) > 0 {
		panic("rtree: impossible basket serialization case")
	}

	var flag byte
	switch {
	case b.header:
		if mustGenOffsets {
			flag = 80
		}
		w.WriteU8(flag)

	default:
		if b.nevbuf > 0 {
			b.computeEntryOffsets()
		}
		flag = 1
		if b.nevbuf <= 0 || len(b.offsets) == 0 {
			flag = 2
		}
		if b.wbuf != nil {
			flag += 10
		}

		if len(b.displ) > 0 {
			flag += 40
		}

		if mustGenOffsets {
			flag += 80
		}
		w.WriteU8(flag)

		if !mustGenOffsets && len(b.offsets) > 0 && b.nevbuf > 0 {
			w.WriteI32(int32(b.nevbuf))
			w.WriteFastArrayI32(b.offsets[:b.nevbuf])
			if len(b.displ) > 0 {
				w.WriteI32(int32(b.nevbuf))
				w.WriteFastArrayI32(b.displ)
			}
		}
		if b.wbuf != nil {
			raw := b.wbuf.Bytes()
			w.Write(raw[:b.last])
		}
	}

	n := w.Pos() - beg
	return int(n), w.Err()
}

func (b *Basket) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if err := b.key.UnmarshalROOT(r); err != nil {
		return err
	}

	if b.Class() != "TBasket" {
		return errors.Errorf("rtree: Key is not a Basket")
	}

	vers := r.ReadI16()
	if vers > rvers.Basket {
		return errors.Errorf("rtree: unknown Basket version (got = %d > %d)", vers, rvers.Basket)
	}

	b.bufsize = int(r.ReadI32())
	b.nevsize = int(r.ReadI32())

	if b.nevsize < 0 {
		b.nevsize = -b.nevsize
		if err := b.iobits.UnmarshalROOT(r); err != nil {
			r.SetErr(errors.Wrapf(err, "rtree: could not read basket I/O bits"))
			return r.Err()
		}
	}

	b.nevbuf = int(r.ReadI32())
	b.last = int(r.ReadI32())

	flag := r.ReadU8()

	if b.last > b.bufsize {
		b.bufsize = b.last
	}

	mustGenOffsets := false
	if flag >= 80 {
		mustGenOffsets = true
		flag -= 80
	}

	switch {
	case !mustGenOffsets && flag != 0 && (flag%10 != 2):
		if b.nevbuf > 0 {
			n := int(r.ReadI32())
			b.offsets = r.ReadFastArrayI32(n)
			if 20 < flag && flag < 40 {
				const displacement uint32 = 0xFF000000
				for i, v := range b.offsets {
					b.offsets[i] = int32(uint32(v) &^ displacement)
				}
			}
		}
		if flag > 40 {
			n := int(r.ReadI32())
			b.displ = r.ReadFastArrayI32(n)
		}
	case mustGenOffsets:
		b.offsets = nil
		if flag <= 40 {
			panic(errors.Errorf("rtree: invalid basket state (flag=%v <= 40)", flag))
		}
	}

	if flag == 1 || flag > 10 {
		// reading raw data
		var sz = int32(b.last)
		if vers <= 1 {
			sz = r.ReadI32()
		}
		buf := make([]byte, int(sz))
		_, err := io.ReadFull(r, buf)
		if err != nil {
			r.SetErr(err)
			return r.Err()
		}
		b.key.SetBuffer(buf)
	}

	return r.Err()
}

func (b *Basket) loadEntry(entry int64) error {
	var err error
	var offset = int64(b.key.KeyLen())
	if len(b.offsets) > 0 {
		offset = int64(b.offsets[int(entry)])
	}
	pos := entry*int64(b.nevsize) + offset
	err = b.rbuf.SetPos(pos)
	return err
}

func (b *Basket) readLeaf(entry int64, leaf Leaf) error {
	var offset int64
	if len(b.offsets) == 0 {
		offset = entry*int64(b.nevsize) + int64(leaf.Offset()) + int64(b.key.KeyLen())
	} else {
		offset = int64(b.offsets[int(entry)]) + int64(leaf.Offset())
	}
	err := b.rbuf.SetPos(offset)
	if err != nil {
		return err
	}
	return leaf.readFromBuffer(b.rbuf)
}

func (b *Basket) computeEntryOffsets() {
	if b.offsets != nil {
		return
	}

	if b.branch == nil {
		panic("rtree: basket with no associated branch")
	}

	if len(b.branch.Leaves()) != 1 {
		panic("rtree: basket's associated branch contains multiple leaves")
	}

	leaf := b.branch.Leaves()[0]
	b.offsets = leaf.computeOffsetArray(int(b.key.KeyLen()), b.nevbuf)
}

func (b *Basket) canGenerateOffsetArray() bool {
	if len(b.branch.Leaves()) != 1 {
		return false
	}
	leaf := b.branch.Leaves()[0]
	return leaf.canGenerateOffsetArray()
}

func init() {
	f := func() reflect.Value {
		o := &Basket{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TBasket", f)
}

var (
	_ root.Object        = (*Basket)(nil)
	_ root.Named         = (*Basket)(nil)
	_ rbytes.Unmarshaler = (*Basket)(nil)
)
