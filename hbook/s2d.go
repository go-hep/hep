// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"sort"
)

// S2D is a collection of 2-dim data points with errors.
type S2D struct {
	pts []Point2D
	ann Annotation
}

// NewS2D creates a new 2-dim scatter with pts as an optional
// initial set of data points.
//
// If n <= 0, the initial capacity is 0.
func NewS2D(pts ...Point2D) *S2D {
	s := &S2D{
		pts: make([]Point2D, len(pts)),
		ann: make(Annotation),
	}
	copy(s.pts, pts)
	return s
}

// NewS2DFrom creates a new 2-dim scatter with x,y data slices.
//
// It panics if the lengths of the 2 slices don't match.
func NewS2DFrom(x, y []float64) *S2D {
	if len(x) != len(y) {
		panic("hbook: len differ")
	}

	s := &S2D{
		pts: make([]Point2D, len(x)),
		ann: make(Annotation),
	}
	for i := range s.pts {
		pt := &s.pts[i]
		pt.X = x[i]
		pt.Y = y[i]
	}
	return s
}

// S2DOpts controls how S2D scatters are created from H1D and P1D.
type S2DOpts struct {
	UseFocus  bool
	UseStdDev bool
}

// NewS2DFromH1D creates a new 2-dim scatter from the given H1D.
// NewS2DFromH1D optionally takes a S2DOpts slice:
// only the first element is considered.
func NewS2DFromH1D(h *H1D, opts ...S2DOpts) *S2D {
	s := NewS2D()
	for k, v := range h.Ann {
		s.ann[k] = v
	}
	var opt S2DOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	// YODA support
	if _, ok := s.ann["Type"]; ok {
		s.ann["Type"] = "Scatter2D"
	}
	for _, bin := range h.Binning.Bins {
		var x float64
		if opt.UseFocus {
			x = bin.XFocus()
		} else {
			x = bin.XMid()
		}
		exm := x - bin.XMin()
		exp := bin.XMax() - x
		var y, ey float64
		if w := bin.XWidth(); w != 0 {
			ww := 1 / w
			y = bin.SumW() * ww
			ey = math.Sqrt(bin.SumW2()) * ww
		} else {
			y = math.NaN()  // FIXME(sbinet): use quiet-NaN ?
			ey = math.NaN() // FIXME(sbinet): use quiet-NaN ?
		}
		s.Fill(Point2D{X: x, Y: y, ErrX: Range{exm, exp}, ErrY: Range{ey, ey}})
	}
	return s
}

// NewS2DFromP1D creates a new 2-dim scatter from the given P1D.
// NewS2DFromP1D optionally takes a S2DOpts slice:
// only the first element is considered.
func NewS2DFromP1D(p *P1D, opts ...S2DOpts) *S2D {
	s := NewS2D()
	for k, v := range p.ann {
		p.ann[k] = v
	}
	var opt S2DOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	// YODA support
	if _, ok := s.ann["Type"]; ok {
		s.ann["Type"] = "Scatter2D"
	}
	for _, bin := range p.bng.bins {
		var x float64
		if opt.UseFocus {
			x = bin.XFocus()
		} else {
			x = bin.XMid()
		}
		exm := x - bin.XMin()
		exp := bin.XMax() - x
		var y, ey float64
		if w := bin.XWidth(); w != 0 {
			ww := 1 / w
			y = bin.SumW() * ww
			if opt.UseStdDev {
				ey = bin.XStdDev()
			} else {
				ey = bin.XStdErr()
			}
		} else {
			y = math.NaN()  // FIXME(sbinet): use quiet-NaN ?
			ey = math.NaN() // FIXME(sbinet): use quiet-NaN ?
		}
		s.Fill(Point2D{X: x, Y: y, ErrX: Range{exm, exp}, ErrY: Range{ey, ey}})
	}
	return s
}

// Annotation returns the annotations attached to the
// scatter. (e.g. name, title, ...)
func (s *S2D) Annotation() Annotation {
	return s.ann
}

// Name returns the name of this scatter
func (s *S2D) Name() string {
	v, ok := s.ann["name"]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

// Rank returns the number of dimensions of this scatter.
func (*S2D) Rank() int {
	return 2
}

// Entries returns the number of entries of this histogram.
func (s *S2D) Entries() int64 {
	return int64(len(s.pts))
}

// Fill adds new points to the scatter.
func (s *S2D) Fill(pts ...Point2D) {
	if len(pts) == 0 {
		return
	}

	i := len(s.pts)
	s.pts = append(s.pts, make([]Point2D, len(pts))...)
	copy(s.pts[i:], pts)
}

// Sort sorts the data points by x,y and x-err,y-err.
func (s *S2D) Sort() {
	sort.Sort(points2D(s.pts))
}

// Points returns the points of the scatter.
//
// Users may not modify the returned slice.
// Users may not rely on the stability of the indices as the slice of points
// may be re-sorted at any point in time.
func (s *S2D) Points() []Point2D {
	return s.pts
}

// Point returns the point at index i.
//
// Point panics if i is out of bounds.
func (s *S2D) Point(i int) Point2D {
	return s.pts[i]
}

// ScaleX rescales the X values by a factor f.
func (s *S2D) ScaleX(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleX(f)
	}
}

