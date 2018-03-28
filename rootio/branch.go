// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
	"strings"
)

type tbranch struct {
	rvers          int16
	named          tnamed
	attfill        attfill
	compress       int         // compression level and algorithm
	basketSize     int         // initial size of Basket buffer
	entryOffsetLen int         // initial length of entryOffset table in the basket buffers
	writeBasket    int         // last basket number written
	entryNumber    int64       // current entry number (last one filled in this branch)
	iofeats        tioFeatures // IO features for newly-created baskets
	offset         int         // offset of this branch
	maxBaskets     int         // maximum number of baskets so far
	splitLevel     int         // branch split level
	entries        int64       // number of entries
	firstEntry     int64       // number of the first entry in this branch
	totBytes       int64       // total number of bytes in all leaves before compression
	zipBytes       int64       // total number of bytes in all leaves after compression
	branches       []Branch    // list of branches of this branch
	leaves         []Leaf      // list of leaves of this branch
	baskets        []Basket    // list of baskets of this branch

	basketBytes []int32 // length of baskets on file
	basketEntry []int64 // table of first entry in each basket
	basketSeek  []int64 // addresses of baskets on file

	fname string // named of file where buffers are stored (empty if in same file as Tree header)

	readbasket  int     // current basket number when reading
	readentry   int64   // current entry number when reading
	firstbasket int64   // first entry in the current basket
	nextbasket  int64   // next entry that will require us to go to the next basket
	basket      *Basket // pointer to the current basket
	basketBuf   []byte  // scratch space for the current basket

	tree Tree        // tree header
	btop Branch      // top-level parent branch in the tree
	bup  Branch      // parent branch
	dir  *tdirectory // directory where this branch's buffers are stored
}

func (b *tbranch) Name() string {
	return b.named.Name()
}

func (b *tbranch) Title() string {
	return b.named.Title()
}

func (b *tbranch) Class() string {
	return "TBranch"
}

func (b *tbranch) getTree() Tree {
	return b.tree
}

func (b *tbranch) setTree(t Tree) {
	b.tree = t
	for _, sub := range b.branches {
		sub.setTree(t)
	}
}

func (b *tbranch) Branches() []Branch {
	return b.branches
}

func (b *tbranch) Leaves() []Leaf {
	return b.leaves
}

func (b *tbranch) Branch(name string) Branch {
	for _, bb := range b.Branches() {
		if bb.Name() == name {
			return bb
		}
	}
	return nil
}

func (b *tbranch) Leaf(name string) Leaf {
	for _, lf := range b.Leaves() {
		if lf.Name() == name {
			return lf
		}
	}
	return nil
}

func (b *tbranch) GoType() reflect.Type {
	if len(b.Leaves()) == 1 {
		return b.leaves[0].Type()
	}
	fields := make([]reflect.StructField, len(b.leaves))
	for i, leaf := range b.leaves {
		ft := &fields[i]
		ft.Name = "ROOT_" + leaf.Name()
		etype := leaf.Type()
		switch {
		case leaf.LeafCount() != nil:
			etype = reflect.SliceOf(etype)
		case leaf.Len() > 1 && leaf.Kind() != reflect.String:
			etype = reflect.ArrayOf(leaf.Len(), etype)
		}
		ft.Type = etype
	}
	return reflect.StructOf(fields)
}

func (b *tbranch) getReadEntry() int64 {
	return b.readentry
}

