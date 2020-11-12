package main

import (
	
)

func main() {


	//pedir libro (solicitar ubicaciones de chunks)
	//recibe direcciones de los chunks
	
	//pide chunk a cada id
	//armar libro con los chunks


	// Reunir fragmentos ~~~

	var archLibro string

	log.Println("Nombre de archivo del libro a registrar: ")
	_, err := fmt.Scanf("%s", &archLibro)

	for (err != nil){

		log.Println("Opción inválida")
		log.Println("Nombre de archivo del libro a registrar: ")
		_, err = fmt.Scanf("%s", &archLibro)

	}

	nombreArchLibro := strings.Split(archLibro, ".")[0]

	neoArchLibro := NombreArchLibro + "reconstruido.pdf"
	_, err = os.Create(neoArchLibro)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//set the neoArchLibro file to APPEND MODE!!
	// open files r and w

	file, err = os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

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

	file.Close()

}
