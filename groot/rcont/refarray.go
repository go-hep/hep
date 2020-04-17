// Copyright ©2017 The go-hep Authors. All rights reserved.
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

type RefArray struct {
	obj  rbase.Object
	name string
	// pid   *rbase.ProcessID // FIXME(sbinet)
	refs  []uint32 // uids of referenced objects
	lower int32    // lower bound of array
	last  int32    // last element in array containing an object
}

func NewRefArray() *RefArray {
	return &RefArray{
		refs: make([]uint32, 0),
		last: -1,
	}
}

func (*RefArray) RVersion() int16 {
	return rvers.RefArray
}

func (*RefArray) Class() string {
	return "TRefArray"
}

func (arr *RefArray) UID() uint32 {
	return arr.obj.UID()
}

func (arr *RefArray) Name() string {
	if arr.name == "" {
		return "TRefArray"
	}
	return arr.name
}

func (*RefArray) Title() string {
	return "An array of references to TObjects"
}

func (arr *RefArray) At(i int) root.Object {
	panic("not implemented")
}

func (arr *RefArray) Last() int {
	return int(arr.last)
}

func (arr *RefArray) Len() int {
	return len(arr.refs)
}

func (arr *RefArray) UIDs() []uint32 {
	return arr.refs
}

func (arr *RefArray) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(arr.RVersion())
	if _, err := arr.obj.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteString(arr.name)
	w.WriteI32(int32(len(arr.refs)))
	w.WriteI32(arr.lower)
	w.WriteI16(0) // FIXME(sbinet): handle fPID ProcessID

	w.WriteFastArrayU32(arr.refs)

	return w.SetByteCount(pos, arr.Class())
}

func (arr *RefArray) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(arr.Class())

	if vers < 1 {
		return fmt.Errorf("rcont: TRefArray version too old (%d < 1)", vers)
	}

	if err := arr.obj.UnmarshalROOT(r); err != nil {
		return err
	}

	arr.name = r.ReadString()
	size := int(r.ReadI32())
	arr.lower = r.ReadI32()
	arr.last = -1
	_ = r.ReadU16() // pid

	arr.refs = make([]uint32, size)
	for i := range arr.refs {
		arr.refs[i] = r.ReadU32()
		if arr.refs[i] != 0 {
			arr.last = int32(i)
		}
	}

	r.CheckByteCount(pos, bcnt, beg, arr.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := NewRefArray()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TRefArray", f)
}

var (
	_ root.Object        = (*RefArray)(nil)
	_ root.Named         = (*RefArray)(nil)
	_ root.SeqCollection = (*RefArray)(nil)
	_ rbytes.Marshaler   = (*RefArray)(nil)
	_ rbytes.Unmarshaler = (*RefArray)(nil)
)
