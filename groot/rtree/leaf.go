// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type tleaf struct {
	named    rbase.Named
	len      int  // number of fixed length elements in the leaf's data.
	etype    int  // number of bytes for this data type
	offset   int  // offset in ClonesArray object
	hasrange bool // whether the leaf has a range
	unsigned bool // whether the leaf holds unsigned data (uint8, uint16, uint32 or uint64)
	shape    []int
	count    leafCount // leaf count if variable length
	branch   Branch    // supporting branch of this leaf
}

func newLeaf(name string, shape []int, etype, offset int, hasrange, unsigned bool, count leafCount, b Branch) tleaf {
	var (
		nelems = 1
		title  = new(strings.Builder)
	)

	title.WriteString(name)
	switch {
	case count != nil:
		fmt.Fprintf(title, "[%s]", count.Name())
	default:
		for _, dim := range shape {
			nelems *= dim
			fmt.Fprintf(title, "[%d]", dim)
		}
	}
	return tleaf{
		named:    *rbase.NewNamed(name, title.String()),
		len:      nelems,
		etype:    etype,
		offset:   offset,
		hasrange: hasrange,
		unsigned: unsigned,
		shape:    shape,
		count:    count,
		branch:   b,
	}
}

func (*tleaf) RVersion() int16 {
	return rvers.Leaf
}

// Name returns the name of the instance
func (leaf *tleaf) Name() string {
	return leaf.named.Name()
}

// Title returns the title of the instance
func (leaf *tleaf) Title() string {
	return leaf.named.Title()
}

func (leaf *tleaf) Class() string {
	return "TLeaf"
}

func (leaf *tleaf) Shape() []int {
	return leaf.shape
}

func (leaf *tleaf) setBranch(b Branch) {
	leaf.branch = b
}

func (leaf *tleaf) Branch() Branch {
	return leaf.branch
}

func (leaf *tleaf) HasRange() bool {
	return leaf.hasrange
}

func (leaf *tleaf) IsUnsigned() bool {
	return leaf.unsigned
}

func (leaf *tleaf) LeafCount() Leaf {
	return leaf.count
}

func (leaf *tleaf) Len() int {
	if leaf.count != nil {
		// variable length array
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		return leaf.len * n
	}
	return leaf.len
}

func (leaf *tleaf) LenType() int {
	return leaf.etype
}

func (leaf *tleaf) Offset() int {
	return leaf.offset
}

func (leaf *tleaf) Kind() reflect.Kind {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) Type() reflect.Type {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) readFromBuffer(r *rbytes.RBuffer) error {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) setAddress(ptr interface{}) error {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) TypeName() string {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.RVersion())
	_, _ = leaf.named.MarshalROOT(w)

	w.WriteI32(int32(leaf.len))
	w.WriteI32(int32(leaf.etype))
	w.WriteI32(int32(leaf.offset))
	w.WriteBool(leaf.hasrange)
	w.WriteBool(leaf.unsigned)
	_ = w.WriteObjectAny(leaf.count)
	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *tleaf) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion(leaf.Class())

	if err := leaf.named.UnmarshalROOT(r); err != nil {
		return err
	}
	leaf.shape = leafDims(leaf.Title())

	leaf.len = int(r.ReadI32())
	leaf.etype = int(r.ReadI32())
	leaf.offset = int(r.ReadI32())
	leaf.hasrange = r.ReadBool()
	leaf.unsigned = r.ReadBool()

	leaf.count = nil
	ptr := r.ReadObjectAny()
	if ptr != nil {
		leaf.count = ptr.(leafCount)
	}

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
	if leaf.len == 0 {
		leaf.len = 1
	}

	return r.Err()
}

func (leaf *tleaf) canGenerateOffsetArray() bool {
	return leaf.count != nil
}

func (leaf *tleaf) computeOffsetArray(base, nevts int) []int32 {
	o := make([]int32, nevts)
	if leaf.count == nil {
		return o
	}

	var (
		hdr           = tleafHdrSize
		origEntry     = maxI64(leaf.branch.getReadEntry(), 0) // -1 indicates to start at the beginning
		origLeafEntry = leaf.count.Branch().getReadEntry()
		sz            int
		offset        = int32(base)
	)
	for i := 0; i < nevts; i++ {
		o[i] = offset
		leaf.count.Branch().getEntry(origEntry + int64(i))
		sz = leaf.count.ivalue()
		offset += int32(leaf.etype*sz + hdr)
	}
	leaf.count.Branch().getEntry(origLeafEntry)

	return o
}

