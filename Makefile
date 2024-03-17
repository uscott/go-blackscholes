.DEFAULT_GOAL := all

all: setup test
.PHONY: all

check:
	golangci-lint run --fix
.PHONY: check

format:
	gofmt -l -w -s .
	golines -w .
.PHONY: format

setup:
	go mod download
.PHONY: setup

test:
	go test -v ./...
.PHONY: test
