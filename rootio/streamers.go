// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
	"sync"

	rstreamerspkg "go-hep.org/x/hep/rootio/internal/rstreamers"
	"golang.org/x/xerrors"
)

var (
	cxxNameSanitizer = strings.NewReplacer(
		"<", "_",
		">", "_",
		":", "_",
		",", "_",
		" ", "_",
	)

	reStdVector = regexp.MustCompile("^vector<(.+)>$")
)

type StreamerInfoContext interface {
	StreamerInfo(name string) (StreamerInfo, error)
}

type streamerInfoStore interface {
	addStreamer(si StreamerInfo)
}

type tstreamerInfo struct {
	rvers  int16
	named  tnamed
	chksum uint32
	clsver int32
	objarr *tobjarray
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

func (tsi *tstreamerInfo) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tsi.rvers)
	tsi.named.MarshalROOT(w)
	w.WriteU32(tsi.chksum)
	w.WriteI32(tsi.clsver)

	if len(tsi.elems) > 0 {
		tsi.objarr.arr = make([]Object, len(tsi.elems))
		for i, v := range tsi.elems {
			tsi.objarr.arr[i] = v
		}
	}
	w.WriteObjectAny(tsi.objarr)
	tsi.objarr.arr = nil

	return w.SetByteCount(pos, "TStreamerInfo")
}

func (tsi *tstreamerInfo) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	tsi.rvers = vers

	if err := tsi.named.UnmarshalROOT(r); err != nil {
		return err
	}

	tsi.chksum = r.ReadU32()
	tsi.clsver = r.ReadI32()
	objs := r.ReadObjectAny()
	if r.err != nil {
		return r.err
	}

	tsi.objarr = objs.(*tobjarray)
	tsi.elems = nil
	if tsi.objarr.Len() > 0 {
		tsi.elems = make([]StreamerElement, tsi.objarr.Len())
		for i := range tsi.elems {
			elem := tsi.objarr.At(i)
			tsi.elems[i] = elem.(StreamerElement)
		}
	}
	tsi.objarr.arr = nil

	r.CheckByteCount(pos, bcnt, start, "TStreamerInfo")
	return r.Err()
}

type tstreamerElement struct {
	rvers  int16
	named  tnamed
	etype  int32    // element type
	esize  int32    // size of element
	arrlen int32    // cumulative size of all array dims
	arrdim int32    // number of array dimensions
	maxidx [5]int32 // maximum array index for array dimension "dim"
	offset int32    // element offset in class
	ename  string   // data type name of data member
	xmin   float64  // minimum of data member if a range is specified [xmin.xmax.nbits]
	xmax   float64  // maximum of data member if a range is specified [xmin.xmax.nbits]
	factor float64  // conversion factor if a range is specified. factor = (1<<nbits/(xmax-xmin))
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
	return uintptr(tse.offset)
}

func (tse *tstreamerElement) Size() uintptr {
	return uintptr(tse.esize)
}

func (tse *tstreamerElement) TypeName() string {
	return tse.ename
}

func (tse *tstreamerElement) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tse.rvers)
	tse.named.MarshalROOT(w)
	w.WriteI32(tse.etype)
	w.WriteI32(tse.esize)
	w.WriteI32(tse.arrlen)
	w.WriteI32(tse.arrdim)
	if tse.rvers == 1 {
		w.WriteStaticArrayI32(tse.maxidx[:])
	} else {
		w.WriteFastArrayI32(tse.maxidx[:])
	}
	w.WriteString(tse.ename)

	switch {
	case tse.rvers == 3:
		w.WriteF64(tse.xmin)
		w.WriteF64(tse.xmax)
		w.WriteF64(tse.factor)
	case tse.rvers > 3:
		// FIXME(sbinet)
		// if (TestBit(kHasRange)) GetRange(GetTitle(),fXmin,fXmax,fFactor)
	}

	return w.SetByteCount(pos, "TStreamerElement")
}

func (tse *tstreamerElement) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	tse.rvers = vers
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

	if tse.rvers <= 2 {
		// FIXME(sbinet)
		// tse.esize = tse.arrlen * gROOT->GetType(GetTypeName())->Size()
	}
	switch {
	default:
		tse.xmin = 0
		tse.xmax = 0
		tse.factor = 0
	case tse.rvers == 3:
		tse.xmin = r.ReadF64()
		tse.xmax = r.ReadF64()
		tse.factor = r.ReadF64()
	case tse.rvers > 3:
		// FIXME(sbinet)
		// if (TestBit(kHasRange)) GetRange(GetTitle(),fXmin,fXmax,fFactor)
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerElement")
	return r.Err()
}

