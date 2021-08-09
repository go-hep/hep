// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package genroot // import "go-hep.org/x/hep/groot/internal/genroot"

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// GoFmt formats a file using gofmt.
func GoFmt(f *os.File) {
	fname := f.Name()
	src, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	src, err = format.Source(src)
	if err != nil {
		log.Fatalf("error formating sources of %q: %v\n", fname, err)
	}

	err = os.WriteFile(fname, src, 0644)
	if err != nil {
		log.Fatalf("error writing back %q: %v\n", fname, err)
	}
}

// GenImports adds the provided imports to the given writer.
func GenImports(year int, pkg string, w io.Writer, imports ...string) {
	if year <= 0 {
		year = gopherYear
	}
	fmt.Fprintf(w, srcHeader, year, pkg)
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

// ExtractYear returns the copyright year of a Go-HEP header file.
func ExtractYear(fname string) int {
	raw, err := os.ReadFile(fname)
	if err != nil {
		return gopherYear
	}
	if !bytes.HasPrefix(raw, []byte("// Copyright")) {
		return gopherYear
	}
	raw = bytes.TrimPrefix(raw, []byte("// Copyright ©"))
	raw = bytes.TrimPrefix(raw, []byte("// Copyright "))
	idx := bytes.Index(raw, []byte("The go-hep Authors"))
	raw = bytes.TrimSpace(raw[:idx])

	year, err := strconv.Atoi(string(raw))
	if err != nil {
		return gopherYear
	}
	return year
}

var gopherYear = time.Now().Year()

const srcHeader = `// Copyright ©%d The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package %s

`
