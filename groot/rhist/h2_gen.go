// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rhist

import (
	"fmt"
	"math"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"go-hep.org/x/hep/hbook"
)

// H2F implements ROOT TH2F
type H2F struct {
	th2
	arr rcont.ArrayF
}

func newH2F() *H2F {
	return &H2F{
		th2: *newH2(),
	}
}

// NewH2FFrom creates a new H2F from hbook 2-dim histogram.
func NewH2FFrom(h *hbook.H2D) *H2F {
	var (
		hroot  = newH2F()
		bins   = h.Binning.Bins
		nxbins = h.Binning.Nx
		nybins = h.Binning.Ny
		xedges = make([]float64, 0, nxbins+1)
		yedges = make([]float64, 0, nybins+1)
	)

	hroot.th2.th1.entries = float64(h.Entries())
	hroot.th2.th1.tsumw = h.SumW()
	hroot.th2.th1.tsumw2 = h.SumW2()
	hroot.th2.th1.tsumwx = h.SumWX()
	hroot.th2.th1.tsumwx2 = h.SumWX2()
	hroot.th2.tsumwy = h.SumWY()
	hroot.th2.tsumwy2 = h.SumWY2()
	hroot.th2.tsumwxy = h.SumWXY()

	ncells := (nxbins + 2) * (nybins + 2)
	hroot.th2.th1.ncells = ncells

	hroot.th2.th1.xaxis.nbins = nxbins
	hroot.th2.th1.xaxis.xmin = h.XMin()
	hroot.th2.th1.xaxis.xmax = h.XMax()

	hroot.th2.th1.yaxis.nbins = nybins
	hroot.th2.th1.yaxis.xmin = h.YMin()
	hroot.th2.th1.yaxis.xmax = h.YMax()

	hroot.arr.Data = make([]float32, ncells)
	hroot.th2.th1.sumw2.Data = make([]float64, ncells)

	ibin := func(ix, iy int) int { return iy*nxbins + ix }

	for ix := 0; ix < h.Binning.Nx; ix++ {
		for iy := 0; iy < h.Binning.Ny; iy++ {
			i := ibin(ix, iy)
			bin := bins[i]
			if ix == 0 {
				yedges = append(yedges, bin.YMin())
			}
			if iy == 0 {
				xedges = append(xedges, bin.XMin())
			}
			hroot.setDist2D(ix+1, iy+1, bin.Dist.SumW(), bin.Dist.SumW2())
		}
	}

	oflows := h.Binning.Outflows[:]
	for i, v := range []struct{ ix, iy int }{
		{0, 0},
		{0, 1},
		{0, nybins + 1},
		{nxbins + 1, 0},
		{nxbins + 1, 1},
		{nxbins + 1, nybins + 1},
		{1, 0},
		{1, nybins + 1},
	} {
		hroot.setDist2D(v.ix, v.iy, oflows[i].SumW(), oflows[i].SumW2())
	}

	xedges = append(xedges, bins[ibin(h.Binning.Nx-1, 0)].XMax())
	yedges = append(yedges, bins[ibin(0, h.Binning.Ny-1)].YMax())

	hroot.th2.th1.SetName(h.Name())
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th2.th1.SetTitle(v.(string))
	}
	hroot.th2.th1.xaxis.xbins.Data = xedges
	hroot.th2.th1.yaxis.xbins.Data = yedges

	return hroot
}

func (*H2F) RVersion() int16 {
	return rvers.H2F
}

func (*H2F) isH2() {}

// Class returns the ROOT class name.
func (*H2F) Class() string {
	return "TH2F"
}

