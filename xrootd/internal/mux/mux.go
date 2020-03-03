// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mux implements the multiplexer that manages access to and writes data to the channels
by corresponding StreamID from xrootd protocol specification.

Example of usage:

	mux := New()
	defer m.Close()

	// Claim channel for response retrieving.
	id, channel, err := m.Claim()
	if err != nil {
		// handle error.
	}

	// Send a request to the server using id as a streamID.

	go func() {
		// Read response from the server.
		// ...

		// Send response to the awaiting caller using streamID from the server.
		err := m.SendData(streamID, want)
		if err != nil {
			// handle error.
		}
	}


	// Fetch response.
	response := <-channel
*/
package mux // import "go-hep.org/x/hep/xrootd/internal/mux"

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"go-hep.org/x/hep/xrootd/xrdproto"
)

// ServerResponse contains slice of bytes Data representing data from
// XRootD server response (see XRootD protocol specification) and
// Err representing error received from server or occurred
// during response decoding.
type ServerResponse struct {
	Data        []byte
	Err         error
	Redirection *Redirection
}

// Redirection represents the redirection request from the server.
// It contains address addr to which client must connect,
// opaque data that must be delivered to the new server as
// opaque information added to the file name, and token that
// must be delivered to the new server as part of login request.
type Redirection struct {
	// Addr is the server address to which client must connect in the format "host:port".
	Addr string

	// Opaque is the data that must be delivered to the new server as
	// opaque information added to the file name
	Opaque string

	// Token is the data that must be delivered to the new server as
	// part of the login request.
	Token string
}

// ParseRedirection parses the Redirection from the XRootD redirect response format.
// See http://xrootd.org/doc/dev45/XRdv310.pdf, p. 33 for details.
func ParseRedirection(raw []byte) (*Redirection, error) {
	port := binary.BigEndian.Uint32(raw)
	parts := strings.Split(string(raw[4:]), "?")
	if len(parts) == 0 {
		return nil, fmt.Errorf("xrootd: could not parse redirect url %q", string(raw))
	}

	var opaque, token string
	if len(parts) > 1 {
		opaque = parts[1]
	}
	if len(parts) > 2 {
		token = parts[2]
	}
	addr := parts[0] + ":" + strconv.Itoa(int(port))
	return &Redirection{Addr: addr, Opaque: opaque, Token: token}, nil
}

type dataSendChan chan<- ServerResponse
type DataRecvChan <-chan ServerResponse

const streamIDPartSize = math.MaxUint8
const streamIDPoolSize = streamIDPartSize * streamIDPartSize

// Mux manages channels by their ids.
// Basically, it's a map[StreamID] chan<-ServerResponse
// with methods to claim, free and pass data to a specific channel by id.
type Mux struct {
	mu          sync.Mutex
	dataWaiters map[xrdproto.StreamID]dataSendChan
	freeIDs     chan uint16
	quit        chan struct{}
	closed      bool
}

// New creates a new Mux.
func New() *Mux {
	const freeIDsBufferSize = 32 // 32 is completely arbitrary ATM and should be refined based on real use cases.

	m := Mux{
		dataWaiters: make(map[xrdproto.StreamID]dataSendChan),
		freeIDs:     make(chan uint16, freeIDsBufferSize),
		quit:        make(chan struct{}),
	}

	go func() {
		var i uint16 = 0
		for {
			select {
			case m.freeIDs <- i:
				i = (i + 1) % streamIDPoolSize
			case <-m.quit:
				close(m.freeIDs)
				return
			}
		}
	}()

	return &m
}

// Close closes the Mux.
func (m *Mux) Close() {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return
	}
	m.closed = true
	m.mu.Unlock()
	close(m.quit)

	response := ServerResponse{Err: errors.New("xrootd: close was called before response was fully received")}
	for streamID := range m.dataWaiters {
		m.SendData(streamID, response)
		m.Unclaim(streamID)
	}
}

// Claim searches for unclaimed id and returns corresponding channel.
func (m *Mux) Claim() (xrdproto.StreamID, DataRecvChan, error) {
	ch := make(chan ServerResponse)

	for {
		id := <-m.freeIDs
		streamId := xrdproto.StreamID{byte(id >> 8), byte(id)}

		m.mu.Lock()
		if m.closed {
			m.mu.Unlock()
			return xrdproto.StreamID{}, nil, errors.New("mux: Claim was called on closed Mux")
		}
		if _, claimed := m.dataWaiters[streamId]; claimed { // Skip id if it was already claimed manually via ClaimWithID
			m.mu.Unlock()
			continue
		}

		m.dataWaiters[streamId] = ch
		m.mu.Unlock()
		return streamId, ch, nil
	}
}

// ClaimWithID checks if id is unclaimed and returns the corresponding channel in case of success.
func (m *Mux) ClaimWithID(id xrdproto.StreamID) (DataRecvChan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, errors.New("mux: ClaimWithID was called on closed Mux")
	}
	ch := make(chan ServerResponse)

	if _, claimed := m.dataWaiters[id]; claimed {
		return nil, fmt.Errorf("mux: channel with id %v is already claimed", id)
	}

	m.dataWaiters[id] = ch

	return ch, nil
}

// Unclaim marks channel with specified id as unclaimed.
func (m *Mux) Unclaim(id xrdproto.StreamID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.dataWaiters[id]; ok {
		close(m.dataWaiters[id])
		delete(m.dataWaiters, id)
	}
}

// SendData sends data to channel with specific id.
func (m *Mux) SendData(id xrdproto.StreamID, data ServerResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.dataWaiters[id]; !ok {
		return fmt.Errorf("mux: cannot find data waiter for id %v", id)
	}

	m.dataWaiters[id] <- data

	return nil
}
