// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"

	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
)

func BenchmarkDumpLocal(b *testing.B) {
	// big-file.root is: rtests.XrdRemote("testdata/SMHiggsToZZTo4L.root")
	benchDump(b, "../../testdata/big-file.root", "Events")
}

func BenchmarkDumpHTTP(b *testing.B) {
	for _, i := range []time.Duration{
		0,
		5 * time.Millisecond,
		10 * time.Millisecond,
		100 * time.Millisecond,
	} {
		func(i time.Duration) {
			srv := httptest.NewServer(server{delay: i, h: http.FileServer(http.Dir("../../testdata"))})
			defer srv.Close()

			b.Run(fmt.Sprintf("delay-%v", i), func(b *testing.B) {
				// big-file.root is: rtests.XrdRemote("testdata/SMHiggsToZZTo4L.root")
				benchDump(b, srv.URL+"/big-file.root", "Events")
			})
		}(i)
	}
}

type server struct {
	delay time.Duration
	h     http.Handler
}

func (srv server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(srv.delay)
	srv.h.ServeHTTP(w, r)
}

func benchDump(b *testing.B, fname, tname string) {
	f, err := groot.Open(fname)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	o, err := riofs.Dir(f).Get(tname)
	if err != nil {
		b.Fatal(err)
	}

	tree := o.(rtree.Tree)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := dump(tree)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func dump(tree rtree.Tree) (n int64, err error) {
	r, err := rtree.NewReader(tree, rtree.NewReadVars(tree))
	if err != nil {
		return 0, err
	}
	defer r.Close()

	err = r.Read(func(rctx rtree.RCtx) error {
		_ = rctx.Entry
		n++
		return nil
	})

	return n, err
}
