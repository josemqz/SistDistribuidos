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
	"strings"

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


//distribuido
func manejarConflictoDist(){

	//Ricart y Agrawala
	
	var state string
	var Ci []int
	
	//On initialization
		state = "RELEASED"

	//To enter the section
		state = "WANTED"							//-\
		//Multicast request to all processes		//---> Request processing deferred here
		//T:= request timestamp;					//_/
		//Wait until (number of replies received = (N-1));
		state = "HELD"

	//On receipt of a request <T_i,p_i> at p_j (i!=j)
		if (state == HELD ) or (state == WANTED and (T_j,p_j) < (T_i,p_i))
		{
			//queue request from p_i without replying
		}
		else{
			//reply immediately to p_i
	}

//To exit the critical section
	state := RELEASED
	//reply to any queued requests
	
}


//enviar chunk a Cliente Downloader
func (s *server) enviarChunkDN(ctx context.Context, *book.Chunk) (*book.Chunk, error){
	
	file, err = os.Open("./Chunks/" + Chunk.NombreArchivo)
	failOnError(err,"No se pudo abrir archivo de chunk a enviar a Cliente Downloader")
	defer file.Close()

	buf := make([]byte, 250000)
	n, err := file.Read(buf)
	failOnError(err,"No se pudo leer del archivo de chunk en buffer")

	file.Close()
	return &book.Chunk{Contenido: buf[:n]}
}


//recibir chunks de DataNode al distribuirlos
func (s *server) recibirChunksDN(ctx context.Context, c *book.Chunk) (*book.ACK, error){

	f, err := os.OpenFile(c.NombreArchivo, os.O_WRONLY|os.O_APPEND, 0644)
	failOnError(err, "Error creando archivo")
	defer f.Close()


	//obtener tamaño del chunk
	chunkStat := c.Contenido.Stat()
	chunkSize := chunkStat.Size()

	//arreglo de bytes para traspasar chunk
	buf := make([]byte, chunkSize)

	//escribir contenido en buffer
	reader := bufio.NewReader(c.Contenido)
	_, err = reader.Read(buf)
	failOnError(err, "Error escribiendo chunk recibido en buffer")

	//escribir contenido de buffer en el archivo
	_, err := f.Write(buf)
	failOnError(err, "Error escribiendo chunk en archivo")


	f.Sync() //flush to disk
	f.Close()
	buf = nil
	fmt.Println("Escritura exitosa")

	return nil
}


