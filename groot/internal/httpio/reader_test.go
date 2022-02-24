// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/iotest"
)

func TestOpen(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("./testdata")))

	for _, tc := range []struct {
		url  string
		opts []Option
		err  error
	}{
		{
			url: srv.URL + "/data.txt",
			opts: []Option{
				func(c *config) error { return fmt.Errorf("option error") },
			},
			err: fmt.Errorf("httpio: could not open %q: %w", srv.URL+"/data.txt", fmt.Errorf("option error")),
		},
		{
			url: srv.URL + "xxx/not-there.txt",
			err: fmt.Errorf("httpio: could not create HTTP request: "),
		},
		{
			url: srv.URL + "000/not-there.txt",
			err: fmt.Errorf("httpio: could not send HEAD request: "),
		},
		{
			url: srv.URL + "/not-there.txt",
			err: fmt.Errorf("httpio: invalid HEAD response code=404"),
		},
	} {
		t.Run(tc.url, func(t *testing.T) {
			_, err := Open(tc.url, tc.opts...)
			switch {
			case err == nil:
				t.Fatalf("expected an error")
			default:
				if got, want := err.Error(), tc.err.Error(); !strings.HasPrefix(got, want) {
					t.Fatalf("invalid error:\ngot= %s\nwant=%s", got, want)
				}
			}
		})
	}

	t.Run("no-range", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("accept-range", "none")
		}))

		want := fmt.Errorf("httpio: invalid HEAD response: httpio: accept-range not supported")

		_, err := Open(srv.URL)
		if got, want := err.Error(), want.Error(); got != want {
			t.Fatalf("invalid error:\ngot= %s\nwant=%s", got, want)
		}
	})
}

func TestReader(t *testing.T) {
	dir := "./testdata"
	srv := httptest.NewServer(http.FileServer(http.Dir(dir)))

	want, err := os.ReadFile("./testdata/data.txt")
	if err != nil {
		t.Fatal(err)
	}

	r, err := Open(
		srv.URL+"/data.txt",
		WithBasicAuth("gopher", "s3cr3t"),
		WithContext(context.Background()),
		WithClient(http.DefaultClient),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	if got, want := r.Name(), srv.URL+"/data.txt"; got != want {
		t.Fatalf("invalid name:\ngot= %q\nwant=%q", got, want)
	}

	if got, want := r.Size(), int64(len(want)); got != want {
		t.Fatalf("invalid size: got=%d, want=%d", got, want)
	}

	got := make([]byte, len(want))
	n, err := r.ReadAt(got, 0)
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("could not read-at: %+v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("invalid read-at:\ngot= %q\nwant=%q", got, want)
	}
	if len(want) != n {
		t.Fatalf("invalid nbytes: got=%d, want=%d", n, len(want))
	}

	t.Run("iotest", func(t *testing.T) {
		r, err := Open(srv.URL + "/data.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		err = iotest.TestReader(r, want)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("etag", func(t *testing.T) {
		r, err := Open(srv.URL + "/data.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		// fake changing resource
		r.etag += "XXX"

		p := make([]byte, len(want))
		n, err := r.ReadAt(p, 0)
		switch {
		case err == nil:
			t.Fatalf("expected an Etag-related error")
		default:
			if got, want := err.Error(), "httpio: resource changed"; got != want {
				t.Fatalf("invalid error:\ngot= %q\nwant=%q", got, want)
			}
		}
		if n != len(want) {
			t.Fatalf("invalid size: got=%d, want=%d", n, len(want))
		}
	})
}
