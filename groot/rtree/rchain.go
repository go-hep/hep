// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
)

type rchain struct {
	ch *tchain

	rvs  []ReadVar
	nrab int
	beg  int64
	end  int64

	ibeg int // first tree to process
	iend int // last-1 tree to process
}

var (
	_ reader = (*rchain)(nil)
)

func newRChain(ch *tchain, rvars []ReadVar, n int, beg, end int64) *rchain {
	r := &rchain{
		ch:   ch,
		rvs:  rvars,
		nrab: n,
		beg:  beg,
		end:  end,
	}

	tbeg, tend := r.findTrees(beg, end)
	if tbeg < 0 || tend < 0 {
		panic(fmt.Errorf(
			"rtree: could not find matching trees in chain for [%d, %d) within [%d, %d)",
			beg, end, 0, ch.Entries(),
		))
	}
	r.ibeg = tbeg
	r.iend = tend

	r.loadRVars()

	return r
}

func (r *rchain) Close() error {
	return nil
}

func (r *rchain) rvars() []ReadVar { return r.rvs }

func (r *rchain) loadRVars() {
	if len(r.ch.trees) == 0 {
		return
	}

	rr := newReader(r.ch.trees[0], r.rvs, r.nrab, 0, 1)
	defer rr.Close()
	r.rvs = rr.rvars()
}

func (r *rchain) run(off, beg, end int64, f func(RCtx) error) error {
	defer r.Close()

	trees := r.ch.trees[r.ibeg:r.iend]
	if len(trees) == 0 {
		return nil
	}

	for i := r.ibeg; i < r.iend; i++ {
		var (
			eoff = r.ch.offs[i]
			tots = r.ch.tots[i]
			ibeg = maxI64(beg-eoff, 0)
			iend = minI64(end, tots-eoff)
			err  = r.runTree(i, eoff+off, ibeg, iend, f)
		)
		if err != nil {
			return fmt.Errorf("rtree: could not process entry %d: %w", i, err)
		}
	}

	return nil
}

func (r *rchain) findTrees(beg, end int64) (int, int) {
	var (
		eoff int64
		ibeg = -1
		iend = -1
	)

	for i, t := range r.ch.trees {
		n := t.Entries()
		if ibeg < 0 && beg <= eoff {
			ibeg = i
		}
		if iend < 0 && end <= eoff+n {
			iend = i + 1
		}
		eoff += n
	}
	if iend < 0 {
		iend = len(r.ch.trees)
	}

	return ibeg, iend
}

func (r *rchain) runTree(itree int, off, beg, end int64, f func(RCtx) error) error {
	rr := newReader(r.ch.trees[itree], r.rvs, r.nrab, beg, end)
	return rr.run(off, beg, end, f)
}

func (r *rchain) reset() {}
