// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package truncate // import "go-hep.org/x/hep/xrootd/requests/truncate"

const RequestID uint16 = 3028

type Response struct {
	Data []byte
}

type Request struct {
	FileHandle [4]byte
	Size       int64
	Reserved   [4]byte
	PathLength int32
	Path       []byte
}

func NewRequestWithHandle(fileHandle [4]byte, size int64) Request {
	return Request{fileHandle, size, [4]byte{}, 0, []byte{}}
}

func NewRequestWithPath(path string, size int64) Request {
	var pathBytes = make([]byte, len(path))
	copy(pathBytes, path)
	return Request{[4]byte{}, size, [4]byte{}, int32(len(path)), pathBytes}
}
