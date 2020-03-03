// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/pierrec/lz4"
	"github.com/pierrec/xxHash/xxHash64"
	"github.com/ulikunitz/xz"
)

type compressAlgType int

// constants for compression/decompression
const (
	kZLIB                          compressAlgType = 1
	kLZMA                          compressAlgType = 2
	kOldCompressionAlgo            compressAlgType = 3
	kLZ4                           compressAlgType = 4
	kUndefinedCompressionAlgorithm compressAlgType = 5
)

var (
	// errNoCompression is returned when the compression algorithm
	// couldn't compress the input or when the compressed output is bigger
	// than the input
	errNoCompression = errors.New("rootio: no compression")
)

// Note: this contains ZL[src][dst] where src and dst are 3 bytes each.
const rootHDRSIZE = 9

// because each zipped block contains:
// - the size of the input data
// - the size of the compressed data
// where each size is saved on 3 bytes, the maximal size
// of each block can not be bigger than 16Mb.
const kMaxCompressedBlockSize = 0xffffff

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
			panic(fmt.Errorf("rootio: unknown compression algorithm: %v", alg))
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
	const (
		blksz = kMaxCompressedBlockSize // 16Mb
	)

	alg, lvl := rootCompressAlgLvl(compr)

	if alg == 0 || lvl == 0 || len(src) < 512 {
		// no compression
		return src, nil
	}

	var (
		nblocks = len(src)/blksz + 1
		dst     = make([]byte, len(src)+nblocks*rootHDRSIZE)
		cur     = 0
		beg     = 0
		end     = 0
	)

	for beg = 0; beg < len(src); beg += blksz {
		end = beg + blksz
		if end > len(src) {
			end = len(src)
		}
		// FIXME(sbinet): split out into compressBlock{Zlib,LZ4,...}
		n, err := compressBlock(alg, lvl, dst[cur:], src[beg:end])
		switch err {
		case nil:
			cur += n
		case errNoCompression:
			return src, nil
		default:
			return nil, err
		}
	}

	return dst[:cur], nil
}

func compressBlock(alg compressAlgType, lvl int, tgt, src []byte) (int, error) {
	// FIXME(sbinet): rework tgt/dst to reduce buffer allocation.

	var (
		err error
		hdr [rootHDRSIZE]byte
		dst []byte

		srcsz = int32(len(src))
		dstsz = srcsz
	)

	switch alg {
	case kZLIB:
		hdr[0] = 'Z'
		hdr[1] = 'L'
		hdr[2] = 8 // zlib deflated
		buf := new(bytes.Buffer)
		buf.Grow(len(src))
		w, err := zlib.NewWriterLevel(buf, lvl)
		if err != nil {
			return 0, fmt.Errorf("rootio: could not create ZLIB compressor: %w", err)
		}
		_, err = w.Write(src)
		if err != nil {
			return 0, fmt.Errorf("rootio: could not write ZLIB compressed bytes: %w", err)
		}
		err = w.Close()
		if err != nil {
			return 0, fmt.Errorf("rootio: could not close ZLIB compressor: %w", err)
		}
		dstsz = int32(buf.Len())
		if dstsz > srcsz {
			return 0, errNoCompression
		}
		dst = append(hdr[:], buf.Bytes()...)

	case kLZMA:
		hdr[0] = 'X'
		hdr[1] = 'Z'
		cfg := xz.WriterConfig{
			CheckSum: xz.CRC32,
		}
		if err := cfg.Verify(); err != nil {
			return 0, fmt.Errorf("rootio: could not create LZMA compressor config: %w", err)
		}
		buf := new(bytes.Buffer)
		buf.Grow(len(src))
		w, err := cfg.NewWriter(buf)
		if err != nil {
			return 0, fmt.Errorf("rootio: could not create LZMA compressor: %w", err)
		}

		_, err = w.Write(src)
		if err != nil {
			return 0, fmt.Errorf("rootio: could not write LZMA compressed bytes: %w", err)
		}

		err = w.Close()
		if err != nil {
			return 0, fmt.Errorf("rootio: could not close LZMA compressor: %w", err)
		}

		dstsz = int32(buf.Len())
		if dstsz > srcsz {
			return 0, errNoCompression
		}
		dst = append(hdr[:], buf.Bytes()...)

	case kLZ4:
		hdr[0] = 'L'
		hdr[1] = '4'
		hdr[2] = lz4.Version

		const chksum = 8
		var room = int(float64(srcsz) * 2e-4) // lz4 needs some extra scratch space
		dst = make([]byte, rootHDRSIZE+chksum+len(src)+room)
		buf := dst[rootHDRSIZE:]
		var n = 0
		switch {
		case lvl >= 4:
			if lvl > 9 {
				lvl = 9
			}
			n, err = lz4.CompressBlockHC(src, buf[chksum:], lvl)
		default:
			ht := make([]int, 1<<16)
			n, err = lz4.CompressBlock(src, buf[chksum:], ht)
		}
		if err != nil {
			return 0, fmt.Errorf("rootio: could not compress with LZ4: %w", err)
		}

		if n == 0 {
			// not compressible.
			return 0, errNoCompression
		}

		buf = buf[:n+chksum]
		dst = dst[:len(buf)+rootHDRSIZE]
		binary.BigEndian.PutUint64(buf[:chksum], xxHash64.Checksum(buf[chksum:], 0))
		dstsz = int32(n + chksum)

	case kOldCompressionAlgo:
		return 0, fmt.Errorf("rootio: old compression algorithm unsupported")

	default:
		return 0, fmt.Errorf("rootio: unknown algorithm %d", alg)
	}

	if dstsz > kMaxCompressedBlockSize {
		return 0, errNoCompression
	}

	hdr[3] = byte(dstsz)
	hdr[4] = byte(dstsz >> 8)
	hdr[5] = byte(dstsz >> 16)

	hdr[6] = byte(srcsz)
	hdr[7] = byte(srcsz >> 8)
	hdr[8] = byte(srcsz >> 16)

	copy(dst, hdr[:])
	n := copy(tgt, dst)

	return n, nil
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
			return fmt.Errorf("rootio: could not read compress header: %w", err)
		}

		srcsz := (int64(hdr[3]) | int64(hdr[4])<<8 | int64(hdr[5])<<16)
		tgtsz := int64(hdr[6]) | int64(hdr[7])<<8 | int64(hdr[8])<<16
		end += int(tgtsz)
		lr := &io.LimitedReader{R: r, N: srcsz}
		switch rootCompressAlg(hdr) {
		case kZLIB:
			rc, err := zlib.NewReader(lr)
			if err != nil {
				return fmt.Errorf("rootio: could not create ZLIB reader: %w", err)
			}
			defer rc.Close()
			_, err = io.ReadFull(rc, buf[beg:end])
			if err != nil {
				return fmt.Errorf("rootio: could not decompress ZLIB buffer: %w", err)
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
			if lr.N > 0 {
				// FIXME(sbinet): LZMA leaves some bytes on the floor...
				lr.Read(make([]byte, lr.N))
			}

		default:
			panic(fmt.Errorf("rootio: unknown compression algorithm %q", hdr[:2]))
		}
		beg = end
	}

	return nil
}
