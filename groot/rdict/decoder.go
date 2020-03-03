// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
)

type decoder struct {
	r    *rbytes.RBuffer
	si   *StreamerInfo
	kind rbytes.StreamKind
	rops []rstreamer
}

func newDecoder(r *rbytes.RBuffer, si *StreamerInfo, kind rbytes.StreamKind, ops []rstreamer) (*decoder, error) {
	return &decoder{r, si, kind, ops}, nil
}

func (dec *decoder) DecodeROOT(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rdict: invalid kind (got=%T, want=pointer)", ptr)
	}

	for i, op := range dec.rops {
		err := op.rstream(dec.r, ptr)
		if err != nil {
			return fmt.Errorf("rdict: could not read element %d from %q: %w", i, dec.si.Name(), err)
		}
	}

	return nil
}

var (
	_ rbytes.Decoder = (*decoder)(nil)
)
