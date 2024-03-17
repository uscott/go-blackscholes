.DEFAULT_GOAL := all

all: setup test
.PHONY: all

setup:
	go mod download
.PHONY: setup

test:
	go test -v ./...
.PHONY: test