// tleafObject is a Leaf for a general object derived from Object.
type tleafObject struct {
	tleaf
	virtual bool
	typ     reflect.Type
}

func (*tleafObject) RVersion() int16 {
	return rvers.LeafObject
}

func (leaf *tleafObject) Class() string {
	return "TLeafObject"
}

func (leaf *tleafObject) Type() reflect.Type {
	return leaf.typ
}

func (leaf *tleafObject) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.RVersion())
	_, _ = leaf.tleaf.MarshalROOT(w)
	w.WriteBool(leaf.virtual)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *tleafObject) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	_ = leaf.tleaf.UnmarshalROOT(r)

	if vers < 4 {
		panic(fmt.Errorf(
			"rtree: TLeafObject %q with version [%v] is not supported (too old)",
			leaf.Name(),
			vers,
		))
	}
	leaf.virtual = r.ReadBool()

	if !rtypes.Factory.HasKey(leaf.Title()) {
		return fmt.Errorf("rtree: could not find type %q for TLeafObject %q", leaf.Title(), leaf.Name())
	}
	leaf.typ = rtypes.Factory.Get(leaf.Title())().Type().Elem()

	r.CheckByteCount(pos, bcnt, beg, leaf.Class())
	return r.Err()
}

const (
	tleafHdrSize        = 0
	tleafElementHdrSize = 1
)

// tleafElement is a Leaf for a general object derived from Object.
type tleafElement struct {
	rvers int16
	tleaf
	id    int32 // element serial number in fInfo
	ltype int32 // leaf type

	ptr       interface{}
	src       reflect.Value
	rstreamer rbytes.RStreamer
	wstreamer rbytes.WStreamer
	streamers []rbytes.StreamerElement
}

func (*tleafElement) RVersion() int16 {
	return rvers.LeafElement
}

func (leaf *tleafElement) Class() string {
	return "TLeafElement"
}

func (leaf *tleafElement) ivalue() int {
	return int(leaf.src.Int())
}

func (leaf *tleafElement) imax() int {
	panic("not implemented")
}

func (leaf *tleafElement) Kind() reflect.Kind {
	return leaf.src.Kind()
}

func (leaf *tleafElement) Type() reflect.Type {
	return leaf.src.Type()
}

func (leaf *tleafElement) TypeName() string {
	name := leaf.src.Type().Name()
	if name == "" {
		panic(fmt.Errorf("rtree: invalid typename for leaf %q", leaf.Name()))
	}
	return name
}

func (leaf *tleafElement) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
	_, _ = leaf.tleaf.MarshalROOT(w)
	w.WriteI32(leaf.id)
	w.WriteI32(leaf.ltype)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *tleafElement) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	_ = leaf.tleaf.UnmarshalROOT(r)
	leaf.id = r.ReadI32()
	leaf.ltype = r.ReadI32()

	r.CheckByteCount(pos, bcnt, beg, leaf.Class())
	return r.Err()
}

func (leaf *tleafElement) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.rstreamer == nil {
		panic(fmt.Errorf("rtree: nil streamer (leaf: %s)", leaf.Name()))
	}

	err := leaf.rstreamer.RStreamROOT(r)
	if err != nil {
		return err
	}

	return nil
}

func (leaf *tleafElement) setAddress(ptr interface{}) error {
	leaf.ptr = ptr
	leaf.src = reflect.ValueOf(leaf.ptr).Elem()

	if leaf.rstreamer != nil {
		return leaf.setReadAddress(ptr)
	}

	if leaf.wstreamer != nil {
		return leaf.setWriteAddress(ptr)
	}

	return fmt.Errorf("rtree: leaf %q is neither read nor write", leaf.Name())
}

