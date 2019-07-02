// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

// A ttree object is a list of Branch.
//   To Create a TTree object one must:
//    - Create the TTree header via the TTree constructor
//    - Call the TBranch constructor for every branch.
//
//   To Fill this object, use member function Fill with no parameters
//     The Fill function loops on all defined TBranch
type ttree struct {
	f *riofs.File // underlying file

	rvers int16
	named rbase.Named

	entries  int64 // Number of entries
	totbytes int64 // Total number of bytes in all branches before compression
	zipbytes int64 // Total number of bytes in all branches after  compression

	iofeats tioFeatures // IO features to define for newly-written baskets and branches

	clusters clusters

	branches []Branch // list of branches
	leaves   []Leaf   // direct pointers to individual branch leaves
}

type clusters struct {
	ranges []int64 // last entry to a cluster range
	sizes  []int64 // number of entries in each cluster for a given range
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

func (tree *ttree) SetFile(f *riofs.File) {
	tree.f = f
}

func (tree *ttree) getFile() *riofs.File {
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
func (tree *ttree) UnmarshalROOT(r *rbytes.RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(tree.Class())
	tree.rvers = vers

	for _, a := range []rbytes.Unmarshaler{
		&tree.named,
		&rbase.AttLine{},
		&rbase.AttFill{},
		&rbase.AttMarker{},
	} {
		err := a.UnmarshalROOT(r)
		if err != nil {
			return err
		}
	}

	if vers < 16 {
		return errors.Errorf(
			"rtree: tree [%s] with version [%v] is not supported (too old)",
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
		tree.clusters.ranges = r.ReadFastArrayI64(nclus) // fClusterRangeEnd
		_ = r.ReadI8()
		tree.clusters.sizes = r.ReadFastArrayI64(nclus) // fClusterSize
	}

	if vers >= 20 {
		if err := tree.iofeats.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	var branches rcont.ObjArray
	if err := branches.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.branches = make([]Branch, branches.Last()+1)
	for i := range tree.branches {
		tree.branches[i] = branches.At(i).(Branch)
		tree.branches[i].setTree(tree)
	}

	var leaves rcont.ObjArray
	if err := leaves.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.leaves = make([]Leaf, leaves.Last()+1)
	for i := range tree.leaves {
		leaf := leaves.At(i).(Leaf)
		tree.leaves[i] = leaf
		// FIXME(sbinet)
		//tree.leaves[i].SetBranch(tree.branches[i])
	}

	for range []string{
		"fAliases", "fIndexValues", "fIndex", "fTreeIndex", "fFriends",
		"fUserInfo", "fBranchRef",
	} {
		_ = r.ReadObjectAny()
	}

	r.CheckByteCount(pos, bcnt, beg, tree.Class())

	// attach streamers to branches
	for i := range tree.branches {
		br := tree.branches[i]
		bre, ok := br.(*tbranchElement)
		if !ok {
			continue
		}
		cls := bre.class
		si, err := r.StreamerInfo(cls, int(bre.clsver))
		if err != nil {
			panic(errors.Errorf("rtree: could not find streamer for branch %q: %v", br.Name(), err))
		}
		// tree.attachStreamer(br, rstreamer, rstreamerCtx)
		tree.attachStreamer(br, si, r)
	}

	return r.Err()
}

func (tree *ttree) attachStreamer(br Branch, info rbytes.StreamerInfo, ctx rbytes.StreamerInfoContext) {
	if info == nil {
		return
	}

	if len(info.Elements()) == 1 {
		switch elem := info.Elements()[0].(type) {
		case *rdict.StreamerBase:
			if elem.Name() == "TObjArray" {
				switch info.Name() {
				case "TClonesArray":
					cls := ""
					version := -1
					if bre, ok := br.(*tbranchElement); ok {
						cls = bre.clones
						version = int(bre.clsver)
					}
					si, err := ctx.StreamerInfo(cls, version)
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
		case *rdict.StreamerSTL:
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
		var se rbytes.StreamerElement
		for _, elmt := range info.Elements() {
			if elmt.Name() == name {
				se = elmt
				break
			}
		}
		tree.attachStreamerElement(sub, se, ctx)
	}
}

func (tree *ttree) attachStreamerElement(br Branch, se rbytes.StreamerElement, ctx rbytes.StreamerInfoContext) {
	if se == nil {
		return
	}

	br.setStreamerElement(se, ctx)
	var members []rbytes.StreamerElement
	switch se.(type) {
	case *rdict.StreamerObject, *rdict.StreamerObjectAny, *rdict.StreamerObjectPointer, *rdict.StreamerObjectAnyPointer:
		typename := strings.TrimRight(se.TypeName(), "*")
		typevers := -1
		// FIXME(sbinet): always load latest version?
		info, err := ctx.StreamerInfo(typename, typevers)
		if err != nil {
			panic(err)
		}
		members = info.Elements()
	case *rdict.StreamerSTL:
		typename := strings.TrimSpace(se.TypeName())
		// FIXME(sbinet): this string manipulation only works for one-parameter templates
		if strings.Contains(typename, "<") {
			typename = typename[strings.Index(typename, "<")+1 : strings.LastIndex(typename, ">")]
			typename = strings.TrimRight(typename, "*")
		}
		typename = strings.TrimSpace(se.TypeName())
		typevers := -1
		// FIXME(sbinet): always load latest version?
		info, err := ctx.StreamerInfo(typename, typevers)
		if err != nil {
			if _, ok := rmeta.CxxBuiltins[typename]; !ok {
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
					case *rdict.StreamerObject, *rdict.StreamerObjectAny, *rdict.StreamerObjectPointer, *rdict.StreamerObjectAnyPointer:
						// FIXME(sbinet): always load latest version?
						subinfo, err := ctx.StreamerInfo(strings.TrimRight(subse.TypeName(), "*"), -1)
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
		var subse rbytes.StreamerElement
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

func (*tntuple) Class() string {
	return "TNtuple"
}

func (nt *tntuple) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion(nt.Class())

	if err := nt.ttree.UnmarshalROOT(r); err != nil {
		return err
	}

	nt.nvars = int(r.ReadI32())

	r.CheckByteCount(pos, bcnt, beg, nt.Class())
	return r.Err()
}

type tioFeatures uint8

func (tio *tioFeatures) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	_ /*vers*/, pos, bcnt := r.ReadVersion("TIOFeatures")

	var buf [4]byte // FIXME(sbinet) where do these 4 bytes come from ?
	r.Read(buf[:])

	*tio = tioFeatures(r.ReadU8())

	r.CheckByteCount(pos, bcnt, beg, "TIOFeatures")
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := &ttree{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TTree", f)
	}
	{
		f := func() reflect.Value {
			o := &tntuple{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TNtuple", f)
	}
}

var (
	_ root.Object        = (*ttree)(nil)
	_ root.Named         = (*ttree)(nil)
	_ Tree               = (*ttree)(nil)
	_ rbytes.Unmarshaler = (*ttree)(nil)

	_ root.Object        = (*tntuple)(nil)
	_ root.Named         = (*tntuple)(nil)
	_ Tree               = (*tntuple)(nil)
	_ rbytes.Unmarshaler = (*tntuple)(nil)
)
