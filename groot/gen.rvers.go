// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"go-hep.org/x/hep/groot/internal/genroot"
)

func main() {
	genVersions()
}

type Type struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

func (t Type) GoName() string {
	if strings.HasPrefix(t.Name, "T") {
		return t.Name[1:]
	}
	return t.Name
}

func genVersions() {
	types, err := genROOTTypes()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./rvers/versions_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rvers", f)

	fmt.Fprintf(f, "\n\n// ROOT version\nconst ROOT = %d\n\n", types[len(types)-1].Version)

	tmpl := template.Must(template.New("types").Parse(tmpl))
	err = tmpl.Execute(f, types[:len(types)-1])
	if err != nil {
		log.Fatalf("error executing template: %v\n", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)
}

func genROOTTypes() ([]Type, error) {
	rclasses := []string{
		// rbase
		"TAttAxis", "TAttFill", "TAttLine", "TAttMarker",
		"TNamed",
		"TObject", "TObjString",
		"TProcessID",
		"TUUID",
		// rcont
		"TArray", "TArrayC", "TArrayD", "TArrayF", "TArrayI", "TArrayL", "TArrayL64", "TArrayS",
		"THashList",
		"TList",
		"TMap",
		"TObjArray",
		"TClonesArray",
		// rdict
		"TStreamerInfo",
		"TStreamerElement",
		"TStreamerBase",
		"TStreamerBasicType",
		"TStreamerBasicPointer",
		"TStreamerLoop",
		"TStreamerObject",
		"TStreamerObjectPointer",
		"TStreamerObjectAny",
		"TStreamerObjectAnyPointer",
		"TStreamerString",
		"TStreamerSTL",
		"TStreamerSTLstring",
		"TStreamerArtificial",

		// rhist
		"TAxis",
		"TGraph", "TGraphErrors", "TGraphAsymmErrors",
		"TH1", "TH1C", "TH1D", "TH1F", "TH1I", "TH1K", "TH1S",
		"TH2", "TH2C", "TH2D", "TH2F", "TH2I", "TH2Poly", "TH2PolyBin", "TH2S",

		// rtree
		"TBasket",
		"TBranch", "TBranchElement",
		"TChain",
		"TLeaf", "TLeafElement",
		"TLeafB", "TLeafC", "TLeafD", "TLeafF", "TLeafI", "TLeafL", "TLeafO", "TLeafS",
		"TNtuple",
		"TTree",
		// riofs
		"TDirectory",
		"TDirectoryFile",
		"TFile",
		"TKey",
	}

	const (
		macro = "genrversions.C"
		oname = "rversions.json"
	)

	froot, err := os.Create(macro)
	if err != nil {
		log.Fatal(err)
	}
	defer froot.Close()
	defer os.Remove(macro)
	defer os.Remove(oname)

	tmpl := template.Must(template.New("genrversions").Parse(`
#include <fstream>

void genrversions(const char* fname) {
	auto f = fopen(fname, "w");
	fprintf(f, "[\n");
{{range .}}
fprintf(f, "{\"name\": \"{{.}}\", \"version\": %d},\n", {{.}}::Class_Version());
{{- end }}
	fprintf(f, "{\"name\": \"TROOT\", \"version\": %d}\n]\n", gROOT->GetVersionInt());

	fclose(f);

	exit(0);
}
`))
	err = tmpl.Execute(froot, rclasses)
	if err != nil {
		log.Fatal(err)
	}

	err = froot.Close()
	if err != nil {
		log.Fatalf("could not close ROOT macro: %v", err)
	}

	cmd := exec.Command("root.exe", "-b", fmt.Sprintf("./%s(%q)", macro, oname))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	f, err := ioutil.ReadFile(oname)
	if err != nil {
		log.Fatal(err)
	}

	var types []Type
	err = json.NewDecoder(bytes.NewReader(f)).Decode(&types)
	return types, err
}

const tmpl = `// ROOT classes versions
const (
{{range .}}
	{{.GoName}} = {{.Version}} // ROOT version for {{.Name}}
{{- end}}
)
`
