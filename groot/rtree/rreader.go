// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
)

type preader struct {
	ievt int64
	nevt int64

	tree   Tree
	rvars  []ReadVar
	rbs    []rbranch
	leaves []rleaf
}

func sanitizeRVars(t Tree, rvars []ReadVar) []ReadVar {
	for i := range rvars {
		rvar := &rvars[i]
		if rvar.Leaf == "" {
			rvar.Leaf = rvar.Name
		}
		if rvar.count == "" {
			leaf := t.Leaf(rvar.Leaf).LeafCount()
			if leaf != nil {
				rvar.count = leaf.Name()
			}
		}
	}
	return rvars
}

func NewPReader(t Tree, rvars []ReadVar) (*preader, error) {
	return newPReader(t, rvars, 128)
}

func newPReader(t Tree, rvars []ReadVar, n int) (*preader, error) {
	rvars = sanitizeRVars(t, rvars)

	pr := &preader{
		nevt:  t.Entries(),
		tree:  t,
		rvars: rvars,
	}
	usr := make(map[string]struct{}, len(rvars))
	for _, rvar := range rvars {
		usr[rvar.Name+"."+rvar.Leaf] = struct{}{}
	}

	var rcounts []ReadVar
	for _, rvar := range rvars {
		if rvar.count == "" {
			continue
		}
		leaf := t.Leaf(rvar.Leaf).LeafCount()
		name := leaf.Branch().Name() + "." + leaf.Name()
		if _, ok := usr[name]; !ok {
			var ptr interface{}
			switch leaf := leaf.(type) {
			case *LeafB:
				ptr = new(int8)
			case *LeafS:
				ptr = new(int16)
			case *LeafI:
				ptr = new(int32)
			case *LeafL:
				ptr = new(int64)
			default:
				panic(fmt.Errorf("unknown Leaf count type %T", leaf))
			}
			rcounts = append(rcounts, ReadVar{
				Name:  leaf.Branch().Name(),
				Leaf:  leaf.Name(),
				Value: ptr,
			})
		}
	}
	pr.rvars = append(rcounts, pr.rvars...)

	pr.leaves = make([]rleaf, len(pr.rvars))
	for i, rvar := range pr.rvars {
		br := t.Branch(rvar.Name)
		leaf := br.Leaf(rvar.Leaf)
		pr.leaves[i] = rleafFrom(leaf, rvar, pr)
	}

	// regroup leaves by holding branch
	set := make(map[string]int)
	brs := make([][]rleaf, 0, len(pr.leaves))
	for _, leaf := range pr.leaves {
		br := leaf.Leaf().Branch().Name()
		if _, ok := set[br]; !ok {
			set[br] = len(brs)
			brs = append(brs, []rleaf{})
		}
		id := set[br]
		brs[id] = append(brs[id], leaf)
	}

	pr.rbs = make([]rbranch, len(brs))
	for i, leaves := range brs {
		branch := leaves[0].Leaf().Branch()
		pr.rbs[i] = newRBranch(branch, n, leaves, pr)
	}

	return pr, nil
}

func (pr *preader) Close() error {
	// FIXME(sbinet)
	return nil
}

func (pr *preader) reset() {
	for i := range pr.rbs {
		rb := &pr.rbs[i]
		rb.reset()
	}
}

func (pr *preader) rcount(name string) func() int {
	for _, leaf := range pr.leaves {
		n := leaf.Leaf().Name()
		if n != name {
			continue
		}
		switch leaf := leaf.(type) {
		case *rleafValI32:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValI64:
			return func() int {
				return int(*leaf.v)
			}
		default:
			panic(fmt.Errorf("rleaf %T not implemented", leaf))
		}
	}
	panic("impossible")
}

func (pr *preader) Read(f func(RCtx) error) error {
	var (
		err  error
		rctx RCtx
	)

	for i := range pr.rbs {
		rb := &pr.rbs[i]
		err = rb.start()
		if err != nil {
			return err
		}
	}
	defer func() {
		for i := range pr.rbs {
			rb := &pr.rbs[i]
			rb.stop()
		}
	}()

	for i := int64(0); i < pr.nevt; i++ {
		err = pr.read(i)
		if err != nil {
			return err
		}
		rctx.Entry = i
		err = f(rctx)
		if err != nil {
			return err
		}
	}

	return err
}

func (pr *preader) next() bool {
	return true
}

func (pr *preader) read(ievt int64) error {
	for i := range pr.rbs {
		rb := &pr.rbs[i]
		err := rb.read(ievt)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	_ rleafCtx = (*preader)(nil)
)
