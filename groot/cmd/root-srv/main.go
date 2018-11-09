// Copyright 2017 The go-hep Authors. All rights reserved.
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
//  $> root-srv -addr :8080 -serv https -host example.com
//  2017/04/06 15:13:59 https server listening on :8080 at example.com
package main // import "go-hep.org/x/hep/groot/cmd/root-srv"

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/crypto/acme/autocert"
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

	log.SetPrefix("root-srv: ")
	log.SetFlags(0)

	dir, err := ioutil.TempDir("", "groot-rsrv-")
	if err != nil {
		log.Panicf("could not create temporary directory: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	run(dir, c)
}

func run(dir string, c chan os.Signal) {
	defer func() {
		log.Printf("shutdown sequence...")
		os.RemoveAll(dir)
	}()

	log.Printf("%s server listening on %s", *servFlag, *addrFlag)

	srv := newServer(*hostFlag == "", dir, http.DefaultServeMux)
	defer srv.Shutdown()

	go func() {
		if *servFlag == "http" {
			log.Fatal(http.ListenAndServe(*addrFlag, nil))
		} else if *servFlag == "https" {
			m := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(*hostFlag),
				Cache:      autocert.DirCache("certs"), //folder for storing certificates
			}
			server := &http.Server{
				Addr: *addrFlag,
				TLSConfig: &tls.Config{
					GetCertificate: m.GetCertificate,
				},
			}
			log.Fatal(server.ListenAndServeTLS("", ""))
		}
	}()
	<-c
}
