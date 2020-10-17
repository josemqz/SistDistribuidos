package camion

import (
	"context"
	"log"
	"time"

	"github.com/josemqz/SistDistribuidos/Lab1/logis"
	"google.golang.org/grpc"

	
	//"os"
	//"encoding/csv"
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

	go initCamionPyme(pkg_time, conn)
	go initCamionRetail(pkg_time, conn)
	go initCamionRetail(pkg_time, conn)

}

func initCamionPyme(pkg_time int32, conn ){

	var RegistroCam []RegPackage
	var tipoCam = "pyme"
	
	camion := logis.NewLogisServiceClient(conn)

	//pedir a logis paquete
	//se guarda en registro
	pkg, geo := PedidoACamion(&logis.tipoCam{tipo: tipoCam})
	
	ResgistroCam =  append(ResgistroCam, RegPackage{id_pkg: pkg.Id,
										tipo: pkg.Tipo,
										valor: pkg.Valor,
										origen: geo.Origen,
										destino: geo.Destino,
										num_intentos: pkg.Num_intentos,
										}
	//time.wait(pkg_time*time.Second)
	//pedir a logis paquete

}

func initCamionRetail(){

}



func Delivery(){
	
	//una vez listo para salir a hacer entregas, esperar un tiempo (puede ser el mismo 
	//independientemente del destino, o variar según este)

}

	//reintento: 10 dp
		//pyme: si <= precio producto + [30% en caso de prioritario] || n_intento > 2
		//retail: n_intento <= 3

	//si paquetes no son entregados:
		//normal: nada
		//prioritario: 30%
		//retail: precio producto

}