// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type tstreamerInfo struct {
	named  tnamed
	chksum uint32
	clsver int32
	elems  []StreamerElement
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

func (tsi *tstreamerInfo) Elements() []StreamerElement {
	return tsi.elems
}

func (tsi *tstreamerInfo) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := tsi.named.UnmarshalROOT(r); err != nil {
		return err
	}

	r.ReadU32(&tsi.chksum)
	r.ReadI32(&tsi.clsver)
	objs := r.ReadObjectAny()
	if r.err != nil {
		return r.err
	}

	elems := objs.(ObjArray)
	tsi.elems = nil
	if elems.Len() > 0 {
		tsi.elems = make([]StreamerElement, elems.Len())
		for i := range tsi.elems {
			elem := elems.At(i)
			tsi.elems[i] = elem.(StreamerElement)
		}
	}

	r.CheckByteCount(pos, bcnt, start, "TStreamerInfo")
	return r.Err()
}

type tstreamerElement struct {
	named  tnamed
	etype  int32    // element type
	esize  int32    // size of element
	arrlen int32    // cumulative size of all array dims
	arrdim int32    // number of array dimensions
	maxidx [5]int32 // maximum array index for array dimension "dim"
	ename  string   // data type name of data member
	xmin   float64  // minimum of data member if a range is specified [xmin.xmax,nbits]
	xmax   float64  // maximum of data member if a range is specified [xmin,xmax,nbits]
	fact   float64  // conversion factor if a range is specified (fact = (1<<nbits/(xmax-xmin)))
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

func (tse *tstreamerElement) ArrayDim() int {
	return int(tse.arrdim)
}

func (tse *tstreamerElement) ArrayLen() int {
	return int(tse.arrlen)
}

func (tse *tstreamerElement) Type() int {
	return int(tse.etype)
}

func (tse *tstreamerElement) Offset() uintptr {
	return 0
}

func (tse *tstreamerElement) Size() uintptr {
	return uintptr(tse.esize)
}

func (tse *tstreamerElement) TypeName() string {
	return tse.ename
}

func (tse *tstreamerElement) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if err := tse.named.UnmarshalROOT(r); err != nil {
		return err
	}

	r.ReadI32(&tse.etype)
	r.ReadI32(&tse.esize)
	r.ReadI32(&tse.arrlen)
	r.ReadI32(&tse.arrdim)
	if vers == 1 {
		copy(tse.maxidx[:], r.ReadStaticArrayI32())
	} else {
		r.ReadFastArrayI32(tse.maxidx[:])
	}
	r.ReadString(&tse.ename)

	if tse.etype == 11 && (tse.ename == "Bool_t" || tse.ename == "bool") {
		tse.etype = 18
	}

	if vers == 3 {
		r.ReadF64(&tse.xmin)
		r.ReadF64(&tse.xmax)
		r.ReadF64(&tse.fact)
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerElement")
	return r.Err()
}

type tstreamerBase struct {
	tstreamerElement
	vbase int32 // version number of the base class
}

func (tsb *tstreamerBase) Class() string {
	return "TStreamerBase"
}

func (tsb *tstreamerBase) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()

	if err := tsb.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	if vers > 2 {
		r.ReadI32(&tsb.vbase)
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerBase")
	return r.Err()
}

type tstreamerBasicType struct {
	tstreamerElement
}

func (tsb *tstreamerBasicType) Class() string {
	return "TStreamerBasicType"
}

func (tsb *tstreamerBasicType) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tsb.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	etype := tsb.tstreamerElement.etype
	if kOffsetL < etype && etype < kOffsetP {
		etype -= kOffsetL
	}

	basic := true
	switch etype {
	case kBool, kUChar, kChar:
		tsb.tstreamerElement.esize = 1
	case kUShort, kShort:
		tsb.tstreamerElement.esize = 2
	case kBits, kUInt, kInt, kCounter:
		tsb.tstreamerElement.esize = 4
	case kULong, kULong64, kLong, kLong64:
		tsb.tstreamerElement.esize = 8
	case kFloat, kFloat16:
		tsb.tstreamerElement.esize = 4
	case kDouble, kDouble32:
		tsb.tstreamerElement.esize = 8
	case kCharStar:
		tsb.tstreamerElement.esize = int32(ptrSize)
	default:
		basic = false
	}
	if basic && tsb.tstreamerElement.arrlen > 0 {
		tsb.tstreamerElement.esize *= tsb.tstreamerElement.arrlen
	}
	r.CheckByteCount(pos, bcnt, beg, "TStreamerBasicType")
	return r.Err()
}

type tstreamerBasicPointer struct {
	tstreamerElement
	cvers int32  // version number of the class with the counter
	cname string // name of data member holding the array count
	ccls  string // name of the class with the counter
}

