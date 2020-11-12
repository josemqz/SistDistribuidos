package main

import (
		"bufio"
		"fmt"
		"io/ioutil"
		"math"
		"os"
		"strconv"
		"strings"
)


func main() {

	var archLibro string

	log.Println("Nombre de archivo del libro a registrar: ")
	_, err = fmt.Scanf("%s", &archLibro)

	for (err != nil){

		log.Println("Opción inválida\n")

		log.Println("Nombre de archivo del libro a registrar: ")
		_, err = fmt.Scanf("%s", &archLibro)

	}

	file, err := os.Open(archLibro)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	defer file.Close()


	/*var filename = "hello.blah"
	var extension = filepath.Ext(filename)
	var name = filename[0:len(filename)-len(extension)]*/

	//nombre de archivo sin extensión
	nombreArchLibro = strings.Split(archLibro, ".")[0]
	log.Println("nombre libro format: ", nombreArchLibro)

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 0.25 * (1 << 20) // 250 kB
	log.Println("chunks tamaño ", fileChunk)
	
	//Calcular cantidad de fragmentos para el archivo

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Cantidad de chunks: %d\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize - int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// nombre de fragmentos
		fileName := "fragmento_" + nombreArchLibro + "_" + strconv.FormatUint(i, 10)
		_, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// guardar en disco
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

		fmt.Println("Fragmento guardado en: ", fileName)
	}


	/*
	// Reunir fragmentos ~~~


	neoArchLibro := NombreArchLibro + "reconstruido.pdf"
	_, err = os.Create(newFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//set the newFileName file to APPEND MODE!!
	// open files r and w

	file, err = os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
	// defer file.Close()

	// just information on which part of the new file we are appending
	var writePosition int64 = 0

	for j := uint64(0); j < totalPartsNum; j++ {

		//read a chunk
		currentChunkFileName := "bigfile_" + strconv.FormatUint(j, 10)

		newFileChunk, err := os.Open(currentChunkFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer newFileChunk.Close()

		chunkInfo, err := newFileChunk.Stat()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant

		var chunkSize int64 = chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Appending at position : [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
		// write/save buffer to disk
		//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

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

		fmt.Println("Recombining part [", j, "] into : ", newFileName)
	}
*/
	// now, we close the newFileName
	file.Close()

}


/*package main

import (

)

func main() {

	//escoger libros?

	//fragmentar libro en chunks (250kb c/u)

	//enviar cada chunk a un datanode (cuál?)

}
*/