package main

import (
	"context"
	"log"
	"time"

	"github.com/josemqz/SistDistribuidos"
	"google.golang.org/grpc"
)

func leercsv(arch String) ([]byte){
	
	csv_arch, err := os.Open(arch)
	
	if err != nil {
		fmt.Println(err)
	}
	defer csv_arch.Close()
	
	r := csv.NewReader(csv_arch)
	records, err := r.ReadAll()
	
	csv_arch.Close()
	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	return records
}

func main() {

	log.Println("Client running ...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := logis.NewLogisServiceClient(conn)

	//pedir tiempo de espera por pedido por input

	//leer CSVs
	records_pyme := leercsv("pymes.csv")
	records_retail := leercsv("retail.csv")

	//enviar pedidos pyme
	for _, rec := range records_pyme {

		request := &logis.Pedido{
			id := rec[0],
			producto := rec[1],
			valor := rec[2],
			tienda := rec[3],
			destino := rec[4],
			prioritario := rec[5],
		}

		//wait
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		response, err := client.PedidoCliente(ctx, request)
		if err != nil {
			log.Fatalln(err)
		}
	
		log.Println("Respuesta:", response.GetRecibo())
	}

	/*
	request := &logis.Pedido{
		id := ,
		producto := ,
		valor := ,
		tienda := ,
		destino := ,
	}
	*/

}