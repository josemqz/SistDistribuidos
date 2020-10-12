//logistica
package main

import (
	
	"fmt"
	"os"
	//"time"
	"log"
	
	"encoding/csv"

	//"bufio"
	"io"
)

	
type paquete struct {
	id int
	seguimiento int
	tipo string
	valor int
	intentos int
	estado string
}


//uso: p:= paquete{id: 1, seguimiento: 12001,...} o simplemente paquete{1, 12001,...} 
//p.valor = 200
//return &p

func fileExists(arch string) bool {
    
    info, err := os.Stat(arch)
    
    if os.IsNotExist(err) {
        return false
    }
    else return true
}

/*
func enqueue(cola []paquete, p paquete) paquete{
	cola = append(cola, p)
	return cola
}

func dequeue(cola []paquete) paquete{
	return cola[1:]
}
*/

func crearRegistro(nombre string) *os.File{ //quizás haya que poner (arch *os.File)
	
	//ya existe el registro
	if fileExists(nombre){

		arch, err := os.Open(nombre)

	    if err != nil {
	    	log.Fatalln("No se pudo abrir el registro" + nombre, csv) //chequear si está bien
	        return
    	}
    }

    //crear registro
    else{

    	arch, err := os.Create(nombre)
	    
	    if err != nil {
        	log.Fatal("No se pudo crear archivo de registro", err)
    	}
    	
    	defer arch.Close()
    }

    return arch
}

func main(){

	var colaR []paquete
	var colaP []paquete
	var colaN []paquete

/*
	//registro de pedidos
	reg := crearRegistro("registro_logistica.csv")
    regWriter := csv.NewWriter(reg)
    defer regWriter.Flush()
*/  
    //escribir cabeceras
	//timestamp, id del paquete, tipo de paquete, nombre de producto*, 
	//valor de producto*, origen (tienda?), destino*, número de seguimiento


    //registro de seguimiento
    seg := crearRegistro("registro_seguimiento.csv")
    segWriter := csv.NewWriter(reg)
    defer segWriter.Flush()

    //escribir cabeceras
    //id del paquete, estado del paquete, id del camión, id de seguimiento, cantidad de intentos


	//for{

		/*gRPC clientes
		if llega mensaje:

			if (tienda == "pyme"):
				generar codigo de seguimiento

			timestamp :=

			escribir en registro (timestamp, id del paquete, tipo de paquete, nombre de producto*, 
									valor de producto*, origen (tienda?), destino*, número de seguimiento)
			
		retail
			id	producto	valor	tienda	destino
		pyme
			id	producto	valor	tienda	destino prioritario


	        err := writer.Write(registro)
		    if err != nil {
	        	log.Fatal("No se pudo escribir en registro", err)
	    	}
			
			enviar codigo de seguimiento a cliente

			//hay que usar un indice para saber en qué producto va (para no repetirlos)

			generar EDD de paquete para enviar a camión 

		*/
		
		/*gRPC camiones

		enviar id de seguimiento
		recibir estado (En bodega, En camino, Recibido o No Recibido)
		actualizar registro de seguimiento (seg)
		enviar a cliente que pidió seguimiento

		si un camión vuelve, hay que entregar datos de ganancias/pérdidas a financiero

		*/
	//}

	//colas rabbitmq con financiero

}