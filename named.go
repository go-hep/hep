// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
)

// The TNamed class is the base class for all named ROOT classes
// A TNamed contains the essential elements (name, title)
// to identify a derived object in containers, directories and files.
// Most member functions defined in this base class are in general
// overridden by the derived classes.
type named struct {
	name  string
	title string
}

// Name returns the name of the instance
func (n *named) Name() string {
	return n.name
}

// Title returns the title of the instance
func (n *named) Title() string {
	return n.title
}

func (n *named) Class() string {
	return "TNamed"
}

func (n *named) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	dec := NewDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	} else {
		myprintf("named: %v %v %v\n", vers, pos, bcnt)
	}

	var id uint32
	err = dec.readBin(&id)
	if err != nil {
		return err
	}

	var bits uint32
	err = dec.readBin(&bits)
	if err != nil {
		return err
	}
	bits |= kIsOnHeap // by definition, de-serialized object is on heap
	if (bits & kIsReferenced) == 0 {
		var trash uint16
		err = dec.readBin(&trash)
		if err != nil {
			return err
		}
	}

	err = dec.readString(&n.name)
	if err != nil {
		return err
	}

	err = dec.readString(&n.title)
	if err != nil {
		return err
	}

	err = dec.checkByteCount(pos, bcnt, start, "TNamed")
	return err
}

var _ Object = (*named)(nil)
var _ ROOTUnmarshaler = (*named)(nil)
