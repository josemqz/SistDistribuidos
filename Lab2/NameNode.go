package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"context"


)


var tipo_al string
var dA string //datanode A
var dB string //datanode B
var dC string //datanode C
var deadline int32
var name string

//DEFINIR NOMBRES DE DATANODES CON IP CORRESPONDIENTES
dA = "" //ip maquina virtual datanode A
dB = "" //ip maquina virtual datanode B
dC = "" //ip maquina virtual datanode C



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
func (s *server) escribirLogD(ctx context.Context, prop *book.PropuestaLibro) (*book.ACK, error) {
	
	f, err := os.OpenFile("logdata.txt", os.O_WRONLY|os.O_APPEND, 0644)

    if err != nil {
        log.Fatalf("error al abrir el archivo: %s", err)
    }

    defer f.Close()

    _, err2 := f.WriteString(prop.NombreLibro + " " + prop.CantChunks + "\n" + prop.Propuesta + "\n")

    if err2 != nil {
        log.Fatalf("error al escribir en archivo: %s", err)
    }
	
	return &book.ACK{Ok: "listo"}, nil
}

//centralizado
func escribirLogCen(prop string, nombreL string, cant int32){
	
	f, err := os.OpenFile("logdata.txt", os.O_WRONLY|os.O_APPEND, 0644)

    if err != nil {
        log.Fatalf("error al abrir el archivo: %s", err)
    }

    defer f.Close()

    _, err2 := f.WriteString(nombreL + " " + cant + "\n" + prop+ "\n")

    if err2 != nil {
        log.Fatalf("error al escribir en archivo: %s", err)
	}
	
	fmt.Println("listo")
	return nil
}


//verifica si hay un datanode caido
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
	return true
}


func (s *server) analizarPropuesta(ctx context.Context, prop *book.PropuestaLibro) (bool, error){
	//retorna true si acepta propuesta y false si rechaza
	fmt.Println("Analizando la propuesta...\n")

	//book.PropuestaLibro sabe cuales datanodes se pretenden usar en la propuesta, por ejemplo
	//si prop.DatanodeA es true entonces se usaría en la propuesta, pero si esta caído se rechaza la propuesta

	//revisa si hay un datanode de la propuesta caído
	if (prop.DatanodeA == true && checkDatanode(dA) == false){
		fmt.Println("Se rechaza la propuesta\n")
		return false, nil
	}
	if (prop.DatanodeB == true && checkDatanode(dB) == false){
		fmt.Println("Se rechaza la propuesta\n")
		return false, nil
	} 
	if (prop.DatanodeC == true && checkDatanode(dC) == false){
		fmt.Println("Se rechaza la propuesta\n")
		return false, nil
	}

	return true, nil
}



func generarNuevaPropuesta(){

}



func (s *server) recibirPropDatanode(ctx context.Context, prop *book.PropuestaLibro) (*book.ACK, error){
	//solo para centralizado (si es distr. el namenode va directo a escribir al log la propuesta
	//aceptada por los otros datanodes - diagrama secuencia)

//Cuando datanode envia propuesta aca se usa
//la funcion analizarPropuesta y si esta da false
//se usa generarNuevaPropuesta y se analiza
//hasta que analizarPropuesta de true
//luego se escribe en el log con escribirLogCen
	
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	//for infinito
	for {
		if analizarPropuesta(ctx,prop) == true {
			break
		}
		generarNuevaPropuesta(ctx,prop) //argumentos correctos?
	}
	escribirLogCen(prop.Propuesta, prop.NombreLibro, prop.CantChunks)
	

}

func conectarCliente(){
	//Cliente Uploader
	//Cliente Downloader
}


//Responde al Cliente Downloader con las ubicaciones de los chunks de libro solicitado
func localizacionChunks(nombreL string){

	f, err := os.Open("logdata.txt")
	failOnError(err, "Error en cargar log")
	defer f.Close()

	// hace Splits por cada linea por defecto.
	scanner := bufio.NewScanner(f)

	line := 1
	
	for scanner.Scan() {
    	if strings.Contains(scanner.Text(), "yourstring") {
        	return line, nil
    	}

    	line++
	}

	if err := scanner.Err(); err != nil {
    	// Handle the error
	}
	
	
	
	
	
	
	/*file, err := os.Open("logdata.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = f.Close(); err != nil {
            log.Fatal(err)
        }
    }()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {             // internally, it advances token based on sperator
        fmt.Println(scanner.Text())  // token in unicode-char
        fmt.Println(scanner.Bytes()) // token in bytes

    }*/
}

//Envia el listado de libros disponibles a los clientes que se lo solicitan
func listadoLibrosAv(){}


func main() {

	//revisar si aca es el input o en cliente uploader unicamente
	log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")

	_, err := fmt.Scanf("%d", &tipo_al)

	for (err != nil){

		log.Println("Tipo ingresado inválido!\n")
		log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")
		
		_, err = fmt.Scanf("%d", &tipo_al)
	}

	if (tipo_al == "c"){
		recibirPropDatanode(/* ? */)
	}else{
		escribirLogD(/* context.Context, *book.PropuestaLibro ???  */)
	}

}

