// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// Command to build a given ROOT version from sources.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	rvers := flag.String("root-version", "6.34.04", "ROOT version to build")
	nproc := flag.Int("j", runtime.NumCPU(), "number of parallel build processes")

	flag.Parse()

	build(*rvers, *nproc)
}

func build(rvers string, nproc int) {
	tmp, err := os.MkdirTemp("", "build-root-")
	if err != nil {
		log.Fatalf("could not create top-level tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	bdir := filepath.Join(tmp, "build")
	err = os.MkdirAll(bdir, 0755)
	if err != nil {
		log.Fatalf("could not create build dir %q: %+v", bdir, err)
	}

	dst := filepath.Join(tmp, "root-"+rvers)
	err = os.MkdirAll(dst, 0755)
	if err != nil {
		log.Fatalf("could not create dst dir %q: %+v", dst, err)
	}

	fname := filepath.Join(tmp, "build.sh")
	err = os.WriteFile(fname, []byte(fmt.Sprintf(docker, rvers, nproc)), 0644)
	if err != nil {
		log.Fatalf("could not create build-script: %+v", err)
	}

	cmd := exec.Command("docker", "run", "--rm",
		"--network=host",
		"-v", fname+":/build.sh",
		"-v", bdir+":/build/src",
		"-v", dst+":/build/install",
		"ubuntu:24.04", "/bin/sh", "/build.sh",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not build ROOT-%s: %+v", rvers, err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get pwd: %+v", err)
	}

	cmd = exec.Command("tar", "zcf", filepath.Join(wd, "root-"+rvers+"-linux_amd64.tar.gz"), "root-"+rvers)
	cmd.Dir = tmp
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("could not tar ROOT-%s: %+v", rvers, err)
	}
}

const docker = `#!/bin/sh

set -e
set -x

export DEBIAN_FRONTEND=noninteractive

apt update  -yq
apt install -yq    \
		cmake curl \
		g++ git    \
		python3    \
;

export ROOT_VERSION="%[1]s"

cd /tmp

curl -O -L https://root.cern.ch/download/root_v${ROOT_VERSION}.source.tar.gz
tar zxf ./root_v${ROOT_VERSION}.source.tar.gz

cd /build/src
cmake /tmp/root-${ROOT_VERSION} \
 -DCMAKE_INSTALL_PREFIX=/build/install -DCMAKE_BUILD_TYPE=Release \
 -Dall=OFF -Dfail-on-missing=OFF \
 -Dalien=OFF \
 -Dastiff=OFF \
 -Dbonjour=OFF \
 -Dbuiltin_afterimage=OFF \
 -Dbuiltin_ftgl=OFF \
 -Dbuiltin_glez=OFF \
 -Dcastor=OFF \
 -Dclad=OFF \
 -Dchirp=OFF \
 -Ddcache=OFF \
 -Dfftw3=OFF \
 -Dfitsio=OFF \
 -Dgenvector=OFF \
 -Dgfal=OFF \
 -Dglite=OFF \
 -Dgnuinstall=OFF \
 -Dgraphics=OFF \
 -Dgviz=OFF \
 -Dhdfs=OFF \
 -Dkrb5=OFF \
 -Dldap=OFF \
 -Dmathmore=OFF \
 -Dmonalisa=OFF \
 -Dmysql=OFF \
 -Dodbc=OFF \
 -Dpython=OFF \
 -Dshared=OFF \
 -Dtmva=OFF \
 -Dvdt=OFF \
 -Dx11=OFF \
 -Dxrootd=OFF \
;

make -j%[2]d
make -j%[2]d install
`
