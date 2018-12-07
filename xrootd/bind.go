// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/bind"
)

func (sess *cliSession) bind(ctx context.Context, sessionID [16]byte) (xrdproto.PathID, error) {
	var resp bind.Response
	_, err := sess.Send(ctx, &resp, &bind.Request{SessionID: sessionID})
	// TODO: should we react somehow to redirection?
	return resp.PathID, err
}
