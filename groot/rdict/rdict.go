// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rdict contains the definition of ROOT streamers and facilities
// to generate new streamers meta data from user types.
package rdict // import "go-hep.org/x/hep/groot/rdict"

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

var (
	ptrSize = 4 << (^uintptr(0) >> 63)
	intSize = int(reflect.TypeOf(int(0)).Size())
)

type StreamerInfo struct {
	named  rbase.Named
	chksum uint32
	clsver int32
	objarr *rcont.ObjArray
	elems  []rbytes.StreamerElement

	init  sync.Once
	descr []elemDescr
	roops []rstreamer // read-stream object-wise operations
	woops []wstreamer // write-stream object-wise operations
	rmops []rstreamer // read-stream member-wise operations
	wmops []wstreamer // write-stream member-wise operations
}

// NewStreamerInfo creates a new StreamerInfo from Go provided informations.
func NewStreamerInfo(name string, version int, elems []rbytes.StreamerElement) *StreamerInfo {
	sinfos := &StreamerInfo{
		named:  *rbase.NewNamed(GoName2Cxx(name), "Go;"+name),
		chksum: genChecksum(name, elems),
		clsver: int32(version),
		objarr: rcont.NewObjArray(),
		elems:  elems,
	}
	return sinfos
}

// NewCxxStreamerInfo creates a new StreamerInfo from C++ provided informations.
func NewCxxStreamerInfo(name string, version int32, chksum uint32, elems []rbytes.StreamerElement) *StreamerInfo {
	sinfos := &StreamerInfo{
		named:  *rbase.NewNamed(name, ""),
		chksum: chksum,
		clsver: version,
		objarr: rcont.NewObjArray(),
		elems:  elems,
	}
	return sinfos
}

func (*StreamerInfo) RVersion() int16 { return rvers.StreamerInfo }

func (tsi *StreamerInfo) Class() string {
	return "TStreamerInfo"
}

func (tsi *StreamerInfo) Name() string {
	return tsi.named.Name()
}

func (tsi *StreamerInfo) Title() string {
	return tsi.named.Title()
}

func (tsi *StreamerInfo) CheckSum() int {
	return int(tsi.chksum)
}

func (tsi *StreamerInfo) ClassVersion() int {
	return int(tsi.clsver)
}

func (tsi *StreamerInfo) Elements() []rbytes.StreamerElement {
	return tsi.elems
}

func (tsi *StreamerInfo) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tsi.Class(), tsi.RVersion())
	w.WriteObject(&tsi.named)
	w.WriteU32(tsi.chksum)
	w.WriteI32(tsi.clsver)

	if len(tsi.elems) > 0 {
		elems := make([]root.Object, len(tsi.elems))
		for i, v := range tsi.elems {
			elems[i] = v
		}
		tsi.objarr.SetElems(elems)
	}
	w.WriteObjectAny(tsi.objarr)
	tsi.objarr.SetElems(nil)

	return w.SetHeader(hdr)
}

func (tsi *StreamerInfo) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tsi.Class())
	if hdr.Vers > rvers.StreamerInfo {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tsi.Class(), hdr.Vers, tsi.RVersion(),
		))
	}

	r.ReadObject(&tsi.named)

	tsi.chksum = r.ReadU32()
	tsi.clsver = r.ReadI32()
	objs := r.ReadObjectAny()
	if r.Err() != nil {
		return r.Err()
	}

	tsi.objarr = objs.(*rcont.ObjArray)
	tsi.elems = nil
	if tsi.objarr.Len() > 0 {
		tsi.elems = make([]rbytes.StreamerElement, tsi.objarr.Len())
		for i := range tsi.elems {
			elem := tsi.objarr.At(i)
			tsi.elems[i] = elem.(rbytes.StreamerElement)
		}
	}
	tsi.objarr.SetElems(nil)

	r.CheckHeader(hdr)
	return r.Err()
}

func (si *StreamerInfo) String() string {
	o := new(strings.Builder)
	fmt.Fprintf(o, " StreamerInfo for %q version=%d title=%q\n", si.Name(), si.ClassVersion(), si.Title())
	w := tabwriter.NewWriter(o, 8, 4, 1, ' ', 0)
	for _, elm := range si.Elements() {
		fmt.Fprintf(w, "  %s\t%s\toffset=%3d\ttype=%3d\tsize=%3d\t %s\n", elm.TypeName(), elm.Name(), elm.Offset(), elm.Type(), elm.Size(), elm.Title())
	}
	w.Flush()
	return o.String()

}

// BuildStreamers builds the r/w streamers.
func (si *StreamerInfo) BuildStreamers() error {
	var err error
	si.init.Do(func() {
		err = si.build(StreamerInfos)
	})
	return err
}

