// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

// dobject is a dummy placeholder object
type dobject struct {
	vers  int16
	size  int32
	class string
}

func (d dobject) Class() string {
	return d.class
}

func (d *dobject) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	d.vers = vers
	d.size = bcnt
	r.setPos(beg + int64(bcnt) + 4)
	r.CheckByteCount(pos, bcnt, beg, d.class)
	return r.err
}
