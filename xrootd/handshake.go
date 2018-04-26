// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/requests/handshake"
	"go-hep.org/x/hep/xrootd/streammanager"
)

func (client *Client) handshake(ctx context.Context) error {
	responseChannel, err := client.sm.ClaimWithID(streammanager.StreamID{0, 0})
	if err != nil {
		return err
	}

	requestBytes, err := encoder.Marshal(handshake.NewRequest())
	if err != nil {
		return err
	}

	responseBytes, err := client.callWithBytesAndResponseChannel(ctx, responseChannel, requestBytes)
	if err != nil {
		return err
	}

	var handshakeResult handshake.Response
	if encoder.Unmarshal(responseBytes, &handshakeResult) != nil {
		return err
	}

	client.protocolVersion = handshakeResult.ProtocolVersion
	logger.Printf("Connected! Protocol version is %d. Server type is %s.", handshakeResult.ProtocolVersion, handshakeResult.ServerType)

	return nil
}
