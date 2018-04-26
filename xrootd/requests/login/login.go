// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package login // import "go-hep.org/x/hep/xrootd/requests/login"

import (
	"os"
)

const RequestID uint16 = 3007

type Response struct {
	SessionID           [16]byte
	SecurityInformation []byte
}

type Request struct {
	Pid           int32
	UsernameBytes [8]byte
	Reserved1     byte
	Ability       byte
	Capabilities  byte
	Role          byte
	TokenLength   int32
	Token         []byte
}

func NewRequest(username string) Request {
	var usernameBytes [8]byte
	copy(usernameBytes[:], username)
	var ability = byte(00000)

	return Request{int32(os.Getpid()), usernameBytes, 0, ability, 1, 0, 0, []byte{}}
}
