package main

import (
	"context"
	"log"
	"net"
	"time"
	"errors"
	"encoding/json"
	"sync"
	"os"
	"fmt"
	
	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"
	"github.com/streadway/amqp"
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

var mutex = &sync.Mutex{}


func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s\n", msg, err)
	}
}

//escribir en csv
func csvData(reg RegPedido) {

	ap, err := os.OpenFile("registro_logistica.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}
	defer ap.Close()
	if _, err := ap.WriteString(reg.timestamp + "," + reg.id + "," + reg.tipo + "," + reg.nombre + "," + fmt.Sprint(reg.valor) + "," + reg.origen + "," + reg.destino + "," + fmt.Sprint(reg.num_seguimiento+1) + "\n"); err != nil {
		log.Println(err)
	}
}


//
//------------------------ FINANCIERO ------------------------
//

func haciaFinanciero(pak financiero){
	//conn, err := amqp.Dial("amqp://birmania:birmania@10.6.40.157:5672")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "error al conectar")
	defer conn.Close()

	ch,err := conn.Channel()
	failOnError(err, "error al abrir canal")
	defer ch.Close()

	
	q, err := ch.QueueDeclare(
		"hello",
		false,   
		false,   
		false,  
		false,   
		nil,    
	)
	failOnError(err,"error al hacer una cola")

	cosa, err := json.Marshal(pak)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body: []byte(cosa),
		})
		failOnError(err, "no se envio mensaje")

}


func paqueteFinanciero(i int){
	
	var pakfi financiero
	var j = 0

	mutex.Lock()
	pakfi.Id = RegistroSeguimiento[i].id_pkg
	pakfi.Estado = RegistroSeguimiento[i].estado_pkg
	pakfi.Intentos = RegistroSeguimiento[i].num_intentos
	
	lR := len(Registros)
	mutex.Unlock()

	for j = 0; j < lR; j += 1{

		mutex.Lock()
		if (RegistroSeguimiento[i].id_pkg == Registros[j].id){
			
			pakfi.Valor = Registros[j].valor
			pakfi.Tipo = Registros[i].tipo
			mutex.Unlock()

			break
		}
		mutex.Unlock()
	}


	haciaFinanciero(pakfi)
}

//
//-------------------------- CAMIÓN --------------------------
//

//actualización de estado de paquetes al llegar camión a la Central
func (s *server) ResEntrega(ctx context.Context, pkg *logis.RegCamion) (*logis.ACK, error){

	//actualizar estados de paquetes en registro
	var i int
	mutex.Lock()
	lRS := len(RegistroSeguimiento)
	mutex.Unlock()

	for i = 0; i < lRS; i++{
		
		mutex.Lock()
		if (pkg.Id == RegistroSeguimiento[i].id_pkg){

			RegistroSeguimiento[i].estado_pkg = pkg.Estado
			RegistroSeguimiento[i].num_intentos = pkg.Intentos
			mutex.Unlock()
			
			paqueteFinanciero(i)
			
			if pkg.Estado == "rec"{
				return &logis.ACK{Ok: "ok"}, nil
			} else {
				return &logis.ACK{Ok: "ok (pero pedido no entregado jaja)"}, nil
			}

		}
		mutex.Unlock()

	}

	return &logis.ACK{Ok: "error"}, errors.New("No se encontró pedido en registro de seguimiento")
}


func checkColasHelper(tipoCam string, prevRetail bool)(Package, bool){

	if tipoCam == "normal"{

		mutex.Lock()
		lp := len(PaquetesPri)
		ln := len(PaquetesNormal)
		
		if (lp > 0){
			
			pkg := PaquetesPri[0]
			PaquetesPri = PaquetesPri[1:]
			mutex.Unlock()

			log.Println("Cargando paquete prioritario...")
			return pkg, true

		} else if (ln > 0){
			
			pkg := PaquetesNormal[0]
			PaquetesNormal = PaquetesNormal[1:]
			mutex.Unlock()
			
			log.Println("Cargando paquete normal...")
			return pkg, true
		}
		mutex.Unlock()
	
	} else {
		
		mutex.Lock()
		lp := len(PaquetesPri)
		lr := len(PaquetesRetail)

		if (lr > 0){

			pkg := PaquetesRetail[0]
			PaquetesRetail = PaquetesRetail[1:]
			mutex.Unlock()
			
			log.Println("Cargando paquete de retail...")
			return pkg, true
		
		// condicion en caso de haber hecho recien un envio de retail
		} else if (lp > 0) && prevRetail{ 
			
			pkg := PaquetesPri[0]
			PaquetesPri = PaquetesPri[1:]
			mutex.Unlock()
			
			log.Println("Cargando paquete prioritario...")
			return pkg, true
		}
		mutex.Unlock()
	}
	
	return Package{}, false
}


