// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	xrdproto "go-hep.org/x/hep/xrootd/protocol"
	"go-hep.org/x/hep/xrootd/protocol/protocol"
)

// ProtocolInfo is a response for the `Protocol` request. See details in the xrootd protocol specification.
type ProtocolInfo struct {
	BinaryProtocolVersion int32
	ServerType            xrdproto.ServerType
	IsManager             bool
	IsServer              bool
	IsMeta                bool
	IsProxy               bool
	IsSupervisor          bool
	SecurityVersion       byte
	ForceSecurity         bool
	SecurityLevel         protocol.SecurityLevel
	SecurityOverrides     []protocol.SecurityOverride
}

// Protocol obtains the protocol version number, type of the server and security information, such as:
// the security version, the security options, the security level, and the list of alterations
// needed to the specified predefined security level.
func (client *Client) Protocol(ctx context.Context) (ProtocolInfo, error) {
	resp, err := client.call(ctx, protocol.RequestID, protocol.NewRequest(client.protocolVersion, true))
	if err != nil {
		return ProtocolInfo{}, err
	}

	var generalResp protocol.GeneralResponse

	if err = xrdproto.Unmarshal(resp, &generalResp); err != nil {
		return ProtocolInfo{}, err
	}

	var securityInfo *protocol.SecurityInfo
	if len(resp) > protocol.GeneralResponseLength {
		securityInfo = &protocol.SecurityInfo{}
		err = xrdproto.Unmarshal(resp[protocol.GeneralResponseLength:], securityInfo)
		if err != nil {
			return ProtocolInfo{}, err
		}
	}

	var info = ProtocolInfo{
		BinaryProtocolVersion: generalResp.BinaryProtocolVersion,
		ServerType:            extractServerType(generalResp.Flags),

		// TODO: implement bit-encoded fields support in Unmarshal.
		IsManager:    generalResp.Flags&protocol.IsManager != 0,
		IsServer:     generalResp.Flags&protocol.IsServer != 0,
		IsMeta:       generalResp.Flags&protocol.IsMeta != 0,
		IsProxy:      generalResp.Flags&protocol.IsProxy != 0,
		IsSupervisor: generalResp.Flags&protocol.IsSupervisor != 0,
	}

	if securityInfo != nil {
		info.SecurityVersion = securityInfo.SecurityVersion
		info.ForceSecurity = securityInfo.SecurityOptions&protocol.ForceSecurity != 0
		info.SecurityLevel = securityInfo.SecurityLevel

		if securityInfo.SecurityOverridesSize > 0 {
			info.SecurityOverrides = make([]protocol.SecurityOverride, securityInfo.SecurityOverridesSize)

			const offset = protocol.GeneralResponseLength + protocol.SecurityInfoLength
			const elementSize = protocol.SecurityOverrideLength

			for i := byte(0); i < securityInfo.SecurityOverridesSize; i++ {
				err = xrdproto.Unmarshal(resp[offset+elementSize*int(i):], &info.SecurityOverrides[i])
				if err != nil {
					return ProtocolInfo{}, err
				}
			}
		}
	}

	return info, nil
}

func extractServerType(flags protocol.Flags) xrdproto.ServerType {
	if int32(flags)&int32(xrdproto.DataServer) != 0 {
		return xrdproto.DataServer
	}
	return xrdproto.LoadBalancingServer
}
