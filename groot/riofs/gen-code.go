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
	"os/exec"
	"text/template"

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
	classes := []string{
		"TObject",
		"TDirectory",
		"TDirectoryFile",
		"TKey",
		"TNamed",
		"TList",
		"THashList",
		"TMap",
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
	f->SetCompressionAlgorithm(1);
	f->SetCompressionLevel(9);

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
		panic(fmt.Errorf("rstreamers: could not decode embedded streamer: %%v", err))
	}
}
`)

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}
