// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"bytes"
	"fmt"
	"math"
	"reflect"

	"go-hep.org/x/hep/hbook"
)

// H2F implements ROOT TH2F
type H2F struct {
	rvers int16
	th2
	arr ArrayF
}

func (*H2F) isH2() {}

// Class returns the ROOT class name.
func (*H2F) Class() string {
	return "TH2F"
}

func (h *H2F) Array() ArrayF {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H2F) Rank() int {
	return 2
}

// NbinsX returns the number of bins in X.
func (h *H2F) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H2F) XAxis() Axis {
	return &h.th1.xaxis
}

// XBinCenter returns the bin center value in X.
func (h *H2F) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *H2F) XBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// XBinError returns the bin error in X.
func (h *H2F) XBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *H2F) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *H2F) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

// NbinsY returns the number of bins in Y.
func (h *H2F) NbinsY() int {
	return h.th1.yaxis.nbins
}

// YAxis returns the axis along Y.
func (h *H2F) YAxis() Axis {
	return &h.th1.yaxis
}

// YBinCenter returns the bin center value in Y.
func (h *H2F) YBinCenter(i int) float64 {
	return float64(h.th1.yaxis.BinCenter(i))
}

// YBinContent returns the bin content value in Y.
func (h *H2F) YBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// YBinError returns the bin error in Y.
func (h *H2F) YBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// YBinLowEdge returns the bin lower edge value in Y.
func (h *H2F) YBinLowEdge(i int) float64 {
	return h.th1.yaxis.BinLowEdge(i)
}

// YBinWidth returns the bin width in Y.
func (h *H2F) YBinWidth(i int) float64 {
	return h.th1.yaxis.BinWidth(i)
}

// bin returns the regularized bin number given an (x,y) bin index pair.
func (h *H2F) bin(ix, iy int) int {
	nx := h.th1.xaxis.nbins + 1 // overflow bin
	ny := h.th1.yaxis.nbins + 1 // overflow bin
	switch {
	case ix < 0:
		ix = 0
	case ix > nx:
		ix = nx
	}
	switch {
	case iy < 0:
		iy = 0
	case iy > ny:
		iy = ny
	}
	return ix + (nx+1)*iy
}

func (h *H2F) dist2D(ix, iy int) hbook.Dist2D {
	i := h.bin(ix, iy)
	vx := h.XBinContent(i)
	xerr := h.XBinError(i)
	nx := h.entries(vx, xerr)
	vy := h.YBinContent(i)
	yerr := h.YBinError(i)
	ny := h.entries(vy, yerr)

	sumw := h.arr.Data[i]
	sumw2 := 0.0
	if len(h.th1.sumw2.Data) > 0 {
		sumw2 = h.th1.sumw2.Data[i]
	}
	return hbook.Dist2D{
		X: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     nx,
				SumW:  float64(sumw),
				SumW2: float64(sumw2),
			},
		},
		Y: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     ny,
				SumW:  float64(sumw),
				SumW2: float64(sumw2),
			},
		},
	}
}

func (h *H2F) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

func (h *H2F) MarshalYODA() ([]byte, error) {
	var (
		nx       = h.NbinsX()
		ny       = h.NbinsY()
		xinrange = 1
		yinrange = 1
		dflow    = [8]hbook.Dist2D{
			h.dist2D(0, 0),
			h.dist2D(0, yinrange),
			h.dist2D(0, ny+1),
			h.dist2D(nx+1, 0),
			h.dist2D(nx+1, yinrange),
			h.dist2D(nx+1, ny+1),
			h.dist2D(xinrange, 0),
			h.dist2D(xinrange, ny+1),
		}
		dtot = hbook.Dist2D{
			X: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:     int64(h.Entries()),
					SumW:  float64(h.SumW()),
					SumW2: float64(h.SumW2()),
				},
				SumWX:  float64(h.SumWX()),
				SumWX2: float64(h.SumWX2()),
			},
			Y: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:     int64(h.Entries()),
					SumW:  float64(h.SumW()),
					SumW2: float64(h.SumW2()),
				},
				SumWX:  float64(h.SumWY()),
				SumWX2: float64(h.SumWY2()),
			},
			SumWXY: h.SumWXY(),
		}
		dists = make([]hbook.Dist2D, int(nx*ny))
	)
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			i := iy*nx + ix
			dists[i] = h.dist2D(ix+1, iy+1)
		}
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "BEGIN YODA_HISTO2D /%s\n", h.Name())
	fmt.Fprintf(buf, "Path=/%s\n", h.Name())
	fmt.Fprintf(buf, "Title=%s\n", h.Title())
	fmt.Fprintf(buf, "Type=Histo2D\n")
	fmt.Fprintf(buf, "# Mean: %e\n", math.NaN())
	fmt.Fprintf(buf, "# Volume: %e\n", math.NaN())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")

	var name = "Total   "
	d := &dtot
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
	)

	if false { // FIXME(sbinet)
		for _, d := range dflow {
			fmt.Fprintf(
				buf,
				"%s\t%s\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				name, name,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
			)

		}
	} else {
		// outflows
		fmt.Fprintf(buf, "# 2D outflow persistency not currently supported until API is stable\n")
	}

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t ylow\t yhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			xmin := h.XBinLowEdge(ix + 1)
			xmax := h.XBinWidth(ix+1) + xmin
			ymin := h.YBinLowEdge(iy + 1)
			ymax := h.YBinWidth(iy+1) + ymin
			i := iy*nx + ix
			d := &dists[i]
			fmt.Fprintf(
				buf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				xmin, xmax, ymin, ymax,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
			)
		}
	}
	fmt.Fprintf(buf, "END YODA_HISTO2D\n\n")
	return buf.Bytes(), nil
}

