// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package rcompress

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNotCompressible(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "groot-rcompress-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	src, err := ioutil.ReadFile("testdata/not-compressible.raw")
	if err != nil {
		t.Fatalf("could not read reference data: %+v", err)
	}
	tgtsz := len(src) + HeaderSize

	for _, tc := range []struct {
		name string
		alg  Kind
		lvl  int
		want int
	}{
		{name: "lz4", alg: LZ4, lvl: 9, want: tgtsz + 8},
		{name: "lzma", alg: LZMA, lvl: 9, want: tgtsz},
		{name: "zlib", alg: ZLIB, lvl: 9, want: tgtsz},
		{name: "zstd", alg: ZSTD, lvl: 1, want: tgtsz},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tgt := make([]byte, len(src)+HeaderSize)
			n, err := compressBlock(tc.alg, tc.lvl, tgt, src)
			if err != nil {
				t.Fatalf("could not compress block: %+v", err)
			}
			if got, want := n, tc.want; got != want {
				t.Fatalf("invalid output size: got=%d, want=%d", got, want)
			}
		})
	}
}
