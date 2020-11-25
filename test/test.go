package main
//CLIENT

import (
	"log"
	"time"
	"context"
	"io"
	"os"
	
	"github.com/josemqz/SistDistribuidos/test/testp"
	"google.golang.org/grpc"
)

func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}

func main(){

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a Testp")
	defer conn.Close()

	client := testp.NewTestpServiceClient(conn)
	log.Println("Conexión realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	bookName := "Alicia_a_traves_del_espejo-Carroll_Lewis"
	fileName := "./srcBooks/" + bookName + ".pdf"

	err = SubirArchivo(client, ctx, fileName, bookName)
	failOnError(err,"Error enviando archivo")

}


func SubirArchivo(client testp.TestpServiceClient, ctx context.Context, fileN string, bookN string) (err error) {

	//abrir archivo a enviar
	file, err := os.Open(fileN)
	if (err != nil) {
		log.Fatalf("%s: %s\n", "no se pudo abrir archivo", err)
	}

	//abrir conexión basada en stream
	stream, err := client.RecibirBytes(ctx)
	if (err != nil) {
		log.Fatalf("%s: %s\n", "error llamando función RecibirBytes", err)
	}


	//buffer de tamaño máximo 250 kB
	buf := make([]byte, 250000)

	for {
		//escribir en buffer tanto como se pueda (<= 250 kB)
		n, err := file.Read(buf)

		if (err != nil){
			
			//terminar ciclo
			if err == io.EOF{
				break
			}
			
			return err
		}

		//enviar chunk
		log.Println("enviando chunk")
		stream.Send(&testp.MSG{
			Nombre: bookN,
			Arch: buf[:n]}) //con [:n] se envía <= 250 kB
	}

	//cerrar
	_, err = stream.CloseAndRecv()

	return nil
}