// ScaleY rescales the Y values by a factor f.
func (s *S2D) ScaleY(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleY(f)
	}
}

// ScaleXY rescales the X and Y values by a factor f.
func (s *S2D) ScaleXY(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleX(f)
		p.ScaleY(f)
	}
}

// Len returns the number of points in the scatter.
//
// Len implements the gonum/plot/plotter.XYer interface.
func (s *S2D) Len() int {
	return len(s.pts)
}

// XY returns the x, y pair at index i.
//
// XY panics if i is out of bounds.
// XY implements the gonum/plot/plotter.XYer interface.
func (s *S2D) XY(i int) (x, y float64) {
	pt := s.pts[i]
	x = pt.X
	y = pt.Y
	return
}

// XError returns the two error values for X data.
//
// XError implements the gonum/plot/plotter.XErrorer interface.
func (s *S2D) XError(i int) (float64, float64) {
	pt := s.pts[i]
	return pt.ErrX.Min, pt.ErrX.Max
}

// YError returns the two error values for Y data.
//
// YError implements the gonum/plot/plotter.YErrorer interface.
func (s *S2D) YError(i int) (float64, float64) {
	pt := s.pts[i]
	return pt.ErrY.Min, pt.ErrY.Max
}

// DataRange returns the minimum and maximum
// x and y values, implementing the gonum/plot.DataRanger
// interface.
func (s *S2D) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(+1)
	ymin = math.Inf(+1)
	xmax = math.Inf(-1)
	ymax = math.Inf(-1)
	for _, p := range s.pts {
		xmin = math.Min(p.XMin(), xmin)
		xmax = math.Max(p.XMax(), xmax)
		ymin = math.Min(p.YMin(), ymin)
		ymax = math.Max(p.YMax(), ymax)
	}
	return
}

// annToYODA creates a new Annotation with fields compatible with YODA
func (s *S2D) annToYODA() Annotation {
	ann := make(Annotation, len(s.ann))
	ann["Type"] = "Scatter2D"
	ann["Path"] = "/" + s.Name()
	ann["Title"] = ""
	for k, v := range s.ann {
		if k == "name" {
			continue
		}
		if k == "title" {
			ann["Title"] = v
			continue
		}
		ann[k] = v
	}
	return ann
}

// annFromYODA creates a new Annotation from YODA compatible fields
func (s *S2D) annFromYODA(ann Annotation) {
	if len(s.ann) == 0 {
		s.ann = make(Annotation, len(ann))
	}
	for k, v := range ann {
		switch k {
		case "Type":
			// noop
		case "Path":
			s.ann["name"] = string(v.(string)[1:]) // skip leading '/'
		case "Title":
			s.ann["title"] = v.(string)
		default:
			s.ann[k] = v
		}
	}
}

// MarshalYODA implements the YODAMarshaler interface.
func (s *S2D) MarshalYODA() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := s.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_SCATTER2D %s\n", ann["Path"])
	data, err := ann.MarshalYODA()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	// TODO: change ordering to {vals} {errs} {errs} ...
	fmt.Fprintf(buf, "# xval\t xerr-\t xerr+\t yval\t yerr-\t yerr+\n")
	s.Sort()
	for _, pt := range s.pts {
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\n",
			pt.X, pt.ErrX.Min, pt.ErrX.Max, pt.Y, pt.ErrY.Min, pt.ErrY.Max,
		)
	}
	fmt.Fprintf(buf, "END YODA_SCATTER2D\n\n")
	return buf.Bytes(), err
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (s *S2D) UnmarshalYODA(data []byte) error {
	r := bytes.NewBuffer(data)
	_, err := readYODAHeader(r, "BEGIN YODA_SCATTER2D")
	if err != nil {
		return err
	}
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# xval\t xerr-\t"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid Scatter2D-YODA data")
	}
	err = ann.UnmarshalYODA(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %v\nhbook: %q", err, string(r.Bytes()[:pos+1]))
	}
	s.annFromYODA(ann)
	r.Next(pos)

	sc := bufio.NewScanner(r)
scanLoop:
	for sc.Scan() {
		buf := sc.Bytes()
		if len(buf) == 0 || buf[0] == '#' {
			continue
		}
		rbuf := bytes.NewReader(buf)
		switch {
		case bytes.HasPrefix(buf, []byte("END YODA_SCATTER2D")):
			break scanLoop
		default:
			var pt Point2D
			fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\n",
				&pt.X, &pt.ErrX.Min, &pt.ErrX.Max, &pt.Y, &pt.ErrY.Min, &pt.ErrY.Max,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
			s.Fill(pt)
		}
	}
	err = sc.Err()
	if err == io.EOF {
		err = nil
	}
	s.Sort()
	return err
}
