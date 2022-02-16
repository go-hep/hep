// Copyright ©2017 The go-hep Authors. All rights reserved.
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

// The TNamed class is the base class for all named ROOT classes
// A TNamed contains the essential elements (name, title)
// to identify a derived object in containers, directories and files.
// Most member functions defined in this base class are in general
// overridden by the derived classes.
type Named struct {
	obj   Object
	name  string
	title string
}

func NewNamed(name, title string) *Named {
	return &Named{
		obj:   *NewObject(),
		name:  name,
		title: title,
	}
}

func (*Named) RVersion() int16 {
	return rvers.Named
}

// Name returns the name of the instance
func (n *Named) Name() string {
	return n.name
}

// Title returns the title of the instance
func (n *Named) Title() string {
	return n.title
}

func (n *Named) SetName(name string)   { n.name = name }
func (n *Named) SetTitle(title string) { n.title = title }

func (*Named) Class() string {
	return "TNamed"
}

func (n *Named) UID() uint32 {
	return n.obj.UID()
}

func (n *Named) Sizeof() int32 {
	return tstringSizeof(n.name) + tstringSizeof(n.title)
}

// tstringSizeof returns the size in bytes of the TString structure.
func tstringSizeof(v string) int32 {
	n := int32(len(v))
	if n > 254 {
		return n + 1 + 4
	}
	return n + 1
}

func (n *Named) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(n.Class())
	if hdr.Vers > rvers.Named {
		panic(fmt.Errorf("rbase: invalid named version=%d > %d", hdr.Vers, rvers.Named))
	}

	r.ReadObject(&n.obj)

	n.name = r.ReadString()
	n.title = r.ReadString()

	r.CheckHeader(hdr)
	return r.Err()
}

func (n *Named) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(n.Class(), n.RVersion())
	w.WriteObject(&n.obj)
	w.WriteString(n.name)
	w.WriteString(n.title)
	return w.SetHeader(hdr)
}

func init() {
	f := func() reflect.Value {
		o := NewNamed("", "")
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TNamed", f)
}

var (
	_ root.Object        = (*Named)(nil)
	_ root.UIDer         = (*Named)(nil)
	_ root.Named         = (*Named)(nil)
	_ rbytes.Marshaler   = (*Named)(nil)
	_ rbytes.Unmarshaler = (*Named)(nil)
)
