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

dA := "" //ip maquina virtual datanode A
dB := "" //ip maquina virtual datanode B
dC := "" //ip maquina virtual datanode C

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

	//nombre de archivo sin extensión
	nombreArchLibro := strings.Split(archLibro, ".")[0]

	return nombreArchLibro
}


func subirLibro(client book.BookServiceClient, ctx context.Context, fileN string, bookN string, tipo_al string) (err error) {

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
	fmt.Println("  2. Descentralizada		|")
	fmt.Println("----------------------------\n")
	fmt.Print("Seleccionar opción: ")

	_, err = fmt.Scanf("%d", &opcion)

	for (err != nil) || (opcion != 1 && opcion != 2) {

		log.Println("Opción inválida\n")
		
		fmt.Println("Escoger tipo de distribución")
		fmt.Println("----------------------------")
		fmt.Println("  1. Centralizada          |")
		fmt.Println("  2. Descentralizada		|")
		fmt.Println("----------------------------\n")
		fmt.Print("Seleccionar opción: ")

		_, err = fmt.Scanf("%d", &opcion)
	}

	if (opcion == 1){
		tipo_al = "c"
	} else{
		tipo_al = "d"
	}

	//escoger datanode al que enviar chunks
	switch dn := rand.Intn(3)
	dn {
		case 0:
			dir := dA
		case 1:
			dir := dB
		case 2:
			dir := dC
	}

	//conn, err := grpc.Dial(dir, grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	
	//si el nodo está caído
	if (err != nil){
		switch dn {
			case 0:
				conn, err := grpc.Dial(dB, grpc.WithInsecure(), grpc.WithBlock())
				failOnError(err,"Error en conexión a NameNode")
			case 1:
				conn, err := grpc.Dial(dC, grpc.WithInsecure(), grpc.WithBlock())
				failOnError(err,"Error en conexión a NameNode")
			case 2:
				conn, err := grpc.Dial(dA, grpc.WithInsecure(), grpc.WithBlock())
				failOnError(err,"Error en conexión a NameNode")
		}
	}
	defer conn.Close()

	client := testp.NewTestpServiceClient(conn)
	log.Println("Conexión realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()


	//nombre de archivo
	nombreLibro := ArchivoLibro()
	nombreArch := "./Libros/" + nombreLibro + ".pdf"


	//subir libro
	err = subirLibro(client, ctx, nombreArch, nombreLibro, tipo_al)
	failOnError(err,"Error subiendo libro")

}

