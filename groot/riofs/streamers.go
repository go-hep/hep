// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

func stdvecSIFrom(name, ename string, ctx rbytes.StreamerInfoContext) rbytes.StreamerInfo {
	if etyp, ok := rmeta.CxxBuiltins[ename]; ok {
		return rdict.StreamerOf(ctx, reflect.SliceOf(etyp))
	}
	esi, err := ctx.StreamerInfo(ename, -1)
	if esi == nil || err != nil {
		return nil
	}
	etyp, err := rdict.TypeFromSI(ctx, esi)
	if err != nil || etyp == nil {
		return nil
	}
	return rdict.StreamerOf(ctx, reflect.SliceOf(etyp))
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
