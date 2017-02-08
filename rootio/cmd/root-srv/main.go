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
	"strings"
	"text/tabwriter"
	"time"

	"go-hep.org/x/hep/rootio"
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
			Path:  strings.Replace(r.URL.Path+"/root-file-upload", "//", "/", -1),
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
		_, err = f.Seek(0, 0)
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
<form id="rootio-form" enctype="multipart/form-data" action={{.Path}} method="post">
      <input id="rootio-file" type="file" name="upload-file" />
      <input type="hidden" name="token" value="{{.Token}}"/>
      <input type="submit" value="upload" />
</form>
<script type="text/javascript">
	document.getElementById("rootio-file").onchange = function() {
		document.getElementById("rootio-form").submit();
	}
</script>
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
	ww := tabwriter.NewWriter(w, 8, 4, 1, ' ', 0)
	for _, k := range f.Keys() {
		obj, err := k.Object()
		if err != nil {
			return "", fmt.Errorf("failed to extract key %q: %v", k.Name(), err)
		}
		switch obj := obj.(type) {
		case rootio.Tree:
			tree := obj
			ww := tabwriter.NewWriter(ww, 8, 4, 1, ' ', 0)
			fmt.Fprintf(ww, "%s\t%s\t%s\t(entries=%d)\n", k.Class(), k.Name(), k.Title(), tree.Entries())
			displayBranches(ww, tree, 2)
			ww.Flush()
		default:
			fmt.Fprintf(ww, "%s\t%s\t%s\t(cycle=%d)\n", k.Class(), k.Name(), k.Title(), k.Cycle())
		}
	}
	ww.Flush()
	return string(w.Bytes()), nil
}

type windent struct {
	hdr []byte
	w   io.Writer
}

func newWindent(n int, w io.Writer) *windent {
	return &windent{
		hdr: bytes.Repeat([]byte(" "), n),
		w:   w,
	}
}

func (w *windent) Write(data []byte) (int, error) {
	return w.w.Write(append(w.hdr, data...))
}

func (w *windent) Flush() error {
	ww, ok := w.w.(flusher)
	if !ok {
		return nil
	}
	return ww.Flush()
}

type flusher interface {
	Flush() error
}

type brancher interface {
	Branches() []rootio.Branch
}

func displayBranches(w io.Writer, bres brancher, indent int) {
	branches := bres.Branches()
	if len(branches) <= 0 {
		return
	}
	ww := newWindent(indent, w)
	for _, b := range branches {
		fmt.Fprintf(ww, "%s\t%q\t%v\n", b.Name(), b.Title(), b.Class())
		displayBranches(ww, b, 2)
	}
	ww.Flush()
}
