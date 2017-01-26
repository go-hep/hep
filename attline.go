// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"reflect"
)

type attline struct {
	color int16
	style int16
	width int16
}

func (a *attline) UnmarshalROOT(data *bytes.Buffer) error {
	dec := newDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt := dec.readVersion()
	myprintf("attline-vers=%v\n", vers)

	dec.readBin(&a.color)
	dec.readBin(&a.style)
	dec.readBin(&a.width)
	dec.checkByteCount(pos, bcnt, start, "TAttLine")

	return dec.err
}

func init() {
	f := func() reflect.Value {
		o := &attline{}
		return reflect.ValueOf(o)
	}
	Factory.add("TAttLine", f)
	Factory.add("*rootio.attline", f)
}

var _ ROOTUnmarshaler = (*attline)(nil)
