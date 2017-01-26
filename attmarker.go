// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"reflect"
)

type attmarker struct {
	color int16
	style int16
	width float32
}

func (a *attmarker) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	dec := NewDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	}

	err = dec.readBin(&a.color)
	if err != nil {
		return err
	}

	err = dec.readBin(&a.style)
	if err != nil {
		return err
	}

	err = dec.readBin(&a.width)
	if err != nil {
		return err
	}

	err = dec.checkByteCount(pos, bcnt, start, "TAttMarker")
	return err
}

func init() {
	f := func() reflect.Value {
		o := &attmarker{}
		return reflect.ValueOf(o)
	}
	Factory.db["TAttMarker"] = f
	Factory.db["*rootio.attmarker"] = f
}

var _ ROOTUnmarshaler = (*attmarker)(nil)
