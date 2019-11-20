// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"go-hep.org/x/hep/groot/internal/rcompress"
)

func (f *File) setCompression(alg rcompress.Kind, lvl int) {
	f.compression = rcompress.Settings{Alg: alg, Lvl: lvl}.Compression()
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
