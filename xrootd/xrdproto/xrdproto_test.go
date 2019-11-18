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
	"time"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"golang.org/x/xerrors"
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

func TestWaitResponse(t *testing.T) {
	for _, want := range []WaitResponse{
		{Duration: 0},
		{Duration: 42 * time.Second},
		{Duration: 42 * time.Hour},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got WaitResponse
			)

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatalf("could not marshal response: %v", err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatalf("could not unmarshal response: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed\ngot = %#v\nwant= %#v\n", got, want)
			}
		})
	}
}

func TestServerError(t *testing.T) {
	for _, want := range []ServerError{
		{Code: IOError, Message: ""},
		{Code: NotAuthorized, Message: "not authorized"},
		{Code: NotFound, Message: "not\nfound"},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got ServerError
			)

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatalf("could not marshal server error: %v", err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatalf("could not unmarshal server error: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed\ngot = %#v\nwant= %#v\n", got, want)
			}

			if got, want := got.Error(), want.Error(); got != want {
				t.Fatalf("error messages differ: got=%q, want=%q", got, want)
			}
		})
	}
}

func TestRequestHeader(t *testing.T) {
	for _, want := range []RequestHeader{
		{StreamID: StreamID{1, 2}, RequestID: 2},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got RequestHeader
			)

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatalf("could not marshal: %v", err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatalf("could not unmarshal: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed\ngot = %#v\nwant= %#v\n", got, want)
			}
		})
	}
}

func TestResponseHeaderError(t *testing.T) {
	get := func(err error) string {
		if err != nil {
			return err.Error()
		}
		return ""
	}

	for _, tc := range []struct {
		hdr  ResponseHeader
		data []byte
		err  error
	}{
		{
			hdr:  ResponseHeader{Status: Ok},
			data: nil,
			err:  nil,
		},
		{
			hdr:  ResponseHeader{Status: OkSoFar},
			data: nil,
			err:  nil,
		},
		{
			hdr: ResponseHeader{Status: Error},
			data: func() []byte {
				w := new(xrdenc.WBuffer)
				err := ServerError{Code: IOError, Message: "boo"}.MarshalXrd(w)
				if err != nil {
					t.Fatal(err)
				}
				return w.Bytes()
			}(),
			err: ServerError{Code: IOError, Message: "boo"},
		},
		{
			hdr:  ResponseHeader{Status: Error},
			data: []byte{1, 2, 3},
			err:  xerrors.Errorf("xrootd: invalid ResponseHeader error: %w", io.ErrShortBuffer),
		},
		{
			hdr:  ResponseHeader{Status: Error},
			data: []byte{1, 2, 3, 4},
			err:  xerrors.Errorf("xrootd: error occurred during unmarshaling of a server error: xrootd: missing error message in server response"),
		},
	} {
		t.Run("", func(t *testing.T) {
			err := tc.hdr.Error(tc.data)
			if get(err) != get(tc.err) {
				t.Fatalf("got=%#v, want=%#v", err, tc.err)
			}
		})
	}
}

func TestSecurityOverride(t *testing.T) {
	for _, want := range []SecurityOverride{
		{RequestIndex: 1, RequestLevel: SignNone},
		{RequestIndex: 2, RequestLevel: SignLikely},
		{RequestIndex: 3, RequestLevel: SignNeeded},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got SecurityOverride
			)

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatalf("could not marshal: %v", err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatalf("could not unmarshal: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed\ngot = %#v\nwant= %#v\n", got, want)
			}
		})
	}
}

func TestOpaque(t *testing.T) {
	for _, tc := range []struct {
		path string
		want string
	}{
		{"hello", "hello"},
		{"hello?", ""},
		{"hello?boo", "boo"},
		{"?boo", "boo"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			got := Opaque(tc.path)
			if got != tc.want {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}

func TestSetOpaque(t *testing.T) {
	for _, tc := range []struct {
		path string
		opaq string
		want string
	}{
		{"", "v", "?v"},
		{"hello", "v", "hello?v"},
		{"hello?", "v", "hello?v"},
		{"hello?boo", "v", "hello?v"},
		{"?boo", "v", "?v"},
		{"hello?boo?", "v", "hello?boo?v"},
		{"hello?boo?bar", "v", "hello?boo?v"},
		{"hello?boo=value?bar=33", "v", "hello?boo=value?v"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			got := tc.path
			SetOpaque(&got, tc.opaq)
			if got != tc.want {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}
