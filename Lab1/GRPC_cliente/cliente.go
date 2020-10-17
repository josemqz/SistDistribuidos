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
	usr_time, _ := strconv.ParseInt(reader.ReadString('\n'), 10, 32)

	fmt.Print("Seleccionar opción: ")
	fmt.Print("-------------------------------------")
	fmt.Print("  1. Realizar pedidos			    |")
	fmt.Print("  2. Solicitar seguimiento de pedido |")
	fmt.Print("-------------------------------------")
	opcion, err := strconv.ParseInt(reader.ReadString('\n'), 10, 32)

	for (err!=nil) || (opcion != 1 && opcion != 2) {

		fmt.Print("Opción inválida")
		fmt.Print("Seleccionar opción: ")
		fmt.Print("-------------------------------------")
		fmt.Print("  1. Realizar pedidos			    |")
		fmt.Print("  2. Solicitar seguimiento de pedido |")
		fmt.Print("-------------------------------------")
		opcion, err := reader.ReadString('\n')

	}

	log.Println("Cliente corriendo...")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := logis.NewLogisServiceClient(conn)


	if opcion == 1{
		//leer CSVs
		records_pyme := leercsv("pymes.csv")
		records_retail := leercsv("retail.csv")
		
		//enviar pedidos pyme
		enviarPedidos(records_pyme, client, usr_time)
		enviarPedidos(records_retail, client, usr_time)
	
	} else if opcion == 2{
		
		fmt.Print("Ingresar código de seguimiento de pedido: ")
		cod, err := strconv.ParseInt(reader.ReadString('\n'), 10, 32)
	
		for (err!=nil) {
			fmt.Print("Ingresar código de seguimiento de pedido: ")
			cod, err := strconv.ParseInt(reader.ReadString('\n'), 10, 32)
		}

		doSeguimiento(ctx, client, cod)
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
