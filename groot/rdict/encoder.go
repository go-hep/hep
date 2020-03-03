// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"go-hep.org/x/hep/groot/rbytes"
)

type encoder struct {
	w    *rbytes.WBuffer
	si   *StreamerInfo
	kind rbytes.StreamKind
	wops []wstreamOp
}

func newEncoder(w *rbytes.WBuffer, si *StreamerInfo, kind rbytes.StreamKind, ops []wstreamOp) (*encoder, error) {
	return &encoder{w, si, kind, ops}, nil
}

func (enc *encoder) EncodeROOT(ptr interface{}) error {
	panic("not implemented")
}

var (
	_ rbytes.Encoder = (*encoder)(nil)
)
