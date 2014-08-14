## simple makefile to log workflow
.PHONY: all test clean build

#GOFLAGS := $(GOFLAGS:-race -v)
GOFLAGS := $(GOFLAGS:-v)

all: clean build test
	@echo "## bye."

build:
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) ./...
	test-fads-app -l INFO -evtmax=-1 -cpu-prof

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
