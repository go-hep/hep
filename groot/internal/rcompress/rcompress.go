// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rcompress provides types and functions to compress and decompress
// ROOT data payloads.
package rcompress // import "go-hep.org/x/hep/groot/internal/rcompress"

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
	"github.com/pierrec/xxHash/xxHash64"
	"github.com/ulikunitz/xz"
	"golang.org/x/xerrors"
)

// Kind specifies the compression algorithm
// to be used during reading or writing ROOT files.
type Kind int

// constants for compression/decompression
const (
	Inherit              Kind = -1
	UseGlobal            Kind = 0
	ZLIB                 Kind = +1
	LZMA                 Kind = +2
	OldCompression       Kind = +3
	LZ4                  Kind = +4
	ZSTD                 Kind = +5
	UndefinedCompression Kind = +6
)

const (
	zstdVersion = 1 // keep in sync with klauspost/compress/zstd and ROOT
)

var (
	// errNoCompression is returned when the compression algorithm
	// couldn't compress the input or when the compressed output is bigger
	// than the input
	errNoCompression = xerrors.Errorf("rcompress: no compression")
)

// Settings encodes the ROOT way of specifying a compression mechanism
// and its compression level.
type Settings struct {
	Alg Kind
	Lvl int
}

// DefaultSettings is the default compression algorithm and level used
// in ROOT files and trees.
var DefaultSettings = Settings{Alg: ZLIB, Lvl: flate.BestSpeed}

func (set Settings) Compression() int32 {
	var (
		lvl = set.Lvl
		alg = set.Alg
	)
	switch {
	case lvl == flate.DefaultCompression:
		switch alg {
		case ZLIB:
			lvl = 6
		case LZ4:
			lvl = 1
		case LZMA:
			lvl = 1
		case ZSTD:
			lvl = 1 // FIXME(sbinet): check with ROOT-6.20.00 default
		default:
			panic(xerrors.Errorf("rcompress: unknown compression algorithm: %v", alg))
		}
	case lvl == flate.BestSpeed:
		if alg == ZSTD {
			lvl = int(zstd.SpeedFastest)
		}
	case lvl == flate.BestCompression:
		if alg == ZSTD {
			lvl = int(zstd.SpeedBestCompression)
		}
	case lvl > 99:
		lvl = 99
	}
	return int32(alg*100) + int32(lvl)

}

// Note: this contains ZL[src][dst] where src and dst are 3 bytes each.
const HeaderSize = 9

// because each zipped block contains:
// - the size of the input data
// - the size of the compressed data
// where each size is saved on 3 bytes, the maximal size
// of each block can not be bigger than 16Mb.
const kMaxCompressedBlockSize = 0xffffff

// kindOf returns the kind of compression algorithm.
func kindOf(buf [HeaderSize]byte) Kind {
	switch {
	case buf[0] == 'Z' && buf[1] == 'L':
		return ZLIB
	case buf[0] == 'X' && buf[1] == 'Z':
		return LZMA
	case buf[0] == 'L' && buf[1] == '4':
		return LZ4
	case buf[0] == 'Z' && buf[1] == 'S':
		return ZSTD
	case buf[0] == 'C' && buf[1] == 'S':
		return OldCompression
	default:
		return UndefinedCompression
	}
}

func rootCompressAlgLvl(v int32) (Kind, int) {
	var (
		alg = Kind(v / 100)
		lvl = int(v % 100)
	)

	return alg, lvl
}

// Compress compresses src, using the compression kind and level encoded into compr.
// Users can provide a non-nil dst to reduce allocation.
func Compress(dst, src []byte, compr int32) ([]byte, error) {
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
		cur     = 0
		beg     int
		end     int
	)

	size := len(src) + nblocks*HeaderSize
	if dst == nil || len(dst) < size {
		dst = append(dst, make([]byte, size-len(dst))...)
	}

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