func (si *StreamerInfo) build(sictx rbytes.StreamerInfoContext) error {
	si.descr = make([]elemDescr, 0, len(si.elems))
	si.roops = make([]rstreamer, 0, len(si.elems))
	si.woops = make([]wstreamer, 0, len(si.elems))
	si.rmops = make([]rstreamer, 0, len(si.elems))
	si.wmops = make([]wstreamer, 0, len(si.elems))

	for i, se := range si.elems {
		class := strings.TrimRight(se.TypeName(), "*")
		method := func() []int {
			var cname string
			switch se := se.(type) {
			case *StreamerBasicPointer:
				cname = se.CountName()
			case *StreamerLoop:
				cname = se.CountName()
				if se.cclass != si.Name() {
					// reaching into another class internals isn't supported (yet?) in groot
					panic(fmt.Errorf("rdict: unsupported StreamerLoop case: si=%q, se=%q, count=%q, class=%q",
						si.Name(), se.Name(), cname, se.cclass,
					))
				}
			default:
				return []int{0}
			}

			return si.findField(sictx, cname, se, nil)
		}()

		if method == nil {
			return fmt.Errorf("rdict: could not find count-offset for element %q in streamer %q (se=%T)", se.Name(), si.Name(), se)
		}

		descr := elemDescr{
			otype:  se.Type(),
			ntype:  se.Type(), // FIXME(sbinet): handle schema evolution
			offset: i,         // FIXME(sbinet): make sure this works (instead of se.Offset())
			length: se.ArrayLen(),
			elem:   se,
			method: method, // FIXME(sbinet): schema evolution (old/new class may not have the same "offsets")
			oclass: class,  // FIXME(sbinet): impl.
			nclass: class,  // FIXME(sbinet): impl. + schema evolution
			mbr:    nil,    // FIXME(sbinet): impl
		}

		si.descr = append(si.descr, descr)
	}

	for i, descr := range si.descr {
		si.roops = append(si.roops, si.makeROp(sictx, i, descr))
		si.woops = append(si.woops, si.makeWOp(sictx, i, descr))
		// FIXME(sbinet): handle member-wise r/w ops
		// si.rmops
		// si.wmops
	}

	return nil
}

func (si *StreamerInfo) NewDecoder(kind rbytes.StreamKind, r *rbytes.RBuffer) (rbytes.Decoder, error) {
	err := si.BuildStreamers()
	if err != nil {
		return nil, fmt.Errorf("rdict: could not build read streamers: %w", err)
	}

	switch kind {
	case rbytes.ObjectWise:
		return newDecoder(r, si, kind, si.roops)
	case rbytes.MemberWise:
		return newDecoder(r, si, kind, si.rmops)
	default:
		return nil, fmt.Errorf("rdict: invalid stream kind %v", kind)
	}
}

func (si *StreamerInfo) NewEncoder(kind rbytes.StreamKind, w *rbytes.WBuffer) (rbytes.Encoder, error) {
	err := si.BuildStreamers()
	if err != nil {
		return nil, fmt.Errorf("rdict: could not build write streamers: %w", err)
	}

	switch kind {
	case rbytes.ObjectWise:
		return newEncoder(w, si, kind, si.woops)
	case rbytes.MemberWise:
		return newEncoder(w, si, kind, si.wmops)
	default:
		return nil, fmt.Errorf("rdict: invalid stream kind %v", kind)
	}
}

func (si *StreamerInfo) NewRStreamer(kind rbytes.StreamKind) (rbytes.RStreamer, error) {
	err := si.BuildStreamers()
	if err != nil {
		return nil, fmt.Errorf("rdict: could not build read streamers: %w", err)
	}

	roops := make([]rstreamer, len(si.descr))
	switch kind {
	case rbytes.ObjectWise:
		copy(roops, si.roops)
	case rbytes.MemberWise:
		copy(roops, si.rmops)
	default:
		return nil, fmt.Errorf("rdict: invalid stream kind %v", kind)
	}
	return newRStreamerInfo(si, kind, roops)
}

func (si *StreamerInfo) NewWStreamer(kind rbytes.StreamKind) (rbytes.WStreamer, error) {
	err := si.BuildStreamers()
	if err != nil {
		return nil, fmt.Errorf("rdict: could not build write streamers: %w", err)
	}

	wops := make([]wstreamer, len(si.descr))
	switch kind {
	case rbytes.ObjectWise:
		copy(wops, si.woops)
	case rbytes.MemberWise:
		copy(wops, si.wmops)
	default:
		return nil, fmt.Errorf("rdict: invalid stream kind %v", kind)
	}
	return newWStreamerInfo(si, kind, wops)
}

