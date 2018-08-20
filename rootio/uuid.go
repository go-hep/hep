// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type tuuid [16]byte

func (*tuuid) Class() string {
	return "TUUID"
}

func (*tuuid) sizeof() int32 { return 18 }

func (uuid *tuuid) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	w.write((*uuid)[:])
	return 16, w.err
}

func (uuid *tuuid) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}
	r.read((*uuid)[:])
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &tuuid{}
		return reflect.ValueOf(o)
	}
	Factory.add("TUUID", f)
	Factory.add("*rootio.tuuid", f)
}

var (
	_ Object          = (*tuuid)(nil)
	_ ROOTMarshaler   = (*tuuid)(nil)
	_ ROOTUnmarshaler = (*tuuid)(nil)
)
