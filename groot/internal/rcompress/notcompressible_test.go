// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package rcompress

import (
	"os"
	"testing"
)

func TestNotCompressible(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "groot-rcompress-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	src, err := os.ReadFile("testdata/not-compressible.raw")
	if err != nil {
		t.Fatalf("could not read reference data: %+v", err)
	}

	srcsz := len(src)
	tgtsz := len(src) + HeaderSize

	for _, tc := range []struct {
		name string
		alg  Kind
		lvl  int
		want int
	}{
		{name: "lz4", alg: LZ4, lvl: 1, want: srcsz},
		{name: "lzma", alg: LZMA, lvl: 1, want: srcsz},
		{name: "zlib", alg: ZLIB, lvl: 1, want: srcsz},
		{name: "zstd", alg: ZSTD, lvl: 1, want: srcsz},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tgt := make([]byte, tgtsz)
			n, err := compressBlock(tc.alg, tc.lvl, tgt, src)
			if err != nil && err != errNoCompression {
				t.Fatalf("could not compress block: %+v", err)
			}
			if got, want := n, tc.want; got != want {
				t.Fatalf("invalid output size: got=%d, want=%d", got, want)
			}
		})
	}
}