func (h *H2F) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(h.rvers)

	for _, v := range []ROOTMarshaler{
		&h.th2,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	return w.SetByteCount(pos, "TH2F")
}

func (h *H2F) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: TH2F version too old (%d<1)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH2F")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &H2F{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH2F", f)
	Factory.add("*rootio.H2F", f)
}

var (
	_ Object          = (*H2F)(nil)
	_ Named           = (*H2F)(nil)
	_ H2              = (*H2F)(nil)
	_ ROOTMarshaler   = (*H2F)(nil)
	_ ROOTUnmarshaler = (*H2F)(nil)
)

// H2D implements ROOT TH2D
type H2D struct {
	rvers int16
	th2
	arr ArrayD
}

func (*H2D) isH2() {}

// Class returns the ROOT class name.
func (*H2D) Class() string {
	return "TH2D"
}

func (h *H2D) Array() ArrayD {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H2D) Rank() int {
	return 2
}

// NbinsX returns the number of bins in X.
func (h *H2D) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H2D) XAxis() Axis {
	return &h.th1.xaxis
}

// XBinCenter returns the bin center value in X.
func (h *H2D) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *H2D) XBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// XBinError returns the bin error in X.
func (h *H2D) XBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *H2D) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *H2D) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

// NbinsY returns the number of bins in Y.
func (h *H2D) NbinsY() int {
	return h.th1.yaxis.nbins
}

// YAxis returns the axis along Y.
func (h *H2D) YAxis() Axis {
	return &h.th1.yaxis
}

// YBinCenter returns the bin center value in Y.
func (h *H2D) YBinCenter(i int) float64 {
	return float64(h.th1.yaxis.BinCenter(i))
}

// YBinContent returns the bin content value in Y.
func (h *H2D) YBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// YBinError returns the bin error in Y.
func (h *H2D) YBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// YBinLowEdge returns the bin lower edge value in Y.
func (h *H2D) YBinLowEdge(i int) float64 {
	return h.th1.yaxis.BinLowEdge(i)
}

// YBinWidth returns the bin width in Y.
func (h *H2D) YBinWidth(i int) float64 {
	return h.th1.yaxis.BinWidth(i)
}

// bin returns the regularized bin number given an (x,y) bin index pair.
func (h *H2D) bin(ix, iy int) int {
	nx := h.th1.xaxis.nbins + 1 // overflow bin
	ny := h.th1.yaxis.nbins + 1 // overflow bin
	switch {
	case ix < 0:
		ix = 0
	case ix > nx:
		ix = nx
	}
	switch {
	case iy < 0:
		iy = 0
	case iy > ny:
		iy = ny
	}
	return ix + (nx+1)*iy
}

func (h *H2D) dist2D(ix, iy int) hbook.Dist2D {
	i := h.bin(ix, iy)
	vx := h.XBinContent(i)
	xerr := h.XBinError(i)
	nx := h.entries(vx, xerr)
	vy := h.YBinContent(i)
	yerr := h.YBinError(i)
	ny := h.entries(vy, yerr)

	sumw := h.arr.Data[i]
	sumw2 := 0.0
	if len(h.th1.sumw2.Data) > 0 {
		sumw2 = h.th1.sumw2.Data[i]
	}
	return hbook.Dist2D{
		X: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     nx,
				SumW:  float64(sumw),
				SumW2: float64(sumw2),
			},
		},
		Y: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     ny,
				SumW:  float64(sumw),
				SumW2: float64(sumw2),
			},
		},
	}
}

func (h *H2D) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

