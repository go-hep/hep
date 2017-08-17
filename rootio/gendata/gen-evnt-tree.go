// +build ignore

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	root  = flag.String("f", "test-small.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	const fname = "gentree.C"

	flag.Parse()

	err := ioutil.WriteFile(fname, []byte(script), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fname)

	cmd := exec.Command("root.exe", "-b", fmt.Sprintf("./gentree.C(%q, %d)", *root, *split))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

const script = `
#include <vector>
#include <string>

const int ARRAYSZ = 10;

struct P3 {
	int32_t Px;
	double  Py;
	int32_t Pz;
};

struct Event {
	TString  Beg;
	int16_t  I16;
	int32_t  I32;
	int64_t  I64;
	uint16_t U16;
	uint32_t U32;
	uint64_t U64;
	float    F32;
	double   F64;
	TString  Str;
	P3       P3;

	int16_t  ArrayI16[ARRAYSZ];
	int32_t  ArrayI32[ARRAYSZ];
	int64_t  ArrayI64[ARRAYSZ];
	uint16_t ArrayU16[ARRAYSZ];
	uint32_t ArrayU32[ARRAYSZ];
	uint64_t ArrayU64[ARRAYSZ];
	float    ArrayF32[ARRAYSZ];
	double   ArrayF64[ARRAYSZ];

	int32_t  N;
	int16_t  *SliceI16;  //[N]
	int32_t  *SliceI32;  //[N]
	int64_t  *SliceI64;  //[N]
	uint16_t *SliceU16;  //[N]
	uint32_t *SliceU32;  //[N]
	uint64_t *SliceU64;  //[N]
	float    *SliceF32;  //[N]
	double   *SliceF64;  //[N]

	std::string StdStr;

	std::vector<int16_t> StlVecI16;
	std::vector<int32_t> StlVecI32;
	std::vector<int64_t> StlVecI64;
	std::vector<uint16_t> StlVecU16;
	std::vector<uint32_t> StlVecU32;
	std::vector<uint64_t> StlVecU64;
	std::vector<float> StlVecF32;
	std::vector<double> StlVecF64;
	std::vector<std::string> StlVecStr;

	TString End;
};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 100;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	Event e;

	t->Branch("evt", &e, bufsize, splitlvl);

	for (int i = 0; i != evtmax; i++) {
		e.Beg = TString::Format("beg-%03d", i);
		e.I16 = i;
		e.I32 = i;
		e.I64 = i;
		e.U16 = i;
		e.U32 = i;
		e.U64 = i;
		e.F32 = float(i);
		e.F64 = double(i);
		e.Str = TString::Format("evt-%03d", i);
		e.P3.Px = i-1;
		e.P3.Py = double(i);
		e.P3.Pz = i-1;

		for (int ii = 0; ii != ARRAYSZ; ii++) {
			e.ArrayI16[ii] = i;
			e.ArrayI32[ii] = i;
			e.ArrayI64[ii] = i;
			e.ArrayU16[ii] = i;
			e.ArrayU32[ii] = i;
			e.ArrayU64[ii] = i;
			e.ArrayF32[ii] = float(i);
			e.ArrayF64[ii] = double(i);
		}

		e.N = int32_t(i) % 10;
		e.SliceI16 = (int16_t*)malloc(sizeof(int16_t)*e.N);
		e.SliceI32 = (int32_t*)malloc(sizeof(int32_t)*e.N);
		e.SliceI64 = (int64_t*)malloc(sizeof(int64_t)*e.N);
		e.SliceU16 = (uint16_t*)malloc(sizeof(uint16_t)*e.N);
		e.SliceU32 = (uint32_t*)malloc(sizeof(uint32_t)*e.N);
		e.SliceU64 = (uint64_t*)malloc(sizeof(uint64_t)*e.N);
		e.SliceF32 = (float*)malloc(sizeof(float)*e.N);
		e.SliceF64 = (double*)malloc(sizeof(double)*e.N);

		for (int ii = 0; ii != e.N; ii++) {
			e.SliceI16[ii] = i;
			e.SliceI32[ii] = i;
			e.SliceI64[ii] = i;
			e.SliceU16[ii] = i;
			e.SliceU32[ii] = i;
			e.SliceU64[ii] = i;
			e.SliceF32[ii] = float(i);
			e.SliceF64[ii] = double(i);
		}

		e.StdStr = std::string(TString::Format("std-%03d", i));
		e.StlVecI16.resize(e.N);
		e.StlVecI32.resize(e.N);
		e.StlVecI64.resize(e.N);
		e.StlVecU16.resize(e.N);
		e.StlVecU32.resize(e.N);
		e.StlVecU64.resize(e.N);
		e.StlVecF32.resize(e.N);
		e.StlVecF64.resize(e.N);
		e.StlVecStr.resize(e.N);
		for (int ii =0; ii != e.N; ii++) {
			e.StlVecI16[ii] = i;
			e.StlVecI32[ii] = i;
			e.StlVecI64[ii] = i;
			e.StlVecU16[ii] = i;
			e.StlVecU32[ii] = i;
			e.StlVecU64[ii] = i;
			e.StlVecF32[ii] = float(i);
			e.StlVecF64[ii] = double(i);
			e.StlVecStr[ii] = std::string(TString::Format("vec-%03d", i));
		}
		e.End = TString::Format("end-%03d", i);
		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
