all: setup test
.PHONY: all

setup:
	go mod download
.PHONY: setup

test:
	go test ./...
.PHONY: test
