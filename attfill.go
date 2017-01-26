// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"reflect"
)

type attfill struct {
	color int16
	style int16
}

func (a *attfill) UnmarshalROOT(data *bytes.Buffer) error {
	dec := newDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt := dec.readVersion()
	myprintf("attfill-vers=%v\n", vers)

	dec.readBin(&a.color)
	dec.readBin(&a.style)
	dec.checkByteCount(pos, bcnt, start, "TAttFill")

	return dec.err
}

func init() {
	f := func() reflect.Value {
		o := &attfill{}
		return reflect.ValueOf(o)
	}
	Factory.add("TAttFill", f)
	Factory.add("*rootio.attfill", f)
}

var _ ROOTUnmarshaler = (*attfill)(nil)
