// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

const (
	kGenerateOffsetMap = 0
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

func newBasketFrom(t Tree, b Branch, cycle int16, bufsize, eoffsetLen int) Basket {
	var (
		f     = FileOf(t)
		name  = b.Name()
		title = t.Name()
		class = "TBasket"
	)

	bkt := Basket{
		key:     riofs.NewKeyForBasketInternal(f, name, title, class, cycle),
		bufsize: bufsize,
		nevsize: eoffsetLen,
		wbuf:    rbytes.NewWBuffer(nil, nil, 0, nil),
		header:  true, // FIXME(sbinet): ROOT default is "false"
		branch:  b,
	}

	bkt.offsets = rbytes.ResizeI32(bkt.offsets, bkt.nevsize)
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
			w.SetErr(fmt.Errorf("rtree: could not marshal iobits basket: %w", err))
			return n, w.Err()
		}
	default:
		w.WriteI32(int32(b.nevsize))
	}
	w.WriteI32(int32(b.nevbuf))
	w.WriteI32(int32(b.last))

	mustGenOffsets := (len(b.offsets) > 0 && b.nevbuf > 0 &&
		(b.iobits&kGenerateOffsetMap != 0) &&
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
		if b.wbuf != nil && b.wbuf.Len() > 0 {
			raw := b.wbuf.Bytes()
			n := b.last
			if len(raw) < n {
				n = len(raw)
			}
			_, err := w.Write(raw[:n])
			if err != nil {
				return int(w.Pos() - beg), err
			}
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
		return fmt.Errorf("rtree: Key is not a Basket")
	}

	vers := r.ReadI16()
	if vers > rvers.Basket {
		return fmt.Errorf("rtree: unknown Basket version (got = %d > %d)", vers, rvers.Basket)
	}

	b.bufsize = int(r.ReadI32())
	b.nevsize = int(r.ReadI32())

	if b.nevsize < 0 {
		b.nevsize = -b.nevsize
		if err := b.iobits.UnmarshalROOT(r); err != nil {
			r.SetErr(fmt.Errorf("rtree: could not read basket I/O bits: %w", err))
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
			b.offsets = rbytes.ResizeI32(b.offsets, n)
			r.ReadArrayI32(b.offsets)
			if 20 < flag && flag < 40 {
				for i, v := range b.offsets {
					b.offsets[i] = int32(uint32(v) &^ rbytes.DisplacementMask)
				}
			}
		}
		if flag > 40 {
			n := int(r.ReadI32())
			b.displ = rbytes.ResizeI32(b.displ, n)
			r.ReadArrayI32(b.displ)
		}
	case mustGenOffsets:
		b.offsets = nil
		if flag <= 40 {
			panic(fmt.Errorf("rtree: invalid basket[%s] state (flag=%v <= 40)", b.Name(), flag))
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

func (b *Basket) loadLeaf(entry int64, leaf Leaf) error {
	var offset int64
	if len(b.offsets) == 0 {
		offset = entry*int64(b.nevsize) + int64(leaf.Offset()) + int64(b.key.KeyLen())
	} else {
		offset = int64(b.offsets[int(entry)]) + int64(leaf.Offset())
	}
	b.rbuf.SetPos(offset)
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

func (b *Basket) update(offset int64) {
	offset += int64(b.key.KeyLen())
	if len(b.offsets) > 0 {
		if b.nevbuf+1 >= b.nevsize {
			nevsize := 10
			if nevsize < 2*b.nevsize {
				nevsize = 2 * b.nevsize
			}
			b.nevsize = nevsize
			delta := len(b.offsets) - nevsize
			if delta < 0 {
				delta = -delta
			}
			b.offsets = append(b.offsets, make([]int32, delta)...)
		}
		b.offsets[b.nevbuf] = int32(offset)
	}
	b.nevbuf++
}

func (b *Basket) grow(n int) {
	b.nevsize = n
	delta := len(b.offsets) - n
	if delta < 0 {
		delta = -delta
	}
	b.offsets = append(b.offsets, make([]int32, delta)...)
}

func (b *Basket) writeFile(f *riofs.File) (totBytes int64, zipBytes int64, err error) {
	header := b.header
	b.header = true
	defer func() {
		b.header = header
	}()

	// we need to handle the case for a basket being created
	// while the file was small, and *then* being flushed while
	// the file is big.
	// ie: the TKey structure switched to 64b offsets, and add an
	// extra 8bytes.
	// we need to propagate to the 'offsets' and 'last' fields.
	adjust := !(b.key.RVersion() > 1000) && f.IsBigFile()

	b.last = int(int64(b.key.KeyLen()) + b.wbuf.Len())
	if b.offsets != nil {
		if adjust {
			for i, v := range b.offsets {
				b.offsets[i] = v + 8
			}
		}
		b.wbuf.WriteI32(int32(b.nevbuf + 1))
		b.wbuf.WriteFastArrayI32(b.offsets[:b.nevbuf])
		b.wbuf.WriteI32(0)
	}
	b.key, err = riofs.NewKey(nil, b.key.Name(), b.key.Title(), b.Class(), int16(b.key.Cycle()), b.wbuf.Bytes(), f)
	if err != nil {
		return 0, 0, fmt.Errorf("rtree: could not create basket-key: %w", err)
	}
	if adjust {
		b.last += 8
	}

	nbytes := b.key.KeyLen() + b.key.ObjLen()
	buf := rbytes.NewWBuffer(make([]byte, nbytes), nil, uint32(b.key.KeyLen()), f)
	_, err = b.MarshalROOT(buf)
	if err != nil {
		return 0, 0, err
	}

	n, err := f.WriteAt(buf.Bytes(), b.key.SeekKey())
	if err != nil {
		return int64(n), int64(n), err
	}
	nn, err := f.WriteAt(b.key.Buffer(), b.key.SeekKey()+int64(b.key.KeyLen()))
	n += nn
	if err != nil {
		return int64(n), int64(n), err
	}
	b.wbuf = nil
	b.key.SetBuffer(nil)

	return int64(nbytes), int64(b.key.Nbytes()), nil
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
	_ rbytes.Marshaler   = (*Basket)(nil)
	_ rbytes.Unmarshaler = (*Basket)(nil)
)
