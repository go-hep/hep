## simple makefile to log workflow
.PHONY: all test clean build

all: build test
	@echo "## bye."

build:
	@echo "build github.com/go-hep/rootio"
	@go get -v .
	@echo "build github.com/go-hep/rootio/root-ls"
	@go get -v ./root-ls

test: build
	@go test -v

## EOF
