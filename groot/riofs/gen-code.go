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
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/riofs"
)

func main() {
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

func genStreamers() {
	classes := []struct {
		Name string
		Type string
	}{
		{Name: "TObject"},
		{Name: "TDirectory"},
		{Name: "TDirectoryFile"},
		{Name: "TKey"},
		{Name: "TNamed"},
		{Name: "TClonesArray"},
		{Name: "TArrayC", Type: "TArray"},
		{Name: "TArrayS", Type: "TArray"},
		{Name: "TArrayI", Type: "TArray"},
		{Name: "TArrayL", Type: "TArray"},
		{Name: "TArrayL64", Type: "TArray"},
		{Name: "TArrayF", Type: "TArray"},
		{Name: "TArrayD", Type: "TArray"},
		{Name: "TList"},
		{Name: "THashList"},
		{Name: "TMap"},
		{Name: "TObjArray"},
		{Name: "TObjString"},
		{Name: "TProcessID"},

		{Name: "TGraph"}, {Name: "TGraphErrors"}, {Name: "TGraphAsymmErrors"},
		{Name: "TH1F"}, {Name: "TH1D"}, {Name: "TH1I"},
		{Name: "TH2F"}, {Name: "TH2D"}, {Name: "TH2I"},

		{Name: "TTree"},
		{Name: "TBranch"}, {Name: "TBranchElement"},
		{Name: "TBasket"},
		{Name: "TLeaf"}, {Name: "TLeafElement"},
		{Name: "TLeafO"},
		{Name: "TLeafB"}, {Name: "TLeafS"}, {Name: "TLeafI"}, {Name: "TLeafL"},
		{Name: "TLeafF"}, {Name: "TLeafD"},
		{Name: "TLeafC"},
	}

	const (
		oname = "streamers.root"
	)
	defer os.Remove(oname)

	tmpl := template.Must(template.New("genstreamers").Parse(`
void genstreamers(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");
	f->SetCompressionAlgorithm(1);
	f->SetCompressionLevel(9);

{{range .}}
{{if ne .Type "TArray"}}
	(({{.Name}}*)(TClass::GetClass("{{.Name}}")->New()))->Write("type-{{.}}");
{{- end}}
	TClass::GetClass("{{.Name}}")->GetStreamerInfo()->Write("streamer-info-{{.Name}}");
{{end }}

	f->Write();
	f->Close();

	exit(0);
}
`))

	macro := new(strings.Builder)
	err = tmpl.Execute(macro, classes)
	if err != nil {
		log.Fatal(err)
	}

	out, err := rtests.RunCxxROOT("genstreamers", []byte(macro.String()), oname)
	if err != nil {
		log.Fatalf("could not run gen-streamers:\n%s\nerror: %+v", out, err)
	}

	root, err := riofs.Open(oname)
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

	fmt.Fprintf(f, `// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rstreamers

import (
	"encoding/base64"

	"golang.org/x/xerrors"
)

var Data []byte

func init() {
	var err error
	Data, err = base64.StdEncoding.DecodeString(`,
	)

	fmt.Fprintf(f, "`%s`)\n", base64.StdEncoding.EncodeToString(raw))
	fmt.Fprintf(f, `
	if err != nil {
		panic(xerrors.Errorf("rstreamers: could not decode embedded streamer: %%w", err))
	}
}
`)

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}
