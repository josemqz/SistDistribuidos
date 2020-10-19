//financiero
package main

import (
	
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"sync"
)

/*
- Los envı́os completados -> se guardan en completados[]
- La cantidad de veces que se intentó entregar un paquete -> paquete.intentos
- Los paquetes que no pudieron ser entregados -> se guardan en norecib[]
- Pérdidas o ganancias de cada paquete (en Dignipesos)-> se calcula en contador(paquete)
*/

//Además mostrará el balance final en dignipesos cuando termine su ejecución.




type EnviosFinanzas struct {
	Id_paquete string `json:"id"`
	Tipo string `json:"tipo"`
	Valor int32 `json:"valor"`
	Intentos int32 `json:"intentos"`
	//Date_of_delivery time.Time
	Estado string `json:"estado"`
	Balance float64

}



var balance_p float64
var ganancia_p float64
var perdidas_p float64

var balance_t float64
var ganancias_t float64
var perdidas_t float64

var completados []EnviosFinanzas
var norecib []EnviosFinanzas

var mutex = &sync.Mutex{}


func csvData(reg EnviosFinanzas) {

	ap, err := os.OpenFile("registro_financiero.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}
	defer ap.Close()
	if _, err := ap.WriteString("ID: " + reg.Id_paquete + "," + "Tipo: " + reg.Tipo + "," + "Valor: " + fmt.Sprint(reg.Valor) + "," + "Intentos: " + fmt.Sprint(reg.Intentos) + "," + "Estado: " + reg.Estado + "," + "Balance: " + fmt.Sprintf("%f", reg.Balance) + "\n"); err != nil {
		log.Println(err)
	}
}


func finSesion(){

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {

		<-c

		mutex.Lock()
		log.Println("Balance total:", balance_t)
		mutex.Unlock()

		ap, err := os.OpenFile("registro_financiero.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println(err)
		}

		defer ap.Close()
		if _, err := ap.WriteString("Número de pedidos completados: " + fmt.Sprint(len(completados)) + ",Número de pedidos no entregados: " + fmt.Sprint(len(norecib)) + "\n"); err != nil {
			log.Println(err)
		}

		os.Exit(0)
	}()
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

//calcula perdidas y ganancias segun tipo de envío
func contador(pak EnviosFinanzas){

	balance_p = 0
	ganancia_p = 0
	perdidas_p  = 0

	if(pak.Tipo == "retail"){
		balance_p = float64(pak.Valor) - float64(10*(pak.Intentos - 1)) 
		ganancia_p = float64(pak.Valor)
		perdidas_p = float64(10*(pak.Intentos - 1))
	}else if(pak.Tipo == "prioritario"){
		if(pak.Estado == "rec"){
			balance_p = float64(pak.Valor - (10*(pak.Intentos - 1)))
			ganancia_p = float64(pak.Valor)
			perdidas_p = float64(10*(pak.Intentos - 1))
		}else{
			balance_p = float64(pak.Valor - (10*(pak.Intentos - 1)))
			ganancia_p = 0.3*(float64(pak.Valor))
			perdidas_p = float64(10*(pak.Intentos - 1))			
		}
	}else{
		if(pak.Estado == "rec"){
			balance_p = float64(pak.Valor - (10*(pak.Intentos-1)))
			ganancia_p = float64(pak.Valor)
			perdidas_p = float64(10*(pak.Intentos - 1))	
		}else{
			balance_p = float64(pak.Valor - (10*(pak.Intentos - 1)))
			perdidas_p = float64(10*(pak.Intentos - 1))
		}
	}
}

//cantidad entregados

//cantidad no entregados

func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func main(){

	balance_t = 0
	ganancias_t = 0
	perdidas_t = 0

	finSesion()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error de conexion")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Error al abrir canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "error al declarar la cola")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	forever := make(chan bool)

	go func() {
		
		var i = 0
		
		for x := range msgs {

			if i > 0{
				log.Printf("Recibiendo informacion, calculando...")
			}
			
			nuevo := convertjson(x.Body)
			if(nuevo.Estado == "nr"){
				norecib = append(norecib, nuevo)
			}
			contador(nuevo)
			nuevo.Balance = balance_p
			ganancias_t += ganancia_p
			perdidas_t += perdidas_p
			csvData(nuevo)
			
			mutex.Lock()
			completados = append(completados, nuevo)
			mutex.Unlock()

	)
	failOnError(err, "error al consumir de la cola")

	forever := make(chan bool)

	go func() {
		for x := range msgs {

			log.Printf("Recibiendo informacion, calculando...")
			
			nuevo := convertjson(x.Body)
			if(nuevo.Estado == "nr"){
				norecib = append(norecib, nuevo)
			}
			contador(nuevo)
			nuevo.Balance = balance_p
			ganancias_t += ganancia_p
			perdidas_t += perdidas_p
			csvData(nuevo)
			
			mutex.Lock()
			completados = append(completados, nuevo)
			mutex.Unlock()

			balance_t += balance_p
			
		}
	}()

	log.Printf("El balance final es: %f dignipesos", balance_t)

	  
	log.Printf(" [*] Esperando mensajes... Para salir CTRL+C")
	<-forever

	log.Printf("El balance final es: %f dignipesos", balance_t)
}