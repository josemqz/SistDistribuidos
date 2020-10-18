package main

import (
	"context"
	"log"
	"net"
	"time"
	"errors"
	//"amqp"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"
)

type server struct {
	logis.UnimplementedLogisServiceServer
}

type RegPedido struct{
	timestamp string
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
	id_camion string
	num_seguimiento int32
	num_intentos int32
}

type financiero struct {
	Id          string `json: "id"`
	Tipo        string `json: "tipo"`
	Valor       int32  `json: "valor"`
	Estado      string `json: "estado"`
	Intentos    int32  `json: "intentos"`
	//Seguimiento string `json: "seguimiento"`
}

//codigo de seguimiento para cada pedido
var cod_tracking int32 = 0

//colas de pedidos
var PaquetesRetail []Package
var PaquetesPri []Package
var PaquetesNormal []Package

//registro de pedidos
var Registros = []RegPedido{}
var RegistroSeguimiento = []PackageSeguimiento{}


func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s", msg, err)
	}
}

//
//------------------------ FINANCIERO ------------------------
//

func haciaFinanciero(pak financiero){
	conn, err := amqp.Dial("amqp:/guest:guest@localhost:5672")
	failOnError(err, "error al conectar")
	defer conn.Close()

	ch,err := conn.Channel()
	failOnError(err, "error al abrir canal")
	defer ch.Close()

	q, err:= ch.QueueDeclare(
		"hello",
		false,   
		false,   
		false,  
		false,   
		nil,    
	)
	failOnError(err,"error al enviar mensaje")

}


func paqueteFinanciero(i int){
	
	var pakfi financiero
	var j = 0
	pakfi.Id = RegistroSeguimiento[i].id_pkg
	pakfi.Estado = RegistroSeguimiento[i].estado_pkg
	pakfi.Intentos = RegistroSeguimiento[i].num_intentos

	for j < len(Registros){
		if RegistroSeguimiento[i].id_pkg == Registros[j].id{
			pakfi.Valor = Registros[j].valor
			pakfi.Tipo = Registros[i].tipo
			break
		}
	}
	haciaFinanciero(pakfi)
}

//
//-------------------------- CAMIÓN --------------------------
//

//actualización de estado de paquetes al llegar camión a la Central
func (s *server) ResEntrega(ctx context.Context, pkg *logis.RegCamion) (*logis.ACK, error){

	//actualizar estados de paquetes en registro
	for i := range RegistroSeguimiento{
		if (pkg.Id == RegistroSeguimiento[i].id_pkg){

			RegistroSeguimiento[i].estado_pkg = pkg.Estado
			RegistroSeguimiento[i].num_intentos = pkg.Intentos

			paqueteFinanciero(i)

			return &logis.ACK{Ok: "ok"}, nil
		}
	}


	
	return &logis.ACK{Ok: "error"}, errors.New("No se encontró pedido en registro de seguimiento")

}


func checkColasHelper(tipoCam string, prevRetail bool)(Package, bool){

	if tipoCam == "normal"{
		if len(PaquetesPri) > 0{
			pkg := PaquetesPri[0]
			PaquetesPri = PaquetesPri[1:]
			log.Println("Cargando paquete prioritario...")
			return pkg, true

		} else if len(PaquetesNormal) > 0{
			pkg := PaquetesNormal[0]
			PaquetesPri = PaquetesNormal[1:]
			log.Println("Cargando paquete normal...")
			return pkg, true
		}
	
	} else {
		if len(PaquetesRetail) > 0{
			pkg := PaquetesRetail[0]
			PaquetesPri = PaquetesRetail[1:]
			log.Println("Cargando paquete de retail...")
			return pkg, true
		
		// condicion en caso de haber hecho recien un envio de retail
		} else if (len(PaquetesPri) > 0) && prevRetail{ 
			pkg := PaquetesPri[0]
			PaquetesPri = PaquetesPri[1:]
			log.Println("Cargando paquete prioritario...")
			return pkg, true
		}
	}
	
	return Package{}, false
}


//revisar colas mediante función auxiliar dependiendo del tipo de camión
//que está solicitando y si es la primera o segunda peticióń de este
func CheckColas(tipoCam string, numPeticion bool, prevRetail bool)(Package, bool){

	var pkg Package
	var flagPeticion bool

	if !numPeticion{
		pkg, flagPeticion = checkColasHelper(tipoCam, prevRetail)

	} else{
		//esperar hasta que llegue algo a prioritario y normal
		for {
			pkg, flagPeticion = checkColasHelper(tipoCam, prevRetail)
			if pkg.id != "" {
				break
			}
			
			time.Sleep(1 * time.Second)
			log.Println("Esperando que hayan paquetes en colas...\n")

		}
	}

	return pkg, flagPeticion
}


