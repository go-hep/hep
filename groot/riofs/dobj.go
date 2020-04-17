// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

// dobject is a dummy placeholder object
type dobject struct {
	rvers int16
	size  int32
	class string
}

func (d dobject) Class() string {
	return d.class
}

func (d *dobject) SetClass(n string) { d.class = n }

func (d *dobject) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(d.class)
	d.rvers = vers
	d.size = bcnt
	r.SetPos(beg + int64(bcnt) + 4)
	r.CheckByteCount(pos, bcnt, beg, d.class)
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := &dobject{class: "*groot.dobject"}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("*groot.dobject", f)
	}
}

var (
	_ root.Object        = (*dobject)(nil)
	_ rbytes.Unmarshaler = (*dobject)(nil)
)
