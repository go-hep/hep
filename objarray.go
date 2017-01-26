// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"reflect"
)

type objarray struct {
	arr []Object
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (arr *objarray) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	panic("not implemented")
	return err
}

func init() {
	f := func() reflect.Value {
		o := &objarray{
			arr: make([]Object, 0),
		}
		return reflect.ValueOf(o)
	}
	Factory.add("TObjArray", f)
	Factory.add("*rootio.objarray", f)
}

//var _ Object = (*objarray)(nil) // FIXME(sbinet)
var _ ROOTUnmarshaler = (*objarray)(nil)
