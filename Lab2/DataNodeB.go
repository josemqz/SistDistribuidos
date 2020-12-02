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


var dNN = "10.6.40.157" //ip namenode
var dActual = "10.6.40.159" //ip nodo actual
var dA = "10.6.40.158" //ip maquina virtual datanode A
var dC = "10.6.40.160" //ip maquina virtual datanode C

/* test local
var dNN = "localhost"
var dA = "localhost"
var dB = "localhost"
var dC = "localhost"
var dActual = "localhost"
*/

var contadorMensajesB int  //mensajes de datanode B
var mutexCB = &sync.Mutex{}

//Ricart & Agrawala
var estado string		  //estado de proceso
var id_proceso = 0		  //id de proceso (para insertar en cola)
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


//funcion para medir el tiempo
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
		id := id_proceso + 1

		//meter id a cola
		colaRA = append(colaRA, id)

		mutexExt.Unlock()
		
	} else{

		mutexExt.Unlock()
		
		//contestar al nodo
		return &book.ACK{Ok: "ok"}, nil	
	}
	
	
	mutexExt.Lock()
	
	mutexCB.Lock()
	contadorMensajesB +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCB.Unlock()
	
	for checkColas(colaRA, id){
		time.Sleep(500 * time.Millisecond)
	}
	mutexExt.Unlock()
	
	return &book.ACK{Ok: "ok"}, nil
	
}


