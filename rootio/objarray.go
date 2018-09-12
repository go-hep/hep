// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
)

type tobjarray struct {
	rvers int16
	obj   tobject
	name  string
	last  int
	arr   []Object
	low   int32
}

func (arr *tobjarray) Class() string {
	return "TObjArray"
}

func (arr *tobjarray) Name() string {
	n := arr.name
	if n == "" {
		return "TObjArray"
	}
	return n
}

func (arr *tobjarray) Title() string {
	return "An array of objects"
}

func (arr *tobjarray) At(i int) Object {
	return arr.arr[i]
}

func (arr *tobjarray) Last() int {
	return arr.last
}

func (arr *tobjarray) Len() int {
	return len(arr.arr)
}

func (arr *tobjarray) LowerBound() int {
	return int(arr.low)
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (arr *tobjarray) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(arr.rvers)
	arr.obj.MarshalROOT(w)
	w.WriteString(arr.name)

	w.WriteI32(int32(len(arr.arr)))
	w.WriteI32(arr.low)

	for _, obj := range arr.arr {
		w.WriteObjectAny(obj)
	}

	return w.SetByteCount(pos, "TObjArray")
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (arr *tobjarray) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	arr.rvers = vers

	if vers > 2 {
		if err := arr.obj.UnmarshalROOT(r); err != nil {
			return err
		}
	}
	if vers > 1 {
		arr.name = r.ReadString()
	}

	nobjs := int(r.ReadI32())
	arr.low = r.ReadI32()

	arr.arr = make([]Object, nobjs)
	arr.last = -1
	for i := range arr.arr {
		obj := r.ReadObjectAny()
		if r.err != nil {
			return r.err
		}
		if obj != nil {
			arr.last = i
			arr.arr[i] = obj
		}
	}

	r.CheckByteCount(pos, bcnt, start, "TObjArray")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &tobjarray{
			arr: make([]Object, 0),
		}
		return reflect.ValueOf(o)
	}
	Factory.add("TObjArray", f)
	Factory.add("*rootio.tobjarray", f)
}

var (
	_ Object          = (*tobjarray)(nil)
	_ Named           = (*tobjarray)(nil)
	_ ObjArray        = (*tobjarray)(nil)
	_ ROOTMarshaler   = (*tobjarray)(nil)
	_ ROOTUnmarshaler = (*tobjarray)(nil)
)