func (si *StreamerInfo) findField(ctx rbytes.StreamerInfoContext, name string, se rbytes.StreamerElement, offset []int) []int {
	for j := range si.elems {
		if si.elems[j].Name() == name {
			offset = append(offset, j)
			return offset
		}
	}

	// look into base classes, if any.
	for j, bse := range si.elems {
		switch bse := bse.(type) {
		case *StreamerBase:
			base, err := ctx.StreamerInfo(bse.Name(), -1)
			if err != nil {
				panic(fmt.Errorf("rdict: could not find base class %q of %q: %w", se.Name(), si.Name(), err))
			}
			boffset := base.(*StreamerInfo).findField(ctx, name, se, nil)
			if boffset != nil {
				return append(append(offset, j), boffset...)
			}
		}
	}

	return nil
}

type Element struct {
	Name   rbase.Named
	Type   rmeta.Enum // element type
	Size   int32      // size of element
	ArrLen int32      // cumulative size of all array dims
	ArrDim int32      // number of array dimensions
	MaxIdx [5]int32   // maximum array index for array dimension "dim"
	Offset int32      // element offset in class
	EName  string     // data type name of data member
	XMin   float64    // minimum of data member if a range is specified [xmin.xmax.nbits]
	XMax   float64    // maximum of data member if a range is specified [xmin.xmax.nbits]
	Factor float64    // conversion factor if a range is specified. factor = (1<<nbits/(xmax-xmin))
}

func (e Element) New() StreamerElement {
	e.parse()
	return StreamerElement{
		named:  e.Name,
		etype:  e.Type,
		esize:  e.Size,
		arrlen: e.ArrLen,
		arrdim: e.ArrDim,
		maxidx: e.MaxIdx,
		offset: e.Offset,
		ename:  e.EName,
		xmin:   e.XMin,
		xmax:   e.XMax,
		factor: e.Factor,
	}
}

// parse parses the element's title for ROOT meta-data information (range, factor, ...)
func (e *Element) parse() {
	switch e.Type {
	case rmeta.Float16, rmeta.OffsetL + rmeta.Float16, rmeta.OffsetP + rmeta.Float16:
		e.XMin, e.XMax, e.Factor = e.getRange(e.Name.Title())
	case rmeta.Double32, rmeta.OffsetL + rmeta.Double32, rmeta.OffsetP + rmeta.Double32:
		e.XMin, e.XMax, e.Factor = e.getRange(e.Name.Title())
	}
}

func (Element) getRange(str string) (xmin, xmax, factor float64) {
	if str == "" {
		return xmin, xmax, factor
	}
	beg := strings.LastIndex(str, "[")
	if beg < 0 {
		return xmin, xmax, factor
	}
	if beg > 0 {
		// make sure there is a '/' in-between
		// make sure the slash is just one position before.
		slash := strings.LastIndex(str[:beg], "/")
		if slash < 0 || slash+2 != beg {
			return xmin, xmax, factor
		}
	}
	end := strings.LastIndex(str, "]")
	if end < 0 {
		return xmin, xmax, factor
	}
	str = str[beg+1 : end]
	if !strings.Contains(str, ",") {
		return xmin, xmax, factor
	}

	toks := strings.Split(str, ",")
	for i, tok := range toks {
		toks[i] = strings.ToLower(strings.TrimSpace(tok))
	}

	switch len(toks) {
	case 2, 3:
	default:
		panic(fmt.Errorf("rdict: invalid ROOT range specification (too many commas): %q", str))
	}

	var nbits uint32 = 32
	if len(toks) == 3 {
		n, err := strconv.ParseUint(toks[2], 10, 32)
		if err != nil {
			panic(fmt.Errorf("rdict: could not parse nbits specification %q: %w", str, err))
		}
		nbits = uint32(n)
		if nbits < 2 || nbits > 32 {
			panic(fmt.Errorf("rdict: illegal nbits specification (nbits=%d outside of range [2,32])", nbits))
		}
	}

	fct := func(s string) float64 {
		switch {
		case strings.Contains(s, "pi"):
			var f float64
			switch {
			case strings.Contains(s, "2pi"), strings.Contains(s, "2*pi"), strings.Contains(s, "twopi"):
				f = 2 * math.Pi
			case strings.Contains(s, "pi/2"):
				f = math.Pi / 2
			case strings.Contains(s, "pi/4"):
				f = math.Pi / 4
			case strings.Contains(s, "pi"):
				f = math.Pi
			}
			if strings.Contains(s, "-") {
				f = -f
			}
			return f
		default:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(fmt.Errorf("rdict: could not parse range value %q: %w", s, err))
			}
			return f
		}
	}

	xmin = fct(toks[0])
	xmax = fct(toks[1])

	var bigint uint32
	switch {
	case nbits < 32:
		bigint = 1 << nbits
	default:
		bigint = 0xffffffff
	}
	if xmin < xmax {
		factor = float64(bigint) / (xmax - xmin)
	}
	if xmin >= xmax && nbits < 15 {
		xmin = float64(nbits) + 0.1
	}
	return xmin, xmax, factor
}

