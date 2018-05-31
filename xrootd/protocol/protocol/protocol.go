// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protocol contains the structures describing request and response
// for protocol request (see XRootD specification).
//
// A response consists of 3 parts:
//
// 1) GeneralResponse - general response that is always returned and specifies protocol version and flags describing server type.
//
// 2) SecurityInfo - a response part that is added to the general response
// if `ReturnSecurityRequirements` is provided and server supports it.
// It contains the security version, the security options, the security level,
// and the number of following security overrides, if any.
//
// 3) A list of SecurityOverride - alterations needed to the specified predefined security level.
package protocol // import "go-hep.org/x/hep/xrootd/protocol/protocol"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3006

// GeneralResponseLength is the length of GeneralResponse in bytes.
const GeneralResponseLength = 8

// General response is the response that is always returned from xrootd server.
// It contains protocol version and flags that describe server type.
type GeneralResponse struct {
	BinaryProtocolVersion int32
	Flags                 Flags
}

// Flags are the flags that define xrootd server type. See xrootd protocol specification for further info.
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
	// None specifies that no security options are provided.
	None SecurityOptions = 0
	// ForceSecurity specifies that signing is required even if the authentication
	// protocol does not support generic encryption.
	ForceSecurity SecurityOptions = 0x02
)

// SecurityInfoLength is the length of SecurityInfo in bytes.
const SecurityInfoLength = 6

// SecurityInfo is a response part that is provided when required (if server supports that).
// It contains the security version, the security options, the security level,
// and the number of following security overrides, if any.
type SecurityInfo struct {
	// FIXME: Rename Reserved* fields to _ when automatically generated (un)marshalling will be available.
	Reserved1             byte
	Reserved2             byte
	SecurityVersion       byte
	SecurityOptions       SecurityOptions
	SecurityLevel         SecurityLevel
	SecurityOverridesSize byte
}

// SecurityLevel is the predefined security level that specifies which requests should be signed.
// See specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, p. 75.
type SecurityLevel byte

const (
	// NoneLevel indicates that no request needs to be signed.
	NoneLevel SecurityLevel = 0
	// Compatible indicates that only potentially destructive requests need to be signed.
	Compatible SecurityLevel = 1
	// Standard indicates that potentially destructive requests
	// as well as certain non-destructive requests need to be signed.
	Standard SecurityLevel = 2
	// Intense indicates that request that may reveal metadata or modify data need to be signed.
	Intense SecurityLevel = 3
	// Pedantic indicates that all requests need to be signed.
	Pedantic SecurityLevel = 4
)

// RequestLevel is the security requirement that the associated request is to have.
type RequestLevel byte

const (
	SignNone   RequestLevel = 0 // SignNone indicates that the request need not to be signed.
	SignLikely RequestLevel = 1 // SignLikely indicates that the request must be signed if it modifies data.
	SignNeeded RequestLevel = 2 // SignNeeded indicates that the request mush be signed.
)

// SecurityOverrideLength is the length of SecurityOverride in bytes.
const SecurityOverrideLength = 2

// SecurityOverride is an alteration needed to the specified predefined security level.
// It consists of the request index and the security requirement the associated request should have.
// Request index is calculated as:
//     (request code) - (request code of Auth request)
// according to xrootd protocol specification.
type SecurityOverride struct {
	RequestIndex byte
	RequestLevel RequestLevel
}

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

// ReqID implements protocol.Request.ReqID
func (req *Request) ReqID() uint16 { return RequestID }

// MarshalXrd implements xrootd/protocol.Marshaler
func (o *Request) MarshalXrd() ([]byte, error) {
	var enc xrdenc.Encoder
	enc.WriteI32(o.ClientProtocolVersion)
	enc.WriteU8(byte(o.Options))
	enc.WriteReserved(15)
	return enc.Bytes(), nil
}

// UnmarshalXrd implements xrootd/protocol.Unmarshaler
func (o *Request) UnmarshalXrd(data []byte) error {
	dec := xrdenc.NewDecoder(data)
	o.ClientProtocolVersion = dec.ReadI32()
	o.Options = RequestOptions(dec.ReadU8())
	dec.Skip(15)
	return nil
}
