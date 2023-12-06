// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"bytes"
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/yodacnv"
)

type tmultigraph struct {
	rbase.Named

	graphs *rcont.List // Pointer to list of TGraphs
	funcs  *rcont.List // Pointer to list of functions (fits and user)
	histo  *H1F        // Pointer to histogram used for drawing axis
	ymax   float64     // Maximum value for plotting along y
	ymin   float64     // Minimum value for plotting along y
}

func newMultiGraph() *tmultigraph {
	return &tmultigraph{
		Named:  *rbase.NewNamed("", ""),
		graphs: rcont.NewList("", nil),
		funcs:  rcont.NewList("", nil),
	}
}

func (*tmultigraph) Class() string {
	return "TMultiGraph"
}

func (*tmultigraph) RVersion() int16 {
	return rvers.MultiGraph
}

func (mg *tmultigraph) Len() int {
	return mg.graphs.Len()
}

func (mg *tmultigraph) Graphs() []Graph {
	o := make([]Graph, mg.Len())
	for i := range o {
		o[i] = mg.graphs.At(i).(Graph)
	}
	return o
}

func (mg *tmultigraph) ROOTMerge(src root.Object) error {
	panic("not implemented")
}

// MarshalROOT implements rbytes.Marshaler
func (o *tmultigraph) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(o.Class(), o.RVersion())

	w.WriteObject(&o.Named)
	w.WriteObjectAny(o.graphs) // obj-ptr
	w.WriteObjectAny(o.funcs)  // obj-ptr
	w.WriteObjectAny(o.histo)  // obj-ptr
	w.WriteF64(o.ymax)
	w.WriteF64(o.ymin)

	return w.SetHeader(hdr)
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (o *tmultigraph) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(o.Class(), o.RVersion())

	r.ReadObject(&o.Named)
	{
		o.graphs = nil
		if oo := r.ReadObjectAny(); oo != nil { // obj-ptr
			o.graphs = oo.(*rcont.List)
		}
	}
	{
		o.funcs = nil
		if oo := r.ReadObjectAny(); oo != nil { // obj-ptr
			o.funcs = oo.(*rcont.List)
		}
	}
	{
		o.histo = nil
		if oo := r.ReadObjectAny(); oo != nil { // obj-ptr
			o.histo = oo.(*H1F)
		}
	}
	o.ymax = r.ReadF64()
	o.ymin = r.ReadF64()

	r.CheckHeader(hdr)
	return r.Err()
}

// MarshalYODA implements the YODAMarshaler interface.
func (mg *tmultigraph) MarshalYODA() ([]byte, error) {
	out := new(bytes.Buffer)
	for i := 0; i < mg.graphs.Len(); i++ {
		g := mg.graphs.At(i).(yodacnv.Marshaler)
		raw, err := g.MarshalYODA()
		if err != nil {
			return nil, fmt.Errorf("rhist: could not marshal multigraph %q: %w", mg.Name(), err)
		}
		_, _ = out.Write(raw)
	}
	return out.Bytes(), nil
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (mg *tmultigraph) UnmarshalYODA(raw []byte) error {
	objs, err := yodacnv.Read(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("rhist: could not unmarshal multigraph: %w", err)
	}
	for i, obj := range objs {
		s2, ok := obj.(*hbook.S2D)
		if !ok {
			return fmt.Errorf("rhist: could not unmarshal multigraph element #%d: got=%T, want=*hbook.S2D", i, obj)
		}
		mg.graphs.Append(NewGraphAsymmErrorsFrom(s2))
	}
	return nil
}

func (mg *tmultigraph) String() string {
	o, err := mg.MarshalYODA()
	if err != nil {
		panic(err)
	}
	return string(o)
}

func init() {
	f := func() reflect.Value {
		o := newMultiGraph()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TMultiGraph", f)
}

var (
	_ root.Object         = (*tmultigraph)(nil)
	_ root.Named          = (*tmultigraph)(nil)
	_ root.Merger         = (*tmultigraph)(nil)
	_ MultiGraph          = (*tmultigraph)(nil)
	_ rbytes.Marshaler    = (*tmultigraph)(nil)
	_ rbytes.Unmarshaler  = (*tmultigraph)(nil)
	_ yodacnv.Marshaler   = (*tmultigraph)(nil)
	_ yodacnv.Unmarshaler = (*tmultigraph)(nil)
)
