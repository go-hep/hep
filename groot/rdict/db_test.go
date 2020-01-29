// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict_test

import (
	"testing"

	"go-hep.org/x/hep/groot/rdict"
)

func TestStreamerInfosDbList(t *testing.T) {
	found := false
	sinfos := rdict.StreamerInfos.Values()
	for _, si := range sinfos {
		if si.Name() != "TObject" {
			continue
		}
		found = true
	}

	if !found {
		t.Fatalf("could not find streamer for TObject")
	}
}
