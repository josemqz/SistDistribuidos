package main

import (
	"fmt"
	"log"
)

var tipo_al string

type RegistroLog struct{  //utilizada en func logNameNode - en caso de aprobar crear mje en proto
	nombreLibro string
	numLibro int32
	ipChunk string
}


func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}

/*
func logNameNode(reg RegistroLog, npartesA int, npartesB int, npartesC int, ipA string, ipB string, ipC string) {

	ap, err := os.OpenFile("registro_log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	failOnError(err,"Error abriendo Log")

	int x := reg.numLibro //numero del libro para identificar parte_x_i

	defer ap.Close()
	
	if (reg.ipChunk == ipA){// escribir en datanode A
		if _, err := ap.WriteString(reg.nombreLibro + " " + npartesA + "\n"); err != nil {
			log.Println(err)
			for (i = 0; i < npartesA; i += 1){
				ap.WriteString("parte_" + x + "_" + i + " " + ipA + "\n")

			}
		}
	}
	if (reg.ipChunk == ipB){// escribir en datanode B
		if _, err := ap.WriteString(reg.nombreLibro + " " + npartesB + "\n"); err != nil {
			log.Println(err)
			for (i = 0; i < npartesB; i += 1){
				ap.WriteString("parte_" + x + "_" + i + " " + ipB + "\n")
	
			}
		}
	}
	if (reg.ipChunk == ipC){// escribir en datanode C
		if _, err := ap.WriteString(reg.nombreLibro + " " + npartesC + "\n"); err != nil {
			log.Println(err)
			for (i = 0; i < npartesC; i += 1){
				ap.WriteString("parte_" + x + "_" + i + " " + ipC + "\n")
	
			}
		}
	}		
}
*/


func main() {


	log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")

	_, err := fmt.Scanf("%d", &tipo_al)

	for (err != nil){

		log.Println("Tipo ingresado inválido!\n")
		log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")
		
		_, err = fmt.Scanf("%d", &tipo_al)
	}


}

//función para recibir entradas para el log
//dependiendo del algoritmo 