// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd_test

import (
	"testing"

	"go-hep.org/x/hep/groot"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/rtree"
)

func BenchmarkRead(b *testing.B) {
	const (
		fname = "root://ccxrootdgotest.in2p3.fr:9001/tmp/rootio/testdata/SMHiggsToZZTo4L.root"
		tname = "Events"
	)

	f, err := groot.Open(fname)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get(tname)
	if err != nil {
		b.Fatal(err)
	}

	tree := o.(rtree.Tree)
	read := func(tree rtree.Tree) error {
		r, err := rtree.NewReader(tree, rtree.NewReadVars(tree))
		if err != nil {
			b.Fatal(err)
		}
		defer r.Close()

		return r.Read(func(rctx rtree.RCtx) error {
			_ = rctx.Entry
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := read(tree)
		if err != nil {
			b.Fatal(err)
		}
	}
}
