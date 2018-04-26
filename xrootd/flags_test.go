// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import "flag"

var Addr = flag.String("addr", "0.0.0.0:9001", "address of xrootd server")
