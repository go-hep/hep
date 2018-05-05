// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package streammanager // import "go-hep.org/x/hep/xrootd/streammanager"

import (
	"sync"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/protocol"
)

type ServerResponse struct {
	Data  []byte
	Error error
}

type dataSendChannel chan<- *ServerResponse
type DataReceiveChannel <-chan *ServerResponse

// StreamManager manages channels by their ids. Basically, it's a map[StreamID] chan<-*ServerResponse with methods
// to claim, free and pass data to some channel by id.
type StreamManager struct {
	dataWaiters map[protocol.StreamID]dataSendChannel
	mutex       sync.Mutex
	freeIds     chan protocol.StreamID
}

// New creates new StreamManager
func New() *StreamManager {
	sm := StreamManager{make(map[protocol.StreamID]dataSendChannel), sync.Mutex{}, make(chan protocol.StreamID, 256*256)}

	for firstByte := uint16(0); firstByte < 256; firstByte++ {
		for secondByte := uint16(0); secondByte < 256; secondByte++ {
			sm.freeIds <- protocol.StreamID{byte(firstByte), byte(secondByte)}
		}
	}

	return &sm
}

// Claim searches for unclaimed id and returns corresponding channel
func (sm *StreamManager) Claim() (id protocol.StreamID, channel DataReceiveChannel, err error) {
	bidirectionalChannel := make(chan *ServerResponse, 1)
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for channel == nil && err == nil {
		select {
		case id = <-sm.freeIds:
			if _, ok := sm.dataWaiters[id]; ok { // Skip id if it was already taken manually via ClaimWithID
				continue
			}

			sm.dataWaiters[id] = bidirectionalChannel
			channel = bidirectionalChannel
		default:
			err = errors.New("Cannot obtain free channel")
		}
	}
	return
}

// ClaimWithID checks if id is unclaimed and returns corresponding channel in case of success
func (sm *StreamManager) ClaimWithID(id protocol.StreamID) (channel DataReceiveChannel, err error) {
	bidirectionalChannel := make(chan *ServerResponse, 1)
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if _, ok := sm.dataWaiters[id]; ok {
		err = errors.Errorf("Channel with id %s is already taken", id)
	} else {
		sm.dataWaiters[id] = bidirectionalChannel
		channel = bidirectionalChannel
	}
	return
}

// Unclaim marks channel with specified id as unclaimed
func (sm *StreamManager) Unclaim(id protocol.StreamID) {
	sm.mutex.Lock()
	close(sm.dataWaiters[id])
	delete(sm.dataWaiters, id)
	sm.mutex.Unlock()
	sm.freeIds <- id
}

// SendData sends data to channel with specific id
func (sm *StreamManager) SendData(id protocol.StreamID, data *ServerResponse) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, ok := sm.dataWaiters[id]; !ok {
		return errors.Errorf("Cannot find data waiter for id %s", id)
	}

	sm.dataWaiters[id] <- data

	return nil
}
