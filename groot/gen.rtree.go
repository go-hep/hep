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
	genLeaves()
}

func genLeaves() {
	f, err := os.Create("./rtree/leaf_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rtree", f,
		"reflect",
		"unsafe", // FIXME(sbinet): needed for signed/unsigned handling
		"",
		"github.com/pkg/errors",
		"",
		"go-hep.org/x/hep/groot/root",
		"go-hep.org/x/hep/groot/rbytes",
		"go-hep.org/x/hep/groot/rtypes",
		"go-hep.org/x/hep/groot/rvers",
	)

	for i, typ := range []struct {
		Name       string
		Type       string
		Kind       string
		LenType    int
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
			LenType:    1,
			RFunc:      "r.ReadBool()",
			RFuncArray: "r.ReadFastArrayBool",
			WFunc:      "w.WriteBool",
			WFuncArray: "w.WriteFastArrayBool",
		},
		{
			Name:       "LeafB",
			Type:       "int8",
			Kind:       "reflect.Int8",
			LenType:    1,
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
			LenType:    2,
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
			LenType:    4,
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
			LenType:    8,
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
			LenType:    4,
			RFunc:      "r.ReadF32()",
			RFuncArray: "r.ReadFastArrayF32",
			WFunc:      "w.WriteF32",
			WFuncArray: "w.WriteFastArrayF32",
		},
		{
			Name:       "LeafD",
			Type:       "float64",
			Kind:       "reflect.Float64",
			LenType:    8,
			RFunc:      "r.ReadF64()",
			RFuncArray: "r.ReadFastArrayF64",
			WFunc:      "w.WriteF64",
			WFuncArray: "w.WriteFastArrayF64",
		},
		{
			Name:       "LeafC",
			Type:       "string",
			Kind:       "reflect.String",
			LenType:    1,
			RFunc:      "r.ReadString()",
			RFuncArray: "r.ReadFastArrayString",
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
}

func new{{.Name}}(b Branch, name string, len int, unsigned bool, count Leaf) *{{.Name}} {
	const etype = {{.LenType}}
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &{{.Name}}{
		rvers: rvers.{{.Name}},
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
	}
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
	{{.WRangeFunc}}(leaf.min)
	{{.WRangeFunc}}(leaf.max)

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
                        entry := leaf.Branch().getReadEntry()
                        if leaf.count.Branch().getReadEntry() != entry {
                                leaf.count.Branch().getEntry(entry)
                        }
                        n := leaf.count.ivalue()
                        max := leaf.count.imax()
                        if n > max {
                                n = max
                        }
                        *leaf.sli = {{.RFuncArray}}(leaf.tleaf.len * n)
                } else {
						copy(*leaf.sli, {{.RFuncArray}}(leaf.tleaf.len))
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
            panic(errors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
    }

    return r.Err()
}

func (leaf *{{.Name}}) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

    if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli :=arr.Slice(0, rv.Len()).Interface().(type) {
		case []{{.Type}}:
			return leaf.setAddress(&sli)
{{- if .DoUnsigned}}
		case []u{{.Type}}:
			return leaf.setAddress(&sli)
{{- end}}
		default:
			panic(errors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
    }

	switch v := ptr.(type) {
    case *{{.Type}}:
		leaf.ptr = v
    case *[]{{.Type}}:
		leaf.sli = v
{{- if .DoUnsigned}}
    case *u{{.Type}}:
		leaf.ptr = (*{{.Type}})(unsafe.Pointer(v))
    case *[]u{{.Type}}:
		leaf.sli = (*[]{{.Type}})(unsafe.Pointer(v))
{{- end}}
	default:
		panic(errors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		{{.WFunc}}(*leaf.ptr)
		nbytes += leaf.tleaf.etype
{{- if eq .Name "LeafC"}}
		sz := len(*leaf.ptr)
		if v := int32(sz); v >= leaf.max {
			leaf.max = v+1
		}
		if sz >= leaf.tleaf.len {
			leaf.tleaf.len = sz+1
		}
{{- else if eq .Name "LeafO"}}
{{- else}}
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
		{{.WFuncArray}}((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		{{.WFuncArray}}((*leaf.sli)[:leaf.tleaf.len])
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
