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
var dirNN = ""  //direccion NameNode
var dirDN1 = "" //dirección DataNodesde PC de
var dirDN1 = "" 
var dirDN1 = "" 


func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}


func verLibros(clienteNN book.BookServiceClient, ctx context.Context){
	clienteNN.EnviarListaLibros(ctx) //string or what

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
	if os.IsNotExist("Neolibros"){
		log.Println("Carpeta para libros descargados no existe")
		os.Mkdir("Neolibros")
		log.Println("Carpeta creada")
	}
	
	//crear nuevo archivo donde escribir los chunks del libro
	neoArchLibro := "./NeoLibros/" + NombreArchLibro + "_reconstruido.pdf"
	_, err = os.Create(neoArchLibro)
	failOnError(err, "Error creando archivo de libro reconstruido")

	//abrir archivo
	file, err = os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	failOnError(err, "Error abriendo archivo de libro reconstruido")

	
//solicitar ubicaciones de chunks al namenode
	chunksInfo := clienteNN.ChunkInfoLog(ctx, &book.ChunksInfo{NombreLibro: nombreArchLibro})

	dirChunks := Strings.Fields(chunksInfo.info)
	totalPartsNum := len(dirChunks)


// Reunir fragmentos ~~~

	//conexión a los tres nodos

	var writePosition int32 = 0

	//for j := uint64(0); j < totalPartsNum; j++ {
	for j in range(dirChunks){

		jInfo = strings.split(j, " ") 	//definir
		archChunk = jInfo[0] 			//definir
		DirChunk = jInfo[1]  			//definir

		//obtiene dirección de chunk j
		//envía mensaje a nodo con chunk j
		//recibe chunk y lo escribe en libro

		//obtener tamaño del chunk

		// calculate the bytes size of each chunk
		//var chunkSize int32 = chunk.Size

		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Escribiendo en byte: [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(chunk.Contenido)
		_, err = reader.Read(chunkBufferBytes)
		failOnError(err, "Error escribiendo chunk en buffer")
		

		n, err := file.Write(chunkBufferBytes)
		failOnError(err, "Error escribiendo chunk en archivo para reconstruir")
		

		file.Sync() //flush to disk

		// free up the buffer for next cycle
		// should not be a problem if the chunk size is small, but
		// can be resource hogging if the chunk size is huge.
		// also a good practice to clean up your own plate after eating

		chunkBufferBytes = nil // reset or empty our buffer

		fmt.Println("Written ", n, " bytes")

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


	if (opcion == 1){
		verLibros(clienteNN, ctx)

	} else{
		descargarLibro(clienteNN, ctx)
	}

}