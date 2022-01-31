// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// gen.rboot generates rdict streamers from the C++/ROOT ones.
// gen.rboot generates ROOT class versions from the C++/ROOT streamers.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"go-hep.org/x/hep/groot/internal/genroot"
	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/riofs"
)

var (
	classes = []string{
		// rbase
		"TAttAxis", "TAttFill", "TAttLine", "TAttMarker",
		"TNamed",
		"TObject", "TObjString",
		"TProcessID", "TProcessUUID", "TRef", "TUUID",
		"TString",

		// rcont
		"TArray", "TArrayC", "TArrayS", "TArrayI", "TArrayL", "TArrayL64", "TArrayF", "TArrayD",
		"TBits",
		"TCollection",
		"TClonesArray",
		"TList",
		"THashList",
		"THashTable",
		"TMap",
		"TObjArray",
		"TRefArray",
		"TRefTable",
		"TSeqCollection",

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
		"TConfidenceLevel",
		"TF1",
		"TF1AbsComposition", "TF1Convolution", "TF1NormSum", "TF1Parameters",
		"TFormula",
		"TGraph", "TGraphErrors", "TGraphAsymmErrors",
		"TH1", "TH1C", "TH1D", "TH1F", "TH1I", "TH1K", "TH1S",
		"TH2", "TH2C", "TH2D", "TH2F", "TH2I", "TH2Poly", "TH2PolyBin", "TH2S",
		"TLimit", "TLimitDataSource",

		// riofs
		"TDirectory",
		"TDirectoryFile",
		"TFile",
		"TKey",

		// rntup
		// "ROOT::Experimental::RNTuple", // FIXME(sbinet): TODO

		// rphys
		"TFeldmanCousins",
		"TLorentzVector",
		"TVector2", "TVector3",

		// rtree
		"ROOT::TIOFeatures",
		"TBasket",
		"TBranch", "TBranchElement", "TBranchObject", "TBranchRef",
		"TChain",
		"TLeaf", "TLeafElement", "TLeafObject",
		"TLeafO",
		"TLeafB", "TLeafS", "TLeafI", "TLeafL",
		"TLeafF", "TLeafD",
		"TLeafF16", "TLeafD32",
		"TLeafC",
		"TNtuple", "TNtupleD",
		"TTree",
	}
)

func main() {
	genStreamers(classes)
}

func genStreamers(classes []string) {
	const (
		oname = "streamers.root"
	)
	defer os.Remove(oname)

	tmpl := template.Must(template.New("genstreamers").Parse(`
void genstreamers(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");

{{range .}}
	TClass::GetClass("{{.}}")->GetStreamerInfo()->Write("streamer-info-{{.}}");
{{end }}

	f->Write();
	f->Close();

	exit(0);
}
`))

	macro := new(strings.Builder)
	err := tmpl.Execute(macro, classes)
	if err != nil {
		log.Fatal(err)
	}

	out, err := rtests.RunCxxROOT("genstreamers", []byte(macro.String()), oname)
	if err != nil {
		log.Fatalf("could not run gen-streamers:\n%s\nerror: %+v", out, err)
	}

	root, err := riofs.Open(oname)
	if err != nil {
		log.Fatalf("could not open ROOT streamers file: %+v", err)
	}
	defer root.Close()

	var streamers []rbytes.StreamerInfo
	for _, k := range root.Keys() {
		if !strings.HasPrefix(k.Name(), "streamer-info-") {
			continue
		}
		obj, err := k.Object()
		if err != nil {
			log.Fatalf("could not load streamer %q: +v", k.Name(), err)
		}
		si, ok := obj.(rbytes.StreamerInfo)
		if !ok {
			log.Fatalf("object %q is not a streamer info: type=%T", k.Name(), obj)
		}
		streamers = append(streamers, si)
	}
	if got, want := len(streamers), len(classes); got != want {
		log.Fatalf("missing streamers: got=%d, want=%d", got, want)
	}

	f, err := os.Create("rdict/cxx_root_streamers_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(`// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rdict

import (
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
)

func init() {
`))
	if err != nil {
		log.Fatalf("could not write groot streamers header: %+v", err)
	}
	for _, sinfo := range streamers {
		fmt.Fprintf(f, "StreamerInfos.Add(")
		err = writeStreamer(f, sinfo)
		if err != nil {
			log.Fatalf("could not generate groot streamer %v: %+v", sinfo, err)
		}
		f.WriteString(")\n")
	}
	fmt.Fprintf(f, "\n}\n")

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)

	genVersions(streamers)
}

func writeStreamer(w io.Writer, si rbytes.StreamerInfo) error {
	const verbose = true
	err := rdict.GenCxxStreamerInfo(w, si, verbose)
	if err != nil {
		return fmt.Errorf("could not write streamer info %q: %w", si.Name(), err)
	}
	return nil
}

type Type struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

func (t Type) GoName() string {
	var (
		namespace = ""
		name      = t.Name
	)
	if strings.HasPrefix(name, "ROOT::Experimental") {
		namespace = "ROOT_Experimental"
		name = name[len("ROOT::Experimental"):]
	}

	if strings.HasPrefix(name, "ROOT::") {
		namespace = "ROOT_"
		name = name[len("ROOT::"):]
	}

	if strings.HasPrefix(name, "T") {
		name = name[1:]
	}
	return namespace + name
}

func genVersions(streamers []rbytes.StreamerInfo) {
	types, err := genROOTTypes(streamers)
	if err != nil {
		log.Fatalf("could not generate ROOT classes versions: %+v", err)
	}

	fname := "./rvers/versions_gen.go"
	year := genroot.ExtractYear(fname)
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports(year, "rvers", f)

	fmt.Fprintf(f, "\n\n// ROOT version\nconst ROOT = %d\n\n", types[len(types)-1].Version)

	tmpl := template.Must(template.New("types").Parse(tmpl))
	err = tmpl.Execute(f, types[:len(types)-1])
	if err != nil {
		log.Fatalf("error executing template: %+v\n", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close generated Go file: %+v", err)
	}
	genroot.GoFmt(f)
}

func genROOTTypes(streamers []rbytes.StreamerInfo) ([]Type, error) {
	const (
		oname = "rversions.json"
	)
	defer os.Remove(oname)

	tmpl := template.Must(template.New("genrversions").Parse(`
#include <fstream>

void genrversions(const char* fname) {
	auto f = fopen(fname, "w");
	fprintf(f, "[\n");
{{range .}}
fprintf(f, "{\"name\": \"{{.Name}}\", \"version\": %d},\n", {{.ClassVersion}});
{{- end }}
	fprintf(f, "{\"name\": \"TROOT\", \"version\": %d}\n]\n", gROOT->GetVersionInt());

	fclose(f);

	exit(0);
}
`))

	macro := new(strings.Builder)
	err := tmpl.Execute(macro, streamers)
	if err != nil {
		log.Fatal(err)
	}

	out, err := rtests.RunCxxROOT("genrversions", []byte(macro.String()), oname)
	if err != nil {
		log.Fatalf("could not run gen-rversions:\n%s\nerror: %+v", out, err)
	}

	f, err := os.ReadFile(oname)
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
