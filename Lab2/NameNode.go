package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"context"
	"math/rand"
)

var tipo_al string
var deadline int32
var name string

//DEFINIR NOMBRES DE DATANODES CON IP CORRESPONDIENTES
dA := "" //ip maquina virtual datanode A
dB := "" //ip maquina virtual datanode B
dC := "" //ip maquina virtual datanode C



type server struct {
	book.UnimplementedBoookServiceServer
}

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


//para recibir entradas para el log
//distribuido
func (s *server) escribirLogD(prop *book.PropuestaLibro) (*book.ACK, error) {
	
	f, err := os.OpenFile("logdata.txt", os.O_WRONLY|os.O_APPEND, 0644)
	failOnError(err, "Error abriendo log")
    defer f.Close()

    _, err2 := f.WriteString(prop.NombreLibro + " " + prop.CantChunks + "\n" + prop.Propuesta + "\n")
	failOnError(err, "Error escribiendo en log")
	
	return &book.ACK{Ok: "listo"}, nil
}


//centralizado
func escribirLogCen(prop string, nombreL string, cant int32){
	
	f, err := os.OpenFile("logdata.txt", os.O_WRONLY|os.O_APPEND, 0644)
	failOnError(err, "Error abriendo log")
    defer f.Close()

    _, err2 := f.WriteString(nombreL + " " + cant + "\n" + prop + "\n")
	failOnError(err, "Error escribiendo en log")

	fmt.Println("Escritura en log exitosa")
	return nil
}


//verifica si hay un DataNode caido
func checkDatanode(dn string){

	deadline = 2 //segundos que espera para ver si hay conexión
	
	if (dn == dA){
		name = "datanode A"
	}
	if (dn == dB){
		name = "datanode B"
	}
	if (dn == dC){
		name = "datanode C"
	}

	conn, err := grpc.Dial(dn, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(deadline)*time.Second)
    if err != nil {
		log.Fatalf("No se pudo conectar con %v: %v", name, err)
		return false
    }
	defer conn.Close()
	
	c := book.NewBookServiceClient(conn) //necesario?? es solo chekear pero no aun conectar 4real

	conn.Close()
	return true
}


func analizarPropuesta(prop *book.PropuestaLibro) (bool, error){
	//retorna true si acepta propuesta y false si rechaza
	fmt.Println("Analizando la propuesta...\n")

	//book.PropuestaLibro sabe cuales datanodes se pretenden usar en la propuesta, por ejemplo
	//si prop.DatanodeA es true entonces se usaría en la propuesta, pero si esta caído se rechaza la propuesta

	//revisa si hay un datanode de la propuesta caído
	if (prop.DatanodeA && !checkDatanode(dA)){
		fmt.Println("Se rechaza la propuesta\n")
		return false, nil
	}
	if (prop.DatanodeB && !checkDatanode(dB)){
		fmt.Println("Se rechaza la propuesta\n")
		return false, nil
	} 
	if (prop.DatanodeC && !checkDatanode(dC)){
		fmt.Println("Se rechaza la propuesta\n")
		return false, nil
	}

	return true, nil
}


//Genera una nueva propuesta con una lista generada aleatoriamente
//Recibe el mensaje con la propuesta original para tener
//el nombre del libro y la cantidad de chunks
func generarNuevaPropuesta(prop *book.PropuestaLibro) string{

	//inicio de propuesta
	n := prop.CantChunks
	Prop := prop.NombreLibro + " " + n + "\n"

	//arreglo con valores aleatorios de DataNodes
	intProp := make([]int, n)

	for i := 0; i < n; i++{
		intProp[i] = rand.Intn(3)
	}
	
	var dAct string
	for i = 0; i < n; i++{ 
		
		switch intProp[i]{
		case 0:
			dAct = dA
		case 1:
			dAct = dB
		case 2:
			dAct = dC
		}
		
		Prop += prop.NombreLibro + "_" + strconv.Itoa(i) + " " + dAct + "\n"
		
		/*"nombre_archivo_chunk_0 ipDNx\n
		 nombre_archivo_chunk_1 ipDNy\n
		 nombre_archivo_chunk_2 ipDNz..." */
	}

	return Prop
}



//solo para centralizado
func (s *server) recibirPropDatanode(ctx context.Context, prop *book.PropuestaLibro) (*book.ACK, error){

//Cuando un DataNode envia una propuesta, se analiza la Propuesta
//Si se rechaza se genera una nueva y se analiza hasta que analizarPropuesta sea true
//Luego se escribe en el log la propuesta

	Prop := prop.PropuestaLibro

	for {
		if analizarPropuesta(Prop) {
			break
		}
		Prop = generarNuevaPropuesta(prop)
	}

	escribirLogCen(Prop, prop.NombreLibro, prop.CantChunks)
}


//Responde al Cliente Downloader con las ubicaciones de los chunks de libro solicitado
func localizacionChunks(nombreL string) string{

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
				return listachunks
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
}


//Lee el log y retorna la lista de libros disponibles
func ListaLibrosLog() string{
	
	var listaLibros string

	f, err := os.Open("logdata.txt")
	failOnError(err, "Error en abrir log")
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

	return listaLibros
}


//Enviar direcciones de chunks desde el Log
func (s *server) ChunkInfoLog(ctx context.Context, libro *book.ChunksInfo) *book.ChunksInfo{
	return &book.ChunksInfo{info: localizacionChunks(libro.NombreLibro)}
}


//Envia el listado de libros disponibles a los clientes que se lo solicitan
func (s *server) EnviarListaLibros(ctx context.Context) *book.ListaLibros{
	return &book.ListaLibros{Lista: ListaLibrosLog()}
}



func exclusionMCen(){
	var q [2]int
	
}


func main() {

//revisar si acá es el input 
	log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")

	_, err := fmt.Scanf("%d", &tipo_al)

	for (err != nil){

		log.Println("Tipo ingresado inválido!\n")
		log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")
		
		_, err = fmt.Scanf("%d", &tipo_al)
	}


	//conexión a cliente Downloader
	listenCD, err := net.Listen("tcp", addressCD)
	failOnError(err, "Error de conexión con cliente downloader")

	srv := grpc.NewServer()
	book.RegisterBookServiceClient(srv, &server{})

	log.Fatalln(srv.Serve(listenCD))


	//conexión a cliente Uploader
	listenCU, err := net.Listen("tcp", addressCU)
	failOnError(err, "Error de conexión con cliente downloader")

	srv := grpc.NewServer()
	book.RegisterBookServiceClient(srv, &server{})

	log.Fatalln(srv.Serve(listenCU))

	//conexión a DataNodes

	//proceder según algoritmo
	if (tipo_al == "c"){
		recibirPropDatanode(/* ? */)
		
	} else{
		escribirLogD(/* context.Context, *book.PropuestaLibro ???  */)
	}

}

