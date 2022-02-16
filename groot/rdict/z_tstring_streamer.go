// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rmeta"
)

func init() {
	si, ok := StreamerInfos.Get("TString", -1)
	if !ok {
		panic(fmt.Errorf("rdict: could not get streamer info for TString"))
	}
	if len(si.Elements()) != 0 {
		return
	}

	// FIXME(sbinet): the ROOT/C++ streamer for TString is a simple placeholder.
	// but groot relies on the actual list of StreamerElements to generate the r/w-streaming code.
	// So, apply this "regularization" and hope for the best.
	sinfo := si.(*StreamerInfo)
	sinfo.elems = append(sinfo.elems, &StreamerBasicType{
		StreamerElement: Element{
			Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
			Type:   rmeta.TString,
			Size:   25,
			MaxIdx: [5]int32{0, 0, 0, 0, 0},
			EName:  "TString",
		}.New(),
	})
}
