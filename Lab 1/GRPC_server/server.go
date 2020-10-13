package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/josemqz/SistDistribuidos/Lab\ 1/GRPC_proto"
	"google.golang.org/grpc"
)

type server struct {
	logis.UnimplementedLogisServiceServer
}

type RegPedido struct{
	timestamp time.Time
	id string
	tipo string
	nombre string
	valor int
	origen string
	destino string
	num_seguimiento int
}

type PackageTracking struct{
	id string
	num_seguimiento int
	tipo string
	valor int
	numintentos int
	estado string
	// estados-> bdg: bodega, tr: en tránsito , rec: recibido , r: no recibido
}

//codigo de seguimiento para cada pedido
cod_tracking := 0

//colas de pedidos
var PaquetesRetail[]
var PaquetesPri[]
var PaquetesNormal[]

//registro de pedidos
var Registros []RegPedido


func main() {
	log.Println("Server running ...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))

/*gRPC clientes
		if llega mensaje:

		retail
			id	producto	valor	tienda	destino
		pyme
			id	producto	valor	tienda	destino prioritario
			
			enviar codigo de seguimiento a cliente

			//hay que usar un indice para saber en qué producto va (para no repetirlos)

			generar EDD de paquete para enviar a camión 

		*/
		
		/*gRPC camiones

		enviar id de seguimiento
		recibir estado (En bodega, En camino, Recibido o No Recibido)
		actualizar registro de seguimiento (seg)
		enviar a cliente que pidió seguimiento

		si un camión vuelve, hay que entregar datos de ganancias/pérdidas a financiero

		*/
	//}

	//colas rabbitmq con financiero

}

func (s *server) PedidoCliente(ctx context.Context, pedido *logis.Pedido) (*logis.Recibo, error) {
	
	log.Println("Pedido recibido")

	//guardar en registro

	//tipo de paquete
	if (pedido.tienda == "pyme"){
		if (pedido.prioritario == 0){
			tipoP := "normal" 
		} else if (pedido.prioritario == 1){
			tipoP := "prioritario"
		}
	} else {
		tipoP := "retail"
	}

	reg := RegPedido{timestamp:time.now(), id:pedido.id, tipo: tipoP, 
					nombre: pedido.producto, valor: pedido.valor, origen: pedido.tienda, 
					destino: pedido.destino, num_seguimiento: cod_tracking}

	Registros := append(Registros, reg)

	cod_tracking += 1

	//guardar paquete en cola 

	return &logis.Recibo{seguimiento: cod_tracking-1}, nil
}

func (s *server) PedidoACamion()