func (h *H2F) Array() rcont.ArrayF {
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

func (h *H2F) setDist2D(ix, iy int, sumw, sumw2 float64) {
	i := h.bin(ix, iy)
	h.arr.Data[i] = float32(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *H2F) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

// AsH2D creates a new hbook.H2D from this ROOT histogram.
func (h *H2F) AsH2D() *hbook.H2D {
	var (
		nx = h.NbinsX()
		ny = h.NbinsY()
		hh = hbook.NewH2D(
			nx, h.XAxis().XMin(), h.XAxis().XMax(),
			ny, h.YAxis().XMin(), h.YAxis().XMax(),
		)
		xinrange = 1
		yinrange = 1
	)
	hh.Ann = hbook.Annotation{
		"name":  h.Name(),
		"title": h.Title(),
	}
	hh.Binning.Outflows = [8]hbook.Dist2D{
		h.dist2D(0, 0),
		h.dist2D(0, yinrange),
		h.dist2D(0, ny+1),
		h.dist2D(nx+1, 0),
		h.dist2D(nx+1, yinrange),
		h.dist2D(nx+1, ny+1),
		h.dist2D(xinrange, 0),
		h.dist2D(xinrange, ny+1),
	}

	hh.Binning.Dist = hbook.Dist2D{
		X: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		},
		Y: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		},
	}
	hh.Binning.Dist.X.Stats.SumWX = float64(h.SumWX())
	hh.Binning.Dist.X.Stats.SumWX2 = float64(h.SumWX2())
	hh.Binning.Dist.Y.Stats.SumWX = float64(h.SumWY())
	hh.Binning.Dist.Y.Stats.SumWX2 = float64(h.SumWY2())
	hh.Binning.Dist.Stats.SumWXY = h.SumWXY()

	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			var (
				i    = iy*nx + ix
				xmin = h.XBinLowEdge(ix + 1)
				xmax = h.XBinWidth(ix+1) + xmin
				ymin = h.YBinLowEdge(iy + 1)
				ymax = h.YBinWidth(iy+1) + ymin
				bin  = &hh.Binning.Bins[i]
			)
			bin.XRange.Min = xmin
			bin.XRange.Max = xmax
			bin.YRange.Min = ymin
			bin.YRange.Max = ymax
			bin.Dist = h.dist2D(ix+1, iy+1)
		}
	}

	return hh
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H2F) MarshalYODA() ([]byte, error) {
	return h.AsH2D().MarshalYODA()
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H2F) UnmarshalYODA(raw []byte) error {
	var hh hbook.H2D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *NewH2FFrom(&hh)
	return nil
}

func (h *H2F) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())

	for _, v := range []rbytes.Marshaler{
		&h.th2,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	return w.SetByteCount(pos, h.Class())
}

func (h *H2F) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers < 1 {
		return fmt.Errorf("rhist: TH2F version too old (%d<1)", vers)
	}

	for _, v := range []rbytes.Unmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := newH2F()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TH2F", f)
}

var (
	_ root.Object        = (*H2F)(nil)
	_ root.Named         = (*H2F)(nil)
	_ H2                 = (*H2F)(nil)
	_ rbytes.Marshaler   = (*H2F)(nil)
	_ rbytes.Unmarshaler = (*H2F)(nil)
)

// H2D implements ROOT TH2D
type H2D struct {
	th2
	arr rcont.ArrayD
}

func newH2D() *H2D {
	return &H2D{
		th2: *newH2(),
	}
}

