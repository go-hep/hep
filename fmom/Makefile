## simple makefile to log workflow
.PHONY: all test clean build

GOFLAGS := $(GOFLAGS:)

all: clean build test

build:
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) -v ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
