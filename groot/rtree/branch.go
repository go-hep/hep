// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

const (
	defaultBasketSize = 32000 // default basket size in bytes
	maxBaskets        = 10    // default number of baskets
)

type tbranch struct {
	named          rbase.Named
	attfill        rbase.AttFill
	compress       int         // compression level and algorithm
	basketSize     int         // initial size of Basket buffer
	entryOffsetLen int         // initial length of entryOffset table in the basket buffers
	writeBasket    int         // last basket number written
	entryNumber    int64       // current entry number (last one filled in this branch)
	iobits         tioFeatures // IO features for newly-created baskets
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

	tree Tree            // tree header
	btop Branch          // top-level parent branch in the tree
	bup  Branch          // parent branch
	dir  riofs.Directory // directory where this branch's buffers are stored
}

func newBranchFromWVars(w *wtree, name string, wvars []WriteVar, parent Branch, compress int) (*tbranch, error) {
	b := &tbranch{
		named:    *rbase.NewNamed(name, ""),
		attfill:  *rbase.NewAttFill(),
		compress: compress,

		basketSize:  defaultBasketSize,
		maxBaskets:  maxBaskets,
		basketBytes: make([]int32, maxBaskets),
		basketEntry: make([]int64, maxBaskets),
		basketSeek:  make([]int64, maxBaskets),

		tree: &w.ttree,
		btop: btopOf(parent),
		bup:  parent,
		dir:  w.dir,
	}

	title := new(strings.Builder)
	for i, wvar := range wvars {
		if i > 0 {
			title.WriteString(":")
		}
		title.WriteString(wvar.Name)
		rt := reflect.TypeOf(wvar.Value).Elem()
		switch k := rt.Kind(); k {
		case reflect.Array:
			fmt.Fprintf(title, "[%d]", rt.Len())
			rt = rt.Elem()
		case reflect.Slice:
			fmt.Fprintf(title, "[N]") // FIXME(sbinet): how to link everything together?
			rt = rt.Elem()
		}
		code := gotypeToROOTTypeCode(rt)
		fmt.Fprintf(title, "/%s", code)

		leaf, err := newLeafFromWVar(b, wvar)
		if err != nil {
			return nil, err
		}
		b.leaves = append(b.leaves, leaf)
		w.ttree.leaves = append(w.ttree.leaves, leaf)
	}

	b.named.SetTitle(title.String())
	return b, nil
}

func (b *tbranch) RVersion() int16 {
	return rvers.Branch
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
		panic(errors.Errorf("rtree: branch [%s] failed to load entry %d: %v", b.Name(), i, err))
	}
}

