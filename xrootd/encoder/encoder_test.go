// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoder // import "go-hep.org/x/hep/xrootd/encoder"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-hep.org/x/hep/xrootd/protocol"
)

type request struct {
	Z int64
	X int32
	A uint8
	C uint16
	D [2]byte
	Y []byte
}

type benchmarkRequest struct {
	X   int32
	A   uint8
	A1  uint8
	A2  uint8
	A3  uint8
	A4  int32
	A5  int32
	A6  int32
	A7  int32
	A8  int32
	A9  int32
	A10 int32
	A11 int32
	A12 int32
	A13 int32
	A14 int32
	A16 int32
	C   uint16
	D   [10]byte
	Z   [10]byte
}

type undecodable struct {
	A float64
}

func TestMarshalRequest(t *testing.T) {
	var requestID uint16 = 1337
	var streamID = protocol.StreamID{42, 37}
	expected := []byte{42, 37, 5, 57, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 1, 2, 0, 3, 6, 7, 11, 13}

	actual, err := MarshalRequest(requestID, streamID, request{7, 1, 2, 3, protocol.StreamID{6, 7}, []byte{11, 13}})

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnmarshal(t *testing.T) {
	var expected = request{7, 1, 2, 3, protocol.StreamID{6, 7}, []byte{11, 13}}

	var actual = &request{}
	err := Unmarshal([]byte{0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 1, 2, 0, 3, 6, 7, 11, 13}, actual)

	assert.Equal(t, expected, *actual)
	assert.NoError(t, err)
}

func TestMarshalRequest_Undecodable(t *testing.T) {
	var requestID uint16 = 1337
	var streamID = protocol.StreamID{42, 37}

	_, err := MarshalRequest(requestID, streamID, undecodable{1})

	assert.Error(t, err)
}

func TestUnmarshal_Undecodable(t *testing.T) {
	var actual = &undecodable{}

	err := Unmarshal([]byte{0, 0, 0, 1, 2, 0, 3, 6, 7, 11, 13}, actual)

	assert.Error(t, err)
}

func BenchmarkMarshal(b *testing.B) {
	br := benchmarkRequest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(br); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkUnMarshal(b *testing.B) {
	br := benchmarkRequest{}
	data := make([]byte, 78)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &br); err != nil {
			b.Error(err)
		}
	}
}
