// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bytes"
	"reflect"
	"testing"
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

	keys := []string{"evt1", "evt2"}
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
