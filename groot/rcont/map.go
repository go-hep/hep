// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcont

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// Map is a ROOT associative array of (key,value) pairs.
// Keys and values must implement the root.Object interface.
type Map struct {
	obj  rbase.Object
	name string
	tbl  map[root.Object]root.Object
}

func NewMap() *Map {
	return &Map{
		obj:  *rbase.NewObject(),
		name: "TMap",
		tbl:  make(map[root.Object]root.Object),
	}
}

func (*Map) RVersion() int16 { return rvers.Map }
func (*Map) Class() string   { return "TMap" }
func (m *Map) UID() uint32   { return m.obj.UID() }

func (m *Map) Name() string        { return m.name }
func (m *Map) Title() string       { return "A (key,value) map" }
func (m *Map) SetName(name string) { m.name = name }

// Table returns the underlying hash table.
func (m *Map) Table() map[root.Object]root.Object { return m.tbl }

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (m *Map) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(m.RVersion())
	_, _ = m.obj.MarshalROOT(w)
	w.WriteString(m.name)

	w.WriteI32(int32(len(m.tbl)))

	for k, v := range m.tbl {
		_ = w.WriteObjectAny(k)
		_ = w.WriteObjectAny(v)
	}

	return w.SetByteCount(pos, m.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (m *Map) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(m.Class())

	if vers > 2 {
		if err := m.obj.UnmarshalROOT(r); err != nil {
			return err
		}
	}
	if vers > 1 {
		m.name = r.ReadString()
	}

	nobjs := int(r.ReadI32())
	m.tbl = make(map[root.Object]root.Object, nobjs)
	for i := 0; i < nobjs; i++ {
		k := r.ReadObjectAny()
		if r.Err() != nil {
			return r.Err()
		}
		v := r.ReadObjectAny()
		if r.Err() != nil {
			return r.Err()
		}
		if k != nil {
			m.tbl[k] = v
		}
	}

	r.CheckByteCount(pos, bcnt, start, m.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := NewMap()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TMap", f)
}

var (
	_ root.Object        = (*Map)(nil)
	_ root.UIDer         = (*Map)(nil)
	_ root.Named         = (*Map)(nil)
	_ rbytes.Marshaler   = (*Map)(nil)
	_ rbytes.Unmarshaler = (*Map)(nil)
)