func (b *tbranch) getEntry(i int64) {
	err := b.loadEntry(i)
	if err != nil {
		panic(errorf("rootio: branch [%s] failed to load entry %d: %v", b.Name(), i, err))
	}
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *tbranch) UnmarshalROOT(r *RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	b.rvers = vers

	b.tree = nil
	b.basket = nil
	b.firstbasket = -1
	b.nextbasket = -1

	if vers < 10 {
		panic(fmt.Errorf("rootio: too old TBanch version (%d<10)", vers))
	}

	if err := b.named.UnmarshalROOT(r); err != nil {
		return err
	}

	if err := b.attfill.UnmarshalROOT(r); err != nil {
		return err
	}

	b.compress = int(r.ReadI32())
	b.basketSize = int(r.ReadI32())
	b.entryOffsetLen = int(r.ReadI32())
	b.writeBasket = int(r.ReadI32())
	b.entryNumber = r.ReadI64()
	if b.rvers >= 13 {
		if err := b.iofeats.UnmarshalROOT(r); err != nil {
			return err
		}
	}
	b.offset = int(r.ReadI32())
	b.maxBaskets = int(r.ReadI32())
	b.splitLevel = int(r.ReadI32())
	b.entries = r.ReadI64()
	if vers >= 11 {
		b.firstEntry = r.ReadI64()
	}
	b.totBytes = r.ReadI64()
	b.zipBytes = r.ReadI64()

	{
		var branches tobjarray
		if err := branches.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
		b.branches = make([]Branch, branches.last+1)
		for i := range b.branches {
			br := branches.At(i).(Branch)
			b.branches[i] = br
		}
	}

	{
		var leaves tobjarray
		if err := leaves.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
		b.leaves = make([]Leaf, leaves.last+1)
		for i := range b.leaves {
			leaf := leaves.At(i).(Leaf)
			leaf.setBranch(b)
			b.leaves[i] = leaf
		}
	}
	{
		var baskets tobjarray
		if err := baskets.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
		b.baskets = make([]Basket, 0, baskets.last+1)
		for i := 0; i < baskets.last+1; i++ {
			bkt := baskets.At(i)
			// FIXME(sbinet) check why some are nil
			if bkt == nil {
				b.baskets = append(b.baskets, Basket{})
				continue
			}
			bk := bkt.(*Basket)
			b.baskets = append(b.baskets, *bk)
		}
	}

	b.basketBytes = nil
	b.basketEntry = nil
	b.basketSeek = nil

	/*isArray*/
	_ = r.ReadI8()
	b.basketBytes = r.ReadFastArrayI32(b.maxBaskets)[:b.writeBasket:b.writeBasket]

	/*isArray*/
	_ = r.ReadI8()
	// FIXME(sbinet) drop when go-1.9 isn't supported anymore.
	// workaround different gofmt rules.
	end := b.writeBasket + 1
	b.basketEntry = r.ReadFastArrayI64(b.maxBaskets)[:end:end]

	/*isArray*/
	_ = r.ReadI8()
	b.basketSeek = r.ReadFastArrayI64(b.maxBaskets)[:b.writeBasket:b.writeBasket]

	b.fname = r.ReadString()

	r.CheckByteCount(pos, bcnt, beg, "TBranch")

	if b.splitLevel == 0 && len(b.branches) > 0 {
		b.splitLevel = 1
	}

	return r.Err()
}

