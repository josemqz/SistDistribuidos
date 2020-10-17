package Camiones

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
	pkg_time, err := reader.ReadString('\n')

	for (err != nil) {
		fmt.Print("Ingresar tiempo de espera entre pedidos: ")
		pkg_time, err := reader.ReadString('\n')	
	}

	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	camion := logis.NewLogisServiceClient(conn)

	//tiene que escuchar a logis hasta que le llegue un request ???

}
	
//se llama a la función cada vez que llega un paquete
func PedidoACamion(pkg *logis.Package /*debería ser stream?*/, geo *logis.Geo){

	//server envía un mensaje con el paquete a camión 
	//se guarda en registro
	RegPackage{
		id_pkg: pkg.Id,
		tipo: pkg.Tipo,
		valor: pkg.Valor,
		origen: geo.Origen,
		destino: geo.Destino,
		num_intentos: pkg.Num_intentos,
	}

	//espera un tiempo, y si al terminar le llega otro paquete, lo acepta
	Delivery()
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