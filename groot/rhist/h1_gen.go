// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rhist

import (
	"bytes"
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

// H1F implements ROOT TH1F
type H1F struct {
	th1
	arr rcont.ArrayF
}

func newH1F() *H1F {
	return &H1F{
		th1: *newH1(),
	}
}

// NewH1FFrom creates a new 1-dim histogram from hbook.
func NewH1FFrom(h *hbook.H1D) *H1F {
	var (
		hroot = newH1F()
		bins  = h.Binning.Bins
		nbins = len(bins)
		edges = make([]float64, 0, nbins+1)
		uflow = h.Binning.Underflow()
		oflow = h.Binning.Overflow()
	)

	hroot.th1.entries = float64(h.Entries())
	hroot.th1.tsumw = h.SumW()
	hroot.th1.tsumw2 = h.SumW2()
	hroot.th1.tsumwx = h.SumWX()
	hroot.th1.tsumwx2 = h.SumWX2()
	hroot.th1.ncells = nbins + 2

	hroot.th1.xaxis.nbins = nbins
	hroot.th1.xaxis.xmin = h.XMin()
	hroot.th1.xaxis.xmax = h.XMax()

	hroot.arr.Data = make([]float32, nbins+2)
	hroot.th1.sumw2.Data = make([]float64, nbins+2)

	for i, bin := range bins {
		if i == 0 {
			edges = append(edges, bin.XMin())
		}
		edges = append(edges, bin.XMax())
		hroot.setDist1D(i+1, bin.Dist.SumW(), bin.Dist.SumW2())
	}
	hroot.setDist1D(0, uflow.SumW(), uflow.SumW2())
	hroot.setDist1D(nbins+1, oflow.SumW(), oflow.SumW2())

	hroot.th1.SetName(h.Name())
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th1.SetTitle(v.(string))
	}
	hroot.th1.xaxis.xbins.Data = edges
	return hroot
}

func (*H1F) RVersion() int16 {
	return rvers.H1F
}

func (*H1F) isH1() {}

// Class returns the ROOT class name.
func (*H1F) Class() string {
	return "TH1F"
}

func (h *H1F) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())

	for _, v := range []rbytes.Marshaler{
		&h.th1,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	return w.SetByteCount(pos, h.Class())
}

func (h *H1F) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers > rvers.H1F {
		panic(fmt.Errorf("rhist: invalid H1F version=%d > %d", vers, rvers.H1F))
	}

	for _, v := range []rbytes.Unmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

func (h *H1F) Array() rcont.ArrayF {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H1F) Rank() int {
	return 1
}

// NbinsX returns the number of bins in X.
func (h *H1F) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H1F) XAxis() Axis {
	return &h.th1.xaxis
}

// bin returns the regularized bin number given an x bin pair.
func (h *H1F) bin(ix int) int {
	nx := h.th1.xaxis.nbins + 1 // overflow bin
	switch {
	case ix < 0:
		ix = 0
	case ix > nx:
		ix = nx
	}
	return ix
}

// XBinCenter returns the bin center value in X.
func (h *H1F) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *H1F) XBinContent(i int) float64 {
	ibin := h.bin(i)
	return float64(h.arr.Data[ibin])
}

// XBinError returns the bin error in X.
func (h *H1F) XBinError(i int) float64 {
	ibin := h.bin(i)
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[ibin]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[ibin])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *H1F) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *H1F) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

func (h *H1F) dist1D(i int) hbook.Dist1D {
	v := h.XBinContent(i)
	err := h.XBinError(i)
	n := h.entries(v, err)
	sumw := h.arr.Data[i]
	sumw2 := 0.0
	if len(h.th1.sumw2.Data) > 0 {
		sumw2 = h.th1.sumw2.Data[i]
	}
	return hbook.Dist1D{
		Dist: hbook.Dist0D{
			N:     n,
			SumW:  float64(sumw),
			SumW2: float64(sumw2),
		},
	}
}