func (leaf *tleafElement) setReadAddress(ptr interface{}) error {
	err := leaf.rstreamer.(rbytes.Binder).Bind(ptr)
	if err != nil {
		return fmt.Errorf("rtree: could not bind read-streamer for leaf=%q (type=%s) to ptr=%T: %w",
			leaf.Name(), leaf.TypeName(), leaf.ptr, err,
		)
	}
	if leaf.count != nil {
		r, ok := leaf.rstreamer.(rbytes.Counter)
		if !ok {
			return fmt.Errorf(
				"rtree: could not set read-streamer counter for leaf=%q (type=%s)",
				leaf.Name(), leaf.TypeName(),
			)
		}
		err = r.Count(leaf.count.ivalue)
		if err != nil {
			return fmt.Errorf(
				"rtree: could not set read-streamer counter for leaf=%q (type=%s): %w",
				leaf.Name(), leaf.TypeName(), err,
			)
		}
	}
	return nil
}

func (leaf *tleafElement) setWriteAddress(ptr interface{}) error {
	err := leaf.wstreamer.(rbytes.Binder).Bind(ptr)
	if err != nil {
		return fmt.Errorf("rtree: could not bind write-streamer for leaf=%q (type=%s) to ptr=%T: %w",
			leaf.Name(), leaf.TypeName(), leaf.ptr, err,
		)
	}
	if leaf.count != nil {
		w, ok := leaf.wstreamer.(rbytes.Counter)
		if !ok {
			return fmt.Errorf(
				"rtree: could not set write-streamer counter for leaf=%q (type=%s)",
				leaf.Name(), leaf.TypeName(),
			)
		}
		err = w.Count(leaf.count.ivalue)
		if err != nil {
			return fmt.Errorf(
				"rtree: could not set write-streamer counter for leaf=%q (type=%s): %w",
				leaf.Name(), leaf.TypeName(), err,
			)
		}
	}
	return nil
}

func (leaf *tleafElement) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	if leaf.wstreamer == nil {
		panic(fmt.Errorf("rtree: nil write-streamer (leaf: %s)", leaf.Name()))
	}

	pos := w.Pos()
	err := leaf.wstreamer.WStreamROOT(w)
	return int(w.Pos() - pos), err
}

func (leaf *tleafElement) canGenerateOffsetArray() bool {
	return leaf.count != nil && leaf.tleaf.etype != 0
}

func (leaf *tleafElement) computeOffsetArray(base, nevts int) []int32 {
	o := make([]int32, nevts)
	if leaf.count == nil {
		return o
	}

	var (
		hdr           = tleafElementHdrSize
		origEntry     = maxI64(leaf.branch.getReadEntry(), 0) // -1 indicates to start at the beginning
		origLeafEntry = leaf.count.Branch().getReadEntry()
		sz            int
		offset        = int32(base)
	)
	for i := 0; i < nevts; i++ {
		o[i] = offset
		leaf.count.Branch().getEntry(origEntry + int64(i))
		sz = leaf.count.ivalue()
		offset += int32(leaf.etype*sz + hdr)
	}
	leaf.count.Branch().getEntry(origLeafEntry)

	return o
}

