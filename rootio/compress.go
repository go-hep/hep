// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"io"

	"github.com/pierrec/lz4"
	"github.com/pierrec/xxHash/xxHash64"
	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
)

func (f *File) setCompression(alg compressAlgType, lvl int) {
	switch {
	case lvl == flate.DefaultCompression:
		switch alg {
		case kLZ4:
			lvl = 1
		case kLZMA:
			lvl = 1
		case kZLIB:
			lvl = 6
		default:
			panic(errors.Errorf("rootio: unknown compression algorithm: %v", alg))
		}
	case lvl > 99:
		lvl = 99
	}
	f.compression = int32(alg*100) + int32(lvl)
}

// WithLZ4 configures a ROOT file to use LZ4 as a compression mechanism.
func WithLZ4(level int) FileOption {
	return func(f *File) error {
		f.setCompression(kLZ4, level)
		return nil
	}
}

// WithLZMA configures a ROOT file to use LZMA as a compression mechanism.
func WithLZMA(level int) FileOption {
	return func(f *File) error {
		f.setCompression(kLZMA, level)
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
		f.setCompression(kZLIB, level)
		return nil
	}
}

func rootCompressAlg(buf [rootHDRSIZE]byte) compressAlgType {
	switch {
	case buf[0] == 'Z' && buf[1] == 'L':
		return kZLIB
	case buf[0] == 'X' && buf[1] == 'Z':
		return kLZMA
	case buf[0] == 'L' && buf[1] == '4':
		return kLZ4
	case buf[0] == 'C' && buf[1] == 'S':
		return kOldCompressionAlgo
	default:
		return kUndefinedCompressionAlgorithm
	}
}

func rootCompressAlgLvl(v int32) (compressAlgType, int) {
	var (
		alg = compressAlgType(v / 100)
		lvl = int(v % 100)
	)

	return alg, lvl
}

func compress(compr int32, src []byte) ([]byte, error) {
	alg, lvl := rootCompressAlgLvl(compr)

	if alg == 0 || lvl == 0 {
		// no compression
		return src, nil
	}

	var (
		err error
		hdr [rootHDRSIZE]byte
		dst []byte

		srcsz = int32(len(src))
		dstsz = srcsz
	)

	// FIXME(sbinet): handle multi-key-payload
	switch {
	case srcsz > 0xffffff || srcsz < 0:
		panic("rootio: invalid src size")
	case dstsz > 0xffffff:
		panic("rootio: invalid dst size")
	}

	switch compressAlgType(alg) {
	case kZLIB:
		hdr[0] = 'Z'
		hdr[1] = 'L'
		hdr[2] = 8 // zlib deflated
		buf := new(bytes.Buffer)
		buf.Grow(len(src))
		w, err := zlib.NewWriterLevel(buf, lvl)
		if err != nil {
			return nil, errors.Wrapf(err, "rootio: could not create ZLIB compressor")
		}
		_, err = w.Write(src)
		if err != nil {
			return nil, errors.Wrapf(err, "rootio: could not write ZLIB compressed bytes")
		}
		err = w.Close()
		if err != nil {
			return nil, errors.Wrapf(err, "rootio: could not close ZLIB compressor")
		}
		dstsz = int32(buf.Len())
		if dstsz > srcsz {
			return src, nil
		}
		dst = append(hdr[:], buf.Bytes()...)

	case kLZMA:
		hdr[0] = 'X'
		hdr[1] = 'Z'
		hdr[2] = 0
		cfg := xz.WriterConfig{
			CheckSum: xz.CRC32,
		}
		if err := cfg.Verify(); err != nil {
			return nil, errors.Wrapf(err, "rootio: could not create LZMA compressor config")
		}
		buf := new(bytes.Buffer)
		buf.Grow(len(src))
		w, err := cfg.NewWriter(buf)
		if err != nil {
			return nil, errors.Wrapf(err, "rootio: could not create LZMA compressor")
		}

		_, err = w.Write(src)
		if err != nil {
			return nil, errors.Wrapf(err, "rootio: could not write LZMA compressed bytes")
		}

		err = w.Close()
		if err != nil {
			return nil, errors.Wrapf(err, "rootio: could not close LZMA compressor")
		}

		dstsz = int32(buf.Len())
		if dstsz > srcsz {
			return src, nil
		}
		dst = append(hdr[:], buf.Bytes()...)

	case kLZ4:
		hdr[0] = 'L'
		hdr[1] = '4'
		hdr[2] = 1 // lz4 version

		dst = make([]byte, rootHDRSIZE+len(src)+8)
		buf := dst[rootHDRSIZE:]
		var n = 0
		switch {
		case lvl >= 4:
			if lvl > 9 {
				lvl = 9
			}
			n, err = lz4.CompressBlockHC(src, buf[8:], lvl)
		default:
			ht := make([]int, 1<<16)
			n, err = lz4.CompressBlock(src, buf[8:], ht)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "rootio: could not compress with LZ4")
		}

		if n == 0 {
			// not compressible.
			return src, nil
		}

		buf = buf[:n+8]
		dst = dst[:len(buf)+rootHDRSIZE]
		binary.BigEndian.PutUint64(buf[:8], xxHash64.Checksum(buf[8:], 0))
		dstsz = int32(n + 8)

	case kOldCompressionAlgo:
		return nil, errors.Errorf("rootio: old compression algorithm unsupported")

	default:
		return nil, errors.Errorf("rootio: unknown algorithm %d", alg)
	}

	hdr[6] = byte(srcsz)
	hdr[7] = byte(srcsz >> 8)
	hdr[8] = byte(srcsz >> 16)

	hdr[3] = byte(dstsz)
	hdr[4] = byte(dstsz >> 8)
	hdr[5] = byte(dstsz >> 16)

	copy(dst, hdr[:])

	return dst, nil
}

func decompress(r io.Reader, buf []byte) error {
	var (
		beg    = 0
		end    = 0
		buflen = len(buf)
	)

	for end < buflen {
		var hdr [rootHDRSIZE]byte
		_, err := io.ReadFull(r, hdr[:])
		if err != nil {
			return errors.Wrapf(err, "rootio: could not read compress header")
		}

		srcsz := (int64(hdr[3]) | int64(hdr[4])<<8 | int64(hdr[5])<<16)
		tgtsz := int64(hdr[6]) | int64(hdr[7])<<8 | int64(hdr[8])<<16
		end += int(tgtsz)
		lr := io.LimitReader(r, srcsz)
		switch rootCompressAlg(hdr) {
		case kZLIB:
			rc, err := zlib.NewReader(lr)
			if err != nil {
				return errors.Wrapf(err, "rootio: could not create ZLIB reader")
			}
			defer rc.Close()
			_, err = io.ReadFull(rc, buf[beg:end])
			if err != nil {
				return errors.Wrapf(err, "rootio: could not decompress ZLIB buffer")
			}

		case kLZ4:
			src := make([]byte, srcsz)
			_, err = io.ReadFull(lr, src)
			if err != nil {
				return err
			}
			const chksum = 8
			// FIXME: we skip the 32b checksum. use it!
			_, err = lz4.UncompressBlock(src[chksum:], buf[beg:end])
			if err != nil {
				return err
			}

		case kLZMA:
			rc, err := xz.NewReader(lr)
			if err != nil {
				return err
			}
			_, err = io.ReadFull(rc, buf[beg:end])
			if err != nil {
				return err
			}

		default:
			panic(errors.Errorf("rootio: unknown compression algorithm %q", hdr[:2]))
		}
		beg = end
	}

	return nil
}