func (h *H2D) MarshalYODA() ([]byte, error) {
	var (
		nx       = h.NbinsX()
		ny       = h.NbinsY()
		xinrange = 1
		yinrange = 1
		dflow    = [8]hbook.Dist2D{
			h.dist2D(0, 0),
			h.dist2D(0, yinrange),
			h.dist2D(0, ny+1),
			h.dist2D(nx+1, 0),
			h.dist2D(nx+1, yinrange),
			h.dist2D(nx+1, ny+1),
			h.dist2D(xinrange, 0),
			h.dist2D(xinrange, ny+1),
		}
		dtot = hbook.Dist2D{
			X: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:     int64(h.Entries()),
					SumW:  float64(h.SumW()),
					SumW2: float64(h.SumW2()),
				},
				SumWX:  float64(h.SumWX()),
				SumWX2: float64(h.SumWX2()),
			},
			Y: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:     int64(h.Entries()),
					SumW:  float64(h.SumW()),
					SumW2: float64(h.SumW2()),
				},
				SumWX:  float64(h.SumWY()),
				SumWX2: float64(h.SumWY2()),
			},
			SumWXY: h.SumWXY(),
		}
		dists = make([]hbook.Dist2D, int(nx*ny))
	)
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			i := iy*nx + ix
			dists[i] = h.dist2D(ix+1, iy+1)
		}
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "BEGIN YODA_HISTO2D /%s\n", h.Name())
	fmt.Fprintf(buf, "Path=/%s\n", h.Name())
	fmt.Fprintf(buf, "Title=%s\n", h.Title())
	fmt.Fprintf(buf, "Type=Histo2D\n")
	fmt.Fprintf(buf, "# Mean: %e\n", math.NaN())
	fmt.Fprintf(buf, "# Volume: %e\n", math.NaN())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")

	var name = "Total   "
	d := &dtot
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
	)

	if false { // FIXME(sbinet)
		for _, d := range dflow {
			fmt.Fprintf(
				buf,
				"%s\t%s\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				name, name,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
			)

		}
	} else {
		// outflows
		fmt.Fprintf(buf, "# 2D outflow persistency not currently supported until API is stable\n")
	}

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t ylow\t yhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			xmin := h.XBinLowEdge(ix + 1)
			xmax := h.XBinWidth(ix+1) + xmin
			ymin := h.YBinLowEdge(iy + 1)
			ymax := h.YBinWidth(iy+1) + ymin
			i := iy*nx + ix
			d := &dists[i]
			fmt.Fprintf(
				buf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				xmin, xmax, ymin, ymax,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
			)
		}
	}
	fmt.Fprintf(buf, "END YODA_HISTO2D\n\n")
	return buf.Bytes(), nil
}

func (h *H2D) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(h.rvers)

	for _, v := range []ROOTMarshaler{
		&h.th2,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	return w.SetByteCount(pos, "TH2D")
}

func (h *H2D) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: TH2D version too old (%d<1)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH2D")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &H2D{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH2D", f)
	Factory.add("*rootio.H2D", f)
}

var (
	_ Object          = (*H2D)(nil)
	_ Named           = (*H2D)(nil)
	_ H2              = (*H2D)(nil)
	_ ROOTMarshaler   = (*H2D)(nil)
	_ ROOTUnmarshaler = (*H2D)(nil)
)

// H2I implements ROOT TH2I
type H2I struct {
	rvers int16
	th2
	arr ArrayI
}

func (*H2I) isH2() {}

// Class returns the ROOT class name.
func (*H2I) Class() string {
	return "TH2I"
}

func (h *H2I) Array() ArrayI {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H2I) Rank() int {
	return 2
}

// NbinsX returns the number of bins in X.
func (h *H2I) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H2I) XAxis() Axis {
	return &h.th1.xaxis
}

// XBinCenter returns the bin center value in X.
func (h *H2I) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *H2I) XBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// XBinError returns the bin error in X.
func (h *H2I) XBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *H2I) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *H2I) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

// NbinsY returns the number of bins in Y.
func (h *H2I) NbinsY() int {
	return h.th1.yaxis.nbins
}

// YAxis returns the axis along Y.
func (h *H2I) YAxis() Axis {
	return &h.th1.yaxis
}

// YBinCenter returns the bin center value in Y.
func (h *H2I) YBinCenter(i int) float64 {
	return float64(h.th1.yaxis.BinCenter(i))
}

// YBinContent returns the bin content value in Y.
func (h *H2I) YBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// YBinError returns the bin error in Y.
func (h *H2I) YBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// YBinLowEdge returns the bin lower edge value in Y.
func (h *H2I) YBinLowEdge(i int) float64 {
	return h.th1.yaxis.BinLowEdge(i)
}

// YBinWidth returns the bin width in Y.
func (h *H2I) YBinWidth(i int) float64 {
	return h.th1.yaxis.BinWidth(i)
}

