// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: logfields.proto

/*
Package logfields is a generated protocol buffer package.

It is generated from these files:
	logfields.proto

It has these top-level messages:
	LogField
*/
package logfields

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type LogField struct {
	// name of the log context field where the value of this field should be
	// recorded. Fields with empty names are ignored.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *LogField) Reset()                    { *m = LogField{} }
func (m *LogField) String() string            { return proto.CompactTextString(m) }
func (*LogField) ProtoMessage()               {}
func (*LogField) Descriptor() ([]byte, []int) { return fileDescriptorLogfields, []int{0} }

func (m *LogField) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

var E_Logfield = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FieldOptions)(nil),
	ExtensionType: (*LogField)(nil),
	Field:         62132,
	Name:          "improbable.logfield",
	Tag:           "bytes,62132,opt,name=logfield",
	Filename:      "logfields.proto",
}

func init() {
	proto.RegisterType((*LogField)(nil), "improbable.LogField")
	proto.RegisterExtension(E_Logfield)
}

func init() { proto.RegisterFile("logfields.proto", fileDescriptorLogfields) }

var fileDescriptorLogfields = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0xc9, 0x4f, 0x4f,
	0xcb, 0x4c, 0xcd, 0x49, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xca, 0xcc, 0x2d,
	0x28, 0xca, 0x4f, 0x4a, 0x4c, 0xca, 0x49, 0x95, 0x52, 0x48, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5,
	0x07, 0xcb, 0x24, 0x95, 0xa6, 0xe9, 0xa7, 0xa4, 0x16, 0x27, 0x17, 0x65, 0x16, 0x94, 0xe4, 0x17,
	0x41, 0x54, 0x2b, 0xc9, 0x71, 0x71, 0xf8, 0xe4, 0xa7, 0xbb, 0x81, 0x0c, 0x10, 0x12, 0xe2, 0x62,
	0xc9, 0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0xad, 0x02, 0xb9,
	0x38, 0x60, 0x16, 0x08, 0xc9, 0xea, 0x41, 0x8c, 0xd3, 0x83, 0x19, 0xa7, 0x07, 0xd6, 0xe7, 0x5f,
	0x50, 0x92, 0x99, 0x9f, 0x57, 0x2c, 0xb1, 0xe5, 0x29, 0xb3, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x88,
	0x1e, 0xc2, 0x05, 0x7a, 0x30, 0xc3, 0x83, 0xe0, 0xc6, 0x38, 0x71, 0x47, 0x71, 0xc2, 0xdd, 0x9c,
	0xc4, 0x06, 0x36, 0xcb, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe5, 0x19, 0xb4, 0x18, 0xc7, 0x00,
	0x00, 0x00,
}
