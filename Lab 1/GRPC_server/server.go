package main

import (
	"context"
	"fmt"
	"log"
	"net"
 
	"github.com/josemqz/SistDistribuidos"
	"google.golang.org/grpc"
)

type server struct {
	logis.UnimplementedLogisServiceServer
}

func main() {
	log.Println("Server running ...")
 
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})
 
	log.Fatalln(srv.Serve(lis))
}

func (s *server) PedidoCliente(ctx context.Context, pedido *logis.Pedido) (*logis.Recibo, error) {
	log.Println("Pedido recibido")
 
	return &logis.Recibo{msg: "recibido"}, nil
}