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

)



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

func obtenerPaquete(coss logis.LogisService, tipoCam string, idCam string) (bool, paquete) {
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
	

}








//__________________________

func main(){

	//pedir tiempo de espera entre pedidos por input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresar tiempo de espera para pedidos: ")
	usr_time, _ := reader.ReadString('\n')
	

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	camion := logis.NewLogisServiceClient(conn)

	//tiene que escuchar a logis hasta que le llegue un request

}
	
func GenerarRegistroCam(ctx context.Context, cs CodSeguimiento) {

//		for _, regseg := range PaquetesRetail{
//			RegCamion := append(RegCamion, regseg // + origen + destino + fecha entrega - seguimiento - estado)

//		}	
//	}
	
func Delivery(ctx context.Context, )

	//reintento: 10 dp
		//pyme: si <= precio producto + [30% en caso de prioritario] || n_intento > 2
		//retail: n_intento <= 3

	//si paquetes no son entregados:
		//normal: nada
		//prioritario: 30%
		//retail: precio producto


	/*


	*/

	
}