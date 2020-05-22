// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
)

type rjoin struct {
	j *join

	rs []*rtree // FIXME(sbinet): handle join of chains?

	rvs  []ReadVar
	nrab int
	beg  int64
	end  int64
}

func newRJoin(t *join, rvars []ReadVar, n int, beg, end int64) *rjoin {
	rvars = bindRVarsTo(t, rvars)
	r := &rjoin{
		j:    t,
		rs:   make([]*rtree, len(t.trees)),
		rvs:  rvars,
		nrab: n,
		beg:  beg,
		end:  end,
	}
	rps := make([][]ReadVar, len(r.rs))
	for i, t := range r.j.trees {
		rps[i] = r.loadRVars(t.(*ttree), rvars)
	}

	r.rvs = r.rvs[:0]
	for i, tree := range t.trees {
		r.rs[i] = newRTree(tree.(*ttree), rps[i], r.nrab, beg, end)
		r.rvs = append(r.rvs, r.rs[i].rvars()...)
	}

	return r
}

func (r *rjoin) loadRVars(t *ttree, rvars []ReadVar) []ReadVar {
	rps := make([]ReadVar, 0, len(rvars))
	for _, rv := range rvars {
		br := asBranch(rv.leaf.Branch())
		if br.tree != t {
			continue
		}
		rps = append(rps, rv)
	}
	return rps
}

func (r *rjoin) Close() error {
	var err error
	for _, rr := range r.rs {
		e := rr.Close()
		if e != nil && err == nil {
			err = e
		}
	}
	return err
}

func (r *rjoin) rvars() []ReadVar { return r.rvs }

func (r *rjoin) reset() {
	for _, rr := range r.rs {
		rr.reset()
	}
}

func (r *rjoin) run(off, beg, end int64, f func(RCtx) error) error {
	var (
		err  error
		rctx RCtx
	)
	defer r.Close()

	err = r.start()
	if err != nil {
		return err
	}
	defer r.stop()

	for i := beg; i < end; i++ {
		err = r.read(i)
		if err != nil {
			return fmt.Errorf("rtree: could not read entry %d: %w", i, err)
		}
		rctx.Entry = i + off
		err = f(rctx)
		if err != nil {
			return fmt.Errorf("rtree: could not process entry %d: %w", i, err)
		}
	}

	return err
}

func (r *rjoin) read(ievt int64) error {
	for _, rr := range r.rs {
		err := rr.read(ievt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rjoin) start() error {
	for _, rr := range r.rs {
		err := rr.start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rjoin) stop() {
	for _, rr := range r.rs {
		rr.stop()
	}
}

var (
	_ reader = (*rjoin)(nil)
)
