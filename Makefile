all: setup
.PHONY: all

setup:
	go mod download
.PHONY: setup

test:
	go test -v ./...
.PHONY: test
