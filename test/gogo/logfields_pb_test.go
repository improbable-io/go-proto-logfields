// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package logfieldstest

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

type msg interface {
	LogFields() map[string]string
	ExtractRequestFields(map[string]interface{})
}

func assertExtracts(t *testing.T, expected map[string]interface{}, msg msg) {
	m := map[string]interface{}{}
	msg.ExtractRequestFields(m)
	assert.Equal(t, expected, m)
}

func TestNilProto3(t *testing.T) {
	assert.Equal(t, map[string]string{}, (*UnloggedTest3)(nil).LogFields())
	assert.Equal(t, map[string]string{}, (*FieldsTest3)(nil).LogFields())
	assert.Equal(t, map[string]string{}, (*EmbeddedTest3)(nil).LogFields())
	assert.Equal(t, map[string]string{}, (*OneOfTest3)(nil).LogFields())
	assert.Equal(t, map[string]string{}, (*MapTest3)(nil).LogFields())

	assertExtracts(t, map[string]interface{}{}, (*UnloggedTest3)(nil))
	assertExtracts(t, map[string]interface{}{}, (*FieldsTest3)(nil))
	assertExtracts(t, map[string]interface{}{}, (*EmbeddedTest3)(nil))
	assertExtracts(t, map[string]interface{}{}, (*OneOfTest3)(nil))
	assertExtracts(t, map[string]interface{}{}, (*MapTest3)(nil))
}

func TestEmptyUnloggedProto3(t *testing.T) {
	assert.Equal(t, map[string]string{}, (&UnloggedTest3{}).LogFields())
	assertExtracts(t, map[string]interface{}{}, (&UnloggedTest3{}))
}

func TestEmptyWithLoggedFields(t *testing.T) {
	assert.Equal(t, map[string]string{
		"an_int":     "0",
		"a_string":   "",
		"some_bytes": "",
	}, (&FieldsTest3{}).LogFields())
	assertExtracts(t, map[string]interface{}{
		"an_int":     int32(0),
		"a_string":   "",
		"some_bytes": []uint8(nil),
	}, &FieldsTest3{})
}

func TestEmptyWithEmbeddedMessages(t *testing.T) {
	assert.Equal(t, map[string]string{
		"a_string": "",
		"log_text": "",
	}, (&EmbeddedTest3{}).LogFields())
	assertExtracts(t, map[string]interface{}{
		"a_string": "",
		"log_text": "",
	}, &EmbeddedTest3{})
}

func TestEmptyWithOneOf(t *testing.T) {
	assert.Equal(t, map[string]string{
		"a_string": "",
	}, (&OneOfTest3{}).LogFields())
	assertExtracts(t, map[string]interface{}{
		"a_string": "",
	}, &OneOfTest3{})
}

func TestEmptyWithMap(t *testing.T) {
	assert.Equal(t, map[string]string{}, (&MapTest3{}).LogFields())
	assertExtracts(t, map[string]interface{}{}, &MapTest3{})
}

func TestEmbeddedMessagesEmpty(t *testing.T) {
	msg := &EmbeddedTest3{
		SingleMessage: &EmbeddedTest3_Inner{},
		RepeatedMessages: []*EmbeddedTest3_Inner{
			&EmbeddedTest3_Inner{},
		},
	}
	assert.Equal(t, map[string]string{
		"a_string": "",
		"log_text": "",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"a_string": "",
		"log_text": "",
	}, msg)
}

func TestOneOfEmbeddedMessageEmpty(t *testing.T) {
	msg := &OneOfTest3{
		AOneof: &OneOfTest3_SingleOneofMessage{
			&OneOfTest3_Inner{},
		},
	}
	assert.Equal(t, map[string]string{
		"a_string": "",
		"log_text": "",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"a_string": "",
		"log_text": "",
	}, msg)
}

func TestOneOfEmbeddedMessage(t *testing.T) {
	msg := &OneOfTest3{
		AOneof: &OneOfTest3_SingleOneofMessage{
			&OneOfTest3_Inner{
				SingleInnerString: "a_text",
			},
		},
	}
	assert.Equal(t, map[string]string{
		"a_string": "",
		"log_text": "a_text",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"a_string": "",
		"log_text": "a_text",
	}, msg)
}

