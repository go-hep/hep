// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdfuse // import "go-hep.org/x/hep/xrootd/xrdfuse"

func init() {
	// add a C++ XRootD test server hosted in CC-Lyon.
	testClientAddrs = append(testClientAddrs, "ccxrootdgotest.in2p3.fr:9001")
}
