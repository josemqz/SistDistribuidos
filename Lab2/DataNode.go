package main

import (
	"fmt"
	"log"
)

var tipo_al string
var num_chunks int32

//verificar si hay un pc caido para aceptar/rechazar la propuesta


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

	//R&A
	/*
	On initialization
		state := RELEASED

	To enter the section
		state := WANTED							//-\
		Multicast request to all processes		//---> Request processing deferred here
		T:= request's timestamp;				//_/
		Wait until (number of replies received = (N-1));
		state := HELD

	On receipt of a request <T_i,p_i> at p_j (i!=j)
		if (state == HELD ) or (state == WANTED and (T_j,p_j) < (T_i,p_i))
		{
			queue request from p_i without replying
		}
		else{
			reply immediately to p_i
	}

To exit the critical section
	state := RELEASED
	reply to any queued requests
	*/

}

func manejarConflictoCentr(){

	/*
	aleatoriamente dar paso a un datanode primero y otro después
	
	*/

}


func main() {

	//chequear IP de máquina en la que estamos

	log.Print("Ingresar tipo de algoritmo - c: centralizado / d: distribuido : ")

	_, err := fmt.Scanf("%d", &tipo_al)

	for (err != nil) || (tipo_al != "c" && tipo_al != "d"){ //chequear

		log.Println("Tipo ingresado inválido!\n")
		log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")
		
		_, err = fmt.Scanf("%d", &tipo_al)
	}

	//conexion
	

}


func (s *server) RecibirChunksInfo(ctx context.Context, ci *book.ChunksInfo) (*book.ACK, error) {
	
	//info de chunks a recibir (para saber la cantidad de chunks 
								//(para que genere la propuesta al recibir el último (ji)))

	num_chunks = ci.cantidadChunks

	return &book.ACK{ok: "ok"}, nil
}


//función rpc para recibir chunks de cliente
func (s *server) RecibirChunk(ctx context.Context, chunk *book.Chunk) (*book.ACK, error) {
	
	//debería guardar los chunks en archivos locales?? (a mí me tinca)

	//if lastChunk && está recibiendo del cliente:
		//generarPropuesta()


	//if algoritmo centralizado: >> enviar propuesta a namenode
	//else >> enviar a demas pcs

	return &book.ACK{ok: "ok"}, nil
}

//dependiendo del algoritmo...

//función para aceptar o rechazar propuesta de distribución de chunks