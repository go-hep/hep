// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protocol contains the XRootD protocol specific types
// and methods to handle them, such as marshalling and unmarshalling requests.
package xrdproto // import "go-hep.org/x/hep/xrootd/xrdproto"

import (
	"encoding/binary"
	"fmt"

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
)

// ServerError is the error returned by the XRootD server as part of response to the request.
type ServerError struct {
	Code    int32
	Message string
}

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

// Error returns an error received from the server or nil if request hasn't failed.
func (hdr ResponseHeader) Error(data []byte) error {
	if hdr.Status == Error {
		// 4 bytes for error code and at least 1 byte for message (in case it is null-terminated empty string)
		if len(data) < 5 {
			return errors.New("xrootd: an server error occurred, but code and message were not provided")
		}
		code := int32(binary.BigEndian.Uint32(data[0:4]))
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
