// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package logfieldstest

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestNilProto3(t *testing.T) {
	assert.Equal(t, map[string][]string{}, (*UnloggedTest3)(nil).LogFields())
	assert.Equal(t, map[string][]string{}, (*FieldsTest3)(nil).LogFields())
	assert.Equal(t, map[string][]string{}, (*EmbeddedTest3)(nil).LogFields())
	assert.Equal(t, map[string][]string{}, (*OneOfTest3)(nil).LogFields())
	assert.Equal(t, map[string][]string{}, (*MapTest3)(nil).LogFields())
}

func TestEmptyUnloggedProto3(t *testing.T) {
	assert.Equal(t, map[string][]string{}, (&UnloggedTest3{}).LogFields())
}

func TestEmptyWithLoggedFields(t *testing.T) {
	assert.Equal(t, map[string][]string{
		"an_int":       {"0"},
		"some_ints":    nil,
		"a_string":     {""},
		"some_strings": nil,
		"some_bytes":   {""},
		"many_bytes":   nil,
	}, (&FieldsTest3{}).LogFields())
}

func TestEmptyWithEmbeddedMessages(t *testing.T) {
	assert.Equal(t, map[string][]string{
		"a_string": {""},
		"log_text": {""},
	}, (&EmbeddedTest3{}).LogFields())
}

func TestEmptyWithOneOf(t *testing.T) {
	assert.Equal(t, map[string][]string{
		"a_string": {""},
		"log_text": {""},
	}, (&OneOfTest3{}).LogFields())
}

func TestEmptyWithMap(t *testing.T) {
	assert.Equal(t, map[string][]string{}, (&MapTest3{}).LogFields())
}

func TestEmbeddedMessagesEmpty(t *testing.T) {
	fields := (&EmbeddedTest3{
		SingleMessage: &EmbeddedTest3_Inner{},
		RepeatedMessages: []*EmbeddedTest3_Inner{
			&EmbeddedTest3_Inner{},
		},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"a_string": {""},
		"log_text": {"", "", ""},
	}, fields)
}

func TestOneOfEmbeddedMessageEmpty(t *testing.T) {
	fields := (&OneOfTest3{
		AOneof: &OneOfTest3_SingleOneofMessage{
			&OneOfTest3_Inner{},
		},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"a_string": {""},
		"log_text": {"", ""},
	}, fields)
}

func TestOneOfUnloggedField(t *testing.T) {
	fields := (&OneOfTest3{
		AOneof: &OneOfTest3_UnloggedOneofString{},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"a_string": {""},
		"log_text": {""},
	}, fields)
}

func TestOneOfStringField(t *testing.T) {
	fields := (&OneOfTest3{
		AOneof: &OneOfTest3_SingleOneofString{},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"a_string": {""},
		"log_text": {"", ""},
	}, fields)
}

func TestMapEntry(t *testing.T) {
	fields := (&MapTest3{
		AStringMap: map[string]string{
			"a_string_key": "a_string_value",
		},
		AStringToInnerMap: map[string]*MapTest3_Inner{
			"a_inner_key": &MapTest3_Inner{
				SingleInnerString: "a_inner_string_value",
			},
		},
	})
	assert.Equal(t, map[string][]string{
		"log_text": {"a_inner_string_value"},
	}, fields.LogFields())
}

func TestProto3Formatting(t *testing.T) {
	fields := (&FieldsTest3{
		SingleInteger:   42,
		RepeatedInteger: []int32{43, 44},
		SingleString:    "23",
		RepeatedString:  []string{"24", "25"},
		SingleBytes:     []byte{1, 'a'},
		RepeatedBytes:   [][]byte{{2}, {'b'}},
	}).LogFields()
	assert.Equal(t, map[string][]string{
		"an_int":       {"42"},
		"some_ints":    {"43", "44"},
		"a_string":     {"23"},
		"some_strings": {"24", "25"},
		"some_bytes":   {"\x01a"},
		"many_bytes":   {"\x02", "b"},
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
