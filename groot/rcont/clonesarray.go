// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcont

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// ClonesArray implements a ROOT TClonesArray.
type ClonesArray struct {
	arr ObjArray
	cls string
}

func NewClonesArray() *ClonesArray {
	arr := &ClonesArray{
		arr: *NewObjArray(),
	}
	arr.BypassStreamer(false)
	arr.arr.obj.SetBits(rbytes.CannotHandleMemberWiseStreaming)
	return arr
}

func (*ClonesArray) RVersion() int16 {
	return rvers.ClonesArray
}

func (arr *ClonesArray) Class() string {
	return "TClonesArray"
}

func (arr *ClonesArray) Name() string {
	n := arr.arr.name
	if n == "" {
		return "TClonesArray"
	}
	return n
}

func (arr *ClonesArray) Title() string {
	return "object title"
}

func (arr *ClonesArray) At(i int) root.Object {
	return arr.arr.At(i)
}

func (arr *ClonesArray) Last() int {
	return arr.arr.Last()
}

func (arr *ClonesArray) Len() int {
	return arr.arr.Len()
}

func (arr *ClonesArray) LowerBound() int {
	return arr.arr.LowerBound()
}

func (arr *ClonesArray) SetElems(v []root.Object) {
	if arr.cls == "" {
		arr.cls = v[0].Class()
	}
	arr.arr.SetElems(v)
}

func (arr *ClonesArray) TestBits(bits uint32) bool {
	return arr.arr.TestBits(bits)
}

func (arr *ClonesArray) BypassStreamer(bypass bool) {
	switch bypass {
	case true:
		arr.arr.obj.SetBit(rbytes.BypassStreamer)
	default:
		arr.arr.obj.ResetBit(rbytes.BypassStreamer)
	}
}

func (arr *ClonesArray) CanBypassStreamer() bool {
	return arr.TestBits(rbytes.BypassStreamer)
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (arr *ClonesArray) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	bypass := false
	// make sure the status of bypass-streamer is part of the buffer.
	if arr.TestBits(rbytes.CannotHandleMemberWiseStreaming) {
		bypass = arr.CanBypassStreamer()
		arr.BypassStreamer(false)
	}

	si, err := w.StreamerInfo(arr.cls, -1)
	if err != nil {
		w.SetErr(errors.Wrapf(err, "rcont: could not find streamer for TClonesArray element %q", arr.cls))
		return 0, w.Err()
	}
	clsv := si.ClassVersion()

	pos := w.WriteVersion(arr.RVersion())

	arr.arr.obj.MarshalROOT(w)
	w.WriteString(arr.arr.name)
	w.WriteString(fmt.Sprintf("%s;%d", arr.cls, clsv))

	w.WriteI32(int32(len(arr.arr.objs)))
	w.WriteI32(arr.arr.low)

	switch {
	case arr.CanBypassStreamer():
		panic("rcont: writing TClonesArray with streamer by-pass not implemented")
	default:
		for i, obj := range arr.arr.objs {
			switch obj {
			case nil:
				w.WriteI8(0)
			default:
				w.WriteI8(1)
				if _, err := obj.(rbytes.Marshaler).MarshalROOT(w); err != nil {
					return 0, errors.Wrapf(err, "rcont: could not marshal TClonesArray element [%d/%d] (%T)", i+1, len(arr.arr.objs), obj)
				}
			}
		}
	}

	n, err := w.SetByteCount(pos, arr.Class())

	if bypass {
		arr.BypassStreamer(true)
	}

	return n, err
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (arr *ClonesArray) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(arr.Class())
	if vers > 2 {
		if err := arr.arr.obj.UnmarshalROOT(r); err != nil {
			return err
		}
	}
	if vers > 1 {
		arr.arr.name = r.ReadString()
	}
	clsv := r.ReadString()
	toks := strings.Split(clsv, ";")
	arr.cls = toks[0]
	clv, err := strconv.Atoi(toks[1])
	if err != nil {
		r.SetErr(errors.Wrapf(err, "rcont: could not extract TClonesArray element version"))
		return r.Err()
	}

	nobjs := int(r.ReadI32())
	if nobjs < 0 {
		nobjs = -nobjs
	}
	arr.arr.low = r.ReadI32()

	arr.arr.objs = make([]root.Object, nobjs)
	arr.arr.last = nobjs - 1
	si, err := r.StreamerInfo(arr.cls, clv)
	if err != nil {
		r.SetErr(errors.Wrapf(err, "rcont: could not find TClonesArray's element streamer %q and version=%d", arr.cls, clv))
		return r.Err()
	}
	fct := rtypes.Factory.Get(si.Name())

	switch {
	case arr.TestBits(rbytes.BypassStreamer) && !arr.TestBits(rbytes.CannotHandleMemberWiseStreaming):
		for i := range arr.arr.objs {
			obj := fct().Interface().(root.Object)
			arr.arr.objs[i] = obj
		}
		panic("rcont: TClonesArray with BypassStreamer not supported")
	default:
		for i := range arr.arr.objs {
			nch := r.ReadI8()
			if nch != 0 {
				obj := fct().Interface().(root.Object)
				if o, ok := obj.(rbytes.Unmarshaler); ok {
					err := o.UnmarshalROOT(r)
					if err != nil {
						return errors.Wrapf(err, "rcont: could not unmarshal TClonesArray element %d/%d", i+1, nobjs)
					}
				}
				if r.Err() != nil {
					return r.Err()
				}
				arr.arr.objs[i] = obj
			}
		}
	}

	r.CheckByteCount(pos, bcnt, start, arr.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := NewClonesArray()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TClonesArray", f)
}

var (
	_ root.Object        = (*ClonesArray)(nil)
	_ root.Named         = (*ClonesArray)(nil)
	_ root.ObjArray      = (*ClonesArray)(nil)
	_ rbytes.Marshaler   = (*ClonesArray)(nil)
	_ rbytes.Unmarshaler = (*ClonesArray)(nil)
)
