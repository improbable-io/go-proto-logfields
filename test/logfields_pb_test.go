// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package logfieldstest

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestNilProto3(t *testing.T) {
	fields := (*TestMessage3)(nil).LogFields()
	assert.Equal(t, map[string][]string{}, fields)
}

func TestEmptyProto3(t *testing.T) {
	fields := (&TestMessage3{}).LogFields()
	assert.Equal(t, map[string][]string{
		"an_int":       {"0"},
		"some_ints":    nil,
		"a_string":     {""},
		"some_strings": nil,
		"some_bytes":   {""},
		"many_bytes":   nil,
	}, fields)
}

func TestEmptyProto3WithChild(t *testing.T) {
	fields := (&TestMessage3{
		SingleMessage:    &TestMessage3_Inner{},
		RepeatedMessages: []*TestMessage3_Inner{&TestMessage3_Inner{}},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"an_int":       {"0"},
		"some_ints":    nil,
		"a_string":     {""},
		"some_strings": nil,
		"some_bytes":   {""},
		"many_bytes":   nil,
		"log_text":     {"", ""},
	}, fields)
}

func TestProto3Formatting(t *testing.T) {
	fields := (&TestMessage3{
		SingleInteger:   42,
		RepeatedInteger: []int32{43, 44},
		SingleString:    "23",
		RepeatedString:  []string{"24", "25"},
		SingleBytes:     []byte{1, 'a'},
		RepeatedBytes:   [][]byte{{2}, {'b'}},
		SingleMessage: &TestMessage3_Inner{
			Text: "a logged text",
		},
		RepeatedMessages: []*TestMessage3_Inner{
			&TestMessage3_Inner{Text: "text 1"},
			&TestMessage3_Inner{Text: "text 2"},
		},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"an_int":       {"42"},
		"some_ints":    {"43", "44"},
		"a_string":     {"23"},
		"some_strings": {"24", "25"},
		"some_bytes":   {"\x01a"},
		"many_bytes":   {"\x02", "b"},
		"log_text":     {"a logged text", "text 1", "text 2"},
	}, fields)
}

func TestNilProto2(t *testing.T) {
	fields := (*TestMessage2)(nil).LogFields()
	assert.Equal(t, map[string][]string{}, fields)
}

func TestEmptyProto2(t *testing.T) {
	fields := (&TestMessage2{}).LogFields()
	assert.Equal(t, map[string][]string{
		"opt_int":    {"0"},
		"req_int":    {"0"},
		"rep_int":    nil,
		"opt_string": {""},
		"req_string": {""},
		"rep_string": nil,
		"opt_bytes":  {""},
		"req_bytes":  {""},
		"rep_bytes":  nil,
	}, fields)
}

func TestEmptyProto2WithChild(t *testing.T) {
	fields := (&TestMessage2{
		RequiredMessage: &TestMessage2_Inner{},
		OptionalMessage: &TestMessage2_Inner{},
		RepeatedMessages: []*TestMessage2_Inner{
			&TestMessage2_Inner{},
		},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"log_text":   {"", "", ""},
		"opt_int":    {"0"},
		"req_int":    {"0"},
		"rep_int":    nil,
		"opt_string": {""},
		"req_string": {""},
		"rep_string": nil,
		"opt_bytes":  {""},
		"req_bytes":  {""},
		"rep_bytes":  nil,
	}, fields)
}

func TestProto2Formatting(t *testing.T) {
	fields := (&TestMessage2{
		RequiredInteger: proto.Int(42),
		OptionalInteger: proto.Int(42),
		RepeatedInteger: []int32{42},
		RequiredString:  proto.String("23"),
		OptionalString:  proto.String("23"),
		RepeatedString:  []string{"23"},
		RequiredBytes:   []byte{1, 'a'},
		OptionalBytes:   []byte{1, 'a'},
		RepeatedBytes:   [][]byte{{1, 'a'}},
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
	assert.Equal(t, map[string][]string{
		"log_text":   {"required message text", "optional message text", "repeated message text"},
		"opt_int":    {"42"},
		"req_int":    {"42"},
		"rep_int":    {"42"},
		"opt_string": {"23"},
		"req_string": {"23"},
		"rep_string": {"23"},
		"opt_bytes":  {"\x01a"},
		"req_bytes":  {"\x01a"},
		"rep_bytes":  {"\x01a"},
	}, fields)
}