type StreamerElement struct {
	named  rbase.Named
	etype  rmeta.Enum // element type
	esize  int32      // size of element
	arrlen int32      // cumulative size of all array dims
	arrdim int32      // number of array dimensions
	maxidx [5]int32   // maximum array index for array dimension "dim"
	offset int32      // element offset in class
	ename  string     // data type name of data member
	xmin   float64    // minimum of data member if a range is specified [xmin.xmax.nbits]
	xmax   float64    // maximum of data member if a range is specified [xmin.xmax.nbits]
	factor float64    // conversion factor if a range is specified. factor = (1<<nbits/(xmax-xmin))
}

func (*StreamerElement) RVersion() int16 { return rvers.StreamerElement }

func (tse *StreamerElement) Class() string {
	return "TStreamerElement"
}

func (tse *StreamerElement) Name() string {
	return tse.named.Name()
}

func (tse *StreamerElement) Title() string {
	return tse.named.Title()
}

func (tse *StreamerElement) ArrayDim() int {
	return int(tse.arrdim)
}

func (tse *StreamerElement) ArrayDims() []int32 {
	return tse.maxidx[:tse.arrdim]
}

func (tse *StreamerElement) ArrayLen() int {
	return int(tse.arrlen)
}

func (tse *StreamerElement) Type() rmeta.Enum {
	return tse.etype
}

func (tse *StreamerElement) Offset() uintptr {
	return uintptr(tse.offset)
}

func (tse *StreamerElement) Size() uintptr {
	return uintptr(tse.esize)
}

func (tse *StreamerElement) TypeName() string {
	return tse.ename
}

func (tse *StreamerElement) XMin() float64 {
	return tse.xmin
}

func (tse *StreamerElement) XMax() float64 {
	return tse.xmax
}

func (tse *StreamerElement) Factor() float64 {
	return tse.factor
}

func (tse *StreamerElement) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tse.Class(), tse.RVersion())
	w.WriteObject(&tse.named)
	w.WriteI32(int32(tse.etype))
	w.WriteI32(tse.esize)
	w.WriteI32(tse.arrlen)
	w.WriteI32(tse.arrdim)
	w.WriteArrayI32(tse.maxidx[:])
	w.WriteString(tse.ename)

	switch {
	case tse.RVersion() == 3:
		w.WriteF64(tse.xmin)
		w.WriteF64(tse.xmax)
		w.WriteF64(tse.factor)
	case tse.RVersion() > 3:
		// FIXME(sbinet)
		// if (TestBit(kHasRange)) GetRange(GetTitle(),fXmin,fXmax,fFactor)
	}

	return w.SetHeader(hdr)
}

func (tse *StreamerElement) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tse.Class())
	if hdr.Vers > rvers.StreamerElement {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tse.Class(), hdr.Vers, tse.RVersion(),
		))
	}

	r.ReadObject(&tse.named)

	tse.etype = rmeta.Enum(r.ReadI32())
	tse.esize = r.ReadI32()
	tse.arrlen = r.ReadI32()
	tse.arrdim = r.ReadI32()
	if hdr.Vers == 1 {
		copy(tse.maxidx[:], r.ReadStaticArrayI32())
	} else {
		r.ReadArrayI32(tse.maxidx[:])
	}
	tse.ename = r.ReadString()

	if tse.etype == 11 && (tse.ename == "Bool_t" || tse.ename == "bool") {
		tse.etype = 18
	}

	// if vers <= 2 {
	// 	// FIXME(sbinet)
	// 	// tse.esize = tse.arrlen * gROOT->GetType(GetTypeName())->Size()
	// }
	switch {
	default:
		tse.xmin = 0
		tse.xmax = 0
		tse.factor = 0
	case hdr.Vers == 3:
		tse.xmin = r.ReadF64()
		tse.xmax = r.ReadF64()
		tse.factor = r.ReadF64()
	case hdr.Vers > 3:
		tse.xmin, tse.xmax, tse.factor = Element{}.getRange(tse.Title())
	}

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerBase struct {
	StreamerElement
	vbase int32 // version number of the base class
}

func NewStreamerBase(se StreamerElement, vbase int32) *StreamerBase {
	return &StreamerBase{StreamerElement: se, vbase: vbase}
}

func (*StreamerBase) RVersion() int16 { return rvers.StreamerBase }

func (tsb *StreamerBase) Class() string {
	return "TStreamerBase"
}

// Base returns the base class' version.
func (tsb *StreamerBase) Base() int {
	return int(tsb.vbase)
}

func (tsb *StreamerBase) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tsb.Class(), tsb.RVersion())
	w.WriteObject(&tsb.StreamerElement)
	w.WriteI32(tsb.vbase)

	return w.SetHeader(hdr)
}

