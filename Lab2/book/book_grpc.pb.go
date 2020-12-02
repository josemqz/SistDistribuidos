// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package book

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// BookServiceClient is the client API for BookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookServiceClient interface {
	RequestRA(ctx context.Context, in *ExMutua, opts ...grpc.CallOption) (*ACK, error)
	EnviarChunkDN(ctx context.Context, in *Chunk, opts ...grpc.CallOption) (*Chunk, error)
	RecibirChunksDN(ctx context.Context, in *Chunk, opts ...grpc.CallOption) (*ACK, error)
	RecibirPropuesta(ctx context.Context, in *PropuestaLibro, opts ...grpc.CallOption) (*RespuestaP, error)
	RecibirChunks(ctx context.Context, opts ...grpc.CallOption) (BookService_RecibirChunksClient, error)
	EscribirLogDes(ctx context.Context, in *PropuestaLibro, opts ...grpc.CallOption) (*ACK, error)
	RecibirPropDatanode(ctx context.Context, in *PropuestaLibro, opts ...grpc.CallOption) (*PropuestaLibro, error)
	ChunkInfoLog(ctx context.Context, in *ChunksInfo, opts ...grpc.CallOption) (*ChunksInfo, error)
	EnviarListaLibros(ctx context.Context, in *ACK, opts ...grpc.CallOption) (*ListaLibros, error)
}

type bookServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBookServiceClient(cc grpc.ClientConnInterface) BookServiceClient {
	return &bookServiceClient{cc}
}

func (c *bookServiceClient) RequestRA(ctx context.Context, in *ExMutua, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/book.BookService/RequestRA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) EnviarChunkDN(ctx context.Context, in *Chunk, opts ...grpc.CallOption) (*Chunk, error) {
	out := new(Chunk)
	err := c.cc.Invoke(ctx, "/book.BookService/EnviarChunkDN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) RecibirChunksDN(ctx context.Context, in *Chunk, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/book.BookService/RecibirChunksDN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) RecibirPropuesta(ctx context.Context, in *PropuestaLibro, opts ...grpc.CallOption) (*RespuestaP, error) {
	out := new(RespuestaP)
	err := c.cc.Invoke(ctx, "/book.BookService/RecibirPropuesta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) RecibirChunks(ctx context.Context, opts ...grpc.CallOption) (BookService_RecibirChunksClient, error) {
	stream, err := c.cc.NewStream(ctx, &_BookService_serviceDesc.Streams[0], "/book.BookService/RecibirChunks", opts...)
	if err != nil {
		return nil, err
	}
	x := &bookServiceRecibirChunksClient{stream}
	return x, nil
}

type BookService_RecibirChunksClient interface {
	Send(*Chunk) error
	CloseAndRecv() (*ACK, error)
	grpc.ClientStream
}

type bookServiceRecibirChunksClient struct {
	grpc.ClientStream
}

func (x *bookServiceRecibirChunksClient) Send(m *Chunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *bookServiceRecibirChunksClient) CloseAndRecv() (*ACK, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ACK)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *bookServiceClient) EscribirLogDes(ctx context.Context, in *PropuestaLibro, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/book.BookService/EscribirLogDes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) RecibirPropDatanode(ctx context.Context, in *PropuestaLibro, opts ...grpc.CallOption) (*PropuestaLibro, error) {
	out := new(PropuestaLibro)
	err := c.cc.Invoke(ctx, "/book.BookService/RecibirPropDatanode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) ChunkInfoLog(ctx context.Context, in *ChunksInfo, opts ...grpc.CallOption) (*ChunksInfo, error) {
	out := new(ChunksInfo)
	err := c.cc.Invoke(ctx, "/book.BookService/ChunkInfoLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) EnviarListaLibros(ctx context.Context, in *ACK, opts ...grpc.CallOption) (*ListaLibros, error) {
	out := new(ListaLibros)
	err := c.cc.Invoke(ctx, "/book.BookService/EnviarListaLibros", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookServiceServer is the server API for BookService service.
// All implementations must embed UnimplementedBookServiceServer
// for forward compatibility
type BookServiceServer interface {
	RequestRA(context.Context, *ExMutua) (*ACK, error)
	EnviarChunkDN(context.Context, *Chunk) (*Chunk, error)
	RecibirChunksDN(context.Context, *Chunk) (*ACK, error)
	RecibirPropuesta(context.Context, *PropuestaLibro) (*RespuestaP, error)
	RecibirChunks(BookService_RecibirChunksServer) error
	EscribirLogDes(context.Context, *PropuestaLibro) (*ACK, error)
	RecibirPropDatanode(context.Context, *PropuestaLibro) (*PropuestaLibro, error)
	ChunkInfoLog(context.Context, *ChunksInfo) (*ChunksInfo, error)
	EnviarListaLibros(context.Context, *ACK) (*ListaLibros, error)
	mustEmbedUnimplementedBookServiceServer()
}

// UnimplementedBookServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBookServiceServer struct {
}

func (UnimplementedBookServiceServer) RequestRA(context.Context, *ExMutua) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestRA not implemented")
}
func (UnimplementedBookServiceServer) EnviarChunkDN(context.Context, *Chunk) (*Chunk, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnviarChunkDN not implemented")
}
func (UnimplementedBookServiceServer) RecibirChunksDN(context.Context, *Chunk) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecibirChunksDN not implemented")
}
func (UnimplementedBookServiceServer) RecibirPropuesta(context.Context, *PropuestaLibro) (*RespuestaP, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecibirPropuesta not implemented")
}
func (UnimplementedBookServiceServer) RecibirChunks(BookService_RecibirChunksServer) error {
	return status.Errorf(codes.Unimplemented, "method RecibirChunks not implemented")
}
func (UnimplementedBookServiceServer) EscribirLogDes(context.Context, *PropuestaLibro) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EscribirLogDes not implemented")
}
func (UnimplementedBookServiceServer) RecibirPropDatanode(context.Context, *PropuestaLibro) (*PropuestaLibro, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecibirPropDatanode not implemented")
}
func (UnimplementedBookServiceServer) ChunkInfoLog(context.Context, *ChunksInfo) (*ChunksInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChunkInfoLog not implemented")
}
func (UnimplementedBookServiceServer) EnviarListaLibros(context.Context, *ACK) (*ListaLibros, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnviarListaLibros not implemented")
}
func (UnimplementedBookServiceServer) mustEmbedUnimplementedBookServiceServer() {}

// UnsafeBookServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookServiceServer will
// result in compilation errors.
type UnsafeBookServiceServer interface {
	mustEmbedUnimplementedBookServiceServer()
}

func RegisterBookServiceServer(s *grpc.Server, srv BookServiceServer) {
	s.RegisterService(&_BookService_serviceDesc, srv)
}

func _BookService_RequestRA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExMutua)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).RequestRA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/RequestRA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).RequestRA(ctx, req.(*ExMutua))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_EnviarChunkDN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Chunk)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).EnviarChunkDN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/EnviarChunkDN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).EnviarChunkDN(ctx, req.(*Chunk))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_RecibirChunksDN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Chunk)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).RecibirChunksDN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/RecibirChunksDN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).RecibirChunksDN(ctx, req.(*Chunk))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_RecibirPropuesta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PropuestaLibro)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).RecibirPropuesta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/RecibirPropuesta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).RecibirPropuesta(ctx, req.(*PropuestaLibro))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_RecibirChunks_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BookServiceServer).RecibirChunks(&bookServiceRecibirChunksServer{stream})
}