//camión solicita un paquete
func (s *server) PedidoACamion(ctx context.Context, tC *logis.InfoCam)(*logis.PackageYGeo, error){

	tipoCam := tC.GetTipo()
	idCam := tC.GetId()
	prevRetail := tC.GetPrevRetail()

	//numPeticion representa la primera o segunda petición de paquete del camión
	//flagPeticion es si tiene exito en la búsqueda en las colas
	numPeticion := tC.GetNumPeticion()
	var flagPeticion = false

	//server revisa en orden de prioridad las colas para saber qué paquete enviar
	pkg, flagPeticion := CheckColas(tipoCam, numPeticion, prevRetail)

	//no hay paquetes en colas luego del tiempo de espera del camión
	if !numPeticion && !flagPeticion{
		return &logis.PackageYGeo{}, nil
	}


	//obtener origen y destino del pedido
	var orig string = ""
	var dest string = ""

	for i := range Registros{
		if (Registros[i].id == pkg.id){
			
			orig = Registros[i].origen
			dest = Registros[i].destino

			break
		}
	}
	
	if orig == "" || len(orig) == 0{
		err := errors.New("No se encontró paquete de cola en registro de pedidos")
		return &logis.PackageYGeo{}, err
	}
	

	//actualizar estado en registro seguimiento
	for i := range RegistroSeguimiento{
		if (RegistroSeguimiento[i].id_pkg == pkg.id){
			RegistroSeguimiento[i].estado_pkg = "tr"
			RegistroSeguimiento[i].id_camion = idCam
			break
		}
	}

	//generar mensaje a enviar a camión con datos
	mensajePkg := &logis.PackageYGeo{Id: pkg.id,
									NumSeguimiento: pkg.num_seguimiento,
									Tipo: pkg.tipo,
									Valor: pkg.valor,
									NumIntentos: pkg.num_intentos,
									Estado: "tr",
									Origen: orig,
									Destino: dest,
	}

	return mensajePkg, nil
}


//
//-------------------------- CLIENTE --------------------------
//

//Cliente solicita estado de paquete
func (s *server) SeguimientoCliente(ctx context.Context, cs *logis.CodSeguimiento) (*logis.EstadoPedido, error) {

	for _, regseg := range RegistroSeguimiento{
		if (cs.GetCodigo() == regseg.num_seguimiento){
			return &logis.EstadoPedido{Estado: regseg.estado_pkg}, nil			
		}
	}
	return &logis.EstadoPedido{Estado: "nulo"}, errors.New("No se encontró el paquete pedido en registro.")
}


func (s *server) PedidoCliente(ctx context.Context, pedido *logis.Pedido) (*logis.CodSeguimiento, error) {
	
	log.Print("Pedido recibido")

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
	var pId = pedido.Id
	var pVal = pedido.Valor

	//guardar en registro
	Registros = append(Registros, RegPedido{timestamp:nownow.Format("2006-01-02 15:04:05"), 
											id:pId, 
											tipo: tipoP, 
											nombre: pedido.Producto, 
											valor: pVal, 
											origen: pedido.Tienda, 
											destino: pedido.Destino, 
											num_seguimiento: cod_tracking,
	})

	//guardar en registro de seguimiento
	RegistroSeguimiento = append(RegistroSeguimiento, PackageSeguimiento{id_pkg: pId,
														estado_pkg: "bdg",
														id_camion: "",
														num_seguimiento: cod_tracking,
														num_intentos: 0,
	})
	
	//transformar pedido en paquete
	pkg := Package{id: pId,
					num_seguimiento: cod_tracking,
					tipo: tipoP,
					valor: pVal,
					num_intentos: 0,
					estado: "bdg"}


	//guardar paquete en cola
	if (tipoP == "normal") {
		PaquetesNormal = append(PaquetesNormal, pkg)
		log.Println("Paquete ingresado a cola normal\n")

	} else if (tipoP == "prioritario") {
		PaquetesPri = append(PaquetesPri, pkg)
		log.Println("Paquete ingresado a cola prioritaria\n")
	
	} else {
		PaquetesRetail = append(PaquetesRetail, pkg)
		log.Println("Paquete ingresado a cola de retail\n")
	}

	//se vuelve al valor original del codigo de seguimiento
	cod_tracking = help_cod_tracking

	if (tipoP != "retail"){
		cod_tracking += 1
	}

	return &logis.CodSeguimiento{Codigo: cod_tracking}, nil
}

/*
func ConectarCamion(){

	listenCamion, err := net.Listen("tcp", ":50055")
	failOnError(err, "error de conexion con camiones")

	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenCamion))
}


func ConectarCliente(){

	listenCliente, err := net.Listen("tcp", ":50051")
	failOnError(err, "error de conexion con cliente")

	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenCliente))
}
*/

func main() {

	log.Println("Server running ...")
	
	//ConectarCliente()
	//go ConectarCamion()
	
	listenCamion, err := net.Listen("tcp", ":50055")
	failOnError(err, "error de conexion con camiones")
	
	listenCliente, err := net.Listen("tcp", ":50051")
	failOnError(err, "error de conexion con cliente")

	log.Fatalln("holi cliente")
	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})
	log.Fatalln("como estas")

	log.Fatalln(srv.Serve(listenCliente))
	log.Fatalln(srv.Serve(listenCamion))
	log.Fatalln("shao")

	//colas rabbitmq con financiero

}