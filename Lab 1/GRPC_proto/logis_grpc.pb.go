// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package GRPC_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// LogisServiceClient is the client API for LogisService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogisServiceClient interface {
	PedidoCliente(ctx context.Context, in *Pedido, opts ...grpc.CallOption) (*CodSeguimiento, error)
	SeguimientoCliente(ctx context.Context, in *CodSeguimiento, opts ...grpc.CallOption) (*EstadoPedido, error)
}

type logisServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogisServiceClient(cc grpc.ClientConnInterface) LogisServiceClient {
	return &logisServiceClient{cc}
}

func (c *logisServiceClient) PedidoCliente(ctx context.Context, in *Pedido, opts ...grpc.CallOption) (*CodSeguimiento, error) {
	out := new(CodSeguimiento)
	err := c.cc.Invoke(ctx, "/logis.LogisService/PedidoCliente", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisServiceClient) SeguimientoCliente(ctx context.Context, in *CodSeguimiento, opts ...grpc.CallOption) (*EstadoPedido, error) {
	out := new(EstadoPedido)
	err := c.cc.Invoke(ctx, "/logis.LogisService/SeguimientoCliente", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogisServiceServer is the server API for LogisService service.
// All implementations must embed UnimplementedLogisServiceServer
// for forward compatibility
type LogisServiceServer interface {
	PedidoCliente(context.Context, *Pedido) (*CodSeguimiento, error)
	SeguimientoCliente(context.Context, *CodSeguimiento) (*EstadoPedido, error)
	mustEmbedUnimplementedLogisServiceServer()
}

// UnimplementedLogisServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLogisServiceServer struct {
}

func (UnimplementedLogisServiceServer) PedidoCliente(context.Context, *Pedido) (*CodSeguimiento, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PedidoCliente not implemented")
}
func (UnimplementedLogisServiceServer) SeguimientoCliente(context.Context, *CodSeguimiento) (*EstadoPedido, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SeguimientoCliente not implemented")
}
func (UnimplementedLogisServiceServer) mustEmbedUnimplementedLogisServiceServer() {}

// UnsafeLogisServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogisServiceServer will
// result in compilation errors.
type UnsafeLogisServiceServer interface {
	mustEmbedUnimplementedLogisServiceServer()
}

func RegisterLogisServiceServer(s *grpc.Server, srv LogisServiceServer) {
	s.RegisterService(&_LogisService_serviceDesc, srv)
}

func _LogisService_PedidoCliente_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pedido)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisServiceServer).PedidoCliente(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logis.LogisService/PedidoCliente",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisServiceServer).PedidoCliente(ctx, req.(*Pedido))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisService_SeguimientoCliente_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CodSeguimiento)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisServiceServer).SeguimientoCliente(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logis.LogisService/SeguimientoCliente",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisServiceServer).SeguimientoCliente(ctx, req.(*CodSeguimiento))
	}
	return interceptor(ctx, in, info, handler)
}

var _LogisService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logis.LogisService",
	HandlerType: (*LogisServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PedidoCliente",
			Handler:    _LogisService_PedidoCliente_Handler,
		},
		{
			MethodName: "SeguimientoCliente",
			Handler:    _LogisService_SeguimientoCliente_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "GRPC_proto/logis.proto",
}
