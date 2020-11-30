package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"net"
	"time"
	"context"
	"math/rand"
	"sync"
	"strings"
	"errors"
	"strconv"

	"google.golang.org/grpc"
	"github.com/josemqz/SistDistribuidos/Lab2/book"
)

//DEFINIR NOMBRES DE DATANODES CON IP CORRESPONDIENTES
var dActual = "10.6.40.157"  //ip NameNode
var dA = "10.6.40.158" //ip maquina virtual datanode A
var dB = "10.6.40.159" //ip maquina virtual datanode B
var dC = "10.6.40.160" //ip maquina virtual datanode C

/* test local
var dA = "localhost"
var dB = "localhost"
var dC = "localhost"
var dActual = "localhost"
*/

type server struct {
	book.UnimplementedBookServiceServer
}

var mutex = &sync.Mutex{}


func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}


//para recibir entradas para el log
//distribuido
func (s *server) escribirLogDes(prop *book.PropuestaLibro) (*book.ACK, error) {
	
	f, err := os.OpenFile("logdata.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if (err != nil){
		return &book.ACK{Ok: "error"}, errors.New("Error abriendo Log en NameNode")
	}
    defer f.Close()

    _, err2 := f.WriteString(prop.Propuesta)
	if (err2 != nil){
		return &book.ACK{Ok: "error"}, errors.New("Error escribiendo en Log en NameNode")
	}
	
	return &book.ACK{Ok: "listo"}, nil
}


//centralizado
func escribirLogCen(prop string, nombreL string, cant int32){

	mutex.Lock()
	f, err := os.OpenFile("logdata.txt", os.O_WRONLY|os.O_APPEND, 0644)
	failOnError(err, "Error abriendo log")
    defer f.Close()

    _, err2 := f.WriteString(prop)
	failOnError(err2, "Error escribiendo en log")

	mutex.Unlock()

	fmt.Println("Escritura en log exitosa")
}


//verifica si hay un DataNode caido
func checkDatanode(dn string, port string, name string) bool{

	deadline := 5 //segundos que espera para ver si hay conexión

	connDN, err := grpc.Dial(dn + port, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(deadline)*time.Second))
    if (err != nil) {
		log.Printf("Se detectó datanode caido %v: %v", name, err)
		return false
    }
	defer connDN.Close()
	
	//c := book.NewBookServiceClient(connDN) //necesario?? es solo chekear pero no aun conectar 4real

	connDN.Close()
	return true
}


//retorna true si acepta propuesta y false si rechaza
func analizarPropuesta(prop *book.PropuestaLibro) bool{
	
	log.Println("Analizando la propuesta...")

	//prop sabe cuales datanodes se pretenden usar en la propuesta
	//por ejemplo, si prop.DatanodeA==true, entonces se usaría en la propuesta
	//pero si esta caído se rechaza la propuesta

	//revisa si hay un datanode de la propuesta caído
	if (prop.DatanodeA && !checkDatanode(dA, ":50509", "DataNode A")){
		fmt.Println("Se rechaza la propuesta\n")
		return false
	}
	if (prop.DatanodeB && !checkDatanode(dB, ":50510", "DataNode B")){
		fmt.Println("Se rechaza la propuesta\n")
		return false
	} 
	if (prop.DatanodeC && !checkDatanode(dC, ":50511", "DataNode C")){
		fmt.Println("Se rechaza la propuesta\n")
		return false
	}

	return true
}


func buscarNodoCaido() string{

	//revisa si hay un datanode caído
	if (!checkDatanode(dA, ":50509", "DataNode A")){
		return dA
	}
	if (!checkDatanode(dB, ":50510", "DataNode B")){
		return dB
	} 
	if (!checkDatanode(dC, ":50511", "DataNode C")){
		return dC

	} else{
		return "no hay nodos caidos"
	}

}


//genera una nueva propuesta considerando solo los nodos activos
func nuevaPropuesta2(dn string, prop *book.PropuestaLibro) (string, bool, bool, bool){

	//estado de nodos
	a := true
	b := true
	c := true

	//inicio de propuesta
	n := prop.CantChunks
	Prop := prop.NombreLibro + " " + strconv.Itoa(int(n)) + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, n)

	for i := 0; int32(i) < n; i++{
		intProp[i] = rand.Intn(2)
	}
	
	var dAct string
	for i := 0; int32(i) < n; i++{ 
		
		if (dn == dA) {
			switch intProp[i]{
			case 0:
				dAct = dB
			case 1:
				dAct = dC
			}
			a = false
		}
		if (dn == dB) {
			switch intProp[i]{
			case 0:
				dAct = dA
			case 1:
				dAct = dC
			}
			b = false
		}
		if (dn == dC) {
			switch intProp[i]{
			case 0:
				dAct = dA
			case 1:
				dAct = dB
			}
			c = false
		}

		Prop += prop.NombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"

	}

	return Prop, a, b, c
}


