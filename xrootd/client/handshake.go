// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/handshake"
)

func (sess *session) handshake(ctx context.Context) error {
	streamID := xrdproto.StreamID{0, 0}
	responseChannel, err := sess.mux.ClaimWithID(streamID)
	if err != nil {
		return err
	}

	req := handshake.NewRequest()
	var wBuffer xrdenc.WBuffer
	err = req.MarshalXrd(&wBuffer)
	if err != nil {
		return err
	}

	resp, _, err := sess.send(ctx, streamID, responseChannel, wBuffer.Bytes(), nil, 0)
	// TODO: should we react somehow to redirection?
	if err != nil {
		return err
	}

	var result handshake.Response
	if err = xrdproto.Unmarshal(resp, &result); err != nil {
		return err
	}

	sess.protocolVersion = result.ProtocolVersion

	return nil
}