func TestOneOfUnloggedField(t *testing.T) {
	msg := &OneOfTest3{
		AOneof: &OneOfTest3_UnloggedOneofString{},
	}
	assert.Equal(t, map[string]string{
		"a_string": "",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"a_string": "",
	}, msg)
}

func TestOneOfStringField(t *testing.T) {
	msg := &OneOfTest3{
		AOneof: &OneOfTest3_SingleOneofString{},
	}
	assert.Equal(t, map[string]string{
		"a_string": "",
		"log_text": "",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"a_string": "",
		"log_text": "",
	}, msg)
}

func TestMapEntry(t *testing.T) {
	msg := &MapTest3{
		AStringMap: map[string]string{
			"a_string_key": "a_string_value",
		},
		AStringToInnerMap: map[string]*MapTest3_Inner{
			"a_inner_key": &MapTest3_Inner{
				SingleInnerString: "a_inner_string_value",
			},
		},
	}
	assert.Equal(t, map[string]string{}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{}, msg)
}

func TestProto3Formatting(t *testing.T) {
	msg := &FieldsTest3{
		SingleInteger: 42,
		SingleString:  "23",
		SingleBytes:   []byte{1, 'a'},
	}
	assert.Equal(t, map[string]string{
		"an_int":     "42",
		"a_string":   "23",
		"some_bytes": "\x01a",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"an_int":     int32(42),
		"a_string":   "23",
		"some_bytes": []byte{1, 'a'},
	}, msg)
}

func TestNilProto2(t *testing.T) {
	assert.Equal(t, map[string]string{}, (*TestMessage2)(nil).LogFields())
	assertExtracts(t, map[string]interface{}{}, (*TestMessage2)(nil))
}

func TestEmptyProto2(t *testing.T) {
	msg := &TestMessage2{}
	assert.Equal(t, map[string]string{
		"opt_int":    "0",
		"req_int":    "0",
		"opt_string": "",
		"req_string": "",
		"opt_bytes":  "",
		"req_bytes":  "",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"opt_int":    int32(0),
		"req_int":    int32(0),
		"opt_string": "",
		"req_string": "",
		"opt_bytes":  []byte(nil),
		"req_bytes":  []byte(nil),
	}, msg)
}

func TestEmptyProto2WithChild(t *testing.T) {
	msg := &TestMessage2{
		RequiredMessage: &TestMessage2_Inner{},
		OptionalMessage: &TestMessage2_Inner{},
	}
	assert.Equal(t, map[string]string{
		"log_text":   "",
		"opt_int":    "0",
		"req_int":    "0",
		"opt_string": "",
		"req_string": "",
		"opt_bytes":  "",
		"req_bytes":  "",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"log_text":   "",
		"opt_int":    int32(0),
		"req_int":    int32(0),
		"opt_string": "",
		"req_string": "",
		"opt_bytes":  []byte(nil),
		"req_bytes":  []byte(nil),
	}, msg)
}

func TestProto2Formatting(t *testing.T) {
	msg := &TestMessage2{
		RequiredInteger: proto.Int(42),
		OptionalInteger: proto.Int(42),
		RequiredString:  proto.String("23"),
		OptionalString:  proto.String("23"),
		RequiredBytes:   []byte{1, 'a'},
		OptionalBytes:   []byte{1, 'a'},
		RequiredMessage: &TestMessage2_Inner{
			Text: proto.String("required message text"),
		},
		OptionalMessage: &TestMessage2_Inner{
			Text: proto.String("optional message text"),
		},
	}
	assert.Equal(t, map[string]string{
		"log_text":   "optional message text",
		"opt_int":    "42",
		"req_int":    "42",
		"opt_string": "23",
		"req_string": "23",
		"opt_bytes":  "\x01a",
		"req_bytes":  "\x01a",
	}, msg.LogFields())
	assertExtracts(t, map[string]interface{}{
		"log_text":   "optional message text",
		"opt_int":    int32(42),
		"req_int":    int32(42),
		"opt_string": "23",
		"req_string": "23",
		"opt_bytes":  []byte{1, 'a'},
		"req_bytes":  []byte{1, 'a'},
	}, msg)
}