func (tsb *StreamerBase) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tsb.Class())
	if hdr.Vers > rvers.StreamerBase {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tsb.Class(), hdr.Vers, tsb.RVersion(),
		))
	}

	r.ReadObject(&tsb.StreamerElement)

	if hdr.Vers > 2 {
		tsb.vbase = r.ReadI32()
	}

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerBasicType struct {
	StreamerElement
}

func (*StreamerBasicType) RVersion() int16 { return rvers.StreamerBasicType }

func (tsb *StreamerBasicType) Class() string {
	return "TStreamerBasicType"
}

func (tsb *StreamerBasicType) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tsb.Class(), tsb.RVersion())
	w.WriteObject(&tsb.StreamerElement)

	return w.SetHeader(hdr)
}

func (tsb *StreamerBasicType) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tsb.Class())
	if hdr.Vers > rvers.StreamerBasicType {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tsb.Class(), hdr.Vers, tsb.RVersion(),
		))
	}

	r.ReadObject(&tsb.StreamerElement)

	etype := tsb.StreamerElement.etype
	if rmeta.OffsetL < etype && etype < rmeta.OffsetP {
		etype -= rmeta.OffsetL
	}

	basic := true
	switch etype {
	case rmeta.Bool, rmeta.UChar, rmeta.Char:
		tsb.StreamerElement.esize = 1
	case rmeta.UShort, rmeta.Short:
		tsb.StreamerElement.esize = 2
	case rmeta.Bits, rmeta.UInt, rmeta.Int, rmeta.Counter:
		tsb.StreamerElement.esize = 4
	case rmeta.ULong, rmeta.ULong64, rmeta.Long, rmeta.Long64:
		tsb.StreamerElement.esize = 8
	case rmeta.Float, rmeta.Float16:
		tsb.StreamerElement.esize = 4
	case rmeta.Double, rmeta.Double32:
		tsb.StreamerElement.esize = 8
	case rmeta.CharStar:
		tsb.StreamerElement.esize = int32(ptrSize)
	default:
		basic = false
	}
	if basic && tsb.StreamerElement.arrlen > 0 {
		tsb.StreamerElement.esize *= tsb.StreamerElement.arrlen
	}

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerBasicPointer struct {
	StreamerElement
	cvers int32  // version number of the class with the counter
	cname string // name of data member holding the array count
	ccls  string // name of the class with the counter
}

func NewStreamerBasicPointer(se StreamerElement, cvers int32, cname, ccls string) *StreamerBasicPointer {
	return &StreamerBasicPointer{
		StreamerElement: se,
		cvers:           cvers,
		cname:           cname,
		ccls:            ccls,
	}
}

func (*StreamerBasicPointer) RVersion() int16 { return rvers.StreamerBasicPointer }

func (tsb *StreamerBasicPointer) Class() string {
	return "TStreamerBasicPointer"
}

func (tsb *StreamerBasicPointer) CountName() string {
	return tsb.cname
}

func (tsb *StreamerBasicPointer) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tsb.Class(), tsb.RVersion())
	w.WriteObject(&tsb.StreamerElement)
	w.WriteI32(tsb.cvers)
	w.WriteString(tsb.cname)
	w.WriteString(tsb.ccls)

	return w.SetHeader(hdr)
}

func (tsb *StreamerBasicPointer) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tsb.Class())
	if hdr.Vers > rvers.StreamerBasicPointer {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tsb.Class(), hdr.Vers, tsb.RVersion(),
		))
	}

	r.ReadObject(&tsb.StreamerElement)

	tsb.cvers = r.ReadI32()
	tsb.cname = r.ReadString()
	tsb.ccls = r.ReadString()

	r.CheckHeader(hdr)
	return r.Err()
}

// StreamerLoop represents a streamer for a var-length array of a non-basic type.
type StreamerLoop struct {
	StreamerElement
	cvers  int32  // version number of the class with the counter
	cname  string // name of data member holding the array count
	cclass string // name of the class with the counter
}

func NewStreamerLoop(se StreamerElement, cvers int32, cname, cclass string) *StreamerLoop {
	se.etype = rmeta.StreamLoop
	return &StreamerLoop{
		StreamerElement: se,
		cvers:           cvers,
		cname:           cname,
		cclass:          cclass,
	}
}

func (*StreamerLoop) RVersion() int16 { return rvers.StreamerLoop }

func (*StreamerLoop) Class() string {
	return "TStreamerLoop"
}

func (tsl *StreamerLoop) CountName() string {
	return tsl.cname
}

func (tsl *StreamerLoop) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tsl.Class(), tsl.RVersion())
	w.WriteObject(&tsl.StreamerElement)
	w.WriteI32(tsl.cvers)
	w.WriteString(tsl.cname)
	w.WriteString(tsl.cclass)

	return w.SetHeader(hdr)
}

