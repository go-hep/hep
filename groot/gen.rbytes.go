// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"go-hep.org/x/hep/groot/internal/genroot"
)

func main() {
	genRBuffer()
	genWBuffer()
}

func genRBuffer() {
	fname := "./rbytes/rbuffer_gen.go"
	year := genroot.ExtractYear(fname)
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports(year, "rbytes", f,
		"encoding/binary",
		"fmt",
		"math",
		"",
		"go-hep.org/x/hep/groot/rvers",
	)

	for i, typ := range []struct {
		Name     string
		Type     string
		Size     int
		REndian  string
		Frombits string
	}{
		{
			Name:     "U16",
			Type:     "uint16",
			Size:     2,
			REndian:  "binary.BigEndian.Uint16",
			Frombits: "",
		},
		{
			Name:     "U32",
			Type:     "uint32",
			Size:     4,
			REndian:  "binary.BigEndian.Uint32",
			Frombits: "",
		},
		{
			Name:     "U64",
			Type:     "uint64",
			Size:     8,
			REndian:  "binary.BigEndian.Uint64",
			Frombits: "",
		},
		{
			Name:     "I16",
			Type:     "int16",
			Size:     2,
			REndian:  "binary.BigEndian.Uint16",
			Frombits: "int16",
		},
		{
			Name:     "I32",
			Type:     "int32",
			Size:     4,
			REndian:  "binary.BigEndian.Uint32",
			Frombits: "int32",
		},
		{
			Name:     "I64",
			Type:     "int64",
			Size:     8,
			REndian:  "binary.BigEndian.Uint64",
			Frombits: "int64",
		},
		{
			Name:     "F32",
			Type:     "float32",
			Size:     4,
			REndian:  "binary.BigEndian.Uint32",
			Frombits: "math.Float32frombits",
		},
		{
			Name:     "F64",
			Type:     "float64",
			Size:     8,
			REndian:  "binary.BigEndian.Uint64",
			Frombits: "math.Float64frombits",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(rbufferTmpl))
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

const rbufferTmpl = `func (r *RBuffer) ReadArray{{.Name}}(sli []{{.Type}}) {
	if r.err != nil {
		return
	}
	n := len(sli)
	if n <= 0 || int64(n) > r.Len() {
		return
	}

	cur := r.r.c
	end := r.r.c + {{.Size}}*len(sli)
	sub := r.r.p[cur:end]
	cur = 0
	for i := range sli {
		beg := cur
		end := cur + {{.Size}}
		cur = end
		v := {{.REndian}}(sub[beg:end])
{{- if eq .Frombits ""}}
		sli[i] = v
{{else}}
		sli[i] = {{.Frombits}}(v)
{{- end}}
	}
	r.r.c = end
}

func (r *RBuffer) Read{{.Name}}() {{.Type}} {
	if r.err != nil {
		return 0
	}
	beg := r.r.c
	r.r.c += {{.Size}}
	v := {{.REndian}}(r.r.p[beg:r.r.c])
{{- if eq .Frombits ""}}
	return v
{{else}}
	return {{.Frombits}}(v)
{{- end}}
}

func (r *RBuffer) ReadStdVector{{.Name}}(sli *[]{{.Type}}) {
	if r.err != nil {
		return
	}
	
	hdr := r.ReadHeader("vector<{{.Type}}>", rvers.StreamerBaseSTL)
	if hdr.Vers > rvers.StreamerBaseSTL {
		r.err = fmt.Errorf(
			"rbytes: invalid %s version: got=%d, want=%d",
			hdr.Name, hdr.Vers, rvers.StreamerBaseSTL,
		)
		return
	}
	n := int(r.ReadI32())
	*sli = Resize{{.Name}}(*sli, n)
	for i := range *sli {
		(*sli)[i] = r.Read{{.Name}}()
	}

	r.CheckHeader(hdr)
}
`

func genWBuffer() {
	fname := "./rbytes/wbuffer_gen.go"
	year := genroot.ExtractYear(fname)
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports(year, "rbytes", f,
		"encoding/binary",
		"math",
		"",
		"go-hep.org/x/hep/groot/rvers",
	)

	for i, typ := range []struct {
		Name    string
		Type    string
		Size    int
		WEndian string
		Tobits  string
	}{
		{
			Name:    "U16",
			Type:    "uint16",
			Size:    2,
			WEndian: "binary.BigEndian.PutUint16",
		},
		{
			Name:    "U32",
			Type:    "uint32",
			Size:    4,
			WEndian: "binary.BigEndian.PutUint32",
		},
		{
			Name:    "U64",
			Type:    "uint64",
			Size:    8,
			WEndian: "binary.BigEndian.PutUint64",
		},
		{
			Name:    "I16",
			Type:    "int16",
			Size:    2,
			WEndian: "binary.BigEndian.PutUint16",
			Tobits:  "uint16",
		},
		{
			Name:    "I32",
			Type:    "int32",
			Size:    4,
			WEndian: "binary.BigEndian.PutUint32",
			Tobits:  "uint32",
		},
		{
			Name:    "I64",
			Type:    "int64",
			Size:    8,
			WEndian: "binary.BigEndian.PutUint64",
			Tobits:  "uint64",
		},
		{
			Name:    "F32",
			Type:    "float32",
			Size:    4,
			WEndian: "binary.BigEndian.PutUint32",
			Tobits:  "math.Float32bits",
		},
		{
			Name:    "F64",
			Type:    "float64",
			Size:    8,
			WEndian: "binary.BigEndian.PutUint64",
			Tobits:  "math.Float64bits",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(wbufferTmpl))
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

const wbufferTmpl = `func (w *WBuffer) WriteArray{{.Name}}(sli []{{.Type}}) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli)*{{.Size}})

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + {{.Size}}
		cur = end
{{- if eq .Tobits ""}}
		{{.WEndian}}(w.w.p[beg:end], v)
{{else}}
		{{.WEndian}}(w.w.p[beg:end], {{.Tobits}}(v))
{{- end}}
	}
	w.w.c += {{.Size}}*len(sli)
}

func (w *WBuffer) Write{{.Name}}(v {{.Type}}) {
	if w.err != nil {
		return
	}
	w.w.grow({{.Size}})
	beg := w.w.c
	end := w.w.c + {{.Size}}
{{- if eq .Tobits ""}}
	{{.WEndian}}(w.w.p[beg:end], v)
{{else}}
	{{.WEndian}}(w.w.p[beg:end], {{.Tobits}}(v))
{{- end}}
	w.w.c += {{.Size}}
}

func (w *WBuffer) WriteStdVector{{.Name}}(sli []{{.Type}}) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<{{.Type}}>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli)*{{.Size}})

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + {{.Size}}
		cur = end
{{- if eq .Tobits ""}}
		{{.WEndian}}(w.w.p[beg:end], v)
{{else}}
		{{.WEndian}}(w.w.p[beg:end], {{.Tobits}}(v))
{{- end}}
	}
	w.w.c += {{.Size}}*len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}
`
