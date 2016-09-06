// Code generated by protoc-gen-go.
// source: date.proto
// DO NOT EDIT!

/*
Package ptypes is a generated protocol buffer package.

It is generated from these files:
	date.proto

It has these top-level messages:
	Date
*/
package ptypes

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Date struct {
	Day   int32 `protobuf:"varint,1,opt,name=day" json:"day,omitempty"`
	Month int32 `protobuf:"varint,2,opt,name=month" json:"month,omitempty"`
	Year  int32 `protobuf:"varint,3,opt,name=year" json:"year,omitempty"`
}

func (m *Date) Reset()                    { *m = Date{} }
func (m *Date) String() string            { return proto.CompactTextString(m) }
func (*Date) ProtoMessage()               {}
func (*Date) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Date)(nil), "altipla.protobuf.Date")
}

func init() { proto.RegisterFile("date.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 148 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x49, 0x2c, 0x49,
	0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x48, 0xcc, 0x29, 0xc9, 0x2c, 0xc8, 0x49, 0x84,
	0x70, 0x93, 0x4a, 0xd3, 0x94, 0x9c, 0xb8, 0x58, 0x5c, 0x80, 0xf2, 0x42, 0x02, 0x5c, 0xcc, 0x29,
	0x89, 0x95, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xac, 0x41, 0x20, 0xa6, 0x90, 0x08, 0x17, 0x6b, 0x6e,
	0x7e, 0x5e, 0x49, 0x86, 0x04, 0x13, 0x58, 0x0c, 0xc2, 0x11, 0x12, 0xe2, 0x62, 0xa9, 0x4c, 0x4d,
	0x2c, 0x92, 0x60, 0x06, 0x0b, 0x82, 0xd9, 0x4e, 0x26, 0x51, 0x46, 0xe9, 0x99, 0x25, 0x19, 0xa5,
	0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x50, 0x2b, 0x74, 0x93, 0xf3, 0xf3, 0x8a, 0x4b, 0x81, 0xec,
	0xbc, 0x74, 0x7d, 0x98, 0x6d, 0xba, 0x29, 0xa9, 0x69, 0xc5, 0xfa, 0x05, 0x25, 0x95, 0x05, 0xa9,
	0xc5, 0x49, 0x6c, 0x60, 0x51, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5e, 0x49, 0x93, 0x98,
	0xa0, 0x00, 0x00, 0x00,
}