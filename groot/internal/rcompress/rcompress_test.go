// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package rcompress_test

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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

	cxxROOT, err := exec.LookPath("root.exe")
	withCxxROOT := err == nil

	macroROOT := `
void testcompress(const char *fname, int size) {
	auto f = TFile::Open(fname, "READ");
	auto str = (TObjString*)f->Get("str");
	if (str == nullptr) { exit(1); }
	if (str->GetString().Length() != size) { exit(2); }

	exit(0);
}
`
	macro := filepath.Join(dir, "testcompress.C")
	err = ioutil.WriteFile(macro, []byte(macroROOT), 0644)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		opt  riofs.FileOption
	}{
		{name: "default", opt: func(*riofs.File) error { return nil }},
		{name: "default-nil", opt: nil},
		{name: "no-compr", opt: riofs.WithoutCompression()},
		{name: "lz4-default", opt: riofs.WithLZ4(flate.DefaultCompression)},
		{name: "lz4-0", opt: riofs.WithLZ4(0)},
		{name: "lz4-1", opt: riofs.WithLZ4(1)},
		{name: "lz4-9", opt: riofs.WithLZ4(9)},
		{name: "lz4-best-speed", opt: riofs.WithLZ4(flate.BestSpeed)},
		{name: "lz4-best-compr", opt: riofs.WithLZ4(flate.BestCompression)},
		{name: "lzma-default", opt: riofs.WithLZMA(flate.DefaultCompression)},
		{name: "lzma-0", opt: riofs.WithLZMA(0)},
		{name: "lzma-1", opt: riofs.WithLZMA(1)},
		{name: "lzma-9", opt: riofs.WithLZMA(9)},
		{name: "lzma-best-speed", opt: riofs.WithLZMA(flate.BestSpeed)},
		{name: "lzma-best-compr", opt: riofs.WithLZMA(flate.BestCompression)},
		{name: "zlib-default", opt: riofs.WithZlib(flate.DefaultCompression)},
		{name: "zlib-0", opt: riofs.WithZlib(0)},
		{name: "zlib-1", opt: riofs.WithZlib(1)},
		{name: "zlib-9", opt: riofs.WithZlib(9)},
		{name: "zlib-best-speed", opt: riofs.WithZlib(flate.BestSpeed)},
		{name: "zlib-best-compr", opt: riofs.WithZlib(flate.BestCompression)},
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
					t.Fatal(err)
				}
				defer w.Close()

				err = w.Put("str", want)
				if err != nil {
					t.Fatal(err)
				}

				err = w.Close()
				if err != nil {
					t.Fatal(err)
				}

				r, err := riofs.Open(fname)
				if err != nil {
					t.Fatal(err)
				}
				defer r.Close()

				obj, err := r.Get("str")
				if err != nil {
					t.Fatal(err)
				}
				str := obj.(root.ObjString)

				if got, want := str.String(), want.(root.ObjString).String(); got != want {
					t.Fatalf("got=%q, want=%q", got, want)
				}

				if !withCxxROOT {
					return
				}

				buf := new(bytes.Buffer)
				cmd := exec.Command(cxxROOT, "-b", fmt.Sprintf("%s(%q, %d)", macro, fname, len(want.(root.ObjString).String())))
				cmd.Stdout = buf
				cmd.Stderr = buf
				err = cmd.Run()
				if err != nil {
					t.Fatalf("error: %v\n%s\n", err, buf.String())
				}
			})
		}
	}
}
