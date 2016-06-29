#!/usr/bin/make -f

SHELL := /bin/bash

export PATH := $(shell echo $${GOPATH//:/\/bin:}/bin:$${PATH})

.PHONY: install regenerate regenerate_examples test

install: regenerate
	go install github.com/improbable-io/go-proto-logfields/protoc-gen-gologfields

regenerate:
	protoc \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src/github.com/gogo/protobuf/protobuf \
	  --proto_path=. \
	  --gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. \
	  logfields.proto

regenerate_examples: install
	protoc \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src/github.com/google/protobuf/src \
	  --proto_path=. \
	  --go_out=. \
	  --gologfields_out=gogoimport=false:. \
	  examples/*.proto

regenerate_test: install
	protoc \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src/github.com/google/protobuf/src \
	  --proto_path=. \
	  --go_out=. \
	  --gologfields_out=gogoimport=false:. \
	  test/*.proto

test: regenerate_test
	go test github.com/improbable-io/go-proto-logfields/test