func (h *H1F) setDist1D(i int, sumw, sumw2 float64) {
	h.arr.Data[i] = float32(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *H1F) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H1F) MarshalYODA() ([]byte, error) {
	var (
		nx    = h.NbinsX()
		dflow = [2]hbook.Dist1D{
			h.dist1D(0),      // underflow
			h.dist1D(nx + 1), // overflow
		}
		dtot = hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		}
		dists = make([]hbook.Dist1D, int(nx))
	)
	dtot.Stats.SumWX = float64(h.SumWX())
	dtot.Stats.SumWX2 = float64(h.SumWX2())

	for i := 0; i < nx; i++ {
		dists[i] = h.dist1D(i + 1)
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "BEGIN YODA_HISTO1D /%s\n", h.Name())
	fmt.Fprintf(buf, "Path=/%s\n", h.Name())
	fmt.Fprintf(buf, "Title=%s\n", h.Title())
	fmt.Fprintf(buf, "Type=Histo1D\n")
	fmt.Fprintf(buf, "# Mean: %e\n", math.NaN())
	fmt.Fprintf(buf, "# Area: %e\n", math.NaN())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t numEntries\n")

	var name = "Total   "
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dtot.SumW(), dtot.SumW2(), dtot.SumWX(), dtot.SumWX2(), dtot.Entries(),
	)

	name = "Underflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[0].SumW(), dflow[0].SumW2(), dflow[0].SumWX(), dflow[0].SumWX2(), dflow[0].Entries(),
	)

	name = "Overflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[1].SumW(), dflow[1].SumW2(), dflow[1].SumWX(), dflow[1].SumWX2(), dflow[1].Entries(),
	)
	fmt.Fprintf(buf, "# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries\n")
	for i, d := range dists {
		xmin := h.XBinLowEdge(i + 1)
		xmax := h.XBinWidth(i+1) + xmin
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			xmin, xmax,
			d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")

	return buf.Bytes(), nil
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H1F) UnmarshalYODA(raw []byte) error {
	var hh hbook.H1D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *NewH1FFrom(&hh)
	return nil
}

func init() {
	f := func() reflect.Value {
		o := newH1F()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TH1F", f)
}

var (
	_ root.Object        = (*H1F)(nil)
	_ root.Named         = (*H1F)(nil)
	_ H1                 = (*H1F)(nil)
	_ rbytes.Marshaler   = (*H1F)(nil)
	_ rbytes.Unmarshaler = (*H1F)(nil)
)

// H1D implements ROOT TH1D
type H1D struct {
	th1
	arr rcont.ArrayD
}

func newH1D() *H1D {
	return &H1D{
		th1: *newH1(),
	}
}

// NewH1DFrom creates a new 1-dim histogram from hbook.
func NewH1DFrom(h *hbook.H1D) *H1D {
	var (
		hroot = newH1D()
		bins  = h.Binning.Bins
		nbins = len(bins)
		edges = make([]float64, 0, nbins+1)
		uflow = h.Binning.Underflow()
		oflow = h.Binning.Overflow()
	)

	hroot.th1.entries = float64(h.Entries())
	hroot.th1.tsumw = h.SumW()
	hroot.th1.tsumw2 = h.SumW2()
	hroot.th1.tsumwx = h.SumWX()
	hroot.th1.tsumwx2 = h.SumWX2()
	hroot.th1.ncells = nbins + 2

	hroot.th1.xaxis.nbins = nbins
	hroot.th1.xaxis.xmin = h.XMin()
	hroot.th1.xaxis.xmax = h.XMax()

	hroot.arr.Data = make([]float64, nbins+2)
	hroot.th1.sumw2.Data = make([]float64, nbins+2)

	for i, bin := range bins {
		if i == 0 {
			edges = append(edges, bin.XMin())
		}
		edges = append(edges, bin.XMax())
		hroot.setDist1D(i+1, bin.Dist.SumW(), bin.Dist.SumW2())
	}
	hroot.setDist1D(0, uflow.SumW(), uflow.SumW2())
	hroot.setDist1D(nbins+1, oflow.SumW(), oflow.SumW2())

	hroot.th1.SetName(h.Name())
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th1.SetTitle(v.(string))
	}
	hroot.th1.xaxis.xbins.Data = edges
	return hroot
}

func (*H1D) RVersion() int16 {
	return rvers.H1D
}

func (*H1D) isH1() {}

// Class returns the ROOT class name.
func (*H1D) Class() string {
	return "TH1D"
}

func (h *H1D) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())

	for _, v := range []rbytes.Marshaler{
		&h.th1,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	return w.SetByteCount(pos, h.Class())
}

func (h *H1D) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers > rvers.H1D {
		panic(fmt.Errorf("rhist: invalid H1D version=%d > %d", vers, rvers.H1D))
	}

	for _, v := range []rbytes.Unmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

func (h *H1D) Array() rcont.ArrayD {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H1D) Rank() int {
	return 1
}

// NbinsX returns the number of bins in X.
func (h *H1D) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H1D) XAxis() Axis {
	return &h.th1.xaxis
}

// bin returns the regularized bin number given an x bin pair.
func (h *H1D) bin(ix int) int {
	nx := h.th1.xaxis.nbins + 1 // overflow bin
	switch {
	case ix < 0:
		ix = 0
	case ix > nx:
		ix = nx
	}
	return ix
}

// XBinCenter returns the bin center value in X.
func (h *H1D) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *H1D) XBinContent(i int) float64 {
	ibin := h.bin(i)
	return float64(h.arr.Data[ibin])
}

