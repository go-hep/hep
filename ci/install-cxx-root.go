// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// Command to install a given (binary) ROOT version.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	rvers := flag.String("root-version", "6.18.04", "ROOT version to install")
	odir := flag.String("o", "", "install directory for ROOT")

	flag.Parse()

	if *odir == "" {
		*odir = fmt.Sprintf("root-%s", *rvers)
	}

	dst, err := filepath.Abs(*odir)
	if err != nil {
		log.Fatalf("could not get absolute path to %q: %+v", *odir, err)
	}
	switch _, err := os.Stat(filepath.Join(dst, "root-"+*rvers, "bin", "root.exe")); err {
	case nil:
		log.Printf("ROOT version %s already installed", *rvers)
		return
	default:
		err = os.MkdirAll(dst, 0755)
		if err != nil {
			log.Fatalf("could not create directory %q: %+v", dst, err)
		}
	}

	install(*rvers, dst)
}

func install(rvers, odir string) {
	log.Printf("installing ROOT %s to: %q...", rvers, odir)
	tmp, err := ioutil.TempDir("", "go-hep-build-")
	if err != nil {
		log.Fatalf("could not create top-level tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	oname := filepath.Join(tmp, "root.tar.gz")
	targz, err := os.Create(oname)
	if err != nil {
		log.Fatalf("could not create ROOT archive destination file %q: %+v", oname, err)
	}
	defer targz.Close()

	binURL := fmt.Sprintf("https://cern.ch/binet/go-hep/cern-root/root-%s-linux_amd64.tar.gz", rvers)
	log.Printf("downloading %q...", binURL)
	resp, err := http.Get(binURL)
	if err != nil {
		log.Fatalf("could not get ROOT %s from %v: %+v", rvers, binURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("could not get ROOT %s from %v: %v", rvers, binURL, resp.Status)
	}

	_, err = io.Copy(targz, resp.Body)
	if err != nil {
		log.Fatalf("could not download ROOT archive from %q: %+v", binURL, err)
	}

	err = targz.Close()
	if err != nil {
		log.Fatalf("could not close ROOT binary archive %q: %+v", oname, err)
	}

	log.Printf("decompressing...")
	cmd := exec.Command("tar", "xf", oname)
	cmd.Dir = odir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not untar ROOT binary archive: %+v", err)
	}
}
