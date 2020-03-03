// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"regexp"
	"strings"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
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
	var (
		typename = obj.Class()
		cxxtype  = rdict.GoName2Cxx(typename)
		vers     = -1
	)

	if o, ok := obj.(rbytes.RVersioner); ok {
		vers = int(o.RVersion())
	}

	si, ok := rdict.StreamerInfos.Get(cxxtype, vers)
	if !ok {
		return nil, fmt.Errorf("riofs: could not find streamer for %q (version=%d)", cxxtype, vers)
	}
	sictx.addStreamer(si)
	return si, nil
}
