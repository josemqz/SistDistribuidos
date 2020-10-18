//camiones
package main

import (

	"context"
	"log"
	"time"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"

	"math"
	"math/rand"
	"strconv"
	"bufio"

)


var tpo1 int //tiempo de espera por segundo paquete - input usuario en main
var tpoEnvio int //tpo de demora de cada envio - input usuario en main

type paquete struct {
	id string
	tipo string
	valor int32
	origen string
	destino string
	intentos int32
	fechaEntrega string
	tipoCam string
	estado string
}

var todos []paquete
var retail1 []paquete
var retail2 []paquete
var normal []paquete

func obtenerPaquete(coss logis.LogisServiceClient, tipoCam string, idCam string) (bool, paquete) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := coss.SolPaquetes(ctx, &logis.TipoCam{Tipo: tipoCam})
	if err != nil {
		log.Fatalf("No se pudo obtener paquete: %v", err)
	}
	var newPak paquete
	if r.GetId() != "null" {
		log.Printf("Paquete recibido")
		newPak.id = r.GetId()
		newPak.tipo = r.GetTipo()
		newPak.valor = r.GetValor()
		newPak.origen = r.GetOrigen()
		newPak.destino = r.GetDestino()
		newPak.intentos = r.GetIntentos()
		newPak.estado = "Procesando"
		newPak.fechaEntrega = ""

		if idCam == "Ret1" {
			newPak.tipoCam = idCam
			retail1 = append(retail, newPak)
		}
		if idCam == "Ret2" {
			newPak.tipoCam = idCam
			retail2 = append(retail2, newPak)
		}
		if idCam == "norm" {
			newPak.tipoCam = idCam
			normal = append(normal, newPak)
		}
		todos = append(todos, newPak)
		return true, newPak
	}
	return false, newPak
}



func siRecibe(xIntentos int) int{
	
	i := 1
	for i < (xIntentos+1){
		time.Sleep(time.Duration(tpoEnvio) * time.Millisecond)
		n := rand.Intn(100)
		if n < 80 {
			return i
		} i++
	}
}


func actualizaP(id string, estado string, fecha string, intento int32) {
	i := 0
	for i < len(todos){
		if todos[i].id == id{
			var aux = todos[i]
			aux.estado = estado
			aux.intentos = intento
			aux.fechaEntrega = fecha
			todos[i] = aux
			return
		}
		i++
	}
}


func enviarEstado(pak paquete){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, _ := logg.UpdateEstado(ctx, &logis.EstadoPedido{Id: pak.id, Estado: pak.estado})
	log.Printf("mensaje: %v", r.GetMessage())

}



func ePyme(pak paquete) {
	
	// definir maximo de intentos 
	max := math.Floor((pak.valor) / 10)	
	
	pak.estado = "tr" //en transito
	enviarEstado(pak)
	var intentos = siRecibe(int(max))
	t := time.Now()
	pak.fechaEntrega = t.Format("2006-01-02 15:04:05")
	pak.intentos = int32(intentos)
	
	if intentos == int(max) {
		
		//ya se intento el maximo de veces y no fue recibido

		pak.estado = "nr" //no recibido
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := logg.ResEntrega(ctx, &logis.PaqRecibido{Id: pak.id, Intentos: pak.intentos, Estado: pak.estado, Tipo: pak.tipo})
		log.Printf("Mje: %v", r.GetMessage())

	} else {
		
		// exito
		
		pak.estado = "rec" //recibido
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := logg.ResEntrega(ctx, &logis.PaqRecibido{Id: pak.id, Intentos: pak.intentos, Estado: pak.estado, Tipo: pak.tipo})
		log.Printf("Mje: %v", r.GetMessage())
	}
	actualizaP(pak.id, pak.estado, pak.fechaEntrega, pak.intentos)
}
	
	
func eRetail(pak paquete) {
	
	max := 3

	pak.estado = "tr" //en transito
	enviarEstado(pak)

	var intentos = siRecibe(max)
	t := time.Now()
	pak.fechaEntrega = t.Format("2006-01-02 15:04:05")
	pak.intentos = int32(intentos)
	
	if intentos == 3 {
		
		//ya se intento el maximo de veces y no fue recibido

		pak.estado = "nr" //no recibido
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := logg.ResEntrega(ctx, &logis.PaqRecibido{Id: pak.id, Intentos: pak.intentos, Estado: pak.estado, Tipo: pak.tipo})
		log.Printf("Mje: %v", r.GetMessage())

	} else {
		
		// exito
		
		pak.estado = "rec" //recibido
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := logg.ResEntrega(ctx, &logis.PaqRecibido{Id: pak.id, Intentos: pak.intentos, Estado: pak.estado, Tipo: pak.tipo})
		log.Printf("Mje: %v", r.GetMessage())
	}
	actualizaP(pak.id, pak.estado, pak.fechaEntrega, pak.intentos)
}


//var logg logis.LogisServiceClient - maybe

func tipoDespacho(pak paquete){
	if pak.origen == "pyme" {
		ePyme(pak)
	}else {
		eRetail(pak)
	}

}

func entrega(pak1 paquete, pak2 paquete){

	if pak2 != nil {
		if (pak1.valor >= pak2.valor) {
			tipoDespacho(pak1)
			tipoDespacho(pak2)
		}else {
			tipoDespacho(pak2)
			tipoDespacho(pak1)
		}
	}else{
		tipoDespacho(pak1)
	}
}


func delivery(lc logis.LogisServiceClient, tipoCam string, idCam string, /*WaitGroup??*/) bool {

	var rs1, primero = obtenerPaquete(lc, tipoCam, idCam)
	
	if rs1 == true {
		var rs2, segundo = obtenerPaquete(lc, tipoCam, idCam)
		// si no hay aun segundo paquete
		if !rs2 {
			//espera el tpo definido por usuario por otro paquete
			time.sleep(tpo1)
			var otro = obtenerPaquete(lc, tipoCam, idCam)
			entrega(primero, otro) //la funcion entrega maneja si existe o es null el segundo
		}else{
			entrega(primero, segundo)
		}
	}
	//waitgroup?	
	return true
}




//




//___________________________________________________

func main(){

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("\ndid not connect with server: %v", err)
	}
	defer conn.Close()
	cl := logis.NewLogisServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese tiempo espera segundo pedido\n")
	texto, _ := reader.ReadString('\n')
	tpo1, _ = strconv.Atoi(strings.TrimSuffix(texto, "\n"))
	fmt.Print("Ingrese tiempo de envio de un pedido\n")
	reader = bufio.NewReader(os.Stdin)
	texto, _ = reader.ReadString('\n')
	tpoEnvio, _ = strconv.Atoi(strings.TrimSuffix(texto, "\n"))

	
	go delivery(lc logis.LogisServiceClient, "retail" , "Ret1" )
	go delivery(lc logis.LogisServiceClient, "retail" , "Ret2")
	go delivery(lc logis.LogisServiceClient, "normal" , "norm")	




	//tiene que escuchar a logis hasta que le llegue un request

	
}