// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

type Event struct {
	Header RunHeader
}

type RunHeader struct {
	RunNbr  int64 // run number
	EvtNbr  int64 // event number
	Trigger int64 // trigger word
}
