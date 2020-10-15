//camiones
package main

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

/*
type Package struct{
	id string
	num_seguimiento int
	tipo string
	valor int
	num_intentos int
	estado string
	// estados-> bdg: bodega, tr: en tránsito , rec: recibido , r: no recibido
}

id-paquete, tipo de paquete, valor, origen, destino, número de intentos
y fecha de entrega
*/

func main(){


	//pedir tiempo de espera entre pedidos por input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresar tiempo de espera para pedidos: ")
	usr_time, _ := reader.ReadString('\n')
	

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	//defer conn.Close()

	camion := logis.NewLogisServiceClient(conn)

	
	func (s *server)GenerarRegistroCam(ctx context.Context, //cs CodSeguimiento) return (*logis.EstadoPedido, error) {

//		for _, regseg := range PaquetesRetail{
//			RegCamion := append(RegCamion, regseg // + origen + destino + fecha entrega - seguimiento - estado)

//		}	
//	}
	


	//reintento: 10 dp
		//pyme: si <= precio producto + [30% en caso de prioritario] || n_intento > 2
		//retail: n_intento <= 3

	//si paquetes no son entregados:
		//normal: nada
		//prioritario: 30%
		//retail: precio producto


	/*
	3 camiones: 2 retail + normal

	while not pedido -> no hacer nada


	if llega un pedido
		agregar
		if pedido
			agregar
		else
			esperar tiempo_input
			if llega un pedido
				agregar
			

	//normal
		if prioritario
			agregar
			if prioritario
				agregar
			else if normal
				agregar
		else if normal
			agregar
			

	//retail
		if retail
			agregar
			if retail
				agregar
			else if proritario
				agregar


	*/

	
}