# go-proto-logfields

A protoc plugin for generating log field extractors.

## Use-case

It is useful to augment log messages with additional context from the execution environment. For example, when handling an RPC, it may be useful to add some fields of the request message to any messages logged while handling the request.

## Install

* Install `protoc`, and ensure it is on your path.
* Install Go.
* `go install github.com/gogo/protobuf/protoc-gen-gogo`
* `go install github.com/improbable-io/go-proto-logfields/protoc-gen-gologfields`

## Usage

Given a an RPC:
```
syntax = "proto3";
package doer_of_something;
import "github.com/improbable-io/go-proto-logfields/logfields.proto";

Service Doer {
  rpc DoSomething (SomethingRequest) returns (SomethingResponse) {}
}
```

The request might be annotated as such:
```
message SomethingRequest {
    string what = 1 [(improbable.logfield) = {name: "what_was_requested"}];
}
```

The logfields extractor can be generated with:
```
protoc \
  --proto_path=${GOPATH//:/\/src --proto_path=}/src \
  --proto_path=${GOPATH//:/\/src --proto_path=}/src/github.com/gogo/protobuf/src \
  --proto_path=. \
  --gogo_out=. \
  --gologfields_out=. \
  *.proto
```

With the generated code, a set of logging fields can be generated as follows:
```
fields := (&SomethingRequest{What: "something"}).LogFields()
fmt.Println(fields)
// prints: map[what_was_requested: something]
```

## Using golang/protobuf instead of gogo/protobuf

By default, the output files use the `github.com/gogo/protobuf` implementation. To use the `github.com/golang/protobuf` implementation, the logfields generator must be passed a `gogoimport=false` flag as follows:
```
protoc \
  --proto_path=${GOPATH//:/\/src --proto_path=}/src \
  --proto_path=${GOPATH//:/\/src --proto_path=}/src/github.com/google/protobuf/src \
  --proto_path=. \
  --go_out=. \
  --gologfields_out=gogoimport=false:. \
  *.proto
```

The changes are:
* The `--proto_path` is adjusted from the gogo repository to the google repository.
* `--gogo_out` is replaced with `--go_out`.
* `--gologfields_out=.` is replaced with `--gologfields_out=gogoimport=false:.`.
