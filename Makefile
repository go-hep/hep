## simple makefile to log workflow
.PHONY: all test clean build gen

all: build test
	@echo "#############################"
	root-ls -sinfos ./testdata/small.root
	@echo "## bye."

build: gen
	@go get -v ./...

gen:
	@go generate

test: build
	@go test -v

## EOF
