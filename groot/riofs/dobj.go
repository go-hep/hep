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

	hdr := r.ReadHeader(d.class)
	d.rvers = hdr.Vers
	d.size = hdr.Len
	r.SetPos(hdr.Pos + int64(hdr.Len) + 4)
	r.CheckHeader(hdr)
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
