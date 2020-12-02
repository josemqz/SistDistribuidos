package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"net"
	"time"
	"context"
	"strconv"
	"math/rand"
	"strings"
	"bytes"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

/*
var dNN = "10.6.40.157" //ip namenode
var dActual = "10.6.40.158" //ip nodo actual
//var dA = "10.6.40.158" //ip maquina virtual datanode A
var dB = "10.6.40.159" //ip maquina virtual datanode B
var dC = "10.6.40.160" //ip maquina virtual datanode C
*/

// test local
var local = true
var dNN = "localhost"
var dA = "localhost"
var dB = "localhost"
var dC = "localhost"
var dActual = "localhost"


var contadorMensajesA int  //mensajes de datanode A
var mutexCA = &sync.Mutex{}

//Ricart & Agrawala
var estado string		  	 //estado de proceso
var id_proceso = 0		  	 //id de proceso (para insertar en cola)
var colaRA []int		 	 //cola de nodos en espera
var tiempo_request string

var mutexLocal = &sync.Mutex{}
var mutexExt = &sync.Mutex{}


type server struct {
	book.UnimplementedBookServiceServer
}


//manejo de errores
func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}


//funcion para medir el tiempo (informe)
func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("tiempo %s : %s", name, elapsed)
}


// RICART Y AGRAWALA

//verificar si un proceso sigue en la cola
func checkColas(queue []int, id int) bool{
	
	for _, elemento := range queue{
		if (elemento == id){
			return true
		}
	}
	return false
}


func (s *server) RequestRA(ctx context.Context, p *book.ExMutua) (*book.ACK, error){

	var id int

	//String a Time
	tiempo_i, _ := time.Parse("2006-01-02 15:04:05", p.Tiempo)
	tiempo_req, _ := time.Parse("2006-01-02 15:04:05", tiempo_request)
	
	mutexExt.Lock()
	if (estado == "HELD") || (estado == "WANTED" && (tiempo_req.Before(tiempo_i))){
		
		//crear id para proceso
		id_proceso += 1
		id = id_proceso + 1

		//meter id a cola
		colaRA = append(colaRA, id)

		mutexExt.Unlock()
		
	} else{

		mutexExt.Unlock()
		
		//contestar al nodo
		return &book.ACK{Ok: "ok"}, nil	
	}
	
	
	mutexExt.Lock()
	
	mutexCA.Lock()
	contadorMensajesA +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCA.Unlock()
	
	for checkColas(colaRA, id){
		time.Sleep(500 * time.Millisecond)
	}
	mutexExt.Unlock()
	
	return &book.ACK{Ok: "ok"}, nil
	
}


//Nodo que quiere escribir en log llama a funcion de Ricart y Agrawala
func RicAwla(prop *book.PropuestaLibro){

	mutexLocal.Lock()

	var resp1 = false
	var resp2 = false

	//para entrar a zona crítica
	estado = "WANTED"

	//tiempo de request
	t := time.Now()
	tiempo_request = t.Format("2006-01-02 15:04:05")

	//request a nodos considerados en propuesta
	if prop.DatanodeB{

		connDNB, err := grpc.Dial(dB + ":50500", grpc.WithInsecure())
		if err != nil{
			log.Printf("No se pudo conectar a DataNode: %v", err)
		}
		defer connDNB.Close()
		
		clientDNB := book.NewBookServiceClient(connDNB)

		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
		
		ACK, err := clientDNB.RequestRA(ctx, &book.ExMutua{Tiempo: tiempo_request})
		if err != nil {
			log.Printf("problema con request: %v", err)
		}

		if ACK.Ok == "ok"{
			resp1 = true
		}

		connDNB.Close()
	}
	
	if prop.DatanodeC{

		connDNC, err := grpc.Dial(dC + ":50502", grpc.WithInsecure())
		if err != nil{
			log.Printf("No se pudo conectar a DataNode: %v", err)
		}
		defer connDNC.Close()
		
		clientDNC := book.NewBookServiceClient(connDNC)

		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		ACK, err := clientDNC.RequestRA(ctx, &book.ExMutua{Tiempo: tiempo_request})
		if err != nil {
			log.Printf("problema con request: %v", err)
		}
		
		if ACK.Ok == "ok"{
			resp2 = true
		}
	
		connDNC.Close()
	}

	//enviar propuesta a Namenode al recibir ambas respuestas
	if (prop.DatanodeB && resp1) || (prop.DatanodeC && resp2){

		estado = "HELD"
		
		enviarPropNNDes(prop)
	}

	//terminar
	estado = "RELEASED"
	
	//sacar nodos de cola
	if (len(colaRA) >= 1){
		colaRA = colaRA[1:]
	} else if (len(colaRA) == 1){
		colaRA = colaRA[:0]
	}

	mutexLocal.Unlock()
}


