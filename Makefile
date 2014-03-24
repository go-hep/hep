## simple makefile to log workflow
.PHONY: all test clean build

all: build test
	@echo "## bye."

build:
	@go get -v ./...

test: build
	@go test -v

## EOF
