package main

import (
/*	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"*/
	"log"
	//"math/rand"
	//"time"
	"reflect"
	
	"google.golang.org/grpc"
	//"github.com/josemqz/SistDistribuidos/Lab2/book"

)

func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func main(){

	connA, err := grpc.Dial("localhost" + ":50513", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err, "Error en conexión con DataNode A, no se podrá descargar libro")
	defer connA.Close()
	
	log.Println(reflect.TypeOf(connA))

	//clienteA = book.NewBookServiceClient(connA)

	connA.Close()


}
