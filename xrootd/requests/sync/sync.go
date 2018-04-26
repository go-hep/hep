// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync // import "go-hep.org/x/hep/xrootd/requests/sync"

const RequestID uint16 = 3016

type Request struct {
	FileHandle [4]byte
	Reserved1  [12]byte
	Reserved2  int32
}

func NewRequest(fileHandle [4]byte) Request {
	return Request{fileHandle, [12]byte{}, 0}
}
