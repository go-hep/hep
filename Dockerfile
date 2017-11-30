##
## a simple Dockerfile where Go, Go-HEP and Neugram are installed
##
from andrewosh/binder-base
maintainer Sebastien Binet <binet@cern.ch>

user root

env GOVERS 1.9.2

# install Go
run apt-get update -y && \
	apt-get install -y curl git pkg-config mercurial && \
	curl -O -L https://golang.org/dl/go${GOVERS}.linux-amd64.tar.gz && \
	tar -C /usr/local -zxf go${GOVERS}.linux-amd64.tar.gz && \
	/bin/rm go${GOVERS}.linux-amd64.tar.gz

# prepare for Go plugin compilation
run mkdir /usr/local/go/pkg/linux_amd64_dynlink && \
	chown -R main:main /usr/local/go

user main

env GOPATH $HOME/gopath
env PATH $GOPATH/bin:/usr/local/go/bin:$PATH

run go get golang.org/x/tools/cmd/goimports && \
	go get neugram.io/ng/... && \
	go get go-hep.org/x/hep/... && \
	go get gonum.org/v1/gonum/...

# install the Go kernel
run git clone https://github.com/neugram/binder $HOME/.local/share/jupyter/kernels/neugram

run mkdir -p $HOME/notebooks

user root
run chown -R main:main /home/main/notebooks
