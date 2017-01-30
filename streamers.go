// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
	"strings"
)

type tstreamerInfo struct {
	named  named
	chksum uint32
	clsver int32
	elems  ObjArray
}

func (tsi *tstreamerInfo) Class() string {
	return "TStreamerInfo"
}

func (tsi *tstreamerInfo) Name() string {
	return tsi.named.Name()
}

func (tsi *tstreamerInfo) Title() string {
	return tsi.named.Title()
}

func (tsi *tstreamerInfo) CheckSum() int {
	return int(tsi.chksum)
}

func (tsi *tstreamerInfo) ClassVersion() int {
	return int(tsi.clsver)
}

func (tsi *tstreamerInfo) Elements() ObjArray {
	return tsi.elems
}

func (tsi *tstreamerInfo) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("tstreamerinfo-vers=%v\n", vers)

	if err := tsi.named.UnmarshalROOT(r); err != nil {
		return err
	}

	tsi.chksum = r.ReadU32()
	tsi.clsver = r.ReadI32()
	objs := r.ReadObjectAny()

	tsi.elems = objs.(ObjArray)

	r.CheckByteCount(pos, bcnt, start, "TStreamerInfo")
	return r.Err()
}

type tstreamerElement struct {
	named  named
	etype  int32    // element type
	esize  int32    // size of element
	arrlen int32    // cumulative size of all array dims
	arrdim int32    // number of array dimensions
	maxidx [5]int32 // maximum array index for array dimension "dim"
	ename  string   // data type name of data member
}

func (tse *tstreamerElement) Class() string {
	return "TStreamerElement"
}

func (tse *tstreamerElement) Name() string {
	return tse.named.Name()
}

func (tse *tstreamerElement) Title() string {
	return tse.named.Title()
}

func (tse *tstreamerElement) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if err := tse.named.UnmarshalROOT(r); err != nil {
		return err
	}

	tse.etype = r.ReadI32()
	tse.esize = r.ReadI32()
	tse.arrlen = r.ReadI32()
	tse.arrdim = r.ReadI32()
	if vers == 1 {
		copy(tse.maxidx[:], r.ReadStaticArrayI32())
	} else {
		copy(tse.maxidx[:], r.ReadFastArrayI32(len(tse.maxidx)))
	}
	tse.ename = r.ReadString()

	if tse.etype == 11 && (tse.ename == "Bool_t" || tse.ename == "bool") {
		tse.etype = 18
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerElement")
	return r.Err()
}

type tstreamerBase struct {
	tse      tstreamerElement
	baseVers int32 // version number of the base class
}

func (tsb *tstreamerBase) Class() string {
	return "TStreamerBase"
}

func (tsb *tstreamerBase) Name() string {
	return tsb.tse.Name()
}

func (tsb *tstreamerBase) Title() string {
	return tsb.tse.Title()
}

func (tsb *tstreamerBase) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()

	if err := tsb.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	if vers > 2 {
		tsb.baseVers = r.ReadI32()
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerBase")
	return r.Err()
}

type tstreamerBasicType struct {
	tse tstreamerElement
}

func (tsb *tstreamerBasicType) Class() string {
	return "TStreamerBasicType"
}

func (tsb *tstreamerBasicType) Name() string {
	return tsb.tse.Name()
}

func (tsb *tstreamerBasicType) Title() string {
	return tsb.tse.Title()
}

func (tsb *tstreamerBasicType) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := tsb.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	etype := tsb.tse.etype
	if kOffsetL < etype && etype < kOffsetP {
		etype -= kOffsetL
	}

	basic := true
	switch etype {
	case kBool, kUChar, kChar:
		tsb.tse.esize = 1
	case kUShort, kShort:
		tsb.tse.esize = 2
	case kBits, kUInt, kInt, kCounter:
		tsb.tse.esize = 4
	case kULong, kULong64, kLong, kLong64:
		tsb.tse.esize = 8
	case kFloat, kFloat16:
		tsb.tse.esize = 4
	case kDouble, kDouble32:
		tsb.tse.esize = 8
	case kCharStar:
		tsb.tse.esize = int32(ptrSize)
	default:
		basic = false
	}
	if basic && tsb.tse.arrlen > 0 {
		tsb.tse.esize *= tsb.tse.arrlen
	}
	r.CheckByteCount(pos, bcnt, beg, "TStreamerBasicType")
	return r.Err()
}

type tstreamerBasicPointer struct {
	tse   tstreamerElement
	cvers int32  // version number of the class with the counter
	cname string // name of data member holding the array count
	ccls  string // name of the class with the counter
}

func (tsb *tstreamerBasicPointer) Class() string {
	return "TStreamerBasicPointer"
}

func (tsb *tstreamerBasicPointer) Name() string {
	return tsb.tse.Name()
}

func (tsb *tstreamerBasicPointer) Title() string {
	return tsb.tse.Title()
}

func (tsb *tstreamerBasicPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := tsb.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	tsb.cvers = r.ReadI32()
	tsb.cname = r.ReadString()
	tsb.ccls = r.ReadString()

	r.CheckByteCount(pos, bcnt, beg, "TStreamerBasicPointer")
	return r.Err()
}

type tstreamerObject struct {
	tse tstreamerElement
}

func (tso *tstreamerObject) Class() string {
	return "TStreamerObject"
}

func (tso *tstreamerObject) Name() string {
	return tso.tse.Name()
}

func (tso *tstreamerObject) Title() string {
	return tso.tse.Title()
}

func (tso *tstreamerObject) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := tso.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObject")
	return r.Err()
}

