package main
//SERVER

import (
	"log"
	"net"
	"io"
	"os"
	"fmt"
	"bytes"
	"strconv"

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


// Función rpc para recibir chunks
func (s *server) RecibirBytes(stream testp.TestpService_RecibirBytesServer) (err error) {
	
	//número de chunk
	numChunk := 0

	for {
		
		log.Println("Recibiendo chunk")
		chunk, err := stream.Recv()
		
		if err != nil {

			if err == io.EOF {
				goto END
			}
			
			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
				return err
		}
		
		//crear archivo de chunk
		strNumChunk := strconv.Itoa(numChunk)
		neoArchLibro := "./NeoLibros/" + chunk.Nombre + "_" + strNumChunk
		_, err = os.Create(neoArchLibro)
		failOnError(err, "Error creando archivo de libro reconstruido")


		//abrir archivo
		file, err := os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		failOnError(err, "Error abriendo archivo de libro reconstruido")
		fmt.Println("Escribiendo chunk en disco")

		/*
		//buffer
		chunkBufferBytes := make([]byte, len(chunk.Arch))

		//escribir chunk en buffer
		reader := bytes.NewReader(chunk.Arch) //cambiar
		_, err = reader.Read(chunkBufferBytes)
		failOnError(err, "Error escribiendo chunk en buffer")
		*/

		n, err := file.Write(chunk.Arch)
		failOnError(err, "Error escribiendo chunk en archivo para reconstruir")
		fmt.Println(n, " bytes escritos")

		file.Close()
		chunk = nil
		file = nil
		//chunkBufferBytes = nil // reset or empty our buffer

		numChunk += 1
	}

	//Enviar confirmación al terminar el stream
	END:

	err = stream.SendAndClose(&testp.ACK{Ok: "ok"})
	return
}