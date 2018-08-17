// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io"
	"reflect"
)

type tlist struct {
	rvers int16
	obj   tobject
	name  string
	objs  []Object
}

func (li *tlist) Class() string {
	return "TList"
}

func (li *tlist) Name() string {
	if li.name == "" {
		return "TList"
	}
	return li.name
}

func (li *tlist) At(i int) Object {
	return li.objs[i]
}

func (li *tlist) Last() int {
	panic("not implemented")
}

func (li *tlist) Len() int {
	return len(li.objs)
}

func (li *tlist) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(li.rvers)
	if _, err := li.obj.MarshalROOT(w); err != nil {
		w.err = err
		return 0, w.err
	}

	w.WriteString(li.name)
	w.WriteI32(int32(len(li.objs)))
	for _, obj := range li.objs {
		err := w.WriteObjectAny(obj)
		if err != nil {
			w.err = err
			return 0, w.err
		}

		w.WriteU8(0) // FIXME(sbinet): properly serialize the 'OPTION'.
	}

	return w.SetByteCount(pos, "TList")
}

func (li *tlist) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	li.rvers = vers

	if vers <= 3 {
		return fmt.Errorf("rootio: TList version too old (%d <= 3)", vers)
	}

	if err := li.obj.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	li.name = r.ReadString()
	size := int(r.ReadI32())

	li.objs = make([]Object, size)

	for i := range li.objs {
		obj := r.ReadObjectAny()
		// obj := r.ReadObjectRef()
		if obj == nil {
			panic(fmt.Errorf("nil obj ref: %v\n", r.Err())) // FIXME(sbinet)
			// return r.Err()
		}
		li.objs[i] = obj

		n := int(r.ReadU8())
		if n > 0 {
			opt := make([]byte, n)
			io.ReadFull(r.r, opt)
			// drop the option on the floor. // FIXME(sbinet)
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TList")
	return r.Err()
}

type thashList struct {
	tlist
}

func (*thashList) Class() string {
	return "THashList"
}

func (li *thashList) MarshalROOT(w *WBuffer) (int, error) {
	return li.tlist.MarshalROOT(w)
}

func (li *thashList) UnmarshalROOT(r *RBuffer) error {
	return li.tlist.UnmarshalROOT(r)
}

func init() {
	{
		f := func() reflect.Value {
			o := &tlist{}
			return reflect.ValueOf(o)
		}
		Factory.add("TList", f)
		Factory.add("*rootio.tlist", f)
	}
	{
		f := func() reflect.Value {
			o := &thashList{}
			return reflect.ValueOf(o)
		}
		Factory.add("THashList", f)
		Factory.add("*rootio.thashList", f)
	}
}

var (
	_ Object          = (*tlist)(nil)
	_ Collection      = (*tlist)(nil)
	_ SeqCollection   = (*tlist)(nil)
	_ List            = (*tlist)(nil)
	_ ROOTMarshaler   = (*tlist)(nil)
	_ ROOTUnmarshaler = (*tlist)(nil)
)

var (
	_ Object          = (*thashList)(nil)
	_ Collection      = (*thashList)(nil)
	_ SeqCollection   = (*thashList)(nil)
	_ List            = (*thashList)(nil)
	_ ROOTMarshaler   = (*thashList)(nil)
	_ ROOTUnmarshaler = (*thashList)(nil)
)
