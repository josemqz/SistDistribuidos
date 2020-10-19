package main
import(
	"log"
	//"amqp"
	//"fmt"
	//"bufio"
	//"reflect"
	//"math"
	//"strconv"
	//"os"
)
/*
func main(){
	log.Println(28/10)
	log.Println(math.Floor(28/10))
	log.Println(math.Floor(13/10))
}
*/


func main() {

	//var u int
	//bla = 10

	/*reader := bufio.NewReader(os.Stdin)

	log.Println("Ingresar tiempo de espera entre pedidos: ")

	read, _ := reader.ReadString('\n')
	u, _ := strconv.Atoi(read)*/

	
	type Package struct{
		id string
		num_seguimiento int32
		tipo string
		valor int32
		num_intentos int32
		estado string
		// estados-> bdg: bodega, tr: en tránsito , rec: recibido , nr: no recibido
	}
	
	var Paquetes []Package

	Paquetes = append(Paquetes, Package{id:"1"})
	Paquetes = append(Paquetes, Package{id:"2"})
	Paquetes = append(Paquetes, Package{id:"3"})
	Paquetes = append(Paquetes, Package{id:"4"})
	Paquetes = append(Paquetes, Package{id:"5"})
	
	
	for i := range(Paquetes){
		log.Println(Paquetes[i])
	}

	Paquetes = Paquetes[1:]

	for i := range(Paquetes){
		log.Println(Paquetes[i])
	}

	//log.Println(len(Paquetes))

	
	
	/*
	log.Println("ingresa:")
	
	fmt.Scanf("%d", &u)
	//log.Println(reflect.TypeOf(u))
	log.Println(reflect.TypeOf(u))


	log.Println(u == 10)
	log.Println(u == 8)
	log.Println(u != 10)
	log.Println(u != 23)
	log.Println(u > 1)
	log.Println(u < 5)
*/



}