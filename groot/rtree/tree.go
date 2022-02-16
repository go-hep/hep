// Copyright ©2017 The go-hep Authors. All rights reserved.
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
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// A ttree object is a list of Branch.
//   To Create a TTree object one must:
//    - Create the TTree header via the TTree constructor
//    - Call the TBranch constructor for every branch.
//
//   To Fill this object, use member function Fill with no parameters
//     The Fill function loops on all defined TBranch
type ttree struct {
	f   *riofs.File     // underlying file
	dir riofs.Directory // directory holding this tree

	rvers     int16
	named     rbase.Named
	attline   rbase.AttLine
	attfill   rbase.AttFill
	attmarker rbase.AttMarker

	entries       int64   // Number of entries
	totBytes      int64   // Total number of bytes in all branches before compression
	zipBytes      int64   // Total number of bytes in all branches after  compression
	savedBytes    int64   // number of autosaved bytes
	flushedBytes  int64   // number of auto-flushed bytes
	weight        float64 // tree weight
	timerInterval int32   // timer interval in milliseconds
	scanField     int32   // number of runs before prompting in Scan
	update        int32   // update frequency for entry-loop

	defaultEntryOffsetLen int32 // initial length of the entry offset table in the basket buffers
	maxEntries            int64 // maximum number of entries in case of circular buffers
	maxEntryLoop          int64 // maximum number of entries to process
	maxVirtualSize        int64 // maximum total size of buffers kept in memory
	autoSave              int64 // autosave tree when autoSave entries written
	autoFlush             int64 // autoflush tree when autoFlush entries written
	estimate              int64 // number of entries to estimate histogram limits

	clusters clusters

	iobits tioFeatures // IO features to define for newly-written baskets and branches

	branches []Branch // list of branches
	leaves   []Leaf   // direct pointers to individual branch leaves

	aliases     *rcont.List   // list of aliases for expressions based on the tree branches
	indexValues *rcont.ArrayD // sorted index values
	index       *rcont.ArrayI // index of sorted values
	treeIndex   root.Object   // pointer to the tree index (if any) // FIXME(sbinet): impl TVirtualIndex?
	friends     *rcont.List   // pointer to the list of firend elements
	userInfo    *rcont.List   // pointer to a list of user objects associated with this tree
	branchRef   root.Object   // branch supporting the reftable (if any) // FIXME(sbinet): impl TBranchRef?
}

type clusters struct {
	ranges []int64 // last entry to a cluster range
	sizes  []int64 // number of entries in each cluster for a given range
}

func (*ttree) RVersion() int16 {
	return rvers.Tree
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
	return tree.totBytes
}

