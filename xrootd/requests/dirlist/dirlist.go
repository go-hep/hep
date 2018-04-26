// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dirlist // import "go-hep.org/x/hep/xrootd/requests/dirlist"

const RequestID uint16 = 3004

type Response struct {
	Data []byte
}

type Request struct {
	Reserved1  [15]byte
	Options    byte
	PathLength int32
	Path       []byte
}

func NewRequest(path string) Request {
	var pathBytes = make([]byte, len(path))
	copy(pathBytes, path)

	return Request{[15]byte{}, 0, int32(len(path)), pathBytes}
}
