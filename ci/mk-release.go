// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	var (
		module  = flag.String("module", "go-hep.org/x/hep", "module name to publish")
		version = flag.String("version", "latest", "module version to publish")
	)

	flag.Parse()

	publish(*module, *version)
}

func publish(module, version string) {
	log.Printf("publishing module=%q, version=%q", module, version)
	modver := module + "@" + version

	tmp, err := ioutil.TempDir("", "go-hep-release-")
	if err != nil {
		log.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	os.Setenv("GO111MODULE", "on")

	log.Printf("## creating modpub module...")
	cmd := exec.Command("go", "mod", "init", "modpub")
	cmd.Dir = tmp
	cmd.Stderr = log.Writer()
	cmd.Stdout = log.Writer()
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not initialize modpub module: %+v", err)
	}

	log.Printf("## get %q...", modver)
	cmd = exec.Command("go", "get", "-u", "-v", modver)
	cmd.Dir = tmp
	cmd.Stderr = log.Writer()
	cmd.Stdout = log.Writer()
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not get %q module: %+v", modver, err)
	}

	log.Printf("## generating main package...")
	const tmpl = `package main

import (
	_ "%s"
)

func main() {}
`

	err = ioutil.WriteFile(filepath.Join(tmp, "main.go"), []byte(fmt.Sprintf(tmpl, module)), 0644)
	if err != nil {
		log.Fatalf("could not generate main: %+v", err)
	}

	log.Printf("## go build...")
	cmd = exec.Command("go", "build", "-v")
	cmd.Dir = tmp
	cmd.Stderr = log.Writer()
	cmd.Stdout = log.Writer()
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not get %q module: %+v", modver, err)
	}
}
