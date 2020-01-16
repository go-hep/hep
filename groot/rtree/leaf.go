// Copyright 2017 The go-hep Authors. All rights reserved.
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
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
)

type tleaf struct {
	named    rbase.Named
	len      int       // number of fixed length elements in the leaf's data.
	etype    int       // number of bytes for this data type
	offset   int       // offset in ClonesArray object
	hasrange bool      // whether the leaf has a range
	unsigned bool      // whether the leaf holds unsigned data (uint8, uint16, uint32 or uint64)
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

func (leaf *tleaf) ArrayDim() int {
	return strings.Count(leaf.named.Title(), "[")
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

func (leaf *tleaf) MaxIndex() []int {
	panic("not implemented: " + leaf.Name())
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

func (leaf *tleaf) Value(i int) interface{} {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) value() interface{} {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) readFromBuffer(r *rbytes.RBuffer) error {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) scan(r *rbytes.RBuffer, ptr interface{}) error {
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
	leaf.named.MarshalROOT(w)

	w.WriteI32(int32(leaf.len))
	w.WriteI32(int32(leaf.etype))
	w.WriteI32(int32(leaf.offset))
	w.WriteBool(leaf.hasrange)
	w.WriteBool(leaf.unsigned)
	w.WriteObjectAny(leaf.count)
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
	streamers []rbytes.StreamerElement
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
		panic(xerrors.Errorf("rtree: invalid typename for leaf %q", leaf.Name()))
	}
	return name
}

func (leaf *tleafElement) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
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

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

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
		panic(xerrors.Errorf("rtree: nil streamer (leaf: %s)", leaf.Name()))
	}

	err := leaf.rstreamer.RStreamROOT(r)
	if err != nil {
		return err
	}

	return nil
}

func (leaf *tleafElement) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	rv := reflect.Indirect(reflect.ValueOf(ptr))
	switch rv.Kind() {
	case reflect.Struct:
		for i := 0; i < rv.Type().NumField(); i++ {
			f := rv.Field(i)
			ft := rv.Type().Field(i)
			f.Set(leaf.src.FieldByName(ft.Name))
		}
	case reflect.Array:
		reflect.Copy(rv, leaf.src)
	case reflect.Slice:
		if rv.UnsafeAddr() != leaf.src.UnsafeAddr() {
			sli := leaf.src
			rv.Set(reflect.MakeSlice(sli.Type(), sli.Len(), sli.Cap()))
			reflect.Copy(rv, sli)
		}
	default:
		rv.Set(leaf.src)
	}
	return r.Err()
}

func (leaf *tleafElement) setAddress(ptr interface{}) error {
	var err error
	leaf.ptr = ptr
	leaf.src = reflect.ValueOf(leaf.ptr).Elem()

	var impl rstreamerImpl
	sictx := leaf.branch.getTree().getFile()
	for _, elt := range leaf.streamers {
		impl.funcs = append(impl.funcs, rstreamerFrom(elt, ptr, leaf.count, sictx))
	}
	leaf.rstreamer = &impl
	return err
}

func (leaf *tleafElement) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	panic("not implemented")
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
				panic(xerrors.Errorf("could not infer leaf array dimension: %w", err))
			}
			dims[i] = dim
		}
	}

	return dims
}

func maxI64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
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

func newLeafFromWVar(b Branch, v WriteVar) (Leaf, error) {
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
	)

	switch kind {
	case reflect.Slice:
		lc := b.Leaf(v.Count)
		if lc == nil {
			leaves := b.Leaves()
			names := make([]string, len(leaves))
			for i, ll := range leaves {
				names[i] = ll.Name()
			}
			return nil, xerrors.Errorf(
				"could not find leaf count %q from branch %q for slice (name=%q, type=%T) among: %q",
				v.Count, b.Name(), v.Name, v.Value, names,
			)
		}
		lcc, ok := lc.(leafCount)
		if !ok {
			return nil, xerrors.Errorf(
				"leaf count %q from branch %q for slice (name=%q, type=%T) is not a LeafCount",
				v.Count, b.Name(), v.Name, v.Value,
			)
		}
		count = lcc
		kind = rt.Elem().Kind()
	case reflect.Struct:
		panic("not implemented")
	}

	switch kind {
	case reflect.Bool:
		leaf = newLeafO(b, v.Name, shape, false, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Uint8:
		leaf = newLeafB(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Uint16:
		leaf = newLeafS(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Uint32:
		leaf = newLeafI(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Uint64:
		leaf = newLeafL(b, v.Name, shape, unsigned, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Int8:
		leaf = newLeafB(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Int16:
		leaf = newLeafS(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Int32:
		leaf = newLeafI(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Int64:
		leaf = newLeafL(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	case reflect.Float32:
		switch rt {
		case reflect.TypeOf(float32(0)), reflect.TypeOf([]float32(nil)):
			leaf = newLeafF(b, v.Name, shape, signed, count)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
		case reflect.TypeOf(root.Float16(0)), reflect.TypeOf([]root.Float16(nil)):
			leaf = newLeafF16(b, v.Name, shape, signed, count, nil)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
		default:
			panic(xerrors.Errorf("invalid type %T", v.Value))
		}
	case reflect.Float64:
		switch rt {
		case reflect.TypeOf(float64(0)), reflect.TypeOf([]float64(nil)):
			leaf = newLeafD(b, v.Name, shape, signed, count)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
		case reflect.TypeOf(root.Double32(0)), reflect.TypeOf([]root.Double32(nil)):
			leaf = newLeafD32(b, v.Name, shape, signed, count, nil)
			err := leaf.setAddress(v.Value)
			if err != nil {
				return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
			}
		default:
			panic(xerrors.Errorf("invalid type %T", v.Value))
		}
	case reflect.String:
		leaf = newLeafC(b, v.Name, shape, signed, count)
		err := leaf.setAddress(v.Value)
		if err != nil {
			return nil, xerrors.Errorf("could not set leaf address for %q: %w", v.Name, err)
		}
	default:
		return nil, xerrors.Errorf("rtree: invalid write-var (name=%q) type %T", v.Name, v.Value)
	}

	return leaf, nil
}

var (
	_ root.Object        = (*tleaf)(nil)
	_ root.Named         = (*tleaf)(nil)
	_ Leaf               = (*tleaf)(nil)
	_ rbytes.Marshaler   = (*tleaf)(nil)
	_ rbytes.Unmarshaler = (*tleaf)(nil)

	_ root.Object        = (*tleafElement)(nil)
	_ root.Named         = (*tleafElement)(nil)
	_ Leaf               = (*tleafElement)(nil)
	_ rbytes.Marshaler   = (*tleafElement)(nil)
	_ rbytes.Unmarshaler = (*tleafElement)(nil)
)
