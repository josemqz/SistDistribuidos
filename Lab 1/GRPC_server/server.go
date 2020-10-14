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

type Package struct{
	id string
	num_seguimiento int
	tipo string
	valor int
	num_intentos int
	estado string
	// estados-> bdg: bodega, tr: en tránsito , rec: recibido , r: no recibido
}

type PackageSeguimiento struct{
	id_pkg string
	estado_pkg string
	id_camion int
	num_seguimiento int
	num_intentos int
}

//codigo de seguimiento para cada pedido
cod_tracking := 0

//colas de pedidos
var PaquetesRetail []Package
var PaquetesPri []Package
var PaquetesNormal []Package

//registro de pedidos
var Registros []RegPedido
var RegistroSeguimiento []PackageSeguimiento

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
	
	//auxiliar para guardar valor de cod_tracking
	help_cod_tracking := cod_tracking

	//se escoge tipo de paquete
	if (pedido.tienda == "pyme"){
		if (pedido.prioritario == 0){
			tipoP := "normal" 
		} else if (pedido.prioritario == 1){
			tipoP := "prioritario"
		}
	} else {
		tipoP := "retail"
		cod_tracking := 0
	}


	reg := RegPedido{timestamp:time.now(), id:pedido.id, tipo: tipoP, 
					nombre: pedido.producto, valor: pedido.valor, origen: pedido.tienda, 
					destino: pedido.destino, num_seguimiento: cod_tracking}

	Registros := append(Registros, reg)


	cod_tracking := help_cod_tracking

	if (tipoP != "retail"){
		cod_tracking += 1
	}

	//guardar paquete en cola 
	if (tipoP = "normal") {
		PaquetesNormal := append(PaquetesNormal, reg)
	} else if (tipoP == "prioritario") {
		PaquetesPri := append(PaquetesPri, reg)
	} else {
		PaquetesRetail := append(PaquetesRetail, reg)
	}

	return &logis.Recibo{seguimiento: cod_tracking}, nil
}

func (s *server)SeguimientoCliente(ctx context.Context, cs CodSeguimiento) return (*logis.EstadoPedido, error) {

	for _, regseg := range RegistroSeguimiento{
		if (cs.GetCodigo() == regseg.num_seguimiento){
			return &logis.EstadoPedido{estado: regseg.estado_pkg}, nil
		} else {
			return _, errors.New("No se encontró el paquete pedido en registro.")
		}
	}	
}


func (s *server) PedidoACamion(){
	
	//cola a camionid-paquete, tipo de paquete, valor, origen, destino, número de intentos
y fecha de entrega

}

/*
compile: ## Compile the proto file.
	protoc -I GRPC_proto/ GRPC_proto/cliente_logis.proto --go_out=plugins=grpc:GRPC_proto/
 
.PHONY: server
server: ## Build and run server.
	go build -race -ldflags "-s -w" -o bin/server server/main.go
	bin/server
 
.PHONY: client
client: ## Build and run client.
	go build -race -ldflags "-s -w" -o bin/client client/main.go
	bin/client
*/