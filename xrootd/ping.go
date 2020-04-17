// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/xrdproto/ping"
)

// Ping determines whether the server is still alive.
func (sess *cliSession) Ping(ctx context.Context) error {
	_, err := sess.Send(ctx, nil, &ping.Request{})
	// TODO: should we react somehow to redirection?
	return err
}
