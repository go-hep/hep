// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htex_test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/htex"
	"gonum.org/v1/plot/vg"
)

func TestGoHandler(t *testing.T) {
	tmp, err := os.MkdirTemp("", "hplot-htex-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	for i, tc := range []struct {
		name  string
		latex *htex.GoHandler
		skip  bool
		want  error
	}{
		{
			name:  "parallel-pdflatex-not-there",
			latex: htex.NewGoHandler(-1, "pdflatex-not-there"),
			want: func() error {
				err := fmt.Errorf("htex: could not generate PDF from vgtex:\n\nerror: exec: \"pdflatex-not-there\": executable file not found in $PATH")
				if runtime.GOOS != "windows" {
					return err
				}
				return fmt.Errorf("htex: could not generate PDF from vgtex:\n\nerror: exec: \"pdflatex-not-there\": executable file not found in %%PATH%%")
			}(),
		},
		{
			name:  "parallel-pdflatex-handler",
			latex: htex.NewGoHandler(-1, "pdflatex"),
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
			p := hplot.New()
			p.X.Min = -10
			p.X.Max = +10
			p.Y.Min = -10
			p.Y.Max = +10

			p.Title.Text = name
			p.X.Label.Text = "X"
			p.Y.Label.Text = "Y"

			fig := hplot.Figure(p, hplot.WithLatexHandler(tc.latex))

			for i := range 10 {
				fname := fmt.Sprintf("%s/%s-%02d.tex", tmp, name, i)
				defer os.RemoveAll(fname)

				err := hplot.Save(fig, 10*vg.Centimeter, 10*vg.Centimeter, fname)
				if err != nil {
					t.Fatalf("error: %+v", err)
				}
			}

			err := tc.latex.Wait()
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