//recibir la propuesta de un DataNode (centralizado)
func (s *server) recibirPropDatanode(ctx context.Context, prop *book.PropuestaLibro) (*book.PropuestaLibro, error){

//Cuando un DataNode envia una propuesta, se analiza la Propuesta
//Si se rechaza se genera una nueva y se analiza hasta que analizarPropuesta sea true
//Luego se escribe en el log la propuesta

	Prop := prop.Propuesta
	a := prop.DatanodeA
	b := prop.DatanodeB
	c := prop.DatanodeC

	for {
		if analizarPropuesta(&book.PropuestaLibro{Propuesta: Prop, DatanodeA: a, DatanodeB: b, DatanodeC: c}) {
			break
		} else{
			dn := buscarNodoCaido()
			Prop, a, b, c = nuevaPropuesta2(dn,prop)
		}
	}

	escribirLogCen(Prop, prop.NombreLibro, prop.CantChunks)

	return &book.PropuestaLibro{Propuesta: Prop, DatanodeA: a, DatanodeB: b, DatanodeC: c}, nil
}


//Responde al Cliente Downloader con las ubicaciones de los chunks del libro solicitado
func localizacionChunks(nombreL string) (string, error){

	f, err := os.Open("logdata.txt")
	failOnError(err, "Error en abrir log")
	defer f.Close()

	// hace Splits por cada linea por defecto.
	scanner := bufio.NewScanner(f)

	var listachunks string
	var info []string
	var t string
	var init int
	var n int
	var mark bool
	mark = false
	
	for scanner.Scan() {
		t = scanner.Text()
		if mark{
			info = strings.Fields(t)
			listachunks += info[1] + " "
			init++
			if (init == n) {
				return listachunks, nil
			}
		}
    	if (strings.Contains(t, nombreL)) {
			words := strings.Fields(t) //es como split por blankspace
			n, _ = strconv.Atoi(words[1])
			init = 0
			mark = true
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("No se pudo localizar chunks correctamente: %v", err)
	}

	return "Error", errors.New("obtención de ubicaciones de chunks incorrecta")
}


//Lee el log y retorna la lista de libros disponibles
func ListaLibrosLog() (string, error){
	
	var listaLibros string

	f, err := os.Open("logdata.txt")
	if (err != nil){
		f.Close()
		return "Error en abrir log", err
	}
	defer f.Close()

	//hace Splits por cada linea por defecto.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		
		words := strings.Fields(scanner.Text())

		//si el segundo elemento en la línea es un número
		//entonces el primero es un título de libro
		if _, err := strconv.Atoi(words[1]); err == nil {
			listaLibros += words[0] + "\n" //(?)
		}
	}

	return listaLibros, nil
}


//Enviar direcciones de chunks desde el Log
func (s *server) ChunkInfoLog(ctx context.Context, libro *book.ChunksInfo) (*book.ChunksInfo, error){
	
	localizacion, err := localizacionChunks(libro.NombreLibro)
	return &book.ChunksInfo{Info: localizacion}, err
}


//Envia el listado de libros disponibles a los clientes que se lo solicitan
func (s *server) EnviarListaLibros(ctx context.Context, ok *book.ACK) (*book.ListaLibros, error){
	
	lista, err := ListaLibrosLog()

	if (ok.Ok == "ok") {
		return &book.ListaLibros{Lista: lista}, err
	}

	return &book.ListaLibros{Lista: "Error en mensaje enviado a NameNode"}, errors.New("ACK corrupto")
}


// CONEXIONES
//cliente Downloader
func serveCD(){
	
	listenCD, err := net.Listen("tcp", dActual + ":50512")
	failOnError(err, "Error de conexión con cliente downloader")
	
	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})
	
	log.Fatalln(srv.Serve(listenCD))
}


//DataNode A
func serveDNA(){

	listenDNA, err := net.Listen("tcp", dActual + ":50506")
	failOnError(err, "Error de conexión con DataNode A")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNA))
}


//DataNode B
func serveDNB(){

	listenDNB, err := net.Listen("tcp", dActual + ":50507")
	failOnError(err, "Error de conexión con DataNode B")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNB))
}


//DataNode C
func serveDNC(){
	
	listenDNC, err := net.Listen("tcp", dActual + ":50508")
	failOnError(err, "Error de conexión con DataNode C")

	srv := grpc.NewServer()
	book.RegisterBookServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenDNC))
}


func main() {
	
//servers
	go serveCD()
	go serveDNA()
	go serveDNB()
	go serveDNC()

}

