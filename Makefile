#!/usr/bin/make -f

SHELL := /bin/bash

export PATH := $(shell echo $${GOPATH//:/\/bin:}/bin:$${PATH})

.PHONY: install regenerate regenerate_examples test

install: regenerate
	go install github.com/improbable-io/go-proto-logfields/protoc-gen-gologfields

regenerate:
	protoc \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src \
	  --proto_path=./vendor/github.com/gogo/protobuf/protobuf \
	  --proto_path=. \
	  --gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. \
	  src/github.com/improbable-io/go-proto-logfields/logfields.proto

regenerate_examples: install
	protoc \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src \
	  --proto_path=./vendor/github.com/gogo/protobuf/protobuf \
	  --proto_path=. \
	  --go_out=. \
	  --gologfields_out=gogoimport=false:. \
	  examples/*.proto

regenerate_test: install
	protoc \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src \
	  --proto_path=./vendor/github.com/gogo/protobuf/protobuf \
	  --proto_path=. \
	  --go_out=. \
	  --gologfields_out=gogoimport=false:. \
	  test/*.proto

regenerate_test_gogo: install
	protoc \
	  --proto_path=$${GOPATH//:/\/src --proto_path=}/src \
	  --proto_path=./vendor/github.com/gogo/protobuf/protobuf \
	  --proto_path=. \
	  --gogo_out=. \
	  --gologfields_out=gogoimport=true:. \
	  test/*.proto

test:
	go test -v .
	$(MAKE) regenerate_test
	go test -v ./test
	$(MAKE) regenerate_test_gogo
	go test -v ./test
