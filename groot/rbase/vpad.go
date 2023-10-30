// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type VirtualPad struct {
	obj     Object
	attline AttLine
	attfill AttFill
	attpad  AttPad
	qobj    QObject
}

func (*VirtualPad) Class() string {
	return "TVirtualPad"
}

func (*VirtualPad) RVersion() int16 {
	return rvers.VirtualPad
}

func (vpad *VirtualPad) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(vpad.Class())
	if hdr.Vers > rvers.VirtualPad {
		panic(fmt.Errorf("rbase: invalid virtualpad version=%d > %d", hdr.Vers, rvers.VirtualPad))
	}

	r.ReadObject(&vpad.obj)
	r.ReadObject(&vpad.attline)
	r.ReadObject(&vpad.attfill)
	r.ReadObject(&vpad.attpad)
	if hdr.Vers > 1 {
		r.ReadObject(&vpad.qobj)
	}

	r.CheckHeader(hdr)
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		var v VirtualPad
		return reflect.ValueOf(&v)
	}
	rtypes.Factory.Add("TVirtualPad", f)
}

var (
	_ root.Object        = (*VirtualPad)(nil)
	_ rbytes.Unmarshaler = (*VirtualPad)(nil)
)
