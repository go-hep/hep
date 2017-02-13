## simple makefile to log workflow
.PHONY: all test clean build gen

all: build test
	@echo "## bye."

build:
	@go get -v ./...

gen:
	@go generate

test: build
	@go test -v

## EOF
