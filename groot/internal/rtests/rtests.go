// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtests

import (
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
)

type ROOTer interface {
	root.Object
	rbytes.Marshaler
	rbytes.Unmarshaler
}

func XrdRemote(fname string) string {
	const remote = "root://ccxrootdgotest.in2p3.fr:9001/tmp/rootio"
	return remote + "/" + fname
}
