// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-uuid"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
)

var (
	gPID = ProcessID{
		named: *NewNamed("ProcessID0", mustGenUUID()),
		objs:  make(map[uint32]root.Object),
	}
)

func mustGenUUID() string {
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(xerrors.Errorf("groot/rbase: could not generate UUID: %+v", err))
	}
	return id
}

// Ref implements a persistent link to a root.Object.
type Ref struct {
	obj Object
	pid *ProcessID
}

func (*Ref) RVersion() int16 {
	return rvers.Ref
}

func (*Ref) Class() string { return "TRef" }
func (*Ref) Name() string  { return "TRef" }
func (*Ref) Title() string { return "Persistent Reference link to a TObject" }

func (ref *Ref) UID() uint32 {
	return ref.obj.UID()
}

func (ref *Ref) String() string {
	return fmt.Sprintf("Ref{id:%d}", ref.obj.ID)
}

// Object returns the root.Object being referenced by this Ref.
func (ref *Ref) Object() root.Object {
	uid := ref.UID()
	if uid == 0 {
		return nil
	}
	obj, ok := ref.pid.objs[uid]
	if !ok {
		return nil
	}
	return obj
}

func (ref *Ref) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if err := ref.obj.UnmarshalROOT(r); err != nil {
		return err
	}

	switch {
	case ref.obj.TestBits(kHasUUID):
		_ = r.ReadString() // UUID string
	default:
		_ = r.ReadU16() // pid
	}

	return nil
}

func (ref *Ref) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	beg := w.Pos()
	if _, err := ref.obj.MarshalROOT(w); err != nil {
		return 0, err
	}

	switch {
	case ref.obj.TestBits(kHasUUID):
		panic("rbase: TRef with UUID not supported")
	default:
		w.WriteU16(uint16(ref.UID())) // FIXME(sbinet): this should go thru TFile.
	}

	return int(w.Pos() - beg), w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &Ref{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TRef", f)
}

var (
	_ root.Object        = (*Ref)(nil)
	_ root.Named         = (*Ref)(nil)
	_ rbytes.Marshaler   = (*Ref)(nil)
	_ rbytes.Unmarshaler = (*Ref)(nil)
)
