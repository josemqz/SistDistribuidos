package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"net"
	"time"
	//"bufio"
	"context"
	"strconv"
	"math/rand"
	"strings"
	"bytes"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

var dNN = "10.6.40.157" //ip namenode
var dActual = "10.6.40.158" //ip nodo actual
//var dA = "10.6.40.158" //ip maquina virtual datanode A
var dB = "10.6.40.159" //ip maquina virtual datanode B
var dC = "10.6.40.160" //ip maquina virtual datanode C

/* test local
var dNN = "localhost"
var dA = "localhost"
var dB = "localhost"
var dC = "localhost"
var dActual = "localhost"
*/

type server struct {
	book.UnimplementedBookServiceServer
}

func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}


//distribuido

/*
//Llama a funcion de Ricart y Agrawala con el datanode interesado en 
//contactar NN para escribir en log
func manejarConflictoDist(dn string){

	t := time.Now()
	tiempo := t.Format(“Mon Jan _2 15:04:05 2006”)
	
	connDN, errDN := grpc.Dial(dn, grpc.WithInsecure())
	if errDN != nil{
		log.Printf("No se pudo conectar a DB: %v", errDN)
	}
	defer connDN.Close()

	clientDN := book.NewBookServiceClient(connDN)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	replyDN, errDN2 := clientDN.RicAwla(cxt,*//* book.ExMutua*//*)
	if errDN2 != nil {
		log.Printf("problema con request: %v", errDN2)
	}
	return nil
}*/

/*
func (s *server) RicAwla(ctx context.Context, p *book.ExMutua) (*//**//*,error){

	//lista de datanodes en espera para escribir ?

	//func Parse(layout, value string) (Time, error)
	t := time.Now()	
	tiempo, et := time.Parse(time.ANSIC, p.tiempo) //mje string lo pasa a time
	
	// ******** if - alguna comparacion del timestamp 

	//
	conn, err := grpc.Dial(dNN + *//*puerto*//*, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("No se pudo conectar a Namenode: %v", err)
	}
	defer conn.Close()

	cliente := book.NewBookServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	d, err2 := cliente.EscribirLogDes(&book.PropuestaLibro) // revisar si se puede
	if err2 != nil {
		log.Printf("problema con request: %v", err2)
	}	

	
}
*/
	/*
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
	
}*/


//enviar chunk a Cliente Downloader
func (s *server) enviarChunkDN(ctx context.Context, chunk *book.Chunk) (*book.Chunk, error){
	
	file, err := os.Open("./Chunks/" + chunk.NombreArchivo)
	failOnError(err, "No se pudo abrir archivo de chunk a enviar a Cliente Downloader")
	defer file.Close()

	buf := make([]byte, 250000)
	n, err := file.Read(buf)
	failOnError(err, "No se pudo leer del archivo de chunk en buffer")

	file.Close()
	return &book.Chunk{Contenido: buf[:n]}, nil
}


//recibir chunks de DataNode al ser distribuidos
func (s *server) recibirChunksDN(ctx context.Context, c *book.Chunk) (*book.ACK, error){

	f, err := os.OpenFile(c.NombreArchivo, os.O_WRONLY|os.O_APPEND, 0644)
	failOnError(err, "Error creando archivo")
	defer f.Close()


	//obtener tamaño del chunk
	/*chunkStat := c.Contenido.Stat()
	chunkSize := chunkStat.Size()*/
	chunkSize := len(c.Contenido)

	//arreglo de bytes para traspasar chunk
	buf := make([]byte, chunkSize)

	//escribir contenido en buffer
	reader := bytes.NewReader(c.Contenido)
	_, err = reader.Read(buf)
	failOnError(err, "Error escribiendo chunk recibido en buffer")

	//escribir contenido de buffer en el archivo
	_, err = f.Write(buf)
	failOnError(err, "Error escribiendo chunk en archivo")


	f.Sync() //flush to disk
	f.Close()
	buf = nil
	fmt.Println("Escritura exitosa")

	return &book.ACK{Ok: "ok"}, nil
}


//enviar 1 chunk a otro DataNode
func enviarChunk(archivoChunk string, ip string){
	
	var port string

	if (ip == dB){
		port = ":50500"
	} else if (ip == dC){
		port = ":50502"
	}

	//conexión
	connDN, err := grpc.Dial(ip + port, grpc.WithInsecure(), grpc.WithBlock())
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


	_, err = clientDN.RecibirChunksDN(ctx, &book.Chunk{Contenido: buf[:n], NombreArchivo: archivoChunk})
	failOnError(err, "Error distribuyendo chunks a DataNode")

	buf = nil
	connDN.Close()
	file.Close()
}


