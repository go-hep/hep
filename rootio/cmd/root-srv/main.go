// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"image/color"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/rootio"
)

var (
	addrFlag = flag.String("addr", ":8080", "server address:port")

	db = dbFiles{
		files: make(map[string]*rootio.File),
	}
)

func main() {
	flag.Parse()

	http.Handle("/", appHandler(rootHandle))
	http.Handle("/root-file-upload", appHandler(uploadHandle))
	http.Handle("/plot/", appHandler(plotH1Handle))
	http.Handle("/plot2d/", appHandler(plotH2Handle))
	log.Printf("server listening on %s", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("error %q: %v\n", r.URL.Path, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		// ok
	default:
		return fmt.Errorf("invalid request %q for /", r.Method)
	}

	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	t, err := template.New("upload").Parse(page)
	if err != nil {
		return err
	}

	return t.Execute(w, struct{ Token string }{token})
}

func uploadHandle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		log.Printf("invalid request %q", r.Method)
		return fmt.Errorf("invalid request %q for /root-file-upload", r.Method)
	}

	r.ParseMultipartForm(500 << 20)
	f, handler, err := r.FormFile("upload-file")
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	rfile, err := rootio.NewReader(f, handler.Filename)
	if err != nil {
		return err
	}
	db.set(session(r, handler.Filename), rfile)

	var nodes []jsNode
	db.RLock()
	defer db.RUnlock()

	for k, rfile := range db.files {
		node, err := fileJsTree(rfile, k)
		if err != nil {
			return err
		}
		nodes = append(nodes, node...)
	}
	sort.Sort(jsNodes(nodes))
	return json.NewEncoder(w).Encode(nodes)
}

func plotH2Handle(w http.ResponseWriter, r *http.Request) error {
	log.Printf(">>> request: %q\n", r.URL.Path)

	return json.NewEncoder(w).Encode(map[string]string{
		"url": r.URL.Path,
	})
}

func plotH1Handle(w http.ResponseWriter, r *http.Request) error {
	log.Printf(">>> request: %q\n", r.URL.Path)
	url := r.URL.Path[len("/plot/"):]
	toks := strings.Split(url, "/")
	fname := toks[0]
	f := db.get(session(r, fname))
	obj, ok := f.Get(toks[1]) // FIXME(sbinet): handle sub-dirs
	if !ok {
		return fmt.Errorf("could not find %q in file %q", toks[1], fname)
	}

	robj, ok := obj.(yodacnv.Marshaler)
	if !ok {
		return fmt.Errorf("object %q could not be converted to hbook.H1D", toks[1])
	}
	h1d, err := rootcnv.H1D(robj)
	if err != nil {
		return err
	}

	plot, err := hplot.New()
	if err != nil {
		return err
	}
	plot.Title.Text = obj.(rootio.Named).Title()

	h, err := hplot.NewH1D(h1d)
	if err != nil {
		return err
	}
	h.Infos.Style = hplot.HInfoSummary
	h.Color = color.RGBA{255, 0, 0, 255}

	plot.Add(h, hplot.NewGrid())

	svg, err := renderSVG(plot)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(string(svg))
}

const page = `<html>
<head>
    <title>go-hep/rootio file inspector</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css" />
	<link rel="stylesheet" href="https://www.w3schools.com/w3css/3/w3.css">
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
	<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/jstree/3.3.3/themes/default/style.min.css" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jstree/3.3.3/jstree.min.js"></script>
	<style>
	input[type=file] {
		display: none;
	}
	input[type=submit] {
		background-color: #F44336;
		padding:5px 15px;
		border:0 none;
		cursor:pointer;
		-webkit-border-radius: 5px;
		border-radius: 5px;
	}
	.flex-container {
		display: -webkit-flex;
		display: flex;
	}
	.flex-item {
		margin: 5px;
	}
	.rootio-file-upload {
		color: white;
		background-color: #0091EA;
		padding:5px 15px;
		border:0 none;
		cursor:pointer;
		-webkit-border-radius: 5px;
	}
	</style>
</head>
<body>

<!-- Sidebar -->
<div class="w3-sidebar w3-bar-block w3-card-4 w3-light-grey" style="width:25%">
	<div class="w3-bar-item w3-card-2 w3-black">
		<h2>go-hep/rootio ROOT file inspector</h2>
	</div>
	<div class="w3-bar-item">
	<form id="rootio-form" enctype="multipart/form-data" action="/root-file-upload" method="post">
		<label for="rootio-file" class="rootio-file-upload" style="font-size:16px">
		<i class="fa fa-cloud-upload" aria-hidden="true" style="font-size:16px"></i> Upload
		</label>
		<input id="rootio-file" type="file" name="upload-file"/>
		<input type="hidden" name="token" value="{{.Token}}"/>
		<input type="hidden" value="upload" />
	</form>
	</div>
	<div id="rootio-file-tree" class="w3-bar-item">
	</div>
</div>

<!-- Page Content -->
<div style="margin-left:25%; height:100%" class="w3-grey">
	<div class="w3-container w3-content w3-cell w3-cell-middle w3-cell-row w3-center w3-justify" id="rootio-display">
	</div>
</div>

<script type="text/javascript">
	document.getElementById("rootio-file").onchange = function() {
		var data = new FormData($("#rootio-form")[0]);
		$.ajax({
			url: "/root-file-upload",
			method: "POST",
			data: data,
			processData: false,
			contentType: false,
			success: function(result){
				// console.log("json-result: "+result);
				$('#rootio-file-tree').jstree(true).settings.core.data = JSON.parse(result);
				$("#rootio-file-tree").jstree(true).refresh();
			},
			error: function(er){
				alert("upload failed: "+er);
			}
		});
	}
	$(function () {
		$('#rootio-file-tree').jstree();
		$("#rootio-file-tree").on("select_node.jstree",
			function(evt, data){
				data.instance.toggle_node(data.node);
				if (data.node.a_attr.plot) {
					$.get(data.node.a_attr.href, plotCallback);
				}
			}
		);
	});

	function plotCallback(data, status) {
		var node = $("<div></div>");
		node.addClass("w3-panel w3-white w3-card-2 w3-display-container w3-content w3-center");
		node.css("width","100%");
		node.html(
			""
			+JSON.parse(data)
			+"<span onclick=\"this.parentElement.style.display='none'\" class=\"w3-button w3-display-topright w3-hover-red w3-tiny\">X</span>"
		);
		$("#rootio-display").append(node);
	};
</script>
</body>
</html>
`
