package main
//SERVER

import (
	"context"
	"log"
	
	"github.com/josemqz/SistDistribuidos/test/testp"
	"google.golang.org/grpc"
)


type server struct {
	testp.UnimplementedTestpServiceServer
}

func main(){

	listenCliente, err := net.Listen("tcp", "localhost:50051")
	failOnError(err, "error de conexion con cliente")

	srv := grpc.NewServer()
	testp.RegisterTestpServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenCliente))

}


// Upload implements the Upload method of the GuploadService interface which is 
// responsible for receiving a stream of chunks that form a complete file.
func (s *ServerGRPC) RecibirBytes(stream testp.TestpService_UploadServer) (err error) {
	
	// while there are messages coming
	for {
		_, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return
		}
	}

	END:
	// once the transmission finished, send the confirmation if nothing went wrong

	err = stream.SendAndClose(&testp.ACK{ok: "ok"})
	// ...

	return
}