//enviar chunk a Cliente Downloader
func (s *server) EnviarChunkDN(ctx context.Context, chunk *book.Chunk) (*book.Chunk, error){
	
	log.Println("Abriendo chunk para enviar a Cliente Downloader")

	file, err := os.Open("./DNA/Chunks/" + chunk.NombreArchivo)
	if (err != nil){
		return &book.Chunk{}, err
	}
	defer file.Close()

	buf := make([]byte, 250000)
	n, err := file.Read(buf)

	if (err != nil){
		file.Close()
		return &book.Chunk{}, err
	}

	file.Close()
	log.Println("Enviando chunk a Cliente Downloader")
	return &book.Chunk{Contenido: buf[:n], NumChunk: chunk.NumChunk}, nil
}


//recibir chunks de DataNode al ser distribuidos
func (s *server) RecibirChunksDN(ctx context.Context, c *book.Chunk) (*book.ACK, error){

	f, err := os.OpenFile("./DNA/Chunks/" + c.NombreArchivo, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	failOnError(err, "Error creando archivo")
	defer f.Close()


	//obtener tamaño del chunk
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
	fmt.Println("Escritura en chunk exitosa\n")

	mutexCA.Lock()
	contadorMensajesA +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCA.Unlock()

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
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100) * time.Second)
	defer cancel()


	//abrir archivo a enviar
	file, err := os.Open("./DNA/Chunks/" + archivoChunk)
	failOnError(err,"No se pudo abrir archivo de chunk a enviar a DataNode")
	defer file.Close()

	//buffer para enviar bytes leídos del archivo
	buf := make([]byte, 250000)
	n, err := file.Read(buf)
	failOnError(err,"No se pudo leer del archivo de chunk en buffer")


	log.Println("Distribuyendo chunk a otro DataNode\n")
	_, err = clientDN.RecibirChunksDN(ctx, &book.Chunk{Contenido: buf[:n], NombreArchivo: archivoChunk})
	failOnError(err, "Error distribuyendo chunks a DataNode")

	
	//eliminar archivo
	file.Close()
	err = os.Remove("./DNA/Chunks/" + archivoChunk)
	failOnError(err, "No se pudo eliminar archivo de chunk")

	buf = nil
	connDN.Close()

	mutexCA.Lock()
	contadorMensajesA +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCA.Unlock()
}


//enviar chunks a otros nodos
func distribuirChunks(prop string){

	var line []string
	var c []string

	line = strings.Split(prop[:(len(prop)-1)],"\n")[1:] //what if está considerando el último elemento como un ""

	fmt.Println(line)

	for _, info := range line {

		c = strings.Fields(info)
		
		if (c[1] != dActual) || local{  // DEBUG < < < < < < < < < < < < < < < <
			//parámetros: nombre de archivo de chunk y dirección de nodo
			enviarChunk(c[0], c[1])
		}
	}
}


//genera una nueva propuesta considerando solo los nodos activos
//dn: nodo caído
func nuevaPropuesta2(dn string, cantChunks int, nombreLibro string) (string, bool, bool, bool){

	a := false
	b := false
	c := false

	//inicio de propuesta
	Prop := nombreLibro + " " + strconv.Itoa(cantChunks) + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, cantChunks)

	for i := 0; i < cantChunks; i++{
		intProp[i] = rand.Intn(2)
	}
	
	//asegurar que los dos nodos sean escogidos para los primeros chunks
	if (cantChunks >= 2){
		rd := rand.Perm(2)
		
		intProp[0] = rd[0]
		intProp[1] = rd[1]
	}
	
	var dAct string
	for i := 0; i < cantChunks; i++{ 
		
		if (dn == dB) {
			switch intProp[i]{
			case 0:
				dAct = dActual
				a = true
			case 1:
				dAct = dC
				c = true
			}
		}

		if (dn == dC) {
			switch intProp[i]{
			case 0:
				dAct = dActual
				a = true
			case 1:
				dAct = dB
				b = true
			}
		}

		Prop += nombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"
	}

	return Prop, a, b, c
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

	mutexCA.Lock()
	contadorMensajesA +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCA.Unlock()

	return true
}


