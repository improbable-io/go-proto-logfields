// Code generated by protoc-gen-gogo.
// source: service.proto
// DO NOT EDIT!

/*
Package example is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	Note
	Request
	Response
*/
package example

import github_com_improbable_io_go_proto_logfields "github.com/improbable-io/go-proto-logfields"
import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/improbable-io/go-proto-logfields"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Note) LogFields() map[string]string {
	// Handle being called on nil message.
	if this == nil {
		return map[string]string{}
	}
	// Generate fields for this message.
	return map[string]string{
		"author": this.Author,
	}
}

func (this *Note) ExtractRequestFields(dst map[string]interface{}) {
	// Handle being called on nil message.
	if this == nil {
		return
	}

	dst["author"] = this.Author
}

func (this *Request) LogFields() map[string]string {
	// Handle being called on nil message.
	if this == nil {
		return map[string]string{}
	}
	// Gather fields from oneofs and child messages.
	var hasInner bool
	noteFields := github_com_improbable_io_go_proto_logfields.ExtractLogFieldsFromMessage(this.Note)
	hasInner = hasInner || len(noteFields) > 0
	if !hasInner {
		// If no inner messages added any fields, avoid merging maps.
		return map[string]string{
			"path": this.Path,
		}
	}
	// Merge all the field maps.
	res := map[string]string{}
	res["path"] = this.Path
	for k, v := range noteFields {
		res[k] = v
	}
	return res
}

func (this *Request) ExtractRequestFields(dst map[string]interface{}) {
	// Handle being called on nil message.
	if this == nil {
		return
	}

	dst["path"] = this.Path
	github_com_improbable_io_go_proto_logfields.ExtractRequestFieldsFromMessage(this.Note, dst)
}

func (this *Response) LogFields() map[string]string {
	// Handle being called on nil message.
	if this == nil {
		return map[string]string{}
	}
	// Gather fields from oneofs and child messages.
	var hasInner bool
	changedNoteFields := github_com_improbable_io_go_proto_logfields.ExtractLogFieldsFromMessage(this.ChangedNote)
	hasInner = hasInner || len(changedNoteFields) > 0
	if !hasInner {
		// If no inner messages added any fields, avoid merging maps.
		return map[string]string{
			"did_it": fmt.Sprintf("%v", this.DidStuff),
		}
	}
	// Merge all the field maps.
	res := map[string]string{}
	res["did_it"] = fmt.Sprintf("%v", this.DidStuff)
	for k, v := range changedNoteFields {
		res[k] = v
	}
	return res
}

func (this *Response) ExtractRequestFields(dst map[string]interface{}) {
	// Handle being called on nil message.
	if this == nil {
		return
	}

	dst["did_it"] = this.DidStuff
	github_com_improbable_io_go_proto_logfields.ExtractRequestFieldsFromMessage(this.ChangedNote, dst)
}
