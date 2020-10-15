package cliente

import (
	"context"
	"log"
	"time"

	"github.com/josemqz/SistDistribuidos"
	"github.com/josemqz/SistDistribuidos/Lab1/logis"
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


func enviarPedidos(records []byte, client_var service /*o logisService*/, usr_time int) (message /*o Recibo*/){

	//seguimiento_opcional := 1 //contador para solicitar info. de paquete cada ciertos pedidos

	//se realiza cada pedido
	for _, rec := range records{

		request := &logis.Pedido{
			id := rec[0],
			producto := rec[1],
			valor := rec[2],
			tienda := rec[3],
			destino := rec[4],
			prioritario := rec[5],
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	
		response, err := client_var.PedidoCliente(ctx, request)
		if err != nil {
			log.Fatalln(err)
		}
	
		cod_seguimiento := response.GetCodigo()
		log.Println("Pedido realizado\nCodigo de seguimiento: ", cod_seguimiento, "\n")
		
		//cómo se debería entregar el codigo de seguimiento al cliente?
		
		/*
		//solicitar estado de paquete
		if (seguimiento_opcional%4 == 0){
			
			seguimientoPkg := &logis.CodSeguimiento{codigo: cod_seguimiento - 2}

			estado_pkg, err := client_var.SeguimientoCliente(ctx, seguimientoPkg)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println(estado_pkg)

		}

*/
		//seguimiento_opcional += 1

		//wait
		time.Sleep(usr_time * time.Second)
	} 
}

func doSeguimiento(client_var service /*o logisService*/, cod int) {

	seguimientoPkg := &logis.CodSeguimiento{codigo: cod}

	estado_pkg, err := client_var.SeguimientoCliente(ctx, seguimientoPkg)
	if err != nil {
		log.Fatalln(err)
	}
	
	log.Println(estado_pkg.GetEstado())
}

/*
export GOROOT=/media/joseesmuyoriginal/opt/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
*/

func main() {
	
	//pedir tiempo de espera entre pedidos por input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresar tiempo de espera entre pedidos: ")
	usr_time, _ := reader.ReadString('\n')

	flag := true
	fmt.Print("Seleccionar opción: ")
	fmt.Print("-------------------------------------")
	fmt.Print("  1. Realizar pedidos			    |")
	fmt.Print("  2. Solicitar seguimiento de pedido |")
	fmt.Print("-------------------------------------")
	opcion, err := reader.ReadString('\n')

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
	
	}else if opcion == 2{

	}




}