// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	evts := makeEvents(10)
	buf := new(bytes.Buffer)
	w, err := NewWriter(buf)
	if err != nil {
		t.Fatalf("could not create rio writer: %v", err)
	}
	defer w.Close()

	for i := range evts {
		id := fmt.Sprintf("evt-%03d", i+1)
		err = w.WriteValue(id, &evts[i])
		if err != nil {
			t.Fatalf("could not write %s: %v", id, err)
		}
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("could not close rio writer: %v", err)
	}

	r, err := NewReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("could not create rio reader: %v", err)
	}
	defer r.Close()

	nevts := 0
	sc := NewScanner(r)
	sc.Select([]Selector{
		{Name: "evt-001", Unpack: true},
		{Name: "evt-002", Unpack: true},
		{Name: "evt-003", Unpack: true},
		{Name: "evt-004", Unpack: true},
		{Name: "evt-005", Unpack: true},
		{Name: "evt-006", Unpack: true},
		{Name: "evt-007", Unpack: true},
		{Name: "evt-008", Unpack: true},
		{Name: "evt-009", Unpack: true},
		{Name: "evt-010", Unpack: true},
	})

	var evt event
	for sc.Scan() {
		rec := sc.Record()
		if rec == nil {
			break
		}
		blk := rec.Block(rec.Name())
		err = blk.Read(&evt)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(evt, evts[nevts]) {
			t.Fatalf("%s: events differ.\ngot= %#v\nwant=%#v\n", rec.Name(), evt, evts[nevts])
		}
		nevts++
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("error during scan: %v", err)
	}

	if nevts != len(evts) {
		t.Fatalf("invalid number of events: got=%d, want=%d", nevts, len(evts))
	}
}

func makeEvents(n int) []event {
	evts := make([]event, 0, n)
	for i := range n {
		evts = append(evts, event{
			runnbr: int64(1000 + n),
			evtnbr: int64(10000 + n),
			id:     fmt.Sprintf("id-%03d", i),
			eles:   []electron{newElectron(float64(i+1)+100, float64(i+2)+100, float64(i+3)+100, float64(i+4)+100)},
			muons:  []muon{newMuon(float64(i+1)+200, float64(i+2)+200, float64(i+3)+200, float64(i+4)+200)},
		})
	}
	return evts
}