type BookService_RecibirChunksServer interface {
	SendAndClose(*ACK) error
	Recv() (*Chunk, error)
	grpc.ServerStream
}

type bookServiceRecibirChunksServer struct {
	grpc.ServerStream
}

func (x *bookServiceRecibirChunksServer) SendAndClose(m *ACK) error {
	return x.ServerStream.SendMsg(m)
}

func (x *bookServiceRecibirChunksServer) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _BookService_EscribirLogDes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PropuestaLibro)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).EscribirLogDes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/EscribirLogDes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).EscribirLogDes(ctx, req.(*PropuestaLibro))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_RecibirPropDatanode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PropuestaLibro)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).RecibirPropDatanode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/RecibirPropDatanode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).RecibirPropDatanode(ctx, req.(*PropuestaLibro))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_ChunkInfoLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChunksInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).ChunkInfoLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/ChunkInfoLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).ChunkInfoLog(ctx, req.(*ChunksInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_EnviarListaLibros_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ACK)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).EnviarListaLibros(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.BookService/EnviarListaLibros",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).EnviarListaLibros(ctx, req.(*ACK))
	}
	return interceptor(ctx, in, info, handler)
}

var _BookService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "book.BookService",
	HandlerType: (*BookServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestRA",
			Handler:    _BookService_RequestRA_Handler,
		},
		{
			MethodName: "EnviarChunkDN",
			Handler:    _BookService_EnviarChunkDN_Handler,
		},
		{
			MethodName: "RecibirChunksDN",
			Handler:    _BookService_RecibirChunksDN_Handler,
		},
		{
			MethodName: "RecibirPropuesta",
			Handler:    _BookService_RecibirPropuesta_Handler,
		},
		{
			MethodName: "EscribirLogDes",
			Handler:    _BookService_EscribirLogDes_Handler,
		},
		{
			MethodName: "RecibirPropDatanode",
			Handler:    _BookService_RecibirPropDatanode_Handler,
		},
		{
			MethodName: "ChunkInfoLog",
			Handler:    _BookService_ChunkInfoLog_Handler,
		},
		{
			MethodName: "EnviarListaLibros",
			Handler:    _BookService_EnviarListaLibros_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RecibirChunks",
			Handler:       _BookService_RecibirChunks_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "book.proto",
}
