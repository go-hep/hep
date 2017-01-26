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
	var err error
	dec := newDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	} else {
		myprintf("attline: %v %v %v\n", vers, pos, bcnt)
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

	err = dec.checkByteCount(pos, bcnt, start, "TAttLine")
	return err
}

func init() {
	f := func() reflect.Value {
		o := &attline{}
		return reflect.ValueOf(o)
	}
	Factory.db["TAttLine"] = f
	Factory.db["*rootio.attline"] = f
}

var _ ROOTUnmarshaler = (*attline)(nil)
