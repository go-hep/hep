// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"bytes"
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
		si := rdict.NewStreamerInfo(name, []rbytes.StreamerElement{
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

	si := rdict.NewStreamerInfo(name, []rbytes.StreamerElement{
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

	// FIXME(sbinet): should we make sure we load the correct version ?
	// or is it okay to just always load the latest one?
	si, err := f.StreamerInfo(obj.Class(), -1)
	if err != nil {
		return nil, err
	}
	rdict.Streamers.Add(si)
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
