// Code generated by protoc-gen-go.
// source: protos/fn.proto
// DO NOT EDIT!

/*
Package fn is a generated protocol buffer package.

It is generated from these files:
	protos/fn.proto

It has these top-level messages:
	ListFunctionsReply
	Function
	DeployFunctionRequest
	DeleteFunctionRequest
*/
package fn

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ListFunctionsReply struct {
	Functions []*Function `protobuf:"bytes,1,rep,name=functions" json:"functions,omitempty"`
}

func (m *ListFunctionsReply) Reset()                    { *m = ListFunctionsReply{} }
func (m *ListFunctionsReply) String() string            { return proto.CompactTextString(m) }
func (*ListFunctionsReply) ProtoMessage()               {}
func (*ListFunctionsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ListFunctionsReply) GetFunctions() []*Function {
	if m != nil {
		return m.Functions
	}
	return nil
}

type Function struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Call    string `protobuf:"bytes,2,opt,name=call" json:"call,omitempty"`
	Trigger string `protobuf:"bytes,3,opt,name=trigger" json:"trigger,omitempty"`
	Method  string `protobuf:"bytes,4,opt,name=method" json:"method,omitempty"`
	// Output only.
	CreatedAt *google_protobuf1.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
}

func (m *Function) Reset()                    { *m = Function{} }
func (m *Function) String() string            { return proto.CompactTextString(m) }
func (*Function) ProtoMessage()               {}
func (*Function) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Function) GetCreatedAt() *google_protobuf1.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

type DeployFunctionRequest struct {
	Function *Function `protobuf:"bytes,1,opt,name=function" json:"function,omitempty"`
}

func (m *DeployFunctionRequest) Reset()                    { *m = DeployFunctionRequest{} }
func (m *DeployFunctionRequest) String() string            { return proto.CompactTextString(m) }
func (*DeployFunctionRequest) ProtoMessage()               {}
func (*DeployFunctionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DeployFunctionRequest) GetFunction() *Function {
	if m != nil {
		return m.Function
	}
	return nil
}

type DeleteFunctionRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *DeleteFunctionRequest) Reset()                    { *m = DeleteFunctionRequest{} }
func (m *DeleteFunctionRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteFunctionRequest) ProtoMessage()               {}
func (*DeleteFunctionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*ListFunctionsReply)(nil), "fn.ListFunctionsReply")
	proto.RegisterType((*Function)(nil), "fn.Function")
	proto.RegisterType((*DeployFunctionRequest)(nil), "fn.DeployFunctionRequest")
	proto.RegisterType((*DeleteFunctionRequest)(nil), "fn.DeleteFunctionRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Fn service

type FnClient interface {
	ListFunctions(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*ListFunctionsReply, error)
	DeployFunction(ctx context.Context, in *DeployFunctionRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	DeleteFunction(ctx context.Context, in *DeleteFunctionRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type fnClient struct {
	cc *grpc.ClientConn
}

func NewFnClient(cc *grpc.ClientConn) FnClient {
	return &fnClient{cc}
}

func (c *fnClient) ListFunctions(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*ListFunctionsReply, error) {
	out := new(ListFunctionsReply)
	err := grpc.Invoke(ctx, "/fn.Fn/ListFunctions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fnClient) DeployFunction(ctx context.Context, in *DeployFunctionRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/fn.Fn/DeployFunction", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fnClient) DeleteFunction(ctx context.Context, in *DeleteFunctionRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/fn.Fn/DeleteFunction", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Fn service

type FnServer interface {
	ListFunctions(context.Context, *google_protobuf.Empty) (*ListFunctionsReply, error)
	DeployFunction(context.Context, *DeployFunctionRequest) (*google_protobuf.Empty, error)
	DeleteFunction(context.Context, *DeleteFunctionRequest) (*google_protobuf.Empty, error)
}

func RegisterFnServer(s *grpc.Server, srv FnServer) {
	s.RegisterService(&_Fn_serviceDesc, srv)
}

func _Fn_ListFunctions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FnServer).ListFunctions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fn.Fn/ListFunctions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FnServer).ListFunctions(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fn_DeployFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FnServer).DeployFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fn.Fn/DeployFunction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FnServer).DeployFunction(ctx, req.(*DeployFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fn_DeleteFunction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFunctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FnServer).DeleteFunction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fn.Fn/DeleteFunction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FnServer).DeleteFunction(ctx, req.(*DeleteFunctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Fn_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fn.Fn",
	HandlerType: (*FnServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListFunctions",
			Handler:    _Fn_ListFunctions_Handler,
		},
		{
			MethodName: "DeployFunction",
			Handler:    _Fn_DeployFunction_Handler,
		},
		{
			MethodName: "DeleteFunction",
			Handler:    _Fn_DeleteFunction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("protos/fn.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 320 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x91, 0x4f, 0x4f, 0xc2, 0x40,
	0x10, 0xc5, 0x2d, 0x20, 0xc2, 0xe0, 0x9f, 0x64, 0x12, 0xc9, 0x5a, 0x0f, 0x9a, 0x9e, 0x88, 0x26,
	0x25, 0xc1, 0x93, 0x37, 0x49, 0x84, 0x93, 0xa7, 0xc6, 0xbb, 0x29, 0x30, 0xad, 0x4d, 0xda, 0xdd,
	0xda, 0x0e, 0x07, 0x3e, 0x8d, 0x9f, 0xca, 0xef, 0x63, 0x77, 0xcb, 0xa2, 0xfc, 0xbb, 0xcd, 0xbe,
	0x79, 0x7d, 0xfb, 0x7e, 0x5d, 0xb8, 0xca, 0x0b, 0xc5, 0xaa, 0x1c, 0x46, 0xd2, 0x37, 0x13, 0x36,
	0x22, 0xe9, 0xde, 0xc6, 0x4a, 0xc5, 0x29, 0x0d, 0x8d, 0x32, 0x5b, 0x46, 0x43, 0xca, 0x72, 0x5e,
	0xd5, 0x06, 0xf7, 0x6e, 0x77, 0xc9, 0x49, 0x46, 0x25, 0x87, 0x59, 0x5e, 0x1b, 0xbc, 0x17, 0xc0,
	0xb7, 0xa4, 0xe4, 0xe9, 0x52, 0xce, 0x39, 0x51, 0xb2, 0x0c, 0x28, 0x4f, 0x57, 0xf8, 0x00, 0xdd,
	0xc8, 0x2a, 0xc2, 0xb9, 0x6f, 0x0e, 0x7a, 0xa3, 0x73, 0xbf, 0xba, 0xd5, 0xda, 0x82, 0xbf, 0xb5,
	0xf7, 0xed, 0x40, 0xc7, 0xea, 0x88, 0xd0, 0x92, 0x61, 0x46, 0xd5, 0x37, 0xce, 0xa0, 0x1b, 0x98,
	0x59, 0x6b, 0xf3, 0x30, 0x4d, 0x45, 0xa3, 0xd6, 0xf4, 0x8c, 0x02, 0xce, 0xb8, 0x48, 0xe2, 0x98,
	0x0a, 0xd1, 0x34, 0xb2, 0x3d, 0x62, 0x1f, 0xda, 0x19, 0xf1, 0xa7, 0x5a, 0x88, 0x96, 0x59, 0xac,
	0x4f, 0xf8, 0x0c, 0x30, 0x2f, 0x28, 0x64, 0x5a, 0x7c, 0x84, 0x2c, 0x4e, 0xab, 0x5d, 0x6f, 0xe4,
	0xfa, 0x35, 0x9e, 0x6f, 0xf1, 0xfc, 0x77, 0x8b, 0x17, 0x74, 0xd7, 0xee, 0x31, 0x7b, 0x63, 0xb8,
	0x7e, 0xad, 0xb0, 0xd4, 0x6a, 0x53, 0x9f, 0xbe, 0x96, 0x95, 0x0d, 0x07, 0xd0, 0xb1, 0x1c, 0xa6,
	0xf1, 0x2e, 0xe5, 0x66, 0xeb, 0x3d, 0xea, 0x88, 0x94, 0x98, 0x76, 0x23, 0x0e, 0x00, 0x8f, 0x7e,
	0x1c, 0x68, 0x4c, 0x25, 0x8e, 0xe1, 0x62, 0xeb, 0xd7, 0x62, 0x7f, 0xaf, 0xee, 0x44, 0x3f, 0x95,
	0xdb, 0xd7, 0x97, 0xee, 0xbf, 0x82, 0x77, 0x82, 0x13, 0xb8, 0xdc, 0x6e, 0x8e, 0x37, 0xda, 0x7b,
	0x90, 0xc6, 0x3d, 0x12, 0x6f, 0x63, 0xfe, 0xb7, 0xb7, 0x31, 0x07, 0x88, 0x8e, 0xc7, 0xcc, 0xda,
	0x46, 0x79, 0xfa, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x34, 0x02, 0x79, 0x85, 0x87, 0x02, 0x00, 0x00,
}
