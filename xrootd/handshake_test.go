// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build xrootd_test_with_server

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandshake(t *testing.T) {
	_, err := New(context.Background(), *Addr)
	assert.NoError(t, err)
}