func (b *tbranch) loadEntry(ientry int64) error {
	var err error

	if len(b.basketBytes) == 0 {
		return nil
	}

	err = b.loadBasket(ientry)
	if err != nil {
		return err
	}

	err = b.basket.loadEntry(ientry - b.firstEntry)
	if err != nil {
		return err
	}

	for _, leaf := range b.leaves {
		err = b.basket.readLeaf(ientry-b.firstEntry, leaf)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *tbranch) loadBasket(entry int64) error {
	ib := b.findBasketIndex(entry)
	if ib < 0 {
		return errorf("rootio: no basket for entry %d", entry)
	}
	b.readentry = entry
	b.readbasket = ib
	b.nextbasket = b.basketEntry[ib+1]
	b.firstbasket = b.basketEntry[ib]
	if ib < len(b.baskets) {
		b.basket = &b.baskets[ib]
		b.firstEntry = b.basketEntry[ib]
		if b.basket.rbuf == nil {
			return b.setupBasket(b.basket, ib, entry)
		}
		return nil
	}

	b.baskets = append(b.baskets, Basket{})
	b.basket = &b.baskets[len(b.baskets)-1]
	return b.setupBasket(b.basket, ib, entry)
}

func (b *tbranch) findBasketIndex(entry int64) int {
	switch {
	case entry == 0:
		return 0
	case b.firstbasket <= entry && entry < b.nextbasket:
		return b.readbasket
	}
	/*
		    // binary search is not efficient for small slices (like basketEntry)
			// TODO(sbinet): test at which length of basketEntry it starts to be efficient.
			entries := b.basketEntry[1:]
			i := sort.Search(len(entries), func(i int) bool { return entries[i] >= entry })
			if b.basketEntry[i+1] == entry {
				return i + 1
			}
			return i
	*/

	for i := b.readbasket; i < len(b.basketEntry); i++ {
		v := b.basketEntry[i]
		if v > entry && v > 0 {
			return i - 1
		}
	}
	if entry == b.basketEntry[len(b.basketEntry)-1] {
		return -2 // len(b.basketEntry) - 1
	}
	return -1
}

func (b *tbranch) setupBasket(bk *Basket, ib int, entry int64) error {
	var err error
	if len(b.basketBuf) < int(b.basketBytes[ib]) {
		b.basketBuf = make([]byte, int(b.basketBytes[ib]))
	}
	buf := b.basketBuf[:int(b.basketBytes[ib])]
	f := b.tree.getFile()
	_, err = f.ReadAt(buf, b.basketSeek[ib])
	if err != nil {
		return err
	}

	sictx := b.tree.getFile()
	err = b.basket.UnmarshalROOT(NewRBuffer(buf, nil, 0, sictx))
	if err != nil {
		return err
	}
	b.basket.key.f = f
	b.firstEntry = entry

	if len(b.basketBuf) < int(b.basket.key.objlen) {
		b.basketBuf = make([]byte, b.basket.key.objlen)
	}
	buf = b.basketBuf[:int(b.basket.key.objlen)]
	_, err = b.basket.key.load(buf)
	if err != nil {
		return err
	}
	b.basket.rbuf = NewRBuffer(buf, nil, uint32(b.basket.key.keylen), sictx)

	for _, leaf := range b.leaves {
		err = leaf.readBasket(b.basket.rbuf)
		if err != nil {
			return err
		}
	}

	if b.entryOffsetLen > 0 {
		last := int64(b.basket.last)
		err = b.basket.rbuf.setPos(last)
		if err != nil {
			return err
		}
		n := b.basket.rbuf.ReadI32()
		b.basket.offsets = b.basket.rbuf.ReadFastArrayI32(int(n))
		if b.basket.rbuf.err != nil {
			return b.basket.rbuf.err
		}
	}
	return err
}

func (b *tbranch) scan(ptr interface{}) error {
	for _, leaf := range b.leaves {
		err := leaf.scan(b.basket.rbuf, ptr)
		if err != nil {
			return err
		}
	}
	return b.basket.rbuf.err
}

func (b *tbranch) setAddress(ptr interface{}) error {
	var err error
	return err
}

func (b *tbranch) setStreamer(s StreamerInfo, ctx StreamerInfoContext) {
	// no op
}

func (b *tbranch) setStreamerElement(s StreamerElement, ctx StreamerInfoContext) {
	// no op
}

// tbranchElement is a Branch for objects.
type tbranchElement struct {
	rvers int16
	tbranch
	class   string          // class name of referenced object
	parent  string          // name of parent class
	clones  string          // named of class in TClonesArray (if any)
	chksum  uint32          // checksum of class
	clsver  uint16          // version number of class
	id      int32           // element serial number in fInfo
	btype   int32           // branch type
	stype   int32           // branch streamer type
	max     int32           // maximum entries for a TClonesArray or variable array
	stltyp  int32           // STL container type
	bcount1 *tbranchElement // pointer to primary branchcount branch
	bcount2 *tbranchElement // pointer to secondary branchcount branch

	streamer  StreamerInfo
	estreamer StreamerElement
	scanfct   func(b *tbranchElement, ptr interface{}) error
}

func (b *tbranchElement) Class() string {
	return "TBranchElement"
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *tbranchElement) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	b.rvers = vers
	if vers < 8 {
		r.err = fmt.Errorf("rootio: TBranchElement version too old (%d < 8)", vers)
		return r.err
	}

	if err := b.tbranch.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	b.class = r.ReadString()
	b.parent = r.ReadString()
	b.clones = r.ReadString()
	b.chksum = r.ReadU32()
	if vers >= 10 {
		b.clsver = r.ReadU16()
	} else {
		b.clsver = uint16(r.ReadU32())
	}
	b.id = r.ReadI32()
	b.btype = r.ReadI32()
	b.stype = r.ReadI32()
	b.max = r.ReadI32()

	bcount1 := r.ReadObjectAny()
	if bcount1 != nil {
		b.bcount1 = bcount1.(*tbranchElement)
	}

	bcount2 := r.ReadObjectAny()
	if bcount2 != nil {
		b.bcount2 = bcount2.(*tbranchElement)
	}

	r.CheckByteCount(pos, bcnt, beg, "TBranchElement")
	return r.err
}

func (b *tbranchElement) loadEntry(ientry int64) error {
	if len(b.branches) > 0 {
		for _, sub := range b.branches {
			err := sub.loadEntry(ientry)
			if err != nil {
				return err
			}
		}
	}
	return b.tbranch.loadEntry(ientry)
}