//enviar chunks a otros nodos
func distribuirChunks(prop string){

	var line []string
	var c []string

	line = strings.Split(prop,"\n")[1:]

	for _, info := range line {

		c = strings.Fields(info)
		
		//parámetros: nombre de archivo de chunk y dirección de nodo
		enviarChunk(c[0], c[1])
	}
}


//genera una nueva propuesta considerando solo los nodos activos
//dn: nodo caído
func nuevaPropuesta2(dn string, cantChunks int, nombreLibro string) (string, bool, bool){

	b := false
	c := false

	//inicio de propuesta
	Prop := nombreLibro + " " + strconv.Itoa(cantChunks) + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, cantChunks)

	for i := 0; i < cantChunks; i++{
		intProp[i] = rand.Intn(2)
	}
	
	var dAct string
	for i := 0; i < cantChunks; i++{ 
		
		if (dn == dB) {
			switch intProp[i]{
			case 0:
				dAct = dActual
			case 1:
				dAct = dC
			}
			c = true
		}

		if (dn == dC) {
			switch intProp[i]{
			case 0:
				dAct = dActual
			case 1:
				dAct = dB
			}
			b = true
		}

		Prop += nombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"
		
	}

	return Prop, b, c
}


//verifica si hay un DataNode caido
//name: nombre de datanode ("datanode X")
func checkDatanode(dn string, port string, name string) bool{

	deadline := 5 //segundos que espera para ver si hay conexión

	connDN, err := grpc.Dial(dn + port, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(deadline)*time.Second))
    if (err != nil) {
		log.Printf("No se pudo conectar con %v: %v", name, err)
		return false
    }
	defer connDN.Close()

	//c := book.NewBookServiceClient(connDN) //necesario?? es solo chekear pero no aun conectar 4real

	connDN.Close()
	return true
}


//retorna true si acepta propuesta y false y el nodo caído si rechaza
func analizarPropuesta(prop *book.PropuestaLibro) (bool, string){
	
	log.Println("Analizando la propuesta...")

	//prop sabe cuales datanodes se pretenden usar en la propuesta
	//por ejemplo, si prop.DatanodeA==true, entonces se usaría en la propuesta
	//pero si esta caído se rechaza la propuesta

	//revisa si hay un DataNode de la propuesta caído
	if (prop.DatanodeB && !checkDatanode(dB, ":50500", "datanode B")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, dB
	} 
	if (prop.DatanodeC && !checkDatanode(dC, "50502", "datanode C")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, dC
	}

	return true, ""
}


//aceptar o rechazar propuesta de distribución de chunks (descentralizado)
func (s *server) RecibirPropuesta(ctx context.Context, prop *book.PropuestaLibro) (*book.RespuestaP, error){
	
	resp, DN := analizarPropuesta(prop)
	return &book.RespuestaP{Respuesta: resp, DNcaido: DN}, nil
}