func init() {
	{
		f := func() reflect.Value {
			o := &tleaf{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TLeaf", f)
	}
	{
		f := func() reflect.Value {
			o := &tleafObject{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TLeafObject", f)
	}
	{
		f := func() reflect.Value {
			o := &tleafElement{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TLeafElement", f)
	}
}

var (
	reLeafDims = regexp.MustCompile(`\w*?\[(\d*)\]+?`)
)

func leafDims(s string) []int {
	out := reLeafDims.FindAllStringSubmatch(s, -1)
	if len(out) == 0 {
		return nil
	}

	dims := make([]int, len(out))
	for i := range out {
		v := out[i][1]
		switch v {
		case "":
			dims[i] = -1 // variable size
		default:
			dim, err := strconv.Atoi(v)
			if err != nil {
				panic(fmt.Errorf("could not infer leaf array dimension: %w", err))
			}
			dims[i] = dim
		}
	}

	return dims
}

func gotypeToROOTTypeCode(rt reflect.Type) string {
	switch rt.Kind() {
	case reflect.Bool:
		return "O"
	case reflect.String:
		return "C"
	case reflect.Int8:
		return "B"
	case reflect.Int16:
		return "S"
	case reflect.Int32:
		return "I"
	case reflect.Int64:
		return "L"
	case reflect.Uint8:
		return "b"
	case reflect.Uint16:
		return "s"
	case reflect.Uint32:
		return "i"
	case reflect.Uint64:
		return "l"
	case reflect.Float32:
		if rt == reflect.TypeOf(root.Float16(0)) {
			return "f"
		}
		return "F"
	case reflect.Float64:
		if rt == reflect.TypeOf(root.Double32(0)) {
			return "d"
		}
		return "D"
	case reflect.Array:
		return gotypeToROOTTypeCode(rt.Elem())
	case reflect.Slice:
		return gotypeToROOTTypeCode(rt.Elem())
	}
	panic("impossible")
}

func newLeafFromWVar(w *wtree, b Branch, v WriteVar, lvl int, cfg wopt) (Leaf, error) {
	const (
		signed   = false
		unsigned = true
	)

	var (
		rv        = reflect.Indirect(reflect.ValueOf(v.Value))
		rt, shape = flattenArrayType(rv.Type())
		kind      = rt.Kind()
		leaf      Leaf
		count     leafCount
		addLeaf   func(leaf Leaf)
	)

	switch b := b.(type) {
	case *tbranch:
		addLeaf = func(leaf Leaf) {
			b.leaves = append(b.leaves, leaf)
			w.ttree.leaves = append(w.ttree.leaves, leaf)
		}
	case *tbranchElement:
		addLeaf = func(leaf Leaf) {
			if _, ok := leaf.(*tleafElement); !ok {
				lb, ltyp := asLeafBase(leaf)

				if b.bup != nil {
					lb.named.SetName(b.bup.Name() + "." + lb.Name())
					lb.named.SetTitle(b.bup.Name() + "." + lb.Title())
				}
				leaf = &tleafElement{
					rvers: rvers.LeafElement,
					tleaf: *lb,
					id:    -1, // FIXME(sbinet): infer correct index
					ltype: int32(ltyp),
				}
			}

			b.leaves = append(b.leaves, leaf)
			w.ttree.leaves = append(w.ttree.leaves, leaf)
		}
	}

	switch kind {
	case reflect.Slice:
		lc := b.Leaf(v.Count)
		switch lc {
		case nil:
			// write as vector<T>.
			const (
				offset   = 0
				hasrange = false
				unsigned = false
			)
			base := newLeaf(v.Name, nil, int(rt.Size()), offset, hasrange, unsigned, count, b)
			leaf := &tleafElement{
				rvers: rvers.LeafElement,
				tleaf: base,
				id:    -1, // FIXME(sbinet): create proper serial number
				ltype: 2,  // FIXME(sbinet)
				ptr:   v.Value,
				src:   reflect.ValueOf(v.Value),
			}
			si := rdict.StreamerOf(w.ttree.f, reflect.TypeOf(v.Value).Elem())

			var err error
			leaf.wstreamer, err = si.NewWStreamer(rbytes.ObjectWise)
			if err != nil {
				return nil, fmt.Errorf("could not create w-streamer for leaf %q: %w", v.Name, err)
			}

			err = leaf.setAddress(v.Value)
			if err != nil {
				return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}

			addLeaf(leaf)
			return leaf, nil

		default:
			lcc, ok := lc.(leafCount)
			if !ok {
				return nil, fmt.Errorf(
					"leaf count %q from branch %q for slice (name=%q, type=%T) is not a LeafCount",
					v.Count, b.Name(), v.Name, v.Value,
				)
			}
			count = lcc
			kind = rt.Elem().Kind()
		}

	case reflect.Struct:
		const (
			offset   = 0
			hasrange = false
			unsigned = false
		)
		base := newLeaf(v.Name, nil, int(rt.Size()), offset, hasrange, unsigned, count, b)
		leaf := &tleafElement{
			rvers: rvers.LeafElement,
			tleaf: base,
			id:    -1, // FIXME(sbinet): create proper serial number
			ltype: -1, // FIXME(sbinet)
			ptr:   v.Value,
			src:   reflect.ValueOf(v.Value),
		}
		si := rdict.StreamerOf(w.ttree.f, reflect.TypeOf(v.Value).Elem())
		wstreamer, err := si.NewWStreamer(rbytes.ObjectWise)
		if err != nil {
			return nil, fmt.Errorf("could not create w-streamer for leaf %q: %w", v.Name, err)
		}
		leaf.wstreamer = wstreamer

		err = leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)
		return leaf, nil
	}

	switch kind {
	case reflect.Bool:
		leaf = newLeafO(b, v.Name, shape, false, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Uint8:
		leaf = newLeafB(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Uint16:
		leaf = newLeafS(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Uint32:
		leaf = newLeafI(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Uint64:
		leaf = newLeafL(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Int8:
		leaf = newLeafB(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Int16:
		leaf = newLeafS(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Int32:
		leaf = newLeafI(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Int64:
		leaf = newLeafL(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	case reflect.Float32:
		switch rt {
		case reflect.TypeOf(float32(0)), reflect.TypeOf([]float32(nil)):
			leaf = newLeafF(b, v.Name, shape, signed, count)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
			addLeaf(leaf)

		case reflect.TypeOf(root.Float16(0)), reflect.TypeOf([]root.Float16(nil)):
			leaf = newLeafF16(b, v.Name, shape, signed, count, nil)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
			addLeaf(leaf)

		default:
			panic(fmt.Errorf("invalid type %T", v.Value))
		}
	case reflect.Float64:
		switch rt {
		case reflect.TypeOf(float64(0)), reflect.TypeOf([]float64(nil)):
			leaf = newLeafD(b, v.Name, shape, signed, count)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
			addLeaf(leaf)

		case reflect.TypeOf(root.Double32(0)), reflect.TypeOf([]root.Double32(nil)):
			leaf = newLeafD32(b, v.Name, shape, signed, count, nil)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
			addLeaf(leaf)

		default:
			panic(fmt.Errorf("invalid type %T", v.Value))
		}
	case reflect.String:
		leaf = newLeafC(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, fmt.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
		addLeaf(leaf)

	default:
		return nil, fmt.Errorf("rtree: invalid write-var (name=%q) type %T", v.Name, v.Value)
	}

	return leaf, nil
}

func asLeafBase(leaf Leaf) (*tleaf, rmeta.Enum) {
	switch leaf := leaf.(type) {
	case *LeafO:
		return &leaf.tleaf, rmeta.Bool
	case *LeafB:
		return &leaf.tleaf, rmeta.Int8
	case *LeafS:
		return &leaf.tleaf, rmeta.Int16
	case *LeafI:
		return &leaf.tleaf, rmeta.Int32
	case *LeafL:
		return &leaf.tleaf, rmeta.Int64
	case *LeafF:
		return &leaf.tleaf, rmeta.Float32
	case *LeafD:
		return &leaf.tleaf, rmeta.Float64
	case *LeafF16:
		return &leaf.tleaf, rmeta.Float16
	case *LeafD32:
		return &leaf.tleaf, rmeta.Double32
	case *LeafC:
		return &leaf.tleaf, rmeta.CharStar // FIXME(sbinet): rmeta.Char?
	default:
		panic(fmt.Errorf("rtree: invalid leaf type %T", leaf))
	}
}

var (
	_ root.Object        = (*tleaf)(nil)
	_ root.Named         = (*tleaf)(nil)
	_ Leaf               = (*tleaf)(nil)
	_ rbytes.Marshaler   = (*tleaf)(nil)
	_ rbytes.Unmarshaler = (*tleaf)(nil)

	_ root.Object        = (*tleafObject)(nil)
	_ root.Named         = (*tleafObject)(nil)
	_ Leaf               = (*tleafObject)(nil)
	_ rbytes.Marshaler   = (*tleafObject)(nil)
	_ rbytes.Unmarshaler = (*tleafObject)(nil)

	_ root.Object        = (*tleafElement)(nil)
	_ root.Named         = (*tleafElement)(nil)
	_ Leaf               = (*tleafElement)(nil)
	_ rbytes.Marshaler   = (*tleafElement)(nil)
	_ rbytes.Unmarshaler = (*tleafElement)(nil)
)
