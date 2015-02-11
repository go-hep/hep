// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"fmt"
)

// errorf returns a new formated error
func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// rioAlignU32 returns sz adjusted to align at 4-byte boundaries
func rioAlignU32(sz uint32) uint32 {
	return sz + (4-(sz&gAlign))&gAlign
}

// rioAlignU64 returns sz adjusted to align at 4-byte boundaries
func rioAlignU64(sz uint64) uint64 {
	return sz + (4-(sz&gAlign))&gAlign
}

// rioAlignI64 returns sz adjusted to align at 4-byte boundaries
func rioAlignI64(sz int64) int64 {
	return sz + (4-(sz&int64(gAlign)))&int64(gAlign)
}

// rioAlign returns sz adjusted to align at 4-byte boundaries
func rioAlign(sz int) int {
	return sz + (4-(sz&int(gAlign)))&int(gAlign)
}

// EOF
