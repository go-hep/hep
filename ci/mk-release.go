// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/sync/errgroup"
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

	buildCmds(module, filepath.Join(tmp, "go.mod"))
}

func buildCmds(modname, fname string) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("could not read file %q: %+v", fname, err)
	}

	f, err := modfile.Parse(fname, data, nil)
	if err != nil {
		log.Fatalf("could not parse modfile: %+v", err)
	}
	found := false
	log.Printf("require:")
	var mod module.Version
loop:
	for _, req := range f.Require {
		log.Printf(" - %v: %v", req.Mod.Path, req.Mod.Version)
		if req.Mod.Path == modname {
			found = true
			mod = req.Mod
			break loop
		}
	}
	if !found {
		log.Fatalf("could not find module %q in modpub", modname)
	}

	tmp, err := ioutil.TempDir("", "go-hep-release-cmds-")
	if err != nil {
		log.Fatalf("could not create tmp dir for build cmds: %+v", err)
	}
	defer os.RemoveAll(tmp)

	allpkgs, err := pkgList(filepath.Dir(fname), modname, OSArch{"linux", "amd64"})
	if err != nil {
		log.Fatalf("could not build package list of module %q: %+v", modname, err)
	}

	for _, osarch := range []struct {
		os, arch string
	}{
		{"linux", "amd64"},
		{"linux", "386"},
		{"windows", "amd64"},
		{"windows", "386"},
		{"darwin", "amd64"},
		{"linux", "arm64"},
		{"freebsd", "amd64"},
	} {
		var (
			grp errgroup.Group
			ctx = osarch
		)
		log.Printf("--> GOOS=%s, GOARCH=%s", ctx.os, ctx.arch)
		cmds := make([]string, 0, len(allpkgs))
		for _, pkg := range allpkgs {
			if !strings.Contains(pkg, "/cmd") {
				continue
			}
			if strings.Contains(pkg, "/internal") {
				continue
			}
			if _, ok := excludeList[ctx][pkg]; ok {
				continue
			}
			cmds = append(cmds, pkg)
		}

		log.Printf("--> found %d commands", len(cmds))
		for i := range cmds {
			cmd := cmds[i]
			grp.Go(func() error {
				name := fmt.Sprintf("%s-%s_%s.exe", filepath.Base(cmd), ctx.os, ctx.arch)
				exe := filepath.Join(tmp, name)
				bld := exec.Command("go", "build", "-o", exe, cmd)
				bld.Env = append([]string{}, os.Environ()...)
				bld.Env = append(bld.Env, fmt.Sprintf("GOOS=%s", ctx.os), fmt.Sprintf("GOARCH=%s", ctx.arch))
				out, err := bld.CombinedOutput()
				if err != nil {
					log.Printf("could not compile %s: %+v\noutput:\n%s", name, err, out)
					return err
				}
				return nil
			})
		}

		err = grp.Wait()
		if err != nil {
			log.Fatalf("could not build commands for %s/%s: %+v", ctx.os, ctx.arch, err)
		}
	}

	upload(tmp, mod)
}

type OSArch struct {
	os, arch string
}

func pkgList(dir, module string, ctx OSArch) ([]string, error) {
	out := new(bytes.Buffer)
	cmd := exec.Command("go", "list", module+"/...")
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append([]string{}, os.Environ()...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", ctx.os), fmt.Sprintf("GOARCH=%s", ctx.arch))

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

func upload(dir string, mod module.Version) {
	cmd := exec.Command("scp", "-r", dir, "root@clrwebgohep.in2p3.fr:/srv/go-hep.org/dist/"+mod.Version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not upload binaries to server: %+v", err)
	}
}

var excludeList = map[OSArch]map[string]struct{}{
	OSArch{"darwin", "amd64"}: map[string]struct{}{
		"go-hep.org/x/hep/hplot/cmd/iplot": struct{}{},
	},
	OSArch{"freebsd", "amd64"}: map[string]struct{}{
		"go-hep.org/x/hep/groot/cmd/root-fuse": struct{}{},
		"go-hep.org/x/hep/xrootd/cmd/xrd-fuse": struct{}{},
	},
	OSArch{"windows", "amd64"}: map[string]struct{}{
		"go-hep.org/x/hep/groot/cmd/root-fuse": struct{}{},
		"go-hep.org/x/hep/xrootd/cmd/xrd-fuse": struct{}{},
	},
	OSArch{"windows", "386"}: map[string]struct{}{
		"go-hep.org/x/hep/groot/cmd/root-fuse": struct{}{},
		"go-hep.org/x/hep/xrootd/cmd/xrd-fuse": struct{}{},
	},
}