func (tsb *tstreamerBasicPointer) Class() string {
	return "TStreamerBasicPointer"
}

func (tsb *tstreamerBasicPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tsb.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.ReadI32(&tsb.cvers)
	r.ReadString(&tsb.cname)
	r.ReadString(&tsb.ccls)

	r.CheckByteCount(pos, bcnt, beg, "TStreamerBasicPointer")
	return r.Err()
}

type tstreamerLoop struct {
	tstreamerElement
	cvers  int32  // version number of the class with the counter
	cname  string // name of data member holding the array count
	cclass string // name of the class with the counter
}

func (*tstreamerLoop) Class() string {
	return "TStreamerLoop"
}

func (tsl *tstreamerLoop) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tsl.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.ReadI32(&tsl.cvers)
	r.ReadString(&tsl.cname)
	r.ReadString(&tsl.cclass)

	r.CheckByteCount(pos, bcnt, beg, "TStreamerLoop")
	return r.Err()
}

type tstreamerObject struct {
	tstreamerElement
}

func (tso *tstreamerObject) Class() string {
	return "TStreamerObject"
}

func (tso *tstreamerObject) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObject")
	return r.Err()
}

type tstreamerObjectPointer struct {
	tstreamerElement
}

func (tso *tstreamerObjectPointer) Class() string {
	return "TStreamerObjectPointer"
}

func (tso *tstreamerObjectPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectPointer")
	return r.Err()
}

type tstreamerObjectAny struct {
	tstreamerElement
}

func (tso *tstreamerObjectAny) Class() string {
	return "TStreamerObjectAny"
}

func (tso *tstreamerObjectAny) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectAny")
	return r.Err()
}

type tstreamerObjectAnyPointer struct {
	tstreamerElement
}

func (tso *tstreamerObjectAnyPointer) Class() string {
	return "TStreamerObjectAnyPointer"
}

func (tso *tstreamerObjectAnyPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectAnyPointer")
	return r.Err()
}

type tstreamerString struct {
	tstreamerElement
}

func (tss *tstreamerString) Class() string {
	return "TStreamerString"
}

func (tss *tstreamerString) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tss.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerString")
	return r.Err()
}

type tstreamerSTL struct {
	tstreamerElement
	vtype int32 // type of STL vector
	ctype int32 // STL contained type
}

func (tss *tstreamerSTL) Class() string {
	return "TStreamerSTL"
}

func (tss *tstreamerSTL) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tss.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.ReadI32(&tss.vtype)
	r.ReadI32(&tss.ctype)

	if tss.vtype == kSTLmultimap || tss.vtype == kSTLset {
		switch {
		case strings.HasPrefix(tss.tstreamerElement.ename, "std::set") || strings.HasPrefix(tss.tstreamerElement.ename, "set"):
			tss.vtype = kSTLset
		case strings.HasPrefix(tss.tstreamerElement.ename, "std::multimap") || strings.HasPrefix(tss.tstreamerElement.ename, "multimap"):
			tss.vtype = kSTLmultimap
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerSTL")
	return r.Err()
}

func (tss *tstreamerSTL) isaPointer() bool {
	tname := tss.tstreamerElement.ename
	return strings.HasSuffix(tname, "*")
}

type tstreamerSTLstring struct {
	tstreamerSTL
}

func (tss *tstreamerSTLstring) Class() string {
	return "TStreamerSTLstring"
}

func (tss *tstreamerSTLstring) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tss.tstreamerSTL.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerSTLstring")
	return r.Err()
}

type tstreamerArtificial struct {
	tstreamerElement
}

func (tss *tstreamerArtificial) Class() string {
	return "TStreamerArtificial"
}

