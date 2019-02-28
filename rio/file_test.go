// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestFile(t *testing.T) {
	buf := new(bytes.Buffer)
	w, err := NewWriter(buf)
	if err != nil {
		t.Fatalf("error creating new rio Writer: %v", err)
	}

	evt1 := event{
		runnbr: 10,
		evtnbr: 11,
		id:     "id-11",
		eles:   []electron{newElectron(11, 12, 13, 14)},
		muons:  []muon{newMuon(11, 12, 13, 14)},
	}
	err = w.WriteValue("evt1", &evt1)
	if err != nil {
		t.Fatalf("error writing value [evt1]: %v", err)
	}

	evt2 := event{
		runnbr: 10,
		evtnbr: 12,
		id:     "id-22",
		eles:   []electron{newElectron(21, 22, 23, 24)},
		muons:  []muon{newMuon(21, 22, 23, 24)},
	}
	err = w.WriteValue("evt2", &evt2)
	if err != nil {
		t.Fatalf("error writing value [evt2]: %v", err)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("error closing rio writer: %v", err)
	}

	r := bytes.NewReader(buf.Bytes())

	f, err := Open(r)
	if err != nil {
		t.Fatalf("error opening file: %v\n", err)
	}
	defer f.Close()

	keys := []RecordDesc{
		{
			Name: "evt1",
			Blocks: []BlockDesc{
				{
					Name: "evt1",
					Type: "*go-hep.org/x/hep/rio.event",
				},
			},
		},
		{
			Name: "evt2",
			Blocks: []BlockDesc{
				{
					Name: "evt2",
					Type: "*go-hep.org/x/hep/rio.event",
				},
			},
		},
	}
	if !reflect.DeepEqual(keys, f.Keys()) {
		t.Fatalf("keys differ.\ngot= %v\nwant=%v\n", f.Keys(), keys)
	}

	if !f.Has("evt1") {
		t.Fatalf("expected 'evt1' to be present in file")
	}
	if f.Has("not-there") {
		t.Fatalf("did not expect a 'not-there' record in file")
	}

	var (
		revt1 event
		revt2 event
	)

	err = f.Get("evt2", &revt2)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}

	err = f.Get("evt1", &revt1)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}

	if !reflect.DeepEqual(evt1, revt1) {
		t.Fatalf("evt1 differ.\ngot= %v\nwant=%v\n", revt1, evt1)
	}

	if !reflect.DeepEqual(evt2, revt2) {
		t.Fatalf("evt2 differ.\ngot= %v\nwant=%v\n", revt2, evt2)
	}
}

func TestInvalidFile(t *testing.T) {
	for _, tc := range []struct {
		r   []byte
		err error
	}{
		{
			r:   nil,
			err: errors.Errorf("rio: error reading magic-header: EOF"),
		},
		{
			r:   []byte{'s', 'i', 'o', '\x00'},
			err: errors.Errorf("rio: not a rio-stream. magic-header=\"sio\\x00\". want=\"rio\\x00\""),
		},
		{
			r:   []byte{'r', 'i', 'o', '\x00'},
			err: errors.Errorf("rio: error seeking footer (err=bytes.Reader.Seek: negative position)"),
		},
	} {
		t.Run("", func(t *testing.T) {
			r := bytes.NewReader(tc.r)
			f, err := Open(r)
			if !reflect.DeepEqual(err, tc.err) {
				t.Fatalf("got=%#v, want=%#v", err, tc.err)
			}
			defer f.Close()
		})
	}
}