// NewH2DFrom creates a new H2D from hbook 2-dim histogram.
func NewH2DFrom(h *hbook.H2D) *H2D {
	var (
		hroot  = newH2D()
		bins   = h.Binning.Bins
		nxbins = h.Binning.Nx
		nybins = h.Binning.Ny
		xedges = make([]float64, 0, nxbins+1)
		yedges = make([]float64, 0, nybins+1)
	)

	hroot.th2.th1.entries = float64(h.Entries())
	hroot.th2.th1.tsumw = h.SumW()
	hroot.th2.th1.tsumw2 = h.SumW2()
	hroot.th2.th1.tsumwx = h.SumWX()
	hroot.th2.th1.tsumwx2 = h.SumWX2()
	hroot.th2.tsumwy = h.SumWY()
	hroot.th2.tsumwy2 = h.SumWY2()
	hroot.th2.tsumwxy = h.SumWXY()

	ncells := (nxbins + 2) * (nybins + 2)
	hroot.th2.th1.ncells = ncells

	hroot.th2.th1.xaxis.nbins = nxbins
	hroot.th2.th1.xaxis.xmin = h.XMin()
	hroot.th2.th1.xaxis.xmax = h.XMax()

	hroot.th2.th1.yaxis.nbins = nybins
	hroot.th2.th1.yaxis.xmin = h.YMin()
	hroot.th2.th1.yaxis.xmax = h.YMax()

	hroot.arr.Data = make([]float64, ncells)
	hroot.th2.th1.sumw2.Data = make([]float64, ncells)

	ibin := func(ix, iy int) int { return iy*nxbins + ix }

	for ix := 0; ix < h.Binning.Nx; ix++ {
		for iy := 0; iy < h.Binning.Ny; iy++ {
			i := ibin(ix, iy)
			bin := bins[i]
			if ix == 0 {
				yedges = append(yedges, bin.YMin())
			}
			if iy == 0 {
				xedges = append(xedges, bin.XMin())
			}
			hroot.setDist2D(ix+1, iy+1, bin.Dist.SumW(), bin.Dist.SumW2())
		}
	}

	oflows := h.Binning.Outflows[:]
	for i, v := range []struct{ ix, iy int }{
		{0, 0},
		{0, 1},
		{0, nybins + 1},
		{nxbins + 1, 0},
		{nxbins + 1, 1},
		{nxbins + 1, nybins + 1},
		{1, 0},
		{1, nybins + 1},
	} {
		hroot.setDist2D(v.ix, v.iy, oflows[i].SumW(), oflows[i].SumW2())
	}

	xedges = append(xedges, bins[ibin(h.Binning.Nx-1, 0)].XMax())
	yedges = append(yedges, bins[ibin(0, h.Binning.Ny-1)].YMax())

	hroot.th2.th1.SetName(h.Name())
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th2.th1.SetTitle(v.(string))
	}
	hroot.th2.th1.xaxis.xbins.Data = xedges
	hroot.th2.th1.yaxis.xbins.Data = yedges

	return hroot
}

func (*H2D) RVersion() int16 {
	return rvers.H2D
}

func (*H2D) isH2() {}

// Class returns the ROOT class name.
func (*H2D) Class() string {
	return "TH2D"
}

func (h *H2D) Array() rcont.ArrayD {
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

func (h *H2D) setDist2D(ix, iy int, sumw, sumw2 float64) {
	i := h.bin(ix, iy)
	h.arr.Data[i] = float64(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *H2D) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

// AsH2D creates a new hbook.H2D from this ROOT histogram.
func (h *H2D) AsH2D() *hbook.H2D {
	var (
		nx = h.NbinsX()
		ny = h.NbinsY()
		hh = hbook.NewH2D(
			nx, h.XAxis().XMin(), h.XAxis().XMax(),
			ny, h.YAxis().XMin(), h.YAxis().XMax(),
		)
		xinrange = 1
		yinrange = 1
	)
	hh.Ann = hbook.Annotation{
		"name":  h.Name(),
		"title": h.Title(),
	}
	hh.Binning.Outflows = [8]hbook.Dist2D{
		h.dist2D(0, 0),
		h.dist2D(0, yinrange),
		h.dist2D(0, ny+1),
		h.dist2D(nx+1, 0),
		h.dist2D(nx+1, yinrange),
		h.dist2D(nx+1, ny+1),
		h.dist2D(xinrange, 0),
		h.dist2D(xinrange, ny+1),
	}

	hh.Binning.Dist = hbook.Dist2D{
		X: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		},
		Y: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		},
	}
	hh.Binning.Dist.X.Stats.SumWX = float64(h.SumWX())
	hh.Binning.Dist.X.Stats.SumWX2 = float64(h.SumWX2())
	hh.Binning.Dist.Y.Stats.SumWX = float64(h.SumWY())
	hh.Binning.Dist.Y.Stats.SumWX2 = float64(h.SumWY2())
	hh.Binning.Dist.Stats.SumWXY = h.SumWXY()

	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			var (
				i    = iy*nx + ix
				xmin = h.XBinLowEdge(ix + 1)
				xmax = h.XBinWidth(ix+1) + xmin
				ymin = h.YBinLowEdge(iy + 1)
				ymax = h.YBinWidth(iy+1) + ymin
				bin  = &hh.Binning.Bins[i]
			)
			bin.XRange.Min = xmin
			bin.XRange.Max = xmax
			bin.YRange.Min = ymin
			bin.YRange.Max = ymax
			bin.Dist = h.dist2D(ix+1, iy+1)
		}
	}

	return hh
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H2D) MarshalYODA() ([]byte, error) {
	return h.AsH2D().MarshalYODA()
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H2D) UnmarshalYODA(raw []byte) error {
	var hh hbook.H2D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *NewH2DFrom(&hh)
	return nil
}

