// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

// start of payload in a TFile (in bytes)
const kBEGIN = 100

// kStartBigFile-1 is the largest position in a ROOT file before switching to
// the "big file" scheme (supporting files bigger than 4Gb) of ROOT.
const kStartBigFile = 2000000000

var (
	rootMagic = []byte("root")
)