type tstreamerBase struct {
	tstreamerElement
	rvers int16
	vbase int32 // version number of the base class
}

func (tsb *tstreamerBase) Class() string {
	return "TStreamerBase"
}

func (tsb *tstreamerBase) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tsb.rvers)
	tsb.tstreamerElement.MarshalROOT(w)

	if tsb.rvers > 2 {
		w.WriteI32(tsb.vbase)
	}

	return w.SetByteCount(pos, "TStreamerBase")
}

func (tsb *tstreamerBase) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	tsb.rvers = vers

	if err := tsb.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	if vers > 2 {
		tsb.vbase = r.ReadI32()
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerBase")
	return r.Err()
}

type tstreamerBasicType struct {
	tstreamerElement
	rvers int16
}

func (tsb *tstreamerBasicType) Class() string {
	return "TStreamerBasicType"
}

func (tsb *tstreamerBasicType) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tsb.rvers)
	tsb.tstreamerElement.MarshalROOT(w)

	return w.SetByteCount(pos, "TStreamerBasicType")
}

func (tsb *tstreamerBasicType) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	tsb.rvers = vers

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
	rvers int16
	cvers int32  // version number of the class with the counter
	cname string // name of data member holding the array count
	ccls  string // name of the class with the counter
}

func (tsb *tstreamerBasicPointer) Class() string {
	return "TStreamerBasicPointer"
}

func (tsb *tstreamerBasicPointer) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tsb.rvers)
	tsb.tstreamerElement.MarshalROOT(w)
	w.WriteI32(tsb.cvers)
	w.WriteString(tsb.cname)
	w.WriteString(tsb.ccls)

	return w.SetByteCount(pos, "TStreamerBasicPointer")
}

func (tsb *tstreamerBasicPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tsb.rvers = vers

	if err := tsb.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	tsb.cvers = r.ReadI32()
	tsb.cname = r.ReadString()
	tsb.ccls = r.ReadString()

	r.CheckByteCount(pos, bcnt, beg, "TStreamerBasicPointer")
	return r.Err()
}

type tstreamerLoop struct {
	tstreamerElement
	rvers  int16
	cvers  int32  // version number of the class with the counter
	cname  string // name of data member holding the array count
	cclass string // name of the class with the counter
}

func (*tstreamerLoop) Class() string {
	return "TStreamerLoop"
}

func (tsl *tstreamerLoop) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tsl.rvers)
	tsl.tstreamerElement.MarshalROOT(w)
	w.WriteI32(tsl.cvers)
	w.WriteString(tsl.cname)
	w.WriteString(tsl.cclass)

	return w.SetByteCount(pos, "TStreamerLoop")
}

func (tsl *tstreamerLoop) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tsl.rvers = vers

	if err := tsl.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	tsl.cvers = r.ReadI32()
	tsl.cname = r.ReadString()
	tsl.cclass = r.ReadString()

	r.CheckByteCount(pos, bcnt, beg, "TStreamerLoop")
	return r.Err()
}

type tstreamerObject struct {
	tstreamerElement
	rvers int16
}

func (tso *tstreamerObject) Class() string {
	return "TStreamerObject"
}

func (tso *tstreamerObject) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tso.rvers)
	tso.tstreamerElement.MarshalROOT(w)
	return w.SetByteCount(pos, "TStreamerObject")
}

func (tso *tstreamerObject) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tso.rvers = vers

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObject")
	return r.Err()
}

type tstreamerObjectPointer struct {
	tstreamerElement
	rvers int16
}

func (tso *tstreamerObjectPointer) Class() string {
	return "TStreamerObjectPointer"
}

func (tso *tstreamerObjectPointer) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tso.rvers)
	tso.tstreamerElement.MarshalROOT(w)
	return w.SetByteCount(pos, "TStreamerObjectPointer")
}

func (tso *tstreamerObjectPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tso.rvers = vers

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectPointer")
	return r.Err()
}

type tstreamerObjectAny struct {
	tstreamerElement
	rvers int16
}

func (tso *tstreamerObjectAny) Class() string {
	return "TStreamerObjectAny"
}

