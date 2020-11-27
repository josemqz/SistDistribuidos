package main

import (
	"os"
	"io"
	"net"
	"log"
	"fmt"
	"bytes"
	"strconv"
	"math/rand"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

dNN := ""
dA := "" //ip maquina virtual datanode A
dB := "" //ip maquina virtual datanode B
dC := "" //ip maquina virtual datanode C
//dAct := "" //ip nodo actual
//borrar la variable del nodo al que el actual corresponda arriba


func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}


func manejarConflictoDist(){

	//implementar Ricart y Agrawala
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



//enviar chunks a otros nodos
func distribuirChunks(dn string, prop string){

	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("no se pudo conectar: %v", err)
	}
	defer conn.Close()
	c := book.NewBookServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	//leer propuesta

}


//Genera una nueva propuesta con una lista generada aleatoriamente
//Recibe el mensaje con la propuesta original para tener
//el nombre del libro y la cantidad de chunks
func generarPropuesta(cantChunks int, nombreLibro string) string{

	//variables para saber qué nodos participan
	var DNA = false
	var DNB = false
	var DNC = false

	//inicio de propuesta
	Prop := NombreLibro + " " + cantChunks + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, cantChunks)

	for i := 0; i < cantChunks; i++{
		intProp[i] = rand.Intn(3)
	}
	
	var dAct string
	for i = 0; i < cantChunks; i++{ 
		
		switch intProp[i]{
		case 0:
			dAct = dA
			DNA = true
		case 1:
			dAct = dB
			DNB = true
		case 2:
			dAct = dC
			DNC = true
		}
		
		Prop += NombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"
		
		/*"nombre_archivo_chunk_0 ipDNx\n
		 nombre_archivo_chunk_1 ipDNy\n
		 nombre_archivo_chunk_2 ipDNz..." */
	}

	return &book.PropuestaLibro{NombreLibro: nombreLibro, 
								Propuesta: Prop, 
								estadoP: false, 
								CantChunks: cantChunks,
								DatanodeA: DNA,
								DatanodeB: DNB,
								DatanodeC: DNC}
}


func main() {

//conexion
//listen (go)
	//DataNode
	//DataNode

	//ClienteUploader //recibir stream de chunks
	//ClienteDownloader //entregar chunk


//clientes
	//DataNode //enviar propuesta | verificar estado de nodos (verificar propuesta) | enviar chunks
	//DataNode

}


//enviar propuesta a Datanode (centralizado)
func EnviarPropNameNode(prop *book.PropuestaLibro) error{

	//conexión
	connNN, err := grpc.Dial("localhost:50050", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()

	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()


	clientNN.recibirPropDatanode(ctx, prop)

	connNN.close()
	return nil
}


//función rpc para recibir chunks de Cliente Uploader
func (s *server) RecibirChunks(ctx context.Context, stream book.BookService_RecibirChunksServer) (*book.ACK, error) {
	
	//tipo de algoritmo a utilizar
	var tipo_al string
	var nombreL string
	
	//número de chunk
	numChunk := 0

	for {
		
		log.Println("Recibiendo chunk")
		chunk, err := stream.Recv()
		
		if err != nil {

			if err == io.EOF {
				goto END
			}
			
			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
				return err
		}

		//crear archivo de chunk
		strNumChunk := strconv.Itoa(numChunk)
		neoArchLibro := "./NeoLibros/" + chunk.NombreLibro + "_" + strNumChunk
		_, err = os.Create(neoArchLibro)
		failOnError(err, "Error creando archivo de libro reconstruido")

		//abrir archivo
		file, err := os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		failOnError(err, "Error abriendo archivo de libro reconstruido")
		fmt.Println("Escribiendo chunk en disco")

		//escribir en archivo
		n, err := file.Write(EnviarPropNameNode(prop)chunk.Arch)
		failOnError(err, "Error escribiendo chunk en archivo para reconstruir")
		fmt.Println(n, " bytes escritos")

		if (numChunk == 0){
			tipo_al = chunk.Algoritmo
			nombreL = chunk.NombreLibro
		}

		file.Close()
		chunk = nil
		file = nil

		numChunk += 1
	}

	//Enviar confirmación al terminar el stream
	END:

	err = stream.SendAndClose(&book.ACK{Ok: "ok"})
	
//Generar propuesta y enviar

	//centralizado
	if (tipo_al == c){
		prop := generarPropuesta(numChunk, nombreLibro)
		EnviarPropNameNode(prop)

	//descentralizado
	} else{
		//enviar a demas pcs
		enviarPropuestaDN(dA, prop.Propuesta) //por mientras, hay que cambiar a futuro
		enviarPropuestaDN(dB, prop.Propuesta) //por mientras, hay que cambiar a futuro
	}
	
	//distribuir chunks

	return &book.ACK{Ok: "ok"}, nil
}


//verifica si hay un DataNode caido
func checkDatanode(dn string, name string){

	deadline = 2 //segundos que espera para ver si hay conexión

	conn, err := grpc.Dial(dn, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(deadline)*time.Second)
    if err != nil {
		log.Printf("No se pudo conectar con %v: %v", name, err)  // CHECK <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
		return false
    }
	defer conn.Close()
	
	c := book.NewBookServiceClient(conn) //necesario?? es solo chekear pero no aun conectar 4real

	conn.Close()
	return true
}


func buscarNodoCaido() (string, error){

	//revisa si hay un datanode caído
	if (!checkDatanode(dA, "datanode A")){
		return dA, nil
	}
	if (!checkDatanode(dB, "datanode B")){
		return dB, nil
	} 
	if (!checkDatanode(dC, "datanode C")){
		return dC, nil

	} else{
		e := "no hay nodos caidos"
	}

	return e, nil
}


//retorna true si acepta propuesta y false si rechaza
func analizarPropuesta(prop *book.PropuestaLibro) (bool, int){
	
	log.Println("Analizando la propuesta...")

	//prop sabe cuales datanodes se pretenden usar en la propuesta
	//por ejemplo, si prop.DatanodeA==true, entonces se usaría en la propuesta
	//pero si esta caído se rechaza la propuesta

	//revisa si hay un datanode de la propuesta caído
	if (prop.DatanodeA && !checkDatanode(dA, "datanode A")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, 0
	}
	if (prop.DatanodeB && !checkDatanode(dB, "datanode B")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, 1
	} 
	if (prop.DatanodeC && !checkDatanode(dC, "datanode C")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, 2
	}

	return true, -1
}


//aceptar o rechazar propuesta de distribución de chunks (descentralizado)
func (s *server) RecibirPropuesta(ctx context.Context, prop *book.PropuestaLibro) (*book.RespuestaP, error){
	
	resp, DN := analizarPropuesta(prop)

	if resp{
		return RespuestaP{Respuesta: true}

	} else{
		return RespuestaP{Respuesta: false, DNcaido: DN}
	}
}


//enviar propuesta a otros DataNodes (descentralizado)
func enviarPropuesta(dir string, prop *book.PropuestaLibro) {

	//conexión
	connDN, err := grpc.Dial(dir, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connDN.Close()

	clientDN := book.NewBookServiceClient(connDN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()


	clientDN.RecibirPropuesta(ctx, prop) //check

	connDN.close()

}