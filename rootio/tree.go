// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
	"strings"
)

// A ttree object is a list of Branch.
//   To Create a TTree object one must:
//    - Create the TTree header via the TTree constructor
//    - Call the TBranch constructor for every branch.
//
//   To Fill this object, use member function Fill with no parameters
//     The Fill function loops on all defined TBranch
type ttree struct {
	f *File // underlying file

	rvers int16
	named tnamed

	entries  int64 // Number of entries
	totbytes int64 // Total number of bytes in all branches before compression
	zipbytes int64 // Total number of bytes in all branches after  compression

	iofeats tioFeatures // IO features to define for newly-written baskets and branches

	branches []Branch // list of branches
	leaves   []Leaf   // direct pointers to individual branch leaves
}

func (tree *ttree) Class() string {
	return "TTree"
}

func (tree *ttree) Name() string {
	return tree.named.Name()
}

func (tree *ttree) Title() string {
	return tree.named.Title()
}

func (tree *ttree) Entries() int64 {
	return tree.entries
}

func (tree *ttree) TotBytes() int64 {
	return tree.totbytes
}

func (tree *ttree) ZipBytes() int64 {
	return tree.zipbytes
}

func (tree *ttree) Branches() []Branch {
	return tree.branches
}

func (tree *ttree) Branch(name string) Branch {
	for _, br := range tree.branches {
		if br.Name() == name {
			return br
		}
		for _, b1 := range br.Branches() {
			if b1.Name() == name {
				return b1
			}

			for _, b2 := range b1.Branches() {
				if b2.Name() == name {
					return b2
				}
			}
		}
	}

	// search using leaves.
	for _, leaf := range tree.leaves {
		b := leaf.Branch()
		if b.Name() == name {
			return b
		}
	}

	return nil
}

func (tree *ttree) Leaves() []Leaf {
	return tree.leaves
}

func (tree *ttree) Leaf(name string) Leaf {
	for _, leaf := range tree.leaves {
		if leaf.Name() == name {
			return leaf
		}
	}
	return nil
}

func (tree *ttree) SetFile(f *File) {
	tree.f = f
}

func (tree *ttree) getFile() *File {
	return tree.f
}

