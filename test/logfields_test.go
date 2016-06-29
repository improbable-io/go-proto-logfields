// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package logfieldstest

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestNilProto3(t *testing.T) {
	fields := (*TestMessage3)(nil).LogFields()
	assert.Equal(t, map[string]string{}, fields)
}

func TestEmptyProto3(t *testing.T) {
	fields := (&TestMessage3{}).LogFields()
	assert.Equal(t, map[string]string{
		"an_int":     "0",
		"a_string":   "",
		"some_bytes": "",
	}, fields)
}

func TestEmptyProto3WithChild(t *testing.T) {
	fields := (&TestMessage3{
		SingleMessage: &TestMessage3_Inner{},
	}).LogFields()
	assert.Equal(t, map[string]string{
		"a_string":   "",
		"an_int":     "0",
		"some_bytes": "",
		"log_text":   "",
	}, fields)
}

func TestProto3Formatting(t *testing.T) {
	fields := (&TestMessage3{
		SingleInteger: 42,
		SingleString:  "23",
		SingleBytes:   []byte{1, 'a'},
		SingleMessage: &TestMessage3_Inner{
			Text: "a logged text",
		},
	}).LogFields()
	assert.Equal(t, map[string]string{
		"a_string":   "23",
		"an_int":     "42",
		"some_bytes": "\x01a",
		"log_text":   "a logged text",
	}, fields)
}

func TestNilProto2(t *testing.T) {
	fields := (*TestMessage2)(nil).LogFields()
	assert.Equal(t, map[string]string{}, fields)
}

func TestEmptyProto2(t *testing.T) {
	fields := (&TestMessage2{}).LogFields()
	assert.Equal(t, map[string]string{
		"opt_int":    "0",
		"req_int":    "0",
		"opt_string": "",
		"req_string": "",
		"opt_bytes":  "",
		"req_bytes":  "",
	}, fields)
}

func TestEmptyProto2WithChild(t *testing.T) {
	fields := (&TestMessage2{
		RequiredMessage: &TestMessage2_Inner{},
		OptionalMessage: &TestMessage2_Inner{},
	}).LogFields()
	assert.Equal(t, map[string]string{
		"log_text":   "",
		"opt_int":    "0",
		"req_int":    "0",
		"opt_string": "",
		"req_string": "",
		"opt_bytes":  "",
		"req_bytes":  "",
	}, fields)
}

func TestProto2Formatting(t *testing.T) {
	fields := (&TestMessage2{
		RequiredInteger: proto.Int(42),
		RequiredString:  proto.String("23"),
		RequiredBytes:   []byte{1, 'a'},
		OptionalInteger: proto.Int(42),
		OptionalString:  proto.String("23"),
		OptionalBytes:   []byte{1, 'a'},
		RequiredMessage: &TestMessage2_Inner{
			Text: proto.String("required message text"),
		},
		OptionalMessage: &TestMessage2_Inner{
			Text: proto.String("optional message text"),
		},
		RepeatedMessages: []*TestMessage2_Inner{{
			Text: proto.String("repeated message text"),
		}},
	}).LogFields()
	assert.Equal(t, map[string]string{
		"log_text":   "optional message text",
		"opt_int":    "42",
		"req_int":    "42",
		"opt_string": "23",
		"req_string": "23",
		"opt_bytes":  "\x01a",
		"req_bytes":  "\x01a",
	}, fields)
}
