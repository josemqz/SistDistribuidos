package main

import (
	"fmt"
	"log"
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

var dNN = ""  //NameNode
var dA= "" 	  //DataNode
var dB = "" 
var dC = "" 


func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}


func verLibros(clienteNN book.BookServiceClient, ctx context.Context){
	
	listaLibros := clienteNN.EnviarListaLibros(ctx) //string or what
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
	if os.IsNotExist("Libros"){
		log.Println("Carpeta para libros descargados no existe")
		os.Mkdir("Libros")
		log.Println("Carpeta creada")
	}
	
	//crear nuevo archivo donde escribir los chunks del libro
	neoArchLibro := "./Libros/" + NombreArchLibro + "_reconstruido.pdf"
	_, err = os.Create(neoArchLibro)
	failOnError(err, "Error creando archivo de libro reconstruido")

	//abrir archivo
	file, err = os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	failOnError(err, "Error abriendo archivo de libro reconstruido")

	
//solicitar ubicaciones de chunks al namenode
	chunksInfo := clienteNN.ChunkInfoLog(ctx, &book.ChunksInfo{NombreLibro: nombreArchLibro})

	dirChunks := Strings.Fields(chunksInfo.info)
	totalPartsNum := len(dirChunks) //se usa?
	
	//estado de DataNodes
	a := chunksInfo.DatanodeA
	b := chunksInfo.DatanodeB
	c := chunksInfo.DatanodeC


// Reunir fragmentos ~~~

//conexión a los tres nodos
	
	//DN A
	if a{
		connA, err := grpc.Dial(dA + ":50501", grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Error en conexión con DataNodeA, no se podrá descargar libro")
		defer conn.Close()
		
		clienteA := book.NewBookServiceClient(connA)
	}

	//DN B
	if b{
		connB, err := grpc.Dial(dB + ":50502", grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Error en conexión con DataNodeB, no se podrá descargar libro")
		defer conn.Close()
		
		clienteB := book.NewBookServiceClient(connB)
	}

	//DN C
	if c{
		connC, err := grpc.Dial(dC + "50503", grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Error en conexión con DataNodeC, no se podrá descargar libro")
		defer conn.Close()
		
		clienteC := book.NewBookServiceClient(connC)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	var writePosition int32 = 0
	var jInfo []string
	var archChunk string
	var dirChunk string
	
	var msgChunk *book.Chunk

	//for j := uint64(0); j < totalPartsNum; j++ {
	for j in range(dirChunks){

		jInfo = strings.split(j, " ")
		archChunk = jInfo[0] 			//nombre de archivo del chunk
		dirChunk = jInfo[1]  			//dirección de DataNode

		//envía mensaje a nodo con chunk j y lo recibe
		switch dirChunk{
			case dA:
				msgChunk = clienteA.enviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
			case dB:
				msgChunk = clienteB.enviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
			case dC:
				msgChunk = clienteC.enviarChunkDN(ctx, &book.Chunk{NombreArchivo: archChunk})
			default:
				log.Fatalf("Dirección obtenida de propuesta inválida: %s", dirChunk)
		}

		log.Println("número de chunk: ", msgChunk.NumChunk) //solo para debug
		
		
		//obtener tamaño del chunk
		chunkStat := msgChunk.Contenido.Stat()
		chunkSize := chunkStat.Size()

		//arreglo de bytes para 
		buf := make([]byte, chunkSize)

		log.Println("Escribiendo en byte: [", writePosition, "]")
		writePosition += chunkSize

		//escribir contenido en buffer
		reader := bufio.NewReader(msgChunk.Contenido)
		_, err = reader.Read(buf)
		failOnError(err, "Error escribiendo chunk en buffer")

		//escribir contenido de buffer en el archivo
		_, err := file.Write(buf)
		failOnError(err, "Error escribiendo chunk en archivo para reconstruir")


		file.Sync() //flush to disk
		buf = nil

		log.Println(n, " bytes escritos")
		fmt.Println("Insertando parte [", j, "] en : ", NombreArchLibro, "_reconstruido.pdf")
	}

	file.Close()
}


func main() {

	var opcion int

	log.Println("-------------------------------------")
	log.Println("  1. Ver libros disponibles          |")
	log.Println("  2. Descargar libro				  |")
	log.Println("-------------------------------------\n")
	log.Print("Seleccionar opción: ")

	_, err = fmt.Scanf("%d", &opcion)

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
	connNN, err := grpc.Dial(dirNN, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()

	clienteNN := book.NewBookServiceClient(conn)
	log.Println("Conexión a NameNode realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//poder ver más veces los libros disponibles
	for (opcion == 1){

		verLibros(clienteNN, ctx)

		log.Println("-------------------------------------")
		log.Println("  1. Ver libros disponibles          |")
		log.Println("  2. Descargar libro				  |")
		log.Println("-------------------------------------\n")
		log.Print("Seleccionar opción: ")
	
		_, err = fmt.Scanf("%d", &opcion)
	
		for (err != nil) || (opcion != 1 && opcion != 2) {
	
			log.Println("Opción inválida\n")
	
			log.Println("-------------------------------------")
			log.Println("  1. Ver libros disponibles          |")
			log.Println("  2. Descargar libro				  |")
			log.Println("-------------------------------------\n")
			log.Print("Seleccionar opción: ")
	
			_, err = fmt.Scanf("%d", &opcion)
		}
	}

	if (opcion == 2){
		descargarLibro(clienteNN, ctx)
	}

}