// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"compress/flate"
	"testing"
)

func TestRecord(t *testing.T) {
	for _, tc := range []struct {
		name  string
		compr CompressorKind
		want  CompressorKind
		level int
		codec int
	}{
		{
			name:  "zlib",
			compr: CompressZlib,
			level: flate.DefaultCompression,
			codec: 3,
		},
		{
			name:  "gzip",
			compr: CompressGzip,
			level: 9,
			codec: 4,
		},
		{
			name:  "lza",
			compr: CompressLZA,
			level: flate.BestCompression,
			codec: 4,
		},
		{
			name:  "lzo",
			compr: CompressLZO,
			level: flate.BestSpeed,
			codec: 0,
		},
		{
			name:  "snappy",
			compr: CompressSnappy,
			level: flate.DefaultCompression,
			codec: 0,
		},
		{
			name:  "none",
			compr: CompressNone,
			level: 0,
			codec: 0,
		},
		{
			name:  "default",
			compr: CompressDefault,
			want:  CompressZlib,
			level: flate.DefaultCompression,
			codec: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {

			rec := newRecord("name", NewOptions(tc.compr, tc.level, tc.codec))
			if got, want := rec.Name(), "name"; got != want {
				t.Fatalf("got=%q, want=%q", got, want)
			}

			switch tc.compr {
			case CompressNone:
				if rec.Compress() {
					t.Fatalf("record compressed")
				}
			default:
				if !rec.Compress() {
					t.Fatalf("record not compressed")
				}
			}

			if rec.Unpack() {
				t.Fatalf("record unpacked")
			}

			rec.SetUnpack(true)
			if !rec.Unpack() {
				t.Fatalf("record not unpacked")
			}

			opts := rec.Options()
			tc.want = tc.compr
			if tc.want == CompressDefault {
				tc.want = CompressZlib
			}
			if got, want := opts.CompressorKind(), tc.want; got != want {
				t.Fatalf("invalid compressor kind: got=%v, want=%v", got, want)
			}

			if got, want := opts.CompressorLevel(), tc.level; got != want {
				t.Fatalf("invalid compressor level: got=%v, want=%v", got, want)
			}

			if got, want := opts.CompressorCodec(), tc.codec; got != want {
				t.Fatalf("invalid compressor codec: got=%v, want=%v", got, want)
			}
		})
	}
}
