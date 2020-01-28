// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"go-hep.org/x/hep/groot/rbytes"
)

type decoder struct {
	si *StreamerInfo
}

func (dec *decoder) DecodeROOT(ptr interface{}) error {
	panic("not implemented")
}

var (
	_ rbytes.Decoder = (*decoder)(nil)
)