//Nodo que quiere escribir en log llama a funcion de Ricart y Agrawala
func RicAwla(prop *book.PropuestaLibro){

	mutexLocal.Lock()

	resp1 := false
	resp2 := false

	//para entrar a zona crítica
	estado = "WANTED"

	//tiempo de request
	t := time.Now()
	tiempo_request = t.Format("2006-01-02 15:04:05")

	//request a nodos considerados en propuesta
	if prop.DatanodeA{

		connDNA, err := grpc.Dial(dA + ":50501", grpc.WithInsecure())
		if err != nil{
			log.Printf("No se pudo conectar a DataNode: %v", err)
		}
		defer connDNA.Close()
		
		clientDNA := book.NewBookServiceClient(connDNA)

		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
		
		ACK, err := clientDNA.RequestRA(ctx, &book.ExMutua{Tiempo: tiempo_request})
		if err != nil {
			log.Printf("problema con request: %v", err)
		}
	
		if ACK.Ok == "ok"{
			resp1 = true
		}

		connDNA.Close()
	}
	
	if prop.DatanodeC{

		connDNC, err := grpc.Dial(dC + ":50504", grpc.WithInsecure())
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
	if (prop.DatanodeA && resp1) || (prop.DatanodeC && resp2){

		estado = "HELD"
		
		enviarPropNNDes(prop)
	}

	//terminar
	estado = "RELEASED"
	//sacar nodos de cola
	colaRA = colaRA[1:]

	mutexLocal.Unlock()
}


//enviar chunk a Cliente Downloader
func (s *server) enviarChunkDN(ctx context.Context, chunk *book.Chunk) (*book.Chunk, error){
	
	file, err := os.Open("./DNB/Chunks/" + chunk.NombreArchivo)
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

	mutexCB.Lock()
	contadorMensajesB +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCB.Unlock()

	return &book.ACK{Ok: "ok"}, nil
}


//enviar 1 chunk a otro DataNode
func enviarChunk(archivoChunk string, ip string){
	
	var port string

	if (ip == dA){
		port = ":50501"
	} else if (ip == dC){
		port = ":50504"
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
	file, err := os.Open("./DNB/Chunks/" + archivoChunk)
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

	mutexCB.Lock()
	contadorMensajesB +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCB.Unlock()
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
		
		if (dn == dA) {
			switch intProp[i]{
			case 0:
				dAct = dActual
				b = true
			case 1:
				dAct = dC
				c = true
			}
		}

		if (dn == dC) {
			switch intProp[i]{
			case 0:
				dAct = dActual
				b = true
			case 1:
				dAct = dA
				a = true
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

	mutexCB.Lock()
	contadorMensajesB +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCB.Unlock()

	return true
}


//retorna true si acepta propuesta y false y el nodo caído si rechaza
func analizarPropuesta(prop *book.PropuestaLibro) (bool, string){
	
	log.Println("Analizando la propuesta...")

	//prop sabe cuales DataNodes son considerados en la propuesta
	//por ejemplo, si prop.DatanodeA == true, entonces el nodo A se usaría en 
	//la propuesta pero si esta caído se rechaza la propuesta

	//revisa si hay un DataNode de la propuesta caído
	if (prop.DatanodeA && !checkDatanode(dA, ":50501", "datanode A")){
		fmt.Println("Se rechaza la propuesta\n")
		return false, dA
	} 
	if (prop.DatanodeC && !checkDatanode(dC, "50504", "datanode C")){
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
	
	mutexCB.Lock()
	contadorMensajesB +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCB.Unlock()
	
	resp, err := clientDN.RecibirPropuesta(ctx, prop) //check
	
	connDN.Close()
	return resp.Respuesta, resp.DNcaido
}


//enviar propuesta a NameNode (descentralizado)
func enviarPropNNDes(prop *book.PropuestaLibro) {

	//conexión
	connNN, err := grpc.Dial(dNN + ":50507", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()
	
	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	mutexCB.Lock()
	contadorMensajesB +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCB.Unlock()
	
	//propuesta final (book.PropuestaLibro)
	ACK, err := clientNN.EscribirLogDes(ctx, prop)
	log.Println(ACK.Ok)
	failOnError(err, "")
	
	connNN.Close()

	log.Println("Cantidad de mensajes: s%", contadorMensajesB)
}


//enviar propuesta a NameNode (centralizado)
func EnviarPropNNCen(prop *book.PropuestaLibro) (string){
	
	//conexión
	connNN, err := grpc.Dial(dNN + ":50507", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a NameNode")
	defer connNN.Close()
	
	clientNN := book.NewBookServiceClient(connNN)
	log.Println("Conexión realizada\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	mutexCB.Lock()
	contadorMensajesB +=1 // SUMA UN MENSAJE PARA METRICAS DEL INFORME
	mutexCB.Unlock()
	
	//propuesta final (book.PropuestaLibro)
	propFinal, err := clientNN.RecibirPropDatanode(ctx, prop)
	failOnError(err, "Error enviando ")
	
	connNN.Close()

	log.Println("Cantidad de mensajes: s%", contadorMensajesB)

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
				dAct = dA
				a = true
			case 1:
				dAct = dActual
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

	contadorMensajesB = 0
	
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
		neoArchLibro := "./DNB/Chunks/" + chunk.NombreLibro + "_" + strNumChunk
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
		respA, DNcaidoA := enviarPropuestaDN(dA, ":50501", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
		respC, DNcaidoC := enviarPropuestaDN(dC, ":50504", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})

		//generar nuevas propuestas hasta que se apruebe una
		for {

			//DNcaidoX es "" cuando respX == true (se aprueba prop y no hay nodo caído)
			//DNcaidoX es "nil" cuando hubo un problema de conexión (no se aprueba ni rechaza prop)
			if (!respA && (DNcaidoA != "nil")) {
				prop, a, b, c = nuevaPropuesta2(DNcaidoA, numChunk, nombreL)
				respA, DNcaidoA = enviarPropuestaDN(dA, ":50501", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})

			} else if !respC {
				prop, a, b, c = nuevaPropuesta2(DNcaidoC, numChunk, nombreL)
				respC, DNcaidoC = enviarPropuestaDN(dC, ":50504", &book.PropuestaLibro{Propuesta: prop, DatanodeA: a, DatanodeB: b, DatanodeC: c})
			
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
	
	listenCD, err := net.Listen("tcp", dActual + ":50514")
	failOnError(err, "Error de conexión con cliente downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCD))
}


//cliente Uploader
func serveCU(){
	
	log.Println("Servidor Cliente Uploader esperando...")
	
	listenCU, err := net.Listen("tcp", dActual + ":50518")
	failOnError(err, "Error de conexión con Cliente Downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCU))
}


//NameNode
func serveNN(){
	
	log.Println("Servidor NameNode esperando...")

	listenNN, err := net.Listen("tcp", dActual + ":50510")
	failOnError(err, "Error de conexión con Cliente Uploader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenNN))
}


//DataNode A
func serveDNA(){

	log.Println("Servidor DataNode A esperando...")

	listenDNA, err := net.Listen("tcp", dActual + ":50500")
	failOnError(err, "Error de conexión con DataNode A")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNA))
}


//DataNode C
func serveDNC(){
	
	log.Println("Servidor DataNode C esperando...")

	listenDNC, err := net.Listen("tcp", dActual + ":50505")
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
	go serveDNA()
	serveDNC()
	
}