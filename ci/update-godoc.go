// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {

	log.SetPrefix("")
	log.SetFlags(0)

	njobs := flag.Int("j", 5, "number of parallel godoc update jobs")

	flag.Parse()

	pkgs, err := pkgList()
	if err != nil {
		log.Fatal(err)
	}

	var (
		n   = len(pkgs) / *njobs
		grp errgroup.Group
	)
	log.Printf("pkgs:  %d", len(pkgs))
	log.Printf("njobs: %d", *njobs)

	for i := 0; i < len(pkgs); i += n {
		beg := i
		end := beg + n
		if end > len(pkgs) {
			end = len(pkgs)
		}
		grp.Go(func() error {
			for _, pkg := range pkgs[beg:end] {
				log.Printf("updating %q...", pkg)
				v := make(url.Values)
				v.Add("path", pkg)
				resp, err := http.PostForm("https://godoc.org/-/refresh", v)
				if err != nil {
					return fmt.Errorf("could not post %q: %w", pkg, err)
				}
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("invalid response status for %q: %v", pkg, resp.Status)
				}
			}
			return nil
		})
	}

	err = grp.Wait()
	if err != nil {
		log.Printf("error running group: %v", err)
	}
}

func pkgList() ([]string, error) {
	out := new(bytes.Buffer)
	cmd := exec.Command("go", "list", "./...")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("could not get package list: %w", err)
	}

	var pkgs []string
	scan := bufio.NewScanner(out)
	for scan.Scan() {
		pkg := scan.Text()
		if strings.Contains(pkg, "vendor") {
			continue
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

/*
for pkg in $(go list ./...); do curl -X POST -d "path=${pkg}" https://godoc.org/-/refresh; done
*/
