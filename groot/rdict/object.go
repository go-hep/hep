// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

var (
	cxxNameSanitizer = strings.NewReplacer(
		"<", "_",
		">", "_",
		":", "_",
		",", "_",
		" ", "_",
	)
)

func ObjectFrom(si rbytes.StreamerInfo, sictx rbytes.StreamerInfoContext) *Object {
	return newObjectFrom(si, sictx)
}

// Object wraps a type created from a Streamer and implements the
// following interfaces:
//   - root.Object
//   - rbytes.RVersioner
//   - rbytes.Marshaler
//   - rbytes.Unmarshaler
type Object struct {
	v interface{}

	si    *StreamerInfo
	rvers int16
	class string
}

func (obj *Object) Class() string {
	return obj.class
}

func (obj *Object) SetClass(name string) {
	si, ok := StreamerInfos.Get(name, -1)
	if !ok {
		panic(fmt.Errorf("rdict: no streamer for %q", name))
	}
	*obj = *newObjectFrom(si, StreamerInfos)
}

func (obj *Object) String() string {
	return fmt.Sprintf("%v", obj.v)
}

func (obj *Object) RVersion() int16 {
	return obj.rvers
}

func (obj *Object) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	dec, err := obj.si.NewDecoder(rbytes.ObjectWise, r)
	if err != nil {
		return fmt.Errorf("rdict: could not create decoder for %q: %w", obj.si.Name(), err)
	}

	err = dec.DecodeROOT(obj.v)
	if err != nil {
		return fmt.Errorf("rdict: could not decode %q: %w", obj.si.Name(), err)
	}

	return r.Err()
}

func (obj *Object) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	enc, err := obj.si.NewEncoder(rbytes.ObjectWise, w)
	if err != nil {
		return 0, fmt.Errorf("rdict: could not create encoder for %q: %w", obj.si.Name(), err)
	}

	pos := w.Pos()

	err = enc.EncodeROOT(obj.v)
	if err != nil {
		return 0, fmt.Errorf("rdict: could not encode %q: %w", obj.si.Name(), err)
	}

	return int(w.Pos() - pos), w.Err()
}

func newObjectFrom(si rbytes.StreamerInfo, sictx rbytes.StreamerInfoContext) *Object {
	err := si.BuildStreamers()
	if err != nil {
		panic(fmt.Errorf("rdict: could not build streamers: %w", err))
	}

	rt, err := TypeFromSI(sictx, si)
	if err != nil {
		panic(fmt.Errorf("rdict: could not build object type: %w", err))
	}

	recv := reflect.New(rt)
	obj := &Object{
		v:     recv.Interface(),
		si:    si.(*StreamerInfo),
		rvers: int16(si.ClassVersion()),
		class: si.Name(),
	}
	return obj
}

var (
	_ root.Object        = (*Object)(nil)
	_ rbytes.RVersioner  = (*Object)(nil)
	_ rbytes.Marshaler   = (*Object)(nil)
	_ rbytes.Unmarshaler = (*Object)(nil)
)

func init() {
	{
		f := func() reflect.Value {
			o := &Object{class: "*rdict.Object"}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("*rdict.Object", f)
	}
}