// XBinError returns the bin error in X.
func (h *H1D) XBinError(i int) float64 {
	ibin := h.bin(i)
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[ibin]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[ibin])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *H1D) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *H1D) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

func (h *H1D) dist1D(i int) hbook.Dist1D {
	v := h.XBinContent(i)
	err := h.XBinError(i)
	n := h.entries(v, err)
	sumw := h.arr.Data[i]
	sumw2 := 0.0
	if len(h.th1.sumw2.Data) > 0 {
		sumw2 = h.th1.sumw2.Data[i]
	}
	return hbook.Dist1D{
		Dist: hbook.Dist0D{
			N:     n,
			SumW:  float64(sumw),
			SumW2: float64(sumw2),
		},
	}
}

func (h *H1D) setDist1D(i int, sumw, sumw2 float64) {
	h.arr.Data[i] = float64(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *H1D) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H1D) MarshalYODA() ([]byte, error) {
	var (
		nx    = h.NbinsX()
		dflow = [2]hbook.Dist1D{
			h.dist1D(0),      // underflow
			h.dist1D(nx + 1), // overflow
		}
		dtot = hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		}
		dists = make([]hbook.Dist1D, int(nx))
	)
	dtot.Stats.SumWX = float64(h.SumWX())
	dtot.Stats.SumWX2 = float64(h.SumWX2())

	for i := 0; i < nx; i++ {
		dists[i] = h.dist1D(i + 1)
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "BEGIN YODA_HISTO1D /%s\n", h.Name())
	fmt.Fprintf(buf, "Path=/%s\n", h.Name())
	fmt.Fprintf(buf, "Title=%s\n", h.Title())
	fmt.Fprintf(buf, "Type=Histo1D\n")
	fmt.Fprintf(buf, "# Mean: %e\n", math.NaN())
	fmt.Fprintf(buf, "# Area: %e\n", math.NaN())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t numEntries\n")

	var name = "Total   "
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dtot.SumW(), dtot.SumW2(), dtot.SumWX(), dtot.SumWX2(), dtot.Entries(),
	)

	name = "Underflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[0].SumW(), dflow[0].SumW2(), dflow[0].SumWX(), dflow[0].SumWX2(), dflow[0].Entries(),
	)

	name = "Overflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[1].SumW(), dflow[1].SumW2(), dflow[1].SumWX(), dflow[1].SumWX2(), dflow[1].Entries(),
	)
	fmt.Fprintf(buf, "# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries\n")
	for i, d := range dists {
		xmin := h.XBinLowEdge(i + 1)
		xmax := h.XBinWidth(i+1) + xmin
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			xmin, xmax,
			d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")

	return buf.Bytes(), nil
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H1D) UnmarshalYODA(raw []byte) error {
	var hh hbook.H1D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *NewH1DFrom(&hh)
	return nil
}

func init() {
	f := func() reflect.Value {
		o := newH1D()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TH1D", f)
}

var (
	_ root.Object        = (*H1D)(nil)
	_ root.Named         = (*H1D)(nil)
	_ H1                 = (*H1D)(nil)
	_ rbytes.Marshaler   = (*H1D)(nil)
	_ rbytes.Unmarshaler = (*H1D)(nil)
)

// H1I implements ROOT TH1I
type H1I struct {
	th1
	arr rcont.ArrayI
}

func newH1I() *H1I {
	return &H1I{
		th1: *newH1(),
	}
}

// NewH1IFrom creates a new 1-dim histogram from hbook.
func NewH1IFrom(h *hbook.H1D) *H1I {
	var (
		hroot = newH1I()
		bins  = h.Binning.Bins
		nbins = len(bins)
		edges = make([]float64, 0, nbins+1)
		uflow = h.Binning.Underflow()
		oflow = h.Binning.Overflow()
	)

	hroot.th1.entries = float64(h.Entries())
	hroot.th1.tsumw = h.SumW()
	hroot.th1.tsumw2 = h.SumW2()
	hroot.th1.tsumwx = h.SumWX()
	hroot.th1.tsumwx2 = h.SumWX2()
	hroot.th1.ncells = nbins + 2

	hroot.th1.xaxis.nbins = nbins
	hroot.th1.xaxis.xmin = h.XMin()
	hroot.th1.xaxis.xmax = h.XMax()

	hroot.arr.Data = make([]int32, nbins+2)
	hroot.th1.sumw2.Data = make([]float64, nbins+2)

	for i, bin := range bins {
		if i == 0 {
			edges = append(edges, bin.XMin())
		}
		edges = append(edges, bin.XMax())
		hroot.setDist1D(i+1, bin.Dist.SumW(), bin.Dist.SumW2())
	}
	hroot.setDist1D(0, uflow.SumW(), uflow.SumW2())
	hroot.setDist1D(nbins+1, oflow.SumW(), oflow.SumW2())

	hroot.th1.SetName(h.Name())
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th1.SetTitle(v.(string))
	}
	hroot.th1.xaxis.xbins.Data = edges
	return hroot
}

