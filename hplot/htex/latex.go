// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htex // import "go-hep.org/x/hep/hplot/htex"

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

// Handler is the interface that handles the generation of PDFs
// from TeX, usually via pdflatex.
type Handler interface {
	// CompileLatex compiles the provided .tex document.
	CompileLatex(fname string) error
}

var (
	// DefaultHandler generates PDFs via the pdflatex executable.
	// A LaTeX installation is required, as well as the pdflatex command.
	DefaultHandler = NewHandler("pdflatex")
)

// NoopLatexHandler is a no-op LaTeX compiler.
type NoopHandler struct{}

func (NoopHandler) CompileLatex(fname string) error { return nil }

type pdfLatex struct {
	cmd string
}

// NewHandler returns a Handler compiling .tex documents
// with the provided cmd executable.
func NewHandler(cmd string) Handler {
	return &pdfLatex{cmd: cmd}
}

// CompileLatex compiles the provided .tex document.
func (pdf *pdfLatex) CompileLatex(fname string) error {
	tmp, err := os.MkdirTemp("", "hplot-htex-")
	if err != nil {
		return fmt.Errorf("htex: could not create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	var (
		stdout = new(bytes.Buffer)
		args   = []string{
			fmt.Sprintf("-output-directory=%s", tmp),
			fname,
		}
	)

	cmd := exec.Command(pdf.cmd, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(
			"htex: could not generate PDF from vgtex:\n%s\nerror: %w",
			stdout.Bytes(),
			err,
		)
	}

	oname := fname[:len(fname)-len(".tex")] + ".pdf"
	o, err := os.Create(oname)
	if err != nil {
		return fmt.Errorf("htex: could not create output PDF file: %w", err)
	}
	defer o.Close()

	f, err := os.Open(path.Join(tmp, path.Base(oname)))
	if err != nil {
		return fmt.Errorf("htex: could not open generated PDF file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(o, f)
	if err != nil {
		return fmt.Errorf("htex: could not copy PDF file: %w", err)
	}

	err = o.Close()
	if err != nil {
		return fmt.Errorf("htex: could not close PDF file: %w", err)
	}

	return nil
}

var (
	_ Handler = (*NoopHandler)(nil)
	_ Handler = (*pdfLatex)(nil)
)
