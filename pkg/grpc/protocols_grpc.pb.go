// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: protocols.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	OpenAIService_CreateChatCompletion_FullMethodName = "/openai.OpenAIService/CreateChatCompletion"
)

// OpenAIServiceClient is the client API for OpenAIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OpenAIServiceClient interface {
	CreateChatCompletion(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (*ChatCompletionResponse, error)
}

type openAIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOpenAIServiceClient(cc grpc.ClientConnInterface) OpenAIServiceClient {
	return &openAIServiceClient{cc}
}

func (c *openAIServiceClient) CreateChatCompletion(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (*ChatCompletionResponse, error) {
	out := new(ChatCompletionResponse)
	err := c.cc.Invoke(ctx, OpenAIService_CreateChatCompletion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OpenAIServiceServer is the server API for OpenAIService service.
// All implementations must embed UnimplementedOpenAIServiceServer
// for forward compatibility
type OpenAIServiceServer interface {
	CreateChatCompletion(context.Context, *ChatCompletionRequest) (*ChatCompletionResponse, error)
	mustEmbedUnimplementedOpenAIServiceServer()
}

// UnimplementedOpenAIServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOpenAIServiceServer struct {
}

func (UnimplementedOpenAIServiceServer) CreateChatCompletion(context.Context, *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChatCompletion not implemented")
}
func (UnimplementedOpenAIServiceServer) mustEmbedUnimplementedOpenAIServiceServer() {}

// UnsafeOpenAIServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OpenAIServiceServer will
// result in compilation errors.
type UnsafeOpenAIServiceServer interface {
	mustEmbedUnimplementedOpenAIServiceServer()
}

func RegisterOpenAIServiceServer(s grpc.ServiceRegistrar, srv OpenAIServiceServer) {
	s.RegisterService(&OpenAIService_ServiceDesc, srv)
}

func _OpenAIService_CreateChatCompletion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatCompletionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpenAIServiceServer).CreateChatCompletion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OpenAIService_CreateChatCompletion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpenAIServiceServer).CreateChatCompletion(ctx, req.(*ChatCompletionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OpenAIService_ServiceDesc is the grpc.ServiceDesc for OpenAIService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OpenAIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "openai.OpenAIService",
	HandlerType: (*OpenAIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChatCompletion",
			Handler:    _OpenAIService_CreateChatCompletion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protocols.proto",
}
