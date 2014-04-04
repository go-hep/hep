## simple makefile to log workflow
.PHONY: all test clean build run

#GOFLAGS := $(GOFLAGS:-race -v)
GOFLAGS := $(GOFLAGS:-v)

all: build test
	@echo "## bye."

build: clean
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
