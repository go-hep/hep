// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io"
	"runtime"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
)

// raBasket is a read-ahead basket reader.
type raBasket struct {
	f     *riofs.File
	spans []rspan

	ready  chan *rbasket // baskets ready to be handed to the reader
	reuse  chan *rbasket // baskets to reuse for input reading
	exit   chan struct{} // closes when finished
	n      int           // number of in-flight baskets
	err    error
	cur    *rbasket      // current buffer being served
	closed chan struct{} // channel is closed when the async reader shuts down

	name string
}

func newRBasket(b Branch, n int) *raBasket {
	if n < 0 {
		n = runtime.NumCPU() + 1
	}
	if n == 0 {
		n = 1
	}
	base := asBranch(b)
	ra := &raBasket{
		f:      b.getTree().f,
		spans:  make([]rspan, len(base.basketSeek)),
		ready:  make(chan *rbasket, n),
		reuse:  make(chan *rbasket, n),
		exit:   make(chan struct{}),
		n:      n,
		closed: make(chan struct{}),
		name:   b.Name(),
	}

	for i, seek := range base.basketSeek {
		ra.spans[i] = rspan{
			pos: seek,
			sz:  base.basketBytes[i],
			beg: base.basketEntry[i],
			end: base.basketEntry[i+1],
		}
	}

	for i := 0; i < n; i++ {
		ra.reuse <- &rbasket{}
	}

	go ra.run(base.entryOffsetLen)

	return ra
}

func (ra *raBasket) run(eoff int) {
	defer close(ra.closed)
	defer close(ra.ready)
	for _, span := range ra.spans {
		select {
		case b := <-ra.reuse:
			err := b.inflate(ra.name, span, eoff, ra.f)
			if err != nil {
				ra.err = err
				return
			}
			ra.ready <- b
		case <-ra.exit:
			return
		}
	}
}

func (ra *raBasket) read() (*rbasket, error) {
	if ra.cur != nil {
		ra.cur.reset()
		ra.reuse <- ra.cur
		ra.cur = nil
	}
	b, ok := <-ra.ready
	if !ok {
		return nil, io.EOF
		//		if ra.err == nil {
		//			ra.err = errors.New("rtree: read-read basket after close")
		//		}
		//		return nil, ra.err
	}
	ra.cur = b

	return ra.cur, nil
}

func (a *raBasket) close() {
	select {
	case <-a.closed:
	case a.exit <- struct{}{}:
		<-a.closed
	}
	return
}

type rspan struct {
	beg int64 // first entry
	end int64 // last entry

	pos int64 // span location on-disk
	sz  int32 // basket size
}

type rbasket struct {
	id    int    // basket number
	entry int64  // current entry number
	span  rspan  // basket entry span
	bk    Basket // current basket
	buf   []byte
}

func (rbk *rbasket) reset() {
	rbk.id = 0
	rbk.entry = 0
	rbk.span = rspan{}
	//	rbk.bk = Basket{}
}

func (rbk *rbasket) inflate(name string, span rspan, eoff int, f *riofs.File) error {
	var (
		bufsz = span.sz
		seek  = span.pos
	)

	rbk.span = span
	//log.Printf("inflate[%s]: [%d, %d)", name, span.pos, span.pos+int64(span.sz))

	var (
		sictx  = f
		err    error
		keylen uint32
	)

	switch {
	case bufsz == 0: // FIXME(sbinet): from trial and error. check this is ok for all cases

		rbk.bk.key.SetFile(f)
		rbk.buf = rbytes.ResizeU8(rbk.buf, int(rbk.bk.key.ObjLen()))
		_, err = rbk.bk.key.Load(rbk.buf)
		if err != nil {
			return err
		}
		rbk.bk.rbuf = rbytes.NewRBuffer(rbk.buf, nil, keylen, sictx)

	default:
		rbk.buf = rbytes.ResizeU8(rbk.buf, int(bufsz))
		_, err = f.ReadAt(rbk.buf, seek)
		if err != nil {
			return fmt.Errorf("rtree: could not read basket buffer from file: %w", err)
		}

		err = rbk.bk.UnmarshalROOT(rbytes.NewRBuffer(rbk.buf, nil, 0, sictx))
		if err != nil {
			return fmt.Errorf("rtree: could not unmarshal basket buffer from file: %w", err)
		}
		rbk.bk.key.SetFile(f)

		rbk.buf = rbytes.ResizeU8(rbk.buf, int(rbk.bk.key.ObjLen()))
		_, err = rbk.bk.key.Load(rbk.buf)
		if err != nil {
			return err
		}
		keylen = uint32(rbk.bk.key.KeyLen())
		rbk.bk.rbuf = rbytes.NewRBuffer(rbk.buf, nil, keylen, sictx)

		if eoff > 0 {
			last := int64(rbk.bk.last)
			err = rbk.bk.rbuf.SetPos(last)
			if err != nil {
				return err
			}
			n := int(rbk.bk.rbuf.ReadI32())
			rbk.bk.offsets = rbytes.ResizeI32(rbk.bk.offsets, n)
			rbk.bk.rbuf.ReadArrayI32(rbk.bk.offsets)
			if err := rbk.bk.rbuf.Err(); err != nil {
				return err
			}
		}
	}

	return nil
}
