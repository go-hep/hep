// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package rcompress_test

import (
	"compress/flate"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func TestCompress(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "groot-rcompress-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	wants := map[string]root.Object{
		"small": rbase.NewObjString("hello"),
		"10mb":  rbase.NewObjString(strings.Repeat("-+", 10*1024*1024)),
		"16mb":  rbase.NewObjString(strings.Repeat("-+", 16*1024*1024)),
	}

	const macroROOT = `
#include "TFile.h"
#include "TObjString.h"

void testcompress(const char *fname, int size) {
	auto f = TFile::Open(fname, "READ");
	auto str = (TObjString*)f->Get("str");
	if (str == nullptr) { exit(1); }
	if (str->GetString().Length() != size) { exit(2); }

	exit(0);
}
`
	for _, tc := range []struct {
		name string
		opt  riofs.FileOption
	}{
		{name: "default", opt: func(*riofs.File) error { return nil }},
		{name: "default-nil", opt: nil},
		{name: "no-compr", opt: riofs.WithoutCompression()},
		// lz4
		{name: "lz4-default", opt: riofs.WithLZ4(flate.DefaultCompression)},
		{name: "lz4-0", opt: riofs.WithLZ4(0)},
		{name: "lz4-1", opt: riofs.WithLZ4(1)},
		{name: "lz4-9", opt: riofs.WithLZ4(9)},
		{name: "lz4-best-speed", opt: riofs.WithLZ4(flate.BestSpeed)},
		{name: "lz4-best-compr", opt: riofs.WithLZ4(flate.BestCompression)},
		// lzma
		{name: "lzma-default", opt: riofs.WithLZMA(flate.DefaultCompression)},
		{name: "lzma-0", opt: riofs.WithLZMA(0)},
		{name: "lzma-1", opt: riofs.WithLZMA(1)},
		{name: "lzma-9", opt: riofs.WithLZMA(9)},
		{name: "lzma-best-speed", opt: riofs.WithLZMA(flate.BestSpeed)},
		{name: "lzma-best-compr", opt: riofs.WithLZMA(flate.BestCompression)},
		// zlib
		{name: "zlib-default", opt: riofs.WithZlib(flate.DefaultCompression)},
		{name: "zlib-0", opt: riofs.WithZlib(0)},
		{name: "zlib-1", opt: riofs.WithZlib(1)},
		{name: "zlib-9", opt: riofs.WithZlib(9)},
		{name: "zlib-best-speed", opt: riofs.WithZlib(flate.BestSpeed)},
		{name: "zlib-best-compr", opt: riofs.WithZlib(flate.BestCompression)},
		// zstd
		{name: "zstd-default", opt: riofs.WithZstd(flate.DefaultCompression)},
		{name: "zstd-0", opt: riofs.WithZstd(0)},
		{name: "zstd-1", opt: riofs.WithZstd(1)},
		{name: "zstd-9", opt: riofs.WithZstd(9)},
		{name: "zstd-best-speed", opt: riofs.WithZstd(flate.BestSpeed)},
		{name: "zstd-best-compr", opt: riofs.WithZstd(flate.BestCompression)},
	} {
		for k, want := range wants {
			if (k == "16mb" || k == "10mb") &&
				!strings.HasSuffix(tc.name, "best-compr") &&
				!strings.HasSuffix(tc.name, "-1") {
				continue
			}
			tname := fmt.Sprintf("%s-%s", k, tc.name)
			t.Run(tname, func(t *testing.T) {
				fname := filepath.Join(dir, "test-"+tname+".root")
				w, err := riofs.Create(fname, tc.opt)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				defer w.Close()

				err = w.Put("str", want)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				err = w.Close()
				if err != nil {
					t.Fatalf("%+v", err)
				}

				r, err := riofs.Open(fname)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				defer r.Close()

				obj, err := r.Get("str")
				if err != nil {
					t.Fatalf("%+v", err)
				}
				str := obj.(root.ObjString)

				if got, want := str.String(), want.(root.ObjString).String(); got != want {
					t.Fatalf("got=%q, want=%q", got, want)
				}

				if !rtests.HasROOT {
					return
				}

				if strings.Contains(tname, "zstd") {
					// FIXME(sbinet): do run test when ROOT-6.20/00 is out.
					return
				}

				out, err := rtests.RunCxxROOT("testcompress", []byte(macroROOT), fname, len(want.(root.ObjString).String()))
				if err != nil {
					t.Fatalf("error: %+v\n%s\n", err, out)
				}
			})
		}
	}
}
