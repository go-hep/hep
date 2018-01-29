// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

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
