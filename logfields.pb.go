// Code generated by protoc-gen-gogo.
// source: logfields.proto
// DO NOT EDIT!

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
const _ = proto.GoGoProtoPackageIsVersion1

type LogField struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *LogField) Reset()                    { *m = LogField{} }
func (m *LogField) String() string            { return proto.CompactTextString(m) }
func (*LogField) ProtoMessage()               {}
func (*LogField) Descriptor() ([]byte, []int) { return fileDescriptorLogfields, []int{0} }

var E_Logfield = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FieldOptions)(nil),
	ExtensionType: (*LogField)(nil),
	Field:         12345,
	Name:          "improbable.logfield",
	Tag:           "bytes,12345,opt,name=logfield",
}

func init() {
	proto.RegisterType((*LogField)(nil), "improbable.LogField")
	proto.RegisterExtension(E_Logfield)
}

var fileDescriptorLogfields = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0xc9, 0x4f, 0x4f,
	0xcb, 0x4c, 0xcd, 0x49, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xca, 0xcc, 0x05,
	0x32, 0x92, 0x12, 0x93, 0x72, 0x52, 0xa5, 0x14, 0xd2, 0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1,
	0x32, 0x49, 0xa5, 0x69, 0xfa, 0x29, 0xa9, 0xc5, 0xc9, 0x45, 0x99, 0x05, 0x25, 0xf9, 0x45, 0x10,
	0xd5, 0x4a, 0x72, 0x5c, 0x1c, 0x3e, 0xf9, 0xe9, 0x6e, 0x20, 0x03, 0x84, 0x84, 0xb8, 0x58, 0xf2,
	0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0xab, 0x00, 0x2e, 0x0e,
	0x98, 0x05, 0x42, 0xb2, 0x7a, 0x10, 0xe3, 0xf4, 0x60, 0xc6, 0xe9, 0x81, 0xf5, 0xf9, 0x17, 0x94,
	0x64, 0xe6, 0xe7, 0x15, 0x4b, 0xec, 0x4c, 0x00, 0xea, 0xe3, 0x36, 0x12, 0xd1, 0x43, 0x38, 0x40,
	0x0f, 0x66, 0x76, 0x10, 0xdc, 0x14, 0x27, 0xee, 0x28, 0x4e, 0xb8, 0x93, 0x93, 0xd8, 0xc0, 0x46,
	0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x9e, 0xff, 0x5c, 0x3b, 0xc6, 0x00, 0x00, 0x00,
}