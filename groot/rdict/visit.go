// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"strings"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
	"golang.org/x/xerrors"
)

// Visit inspects a streamer info and visits all its elements, once.
func Visit(ctx rbytes.StreamerInfoContext, si rbytes.StreamerInfo, f func(depth int, se rbytes.StreamerElement) error) error {
	v := newVisitor(ctx, f)
	return v.run(0, si)
}

type visitor struct {
	ctx rbytes.StreamerInfoContext
	set map[rbytes.StreamerElement]struct{}
	f   func(depth int, se rbytes.StreamerElement) error
}

func newVisitor(ctx rbytes.StreamerInfoContext, f func(depth int, se rbytes.StreamerElement) error) *visitor {
	if ctx == nil {
		ctx = StreamerInfos
	}
	return &visitor{
		ctx: ctx,
		set: make(map[rbytes.StreamerElement]struct{}),
		f:   f,
	}
}

func (v *visitor) seen(se rbytes.StreamerElement) bool {
	if _, seen := v.set[se]; seen {
		return true
	}
	v.set[se] = struct{}{}
	return false
}

func (v *visitor) run(depth int, si rbytes.StreamerInfo) error {
	for _, se := range si.Elements() {
		err := v.visitSE(depth, se)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *visitor) visitSE(depth int, se rbytes.StreamerElement) error {
	if v.seen(se) {
		return nil
	}

	err := v.f(depth, se)
	if err != nil {
		return err
	}

	switch se.TypeName() {
	case "TVirtualIndex", "TVirtualIndex*":
		return nil
	}

	switch se := se.(type) {
	case *StreamerBasicType:
		// no-op
	case *StreamerBasicPointer:
		// no-op

	case *StreamerBase:
		base, err := v.ctx.StreamerInfo(se.Name(), -1)
		if err != nil {
			return xerrors.Errorf("could not find base %q: %w", se.Name(), err)
		}
		return v.run(depth+1, base)
	case *StreamerObject:
		si, err := v.ctx.StreamerInfo(se.TypeName(), -1)
		if err != nil {
			return xerrors.Errorf("could not find object %q: %w", se.TypeName(), err)
		}
		return v.run(depth+1, si)

	case *StreamerObjectPointer:
		tname := strings.TrimRight(se.TypeName(), "*")
		si, err := v.ctx.StreamerInfo(tname, -1)
		if err != nil {
			return xerrors.Errorf("could not find object-pointer %q: %w", tname, err)
		}
		return v.run(depth+1, si)

	case *StreamerObjectAny:
		tname := se.TypeName()
		si, err := v.ctx.StreamerInfo(tname, -1)
		if err != nil {
			return xerrors.Errorf("could not find object-any %q: %w", tname, err)
		}
		return v.run(depth+1, si)

	case *StreamerObjectAnyPointer:
		tname := strings.TrimRight(se.TypeName(), "*")
		si, err := v.ctx.StreamerInfo(tname, -1)
		if err != nil {
			return xerrors.Errorf("could not find object-any-pointer %q: %w", tname, err)
		}
		return v.run(depth+1, si)

	case *StreamerString, *StreamerSTLstring:
		// no-op

	case *StreamerSTL:
		switch se.STLType() {
		case rmeta.STLdeque, rmeta.STLforwardlist, rmeta.STLlist,
			rmeta.STLset, rmeta.STLunorderedset, rmeta.STLunorderedmultiset,
			rmeta.STLvector:

			etn := se.ElemTypeName()
			tname := strings.TrimRight(etn[0], "*")
			if _, ok := rmeta.CxxBuiltins[tname]; ok {
				// no-op: C++ builtin.
				return nil
			}
			si, err := v.ctx.StreamerInfo(tname, -1)
			if err != nil {
				return xerrors.Errorf("could not find std::container<T> element %q: %w", tname, err)
			}
			return v.run(depth+1, si)

		default:
			return xerrors.Errorf("rdict: cant visit non-vector-like STL streamers %#v", se)
		}

	default:
		panic(xerrors.Errorf("rdict: unknown visit streamer %T", se))
	}

	return nil
}
