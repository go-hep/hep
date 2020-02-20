// Copyright 2020 The go-hep Authors. All rights reserved.
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

// Ref implements a persistent link to a root.Object.
type Ref struct {
	// FIXME(sbinet): we should use root.Object instead of *rbase.Object
	obj *Object
}

func (*Ref) RVersion() int16 {
	return rvers.Ref
}

func (*Ref) Class() string { return "TRef" }
func (*Ref) Name() string  { return "TRef" }
func (*Ref) Title() string { return "Persistent Reference link to a TObject" }

func (ref *Ref) String() string {
	if ref.obj == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Ref{id:%d}", ref.obj.ID)
}

// Object returns the root.Object being referenced by this Ref.
func (ref *Ref) Object() root.Object {
	if ref.obj == nil {
		return nil
	}
	return ref.obj
}

func (ref *Ref) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}
	var (
		obj Object
		err = obj.UnmarshalROOT(r)
	)
	if err != nil {
		return err
	}

	ref.obj = &obj

	switch {
	case obj.TestBits(kHasUUID):
		_ = r.ReadString() // UUID string
	default:
		_ = r.ReadU16() // pid
	}

	return nil
}

func (ref *Ref) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	//	if w.Err() != nil {
	//		return 0, w.Err()
	//	}
	panic("not implemented")
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
