## simple makefile to log workflow
.PHONY: all test clean build run

#GOFLAGS := $(GOFLAGS:-race -v)
GOFLAGS := $(GOFLAGS:-v)

all: build test

build: clean
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) -v ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
