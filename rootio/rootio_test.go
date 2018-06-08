// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

func XrdRemote(fname string) string {
	const remote = "root://ccxrootdgotest.in2p3.fr:9001/tmp/rootio"
	return remote + "/" + fname
}