func (b *tbranch) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(b.RVersion())
	if n, err := b.named.MarshalROOT(w); err != nil {
		return n, err
	}
	if n, err := b.attfill.MarshalROOT(w); err != nil {
		return n, err
	}
	w.WriteI32(int32(b.compress))
	w.WriteI32(int32(b.basketSize))
	w.WriteI32(int32(b.entryOffsetLen))
	w.WriteI32(int32(b.writeBasket))
	w.WriteI64(b.entryNumber)

	if n, err := b.iobits.MarshalROOT(w); err != nil {
		return n, err
	}

	w.WriteI32(int32(b.offset))
	w.WriteI32(int32(b.maxBaskets))
	w.WriteI32(int32(b.splitLevel))
	w.WriteI64(b.entries)
	w.WriteI64(b.firstEntry)
	w.WriteI64(b.totBytes)
	w.WriteI64(b.zipBytes)

	{
		branches := rcont.NewObjArray()
		if len(b.branches) > 0 {
			elems := make([]root.Object, len(b.branches))
			for i, v := range b.branches {
				elems[i] = v
			}
			branches.SetElems(elems)
		}
		if n, err := branches.MarshalROOT(w); err != nil {
			return n, err
		}
	}
	{
		leaves := rcont.NewObjArray()
		if len(b.leaves) > 0 {
			elems := make([]root.Object, len(b.leaves))
			for i, v := range b.leaves {
				elems[i] = v
			}
			leaves.SetElems(elems)
		}
		if n, err := leaves.MarshalROOT(w); err != nil {
			return n, err
		}
	}
	{
		baskets := rcont.NewObjArray()
		if len(b.baskets) > 0 {
			elems := make([]root.Object, len(b.baskets))
			for i := range b.baskets {
				elems[i] = &b.baskets[i]
			}
			baskets.SetElems(elems)
		}
		if n, err := baskets.MarshalROOT(w); err != nil {
			return n, err
		}
		baskets.SetElems(nil)
	}

	w.WriteI8(0)
	w.WriteFastArrayI32(b.basketBytes)
	if len(b.basketBytes) < b.maxBaskets {
		// fill up with zeros.
		w.WriteFastArrayI32(make([]int32, b.maxBaskets-len(b.basketBytes)))
	}

	w.WriteI8(0)
	w.WriteFastArrayI64(b.basketEntry)
	if len(b.basketEntry) < b.maxBaskets {
		// fill up with zeros.
		w.WriteFastArrayI64(make([]int64, b.maxBaskets-len(b.basketEntry)))
	}

	w.WriteI8(0)
	w.WriteFastArrayI64(b.basketSeek[:b.writeBasket])
	if len(b.basketSeek) < b.maxBaskets {
		// fill up with zeros.
		w.WriteFastArrayI64(make([]int64, b.maxBaskets-len(b.basketSeek)))
	}

	w.WriteString(b.fname)

	return w.SetByteCount(pos, b.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *tbranch) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(b.Class())

	b.tree = nil
	b.basket = nil
	b.firstbasket = -1
	b.nextbasket = -1

	if vers < 10 {
		panic(errors.Errorf("rtree: too old TBanch version (%d<10)", vers))
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
	if vers >= 13 {
		if err := b.iobits.UnmarshalROOT(r); err != nil {
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
		var branches rcont.ObjArray
		if err := branches.UnmarshalROOT(r); err != nil {
			return err
		}
		b.branches = make([]Branch, branches.Last()+1)
		for i := range b.branches {
			br := branches.At(i).(Branch)
			b.branches[i] = br
		}
	}

	{
		var leaves rcont.ObjArray
		if err := leaves.UnmarshalROOT(r); err != nil {
			return err
		}
		b.leaves = make([]Leaf, leaves.Last()+1)
		for i := range b.leaves {
			leaf := leaves.At(i).(Leaf)
			leaf.setBranch(b)
			b.leaves[i] = leaf
		}
	}
	{
		var baskets rcont.ObjArray
		if err := baskets.UnmarshalROOT(r); err != nil {
			return err
		}
		b.baskets = make([]Basket, 0, baskets.Last()+1)
		for i := 0; i < baskets.Last()+1; i++ {
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
	b.basketEntry = r.ReadFastArrayI64(b.maxBaskets)[: b.writeBasket+1 : b.writeBasket+1]

	/*isArray*/
	_ = r.ReadI8()
	b.basketSeek = r.ReadFastArrayI64(b.maxBaskets)[:b.writeBasket:b.writeBasket]

	b.fname = r.ReadString()

	r.CheckByteCount(pos, bcnt, beg, b.Class())

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
		return errors.Errorf("rtree: no basket for entry %d", entry)
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
	err = bk.UnmarshalROOT(rbytes.NewRBuffer(buf, nil, 0, sictx))
	if err != nil {
		return err
	}
	bk.key.SetFile(f)
	b.firstEntry = b.basketEntry[ib]

	buf = make([]byte, int(b.basket.key.ObjLen()))
	_, err = bk.key.Load(buf)
	if err != nil {
		return err
	}
	bk.rbuf = rbytes.NewRBuffer(buf, nil, uint32(bk.key.KeyLen()), sictx)

	for _, leaf := range b.leaves {
		err = leaf.readFromBuffer(bk.rbuf)
		if err != nil {
			return err
		}
	}

	if b.entryOffsetLen > 0 {
		last := int64(b.basket.last)
		err = bk.rbuf.SetPos(last)
		if err != nil {
			return err
		}
		n := bk.rbuf.ReadI32()
		bk.offsets = bk.rbuf.ReadFastArrayI32(int(n))
		if err := bk.rbuf.Err(); err != nil {
			return err
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
	return b.basket.rbuf.Err()
}

func (b *tbranch) setAddress(ptr interface{}) error {
	for i, leaf := range b.leaves {
		// FIXME(sbinet): adjust ptr for e.g. a "f1/F;f2/I;f3/i" branch
		err := leaf.setAddress(ptr)
		if err != nil {
			return errors.Wrapf(err, "rtree: could not set address for leaf[%d][%s]", i, leaf.Name())
		}
	}
	return nil
}

func (b *tbranch) setStreamer(s rbytes.StreamerInfo, ctx rbytes.StreamerInfoContext) {
	// no op
}

func (b *tbranch) setStreamerElement(s rbytes.StreamerElement, ctx rbytes.StreamerInfoContext) {
	// no op
}

func (b *tbranch) write() (int, error) {
	basket := b.basket
	if basket == nil {
		b.writeBasket = len(b.baskets)
		b.baskets = append(b.baskets, newBasketFrom(b.tree, b))
		basket = &b.baskets[b.writeBasket]
	}

	wbuf := basket.wbuf
	b.entries++
	b.entryNumber++

	n, err := b.writeToBuffer(wbuf)
	if err != nil {
		return n, errors.Wrapf(err, "could not write to buffer (branch=%q)", b.Name())
	}

	return n, nil
}

func (b *tbranch) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	var tot int
	for i, leaf := range b.leaves {
		n, err := leaf.writeToBuffer(w)
		if err != nil {
			return tot, errors.Wrapf(err, "could not write leaf[%d] name=%q of branch %q", i, leaf.Name(), b.Name())
		}
		tot += n
	}
	return tot, nil
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

	streamer  rbytes.StreamerInfo
	estreamer rbytes.StreamerElement
	scanfct   func(b *tbranchElement, ptr interface{}) error
}

func (b *tbranchElement) Class() string {
	return "TBranchElement"
}

func (b *tbranchElement) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}
	panic("not implemented")
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *tbranchElement) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(b.Class())
	b.rvers = vers
	if vers < 8 {
		r.SetErr(errors.Errorf("rtree: TBranchElement version too old (%d < 8)", vers))
		return r.Err()
	}

	if err := b.tbranch.UnmarshalROOT(r); err != nil {
		return err
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

	r.CheckByteCount(pos, bcnt, beg, b.Class())
	return r.Err()
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
	var sictx rbytes.StreamerInfoContext = b.getTree().getFile()
	var err error
	err = b.setupReadStreamer(sictx)
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
			var elt rbytes.StreamerElement
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
				return errors.Errorf("rtree: failed to find StreamerElement for leaf %q", leaf.Name())
			}
			leaf.streamers = []rbytes.StreamerElement{elt}
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

func (b *tbranchElement) setupReadStreamer(sictx rbytes.StreamerInfoContext) error {
	streamer, err := sictx.StreamerInfo(b.class, int(b.clsver))
	if err != nil {
		streamer, err = sictx.StreamerInfo(b.class, -1)
		if err != nil {
			return errors.Errorf("rtree: no StreamerInfo for class=%q version=%d checksum=%d", b.class, b.clsver, b.chksum)
		}
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
		err := sub.setupReadStreamer(sictx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *tbranchElement) GoType() reflect.Type {
	return gotypeFromSI(b.streamer, b.tree.getFile())
}

func (b *tbranchElement) setStreamer(s rbytes.StreamerInfo, ctx rbytes.StreamerInfoContext) {
	b.streamer = s
	if len(b.tbranch.leaves) == 1 {
		tle := b.tbranch.leaves[0].(*tleafElement)
		tle.streamers = s.Elements()
		tle.src = reflect.New(gotypeFromSI(s, ctx)).Elem()
	}
	err := b.setupReadStreamer(ctx)
	if err != nil {
		panic(err)
	}
}

func (b *tbranchElement) setStreamerElement(se rbytes.StreamerElement, ctx rbytes.StreamerInfoContext) {
	b.estreamer = se
	if len(b.Leaves()) == 1 {
		tle := b.Leaves()[0].(*tleafElement)
		tle.streamers = []rbytes.StreamerElement{se}
		tle.src = reflect.New(gotypeFromSE(se, tle.LeafCount(), ctx)).Elem()
	}
	err := b.setupReadStreamer(ctx)
	if err != nil {
		panic(err)
	}
}

func (b *tbranchElement) write() (int, error)                          { panic("not implemented") }
func (b *tbranchElement) writeToBuffer(w *rbytes.WBuffer) (int, error) { panic("not implemented") }

func btopOf(b Branch) Branch {
	if b == nil {
		return nil
	}
	const max = 1 << 32
	for i := 0; i < max; i++ {
		switch bb := b.(type) {
		case *tbranch:
			if bb.bup == nil {
				return bb
			}
			b = bb.bup
		case *tbranchElement:
			if bb.bup == nil {
				return bb
			}
			b = bb.bup
		default:
			panic(errors.Errorf("rtree: unknown branch type %T", b))
		}
	}
	panic("impossible")
}

func init() {
	{
		f := func() reflect.Value {
			o := &tbranch{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TBranch", f)
	}
	{
		f := func() reflect.Value {
			o := &tbranchElement{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TBranchElement", f)
	}
}

var (
	_ root.Object        = (*tbranch)(nil)
	_ root.Named         = (*tbranch)(nil)
	_ Branch             = (*tbranch)(nil)
	_ rbytes.Marshaler   = (*tbranch)(nil)
	_ rbytes.Unmarshaler = (*tbranch)(nil)

	_ root.Object        = (*tbranchElement)(nil)
	_ root.Named         = (*tbranchElement)(nil)
	_ Branch             = (*tbranchElement)(nil)
	_ rbytes.Marshaler   = (*tbranchElement)(nil)
	_ rbytes.Unmarshaler = (*tbranchElement)(nil)
)
