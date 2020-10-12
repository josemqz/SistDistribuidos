package main

import (
	"context"
	"log"
	"time"

	"github.com/josemqz/SistDistribuidos"
	"google.golang.org/grpc"
)

func main() {

	log.Println("Client running ...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := logis.NewLogisServiceClient(conn)
	

	//leer csv


	//if tienda == "pyme"{} //habiendo leido el archivo csv
	request := &logis.Pedido{
		id := ,
		producto := ,
		valor := ,
		tienda := ,
		destino := ,
		prioritario := ,
	}

	//else
	request := &logis.Pedido{
		id := ,
		producto := ,
		valor := ,
		tienda := ,
		destino := ,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.PedidoCliente(ctx, request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Respuesta:", response.GetRecibo())
}