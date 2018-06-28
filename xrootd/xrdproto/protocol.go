// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protocol contains the XRootD protocol specific types
// and methods to handle them, such as marshalling and unmarshalling requests.
package xrdproto // import "go-hep.org/x/hep/xrootd/xrdproto"

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// ResponseStatus is the status code indicating how the request completed.
type ResponseStatus uint16

const (
	// Ok indicates that request fully completed and no addition responses will be forthcoming.
	Ok ResponseStatus = 0
	// OkSoFar indicates that server provides partial response and client should be prepared
	// to receive additional responses on same stream.
	OkSoFar ResponseStatus = 4000
	// Error indicates that an error occurred during request handling.
	// Error code and error message are sent as part of response (see xrootd protocol specification v3.1.0, p. 27).
	Error ResponseStatus = 4003
	// Redirect indicates that the client must re-issue the request to another server.
	Redirect ResponseStatus = 4004
	// Wait indicates that the client must wait the indicated number of seconds and retry the request.
	Wait ResponseStatus = 4005
)

// ServerError is the error returned by the XRootD server as part of response to the request.
type ServerError struct {
	Code    ServerErrorCode
	Message string
}

// ServerErrorCode is the code of the error returned by the XRootD server as part of response to the request.
type ServerErrorCode int32

const (
	InvalidRequestCode ServerErrorCode = 3006 // InvalidRequestCode indicates that request is invalid.
	NotAuthorized      ServerErrorCode = 3010 // NotAuthorized indicates that user was not authorized for operation.
	NotFoundCode       ServerErrorCode = 3011 // NotFoundCode indicates that path was not found on the remote server.
)

func (err ServerError) Error() string {
	return fmt.Sprintf("xrootd: error %d: %s", err.Code, err.Message)
}

// StreamID is the binary identifier associated with a request stream.
type StreamID [2]byte

// ResponseHeaderLength is the length of the ResponseHeader in bytes.
const ResponseHeaderLength = 2 + 2 + 4

// ResponseHeader is the header that precedes all responses (see xrootd protocol specification).
type ResponseHeader struct {
	StreamID   StreamID
	Status     ResponseStatus
	DataLength int32
}

// MarshalXrd implements xrdproto.Marshaler
func (o ResponseHeader) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.StreamID[:])
	wBuffer.WriteU16(uint16(o.Status))
	wBuffer.WriteI32(o.DataLength)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
func (o *ResponseHeader) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.StreamID[:])
	o.Status = ResponseStatus(rBuffer.ReadU16())
	o.DataLength = rBuffer.ReadI32()
	return nil
}

// RequestHeaderLength is the length of the RequestHeader in bytes.
const RequestHeaderLength = 2 + 2

// ResponseHeader is the header that precedes all requests (we are interested in StreamID and RequestID, actual request
// parameters are a part of specific request).
type RequestHeader struct {
	StreamID  StreamID
	RequestID uint16
}

// MarshalXrd implements Marshaler.
func (o RequestHeader) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.StreamID[:])
	wBuffer.WriteU16(o.RequestID)
	return nil
}

// UnmarshalXrd implements Unmarshaler.
func (o *RequestHeader) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.StreamID[:])
	o.RequestID = rBuffer.ReadU16()
	return nil
}

// Error returns an error received from the server or nil if request hasn't failed.
func (hdr ResponseHeader) Error(data []byte) error {
	if hdr.Status == Error {
		// 4 bytes for error code and at least 1 byte for message (in case it is null-terminated empty string)
		if len(data) < 5 {
			return errors.New("xrootd: an server error occurred, but code and message were not provided")
		}
		code := ServerErrorCode(binary.BigEndian.Uint32(data[0:4]))
		message := string(data[4 : len(data)-1]) // Skip \0 character at the end

		return ServerError{code, message}
	}
	return nil
}

// ServerType is the general server type kept for compatibility
// with 2.0 protocol version (see xrootd protocol specification v3.1.0, p. 5).
type ServerType int32

const (
	// LoadBalancingServer indicates whether this is a load-balancing server.
	LoadBalancingServer ServerType = iota
	// DataServer indicates whether this is a data server.
	DataServer
)

// FilepathRequest is a request that contains file paths.
// This interface is used to append opaque data to the request.
// Opaque data is received as part of the redirect response.
type FilepathRequest interface {
	Opaque() string          // Opaque returns opaque data from this request.
	SetOpaque(opaque string) // SetOpaque sets opaque data for this request.
}

// RequestLevel is the security requirement that the associated request is to have.
type RequestLevel byte

const (
	SignNone   RequestLevel = 0 // SignNone indicates that the request need not to be signed.
	SignLikely RequestLevel = 1 // SignLikely indicates that the request must be signed if it modifies data.
	SignNeeded RequestLevel = 2 // SignNeeded indicates that the request mush be signed.
)

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

// MarshalXrd implements xrdproto.Marshaler
func (o SecurityOverride) MarshalXrd(enc *xrdenc.WBuffer) error {
	enc.WriteU8(o.RequestIndex)
	enc.WriteU8(byte(o.RequestLevel))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
func (o *SecurityOverride) UnmarshalXrd(dec *xrdenc.RBuffer) error {
	o.RequestIndex = dec.ReadU8()
	o.RequestLevel = RequestLevel(dec.ReadU8())
	return nil
}

// SetOpaque sets opaque data part in the provided path.
func SetOpaque(path *string, opaque string) {
	pos := strings.LastIndex(*path, "?")
	if pos != -1 {
		*path = (*path)[:pos]
	}
	*path = *path + "?" + opaque
}

// Opaque returns opaque data from provided path.
func Opaque(path string) string {
	pos := strings.LastIndex(path, "?")
	return path[pos+1:]
}
