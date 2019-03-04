// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"go-hep.org/x/hep/groot/internal/genroot"
)

func main() {
	genArrays()
	genTClonesArrayData()
}

func genArrays() {
	f, err := os.Create("./rcont/array_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rcont", f,
		"reflect",
		"",
		"go-hep.org/x/hep/groot/root",
		"go-hep.org/x/hep/groot/rbytes",
		"go-hep.org/x/hep/groot/rtypes",
		"go-hep.org/x/hep/groot/rvers",
	)

	for i, typ := range []struct {
		Name  string
		Type  string
		RFunc string
		WFunc string
	}{
		{
			Name:  "ArrayC",
			Type:  "int8",
			RFunc: "r.ReadFastArrayI8",
			WFunc: "w.WriteFastArrayI8",
		},
		{
			Name:  "ArrayS",
			Type:  "int16",
			RFunc: "r.ReadFastArrayI16",
			WFunc: "w.WriteFastArrayI16",
		},
		{
			Name:  "ArrayI",
			Type:  "int32",
			RFunc: "r.ReadFastArrayI32",
			WFunc: "w.WriteFastArrayI32",
		},
		{
			Name:  "ArrayL",
			Type:  "int64",
			RFunc: "r.ReadFastArrayI64",
			WFunc: "w.WriteFastArrayI64",
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
	genroot.GoFmt(f)
}

const arrayTmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	Data []{{.Type}}
}

func (*{{.Name}}) RVersion() int16 {
	return rvers.{{.Name}}
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

func (arr *{{.Name}}) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	{{.WFunc}}(arr.Data)

	return int(w.Pos()-pos), w.Err()
}

func (arr *{{.Name}}) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
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
	rtypes.Factory.Add("T{{.Name}}", f)
}

var (
	_ root.Array         = (*{{.Name}})(nil)
	_ rbytes.Marshaler   = (*{{.Name}})(nil)
	_ rbytes.Unmarshaler = (*{{.Name}})(nil)
)
`

func genTClonesArrayData() {
	const (
		macro = "gen-tclonesarray.C"
	)

	err := ioutil.WriteFile(macro, []byte(`
void gen_tclonesarray(const char *fname, bool bypass) {
	auto f = TFile::Open(fname, "RECREATE");
	auto c = new TClonesArray("TObjString", 3);
	(*c)[0] = new TObjString("Elem-0");
	(*c)[1] = new TObjString("elem-1");
	(*c)[2] = new TObjString("Elem-20");

	c->BypassStreamer(bypass);
	f->WriteObjectAny(c, "TClonesArray", "clones");
	f->Write();
	f->Close();
}
`), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(macro)

	for _, v := range []struct {
		name   string
		bypass bool
	}{
		{
			name:   "testdata/tclonesarray-with-streamerbypass.root",
			bypass: true,
		},
		{
			name:   "testdata/tclonesarray-no-streamerbypass.root",
			bypass: false,
		},
	} {
		cmd := exec.Command("root.exe", "-b", "-q", fmt.Sprintf("./%s(%q, %v)", macro, v.name, v.bypass))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
