// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rtypes contains the means to register types (ROOT ones and user
// defined ones) with the ROOT type factory.
package rtypes // import "go-hep.org/x/hep/groot/rtypes"

// Bit returns a uint32 with v-th bit set to 1.
func Bit(v int) uint32 {
	return uint32(1) << uint32(v)
}
