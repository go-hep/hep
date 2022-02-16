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

type ObjArray struct {
	obj  rbase.Object
	name string
	last int
	objs []root.Object
	low  int32
}

func NewObjArray() *ObjArray {
	return &ObjArray{
		obj:  *rbase.NewObject(),
		objs: make([]root.Object, 0),
	}
}

func (*ObjArray) RVersion() int16 {
	return rvers.ObjArray
}

func (*ObjArray) Class() string {
	return "TObjArray"
}

func (arr *ObjArray) UID() uint32 {
	return arr.obj.UID()
}

func (arr *ObjArray) Name() string {
	n := arr.name
	if n == "" {
		return "TObjArray"
	}
	return n
}

func (arr *ObjArray) Title() string {
	return "An array of objects"
}

func (arr *ObjArray) TestBits(bits uint32) bool {
	return arr.obj.TestBits(bits)
}

func (arr *ObjArray) At(i int) root.Object {
	return arr.objs[i]
}

func (arr *ObjArray) Last() int {
	return arr.last
}

func (arr *ObjArray) Len() int {
	return len(arr.objs)
}

func (arr *ObjArray) LowerBound() int {
	return int(arr.low)
}

func (arr *ObjArray) SetElems(v []root.Object) {
	arr.objs = v
	arr.last = len(v) - 1
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (arr *ObjArray) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(arr.Class(), arr.RVersion())
	w.WriteObject(&arr.obj)
	w.WriteString(arr.name)

	w.WriteI32(int32(len(arr.objs)))
	w.WriteI32(arr.low)

	for _, obj := range arr.objs {
		w.WriteObjectAny(obj)
	}

	return w.SetHeader(hdr)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (arr *ObjArray) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(arr.Class())
	if hdr.Vers > rvers.ObjArray {
		panic(fmt.Errorf("rcont: invalid TObjArray version=%d > %d", hdr.Vers, rvers.ObjArray))
	}

	if hdr.Vers > 2 {
		r.ReadObject(&arr.obj)
	}
	if hdr.Vers > 1 {
		arr.name = r.ReadString()
	}

	nobjs := int(r.ReadI32())
	arr.low = r.ReadI32()

	arr.objs = make([]root.Object, nobjs)
	arr.last = -1
	for i := range arr.objs {
		obj := r.ReadObjectAny()
		if r.Err() != nil {
			return r.Err()
		}
		if obj != nil {
			arr.last = i
			arr.objs[i] = obj
		}
	}

	r.CheckHeader(hdr)
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := NewObjArray()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TObjArray", f)
}

var (
	_ root.Object        = (*ObjArray)(nil)
	_ root.UIDer         = (*ObjArray)(nil)
	_ root.Named         = (*ObjArray)(nil)
	_ root.ObjArray      = (*ObjArray)(nil)
	_ rbytes.Marshaler   = (*ObjArray)(nil)
	_ rbytes.Unmarshaler = (*ObjArray)(nil)
)
