package main

import (
	"fmt"
)
var tipo_al int


func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func generarPropuesta(){
	/*
	propuesta inicial:
	cantChunks = tamañoLibro/cantidad de datanodes
	reparte equitativamente y en orden los chunks, por ej si son 3 datanodes: [0,3,6,...],[1,4,7,...],[2,5,8,...]
	
	propuesta alternativa si es rechazada la anterior:
	cantChunks = tamañoLibro/cantidad de datanodes
	reparte aleatoriamente pero equitativamente los chunks (de esta forma al ser aleatorio será diferente la propuesta cada vez que sea rechazada)
	for cantChunks/3:
	designarChunkAleatorio(datanode1) --> con cada datanode, si no es multiplo de 3 un datanode quedará con un chunk de más/menos
		--> con cada iteración se elimina el string/id del chunk de la lista
	*/
}


func manejarConflictoDist(){

	/*
	implementar Ricart y Agrawala
	
	*/

}

func manejarConflictoCentr(){

	/*
	aleatoriamente dar paso a un datanode primero y otro después
	
	*/

}


func main() {

	log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")

	_, err := fmt.Scanf("%d", &tipo_al)

	for (err != nil){

		log.Println("Tipo ingresado inválido!\n")
		log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")
		
		_, err = fmt.Scanf("%d", &tipo_al)
	}
	

	//

}


//función para recibir chunks de cliente
//debe proponer forma de distribución de los chunks
//dependiendo del algoritmo...

//función para aceptar o rechazar propuesta de distribución de chunks