//enviar 1 chunk a otro DataNode
func enviarChunk(archivoChunk string, ip string){
	
	//conexión
	connDN, err := grpc.Dial(ip + "", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a DataNode")
	defer connDN.Close()
	
	clientDN := book.NewBookServiceClient(connDN)
	log.Println("Conexión a DataNode realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()


	//abrir archivo a enviar
	file, err := os.Open("./Chunks/" + archivoChunk)
	failOnError(err,"No se pudo abrir archivo de chunk a enviar a DataNode")
	defer file.Close()

	buf := make([]byte, 250000)
	n, err := file.Read(buf)
	failOnError(err,"No se pudo leer del archivo de chunk en buffer")


	okEnviar, err := clientDN.recibirChunksDN(ctx, &book.Chunk{Contenido: buf[:n], NombreArchivo: archivoChunk})

	connDN.close()
	return nil

}


//enviar chunks a otros nodos
func distribuirChunks(prop string){

	var line []string
	var c string

	line = strings.Split(prop,"\n")[1:]

	for info := range line {

		c = strings.split(info, " ")

		enviarChunk(c[0], c[1])
	}
}


//genera una nueva propuesta considerando solo los nodos activos
func nuevaPropuesta2(dn string, cantChunks int, nombreLibro string) (string, bool, bool){

	b := false
	c := false

	//inicio de propuesta
	Prop := nombreLibro + " " + cantChunks + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, cantChunks)

	for i := 0; i < cantChunks; i++{
		intProp[i] = rand.Intn(2)
	}
	
	var dAct string
	for i = 0; i < cantChunks; i++{ 
		
		if (dn == dB) {
			switch intProp[i]{
			case 0:
				dAct = dA
			case 1:
				dAct = dC
			}
			b = true
		}

		if (dn == dC) {
			switch intProp[i]{
			case 0:
				dAct = dA
			case 1:
				dAct = dB
			}
			c = true
		}

		Prop += nombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"
		
		/*"nombre_archivo_chunk_0 ipDNx\n
		 nombre_archivo_chunk_1 ipDNy\n
		 nombre_archivo_chunk_2 ipDNz..." */
	}

	return Prop, b, c
}


//verifica si hay un DataNode caido
//name: nombre de datanode ("datanode X")
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


//retorna true si acepta propuesta y false si rechaza
func analizarPropuesta(prop *book.PropuestaLibro) (bool, string){
	
	log.Println("Analizando la propuesta...")

	//prop sabe cuales datanodes se pretenden usar en la propuesta
	//por ejemplo, si prop.DatanodeA==true, entonces se usaría en la propuesta
	//pero si esta caído se rechaza la propuesta

	//revisa si hay un DataNode de la propuesta caído
	if (prop.DatanodeB && !checkDatanode(dB, "datanode B")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, dB
	} 
	if (prop.DatanodeC && !checkDatanode(dC, "datanode C")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, dC
	}

	return true, ""
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


//enviar una propuesta a otro DataNode (descentralizado)
//retorna la respuesta de este y un DataNode caído, en caso de hallar uno
func enviarPropuestaDN(dir string, prop *book.PropuestaLibro) (bool, string) {

	//conexión
	connDN, err := grpc.Dial(dir, grpc.WithInsecure(), grpc.WithBlock())
	//failOnError(err, "Error en conexión a NameNode")
	defer connDN.Close()
	if (err != nil){
		return false, "nil"
	}

	clientDN := book.NewBookServiceClient(connDN)
	log.Println("Conexión realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	
	resp, DNcaido := clientDN.RecibirPropuesta(ctx, prop) //check
	
	connDN.close()
	return resp, DNcaido
}


//enviar propuesta a NameNode (descentralizado)
func EnviarPropNNDes(prop *book.PropuestaLibro) {
	//conexión
	connNN, err := grpc.Dial("localhost:50050", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()
	
	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	//propuesta final (book.PropuestaLibro)
	ACK, err = clientNN.escribirLogDes(ctx, prop)
	log.Println(ACK.Ok)
	failOnError(err, "")
	
	connNN.close()
}


//enviar propuesta a NameNode (centralizado)
func EnviarPropNNCen(prop *book.PropuestaLibro) (string, bool, bool, bool){
	
	//conexión
	connNN, err := grpc.Dial("localhost:50050", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()
	
	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	//propuesta final (book.PropuestaLibro)
	propFinal := clientNN.recibirPropDatanode(ctx, prop)
	
	connNN.close()
	return propFinal.Propuesta, propFinal.DatanodeA, propFinal.DatanodeB, propFinal.DatanodeC
}


//Genera una nueva propuesta con una lista generada aleatoriamente
//Recibe el mensaje con la propuesta original para tener
//el nombre del libro y la cantidad de chunks
func generarPropuesta(cantChunks int, nombreLibro string) (string, bool, bool, bool) {

	//estados de DataNodes
	var a = false
	var b = false
	var c = false

	//inicio de propuesta
	Prop := NombreLibro + " " + cantChunks + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, cantChunks)

	for i := 0; i < cantChunks; i++{
		intProp[i] = rand.Intn(3)
	}
	
	var dAct string //dirección correspondiente 
	for i = 0; i < cantChunks; i++{ 
		
		switch intProp[i]{
			case 0:
				dAct = dA
				a = true
			case 1:
				dAct = dB
				b = true
			case 2:
				dAct = dC
				c = true
		}
		
		Prop += NombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"
		
		/*"nombre_archivo_chunk_0 ipDNx\n
		 nombre_archivo_chunk_1 ipDNy\n
		 nombre_archivo_chunk_2 ipDNz..." */
	}

	/*return &book.PropuestaLibro{NombreLibro: nombreLibro,
								Propuesta: Prop, 
								estadoP: false,  			//es necesario?? < < < < < < < < < <
								CantChunks: cantChunks,
								DatanodeA: a,
								DatanodeB: b,
								DatanodeC: c}*/
	return Prop, a, b, c
}


//Recibir chunks de Cliente Uploader
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
		neoArchLibro := "./Chunks/" + chunk.NombreLibro + "_" + strNumChunk
		_, err = os.Create(neoArchLibro)
		failOnError(err, "Error creando archivo de libro reconstruido")

		//abrir archivo
		file, err := os.OpenFile(neoArchLibro, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		failOnError(err, "Error abriendo archivo de libro reconstruido")
		fmt.Println("Escribiendo chunk en disco")

		//escribir en archivo
		n, err := file.Write(chunk.Arch)
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

	_ = stream.SendAndClose(&book.ACK{Ok: "ok"})
	

//Generar propuesta, escribir en log según algoritmo y distribuir entre DataNodes
	
	//prop: propuesta; b, c: estados de DataNodes B y C
	prop, a, b, c := generarPropuesta(numChunk, nombreL)

//CENTRALIZADO
	if (tipo_al == c){
		
		//recibe prop, que puede ser la misma propuesta o una nueva del NameNode
		//CHECK < < < < < < son necesarios a, b, c?
		prop, a, b, c = EnviarPropNNCen(&book.PropuestaLibro{NombreLibro: NombreL, Propuesta: prop, CantChunks: numChunk)

//DESCENTRALIZADO
	} else {
		
		//enviar a demás nodos
		respB, DNcaidoB := enviarPropuestaDN(dB, &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
		respC, DNcaidoC := enviarPropuestaDN(dC, &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})

		//CHECK !! <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
		//generar nuevas propuestas hasta que se apruebe una
		for {

			//DNcaidoX es "" cuando respX == true (se aprueba prop y no hay nodo caído)
			//DNcaidoX es "nil" cuando hubo un problema de conexión (no se aprueba ni rechaza prop)
			if !respB && (DNcaidoB != "nil") {
				prop, b, c = nuevaPropuesta2(DNcaidoB, numChunk, nombreLibro)
				respB, DNcaidoB := enviarPropuestaDN(dB, &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})

			} else if !respC {
				prop, b, c = nuevaPropuesta2(DNcaidoC, numChunk, nombreLibro)
				respC, DNcaidoC := enviarPropuestaDN(dC, &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
			
			} else {
				break
			}
		}
		
		//enviar propuesta ya aprobada
		enviarPropNNDes(prop)
	}
	
	distribuirChunks(prop)

	return &book.ACK{Ok: "ok"}, nil
}


//funciones de conexión


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