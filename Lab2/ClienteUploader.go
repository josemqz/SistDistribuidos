package main

import (
	"os"
	"fmt"
	"log"
	"io"
	"math/rand"
	"strings"
	"context"
	"time"
	
	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

var dA = "10.6.40.158" //ip maquina virtual datanode A
var dB = "10.6.40.159" //ip maquina virtual datanode B
var dC = "10.6.40.160" //ip maquina virtual datanode C


/* test local
var dA = "localhost"
var dB = "localhost"
var dC = "localhost"
*/


//manejo de errores
func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}


//ingresar nombre de libro a subir
func ArchivoLibro() string {

	var archLibro string

	log.Println("Nombre de archivo del libro a subir (incluyendo extensión): ")
	_, _ = fmt.Scanf("%s", &archLibro)
	_, err := os.Stat("./Libros/" + archLibro)

	//en caso que no exista el archivo
	for os.IsNotExist(err){

		log.Println("Archivo no existe")
		log.Println("Nombre de archivo del libro a registrar: ")
		_, _ = fmt.Scanf("%s", &archLibro)
		_, err = os.Stat("./Libros/" + archLibro)

	}

	//nombre de archivo sin extensión
	nombreArchLibro := strings.Split(archLibro, ".")[0]

	return nombreArchLibro
}


//función para enviar libro a DataNode, recibe el nombre del archivo del libro, 
//el nombre del archivo sin extensión y el tipo de algoritmo a utilizar los nodos
func subirLibro(client book.BookServiceClient, ctx context.Context, fileN string, bookN string, tipo_al string) error {

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
			Algoritmo: tipo_al,
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
		
//escoger tipo de algoritmo
	var opcion int

	fmt.Println("Escoger tipo de distribución")
	fmt.Println("----------------------------")
	fmt.Println("  1. Centralizada          |")
	fmt.Println("  2. Descentralizada	|")
	fmt.Println("----------------------------\n")
	fmt.Print("Seleccionar opción: ")

	_, err := fmt.Scanf("%d", &opcion)

	for (err != nil) || (opcion != 1 && opcion != 2) {

		log.Println("Opción inválida\n")
		
		fmt.Println("Escoger tipo de distribución")
		fmt.Println("----------------------------")
		fmt.Println("  1. Centralizada          |")
		fmt.Println("  2. Descentralizada	|")
		fmt.Println("----------------------------\n")
		fmt.Print("Seleccionar opción: ")

		_, err = fmt.Scanf("%d", &opcion)
	}

	var tipo_al string
	if (opcion == 1){
		tipo_al = "c"
	} else{
		tipo_al = "d"
	}

	var dir string
	var port string
	//escoger datanode al que enviar chunks
	dn := rand.Intn(3)
	switch dn {
		case 0:
			dir = dA
			port = ":50517"
		case 1:
			dir = dB
			port = ":50518"
		case 2:
			dir = dC
			port = ":50519"
	}

	//conn, err := grpc.Dial(dir, grpc.WithInsecure(), grpc.WithBlock())
	connDN, err := grpc.Dial(dir + port, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connDN.Close()
	
	//si el nodo está caído, intenta conectarse con otro arbitrario.
	//según los supuestos solo un nodo puede caerse, por lo 
	//que solo se maneja una vez un error de conexión
	if (err != nil){
		switch dn {
			case 0:
				connDN, err = grpc.Dial(dB + ":50518", grpc.WithInsecure(), grpc.WithBlock())
				failOnError(err,"Error en conexión a NameNode")
				defer connDN.Close()
			case 1:
				connDN, err = grpc.Dial(dC + ":50519", grpc.WithInsecure(), grpc.WithBlock())
				failOnError(err,"Error en conexión a NameNode")
				defer connDN.Close()
			case 2:
				connDN, err = grpc.Dial(dA + ":50517", grpc.WithInsecure(), grpc.WithBlock())
				failOnError(err,"Error en conexión a NameNode")
				defer connDN.Close()
		}
	}

	clientDN := book.NewBookServiceClient(connDN)
	log.Println("Conexión realizada con nodo (" + dir + ")\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()


	//nombre de archivo
	nombreLibro := ArchivoLibro()
	nombreArch := "./Libros/" + nombreLibro + ".pdf"


	//subir libro
	err = subirLibro(clientDN, ctx, nombreArch, nombreLibro, tipo_al)
	failOnError(err,"Error subiendo libro")

}

