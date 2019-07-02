// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type tleaf struct {
	rvers    int16
	named    rbase.Named
	len      int
	etype    int
	offset   int
	hasrange bool
	unsigned bool
	count    leafCount
	branch   Branch
}

func newLeaf(name string, len, etype, offset int, hasrange, unsigned bool, count leafCount, b Branch) tleaf {
	return tleaf{
		rvers:    rvers.Leaf,
		named:    *rbase.NewNamed(name, ""),
		len:      len,
		etype:    etype,
		offset:   offset,
		hasrange: hasrange,
		unsigned: unsigned,
		count:    count,
		branch:   b,
	}
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

func (leaf *tleaf) readFromBasket(r *rbytes.RBuffer) error {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) scan(r *rbytes.RBuffer, ptr interface{}) error {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) setAddress(ptr interface{}) error {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) writeToBasket(w *rbytes.WBuffer) error {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) TypeName() string {
	panic("not implemented: " + leaf.Name())
}

func (leaf *tleaf) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
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
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

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
		panic(errors.Errorf("rtree: invalid typename for leaf %q", leaf.Name()))
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

func (leaf *tleafElement) readFromBasket(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.rstreamer == nil {
		panic(errors.Errorf("rtree: nil streamer (leaf: %s)", leaf.Name()))
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

func (leaf *tleafElement) writeToBasket(w *rbytes.WBuffer) error {
	panic("not implemented")
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
				panic(errors.Wrap(err, "could not infer leaf array dimension"))
			}
			dims[i] = dim
		}
	}

	return dims
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
