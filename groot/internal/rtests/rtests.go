// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtests

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
)

type ROOTer interface {
	root.Object
	rbytes.Marshaler
	rbytes.Unmarshaler
}

func XrdRemote(fname string) string {
	const remote = "root://ccxrootdgotest.in2p3.fr:9001/tmp/rootio"
	return remote + "/" + fname
}

var (
	HasROOT   = false // HasROOT is true when a C++ ROOT installation could be detected.
	ErrNoROOT = errors.New("rtests: no C++ ROOT installed")
	rootCmd   = ""
)

// RunCxxROOT executes the function fct in the provided C++ code with optional arguments args.
// RunCxxROOT creates a temporary file named '<fct>.C' from the provided C++ code and
// executes it via ROOT C++.
// RunCxxROOT returns the combined stdout/stderr output and an error, if any.
func RunCxxROOT(fct string, code []byte, args ...interface{}) ([]byte, error) {
	tmp, err := ioutil.TempDir("", "groot-rtests-")
	if err != nil {
		return nil, fmt.Errorf("could not create tmpdir: %w", err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, fct+".C")
	err = ioutil.WriteFile(fname, []byte(code), 0644)
	if err != nil {
		return nil, fmt.Errorf("could not generate ROOT macro %q: %w", fname, err)
	}

	o := new(strings.Builder)
	fmt.Fprintf(o, "%s(", fname)
	for i, arg := range args {
		format := ""
		if i > 0 {
			format = ", "
		}
		switch arg.(type) {
		case string:
			format += "%q"
		default:
			format += "%v"
		}
		fmt.Fprintf(o, format, arg)
	}
	fmt.Fprintf(o, ")")

	if !HasROOT {
		return nil, ErrNoROOT
	}

	cmd := exec.Command(rootCmd, "-l", "-b", "-x", "-q", o.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, ROOTError{Err: err, Cmd: cmd.Path, Args: cmd.Args, Out: out}
	}

	return out, nil
}

type ROOTError struct {
	Err  error
	Cmd  string
	Args []string
	Out  []byte
}

func (err ROOTError) Error() string {
	return fmt.Sprintf(
		"could not run '%s': %v\noutput:\n%s",
		strings.Join(append([]string{err.Cmd}, err.Args...), " "),
		err.Err,
		err.Out,
	)
}

func (err ROOTError) Unwrap() error {
	return err.Err
}

func init() {
	cmd, err := exec.LookPath("root.exe")
	if err != nil {
		return
	}
	HasROOT = true
	rootCmd = cmd
}

var (
	_ error                       = (*ROOTError)(nil)
	_ interface{ Unwrap() error } = (*ROOTError)(nil)
)
