package main
import(
	"log"
	"fmt"
	//"bufio"
	"reflect"
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

	var u int
	//bla = 10

	/*reader := bufio.NewReader(os.Stdin)

	log.Println("Ingresar tiempo de espera entre pedidos: ")

	read, _ := reader.ReadString('\n')
	u, _ := strconv.Atoi(read)*/

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


}