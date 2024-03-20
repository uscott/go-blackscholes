.DEFAULT_GOAL := all

all: setup check test
.PHONY: all

check:
	golangci-lint run --fix
	test -z "$(shell golines -l .)"
.PHONY: check

format:
	gofmt -l -w -s .
	golines -w .
.PHONY: format

setup:
	go mod download
	pre-commit install
.PHONY: setup

test:
	go test -v ./...
.PHONY: test
