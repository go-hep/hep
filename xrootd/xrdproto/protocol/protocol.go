// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protocol contains the structures describing request and response
// for protocol request (see XRootD specification).
//
// A response consists of 3 parts:
//
// 1) A general response that is always returned and specifies protocol version and flags describing server type.
//
// 2) A response part that is added to the general response
// if `ReturnSecurityRequirements` is provided and server supports it.
// It contains the security version, the security options, the security level,
// and the number of following security overrides, if any.
//
// 3) A list of SecurityOverride - alterations needed to the specified predefined security level.
package protocol // import "go-hep.org/x/hep/xrootd/xrdproto/protocol"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3006

// Flags are the Flags that define xrootd server type. See xrootd protocol specification for further info.
type Flags int32

const (
	IsServer     Flags = 0x00000001 // IsServer indicates whether this server has server role.
	IsManager    Flags = 0x00000002 // IsManager indicates whether this server has manager role.
	IsMeta       Flags = 0x00000100 // IsMeta indicates whether this server has meta attribute.
	IsProxy      Flags = 0x00000200 // IsProxy indicates whether this server has proxy attribute.
	IsSupervisor Flags = 0x00000400 // IsSupervisor indicates whether this server has supervisor attribute.
)

// SecurityOptions are the security-related options.
// See specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, p. 72.
type SecurityOptions byte

const (
	// ForceSecurity specifies that signing is required even if the authentication
	// protocol does not support generic encryption.
	ForceSecurity SecurityOptions = 0x02
)

// RequestOptions specifies what should be returned as part of response.
type RequestOptions byte

const (
	// RequestOptionsNone specifies that only general response should be returned.
	RequestOptionsNone RequestOptions = 0
	// ReturnSecurityRequirements specifies that security requirements should be returned
	// if that's supported by the server.
	ReturnSecurityRequirements RequestOptions = 1
)

// Request holds protocol request parameters.
type Request struct {
	ClientProtocolVersion int32
	Options               RequestOptions
	_                     [11]byte
	_                     int32
}

// NewRequest forms a Request according to provided parameters.
func NewRequest(protocolVersion int32, withSecurityRequirements bool) *Request {
	var options = RequestOptionsNone
	if withSecurityRequirements {
		options |= ReturnSecurityRequirements
	}
	return &Request{ClientProtocolVersion: protocolVersion, Options: options}
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteI32(o.ClientProtocolVersion)
	wBuffer.WriteU8(byte(o.Options))
	wBuffer.Next(15)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.ClientProtocolVersion = rBuffer.ReadI32()
	o.Options = RequestOptions(rBuffer.ReadU8())
	rBuffer.Skip(15)
	return nil
}

// Response is a response for the `Protocol` request. See details in the xrootd protocol specification.
type Response struct {
	BinaryProtocolVersion int32
	Flags                 Flags
	HasSecurityInfo       bool
	_                     byte
	_                     byte
	SecurityVersion       byte
	SecurityOptions       SecurityOptions
	SecurityLevel         xrdproto.SecurityLevel
	SecurityOverrides     []xrdproto.SecurityOverride
}

// IsManager indicates whether this server has manager role.
func (resp *Response) IsManager() bool {
	return resp.Flags&IsManager != 0
}

// IsServer indicates whether this server has server role.
func (resp *Response) IsServer() bool {
	return resp.Flags&IsServer != 0
}

// IsMeta indicates whether this server has meta attribute.
func (resp *Response) IsMeta() bool {
	return resp.Flags&IsMeta != 0
}

// IsProxy indicates whether this server has proxy attribute.
func (resp *Response) IsProxy() bool {
	return resp.Flags&IsProxy != 0
}

// IsSupervisor indicates whether this server has supervisor attribute.
func (resp *Response) IsSupervisor() bool {
	return resp.Flags&IsSupervisor != 0
}

// ForceSecurity indicates whether signing is required even if the authentication
// protocol does not support generic encryption.
func (resp *Response) ForceSecurity() bool {
	return resp.SecurityOptions&ForceSecurity != 0
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteI32(o.BinaryProtocolVersion)
	wBuffer.WriteI32(int32(o.Flags))
	if !o.HasSecurityInfo {
		return nil
	}
	wBuffer.WriteU8('S')
	wBuffer.Next(1)
	wBuffer.WriteU8(o.SecurityVersion)
	wBuffer.WriteU8(byte(o.SecurityOptions))
	wBuffer.WriteU8(byte(o.SecurityLevel))
	wBuffer.WriteU8(uint8(len(o.SecurityOverrides)))
	for _, x := range o.SecurityOverrides {
		err := x.MarshalXrd(wBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.BinaryProtocolVersion = rBuffer.ReadI32()
	o.Flags = Flags(rBuffer.ReadI32())
	if rBuffer.Len() == 0 {
		return nil
	}
	o.HasSecurityInfo = true
	rBuffer.Skip(1)
	rBuffer.Skip(1)
	o.SecurityVersion = rBuffer.ReadU8()
	o.SecurityOptions = SecurityOptions(rBuffer.ReadU8())
	o.SecurityLevel = xrdproto.SecurityLevel(rBuffer.ReadU8())
	o.SecurityOverrides = make([]xrdproto.SecurityOverride, rBuffer.ReadU8())
	for i := 0; i < len(o.SecurityOverrides); i++ {
		err := o.SecurityOverrides[i].UnmarshalXrd(rBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}

// RespID implements xrdproto.Response.RespID.
func (resp *Response) RespID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }
