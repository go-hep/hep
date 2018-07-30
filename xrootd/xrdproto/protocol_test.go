// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdproto // import "go-hep.org/x/hep/xrootd/xrdproto"

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"reflect"
	"testing"
)

func TestReadRequest(t *testing.T) {
	header := make([]byte, RequestHeaderLength+16+4)
	data := make([]byte, 10)
	rand.Read(data)
	binary.BigEndian.PutUint32(header[RequestHeaderLength+16:], 10)

	for _, tc := range []struct {
		name string
		data []byte
		want []byte
		err  error
	}{
		{
			name: "EOF",
			err:  io.EOF,
			data: []byte{},
		},
		{
			name: "Without data",
			data: make([]byte, RequestHeaderLength+16+4),
			want: make([]byte, RequestHeaderLength+16+4),
		},
		{
			name: "With data",
			data: append(header, data...),
			want: append(header, data...),
		},
		{
			name: "Header with non-zero length but without data",
			err:  io.EOF,
			data: header,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewBuffer(tc.data)
			got, err := ReadRequest(reader)
			if err != tc.err {
				t.Errorf("error doesn't match:\ngot = %v\nwant = %v", err, tc.err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("data doesn't match:\ngot = %v\nwant = %v", got, tc.want)
			}
			if reader.Len() != 0 {
				t.Errorf("reader was not read to the end: %v", reader.Bytes())
			}
		})
	}
}

func TestReadResponse(t *testing.T) {
	header := make([]byte, RequestHeaderLength+16+4)
	data := make([]byte, 10)
	rand.Read(data)
	binary.BigEndian.PutUint32(header[RequestHeaderLength+16:], 10)

	for _, tc := range []struct {
		name       string
		data       []byte
		wantHeader ResponseHeader
		wantData   []byte
		err        error
	}{
		{
			name: "EOF",
			err:  io.EOF,
			data: []byte{},
		},
		{
			name:       "Without data",
			data:       []byte{1, 2, 0, 0, 0, 0, 0, 0},
			wantHeader: ResponseHeader{StreamID: StreamID{1, 2}},
		},
		{
			name:       "With data",
			data:       []byte{1, 2, 0, 0, 0, 0, 0, 5, 1, 2, 3, 4, 5},
			wantHeader: ResponseHeader{StreamID: StreamID{1, 2}, DataLength: 5},
			wantData:   []byte{1, 2, 3, 4, 5},
		},
		{
			name:       "Header with non-zero length but without data",
			err:        io.EOF,
			data:       []byte{1, 2, 0, 0, 0, 0, 0, 5},
			wantHeader: ResponseHeader{StreamID: StreamID{1, 2}, DataLength: 5},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewBuffer(tc.data)
			gotHeader, gotData, err := ReadResponse(reader)
			if err != tc.err {
				t.Errorf("error doesn't match:\ngot = %v\nwant = %v", err, tc.err)
			}
			if !reflect.DeepEqual(gotHeader, tc.wantHeader) {
				t.Errorf("header doesn't match:\ngot = %v\nwant = %v", gotHeader, tc.wantHeader)
			}
			if !reflect.DeepEqual(gotData, tc.wantData) {
				t.Errorf("data doesn't match:\ngot = %v\nwant = %v", gotData, tc.wantData)
			}
			if reader.Len() != 0 {
				t.Errorf("reader was not read to the end: %v", reader.Bytes())
			}
		})
	}
}

func TestWriteResponse(t *testing.T) {
	header := make([]byte, RequestHeaderLength+16+4)
	data := make([]byte, 10)
	rand.Read(data)
	binary.BigEndian.PutUint32(header[RequestHeaderLength+16:], 10)

	for _, tc := range []struct {
		name     string
		wantData []byte
		header   ResponseHeader
		err      error
		streamID StreamID
		status   ResponseStatus
		resp     Marshaler
	}{
		{
			name:     "With data",
			wantData: []byte{1, 2, 15, 163, 0, 0, 0, 5, 0, 0, 0, 12, 0},
			status:   Error,
			streamID: StreamID{1, 2},
			resp:     &ServerError{Code: 12},
		},
		{
			name:     "Without data",
			wantData: []byte{1, 2, 0, 0, 0, 0, 0, 0},
			status:   Ok,
			streamID: StreamID{1, 2},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var writer bytes.Buffer
			err := WriteResponse(&writer, tc.streamID, tc.status, tc.resp)
			if err != tc.err {
				t.Errorf("error doesn't match:\ngot = %v\nwant = %v", err, tc.err)
			}
			if !reflect.DeepEqual(writer.Bytes(), tc.wantData) {
				t.Errorf("data doesn't match:\ngot = %v\nwant = %v", writer.Bytes(), tc.wantData)
			}
		})
	}
}