func (*H1I) RVersion() int16 {
	return rvers.H1I
}

func (*H1I) isH1() {}

// Class returns the ROOT class name.
func (*H1I) Class() string {
	return "TH1I"
}

func (h *H1I) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())

	for _, v := range []rbytes.Marshaler{
		&h.th1,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	return w.SetByteCount(pos, h.Class())
}

func (h *H1I) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers > rvers.H1I {
		panic(fmt.Errorf("rhist: invalid H1I version=%d > %d", vers, rvers.H1I))
	}

	for _, v := range []rbytes.Unmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

func (h *H1I) Array() rcont.ArrayI {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H1I) Rank() int {
	return 1
}

// NbinsX returns the number of bins in X.
func (h *H1I) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H1I) XAxis() Axis {
	return &h.th1.xaxis
}

// bin returns the regularized bin number given an x bin pair.
func (h *H1I) bin(ix int) int {
	nx := h.th1.xaxis.nbins + 1 // overflow bin
	switch {
	case ix < 0:
		ix = 0
	case ix > nx:
		ix = nx
	}
	return ix
}

// XBinCenter returns the bin center value in X.
func (h *H1I) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *H1I) XBinContent(i int) float64 {
	ibin := h.bin(i)
	return float64(h.arr.Data[ibin])
}

// XBinError returns the bin error in X.
func (h *H1I) XBinError(i int) float64 {
	ibin := h.bin(i)
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[ibin]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[ibin])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *H1I) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *H1I) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

func (h *H1I) dist1D(i int) hbook.Dist1D {
	v := h.XBinContent(i)
	err := h.XBinError(i)
	n := h.entries(v, err)
	sumw := h.arr.Data[i]
	sumw2 := 0.0
	if len(h.th1.sumw2.Data) > 0 {
		sumw2 = h.th1.sumw2.Data[i]
	}
	return hbook.Dist1D{
		Dist: hbook.Dist0D{
			N:     n,
			SumW:  float64(sumw),
			SumW2: float64(sumw2),
		},
	}
}

func (h *H1I) setDist1D(i int, sumw, sumw2 float64) {
	h.arr.Data[i] = int32(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *H1I) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H1I) MarshalYODA() ([]byte, error) {
	var (
		nx    = h.NbinsX()
		dflow = [2]hbook.Dist1D{
			h.dist1D(0),      // underflow
			h.dist1D(nx + 1), // overflow
		}
		dtot = hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:     int64(h.Entries()),
				SumW:  float64(h.SumW()),
				SumW2: float64(h.SumW2()),
			},
		}
		dists = make([]hbook.Dist1D, int(nx))
	)
	dtot.Stats.SumWX = float64(h.SumWX())
	dtot.Stats.SumWX2 = float64(h.SumWX2())

	for i := 0; i < nx; i++ {
		dists[i] = h.dist1D(i + 1)
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "BEGIN YODA_HISTO1D /%s\n", h.Name())
	fmt.Fprintf(buf, "Path=/%s\n", h.Name())
	fmt.Fprintf(buf, "Title=%s\n", h.Title())
	fmt.Fprintf(buf, "Type=Histo1D\n")
	fmt.Fprintf(buf, "# Mean: %e\n", math.NaN())
	fmt.Fprintf(buf, "# Area: %e\n", math.NaN())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t numEntries\n")

	var name = "Total   "
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dtot.SumW(), dtot.SumW2(), dtot.SumWX(), dtot.SumWX2(), dtot.Entries(),
	)

	name = "Underflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[0].SumW(), dflow[0].SumW2(), dflow[0].SumWX(), dflow[0].SumWX2(), dflow[0].Entries(),
	)

	name = "Overflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[1].SumW(), dflow[1].SumW2(), dflow[1].SumWX(), dflow[1].SumWX2(), dflow[1].Entries(),
	)
	fmt.Fprintf(buf, "# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries\n")
	for i, d := range dists {
		xmin := h.XBinLowEdge(i + 1)
		xmax := h.XBinWidth(i+1) + xmin
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			xmin, xmax,
			d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")

	return buf.Bytes(), nil
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H1I) UnmarshalYODA(raw []byte) error {
	var hh hbook.H1D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *NewH1IFrom(&hh)
	return nil
}

func init() {
	f := func() reflect.Value {
		o := newH1I()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TH1I", f)
}

var (
	_ root.Object        = (*H1I)(nil)
	_ root.Named         = (*H1I)(nil)
	_ H1                 = (*H1I)(nil)
	_ rbytes.Marshaler   = (*H1I)(nil)
	_ rbytes.Unmarshaler = (*H1I)(nil)
)
