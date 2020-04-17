// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux // import "go-hep.org/x/hep/xrootd/internal/mux"

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto"
)

func TestMux_Claim(t *testing.T) {
	m := New()
	defer m.Close()
	claimedIds := map[xrdproto.StreamID]bool{}

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

	for id := range claimedIds {
		m.Unclaim(id)
	}
}

func TestMux_Claim_AfterUnclaim(t *testing.T) {
	m := New()
	defer m.Close()
	claimedIds := map[xrdproto.StreamID]bool{}

	for i := 0; i < streamIDPoolSize; i++ {
		id, _, _ := m.Claim()
		claimedIds[id] = true
	}

	wantID := xrdproto.StreamID{13, 14}
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

	for id := range claimedIds {
		m.Unclaim(id)
	}
}

func TestMux_ClaimWithID_WhenIDIsFree(t *testing.T) {
	m := New()
	defer m.Close()

	streamID := xrdproto.StreamID{13, 14}
	channel, err := m.ClaimWithID(streamID)

	if err != nil {
		t.Fatalf("could not Claim streamID: %v", err)
	}

	if channel == nil {
		t.Fatal("channel should not be nil")
	}

	m.Unclaim(streamID)
}

func TestMux_ClaimWithID_WhenIDIsTakenByClaimWithID(t *testing.T) {
	m := New()
	defer m.Close()
	streamID := xrdproto.StreamID{13, 14}
	m.ClaimWithID(streamID)

	_, err := m.ClaimWithID(streamID)

	if err == nil {
		t.Fatal("should not be able to ClaimWithID when that id is already claimed")
	}

	m.Unclaim(streamID)
}

func TestMux_ClaimWithID_WhenIDIsTakenByClaim(t *testing.T) {
	m := New()
	defer m.Close()
	id, _, _ := m.Claim()

	_, err := m.ClaimWithID(id)

	if err == nil {
		t.Fatal("should not be able to ClaimWithID when that id is already claimed")
	}

	m.Unclaim(id)
}

func TestMux_Claim_WhenIDIsTakenByClaimWithID(t *testing.T) {
	m := New()
	defer m.Close()
	takenID := xrdproto.StreamID{0, 0}
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

	m.Unclaim(takenID)
	m.Unclaim(id)
}

func TestMux_SendData_WhenIDIsTaken(t *testing.T) {
	m := New()
	defer m.Close()
	takenID := xrdproto.StreamID{0, 0}
	want := ServerResponse{}
	var got ServerResponse

	channel, _ := m.ClaimWithID(takenID)
	go func() {
		err := m.SendData(takenID, want)
		if err != nil {
			t.Fatalf("could not SendData: %v", err)
		}
	}()

	got = <-channel
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid data\ngot = %v\nwant = %v", got, want)
	}

	m.Unclaim(takenID)
}

func TestMux_SendData_WhenIDIsNotTaken(t *testing.T) {
	m := New()
	defer m.Close()
	notTakenID := xrdproto.StreamID{0, 0}

	err := m.SendData(notTakenID, ServerResponse{})

	if err == nil {
		t.Fatal("should not be able to SenData when id is unclaimed")
	}
}

func TestMux_Close_WhenAlreadyClosed(t *testing.T) {
	m := New()
	m.Close()
	m.Close()
}

func TestMux_Unclaim_WhenNotClaimed(t *testing.T) {
	m := New()
	defer m.Close()
	m.Unclaim(xrdproto.StreamID{0, 0})
}

func TestMux_Claim_WhenClosed(t *testing.T) {
	m := New()
	m.Close()
	_, _, err := m.Claim()
	if err == nil {
		t.Fatal("should not be able to Claim when mux is closed")
	}
}

func TestMux_ClaimWithID_WhenClosed(t *testing.T) {
	m := New()
	m.Close()
	_, err := m.ClaimWithID(xrdproto.StreamID{0, 0})
	if err == nil {
		t.Fatal("should not be able to ClaimWithID when mux is closed")
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
	response := ServerResponse{Data: []byte{0, 1, 2, 3, 4, 5}}

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
