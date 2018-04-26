// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package write // import "go-hep.org/x/hep/xrootd/requests/write"

const RequestID uint16 = 3019

type Request struct {
	FileHandle [4]byte
	Offset     int64
	PathID     byte
	Reserved   [3]byte
	DataLength int32
	Data       []byte
}

func NewRequest(fileHandle [4]byte, offset int64, pathID byte, data []byte) Request {
	return Request{fileHandle, offset, pathID, [3]byte{}, int32(len(data)), data}
}
