package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"bytes"
	"context"
	"time"

	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

/*
var dNN = "10.6.40.157" //ip namenode
var dA = "10.6.40.158" //ip maquina virtual datanode A
var dB = "10.6.40.159" //ip maquina virtual datanode B
var dC = "10.6.40.160" //ip maquina virtual datanode C
*/
// test local
var dNN = "localhost"
var dA = "localhost"
var dB = "localhost"
var dC = "localhost"



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

	log.Println("Nombre de archivo del libro a descargar (con extensión): ")
	_, err := fmt.Scanf("%s", &archLibro)

	for (err != nil){

		log.Println("Nombre ingresado inválido")
		log.Println("Nombre de archivo del libro a descargar (con extensión): ")
		_, err = fmt.Scanf("%s", &archLibro)

	}

	//nombre sin extensión
	nombreArchLibro := strings.Split(archLibro, ".")[0]

	
	//verificar si existe carpeta con libros a descargar
	_, err = os.Stat("./CD/Libros")
	if os.IsNotExist(err){
		log.Println("Carpeta para libros descargados no existe")
		os.Mkdir("./CD/Libros", 0777)
		log.Println("Carpeta creada")
	}
	
	//crear nuevo archivo donde escribir los chunks del libro
	neoArchLibro := "./CD/Libros/" + nombreArchLibro + "_reconstruido.pdf"
	_, err = os.Create(neoArchLibro)
	failOnError(err, "Error creando archivo de libro reconstruido")

	//abrir archivo
	file, err := os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	failOnError(err, "Error abriendo archivo de libro reconstruido")

	
//solicitar ubicaciones de chunks al namenode
	chunksInfo, err := clienteNN.ChunkInfoLog(ctx, &book.ChunksInfo{NombreLibro: nombreArchLibro})
	failOnError(err, chunksInfo.Info)

	dirChunks := strings.Fields(chunksInfo.Info)
	//totalPartsNum := len(dirChunks) //se usa?
	
	//estado de DataNodes
	var a = false
	var b = false
	var c = false

	//revisar nodos que son parte de la propuesta previamente
	for _, j := range dirChunks{
		
		jInfo := strings.Split(j, " ")

		if (jInfo[1] == dA){
			a = true
		} else if (jInfo[1] == dB){
			b = true
		} else if (jInfo[1] == dC){
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
	
	//DNA
	if a{
		connA, err := grpc.Dial(dA + ":50513", grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Error en conexión con DataNode A, no se podrá descargar libro")
		defer connA.Close()
		
		clienteA = book.NewBookServiceClient(connA)
	}

	//DNB
	if b{
		connB, err := grpc.Dial(dB + ":50514", grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Error en conexión con DataNode B, no se podrá descargar libro")
		defer connB.Close()
		
		clienteB = book.NewBookServiceClient(connB)
	}

	//DNC
	if c{
		connC, err := grpc.Dial(dC + "50515", grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Error en conexión con DataNode C, no se podrá descargar libro")
		defer connC.Close()
		
		clienteC = book.NewBookServiceClient(connC)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	var writePosition int32 = 0
	var jInfo []string
	var archChunk string
	var dirChunk string
	
	var msgChunk *book.Chunk

	//for j := uint64(0); j < totalPartsNum; j++ {
	for _, j := range dirChunks{

		jInfo = strings.Split(j, " ")
		archChunk = jInfo[0] 			//nombre de archivo del chunk
		dirChunk = jInfo[1]  			//dirección de DataNode

		//envía mensaje a nodo con chunk j y lo recibe
		switch dirChunk{
			case dA:
				msgChunk, _ = clienteA.EnviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
			case dB:
				msgChunk, _ = clienteB.EnviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
			case dC:
				msgChunk, _ = clienteC.EnviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
			default:
				log.Fatalf("Dirección obtenida de propuesta inválida: %s", dirChunk)
		}

		log.Println("número de chunk: ", msgChunk.NumChunk) //solo para debug
		
		
		//obtener tamaño del chunk
		/*chunkStat := msgChunk.Contenido.Stat()
		chunkSize := chunkStat.Size()*/
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
		fmt.Println("Insertando parte [", j, "] en : ", nombreArchLibro, "_reconstruido.pdf")
	}

	file.Close()
	if a{connA.Close()}
	if b{connB.Close()}
	if c{connC.Close()}
}


func main() {

	var opcion int

	log.Println("-------------------------------------")
	log.Println("  1. Ver libros disponibles          |")
	log.Println("  2. Descargar libro				  |")
	log.Println("-------------------------------------\n")
	log.Print("Seleccionar opción: ")

	_, err := fmt.Scanf("%d", &opcion)

	for (err != nil) || (opcion != 1 && opcion != 2) {

		log.Println("Opción inválida\n")

		log.Println("-------------------------------------")
		log.Println("  1. Ver libros disponibles          |")
		log.Println("  2. Descargar libro				  |")
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5) * time.Second)
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