func (tree *ttree) ZipBytes() int64 {
	return tree.zipBytes
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

func (tree *ttree) SetFile(f *riofs.File) { tree.f = f }
func (tree *ttree) getFile() *riofs.File  { return tree.f }

func (tree *ttree) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(tree.Class(), tree.RVersion())
	w.WriteObject(&tree.named)
	w.WriteObject(&tree.attline)
	w.WriteObject(&tree.attfill)
	w.WriteObject(&tree.attmarker)

	w.WriteI64(tree.entries)
	w.WriteI64(tree.totBytes)
	w.WriteI64(tree.zipBytes)
	w.WriteI64(tree.savedBytes)
	w.WriteI64(tree.flushedBytes)
	w.WriteF64(tree.weight)
	w.WriteI32(tree.timerInterval)
	w.WriteI32(tree.scanField)
	w.WriteI32(tree.update)
	w.WriteI32(tree.defaultEntryOffsetLen)
	w.WriteI32(int32(len(tree.clusters.ranges)))

	w.WriteI64(tree.maxEntries)
	w.WriteI64(tree.maxEntryLoop)
	w.WriteI64(tree.maxVirtualSize)
	w.WriteI64(tree.autoSave)

	w.WriteI64(tree.autoFlush)
	w.WriteI64(tree.estimate)

	w.WriteI8(0)
	w.WriteArrayI64(tree.clusters.ranges)
	w.WriteI8(0)
	w.WriteArrayI64(tree.clusters.sizes)

	w.WriteObject(&tree.iobits)

	{
		branches := rcont.NewObjArray()
		if len(tree.branches) > 0 {
			elems := make([]root.Object, len(tree.branches))
			for i, v := range tree.branches {
				elems[i] = v
			}
			branches.SetElems(elems)
		}
		w.WriteObject(branches)
	}
	{
		leaves := rcont.NewObjArray()
		if len(tree.leaves) > 0 {
			elems := make([]root.Object, len(tree.leaves))
			for i, v := range tree.leaves {
				elems[i] = v
			}
			leaves.SetElems(elems)
		}
		w.WriteObject(leaves)
	}

	{
		var obj root.Object
		if tree.aliases != nil {
			obj = tree.aliases
		}
		w.WriteObjectAny(obj)
	}

	{
		var obj root.Object
		if tree.indexValues != nil {
			obj = tree.indexValues
		}
		w.WriteObjectAny(obj)
	}

	{
		var obj root.Object
		if tree.index != nil {
			obj = tree.index
		}
		w.WriteObjectAny(obj)
	}

	{
		var obj root.Object
		if tree.treeIndex != nil {
			obj = tree.treeIndex
		}
		w.WriteObjectAny(obj)
	}

	{
		var obj root.Object
		if tree.friends != nil {
			obj = tree.friends
		}
		w.WriteObjectAny(obj)
	}

	{
		var obj root.Object
		if tree.userInfo != nil {
			obj = tree.userInfo
		}
		w.WriteObjectAny(obj)
	}

	{
		var obj root.Object
		if tree.branchRef != nil {
			obj = tree.branchRef
		}
		w.WriteObjectAny(obj)
	}

	return w.SetHeader(hdr)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (tree *ttree) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tree.Class())
	if hdr.Vers > rvers.Tree {
		panic(fmt.Errorf("rtree: invalid TTree version=%d > %d", hdr.Vers, rvers.Tree))
	}

	tree.rvers = hdr.Vers

	r.ReadObject(&tree.named)
	r.ReadObject(&tree.attline)
	r.ReadObject(&tree.attfill)
	r.ReadObject(&tree.attmarker)

	switch {
	default:
		panic(fmt.Errorf(
			"rtree: tree [%s] with version [%v] is not supported (too old)",
			tree.Name(),
			hdr.Vers,
		))
	case hdr.Vers > 4:

		switch {
		case hdr.Vers > 5:
			tree.entries = r.ReadI64()
			tree.totBytes = r.ReadI64()
			tree.zipBytes = r.ReadI64()
			tree.savedBytes = r.ReadI64()
		default:
			tree.entries = int64(r.ReadF64())
			tree.totBytes = int64(r.ReadF64())
			tree.zipBytes = int64(r.ReadF64())
			tree.savedBytes = int64(r.ReadF64())
		}

		if hdr.Vers >= 18 {
			tree.flushedBytes = r.ReadI64()
		}

		if hdr.Vers >= 16 {
			tree.weight = r.ReadF64()
		}
		tree.timerInterval = r.ReadI32()
		tree.scanField = r.ReadI32()
		tree.update = r.ReadI32()

		if hdr.Vers >= 17 {
			tree.defaultEntryOffsetLen = r.ReadI32()
		}

		nclus := 0
		if hdr.Vers >= 19 { // FIXME
			nclus = int(r.ReadI32()) // fNClusterRange
		}

		if hdr.Vers > 5 {
			tree.maxEntries = r.ReadI64()
		}
		switch {
		case hdr.Vers > 5:
			tree.maxEntryLoop = r.ReadI64()
			tree.maxVirtualSize = r.ReadI64()
			tree.autoSave = r.ReadI64()
		default:
			tree.maxEntryLoop = int64(r.ReadI32())
			tree.maxVirtualSize = int64(r.ReadI32())
			tree.autoSave = int64(r.ReadI32())
		}

		if hdr.Vers >= 18 {
			tree.autoFlush = r.ReadI64()
		}

		switch {
		case hdr.Vers > 5:
			tree.estimate = r.ReadI64()
		default:
			tree.estimate = int64(r.ReadI32())
		}

		if hdr.Vers >= 19 { // FIXME
			tree.clusters.ranges = make([]int64, nclus)
			tree.clusters.sizes = make([]int64, nclus)
			_ = r.ReadI8()
			r.ReadArrayI64(tree.clusters.ranges) // fClusterRangeEnd
			_ = r.ReadI8()
			r.ReadArrayI64(tree.clusters.sizes) // fClusterSize
		}

		if hdr.Vers >= 20 {
			if err := tree.iobits.UnmarshalROOT(r); err != nil {
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

		if hdr.Vers > 5 {
			if v := r.ReadObjectAny(); v != nil {
				tree.aliases = v.(*rcont.List)
			}
		}
		if v := r.ReadObjectAny(); v != nil {
			tree.indexValues = v.(*rcont.ArrayD)
		}
		if v := r.ReadObjectAny(); v != nil {
			tree.index = v.(*rcont.ArrayI)
		}
		if hdr.Vers > 5 {
			if v := r.ReadObjectAny(); v != nil {
				tree.treeIndex = v
			}
			if v := r.ReadObjectAny(); v != nil {
				tree.friends = v.(*rcont.List)
			}
			if v := r.ReadObjectAny(); v != nil {
				tree.userInfo = v.(*rcont.List)
			}
			if v := r.ReadObjectAny(); v != nil {
				tree.branchRef = v
			}
		}
	}

	r.CheckHeader(hdr)

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
			panic(fmt.Errorf("rtree: could not find streamer (type=%q, vers=%d) for branch %q: %w", cls, bre.clsver, br.Name(), err))
		}
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
	if len(br.Branches()) == 0 {
		return
	}

	var members []rbytes.StreamerElement
	switch se := se.(type) {
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
		var (
			typename = strings.TrimSpace(se.TypeName())
			enames   = rmeta.CxxTemplateFrom(typename).Args
			typevers = -1
		)
		switch se.STLType() {
		case rmeta.STLbitset:
			members = []rbytes.StreamerElement{se}

		case rmeta.STLvector, rmeta.STLlist,
			rmeta.STLset, rmeta.STLunorderedset,
			rmeta.STLdeque:
			typename = strings.TrimSpace(strings.TrimRight(enames[0], "*"))
			info, err := ctx.StreamerInfo(typename, typevers)
			if err != nil {
				if _, ok := rmeta.CxxBuiltins[typename]; !ok {
					panic(err)
				}
			}
			if err == nil {
				members = info.Elements()
			}

		case rmeta.STLmap, rmeta.STLunorderedmap:
			// FIXME(sbinet): split it further or not?
			members = []rbytes.StreamerElement{se}

		default:
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

func (*tntuple) RVersion() int16 {
	return rvers.Ntuple
}

func (*tntuple) Class() string {
	return "TNtuple"
}

func (nt *tntuple) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(nt.Class())
	if hdr.Vers > rvers.Ntuple {
		panic(fmt.Errorf(
			"rtree: invalid %s version=%d > %d",
			nt.Class(), hdr.Vers, nt.RVersion(),
		))
	}

	r.ReadObject(&nt.ttree)
	nt.nvars = int(r.ReadI32())

	r.CheckHeader(hdr)
	return r.Err()
}

type tntupleD struct {
	ttree
	nvars int
}

func (*tntupleD) RVersion() int16 {
	return rvers.NtupleD
}

func (*tntupleD) Class() string {
	return "TNtupleD"
}

func (nt *tntupleD) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(nt.Class())
	if hdr.Vers > rvers.NtupleD {
		panic(fmt.Errorf(
			"rtree: invalid %s version=%d > %d",
			nt.Class(), hdr.Vers, nt.RVersion(),
		))
	}

	r.ReadObject(&nt.ttree)
	nt.nvars = int(r.ReadI32())

	r.CheckHeader(hdr)
	return r.Err()
}

type tioFeatures uint8

func (*tioFeatures) RVersion() int16 {
	return rvers.ROOT_IOFeatures
}

func (*tioFeatures) Class() string { return "TIOFeatures" }

func (tio *tioFeatures) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}
	hdr := w.WriteHeader(tio.Class(), tio.RVersion())

	if *tio != 0 {
		var buf = [4]byte{0x1a, 0xa1, 0x2f, 0x10} // FIXME(sbinet) where do these 4 bytes come from ?
		n, err := w.Write(buf[:])
		if err != nil {
			return n, fmt.Errorf("could not write tio marshaled buffer: %w", err)
		}
	}

	w.WriteU8(uint8(*tio))

	return w.SetHeader(hdr)
}