func (tsa *tstreamerArtificial) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := tsa.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerArtificial")
	return r.Err()
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
			o := &tstreamerLoop{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerLoop", f)
		Factory.add("*rootio.tstreamerLoop", f)
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
			o := &tstreamerObjectAnyPointer{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerObjectAnyPointer", f)
		Factory.add("*rootio.tstreamerObjectAnyPointer", f)
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
	{
		f := func() reflect.Value {
			o := &tstreamerSTLstring{
				tstreamerSTL: tstreamerSTL{
					vtype: kSTLstring,
					ctype: kSTLstring,
				},
			}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerSTLstring", f)
		Factory.add("*rootio.tstreamerSTLstring", f)
	}
	{
		f := func() reflect.Value {
			o := &tstreamerArtificial{}
			return reflect.ValueOf(o)
		}
		Factory.add("TStreamerArtificial", f)
		Factory.add("*rootio.tstreamerArtificial", f)
	}
}

var streamers = streamerDb{
	db: make(map[streamerDbKey]StreamerInfo),
}

type streamerDbKey struct {
	class    string
	version  int
	checksum int
}

type streamerDb struct {
	sync.RWMutex
	db map[streamerDbKey]StreamerInfo
}

func (db *streamerDb) getAny(class string) (StreamerInfo, bool) {
	db.RLock()
	defer db.RUnlock()
	for k, v := range db.db {
		if k.class == class {
			return v, true
		}
	}
	return nil, false
}

func (db *streamerDb) get(class string, vers int, chksum int) (StreamerInfo, bool) {
	db.RLock()
	defer db.RUnlock()
	key := streamerDbKey{
		class:    class,
		version:  vers,
		checksum: chksum,
	}

	streamer, ok := db.db[key]
	if !ok {
		return nil, false
	}
	return streamer, true
}

func (db *streamerDb) add(streamer StreamerInfo) {
	db.Lock()
	defer db.Unlock()

	key := streamerDbKey{
		class:    streamer.Name(),
		version:  streamer.ClassVersion(),
		checksum: streamer.CheckSum(),
	}

	old, dup := db.db[key]
	if dup {
		if old.CheckSum() != streamer.CheckSum() {
			panic(fmt.Errorf("rootio: StreamerInfo class=%q version=%d with checksum=%d (got checksum=%d)",
				streamer.Name(), streamer.ClassVersion(), streamer.CheckSum(), old.CheckSum(),
			))
		}
		return
	}

	db.db[key] = streamer
}

var _ Object = (*tstreamerInfo)(nil)
var _ Named = (*tstreamerInfo)(nil)
var _ StreamerInfo = (*tstreamerInfo)(nil)
var _ ROOTUnmarshaler = (*tstreamerInfo)(nil)

var _ Object = (*tstreamerElement)(nil)
var _ Named = (*tstreamerElement)(nil)
var _ StreamerElement = (*tstreamerElement)(nil)
var _ ROOTUnmarshaler = (*tstreamerElement)(nil)

var _ Object = (*tstreamerBase)(nil)
var _ Named = (*tstreamerBase)(nil)
var _ StreamerElement = (*tstreamerBase)(nil)
var _ ROOTUnmarshaler = (*tstreamerBase)(nil)

var _ Object = (*tstreamerBasicType)(nil)
var _ Named = (*tstreamerBasicType)(nil)
var _ StreamerElement = (*tstreamerBasicType)(nil)
var _ ROOTUnmarshaler = (*tstreamerBasicType)(nil)

var _ Object = (*tstreamerBasicPointer)(nil)
var _ Named = (*tstreamerBasicPointer)(nil)
var _ StreamerElement = (*tstreamerBasicPointer)(nil)
var _ ROOTUnmarshaler = (*tstreamerBasicPointer)(nil)

var _ Object = (*tstreamerLoop)(nil)
var _ Named = (*tstreamerLoop)(nil)
var _ StreamerElement = (*tstreamerLoop)(nil)
var _ ROOTUnmarshaler = (*tstreamerLoop)(nil)

var _ Object = (*tstreamerObject)(nil)
var _ Named = (*tstreamerObject)(nil)
var _ StreamerElement = (*tstreamerObject)(nil)
var _ ROOTUnmarshaler = (*tstreamerObject)(nil)

var _ Object = (*tstreamerObjectPointer)(nil)
var _ Named = (*tstreamerObjectPointer)(nil)
var _ StreamerElement = (*tstreamerObjectPointer)(nil)
var _ ROOTUnmarshaler = (*tstreamerObjectPointer)(nil)

var _ Object = (*tstreamerObjectAny)(nil)
var _ Named = (*tstreamerObjectAny)(nil)
var _ StreamerElement = (*tstreamerObjectAny)(nil)
var _ ROOTUnmarshaler = (*tstreamerObjectAny)(nil)

var _ Object = (*tstreamerString)(nil)
var _ Named = (*tstreamerString)(nil)
var _ StreamerElement = (*tstreamerString)(nil)
var _ ROOTUnmarshaler = (*tstreamerString)(nil)

var _ Object = (*tstreamerSTL)(nil)
var _ Named = (*tstreamerSTL)(nil)
var _ StreamerElement = (*tstreamerSTL)(nil)
var _ ROOTUnmarshaler = (*tstreamerSTL)(nil)

var _ Object = (*tstreamerSTLstring)(nil)
var _ Named = (*tstreamerSTLstring)(nil)
var _ StreamerElement = (*tstreamerSTLstring)(nil)
var _ ROOTUnmarshaler = (*tstreamerSTLstring)(nil)

var _ Object = (*tstreamerArtificial)(nil)
var _ Named = (*tstreamerArtificial)(nil)
var _ StreamerElement = (*tstreamerArtificial)(nil)
var _ ROOTUnmarshaler = (*tstreamerArtificial)(nil)
