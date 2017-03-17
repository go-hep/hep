// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rootcnv provides tools to convert ROOT histograms to go-hep/hbook ones.
package rootcnv

import (
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/yodacnv"
)

// NewH1D creates a new H1D from a TH1x.
func NewH1D(r yodacnv.Marshaler) (*hbook.H1D, error) {
	raw, err := r.MarshalYODA()
	if err != nil {
		return nil, err
	}
	var h hbook.H1D
	err = h.UnmarshalYODA(raw)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
