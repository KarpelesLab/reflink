#!/bin/make
GOPATH:=$(shell go env GOPATH)

.PHONY: test deps

all:
	$(GOPATH)/bin/goimports -w -l .
	go build -v

deps:
	go get -v -t .

test:
	go test -v