func (b *tbranchElement) setAddress(ptr interface{}) error {
	var err error
	err = b.setupReadStreamer()
	if err != nil {
		return err
	}

	b.scanfct = func(b *tbranchElement, ptr interface{}) error {
		return b.tbranch.scan(ptr)
	}
	if len(b.branches) > 0 {
		var ids []int
		rv := reflect.ValueOf(ptr).Elem()
		for _, sub := range b.branches {
			i := int(sub.(*tbranchElement).id)
			ids = append(ids, i)
			fptr := rv.Field(i).Addr().Interface()
			err = sub.setAddress(fptr)
			if err != nil {
				return err
			}
		}
		b.scanfct = func(b *tbranchElement, ptr interface{}) error {
			rv := reflect.ValueOf(ptr).Elem()
			for i, sub := range b.branches {
				id := ids[i]
				fptr := rv.Field(id).Addr().Interface()
				err := sub.scan(fptr)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
	if b.id < 0 {
		for _, leaf := range b.tbranch.leaves {
			leaf, ok := leaf.(*tleafElement)
			if !ok {
				continue
			}
			err = leaf.setAddress(ptr)
			if err != nil {
				return err
			}
		}
	} else {
		elts := b.streamer.Elements()
		for _, leaf := range b.tbranch.leaves {
			leaf, ok := leaf.(*tleafElement)
			if !ok {
				continue
			}
			var elt StreamerElement
			leafName := leaf.Name()
			if strings.Contains(leafName, ".") {
				idx := strings.LastIndex(leafName, ".")
				leafName = string(leafName[idx+1:])
			}
			for _, ee := range elts {
				if ee.Name() == leafName {
					elt = ee
					break
				}
			}
			if elt == nil {
				return fmt.Errorf("rootio: failed to find StreamerElement for leaf %q", leaf.Name())
			}
			leaf.streamers = []StreamerElement{elt}
			err = leaf.setAddress(ptr)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (b *tbranchElement) scan(ptr interface{}) error {
	return b.scanfct(b, ptr)
}

func (b *tbranchElement) setupReadStreamer() error {
	streamer, ok := streamers.get(b.class, int(b.clsver), int(b.chksum))
	if !ok {
		return fmt.Errorf("rootio: no StreamerInfo for class=%q version=%d checksum=%d", b.class, b.clsver, b.chksum)
	}
	b.streamer = streamer

	for _, leaf := range b.tbranch.leaves {
		leaf, ok := leaf.(*tleafElement)
		if !ok {
			continue
		}
		leaf.streamers = b.streamer.Elements()
	}

	for _, sub := range b.branches {
		sub, ok := sub.(*tbranchElement)
		if !ok {
			continue
		}
		err := sub.setupReadStreamer()
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *tbranchElement) GoType() reflect.Type {
	return gotypeFromSI(b.streamer, b.tree.getFile())
}

func (b *tbranchElement) setStreamer(s StreamerInfo, ctx StreamerInfoContext) {
	b.streamer = s
	if len(b.tbranch.leaves) == 1 {
		tle := b.tbranch.leaves[0].(*tleafElement)
		tle.streamers = s.Elements()
		tle.src = reflect.New(gotypeFromSI(s, ctx)).Elem()
	}
	err := b.setupReadStreamer()
	if err != nil {
		panic(err)
	}
}

func (b *tbranchElement) setStreamerElement(se StreamerElement, ctx StreamerInfoContext) {
	b.estreamer = se
	if len(b.Leaves()) == 1 {
		tle := b.Leaves()[0].(*tleafElement)
		tle.streamers = []StreamerElement{se}
		tle.src = reflect.New(gotypeFromSE(se, tle.LeafCount(), ctx)).Elem()
	}
	err := b.setupReadStreamer()
	if err != nil {
		panic(err)
	}
}

func init() {
	{
		f := func() reflect.Value {
			o := &tbranch{}
			return reflect.ValueOf(o)
		}
		Factory.add("TBranch", f)
		Factory.add("*rootio.tbranch", f)
	}
	{
		f := func() reflect.Value {
			o := &tbranchElement{}
			return reflect.ValueOf(o)
		}
		Factory.add("TBranchElement", f)
		Factory.add("*rootio.tbranchElement", f)
	}
}

var _ Object = (*tbranch)(nil)
var _ Named = (*tbranch)(nil)
var _ Branch = (*tbranch)(nil)
var _ ROOTUnmarshaler = (*tbranch)(nil)

var _ Object = (*tbranchElement)(nil)
var _ Named = (*tbranchElement)(nil)
var _ Branch = (*tbranchElement)(nil)
var _ ROOTUnmarshaler = (*tbranchElement)(nil)
