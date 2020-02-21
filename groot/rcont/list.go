// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcont

import (
	"io"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
)

type List struct {
	obj  rbase.Object
	name string
	objs []root.Object
}

func NewList(name string, objs []root.Object) *List {
	list := &List{
		obj:  rbase.Object{ID: 0x0, Bits: 0x3000000},
		name: name,
		objs: objs,
	}
	return list
}

func (*List) RVersion() int16 {
	return rvers.List
}

func (*List) Class() string {
	return "TList"
}

func (li *List) UID() uint32 {
	return li.obj.UID()
}

func (li *List) Name() string {
	if li.name == "" {
		return "TList"
	}
	return li.name
}

func (*List) Title() string {
	return "Doubly linked list"
}

func (li *List) At(i int) root.Object {
	return li.objs[i]
}

func (li *List) Last() int {
	panic("not implemented")
}

func (li *List) Len() int {
	return len(li.objs)
}

func (li *List) Append(obj root.Object) {
	li.objs = append(li.objs, obj)
}

func (li *List) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(li.RVersion())
	if _, err := li.obj.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteString(li.name)
	w.WriteI32(int32(len(li.objs)))
	for _, obj := range li.objs {
		err := w.WriteObjectAny(obj)
		if err != nil {
			return 0, err
		}

		w.WriteU8(0) // FIXME(sbinet): properly serialize the 'OPTION'.
	}

	return w.SetByteCount(pos, li.Class())
}

func (li *List) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(li.Class())

	if vers <= 3 {
		return xerrors.Errorf("rcont: TList version too old (%d <= 3)", vers)
	}

	if err := li.obj.UnmarshalROOT(r); err != nil {
		return err
	}

	li.name = r.ReadString()
	size := int(r.ReadI32())

	li.objs = make([]root.Object, size)

	for i := range li.objs {
		obj := r.ReadObjectAny()
		// obj := r.ReadObjectRef()
		if obj == nil {
			panic(xerrors.Errorf("nil obj ref: %w", r.Err())) // FIXME(sbinet)
			// return r.Err()
		}
		li.objs[i] = obj

		n := int(r.ReadU8())
		if n > 0 {
			opt := make([]byte, n)
			io.ReadFull(r, opt)
			// drop the option on the floor. // FIXME(sbinet)
		}
	}

	r.CheckByteCount(pos, bcnt, beg, li.Class())
	return r.Err()
}

type HashList struct {
	List
}

func (*HashList) RVersion() int16 {
	return rvers.HashList
}

func (*HashList) Class() string {
	return "THashList"
}

func (li *HashList) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	return li.List.MarshalROOT(w)
}

func (li *HashList) UnmarshalROOT(r *rbytes.RBuffer) error {
	return li.List.UnmarshalROOT(r)
}

func init() {
	{
		f := func() reflect.Value {
			o := NewList("", nil)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TList", f)
	}
	{
		f := func() reflect.Value {
			o := &HashList{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("THashList", f)
	}
}

var (
	_ root.Object        = (*List)(nil)
	_ root.UIDer         = (*List)(nil)
	_ root.Collection    = (*List)(nil)
	_ root.SeqCollection = (*List)(nil)
	_ root.List          = (*List)(nil)
	_ rbytes.Marshaler   = (*List)(nil)
	_ rbytes.Unmarshaler = (*List)(nil)
)

var (
	_ root.Object        = (*HashList)(nil)
	_ root.UIDer         = (*HashList)(nil)
	_ root.Collection    = (*HashList)(nil)
	_ root.SeqCollection = (*HashList)(nil)
	_ root.List          = (*HashList)(nil)
	_ rbytes.Marshaler   = (*HashList)(nil)
	_ rbytes.Unmarshaler = (*HashList)(nil)
)
