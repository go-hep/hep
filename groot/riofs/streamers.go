// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"bytes"
	"log"
	"regexp"
	"strings"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	rstreamerspkg "go-hep.org/x/hep/groot/riofs/internal/rstreamers"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

var (
	reStdVector = regexp.MustCompile("^vector<(.+)>$")
)

func stdvecSIFrom(name, ename string, ctx rbytes.StreamerInfoContext) rbytes.StreamerInfo {
	ename = strings.TrimSpace(ename)
	if etyp, ok := rmeta.CxxBuiltins[ename]; ok {
		si := rdict.NewStreamerInfo(name, 1, []rbytes.StreamerElement{
			rdict.NewStreamerSTL(
				name, rmeta.STLvector, rmeta.GoType2ROOTEnum[etyp],
			),
		})
		return si
	}
	esi, err := ctx.StreamerInfo(ename, -1)
	if esi == nil || err != nil {
		return nil
	}

	si := rdict.NewStreamerInfo(name, 1, []rbytes.StreamerElement{
		rdict.NewStreamerSTL(name, rmeta.STLvector, rmeta.Object),
	})
	return si
}

type streamerInfoStore interface {
	addStreamer(si rbytes.StreamerInfo)
}

func streamerInfoFrom(obj root.Object, sictx streamerInfoStore) (rbytes.StreamerInfo, error) {
	r := &memFile{bytes.NewReader(rstreamerspkg.Data)}
	f, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		typename = obj.Class()
		cxxtype  = rdict.GoName2Cxx(typename)
		vers     = -1
	)

	if o, ok := obj.(rbytes.RVersioner); ok {
		vers = int(o.RVersion())
	}

	si, err := f.StreamerInfo(cxxtype, vers)
	if err != nil {
		return nil, err
	}
	rdict.Streamers.Add(si)
	sictx.addStreamer(si)
	return si, nil
}

func init() {
	// load bootstrap streamers (core ROOT types, such as TObject, TFile, ...)
	r := &memFile{bytes.NewReader(rstreamerspkg.Data)}
	f, err := NewReader(r)
	if err != nil {
		return
	}
	for _, k := range f.Keys() {
		if !strings.HasPrefix(k.Name(), "streamer-info-") {
			continue
		}
		o, err := k.Object()
		if err != nil {
			log.Printf("riofs: could not load streamer info for %q: %v", k.Name(), err)
		}
		rdict.Streamers.Add(o.(rbytes.StreamerInfo))
	}
	defer f.Close()
}
