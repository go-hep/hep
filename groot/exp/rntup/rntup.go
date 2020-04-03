// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rntup contains types to handle RNTuple-related data.
package rntup // import "go-hep.org/x/hep/groot/exp/rntup"

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

type span struct {
	seek   uint64
	nbytes uint32
	length uint32
}

type NTuple struct {
	rvers uint32
	size  uint32

	header span
	footer span

	reserved uint64
}

func (*NTuple) Class() string {
	return "ROOT::Experimental::RNTuple"
}

func (*NTuple) RVersion() int16 {
	return 0 // FIXME(sbinet): generate through gen.rboot
}

func (nt *NTuple) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(nt.RVersion())

	w.WriteU32(nt.rvers)
	w.WriteU32(nt.size)

	w.WriteU64(nt.header.seek)
	w.WriteU32(nt.header.nbytes)
	w.WriteU32(nt.header.length)

	w.WriteU64(nt.footer.seek)
	w.WriteU32(nt.footer.nbytes)
	w.WriteU32(nt.footer.length)

	w.WriteU64(nt.reserved)

	return w.SetByteCount(pos, nt.Class())
}

func (nt *NTuple) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	_ /*vers*/, pos, bcnt := r.ReadVersion(nt.Class())

	nt.rvers = r.ReadU32()
	nt.size = r.ReadU32()

	nt.header.seek = r.ReadU64()
	nt.header.nbytes = r.ReadU32()
	nt.header.length = r.ReadU32()

	nt.footer.seek = r.ReadU64()
	nt.footer.nbytes = r.ReadU32()
	nt.footer.length = r.ReadU32()

	nt.reserved = r.ReadU64()

	r.CheckByteCount(pos, bcnt, beg, nt.Class())
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := &NTuple{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("ROOT::Experimental::RNTuple", f)
	}
}

var (
	_ root.Object        = (*NTuple)(nil)
	_ rbytes.RVersioner  = (*NTuple)(nil)
	_ rbytes.Marshaler   = (*NTuple)(nil)
	_ rbytes.Unmarshaler = (*NTuple)(nil)
)
