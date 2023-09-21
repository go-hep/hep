// Copyright ©2024 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"fmt"
	"math"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// Scatter implements ROOT's TScatter.
// A scatter plot able to draw four variables on a single plot.
type Scatter struct {
	rbase.Named
	attline   rbase.AttLine
	attfill   rbase.AttFill
	attmarker rbase.AttMarker

	npoints int32     // Number of points <= fMaxSize
	histo   *H2F      // Pointer to histogram used for drawing axis
	graph   *tgraph   // Pointer to graph holding X and Y positions
	color   []float64 // [fNpoints] array of colors
	size    []float64 // [fNpoints] array of marker sizes

	maxMarkerSize float64 // Largest marker size used to paint the markers
	minMarkerSize float64 // Smallest marker size used to paint the markers
	margin        float64 // Margin around the plot in %
}

func newScatter(n int) *Scatter {
	return &Scatter{
		Named:         *rbase.NewNamed("", ""),
		attline:       *rbase.NewAttLine(),
		attfill:       *rbase.NewAttFill(),
		attmarker:     *rbase.NewAttMarker(),
		npoints:       int32(n),
		color:         make([]float64, n),
		size:          make([]float64, n),
		maxMarkerSize: 5,
		minMarkerSize: 1,
		margin:        0.1,
	}
}

func (*Scatter) RVersion() int16 {
	return rvers.Scatter
}

func (*Scatter) Class() string {
	return "TScatter"
}

func (s *Scatter) ROOTMerge(src root.Object) error {
	switch src := src.(type) {
	case *Scatter:
		var err error
		s.npoints += src.npoints
		// FIXME(sbinet): implement ROOTMerge for TH2x
		//	err = s.histo.ROOTMerge(src.histo)
		//	if err != nil {
		//		return fmt.Errorf("rhist: could not merge Scatter's underlying H2F: %w", err)
		//	}
		err = s.graph.ROOTMerge(src.graph)
		if err != nil {
			return fmt.Errorf("rhist: could not merge Scatter's underlying Graph: %w", err)
		}
		s.color = append(s.color, src.color...)
		s.size = append(s.size, src.size...)
		s.maxMarkerSize = math.Max(s.maxMarkerSize, src.maxMarkerSize)
		s.minMarkerSize = math.Min(s.minMarkerSize, src.minMarkerSize)
		// FIXME(sbinet): handle margin
		return nil
	default:
		return fmt.Errorf("rhist: can not merge %T into %T", src, s)
	}
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (s *Scatter) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(s.Class(), s.RVersion())

	w.WriteObject(&s.Named)
	w.WriteObject(&s.attline)
	w.WriteObject(&s.attfill)
	w.WriteObject(&s.attmarker)

	w.WriteI32(s.npoints)
	w.WriteObjectAny(s.histo)
	w.WriteObjectAny(s.graph)

	w.WriteI8(1)
	w.WriteArrayF64(s.color)
	w.WriteI8(1)
	w.WriteArrayF64(s.size)

	w.WriteF64(s.maxMarkerSize)
	w.WriteF64(s.minMarkerSize)
	w.WriteF64(s.margin)

	return w.SetHeader(hdr)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (s *Scatter) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(s.Class(), s.RVersion())

	r.ReadObject(&s.Named)
	r.ReadObject(&s.attline)
	r.ReadObject(&s.attfill)
	r.ReadObject(&s.attmarker)

	s.npoints = r.ReadI32()
	if hdr.Vers < 2 {
		r.SetErr(fmt.Errorf("rhist: invalid TScatter version %d", hdr.Vers))
		return r.Err()
	}

	histo := r.ReadObjectAny()
	if histo != nil {
		s.histo = histo.(*H2F)
	}
	graph := r.ReadObjectAny()
	if graph != nil {
		s.graph = graph.(*tgraph)
	}

	_ = r.ReadI8()
	s.color = make([]float64, s.npoints)
	r.ReadArrayF64(s.color)
	_ = r.ReadI8()
	s.size = make([]float64, s.npoints)
	r.ReadArrayF64(s.size)

	s.maxMarkerSize = r.ReadF64()
	s.minMarkerSize = r.ReadF64()
	s.margin = r.ReadF64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (g *Scatter) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, g.Named.RMembers()...)
	mbrs = append(mbrs, g.attline.RMembers()...)
	mbrs = append(mbrs, g.attfill.RMembers()...)
	mbrs = append(mbrs, g.attmarker.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fNpoints", Value: &g.npoints},
		{Name: "fHistogram", Value: &g.histo},
		{Name: "fGraph", Value: &g.graph},
		{Name: "fColor", Value: &g.color},
		{Name: "fSize", Value: &g.size},
		{Name: "fMaxMarkerSize", Value: &g.maxMarkerSize},
		{Name: "fMinMarkerSize", Value: &g.minMarkerSize},
		{Name: "fMargin", Value: &g.margin},
	}...)

	return mbrs
}

func init() {
	{
		f := func() reflect.Value {
			o := newScatter(0)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TScatter", f)
	}
}

var (
	_ root.Object        = (*Scatter)(nil)
	_ root.Named         = (*Scatter)(nil)
	_ root.Merger        = (*Scatter)(nil)
	_ rbytes.Marshaler   = (*Scatter)(nil)
	_ rbytes.Unmarshaler = (*Scatter)(nil)
	_ rbytes.RSlicer     = (*Scatter)(nil)
)
