// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"go-hep.org/x/hep/groot/internal/rmeta"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	rstreamerspkg "go-hep.org/x/hep/groot/riofs/internal/rstreamers"
	"go-hep.org/x/hep/groot/root"
)

var (
	reStdVector = regexp.MustCompile("^vector<(.+)>$")
)

func stdvecSIFrom(name, ename string, ctx rbytes.StreamerInfoContext) rbytes.StreamerInfo {
	ename = strings.TrimSpace(ename)
	if etyp, ok := rmeta.CxxBuiltins[ename]; ok {
		si := rdict.NewStreamerInfo(name, []rbytes.StreamerElement{
			rdict.NewStreamerSTL(
				name, rmeta.STLvector, rmeta.GoType2ROOTEnum[etyp],
			),
		})
		return si
	}
	esi, err := ctx.StreamerInfo(ename)
	if esi == nil || err != nil {
		return nil
	}

	si := rdict.NewStreamerInfo(name, []rbytes.StreamerElement{
		rdict.NewStreamerSTL(name, rmeta.STLvector, rmeta.Object),
	})
	return si
}

type streamerInfoStore interface {
	addStreamer(si rbytes.StreamerInfo)
}

var Streamers = streamerDb{
	db: make(map[streamerDbKey]rbytes.StreamerInfo),
}

type streamerDbKey struct {
	class    string
	version  int
	checksum int
}

type streamerDb struct {
	sync.RWMutex
	db map[streamerDbKey]rbytes.StreamerInfo
}

func (db *streamerDb) GetAny(class string) (rbytes.StreamerInfo, bool) {
	db.RLock()
	defer db.RUnlock()
	for k, v := range db.db {
		if k.class == class {
			return v, true
		}
	}
	return nil, false
}

func (db *streamerDb) Get(class string, vers int, chksum int) (rbytes.StreamerInfo, bool) {
	db.RLock()
	defer db.RUnlock()
	key := streamerDbKey{
		class:    class,
		version:  vers,
		checksum: chksum,
	}

	streamer, ok := db.db[key]
	if !ok {
		return nil, false
	}
	return streamer, true
}

func (db *streamerDb) add(streamer rbytes.StreamerInfo) {
	db.Lock()
	defer db.Unlock()

	key := streamerDbKey{
		class:    streamer.Name(),
		version:  streamer.ClassVersion(),
		checksum: streamer.CheckSum(),
	}

	old, dup := db.db[key]
	if dup {
		if old.CheckSum() != streamer.CheckSum() {
			panic(fmt.Errorf("riofs: StreamerInfo class=%q version=%d with checksum=%d (got checksum=%d)",
				streamer.Name(), streamer.ClassVersion(), streamer.CheckSum(), old.CheckSum(),
			))
		}
		return
	}

	db.db[key] = streamer
}

func streamerInfoFrom(obj root.Object, sictx streamerInfoStore) (rbytes.StreamerInfo, error) {
	r := &memFile{bytes.NewReader(rstreamerspkg.Data)}
	f, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	si, err := f.StreamerInfo(obj.Class())
	if err != nil {
		return nil, err
	}
	Streamers.add(si)
	sictx.addStreamer(si)
	return si, nil
}

var defaultStreamerInfos []rbytes.StreamerInfo

func init() {
	r := &memFile{bytes.NewReader(rstreamerspkg.Data)}
	f, err := NewReader(r)
	if err != nil {
		return
	}
	defer f.Close()

	defaultStreamerInfos = f.StreamerInfos()
}
