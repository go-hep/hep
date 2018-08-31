// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// gen-code generates code for simple ROOT classes hierarchies.
package main

import (
	"encoding/base64"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"go-hep.org/x/hep/rootio"
)

func main() {
	genLeaves()
	genArrays()
	genH1()
	genH2()

	genStreamers()
}

func gofmt(f *os.File) {
	fname := f.Name()
	src, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	src, err = format.Source(src)
	if err != nil {
		log.Fatalf("error formating sources of %q: %v\n", fname, err)
	}

	err = ioutil.WriteFile(fname, src, 0644)
	if err != nil {
		log.Fatalf("error writing back %q: %v\n", fname, err)
	}
}

func genImports(w io.Writer, imports ...string) {
	fmt.Fprintf(w, srcHeader)
	for _, imp := range imports {
		if imp == "" {
			fmt.Fprintf(w, "\n")
			continue
		}
		fmt.Fprintf(w, "\t%q\n", imp)
	}
	fmt.Fprintf(w, ")\n\n")
}

func genLeaves() {
	f, err := os.Create("leaf_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genImports(f)

	for i, typ := range []struct {
		Name       string
		Type       string
		Kind       string
		DoUnsigned bool
		RFunc      string
		RFuncArray string
		WFunc      string
		WFuncArray string
		RangeType  string
		RRangeFunc string
		WRangeFunc string
		Count      bool
	}{
		{
			Name:       "LeafO",
			Type:       "bool",
			Kind:       "reflect.Bool",
			RFunc:      "r.ReadBool()",
			RFuncArray: "r.ReadFastArrayBool",
			WFunc:      "w.WriteBool",
			WFuncArray: "w.WriteFastArrayBool",
		},
		{
			Name:       "LeafB",
			Type:       "int8",
			Kind:       "reflect.Int8",
			DoUnsigned: true,
			RFunc:      "r.ReadI8()",
			RFuncArray: "r.ReadFastArrayI8",
			WFunc:      "w.WriteI8",
			WFuncArray: "w.WriteFastArrayI8",
			Count:      true,
		},
		{
			Name:       "LeafS",
			Type:       "int16",
			Kind:       "reflect.Int16",
			DoUnsigned: true,
			RFunc:      "r.ReadI16()",
			RFuncArray: "r.ReadFastArrayI16",
			WFunc:      "w.WriteI16",
			WFuncArray: "w.WriteFastArrayI16",
			Count:      true,
		},
		{
			Name:       "LeafI",
			Type:       "int32",
			Kind:       "reflect.Int32",
			DoUnsigned: true,
			RFunc:      "r.ReadI32()",
			RFuncArray: "r.ReadFastArrayI32",
			WFunc:      "w.WriteI32",
			WFuncArray: "w.WriteFastArrayI32",
			Count:      true,
		},
		{
			Name:       "LeafL",
			Type:       "int64",
			Kind:       "reflect.Int64",
			DoUnsigned: true,
			RFunc:      "r.ReadI64()",
			RFuncArray: "r.ReadFastArrayI64",
			WFunc:      "w.WriteI64",
			WFuncArray: "w.WriteFastArrayI64",
			Count:      true,
		},
		{
			Name:       "LeafF",
			Type:       "float32",
			Kind:       "reflect.Float32",
			RFunc:      "r.ReadF32()",
			RFuncArray: "r.ReadFastArrayF32",
			WFunc:      "w.WriteF32",
			WFuncArray: "w.WriteFastArrayF32",
		},
		{
			Name:       "LeafD",
			Type:       "float64",
			Kind:       "reflect.Float64",
			RFunc:      "r.ReadF64()",
			RFuncArray: "r.ReadFastArrayF64",
			WFunc:      "w.WriteF64",
			WFuncArray: "w.WriteFastArrayF64",
		},
		{
			Name:       "LeafC",
			Type:       "string",
			Kind:       "reflect.String",
			RFunc:      "r.ReadString()",
			RFuncArray: "r.ReadFastArrayString",
			WFunc:      "w.WriteString()",
			WFuncArray: "w.WriteFastArrayString",
			RangeType:  "int32",
			RRangeFunc: "r.ReadI32()",
			WRangeFunc: "w.WriteI32",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		if typ.RangeType == "" {
			typ.RangeType = typ.Type
			typ.RRangeFunc = typ.RFunc
			typ.WRangeFunc = typ.WFunc
		}
		tmpl := template.Must(template.New(typ.Name).Parse(leafTmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}

func genArrays() {
	f, err := os.Create("array_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genImports(f)

	for i, typ := range []struct {
		Name  string
		Type  string
		RFunc string
		WFunc string
	}{
		{
			Name:  "ArrayI",
			Type:  "int32",
			RFunc: "r.ReadFastArrayI32",
			WFunc: "w.WriteFastArrayI32",
		},
		{
			Name:  "ArrayL64",
			Type:  "int64",
			RFunc: "r.ReadFastArrayI64",
			WFunc: "w.WriteFastArrayI64",
		},
		{
			Name:  "ArrayF",
			Type:  "float32",
			RFunc: "r.ReadFastArrayF32",
			WFunc: "w.WriteFastArrayF32",
		},
		{
			Name:  "ArrayD",
			Type:  "float64",
			RFunc: "r.ReadFastArrayF64",
			WFunc: "w.WriteFastArrayF64",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(arrayTmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}

func genH1() {
	f, err := os.Create("h1_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genImports(f, "bytes", "fmt", "math", "", "go-hep.org/x/hep/hbook")

	for i, typ := range []struct {
		Name string
		Type string
		Elem string
	}{
		{
			Name: "H1F",
			Type: "ArrayF",
			Elem: "float32",
		},
		{
			Name: "H1D",
			Type: "ArrayD",
			Elem: "float64",
		},
		{
			Name: "H1I",
			Type: "ArrayI",
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
	gofmt(f)
}

func genH2() {
	f, err := os.Create("h2_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genImports(f, "bytes", "fmt", "math", "", "go-hep.org/x/hep/hbook")

	for i, typ := range []struct {
		Name string
		Type string
		Elem string
	}{
		{
			Name: "H2F",
			Type: "ArrayF",
			Elem: "float32",
		},
		{
			Name: "H2D",
			Type: "ArrayD",
			Elem: "float64",
		},
		{
			Name: "H2I",
			Type: "ArrayI",
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
	gofmt(f)
}

func genStreamers() {
	classes := []string{
		"TObject",
		"TDirectory",
		"TDirectoryFile",
		"TKey",
		"TNamed",
		"TList",
		"THashList",
		"TObjArray",
		"TObjString",
		"TGraph", "TGraphErrors", "TGraphAsymmErrors",
		"TH1F", "TH1D", "TH1I",
		"TH2F", "TH2D", "TH2I",
		"TTree",
	}

	const (
		macro = "genstreamers.C"
		oname = "streamers.root"
	)

	froot, err := os.Create(macro)
	if err != nil {
		log.Fatal(err)
	}
	defer froot.Close()
	defer os.Remove(macro)
	defer os.Remove(oname)

	tmpl := template.Must(template.New("genstreamers").Parse(`
void genstreamers(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");

{{range .}}
	(({{.}}*)(TClass::GetClass("{{.}}")->New()))->Write("type-{{.}}");
{{end }}

	f->Write();
	f->Close();

	exit(0);
}
`))
	err = tmpl.Execute(froot, classes)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("root.exe", "-b", fmt.Sprintf("./%s(%q)", macro, oname))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	root, err := rootio.Open(oname)
	if err != nil {
		log.Fatalf("could not open ROOT streamers file: %v", err)
	}
	err = root.Close()
	if err != nil {
		log.Fatalf("could not close ROOT streamers file: %v", err)
	}

	raw, err := ioutil.ReadFile(oname)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("internal/rstreamers/pkg_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, `// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rstreamers

import (
	"encoding/base64"
	"fmt"
)

var Data []byte

func init() {
	var err error
	Data, err = base64.StdEncoding.DecodeString(`,
	)

	fmt.Fprintf(f, "`%s`)\n", base64.StdEncoding.EncodeToString(raw))
	fmt.Fprintf(f, `
	if err != nil {
		panic(fmt.Errorf("rootio: could not decode embedded streamer: %%v", err))
	}
}
`)

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}

const srcHeader = `// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"reflect"
`

const leafTmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	rvers int16
	tleaf
	val []{{.Type}}
	min {{.RangeType}}
	max {{.RangeType}}
}

// Class returns the ROOT class name.
func (leaf *{{.Name}}) Class() string {
	return "T{{.Name}}"
}

// Minimum returns the minimum value of the leaf.
func (leaf *{{.Name}}) Minimum() {{.RangeType}} {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *{{.Name}}) Maximum() {{.RangeType}} {
	return leaf.max
}

// Kind returns the leaf's kind.
func (*{{.Name}}) Kind() reflect.Kind {
	return {{.Kind}}
}

// Type returns the leaf's type.
func (*{{.Name}}) Type() reflect.Type {
	var v {{.Type}}
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *{{.Name}}) Value(i int) interface{} {
	return leaf.val[i]
}

// value returns the leaf value.
func (leaf *{{.Name}}) value() interface{} {
	return leaf.val
}

{{if .Count}}
// ivalue returns the first leaf value as int
func (leaf *{{.Name}}) ivalue() int {
	return int(leaf.val[0])
}

// imax returns the leaf maximum value as int
func (leaf *{{.Name}}) imax() int {
	return int(leaf.max)
}
{{end}}

func (leaf *{{.Name}}) TypeName() string {
	return "{{.Type}}"
}

func (leaf *{{.Name}}) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	{{.WRangeFunc}}(leaf.min)
	{{.WRangeFunc}}(leaf.max)

	return w.SetByteCount(pos, "T{{.Name}}")
}

func (leaf *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.min = {{.RRangeFunc}}
	leaf.max = {{.RRangeFunc}}

	r.CheckByteCount(pos, bcnt, start, "T{{.Name}}")
	return r.Err()
}

func (leaf *{{.Name}}) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && len(leaf.val) == 1 {
		leaf.val[0] = {{.RFunc}}
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.val = {{.RFuncArray}}(leaf.tleaf.len * n)
		} else {
			leaf.val = {{.RFuncArray}}(leaf.tleaf.len)
		}
	}
	return r.err
}

func (leaf *{{.Name}}) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *{{.Type}}:
		*v = leaf.val[0]
	case *[]{{.Type}}:
		if len(*v) < len(leaf.val) || *v == nil {
			*v = make([]{{.Type}}, len(leaf.val))
		}
		copy(*v, leaf.val)
		*v = (*v)[:leaf.count.ivalue()]
	case []{{.Type}}:
		copy(v, leaf.val)
{{if .DoUnsigned}}
	case *u{{.Type}}:
		*v = u{{.Type}}(leaf.val[0])
	case *[]u{{.Type}}:
		if len(*v) < len(leaf.val) || *v == nil {
			*v = make([]u{{.Type}}, len(leaf.val))
		}
		for i, u := range leaf.val {
			(*v)[i] = u{{.Type}}(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []u{{.Type}}:
		for i := range v {
			v[i] = u{{.Type}}(leaf.val[i])
		}
{{end}}
	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &{{.Name}}{}
		return reflect.ValueOf(o)
	}
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var (
	_ Object          = (*{{.Name}})(nil)
	_ Named           = (*{{.Name}})(nil)
	_ Leaf            = (*{{.Name}})(nil)
	_ ROOTMarshaler   = (*{{.Name}})(nil)
	_ ROOTUnmarshaler = (*{{.Name}})(nil)
)
`

const arrayTmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	Data []{{.Type}}
}

// Class returns the ROOT class name.
func (*{{.Name}}) Class() string {
	return "T{{.Name}}"
}

func (arr *{{.Name}}) Len() int {
	return len(arr.Data)
}

func (arr *{{.Name}}) At(i int) {{.Type}} {
	return arr.Data[i]
}

func (arr *{{.Name}}) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *{{.Name}}) Set(i int, v interface{}) {
	arr.Data[i] = v.({{.Type}})
}

func (arr *{{.Name}}) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	{{.WFunc}}(arr.Data)

	return int(w.Pos()-pos), w.err
}

func (arr *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	n := int(r.ReadI32())
	arr.Data = {{.RFunc}}(n)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &{{.Name}}{}
		return reflect.ValueOf(o)
	}
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var (
	_ Array = (*{{.Name}})(nil)
	_ ROOTMarshaler = (*{{.Name}})(nil)
	_ ROOTUnmarshaler = (*{{.Name}})(nil)
)
`

const h1Tmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	rvers int16
	th1
	arr {{.Type}}
}

func new{{.Name}}() *{{.Name}} {
	return &{{.Name}}{
		rvers: 2, // FIXME(sbinet): harmonize versions
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

	hroot.th1.name = h.Name()
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th1.title = v.(string)
	}
	hroot.th1.xaxis.xbins.Data = edges
	return hroot
}


func (*{{.Name}}) isH1() {}

// Class returns the ROOT class name.
func (*{{.Name}}) Class() string {
	return "T{{.Name}}"
}

func (h *{{.Name}}) MarshalROOT(w *WBuffer) (int, error) {
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

	return w.SetByteCount(pos, "T{{.Name}}")
}

func (h *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: T{{.Name}} version too old (%d<1)", vers)
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

	r.CheckByteCount(pos, bcnt, beg, "T{{.Name}}")
	return r.err
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
			SumWX:  float64(h.SumWX()),
			SumWX2: float64(h.SumWX2()),
		}
		dists = make([]hbook.Dist1D, int(nx))
	)

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
		xmin := h.XBinLowEdge(i+1)
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
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var (
	_ Object          = (*{{.Name}})(nil)
	_ Named           = (*{{.Name}})(nil)
	_ H1              = (*{{.Name}})(nil)
	_ ROOTMarshaler   = (*{{.Name}})(nil)
	_ ROOTUnmarshaler = (*{{.Name}})(nil)
)
`

const h2Tmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	rvers int16
	th2
	arr {{.Type}}
}

func new{{.Name}}() *{{.Name}} {
	return &{{.Name}}{
		rvers: 3, // FIXME(sbinet): harmonize versions
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

	hroot.th2.th1.name = h.Name()
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th2.th1.title = v.(string)
	}
	hroot.th2.th1.xaxis.xbins.Data = xedges
	hroot.th2.th1.yaxis.xbins.Data = yedges

	return hroot
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
				SumWX:  float64(h.SumWX()),
				SumWX2: float64(h.SumWX2()),
			},
			Y: hbook.Dist1D{
				Dist: hbook.Dist0D{
					N:      int64(h.Entries()),
					SumW:   float64(h.SumW()),
					SumW2:  float64(h.SumW2()),
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
			i := iy*nx+ix
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

func (h *{{.Name}}) MarshalROOT(w *WBuffer) (int, error) {
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

	return w.SetByteCount(pos, "T{{.Name}}")
}

func (h *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 1 {
		return errorf("rootio: T{{.Name}} version too old (%d<1)", vers)
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

	r.CheckByteCount(pos, bcnt, beg, "T{{.Name}}")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := new{{.Name}}()
		return reflect.ValueOf(o)
	}
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var (
	_ Object          = (*{{.Name}})(nil)
	_ Named           = (*{{.Name}})(nil)
	_ H2              = (*{{.Name}})(nil)
	_ ROOTMarshaler   = (*{{.Name}})(nil)
	_ ROOTUnmarshaler = (*{{.Name}})(nil)
)
`
