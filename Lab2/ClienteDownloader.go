package main

import (
	"fmt"
	"log"
	"os"
	"context"
	"time"
	"strings"
	"bytes"
	"strconv"
	
	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)


var dNN = "10.6.40.157" //ip namenode
var dA = "10.6.40.158" //ip maquina virtual datanode A
var dB = "10.6.40.159" //ip maquina virtual datanode B
var dC = "10.6.40.160" //ip maquina virtual datanode C

/* test local
var dNN = "localhost"
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


func verLibros(clienteNN book.BookServiceClient, ctx context.Context){
	
	listaLibros, err := clienteNN.EnviarListaLibros(ctx, &book.ACK{Ok: "ok"})
	failOnError(err, listaLibros.Lista)

	fmt.Println("Lista de libros disponibles\n")
	fmt.Println(listaLibros.Lista)
}


func descargarLibro(clienteNN book.BookServiceClient, ctx context.Context){

	var archLibro string

	log.Println("Nombre de archivo del libro a descargar (sin extensión): ")
	_, err := fmt.Scanf("%s", &archLibro)

	for (err != nil){

		log.Println("Nombre ingresado inválido")
		log.Println("Nombre de archivo del libro a descargar (sin extensión): ")
		_, err = fmt.Scanf("%s", &archLibro)

	}

	//nombre sin extensión
	nombreArchLibro := strings.Split(archLibro, ".")[0]
	fmt.Println("Comenzando descarga de", nombreArchLibro, "\n")
	
	//verificar si existe carpeta con libros a descargar
	_, err = os.Stat("./Libros")
	if os.IsNotExist(err){
		log.Println("Carpeta para libros descargados no existe")
		os.Mkdir("./Libros", 0777)
		log.Println("Carpeta creada")
	}
	
	//crear nuevo archivo donde escribir los chunks del libro
	neoArchLibro := "./Libros/" + nombreArchLibro + "_reconstruido.pdf"
	file, err := os.OpenFile(neoArchLibro, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	failOnError(err, "Error abriendo archivo de libro reconstruido")

	//solicitar ubicaciones de chunks al namenode
	chunksInfo, err := clienteNN.ChunkInfoLog(ctx, &book.ChunksInfo{NombreLibro: nombreArchLibro})
	failOnError(err, "Error obteniendo direcciones del NameNode (quizás no está en Log)")
	fmt.Println("Ubicaciones de chunks obtenidas")

	dirChunks := strings.Fields(chunksInfo.Info)

	//estado de DataNodes
	var a = false
	var b = false
	var c = false

	//revisar previamente nodos que son parte de la propuesta
	for _, j := range dirChunks{

		if (j == dA){
			a = true
		} else if (j == dB){
			b = true
		} else if (j == dC){
			c = true
		} else {
			log.Fatalf("Dirección de nodo desconocido en Log de distribución de chunks de Libro")
		}
	}


// Reunir fragmentos ~~~

//conexión a los tres nodos
	
	var clienteA book.BookServiceClient
	var clienteB book.BookServiceClient
	var clienteC book.BookServiceClient

	var connA *grpc.ClientConn
	var connB *grpc.ClientConn
	var connC *grpc.ClientConn
	
	fmt.Println("Iniciando conexión con Datanodes para descarga de chunks\n")

	//DNA
	if a{
		
		fmt.Println("Conectando a DataNode A\n")
		
		connA, err = grpc.Dial(dA + ":50513", grpc.WithInsecure(), grpc.WithTimeout(80 * time.Second))
		failOnError(err, "Error en conexión con DataNode A, no se podrá descargar libro")
		defer connA.Close()
		
		clienteA = book.NewBookServiceClient(connA)
	}

	//DNB
	if b{
		
		fmt.Println("Conectando a DataNode B\n")
		
		connB, err = grpc.Dial(dB + ":50514", grpc.WithInsecure(), grpc.WithTimeout(80 * time.Second))
		failOnError(err, "Error en conexión con DataNode B, no se podrá descargar libro")
		defer connB.Close()
		
		clienteB = book.NewBookServiceClient(connB)
	}

	//DNC
	if c{
		
		fmt.Println("Conectando a DataNode C\n")
		
		connC, err = grpc.Dial(dC + ":50515", grpc.WithInsecure(), grpc.WithTimeout(80 * time.Second))
		failOnError(err, "Error en conexión con DataNode C, no se podrá descargar libro")
		defer connC.Close()
		
		clienteC = book.NewBookServiceClient(connC)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100) * time.Second)
	defer cancel()


	var writePosition int32 = 0
	var archChunk string
	
	var msgChunk *book.Chunk

	//número de chunk, dirección de chunk
	for i, dirChunk := range dirChunks{

		archChunk = nombreArchLibro + "_" +	strconv.Itoa(i) //nombre de archivo del chunk

		//envía mensaje a nodo con chunk j y lo recibe
		switch dirChunk{

			case dA:
				msgChunk, err = clienteA.EnviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
				failOnError(err, "Error en envío de un chunk desde un Datanode")

			case dB:
				msgChunk, err = clienteB.EnviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
				failOnError(err, "Error en envío de un chunk desde un Datanode")

			case dC:
				msgChunk, err = clienteC.EnviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
				failOnError(err, "Error en envío de un chunk desde un Datanode")

			default:
				log.Fatalf("Dirección obtenida de propuesta inválida: %s", dirChunk)
		}

		
		//obtener tamaño del chunk
		chunkSize := int32(len(msgChunk.Contenido))

		//arreglo de bytes para 
		buf := make([]byte, chunkSize)

		log.Println("Escribiendo en byte: [", writePosition, "]")
		writePosition += chunkSize

		//escribir contenido en buffer
		reader := bytes.NewReader(msgChunk.Contenido)
		_, err = reader.Read(buf)
		failOnError(err, "Error escribiendo chunk en buffer")

		//escribir contenido de buffer en el archivo
		n, err := file.Write(buf)
		failOnError(err, "Error escribiendo chunk en archivo para reconstruir")


		file.Sync() //flush to disk
		buf = nil

		log.Println(n, " bytes escritos")
		fmt.Println("Insertando parte [", i, "] en : ", nombreArchLibro, "_reconstruido.pdf")
	}

	file.Close()
	if a{
		connA.Close()
		log.Println("Conexión a Datanode A cerrada")
	}
	if b{
		connB.Close()
		log.Println("Conexión a Datanode B cerrada")
	}
	if c{
		connC.Close()
		log.Println("Conexión a Datanode C cerrada")
	}
}


func main() {

	var opcion int

	log.Println("-------------------------------------")
	log.Println("  1. Ver libros disponibles          |")
	log.Println("  2. Descargar libro                 |")
	log.Println("-------------------------------------\n")
	log.Print("Seleccionar opción: ")

	_, err := fmt.Scanf("%d", &opcion)

	for (err != nil) || (opcion != 1 && opcion != 2) {

		log.Println("Opción inválida\n")

		log.Println("-------------------------------------")
		log.Println("  1. Ver libros disponibles          |")
		log.Println("  2. Descargar libro                 |")
		log.Println("-------------------------------------\n")
		log.Print("Seleccionar opción: ")

		_, err = fmt.Scanf("%d", &opcion)
	}


	//conexión con NameNode
	connNN, err := grpc.Dial(dNN + ":50512", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()

	clienteNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión a NameNode realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100) * time.Second)
	defer cancel()

	//poder ver más veces los libros disponibles
	for (opcion == 1){

		verLibros(clienteNN, ctx)

		log.Println("-------------------------------------")
		log.Println("  1. Ver libros disponibles          |")
		log.Println("  2. Descargar libro                 |")
		log.Println("-------------------------------------\n")
		log.Print("Seleccionar opción: ")
	
		_, err = fmt.Scanf("%d", &opcion)
	
		for (err != nil) || (opcion != 1 && opcion != 2) {
	
			log.Println("Opción inválida\n")
	
			log.Println("-------------------------------------")
			log.Println("  1. Ver libros disponibles          |")
			log.Println("  2. Descargar libro                 |")
			log.Println("-------------------------------------\n")
			log.Print("Seleccionar opción: ")
	
			_, err = fmt.Scanf("%d", &opcion)
		}
	}

	if (opcion == 2){
		descargarLibro(clienteNN, ctx)
	}	
	
	connNN.Close()
}