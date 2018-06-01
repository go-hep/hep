// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/protocol"
	"go-hep.org/x/hep/xrootd/protocol/handshake"
)

func (client *Client) handshake(ctx context.Context) error {
	responseChannel, err := client.mux.ClaimWithID(protocol.StreamID{0, 0})
	if err != nil {
		return err
	}

	req := handshake.NewRequest()
	raw, err := req.MarshalXrd()
	if err != nil {
		return err
	}

	resp, err := client.send(ctx, responseChannel, raw)
	if err != nil {
		return err
	}

	var result handshake.Response
	if err = protocol.Unmarshal(resp, &result); err != nil {
		return err
	}

	client.protocolVersion = result.ProtocolVersion

	return nil
}