func compressBlock(alg Kind, lvl int, tgt, src []byte) (int, error) {
	// FIXME(sbinet): rework tgt/dst to reduce buffer allocation.

	var (
		err error
		hdr [HeaderSize]byte
		dst []byte

		srcsz = int32(len(src))
		dstsz int32
	)

	switch alg {
	case ZLIB:
		hdr[0] = 'Z'
		hdr[1] = 'L'
		hdr[2] = 8 // zlib deflated
		buf := new(bytes.Buffer)
		buf.Grow(len(src))
		w, err := zlib.NewWriterLevel(buf, lvl)
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not create ZLIB compressor: %w", err)
		}
		_, err = w.Write(src)
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not write ZLIB compressed bytes: %w", err)
		}
		err = w.Close()
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not close ZLIB compressor: %w", err)
		}
		dstsz = int32(buf.Len())
		if dstsz > srcsz {
			return 0, errNoCompression
		}
		dst = append(hdr[:], buf.Bytes()...)

	case LZMA:
		hdr[0] = 'X'
		hdr[1] = 'Z'
		cfg := xz.WriterConfig{
			CheckSum: xz.CRC32,
		}
		if err := cfg.Verify(); err != nil {
			return 0, xerrors.Errorf("riofs: could not create LZMA compressor config: %w", err)
		}
		buf := new(bytes.Buffer)
		buf.Grow(len(src))
		w, err := cfg.NewWriter(buf)
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not create LZMA compressor: %w", err)
		}

		_, err = w.Write(src)
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not write LZMA compressed bytes: %w", err)
		}

		err = w.Close()
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not close LZMA compressor: %w", err)
		}

		dstsz = int32(buf.Len())
		if dstsz > srcsz {
			return 0, errNoCompression
		}
		dst = append(hdr[:], buf.Bytes()...)

	case LZ4:
		hdr[0] = 'L'
		hdr[1] = '4'
		hdr[2] = lz4.Version

		const chksum = 8
		var room = int(float64(srcsz) * 2e-4) // lz4 needs some extra scratch space
		dst = make([]byte, HeaderSize+chksum+len(src)+room)
		buf := dst[HeaderSize:]
		var n int
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
			return 0, xerrors.Errorf("riofs: could not compress with LZ4: %w", err)
		}

		if n == 0 {
			// not compressible.
			return 0, errNoCompression
		}

		buf = buf[:n+chksum]
		dst = dst[:len(buf)+HeaderSize]
		binary.BigEndian.PutUint64(buf[:chksum], xxHash64.Checksum(buf[chksum:], 0))
		dstsz = int32(n + chksum)

	case ZSTD:
		hdr[0] = 'Z'
		hdr[1] = 'S'
		hdr[2] = zstdVersion

		buf := new(bytes.Buffer)
		buf.Grow(len(src))
		w, err := zstd.NewWriter(buf, zstd.WithEncoderLevel(zstd.EncoderLevel(lvl)))
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not create ZSTD compressor: %w", err)
		}

		_, err = w.Write(src)
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not write ZSTD compressed bytes: %w", err)
		}

		err = w.Close()
		if err != nil {
			return 0, xerrors.Errorf("riofs: could not close ZSTD compressor: %w", err)
		}

		dstsz = int32(buf.Len())
		if dstsz > srcsz {
			return 0, errNoCompression
		}
		dst = append(hdr[:], buf.Bytes()...)

	case OldCompression:
		return 0, xerrors.Errorf("riofs: old compression algorithm unsupported")

	default:
		return 0, xerrors.Errorf("riofs: unknown algorithm %d", alg)
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

// Decompress decompresses src into dst.
func Decompress(dst []byte, src io.Reader) error {
	var (
		beg    = 0
		end    = 0
		buflen = len(dst)
	)

	for end < buflen {
		var hdr [HeaderSize]byte
		_, err := io.ReadFull(src, hdr[:])
		if err != nil {
			return xerrors.Errorf("riofs: could not read compress header: %w", err)
		}

		srcsz := (int64(hdr[3]) | int64(hdr[4])<<8 | int64(hdr[5])<<16)
		tgtsz := int64(hdr[6]) | int64(hdr[7])<<8 | int64(hdr[8])<<16
		end += int(tgtsz)
		lr := &io.LimitedReader{R: src, N: srcsz}
		switch kindOf(hdr) {
		case ZLIB:
			rc, err := zlib.NewReader(lr)
			if err != nil {
				return xerrors.Errorf("riofs: could not create ZLIB reader: %w", err)
			}
			defer rc.Close()
			_, err = io.ReadFull(rc, dst[beg:end])
			if err != nil {
				return xerrors.Errorf("riofs: could not decompress ZLIB buffer: %w", err)
			}

		case LZ4:
			src := make([]byte, srcsz)
			_, err = io.ReadFull(lr, src)
			if err != nil {
				return err
			}
			const chksum = 8
			// FIXME: we skip the 32b checksum. use it!
			_, err = lz4.UncompressBlock(src[chksum:], dst[beg:end])
			if err != nil {
				return err
			}

		case LZMA:
			rc, err := xz.NewReader(lr)
			if err != nil {
				return err
			}
			_, err = io.ReadFull(rc, dst[beg:end])
			if err != nil {
				return err
			}
			if lr.N > 0 {
				// FIXME(sbinet): LZMA leaves some bytes on the floor...
				_, err = lr.Read(make([]byte, lr.N))
				if err != nil {
					return err
				}
			}

		case ZSTD:
			rc, err := zstd.NewReader(lr)
			if err != nil {
				return err
			}
			_, err = io.ReadFull(rc, dst[beg:end])
			if err != nil {
				return err
			}
			if lr.N > 0 {
				panic("zstd extra bytes")
			}

		default:
			panic(xerrors.Errorf("riofs: unknown compression algorithm %q", hdr[:2]))
		}
		beg = end
	}

	return nil
}
