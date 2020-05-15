// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

type rbranch struct {
	b      Branch
	rb     *bkreader
	cur    *rbasket
	leaves []rleaf
}

func newRBranch(b Branch, n int, beg, end int64, leaves []rleaf, rctx rleafCtx) rbranch {
	rb := rbranch{
		b:      b,
		rb:     newBkReader(b, n, beg, end),
		leaves: leaves,
	}
	return rb
}

func (rb *rbranch) start() error {
	var err error
	rb.cur, err = rb.rb.read()
	return err
}

func (rb *rbranch) stop() error {
	if rb.cur == nil {
		return nil
	}
	rb.rb.close()
	return nil
}

func (rb *rbranch) reset() {
	rb.rb.close()
	rb.rb = newBkReader(rb.b, rb.rb.n, rb.rb.beg, rb.rb.end)
}

func (rb *rbranch) read(i int64) error {
	var err error
	if i >= rb.cur.span.end {
		rb.cur, err = rb.rb.read()
		if err != nil {
			return err
		}
	}

	j := i - rb.cur.span.beg
	switch len(rb.leaves) {
	case 1:
		err = rb.cur.bk.loadRLeaf(j, rb.leaves[0])
		if err != nil {
			return err
		}

	default:
		for _, leaf := range rb.leaves {
			err = rb.cur.bk.loadRLeaf(j, leaf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func asBranch(b Branch) *tbranch {
	switch b := b.(type) {
	case *tbranch:
		return b
	case *tbranchElement:
		return &b.tbranch
	}
	panic("impossible")
}
