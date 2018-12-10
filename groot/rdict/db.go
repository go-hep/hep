// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbytes"
)

// Streamers stores all the streamers available at runtime.
var Streamers = &streamerDb{
	db: make(map[streamerDbKey]rbytes.StreamerInfo),
}

type streamerDbKey struct {
	class   string
	version int
}

type streamerDb struct {
	sync.RWMutex
	db map[streamerDbKey]rbytes.StreamerInfo
}

func (db *streamerDb) StreamerInfo(name string, vers int) (rbytes.StreamerInfo, error) {
	si, ok := db.Get(name, vers)
	if !ok {
		return nil, errors.Errorf("rdict: no streamer for %q", name)
	}
	return si, nil
}

func (db *streamerDb) Get(class string, vers int) (rbytes.StreamerInfo, bool) {
	db.RLock()
	defer db.RUnlock()
	switch {
	case vers < 0:
		var slice []rbytes.StreamerInfo
		for k, v := range db.db {
			if k.class == class {
				slice = append(slice, v)
			}
		}
		if len(slice) == 0 {
			return nil, false
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].ClassVersion() < slice[j].ClassVersion()
		})
		return slice[len(slice)-1], true
	default:
		key := streamerDbKey{
			class:   class,
			version: vers,
		}

		streamer, ok := db.db[key]
		if !ok {
			return nil, false
		}
		return streamer, true
	}
}

// FIXME(sbinet): ROOT changed its checksum behaviour at some point.
// our reference ROOT files have been caught in the middle of this migration.
// disable the check for duplicate streamers with different checksums for now.
const checkdups = false

func (db *streamerDb) Add(streamer rbytes.StreamerInfo) {
	db.Lock()
	defer db.Unlock()

	key := streamerDbKey{
		class:   streamer.Name(),
		version: streamer.ClassVersion(),
	}

	if checkdups {
		old, dup := db.db[key]
		if dup {
			if old.CheckSum() != streamer.CheckSum() {
				panic(fmt.Errorf("rdict: StreamerInfo class=%q version=%d with checksum=%d (got checksum=%d)",
					streamer.Name(), streamer.ClassVersion(), streamer.CheckSum(), old.CheckSum(),
				))
			}
			return
		}
	}

	db.db[key] = streamer
}

var (
	_ rbytes.StreamerInfoContext = (*streamerDb)(nil)
)
