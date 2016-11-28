## simple makefile to log workflow
.PHONY: all test clean build

all: build test
	@echo "#############################"
	root-ls ./testdata/small.root
	@echo "## bye."

build:
	@echo "build github.com/go-hep/rootio"
	@go get -v .
	@echo "build github.com/go-hep/rootio/cmd/root-ls"
	@go get -v ./cmd/root-ls

test: build
	@go test -v

## EOF
