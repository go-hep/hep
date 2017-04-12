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
//	"crypto/tls"
	"os"

	"golang.org/x/crypto/acme/autocert"
	
	"go-hep.org/x/hep/rootio/cmd/root-srv/server"
	
)

var (
	addrFlag = flag.String("addr", ":8080", "server address:port")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-srv [options]

ex:

 $> root-srv -addr=:8080
 2017/04/06 15:13:59 server listening on :8080

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()
	server.Init()

//	certManager := autocert.Manager{
//		Prompt:     autocert.AcceptTOS,
//		HostPolicy: autocert.HostWhitelist("example.com"), //your domain here
//		Cache:      autocert.DirCache("certs"), //folder for storing certificates
//	}
//	
//
//	// piqué sur sx
//	server := &http.Server{
//		Addr: ":443",
//		TLSConfig: &tls.Config{
//			GetCertificate: certManager.GetCertificate,
//		},
//	}
//
//	server.ListenAndServeTLS("", "")
	
	log.Printf("Http server listening on %s", *addrFlag)
//	log.Fatal(http.ListenAndServe(*addrFlag, nil))
	//	log.Fatal(http.Serve(autocert.NewListener("localhost"),nil ))

	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("127.0.0.1"),
		Cache:      autocert.DirCache("certs"), //folder for storing certificates
	}

	log.Fatal(http.Serve(m.Listener(),nil ))
	
}
