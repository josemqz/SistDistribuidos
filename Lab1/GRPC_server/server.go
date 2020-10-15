package server

import (
	"context"
	"log"
	"net"
	"time"
	"errors"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
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
	valor int32
	origen string
	destino string
	num_seguimiento int32
}

type Package struct{
	id string
	num_seguimiento int32
	tipo string
	valor int32
	num_intentos int32
	estado string
	// estados-> bdg: bodega, tr: en tránsito , rec: recibido , r: no recibido
}

type PackageSeguimiento struct{
	id_pkg string
	estado_pkg string
	id_camion int32
	num_seguimiento int32
	num_intentos int32
}

//codigo de seguimiento para cada pedido
var cod_tracking int32 = 0 //constante

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

func (s *server) PedidoCliente(ctx context.Context, pedido *logis.Pedido) (*logis.CodSeguimiento, error) {
	
	log.Println("Pedido recibido")

	//guardar en registro
	
	//auxiliar para guardar valor de cod_tracking
	help_cod_tracking := cod_tracking
	var tipoP string = "1"

	//se escoge tipo de paquete
	if (pedido.Tienda == "pyme"){
		if (pedido.Prioritario == 0){
			tipoP := "normal" 
		} else if (pedido.Prioritario == 1){
			tipoP := "prioritario"
		}
	} else {
		tipoP := "retail"
		cod_tracking = 0
	}

	nownow := time.Now()
	log.Println(nownow)

	reg := RegPedido{timestamp:nownow, id:pedido.Id, tipo: tipoP, 
					nombre: pedido.Producto, valor: pedido.Valor, origen: pedido.Tienda, 
					destino: pedido.Destino, num_seguimiento: cod_tracking}

	Registros := append(Registros, reg)


	cod_tracking = help_cod_tracking

	if (tipoP != "retail"){
		cod_tracking += 1
	}

	//guardar paquete en cola 
	/*
	if (tipoP == "normal") {
		PaquetesNormal := append(PaquetesNormal, reg)
	} else if (tipoP == "prioritario") {
		PaquetesPri := append(PaquetesPri, reg)
	} else {
		PaquetesRetail := append(PaquetesRetail, reg)
	}
*/
	return &logis.CodSeguimiento{Codigo: cod_tracking}, nil
}

func (s *server) SeguimientoCliente(ctx context.Context, cs *logis.CodSeguimiento) (*logis.EstadoPedido, error) {

	for _, regseg := range RegistroSeguimiento{
		if (cs.GetCodigo() == regseg.num_seguimiento){
			return &logis.EstadoPedido{Estado: regseg.estado_pkg}, nil			
		}
	}
	return &logis.EstadoPedido{Estado: "nulo"}, errors.New("No se encontró el paquete pedido en registro.")
}


func (s *server) PedidoACamion(){
	
	//cola a camionid-paquete, tipo de paquete, valor, origen, destino, número de intentos
	//y fecha de entrega

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