// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !go1.7

package rootio

import "os"

const (
	ioSeekCurrent = os.SEEK_CUR
	ioSeekStart   = os.SEEK_SET
	ioSeekEnd     = os.SEEK_END
)
