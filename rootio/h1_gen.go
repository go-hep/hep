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

// H1F implements ROOT TH1F
type H1F struct {
	rvers int16
	th1
	arr ArrayF
}

func (*H1F) isH1() {}

// Class returns the ROOT class name.
func (*H1F) Class() string {
	return "TH1F"
}

func (h *H1F) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(h.rvers)

	for _, v := range []ROOTMarshaler{
		&h.th1,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	return w.SetByteCount(pos, "TH1F")
}

func (h *H1F) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: TH1F version too old (%d<1)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH1F")
	return r.err
}

func (h *H1F) Array() ArrayF {
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
			SumWX:  float64(h.SumWX()),
			SumWX2: float64(h.SumWX2()),
		}
		dists = make([]hbook.Dist1D, int(nx))
	)

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
		dtot.SumW(), dtot.SumW2(), dtot.SumWX, dtot.SumWX2, dtot.Entries(),
	)

	name = "Underflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[0].SumW(), dflow[0].SumW2(), dflow[0].SumWX, dflow[0].SumWX2, dflow[0].Entries(),
	)

	name = "Overflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[1].SumW(), dflow[1].SumW2(), dflow[1].SumWX, dflow[1].SumWX2, dflow[1].Entries(),
	)
	fmt.Fprintf(buf, "# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries\n")
	for i, d := range dists {
		xmin := h.XBinLowEdge(i + 1)
		xmax := h.XBinWidth(i+1) + xmin
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			xmin, xmax,
			d.SumW(), d.SumW2(), d.SumWX, d.SumWX2, d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")

	return buf.Bytes(), nil
}

func init() {
	f := func() reflect.Value {
		o := &H1F{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH1F", f)
	Factory.add("*rootio.H1F", f)
}

var (
	_ Object          = (*H1F)(nil)
	_ Named           = (*H1F)(nil)
	_ H1              = (*H1F)(nil)
	_ ROOTMarshaler   = (*H1F)(nil)
	_ ROOTUnmarshaler = (*H1F)(nil)
)

// H1D implements ROOT TH1D
type H1D struct {
	rvers int16
	th1
	arr ArrayD
}

func (*H1D) isH1() {}

// Class returns the ROOT class name.
func (*H1D) Class() string {
	return "TH1D"
}

func (h *H1D) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(h.rvers)

	for _, v := range []ROOTMarshaler{
		&h.th1,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	return w.SetByteCount(pos, "TH1D")
}

func (h *H1D) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: TH1D version too old (%d<1)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH1D")
	return r.err
}

func (h *H1D) Array() ArrayD {
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
			SumWX:  float64(h.SumWX()),
			SumWX2: float64(h.SumWX2()),
		}
		dists = make([]hbook.Dist1D, int(nx))
	)

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
		dtot.SumW(), dtot.SumW2(), dtot.SumWX, dtot.SumWX2, dtot.Entries(),
	)

	name = "Underflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[0].SumW(), dflow[0].SumW2(), dflow[0].SumWX, dflow[0].SumWX2, dflow[0].Entries(),
	)

	name = "Overflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[1].SumW(), dflow[1].SumW2(), dflow[1].SumWX, dflow[1].SumWX2, dflow[1].Entries(),
	)
	fmt.Fprintf(buf, "# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries\n")
	for i, d := range dists {
		xmin := h.XBinLowEdge(i + 1)
		xmax := h.XBinWidth(i+1) + xmin
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			xmin, xmax,
			d.SumW(), d.SumW2(), d.SumWX, d.SumWX2, d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")

	return buf.Bytes(), nil
}

func init() {
	f := func() reflect.Value {
		o := &H1D{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH1D", f)
	Factory.add("*rootio.H1D", f)
}

var (
	_ Object          = (*H1D)(nil)
	_ Named           = (*H1D)(nil)
	_ H1              = (*H1D)(nil)
	_ ROOTMarshaler   = (*H1D)(nil)
	_ ROOTUnmarshaler = (*H1D)(nil)
)

// H1I implements ROOT TH1I
type H1I struct {
	rvers int16
	th1
	arr ArrayI
}

func (*H1I) isH1() {}

// Class returns the ROOT class name.
func (*H1I) Class() string {
	return "TH1I"
}

func (h *H1I) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(h.rvers)

	for _, v := range []ROOTMarshaler{
		&h.th1,
		&h.arr,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	return w.SetByteCount(pos, "TH1I")
}

func (h *H1I) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: TH1I version too old (%d<1)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH1I")
	return r.err
}

func (h *H1I) Array() ArrayI {
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
			SumWX:  float64(h.SumWX()),
			SumWX2: float64(h.SumWX2()),
		}
		dists = make([]hbook.Dist1D, int(nx))
	)

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
		dtot.SumW(), dtot.SumW2(), dtot.SumWX, dtot.SumWX2, dtot.Entries(),
	)

	name = "Underflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[0].SumW(), dflow[0].SumW2(), dflow[0].SumWX, dflow[0].SumWX2, dflow[0].Entries(),
	)

	name = "Overflow"
	fmt.Fprintf(
		buf,
		"%s\t%s\t%e\t%e\t%e\t%e\t%d\n",
		name, name,
		dflow[1].SumW(), dflow[1].SumW2(), dflow[1].SumWX, dflow[1].SumWX2, dflow[1].Entries(),
	)
	fmt.Fprintf(buf, "# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries\n")
	for i, d := range dists {
		xmin := h.XBinLowEdge(i + 1)
		xmax := h.XBinWidth(i+1) + xmin
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			xmin, xmax,
			d.SumW(), d.SumW2(), d.SumWX, d.SumWX2, d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")

	return buf.Bytes(), nil
}

func init() {
	f := func() reflect.Value {
		o := &H1I{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH1I", f)
	Factory.add("*rootio.H1I", f)
}

var (
	_ Object          = (*H1I)(nil)
	_ Named           = (*H1I)(nil)
	_ H1              = (*H1I)(nil)
	_ ROOTMarshaler   = (*H1I)(nil)
	_ ROOTUnmarshaler = (*H1I)(nil)
)
