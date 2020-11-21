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


func checkDatanode(dn string){
	//REVISAR que use parametro
	var estado bool

	//set deadline
	nodeDeadline := time.Now().Add(time.Duration(2) * time.Second)
	ctx, cancel := context.WithDeadline(ctx, nodeDeadline)
	
	// if caido retorna false / si online retorna true
	if ctx.Err() == context.Canceled{
		estado = false
	}else{
		estado = true
	}

	return estado
}

func (s *server) analizarPropuesta(ctx context.Context, prop *book.PropuestaLibro) (bool, error){
	//retorna true si acepta propuesta y false si rechaza

	//revisa si hay un datanode caido que se usa en la propuesta
	if (prop.DatanodeA == true && checkDatanode(dA) == false){
		return false, nil
	}
	if (prop.DatanodeB == true && checkDatanode(dB) == false){
		return false, nil
	} 
	if (prop.DatanodeC == true && checkDatanode(dC) == false){
		return false, nil
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


func main() {


	log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")

	_, err := fmt.Scanf("%d", &tipo_al)

	for (err != nil){

		log.Println("Tipo ingresado inválido!\n")
		log.Print("Ingresar tipo de algoritmo - c:centralizado / d:distribuido : ")
		
		_, err = fmt.Scanf("%d", &tipo_al)
	}


}

