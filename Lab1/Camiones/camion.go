package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"math"
	"math/rand"
	"sync"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"

	//"bufio"
	//"io"
)

//preguntar si en memoria o csv
type RegPackage struct{
	id_pkg string
	tipo string
	valor int32
	origen string
	destino string
	num_intentos int32
	fecha_entrega string
}

var RegistroCN []RegPackage
var RegistroCR1 []RegPackage
var RegistroCR2 []RegPackage

var pkg_time int
var dlvr_time int

//si última entrega tuvo algún paquete de retail
var prevRetail1 = false
var prevRetail2 = false

//si camiones están en la central
var CentralCR1 = true
var CentralCR2 = true
var CentralCN = true

var mutex = &sync.Mutex{}


func failOnError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s: %s", msg, err)
	}
}


//Actualizar registro de camión
func actualizarReg(idCam string, id string, estado string, fecha string, intento int32) {

	var i int

	if idCam == "CR1"{

		mutex.Lock()
		CentralCR1 = true
		mutex.Unlock()

		for i = 0; i < len(RegistroCR1); i++{
			if RegistroCR1[i].id_pkg == id{

				mutex.Lock()
				RegistroCR1[i].num_intentos = intento
				RegistroCR1[i].fecha_entrega = fecha
				mutex.Unlock()

				return
			}
		}

	} else if idCam == "CR2"{
		
		mutex.Lock()
		CentralCR2 = true
		mutex.Unlock()

		for i = 0; i < len(RegistroCR2); i++{
			if RegistroCR2[i].id_pkg == id{

				mutex.Lock()
				RegistroCR2[i].num_intentos = intento
				RegistroCR2[i].fecha_entrega = fecha
				mutex.Unlock()

				return
			}
		}

	} else if idCam == "CN"{
		
		mutex.Lock()
		CentralCN = true
		mutex.Unlock()

		for i = 0; i < len(RegistroCN); i++{
			if RegistroCN[i].id_pkg == id{

				mutex.Lock()
				RegistroCN[i].num_intentos = intento
				RegistroCN[i].fecha_entrega = fecha
				mutex.Unlock()
				
				return
			}
		}
	}
}


//Enviar estado actualizado a logística
func EstadoCamion(idCam string, exito bool, reg RegPackage, camion logis.LogisServiceClient){

	var estado string
	if exito{
		estado = "rec"
	} else {
		estado = "nr"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1000 * time.Second)
	defer cancel()

	r, err := camion.ResEntrega(ctx, &logis.RegCamion{Id: reg.id_pkg, Intentos: reg.num_intentos, Estado: estado})
	failOnError(err, "Error en comunicación de actualización de pedido desde Camión")

	log.Printf("ok camion->logistico: %v", r.GetOk())

	actualizarReg(idCam, reg.id_pkg, estado, reg.fecha_entrega, reg.num_intentos)
}


//probabilidad de éxito en entrega
func intentar() bool{

	n := rand.Intn(100)

	if n < 80 {
		return true
	} else {
		return false
	}
}


//definir maximo de intentos
func maxIntentos(reg RegPackage) int{

	if (reg.tipo == "pyme"){
		return int(math.Floor(float64((reg.valor) / 10)) + 1)

	} else {
		return 3
	}
}


//regA > regB
func Delivery(idCam string, regA RegPackage, regB RegPackage, camion logis.LogisServiceClient){

	//entrega exitosa o no
	var estadoA = false
	var estadoB = false

	//cantidad de intentos
	var tryA = 1
	var tryB = 1

	//cantidad máxima de intentos
	var maxA = maxIntentos(regA)
	var maxB = maxIntentos(regB)


	//solo una entrega
	if (regB.id_pkg == ""){
		for (tryA <= maxA) && !estadoA{
			
			time.Sleep(time.Duration(dlvr_time) * time.Second)

			if (intentar()){

				t := time.Now()
				regA.fecha_entrega = t.Format("2006-01-02 15:04:05")
				regA.num_intentos = int32(tryA)
				
				EstadoCamion(idCam, true, regA, camion)
				estadoA = true

			} else{
				tryA += 1
			}
		}

		//se pasó de maxA y no fue entregado
		if (tryA > maxA){
			
			regA.fecha_entrega = "0"
			regA.num_intentos = int32(tryA)

			EstadoCamion(idCam, false, regA, camion)
		}


	//dos entregas
	} else {
		
		for (!estadoA && tryA <= maxA) && (!estadoB && tryB <= maxB){	//CHECK
			
			if !estadoA && (tryA <= maxA){

				time.Sleep(time.Duration(dlvr_time) * time.Second)

				if (intentar()){
					
					t := time.Now()
					regA.fecha_entrega = t.Format("2006-01-02 15:04:05")
					regA.num_intentos = int32(tryA)
					
					EstadoCamion(idCam, true, regA, camion)
					estadoA = true
	
				} else{
					tryA += 1
				}
			}

			if (tryA > maxA){
				
				regA.fecha_entrega = "0"
				regA.num_intentos = int32(tryA)
				
				EstadoCamion(idCam, false, regA, camion)
			}
			

			if !estadoB && (tryB <= maxB){

				time.Sleep(time.Duration(dlvr_time) * time.Second)
				
				if intentar(){
					
					t := time.Now()
					regB.fecha_entrega = t.Format("2006-01-02 15:04:05")
					regB.num_intentos = int32(tryB)
					
					EstadoCamion(idCam, true, regB, camion)
					estadoB = true
	
				} else{
					tryB += 1
				}
			}

			if (tryB > maxB){

				regB.fecha_entrega = "0"
				regB.num_intentos = int32(tryB)

				EstadoCamion(idCam, false, regB, camion)
			}
		}
	}
}


