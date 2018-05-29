// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protocol contains the XRootD protocol specific types
// and methods to handle them, such as marshalling and unmarshalling requests.
package protocol // import "go-hep.org/x/hep/xrootd/protocol"

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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

// EntryStat holds the entry name and the entry stat information.
type EntryStat struct {
	Name          string // Name is the name of entry.
	HasStatInfo   bool   // HasStatInfo indicates if the following stat information is valid.
	Id            int64  // Id is the OS-dependent identifier assigned to this entry.
	Size          uint64 // Size is the decimal size of the entry.
	IsExecutable  bool   // IsExecutable indicates that entry is either an executable file or a searchable directory.
	IsDir         bool   // IsDir indicates that entry is a directory.
	IsOther       bool   // IsOther indicates that entry is neither a file nor a directory.
	IsOffline     bool   // IsOffline indicates that the file is not online (i. e., on disk).
	IsReadable    bool   // IsReadable indicates that read access to that entry is allowed.
	IsWritable    bool   // IsWritable indicates that write access to that entry is allowed.
	IsPoscPending bool   // IsPoscPending indicates that the file was created with kXR_posc and has not yet been successfully closed.
	Mtime         int64  // Mtime is the last modification time in Unix time units.
}

// NewEntryStat creates a new EntryStat according to the name of the entry
// and statinfo following XRootD protocol specification (see page 114).
func NewEntryStat(name, statinfo string) (EntryStat, error) {
	stats := strings.Split(statinfo, " ")

	if len(stats) < 4 {
		return EntryStat{}, errors.Errorf("xrootd: statinfo \"%s\" doesn't have enough fields, expected format is: \"id size flags modtime\"", statinfo)
	}

	id, err := strconv.Atoi(stats[0])
	if err != nil {
		return EntryStat{}, err
	}
	size, err := strconv.Atoi(stats[1])
	if err != nil {
		return EntryStat{}, err
	}
	flags, err := strconv.Atoi(stats[2])
	if err != nil {
		return EntryStat{}, err
	}
	mtime, err := strconv.Atoi(stats[3])
	if err != nil {
		return EntryStat{}, err
	}

	result := EntryStat{
		Name:          name,
		HasStatInfo:   true,
		Id:            int64(id),
		Size:          uint64(size),
		Mtime:         int64(mtime),
		IsExecutable:  flags&1 > 0,
		IsDir:         flags&2 > 0,
		IsOther:       flags&4 > 0,
		IsOffline:     flags&8 > 0,
		IsReadable:    flags&16 > 0,
		IsWritable:    flags&32 > 0,
		IsPoscPending: flags&64 > 0,
	}

	return result, nil
}
