// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package read // import "go-hep.org/x/hep/xrootd/requests/read"

const RequestID uint16 = 3013

type Response struct {
	Data []byte
}

type Request struct {
	FileHandle [4]byte
	Offset     int64
	Length     int32
	ArgsLength int32
}

func NewRequest(fileHandle [4]byte, offset int64, length int32) Request {
	return Request{fileHandle, offset, length, 0}
}