func (tio *tioFeatures) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(tio.Class())
	if hdr.Vers > rvers.ROOT_IOFeatures {
		panic(fmt.Errorf(
			"rtree: invalid %s version=%d > %d",
			tio.Class(), hdr.Vers, tio.RVersion(),
		))
	}

	var buf [4]byte // FIXME(sbinet) where do these 4 bytes come from ?
	_, err := r.Read(buf[:1])
	if err != nil {
		return err
	}

	var u8 uint8
	switch buf[0] {
	case 0:
		// nothing to do.
	default:
		_, err := r.Read(buf[1:])
		if err != nil {
			return err
		}
		u8 = r.ReadU8()
	}

	*tio = tioFeatures(u8)

	r.CheckHeader(hdr)
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
	{
		f := func() reflect.Value {
			o := &tntupleD{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TNtupleD", f)
	}
}

var (
	_ root.Object        = (*ttree)(nil)
	_ root.Named         = (*ttree)(nil)
	_ Tree               = (*ttree)(nil)
	_ rbytes.Marshaler   = (*ttree)(nil)
	_ rbytes.Unmarshaler = (*ttree)(nil)

	_ root.Object        = (*tntuple)(nil)
	_ root.Named         = (*tntuple)(nil)
	_ Tree               = (*tntuple)(nil)
	_ rbytes.Unmarshaler = (*tntuple)(nil)

	_ root.Object        = (*tntupleD)(nil)
	_ root.Named         = (*tntupleD)(nil)
	_ Tree               = (*tntupleD)(nil)
	_ rbytes.Unmarshaler = (*tntupleD)(nil)

	_ root.Object        = (*tioFeatures)(nil)
	_ rbytes.RVersioner  = (*tioFeatures)(nil)
	_ rbytes.Marshaler   = (*tioFeatures)(nil)
	_ rbytes.Unmarshaler = (*tioFeatures)(nil)
)
