//financiero
package main

import (
	
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	"fmt"

	//"time"
	//"bufio"
	//"io"
	//"os"
)

/*
- Los envı́os completados -> se guardan en completados[]
- La cantidad de veces que se intentó entregar un paquete -> paquete.intentos
- Los paquetes que no pudieron ser entregados -> se guardan en norecib[]
- Pérdidas o ganancias de cada paquete (en Dignipesos)-> se calcula en contador(paquete)
*/

//Además mostrará el balance final en dignipesos durante la ejecución y al finalizar.




type EnviosFinanzas struct {
	Id_paquete string `json:"id"`
	Tipo string `json:"tipo"`
	Valor int32 `json:"valor"`
	Intentos int32 `json:"intentos"`
	Estado string `json:"estado"`
	Balance float64

}


var ganancias float64
var perdidas float64
var balance_final float64
var balance_producto float64
var balance_t float64
var ganancias_f float64
var perdidas_f float64

var completados []EnviosFinanzas
var norecib []EnviosFinanzas


//calcula perdidas y ganancias segun tipo de envío
func contador(pak EnviosFinanzas){

	balance_producto = 0
	if(pak.Tipo == "retail"){
		balance_producto = float64(pak.Valor) - float64(10*(pak.Intentos-1)) 
		ganancias = float64(pak.Valor)
		perdidas = float64(10*(pak.Intentos-1))
	}else if(pak.Tipo == "prioritario"){
		if(pak.Estado == "rec"){
			balance_producto = float64(pak.Valor - (10*(pak.Intentos-1)))
			ganancias = float64(pak.Valor)
			perdidas = float64(10*(pak.Intentos-1))
		}else{
			balance_producto = float64(pak.Valor - (10*(pak.Intentos-1)))
			ganancias = 0.3*(float64(pak.Valor))
			perdidas = float64(10*(pak.Intentos-1))	
			//ojo que la penalización no sea mayor que la ganancia		
		}
	}else{
		if(pak.Estado == "rec"){
			balance_producto = float64(pak.Valor - (10*(pak.Intentos-1)))
			ganancias = float64(pak.Valor)
			perdidas = float64(10*(pak.Intentos-1))	
		}else{
			balance_producto = float64(pak.Valor - (10*(pak.Intentos-1)))
			perdidas = float64(10*(pak.Intentos-1))
		}
	}
}


//recibe y convierte json desde logistico
func convertjson(inf []byte) EnviosFinanzas{
	var sstruc EnviosFinanzas

	err := json.Unmarshal([]byte(inf), &sstruc)
	if err != nil{
		fmt.Println(err)
	}
	return sstruc
}


//cantidad no entregados


func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}



func main(){
	balance_t = 0
	ganancias_f = 0
	perdidas_f = 0

	//log.Printf("...")

	conn, err := amqp.Dial("amqp://birmania:birmania@10.6.40.157:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "fallo al abrir el cananl")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "fallo al declarar la cola")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "fallo al consumir datos")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Recibiendo informacion...")
			nuevo := convertjson(d.Body)
			if(nuevo.Estado == "nr"){
				norecib = append(norecib, nuevo)
			}
			contador(nuevo)
			nuevo.Balance = balance_producto
			ganancias_f += ganancias
			perdidas_f += perdidas
			/*csv?*/
			completados = append(completados, nuevo)
			balance_t += balance_final + ganancias - perdidas
			

		}
	}()

	  
	log.Printf("Esperando mensajes... Para salir CTRL+C")
	<-forever

log.Printf("El balance final es: %f dignipesos", balance_t)
}