// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htex_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/htex"
)

func ExampleGoHandler() {
	hdlr := htex.NewGoHandler(-1, "pdflatex")

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("plot-%0d", i)
		p := hplot.New()
		p.Title.Text = name
		p.X.Label.Text = "x"
		p.Y.Label.Text = "y"

		err := hplot.Save(
			hplot.Wrap(p, hplot.WithLatexHandler(hdlr)),
			-1, -1, name+".tex",
		)
		if err != nil {
			log.Fatalf("could not save plot: %+v", err)
		}
	}

	err := hdlr.Wait()
	if err != nil {
		log.Fatalf("error compiling latex: %+v", err)
	}
}