// bin returns the regularized bin number given an (x,y) bin index pair.
func (h *H2I) bin(ix, iy int) int {
	nx := h.th1.xaxis.nbins + 1 // overflow bin
	ny := h.th1.yaxis.nbins + 1 // overflow bin
	switch {
	case ix < 0:
		ix = 0
	case ix > nx:
		ix = nx
	}
	switch {
	case iy < 0:
		iy = 0
	case iy > ny:
		iy = ny
	}
	return ix + (nx+1)*iy
}

func (h *H2I) dist2D(ix, iy int) hbook.Dist2D {
	i := h.bin(ix, iy)
	vx := h.XBinContent(i)
	xerr := h.XBinError(i)
	nx := h.entries(vx, xerr)
	vy := h.YBinContent(i)
	yerr := h.YBinError(i)
	ny := h.entries(vy, yerr)

	sumw := h.arr.Data[i]
	sumw2 := 0.0
	if len(h.th1.sumw2.Data) > 0 {
		sumw2 = h.th1.sumw2.Data[i]
	}
	return hbook.Dist2D{
		X: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     nx,
				SumW:  float64(sumw),
				SumW2: float64(sumw2),
			},
		},
		Y: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     ny,
				SumW:  float64(sumw),
				SumW2: float64(sumw2),
			},
		},
	}
}

func (h *H2I) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

func (h *H2I) MarshalYODA() ([]byte, error) {
	var (
		nx       = h.NbinsX()
		ny       = h.NbinsY()
		xinrange = 1
		yinrange = 1
		dflow    = [8]hbook.Dist2D{
			h.dist2D(0, 0),
			h.dist2D(0, yinrange),
			h.dist2D(0, ny+1),
			h.dist2D(nx+1, 0),
			h.dist2D(nx+1, yinrange),
			h.dist2D(nx+1, ny+1),
			h.dist2D(xinrange, 0),
			h.dist2D(xinrange, ny+1),
		}
		dtot = hbook.Dist2D{
			X: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:     int64(h.Entries()),
					SumW:  float64(h.SumW()),
					SumW2: float64(h.SumW2()),
				},
				SumWX:  float64(h.SumWX()),
				SumWX2: float64(h.SumWX2()),
			},
			Y: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:     int64(h.Entries()),
					SumW:  float64(h.SumW()),
					SumW2: float64(h.SumW2()),
				},
				SumWX:  float64(h.SumWY()),
				SumWX2: float64(h.SumWY2()),
			},
			SumWXY: h.SumWXY(),
		}
		dists = make([]hbook.Dist2D, int(nx*ny))
	)
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			i := iy*nx + ix
			dists[i] = h.dist2D(ix+1, iy+1)
		}
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "BEGIN YODA_HISTO2D /%s\n", h.Name())
	fmt.Fprintf(buf, "Path=/%s\n", h.Name())
	fmt.Fprintf(buf, "Title=%s\n", h.Title())
	fmt.Fprintf(buf, "Type=Histo2D\n")
	fmt.Fprintf(buf, "# Mean: %e\n", math.NaN())
	fmt.Fprintf(buf, "# Volume: %e\n", math.NaN())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")

	var name = "Total   "
	d := &dtot
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
	)

	if false { // FIXME(sbinet)
		for _, d := range dflow {
			fmt.Fprintf(
				buf,
				"%s\t%s\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				name, name,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
			)

		}
	} else {
		// outflows
		fmt.Fprintf(buf, "# 2D outflow persistency not currently supported until API is stable\n")
	}

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t ylow\t yhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			xmin := h.XBinLowEdge(ix + 1)
			xmax := h.XBinWidth(ix+1) + xmin
			ymin := h.YBinLowEdge(iy + 1)
			ymax := h.YBinWidth(iy+1) + ymin
			i := iy*nx + ix
			d := &dists[i]
			fmt.Fprintf(
				buf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				xmin, xmax, ymin, ymax,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY, d.Entries(),
			)
		}
	}
	fmt.Fprintf(buf, "END YODA_HISTO2D\n\n")
	return buf.Bytes(), nil
}

func (h *H2I) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(h.rvers)

	for _, v := range []ROOTMarshaler{
		&h.th2,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	return w.SetByteCount(pos, "TH2I")
}

func (h *H2I) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: TH2I version too old (%d<1)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH2I")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &H2I{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH2I", f)
	Factory.add("*rootio.H2I", f)
}

var (
	_ Object          = (*H2I)(nil)
	_ Named           = (*H2I)(nil)
	_ H2              = (*H2I)(nil)
	_ ROOTMarshaler   = (*H2I)(nil)
	_ ROOTUnmarshaler = (*H2I)(nil)
)