func (tsl *StreamerLoop) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tsl.Class())
	if hdr.Vers > rvers.StreamerLoop {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tsl.Class(), hdr.Vers, tsl.RVersion(),
		))
	}

	r.ReadObject(&tsl.StreamerElement)
	tsl.cvers = r.ReadI32()
	tsl.cname = r.ReadString()
	tsl.cclass = r.ReadString()

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerObject struct {
	StreamerElement
}

func (*StreamerObject) RVersion() int16 { return rvers.StreamerObject }

func (tso *StreamerObject) Class() string {
	return "TStreamerObject"
}

func (tso *StreamerObject) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tso.Class(), tso.RVersion())
	w.WriteObject(&tso.StreamerElement)
	return w.SetHeader(hdr)
}

func (tso *StreamerObject) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tso.Class())
	if hdr.Vers > rvers.StreamerObject {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tso.Class(), hdr.Vers, tso.RVersion(),
		))
	}

	r.ReadObject(&tso.StreamerElement)

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerObjectPointer struct {
	StreamerElement
}

func (*StreamerObjectPointer) RVersion() int16 { return rvers.StreamerObjectPointer }

func (tso *StreamerObjectPointer) Class() string {
	return "TStreamerObjectPointer"
}

func (tso *StreamerObjectPointer) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tso.Class(), tso.RVersion())
	w.WriteObject(&tso.StreamerElement)
	return w.SetHeader(hdr)
}

func (tso *StreamerObjectPointer) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tso.Class())
	if hdr.Vers > rvers.StreamerObjectPointer {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tso.Class(), hdr.Vers, tso.RVersion(),
		))
	}

	r.ReadObject(&tso.StreamerElement)

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerObjectAny struct {
	StreamerElement
}

func (*StreamerObjectAny) RVersion() int16 { return rvers.StreamerObjectAny }

func (tso *StreamerObjectAny) Class() string {
	return "TStreamerObjectAny"
}

func (tso *StreamerObjectAny) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tso.Class(), tso.RVersion())
	w.WriteObject(&tso.StreamerElement)

	return w.SetHeader(hdr)
}

func (tso *StreamerObjectAny) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tso.Class())
	if hdr.Vers > rvers.StreamerObjectAny {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tso.Class(), hdr.Vers, tso.RVersion(),
		))
	}

	r.ReadObject(&tso.StreamerElement)

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerObjectAnyPointer struct {
	StreamerElement
}

func (*StreamerObjectAnyPointer) RVersion() int16 { return rvers.StreamerObjectAnyPointer }

func (tso *StreamerObjectAnyPointer) Class() string {
	return "TStreamerObjectAnyPointer"
}

func (tso *StreamerObjectAnyPointer) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tso.Class(), tso.RVersion())
	w.WriteObject(&tso.StreamerElement)

	return w.SetHeader(hdr)
}

func (tso *StreamerObjectAnyPointer) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tso.Class())
	if hdr.Vers > rvers.StreamerObjectAnyPointer {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tso.Class(), hdr.Vers, tso.RVersion(),
		))
	}

	r.ReadObject(&tso.StreamerElement)

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerString struct {
	StreamerElement
}

func (*StreamerString) RVersion() int16 { return rvers.StreamerString }

func (tss *StreamerString) Class() string {
	return "TStreamerString"
}

func (tss *StreamerString) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tss.Class(), tss.RVersion())
	w.WriteObject(&tss.StreamerElement)

	return w.SetHeader(hdr)
}

func (tss *StreamerString) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tss.Class())
	if hdr.Vers > rvers.StreamerString {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tss.Class(), hdr.Vers, tss.RVersion(),
		))
	}

	r.ReadObject(&tss.StreamerElement)

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerSTL struct {
	StreamerElement
	vtype rmeta.ESTLType // type of STL vector
	ctype rmeta.Enum     // STL contained type
}

func NewStreamerSTL(name string, vtype rmeta.ESTLType, ctype rmeta.Enum) *StreamerSTL {
	return &StreamerSTL{
		StreamerElement: StreamerElement{
			named: *rbase.NewNamed(name, ""),
			esize: int32(ptrSize + 2*intSize),
			ename: rmeta.STLNameFrom(name, vtype, ctype),
			etype: rmeta.Streamer,
		},
		vtype: vtype,
		ctype: ctype,
	}
}

// NewCxxStreamerSTL creates a new StreamerSTL from C++ informations.
func NewCxxStreamerSTL(se StreamerElement, vtype rmeta.ESTLType, ctype rmeta.Enum) *StreamerSTL {
	return &StreamerSTL{
		StreamerElement: se,
		vtype:           vtype,
		ctype:           ctype,
	}
}

func (*StreamerSTL) RVersion() int16 { return rvers.StreamerSTL }

func (tss *StreamerSTL) Class() string {
	return "TStreamerSTL"
}

