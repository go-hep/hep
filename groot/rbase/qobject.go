// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type QObject struct {
}

func (*QObject) Class() string {
	return "TQObject"
}

func (*QObject) RVersion() int16 {
	return rvers.QObject
}

func (qo *QObject) UnmarshalROOT(r *rbytes.RBuffer) error {
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		var v QObject
		return reflect.ValueOf(&v)
	}
	rtypes.Factory.Add("TQObject", f)
}

var (
	_ root.Object        = (*QObject)(nil)
	_ rbytes.Unmarshaler = (*QObject)(nil)
)
