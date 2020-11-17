package main

import (
	"fmt"
	"log"
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func descargarLibro(){

	var archLibro string

	log.Println("Nombre de archivo del libro a descargar: ")
	_, err := fmt.Scanf("%s", &archLibro)

	for (err != nil){

		log.Println("Texto ingresado inválido")
		log.Println("Nombre de archivo del libro a descargar: ")
		_, err = fmt.Scanf("%s", &archLibro)

	}

	//nombre sin extensión
	nombreArchLibro := strings.Split(archLibro, ".")[0]

	
	//verificar si existe carpeta con libros descargados
	if os.IsNotExist("Neolibros"){
		log.Println("Carpeta para libros descargados no existe")
		os.Mkdir("Neolibros")
		log.Println("Carpeta creada")
	}
	
	//crear nuevo archivo donde escribir los chunks del libro
	neoArchLibro := "./NeoLibros/" + NombreArchLibro + "_reconstruido.pdf"
	_, err = os.Create(neoArchLibro)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//abrir archivo
	file, err = os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	
//solicitar ubicaciones de chunks al namenode
	
	//totalPartsNum := chunksinfo.size
	//chunksinfo.info : "chunck1 ip1\nchunk2 ip2\n..."

// Reunir fragmentos ~~~

	// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
	// defer file.Close()

	// just information on which part of the new file we are appending
	var writePosition int32 = 0

	for j := uint64(0); j < totalPartsNum; j++ {

		//obtiene dirección de chunk j
		//envía mensaje a nodo con chunk j
		//recibe chunk y lo escribe en libro

		//obtener tamaño del chunk
		//


		/*
		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant // pues noso3 sí jaja

		var chunkSize int32 = chunk.Size
		*/

		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Escribiendo en byte: [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(chunk.Contenido)
		_, err = reader.Read(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		n, err := file.Write(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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




}
