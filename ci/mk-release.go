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
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	var (
		module  = flag.String("module", "go-hep.org/x/hep", "module name to publish")
		version = flag.String("version", "latest", "module version to publish")
		repo    = flag.String("repo", "git@github.com:go-hep/hep", "VCS URL of repository")
	)

	flag.Parse()

	publish(*module, *version, *repo)
}

func publish(module, version, repo string) {
	doLatest := version == "latest"
	log.Printf("publishing module=%q, version=%q", module, version)
	modver := module + "@" + version

	tmp, err := os.MkdirTemp("", "go-hep-release-")
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
	cmd = exec.Command("go", "get", "-v", modver)
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

	err = os.WriteFile(filepath.Join(tmp, "main.go"), []byte(fmt.Sprintf(tmpl, module)), 0644)
	if err != nil {
		log.Fatalf("could not generate main: %+v", err)
	}

	log.Printf("## mod download %q...", modver)
	cmd = exec.Command("go", "mod", "download")
	cmd.Dir = tmp
	cmd.Stderr = log.Writer()
	cmd.Stdout = log.Writer()
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not mod-tidy %q module: %+v", modver, err)
	}

	log.Printf("## mod tidy %q...", modver)
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tmp
	cmd.Stderr = log.Writer()
	cmd.Stdout = log.Writer()
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not mod-tidy %q module: %+v", modver, err)
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

	version, err = extractVersion(filepath.Join(tmp, "go.mod"), module)
	if err != nil {
		log.Fatalf("could not extract version from modpub module file: %+v", err)
	}

	buildCmds(module, version, repo)
	if doLatest {
		setLatest(version)
	}
}

func extractVersion(fname, modname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", fmt.Errorf("could not open module file %q: %w", fname, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if !strings.Contains(line, modname+" ") {
			continue
		}

		_, after, ok := strings.Cut(line, modname)
		if ok {
			return strings.TrimSpace(after), nil
		}
	}

	return "", fmt.Errorf("could not find module %q in modpub %q", modname, fname)
}

func buildCmds(modname, version, repo string) {
	top, err := os.MkdirTemp("", "go-hep-release-")
	if err != nil {
		log.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(top)

	src := filepath.Join(top, "hep")
	cmd := exec.Command(
		"git", "clone",
		"-b", version, "--depth", "1",
		repo,
		src,
	)
	cmd.Dir = top
	cmd.Stderr = log.Writer()
	cmd.Stdout = log.Writer()
	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not clone %q: %+v", repo, err)
	}

	tmp, err := os.MkdirTemp("", "go-hep-release-cmds-")
	if err != nil {
		log.Fatalf("could not create tmp dir for build cmds: %+v", err)
	}
	defer os.RemoveAll(tmp)

	allpkgs, err := pkgList(src, modname, OSArch{"linux", "amd64"})
	if err != nil {
		log.Fatalf("could not build package list of module %q: %+v", modname, err)
	}

	for _, osarch := range []struct {
		os, arch string
	}{
		{"linux", "amd64"},
		{"linux", "386"},
		{"linux", "arm64"},
		{"windows", "amd64"},
		{"windows", "386"},
		{"darwin", "amd64"},
		{"freebsd", "amd64"},
	} {
		var (
			grp errgroup.Group
			ctx = osarch
		)
		grp.SetLimit(4)
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

		tags := "-tags=netgo"
		log.Printf("--> found %d commands", len(cmds))
		for i := range cmds {
			cmd := cmds[i]
			grp.Go(func() error {
				name := fmt.Sprintf("%s-%s_%s.exe", filepath.Base(cmd), ctx.os, ctx.arch)
				exe := filepath.Join(tmp, name)
				bld := exec.Command(
					"go", "build",
					"-trimpath",
					"-buildvcs=true",
					"-o", exe,
					tags,
					strings.Replace(cmd, modname, ".", 1),
				)
				bld.Dir = src
				bld.Env = append([]string{}, os.Environ()...)
				bld.Env = append(bld.Env, fmt.Sprintf("GOOS=%s", ctx.os), fmt.Sprintf("GOARCH=%s", ctx.arch))
				if _, ok := needCgo[filepath.Base(cmd)]; !ok {
					bld.Env = append(bld.Env, "CGO_ENABLED=0")
				}
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

	upload(tmp, version)
}

type OSArch struct {
	os, arch string
}

func pkgList(dir, module string, ctx OSArch) ([]string, error) {
	env := append([]string{}, os.Environ()...)
	env = append(env, fmt.Sprintf("GOOS=%s", ctx.os), fmt.Sprintf("GOARCH=%s", ctx.arch))

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("could not initialize list: %w", err)
	}

	out := new(bytes.Buffer)
	cmd = exec.Command("go", "list", "./...")
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = env

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("could not get package list (%s-%s): %w", ctx.os, ctx.arch, err)
	}

	var pkgs []string
	scan := bufio.NewScanner(out)
	for scan.Scan() {
		pkg := scan.Text()
		if strings.Contains(pkg, "vendor") {
			continue
		}
		if !strings.HasPrefix(pkg, module) {
			continue
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func upload(dir string, version string) {
	cmd := exec.Command("scp", "-r", dir, "root@clrwebgohep.in2p3.fr:/srv/go-hep.org/dist/"+version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not upload binaries to server: %+v", err)
	}
}

func setLatest(version string) {
	cmd := exec.Command("ssh", "root@clrwebgohep.in2p3.fr",
		"--",
		fmt.Sprintf("/bin/sh -c 'cd /srv/go-hep.org/dist && /bin/rm ./latest && ln -s %s latest'", version),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not set latest to %q: %+v", version, err)
	}
}

var excludeList = map[OSArch]map[string]struct{}{
	{"linux", "386"}: {
		"go-hep.org/x/hep/hplot/cmd/iplot": struct{}{},
	},
	{"linux", "arm64"}: {
		"go-hep.org/x/hep/hplot/cmd/iplot": struct{}{},
	},
	{"darwin", "amd64"}: {
		"go-hep.org/x/hep/hplot/cmd/iplot": {},
	},
	{"freebsd", "amd64"}: {
		"go-hep.org/x/hep/hplot/cmd/iplot":     struct{}{},
		"go-hep.org/x/hep/groot/cmd/root-fuse": struct{}{},
		"go-hep.org/x/hep/xrootd/cmd/xrd-fuse": struct{}{},
	},
	{"windows", "amd64"}: {
		"go-hep.org/x/hep/groot/cmd/root-fuse": struct{}{},
		"go-hep.org/x/hep/xrootd/cmd/xrd-fuse": struct{}{},
	},
	{"windows", "386"}: {
		"go-hep.org/x/hep/groot/cmd/root-fuse": struct{}{},
		"go-hep.org/x/hep/xrootd/cmd/xrd-fuse": struct{}{},
	},
}

var needCgo = map[string]struct{}{
	"pawgo": {},
	"iplot": {},
	"hplot": {},
}