func (tso *tstreamerObjectAny) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tso.rvers)
	tso.tstreamerElement.MarshalROOT(w)

	return w.SetByteCount(pos, "TStreamerObjectAny")
}

func (tso *tstreamerObjectAny) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tso.rvers = vers

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectAny")
	return r.Err()
}

type tstreamerObjectAnyPointer struct {
	tstreamerElement
	rvers int16
}

func (tso *tstreamerObjectAnyPointer) Class() string {
	return "TStreamerObjectAnyPointer"
}

func (tso *tstreamerObjectAnyPointer) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tso.rvers)
	tso.tstreamerElement.MarshalROOT(w)

	return w.SetByteCount(pos, "TStreamerObjectAnyPointer")
}

func (tso *tstreamerObjectAnyPointer) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tso.rvers = vers

	if err := tso.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerObjectAnyPointer")
	return r.Err()
}

type tstreamerString struct {
	tstreamerElement
	rvers int16
}

func (tss *tstreamerString) Class() string {
	return "TStreamerString"
}

func (tss *tstreamerString) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tss.rvers)
	tss.tstreamerElement.MarshalROOT(w)

	return w.SetByteCount(pos, "TStreamerString")
}

func (tss *tstreamerString) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tss.rvers = vers

	if err := tss.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerString")
	return r.Err()
}

type tstreamerSTL struct {
	tstreamerElement
	rvers int16
	vtype int32 // type of STL vector
	ctype int32 // STL contained type
}

func (tss *tstreamerSTL) Class() string {
	return "TStreamerSTL"
}

func (tss *tstreamerSTL) elemTypeName() string {
	o := reStdVector.FindStringSubmatch(tss.ename)
	if o == nil {
		return ""
	}
	return strings.TrimSpace(o[1])
}

func (tss *tstreamerSTL) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tss.rvers)
	tss.tstreamerElement.MarshalROOT(w)
	w.WriteI32(tss.vtype)
	w.WriteI32(tss.ctype)

	return w.SetByteCount(pos, "TStreamerSTL")
}

func (tss *tstreamerSTL) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tss.rvers = vers

	if err := tss.tstreamerElement.UnmarshalROOT(r); err != nil {
		return err
	}

	tss.vtype = r.ReadI32()
	tss.ctype = r.ReadI32()

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
	rvers int16
}

func (tss *tstreamerSTLstring) Class() string {
	return "TStreamerSTLstring"
}

func (tss *tstreamerSTLstring) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tss.rvers)
	tss.tstreamerSTL.MarshalROOT(w)

	return w.SetByteCount(pos, "TStreamerSTLstring")
}

func (tss *tstreamerSTLstring) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tss.rvers = vers

	if err := tss.tstreamerSTL.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, beg, "TStreamerSTLstring")
	return r.Err()
}

type tstreamerArtificial struct {
	tstreamerElement
	rvers int16
}

func (tss *tstreamerArtificial) Class() string {
	return "TStreamerArtificial"
}

func (tsa *tstreamerArtificial) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(tsa.rvers)
	tsa.tstreamerElement.MarshalROOT(w)

	return w.SetByteCount(pos, "TStreamerArtificial")
}

func (tsa *tstreamerArtificial) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tsa.rvers = vers

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
			panic(xerrors.Errorf("rootio: StreamerInfo class=%q version=%d with checksum=%d (got checksum=%d)",
				streamer.Name(), streamer.ClassVersion(), streamer.CheckSum(), old.CheckSum(),
			))
		}
		return
	}

	db.db[key] = streamer
}

func streamerInfoFrom(obj Object, sictx streamerInfoStore) (StreamerInfo, error) {
	r := &memFile{bytes.NewReader(rstreamerspkg.Data)}
	f, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	si, err := f.StreamerInfo(obj.Class())
	if err != nil {
		return nil, err
	}
	streamers.add(si)
	sictx.addStreamer(si)
	return si, nil
}

var defaultStreamerInfos []StreamerInfo

func init() {
	r := &memFile{bytes.NewReader(rstreamerspkg.Data)}
	f, err := NewReader(r)
	if err != nil {
		return
	}
	defer f.Close()

	defaultStreamerInfos = f.StreamerInfos()
}

