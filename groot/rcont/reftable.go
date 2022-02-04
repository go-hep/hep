// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcont

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type RefTable struct {
	obj     rbase.Object
	size    int32       // dummy, for backward compatibility
	parents *ObjArray   // array of Parent objects (eg TTree branch) holding the referenced objects
	owner   root.Object // Object owning this TRefTable
	guids   []string    // UUIDs of TProcessIDs used in fParentIDs
}

func NewRefTable(owner root.Object) *RefTable {
	return &RefTable{
		obj:   *rbase.NewObject(),
		owner: owner,
	}
}

func (*RefTable) RVersion() int16 {
	return rvers.RefTable
}

func (*RefTable) Class() string {
	return "TRefTable"
}

func (tbl *RefTable) UID() uint32 {
	return tbl.obj.UID()
}

func (tbl *RefTable) At(i int) root.Object {
	panic("not implemented")
}

func (tbl *RefTable) UIDs() []string {
	return tbl.guids
}

func (tbl *RefTable) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(tbl.RVersion())
	w.WriteObject(&tbl.obj)
	w.WriteI32(tbl.size)
	{
		var obj root.Object
		if tbl.parents != nil {
			obj = tbl.parents
		}
		w.WriteObjectAny(obj)
	}
	{
		var obj root.Object
		if tbl.owner != nil {
			obj = tbl.owner
		}
		w.WriteObjectAny(obj)
	}
	w.WriteStdVectorStrs(tbl.guids)

	return w.SetByteCount(pos, tbl.Class())
}

func (tbl *RefTable) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(tbl.Class())

	if vers < 3 {
		return fmt.Errorf("rcont: TRefTable version too old (%d < 3)", vers)
	}

	r.ReadObject(&tbl.obj)
	tbl.size = r.ReadI32()
	{
		obj := r.ReadObjectAny()
		tbl.parents = nil
		if obj != nil {
			tbl.parents = obj.(*ObjArray)
		}
	}
	{
		obj := r.ReadObjectAny()
		tbl.owner = nil
		if obj != nil {
			tbl.owner = obj
		}
	}
	r.ReadStdVectorStrs(&tbl.guids)

	r.CheckByteCount(pos, bcnt, beg, tbl.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := NewRefTable(nil)
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TRefTable", f)
}

var (
	_ root.Object        = (*RefTable)(nil)
	_ rbytes.Marshaler   = (*RefTable)(nil)
	_ rbytes.Unmarshaler = (*RefTable)(nil)
)
