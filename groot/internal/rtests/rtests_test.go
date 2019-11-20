// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtests_test

import (
	"bytes"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"golang.org/x/xerrors"
)

func TestRunCxxROOTWithoutROOT(t *testing.T) {
	hasROOT := rtests.HasROOT
	rtests.HasROOT = false
	defer func() {
		rtests.HasROOT = hasROOT
	}()

	out, err := rtests.RunCxxROOT("hello", []byte(`void hello(const char* name) { std::cout << name << std::endl; }`), "hello")
	if !xerrors.Is(err, rtests.ErrNoROOT) {
		t.Fatalf("unexpected error: got=%v, want=%v\noutput:\n%s", err, rtests.ErrNoROOT, out)
	}
}

func TestRunCxxROOTInvalidMacro(t *testing.T) {
	out, err := rtests.RunCxxROOT("hello", []byte(`void hello(const char* name) { std::cout << nameXXX << std::endl; }`), "hello")
	if err == nil {
		t.Fatalf("expected C++ ROOT macro to fail")
	}
	if !rtests.HasROOT {
		return
	}
	var dst rtests.ROOTError
	if !xerrors.As(err, &dst) {
		t.Fatalf("unexpected error-type (%T): got=%+v", err, err)
	}
	const suffix = `hello.C:1:45: error: use of undeclared identifier 'nameXXX'
void hello(const char* name) { std::cout << nameXXX << std::endl; }
                                            ^
`
	if !bytes.HasSuffix(out, []byte(suffix)) {
		t.Fatalf("unexpected error: got=%+v\noutput:\n%s", err, out)
	}
}

func TestRunCxxROOT(t *testing.T) {
	out, err := rtests.RunCxxROOT("hello", []byte(`void hello(const char* name, int d) { std::cout << name << "-" << d << std::endl; }`), "hello", 42)
	if err != nil {
		switch {
		case rtests.HasROOT:
			t.Fatalf("expected C++ ROOT macro to run correctly: %+v\noutput:\n%s", err, out)
		default:
			if !xerrors.Is(err, rtests.ErrNoROOT) {
				t.Fatalf("unexpected error: got=%v, want=%v", err, rtests.ErrNoROOT)
			}
		}
		return
	}

	// ROOT macros start with printing out:
	// \nProcessing /tmp/groot-rtests-516158679/hello.C(\"hello\")...\n
	if i := bytes.Index(out, []byte("...\n")); i > 0 {
		out = out[i+len([]byte("...\n")):]
	}
	if got, want := string(out), string("hello-42\n"); got != want {
		t.Fatalf("invalid ROOT macro result. got=%q, want=%q", got, want)
	}
}
