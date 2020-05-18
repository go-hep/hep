// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
)

type rbasket struct {
	id   int    // basket number
	span rspan  // basket entry span
	bk   Basket // current basket
	buf  []byte
}

func (rbk *rbasket) reset() {
	rbk.id = 0
	rbk.span = rspan{}
	//	rbk.bk = Basket{}
}

func (rbk *rbasket) inflate(name string, id int, span rspan, eoff int, f *riofs.File) error {
	var (
		bufsz = span.sz
		seek  = span.pos
	)

	rbk.id = id
	rbk.span = span

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
		rbk.bk.rbuf = rbk.bk.rbuf.Reset(rbk.buf, nil, keylen, sictx)

	default:
		rbk.buf = rbytes.ResizeU8(rbk.buf, int(bufsz))
		_, err = f.ReadAt(rbk.buf, seek)
		if err != nil {
			return fmt.Errorf("rtree: could not read basket buffer from file: %w", err)
		}

		rbk.bk.rbuf = rbk.bk.rbuf.Reset(rbk.buf, nil, 0, sictx)
		err = rbk.bk.UnmarshalROOT(rbk.bk.rbuf)
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
		rbk.bk.rbuf = rbk.bk.rbuf.Reset(rbk.buf, nil, keylen, sictx)

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
