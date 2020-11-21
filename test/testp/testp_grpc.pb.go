// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package testp

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// TestpServiceClient is the client API for TestpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestpServiceClient interface {
	RecibirBytes(ctx context.Context, opts ...grpc.CallOption) (TestpService_RecibirBytesClient, error)
}

type testpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTestpServiceClient(cc grpc.ClientConnInterface) TestpServiceClient {
	return &testpServiceClient{cc}
}

func (c *testpServiceClient) RecibirBytes(ctx context.Context, opts ...grpc.CallOption) (TestpService_RecibirBytesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TestpService_serviceDesc.Streams[0], "/testp.TestpService/RecibirBytes", opts...)
	if err != nil {
		return nil, err
	}
	x := &testpServiceRecibirBytesClient{stream}
	return x, nil
}

type TestpService_RecibirBytesClient interface {
	Send(*MSG) error
	CloseAndRecv() (*ACK, error)
	grpc.ClientStream
}

type testpServiceRecibirBytesClient struct {
	grpc.ClientStream
}

func (x *testpServiceRecibirBytesClient) Send(m *MSG) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testpServiceRecibirBytesClient) CloseAndRecv() (*ACK, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ACK)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestpServiceServer is the server API for TestpService service.
// All implementations must embed UnimplementedTestpServiceServer
// for forward compatibility
type TestpServiceServer interface {
	RecibirBytes(TestpService_RecibirBytesServer) error
	mustEmbedUnimplementedTestpServiceServer()
}

// UnimplementedTestpServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTestpServiceServer struct {
}

func (UnimplementedTestpServiceServer) RecibirBytes(TestpService_RecibirBytesServer) error {
	return status.Errorf(codes.Unimplemented, "method RecibirBytes not implemented")
}
func (UnimplementedTestpServiceServer) mustEmbedUnimplementedTestpServiceServer() {}

// UnsafeTestpServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestpServiceServer will
// result in compilation errors.
type UnsafeTestpServiceServer interface {
	mustEmbedUnimplementedTestpServiceServer()
}

func RegisterTestpServiceServer(s *grpc.Server, srv TestpServiceServer) {
	s.RegisterService(&_TestpService_serviceDesc, srv)
}

func _TestpService_RecibirBytes_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestpServiceServer).RecibirBytes(&testpServiceRecibirBytesServer{stream})
}

type TestpService_RecibirBytesServer interface {
	SendAndClose(*ACK) error
	Recv() (*MSG, error)
	grpc.ServerStream
}

type testpServiceRecibirBytesServer struct {
	grpc.ServerStream
}

func (x *testpServiceRecibirBytesServer) SendAndClose(m *ACK) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testpServiceRecibirBytesServer) Recv() (*MSG, error) {
	m := new(MSG)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _TestpService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "testp.TestpService",
	HandlerType: (*TestpServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RecibirBytes",
			Handler:       _TestpService_RecibirBytes_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "testp.proto",
}