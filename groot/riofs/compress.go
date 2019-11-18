// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"compress/flate"

	"go-hep.org/x/hep/groot/internal/rcompress"
	"golang.org/x/xerrors"
)

func (f *File) setCompression(alg rcompress.Kind, lvl int) {
	switch {
	case lvl == flate.DefaultCompression:
		switch alg {
		case rcompress.LZ4:
			lvl = 1
		case rcompress.LZMA:
			lvl = 1
		case rcompress.ZLIB:
			lvl = 6
		default:
			panic(xerrors.Errorf("riofs: unknown compression algorithm: %v", alg))
		}
	case lvl > 99:
		lvl = 99
	}
	f.compression = int32(alg*100) + int32(lvl)
}

// WithLZ4 configures a ROOT file to use LZ4 as a compression mechanism.
func WithLZ4(level int) FileOption {
	return func(f *File) error {
		f.setCompression(rcompress.LZ4, level)
		return nil
	}
}

// WithLZMA configures a ROOT file to use LZMA as a compression mechanism.
func WithLZMA(level int) FileOption {
	return func(f *File) error {
		f.setCompression(rcompress.LZMA, level)
		return nil
	}
}

// WithoutCompression configures a ROOT file to not use any compression mechanism.
func WithoutCompression() FileOption {
	return func(f *File) error {
		f.setCompression(0, 0)
		return nil
	}
}

// WithZlib configures a ROOT file to use zlib as a compression mechanism.
func WithZlib(level int) FileOption {
	return func(f *File) error {
		f.setCompression(rcompress.ZLIB, level)
		return nil
	}
}
