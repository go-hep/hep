// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux // import "go-hep.org/x/hep/xrootd/mux"

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/protocol"
)

func TestMux_Claim(t *testing.T) {
	m := New()
	defer m.Close()
	claimedIds := map[protocol.StreamID]bool{}

	for i := 0; i < streamIDPoolSize; i++ {
		id, channel, err := m.Claim()

		if err != nil {
			t.Fatalf("could not Claim streamID: %v", err)
		}

		if channel == nil {
			t.Fatal("channel should not be nil")
		}

		if claimedIds[id] {
			t.Fatalf("should not claim id %s, which was already claimed", id)
		}

		claimedIds[id] = true
	}
}

func TestMux_Claim_AfterUnclaim(t *testing.T) {
	m := New()
	defer m.Close()
	claimedIds := map[protocol.StreamID]bool{}

	for i := 0; i < streamIDPoolSize; i++ {
		id, _, _ := m.Claim()
		claimedIds[id] = true
	}

	wantID := protocol.StreamID{13, 14}
	m.Unclaim(wantID)

	gotID, channel, err := m.Claim()

	if err != nil {
		t.Fatalf("could not Claim streamID: %v", err)
	}

	if channel == nil {
		t.Fatal("channel should not be nil")
	}

	if !reflect.DeepEqual(gotID, wantID) {
		t.Fatalf("invalid claim\ngot = %v\nwant = %v", gotID, wantID)
	}
}

func TestMux_ClaimWithID_WhenIDIsFree(t *testing.T) {
	m := New()
	defer m.Close()

	channel, err := m.ClaimWithID(protocol.StreamID{13, 14})

	if err != nil {
		t.Fatalf("could not Claim streamID: %v", err)
	}

	if channel == nil {
		t.Fatal("channel should not be nil")
	}
}

func TestMux_ClaimWithID_WhenIDIsTakenByClaimWithID(t *testing.T) {
	m := New()
	defer m.Close()
	m.ClaimWithID(protocol.StreamID{13, 14})

	_, err := m.ClaimWithID(protocol.StreamID{13, 14})

	if err == nil {
		t.Fatal("should not be able to ClaimWithID when that id is already claimed")
	}
}

func TestMux_ClaimWithID_WhenIDIsTakenByClaim(t *testing.T) {
	m := New()
	defer m.Close()
	id, _, _ := m.Claim()

	_, err := m.ClaimWithID(id)

	if err == nil {
		t.Fatal("should not be able to ClaimWithID when that id is already claimed")
	}
}

func TestMux_Claim_WhenIDIsTakenByClaimWithID(t *testing.T) {
	m := New()
	defer m.Close()
	takenID := protocol.StreamID{0, 0}
	m.ClaimWithID(takenID)

	id, channel, err := m.Claim()

	if err != nil {
		t.Fatalf("could not Claim streamID: %v", err)
	}

	if channel == nil {
		t.Fatal("channel should not be nil")
	}

	if reflect.DeepEqual(id, takenID) {
		t.Fatalf("invalid claim: id %v was already taken", takenID)
	}
}

func TestMux_SendData_WhenIDIsTaken(t *testing.T) {
	m := New()
	defer m.Close()
	takenID := protocol.StreamID{0, 0}
	want := ServerResponse{}

	channel, _ := m.ClaimWithID(takenID)
	err := m.SendData(takenID, want)

	if err != nil {
		t.Fatalf("could not SendData: %v", err)
	}

	got := <-channel
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid data\ngot = %v\nwant = %v", got, want)
	}
}

func TestMux_SendData_WhenIDIsNotTaken(t *testing.T) {
	m := New()
	defer m.Close()
	notTakenID := protocol.StreamID{0, 0}

	err := m.SendData(notTakenID, ServerResponse{})

	if err == nil {
		t.Fatal("should not be able to SenData when id is unclaimed")
	}
}

func BenchmarkMux_Claim(b *testing.B) {
	m := New()
	defer m.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, _, err := m.Claim()
		if err != nil {
			b.Error(err)
		}
		m.Unclaim(id)
	}
}

func BenchmarkMux_SendData(b *testing.B) {
	m := New()
	defer m.Close()
	id, ch, _ := m.Claim()
	done := make(chan struct{})
	response := ServerResponse{[]byte{0, 1, 2, 3, 4, 5}, nil}

	go func() {
		for {
			select {
			case <-ch:
			case <-done:
				return
			}
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := m.SendData(id, response)
		if err != nil {
			b.Error(err)
		}
	}

	m.Unclaim(id)
	close(done)
}
