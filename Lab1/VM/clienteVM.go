package main

import (
	"log"
	"fmt"
	"bufio"
	"io"
	"os"
	"time"
	"context"
	"strconv"

	"encoding/csv"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"
)

var usr_time int


func main() {
	
	log.Print("Ingresar tiempo de espera entre pedidos: ")

	_, err := fmt.Scanf("%d", &usr_time)

	for (err != nil){

		log.Println("Tiempo ingresado inválido!\n")
		log.Print("Ingresar tiempo de espera entre pedidos: ")
		
		_, err = fmt.Scanf("%d", &usr_time)
	}

	var opcion int

	log.Println("-------------------------------------")
	log.Println("  1. Realizar pedidos			      |")
	log.Println("  2. Solicitar seguimiento de pedido |")
	log.Println("-------------------------------------\n")
	log.Print("Seleccionar opción: ")

	_, err = fmt.Scanf("%d", &opcion)

	for (err != nil) || (opcion != 1 && opcion != 2) {

		log.Println("Opción inválida\n")

		log.Println("-------------------------------------")
		log.Println("  1. Realizar pedidos			      |")
		log.Println("  2. Solicitar seguimiento de pedido |")
		log.Println("-------------------------------------\n")
		log.Print("Seleccionar opción: ")

		_, err = fmt.Scanf("%d", &opcion)

	}

	log.Println("Cliente corriendo...")

	conn, err := grpc.Dial("10.6.40.157:50051", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err, "Error en conexión a servidor")
	defer conn.Close()

	cliente := logis.NewLogisServiceClient(conn)

	if (opcion == 1){

		var tipoT int = 1

		log.Println("-------------")
		log.Println("  1. Pyme   |")
		log.Println("  2. Retail |")
		log.Println("-------------\n")
		log.Print("Seleccionar tipo de tienda: ")

		_, err = fmt.Scanf("%d", &tipoT)

		for ((err != nil) || (tipoT != 1 && tipoT != 2)) {

			log.Println("Opción inválida\n")

			log.Println("Seleccionar tipo de tienda:")
			log.Println("-------------")
			log.Println("  1. Pyme   |")
			log.Println("  2. Retail |")
			log.Println("-------------\n")
			
			_, err = fmt.Scanf("%d", &tipoT)

		}

		if (tipoT == 1) {

			log.Println("Leyendo csv pymes...")
			records_pyme := leercsv("pymes.csv")
			enviarPedidos("pyme", records_pyme, cliente)

		} else if (tipoT == 2){

			log.Println("Leyendo csv retail...")
			records_retail := leercsv("retail.csv")
			enviarPedidos("retail", records_retail, cliente)
		}
	
	} else if (opcion == 2){
		
		var cod int32
		var codHelp int

		log.Println("Ingresar código de seguimiento de pedido: ")
		
		_, err = fmt.Scanf("%d", &codHelp)
		cod = int32(codHelp)
	
		for (err != nil) {
			log.Println("Ingresar código de seguimiento de pedido: ")
			
			_, err = fmt.Scanf("%d", &codHelp)
			cod = int32(codHelp)
		}

		doSeguimiento(cliente, cod)
	}

}


func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}


func leercsv(arch string) ([][]string){
	
	csv_arch, err := os.Open(arch)
	failOnError(err, "Error abriendo archivo csv")
	defer csv_arch.Close()
	
	// Skip first row (line)
	row1, err := bufio.NewReader(csv_arch).ReadSlice('\n')
	failOnError(err, "Error interpretando csv (lectura primera línea)")

	_, err = csv_arch.Seek(int64(len(row1)), io.SeekStart)
	failOnError(err, "Error interpretando csv (lectura primera línea)")
	
	r := csv.NewReader(csv_arch)
	records, err := r.ReadAll()
	failOnError(err, "Error interpretando archivo csv")
	
	csv_arch.Close()
	log.Println("csv leído")

	return records
}


func enviarPedidos(tipoT string, records [][]string, client_var logis.LogisServiceClient){

	log.Println("funcion: enviarPedidos")

	//se realiza cada pedido
	for _, rec := range records{

		valorHelp, _ := strconv.Atoi(rec[2])
		
		request := &logis.Pedido{
			Id: rec[0],
			Producto: rec[1],
			Valor: int32(valorHelp),
			Tienda: rec[3],
			Destino: rec[4],
		}

		if (tipoT == "pyme"){
			priorHelp, _ := strconv.Atoi(rec[5])
			request.Prioritario = int32(priorHelp)
		}

		log.Println("Pedido: ", request)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
	
		response, err := client_var.PedidoCliente(ctx, request)
		failOnError(err, "Error enviando pedido al servidor")		
	
		cod_seguimiento := response.GetCodigo()
		log.Println("Pedido realizado\nCodigo de seguimiento: ", cod_seguimiento, "\n")
		
		//wait
		time.Sleep(time.Duration(usr_time) * time.Second)
	}
}


func doSeguimiento(client_var logis.LogisServiceClient, cod int32) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	seguimientoPkg := &logis.CodSeguimiento{Codigo: cod}

	estado_pkg, err := client_var.SeguimientoCliente(ctx, seguimientoPkg)
	
	failOnError(err, "Error de seguimiento de pedido")

	
	log.Println("Estado de pedido: " + estado_pkg.GetEstado())
}