func (h *H2D) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())

	for _, v := range []rbytes.Marshaler{
		&h.th2,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	return w.SetByteCount(pos, h.Class())
}

func (h *H2D) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers < 1 {
		return fmt.Errorf("rhist: TH2D version too old (%d<1)", vers)
	}

	for _, v := range []rbytes.Unmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := newH2D()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TH2D", f)
}

var (
	_ root.Object        = (*H2D)(nil)
	_ root.Named         = (*H2D)(nil)
	_ H2                 = (*H2D)(nil)
	_ rbytes.Marshaler   = (*H2D)(nil)
	_ rbytes.Unmarshaler = (*H2D)(nil)
)

// H2I implements ROOT TH2I
type H2I struct {
	th2
	arr rcont.ArrayI
}

func newH2I() *H2I {
	return &H2I{
		th2: *newH2(),
	}
}

// NewH2IFrom creates a new H2I from hbook 2-dim histogram.
func NewH2IFrom(h *hbook.H2D) *H2I {
	var (
		hroot  = newH2I()
		bins   = h.Binning.Bins
		nxbins = h.Binning.Nx
		nybins = h.Binning.Ny
		xedges = make([]float64, 0, nxbins+1)
		yedges = make([]float64, 0, nybins+1)
	)

	hroot.th2.th1.entries = float64(h.Entries())
	hroot.th2.th1.tsumw = h.SumW()
	hroot.th2.th1.tsumw2 = h.SumW2()
	hroot.th2.th1.tsumwx = h.SumWX()
	hroot.th2.th1.tsumwx2 = h.SumWX2()
	hroot.th2.tsumwy = h.SumWY()
	hroot.th2.tsumwy2 = h.SumWY2()
	hroot.th2.tsumwxy = h.SumWXY()

	ncells := (nxbins + 2) * (nybins + 2)
	hroot.th2.th1.ncells = ncells

	hroot.th2.th1.xaxis.nbins = nxbins
	hroot.th2.th1.xaxis.xmin = h.XMin()
	hroot.th2.th1.xaxis.xmax = h.XMax()

	hroot.th2.th1.yaxis.nbins = nybins
	hroot.th2.th1.yaxis.xmin = h.YMin()
	hroot.th2.th1.yaxis.xmax = h.YMax()

	hroot.arr.Data = make([]int32, ncells)
	hroot.th2.th1.sumw2.Data = make([]float64, ncells)

	ibin := func(ix, iy int) int { return iy*nxbins + ix }

	for ix := 0; ix < h.Binning.Nx; ix++ {
		for iy := 0; iy < h.Binning.Ny; iy++ {
			i := ibin(ix, iy)
			bin := bins[i]
			if ix == 0 {
				yedges = append(yedges, bin.YMin())
			}
			if iy == 0 {
				xedges = append(xedges, bin.XMin())
			}
			hroot.setDist2D(ix+1, iy+1, bin.Dist.SumW(), bin.Dist.SumW2())
		}
	}

	oflows := h.Binning.Outflows[:]
	for i, v := range []struct{ ix, iy int }{
		{0, 0},
		{0, 1},
		{0, nybins + 1},
		{nxbins + 1, 0},
		{nxbins + 1, 1},
		{nxbins + 1, nybins + 1},
		{1, 0},
		{1, nybins + 1},
	} {
		hroot.setDist2D(v.ix, v.iy, oflows[i].SumW(), oflows[i].SumW2())
	}

	xedges = append(xedges, bins[ibin(h.Binning.Nx-1, 0)].XMax())
	yedges = append(yedges, bins[ibin(0, h.Binning.Ny-1)].YMax())

	hroot.th2.th1.SetName(h.Name())
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th2.th1.SetTitle(v.(string))
	}
	hroot.th2.th1.xaxis.xbins.Data = xedges
	hroot.th2.th1.yaxis.xbins.Data = yedges

	return hroot
}

