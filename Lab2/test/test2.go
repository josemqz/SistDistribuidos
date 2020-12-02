//SERVER
package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

type server struct {
	book.UnimplementedBookServiceServer
}

func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func main(){

	log.Println("server corriendo")
	listenCD, err := net.Listen("tcp", "localhost" + ":50513")
	failOnError(err, "Error de conexión con cliente downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCD))
}
