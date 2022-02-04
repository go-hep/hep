// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"text/template"

	"go-hep.org/x/hep/groot/internal/genroot"
	"go-hep.org/x/hep/groot/root"
)

func main() {
	genLeaves()
	genRLeaves()
}

func genLeaves() {
	fname := "./rtree/leaf_gen.go"
	year := genroot.ExtractYear(fname)
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports(year, "rtree", f,
		"fmt",
		"reflect",
		"strings",
		"unsafe", // FIXME(sbinet): needed for signed/unsigned handling
		"",
		"go-hep.org/x/hep/groot/root",
		"go-hep.org/x/hep/groot/rbase",
		"go-hep.org/x/hep/groot/rbytes",
		"go-hep.org/x/hep/groot/rdict",
		"go-hep.org/x/hep/groot/rmeta",
		"go-hep.org/x/hep/groot/rtypes",
		"go-hep.org/x/hep/groot/rvers",
	)

	for i, typ := range []struct {
		Name       string
		Type       string
		Kind       string
		UKind      string
		LenType    int
		GoLenType  int
		DoUnsigned bool
		RFunc      string
		RFuncArray string
		ResizeFunc string
		WFunc      string
		WFuncArray string
		RangeType  string
		RRangeFunc string
		WRangeFunc string
		Count      bool

		WithStreamerElement bool   // for TLeaf{F16,D32}
		Meta                string // name of rmeta.Enum to use (for TLeaf{F16,D32})
	}{
		{
			Name:       "LeafO",
			Type:       "bool",
			Kind:       "reflect.Bool",
			LenType:    1,
			GoLenType:  int(reflect.TypeOf(true).Size()),
			RFunc:      "r.ReadBool()",
			RFuncArray: "r.ReadArrayBool",
			ResizeFunc: "rbytes.ResizeBool",
			WFunc:      "w.WriteBool",
			WFuncArray: "w.WriteArrayBool",
		},
		{
			Name:       "LeafB",
			Type:       "int8",
			Kind:       "reflect.Int8",
			UKind:      "reflect.Uint8",
			LenType:    1,
			GoLenType:  int(reflect.TypeOf(int8(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI8()",
			RFuncArray: "r.ReadArrayI8",
			ResizeFunc: "rbytes.ResizeI8",
			WFunc:      "w.WriteI8",
			WFuncArray: "w.WriteArrayI8",
			Count:      true,
		},
		{
			Name:       "LeafS",
			Type:       "int16",
			Kind:       "reflect.Int16",
			UKind:      "reflect.Uint16",
			LenType:    2,
			GoLenType:  int(reflect.TypeOf(int16(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI16()",
			RFuncArray: "r.ReadArrayI16",
			ResizeFunc: "rbytes.ResizeI16",
			WFunc:      "w.WriteI16",
			WFuncArray: "w.WriteArrayI16",
			Count:      true,
		},
		{
			Name:       "LeafI",
			Type:       "int32",
			Kind:       "reflect.Int32",
			UKind:      "reflect.Uint32",
			LenType:    4,
			GoLenType:  int(reflect.TypeOf(int32(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI32()",
			RFuncArray: "r.ReadArrayI32",
			ResizeFunc: "rbytes.ResizeI32",
			WFunc:      "w.WriteI32",
			WFuncArray: "w.WriteArrayI32",
			Count:      true,
		},
		{
			Name:       "LeafL",
			Type:       "int64",
			Kind:       "reflect.Int64",
			UKind:      "reflect.Uint64",
			LenType:    8,
			GoLenType:  int(reflect.TypeOf(int64(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI64()",
			RFuncArray: "r.ReadArrayI64",
			ResizeFunc: "rbytes.ResizeI64",
			WFunc:      "w.WriteI64",
			WFuncArray: "w.WriteArrayI64",
			Count:      true,
		},
		{
			Name:       "LeafF",
			Type:       "float32",
			Kind:       "reflect.Float32",
			LenType:    4,
			GoLenType:  int(reflect.TypeOf(float32(0)).Size()),
			RFunc:      "r.ReadF32()",
			RFuncArray: "r.ReadArrayF32",
			ResizeFunc: "rbytes.ResizeF32",
			WFunc:      "w.WriteF32",
			WFuncArray: "w.WriteArrayF32",
		},
		{
			Name:       "LeafD",
			Type:       "float64",
			Kind:       "reflect.Float64",
			LenType:    8,
			GoLenType:  int(reflect.TypeOf(float64(0)).Size()),
			RFunc:      "r.ReadF64()",
			RFuncArray: "r.ReadArrayF64",
			ResizeFunc: "rbytes.ResizeF64",
			WFunc:      "w.WriteF64",
			WFuncArray: "w.WriteArrayF64",
		},
		{
			Name:                "LeafF16",
			Type:                "root.Float16",
			Kind:                "reflect.Float32",
			LenType:             4,
			GoLenType:           int(reflect.TypeOf(root.Float16(0)).Size()),
			RFunc:               "r.ReadF16(leaf.elm)",
			RFuncArray:          "r.ReadArrayF16",
			ResizeFunc:          "rbytes.ResizeF16",
			WFunc:               "w.WriteF16",
			WFuncArray:          "w.WriteArrayF16",
			WithStreamerElement: true,
			Meta:                "rmeta.Float16",
		},
		{
			Name:                "LeafD32",
			Type:                "root.Double32",
			Kind:                "reflect.Float64",
			LenType:             8,
			GoLenType:           int(reflect.TypeOf(root.Double32(0)).Size()),
			RFunc:               "r.ReadD32(leaf.elm)",
			RFuncArray:          "r.ReadArrayD32",
			ResizeFunc:          "rbytes.ResizeD32",
			WFunc:               "w.WriteD32",
			WFuncArray:          "w.WriteArrayD32",
			WithStreamerElement: true,
			Meta:                "rmeta.Double32",
		},
		{
			Name:       "LeafC",
			Type:       "string",
			Kind:       "reflect.String",
			LenType:    1,
			GoLenType:  int(reflect.TypeOf("").Size()),
			RFunc:      "r.ReadString()",
			RFuncArray: "r.ReadArrayString",
			ResizeFunc: "rbytes.ResizeStr",
			WFunc:      "w.WriteString",
			WFuncArray: "w.WriteArrayString",
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
	genroot.GoFmt(f)
}

const leafTmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	rvers int16
	tleaf
	ptr *{{.Type}}
	sli *[]{{.Type}}
	min {{.RangeType}}
	max {{.RangeType}}
{{- if .WithStreamerElement}}
	elm rbytes.StreamerElement
{{- end}}
}

{{- if .WithStreamerElement}}
func new{{.Name}}(b Branch, name string, shape []int, unsigned bool, count Leaf, elm rbytes.StreamerElement) *{{.Name}} {
	const etype = {{.LenType}}
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &{{.Name}}{
		rvers: rvers.{{.Name}},
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
		elm:   elm,
	}
}
{{- else}}
func new{{.Name}}(b Branch, name string, shape []int, unsigned bool, count Leaf) *{{.Name}} {
	const etype = {{.LenType}}
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &{{.Name}}{
		rvers: rvers.{{.Name}},
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}
{{- end}}

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
func (leaf *{{.Name}}) Kind() reflect.Kind {
{{- if .DoUnsigned}}
	if leaf.IsUnsigned() {
		return {{.UKind}}
	}
	return {{.Kind}}
{{- else}}
	return {{.Kind}}
{{- end}}
}

// Type returns the leaf's type.
func (leaf *{{.Name}}) Type() reflect.Type {
{{- if .DoUnsigned}}
	if leaf.IsUnsigned() {
		var v u{{.Type}}
		return reflect.TypeOf(v)
	}
	var v {{.Type}}
	return reflect.TypeOf(v)
{{- else}}
	var v {{.Type}}
	return reflect.TypeOf(v)
{{- end}}
}

{{- if .Count}}
// ivalue returns the first leaf value as int
func (leaf *{{.Name}}) ivalue() int {
	return int(*leaf.ptr)
}

// imax returns the leaf maximum value as int
func (leaf *{{.Name}}) imax() int {
	return int(leaf.max)
}
{{- end}}

func (leaf *{{.Name}}) TypeName() string {
{{- if .DoUnsigned}}
	if leaf.IsUnsigned() {
		return "u{{.Type}}"
	}
	return "{{.Type}}"
{{- else}}
	return "{{.Type}}"
{{- end}}
}

func (leaf *{{.Name}}) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
	w.WriteObject(&leaf.tleaf)
{{- if .WithStreamerElement}}
	{{.WRangeFunc}}(leaf.min, leaf.elm)
	{{.WRangeFunc}}(leaf.max, leaf.elm)
{{- else}}
	{{.WRangeFunc}}(leaf.min)
	{{.WRangeFunc}}(leaf.max)
{{- end}}

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *{{.Name}}) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = {{.RRangeFunc}}
	leaf.max = {{.RRangeFunc}}

{{if .WithStreamerElement}}
	if strings.Contains(leaf.Title(), "[") {
		elm := rdict.Element{
			Name:   *rbase.NewNamed(fmt.Sprintf("%s_Element", leaf.Name()), leaf.Title()),
			Offset: 0,
			Type:   {{.Meta}},
		}.New()
		leaf.elm = &elm
	}
{{- end}}

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
	return r.Err()
}

func (leaf *{{.Name}}) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = {{.RFunc}}
	} else {
            if leaf.count != nil {
                n := leaf.count.ivalue()
                max := leaf.count.imax()
                if n > max {
                        n = max
                }
				nn := leaf.tleaf.len * n
				*leaf.sli = {{.ResizeFunc}}(*leaf.sli, nn)
{{- if .WithStreamerElement}}
                {{.RFuncArray}}(*leaf.sli, leaf.elm)
{{- else}}
                {{.RFuncArray}}(*leaf.sli)
{{- end}}
            } else {
				nn := leaf.tleaf.len
				*leaf.sli = {{.ResizeFunc}}(*leaf.sli, nn)
{{- if .WithStreamerElement}}
				{{.RFuncArray}}(*leaf.sli, leaf.elm)
{{- else}}
				{{.RFuncArray}}(*leaf.sli)
{{- end}}
            }
    }
    return r.Err()
}

func (leaf *{{.Name}}) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / {{.GoLenType}}
	arr := (*[0]{{.Type}})(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *{{.Name}}) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

    if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]{{.Type}}:
			return leaf.setAddress(sli)
{{- if .DoUnsigned}}
		case *[]u{{.Type}}:
			return leaf.setAddress(sli)
{{- end}}
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
    }

	switch v := ptr.(type) {
    case *{{.Type}}:
		leaf.ptr = v
    case *[]{{.Type}}:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]{{.Type}}, 0)
		}
{{- if .DoUnsigned}}
    case *u{{.Type}}:
		leaf.ptr = (*{{.Type}})(unsafe.Pointer(v))
    case *[]u{{.Type}}:
		leaf.sli = (*[]{{.Type}})(unsafe.Pointer(v))
		if *v == nil {
			*leaf.sli = make([]{{.Type}}, 0)
		}
{{- end}}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *{{.Name}}) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
{{- if .WithStreamerElement}}
		{{.WFunc}}(*leaf.ptr, leaf.elm)
{{- else}}
		{{.WFunc}}(*leaf.ptr)
{{- end}}
{{- if eq .Name "LeafC"}}
		sz := len(*leaf.ptr)
		nbytes += sz
		if v := int32(sz); v >= leaf.max {
			leaf.max = v+1
		}
		if sz >= leaf.tleaf.len {
			leaf.tleaf.len = sz+1
		}
{{- else if eq .Name "LeafO"}}
		nbytes += leaf.tleaf.etype
{{- else}}
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
{{- end}}
	case leaf.count != nil:
		n := leaf.count.ivalue()
        max := leaf.count.imax()
        if n > max {
			n = max
		}
		end := leaf.tleaf.len*n
{{- if .WithStreamerElement}}
		{{.WFuncArray}}((*leaf.sli)[:end], leaf.elm)
{{- else}}
		{{.WFuncArray}}((*leaf.sli)[:end])
{{- end}}
		nbytes += leaf.tleaf.etype * end
	default:
{{- if .WithStreamerElement}}
		{{.WFuncArray}}((*leaf.sli)[:leaf.tleaf.len], leaf.elm)
{{- else}}
		{{.WFuncArray}}((*leaf.sli)[:leaf.tleaf.len])
{{- end}}
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &{{.Name}}{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("T{{.Name}}", f)
}

var (
	_ root.Object        = (*{{.Name}})(nil)
	_ root.Named         = (*{{.Name}})(nil)
	_ Leaf               = (*{{.Name}})(nil)
	_ rbytes.Marshaler   = (*{{.Name}})(nil)
	_ rbytes.Unmarshaler = (*{{.Name}})(nil)
)
`

func genRLeaves() {
	fname := "./rtree/rleaf_gen.go"
	year := genroot.ExtractYear(fname)
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	type Kind string
	const (
		Val = Kind("Val")
		Arr = Kind("Arr")
		Sli = Kind("Sli")
	)

	genroot.GenImports(year, "rtree", f,
		"reflect",
		"unsafe", // for unsafeDecayArrayXXX
		"",
		"go-hep.org/x/hep/groot/rbytes",
		"go-hep.org/x/hep/groot/root",
	)

	for i, typ := range []struct {
		Name  string
		Base  string
		Type  string
		Size  int
		Kind  Kind
		Func  string
		Decay string
		Count bool

		WithStreamerElement bool // for TLeaf{F16,D32}
	}{
		{
			Name: "Bool",
			Base: "LeafO",
			Type: "bool",
			Size: int(reflect.TypeOf(true).Size()),
		},
		{
			Name:  "I8",
			Base:  "LeafB",
			Type:  "int8",
			Size:  int(reflect.TypeOf(int8(0)).Size()),
			Count: true,
		},
		{
			Name:  "I16",
			Base:  "LeafS",
			Type:  "int16",
			Size:  int(reflect.TypeOf(int16(0)).Size()),
			Count: true,
		},
		{
			Name:  "I32",
			Base:  "LeafI",
			Type:  "int32",
			Size:  int(reflect.TypeOf(int32(0)).Size()),
			Count: true,
		},
		{
			Name:  "I64",
			Base:  "LeafL",
			Type:  "int64",
			Size:  int(reflect.TypeOf(int64(0)).Size()),
			Count: true,
		},
		{
			Name:  "U8",
			Base:  "LeafB",
			Type:  "uint8",
			Size:  int(reflect.TypeOf(uint8(0)).Size()),
			Count: true,
		},
		{
			Name:  "U16",
			Base:  "LeafS",
			Type:  "uint16",
			Size:  int(reflect.TypeOf(uint16(0)).Size()),
			Count: true,
		},
		{
			Name:  "U32",
			Base:  "LeafI",
			Type:  "uint32",
			Size:  int(reflect.TypeOf(uint32(0)).Size()),
			Count: true,
		},
		{
			Name:  "U64",
			Base:  "LeafL",
			Type:  "uint64",
			Size:  int(reflect.TypeOf(uint64(0)).Size()),
			Count: true,
		},
		{
			Name: "F32",
			Base: "LeafF",
			Type: "float32",
			Size: int(reflect.TypeOf(float32(0)).Size()),
		},
		{
			Name: "F64",
			Base: "LeafD",
			Type: "float64",
			Size: int(reflect.TypeOf(float64(0)).Size()),
		},
		{
			Name: "D32",
			Base: "LeafD32",
			Type: "root.Double32",
			Size: int(reflect.TypeOf(root.Double32(0)).Size()),

			WithStreamerElement: true,
		},
		{
			Name: "F16",
			Base: "LeafF16",
			Type: "root.Float16",
			Size: int(reflect.TypeOf(root.Float16(0)).Size()),

			WithStreamerElement: true,
		},
		{
			Name: "Str",
			Base: "LeafC",
			Type: "string",
			Size: int(reflect.TypeOf("").Size()),
		},
	} {
		for j, kind := range []Kind{Val, Arr, Sli} {
			if i > 0 || j > 0 {
				fmt.Fprintf(f, "\n")
			}
			typ.Kind = kind
			switch typ.Name {
			case "Str":
				typ.Func = "String"
			default:
				typ.Func = typ.Name
			}

			switch typ.Type[0] {
			case 'u':
				typ.Decay = "unsafeDecayArrayU"
			default:
				typ.Decay = "unsafeDecayArray"
			}

			tmpl := template.Must(template.New(typ.Name).Parse(rleafTmpl))
			err = tmpl.Execute(f, typ)
			if err != nil {
				log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
			}
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)
}

const rleafTmpl = `// rleaf{{.Kind}}{{.Name}} implements rleaf for ROOT T{{.Base}}
type rleaf{{.Kind}}{{.Name}} struct {
	base *{{.Base}}
{{- if eq .Kind "Val" }}
	v *{{.Type}}
{{- else if eq .Kind "Arr" }}
	v []{{.Type}}
{{- else if eq .Kind "Sli" }}
	n   func() int
	v   *[]{{.Type}}
	set func() // reslice underlying slice
{{- end}}
{{- if .WithStreamerElement}}
	elm rbytes.StreamerElement
{{- end}}
}

{{- if eq .Kind "Val"}}
func newRLeaf{{.Name}}(leaf *{{.Base}}, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		switch len(leaf.Shape()) {
		case 0:
			slice := reflect.ValueOf(rvar.Value).Interface().(*[]{{.Type}})
			if *slice == nil {
				*slice = make([]{{.Type}}, 0, rleafDefaultSliceCap)
			}
			return &rleafSli{{.Name}}{
				base: leaf,
				n:    rctx.rcountFunc(leaf.count.Name()),
				v:    slice,
				set:  func() {},
			}
		default:
			sz := 1
			for _, v := range leaf.Shape() {
				sz *= v
			}
			sli := reflect.ValueOf(rvar.Value).Elem()
			ptr := (*[]{{.Type}})(unsafe.Pointer(sli.UnsafeAddr()))
			hdr := unsafeDecaySliceArray{{.Name}}(ptr, sz).(*[]{{.Type}})
			if *hdr == nil {
				*hdr = make([]{{.Type}}, 0, rleafDefaultSliceCap*sz)
			}
			rleaf := &rleafSli{{.Name}}{
				base: leaf,
				n:    rctx.rcountFunc(leaf.count.Name()),
				v:    hdr,
			}
			rawSli := (*reflect.SliceHeader)(unsafe.Pointer(sli.UnsafeAddr()))
			rawHdr := (*reflect.SliceHeader)(unsafe.Pointer(hdr))

			// alias slices
			rawSli.Data = rawHdr.Data

			rleaf.set = func() {
				n := rleaf.n()
				rawSli.Len = n
				rawSli.Cap = n
			}

			return rleaf
}

	case leaf.len > 1:
		return &rleafArr{{.Name}}{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArray{{.Name}}(rvar.Value)).Elem().Interface().([]{{.Type}}),
		}

	default:
		return &rleafVal{{.Name}}{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*{{.Type}}),
		}
	}
}
{{- end}}

func (leaf *rleaf{{.Kind}}{{.Name}}) Leaf() Leaf { return leaf.base }

func (leaf *rleaf{{.Kind}}{{.Name}}) Offset() int64 {
	return int64(leaf.base.Offset())
}

{{if eq .Kind "Val"}}
{{if .Count}}
func (leaf *rleaf{{.Kind}}{{.Name}}) ivalue() int { return int(*leaf.v) }
{{- end}}
{{- end}}

{{if eq .Kind "Arr"}}
func unsafeDecayArray{{.Name}}(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / {{.Size}}
	arr := (*[0]{{.Type}})(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func unsafeDecaySliceArray{{.Name}}(ptr *[]{{.Type}}, size int) interface{} {
	var sli []{{.Type}}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(size)
	hdr.Cap = int(size)
	hdr.Data = (*reflect.SliceHeader)(unsafe.Pointer(ptr)).Data
	return &sli
}
{{- end}}

{{if .WithStreamerElement}}
func (leaf *rleaf{{.Kind}}{{.Name}}) readFromBuffer(r *rbytes.RBuffer) error {
{{- if eq .Kind "Val" }}
	*leaf.v = r.Read{{.Func}}(leaf.elm)
{{- else if eq .Kind "Arr" }}
	r.ReadArray{{.Func}}(leaf.v, leaf.elm)
{{- else if eq .Kind "Sli" }}
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.Resize{{.Name}}(*leaf.v, n)
	r.ReadArray{{.Func}}(sli, leaf.elm)
	*leaf.v = sli
	leaf.set()
{{- end}}
	return r.Err()
}
{{else}}
func (leaf *rleaf{{.Kind}}{{.Name}}) readFromBuffer(r *rbytes.RBuffer) error {
{{- if eq .Kind "Val" }}
	*leaf.v = r.Read{{.Func}}()
{{- else if eq .Kind "Arr" }}
	r.ReadArray{{.Func}}(leaf.v)
{{- else if eq .Kind "Sli" }}
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.Resize{{.Name}}(*leaf.v, n)
	r.ReadArray{{.Func}}(sli)
	*leaf.v = sli
	leaf.set()
{{- end}}
	return r.Err()
}
{{- end}}

var (
	_ rleaf = (*rleaf{{.Kind}}{{.Name}})(nil)
)
`
