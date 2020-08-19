package logfields

import "google.golang.org/protobuf/proto"

type protoWithLogFields interface {
	proto.Message
	LogFields() map[string]string
	ExtractRequestFields(prefixes []string, dst map[string]interface{})
}

func ExtractLogFieldsFromMessage(message proto.Message) map[string]string {
	if m, ok := message.(protoWithLogFields); ok {
		return m.LogFields()
	}

	return map[string]string{}
}

func ExtractRequestFieldsFromMessage(message proto.Message, prefixes []string, dst map[string]interface{}) {
	if m, ok := message.(protoWithLogFields); ok {
		m.ExtractRequestFields(prefixes, dst)
	}
}
