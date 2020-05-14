// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import "fmt"

type preader struct {
	ievt int64
	nevt int64

	tree   Tree
	rvars  []ReadVar
	rbs    []rbranch
	leaves []rleaf
}

func NewPReader(t Tree, rvars []ReadVar) (*preader, error) {
	pr := &preader{
		nevt:   t.Entries(),
		tree:   t,
		rvars:  rvars,
		leaves: make([]rleaf, len(rvars)),
	}

	//brs := make(map[string]struct{}, len(rvars))
	for i, rvar := range rvars {
		name := rvar.Leaf
		if name == "" {
			name = rvar.Name
		}
		pr.leaves[i] = rleafFrom(t.Leaf(name), rvar, pr)

		//brs[rvar.Name] = struct{}{}
	}

	//bnames := make([]string, 0, len(brs))
	//for k := range brs {
	//	bnames = append(bnames, k)
	//}

	const n = 128
	pr.rbs = make([]rbranch, len(rvars))
	for i, rvar := range rvars {
		pr.rbs[i] = newRBranch(t.Branch(rvar.Name), n, pr.leaves[i:i+1], pr)
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
