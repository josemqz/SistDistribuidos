package camion

import (
	"context"
	"log"
	"time"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"

	//"os"
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

var pkg_time int32

func main(){
	
	//pedir tiempo de espera entre pedidos por input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresar tiempo de espera entre pedidos: ")
	pkg_time, err := strconv.ParseInt(reader.ReadString('\n'), 10, 32)

	for (err != nil) {
		fmt.Print("Ingresar tiempo de espera entre pedidos: ")
		pkg_time, err := strconv.ParseInt(reader.ReadString('\n'), 10, 32)
	}

	//conexión
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	//instanciaciones
	go initCamionPyme("CN", conn)
	go initCamionRetail("CR1", conn)
	go initCamionRetail("CR2", conn)

}

func initCamion(idCam string, conn ClientConn /*o LogisServiceClientConn*/) (){

	camion := logis.NewLogisServiceClient(conn)

	pkg1, geo1 := GuardarPedido(camion, idCam, true)

	//time.wait(pkg_time*time.Second)
	time.Sleep(time.Duration(pkg_time) * time.Second)

	//pedir a logis paquete
	pkg2, geo2 := GuardarPedido(camion, idCam, false)

	Delivery(pkg1,geo1,pkg2,geo2)

	//enviar estado de paquetes (no) enviados
	camion.EstadoCamion()
	//if dos paquetes
	camion.EstadoCamion()

}

//se pide un paquete y se guarda en el registro. intento representa si es
//la primera o segunda vez que se pide, para saber si logístico debe esperar
func GuardarPedido(camion logis.LogisServiceClient, idCam string, intento bool) (*logis.Package, *logis.Geo){

	if idCam == "CN"{
		tipoCam := "normal"
	} else{
		tipoCam := "retail"
	}

	//pedir a logis paquete
	pkg, geo := camion.PedidoACamion(&logis.tipoCam{tipo: tipoCam}, &logis.mIntento{intento: intento})

	if (intento == false && pkg/*pkg == nil?*/){
		return nil,nil
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

	return (pkg, geo)
}


func Delivery(){
	
	//una vez listo para salir a hacer entregas, esperar un tiempo (puede ser el mismo 
	//independientemente del destino, o variar según este)
	
}

	//reintento: 10 dp
		//pyme: si <= precio producto + [30% en caso de prioritario] || n_intento <= 3
		//retail: n_intento <= 3

	//si paquetes no son entregados:
		//normal: nada
		//prioritario: 30%
		//retail: precio producto

