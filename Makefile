## simple makefile to log workflow
.PHONY: all test clean build

#GOFLAGS := $(GOFLAGS:-race -v)

all: build test
	@# done

build: clean tabledata.go
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) -v ./...

clean:
	@go clean $(GOFLAGS) -i ./...

tabledata.go: tabledata.tbl
	@cat tabledata.header tabledata.tbl tabledata.footer >| tabledata.go

## EOF