var (
	_ Object          = (*tstreamerInfo)(nil)
	_ Named           = (*tstreamerInfo)(nil)
	_ StreamerInfo    = (*tstreamerInfo)(nil)
	_ ROOTMarshaler   = (*tstreamerInfo)(nil)
	_ ROOTUnmarshaler = (*tstreamerInfo)(nil)

	_ Object          = (*tstreamerElement)(nil)
	_ Named           = (*tstreamerElement)(nil)
	_ StreamerElement = (*tstreamerElement)(nil)
	_ ROOTMarshaler   = (*tstreamerElement)(nil)
	_ ROOTUnmarshaler = (*tstreamerElement)(nil)

	_ Object          = (*tstreamerBase)(nil)
	_ Named           = (*tstreamerBase)(nil)
	_ StreamerElement = (*tstreamerBase)(nil)
	_ ROOTMarshaler   = (*tstreamerBase)(nil)
	_ ROOTUnmarshaler = (*tstreamerBase)(nil)

	_ Object          = (*tstreamerBasicType)(nil)
	_ Named           = (*tstreamerBasicType)(nil)
	_ StreamerElement = (*tstreamerBasicType)(nil)
	_ ROOTMarshaler   = (*tstreamerBasicType)(nil)
	_ ROOTUnmarshaler = (*tstreamerBasicType)(nil)

	_ Object          = (*tstreamerBasicPointer)(nil)
	_ Named           = (*tstreamerBasicPointer)(nil)
	_ StreamerElement = (*tstreamerBasicPointer)(nil)
	_ ROOTMarshaler   = (*tstreamerBasicPointer)(nil)
	_ ROOTUnmarshaler = (*tstreamerBasicPointer)(nil)

	_ Object          = (*tstreamerLoop)(nil)
	_ Named           = (*tstreamerLoop)(nil)
	_ StreamerElement = (*tstreamerLoop)(nil)
	_ ROOTMarshaler   = (*tstreamerLoop)(nil)
	_ ROOTUnmarshaler = (*tstreamerLoop)(nil)

	_ Object          = (*tstreamerObject)(nil)
	_ Named           = (*tstreamerObject)(nil)
	_ StreamerElement = (*tstreamerObject)(nil)
	_ ROOTMarshaler   = (*tstreamerObject)(nil)
	_ ROOTUnmarshaler = (*tstreamerObject)(nil)

	_ Object          = (*tstreamerObjectPointer)(nil)
	_ Named           = (*tstreamerObjectPointer)(nil)
	_ StreamerElement = (*tstreamerObjectPointer)(nil)
	_ ROOTMarshaler   = (*tstreamerObjectPointer)(nil)
	_ ROOTUnmarshaler = (*tstreamerObjectPointer)(nil)

	_ Object          = (*tstreamerObjectAny)(nil)
	_ Named           = (*tstreamerObjectAny)(nil)
	_ StreamerElement = (*tstreamerObjectAny)(nil)
	_ ROOTMarshaler   = (*tstreamerObjectAny)(nil)
	_ ROOTUnmarshaler = (*tstreamerObjectAny)(nil)

	_ Object          = (*tstreamerString)(nil)
	_ Named           = (*tstreamerString)(nil)
	_ StreamerElement = (*tstreamerString)(nil)
	_ ROOTMarshaler   = (*tstreamerString)(nil)
	_ ROOTUnmarshaler = (*tstreamerString)(nil)

	_ Object          = (*tstreamerSTL)(nil)
	_ Named           = (*tstreamerSTL)(nil)
	_ StreamerElement = (*tstreamerSTL)(nil)
	_ ROOTMarshaler   = (*tstreamerSTL)(nil)
	_ ROOTUnmarshaler = (*tstreamerSTL)(nil)

	_ Object          = (*tstreamerSTLstring)(nil)
	_ Named           = (*tstreamerSTLstring)(nil)
	_ StreamerElement = (*tstreamerSTLstring)(nil)
	_ ROOTMarshaler   = (*tstreamerSTLstring)(nil)
	_ ROOTUnmarshaler = (*tstreamerSTLstring)(nil)

	_ Object          = (*tstreamerArtificial)(nil)
	_ Named           = (*tstreamerArtificial)(nil)
	_ StreamerElement = (*tstreamerArtificial)(nil)
	_ ROOTMarshaler   = (*tstreamerArtificial)(nil)
	_ ROOTUnmarshaler = (*tstreamerArtificial)(nil)
)
