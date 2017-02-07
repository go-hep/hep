## simple makefile to log workflow
.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install test


build:
	@go build $(GOFLAGS) ./...

install:
	@go get $(GOFLAGS) ./...

test: install
	@go test $(GOFLAGS) ./...
	fads-app -l INFO -evtmax=-1 -cpu-prof

bench: install
	@go test -bench=. $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
