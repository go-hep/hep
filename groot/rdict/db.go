// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"sort"
	"sync"

	"go-hep.org/x/hep/groot/rbytes"
	"golang.org/x/xerrors"
)

// StreamerInfos stores all the streamers available at runtime.
var StreamerInfos = &streamerDb{
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
		return nil, xerrors.Errorf("rdict: no streamer for %q (version=%d)", name, vers)
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
				continue
			}
			title := v.Title()
			if title == "" {
				continue
			}
			if name, ok := Typename(class, title); ok && name == class {
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
				panic(xerrors.Errorf("rdict: StreamerInfo class=%q version=%d with checksum=%d (got checksum=%d)",
					streamer.Name(), streamer.ClassVersion(), streamer.CheckSum(), old.CheckSum(),
				))
			}
			return
		}
	}

	db.db[key] = streamer
}

// Values returns all the known StreamerInfos.
func (db *streamerDb) Values() []rbytes.StreamerInfo {
	db.RLock()
	defer db.RUnlock()

	var sinfos = make([]rbytes.StreamerInfo, 0, len(db.db))
	for _, si := range db.db {
		sinfos = append(sinfos, si)
	}

	sort.Slice(sinfos, func(i, j int) bool {
		si := sinfos[i]
		sj := sinfos[j]
		if si.Name() == sj.Name() {
			return si.ClassVersion() < sj.ClassVersion()
		}
		return si.Name() < sj.Name()
	})

	return sinfos
}

var (
	_ rbytes.StreamerInfoContext = (*streamerDb)(nil)
)
