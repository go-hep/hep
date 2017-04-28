// +build !go1.7
// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-srv runs a web server that can inspect and browse ROOT files.
// root-srv can also display ROOT objects (TH1x, TH2x, TGraphs, TGraphErrors,
// TGraphAsymmErrors, TDirectories, TTrees, ...).
//
// Usage: root-srv [options]
//
// ex:
//
//  $> root-srv -addr=:8080 &
//  2017/04/06 15:13:59 server listening on :8080
//

//  $> open localhost:8080
//
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-hep.org/x/hep/rootio/cmd/root-srv/server"
)

var (
	addrFlag = flag.String("addr", ":8080", "server address:port")
	servFlag = flag.String("serv", "http", "server protocol")
	hostFlag = flag.String("host", "", "server domain name for TLS ")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-srv [options]

ex:


 $> root-srv -addr :8080 -serv https -host example.com
 2017/04/06 15:13:59 https server listening on :8080 at example.com

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()
	server.Init()

	log.Printf("%s server listening on %s", *servFlag, *addrFlag)

	if *servFlag == "http" {
		log.Fatal(http.ListenAndServe(*addrFlag, nil))
	} else if *servFlag == "https" {
		//		m := autocert.Manager{
		//			Prompt:     autocert.AcceptTOS,
		//			HostPolicy: autocert.HostWhitelist(*hostFlag),
		//			Cache:      autocert.DirCache("certs"), //folder for storing certificates
		//		}
		//		server := &http.Server{
		//			Addr: *addrFlag,
		//			TLSConfig: &tls.Config{
		//				GetCertificate: m.GetCertificate,
		//			},
		//		}
		log.Fatal(server.ListenAndServeTLS("", ""))
	}
}
