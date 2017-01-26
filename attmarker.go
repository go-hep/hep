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
	dec := newDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt := dec.readVersion()
	myprintf("attmarker-vers=%v\n", vers)
	dec.readBin(&a.color)
	dec.readBin(&a.style)
	dec.readBin(&a.width)
	dec.checkByteCount(pos, bcnt, start, "TAttMarker")

	return dec.err
}

func init() {
	f := func() reflect.Value {
		o := &attmarker{}
		return reflect.ValueOf(o)
	}
	Factory.add("TAttMarker", f)
	Factory.add("*rootio.attmarker", f)
}

var _ ROOTUnmarshaler = (*attmarker)(nil)