//enviar una propuesta a otro DataNode (descentralizado)
//retorna la respuesta de este y un DataNode caído, en caso de hallar uno
func enviarPropuestaDN(dir string, port string, prop *book.PropuestaLibro) (bool, string) {

	//conexión
	connDN, err := grpc.Dial(dir + port, grpc.WithInsecure(), grpc.WithBlock())
	defer connDN.Close()
	//en caso de error de conexión, en vez de terminar el programa, se maneja de manera de ignorarlo 
	//y que el otro nodo sea el encargado de rechazar la propuesta
	if (err != nil){ 
		return false, "nil"
	}

	clientDN := book.NewBookServiceClient(connDN)
	log.Println("Conexión realizada\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	
	resp, err := clientDN.RecibirPropuesta(ctx, prop) //check
	
	connDN.Close()
	return resp.Respuesta, resp.DNcaido
}


//enviar propuesta a NameNode (descentralizado)
func enviarPropNNDes(prop *book.PropuestaLibro) {

	//conexión
	connNN, err := grpc.Dial(dNN + ":50506", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()
	
	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	//propuesta final (book.PropuestaLibro)
	ACK, err := clientNN.EscribirLogDes(ctx, prop)
	log.Println(ACK.Ok)
	failOnError(err, "")
	
	connNN.Close()
}


//enviar propuesta a NameNode (centralizado)
func EnviarPropNNCen(prop *book.PropuestaLibro) (string, bool, bool, bool){
	
	//conexión
	connNN, err := grpc.Dial(dNN + ":50506", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()
	
	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	//propuesta final (book.PropuestaLibro)
	propFinal, err := clientNN.RecibirPropDatanode(ctx, prop)
	failOnError(err, "Error enviando ")
	
	connNN.Close()
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
	Prop := nombreLibro + " " + strconv.Itoa(cantChunks) + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, cantChunks)

	for i := 0; i < cantChunks; i++{
		intProp[i] = rand.Intn(3)
	}
	
	var dAct string //dirección correspondiente 
	for i := 0; i < cantChunks; i++{ 
		
		switch intProp[i]{
			case 0:
				dAct = dActual
				a = true
			case 1:
				dAct = dB
				b = true
			case 2:
				dAct = dC
				c = true
		}
		
		Prop += nombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"
		
	}

	return Prop, a, b, c
}


//Recibir chunks de Cliente Uploader
func (s *server) RecibirChunks(stream book.BookService_RecibirChunksServer) error {
	
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
			
			err = errors.Wrapf(err, "failed unexpectadely while reading chunks from stream")
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
		n, err := file.Write(chunk.Contenido)
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

//Generar propuesta, escribir en log según algoritmo y distribuir entre DataNodes
	
	//prop: propuesta; b, c: estados de DataNodes B y C
	prop, a, b, c := generarPropuesta(numChunk, nombreL)

//CENTRALIZADO
	if (tipo_al == "c"){
		
		//recibe prop, que puede ser la misma propuesta o una nueva del NameNode
		//CHECK < < < < < < son necesarios a, b, c?
		prop, a, b, c = EnviarPropNNCen(&book.PropuestaLibro{NombreLibro: nombreL, Propuesta: prop, CantChunks: int32(numChunk), DatanodeA: a, DatanodeB: b, DatanodeC: c})

//DESCENTRALIZADO
	} else {
		
		//enviar a demás nodos
		respB, DNcaidoB := enviarPropuestaDN(dB, ":50500", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
		respC, DNcaidoC := enviarPropuestaDN(dC, ":50502", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})

		//generar nuevas propuestas hasta que se apruebe una
		for {

			//DNcaidoX es "" cuando respX == true (se aprueba prop y no hay nodo caído)
			//DNcaidoX es "nil" cuando hubo un problema de conexión (no se aprueba ni rechaza prop)
			if (!respB && (DNcaidoB != "nil")) {
				prop, b, c = nuevaPropuesta2(DNcaidoB, numChunk, nombreL)
				respB, DNcaidoB = enviarPropuestaDN(dB, ":50500", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})

			} else if !respC {
				prop, b, c = nuevaPropuesta2(DNcaidoC, numChunk, nombreL)
				respC, DNcaidoC = enviarPropuestaDN(dC, ":50502", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
			
			} else {
				break
			}
		}
		
		//enviar propuesta ya aprobada a NameNode para que lo escriba en Log
		enviarPropNNDes(&book.PropuestaLibro{Propuesta: prop})
	}
	
	//enviar chunks a Datanodes de la propuesta
	distribuirChunks(prop)

	_ = stream.SendAndClose(&book.ACK{Ok: "ok"})
	return nil
}


// CONEXIONES
//cliente Downloader
func serveCD(){
	
	listenCD, err := net.Listen("tcp", dActual + ":50513")
	failOnError(err, "Error de conexión con cliente downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCD))
}


//cliente Uploader
func serveCU(){
	
	listenCU, err := net.Listen("tcp", dActual + "50517")
	failOnError(err, "Error de conexión con cliente downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCU))
}


//NameNode
func serveNN(){
	
	listenNN, err := net.Listen("tcp", dActual + "50509")
	failOnError(err, "Error de conexión con cliente downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenNN))
}


//DataNode B
func serveDNB(){

	listenDNB, err := net.Listen("tcp", dActual + ":50501")
	failOnError(err, "Error de conexión con DataNode B")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNB))
}


//DataNode C
func serveDNC(){
	
	listenDNC, err := net.Listen("tcp", dActual + ":50503")
	failOnError(err, "Error de conexión con DataNode C")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNC))
}


func main() {

//servers
	go serveCD()
	go serveCU()
	go serveNN()
	go serveDNB()
	go serveDNC()
	
}