//retorna true si acepta propuesta y false y el nodo caído si rechaza
func analizarPropuesta(prop *book.PropuestaLibro) (bool, string){
	
	log.Println("Analizando la propuesta...")

	//prop sabe cuales DataNodes son considerados en la propuesta
	//por ejemplo, si prop.DatanodeA == true, entonces el nodo A se usaría en 
	//la propuesta pero si esta caído se rechaza la propuesta

	//revisa si hay un DataNode de la propuesta caído
	if (prop.DatanodeB && !checkDatanode(dB, ":50500", "datanode B")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, dB
	} 
	if (prop.DatanodeC && !checkDatanode(dC, "50502", "datanode C")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, dC
	}

	fmt.Println("Se aprueba la propuesta\n")
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
	
	mutexCA.Lock()
	contadorMensajesA +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCA.Unlock()
	
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
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100) * time.Second)
	defer cancel()

	mutexCA.Lock()
	contadorMensajesA +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCA.Unlock()
		

	//propuesta final (book.PropuestaLibro)
	ACK, err := clientNN.EscribirLogDes(ctx, prop)
	failOnError(err, "")
	
	log.Println(ACK.Ok)
	
	connNN.Close()

	log.Println("Cantidad de mensajes:", contadorMensajesA)
}


//enviar propuesta a NameNode (centralizado)
func EnviarPropNNCen(prop *book.PropuestaLibro) (string){
	
	//conexión
	connNN, err := grpc.Dial(dNN + ":50506", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()
	
	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	mutexCA.Lock()
	contadorMensajesA +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCA.Unlock()
	
	//propuesta final (book.PropuestaLibro)
	propFinal, err := clientNN.RecibirPropDatanode(ctx, prop)
	failOnError(err, "Error enviando ")
	
	connNN.Close()

	log.Println("Cantidad de mensajes:", contadorMensajesA)
	return propFinal.Propuesta
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

	//asegurar que los tres (o dos) nodos sean escogidos para los primeros chunks
	if (cantChunks >= 3){
		rd := rand.Perm(3)
		
		intProp[0] = rd[0]
		intProp[1] = rd[1]
		intProp[2] = rd[2]

	} else if (cantChunks == 2){
		rd := rand.Perm(2)
		
		intProp[0] = rd[0]
		intProp[1] = rd[1]
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

	defer timeTrack(time.Now(), "Total") //entrega el tiempo de ejecucion de la funcion

	contadorMensajesA = 0
	
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
		neoArchLibro := "./DNA/Chunks/" + chunk.NombreLibro + "_" + strNumChunk
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
	
	//prop: propuesta; a, b, c: estados de DataNodes A, B y C
	prop, a, b, c := generarPropuesta(numChunk, nombreL)

//CENTRALIZADO
	if (tipo_al == "c"){
		
		//recibe prop, que puede ser la misma propuesta o una nueva del NameNode
		prop = EnviarPropNNCen(&book.PropuestaLibro{NombreLibro: nombreL, Propuesta: prop, CantChunks: int32(numChunk), DatanodeA: a, DatanodeB: b, DatanodeC: c})

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
				prop, a, b, c = nuevaPropuesta2(DNcaidoB, numChunk, nombreL)
				respB, DNcaidoB = enviarPropuestaDN(dB, ":50500", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})

			} else if !respC {
				prop, a, b, c = nuevaPropuesta2(DNcaidoC, numChunk, nombreL)
				respC, DNcaidoC = enviarPropuestaDN(dC, ":50502", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
			
			} else {
				break
			}
		}
		
		//enviar propuesta ya aprobada a NameNode para que lo escriba en Log
		RicAwla(&book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
	}
	
	//enviar chunks a Datanodes de la propuesta
	distribuirChunks(prop)

	_ = stream.SendAndClose(&book.ACK{Ok: "ok"})
	return nil
}


// CONEXIONES
//cliente Downloader
func serveCD(){
	
	log.Println("Servidor Cliente Downloader esperando...")
	
	listenCD, err := net.Listen("tcp", dActual + ":50513")
	failOnError(err, "Error de conexión con Cliente Downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCD))
}


//cliente Uploader
func serveCU(){
	
	log.Println("Servidor Cliente Uploader esperando...")

	listenCU, err := net.Listen("tcp", dActual + ":50517")
	failOnError(err, "Error de conexión con Cliente Uploader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCU))
}


//NameNode
func serveNN(){
	
	log.Println("Servidor NameNode esperando...")

	listenNN, err := net.Listen("tcp", dActual + ":50509")
	failOnError(err, "Error de conexión con Cliente Downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenNN))
}


//DataNode B
func serveDNB(){

	log.Println("Servidor DataNode B esperando...")

	listenDNB, err := net.Listen("tcp", dActual + ":50501")
	failOnError(err, "Error de conexión con DataNode B")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNB))
}


//DataNode C
func serveDNC(){
	
	log.Println("Servidor DataNode C esperando...")

	listenDNC, err := net.Listen("tcp", dActual + ":50503")
	failOnError(err, "Error de conexión con DataNode C")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNC))
}


func main() {

	rand.Seed(time.Now().Unix())

//servers
	go serveCD()
	go serveCU()
	go serveNN()
	go serveDNB()
	serveDNC()
	
}