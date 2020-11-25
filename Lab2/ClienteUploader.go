package main

import (
	"os"
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"context"
	"time"
	
	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)


func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}


//ingresar nombre de libro a subir
func ArchivoLibro() string, string {

	var archLibro string

	log.Println("Nombre de archivo del libro a subir: ")
	_, _ := fmt.Scanf("%s", &archLibro)
	_, err = os.Stat("./Libros/" + archLibro)

	//en caso que no exista el archivo
	for os.IsNotExist(err){

		log.Println("Archivo no existe")
		log.Println("Nombre de archivo del libro a registrar: ")
		_, _ = fmt.Scanf("%s", &archLibro)
		_, err = os.Stat("./Libros/" + archLibro)

	}

	//Nombre de archivo sin extensión
	nombreArchLibro := strings.Split(archLibro, ".")[0]

	return nombreArchLibro
}


func subirLibro(client book.BookServiceClient, ctx context.Context, fileN string, bookN string) (err error) {

	//abrir archivo a enviar
	file, err := os.Open(fileN)
	failOnError(err,"No se pudo abrir archivo de libro")
	defer file.Close()

	//abrir conexión basada en stream
	stream, err := client.RecibirChunks(ctx)
	if (err != nil) {
		log.Fatalf("%s: %s\n", "error llamando función RecibirBytes", err)
	}
	
	
	//buffer de tamaño máximo 250 kB
	buf := make([]byte, 250000)
	
	//Calcular cantidad de fragmentos para el archivo
	//totalPartsNum := uint32(math.Ceil(float64(fileSize) / float64(fileChunk)))
	//fmt.Printf("Cantidad de chunks: %d\n", totalPartsNum)

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
		stream.Send(&book.Chunk{
			Remitente: "cliente",
			NombreLibro: bookN,
			//NumChunk: 
			Contenido: buf[:n]}) //con [:n] se envía <= 250 kB
	}

	//cerrar
	file.Close()
	_, err = stream.CloseAndRecv()

	return nil
}


func main() {
		
	//escoger datanode para enviar chunks
	switch dn := rand.Intn(3)
	dn {
		case 0:
			dir := ""
		case 1:
			dir := ""
		case 2:
			dir := ""
		default:
			//error fatal wtf, no debería pasar
			failOnError(err,"Error fatal interno de Go (apareció un 3 en un random entre 0 y 2)")
	}

	//(cambiar dirección) -> dir
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a Testp")
	defer conn.Close()

	client := testp.NewTestpServiceClient(conn)
	log.Println("Conexión realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()


	//nombre de archivo
	nombreLibro := ArchivoLibro()
	nombreArch := "./Libros/" + nombreLibro + ".pdf"


	//subir libro
	err = subirLibro(client, ctx, nombreArch, nombreLibro)
	failOnError(err,"Error subiendo libro")

}

