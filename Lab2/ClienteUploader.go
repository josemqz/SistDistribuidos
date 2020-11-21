package main

import (
	//"bufio"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
)


func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}


func enviarChunk(chunk byte[], dir string){

	//(cambiar dirección)
	conn, err := grpc.Dial(dir, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err, "Error en conexión a servidor")
	defer conn.Close()
	
	cliente := logis.NewBookServiceClient(conn)
	
	//stream??
	//funcion cliente.enviarChunk(chunk)

}


func subirLibro(){

	var archLibro string

	log.Println("Nombre de archivo del libro a registrar: ")
	_, _ := fmt.Scanf("%s", &archLibro)
	_, err = os.Stat("./Libros/" + archLibro)

	for os.IsNotExist(err){

		log.Println("Archivo no existe")
		log.Println("Nombre de archivo del libro a registrar: ")
		_, _ = fmt.Scanf("%s", &archLibro)
		_, err = os.Stat("./Libros/" + archLibro)

	}

	file, err := os.Open(archLibro)
	failOnError(err,"No se pudo abrir archivo de libro")
	
	defer file.Close()

	//escoger datanode para enviar chunks
	switch dn := rand.Intn(3)
	dn {
		case 0:
			dir := ""
		case 1:
			dir := ""
		case 2:
			dir := ""
		default:
			//error fatal wtf, no debería pasar
			failOnError(err,"Error fatal interno de Go (apareció un 3 en un random entre 0 y 2)")
	}

	
	/*var filename = "hello.blah"
	var extension = filepath.Ext(filename)
	var name = filename[0:len(filename)-len(extension)]*/

	//Nombre de archivo sin extensión
	nombreArchLibro := strings.Split(archLibro, ".")[0]

	//Tamaño de archivo
	fileInfo, _ := file.Stat()
	var fileSize int32 = fileInfo.Size()

	const fileChunk = 250000 // 250 kB

	//Calcular cantidad de fragmentos para el archivo
	totalPartsNum := uint32(math.Ceil(float64(fileSize) / float64(fileChunk)))
	fmt.Printf("Cantidad de chunks: %d\n", totalPartsNum)


	for i := uint32(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize - int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		enviarChunk(partBuffer, dir) //sure??

		fmt.Println("Fragmento", i, "enviado a DataNode", dn)

	/*
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
	*/
	}

	file.Close()
}


func main() {

	subirLibro()

}