func (tss *StreamerSTL) ElemTypeName() []string {
	return rmeta.CxxTemplateFrom(tss.ename).Args
}

func (tss *StreamerSTL) ContainedType() rmeta.Enum {
	return tss.ctype
}

func (tss *StreamerSTL) STLType() rmeta.ESTLType {
	return tss.vtype
}

func (tss *StreamerSTL) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tss.Class(), tss.RVersion())
	w.WriteObject(&tss.StreamerElement)
	w.WriteI32(int32(tss.vtype))
	w.WriteI32(int32(tss.ctype))

	return w.SetHeader(hdr)
}

func (tss *StreamerSTL) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tss.Class())
	if hdr.Vers > rvers.StreamerSTL {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tss.Class(), hdr.Vers, tss.RVersion(),
		))
	}

	r.ReadObject(&tss.StreamerElement)

	tss.vtype = rmeta.ESTLType(r.ReadI32())
	tss.ctype = rmeta.Enum(r.ReadI32())

	if tss.vtype == rmeta.STLmultimap || tss.vtype == rmeta.STLset {
		switch {
		case strings.HasPrefix(tss.StreamerElement.ename, "std::set") || strings.HasPrefix(tss.StreamerElement.ename, "set"):
			tss.vtype = rmeta.STLset
		case strings.HasPrefix(tss.StreamerElement.ename, "std::multimap") || strings.HasPrefix(tss.StreamerElement.ename, "multimap"):
			tss.vtype = rmeta.STLmultimap
		}
	}

	r.CheckHeader(hdr)
	return r.Err()
}

// func (tss *StreamerSTL) isaPointer() bool {
// 	tname := tss.StreamerElement.ename
// 	return strings.HasSuffix(tname, "*")
// }

type StreamerSTLstring struct {
	StreamerSTL
}

func (*StreamerSTLstring) RVersion() int16 { return rvers.StreamerSTLstring }

func (tss *StreamerSTLstring) Class() string {
	return "TStreamerSTLstring"
}

func (tss *StreamerSTLstring) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tss.Class(), tss.RVersion())
	w.WriteObject(&tss.StreamerSTL)

	return w.SetHeader(hdr)
}

func (tss *StreamerSTLstring) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tss.Class())
	if hdr.Vers > rvers.StreamerSTLstring {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tss.Class(), hdr.Vers, tss.RVersion(),
		))
	}

	r.ReadObject(&tss.StreamerSTL)

	r.CheckHeader(hdr)
	return r.Err()
}

type StreamerArtificial struct {
	StreamerElement
}

func (*StreamerArtificial) RVersion() int16 { return rvers.StreamerArtificial }

func (tss *StreamerArtificial) Class() string {
	return "TStreamerArtificial"
}

func (tsa *StreamerArtificial) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tsa.Class(), tsa.RVersion())
	w.WriteObject(&tsa.StreamerElement)

	return w.SetHeader(hdr)
}

func (tsa *StreamerArtificial) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tsa.Class())
	if hdr.Vers > rvers.StreamerArtificial {
		panic(fmt.Errorf(
			"rdict: invalid %s version=%d > %d",
			tsa.Class(), hdr.Vers, tsa.RVersion(),
		))
	}

	r.ReadObject(&tsa.StreamerElement)

	r.CheckHeader(hdr)
	return r.Err()
}

func genChecksum(name string, elems []rbytes.StreamerElement) uint32 {
	var (
		id   uint32
		hash = func(s string) {
			for _, v := range []byte(s) {
				id = id*3 + uint32(v)
			}
		}
	)

	hash(name)

	// FIXME(sbinet): handle base-classes for std::pair<K,V>
	for _, se := range elems {
		//if se, ok := se.(*StreamerBase); ok {
		//	// FIXME(sbinet): get base checksum.
		//}

		// FIXME(sbinet): add enum handling.
		hash(se.Name())
		hash(se.TypeName())

		for _, v := range se.ArrayDims() {
			id = id*3 + uint32(v)
		}
		title := se.Title()
		beg := strings.Index(title, "[")
		if beg != 0 {
			continue
		}
		end := strings.Index(title, "]")
		hash(title[1:end])
	}

	return id
}

