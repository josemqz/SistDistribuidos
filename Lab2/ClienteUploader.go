package main

import (
	//"bufio"
	"fmt"
	"log"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)


func main() {

	var archLibro string

	log.Println("Nombre de archivo del libro a registrar: ")
	_, err := fmt.Scanf("%s", &archLibro)

	for (err != nil){

		log.Println("Opción inválida")
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
	nombreArchLibro := strings.Split(archLibro, ".")[0]

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 250000 // 250 kB
	
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


	file.Close()

}

	//enviar cada chunk a un datanode (cuál?)