type tstreamerObjectPointer struct {
	tse tstreamerElement
}

func (tso *tstreamerObjectPointer) Class() string {
	return "TStreamerObjectPointer"
}

func (tso *tstreamerObjectPointer) Name() string {
	return tso.tse.Name()
}

func (tso *tstreamerObjectPointer) Title() string {
	return tso.tse.Title()
}

func (tso *tstreamerObjectPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := tso.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectPointer")
	return r.Err()
}

type tstreamerObjectAny struct {
	tse tstreamerElement
}

func (tso *tstreamerObjectAny) Class() string {
	return "TStreamerObjectAny"
}

func (tso *tstreamerObjectAny) Name() string {
	return tso.tse.Name()
}

func (tso *tstreamerObjectAny) Title() string {
	return tso.tse.Title()
}

func (tso *tstreamerObjectAny) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := tso.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectAny")
	return r.Err()
}

type tstreamerString struct {
	tse tstreamerElement
}

func (tss *tstreamerString) Class() string {
	return "TStreamerString"
}

func (tss *tstreamerString) Name() string {
	return tss.tse.Name()
}

func (tss *tstreamerString) Title() string {
	return tss.tse.Title()
}

func (tss *tstreamerString) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := tss.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerString")
	return r.Err()
}

type tstreamerSTL struct {
	tse   tstreamerElement
	vtype int32 // type of STL vector
	ctype int32 // STL contained type
}

func (tss *tstreamerSTL) Class() string {
	return "TStreamerSTL"
}

func (tss *tstreamerSTL) Name() string {
	return tss.tse.Name()
}

func (tss *tstreamerSTL) Title() string {
	return tss.tse.Title()
}

func (tss *tstreamerSTL) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := tss.tse.UnmarshalROOT(r); err != nil {
		return err
	}

	tss.vtype = r.ReadI32()
	tss.ctype = r.ReadI32()

	if tss.vtype == kSTLmultimap || tss.vtype == kSTLset {
		switch {
		case strings.HasPrefix(tss.tse.ename, "std::set") || strings.HasPrefix(tss.tse.ename, "set"):
			tss.vtype = kSTLset
		case strings.HasPrefix(tss.tse.ename, "std::multimap") || strings.HasPrefix(tss.tse.ename, "multimap"):
			tss.vtype = kSTLmultimap
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerSTL")
	return r.Err()
}

func (tss *tstreamerSTL) isaPointer() bool {
	tname := tss.tse.ename
	return strings.HasSuffix(tname, "*")
}

func init() {
	{
		f := func() reflect.Value {
			o := &tstreamerInfo{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerInfo", f)
		Factory.add("*rootio.tstreamerInfo", f)
	}

	{
		f := func() reflect.Value {
			o := &tstreamerElement{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerElement", f)
		Factory.add("*rootio.tstreamerElement", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerBase{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerBase", f)
		Factory.add("*rootio.tstreamerBase", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerBasicType{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerBasicType", f)
		Factory.add("*rootio.tstreamerBasicType", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerBasicPointer{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerBasicPointer", f)
		Factory.add("*rootio.tstreamerBasicPointer", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerObject{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerObject", f)
		Factory.add("*rootio.tstreamerObject", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerObjectPointer{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerObjectPointer", f)
		Factory.add("*rootio.tstreamerObjectPointer", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerObjectAny{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerObjectAny", f)
		Factory.add("*rootio.tstreamerObjectAny", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerString{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerString", f)
		Factory.add("*rootio.tstreamerString", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerSTL{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerSTL", f)
		Factory.add("*rootio.tstreamerSTL", f)
	}
}

var _ Object = (*tstreamerInfo)(nil)
var _ Named = (*tstreamerInfo)(nil)
var _ StreamerInfo = (*tstreamerInfo)(nil)
var _ ROOTUnmarshaler = (*tstreamerInfo)(nil)

var _ Object = (*tstreamerElement)(nil)
var _ Named = (*tstreamerElement)(nil)
var _ ROOTUnmarshaler = (*tstreamerElement)(nil)

var _ Object = (*tstreamerBase)(nil)
var _ Named = (*tstreamerBase)(nil)
var _ ROOTUnmarshaler = (*tstreamerBase)(nil)

var _ Object = (*tstreamerBasicType)(nil)
var _ Named = (*tstreamerBasicType)(nil)
var _ ROOTUnmarshaler = (*tstreamerBasicType)(nil)

var _ Object = (*tstreamerBasicPointer)(nil)
var _ Named = (*tstreamerBasicPointer)(nil)
var _ ROOTUnmarshaler = (*tstreamerBasicPointer)(nil)

var _ Object = (*tstreamerObject)(nil)
var _ Named = (*tstreamerObject)(nil)
var _ ROOTUnmarshaler = (*tstreamerObject)(nil)

var _ Object = (*tstreamerObjectPointer)(nil)
var _ Named = (*tstreamerObjectPointer)(nil)
var _ ROOTUnmarshaler = (*tstreamerObjectPointer)(nil)

var _ Object = (*tstreamerObjectAny)(nil)
var _ Named = (*tstreamerObjectAny)(nil)
var _ ROOTUnmarshaler = (*tstreamerObjectAny)(nil)

var _ Object = (*tstreamerString)(nil)
var _ Named = (*tstreamerString)(nil)
var _ ROOTUnmarshaler = (*tstreamerString)(nil)

var _ Object = (*tstreamerSTL)(nil)
var _ Named = (*tstreamerSTL)(nil)
var _ ROOTUnmarshaler = (*tstreamerSTL)(nil)
