// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-hep/rootio"
)

var (
	addrFlag = flag.String("addr", ":8080", "server address:port")
)

func main() {
	flag.Parse()

	http.Handle("/", appHandler(rootHandle))
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.New("upload").Parse(page)
		if err != nil {
			return err
		}

		var data = struct {
			Token string
			Path  string
		}{
			Token: token,
			Path:  r.URL.Path + "root-file-upload",
		}

		err = t.Execute(w, data)
		if err != nil {
			return err
		}

	case "POST":
		r.ParseMultipartForm(500 << 20)
		f, handler, err := r.FormFile("upload-file")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		out, err := inspect(f, handler.Filename)
		if err != nil {
			return err
		}

		fmt.Fprintf(w, out)

	default:
		return fmt.Errorf("invalid request %q", r.Method)
	}

	return nil
}

const page = `<html>
<head>
    <title>go-hep/rootio file inspector</title>
</head>
<body>
<h2>go-hep/rootio ROOT file inspector</h2>
<form enctype="multipart/form-data" action={{.Path}} method="post">
      <input type="file" name="upload-file" />
      <input type="hidden" name="token" value="{{.Token}}"/>
      <input type="submit" value="upload" />
</form>
</body>
</html>
`

func inspect(r rootio.Reader, fname string) (string, error) {
	f, err := rootio.NewReader(r, fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "=== inspecting file %q...\n", fname)
	fmt.Fprintf(w, "version: %v\n", f.Version())
	for _, k := range f.Keys() {
		obj, err := k.Object()
		if err != nil {
			return "", fmt.Errorf("failed to extract key %q: %v", k.Name(), err)
		}
		switch obj := obj.(type) {
		case rootio.Tree:
			tree := obj
			fmt.Fprintf(w, "%-8s %-40s %s (entries=%d)\n", k.Class(), k.Name(), k.Title(), tree.Entries())
			for _, b := range tree.Branches() {
				fmt.Fprintf(w, "  %-20s %-20q %v\n", b.Name(), b.Title(), b.Class())
			}
		default:
			fmt.Fprintf(w, "%-8s %-40s %s (cycle=%d)\n", k.Class(), k.Name(), k.Title(), k.Cycle())
		}
	}

	return string(w.Bytes()), nil
}
