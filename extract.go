package logfields

import "github.com/gogo/protobuf/proto"

type protoWithLogFields interface {
	proto.Message
	LogFields() map[string]string
	ExtractRequestFields(dst map[string]interface{})
}

func ExtractLogFieldsFromMessage(message proto.Message) map[string]string {
	if m, ok := message.(protoWithLogFields); ok {
		return m.LogFields()
	}

	return map[string]string{}
}

func ExtractRequestFieldsFromMessage(message proto.Message, dst map[string]interface{}) {
	if m, ok := message.(protoWithLogFields); ok {
		m.ExtractRequestFields(dst)
	}
}