func (tree *ttree) loadEntry(entry int64) error {
	for _, b := range tree.branches {
		err := b.loadEntry(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (tree *ttree) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	tree.rvers = vers

	for _, a := range []ROOTUnmarshaler{
		&tree.named,
		&attline{},
		&attfill{},
		&attmarker{},
	} {
		err := a.UnmarshalROOT(r)
		if err != nil {
			return err
		}
	}

	if vers < 16 {
		return fmt.Errorf(
			"rootio.Tree: tree [%s] with version [%v] is not supported (too old)",
			tree.Name(),
			vers,
		)
	}

	tree.entries = r.ReadI64()
	tree.totbytes = r.ReadI64()
	tree.zipbytes = r.ReadI64()
	if vers >= 16 {
		_ = r.ReadI64() // fSavedBytes
	}
	if vers >= 18 {
		_ = r.ReadI64() // flushed bytes
	}

	_ = r.ReadF64() // fWeight
	_ = r.ReadI32() // fTimerInterval
	_ = r.ReadI32() // fScanField
	_ = r.ReadI32() // fUpdate

	if vers >= 17 {
		_ = r.ReadI32() // fDefaultEntryOffsetLen
	}
	nclus := 0
	if vers >= 19 { // FIXME
		nclus = int(r.ReadI32()) // fNClusterRange
	}

	_ = r.ReadI64() // fMaxEntries
	_ = r.ReadI64() // fMaxEntryLoop
	_ = r.ReadI64() // fMaxVirtualSize
	_ = r.ReadI64() // fAutoSave

	if vers >= 18 {
		_ = r.ReadI64() // fAutoFlush
	}

	_ = r.ReadI64() // fEstimate

	if vers >= 19 { // FIXME
		_ = r.ReadI8()
		_ = r.ReadFastArrayI64(nclus) // fClusterRangeEnd
		_ = r.ReadI8()
		_ = r.ReadFastArrayI64(nclus) // fClusterSize
	}

	if vers >= 20 {
		if err := tree.iofeats.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	var branches tobjarray
	if err := branches.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.branches = make([]Branch, branches.last+1)
	for i := range tree.branches {
		tree.branches[i] = branches.At(i).(Branch)
		tree.branches[i].setTree(tree)
	}

	var leaves tobjarray
	if err := leaves.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.leaves = make([]Leaf, leaves.last+1)
	for i := range tree.leaves {
		tree.leaves[i] = leaves.At(i).(Leaf)
		// FIXME(sbinet)
		//tree.leaves[i].SetBranch(tree.branches[i])
	}

	for range []string{
		"fAliases", "fIndexValues", "fIndex", "fTreeIndex", "fFriends",
		"fUserInfo", "fBranchRef",
	} {
		_ = r.ReadObjectAny()
	}

	r.CheckByteCount(pos, bcnt, beg, "TTree")

	// attach streamers to branches
	for i := range tree.branches {
		br := tree.branches[i]
		bre, ok := br.(*tbranchElement)
		if !ok {
			continue
		}
		cls := bre.class
		si, err := r.StreamerInfo(cls)
		if err != nil {
			panic(fmt.Errorf("rootio: could not find streamer for branch %q: %v", br.Name(), err))
		}
		// tree.attachStreamer(br, rstreamer, rstreamerCtx)
		tree.attachStreamer(br, si, r)
	}

	return r.Err()
}

func (tree *ttree) attachStreamer(br Branch, info StreamerInfo, ctx StreamerInfoContext) {
	if info == nil {
		return
	}

	if len(info.Elements()) == 1 {
		switch elem := info.Elements()[0].(type) {
		case *tstreamerBase:
			if elem.Name() == "TObjArray" {
				switch info.Name() {
				case "TClonesArray":
					cls := ""
					if bre, ok := br.(*tbranchElement); ok {
						cls = bre.clones
					}
					si, err := ctx.StreamerInfo(cls)
					if err != nil {
						panic(err)
					}
					tree.attachStreamer(br, si, ctx)
					return
				default:
					// FIXME(sbinet): can only determine streamer by reading some value?
					return
				}
			}
		case *tstreamerSTL:
			if elem.Name() == "This" {
				tree.attachStreamerElement(br, elem, ctx)
				return
			}
		}
	}

	br.setStreamer(info, ctx)

	for _, sub := range br.Branches() {
		name := sub.Name()
		if strings.HasPrefix(name, br.Name()+".") {
			name = name[len(br.Name())+1:]
		}

		if strings.Contains(name, "[") {
			idx := strings.Index(name, "[")
			name = name[:idx]
		}
		var se StreamerElement
		for _, elmt := range info.Elements() {
			if elmt.Name() == name {
				se = elmt
				break
			}
		}
		tree.attachStreamerElement(sub, se, ctx)
	}
}

func (tree *ttree) attachStreamerElement(br Branch, se StreamerElement, ctx StreamerInfoContext) {
	if se == nil {
		return
	}

	br.setStreamerElement(se, ctx)
	var members []StreamerElement
	switch se.(type) {
	case *tstreamerObject, *tstreamerObjectAny, *tstreamerObjectPointer, *tstreamerObjectAnyPointer:
		typename := strings.TrimRight(se.TypeName(), "*")
		info, err := ctx.StreamerInfo(typename)
		if err != nil {
			panic(err)
		}
		members = info.Elements()
	case *tstreamerSTL:
		typename := se.TypeName()
		// FIXME(sbinet): this string manipulation only works for one-parameter templates
		if strings.Contains(typename, "<") {
			typename = typename[strings.Index(typename, "<")+1 : strings.LastIndex(typename, ">")]
			typename = strings.TrimRight(typename, "*")
		}
		info, err := ctx.StreamerInfo(typename)
		if err != nil {
			if _, ok := cxxbuiltins[typename]; !ok {
				panic(err)
			}
		}
		if err == nil {
			members = info.Elements()
		}
	}

	if members == nil {
		return
	}

	for _, sub := range br.Branches() {
		name := sub.Name()
		if strings.HasPrefix(name, br.Name()+".") { // drop parent branch's name
			name = name[len(br.Name())+1:]
		}
		submembers := members
		for strings.Contains(name, ".") { // drop nested struct names, one at a time
			dot := strings.Index(name, ".")
			base := name[:dot]
			name = name[dot+1:]
			for _, subse := range submembers {
				if subse.Name() == base {
					switch subse.(type) {
					case *tstreamerObject, *tstreamerObjectAny, *tstreamerObjectPointer, *tstreamerObjectAnyPointer:
						subinfo, err := ctx.StreamerInfo(strings.TrimRight(subse.TypeName(), "*"))
						if err != nil {
							panic(err)
						}
						submembers = subinfo.Elements()
					}
				}
			}
		}

		if strings.Contains(name, "[") {
			idx := strings.Index(name, "[")
			name = name[:idx]
		}
		var subse StreamerElement
		for _, elmt := range members {
			if elmt.Name() == name {
				subse = elmt
				break
			}
		}
		tree.attachStreamerElement(sub, subse, ctx)
	}
}

type tntuple struct {
	ttree
	nvars int
}

func (nt *tntuple) Class() string {
	return "TNtuple"
}

func (nt *tntuple) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := nt.ttree.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	nt.nvars = int(r.ReadI32())

	r.CheckByteCount(pos, bcnt, beg, "TNtuple")
	return r.err
}

type tioFeatures uint8

func (tio *tioFeatures) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	_ /*vers*/, pos, bcnt := r.ReadVersion()

	var buf [4]byte // FIXME(sbinet) where do these 4 bytes come from ?
	r.read(buf[:])

	*tio = tioFeatures(r.ReadU8())

	r.CheckByteCount(pos, bcnt, beg, "TIOFeatures")
	return r.err
}

func init() {
	{
		f := func() reflect.Value {
			o := &ttree{}
			return reflect.ValueOf(o)
		}
		Factory.add("TTree", f)
		Factory.add("*rootio.ttree", f)
	}
	{
		f := func() reflect.Value {
			o := &tntuple{}
			return reflect.ValueOf(o)
		}
		Factory.add("TNtuple", f)
		Factory.add("*rootio.tntuple", f)
	}
}

var _ Object = (*ttree)(nil)
var _ Named = (*ttree)(nil)
var _ Tree = (*ttree)(nil)
var _ ROOTUnmarshaler = (*ttree)(nil)

var _ Object = (*tntuple)(nil)
var _ Named = (*tntuple)(nil)
var _ Tree = (*tntuple)(nil)
var _ ROOTUnmarshaler = (*tntuple)(nil)
