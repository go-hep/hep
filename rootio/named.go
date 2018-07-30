// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

// The TNamed class is the base class for all named ROOT classes
// A TNamed contains the essential elements (name, title)
// to identify a derived object in containers, directories and files.
// Most member functions defined in this base class are in general
// overridden by the derived classes.
type tnamed struct {
	rvers int16
	obj   tobject
	name  string
	title string
}

// Name returns the name of the instance
func (n *tnamed) Name() string {
	return n.name
}

// Title returns the title of the instance
func (n *tnamed) Title() string {
	return n.title
}

func (n *tnamed) Class() string {
	return "TNamed"
}

func (n *tnamed) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	n.rvers = vers

	if err := n.obj.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	n.name = r.ReadString()
	n.title = r.ReadString()

	r.CheckByteCount(pos, bcnt, beg, "TNamed")
	return r.Err()
}

func (n *tnamed) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(n.rvers)
	if _, err := n.obj.MarshalROOT(w); err != nil {
		w.err = err
		return 0, w.err
	}

	w.WriteString(n.name)
	w.WriteString(n.title)

	return w.SetByteCount(pos, "TNamed")
}

func init() {
	f := func() reflect.Value {
		o := &tnamed{}
		return reflect.ValueOf(o)
	}
	Factory.add("TNamed", f)
	Factory.add("*rootio.tnamed", f)
}

var (
	_ Object          = (*tnamed)(nil)
	_ Named           = (*tnamed)(nil)
	_ ROOTMarshaler   = (*tnamed)(nil)
	_ ROOTUnmarshaler = (*tnamed)(nil)
)
