// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"go-hep.org/x/hep/groot/internal/genroot"
)

func main() {
	genH1()
	genH2()
}

func genH1() {
	f, err := os.Create("./rhist/h1_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rhist", f,
		"bytes", "fmt", "math", "reflect",
		"",
		"go-hep.org/x/hep/hbook",
		"go-hep.org/x/hep/groot/root",
		"go-hep.org/x/hep/groot/rcont",
		"go-hep.org/x/hep/groot/rbytes",
		"go-hep.org/x/hep/groot/rtypes",
		"go-hep.org/x/hep/groot/rvers",
	)

	for i, typ := range []struct {
		Name string
		Type string
		Elem string
	}{
		{
			Name: "H1F",
			Type: "rcont.ArrayF",
			Elem: "float32",
		},
		{
			Name: "H1D",
			Type: "rcont.ArrayD",
			Elem: "float64",
		},
		{
			Name: "H1I",
			Type: "rcont.ArrayI",
			Elem: "int32",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(h1Tmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)
}

func genH2() {
	f, err := os.Create("./rhist/h2_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rhist", f,
		"bytes", "fmt", "math", "reflect",
		"",
		"go-hep.org/x/hep/hbook",
		"go-hep.org/x/hep/groot/root",
		"go-hep.org/x/hep/groot/rcont",
		"go-hep.org/x/hep/groot/rbytes",
		"go-hep.org/x/hep/groot/rtypes",
		"go-hep.org/x/hep/groot/rvers",
	)

	for i, typ := range []struct {
		Name string
		Type string
		Elem string
	}{
		{
			Name: "H2F",
			Type: "rcont.ArrayF",
			Elem: "float32",
		},
		{
			Name: "H2D",
			Type: "rcont.ArrayD",
			Elem: "float64",
		},
		{
			Name: "H2I",
			Type: "rcont.ArrayI",
			Elem: "int32",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(h2Tmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)
}

const h1Tmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	th1
	arr {{.Type}}
}

func new{{.Name}}() *{{.Name}} {
	return &{{.Name}}{
		th1:   *newH1(),
	}
}

// New{{.Name}}From creates a new 1-dim histogram from hbook.
func New{{.Name}}From(h *hbook.H1D) *{{.Name}} {
	var (
		hroot = new{{.Name}}()
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
	hroot.th1.ncells = nbins+2

	hroot.th1.xaxis.nbins = nbins
	hroot.th1.xaxis.xmin = h.XMin()
	hroot.th1.xaxis.xmax = h.XMax()

	hroot.arr.Data = make([]{{.Elem}}, nbins+2)
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

func (*{{.Name}}) RVersion() int16 {
	return rvers.{{.Name}}
}

func (*{{.Name}}) isH1() {}

// Class returns the ROOT class name.
func (*{{.Name}}) Class() string {
	return "T{{.Name}}"
}

func (h *{{.Name}}) MarshalROOT(w *rbytes.WBuffer) (int, error) {
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

func (h *{{.Name}}) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers > rvers.{{.Name}} {
		panic(fmt.Errorf("rhist: invalid {{.Name}} version=%d > %d", vers, rvers.{{.Name}}))
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

func (h *{{.Name}}) Array() {{.Type}} {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *{{.Name}}) Rank() int {
	return 1
}

// NbinsX returns the number of bins in X.
func (h *{{.Name}}) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h*{{.Name}}) XAxis() Axis {
	return &h.th1.xaxis
}

// bin returns the regularized bin number given an x bin pair.
func (h *{{.Name}}) bin(ix int) int {
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
func (h *{{.Name}}) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *{{.Name}}) XBinContent(i int) float64 {
	ibin := h.bin(i)
	return float64(h.arr.Data[ibin])
}

// XBinError returns the bin error in X.
func (h *{{.Name}}) XBinError(i int) float64 {
	ibin := h.bin(i)
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[ibin]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[ibin])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *{{.Name}}) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *{{.Name}}) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

func (h *{{.Name}}) dist1D(i int) hbook.Dist1D {
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

func (h *{{.Name}}) setDist1D(i int, sumw, sumw2 float64) {
	h.arr.Data[i] = {{.Elem}}(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *{{.Name}}) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v+0.5)
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *{{.Name}}) MarshalYODA() ([]byte, error) {
	var (
		nx    = h.NbinsX()
		dflow = [2]hbook.Dist1D{
			h.dist1D(0),    // underflow
			h.dist1D(nx+1), // overflow
		}
		dtot = hbook.Dist1D{
			Dist: hbook.Dist0D{
				N:      int64(h.Entries()),
				SumW:   float64(h.SumW()),
				SumW2:  float64(h.SumW2()),
			},
		}
		dists = make([]hbook.Dist1D, int(nx))
	)
	dtot.Stats.SumWX = float64(h.SumWX())
	dtot.Stats.SumWX2 = float64(h.SumWX2())

	for i := 0; i < nx; i++ {
		dists[i] = h.dist1D(i+1)
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
		xmin := h.XBinLowEdge(i+1)
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
func (h *{{.Name}}) UnmarshalYODA(raw []byte) error {
	var hh hbook.H1D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *New{{.Name}}From(&hh)
	return nil
}

func init() {
	f := func() reflect.Value {
		o := new{{.Name}}()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("T{{.Name}}", f)
}

var (
	_ root.Object        = (*{{.Name}})(nil)
	_ root.Named         = (*{{.Name}})(nil)
	_ H1                 = (*{{.Name}})(nil)
	_ rbytes.Marshaler   = (*{{.Name}})(nil)
	_ rbytes.Unmarshaler = (*{{.Name}})(nil)
)
`

const h2Tmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	th2
	arr {{.Type}}
}

func new{{.Name}}() *{{.Name}} {
	return &{{.Name}}{
		th2:   *newH2(),
	}
}

// New{{.Name}}From creates a new {{.Name}} from hbook 2-dim histogram.
func New{{.Name}}From(h *hbook.H2D) *{{.Name}} {
	var (
		hroot  = new{{.Name}}()
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

	hroot.arr.Data = make([]{{.Elem}}, ncells)
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
	for i, v := range []struct{ix,iy int}{
		{0, 0},
		{0, 1},
		{0, nybins+1},
		{nxbins + 1, 0},
		{nxbins + 1, 1},
		{nxbins + 1, nybins + 1},
		{1, 0},
		{1, nybins + 1},
	}{
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

func (*{{.Name}}) RVersion() int16 {
	return rvers.{{.Name}}
}

func (*{{.Name}}) isH2() {}

// Class returns the ROOT class name.
func (*{{.Name}}) Class() string {
	return "T{{.Name}}"
}

func (h *{{.Name}}) Array() {{.Type}} {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *{{.Name}}) Rank() int {
	return 2
}

// NbinsX returns the number of bins in X.
func (h *{{.Name}}) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h*{{.Name}}) XAxis() Axis {
	return &h.th1.xaxis
}

// XBinCenter returns the bin center value in X.
func (h *{{.Name}}) XBinCenter(i int) float64 {
	return float64(h.th1.xaxis.BinCenter(i))
}

// XBinContent returns the bin content value in X.
func (h *{{.Name}}) XBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// XBinError returns the bin error in X.
func (h *{{.Name}}) XBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// XBinLowEdge returns the bin lower edge value in X.
func (h *{{.Name}}) XBinLowEdge(i int) float64 {
	return h.th1.xaxis.BinLowEdge(i)
}

// XBinWidth returns the bin width in X.
func (h *{{.Name}}) XBinWidth(i int) float64 {
	return h.th1.xaxis.BinWidth(i)
}

// NbinsY returns the number of bins in Y.
func (h *{{.Name}}) NbinsY() int {
	return h.th1.yaxis.nbins
}

// YAxis returns the axis along Y.
func (h*{{.Name}}) YAxis() Axis {
	return &h.th1.yaxis
}

// YBinCenter returns the bin center value in Y.
func (h *{{.Name}}) YBinCenter(i int) float64 {
	return float64(h.th1.yaxis.BinCenter(i))
}

// YBinContent returns the bin content value in Y.
func (h *{{.Name}}) YBinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// YBinError returns the bin error in Y.
func (h *{{.Name}}) YBinError(i int) float64 {
	if len(h.th1.sumw2.Data) > 0 {
		return math.Sqrt(float64(h.th1.sumw2.Data[i]))
	}
	return math.Sqrt(math.Abs(float64(h.arr.Data[i])))
}

// YBinLowEdge returns the bin lower edge value in Y.
func (h *{{.Name}}) YBinLowEdge(i int) float64 {
	return h.th1.yaxis.BinLowEdge(i)
}

// YBinWidth returns the bin width in Y.
func (h *{{.Name}}) YBinWidth(i int) float64 {
	return h.th1.yaxis.BinWidth(i)
}

// bin returns the regularized bin number given an (x,y) bin index pair.
func (h *{{.Name}}) bin(ix, iy int) int {
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

func (h *{{.Name}}) dist2D(ix, iy int) hbook.Dist2D {
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

func (h *{{.Name}}) setDist2D(ix, iy int, sumw, sumw2 float64) {
	i := h.bin(ix, iy)
	h.arr.Data[i] = {{.Elem}}(sumw)
	h.th1.sumw2.Data[i] = sumw2
}

func (h *{{.Name}}) entries(height, err float64) int64 {
	if height <= 0 {
		return 0
	}
	v := height / err
	return int64(v*v + 0.5)
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *{{.Name}}) MarshalYODA() ([]byte, error) {
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
				Dist: hbook.Dist0D {
					N:      int64(h.Entries()),
					SumW:   float64(h.SumW()),
					SumW2:  float64(h.SumW2()),
				},
			},
			Y: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:      int64(h.Entries()),
					SumW:   float64(h.SumW()),
					SumW2:  float64(h.SumW2()),
				},
			},
		}
		dists = make([]hbook.Dist2D, int(nx*ny))
	)
	dtot.X.Stats.SumWX = float64(h.SumWX())
	dtot.X.Stats.SumWX2 = float64(h.SumWX2())
	dtot.Y.Stats.SumWX = float64(h.SumWY())
	dtot.Y.Stats.SumWX2 = float64(h.SumWY2())
	dtot.Stats.SumWXY = h.SumWXY()

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
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY(), d.Entries(),
	)

	if false { // FIXME(sbinet)
		for _, d := range dflow {
			fmt.Fprintf(
				buf,
				"%s\t%s\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				name, name,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY(), d.Entries(),
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
			i := iy*nx+ix
			d := &dists[i]
			fmt.Fprintf(
				buf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				xmin, xmax, ymin, ymax,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY(), d.Entries(),
			)
		}
	}
	fmt.Fprintf(buf, "END YODA_HISTO2D\n\n")
	return buf.Bytes(), nil
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *{{.Name}}) UnmarshalYODA(raw []byte) error {
	var hh hbook.H2D
	err := hh.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h = *New{{.Name}}From(&hh)
	return nil
}

func (h *{{.Name}}) MarshalROOT(w *rbytes.WBuffer) (int, error) {
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

func (h *{{.Name}}) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers < 1 {
		return fmt.Errorf("rhist: T{{.Name}} version too old (%d<1)", vers)
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
		o := new{{.Name}}()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("T{{.Name}}", f)
}

var (
	_ root.Object        = (*{{.Name}})(nil)
	_ root.Named         = (*{{.Name}})(nil)
	_ H2                 = (*{{.Name}})(nil)
	_ rbytes.Marshaler   = (*{{.Name}})(nil)
	_ rbytes.Unmarshaler = (*{{.Name}})(nil)
)
`
