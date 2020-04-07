// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// LatexHandler is the interface that handles the generation of PDFs
// from TeX, usually via pdflatex.
type LatexHandler interface {
	CompileLatex(fname string) error
}

var (
	// DefaultLatexHandler does not generate PDFs
	DefaultLatexHandler = NoopLatexHandler{}
	//
	// PDFLatexHandler generates PDFs via the pdflatex executable.
	// A LaTeX installation is required, as well as the pdflatex command.
	PDFLatexHandler = &pdfLatex{cmd: "pdflatex"}
)

// NoopLatexHandler is a no-op LaTeX compiler.
type NoopLatexHandler struct{}

func (NoopLatexHandler) CompileLatex(fname string) error { return nil }

type pdfLatex struct {
	cmd string
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
