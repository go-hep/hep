// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"reflect"
	"sort"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

type freeSegment struct {
	first int64 // first free word of segment
	last  int64 // last free word of segment
}

func (freeSegment) Class() string {
	return "TFree"
}

func (seg freeSegment) free() int64 {
	return seg.last - seg.first + 1
}

func (seg freeSegment) sizeof() int32 {
	if seg.last > kStartBigFile {
		return 18
	}
	return 10
}

func (seg freeSegment) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()

	w.Grow(int(seg.sizeof()))

	vers := int16(1)
	if seg.last > kStartBigFile {
		vers += 1000
	}
	w.WriteI16(vers)
	switch {
	case vers > 1000:
		w.WriteI64(seg.first)
		w.WriteI64(seg.last)
	default:
		w.WriteI32(int32(seg.first))
		w.WriteI32(int32(seg.last))
	}

	end := w.Pos()
	return int(end - pos), w.Err()
}

func (seg *freeSegment) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	vers := r.ReadI16()
	switch {
	case vers > 1000:
		seg.first = r.ReadI64()
		seg.last = r.ReadI64()
	default:
		seg.first = int64(r.ReadI32())
		seg.last = int64(r.ReadI32())
	}

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &freeSegment{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TFree", f)
}

var (
	_ root.Object        = (*freeSegment)(nil)
	_ rbytes.Marshaler   = (*freeSegment)(nil)
	_ rbytes.Unmarshaler = (*freeSegment)(nil)
)

// freeList describes the list of free segments on a ROOT file.
//
// Each ROOT file has a linked list of free segments.
// Each free segment is described by its first and last addresses.
// When an object is written to a file, a new Key is created. The first free
// segment big enough to accomodate the object is used.
//
// If the object size has a length corresponding to the size of the free segment,
// the free segment is deleted from the list of free segments.
// When an object is deleted from a file, a new freeList object is generated.
// If the deleted object is contiguous to an already deleted object, the free
// segments are merged in one single segment.
type freeList []freeSegment

func (p freeList) Len() int      { return len(p) }
func (p freeList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p freeList) Less(i, j int) bool {
	pi := p[i]
	pj := p[j]
	if pi.first < pj.first {
		return true
	}
	if pi.first == pj.first {
		return pi.last < pj.last
	}
	return false
}

func (fl *freeList) add(first, last int64) *freeSegment {
	elmt := freeSegment{first, last}
	*fl = append(*fl, elmt)
	sort.Sort(*fl)
	fl.consolidate()
	return fl.find(elmt)
}

func (fl *freeList) find(elmt freeSegment) *freeSegment {
	// FIXME(sbinet): use sort.Search
	for i := range *fl {
		cur := &(*fl)[i]
		if elmt.last < cur.first || cur.last < elmt.first {
			continue
		}
		if cur.first <= elmt.first && elmt.first <= cur.last &&
			elmt.last <= cur.last {
			return cur
		}
	}
	return nil
}

func (fl *freeList) consolidate() {
	for i := len(*fl) - 1; i >= 1; i-- {
		cur := &(*fl)[i]
		prev := &(*fl)[i-1]
		if prev.last+1 < cur.first {
			continue
		}
		if cur.last >= prev.last {
			prev.last = cur.last
		}
		fl.remove(i)
	}
}

func (fl *freeList) remove(i int) {
	list := *fl
	*fl = append(list[:i], list[i+1:]...)
}

// best returns the best free segment where to store nbytes.
func (fl freeList) best(nbytes int64) *freeSegment {
	var best *freeSegment

	if len(fl) == 0 {
		return best
	}

	for i, cur := range fl {
		nleft := cur.free()
		if nleft == nbytes {
			// exact match.
			return &fl[i]
		}
		if nleft >= nbytes+4 && best == nil {
			best = &fl[i]
		}
	}

	if best != nil {
		return best
	}

	// try big file
	best = &fl[len(fl)-1]
	best.last += 1000000000
	return best
}

func (fl freeList) last() *freeSegment {
	if len(fl) == 0 {
		return nil
	}
	return &fl[len(fl)-1]
}