func (*H2I) RVersion() int16 {
	return rvers.H2I
}

func (*H2I) isH2() {}

// Class returns the ROOT class name.
func (*H2I) Class() string {
	return "TH2I"
}

func (h *H2I) Array() rcont.ArrayI {
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

func (h *H2I) setDist2D(ix, iy int, sumw, sumw2 float64) {
	i := h.bin(ix, iy)
	h.arr.Data[i] = int32(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *H2I) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

// AsH2D creates a new hbook.H2D from this ROOT histogram.
func (h *H2I) AsH2D() *hbook.H2D {
	var (
		nx = h.NbinsX()
		ny = h.NbinsY()
		hh = hbook.NewH2D(
			nx, h.XAxis().XMin(), h.XAxis().XMax(),
			ny, h.YAxis().XMin(), h.YAxis().XMax(),
		)
		xinrange = 1
		yinrange = 1
	)
	hh.Ann = hbook.Annotation{
		"name":  h.Name(),
		"title": h.Title(),
	}
	hh.Binning.Outflows = [8]hbook.Dist2D{
		h.dist2D(0, 0),
		h.dist2D(0, yinrange),
		h.dist2D(0, ny+1),
		h.dist2D(nx+1, 0),
		h.dist2D(nx+1, yinrange),
		h.dist2D(nx+1, ny+1),
		h.dist2D(xinrange, 0),
		h.dist2D(xinrange, ny+1),
	}

	hh.Binning.Dist = hbook.Dist2D{
		X: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		},
		Y: hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		},
	}
	hh.Binning.Dist.X.Stats.SumWX = float64(h.SumWX())
	hh.Binning.Dist.X.Stats.SumWX2 = float64(h.SumWX2())
	hh.Binning.Dist.Y.Stats.SumWX = float64(h.SumWY())
	hh.Binning.Dist.Y.Stats.SumWX2 = float64(h.SumWY2())
	hh.Binning.Dist.Stats.SumWXY = h.SumWXY()

	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			var (
				i    = iy*nx + ix
				xmin = h.XBinLowEdge(ix + 1)
				xmax = h.XBinWidth(ix+1) + xmin
				ymin = h.YBinLowEdge(iy + 1)
				ymax = h.YBinWidth(iy+1) + ymin
				bin  = &hh.Binning.Bins[i]
			)
			bin.XRange.Min = xmin
			bin.XRange.Max = xmax
			bin.YRange.Min = ymin
			bin.YRange.Max = ymax
			bin.Dist = h.dist2D(ix+1, iy+1)
		}
	}

	return hh
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H2I) MarshalYODA() ([]byte, error) {
	return h.AsH2D().MarshalYODA()
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H2I) UnmarshalYODA(raw []byte) error {
	var hh hbook.H2D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *NewH2IFrom(&hh)
	return nil
}

func (h *H2I) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())

	for _, v := range []rbytes.Marshaler{
		&h.th2,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	return w.SetByteCount(pos, h.Class())
}

func (h *H2I) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers < 1 {
		return fmt.Errorf("rhist: TH2I version too old (%d<1)", vers)
	}

	for _, v := range []rbytes.Unmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := newH2I()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TH2I", f)
}

var (
	_ root.Object        = (*H2I)(nil)
	_ root.Named         = (*H2I)(nil)
	_ H2                 = (*H2I)(nil)
	_ rbytes.Marshaler   = (*H2I)(nil)
	_ rbytes.Unmarshaler = (*H2I)(nil)
)
