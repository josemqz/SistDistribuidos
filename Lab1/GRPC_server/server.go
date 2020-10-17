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
	// estados-> bdg: bodega, tr: en tránsito , rec: recibido , nr: no recibido
}

type PackageSeguimiento struct{
	id_pkg string
	estado_pkg string
	id_camion int32
	num_seguimiento int32
	num_intentos int32
}

//codigo de seguimiento para cada pedido
var cod_tracking int32 = 0

//colas de pedidos
var PaquetesRetail []Package
var PaquetesPri []Package
var PaquetesNormal []Package

//registro de pedidos
var Registros []RegPedido
var RegistroSeguimiento []PackageSeguimiento


func main() {

	log.Println("Server running ...")


//conexión clientes
	listenCliente, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

//conexión camiones
	listenCamion, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalln(err)
	}


	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenCliente))
	log.Fatalln(srv.Serve(listenCamion))
	
	//colas rabbitmq con financiero

}


func (s *server) PedidoCliente(ctx context.Context, pedido *logis.Pedido) (*logis.CodSeguimiento, error) {
	
	log.Println("Pedido recibido")

	//auxiliar para guardar valor de cod_tracking
	help_cod_tracking := cod_tracking
	var tipoP string

	//se escoge tipo de paquete
	if (pedido.Tienda == "pyme"){
		if (pedido.Prioritario == 0){
			tipoP = "normal" 
		} else if (pedido.Prioritario == 1){
			tipoP = "prioritario"
		}
	} else {
		tipoP = "retail"
		cod_tracking = 0
	}

	nownow := time.Now()
	log.Println(nownow)
	var pId = pedido.Id
	var pVal = pedido.Valor

	//guardar en registro
	Registros = append(Registros, RegPedido{timestamp:nownow, 
											id:pId, 
											tipo: tipoP, 
											nombre: pedido.Producto, 
											valor: pVal, 
											origen: pedido.Tienda, 
											destino: pedido.Destino, 
											num_seguimiento: cod_tracking})

	//transformar pedido en paquete
	pkg := Package{id: pId
					num_seguimiento: cod_tracking
					tipo: tipoP
					valor: pVal
					num_intentos: 0
					estado: "bdg"}

	//guardar paquete en cola
	if (tipoP == "normal") {
		PaquetesNormal := append(PaquetesNormal, pkg)
	} else if (tipoP == "prioritario") {
		PaquetesPri := append(PaquetesPri, pkg)
	} else {
		PaquetesRetail := append(PaquetesRetail, pkg)
	}

	//se vuelve al valor original del codigo de seguimiento
	cod_tracking = help_cod_tracking

	if (tipoP != "retail"){
		cod_tracking += 1
	}

	return &logis.CodSeguimiento{Codigo: cod_tracking}, nil
}


//Cliente solicita estado de paquete
func (s *server) SeguimientoCliente(ctx context.Context, cs *logis.CodSeguimiento) (*logis.EstadoPedido, error) {

	for _, regseg := range RegistroSeguimiento{
		if (cs.GetCodigo() == regseg.num_seguimiento){
			return &logis.EstadoPedido{Estado: regseg.estado_pkg}, nil			
		}
	}
	return &logis.EstadoPedido{Estado: "nulo"}, errors.New("No se encontró el paquete pedido en registro.")
}


//camión solicita un paquete
func (s *server) PedidoACamion(tC *logis.tipoCam, mIntento *logis.mIntento)(*logis.Package, *logis.Geo){

	tipoCam := tC.GetTipo()
	intento := mIntento.GetIntento()
	var flagIntento = false

	//server revisa en orden de prioridad las colas para saber qué paquete enviar
	if tipoCam == "normal"{

		if len(PaquetesPri > 0){
			pkg := PaquetesPri[0]
			PaquetesPri = PaquetesPri[1:]
			flagIntento = true

		} else if PaquetesNormal > 0{
			pkg := PaquetesNormal[0]
			PaquetesPri = PaquetesNormal[1:]
			flagIntento = true

		} else if intento{
			//esperar hasta que llegue algo a prioritario y normal
			flagIntento = true
		}
	
	//retail
	} else {
		
		if len(PaquetesRetail > 0){
			pkg := PaquetesRetail[0]
			PaquetesPri = PaquetesRetail[1:]
			flagIntento = true

		} else if len(PaquetesPri > 0){
			pkg := PaquetesPri[0]
			PaquetesPri = PaquetesPri[1:]
			flagIntento = true

		} else if intento{
			//esperar hasta que llegue algo a retail y prioritario
			flagIntento = true
		}
	}
	
	//no hay paquetes luego de espera de camión
	if !intento && !flagIntento{
		return nil, nil
	}
	
	mensajePkg := &logis.Package{id: pkg.id;
								num_seguimiento: pkg.num_seguimiento;
								tipo: pkg.tipo;
								valor: pkg.valor;
								num_intentos: pkg.num_intentos;
								estado: pkg.estado;
	}
	
	//obtener origen y destino del registro
	for i := range PackageSeguimiento{
		if PackageSeguimiento[i].id_pkg == pkg.id{
			orig := PackageSeguimiento[i].origen
			dest := PackageSeguimiento[i].destino
			break
		}
	}
	
	if orig == nil || len(orig) == 0{
		log.Println("No se encontró paquete de cola en registro de pedidos")
		return (&logis.Package{}, &logis.Geo{})
	}

	mensajeGeo := &logis.Geo{origen: orig, destino: dest}
	
	return (mensajePkg, mensajeGeo)
}


//actualización de estado de paquetes al llegar camión a la Central
func (s *server) EstadoCamion(idCamion int32, dataPkg *logis.regCamion){

	//actualizar estados de paquetes en registro
	for i := range RegistroSeguimiento{
		if (dataPkg.Id_pkg == RegistroSeguimiento[i].Id){
			//condiciones para obtener estado de paquete
			//cantidad de intentos y origen (retail o pyme)
			RegistroSeguimiento[i].estado = dataPkg.
		}
	}

	//enviar resultados económicos a financiero
}