//revisar colas mediante función auxiliar dependiendo del tipo de camión
//que está solicitando y si es la primera o segunda peticióń de este
func CheckColas(tipoCam string, numPeticion int32, prevRetail bool)(Package, bool){

	var pkg Package
	var flagPeticion = false

	if (numPeticion == 2) {
		pkg, flagPeticion = checkColasHelper(tipoCam, prevRetail)

	} else{

		//esperar hasta que llegue algo a prioritario y normal
		for (pkg.id == "") || !flagPeticion{

			log.Println("algun camion", tipoCam, "está en un ciclo, más vale que no salga")
			
			pkg, flagPeticion = checkColasHelper(tipoCam, prevRetail)
			
			time.Sleep(5 * time.Second)
		}

		log.Println("\nPaquete encontrado en cola de tipo", pkg.tipo)
	}

	log.Println("algun camión", tipoCam, "obtuvo", pkg, flagPeticion, "\n")
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
	//deberia quedarse esperando en caso de ser el primer envío y que no haya nada
	pkg, flagPeticion := CheckColas(tipoCam, numPeticion, prevRetail)
	

	//no hay paquetes en colas luego del tiempo de espera del camión
	if ((numPeticion == 2) && !flagPeticion) {
		return &logis.PackageYGeo{}, nil
	}


	//obtener origen y destino del pedido
	var orig string = ""
	var dest string = ""
	
	mutex.Lock()
	for i := range Registros{
		if (Registros[i].id == pkg.id){
			
			orig = Registros[i].origen
			dest = Registros[i].destino
			
			break
		}
	}
	mutex.Unlock()
	
	if (orig == "") || (len(orig) == 0){
		log.Println("origen está vacío, porque de seguro ALGUIEN dejó salir algún camión de un loop supuestamente infinito >:c")
		err := errors.New("No se encontró paquete de cola en registro de pedidos")
		return &logis.PackageYGeo{}, err
	}
	

	//actualizar estado en registro seguimiento
	mutex.Lock()
	for i := range RegistroSeguimiento{
		if (RegistroSeguimiento[i].id_pkg == pkg.id){
			
			RegistroSeguimiento[i].estado_pkg = "tr"
			RegistroSeguimiento[i].id_camion = idCam
			
			break
		}
	}
	mutex.Unlock()
	
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
	
	mutex.Lock()
	for _, regseg := range RegistroSeguimiento{
		if (cs.GetCodigo() == regseg.num_seguimiento){
			mutex.Unlock()
			return &logis.EstadoPedido{Estado: regseg.estado_pkg}, nil			
		}
	}

	mutex.Unlock()
	return &logis.EstadoPedido{Estado: "nulo"}, errors.New("No se encontró el paquete pedido en registro.")
}


func (s *server) PedidoCliente(ctx context.Context, pedido *logis.Pedido) (*logis.CodSeguimiento, error) {
	
	log.Print("Pedido recibido")

	//auxiliar para guardar valor de cod_tracking
	mutex.Lock()
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
	reg := RegPedido{timestamp:nownow.Format("2006-01-02 15:04:05"), 
					id:pId, 
					tipo: tipoP, 
					nombre: pedido.Producto, 
					valor: pVal, 
					origen: pedido.Tienda, 
					destino: pedido.Destino, 
					num_seguimiento: cod_tracking}

	Registros = append(Registros, reg)

	//guardar en registro csv
	csvData(reg)

	//guardar en registro de seguimiento
	RegistroSeguimiento = append(RegistroSeguimiento, PackageSeguimiento{id_pkg: pId,
														estado_pkg: "bdg",
														id_camion: "",
														num_seguimiento: cod_tracking,
														num_intentos: 0})
	
	//transformar pedido en paquete
	pkg := Package{id: pId,
					num_seguimiento: cod_tracking,
					tipo: tipoP,
					valor: pVal,
					num_intentos: 0,
					estado: "bdg"}

	mutex.Unlock()
					
	
	//guardar paquete en cola
	if (tipoP == "normal") {

		mutex.Lock()
		PaquetesNormal = append(PaquetesNormal, pkg)
		
		log.Println("Paquete ingresado a cola normal\n")
		
		log.Println("largo PaquetesNormal:", len(PaquetesNormal))
		mutex.Unlock()


	} else if (tipoP == "prioritario") {

		mutex.Lock()
		PaquetesPri = append(PaquetesPri, pkg)
		
		log.Println("Paquete ingresado a cola prioritaria\n")
		
		log.Println("largo PaquetesPri:", len(PaquetesPri))
		mutex.Unlock()
		
	} else {

		mutex.Lock()
		PaquetesRetail = append(PaquetesRetail, pkg)
		
		log.Println("Paquete ingresado a cola de retail\n")
		
		log.Println("largo PaquetesRetail:", len(PaquetesRetail))
		mutex.Unlock()
	}

	//se vuelve al valor original del codigo de seguimiento
	mutex.Lock()
	cod_tracking = help_cod_tracking
	mutex.Unlock()

	if (tipoP != "retail"){
		mutex.Lock()
		cod_tracking += 1
		mutex.Unlock()
	}

	return &logis.CodSeguimiento{Codigo: cod_tracking}, nil
}


//
//------------------------ CONEXIONES ------------------------
//
func ConectarCliente(){

	listenCliente, err := net.Listen("tcp", "10.6.40.157:50051")
	//listenCliente, err := net.Listen("tcp", "localhost:50051")
	failOnError(err, "error de conexion con cliente")

	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenCliente))
}


func ConectarCamion(){

	listenCamion, err := net.Listen("tcp", "10.6.40.157:50055")
	//listenCamion, err := net.Listen("tcp", "localhost:50055")
	failOnError(err, "error de conexion con camiones")

	srv := grpc.NewServer()
	logis.RegisterLogisServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(listenCamion))
}


func main() {

	log.Println("Server corriendo...")
	
	go ConectarCamion()
	ConectarCliente()

}
