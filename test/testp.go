package main
//SERVER

import (
	"log"
	"net"
	"io"
	"os"
	"fmt"
	"bytes"
	
	"github.com/pkg/errors"
	"github.com/josemqz/SistDistribuidos/test/testp"
	"google.golang.org/grpc"
)


type server struct {
	testp.UnimplementedTestpServiceServer
}

func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}

func main(){

	listenCliente, err := net.Listen("tcp", "localhost:50051")
	failOnError(err, "error de conexion con cliente")
	log.Println("Esperando conexión...")

	srv := grpc.NewServer()
	testp.RegisterTestpServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenCliente))

}


// Upload implements the Upload method of the GuploadService interface which is 
// responsible for receiving a stream of chunks that form a complete file.
func (s *server) RecibirBytes(stream testp.TestpService_RecibirBytesServer) (err error) {
	
	// while there are messages coming
	for {
		
		log.Println("recibiendo chunk")
		chunk, err := stream.Recv()
		
		if err != nil {

			if err == io.EOF {
				goto END
			}
			
			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
				return err
			}
		
		writePosition := 0

		neoArchLibro := "./NeoLibros/" + "blablabla" + "_reconstruido.pdf"
		_, err = os.Create(neoArchLibro)
		failOnError(err, "Error creando archivo de libro reconstruido")
	
		//abrir archivo
		file, err := os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		failOnError(err, "Error abriendo archivo de libro reconstruido")

		fmt.Println("Escribiendo en byte: [", writePosition, "] bytes")

		chunkBufferBytes := make([]byte, len(chunk.Arch))

		// read into chunkBufferBytes
		reader := bytes.NewReader(chunk.Arch) //cambiar
		_, err = reader.Read(chunkBufferBytes)
		failOnError(err, "Error escribiendo chunk en buffer")


		n, err := file.Write(chunkBufferBytes)
		failOnError(err, "Error escribiendo chunk en archivo para reconstruir")
		

		file.Sync() //flush to disk

		//chunkBufferBytes = nil // reset or empty our buffer
		
		chunk = nil
		
		fmt.Println("Written ", n, " bytes")

		//fmt.Println("Insertando parte [", j, "] en : ", NombreArchLibro, "_reconstruido.pdf")

	}

	END:
	// once the transmission finished, send the confirmation if nothing went wrong

	err = stream.SendAndClose(&testp.ACK{Ok: "ok"})
	// ...

	return
}