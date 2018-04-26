// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package close // import "go-hep.org/x/hep/xrootd/requests/close"

const RequestID uint16 = 3003

type Request struct {
	FileHandle [4]byte
	FileSize   int64
	Reserved1  [4]byte
	Reserved2  int32
}

func NewRequest(fileHandle [4]byte, fileSize int64) Request {
	return Request{fileHandle, fileSize, [4]byte{}, 0}
}
