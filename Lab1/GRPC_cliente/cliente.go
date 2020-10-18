package cliente

import (
	"log"
	"fmt"
	"bufio"
	"os"
	"time"
	"context"
	"strconv"
	
	"encoding/csv"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"
)


func main() {
	
	//pedir tiempo de espera entre pedidos por input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresar tiempo de espera entre pedidos: ")

	read, err := reader.ReadString('\n')
	usr_time, _ := strconv.Atoi(read)
	
	for err != nil || usr_time <= 0{

		fmt.Print("Tiempo ingresado inválido!\n")
		fmt.Print("Ingresar tiempo de espera entre pedidos: ")

		read, err = reader.ReadString('\n')
		usr_time, _ = strconv.Atoi(read)
	}

	fmt.Print("Seleccionar opción: ")
	fmt.Print("-------------------------------------")
	fmt.Print("  1. Realizar pedidos			    |")
	fmt.Print("  2. Solicitar seguimiento de pedido |")
	fmt.Print("-------------------------------------\n")

	read, err = reader.ReadString('\n')
	opcion, _ := strconv.Atoi(read)

	for err != nil || (opcion != 1 && opcion != 2) {

		fmt.Print("Opción inválida\n")

		fmt.Print("Seleccionar opción: ")
		fmt.Print("-------------------------------------")
		fmt.Print("  1. Realizar pedidos			    |")
		fmt.Print("  2. Solicitar seguimiento de pedido |")
		fmt.Print("-------------------------------------\n")

		read, err = reader.ReadString('\n')
		opcion, _ = strconv.Atoi(read)

	}

	log.Println("Cliente corriendo...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	cliente := logis.NewLogisServiceClient(conn)


	if opcion == 1{

		var tipoT int = 1

		fmt.Print("Seleccionar tipo de tienda:")
		fmt.Print("-------------")
		fmt.Print("  1. Pyme   |")
		fmt.Print("  2. Retail |")
		fmt.Print("-------------\n")
		read, err = reader.ReadString('\n')
		tipoT, _ = strconv.Atoi(read)

		for err != nil || (tipoT != 1 && tipoT != 2) {

			fmt.Print("Opción inválida\n")

			fmt.Print("Seleccionar tipo de tienda:")
			fmt.Print("-------------")
			fmt.Print("  1. Pyme   |")
			fmt.Print("  2. Retail |")
			fmt.Print("-------------\n")
			
			read, err = reader.ReadString('\n')
			tipoT, _ = strconv.Atoi(read)

		}

		if (tipoT == 1) {

			records_pyme := leercsv("pymes.csv")
			enviarPedidos(records_pyme, cliente, time.Duration(usr_time))

		} else if (tipoT == 2){

			records_retail := leercsv("retail.csv")
			enviarPedidos(records_retail, cliente, time.Duration(usr_time))
		}
	
	} else if opcion == 2{
		
		fmt.Print("Ingresar código de seguimiento de pedido: ")
			
		read, err = reader.ReadString('\n')
		cod, _ := strconv.ParseInt(read, 10, 32)
	
		for (err!=nil) {
			fmt.Print("Ingresar código de seguimiento de pedido: ")
			
			read, err = reader.ReadString('\n')
			cod, _ = strconv.ParseInt(read, 10, 32)
		}

		doSeguimiento(cliente, cod)
	}

}


func leercsv(arch string) ([][]string){
	
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


func enviarPedidos(records [][]string, client_var logis.LogisServiceClient, usr_time time.Duration){

	//se realiza cada pedido
	for _, rec := range records{

		request := &logis.Pedido{
			Id: rec[0],
			Producto: rec[1],
			Valor: strconv.ParseInt(rec[2], 10, 32),
			Tienda: rec[3],
			Destino: rec[4],
			Prioritario: strconv.ParseInt(rec[5], 10, 32),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		response, err := client_var.PedidoCliente(ctx, request)
		if err != nil {
			log.Fatalln(err)
		}
	
		cod_seguimiento := response.GetCodigo()
		log.Println("Pedido realizado\nCodigo de seguimiento: ", cod_seguimiento, "\n")
		
		//wait
		time.Sleep(time.Duration(usr_time) * time.Second)
	}
}

func doSeguimiento(ctx context.Context, client_var logis.LogisServiceClient, cod int32) {

	seguimientoPkg := &logis.CodSeguimiento{Codigo: cod}

	estado_pkg, err := client_var.SeguimientoCliente(ctx, seguimientoPkg)
	if err != nil {
		log.Fatalln(err)
	}
	
	log.Println("Estado de pedido: " + estado_pkg.GetEstado())
}
