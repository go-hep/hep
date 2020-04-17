// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package genroot

import (
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func GoFmt(f *os.File) {
	fname := f.Name()
	src, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	src, err = format.Source(src)
	if err != nil {
		log.Fatalf("error formating sources of %q: %v\n", fname, err)
	}

	err = ioutil.WriteFile(fname, src, 0644)
	if err != nil {
		log.Fatalf("error writing back %q: %v\n", fname, err)
	}
}

func GenImports(pkg string, w io.Writer, imports ...string) {
	fmt.Fprintf(w, srcHeader, pkg)
	if len(imports) == 0 {
		return
	}

	fmt.Fprintf(w, "import (\n")
	for _, imp := range imports {
		if imp == "" {
			fmt.Fprintf(w, "\n")
			continue
		}
		fmt.Fprintf(w, "\t%q\n", imp)
	}
	fmt.Fprintf(w, ")\n\n")
}

const srcHeader = `// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package %s

`
