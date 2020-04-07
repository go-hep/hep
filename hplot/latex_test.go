// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"gonum.org/v1/plot/vg"
)

func TestPDFLatex(t *testing.T) {
	for i, tc := range []struct {
		name  string
		latex LatexHandler
		skip  bool
		want  error
	}{
		{
			name:  "default-latex-handler",
			latex: DefaultLatexHandler,
			want:  nil,
		},
		{
			name:  "pdflatex-not-there",
			latex: &pdfLatex{cmd: "pdflatex-not-there"},
			want:  fmt.Errorf("hplot: could not generate PDF: hplot: could not generate PDF from vgtex:\n\nerror: exec: \"pdflatex-not-there\": executable file not found in $PATH"),
		},
		{
			name:  "pdflatex-handler",
			latex: PDFLatexHandler,
			skip: func() bool {
				_, err := exec.LookPath("pdflatex")
				return err != nil
			}(),
			want: nil,
		},
	} {
		name := fmt.Sprintf("pdflatex-%d", i)
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("skipping %q", tc.name)
			}
			p := New()
			p.X.Min = -10
			p.X.Max = +10
			p.Y.Min = -10
			p.Y.Max = +10

			p.Title.Text = name
			p.X.Label.Text = "X"
			p.Y.Label.Text = "Y"
			p.Latex = tc.latex

			fname := fmt.Sprintf("testdata/%s.tex", name)
			defer os.RemoveAll(fname)

			err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, fname)
			switch {
			case err != nil && tc.want != nil:
				if got, want := err.Error(), tc.want.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %q\nwant=%v", got, want)
				}
			case err != nil && tc.want == nil:
				t.Fatalf("unexpected error: %+v", err)
			case err == nil && tc.want == nil:
				// ok.
			case err == nil && tc.want != nil:
				t.Fatalf("error:\ngot= %v\nwant=%v", err, tc.want)
			}
		})
	}
}