func init() {
	{
		f := func() reflect.Value {
			o := &StreamerInfo{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerInfo", f)
	}

	{
		f := func() reflect.Value {
			o := &StreamerElement{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerElement", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerBase{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerBase", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerBasicType{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerBasicType", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerBasicPointer{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerBasicPointer", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerLoop{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerLoop", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerObject{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerObject", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerObjectPointer{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerObjectPointer", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerObjectAny{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerObjectAny", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerObjectAnyPointer{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerObjectAnyPointer", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerString{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerString", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerSTL{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerSTL", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerSTLstring{
				StreamerSTL: StreamerSTL{
					vtype: rmeta.STLany,
					ctype: rmeta.STLstring,
				},
			}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerSTLstring", f)
	}
	{
		f := func() reflect.Value {
			o := &StreamerArtificial{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TStreamerArtificial", f)
	}
}

var (
	_ root.Object         = (*StreamerInfo)(nil)
	_ root.Named          = (*StreamerInfo)(nil)
	_ rbytes.StreamerInfo = (*StreamerInfo)(nil)
	_ rbytes.Marshaler    = (*StreamerInfo)(nil)
	_ rbytes.Unmarshaler  = (*StreamerInfo)(nil)

	_ root.Object            = (*StreamerElement)(nil)
	_ root.Named             = (*StreamerElement)(nil)
	_ rbytes.StreamerElement = (*StreamerElement)(nil)
	_ rbytes.Marshaler       = (*StreamerElement)(nil)
	_ rbytes.Unmarshaler     = (*StreamerElement)(nil)

	_ root.Object            = (*StreamerBase)(nil)
	_ root.Named             = (*StreamerBase)(nil)
	_ rbytes.StreamerElement = (*StreamerBase)(nil)
	_ rbytes.Marshaler       = (*StreamerBase)(nil)
	_ rbytes.Unmarshaler     = (*StreamerBase)(nil)

	_ root.Object            = (*StreamerBasicType)(nil)
	_ root.Named             = (*StreamerBasicType)(nil)
	_ rbytes.StreamerElement = (*StreamerBasicType)(nil)
	_ rbytes.Marshaler       = (*StreamerBasicType)(nil)
	_ rbytes.Unmarshaler     = (*StreamerBasicType)(nil)

	_ root.Object            = (*StreamerBasicPointer)(nil)
	_ root.Named             = (*StreamerBasicPointer)(nil)
	_ rbytes.StreamerElement = (*StreamerBasicPointer)(nil)
	_ rbytes.Marshaler       = (*StreamerBasicPointer)(nil)
	_ rbytes.Unmarshaler     = (*StreamerBasicPointer)(nil)

	_ root.Object            = (*StreamerLoop)(nil)
	_ root.Named             = (*StreamerLoop)(nil)
	_ rbytes.StreamerElement = (*StreamerLoop)(nil)
	_ rbytes.Marshaler       = (*StreamerLoop)(nil)
	_ rbytes.Unmarshaler     = (*StreamerLoop)(nil)

	_ root.Object            = (*StreamerObject)(nil)
	_ root.Named             = (*StreamerObject)(nil)
	_ rbytes.StreamerElement = (*StreamerObject)(nil)
	_ rbytes.Marshaler       = (*StreamerObject)(nil)
	_ rbytes.Unmarshaler     = (*StreamerObject)(nil)

	_ root.Object            = (*StreamerObjectPointer)(nil)
	_ root.Named             = (*StreamerObjectPointer)(nil)
	_ rbytes.StreamerElement = (*StreamerObjectPointer)(nil)
	_ rbytes.Marshaler       = (*StreamerObjectPointer)(nil)
	_ rbytes.Unmarshaler     = (*StreamerObjectPointer)(nil)

	_ root.Object            = (*StreamerObjectAny)(nil)
	_ root.Named             = (*StreamerObjectAny)(nil)
	_ rbytes.StreamerElement = (*StreamerObjectAny)(nil)
	_ rbytes.Marshaler       = (*StreamerObjectAny)(nil)
	_ rbytes.Unmarshaler     = (*StreamerObjectAny)(nil)

	_ root.Object            = (*StreamerString)(nil)
	_ root.Named             = (*StreamerString)(nil)
	_ rbytes.StreamerElement = (*StreamerString)(nil)
	_ rbytes.Marshaler       = (*StreamerString)(nil)
	_ rbytes.Unmarshaler     = (*StreamerString)(nil)

	_ root.Object            = (*StreamerSTL)(nil)
	_ root.Named             = (*StreamerSTL)(nil)
	_ rbytes.StreamerElement = (*StreamerSTL)(nil)
	_ rbytes.Marshaler       = (*StreamerSTL)(nil)
	_ rbytes.Unmarshaler     = (*StreamerSTL)(nil)

	_ root.Object            = (*StreamerSTLstring)(nil)
	_ root.Named             = (*StreamerSTLstring)(nil)
	_ rbytes.StreamerElement = (*StreamerSTLstring)(nil)
	_ rbytes.Marshaler       = (*StreamerSTLstring)(nil)
	_ rbytes.Unmarshaler     = (*StreamerSTLstring)(nil)

	_ root.Object            = (*StreamerArtificial)(nil)
	_ root.Named             = (*StreamerArtificial)(nil)
	_ rbytes.StreamerElement = (*StreamerArtificial)(nil)
	_ rbytes.Marshaler       = (*StreamerArtificial)(nil)
	_ rbytes.Unmarshaler     = (*StreamerArtificial)(nil)
)
