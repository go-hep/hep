// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package rcompress_test

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rcompress"
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
#include <iostream>

void testcompress(const char *fname, int size) {
	auto f = TFile::Open(fname, "READ");
	auto str = (TObjString*)f->Get("str");
	if (str == nullptr) { exit(1); }
	if (str->GetString().Length() != size) {
		std::cerr << "invalid length: got=" << str->GetString().Length()
				  << " want=" << size << "\n";
		exit(2);
	}

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
					t.Fatalf("got:\n%s\nwant:\n%s", got, want)
				}

				if !rtests.HasROOT {
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

func TestRoundtrip(t *testing.T) {
	wants := map[string][]byte{
		"00-10kb": []byte(strings.Repeat("-+", 10*1024)),
		"01-16mb": []byte(strings.Repeat("-+", 16*1024*1024-2)), // remove 2 so divisible by kMaxCompressedBlockSize
	}
	keysOf := func(kvs map[string][]byte) []string {
		keys := make([]string, 0, len(kvs))
		for k := range kvs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	for _, tc := range []struct {
		name string
		opt  rcompress.Settings
	}{
		// lz4
		{name: "lz4-default", opt: rcompress.Settings{Alg: rcompress.LZ4, Lvl: flate.DefaultCompression}},
		// lzma
		{name: "lzma-default", opt: rcompress.Settings{Alg: rcompress.LZMA, Lvl: flate.DefaultCompression}},
		// zlib
		{name: "zlib-default", opt: rcompress.Settings{Alg: rcompress.ZLIB, Lvl: flate.DefaultCompression}},
		// zstd
		{name: "zstd-default", opt: rcompress.Settings{Alg: rcompress.ZSTD, Lvl: flate.DefaultCompression}},
	} {
		for _, k := range keysOf(wants) {
			tname := fmt.Sprintf("%s-%s", tc.name, k)
			t.Run(tname, func(t *testing.T) {
				defer func() {
					err := recover()
					if err != nil {
						t.Fatalf("test panicked: %q", err)
					}
				}()

				want := wants[k]
				xsrc, err := rcompress.Compress(nil, want, tc.opt.Compression())
				if err != nil {
					t.Fatalf("could not create compressed source: %+v", err)
				}
				xdst := make([]byte, len(want))
				err = rcompress.Decompress(xdst, bytes.NewReader(xsrc))
				if err != nil {
					t.Fatalf("could not decompress xsrc: %+v", err)
				}
				if !bytes.Equal(xdst, want) {
					t.Fatalf("round-trip failed: %+v", err)
				}
			})
		}
	}
}
func BenchmarkCompression(b *testing.B) {
	b.ReportAllocs()

	wants := map[string][]byte{
		"00-10kb": []byte(strings.Repeat("-+", 10*1024)),
		"01-10mb": []byte(strings.Repeat("-+", 10*1024*1024)),
		"02-16mb": []byte(strings.Repeat("-+", 16*1024*1024)),
	}
	keysOf := func(kvs map[string][]byte) []string {
		keys := make([]string, 0, len(kvs))
		for k := range kvs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	for _, tc := range []struct {
		name string
		opt  rcompress.Settings
	}{
		{name: "default", opt: rcompress.DefaultSettings},
		// lz4
		{name: "lz4-default", opt: rcompress.Settings{Alg: rcompress.LZ4, Lvl: flate.DefaultCompression}},
		{name: "lz4-1", opt: rcompress.Settings{Alg: rcompress.LZ4, Lvl: 1}},
		{name: "lz4-9", opt: rcompress.Settings{Alg: rcompress.LZ4, Lvl: 9}},
		{name: "lz4-best-speed", opt: rcompress.Settings{Alg: rcompress.LZ4, Lvl: flate.BestSpeed}},
		{name: "lz4-best-compr", opt: rcompress.Settings{Alg: rcompress.LZ4, Lvl: flate.BestCompression}},
		// lzma
		{name: "lzma-default", opt: rcompress.Settings{Alg: rcompress.LZMA, Lvl: flate.DefaultCompression}},
		{name: "lzma-1", opt: rcompress.Settings{Alg: rcompress.LZMA, Lvl: 1}},
		{name: "lzma-9", opt: rcompress.Settings{Alg: rcompress.LZMA, Lvl: 9}},
		{name: "lzma-best-speed", opt: rcompress.Settings{Alg: rcompress.LZMA, Lvl: flate.BestSpeed}},
		{name: "lzma-best-compr", opt: rcompress.Settings{Alg: rcompress.LZMA, Lvl: flate.BestCompression}},
		// zlib
		{name: "zlib-default", opt: rcompress.Settings{Alg: rcompress.ZLIB, Lvl: flate.DefaultCompression}},
		{name: "zlib-1", opt: rcompress.Settings{Alg: rcompress.ZLIB, Lvl: 1}},
		{name: "zlib-9", opt: rcompress.Settings{Alg: rcompress.ZLIB, Lvl: 9}},
		{name: "zlib-best-speed", opt: rcompress.Settings{Alg: rcompress.ZLIB, Lvl: flate.BestSpeed}},
		{name: "zlib-best-compr", opt: rcompress.Settings{Alg: rcompress.ZLIB, Lvl: flate.BestCompression}},
		// zstd
		{name: "zstd-default", opt: rcompress.Settings{Alg: rcompress.ZSTD, Lvl: flate.DefaultCompression}},
		{name: "zstd-1", opt: rcompress.Settings{Alg: rcompress.ZSTD, Lvl: 1}},
		{name: "zstd-9", opt: rcompress.Settings{Alg: rcompress.ZSTD, Lvl: 9}},
		{name: "zstd-best-speed", opt: rcompress.Settings{Alg: rcompress.ZSTD, Lvl: flate.BestSpeed}},
		{name: "zstd-best-compr", opt: rcompress.Settings{Alg: rcompress.ZSTD, Lvl: flate.BestCompression}},
	} {
		for _, k := range keysOf(wants) {
			want := wants[k]
			tname := fmt.Sprintf("%s-%s", tc.name, k)
			compr := tc.opt.Compression()
			defer func() {
				err := recover()
				if err != nil {
					b.Fatalf("%s panicked: %+v", tname, err)
				}
			}()
			xsrc, err := rcompress.Compress(nil, want, compr)
			if err != nil {
				b.Fatalf("could not create compressed source: %+v", err)
			}
			xdst := make([]byte, len(want))
			err = rcompress.Decompress(xdst, bytes.NewReader(xsrc))
			if err != nil {
				b.Fatalf("could not decompress xsrc: %+v", err)
			}
			if !bytes.Equal(xdst, want) {
				b.Fatalf("round-trip failed: %+v", err)
			}

			b.Run("enc-"+tname, func(b *testing.B) {
				src := want
				dst := make([]byte, len(src))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := rcompress.Compress(dst, src, compr)
					if err != nil {
						b.Fatalf("%+v", err)
					}
				}
			})

			b.Run("dec-"+tname, func(b *testing.B) {
				dst := make([]byte, len(want))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					src := bytes.NewReader(xsrc)
					err := rcompress.Decompress(dst, src)
					if err != nil {
						b.Fatalf("%+v", err)
					}
				}
			})
		}
	}
}
