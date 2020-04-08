// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htex // import "go-hep.org/x/hep/hplot/htex"

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// Handler is the interface that handles the generation of PDFs
// from TeX, usually via pdflatex.
type Handler interface {
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

func (pdf *pdfLatex) CompileLatex(fname string) error {
	var (
		stdout = new(bytes.Buffer)
		args   = []string{
			fmt.Sprintf("-output-directory=%s", filepath.Dir(fname)),
			fname,
		}
	)

	cmd := exec.Command(pdf.cmd, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(
			"hplot: could not generate PDF from vgtex:\n%s\nerror: %w",
			stdout.Bytes(),
			err,
		)
	}

	return nil
}

var (
	_ Handler = (*NoopHandler)(nil)
	_ Handler = (*pdfLatex)(nil)
)
