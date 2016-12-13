// Copyright 2016 The go-hep Authors.  All rights reserved.
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

// Scatter2D is a collection of 2-dim data points with errors.
type Scatter2D struct {
	pts []Point2D
	ann Annotation
}

// NewScatter2D creates a new 2-dim scatter with pts as an optional
// initial set of data points.
//
// If n <= 0, the initial capacity is 0.
func NewScatter2D(pts ...Point2D) *Scatter2D {
	s := &Scatter2D{
		pts: make([]Point2D, len(pts)),
		ann: make(Annotation),
	}
	copy(s.pts, pts)
	return s
}

// NewScatter2DFrom creates a new 2-dim scatter with x,y data slices.
//
// It panics if the lengths of the 2 slices don't match.
func NewScatter2DFrom(x, y []float64) *Scatter2D {
	if len(x) != len(y) {
		panic("hbook: len differ")
	}

	s := &Scatter2D{
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

// NewScatter2DFromH1D creates a new 2-dim scatter from the given H1D.
func NewScatter2DFromH1D(h *H1D) *Scatter2D {
	s := NewScatter2D()
	for k, v := range h.ann {
		s.ann[k] = v
	}
	for _, bin := range h.bng.bins {
		x := bin.XMid()
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

// Annotation returns the annotations attached to the
// scatter. (e.g. name, title, ...)
func (s *Scatter2D) Annotation() Annotation {
	return s.ann
}

// Name returns the name of this scatter
func (s *Scatter2D) Name() string {
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
func (*Scatter2D) Rank() int {
	return 2
}

// Entries returns the number of entries of this histogram.
func (s *Scatter2D) Entries() int64 {
	return int64(len(s.pts))
}

// Fill adds new points to the scatter.
func (s *Scatter2D) Fill(pts ...Point2D) {
	if len(pts) == 0 {
		return
	}

	i := len(s.pts)
	s.pts = append(s.pts, make([]Point2D, len(pts))...)
	copy(s.pts[i:], pts)
}

// Sort sorts the data points by x,y and x-err,y-err.
func (s *Scatter2D) Sort() {
	sort.Sort(points2D(s.pts))
}

// Points returns the points of the scatter.
//
// Users may not modify the returned slice.
// Users may not rely on the stability of the indices as the slice of points
// may be re-sorted at any point in time.
func (s *Scatter2D) Points() []Point2D {
	return s.pts
}

// Point returns the point at index i.
//
// Point panics if i is out of bounds.
func (s *Scatter2D) Point(i int) Point2D {
	return s.pts[i]
}

// ScaleX rescales the X values by a factor f.
func (s *Scatter2D) ScaleX(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleX(f)
	}
}

// ScaleY rescales the Y values by a factor f.
func (s *Scatter2D) ScaleY(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleY(f)
	}
}

// ScaleXY rescales the X and Y values by a factor f.
func (s *Scatter2D) ScaleXY(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleX(f)
		p.ScaleY(f)
	}
}

// Len returns the number of points in the scatter.
//
// Len implements the gonum/plot/plotter.XYer interface.
func (s *Scatter2D) Len() int {
	return len(s.pts)
}

// XY returns the x, y pair at index i.
//
// XY panics if i is out of bounds.
// XY implements the gonum/plot/plotter.XYer interface.
func (s *Scatter2D) XY(i int) (x, y float64) {
	pt := s.pts[i]
	x = pt.X
	y = pt.Y
	return
}

// XError returns the two error values for X data.
//
// XError implements the gonum/plot/plotter.XErrorer interface.
func (s *Scatter2D) XError(i int) (float64, float64) {
	pt := s.pts[i]
	return pt.ErrX.Min, pt.ErrX.Max
}

// YError returns the two error values for Y data.
//
// YError implements the gonum/plot/plotter.YErrorer interface.
func (s *Scatter2D) YError(i int) (float64, float64) {
	pt := s.pts[i]
	return pt.ErrY.Min, pt.ErrY.Max
}

// DataRange returns the minimum and maximum
// x and y values, implementing the gonum/plot.DataRanger
// interface.
func (s *Scatter2D) DataRange() (xmin, xmax, ymin, ymax float64) {
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

// annYODA creates a new Annotation with fields compatible with YODA
func (s *Scatter2D) annYODA() Annotation {
	ann := make(Annotation, len(s.ann))
	ann["Type"] = "Scatter2D"
	ann["Path"] = "/" + s.Name()
	ann["Title"] = ""
	for k, v := range s.ann {
		ann[k] = v
	}
	return ann
}

// MarshalYODA implements the YODAMarshaler interface.
func (s *Scatter2D) MarshalYODA() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := s.annYODA()
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
func (s *Scatter2D) UnmarshalYODA(data []byte) error {
	var err error
	var path string
	r := bytes.NewBuffer(data)
	_, err = fmt.Fscanf(r, "BEGIN YODA_SCATTER2D %s\n", &path)
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
	s.ann = ann
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
