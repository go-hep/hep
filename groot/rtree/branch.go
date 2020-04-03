// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

const (
	defaultBasketSize = 32 * 1024 // default basket size in bytes
	defaultSplitLevel = 99        // default split-level for branches
	defaultMaxBaskets = 10        // default number of baskets
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

func newBranchFromWVars(w *wtree, name string, wvars []WriteVar, parent Branch, cfg wopt) (*tbranch, error) {
	b := &tbranch{
		named:    *rbase.NewNamed(name, ""),
		attfill:  *rbase.NewAttFill(),
		compress: int(cfg.compress),

		iobits:      w.ttree.iobits,
		basketSize:  int(cfg.bufsize),
		maxBaskets:  defaultMaxBaskets,
		basketBytes: make([]int32, 0, defaultMaxBaskets),
		basketEntry: make([]int64, 1, defaultMaxBaskets),
		basketSeek:  make([]int64, 0, defaultMaxBaskets),

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
			et, shape := flattenArrayType(rt)
			for _, dim := range shape {
				fmt.Fprintf(title, "[%d]", dim)
			}
			rt = et
		case reflect.Slice:
			if wvar.Count == "" {
				return nil, fmt.Errorf("rtree: empty name for count-leaf of slice %q", wvar.Name)
			}
			fmt.Fprintf(title, "[%s]", wvar.Count)
			rt = rt.Elem()
			b.entryOffsetLen = 1000 // slice, so we need an offset array
		case reflect.String:
			b.entryOffsetLen = 1000 // string, so we need an offset array
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
	b.createNewBasket()

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
	switch {
	case b.bup != nil:
		return b.bup.Leaf(name)
	case b.btop != nil:
		return b.btop.Leaf(name)
	case b.tree != nil:
		return b.tree.Leaf(name)
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
		panic(fmt.Errorf("rtree: branch [%s] failed to load entry %d: %w", b.Name(), i, err))
	}
}

func (b *tbranch) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	maxBaskets := b.maxBaskets
	defer func() { b.maxBaskets = maxBaskets }()
	b.maxBaskets = b.writeBasket + 1
	if b.maxBaskets < defaultMaxBaskets {
		b.maxBaskets = defaultMaxBaskets
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

	{
		sli := b.basketBytes[:b.writeBasket]
		w.WriteI8(1)
		w.WriteFastArrayI32(sli)
		if n := b.maxBaskets - len(sli); n > 0 {
			// fill up with zeros.
			w.WriteFastArrayI32(make([]int32, n))
		}
	}

	{
		sli := b.basketEntry[:b.writeBasket+1]
		w.WriteI8(1)
		w.WriteFastArrayI64(sli)
		if n := b.maxBaskets - len(sli); n > 0 {
			// fill up with zeros.
			w.WriteFastArrayI64(make([]int64, n))
		}
	}

	{
		sli := b.basketSeek[:b.writeBasket]
		w.WriteI8(1)
		w.WriteFastArrayI64(sli)
		if n := b.maxBaskets - len(sli); n > 0 {
			// fill up with zeros.
			w.WriteFastArrayI64(make([]int64, n))
		}
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
	b.readentry = -1
	b.firstbasket = -1
	b.nextbasket = -1

	const minVers = 6
	switch {
	case vers >= 10:

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

	case vers >= 6:
		if err := b.named.UnmarshalROOT(r); err != nil {
			return err
		}

		if vers > 7 {
			if err := b.attfill.UnmarshalROOT(r); err != nil {
				return err
			}
		}

		b.compress = int(r.ReadI32())
		b.basketSize = int(r.ReadI32())
		b.entryOffsetLen = int(r.ReadI32())
		b.writeBasket = int(r.ReadI32())
		b.entryNumber = int64(r.ReadI32())
		b.offset = int(r.ReadI32())
		b.maxBaskets = int(r.ReadI32())
		if vers > 6 {
			b.splitLevel = int(r.ReadI32())
		}
		b.entries = int64(r.ReadF64())
		b.totBytes = int64(r.ReadF64())
		b.zipBytes = int64(r.ReadF64())

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
		b.basketBytes = r.ReadFastArrayI32(b.maxBaskets)

		/*isArray*/
		_ = r.ReadI8()
		{
			slice := r.ReadFastArrayI32(b.maxBaskets)
			b.basketEntry = make([]int64, len(slice))
			for i, v := range slice {
				b.basketEntry[i] = int64(v)
			}
		}

		switch r.ReadI8() {
		case 2:
			b.basketSeek = r.ReadFastArrayI64(b.maxBaskets)
		default:
			slice := r.ReadFastArrayI32(b.maxBaskets)
			b.basketSeek = make([]int64, len(slice))
			for i, v := range slice {
				b.basketSeek[i] = int64(v)
			}
		}

		b.fname = r.ReadString()

	default:
		panic(fmt.Errorf("rtree: too old TBranch version (%d<%d)", vers, minVers))
	}

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
		return fmt.Errorf("rtree: no basket for entry %d", entry)
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
	var (
		buf   = b.basketBuf[:int(b.basketBytes[ib])]
		seek  = b.basketSeek[ib]
		f     = b.tree.getFile()
		sictx = f
	)

	switch {
	case len(buf) == 0 && b.basket != nil: // FIXME(sbinet): from trial and error. check this is ok for all cases
		bk = b.basket
		bk.key.SetFile(f)
		b.firstEntry = b.basketEntry[ib]
		b.basketEntry[ib] = 0
		b.basketEntry[ib+1] = int64(bk.nevbuf)

		buf = make([]byte, int(bk.key.ObjLen()))
		_, err = bk.key.Load(buf)
		if err != nil {
			return err
		}

		bk.rbuf = rbytes.NewRBuffer(buf, nil, 0, sictx)
		for _, leaf := range b.leaves {
			err = leaf.readFromBuffer(bk.rbuf)
			if err != nil {
				return err
			}
		}

	default:
		_, err = f.ReadAt(buf, seek)
		if err != nil {
			return err
		}

		err = bk.UnmarshalROOT(rbytes.NewRBuffer(buf, nil, 0, sictx))
		if err != nil {
			return err
		}
		bk.key.SetFile(f)
		b.firstEntry = b.basketEntry[ib]

		buf = make([]byte, int(bk.key.ObjLen()))
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

	}

	return err
}

func (b *tbranch) scan(ptr interface{}) error {
	switch len(b.leaves) {
	case 1:
		leaf := b.leaves[0]
		err := leaf.scan(b.basket.rbuf, ptr)
		if err != nil {
			return err
		}
	default:
		rv := reflect.ValueOf(ptr).Elem()
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			var (
				leaf = b.leaves[i]
				fv   = rv.Field(i)
				ptr  = fv.Addr().Interface()
			)
			err := leaf.scan(b.basket.rbuf, ptr)
			if err != nil {
				return err
			}
		}
	}
	return b.basket.rbuf.Err()
}

func (b *tbranch) setAddress(ptr interface{}) error {
	switch len(b.leaves) {
	case 0:
		return fmt.Errorf("rtree: can not set address for a leaf-less branch (name=%q)", b.Name())

	case 1:
		err := b.leaves[0].setAddress(ptr)
		if err != nil {
			return fmt.Errorf("rtree: could not set address for leaf[%d][%s]: %w", 0, b.leaves[0].Name(), err)
		}

	default:
		rv := reflect.Indirect(reflect.ValueOf(ptr))
		rt := rv.Type()
		switch kind := rv.Kind(); kind {
		case reflect.Struct:
			if len(b.leaves) != rt.NumField() {
				// FIXME(sbinet): be more lenient and clever about this?
				return fmt.Errorf("rtree: fields/leaves number mismatch (name=%q, fields=%d, leaves=%d)", b.Name(), rt.NumField(), len(b.leaves))
			}
			for i, leaf := range b.leaves {
				fv := rv.Field(i)
				err := leaf.setAddress(fv.Addr().Interface())
				if err != nil {
					return fmt.Errorf("rtree: could not set address for leaf[%d][%s]: %w", i, leaf.Name(), err)
				}
			}
		default:
			// TODO(sbinet): also support map[string]*T ?
			// TODO(sbinet): also support []*T ?
			return fmt.Errorf("rtree: multi-leaf branches need a pointer-to-struct (got=%T)", ptr)
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

func (b *tbranch) createNewBasket() {
	b.writeBasket = len(b.baskets)
	cycle := int16(b.writeBasket)
	b.baskets = append(b.baskets, newBasketFrom(b.tree, b, cycle, b.basketSize, b.entryOffsetLen))
	b.basket = &b.baskets[b.writeBasket]
	if n := len(b.baskets); n > b.maxBaskets {
		b.maxBaskets = n
	}
}

func (b *tbranch) write() (int, error) {
	b.entries++
	b.entryNumber++

	szOld := b.basket.wbuf.Len()
	b.basket.update(szOld)
	_, err := b.writeToBuffer(b.basket.wbuf)
	szNew := b.basket.wbuf.Len()
	n := int(szNew - szOld)
	if err != nil {
		return n, fmt.Errorf("could not write to buffer (branch=%q): %w", b.Name(), err)
	}
	if n > b.basket.nevsize {
		b.basket.nevsize = n
	}

	// FIXME(sbinet): harmonize or drive via "auto-flush" ?
	if szNew+int64(n) >= int64(b.basketSize) {
		err = b.flush()
		if err != nil {
			return n, fmt.Errorf("could not flush branch (auto-flush): %w", err)
		}

		b.createNewBasket()
	}
	return n, nil
}

func (b *tbranch) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	var tot int
	for i, leaf := range b.leaves {
		n, err := leaf.writeToBuffer(w)
		if err != nil {
			return tot, fmt.Errorf("could not write leaf[%d] name=%q of branch %q: %w", i, leaf.Name(), b.Name(), err)
		}
		tot += n
	}
	return tot, nil
}

func (b *tbranch) flush() error {
	for i, sub := range b.branches {
		err := sub.flush()
		if err != nil {
			return fmt.Errorf("could not flush subbranch[%d]=%q of branch %q: %w", i, sub.Name(), b.Name(), err)
		}
	}

	f := b.tree.getFile()
	totBytes, zipBytes, err := b.basket.writeFile(f)
	if err != nil {
		return fmt.Errorf("could not marshal basket[%d] (branch=%q): %w", b.writeBasket, b.Name(), err)
	}
	b.totBytes += totBytes
	b.zipBytes += zipBytes

	b.basketBytes = append(b.basketBytes, b.basket.key.Nbytes())
	b.basketEntry = append(b.basketEntry, b.entryNumber)
	b.basketSeek = append(b.basketSeek, b.basket.key.SeekKey())
	b.writeBasket++
	b.basket = nil

	return nil
}

// tbranchElement is a Branch for objects.
type tbranchElement struct {
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

func (b *tbranchElement) RVersion() int16 {
	return rvers.BranchElement
}

func (b *tbranchElement) Class() string {
	return "TBranchElement"
}

func (b *tbranchElement) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(b.RVersion())
	if n, err := b.tbranch.MarshalROOT(w); err != nil {
		return n, err
	}
	w.WriteString(b.class)
	w.WriteString(b.parent)
	w.WriteString(b.clones)
	w.WriteU32(b.chksum)
	w.WriteU16(b.clsver)
	w.WriteI32(b.id)
	w.WriteI32(b.btype)
	w.WriteI32(b.stype)
	w.WriteI32(b.max)

	{
		var obj root.Object
		if b.bcount1 != nil {
			obj = b.bcount1
		}
		if err := w.WriteObjectAny(obj); err != nil {
			return int(w.Pos() - pos), err
		}
	}
	{
		var obj root.Object
		if b.bcount2 != nil {
			obj = b.bcount2
		}
		if err := w.WriteObjectAny(obj); err != nil {
			return int(w.Pos() - pos), err
		}
	}

	return w.SetByteCount(pos, b.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *tbranchElement) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(b.Class())
	if vers < 1 {
		r.SetErr(fmt.Errorf("rtree: TBranchElement version too old (%d < 8)", vers))
		return r.Err()
	}

	if err := b.tbranch.UnmarshalROOT(r); err != nil {
		return err
	}

	b.class = r.ReadString()
	if vers > 1 {
		b.parent = r.ReadString()
		b.clones = r.ReadString()
		b.chksum = r.ReadU32()
	}
	if vers >= 10 {
		b.clsver = r.ReadU16()
	} else {
		b.clsver = uint16(r.ReadU32())
	}
	b.id = r.ReadI32()
	b.btype = r.ReadI32()
	b.stype = r.ReadI32()
	if vers > 1 {
		b.max = r.ReadI32()

		bcount1 := r.ReadObjectAny()
		if bcount1 != nil {
			b.bcount1 = bcount1.(*tbranchElement)
		}

		bcount2 := r.ReadObjectAny()
		if bcount2 != nil {
			b.bcount2 = bcount2.(*tbranchElement)
		}
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
				return fmt.Errorf("rtree: failed to find StreamerElement for leaf %q", leaf.Name())
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
			return fmt.Errorf("rtree: no StreamerInfo for class=%q version=%d checksum=%d", b.class, b.clsver, b.chksum)
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
	const max = 1<<31 - 1
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
			panic(fmt.Errorf("rtree: unknown branch type %T", b))
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
