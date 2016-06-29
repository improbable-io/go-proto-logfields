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

With the generated code, a set of logging fields can be generated as follows:
```
fields := (&SomethingRequest{What: "something"}).LogFields()
fmt.Println(fields)
// prints: map[what_was_requested: something]
```