//orden de pedidos
func OrdenP(idCam string, reg1 RegPackage, reg2 RegPackage, camion logis.LogisServiceClient){

	//si hay un segundo paquete y su valor es mayor
	if (reg2.id_pkg != "") && (reg1.valor < reg2.valor) {
		Delivery(idCam, reg2, reg1, camion)

	} else {
		Delivery(idCam, reg1, reg2, camion)
	}
}


//reintento: 10 dp
	//pyme: si <= precio producto + [30% en caso de prioritario] || n_intento <= 3
	//retail: n_intento <= 3

//si paquetes no son entregados:
	//normal: nada
	//prioritario: 30%
	//retail: precio producto



//se pide un paquete y se guarda en el registro. numPeticion representa si es
//la primera (true) o segunda vez (false) que se pide, para saber si logístico debe esperar
func RegistrarPedido(camion logis.LogisServiceClient, idCam string, numPeticion bool) (RegPackage){
	log.Println("funcion: registrarpedido\n")

	var tipoCam string
	var pedido *logis.PackageYGeo
	var err error

	//tipo de camión
	if idCam == "CN"{
		tipoCam = "normal"
	} else{
		tipoCam = "retail"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1000 * time.Second)
	defer cancel()

	//pedir a logística un paquete
	if idCam == "CR1" {
		pedido, err = camion.PedidoACamion(ctx, &logis.InfoCam{Tipo: tipoCam, NumPeticion: numPeticion, PrevRetail: prevRetail1})
		failOnError(err, "")

	} else if idCam == "CR2"{
		pedido, err = camion.PedidoACamion(ctx, &logis.InfoCam{Tipo: tipoCam, NumPeticion: numPeticion, PrevRetail: prevRetail2})
		failOnError(err, "")

	} else if idCam == "CN" {
		pedido, err = camion.PedidoACamion(ctx, &logis.InfoCam{Tipo: tipoCam, NumPeticion: numPeticion, PrevRetail: false})
		failOnError(err, "")
	}


	//si segunda petición de paquete no fue exitosa se retorna nil y ya uwu
	if (!numPeticion && pedido == &logis.PackageYGeo{}){
		log.Println("No hay más paquetes. Iniciando viaje con uno solo...")
		return RegPackage{}
	}
	
	//guardar en registro
	reg := RegPackage{id_pkg: pedido.Id,
						tipo: pedido.Tipo,
						valor: pedido.Valor,
						origen: pedido.Origen,
						destino: pedido.Destino,
						num_intentos: pedido.NumIntentos,
	}

	if idCam == "CR1" {
		mutex.Lock()
		RegistroCR1 = append(RegistroCR1, reg)
		mutex.Unlock()
		log.Println("Pedido agregado a registro")
	
	} else if idCam == "CR2"{
		mutex.Lock()
		RegistroCR2 = append(RegistroCR2, reg)
		mutex.Unlock()
		log.Println("Pedido agregado a registro")

	} else if idCam == "CN" {
		mutex.Lock()
		RegistroCN = append(RegistroCN, reg)
		mutex.Unlock()
		log.Println("Pedido agregado a registro")
	}

	return reg
}


func initCamion(idCam string, camion logis.LogisServiceClient) (){

	log.Println("Camión iniciado\n")

	reg1 := RegistrarPedido(camion, idCam, true)
	log.Println("Paquete 1 cargado en ", idCam, "\n")
	
	time.Sleep(time.Duration(pkg_time) * time.Second)
	
	reg2 := RegistrarPedido(camion, idCam, false)
	if (reg2.id_pkg != ""){
		log.Println("Paquete 2 cargado en ", idCam, "\n")
	}

	//si se envió retail la última vez
	if (reg1.tipo == "retail") || (reg2.tipo != "retail"){
		
		if idCam == "CR1"{
			prevRetail1 = true
		} else if idCam == "CR2"{
			prevRetail1 = true
		}
	}

	//envíos
	OrdenP(idCam, reg1, reg2, camion)

}


func main(){
	
	//pedir tiempo de espera entre pedidos por input
	fmt.Print("Ingresar tiempo de espera entre pedidos de camiones: ")

	_, err := fmt.Scanf("%d", &pkg_time)

	for (err != nil){

		log.Println("Tiempo ingresado inválido!\n")
		log.Print("Ingresar tiempo de espera entre pedidos: ")
		
		_, err = fmt.Scanf("%d", &pkg_time)
	}


	//pedir tiempo de envíos
	fmt.Print("Ingresar tiempo de viaje para envíos: ")

	_, err = fmt.Scanf("%d", &dlvr_time)

	for (err != nil) || (dlvr_time <= 0) {

		log.Println("Tiempo ingresado inválido!")
		fmt.Print("Ingresar tiempo de viaje para envíos: ")

		_, err = fmt.Scanf("%d", &dlvr_time)
	}


	//conexión
	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure(), grpc.WithBlock())
	failOnError(err,"Error en conexión a servidor")
	defer conn.Close()

	camion := logis.NewLogisServiceClient(conn)
	log.Println("Conexión realizada\n")

	//instanciaciones
	for {
		mutex.Lock()
		if CentralCR1{
			CentralCR1 = false
			go initCamion("CR1", camion)
		}
		mutex.Unlock()
		mutex.Lock()
		if CentralCR2{
			CentralCR2 = false
			go initCamion("CR2", camion)
		}
		mutex.Unlock()
		
		mutex.Lock()
		if CentralCN{
			CentralCN = false
			go initCamion("CN", camion)
		}
		mutex.Unlock()
	}

}
