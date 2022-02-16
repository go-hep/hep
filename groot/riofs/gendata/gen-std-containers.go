// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root  = flag.String("f", "std-containers.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentree", []byte(script), *root, *split)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include "TInterpreter.h"
#include "TFile.h"
#include "TTree.h"
#include "TBranch.h"
#include "TString.h"

#include <iostream>
#include <vector>
#include <set>
#include <map>
#include <unordered_set>
#include <unordered_map>
#include <list>
#include <deque>

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	gInterpreter->GenerateDictionary("list<int>", "list");
	gInterpreter->GenerateDictionary("deque<int>", "deque");
	gInterpreter->GenerateDictionary("vector<vector<int> >", "vector");
	gInterpreter->GenerateDictionary("vector<vector<unsigned int> >", "vector");
	gInterpreter->GenerateDictionary("vector<vector<string> >", "vector;string");
	gInterpreter->GenerateDictionary("vector<vector<TString> >", "vector");
	gInterpreter->GenerateDictionary("vector<set<int> >", "vector;set");
	gInterpreter->GenerateDictionary("vector<set<unsigned int> >", "vector;set");
	gInterpreter->GenerateDictionary("vector<set<string> >", "vector;set;string");
	gInterpreter->GenerateDictionary("vector<set<TString> >", "vector;set");
	gInterpreter->GenerateDictionary("set<int>", "set");
	gInterpreter->GenerateDictionary("set<unsigned int>", "set");
	gInterpreter->GenerateDictionary("set<short>", "set");
	gInterpreter->GenerateDictionary("set<string>", "set;string");
	gInterpreter->GenerateDictionary("set<TString>", "set");
	gInterpreter->GenerateDictionary("unordered_set<string>", "unordered_set;string");
	gInterpreter->GenerateDictionary("pair<int,short>", "pair");
	gInterpreter->GenerateDictionary("map<int,short>", "map");
	gInterpreter->GenerateDictionary("map<unsigned int,unsigned short>", "map");
	gInterpreter->GenerateDictionary("map<int,vector<short> >", "map;vector");
	gInterpreter->GenerateDictionary("map<unsigned int,vector<unsigned short> >", "map;vector");
	gInterpreter->GenerateDictionary("map<int,vector<string> >", "map;vector;string");
	gInterpreter->GenerateDictionary("pair<int,set<short> >", "pair;set");
	gInterpreter->GenerateDictionary("pair<int,set<string> >", "pair;set;string");
	gInterpreter->GenerateDictionary("map<int,set<short> >", "map;set");
	gInterpreter->GenerateDictionary("map<int,set<string> >", "map;set;string");
	gInterpreter->GenerateDictionary("map<string,short>", "map;string");
	gInterpreter->GenerateDictionary("map<string,vector<short> >", "map;vector;string");
	gInterpreter->GenerateDictionary("map<string,vector<string> >", "map;vector;string");
	gInterpreter->GenerateDictionary("map<string,set<short> >", "map;set;string");
	gInterpreter->GenerateDictionary("map<string,set<string> >", "map;set;string");
	gInterpreter->GenerateDictionary("map<int,vector<vector<short> > >", "map;vector");
	gInterpreter->GenerateDictionary("map<int,vector<set<short> > >", "map;vector;set");
	gInterpreter->GenerateDictionary("map<string,string>", "map;string");
	gInterpreter->GenerateDictionary("map<string,TString>", "map;string");
	gInterpreter->GenerateDictionary("map<TString,TString>", "map");
	gInterpreter->GenerateDictionary("map<TString,string>", "map;string");
	gInterpreter->GenerateDictionary("unordered_map<string,string>", "unordered_map;string");
	gInterpreter->GenerateDictionary("unordered_map<string,TString>", "unordered_map;string");


	std::string str;
	TString tstr;
	std::list<int32_t> lst_i32;
	std::deque<int32_t> deq_i32;
	std::vector<int32_t> vec_i32;
	std::vector<uint32_t> vec_u32;
	std::vector<std::string> vec_str;
	std::vector<TString> vec_tstr;
	std::vector<std::vector<int32_t>> vec_vec_i32;
	std::vector<std::vector<uint32_t>> vec_vec_u32;
	std::vector<std::vector<std::string>> vec_vec_str;
	std::vector<std::vector<TString>> vec_vec_tstr;
	std::vector<std::set<int32_t>> vec_set_i32;
	std::vector<std::set<uint32_t>> vec_set_u32;
	std::vector<std::set<std::string>> vec_set_str;
	std::vector<std::set<TString>> vec_set_tstr;
	std::set<int32_t> set_i32;
	std::set<uint32_t> set_u32;
	std::set<std::string> set_str;
	std::set<TString> set_tstr;
	std::unordered_set<std::string> uset_str;
	std::map<int32_t, int16_t> map_i32_i16;
	std::map<uint32_t, uint16_t> map_u32_u16;
	std::map<int32_t, std::vector<int16_t> > map_i32_vec_i16;
	std::map<uint32_t, std::vector<uint16_t> > map_u32_vec_u16;
	std::map<int32_t, std::vector<std::string> > map_i32_vec_str;
	std::map<int32_t, std::set<int16_t> > map_i32_set_i16;
	std::map<int32_t, std::set<std::string> > map_i32_set_str;
	std::map<std::string, int16_t> map_str_i16;
	std::map<std::string, std::vector<int16_t> > map_str_vec_i16;
	std::map<std::string, std::vector<std::string> > map_str_vec_str;
	std::map<std::string, std::set<int16_t> > map_str_set_i16;
	std::map<std::string, std::set<std::string> > map_str_set_str;
	std::map<int32_t, std::vector<std::vector<int16_t> > > map_i32_vec_vec_i16;
	std::map<int32_t, std::vector<std::set<int16_t> > > map_i32_vec_set_i16;
	std::map<std::string, std::string> map_str_str;
	std::map<std::string, TString> map_str_tstr;
	std::map<TString, TString> map_tstr_tstr;
	std::map<TString, std::string> map_tstr_str;
	std::unordered_map<std::string, std::string> umap_str_str;
	std::unordered_map<std::string, TString> umap_str_tstr;

	t->Branch("str", &str);
	t->Branch("tstr", &tstr);
	t->Branch("lst_i32", &lst_i32);
	t->Branch("deq_i32", &deq_i32);
	t->Branch("vec_i32", &vec_i32);
	t->Branch("vec_u32", &vec_u32);
	t->Branch("vec_str", &vec_str);
	t->Branch("vec_tstr", &vec_tstr);
	t->Branch("vec_vec_i32", &vec_vec_i32);
	t->Branch("vec_vec_u32", &vec_vec_u32);
	t->Branch("vec_vec_str", &vec_vec_str);
	t->Branch("vec_vec_tstr", &vec_vec_tstr);
	t->Branch("vec_set_i32", &vec_set_i32);
	t->Branch("vec_set_u32", &vec_set_u32);
	t->Branch("vec_set_str", &vec_set_str);
	t->Branch("vec_set_tstr", &vec_set_tstr);
	t->Branch("set_i32", &set_i32);
	t->Branch("set_u32", &set_u32);
	t->Branch("set_str", &set_str);
	t->Branch("set_tstr", &set_tstr);
	t->Branch("uset_str", &uset_str);
	t->Branch("map_i32_i16", &map_i32_i16, bufsize, splitlvl);
	t->Branch("map_u32_u16", &map_u32_u16, bufsize, splitlvl);
	t->Branch("map_i32_vec_i16", &map_i32_vec_i16, bufsize, splitlvl);
	t->Branch("map_u32_vec_u16", &map_u32_vec_u16, bufsize, splitlvl);
	t->Branch("map_i32_vec_str", &map_i32_vec_str, bufsize, splitlvl);
	t->Branch("map_i32_set_i16", &map_i32_set_i16, bufsize, splitlvl);
	t->Branch("map_i32_set_str", &map_i32_set_str, bufsize, splitlvl);
	t->Branch("map_str_i16", &map_str_i16, bufsize, splitlvl);
	t->Branch("map_str_vec_i16", &map_str_vec_i16, bufsize, splitlvl);
	t->Branch("map_str_vec_str", &map_str_vec_str, bufsize, splitlvl);
	t->Branch("map_str_set_i16", &map_str_set_i16, bufsize, splitlvl);
	t->Branch("map_str_set_str", &map_str_set_str, bufsize, splitlvl);
	t->Branch("map_i32_vec_vec_i16", &map_i32_vec_vec_i16, bufsize, splitlvl);
	t->Branch("map_i32_vec_set_i16", &map_i32_vec_set_i16, bufsize, splitlvl);
	t->Branch("map_str_str", &map_str_str, bufsize, splitlvl);
	t->Branch("map_str_tstr", &map_str_tstr, bufsize, splitlvl);
	t->Branch("map_tstr_tstr", &map_tstr_tstr, bufsize, splitlvl);
	t->Branch("map_tstr_str", &map_tstr_str, bufsize, splitlvl);
	t->Branch("umap_str_str", &umap_str_str, bufsize, splitlvl);
	
	str.clear();
	str.assign("one");
	tstr.Clear();
	tstr.Append("one");
	lst_i32.clear();
	lst_i32.push_back(-1);
	deq_i32.clear();
	deq_i32.push_back(-1);
	vec_i32.clear();
	vec_i32.push_back(-1);
	vec_u32.clear();
	vec_u32.push_back(1);
	vec_str.clear();
	vec_str.push_back("one");
	vec_tstr.clear();
	vec_tstr.push_back("one");
	vec_vec_i32.clear();
	vec_vec_i32.push_back(std::vector<int32_t>{ -1 });
	vec_vec_u32.clear();
	vec_vec_u32.push_back(std::vector<uint32_t>{ 1 });
	vec_vec_str.clear();
	vec_vec_str.push_back(std::vector<std::string>{ "one" });
	vec_vec_tstr.clear();
	vec_vec_tstr.push_back(std::vector<TString>{ "one" });
	vec_set_i32.clear();
	vec_set_i32.push_back(std::set<int32_t>{ -1 });
	vec_set_u32.clear();
	vec_set_u32.push_back(std::set<uint32_t>{ 1 });
	vec_set_str.clear();
	vec_set_str.push_back(std::set<std::string>{ "one" });
	vec_set_tstr.clear();
	vec_set_tstr.push_back(std::set<TString>{ "one" });
	set_i32.clear();
	set_i32.insert(-1);
	set_u32.clear();
	set_u32.insert(1);
	set_str.clear();
	set_str.insert("one");
	set_tstr.clear();
	set_tstr.insert("one");
	uset_str.clear();
	uset_str.insert("one");
	map_i32_i16.clear();
	map_i32_i16[-1] = -1;
	map_u32_u16.clear();
	map_u32_u16[1] = 1;
	map_i32_vec_i16.clear();
	map_i32_vec_i16[-1] = std::vector<int16_t>({ -1 });
	map_u32_vec_u16.clear();
	map_u32_vec_u16[1] = std::vector<uint16_t>({ 1 });
	map_i32_vec_str.clear();
	map_i32_vec_str[-1] = std::vector<std::string>({ "one" });
	map_i32_set_i16.clear();
	map_i32_set_i16[-1] = std::set<int16_t>({ -1 });
	map_i32_set_str.clear();
	map_i32_set_str[-1] = std::set<std::string>({ "one" });
	map_str_i16.clear();
	map_str_i16["one"] = -1;
	map_str_vec_i16.clear();
	map_str_vec_i16["one"] = std::vector<int16_t>({ -1 });
	map_str_vec_str.clear();
	map_str_vec_str["one"] = std::vector<std::string>({ "one" });
	map_str_set_i16.clear();
	map_str_set_i16["one"] = std::set<int16_t>({ -1 });
	map_str_set_str.clear();
	map_str_set_str["one"] = std::set<std::string>({ "one" });
	map_i32_vec_vec_i16.clear();
	map_i32_vec_vec_i16[-1] = std::vector<std::vector<int16_t>>{ std::vector<int16_t>{ -1 } };
	map_i32_vec_set_i16.clear();
	map_i32_vec_set_i16[-1] = std::vector<std::set<int16_t>>{ std::set<int16_t>{ -1 } };
	map_str_str.clear();
	map_str_str["one"] = "ONE";
	map_str_tstr.clear();
	map_str_tstr["one"] = "ONE";
	map_tstr_tstr.clear();
	map_tstr_tstr["one"] = "ONE";
	map_tstr_str.clear();
	map_tstr_str["one"] = "ONE";
	umap_str_str.clear();
	umap_str_str["one"] = "ONE";
	
	t->Fill();
	
	str.clear();
	str.assign("two");
	tstr.Clear();
	tstr.Append("two");
	lst_i32.clear();
	lst_i32.push_back(-1);
	lst_i32.push_back(-2);
	deq_i32.clear();
	deq_i32.push_back(-1);
	deq_i32.push_back(-2);
	vec_i32.clear();
	vec_i32.push_back(-1);
	vec_i32.push_back(-2);
	vec_u32.clear();
	vec_u32.push_back(1);
	vec_u32.push_back(2);
	vec_str.clear();
	vec_str.push_back("one");
	vec_str.push_back("two");
	vec_tstr.clear();
	vec_tstr.push_back("one");
	vec_tstr.push_back("two");
	vec_vec_i32.clear();
	vec_vec_i32.push_back(std::vector<int32_t>{ -1 });
	vec_vec_i32.push_back(std::vector<int32_t>{ -1, -2 });
	vec_vec_u32.clear();
	vec_vec_u32.push_back(std::vector<uint32_t>{ 1 });
	vec_vec_u32.push_back(std::vector<uint32_t>{ 1, 2 });
	vec_vec_str.clear();
	vec_vec_str.push_back(std::vector<std::string>{ "one" });
	vec_vec_str.push_back(std::vector<std::string>{ "one", "two" });
	vec_vec_tstr.clear();
	vec_vec_tstr.push_back(std::vector<TString>{ "one" });
	vec_vec_tstr.push_back(std::vector<TString>{ "one", "two" });
	vec_set_i32.clear();
	vec_set_i32.push_back(std::set<int32_t>{ -1 });
	vec_set_i32.push_back(std::set<int32_t>{ -1, -2 });
	vec_set_u32.clear();
	vec_set_u32.push_back(std::set<uint32_t>{ 1 });
	vec_set_u32.push_back(std::set<uint32_t>{ 1, 2 });
	vec_set_str.clear();
	vec_set_str.push_back(std::set<std::string>{ "one" });
	vec_set_str.push_back(std::set<std::string>{ "one", "two" });
	vec_set_tstr.clear();
	vec_set_tstr.push_back(std::set<TString>{ "one" });
	vec_set_tstr.push_back(std::set<TString>{ "one", "two" });
	set_i32.clear();
	set_i32.insert(-1);
	set_i32.insert(-2);
	set_u32.clear();
	set_u32.insert(1);
	set_u32.insert(2);
	set_str.clear();
	set_str.insert("one");
	set_str.insert("two");
	set_tstr.clear();
	set_tstr.insert("one");
	set_tstr.insert("two");
	uset_str.clear();
	uset_str.insert("one");
	uset_str.insert("two");
	map_i32_i16.clear();
	map_i32_i16[-1] = -1;
	map_i32_i16[-2] = -2;
	map_u32_u16.clear();
	map_u32_u16[1] = 1;
	map_u32_u16[2] = 2;
	map_i32_vec_i16.clear();
	map_i32_vec_i16[-1] = std::vector<int16_t>({ -1 });
	map_i32_vec_i16[-2] = std::vector<int16_t>({ -1, -2 });
	map_u32_vec_u16.clear();
	map_u32_vec_u16[1] = std::vector<uint16_t>({ 1 });
	map_u32_vec_u16[2] = std::vector<uint16_t>({ 1, 2 });
	map_i32_vec_str.clear();
	map_i32_vec_str[-1] = std::vector<std::string>({ "one" });
	map_i32_vec_str[-2] = std::vector<std::string>({ "one", "two" });
	map_i32_set_i16.clear();
	map_i32_set_i16[-1] = std::set<int16_t>({ -1 });
	map_i32_set_i16[-2] = std::set<int16_t>({ -1, -2 });
	map_i32_set_str.clear();
	map_i32_set_str[-1] = std::set<std::string>({ "one" });
	map_i32_set_str[-2] = std::set<std::string>({ "one", "two" });
	map_str_i16.clear();
	map_str_i16["one"] = -1;
	map_str_i16["two"] = -2;
	map_str_vec_i16.clear();
	map_str_vec_i16["one"] = std::vector<int16_t>({ -1 });
	map_str_vec_i16["two"] = std::vector<int16_t>({ -1, -2 });
	map_str_vec_str.clear();
	map_str_vec_str["one"] = std::vector<std::string>({ "one" });
	map_str_vec_str["two"] = std::vector<std::string>({ "one", "two" });
	map_str_set_i16.clear();
	map_str_set_i16["one"] = std::set<int16_t>({ -1 });
	map_str_set_i16["two"] = std::set<int16_t>({ -1, -2 });
	map_str_set_str.clear();
	map_str_set_str["one"] = std::set<std::string>({ "one" });
	map_str_set_str["two"] = std::set<std::string>({ "one", "two" });
	map_i32_vec_vec_i16.clear();
	map_i32_vec_vec_i16[-1] = std::vector<std::vector<int16_t>>{ std::vector<int16_t>{ -1 } };
	map_i32_vec_vec_i16[-2] = std::vector<std::vector<int16_t>>{ std::vector<int16_t>{ -1 }, std::vector<int16_t>{-1, -2} };
	map_i32_vec_set_i16.clear();
	map_i32_vec_set_i16[-1] = std::vector<std::set<int16_t>>{ std::set<int16_t>{ -1 } };
	map_i32_vec_set_i16[-2] = std::vector<std::set<int16_t>>{ std::set<int16_t>{ -1 }, std::set<int16_t>{ -1, -2 } };
	map_str_str.clear();
	map_str_str["one"] = "ONE";
	map_str_str["two"] = "TWO";
	map_str_tstr.clear();
	map_str_tstr["one"] = "ONE";
	map_str_tstr["two"] = "TWO";
	map_tstr_tstr.clear();
	map_tstr_tstr["one"] = "ONE";
	map_tstr_tstr["two"] = "TWO";
	map_tstr_str.clear();
	map_tstr_str["one"] = "ONE";
	map_tstr_str["two"] = "TWO";
	umap_str_str.clear();
	umap_str_str["one"] = "ONE";
	umap_str_str["two"] = "TWO";
	
	t->Fill();

	f->Write();
	f->Close();

	exit(0);
}
`
