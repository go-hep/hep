// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	f, err := os.Create("./rtree/leaf_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rtree", f,
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
			WFuncArray: "w.WriteFastArrayBool",
		},
		{
			Name:       "LeafB",
			Type:       "int8",
			Kind:       "reflect.Int8",
			LenType:    1,
			GoLenType:  int(reflect.TypeOf(int8(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI8()",
			RFuncArray: "r.ReadArrayI8",
			ResizeFunc: "rbytes.ResizeI8",
			WFunc:      "w.WriteI8",
			WFuncArray: "w.WriteFastArrayI8",
			Count:      true,
		},
		{
			Name:       "LeafS",
			Type:       "int16",
			Kind:       "reflect.Int16",
			LenType:    2,
			GoLenType:  int(reflect.TypeOf(int16(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI16()",
			RFuncArray: "r.ReadArrayI16",
			ResizeFunc: "rbytes.ResizeI16",
			WFunc:      "w.WriteI16",
			WFuncArray: "w.WriteFastArrayI16",
			Count:      true,
		},
		{
			Name:       "LeafI",
			Type:       "int32",
			Kind:       "reflect.Int32",
			LenType:    4,
			GoLenType:  int(reflect.TypeOf(int32(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI32()",
			RFuncArray: "r.ReadArrayI32",
			ResizeFunc: "rbytes.ResizeI32",
			WFunc:      "w.WriteI32",
			WFuncArray: "w.WriteFastArrayI32",
			Count:      true,
		},
		{
			Name:       "LeafL",
			Type:       "int64",
			Kind:       "reflect.Int64",
			LenType:    8,
			GoLenType:  int(reflect.TypeOf(int64(0)).Size()),
			DoUnsigned: true,
			RFunc:      "r.ReadI64()",
			RFuncArray: "r.ReadArrayI64",
			ResizeFunc: "rbytes.ResizeI64",
			WFunc:      "w.WriteI64",
			WFuncArray: "w.WriteFastArrayI64",
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
			WFuncArray: "w.WriteFastArrayF32",
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
			WFuncArray: "w.WriteFastArrayF64",
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
			WFuncArray:          "w.WriteFastArrayF16",
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
			WFuncArray:          "w.WriteFastArrayD32",
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
func (*{{.Name}}) Kind() reflect.Kind {
	return {{.Kind}}
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

// Value returns the leaf value at index i.
func (leaf *{{.Name}}) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *{{.Name}}) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
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
	leaf.tleaf.MarshalROOT(w)
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

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

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

func (leaf *{{.Name}}) scan(r *rbytes.RBuffer, ptr interface{}) error {
    if r.Err() != nil {
            return r.Err()
    }

    if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
            return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
    }

    switch v := ptr.(type) {
    case *{{.Type}}:
            *v = *leaf.ptr
    case *[]{{.Type}}:
            if len(*v) < len(*leaf.sli) || *v == nil {
                    *v = make([]{{.Type}}, len(*leaf.sli))
            }
            copy(*v, *leaf.sli)
            *v = (*v)[:leaf.count.ivalue()]
    case []{{.Type}}:
            copy(v, *leaf.sli)
{{- if .DoUnsigned}}
    case *u{{.Type}}:
            *v = u{{.Type}}(*leaf.ptr)
    case *[]u{{.Type}}:
            if len(*v) < len(*leaf.sli) || *v == nil {
                    *v = make([]u{{.Type}}, len(*leaf.sli))
            }
            for i, u := range (*leaf.sli) {
                    (*v)[i] = u{{.Type}}(u)
            }
            *v = (*v)[:leaf.count.ivalue()]
    case []u{{.Type}}:
            for i := range v {
                    v[i] = u{{.Type}}((*leaf.sli)[i])
            }
{{- end}}
    default:
            panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
	f, err := os.Create("./rtree/rleaf_gen.go")
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

	genroot.GenImports("rtree", f,
		"reflect",
		"",
		"go-hep.org/x/hep/groot/rbytes",
		"go-hep.org/x/hep/groot/root",
	)

	for i, typ := range []struct {
		Name string
		Base string
		Type string
		Kind Kind
		Func string

		WithStreamerElement bool // for TLeaf{F16,D32}
	}{
		{
			Name: "Bool",
			Base: "LeafO",
			Type: "bool",
		},
		{
			Name: "I8",
			Base: "LeafB",
			Type: "int8",
		},
		{
			Name: "I16",
			Base: "LeafS",
			Type: "int16",
		},
		{
			Name: "I32",
			Base: "LeafI",
			Type: "int32",
		},
		{
			Name: "I64",
			Base: "LeafL",
			Type: "int64",
		},
		{
			Name: "U8",
			Base: "LeafB",
			Type: "uint8",
		},
		{
			Name: "U16",
			Base: "LeafS",
			Type: "uint16",
		},
		{
			Name: "U32",
			Base: "LeafI",
			Type: "uint32",
		},
		{
			Name: "U64",
			Base: "LeafL",
			Type: "uint64",
		},
		{
			Name: "F32",
			Base: "LeafF",
			Type: "float32",
		},
		{
			Name: "F64",
			Base: "LeafD",
			Type: "float64",
		},
		{
			Name: "D32",
			Base: "LeafD32",
			Type: "root.Double32",

			WithStreamerElement: true,
		},
		{
			Name: "F16",
			Base: "LeafF16",
			Type: "root.Float16",

			WithStreamerElement: true,
		},
		{
			Name: "Str",
			Base: "LeafC",
			Type: "string",
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
	n func() int
	v *[]{{.Type}}
{{- end}}
{{- if .WithStreamerElement}}
	elm rbytes.StreamerElement
{{- end}}
}

{{- if eq .Kind "Val"}}
func newRLeaf{{.Name}}(leaf *{{.Base}}, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]{{.Type}})
		if *slice == nil {
			*slice = make([]{{.Type}}, 0, rleafDefaultSliceCap)
		}
		return &rleafSli{{.Name}}{
			base: leaf,
			n:    rctx.rcount(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArr{{.Name}}{
			base: leaf,
			v:    reflect.ValueOf(leaf.unsafeDecayArray(rvar.Value)).Elem().Interface().([]{{.Type}}),
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

{{if .WithStreamerElement}}
func (leaf *rleaf{{.Kind}}{{.Name}}) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}
{{- if eq .Kind "Val" }}
	*leaf.v = r.Read{{.Func}}(leaf.elm)
{{- else if eq .Kind "Arr" }}
	r.ReadArray{{.Func}}(leaf.v, leaf.elm)
{{- else if eq .Kind "Sli" }}
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.Resize{{.Name}}(*leaf.v, n)
	r.ReadArray{{.Func}}(sli, leaf.elm)
	*leaf.v = sli
{{- end}}
	return r.Err()
}
{{else}}
func (leaf *rleaf{{.Kind}}{{.Name}}) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}
{{- if eq .Kind "Val" }}
	*leaf.v = r.Read{{.Func}}()
{{- else if eq .Kind "Arr" }}
	r.ReadArray{{.Func}}(leaf.v)
{{- else if eq .Kind "Sli" }}
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.Resize{{.Name}}(*leaf.v, n)
	r.ReadArray{{.Func}}(sli)
	*leaf.v = sli
{{- end}}
	return r.Err()
}
{{- end}}

var (
	_ rleaf = (*rleaf{{.Kind}}{{.Name}})(nil)
)
`
