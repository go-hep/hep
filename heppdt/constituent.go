// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

// Constituent holds a particle constituent
// (e.g. quark type and number of quarks of this type)
type Constituent struct {
	ID  PID // particle ID
	Mul int // multiplicity
}

// IsUp returns whether this is an up-quark
func (c Constituent) IsUp() bool {
	panic("not implemented")
}

// IsDown returns whether this is a down-quark
func (c Constituent) IsDown() bool {
	return false
}

// IsStrange returns whether this is a strqnge-quark
func (c Constituent) IsStrange() bool {
	panic("not implemented")
}

// IsCharm returns whether this is a charm-quark
func (c Constituent) IsCharm() bool {
	panic("not implemented")
}

// IsBottom returns whether this is a bottom-quark
func (c Constituent) IsBottom() bool {
	panic("not implemented")
}

// IsTop returns whether this is a top-quark
func (c Constituent) IsTop() bool {
	panic("not implemented")
}
