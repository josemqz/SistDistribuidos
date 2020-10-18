package camion

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"
	"strconv"

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

//si camiones están en la central
var CentralCR1 = true
var CentralCR2 = true
var CentralCN = true


func main(){
	
	//pedir tiempo de espera entre pedidos por input
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Ingresar tiempo de espera entre pedidos de camiones: ")

	read, err := reader.ReadString('\n')
	pkg_time := strconv.Atoi(read)

	for (err != nil) {
		fmt.Print("Ingresar tiempo de espera entre pedidos de camiones: ")

		read, err := reader.ReadString('\n')
		pkg_time := strconv.Atoi(read)
	}

	//tiempo envíos
	fmt.Print("Ingresar tiempo de viaje para envíos: ")

	read, err := reader.ReadString('\n')
	dlvr_time := strconv.Atoi(read)

	for ((err != nil) || dlvr_time <= 0) {

		fmt.Print("Ingresar tiempo de viaje para envíos: ")

		read, err := reader.ReadString('\n')
		dlvr_time := strconv.Atoi(read)
	}


	//conexión
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if (err != nil) {
		log.Fatalln(err)
	}
	defer conn.Close()


	//instanciaciones
	for {
		if CentralCR1{
			CentralCR1 = false
			go initCamion("CR1", conn)
		}
		if CentralCR2{
			CentralCR2 = false
			go initCamion("CR2", conn)
		}
		if CentralCN{
			CentralCN = false
			go initCamion("CN", conn)
		}
	}

}

func initCamion(idCam string, conn ClientConn /*o LogisServiceClientConn*/) (){

	camion := logis.NewLogisServiceClient(conn)

	reg1 := GuardarPedido(camion, idCam, true)

	//time.wait(pkg_time*time.Second)
	time.Sleep(time.Duration(pkg_time) * time.Second)

	//pedir a logis paquete
	reg2 := GuardarPedido(camion, idCam, false)

	//envíos
	Delivery(reg1, reg2)

	//enviar estado de paquetes (no) enviados
	camion.EstadoCamion()
	if idCam == "CR1"{
		CentralCR1 = true
	}
	else if idCam == "CR2"{
		CentralCR2 = true
	}
	else{
		CentralCN = true
	}

	//if dos paquetes | if reg2 != nil (?)
	camion.EstadoCamion()
	if idCam == "CR1"{
		CentralCR1 = true
	}
	else if idCam == "CR2"{
		CentralCR2 = true
	}
	else{
		CentralCN = true
	}

	//hacer funcion ^ ^ ^

}


//se pide un paquete y se guarda en el registro. numPeticion representa si es
//la primera (true) o segunda vez (false) que se pide, para saber si logístico debe esperar
func GuardarPedido(camion logis.LogisServiceClient, idCam string, numPeticion bool) (RegPackage){

	//tipo de camión
	if idCam == "CN"{
		tipoCam := "normal"
	} else{
		tipoCam := "retail"
	}

	//pedir a logis paquete
	pkg, geo := camion.PedidoACamion(&logis.tipoCam{Tipo: tipoCam}, numPeticion)

	//si segunda petición de paquete no fue exitosa se retorna nil y ya uwu
	if (numPeticion == false && pkg == nil){
		return nil
	}
	
	//guardar en registro
	reg := RegPackage{id_pkg: pkg.Id,
						tipo: pkg.Tipo,
						valor: pkg.Valor,
						origen: geo.Origen,
						destino: geo.Destino,
						num_intentos: pkg.Num_intentos,
	}

	if idCam == "CR1" {
		RegistroCR1 = append(RegistroCR1, reg)
	} else if idCam == "CR2" {
		RegistroCR2 = append(RegistroCR2, reg)
	} else if idCam == "CN" {
		RegistroCN = append(RegistroCN, reg)
	}

	return reg
}


func Delivery(reg1, reg2){
	
	//una vez listo para salir a hacer entregas, esperar un tiempo (puede ser el mismo 
	//independientemente del destino, o variar según este)
	
	//if reg2 != nil
		//comparar ingresos de cada paquete
	
	//envio(paqueteActual) -> intenta entregar
	//if not entregado; estado paquete 2 = false
	
	//if reg2 != nil
		//envio(paqueteActual)

	
}

	//reintento: 10 dp
		//pyme: si <= precio producto + [30% en caso de prioritario] || n_intento <= 3
		//retail: n_intento <= 3

	//si paquetes no son entregados:
		//normal: nada
		//prioritario: 30%
		//retail: precio producto

