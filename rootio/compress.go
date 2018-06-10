// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"compress/zlib"
	"io"

	"github.com/pierrec/lz4"
	"github.com/ulikunitz/xz"
)

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
			return err
		}

		srcsz := (int64(hdr[3]) | int64(hdr[4])<<8 | int64(hdr[5])<<16)
		tgtsz := int64(hdr[6]) | int64(hdr[7])<<8 | int64(hdr[8])<<16
		end += int(tgtsz)
		lr := io.LimitReader(r, srcsz)
		switch rootCompressAlg(hdr) {
		case kZLIB:
			rc, err := zlib.NewReader(lr)
			if err != nil {
				return err
			}
			defer rc.Close()
			_, err = io.ReadFull(rc, buf[beg:end])
			if err != nil {
				return err
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
			panic("rootio: unknown compression algorithm")
		}
		beg = end
	}

	return nil
}
