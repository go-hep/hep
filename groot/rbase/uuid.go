// Copyright 2018 The go-hep Authors. All rights reserved.
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

type UUID [16]byte

func (*UUID) Class() string {
	return "TUUID"
}

func (*UUID) RVersion() int16 {
	return rvers.UUID
}

func (*UUID) Sizeof() int32 {
	// UUID-version (u16)
	// UUID-payload (16-bytes)
	return 18
}

func (uuid *UUID) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}
	pos := w.Pos()
	w.WriteU16(uint16(uuid.RVersion()))
	w.Write((*uuid)[:])
	end := w.Pos()

	return int(end - pos), w.Err()
}

func (uuid *UUID) UnmarshalROOTv1(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	r.Read((*uuid)[:])

	return r.Err()
}

func (uuid *UUID) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	_ = r.ReadU16() // version
	r.Read((*uuid)[:])

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &UUID{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TUUID", f)
}

var (
	_ root.Object        = (*UUID)(nil)
	_ rbytes.Marshaler   = (*UUID)(nil)
	_ rbytes.Unmarshaler = (*UUID)(nil)
)
