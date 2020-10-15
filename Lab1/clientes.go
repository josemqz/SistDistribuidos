//clientes
package main

import (
	"fmt"
	"os"

	//"time"
	"log"

	"encoding/csv"

	//"bufio"

	//csv to json
	"encoding/json"
)

type Retailparcel struct {
	ID       string `json:"id"`
	Producto string `json:"producto"`
	Valor    string `json:"valor"`
	Tienda   string `json:"tienda"`
	Destino  string `json:"destino"`
}

func leer(nombre string) *os.File {

	arch, err := os.Open(nombre)

	if err != nil {
		log.Fatalln("No se pudo abrir el archivo"+nombre+"para lectura", err) //chequear si está bien, o es con coma o qué
	}

	return arch
}

func main() {

	//declarar variable
	//var [nombre] [tipo]

	//time:
	//time.sleep([segundos] time.Second)

	//leer csv
	//pyme_csv := leer("pymes.csv")
	//rPyme := csv.NewReader(pyme_csv)

	//retail_csv := leer("retail.csv")
	//rRetail := csv.NewReader(retail_csv)

	//for {

	// Lee cada entrada (pedido) del csv
	//	pedidoP, err_read := rPyme.Read()

	//	if err_read == io.EOF {
	//		break
	//	}

	//	if err_read != nil {
	//		log.Fatal(err_read)
	//		return
	//	}

	//	fmt.Printf("%s, %s\n", pedidoP[0], pedidoP[1])
	//}

	csv_file, err := os.Open("retail.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csv_file.Close()
	r := csv.NewReader(csv_file)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var rparcel Retailparcel
	var rparcels []Retailparcel
	for _, rec := range records {
		rparcel.ID = rec[0]
		rparcel.Producto = rec[1]
		rparcel.Valor, _ = rec[2]
		rparcel.Tienda, _ = rec[3]
		rparcel.Destino = rec[4]
		rparcels = append(rparcels, rparcel)
	}
	// Convert to JSON
	json_data, err := json.Marshal(rparcels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//print json data
	fmt.Println(string(json_data))

	//create json file
	json_file, err := os.Create("sample.json")
	if err != nil {
		fmt.Println(err)
	}
	defer json_file.Close()

	json_file.Write(json_data)
	json_file.Close()

	//cada cierto tiempo (definido por input (creo)) tiene que leer
	//una entrada y enviar el pedido a logística

	//realizar conexion gRPC como client y enviar datos
	//hay que crear tres colas (se crean en client y server porsiaca)

	//tiene que ser capaz de pedir seguimiento de pedido, mediante el id de compra
	//una vez hecho esto tiene que recibir el mensaje de seguimiento e imprimir

}
